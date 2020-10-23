package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/earthfile2llb"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
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
	NoOutput     bool
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

func (b *Builder) BuildTarget(ctx context.Context, target domain.Target, opt BuildOpt) error {
	mts, err := b.buildCommon2(ctx, target, opt)
	if err != nil {
		return err
	}
	if opt.PrintSuccess {
		b.opt.Console.PrintSuccess()
	}
	if !opt.NoOutput {
		for _, states := range mts.All() {
			err = b.buildOutputs(ctx, states, opt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Builder) buildCommon2(ctx context.Context, target domain.Target, opt BuildOpt) (*states.MultiTarget, error) {
	var mts *states.MultiTarget
	bf := func(ctx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		var err error
		mts, err = earthfile2llb.Earthfile2LLB(ctx, target, earthfile2llb.ConvertOpt{
			Resolver:             b.resolver,
			ImageResolveMode:     b.opt.ImageResolveMode,
			DockerBuilderFun:     b.MakeImageAsTarBuilderFun(),
			ArtifactBuilderFun:   b.MakeArtifactBuilderFun(),
			CleanCollection:      b.opt.CleanCollection,
			VarCollection:        b.opt.VarCollection,
			BuildContextProvider: b.opt.BuildContextProvider,
			MetaResolver:         gwClient,
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
		config, err := json.Marshal(mts.Final.MainImage)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal image config")
		}

		res := gwclient.NewResult()
		res.AddMeta(exptypes.ExporterImageConfigKey, config)
		res.SetRef(ref)
		return res, nil
	}
	err := b.s.buildMain(ctx, bf)
	if err != nil {
		return nil, errors.Wrapf(err, "build main")
	}
	if opt.PrintSuccess {
		targetConsole := b.opt.Console.WithPrefixAndSalt(target.String(), mts.Final.Salt)
		targetConsole.Printf("Target %s built successfully\n", target.StringCanonical())
	}
	return mts, nil
}

// Build performs the build for the given multi target states, outputting images for
// all sub-targets and artifacts for all local sub-targets.
func (b *Builder) Build(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	// Start with final side-effects. This will automatically trigger the dependency builds too,
	// in parallel.
	err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	if opt.PrintSuccess {
		b.opt.Console.PrintSuccess()
	}

	// Then output images and artifacts.
	if !opt.NoOutput {
		for _, states := range mts.All() {
			err = b.buildOutputs(ctx, states, opt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// BuildOnlyLastImageAsTar performs the build for the given multi target states,
// and outputs only a docker tar of the last saved image.
func (b *Builder) BuildOnlyLastImageAsTar(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string, opt BuildOpt) error {
	saveImage := mts.Final.LastSaveImage()
	err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	if opt.PrintSuccess {
		b.opt.Console.PrintSuccess()
	}

	err = b.buildImageTar(ctx, saveImage, dockerTag, outFile)
	if err != nil {
		return err
	}
	return nil
}

// MakeImageAsTarBuilderFun returns a fun which can be used to build an image as a tar.
func (b *Builder) MakeImageAsTarBuilderFun() states.DockerBuilderFun {
	return func(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string) error {
		return b.BuildOnlyLastImageAsTar(ctx, mts, dockerTag, outFile, BuildOpt{})
	}
}

// BuildOnlyImages performs the build for the given multi target states, outputting only images
// of the final states.
func (b *Builder) BuildOnlyImages(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	// Start with final side-effects. This will automatically trigger the dependency builds too,
	// in parallel.
	err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	if opt.PrintSuccess {
		b.opt.Console.PrintSuccess()
	}

	err = b.buildImages(ctx, mts.Final, opt)
	if err != nil {
		return err
	}
	return nil
}

// BuildOnlyArtifact performs the build for the given multi target states, outputting only
// the provided artifact of the final states.
func (b *Builder) BuildOnlyArtifact(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string, opt BuildOpt) error {
	// Start with final side-effects. This will automatically trigger the dependency builds too,
	// in parallel.
	err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	if opt.PrintSuccess {
		b.opt.Console.PrintSuccess()
	}

	// TODO: Should double check that the last state is the same as the one
	//       referenced in artifact.Target.
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	solvedStates := make(map[int]bool)
	err = b.buildSpecifiedArtifact(
		ctx, artifact, destPath, outDir, solvedStates, mts.Final, opt)
	if err != nil {
		return err
	}

	return nil
}

// MakeArtifactBuilderFun returns a function that can be used to build artifacts.
func (b *Builder) MakeArtifactBuilderFun() states.ArtifactBuilderFun {
	return func(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string) error {
		return b.BuildOnlyArtifact(ctx, mts, artifact, destPath, BuildOpt{})
	}
}

func (b *Builder) buildCommon(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
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

func (b *Builder) buildOutputs(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	// Run --push commands.
	err := b.buildRunPush(ctx, states, opt)
	if err != nil {
		return err
	}

	// Images.
	err = b.buildImages(ctx, states, opt)
	if err != nil {
		return err
	}

	// Artifacts.
	if !states.Target.IsRemote() {
		// Don't output artifacts for remote images.
		err = b.buildArtifacts(ctx, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildRunPush(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
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

func (b *Builder) buildImages(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	for _, imageToSave := range states.SaveImages {
		if imageToSave.DockerTag == "" {
			// Not a docker export. Skip.
			continue
		}
		err := b.buildImage(ctx, imageToSave, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildImage(ctx context.Context, imageToSave states.SaveImage, states *states.SingleTarget, opt BuildOpt) error {
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

func (b *Builder) buildImageTar(ctx context.Context, saveImage states.SaveImage, dockerTag string, outFile string) error {
	err := b.s.solveDockerTar(ctx, saveImage.State, saveImage.Image, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "solve image tar %s", outFile)
	}
	return nil
}

func (b *Builder) buildArtifacts(ctx context.Context, states *states.SingleTarget, opt BuildOpt) error {
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	solvedStates := make(map[int]bool)
	for _, artifactToSaveLocally := range states.SaveLocals {
		err = b.buildArtifact(
			ctx, artifactToSaveLocally, outDir, solvedStates, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildSpecifiedArtifact(ctx context.Context, artifact domain.Artifact, destPath string, outDir string, solvedStates map[int]bool, states *states.SingleTarget, opt BuildOpt) error {
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

func (b *Builder) buildArtifact(ctx context.Context, artifactToSaveLocally states.SaveLocal, outDir string, solvedStates map[int]bool, states *states.SingleTarget, opt BuildOpt) error {
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
