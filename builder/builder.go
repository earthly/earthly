package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/variables"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/entitlements"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
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
	CacheImports         map[string]bool
	CacheExport          string
	MaxCacheExport       string
	UseInlineCache       bool
	SaveInlineCache      bool
	ImageResolveMode     llb.ResolveMode
	CleanCollection      *cleanup.Collection
	VarCollection        *variables.Collection
	BuildContextProvider *provider.BuildContextProvider
	GitLookup            *buildcontext.GitLookup
	UseFakeDep           bool
}

// BuildOpt is a collection of build options.
type BuildOpt struct {
	Platform              *specs.Platform
	PrintSuccess          bool
	Push                  bool
	NoOutput              bool
	OnlyFinalTargetImages bool
	OnlyArtifact          *domain.Artifact
	OnlyArtifactDestPath  string
}

// Builder executes Earthly builds.
type Builder struct {
	s         *solver
	opt       Opt
	resolver  *buildcontext.Resolver
	builtMain bool
}

// NewBuilder returns a new earthly Builder.
func NewBuilder(ctx context.Context, opt Opt) (*Builder, error) {
	b := &Builder{
		s: &solver{
			sm:              newSolverMonitor(opt.Console, opt.Verbose),
			bkClient:        opt.BkClient,
			cacheImports:    opt.CacheImports,
			cacheExport:     opt.CacheExport,
			maxCacheExport:  opt.MaxCacheExport,
			attachables:     opt.Attachables,
			enttlmnts:       opt.Enttlmnts,
			saveInlineCache: opt.SaveInlineCache,
		},
		opt:      opt,
		resolver: nil, // initialized below
	}
	b.resolver = buildcontext.NewResolver(opt.SessionID, opt.CleanCollection, opt.GitLookup)
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

type successPrinter struct {
	printedOnce []sync.Once
	printFunc   []func()
	printIndex  int
}

func newSuccessPrinter(funcs ...func()) *successPrinter {
	printedOnce := []sync.Once{}
	for i := 0; i < len(funcs); i++ {
		printedOnce = append(printedOnce, sync.Once{})
	}

	return &successPrinter{printedOnce, funcs, 0}
}

func (sp *successPrinter) printCurrentSuccess() {
	sp.printedOnce[sp.printIndex].Do(sp.printFunc[sp.printIndex])
}

func (sp *successPrinter) incrementIndex() {
	sp.printIndex++
}

func (b *Builder) convertAndBuild(ctx context.Context, target domain.Target, opt BuildOpt) (*states.MultiTarget, error) {
	outDir, err := ioutil.TempDir(".", ".tmp-earthly-out")
	if err != nil {
		return nil, errors.Wrap(err, "mk temp dir for artifacts")
	}
	defer os.RemoveAll(outDir)

	successFun := func(msg string) func() {
		return func() {
			if opt.PrintSuccess {
				b.s.sm.SetSuccess(msg)
			}
		}
	}
	sp := newSuccessPrinter(successFun(""), successFun("--push"))

	destPathWhitelist := make(map[string]bool)
	manifestLists := make(map[string][]manifest) // parent image -> child images
	var mts *states.MultiTarget
	depIndex := 0
	imageIndex := 0
	dirIndex := 0
	bf := func(ctx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		var err error
		mts, err = earthfile2llb.Earthfile2LLB(ctx, target, earthfile2llb.ConvertOpt{
			GwClient:             gwClient,
			Resolver:             b.resolver,
			ImageResolveMode:     b.opt.ImageResolveMode,
			DockerBuilderFun:     b.MakeImageAsTarBuilderFun(),
			CleanCollection:      b.opt.CleanCollection,
			Platform:             opt.Platform,
			VarCollection:        b.opt.VarCollection,
			BuildContextProvider: b.opt.BuildContextProvider,
			CacheImports:         b.opt.CacheImports,
			UseInlineCache:       b.opt.UseInlineCache,
			UseFakeDep:           b.opt.UseFakeDep,
		})
		if err != nil {
			return nil, err
		}
		res := gwclient.NewResult()
		ref, err := b.stateToRef(ctx, gwClient, b.targetPhaseState(mts.Final), mts.Final.Platform)
		if err != nil {
			return nil, err
		}
		res.AddRef("main", ref)
		if !opt.NoOutput && opt.OnlyArtifact != nil && !opt.OnlyFinalTargetImages {
			ref, err = b.stateToRef(ctx, gwClient, mts.Final.ArtifactsState, mts.Final.Platform)
			if err != nil {
				return nil, err
			}
			refKey := "final-artifact"
			refPrefix := fmt.Sprintf("ref/%s", refKey)
			res.AddRef(refKey, ref)
			res.AddMeta(fmt.Sprintf("%s/export-dir", refPrefix), []byte("true"))
			res.AddMeta(fmt.Sprintf("%s/final-artifact", refPrefix), []byte("true"))
		}

		for _, sts := range mts.All() {
			if (sts.HasDangling && !b.opt.UseFakeDep) || b.builtMain {
				depRef, err := b.stateToRef(ctx, gwClient, b.targetPhaseState(sts), sts.Platform)
				if err != nil {
					return nil, err
				}
				refKey := fmt.Sprintf("dep-%d", depIndex)
				res.AddRef(refKey, depRef)
				depIndex++
			}

			for _, saveImage := range sts.SaveImages {
				shouldPush := opt.Push && saveImage.Push && !sts.Target.IsRemote() && saveImage.DockerTag != ""
				shouldExport := !opt.NoOutput && opt.OnlyArtifact == nil && !(opt.OnlyFinalTargetImages && sts != mts.Final) && saveImage.DockerTag != ""
				useCacheHint := saveImage.CacheHint && b.opt.CacheExport != ""
				if !shouldPush && !shouldExport && !useCacheHint {
					// Short-circuit.
					continue
				}
				ref, err := b.stateToRef(ctx, gwClient, saveImage.State, sts.Platform)
				if err != nil {
					return nil, err
				}
				config, err := json.Marshal(saveImage.Image)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal save image config")
				}

				if sts.Platform == nil {
					refKey := fmt.Sprintf("image-%d", imageIndex)
					refPrefix := fmt.Sprintf("ref/%s", refKey)
					imageIndex++

					res.AddMeta(fmt.Sprintf("%s/image.name", refPrefix), []byte(saveImage.DockerTag))
					if shouldPush {
						res.AddMeta(fmt.Sprintf("%s/export-image-push", refPrefix), []byte("true"))
						if saveImage.InsecurePush {
							res.AddMeta(fmt.Sprintf("%s/insecure-push", refPrefix), []byte("true"))
						}
					}
					res.AddMeta(fmt.Sprintf("%s/%s", refPrefix, exptypes.ExporterImageConfigKey), config)
					if shouldExport {
						res.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
					}
					res.AddMeta(fmt.Sprintf("%s/image-index", refPrefix), []byte(fmt.Sprintf("%d", imageIndex)))
					res.AddRef(refKey, ref)
				} else {
					// Image has platform set - need to use manifest lists.
					// Need to push as a single multi-manifest image, but output locally as
					// separate images.
					// (docker load does not support tars with manifest lists).

					// For push.
					if shouldPush {
						refKey := fmt.Sprintf("image-%d", imageIndex)
						refPrefix := fmt.Sprintf("ref/%s", refKey)
						imageIndex++

						res.AddMeta(fmt.Sprintf("%s/image.name", refPrefix), []byte(saveImage.DockerTag))
						res.AddMeta(fmt.Sprintf("%s/platform", refPrefix), []byte(llbutil.PlatformToString(sts.Platform)))
						res.AddMeta(fmt.Sprintf("%s/export-image-push", refPrefix), []byte("true"))
						if saveImage.InsecurePush {
							res.AddMeta(fmt.Sprintf("%s/insecure-push", refPrefix), []byte("true"))
						}
						res.AddMeta(fmt.Sprintf("%s/%s", refPrefix, exptypes.ExporterImageConfigKey), config)
						res.AddMeta(fmt.Sprintf("%s/image-index", refPrefix), []byte(fmt.Sprintf("%d", imageIndex)))
						res.AddRef(refKey, ref)
					}

					// For local.
					if shouldExport {
						refKey := fmt.Sprintf("image-%d", imageIndex)
						refPrefix := fmt.Sprintf("ref/%s", refKey)
						imageIndex++

						platformImgName, err := platformSpecificImageName(saveImage.DockerTag, *sts.Platform)
						if err != nil {
							return nil, err
						}
						res.AddMeta(fmt.Sprintf("%s/image.name", refPrefix), []byte(platformImgName))
						res.AddMeta(fmt.Sprintf("%s/%s", refPrefix, exptypes.ExporterImageConfigKey), config)
						res.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
						res.AddMeta(fmt.Sprintf("%s/image-index", refPrefix), []byte(fmt.Sprintf("%d", imageIndex)))
						res.AddRef(refKey, ref)
						manifestLists[saveImage.DockerTag] = append(
							manifestLists[saveImage.DockerTag], manifest{
								imageName: platformImgName,
								platform:  *sts.Platform,
							})
					}
				}
			}
			if !sts.Target.IsRemote() && !opt.NoOutput && !opt.OnlyFinalTargetImages && opt.OnlyArtifact == nil {
				for _, saveLocal := range b.targetPhaseArtifacts(sts) {
					ref, err := b.artifactStateToRef(ctx, gwClient, sts.SeparateArtifactsState[saveLocal.Index], sts.Platform)
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
					res.AddMeta(fmt.Sprintf("%s/export-dir", refPrefix), []byte("true"))
					res.AddMeta(fmt.Sprintf("%s/dir-index", refPrefix), []byte(fmt.Sprintf("%d", dirIndex)))
					destPathWhitelist[saveLocal.DestPath] = true
					dirIndex++
				}
			}
		}
		return res, nil
	}
	onImage := func(ctx context.Context, eg *errgroup.Group, imageName string) (io.WriteCloser, error) {
		sp.printCurrentSuccess()
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
		sp.printCurrentSuccess()
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
		sp.printCurrentSuccess()
		return outDir, nil
	}
	err = b.s.buildMainMulti(ctx, bf, onImage, onArtifact, onFinalArtifact, "main")
	if err != nil {
		return nil, errors.Wrapf(err, "build main")
	}
	sp.printCurrentSuccess()
	sp.incrementIndex()
	b.builtMain = true

	if opt.NoOutput {
		// Nothing.
	} else if opt.OnlyArtifact != nil {
		err := b.saveArtifactLocally(ctx, *opt.OnlyArtifact, outDir, opt.OnlyArtifactDestPath, mts.Final.Salt, opt, false)
		if err != nil {
			return nil, err
		}
	} else if opt.OnlyFinalTargetImages {
		for _, saveImage := range mts.Final.SaveImages {
			shouldPush := opt.Push && saveImage.Push && saveImage.DockerTag != ""
			shouldExport := !opt.NoOutput && saveImage.DockerTag != ""
			if !shouldPush && !shouldExport {
				continue
			}
			console := b.opt.Console.WithPrefixAndSalt(mts.Final.Target.String(), mts.Final.Salt)
			pushStr := ""
			if shouldPush {
				pushStr = " (pushed)"
			}
			console.Printf("Image %s as %s%s\n", mts.Final.Target.StringCanonical(), saveImage.DockerTag, pushStr)
			if saveImage.Push && !opt.Push {
				console.Printf("Did not push %s. Use earthly --push to enable pushing\n", saveImage.DockerTag)
			}
		}
	} else {
		// This needs to match with the same index used during output.
		// TODO: This is a little brittle to future code changes.
		dirIndex := 0
		for _, sts := range mts.All() {
			console := b.opt.Console.WithPrefixAndSalt(sts.Target.String(), sts.Salt)

			for _, saveImage := range sts.SaveImages {
				shouldPush := opt.Push && saveImage.Push && !sts.Target.IsRemote() && saveImage.DockerTag != ""
				shouldExport := !opt.NoOutput && saveImage.DockerTag != ""
				if !shouldPush && !shouldExport {
					continue
				}
				pushStr := ""
				if shouldPush {
					pushStr = " (pushed)"
				}
				console.Printf("Image %s as %s%s\n", sts.Target.StringCanonical(), saveImage.DockerTag, pushStr)

				if saveImage.Push && !opt.Push && !sts.Target.IsRemote() {
					console.Printf("Did not push %s. Use earthly --push to enable pushing\n", saveImage.DockerTag)
				}
			}
			if !sts.Target.IsRemote() {
				for _, saveLocal := range sts.SaveLocals {
					artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%d", dirIndex))
					artifact := domain.Artifact{
						Target:   sts.Target,
						Artifact: saveLocal.ArtifactPath,
					}
					err := b.saveArtifactLocally(ctx, artifact, artifactDir, saveLocal.DestPath, sts.Salt, opt, saveLocal.IfExists)
					if err != nil {
						return nil, err
					}
					dirIndex++
				}

				if sts.RunPush.Initialized {
					if opt.Push {
						err = b.s.buildMainMulti(ctx, bf, onImage, onArtifact, onFinalArtifact, "--push")
						if err != nil {
							return nil, errors.Wrapf(err, "build push")
						}
						sp.printCurrentSuccess()

						for _, saveLocal := range sts.RunPush.SaveLocals {
							artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%d", dirIndex))
							artifact := domain.Artifact{
								Target:   sts.Target,
								Artifact: saveLocal.ArtifactPath,
							}
							err := b.saveArtifactLocally(ctx, artifact, artifactDir, saveLocal.DestPath, sts.Salt, opt, saveLocal.IfExists)
							if err != nil {
								return nil, err
							}
							dirIndex++
						}
					} else {
						for _, commandStr := range sts.RunPush.CommandStrs {
							console.Printf("Did not execute push command %s. Use earthly --push to enable pushing\n", commandStr)
						}
					}
				}
			}
		}
	}
	for parentImageName, children := range manifestLists {
		err = loadDockerManifest(ctx, b.opt.Console, parentImageName, children)
		if err != nil {
			return nil, err
		}
	}

	return mts, nil
}

func (b *Builder) targetPhaseState(sts *states.SingleTarget) llb.State {
	if b.builtMain {
		return sts.RunPush.State
	}
	return sts.MainState
}

func (b *Builder) targetPhaseArtifacts(sts *states.SingleTarget) []states.SaveLocal {
	if b.builtMain {
		return sts.RunPush.SaveLocals
	}
	return sts.SaveLocals
}

func (b *Builder) stateToRef(ctx context.Context, gwClient gwclient.Client, state llb.State, platform *specs.Platform) (gwclient.Reference, error) {
	if b.opt.NoCache && !b.builtMain {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	return llbutil.StateToRef(ctx, gwClient, state, platform, b.opt.CacheImports)
}

func (b *Builder) artifactStateToRef(ctx context.Context, gwClient gwclient.Client, state llb.State, platform *specs.Platform) (gwclient.Reference, error) {
	if b.opt.NoCache || b.builtMain {
		state = state.SetMarshalDefaults(llb.IgnoreCache)
	}
	return llbutil.StateToRef(ctx, gwClient, state, platform, b.opt.CacheImports)
}

func (b *Builder) buildOnlyLastImageAsTar(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string, opt BuildOpt) error {
	saveImage := mts.Final.LastSaveImage()
	err := b.buildMain(ctx, mts, opt)
	if err != nil {
		return err
	}

	platform, err := llbutil.ResolvePlatform(mts.Final.Platform, opt.Platform)
	if err != nil {
		platform = mts.Final.Platform
	}
	plat := llbutil.PlatformWithDefault(platform)
	err = b.outputImageTar(ctx, saveImage, plat, dockerTag, outFile)
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
	platform, err := llbutil.ResolvePlatform(mts.Final.Platform, opt.Platform)
	if err != nil {
		platform = mts.Final.Platform
	}
	plat := llbutil.PlatformWithDefault(platform)
	err = b.s.solveMain(ctx, state, plat)
	if err != nil {
		return errors.Wrapf(err, "solve side effects")
	}
	return nil
}

func (b *Builder) executeRunPush(ctx context.Context, sts *states.SingleTarget, opt BuildOpt) error {
	if !sts.RunPush.Initialized {
		// No run --push commands here. Quick way out.
		return nil
	}
	console := b.opt.Console.WithPrefixAndSalt(sts.Target.String(), sts.Salt)
	if !opt.Push {
		for _, commandStr := range sts.RunPush.CommandStrs {
			console.Printf("Did not execute push command %s. Use earthly --push to enable pushing\n", commandStr)
		}
		return nil
	}
	platform, err := llbutil.ResolvePlatform(sts.Platform, opt.Platform)
	if err != nil {
		platform = sts.Platform
	}
	plat := llbutil.PlatformWithDefault(platform)
	err = b.s.solveMain(ctx, sts.RunPush.State, plat)
	if err != nil {
		return errors.Wrapf(err, "solve run-push")
	}
	return nil
}

func (b *Builder) outputImageTar(ctx context.Context, saveImage states.SaveImage, platform specs.Platform, dockerTag string, outFile string) error {
	err := b.s.solveDockerTar(ctx, saveImage.State, platform, saveImage.Image, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "solve image tar %s", outFile)
	}
	return nil
}

func (b *Builder) saveArtifactLocally(ctx context.Context, artifact domain.Artifact, indexOutDir string, destPath string, salt string, opt BuildOpt, ifExists bool) error {
	console := b.opt.Console.WithPrefixAndSalt(artifact.Target.String(), salt)
	fromPattern := filepath.Join(indexOutDir, filepath.FromSlash(artifact.Artifact))
	// Resolve possible wildcards.
	// TODO: Note that this is not very portable, as the glob is host-platform dependent,
	//       while the pattern is also guest-platform dependent.
	fromGlobMatches, err := filepath.Glob(fromPattern)
	if err != nil {
		return errors.Wrapf(err, "glob")
	} else if !artifact.Target.IsRemote() && len(fromGlobMatches) <= 0 {
		if ifExists {
			return nil
		}
		return fmt.Errorf("cannot save artifact %s, since it does not exist", artifact.StringCanonical())
	}
	isWildcard := strings.ContainsAny(fromPattern, `*?[`)
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
		}
		if destIsDir {
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
