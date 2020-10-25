package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/earthfile2llb/variables"
	"github.com/earthly/earthly/states"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/entitlements"
	reccopy "github.com/otiai10/copy"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Opt represent builder options.
type Opt struct {
	SessionID            string
	BkClient             *client.Client
	Console              conslogging.ConsoleLogger
	Verbose              bool
	Attachables          []session.Attachable
	Enttlmnts            []entitlements.Entitlement
	NoCache              bool
	RemoteCache          string
	ImageResolveMode     llb.ResolveMode
	CleanCollection      *cleanup.Collection
	VarCollection        *variables.Collection
	BuildContextProvider *provider.BuildContextProvider
}

// BuildOpt is a collection of build options.
type BuildOpt struct {
	PrintSuccess bool
	Push         bool
}

// Builder executes Earthly builds.
type Builder struct {
	s        *solver
	opt      Opt
	resolver *buildcontext.Resolver
}

// NewBuilder returns a new earth Builder.
func NewBuilder(ctx context.Context, opt Opt) (*Builder, error) {
	b := &Builder{
		s: &solver{
			sm:          newSolverMonitor(opt.Console, opt.Verbose),
			bkClient:    opt.BkClient,
			remoteCache: opt.RemoteCache,
			attachables: opt.Attachables,
			enttlmnts:   opt.Enttlmnts,
		},
		opt:      opt,
		resolver: nil, // initialized below
	}
	b.resolver = buildcontext.NewResolver(
		opt.SessionID, opt.CleanCollection, b.MakeArtifactBuilderFun())
	return b, nil
}

// BuildTarget executes the build of a given Earthly target.
func (b *Builder) BuildTarget(ctx context.Context, target domain.Target, opt BuildOpt) (*states.MultiTarget, error) {
	mts, err := b.convertAndBuild(ctx, target, opt)
	if err != nil {
		return nil, err
	}
	return mts, nil
}

// Output produces the output for all targets provided.
func (b *Builder) Output(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	for _, states := range mts.All() {
		err := b.outputs(ctx, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputArtifact outputs a specific artifact found in an mts's final target.
func (b *Builder) OutputArtifact(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string, opt BuildOpt) error {
	// TODO: Should double check that the last state is the same as the one
	//       referenced in artifact.Target.
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	solvedStates := make(map[int]bool)
	return b.outputSpecifiedArtifact(
		ctx, artifact, destPath, outDir, solvedStates, mts.Final, opt)
}

// MakeArtifactBuilderFun returns a function that can be used to build artifacts.
func (b *Builder) MakeArtifactBuilderFun() states.ArtifactBuilderFun {
	return func(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string) error {
		return b.buildOnlyArtifact(ctx, mts, artifact, destPath, BuildOpt{})
	}
}

// MakeImageAsTarBuilderFun returns a function which can be used to build an image as a tar.
func (b *Builder) MakeImageAsTarBuilderFun() states.DockerBuilderFun {
	return func(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string) error {
		return b.buildOnlyLastImageAsTar(ctx, mts, dockerTag, outFile, BuildOpt{})
	}
}

// OutputImages outputs the images of a single target.
func (b *Builder) OutputImages(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	for _, imageToSave := range states.SaveImages {
		if imageToSave.DockerTag == "" {
			// Not a docker export. Skip.
			continue
		}
		err := b.outputImage(ctx, imageToSave, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputArtifacts outputs the artifacts of a single target.
func (b *Builder) OutputArtifacts(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	solvedStates := make(map[int]bool)
	for _, artifactToSaveLocally := range states.SaveLocals {
		err = b.outputArtifact(
			ctx, artifactToSaveLocally, outDir, solvedStates, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) convertAndBuild(ctx context.Context, target domain.Target, opt BuildOpt) (*states.MultiTarget, error) {
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return nil, errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	var successOnce sync.Once
	successFun := func() {
		if opt.PrintSuccess {
			b.opt.Console.PrintSuccess()
		}
	}
	destPathWhitelist := make(map[string]bool)
	var mts *states.MultiTarget
	bf := func(ctx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		var err error
		mts, err = earthfile2llb.Earthfile2LLB(ctx, target, earthfile2llb.ConvertOpt{
			GwClient:             gwClient,
			Resolver:             b.resolver,
			ImageResolveMode:     b.opt.ImageResolveMode,
			DockerBuilderFun:     b.MakeImageAsTarBuilderFun(),
			ArtifactBuilderFun:   b.MakeArtifactBuilderFun(),
			CleanCollection:      b.opt.CleanCollection,
			VarCollection:        b.opt.VarCollection,
			BuildContextProvider: b.opt.BuildContextProvider,
		})
		if err != nil {
			return nil, err
		}
		state := mts.Final.MainState
		if b.opt.NoCache {
			state = state.SetMarshalDefaults(llb.IgnoreCache)
		}
		def, err := state.Marshal(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "marshal main state")
		}
		r, err := gwClient.Solve(ctx, gwclient.SolveRequest{
			Definition: def.ToPB(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "solve main state")
		}
		ref, err := r.SingleRef()
		if err != nil {
			return nil, err
		}
		res := gwclient.NewResult()
		res.AddRef("earthly-main", ref)

		imageIndex := 0
		dirIndex := 0
		for _, sts := range mts.All() {
			for _, saveImage := range sts.SaveImages {
				def, err := saveImage.State.Marshal(ctx)
				if err != nil {
					return nil, errors.Wrap(err, "marshal save img state")
				}
				r, err := gwClient.Solve(ctx, gwclient.SolveRequest{
					Definition: def.ToPB(),
				})
				ref, err := r.SingleRef()
				if err != nil {
					return nil, err
				}
				config, err := json.Marshal(saveImage.Image)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal save image config")
				}
				// TODO: Support multiple docker tags at the same time (improves export speed).
				res.AddMeta(fmt.Sprintf("image.name/%d", imageIndex), []byte(saveImage.DockerTag))
				res.AddMeta(fmt.Sprintf("%s/%d", exptypes.ExporterImageConfigKey, imageIndex), config)
				refKey := fmt.Sprintf("earthly-image-%d", imageIndex)
				res.AddRef(refKey, ref)
				imageIndex++
			}
			for _, saveLocal := range sts.SaveLocals {
				def, err := sts.SeparateArtifactsState[saveLocal.Index].Marshal(ctx)
				if err != nil {
					return nil, errors.Wrap(err, "marshal save img state")
				}
				r, err := gwClient.Solve(ctx, gwclient.SolveRequest{
					Definition: def.ToPB(),
				})
				ref, err := r.SingleRef()
				if err != nil {
					return nil, err
				}
				refKey := fmt.Sprintf("earthly-dir-%d", dirIndex)
				res.AddRef(refKey, ref)
				artifact := domain.Artifact{
					Target:   sts.Target,
					Artifact: saveLocal.ArtifactPath,
				}
				res.AddMeta(fmt.Sprintf("earthly-artifact/%d", dirIndex), []byte(artifact.String()))
				res.AddMeta(fmt.Sprintf("earthly-src-path/%d", dirIndex), []byte(saveLocal.ArtifactPath))
				res.AddMeta(fmt.Sprintf("earthly-dest-path/%d", dirIndex), []byte(saveLocal.DestPath))
				destPathWhitelist[saveLocal.DestPath] = true
				dirIndex++
			}
		}
		res.AddMeta("earthly-num-images", []byte(fmt.Sprintf("%d", imageIndex)))
		res.AddMeta("earthly-num-dirs", []byte(fmt.Sprintf("%d", dirIndex)))
		return res, nil
	}
	onImage := func(ctx context.Context, eg *errgroup.Group, index int, imageName string, digest string) (io.WriteCloser, error) {
		successOnce.Do(successFun)
		pipeR, pipeW := io.Pipe()
		eg.Go(func() error {
			defer pipeR.Close()
			err := loadDockerTar(ctx, pipeR)
			if err != nil {
				return errors.Wrapf(err, "load docker tar")
			}
			return nil
		})
		return pipeW, nil
	}
	onArtifact := func(ctx context.Context, index int, artifact domain.Artifact, artifactPath string, destPath string) (string, error) {
		successOnce.Do(successFun)
		if !destPathWhitelist[destPath] {
			return "", errors.Errorf("dest path %s is not in the whitelist: %+v", destPath, destPathWhitelist)
		}
		fmt.Printf("@#@#@# received onArtifact(%d, %s, %s, %s)\n", index, artifact.String(), artifactPath, destPath)
		return outDir, nil
	}
	err = b.s.buildMainMulti(ctx, bf, onImage, onArtifact)
	if err != nil {
		return nil, errors.Wrapf(err, "build main")
	}
	successOnce.Do(successFun)
	return mts, nil
}

func (b *Builder) buildOnlyLastImageAsTar(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string, opt BuildOpt) error {
	saveImage := mts.Final.LastSaveImage()
	err := b.buildMain(ctx, mts, opt)
	if err != nil {
		return err
	}

	err = b.outputImageTar(ctx, saveImage, dockerTag, outFile)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) buildOnlyArtifact(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string, opt BuildOpt) error {
	err := b.buildMain(ctx, mts, opt)
	if err != nil {
		return err
	}

	return b.OutputArtifact(ctx, mts, artifact, destPath, opt)
}

func (b *Builder) buildMain(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	finalTarget := mts.Final.Target
	finalTargetConsole := b.opt.Console.WithPrefixAndSalt(finalTarget.String(), mts.Final.Salt)
	state := mts.Final.MainState
	if b.opt.NoCache {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	err := b.s.solveMain(ctx, state)
	if err != nil {
		return errors.Wrapf(err, "solve side effects")
	}
	if opt.PrintSuccess {
		finalTargetConsole.Printf("Target %s built successfully\n", finalTarget.StringCanonical())
	}
	return nil
}

func (b *Builder) outputs(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	err := b.executeRunPush(ctx, states, opt)
	if err != nil {
		return err
	}
	// @#
	// err = b.OutputImages(ctx, states, opt)
	// if err != nil {
	// 	return err
	// }
	// if !states.Target.IsRemote() {
	// 	// Don't output artifacts for remote images.
	// 	err = b.OutputArtifacts(ctx, states, opt)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (b *Builder) executeRunPush(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	if !states.RunPush.Initialized {
		// No run --push commands here. Quick way out.
		return nil
	}
	console := b.opt.Console.WithPrefixAndSalt(states.Target.String(), states.Salt)
	if !opt.Push {
		for _, commandStr := range states.RunPush.CommandStrs {
			console.Printf("Did not execute push command %s. Use earth --push to enable pushing\n", commandStr)
		}
		return nil
	}
	err := b.s.solveMain(ctx, states.RunPush.State)
	if err != nil {
		return errors.Wrapf(err, "solve run-push")
	}
	return nil
}

func (b *Builder) outputImage(ctx context.Context, imageToSave states.SaveImage, states *states.SingleTarget, opt BuildOpt) error {
	shouldPush := opt.Push && imageToSave.Push
	console := b.opt.Console.WithPrefixAndSalt(states.Target.String(), states.Salt)
	err := b.s.solveDocker(ctx, imageToSave.State, imageToSave.Image, imageToSave.DockerTag, shouldPush)
	if err != nil {
		return errors.Wrapf(err, "solve image %s", imageToSave.DockerTag)
	}
	pushStr := ""
	if shouldPush {
		pushStr = " (pushed)"
	}
	console.Printf("Image %s as %s%s\n", states.Target.StringCanonical(), imageToSave.DockerTag, pushStr)
	if imageToSave.Push && !opt.Push {
		console.Printf("Did not push %s. Use earth --push to enable pushing\n", imageToSave.DockerTag)
	}
	return nil
}

func (b *Builder) outputImageTar(ctx context.Context, saveImage states.SaveImage, dockerTag string, outFile string) error {
	err := b.s.solveDockerTar(ctx, saveImage.State, saveImage.Image, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "solve image tar %s", outFile)
	}
	return nil
}

func (b *Builder) outputSpecifiedArtifact(ctx context.Context, artifact domain.Artifact, destPath string, outDir string, solvedStates map[int]bool, states *states.SingleTarget, opt BuildOpt) error {
	indexOutDir := filepath.Join(outDir, "combined")
	err := os.Mkdir(indexOutDir, 0755)
	if err != nil {
		return errors.Wrap(err, "mk index dir")
	}
	artifactsState := states.ArtifactsState
	err = b.s.solveArtifacts(ctx, artifactsState, indexOutDir)
	if err != nil {
		return errors.Wrap(err, "solve combined artifacts")
	}

	err = b.saveArtifactLocally(ctx, artifact, indexOutDir, destPath, states.Salt, opt)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) outputArtifact(ctx context.Context, artifactToSaveLocally states.SaveLocal, outDir string, solvedStates map[int]bool, states *states.SingleTarget, opt BuildOpt) error {
	index := artifactToSaveLocally.Index
	indexOutDir := filepath.Join(outDir, fmt.Sprintf("index-%d", index))
	if !solvedStates[index] {
		solvedStates[index] = true
		artifactsState := states.SeparateArtifactsState[artifactToSaveLocally.Index]
		err := os.Mkdir(indexOutDir, 0755)
		if err != nil {
			return errors.Wrap(err, "mk index dir")
		}
		err = b.s.solveArtifacts(ctx, artifactsState, indexOutDir)
		if err != nil {
			return errors.Wrap(err, "solve artifacts")
		}
	}

	artifact := domain.Artifact{
		Target:   states.Target,
		Artifact: artifactToSaveLocally.ArtifactPath,
	}
	err := b.saveArtifactLocally(ctx, artifact, indexOutDir, artifactToSaveLocally.DestPath, states.Salt, opt)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) saveArtifactLocally(ctx context.Context, artifact domain.Artifact, indexOutDir string, destPath string, salt string, opt BuildOpt) error {
	console := b.opt.Console.WithPrefixAndSalt(artifact.Target.String(), salt)
	fromPattern := filepath.Join(indexOutDir, filepath.FromSlash(artifact.Artifact))
	// Resolve possible wildcards.
	// TODO: Note that this is not very portable, as the glob is host-platform dependent,
	//       while the pattern is also guest-platform dependent.
	fromGlobMatches, err := filepath.Glob(fromPattern)
	if err != nil {
		return errors.Wrapf(err, "glob")
	}
	isWildcard := (len(fromGlobMatches) > 1)
	for _, from := range fromGlobMatches {
		fiSrc, err := os.Stat(from)
		if err != nil {
			return errors.Wrapf(err, "os stat %s", from)
		}
		srcIsDir := fiSrc.IsDir()
		to := destPath
		destIsDir := strings.HasSuffix(to, "/")
		if artifact.Target.IsLocalExternal() && !filepath.IsAbs(to) {
			// Place within external dir.
			to = path.Join(artifact.Target.LocalPath, to)
		} else {
		}
		if !srcIsDir && destIsDir {
			// Place within dest dir.
			to = path.Join(to, path.Base(from))
		}
		destExists := false
		fiDest, err := os.Stat(to)
		if err != nil {
			// Ignore err. Likely dest path does not exist.
			if isWildcard && !destIsDir {
				return errors.New(
					"Artifact is a wildcard, but AS LOCAL destination does not end with /")
			}
			destIsDir = fiSrc.IsDir()
		} else {
			destExists = true
			destIsDir = fiDest.IsDir()
		}
		if srcIsDir && !destIsDir {
			return errors.New(
				"Artifact is a directory, but existing AS LOCAL destination is a file")
		}
		if destExists {
			if !srcIsDir {
				// Remove pre-existing dest file.
				err = os.Remove(to)
				if err != nil {
					return errors.Wrapf(err, "rm %s", to)
				}
			} else {
				// Remove pre-existing dest dir.
				err = os.RemoveAll(to)
				if err != nil {
					return errors.Wrapf(err, "rm -rf %s", to)
				}
			}
		}

		toDir := path.Dir(to)
		err = os.MkdirAll(toDir, 0755)
		if err != nil {
			return errors.Wrapf(err, "mkdir all for artifact %s", toDir)
		}
		err = os.Link(from, to)
		if err != nil {
			// Hard linking did not work. Try recursive copy.
			errCopy := reccopy.Copy(from, to)
			if errCopy != nil {
				return errors.Wrapf(errCopy, "copy artifact %s", from)
			}
		}

		// Write to console about this artifact.
		parts := strings.Split(filepath.ToSlash(from), "/")
		artifactPath := artifact.Artifact
		if len(parts) >= 2 {
			artifactPath = filepath.FromSlash(strings.Join(parts[2:], "/"))
		}
		artifact2 := domain.Artifact{
			Target:   artifact.Target,
			Artifact: artifactPath,
		}
		destPath2 := filepath.FromSlash(destPath)
		if strings.HasSuffix(destPath, "/") {
			destPath2 = filepath.Join(destPath2, filepath.Base(artifactPath))
		}
		if opt.PrintSuccess {
			console.Printf("Artifact %s as local %s\n", artifact2.StringCanonical(), destPath2)
		}
	}
	return nil
}

func loadDockerTar(ctx context.Context, r io.ReadCloser) error {
	// TODO: This is a gross hack - should use proper docker client.
	cmd := exec.CommandContext(ctx, "docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker load")
	}
	return nil
}
