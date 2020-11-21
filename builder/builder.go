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
	"github.com/earthly/earthly/llbutil"
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
	PrintSuccess          bool
	Push                  bool
	NoOutput              bool
	OnlyFinalTargetImages bool
	OnlyArtifact          *domain.Artifact
	OnlyArtifactDestPath  string
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
	b.resolver = buildcontext.NewResolver(opt.SessionID, opt.CleanCollection)
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

// MakeImageAsTarBuilderFun returns a function which can be used to build an image as a tar.
func (b *Builder) MakeImageAsTarBuilderFun() states.DockerBuilderFun {
	return func(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string) error {
		return b.buildOnlyLastImageAsTar(ctx, mts, dockerTag, outFile, BuildOpt{})
	}
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
			b.s.sm.PrintTiming()
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
			CleanCollection:      b.opt.CleanCollection,
			VarCollection:        b.opt.VarCollection,
			BuildContextProvider: b.opt.BuildContextProvider,
		})
		if err != nil {
			return nil, err
		}
		ref, err := b.stateToRef(ctx, gwClient, mts.Final.MainState)
		if err != nil {
			return nil, err
		}
		res := gwclient.NewResult()
		res.AddRef("main", ref)
		ref, err = b.stateToRef(ctx, gwClient, mts.Final.ArtifactsState)
		if err != nil {
			return nil, err
		}
		refKey := "final-artifact"
		refPrefix := fmt.Sprintf("ref/%s", refKey)
		res.AddRef(refKey, ref)
		if !opt.NoOutput && opt.OnlyArtifact != nil && !opt.OnlyFinalTargetImages {
			res.AddMeta(fmt.Sprintf("%s/export-dir", refPrefix), []byte("true"))
		}
		res.AddMeta(fmt.Sprintf("%s/final-artifact", refPrefix), []byte("true"))

		depIndex := 0
		imageIndex := 0
		dirIndex := 0
		for _, sts := range mts.All() {
			for _, depRef := range sts.DepsRefs {
				refKey := fmt.Sprintf("dep-%d", depIndex)
				res.AddRef(refKey, depRef)
				depIndex++
			}
			for _, saveImage := range sts.SaveImages {
				ref, err := b.stateToRef(ctx, gwClient, saveImage.State)
				if err != nil {
					return nil, err
				}
				config, err := json.Marshal(saveImage.Image)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal save image config")
				}
				// TODO: Support multiple docker tags at the same time (improves export speed).
				refKey := fmt.Sprintf("image-%d", imageIndex)
				refPrefix := fmt.Sprintf("ref/%s", refKey)
				res.AddMeta(fmt.Sprintf("%s/image.name", refPrefix), []byte(saveImage.DockerTag))
				res.AddMeta(fmt.Sprintf("%s/%s", refPrefix, exptypes.ExporterImageConfigKey), config)
				if !opt.NoOutput && opt.OnlyArtifact == nil && !(opt.OnlyFinalTargetImages && sts != mts.Final) {
					res.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
				}
				res.AddMeta(fmt.Sprintf("%s/image-index", refPrefix), []byte(fmt.Sprintf("%d", imageIndex)))
				res.AddRef(refKey, ref)
				imageIndex++
			}
			if sts.Target.IsRemote() {
				// Don't do save local's for remote targets.
				continue
			}
			for _, saveLocal := range sts.SaveLocals {
				ref, err := b.stateToRef(ctx, gwClient, sts.SeparateArtifactsState[saveLocal.Index])
				if err != nil {
					return nil, err
				}
				refKey := fmt.Sprintf("dir-%d", dirIndex)
				refPrefix := fmt.Sprintf("ref/%s", refKey)
				res.AddRef(refKey, ref)
				artifact := domain.Artifact{
					Target:   sts.Target,
					Artifact: saveLocal.ArtifactPath,
				}
				res.AddMeta(fmt.Sprintf("%s/artifact", refPrefix), []byte(artifact.String()))
				res.AddMeta(fmt.Sprintf("%s/src-path", refPrefix), []byte(saveLocal.ArtifactPath))
				res.AddMeta(fmt.Sprintf("%s/dest-path", refPrefix), []byte(saveLocal.DestPath))
				if !opt.NoOutput && !opt.OnlyFinalTargetImages && opt.OnlyArtifact == nil {
					res.AddMeta(fmt.Sprintf("%s/export-dir", refPrefix), []byte("true"))
				}
				res.AddMeta(fmt.Sprintf("%s/dir-index", refPrefix), []byte(fmt.Sprintf("%d", dirIndex)))
				destPathWhitelist[saveLocal.DestPath] = true
				dirIndex++
			}
		}
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
		artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%d", index))
		err := os.MkdirAll(artifactDir, 0755)
		if err != nil {
			return "", errors.Wrapf(err, "create dir %s", artifactDir)
		}
		return artifactDir, nil
	}
	onFinalArtifact := func(ctx context.Context) (string, error) {
		successOnce.Do(successFun)
		return outDir, nil
	}
	err = b.s.buildMainMulti(ctx, bf, onImage, onArtifact, onFinalArtifact)
	if err != nil {
		return nil, errors.Wrapf(err, "build main")
	}
	successOnce.Do(successFun)
	if opt.NoOutput {
		// Nothing.
	} else if opt.OnlyArtifact != nil {
		err := b.saveArtifactLocally(ctx, *opt.OnlyArtifact, outDir, opt.OnlyArtifactDestPath, mts.Final.Salt, opt)
		if err != nil {
			return nil, err
		}
	} else if opt.OnlyFinalTargetImages {
		for _, saveImage := range mts.Final.SaveImages {
			shouldPush := opt.Push && saveImage.Push
			console := b.opt.Console.WithPrefixAndSalt(mts.Final.Target.String(), mts.Final.Salt)
			if shouldPush {
				err := pushDockerImage(ctx, saveImage.DockerTag)
				if err != nil {
					return nil, err
				}
			}
			pushStr := ""
			if shouldPush {
				pushStr = " (pushed)"
			}
			console.Printf("Image %s as %s%s\n", mts.Final.Target.StringCanonical(), saveImage.DockerTag, pushStr)
			if saveImage.Push && !opt.Push {
				console.Printf("Did not push %s. Use earth --push to enable pushing\n", saveImage.DockerTag)
			}
		}
	} else {
		// This needs to match with the same index used during output.
		// TODO: This is a little brittle to future code changes.
		dirIndex := 0
		for _, sts := range mts.All() {
			for _, saveImage := range sts.SaveImages {
				shouldPush := opt.Push && saveImage.Push && !sts.Target.IsRemote()
				console := b.opt.Console.WithPrefixAndSalt(sts.Target.String(), sts.Salt)
				if shouldPush {
					err := pushDockerImage(ctx, saveImage.DockerTag)
					if err != nil {
						return nil, err
					}
				}
				pushStr := ""
				if shouldPush {
					pushStr = " (pushed)"
				}
				console.Printf("Image %s as %s%s\n", sts.Target.StringCanonical(), saveImage.DockerTag, pushStr)
				if saveImage.Push && !opt.Push && !sts.Target.IsRemote() {
					console.Printf("Did not push %s. Use earth --push to enable pushing\n", saveImage.DockerTag)
				}
			}
			for _, saveLocal := range sts.SaveLocals {
				artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%d", dirIndex))
				artifact := domain.Artifact{
					Target:   sts.Target,
					Artifact: saveLocal.ArtifactPath,
				}
				err := b.saveArtifactLocally(ctx, artifact, artifactDir, saveLocal.DestPath, sts.Salt, opt)
				if err != nil {
					return nil, err
				}
				dirIndex++
			}
			if !sts.Target.IsRemote() {
				err = b.executeRunPush(ctx, sts, opt)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return mts, nil
}

func (b *Builder) stateToRef(ctx context.Context, gwClient gwclient.Client, state llb.State) (gwclient.Reference, error) {
	if b.opt.NoCache {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	return llbutil.StateToRef(ctx, gwClient, state)
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

func (b *Builder) buildMain(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	state := mts.Final.MainState
	if b.opt.NoCache {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	err := b.s.solveMain(ctx, state)
	if err != nil {
		return errors.Wrapf(err, "solve side effects")
	}
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

func (b *Builder) outputImageTar(ctx context.Context, saveImage states.SaveImage, dockerTag string, outFile string) error {
	err := b.s.solveDockerTar(ctx, saveImage.State, saveImage.Image, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "solve image tar %s", outFile)
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

func pushDockerImage(ctx context.Context, imageName string) error {
	cmd := exec.CommandContext(ctx, "docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker push %s", imageName)
	}
	return nil
}
