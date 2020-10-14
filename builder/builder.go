package builder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/logging"
	"github.com/earthly/earthly/states"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/entitlements"
	reccopy "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

// BuildOpt is a collection of build options.
type BuildOpt struct {
	PrintSuccess bool
	NoOutput     bool
	Push         bool
}

// Builder provides a earth commands executor.
type Builder struct {
	s           *solver
	bkClient    *client.Client
	console     conslogging.ConsoleLogger
	attachables []session.Attachable
	enttlmnts   []entitlements.Entitlement
	noCache     bool
}

// NewBuilder returns a new earth Builder.
func NewBuilder(ctx context.Context, bkClient *client.Client, console conslogging.ConsoleLogger, verbose bool, attachables []session.Attachable, enttlmnts []entitlements.Entitlement, noCache bool, remoteCache string) (*Builder, error) {
	return &Builder{
		s: &solver{
			sm:          newSolverMonitor(console, verbose),
			bkClient:    bkClient,
			remoteCache: remoteCache,
			attachables: attachables,
			enttlmnts:   enttlmnts,
		},
		console: console,
		noCache: noCache,
	}, nil
}

// Build performs the build for the given multi target states, outputting images for
// all sub-targets and artifacts for all local sub-targets.
func (b *Builder) Build(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	// Start with final side-effects. This will automatically trigger the dependency builds too,
	// in parallel.
	cacheLocalDir, localDirs, err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cacheLocalDir)
	if opt.PrintSuccess {
		b.console.PrintSuccess()
	}

	// Then output images and artifacts.
	if !opt.NoOutput {
		for _, states := range mts.All() {
			err = b.buildOutputs(ctx, localDirs, states, opt)
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
	cacheLocalDir, localDirs, err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cacheLocalDir)
	if opt.PrintSuccess {
		b.console.PrintSuccess()
	}

	err = b.buildImageTar(ctx, localDirs, saveImage, dockerTag, outFile)
	if err != nil {
		return err
	}
	return nil
}

// MakeImageAsTarBuilderFun returns a fun which can be used to build an image as a tar.
func (b *Builder) MakeImageAsTarBuilderFun() func(context.Context, *states.MultiTarget, string, string) error {
	return func(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string) error {
		return b.BuildOnlyLastImageAsTar(ctx, mts, dockerTag, outFile, BuildOpt{})
	}
}

// BuildOnlyImages performs the build for the given multi target states, outputting only images
// of the final states.
func (b *Builder) BuildOnlyImages(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) error {
	// Start with final side-effects. This will automatically trigger the dependency builds too,
	// in parallel.
	cacheLocalDir, localDirs, err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cacheLocalDir)
	if opt.PrintSuccess {
		b.console.PrintSuccess()
	}

	err = b.buildImages(ctx, localDirs, mts.Final, opt)
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
	cacheLocalDir, localDirs, err := b.buildCommon(ctx, mts, opt)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cacheLocalDir)
	if opt.PrintSuccess {
		b.console.PrintSuccess()
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
		ctx, artifact, destPath, outDir, solvedStates, localDirs, mts.Final, opt)
	if err != nil {
		return err
	}

	return nil
}

// MakeArtifactBuilderFun returns a function that can be used to build artifacts.
func (b *Builder) MakeArtifactBuilderFun() func(context.Context, *states.MultiTarget, domain.Artifact, string) error {
	return func(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact, destPath string) error {
		return b.BuildOnlyArtifact(ctx, mts, artifact, destPath, BuildOpt{})
	}
}

func (b *Builder) buildCommon(ctx context.Context, mts *states.MultiTarget, opt BuildOpt) (string, map[string]string, error) {
	cacheLocalDir, err := ioutil.TempDir("/tmp", "earthly-cache")
	if err != nil {
		return "", nil, errors.Wrap(err, "make temp dir for cache")
	}
	// Collect all local dirs.
	localDirs := make(map[string]string)
	localDirs["earthly-cache"] = cacheLocalDir
	for _, states := range mts.All() {
		for key, value := range states.LocalDirs {
			existingValue, alreadyExists := localDirs[key]
			if alreadyExists && existingValue != value {
				return "", nil, fmt.Errorf(
					"inconsistent local dirs. For dir entry %s found both %s and %s",
					key, value, existingValue)
			}
			localDirs[key] = value
		}
	}

	finalTarget := mts.Final.Target
	finalTargetConsole := b.console.WithPrefixAndSalt(finalTarget.String(), mts.Final.Salt)
	err = b.buildMain(ctx, localDirs, mts.Final)
	if err != nil {
		return "", nil, err
	}
	if opt.PrintSuccess {
		finalTargetConsole.Printf("Target %s built successfully\n", finalTarget.StringCanonical())
	}
	return cacheLocalDir, localDirs, nil
}

func (b *Builder) buildMain(ctx context.Context, localDirs map[string]string, states *states.SingleTarget) error {
	targetCtx := logging.With(ctx, "target", states.Target.String())
	solveCtx := logging.With(targetCtx, "solve", "side-effects")
	state := states.MainState
	if b.noCache {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	err := b.s.solveMain(solveCtx, localDirs, state)
	if err != nil {
		return errors.Wrapf(err, "solve side effects")
	}
	return nil
}

func (b *Builder) buildOutputs(ctx context.Context, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	targetCtx := logging.With(ctx, "target", states.Target.String())

	// Run --push commands.
	err := b.buildRunPush(targetCtx, localDirs, states, opt)
	if err != nil {
		return err
	}

	// Images.
	err = b.buildImages(targetCtx, localDirs, states, opt)
	if err != nil {
		return err
	}

	// Artifacts.
	if !states.Target.IsRemote() {
		// Don't output artifacts for remote images.
		err = b.buildArtifacts(targetCtx, localDirs, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildRunPush(ctx context.Context, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	if !states.RunPush.Initialized {
		// No run --push commands here. Quick way out.
		return nil
	}
	console := b.console.WithPrefixAndSalt(states.Target.String(), states.Salt)
	if !opt.Push {
		for _, commandStr := range states.RunPush.CommandStrs {
			console.Printf("Did not execute push command %s. Use earth --push to enable pushing\n", commandStr)
		}
		return nil
	}
	targetCtx := logging.With(ctx, "target", states.Target.String())
	solveCtx := logging.With(targetCtx, "solve", "run-push")
	err := b.s.solveMain(solveCtx, localDirs, states.RunPush.State)
	if err != nil {
		return errors.Wrapf(err, "solve run-push")
	}
	return nil
}

func (b *Builder) buildImages(ctx context.Context, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	for _, imageToSave := range states.SaveImages {
		if imageToSave.DockerTag == "" {
			// Not a docker export. Skip.
			continue
		}
		err := b.buildImage(ctx, imageToSave, localDirs, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildImage(ctx context.Context, imageToSave states.SaveImage, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	shouldPush := opt.Push && imageToSave.Push
	console := b.console.WithPrefixAndSalt(states.Target.String(), states.Salt)
	solveCtx := logging.With(ctx, "image", imageToSave.DockerTag)
	solveCtx = logging.With(solveCtx, "solve", "image")
	err := b.s.solveDocker(solveCtx, localDirs, imageToSave.State, imageToSave.Image, imageToSave.DockerTag, shouldPush)
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

func (b *Builder) buildImageTar(ctx context.Context, localDirs map[string]string, saveImage states.SaveImage, dockerTag string, outFile string) error {
	solveCtx := logging.With(ctx, "image", outFile)
	solveCtx = logging.With(solveCtx, "solve", "image-tar")
	err := b.s.solveDockerTar(solveCtx, localDirs, saveImage.State, saveImage.Image, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "solve image tar %s", outFile)
	}
	return nil
}

func (b *Builder) buildArtifacts(ctx context.Context, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	outDir, err := ioutil.TempDir(".", ".tmp-earth-out")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)
	solvedStates := make(map[int]bool)
	for _, artifactToSaveLocally := range states.SaveLocals {
		err = b.buildArtifact(
			ctx, artifactToSaveLocally, outDir, solvedStates, localDirs, states, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildSpecifiedArtifact(ctx context.Context, artifact domain.Artifact, destPath string, outDir string, solvedStates map[int]bool, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	solveCtx := logging.With(ctx, "solve", "artifacts")
	solveCtx = logging.With(solveCtx, "index", "combined")
	indexOutDir := filepath.Join(outDir, "combined")
	err := os.Mkdir(indexOutDir, 0755)
	if err != nil {
		return errors.Wrap(err, "mk index dir")
	}
	artifactsState := states.ArtifactsState
	err = b.s.solveArtifacts(solveCtx, localDirs, artifactsState, indexOutDir)
	if err != nil {
		return errors.Wrap(err, "solve combined artifacts")
	}

	err = b.saveArtifactLocally(ctx, artifact, indexOutDir, destPath, states.Salt, opt)
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) buildArtifact(ctx context.Context, artifactToSaveLocally states.SaveLocal, outDir string, solvedStates map[int]bool, localDirs map[string]string, states *states.SingleTarget, opt BuildOpt) error {
	index := artifactToSaveLocally.Index
	solveCtx := logging.With(ctx, "solve", "artifacts")
	solveCtx = logging.With(solveCtx, "index", index)
	indexOutDir := filepath.Join(outDir, fmt.Sprintf("index-%d", index))
	if !solvedStates[index] {
		solvedStates[index] = true
		artifactsState := states.SeparateArtifactsState[artifactToSaveLocally.Index]
		err := os.Mkdir(indexOutDir, 0755)
		if err != nil {
			return errors.Wrap(err, "mk index dir")
		}
		err = b.s.solveArtifacts(solveCtx, localDirs, artifactsState, indexOutDir)
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
	console := b.console.WithPrefixAndSalt(artifact.Target.String(), salt)
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

		logging.GetLogger(ctx).
			With("from", from).
			With("to", to).
			Info("Copying artifact")
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
