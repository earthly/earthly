package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/cmd/earthly/bk"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/regproxy"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/dockerutil"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/gwclientlogger"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/llbutil/secretprovider"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/saveartifactlocally"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/variables"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/moby/buildkit/util/entitlements"
	buildkitgitutil "github.com/moby/buildkit/util/gitutil"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	// PhaseInit is the phase text for the init phase.
	PhaseInit = "Init 🚀"
	// PhaseBuild is the phase text for the build phase.
	PhaseBuild = "Build 🔧"
	// PhasePush is the phase text for the push phase.
	PhasePush = "Push Summary ⏫"
	// PhaseOutput is the phase text for the output phase.
	PhaseOutput = "Local Output Summary 🎁"
)

// Opt represent builder options.
type Opt struct {
	BkClient                              *client.Client
	LogBusSolverMonitor                   *solvermon.SolverMonitor
	UseLogstream                          bool
	Console                               conslogging.ConsoleLogger
	Verbose                               bool
	Attachables                           []session.Attachable
	Enttlmnts                             []entitlements.Entitlement
	NoCache                               bool
	CacheImports                          *states.CacheImports
	CacheExport                           string
	MaxCacheExport                        string
	UseInlineCache                        bool
	SaveInlineCache                       bool
	ImageResolveMode                      llb.ResolveMode
	CleanCollection                       *cleanup.Collection
	OverridingVars                        *variables.Scope
	BuildContextProvider                  *provider.BuildContextProvider
	GitLookup                             *buildcontext.GitLookup
	GitBranchOverride                     string
	UseFakeDep                            bool
	Strict                                bool
	DisableNoOutputUpdates                bool
	ParallelConversion                    bool
	Parallelism                           semutil.Semaphore
	DarwinProxyImage                      string
	DarwinProxyWait                       time.Duration
	LocalRegistryAddr                     string
	DisableRemoteRegistryProxy            bool
	FeatureFlagOverrides                  string
	ContainerFrontend                     containerutil.ContainerFrontend
	InternalSecretStore                   *secretprovider.MutableMapStore
	InteractiveDebugging                  bool
	InteractiveDebuggingDebugLevelLogging bool
	GitImage                              string
	GitLFSInclude                         string
	GitLogLevel                           buildkitgitutil.GitLogLevel
	BuildkitSkipper                       bk.BuildkitSkipper
	NoAutoSkip                            bool
}

type ProjectAdder interface {
	AddProject(org, project string)
}

// BuildOpt is a collection of build options.
type BuildOpt struct {
	PlatformResolver           *platutil.Resolver
	AllowPrivileged            bool
	PrintPhases                bool
	Push                       bool
	CI                         bool
	NoOutput                   bool
	OnlyFinalTargetImages      bool
	OnlyArtifact               *domain.Artifact
	OnlyArtifactDestPath       string
	EnableGatewayClientLogging bool
	BuiltinArgs                variables.DefaultArgs
	GlobalWaitBlockFtr         bool
	LocalArtifactWhiteList     *gatewaycrafter.LocalArtifactWhiteList
	Logbus                     *logbus.Bus
	MainTargetDetailsFunc      func(earthfile2llb.TargetDetails) error
	Runner                     string
	ProjectAdder               ProjectAdder
	EarthlyCIRunner            bool
}

// Builder executes Earthly builds.
type Builder struct {
	s         *solver
	opt       Opt
	resolver  *buildcontext.Resolver
	builtMain bool

	outDirOnce sync.Once
	outDir     string
}

// NewBuilder returns a new earthly Builder.
func NewBuilder(ctx context.Context, opt Opt) (*Builder, error) {
	b := &Builder{
		s: &solver{
			sm:              outmon.NewSolverMonitor(opt.Console, opt.Verbose, opt.DisableNoOutputUpdates || opt.UseLogstream),
			logbusSM:        opt.LogBusSolverMonitor,
			useLogstream:    opt.UseLogstream,
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
	b.resolver = buildcontext.NewResolver(opt.CleanCollection, opt.GitLookup, opt.Console, opt.FeatureFlagOverrides, opt.GitBranchOverride, opt.GitLFSInclude, opt.GitLogLevel, opt.GitImage)
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

func (b *Builder) startRegistryProxy(ctx context.Context, caps apicaps.CapSet) (func(), bool) {
	cons := b.opt.Console.WithPrefix("registry-proxy")

	if b.opt.DisableRemoteRegistryProxy {
		cons.VerbosePrintf("Registry proxy disabled via --disable-remote-registry-proxy")
		return nil, false
	}

	if err := caps.Supports(pb.CapEarthlyRegistryProxy); err != nil {
		cons.Printf(err.Error())
		return nil, false
	}

	// Podman does not support the insecure localhost
	if b.opt.ContainerFrontend.Scheme() == "podman-container" {
		cons.Printf("Registry proxy not supported on Podman. Falling back to tar-based outputs.")
		return nil, false
	}

	useProxy, err := useSecondaryProxy()
	if err != nil {
		cons.Printf("Failed to check for registry proxy support: %v", err)
		return nil, false
	}

	controller := regproxy.NewController(
		b.s.bkClient.RegistryClient(),
		b.opt.ContainerFrontend,
		useProxy,
		b.opt.DarwinProxyImage,
		b.opt.DarwinProxyWait,
		cons,
	)
	addr, closeFn, err := controller.Start(ctx)
	if err != nil {
		cons.Printf("Failed to start registry proxy: %v", err)
		return nil, false
	}

	b.opt.LocalRegistryAddr = addr

	return closeFn, true
}

// useSecondaryProxy detects if we're on Mac (Darwin) or running on Windows in WSL2 or otherwise.
func useSecondaryProxy() (bool, error) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return true, nil
	}
	versionFile := "/proc/version"
	_, err := os.Stat(versionFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to stat %s", versionFile)
	}
	f, err := os.Open(versionFile)
	if err != nil {
		return false, errors.Wrapf(err, "failed to open %s", versionFile)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return false, errors.Wrapf(err, "failed to read %s", versionFile)
	}
	s := string(data)
	return strings.Contains(s, "WSL2"), nil
}

func (b *Builder) convertAndBuild(ctx context.Context, target domain.Target, opt BuildOpt) (*states.MultiTarget, error) {
	var (
		sharedLocalStateCache = earthfile2llb.NewSharedLocalStateCache()
		featureFlagOverrides  = b.opt.FeatureFlagOverrides
		manifestLists         = make(map[string][]dockerutil.Manifest) // parent image -> child images
		platformImgNames      = make(map[string]bool)                  // ensure that these are unique
		singPlatImgNames      = make(map[string]bool)                  // ensure that these are unique
		exportCoordinator     = gatewaycrafter.NewExportCoordinator()

		// dirIDs maps a dirIndex to a dirID; the "dir-id" field was introduced
		// to accommodate parallelism in the WAIT/END PopWaitBlock handling
		dirIDs = map[int]string{}
	)

	var (
		depIndex   = 0
		imageIndex = 0
		dirIndex   = 0
	)

	// Delay closing the registry proxy server until after the build function
	// returns. This can be deferred within the build function once global wait
	// block support is enabled.
	stopRegistryProxyFunc := func() {}
	defer func() {
		stopRegistryProxyFunc()
	}()

	var mts *states.MultiTarget
	buildFunc := func(childCtx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		if opt.EnableGatewayClientLogging {
			gwClient = gwclientlogger.New(gwClient)
		}
		var err error
		caps := gwClient.BuildOpts().LLBCaps

		if stopProxy, ok := b.startRegistryProxy(ctx, caps); ok {
			stopRegistryProxyFunc = stopProxy
		}

		if !b.builtMain {
			opt := earthfile2llb.ConvertOpt{
				GwClient:                             gwClient,
				Resolver:                             b.resolver,
				ImageResolveMode:                     b.opt.ImageResolveMode,
				CleanCollection:                      b.opt.CleanCollection,
				PlatformResolver:                     opt.PlatformResolver.SubResolver(opt.PlatformResolver.Current()),
				DockerImageSolverTar:                 newTarImageSolver(b.opt, b.s.sm),
				MultiImageSolver:                     newMultiImageSolver(b.opt, b.s.sm),
				OverridingVars:                       b.opt.OverridingVars,
				BuildContextProvider:                 b.opt.BuildContextProvider,
				CacheImports:                         b.opt.CacheImports,
				UseInlineCache:                       b.opt.UseInlineCache,
				UseFakeDep:                           b.opt.UseFakeDep,
				AllowLocally:                         !b.opt.Strict,
				AllowInteractive:                     !b.opt.Strict,
				AllowPrivileged:                      opt.AllowPrivileged,
				ParallelConversion:                   b.opt.ParallelConversion,
				Parallelism:                          b.opt.Parallelism,
				Console:                              b.opt.Console,
				GitLookup:                            b.opt.GitLookup,
				FeatureFlagOverrides:                 featureFlagOverrides,
				LocalStateCache:                      sharedLocalStateCache,
				BuiltinArgs:                          opt.BuiltinArgs,
				NoCache:                              b.opt.NoCache,
				ContainerFrontend:                    b.opt.ContainerFrontend,
				UseLocalRegistry:                     (b.opt.LocalRegistryAddr != ""),
				LocalRegistryAddr:                    b.opt.LocalRegistryAddr,
				DoSaves:                              !opt.NoOutput,
				OnlyFinalTargetImages:                opt.OnlyFinalTargetImages,
				DoPushes:                             opt.Push,
				IsCI:                                 opt.CI,
				EarthlyCIRunner:                      opt.EarthlyCIRunner,
				ExportCoordinator:                    exportCoordinator,
				LocalArtifactWhiteList:               opt.LocalArtifactWhiteList,
				InternalSecretStore:                  b.opt.InternalSecretStore,
				TempEarthlyOutDir:                    b.tempEarthlyOutDir,
				GlobalWaitBlockFtr:                   opt.GlobalWaitBlockFtr,
				LLBCaps:                              &caps,
				InteractiveDebuggerEnabled:           b.opt.InteractiveDebugging,
				InteractiveDebuggerDebugLevelLogging: b.opt.InteractiveDebuggingDebugLevelLogging,
				Logbus:                               opt.Logbus,
				MainTargetDetailsFunc:                opt.MainTargetDetailsFunc,
				Runner:                               opt.Runner,
				ProjectAdder:                         opt.ProjectAdder,
				FilesWithCommandRenameWarning:        make(map[string]bool),
				BuildkitSkipper:                      b.opt.BuildkitSkipper,
				NoAutoSkip:                           b.opt.NoAutoSkip,
			}
			mts, err = earthfile2llb.Earthfile2LLB(childCtx, target, opt, true)
			if err != nil {
				return nil, err
			}
		}

		//defer func() {
		//	sboms := earthfile2llb.GetSBOMs()
		//	sbomOutputPath := "earthly.sboms"
		//	f, err := os.Create(sbomOutputPath)
		//	if err != nil {
		//		b.opt.Console.Warnf("Error: failed to open %s for sbom output", sbomOutputPath)
		//		return
		//	}
		//	b.opt.Console.Printf("%d SBOM(s) were encountered during the build; contents have been writen to %s\n", len(sboms), sbomOutputPath)
		//	defer f.Close()
		//	for _, x := range sboms {
		//		f.WriteString(fmt.Sprintf("=== %s ===\n%s", x.Target, x.SBOM))
		//	}
		//}()

		if opt.GlobalWaitBlockFtr {
			if opt.OnlyArtifact != nil || opt.OnlyFinalTargetImages {
				b.opt.Console.Printf("builder.go bf code is still required for OnlyArtifact or OnlyFinalTargetImages modes (GlobalWaitBlockFtr has no effect)\n")
			} else {
				b.opt.Console.Printf("skipping builder.go bf code due to GlobalWaitBlockFtr\n")
				return nil, nil
			}
		}

		// WARNING: the code below is deprecated, and will eventually be removed, in favour of wait_block.go
		// This code is only used when dealing with VERSION 0.5 and 0.6; once these reach end-of-life, we can
		// delete the code below.

		// NOTE: this code is still required to support remote caching; it can't be removed until
		// https://github.com/earthly/earthly/issues/2178 is fixed.

		// *** DO NOT ADD CODE TO THE bf BELOW ***

		gwCrafter := gatewaycrafter.NewGatewayCrafter()
		if !b.builtMain {
			ref, err := b.stateToRef(childCtx, gwClient, mts.Final.MainState, mts.Final.PlatformResolver)
			if err != nil {
				return nil, err
			}
			gwCrafter.AddRef("main", ref)
		}
		if !opt.NoOutput && opt.OnlyArtifact != nil && !opt.OnlyFinalTargetImages {
			ref, err := b.stateToRef(childCtx, gwClient, mts.Final.ArtifactsState, mts.Final.PlatformResolver)
			if err != nil {
				return nil, err
			}
			refKey := "final-artifact"
			refPrefix := fmt.Sprintf("ref/%s", refKey)
			gwCrafter.AddRef(refKey, ref)
			gwCrafter.AddMeta(fmt.Sprintf("%s/export-dir", refPrefix), []byte("true"))
			gwCrafter.AddMeta(fmt.Sprintf("%s/final-artifact", refPrefix), []byte("true"))
		}

		isMultiPlatform := make(map[string]bool)    // DockerTag -> bool
		noManifestListImgs := make(map[string]bool) // DockerTag -> bool
		for _, sts := range mts.All() {
			if sts.PlatformResolver.Current() == platutil.DefaultPlatform {
				continue
			}
			for _, saveImage := range b.targetPhaseImages(sts) {
				doSaveOrPush := (sts.GetDoSaves() || sts.GetDoPushes() || saveImage.ForceSave)
				if !saveImage.SkipBuilder && saveImage.DockerTag != "" && doSaveOrPush {
					if saveImage.NoManifestList {
						noManifestListImgs[saveImage.DockerTag] = true
					} else {
						isMultiPlatform[saveImage.DockerTag] = true
					}
					if isMultiPlatform[saveImage.DockerTag] && noManifestListImgs[saveImage.DockerTag] {
						return nil, fmt.Errorf("cannot save image %s defined multiple times, but declared as SAVE IMAGE --no-manifest-list", saveImage.DockerTag)
					}
				}
			}
		}

		for _, sts := range mts.All() {
			hasRunPush := (sts.GetDoPushes() && sts.RunPush.HasState)
			if (sts.HasDangling && !b.opt.UseFakeDep) || (b.builtMain && hasRunPush) {
				depRef, err := b.stateToRef(childCtx, gwClient, b.targetPhaseState(sts), sts.PlatformResolver)
				if err != nil {
					return nil, err
				}
				refKey := fmt.Sprintf("dep-%d", depIndex)
				gwCrafter.AddRef(refKey, depRef)
				depIndex++
			}

			for _, saveImage := range b.targetPhaseImages(sts) {
				doSave := (sts.GetDoSaves() || saveImage.ForceSave)
				shouldExport := !opt.NoOutput && opt.OnlyArtifact == nil && !(opt.OnlyFinalTargetImages && sts != mts.Final) && saveImage.DockerTag != "" && doSave
				shouldPush := opt.Push && saveImage.Push && !sts.Target.IsRemote() && saveImage.DockerTag != "" && sts.GetDoPushes()
				useCacheHint := saveImage.CacheHint && b.opt.CacheExport != ""
				if (saveImage.SkipBuilder || !shouldPush && !shouldExport && !useCacheHint) || (!shouldPush && saveImage.HasPushDependencies) {
					// Short-circuit.
					continue
				}
				ref, err := b.stateToRef(childCtx, gwClient, saveImage.State, sts.PlatformResolver)
				if err != nil {
					return nil, err
				}

				if !isMultiPlatform[saveImage.DockerTag] {
					if saveImage.CheckDuplicate && saveImage.DockerTag != "" {
						if _, found := singPlatImgNames[saveImage.DockerTag]; found {
							return nil, errors.Errorf(
								"image %s is defined multiple times for the same default platform",
								saveImage.DockerTag)
						}
						singPlatImgNames[saveImage.DockerTag] = true
					}
					localRegPullID := exportCoordinator.AddImage(gwClient.BuildOpts().SessionID, saveImage.DockerTag, nil)
					refPrefix, err := gwCrafter.AddPushImageEntry(ref, imageIndex, saveImage.DockerTag, shouldPush, saveImage.InsecurePush, saveImage.Image, nil)
					if err != nil {
						return nil, err
					}
					imageIndex++

					if shouldExport {
						if b.opt.LocalRegistryAddr != "" {
							gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(localRegPullID))
						} else {
							gwCrafter.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
						}
					}
				} else {
					resolvedPlat := sts.PlatformResolver.Materialize(sts.PlatformResolver.Current())
					platformStr := resolvedPlat.String()
					platformImgName, err := llbutil.PlatformSpecificImageName(saveImage.DockerTag, resolvedPlat)
					if err != nil {
						return nil, err
					}
					if saveImage.CheckDuplicate && saveImage.DockerTag != "" {
						if _, found := platformImgNames[platformImgName]; found {
							return nil, errors.Errorf(
								"image %s is defined multiple times for the same platform (%s)",
								saveImage.DockerTag, platformImgName)
						}
						platformImgNames[platformImgName] = true
					}
					// Image has platform set - need to use manifest lists.
					// Need to push as a single multi-manifest image, but output locally as
					// separate images.
					// (docker load does not support tars with manifest lists).

					// For push.
					if shouldPush {
						_, err = gwCrafter.AddPushImageEntry(ref, imageIndex, saveImage.DockerTag, shouldPush, saveImage.InsecurePush, saveImage.Image, []byte(platformStr))
						if err != nil {
							return nil, err
						}
						imageIndex++
					}

					// For local.
					if shouldExport {
						refPrefix, err := gwCrafter.AddPushImageEntry(ref, imageIndex, platformImgName, false, false, saveImage.Image, nil)
						if err != nil {
							return nil, err
						}
						imageIndex++

						localRegPullID := exportCoordinator.AddImage(gwClient.BuildOpts().SessionID, platformImgName, nil)
						if b.opt.LocalRegistryAddr != "" {
							gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(localRegPullID))
						} else {
							gwCrafter.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
						}

						manifestLists[saveImage.DockerTag] = append(
							manifestLists[saveImage.DockerTag], dockerutil.Manifest{
								ImageName: platformImgName,
								Platform:  resolvedPlat,
							})
					}
				}
			}
			performSaveLocals := (!opt.NoOutput &&
				!opt.OnlyFinalTargetImages &&
				opt.OnlyArtifact == nil &&
				sts.GetDoSaves())
			if performSaveLocals {
				for _, saveLocal := range b.targetPhaseArtifacts(sts) {
					ref, err := b.artifactStateToRef(
						childCtx, gwClient, sts.SeparateArtifactsState[saveLocal.Index],
						sts.PlatformResolver)
					if err != nil {
						return nil, err
					}
					artifact := domain.Artifact{
						Target:   sts.Target,
						Artifact: saveLocal.ArtifactPath,
					}
					dirID, err := gwCrafter.AddSaveArtifactLocal(ref, dirIndex, artifact.String(), saveLocal.ArtifactPath, saveLocal.DestPath)
					if err != nil {
						return nil, err
					}
					dirIDs[dirIndex] = dirID
					opt.LocalArtifactWhiteList.Add(saveLocal.DestPath)
					dirIndex++
				}
			}

			targetInteractiveSession := b.targetPhaseInteractiveSession(sts)
			if targetInteractiveSession.Initialized && targetInteractiveSession.Kind == states.SessionEphemeral {
				ref, err := b.stateToRef(ctx, gwClient, targetInteractiveSession.State, sts.PlatformResolver)
				gwCrafter.AddRef("ephemeral", ref)
				if err != nil {
					return nil, err
				}
			}
		}
		return gwCrafter.GetResult(), nil
	}
	exportedTarImageManifestKeys := map[string]struct{}{}
	var exportedImagesMutex sync.Mutex
	onImageDone := func(manifestKey, waitFor string) error {
		exportedImagesMutex.Lock()
		defer exportedImagesMutex.Unlock()
		exportedTarImageManifestKeys[manifestKey] = struct{}{}
		manifests := make(map[string][]dockerutil.Manifest)
		for _, manifestKey := range strings.Split(waitFor, " ") {
			_, ok := exportedTarImageManifestKeys[manifestKey]
			if !ok {
				return nil
			}
			manifest, dockerTag, ok := exportCoordinator.GetImage(manifestKey)
			if !ok {
				return fmt.Errorf("failed to lookup %s in onImageDone", manifestKey)
			}
			manifests[dockerTag] = append(manifests[dockerTag], *manifest)
		}
		for parentImageName, children := range manifests {
			if opt.PlatformResolver == nil {
				panic("platform resolver is nil")
			}
			err := dockerutil.LoadDockerManifest(ctx, b.opt.Console, b.opt.ContainerFrontend, parentImageName, children, opt.PlatformResolver)
			if err != nil {
				return err
			}
		}
		return nil
	}
	onImage := func(childCtx context.Context, eg *errgroup.Group, imageName, waitFor, manifestKey string) (io.WriteCloser, error) {
		pipeR, pipeW := io.Pipe()
		eg.Go(func() error {
			defer pipeR.Close()
			err := dockerutil.LoadDockerTar(childCtx, b.opt.ContainerFrontend, pipeR)
			if err != nil {
				return errors.Wrapf(err, "load docker tar")
			}
			if manifestKey == "" {
				return nil
			}
			return onImageDone(manifestKey, waitFor)
		})
		return pipeW, nil
	}
	onArtifact := func(childCtx context.Context, index string, artifact domain.Artifact, artifactPath string, destPath string) (string, error) {
		if !opt.LocalArtifactWhiteList.Exists(destPath) {
			return "", errors.Errorf("dest path %s is not in the whitelist: %+v", destPath, opt.LocalArtifactWhiteList.AsList())
		}
		outDir, err := b.tempEarthlyOutDir()
		if err != nil {
			return "", err
		}
		artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%s", index))
		err = os.MkdirAll(artifactDir, 0755)
		if err != nil {
			return "", errors.Wrapf(err, "create dir %s", artifactDir)
		}
		return artifactDir, nil
	}
	onFinalArtifact := func(childCtx context.Context) (string, error) {
		return b.tempEarthlyOutDir()
	}
	onPull := func(childCtx context.Context, imagesToPull []string, resp map[string]string) error {
		if b.opt.LocalRegistryAddr == "" {
			return nil
		}
		manifests := make(map[string][]dockerutil.Manifest)
		pullMap := make(map[string]string)
		for _, imgToPull := range imagesToPull {
			manifest, dockerTag, ok := exportCoordinator.GetImage(imgToPull)
			if !ok {
				return errors.Errorf("unrecognized image to pull %s", imgToPull)
			}
			if manifest != nil {
				manifests[dockerTag] = append(manifests[dockerTag], *manifest)
				pullMap[imgToPull] = manifest.ImageName
			} else {
				pullMap[imgToPull] = dockerTag
			}
		}
		err := dockerutil.DockerPullLocalImages(childCtx, b.opt.ContainerFrontend, b.opt.LocalRegistryAddr, pullMap)
		if err != nil {
			return err
		}
		for parentImageName, children := range manifests {
			if opt.PlatformResolver == nil {
				panic("platform resolver is nil")
			}
			err = dockerutil.LoadDockerManifest(ctx, b.opt.Console, b.opt.ContainerFrontend, parentImageName, children, opt.PlatformResolver)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if opt.PrintPhases {
		b.opt.Console.PrintPhaseHeader(PhaseBuild, false, "")
	}
	err := b.s.buildMainMulti(ctx, buildFunc, onImage, onArtifact, onFinalArtifact, onPull, PhaseBuild, b.opt.Console)
	if err != nil {
		return nil, errors.Wrapf(err, "build main")
	}
	if opt.PrintPhases {
		b.opt.Console.PrintPhaseFooter(PhaseBuild, false, "")
	}
	b.builtMain = true

	if opt.PrintPhases {
		b.opt.Console.PrintPhaseHeader(PhasePush, !opt.Push, "")
		if !opt.Push {
			b.opt.Console.Printf("To enable pushing use earthly --push\n")
		}
	}
	if opt.Push && opt.OnlyArtifact == nil && !opt.OnlyFinalTargetImages {
		hasRunPush := false
		for _, sts := range mts.All() {
			if sts.GetDoPushes() && sts.RunPush.HasState {
				hasRunPush = true
				break
			}
		}
		if hasRunPush {
			err = b.s.buildMainMulti(ctx, buildFunc, onImage, onArtifact, onFinalArtifact, onPull, PhasePush, b.opt.Console)
			if err != nil {
				return nil, errors.Wrapf(err, "build push")
			}
		}
	}

	pushConsole := conslogging.NewBufferedLogger(&b.opt.Console)
	outputConsole := conslogging.NewBufferedLogger(&b.opt.Console)
	outputPhaseSpecial := ""

	if opt.NoOutput {
		// Nothing.
	} else if opt.OnlyArtifact != nil {
		if mts.Final.GetDoSaves() {
			outputPhaseSpecial = "single artifact"
			outDir, err := b.tempEarthlyOutDir()
			if err != nil {
				return nil, err
			}
			err = saveartifactlocally.SaveArtifactLocally(
				ctx, exportCoordinator, b.opt.Console, *opt.OnlyArtifact, outDir, opt.OnlyArtifactDestPath, mts.Final.ID, false)
			if err != nil {
				return nil, err
			}
		}
	} else if opt.OnlyFinalTargetImages {
		outputPhaseSpecial = "single image"
		for _, saveImage := range mts.Final.SaveImages {
			doSave := (mts.Final.GetDoSaves() || saveImage.ForceSave)
			shouldExport := !opt.NoOutput && saveImage.DockerTag != "" && doSave
			shouldPush := opt.Push && saveImage.Push && saveImage.DockerTag != "" && mts.Final.GetDoPushes()
			if saveImage.SkipBuilder || !shouldPush && !shouldExport {
				continue
			}

			if shouldPush {
				exportCoordinator.AddPushedImageSummary(mts.Final.Target.StringCanonical(), saveImage.DockerTag, b.opt.Console.Salt(), true)
			}
			if saveImage.Push && !opt.Push {
				exportCoordinator.AddPushedImageSummary(mts.Final.Target.StringCanonical(), saveImage.DockerTag, b.opt.Console.Salt(), false)
			}
			exportCoordinator.AddLocalOutputSummary(mts.Final.Target.StringCanonical(), saveImage.DockerTag, b.opt.Console.Salt())
		}
	} else {
		// This needs to match with the same index used during output.
		// TODO: This is a little brittle to future code changes.
		dirIndex := 0
		for _, sts := range mts.All() {
			for _, saveImage := range sts.SaveImages {
				doSave := (sts.GetDoSaves() || saveImage.ForceSave)
				shouldPush := opt.Push && saveImage.Push && !sts.Target.IsRemote() && saveImage.DockerTag != "" && sts.GetDoPushes()
				shouldExport := !opt.NoOutput && saveImage.DockerTag != "" && doSave
				if saveImage.SkipBuilder || !shouldPush && !shouldExport {
					continue
				}
				if shouldPush {
					exportCoordinator.AddPushedImageSummary(sts.Target.StringCanonical(), saveImage.DockerTag, sts.ID, true)
				}
				if saveImage.Push && !opt.Push && !sts.Target.IsRemote() {
					exportCoordinator.AddPushedImageSummary(sts.Target.StringCanonical(), saveImage.DockerTag, sts.ID, false)
				}
				exportCoordinator.AddLocalOutputSummary(sts.Target.StringCanonical(), saveImage.DockerTag, sts.ID)
			}
			if sts.GetDoSaves() {
				for _, saveLocal := range sts.SaveLocals {
					outDir, err := b.tempEarthlyOutDir()
					if err != nil {
						return nil, err
					}
					dirID, ok := dirIDs[dirIndex]
					if !ok {
						return nil, fmt.Errorf("failed to map dir index %d", dirIndex)
					}
					artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%s", dirID))
					artifact := domain.Artifact{
						Target:   sts.Target,
						Artifact: saveLocal.ArtifactPath,
					}
					err = saveartifactlocally.SaveArtifactLocally(
						ctx, exportCoordinator, b.opt.Console, artifact, artifactDir, saveLocal.DestPath, sts.ID, saveLocal.IfExists)
					if err != nil {
						return nil, err
					}
					dirIndex++
				}
			}

			if sts.GetDoSaves() && sts.RunPush.HasState {
				if opt.Push {
					for _, saveLocal := range sts.RunPush.SaveLocals {
						outDir, err := b.tempEarthlyOutDir()
						if err != nil {
							return nil, err
						}
						dirID, ok := dirIDs[dirIndex]
						if !ok {
							return nil, fmt.Errorf("failed to map dir index %d", dirIndex)
						}
						artifactDir := filepath.Join(outDir, fmt.Sprintf("index-%s", dirID))
						artifact := domain.Artifact{
							Target:   sts.Target,
							Artifact: saveLocal.ArtifactPath,
						}
						err = saveartifactlocally.SaveArtifactLocally(
							ctx, exportCoordinator, b.opt.Console, artifact, artifactDir, saveLocal.DestPath, sts.ID, saveLocal.IfExists)
						if err != nil {
							return nil, err
						}
						dirIndex++
					}
				} else {
					for _, commandStr := range sts.RunPush.CommandStrs {
						pushConsole.Printf("Did not execute push command %s\n", commandStr)
					}

					for _, saveImage := range sts.RunPush.SaveImages {
						pushConsole.Printf(
							"Did not push image %s as evaluating the image would "+
								"have caused a RUN --push to execute", saveImage.DockerTag)
						outputConsole.Printf("Did not output image %s locally, "+
							"as evaluating the image would have caused a "+
							"RUN --push to execute", saveImage.DockerTag)
					}

					if sts.RunPush.InteractiveSession.Initialized {
						pushConsole.Printf("Did not start an %s interactive session "+
							"with command %s\n", sts.RunPush.InteractiveSession.Kind,
							sts.RunPush.InteractiveSession.CommandStr)
					}
				}
			}
		}
	}

	for _, artifactEntry := range exportCoordinator.GetArtifactSummary() {
		console := b.opt.Console.WithPrefixAndSalt(artifactEntry.Target, artifactEntry.Salt)
		targetStr := console.PrefixColor().Sprintf("%s", artifactEntry.Target)
		outputConsole.Printf("Artifact %s output as %s\n", targetStr, artifactEntry.Path)
	}
	for _, outputEntry := range exportCoordinator.GetLocalOutputSummary() {
		console := b.opt.Console.WithPrefixAndSalt(outputEntry.Target, outputEntry.Salt)
		targetStr := console.PrefixColor().Sprintf("%s", outputEntry.Target)
		outputConsole.Printf("Image %s output as %s\n", targetStr, outputEntry.DockerTag)
	}
	for _, pushEntry := range exportCoordinator.GetPushedImageSummary() {
		console := b.opt.Console.WithPrefixAndSalt(pushEntry.Target, pushEntry.Salt)
		targetStr := console.PrefixColor().Sprintf("%s", pushEntry.Target)
		if pushEntry.Pushed {
			pushConsole.Printf("Pushed image %s as %s\n", targetStr, pushEntry.DockerTag)
		} else {
			pushConsole.Printf("Did not push image %s\n", pushEntry.DockerTag)
		}
	}

	pushConsole.Flush()
	if opt.PrintPhases {
		b.opt.Console.PrintPhaseFooter(PhasePush, !opt.Push, "")
		b.opt.Console.PrintPhaseHeader(PhaseOutput, opt.NoOutput, outputPhaseSpecial)
	}
	outputConsole.Flush()

	for parentImageName, children := range manifestLists {
		err = dockerutil.LoadDockerManifest(ctx, b.opt.Console, b.opt.ContainerFrontend, parentImageName, children, opt.PlatformResolver)
		if err != nil {
			return nil, err
		}
	}
	if opt.PrintPhases {
		b.opt.Console.PrintPhaseFooter(PhaseOutput, false, "")
		b.opt.Console.PrintSuccess()
	}
	return mts, nil
}

func (b *Builder) targetPhaseState(sts *states.SingleTarget) pllb.State {
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

func (b *Builder) targetPhaseImages(sts *states.SingleTarget) []states.SaveImage {
	if b.builtMain {
		return sts.RunPush.SaveImages
	}
	return sts.SaveImages
}

func (b *Builder) targetPhaseInteractiveSession(sts *states.SingleTarget) states.InteractiveSession {
	if b.builtMain {
		return sts.RunPush.InteractiveSession
	}
	return sts.InteractiveSession
}

func (b *Builder) stateToRef(ctx context.Context, gwClient gwclient.Client, state pllb.State, platr *platutil.Resolver) (gwclient.Reference, error) {
	noCache := b.opt.NoCache && !b.builtMain
	return llbutil.StateToRef(
		ctx, gwClient, state, noCache,
		platr, b.opt.CacheImports.AsSlice())
}

func (b *Builder) artifactStateToRef(ctx context.Context, gwClient gwclient.Client, state pllb.State, platr *platutil.Resolver) (gwclient.Reference, error) {
	noCache := b.opt.NoCache || b.builtMain
	return llbutil.StateToRef(
		ctx, gwClient, state, noCache,
		platr, b.opt.CacheImports.AsSlice())
}

func (b *Builder) tempEarthlyOutDir() (string, error) {
	var err error
	b.outDirOnce.Do(func() {
		tmpParentDir := ".tmp-earthly-out"
		err = os.MkdirAll(tmpParentDir, 0755)
		if err != nil {
			err = errors.Wrapf(err, "unable to create dir %s", tmpParentDir)
			return
		}
		b.outDir, err = os.MkdirTemp(tmpParentDir, "tmp")
		if err != nil {
			err = errors.Wrap(err, "mk temp dir for artifacts")
			return
		}
		b.opt.CleanCollection.Add(func() error {
			remErr := os.RemoveAll(b.outDir)
			// Remove the parent dir only if it's empty.
			_ = os.Remove(tmpParentDir)
			return remErr
		})
	})
	return b.outDir, err
}
