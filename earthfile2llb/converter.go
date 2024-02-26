package earthfile2llb

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alessio/shellescape"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution/reference"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast/commandflag"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/inputgraph"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/inodeutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/shell"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/earthly/earthly/variables"
	"github.com/earthly/earthly/variables/reserved"
	"github.com/moby/buildkit/client/llb"
	dockerimage "github.com/moby/buildkit/exporter/containerimage/image"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/dockerui"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session/localhost"
	solverpb "github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cmdType int

const (
	argCmd            cmdType = iota + 1 // "ARG"
	buildCmd                             // "BUILD"
	cmdCmd                               // "CMD"
	copyCmd                              // "COPY"
	enterScopeDoCmd                      // "ENTER-SCOPE-DO"
	entrypointCmd                        // "ENTRYPOINT"
	envCmd                               // "ENV"
	exposeCmd                            // "EXPOSE"
	fromCmd                              // "FROM"
	fromDockerfileCmd                    // "FROM DOCKERFILE"
	gitCloneCmd                          // "GIT CLONE"
	healthcheckCmd                       // "HEALTHCHECK"
	importCmd                            // "IMPORT"
	labelCmd                             // "LABEL"
	loadCmd                              // "LOAD"
	locallyCmd                           // "LOCALLY"
	runCmd                               // "RUN"
	saveArtifactCmd                      // "SAVE ARTIFACT"
	saveImageCmd                         // "SAVE IMAGE"
	userCmd                              // "USER"
	volumeCmd                            // "VOLUME"
	workdirCmd                           // "WORKDIR"
	cacheCmd                             // "CACHE"
	hostCmd                              // "HOST"
	projectCmd                           // "PROJECT"
	pipelineCmd                          // "PIPELINE"
	triggerCmd                           // "TRIGGER"
	setCmd                               // "SET"
	letCmd                               // "LET"
)

// Converter turns earthly commands to buildkit LLB representation.
type Converter struct {
	target              domain.Target
	gitMeta             *gitutil.GitMetadata
	platr               *platutil.Resolver
	opt                 ConvertOpt
	mts                 *states.MultiTarget
	directDeps          []*states.SingleTarget
	buildContextFactory llbfactory.Factory
	cacheContext        pllb.State
	persistentCacheDirs map[string]states.CacheMount // maps path->mount
	varCollection       *variables.Collection
	ranSave             bool
	cmdSet              bool
	ftrs                *features.Features
	localWorkingDir     string
	containerFrontend   containerutil.ContainerFrontend
	waitBlockStack      []*waitBlock
	isPipeline          bool
	logbusTarget        *logbus.Target
	nextCmdID           int
}

// NewConverter constructs a new converter for a given earthly target.
func NewConverter(ctx context.Context, target domain.Target, bc *buildcontext.Data, sts *states.SingleTarget, opt ConvertOpt) (*Converter, error) {
	opt.BuildContextProvider.AddDirs(bc.LocalDirs)
	sts.HasDangling = opt.HasDangling
	mts := &states.MultiTarget{
		Final:   sts,
		Visited: opt.Visited,
	}
	newCollOpt := variables.NewCollectionOpt{
		Console:          opt.Console,
		Target:           target,
		Push:             opt.DoPushes,
		CI:               opt.IsCI,
		EarthlyCIRunner:  opt.EarthlyCIRunner,
		PlatformResolver: opt.PlatformResolver,
		GitMeta:          bc.GitMetadata,
		BuiltinArgs:      opt.BuiltinArgs,
		OverridingVars:   opt.OverridingVars,
		GlobalImports:    opt.GlobalImports,
		Features:         opt.Features,
	}
	ovVarsKeysSorted := opt.OverridingVars.Sorted()
	ovVars := make([]string, 0, len(ovVarsKeysSorted))
	for _, k := range ovVarsKeysSorted {
		v, _ := opt.OverridingVars.Get(k)
		ovVars = append(ovVars, fmt.Sprintf("%s=%s", k, v))
	}

	run := opt.Logbus.Run()

	logbusTarget, err := run.NewTarget(
		sts.ID,
		target,
		ovVars,
		opt.PlatformResolver.Current().String(),
		opt.Runner,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new logbus target")
	}

	logbusTarget.SetStart(time.Now())

	c := &Converter{
		target:              target,
		gitMeta:             bc.GitMetadata,
		platr:               opt.PlatformResolver,
		opt:                 opt,
		mts:                 mts,
		buildContextFactory: bc.BuildContextFactory,
		cacheContext:        pllb.Scratch(),
		persistentCacheDirs: make(map[string]states.CacheMount),
		varCollection:       variables.NewCollection(newCollOpt),
		ftrs:                bc.Features,
		localWorkingDir:     filepath.Dir(bc.BuildFilePath),
		containerFrontend:   opt.ContainerFrontend,
		waitBlockStack:      []*waitBlock{opt.waitBlock},
		logbusTarget:        logbusTarget,
	}

	if c.opt.GlobalWaitBlockFtr {
		c.ftrs.WaitBlock = true
	}
	return c, nil
}

// From applies the earthly FROM command.
func (c *Converter) From(ctx context.Context, imageName string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string) error {
	err := c.checkAllowed(fromCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	if len(c.persistentCacheDirs) > 0 {
		c.persistentCacheDirs = make(map[string]states.CacheMount)
	}
	c.cmdSet = false
	err = c.checkOldPlatformIncompatibility(platform)
	if err != nil {
		return err
	}
	c.varCollection.SetLocally(false) // FIXME this will have to change once https://github.com/earthly/earthly/issues/2044 is fixed
	platform = c.setPlatform(platform)
	if strings.Contains(imageName, "+") {
		// Target-based FROM.
		return c.fromTarget(ctx, imageName, platform, allowPrivileged, passArgs, buildArgs)
	}

	// Docker image based FROM.
	if len(buildArgs) != 0 {
		return errors.New("--build-arg not supported in non-target FROM")
	}
	return c.fromClassical(ctx, imageName, platform, false)
}

func (c *Converter) fromClassical(ctx context.Context, imageName string, platform platutil.Platform, local bool) error {
	var internal bool
	if local {
		// local mode uses a fake image containing /bin/true
		// we want to prefix this as internal so it doesn't show up in the output
		internal = true
	} else {
		internal = false
	}
	prefix, _, err := c.newVertexMeta(ctx, local, false, internal, nil)
	if err != nil {
		return err
	}
	state, img, envVars, err := c.internalFromClassical(
		ctx, imageName, platform,
		llb.WithCustomNamef("%sFROM %s", prefix, imageName))
	if err != nil {
		return err
	}
	c.mts.Final.MainState = state
	c.mts.Final.MainImage = img
	c.mts.Final.RanFromLike = true
	c.varCollection.ResetEnvVars(envVars)
	return nil
}

func (c *Converter) fromTarget(ctx context.Context, targetName string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string) (retErr error) {
	cmdID, cmd, err := c.newLogbusCommand(ctx, fmt.Sprintf("FROM %s", targetName))
	if err != nil {
		return errors.Wrap(err, "failed to create command")
	}

	defer func() {
		cmd.SetEndError(retErr)
	}()

	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}

	mts, err := c.buildTarget(ctx, depTarget.String(), platform, allowPrivileged, passArgs, buildArgs, false, fromCmd, cmdID)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", depTarget.String())
	}

	if mts.Final.RanInteractive {
		return errors.New("cannot FROM a target ending with an --interactive")
	}

	if depTarget.IsLocalInternal() {
		depTarget.LocalPath = c.mts.Final.Target.LocalPath
	}

	// Look for the built state in the dep states, after we've built it.
	relevantDepState := mts.Final
	saveImage := relevantDepState.LastSaveImage()

	// Pass on dep state over to this state.
	c.mts.Final.MainState = relevantDepState.MainState
	c.varCollection.ResetEnvVars(mts.Final.VarCollection.EnvVars())
	c.mts.Final.MainImage = saveImage.Image.Clone()
	c.mts.Final.RanFromLike = mts.Final.RanFromLike
	c.mts.Final.RanInteractive = mts.Final.RanInteractive
	c.platr.UpdatePlatform(mts.Final.PlatformResolver.Current())

	return nil
}

// FromDockerfile applies the earthly FROM DOCKERFILE command.
func (c *Converter) FromDockerfile(ctx context.Context, contextPath string, dfPath string, dfTarget string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) (retErr error) {
	var err error
	ctx, err = c.ftrs.WithContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add feature flags to context")
	}
	err = c.checkAllowed(fromDockerfileCmd)
	if err != nil {
		return err
	}
	err = c.checkOldPlatformIncompatibility(platform)
	if err != nil {
		return err
	}
	platform = c.setPlatform(platform)
	plat := c.platr.ToLLBPlatform(platform)
	c.nonSaveCommand()
	cmdID, cmd, err := c.newLogbusCommand(ctx, fmt.Sprintf("FROM DOCKERFILE %s", dfPath))
	if err != nil {
		return errors.Wrap(err, "failed to create command")
	}
	defer func() {
		cmd.SetEndError(retErr)
	}()
	var dfData []byte
	if dfPath != "" {
		dfArtifact, parseErr := domain.ParseArtifact(dfPath)
		if parseErr == nil {
			// The Dockerfile is from a target's artifact.
			mts, err := c.buildTarget(ctx, dfArtifact.Target.String(), platform, allowPrivileged, false, buildArgs, false, fromDockerfileCmd, cmdID)
			if err != nil {
				return err
			}
			dfData, err = c.readArtifact(ctx, mts, dfArtifact)
			if err != nil {
				return err
			}
		} else {
			// The Dockerfile is from the host.
			dockerfileMetaTarget := domain.Target{
				Target:    fmt.Sprintf("%s%s", buildcontext.DockerfileMetaTarget, path.Base(dfPath)),
				LocalPath: path.Dir(dfPath),
			}
			dockerfileMetaTargetRef, err := c.joinRefs(dockerfileMetaTarget)
			if err != nil {
				return errors.Wrap(err, "join targets")
			}
			dockerfileMetaTarget = dockerfileMetaTargetRef.(domain.Target)
			data, err := c.opt.Resolver.Resolve(ctx, c.opt.GwClient, c.platr, dockerfileMetaTarget)
			if err != nil {
				return errors.Wrap(err, "resolve build context for dockerfile")
			}
			c.opt.BuildContextProvider.AddDirs(data.LocalDirs)
			dfData, err = os.ReadFile(data.BuildFilePath)
			if err != nil {
				return errors.Wrapf(err, "read file %s", data.BuildFilePath)
			}
		}
	}
	var BuildContextFactory llbfactory.Factory
	contextArtifact, parseErr := domain.ParseArtifact(contextPath)
	if parseErr == nil {
		prefix, cmdID, err := c.newVertexMeta(ctx, false, false, true, nil)
		if err != nil {
			return err
		}
		// The build context is from a target's artifact.
		// TODO: The build args are used for both the artifact and the Dockerfile. This could be
		//       confusing to the user.
		mts, err := c.buildTarget(ctx, contextArtifact.Target.String(), platform, allowPrivileged, false, buildArgs, false, fromDockerfileCmd, cmdID)
		if err != nil {
			return err
		}
		if dfPath == "" {
			// Imply dockerfile as being ./Dockerfile in the root of the build context.
			dfArtifact := contextArtifact
			dfArtifact.Artifact = path.Join(dfArtifact.Artifact, "Dockerfile")
			dfData, err = c.readArtifact(ctx, mts, dfArtifact)
			if err != nil {
				return err
			}
		}
		copyState, err := llbutil.CopyOp(ctx,
			mts.Final.ArtifactsState, []string{contextArtifact.Artifact},
			c.platr.Scratch(), "/", true, true, false, "", nil, false, false,
			c.ftrs.UseCopyLink,
			llb.WithCustomNamef(
				"%sFROM DOCKERFILE (copy build context from) %s%s",
				prefix,
				joinWrap(buildArgs, "(", " ", ") "), contextArtifact.String()))
		if err != nil {
			return errors.Wrapf(err, "copyOp FROM DOCKERFILE")
		}
		BuildContextFactory = llbfactory.PreconstructedState(copyState)
	} else {
		// The build context is from the host.
		if contextPath != "." &&
			!strings.HasPrefix(contextPath, "./") &&
			!strings.HasPrefix(contextPath, "../") &&
			!strings.HasPrefix(contextPath, "/") {
			contextPath = fmt.Sprintf("./%s", contextPath)
		}
		dockerfileMetaTarget := domain.Target{
			Target:    fmt.Sprintf("%s%s", buildcontext.DockerfileMetaTarget, stringutil.StrOrDefault(dfPath, "Dockerfile")),
			LocalPath: path.Join(contextPath),
		}
		dockerfileMetaTargetRef, err := c.joinRefs(dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "join targets")
		}
		dockerfileMetaTarget = dockerfileMetaTargetRef.(domain.Target)
		data, err := c.opt.Resolver.Resolve(ctx, c.opt.GwClient, c.platr, dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "resolve build context for dockerfile")
		}
		c.opt.BuildContextProvider.AddDirs(data.LocalDirs)
		if dfPath == "" {
			// Imply dockerfile as being ./Dockerfile in the root of the build context.
			dfData, err = os.ReadFile(data.BuildFilePath)
			if err != nil {
				return errors.Wrapf(err, "read file %s", data.BuildFilePath)
			}
		}
		BuildContextFactory = data.BuildContextFactory
	}
	bc, err := dockerui.NewClient(c.opt.GwClient)
	if err != nil {
		return errors.Wrap(err, "dockerui.NewClient")
	}
	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}
	overriding, err := variables.ParseArgs(buildArgs, pncvf, c.varCollection)
	if err != nil {
		return err
	}
	bcRawState, done := BuildContextFactory.Construct().RawState()
	bc.SetBuildContext(&bcRawState, c.mts.FinalTarget().String())
	state, dfImg, _, err := dockerfile2llb.Dockerfile2LLB(ctx, dfData, dockerfile2llb.ConvertOpt{
		MetaResolver: c.opt.MetaResolver,
		LLBCaps:      c.opt.LLBCaps,
		Config: dockerui.Config{
			BuildArgs:        overriding.Map(),
			Target:           dfTarget,
			ImageResolveMode: c.opt.ImageResolveMode,
		},
		TargetPlatform: &plat,
		Client:         bc,
	})
	done()
	if err != nil {
		return errors.Wrapf(err, "dockerfile2llb %s", dfPath)
	}
	// Convert dockerfile2llb image into earthfile2llb image via JSON.
	imgDt, err := json.Marshal(dfImg)
	if err != nil {
		return errors.Wrap(err, "marshal dockerfile image")
	}
	var img image.Image
	err = json.Unmarshal(imgDt, &img)
	if err != nil {
		return errors.Wrap(err, "unmarshal dockerfile image")
	}
	state2, img2, envVars := c.applyFromImage(pllb.FromRawState(*state), &img)
	c.mts.Final.MainState = state2
	c.mts.Final.MainImage = img2
	c.mts.Final.RanFromLike = true
	c.varCollection.ResetEnvVars(envVars)
	return nil
}

// Locally applies the earthly Locally command.
func (c *Converter) Locally(ctx context.Context) error {
	err := c.checkAllowed(locallyCmd)
	if err != nil {
		return err
	}
	if !c.opt.AllowLocally {
		return errors.New("LOCALLY cannot be used when --strict is specified or otherwise implied")
	}

	err = c.fromClassical(ctx, "scratch", platutil.NativePlatform, true)
	if err != nil {
		return err
	}

	workingDir, err := filepath.Abs(c.localWorkingDir)
	if err != nil {
		return errors.Wrapf(err, "unable to get abs path of %s", c.localWorkingDir)
	}

	c.varCollection.SetLocally(true)

	// reset WORKDIR to current directory where Earthfile is
	c.mts.Final.MainState = c.mts.Final.MainState.Dir(workingDir)
	c.mts.Final.MainImage.Config.WorkingDir = workingDir
	c.setPlatform(platutil.UserPlatform)
	return nil
}

// CopyArtifactLocal applies the earthly COPY artifact command which are invoked under a LOCALLY target.
func (c *Converter) CopyArtifactLocal(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string, isDir bool) error {
	err := c.checkAllowed(copyCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	prefix, cmdID, err := c.newVertexMeta(ctx, false, false, false, nil)
	if err != nil {
		return err
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, passArgs, buildArgs, false, copyCmd, cmdID)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", artifact.Target.String())
	}
	if artifact.Target.IsLocalInternal() {
		artifact.Target.LocalPath = c.mts.Final.Target.LocalPath
	}
	// Grab the artifacts state in the dep states, after we've built it.
	relevantDepState := mts.Final

	finalArgs := []string{localhost.SendFileMagicStr}
	if isDir {
		finalArgs = append(finalArgs, "--dir")
	}
	finalArgs = append(finalArgs, artifact.Artifact, dest)

	opts := []llb.RunOption{
		llb.Args(finalArgs),
		llb.IgnoreCache,
		pllb.AddMount("/"+localhost.SendFileMagicStr, relevantDepState.ArtifactsState, llb.Readonly),
		llb.WithCustomNamef(
			"%sCOPY %s%s%s %s",
			prefix,
			strIf(isDir, "--dir "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest),
	}
	c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()
	err = c.forceExecution(ctx, c.mts.Final.MainState, c.platr)
	if err != nil {
		return err
	}
	return nil
}

// CopyArtifact applies the earthly COPY artifact command.
func (c *Converter) CopyArtifact(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string, isDir bool, keepTs bool, keepOwn bool, chown string, chmod *fs.FileMode, ifExists, symlinkNoFollow bool) error {
	err := c.checkAllowed(copyCmd)
	if err != nil {
		return err
	}
	if chmod != nil && !c.ftrs.UseChmod {
		return fmt.Errorf("COPY --chmod is not supported in this version")
	}
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	prefix, cmdID, err := c.newVertexMeta(ctx, false, false, false, nil)
	if err != nil {
		return err
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, passArgs, buildArgs, false, copyCmd, cmdID)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", artifact.Target.String())
	}
	if artifact.Target.IsLocalInternal() {
		artifact.Target.LocalPath = c.mts.Final.Target.LocalPath
	}
	// Grab the artifacts state in the dep states, after we've built it.
	relevantDepState := mts.Final
	// Copy.
	c.mts.Final.MainState, err = llbutil.CopyOp(ctx,
		relevantDepState.ArtifactsState, []string{artifact.Artifact},
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), chmod, ifExists, symlinkNoFollow,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s%s%s %s",
			prefix,
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest))
	if err != nil {
		return errors.Wrapf(err, "copyOp CopyArtifact")
	}
	return nil
}

// CopyClassical applies the earthly COPY command, with classical args.
func (c *Converter) CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool, keepTs bool, keepOwn bool, chown string, chmod *fs.FileMode, ifExists bool) error {
	err := c.checkAllowed(copyCmd)
	if err != nil {
		return err
	}

	if chmod != nil && !c.ftrs.UseChmod {
		return fmt.Errorf("COPY --chmod is not supported in this version")
	}

	var srcState pllb.State
	if c.ftrs.UseCopyIncludePatterns {
		// create a new src state with the include patterns set (if this isn't done the entire context will be copied)
		srcStateFactory := addIncludePathAndSharedKeyHint(c.buildContextFactory, srcs)
		srcState = c.opt.LocalStateCache.getOrConstruct(srcStateFactory)
	} else {
		srcState = c.buildContextFactory.Construct()
	}

	c.nonSaveCommand()
	prefix, _, err := c.newVertexMeta(ctx, false, false, false, nil)
	if err != nil {
		return err
	}
	c.mts.Final.MainState, err = llbutil.CopyOp(ctx,
		srcState,
		srcs,
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), chmod, ifExists, false,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s %s",
			prefix,
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			strings.Join(srcs, " "),
			dest))
	if err != nil {
		return errors.Wrapf(err, "copyOp CopyClassical")
	}
	return nil
}

// ConvertRunOpts represents a set of options needed for the RUN command.
type ConvertRunOpts struct {
	CommandName          string
	Args                 []string
	Locally              bool
	Mounts               []string
	Secrets              []string
	WithEntrypoint       bool
	WithShell            bool
	Privileged           bool
	NoNetwork            bool
	Push                 bool
	Transient            bool
	WithSSH              bool
	NoCache              bool
	Interactive          bool
	InteractiveKeep      bool
	InteractiveSaveFiles []debuggercommon.SaveFilesSettings

	// Internal.
	shellWrap    shellWrapFun
	extraRunOpts []llb.RunOption
	statePrep    func(context.Context, pllb.State) (pllb.State, error)
}

// Run applies the earthly RUN command.
func (c *Converter) Run(ctx context.Context, opts ConvertRunOpts) error {
	err := c.checkAllowed(runCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()

	for _, state := range c.persistentCacheDirs {
		opts.extraRunOpts = append(opts.extraRunOpts, state.RunOption)
	}
	_, err = c.internalRun(ctx, opts)
	return err
}

// RunExitCode executes a run for the purpose of determining the exit code of the command. This can be used in conditionals.
func (c *Converter) RunExitCode(ctx context.Context, opts ConvertRunOpts) (int, error) {
	err := c.checkAllowed(runCmd)
	if err != nil {
		return 0, err
	}
	c.nonSaveCommand()
	for _, state := range c.persistentCacheDirs {
		opts.extraRunOpts = append(opts.extraRunOpts, state.RunOption)
	}

	var exitCodeFile string
	if opts.Locally {
		exitCodeDir, err := os.MkdirTemp(os.TempDir(), "earthlyexitcode")
		if err != nil {
			return 0, errors.Wrap(err, "create temp dir")
		}
		exitCodeFile = filepath.Join(exitCodeDir, "/exit_code")
		c.opt.CleanCollection.Add(func() error {
			return os.RemoveAll(exitCodeDir)
		})
	} else {
		exitCodeFile = "/tmp/earthly_if_statement_exit_code"
		prefix, _, err := c.newVertexMeta(ctx, false, false, true, nil)
		if err != nil {
			return 0, err
		}
		opts.statePrep = func(ctx context.Context, state pllb.State) (pllb.State, error) {
			return state.File(
				pllb.Mkdir("/run", 0755, llb.WithParents(true)),
				llb.WithCustomNamef(
					"%smkdir %s",
					prefix, "/run"),
			), nil
		}
	}

	// Perform execution, but append the command with the right shell incantation that
	// causes it to output the exit code to a file. This is done via the shellWrap.
	opts.shellWrap = withShellAndEnvVarsExitCode(exitCodeFile)
	opts.WithShell = true // force shell wrapping
	state, err := c.internalRun(ctx, opts)
	if err != nil {
		return 0, err
	}
	var codeDt []byte
	if opts.Locally {
		codeDt, err = os.ReadFile(exitCodeFile)
		if err != nil {
			return 0, errors.Wrap(err, "read exit code file")
		}
	} else {
		ref, err := llbutil.StateToRef(
			ctx, c.opt.GwClient, state, c.opt.NoCache,
			c.platr, c.opt.CacheImports.AsSlice())
		if err != nil {
			return 0, errors.Wrap(err, "run exit code state to ref")
		}
		codeDt, err = ref.ReadFile(ctx, gwclient.ReadRequest{
			Filename: exitCodeFile,
		})
		if err != nil {
			return 0, errors.Wrap(err, "read exit code")
		}
	}
	exitCode, err := strconv.ParseInt(string(bytes.TrimSpace(codeDt)), 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse exit code as int")
	}
	return int(exitCode), err
}

// RunExpression runs an expression and returns its output. The run is transient - any state created
// is not used in subsequent commands.
func (c *Converter) RunExpression(ctx context.Context, expressionName string, opts ConvertRunOpts) (string, error) {
	for _, state := range c.persistentCacheDirs {
		opts.extraRunOpts = append(opts.extraRunOpts, state.RunOption)
	}
	return c.runCommand(ctx, expressionName, true, opts)
}

// RunCommand runs a command and returns its output. The run is transient - any state created
// is not used in subsequent commands.
func (c *Converter) RunCommand(ctx context.Context, commandName string, opts ConvertRunOpts) (string, error) {
	for _, state := range c.persistentCacheDirs {
		opts.extraRunOpts = append(opts.extraRunOpts, state.RunOption)
	}
	return c.runCommand(ctx, commandName, false, opts)
}

func (c *Converter) runCommand(ctx context.Context, outputFileName string, isExpression bool, opts ConvertRunOpts) (string, error) {
	err := c.checkAllowed(runCmd)
	if err != nil {
		return "", err
	}
	c.nonSaveCommand()

	if !opts.WithShell {
		panic("runCommand must be called WithShell")
	}
	if opts.Locally == opts.Transient {
		panic("runCommand Transient xor Locally must be true")
	}
	if opts.shellWrap != nil {
		panic("runCommand expects shellWrap to be nil (as it is overridden)")
	}

	var outputFile string
	if opts.Locally {
		outputDir, err := os.MkdirTemp(os.TempDir(), "earthlyexproutput")
		if err != nil {
			return "", errors.Wrap(err, "create temp dir")
		}
		outputFile = filepath.Join(outputDir, "/output")
		c.opt.CleanCollection.Add(func() error {
			return os.RemoveAll(outputDir)
		})
	} else {
		srcBuildArgDir := "/run/buildargs"
		prefix, _, err := c.newVertexMeta(ctx, false, false, true, nil)
		if err != nil {
			return "", err
		}
		outputFile = path.Join(srcBuildArgDir, outputFileName)
		opts.statePrep = func(ctx context.Context, state pllb.State) (pllb.State, error) {
			return state.File(
				pllb.Mkdir(srcBuildArgDir, 0777, llb.WithParents(true)), // Mkdir is performed as root even when USER is set; we must use 0777
				llb.WithCustomNamef(
					"%smkdir %s",
					prefix, srcBuildArgDir),
			), nil
		}
	}

	if isExpression {
		opts.shellWrap = expressionWithShellAndEnvVarsOutput(outputFile)
	} else {
		opts.shellWrap = withShellAndEnvVarsOutput(outputFile)
	}

	state, err := c.internalRun(ctx, opts)
	if err != nil {
		return "", err
	}

	var outputDt []byte
	if opts.Locally {
		outputDt, err = os.ReadFile(outputFile)
		if err != nil {
			return "", errors.Wrap(err, "read output file")
		}
	} else {
		ref, err := llbutil.StateToRef(
			ctx, c.opt.GwClient, state, c.opt.NoCache,
			c.platr, c.opt.CacheImports.AsSlice())
		if err != nil {
			return "", errors.Wrapf(err, "build arg state to ref")
		}
		outputDt, err = ref.ReadFile(ctx, gwclient.ReadRequest{Filename: outputFile})
		if err != nil {
			return "", errors.Wrapf(err, "non constant build arg read request")
		}
	}
	// echo adds a trailing \n.
	outputDt = bytes.TrimSuffix(outputDt, []byte("\n"))
	return string(outputDt), nil
}

// SaveArtifact applies the earthly SAVE ARTIFACT command.
func (c *Converter) SaveArtifact(ctx context.Context, saveFrom, saveTo, saveAsLocalTo string, keepTs, keepOwn, ifExists, symlinkNoFollow, force, isPush bool) (retErr error) {
	err := c.checkAllowed(saveArtifactCmd)
	if err != nil {
		return err
	}
	absSaveFrom, err := llbutil.Abs(ctx, c.mts.Final.MainState, saveFrom)
	if err != nil {
		return err
	}
	if absSaveFrom == "/" || absSaveFrom == "" {
		return errors.New("cannot save root dir as artifact")
	}
	saveToAdjusted := saveTo
	if saveTo == "" || saveTo == "." || strings.HasSuffix(saveTo, "/") {
		saveFromRelative := path.Join(".", absSaveFrom)
		saveToAdjusted = path.Join(saveTo, path.Base(saveFromRelative))
	}
	saveToD, saveToF := splitWildcards(saveToAdjusted)
	var artifactPath string
	if saveToF == "" {
		artifactPath = saveToAdjusted
	} else {
		saveToAdjusted = fmt.Sprintf("%s/", saveToD)
		artifactPath = path.Join(saveToAdjusted, saveToF)
	}
	artifact := domain.Artifact{
		Target:   c.mts.Final.Target,
		Artifact: artifactPath,
	}
	own := "root:root"
	if keepOwn {
		own = ""
	}

	// pcState is a separate state from the main state
	// which persists any files cached via CACHE command.
	// This is necessary so those cached files can be
	// accessed within the CopyOps below.
	pcState := c.persistCache(c.mts.Final.MainState)

	prefix, cmdID, err := c.newVertexMeta(ctx, false, false, false, nil)
	if err != nil {
		return err
	}

	cmd, ok := c.opt.Logbus.Run().Command(cmdID)
	if !ok {
		return errors.New("command not found")
	}

	cmd.SetName(fmt.Sprintf("SAVE ARTIFACT %s", saveFrom))

	defer func() {
		cmd.SetEndError(retErr)
	}()

	c.mts.Final.ArtifactsState, err = llbutil.CopyOp(ctx,
		pcState, []string{saveFrom}, c.mts.Final.ArtifactsState,
		saveToAdjusted, true, true, keepTs, own, nil, ifExists, symlinkNoFollow,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sSAVE ARTIFACT %s%s%s %s",
			prefix,
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			saveFrom,
			artifact.String()))
	if err != nil {
		return errors.Wrapf(err, "copyOp save artifact")
	}
	if saveAsLocalTo != "" {
		separateArtifactsState := c.platr.Scratch()
		if isPush {
			pushState := c.persistCache(c.mts.Final.RunPush.State)
			prefix, _, err := c.newVertexMeta(ctx, false, false, false, nil)
			if err != nil {
				return err
			}
			separateArtifactsState, err = llbutil.CopyOp(ctx,
				pushState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", nil, ifExists, symlinkNoFollow,
				c.ftrs.UseCopyLink,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					prefix,
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
			if err != nil {
				return errors.Wrapf(err, "copyOp save artifact as local")
			}
		} else {
			prefix, _, err := c.newVertexMeta(ctx, false, false, false, nil)
			if err != nil {
				return err
			}
			separateArtifactsState, err = llbutil.CopyOp(ctx,
				pcState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", nil, ifExists, symlinkNoFollow,
				c.ftrs.UseCopyLink,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					prefix,
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
			if err != nil {
				return errors.Wrapf(err, "copyOp save artifact as local")
			}
		}
		c.mts.Final.SeparateArtifactsState = append(c.mts.Final.SeparateArtifactsState, separateArtifactsState)

		saveAsLocalToAdj := saveAsLocalTo
		if saveAsLocalToAdj == "." {
			saveAsLocalToAdj = "./"
		}

		if !force {
			canSave, err := c.canSave(ctx, saveAsLocalToAdj)
			if err != nil {
				return err
			}
			if !canSave {
				if c.ftrs.RequireForceForUnsafeSaves {
					return fmt.Errorf("unable to save to %s; path must be located under %s", saveAsLocalTo, c.target.LocalPath)
				}
				analytics.Count("breaking-change", "save-artifact-force-flag-required-warning")
				c.opt.Console.Warnf("saving to path (%s) outside of current directory (%s) will require a --force flag in a future version", saveAsLocalTo, c.target.LocalPath)
			}
		}

		saveLocal := states.SaveLocal{
			DestPath:     saveAsLocalToAdj,
			ArtifactPath: artifactPath,
			Index:        len(c.mts.Final.SeparateArtifactsState) - 1,
			IfExists:     ifExists,
		}

		if c.ftrs.WaitBlock {
			waitItem := newSaveArtifactLocal(saveLocal, c, c.opt.DoSaves)
			c.waitBlock().AddItem(waitItem)
			c.mts.Final.WaitItems = append(c.mts.Final.WaitItems, waitItem)
		} else {
			if isPush {
				c.mts.Final.RunPush.SaveLocals = append(c.mts.Final.RunPush.SaveLocals, saveLocal)
			} else {
				c.mts.Final.SaveLocals = append(c.mts.Final.SaveLocals, saveLocal)
			}
		}

	}
	c.ranSave = true
	c.markFakeDeps()
	return nil
}

func (c *Converter) canSave(ctx context.Context, saveAsLocalTo string) (bool, error) {
	basepath, err := filepath.Abs(c.target.LocalPath)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get absolute path of %s", basepath)
	}
	basePathExists, err := fileutil.DirExists(basepath)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check if %s exists", basepath)
	}
	if !basePathExists {
		return false, fmt.Errorf("no such directory: %s", basepath)
	}
	basepath += string(filepath.Separator)

	hasTrailingSlash := strings.HasSuffix(saveAsLocalTo, "/") && saveAsLocalTo != "/"
	saveAsLocalToAdj := saveAsLocalTo
	if !strings.HasPrefix(saveAsLocalTo, "/") {
		saveAsLocalToAdj = path.Join(c.target.LocalPath, saveAsLocalTo)
	}
	saveAsLocalToAdj, err = filepath.Abs(saveAsLocalToAdj)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get absolute path of %q", saveAsLocalTo)
	}
	if hasTrailingSlash {
		saveAsLocalToAdj += "/"
	}
	return strings.HasPrefix(saveAsLocalToAdj, basepath), nil
}

// SaveArtifactFromLocal saves a local file into the ArtifactsState
func (c *Converter) SaveArtifactFromLocal(ctx context.Context, saveFrom, saveTo string, keepTs, keepOwn bool, chown string) error {
	err := c.checkAllowed(saveArtifactCmd)
	if err != nil {
		return err
	}
	src, err := filepath.Abs(saveFrom)
	if err != nil {
		return err
	}

	if saveTo == "" || saveTo == "." || strings.HasSuffix(saveTo, "/") {
		saveTo = path.Join(saveTo, path.Base(src))
	}

	// first load the files into a snapshot
	prefix, _, err := c.newVertexMeta(ctx, true, false, true, nil)
	if err != nil {
		return err
	}
	opts := []llb.RunOption{
		llb.Args([]string{localhost.CopyFileMagicStr, saveFrom, saveTo}),
		llb.IgnoreCache,
		llb.WithCustomNamef(
			"%sCopyFileMagicStr %s %s",
			prefix, saveFrom, saveTo),
	}
	c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()

	// then save it via the regular SaveArtifact code since it's now in a snapshot
	absSaveTo := fmt.Sprintf("/%s", saveTo)
	own := "root:root"
	if keepOwn {
		own = ""
	} else if chown != "" {
		own = chown
	}
	ifExists := false
	c.mts.Final.ArtifactsState, err = llbutil.CopyOp(ctx,
		c.mts.Final.MainState, []string{absSaveTo}, c.mts.Final.ArtifactsState,
		absSaveTo, true, true, keepTs, own, nil, ifExists, false,
		c.ftrs.UseCopyLink,
	)
	if err != nil {
		return errors.Wrapf(err, "copyOp save artifact from local")
	}
	err = c.forceExecution(ctx, c.mts.Final.ArtifactsState, c.platr)
	if err != nil {
		return err
	}
	c.ranSave = true
	c.markFakeDeps()
	return nil
}

func (c *Converter) waitBlock() *waitBlock {
	n := len(c.waitBlockStack)
	if n == 0 {
		panic("waitBlock() called on empty stack") // shouldn't happen
	}
	return c.waitBlockStack[n-1]
}

// PushWaitBlock should be called when a WAIT block starts, all commands will be added to this new block
func (c *Converter) PushWaitBlock(ctx context.Context) error {
	waitBlock := newWaitBlock()
	c.waitBlockStack = append(c.waitBlockStack, waitBlock)
	c.mts.Final.AddWaitBlock(waitBlock)
	return nil
}

// PopWaitBlock should be called when an END is encountered, it will block until all commands within the block complete
func (c *Converter) PopWaitBlock(ctx context.Context) error {
	n := len(c.waitBlockStack)
	if n == 0 {
		return fmt.Errorf("waitBlockStack is empty") // shouldn't happen
	}

	if c.ftrs.WaitBlock {
		// an END is only ever encountered by the converter that created the WAIT block, this is the only special
		// instance where we reference mts.Final.MainState before calling FinalizeStates; this can be done here
		// as the waitBlock belongs to the current Converter
		c.waitBlock().AddItem(newStateWaitItem(&c.mts.Final.MainState, c))
	}

	i := n - 1
	waitBlock := c.waitBlockStack[i]
	c.waitBlockStack = c.waitBlockStack[:i]

	return waitBlock.Wait(ctx, c.opt.DoPushes, c.opt.DoSaves)
}

// SaveImage applies the earthly SAVE IMAGE command.
func (c *Converter) SaveImage(ctx context.Context, imageNames []string, hasPushFlag bool, insecurePush bool, cacheHint bool, cacheFrom []string, noManifestList bool) (retErr error) {
	err := c.checkAllowed(saveImageCmd)
	if err != nil {
		return err
	}
	if noManifestList && !c.ftrs.UseNoManifestList {
		return fmt.Errorf("SAVE IMAGE --no-manifest-list is not supported in this version")
	}
	_, cmd, err := c.newLogbusCommand(ctx, fmt.Sprintf("SAVE IMAGE %s", strings.Join(imageNames, " ")))
	if err != nil {
		return errors.Wrap(err, "failed to create command")
	}
	defer func() {
		cmd.SetEndError(retErr)
	}()
	for _, cf := range cacheFrom {
		c.opt.CacheImports.Add(cf)
	}
	justCacheHint := false
	if len(imageNames) == 0 && cacheHint {
		imageNames = []string{""}
		justCacheHint = true
	}
	for _, imageName := range imageNames {
		if c.mts.Final.RunPush.HasState {
			if c.ftrs.WaitBlock {
				panic("RunPush.HasState should never be true when --wait-block is used")
			}
			// pcState persists any files that may be cached via CACHE command.
			pcState := c.persistCache(c.mts.Final.RunPush.State)
			// SAVE IMAGE --push when it comes before any RUN --push should be treated as if they are in the main state,
			// since thats their only dependency. It will still be marked as a push.
			c.mts.Final.RunPush.SaveImages = append(c.mts.Final.RunPush.SaveImages,
				states.SaveImage{
					State:               pcState,
					Image:               c.mts.Final.MainImage.Clone(), // We can get away with this because no Image details can vary in a --push. This should be fixed before then.
					DockerTag:           imageName,
					Push:                hasPushFlag,
					InsecurePush:        insecurePush,
					CacheHint:           cacheHint,
					HasPushDependencies: true,
					ForceSave:           c.opt.ForceSaveImage,
					CheckDuplicate:      c.ftrs.CheckDuplicateImages,
					NoManifestList:      noManifestList,
				})
		} else {
			si := states.SaveImage{
				State:               c.persistCache(c.mts.Final.MainState),
				Image:               c.mts.Final.MainImage.Clone(),
				DockerTag:           imageName,
				Push:                hasPushFlag,
				InsecurePush:        insecurePush,
				CacheHint:           cacheHint,
				HasPushDependencies: false,
				ForceSave:           c.opt.ForceSaveImage,
				CheckDuplicate:      c.ftrs.CheckDuplicateImages,
				NoManifestList:      noManifestList,

				Platform:    c.platr.Materialize(c.platr.Current()),
				HasPlatform: platutil.IsPlatformDefined(c.platr.Current()),
			}

			if c.ftrs.WaitBlock {
				shouldPush := hasPushFlag && si.DockerTag != ""
				shouldExportLocally := si.DockerTag != "" && c.opt.DoSaves
				waitItem := newSaveImage(si, c, shouldPush, shouldExportLocally)
				c.waitBlock().AddItem(waitItem)
				c.mts.Final.WaitItems = append(c.mts.Final.WaitItems, waitItem)
				if hasPushFlag {
					// only add summary for `SAVE IMAGE --push` commands
					c.opt.ExportCoordinator.AddPushedImageSummary(c.target.StringCanonical(), si.DockerTag, c.mts.Final.ID, c.opt.DoPushes)
				}

				// TODO this is here as a work-around for https://github.com/earthly/earthly/issues/2178
				// ideally we should always set SkipBuilder = true even when we are under the first implicit wait block
				// however we don't want to break inline caching for users who are using VERSION 0.7 without any explicit WAIT blocks
				if !c.opt.UseInlineCache || len(c.waitBlockStack) > 1 {
					si.SkipBuilder = true
				}
			}
			c.mts.Final.SaveImages = append(c.mts.Final.SaveImages, si)
		}

		if hasPushFlag && imageName != "" && c.opt.UseInlineCache {
			// Use this image tag as cache import too.
			c.opt.CacheImports.Add(imageName)
		}
	}
	if !justCacheHint {
		c.ranSave = true
		c.markFakeDeps()
	}
	return nil
}

// Build applies the earthly BUILD command.
func (c *Converter) Build(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string) error {
	err := c.checkAllowed(buildCmd)
	if err != nil {
		return err
	}

	c.nonSaveCommand()

	cmdID, cmd, err := c.newLogbusCommand(ctx, fmt.Sprintf("BUILD %s", fullTargetName))
	if err != nil {
		return errors.Wrap(err, "failed to create command")
	}

	_, err = c.buildTarget(ctx, fullTargetName, platform, allowPrivileged, passArgs, buildArgs, true, buildCmd, cmdID)

	cmd.SetEndError(err)

	return err
}

type afterParallelFunc func(context.Context, *states.MultiTarget) error

// BuildAsync applies the earthly BUILD command asynchronously.
func (c *Converter) BuildAsync(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string, cmdT cmdType, apf afterParallelFunc, sem semutil.Semaphore) error {
	target, opt, _, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, passArgs, buildArgs, true, cmdT, "")
	if err != nil {
		return err
	}
	c.opt.ErrorGroup.Go(func() error {
		if sem == nil {
			sem = c.opt.Parallelism
		}
		rel, err := sem.Acquire(ctx, 1)
		if err != nil {
			return errors.Wrapf(err, "acquiring parallelism semaphore for %s", fullTargetName)
		}
		defer rel()
		mts, err := Earthfile2LLB(ctx, target, opt, false)
		if err != nil {
			return errors.Wrapf(err, "async earthfile2llb for %s", fullTargetName)
		}
		if apf != nil {
			if c.ftrs.ExecAfterParallel && mts != nil && mts.Final != nil {
				// TODO: This is a duplication from the forceExecution taking place
				//       from FinalizeStates. However, this is necessary for apf
				//       synchronization (needs to be run after target has executed).
				err := c.forceExecution(ctx, mts.Final.MainState, mts.Final.PlatformResolver)
				if err != nil {
					return errors.Wrapf(err, "async force execution for %s", fullTargetName)
				}
			}
			err := apf(ctx, mts)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

// Workdir applies the WORKDIR command.
func (c *Converter) Workdir(ctx context.Context, workdirPath string) error {
	err := c.checkAllowed(workdirCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.mts.Final.MainState = c.mts.Final.MainState.Dir(workdirPath)
	workdirAbs := workdirPath
	if !path.IsAbs(workdirAbs) {
		workdirAbs = path.Join("/", c.mts.Final.MainImage.Config.WorkingDir, workdirAbs)
	}
	c.mts.Final.MainImage.Config.WorkingDir = workdirAbs
	if workdirAbs != "/" {
		// Mkdir.
		mkdirOpts := []llb.MkdirOption{
			llb.WithParents(true),
		}
		if c.mts.Final.MainImage.Config.User != "" {
			mkdirOpts = append(mkdirOpts, llb.WithUser(c.mts.Final.MainImage.Config.User))
		}
		prefix, _, err := c.newVertexMeta(ctx, false, false, false, nil)
		if err != nil {
			return err
		}
		opts := []llb.ConstraintsOpt{
			llb.WithCustomNamef("%sWORKDIR %s", prefix, workdirPath),
		}
		c.mts.Final.MainState = c.mts.Final.MainState.File(
			pllb.Mkdir(workdirAbs, 0755, mkdirOpts...), opts...)
	}
	return nil
}

// User applies the USER command.
func (c *Converter) User(ctx context.Context, user string) error {
	err := c.checkAllowed(userCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.mts.Final.MainState = c.mts.Final.MainState.User(user)
	c.mts.Final.MainImage.Config.User = user
	return nil
}

// Cmd applies the CMD command.
func (c *Converter) Cmd(ctx context.Context, cmdArgs []string, isWithShell bool) error {
	err := c.checkAllowed(cmdCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.mts.Final.MainImage.Config.Cmd = withShell(cmdArgs, isWithShell)
	c.cmdSet = true
	return nil
}

// Entrypoint applies the ENTRYPOINT command.
func (c *Converter) Entrypoint(ctx context.Context, entrypointArgs []string, isWithShell bool) error {
	err := c.checkAllowed(entrypointCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.mts.Final.MainImage.Config.Entrypoint = withShell(entrypointArgs, isWithShell)
	if !c.cmdSet {
		c.mts.Final.MainImage.Config.Cmd = nil
	}
	return nil
}

// Expose applies the EXPOSE command.
func (c *Converter) Expose(ctx context.Context, ports []string) error {
	err := c.checkAllowed(exposeCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	for _, port := range ports {
		c.mts.Final.MainImage.Config.ExposedPorts[port] = struct{}{}
	}
	return nil
}

// Volume applies the VOLUME command.
func (c *Converter) Volume(ctx context.Context, volumes []string) error {
	err := c.checkAllowed(volumeCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	for _, volume := range volumes {
		c.mts.Final.MainImage.Config.Volumes[volume] = struct{}{}
	}
	return nil
}

// Env applies the ENV command.
func (c *Converter) Env(ctx context.Context, envKey string, envValue string) error {
	err := c.checkAllowed(envCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.varCollection.DeclareEnv(envKey, envValue)
	c.mts.Final.MainState = c.mts.Final.MainState.AddEnv(envKey, envValue)
	c.mts.Final.MainImage.Config.Env = variables.AddEnv(
		c.mts.Final.MainImage.Config.Env, envKey, envValue)
	return nil
}

// Arg applies the ARG command.
func (c *Converter) Arg(ctx context.Context, argKey string, defaultArgValue string, opts commandflag.ArgOpts) error {
	err := c.checkAllowed(argCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()

	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}

	declOpts := []variables.DeclareOpt{
		variables.AsArg(),
		variables.WithValue(defaultArgValue),
		variables.WithPNCVFunc(pncvf),
	}
	if opts.Global {
		declOpts = append(declOpts, variables.AsGlobal())
	}
	effective, effectiveDefault, err := c.varCollection.DeclareVar(argKey, declOpts...)
	if err != nil {
		return err
	}
	if opts.Required && len(effective) == 0 {
		return fmt.Errorf("value not supplied for required ARG: %s", argKey)
	}
	if len(defaultArgValue) > 0 && reserved.IsBuiltIn(argKey) {
		return fmt.Errorf("arg default value supplied for built-in ARG: %s", argKey)
	}
	if c.varCollection.IsStackAtBase() { // Only when outside of UDC.
		c.mts.Final.AddBuildArgInput(dedup.BuildArgInput{
			Name:          argKey,
			DefaultValue:  effectiveDefault,
			ConstantValue: effective,
		})
	}
	return nil
}

// Let applies the LET command.
func (c *Converter) Let(ctx context.Context, key string, value string) error {
	err := c.checkAllowed(letCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()

	if reserved.IsBuiltIn(key) {
		return fmt.Errorf("LET cannot override built-in variable %q", key)
	}

	effective, effectiveDefault, err := c.varCollection.DeclareVar(key, variables.WithValue(value))
	if err != nil {
		return err
	}
	if c.varCollection.IsStackAtBase() {
		c.mts.Final.AddBuildArgInput(dedup.BuildArgInput{
			Name:          key,
			DefaultValue:  effectiveDefault,
			ConstantValue: effective,
		})
	}
	return nil
}

// UpdateArg updates an existing arg to a new value. It errors if the arg could
// not be found.
func (c *Converter) UpdateArg(ctx context.Context, argKey string, argValue string, isBase bool) error {
	if err := c.checkAllowed(setCmd); err != nil {
		return err
	}
	c.nonSaveCommand()

	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}

	err := c.varCollection.UpdateVar(argKey, argValue, pncvf, isBase)
	if err != nil {
		return err
	}
	return nil
}

// SetArg sets an arg to a specific value.
func (c *Converter) SetArg(ctx context.Context, argKey string, argValue string) error {
	err := c.checkAllowed(argCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.varCollection.SetArg(argKey, argValue)
	return nil
}

// UnsetArg unsets a previously declared arg. If the arg does not exist this operation is a no-op.
func (c *Converter) UnsetArg(ctx context.Context, argKey string) error {
	err := c.checkAllowed(argCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.varCollection.UnsetArg(argKey)
	return nil
}

// Label applies the LABEL command.
func (c *Converter) Label(ctx context.Context, labels map[string]string) error {
	err := c.checkAllowed(labelCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	for key, value := range labels {
		c.mts.Final.MainImage.Config.Labels[key] = value
	}
	return nil
}

// GitClone applies the GIT CLONE command.
func (c *Converter) GitClone(ctx context.Context, gitURL string, sshCommand string, branch string, dest string, keepTs bool) error {
	err := c.checkAllowed(gitCloneCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	gitURLScrubbed := stringutil.ScrubCredentials(gitURL)
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef(
			"%sGIT CLONE (--branch %s) %s", c.vertexMetaWithURL(gitURLScrubbed), branch, gitURLScrubbed),
		llb.KeepGitDir(),
	}
	if sshCommand != "" {
		gitOpts = append(gitOpts, llb.SSHCommand(sshCommand))
	}
	gitState := pllb.Git(gitURL, branch, gitOpts...)
	prefix, _, err := c.newVertexMeta(ctx, false, false, false, nil)
	if err != nil {
		return err
	}
	c.mts.Final.MainState, err = llbutil.CopyOp(ctx,
		gitState, []string{"."}, c.mts.Final.MainState, dest, false, false, keepTs,
		c.mts.Final.MainImage.Config.User, nil, false, false, c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY GIT CLONE (--branch %s) %s TO %s", prefix,
			branch, gitURLScrubbed, dest))
	if err != nil {
		return errors.Wrapf(err, "copyOp git clone")
	}
	return nil
}

// WithDockerRun applies an entire WITH DOCKER ... RUN ... END clause.
func (c *Converter) WithDockerRun(ctx context.Context, args []string, opt WithDockerOpt, allowParallel bool) error {
	err := c.checkAllowed(runCmd)
	if err != nil {
		return err
	}

	c.nonSaveCommand()

	enableParallel := allowParallel && c.opt.ParallelConversion && c.ftrs.ParallelLoad

	for _, state := range c.persistentCacheDirs {
		opt.extraRunOpts = append(opt.extraRunOpts, state.RunOption)
	}

	if c.ftrs.UseRegistryForWithDocker {
		wdr := newWithDockerRunRegistry(c, enableParallel)
		return wdr.Run(ctx, args, opt)
	}

	wdr := newWithDockerRunTar(c, enableParallel)
	return wdr.Run(ctx, args, opt)
}

// WithDockerRunLocal applies an entire WITH DOCKER ... RUN ... END clause.
func (c *Converter) WithDockerRunLocal(ctx context.Context, args []string, opt WithDockerOpt, allowParallel bool) error {
	err := c.checkAllowed(runCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	enableParallel := allowParallel && c.opt.ParallelConversion && c.ftrs.ParallelLoad

	if c.ftrs.UseRegistryForWithDocker && c.opt.UseLocalRegistry {
		wdr := newWithDockerRunLocalReg(c, enableParallel)
		return wdr.Run(ctx, args, opt)
	}

	wdrl := newWithDockerRunLocal(c, enableParallel)
	return wdrl.Run(ctx, args, opt)
}

// Healthcheck applies the HEALTHCHECK command.
func (c *Converter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int, startInterval time.Duration) error {
	err := c.checkAllowed(healthcheckCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	hc := &dockerimage.HealthConfig{}
	if isNone {
		hc.Test = []string{"NONE"}
	} else {
		// TODO: Should support also CMD without shell (exec form).
		//       See https://github.com/moby/buildkit/blob/master/frontend/dockerfile/dockerfile2llb/image.go#L18
		hc.Test = []string{"CMD-SHELL", strings.Join(cmdArgs, " ")}
		hc.Interval = interval
		hc.Timeout = timeout
		hc.StartPeriod = startPeriod
		hc.Retries = retries
		hc.StartInterval = startInterval
	}
	c.mts.Final.MainImage.Config.Healthcheck = hc
	return nil
}

// Import applies the IMPORT command.
func (c *Converter) Import(ctx context.Context, importStr, as string, isGlobal, currentlyPrivileged, allowPrivilegedFlag bool) error {
	err := c.checkAllowed(importCmd)
	if err != nil {
		return err
	}
	return c.varCollection.Imports().Add(importStr, as, isGlobal, currentlyPrivileged, allowPrivilegedFlag)
}

// Cache handles a `CACHE` command in a Target.
// It appends run options to the Converter which will mount a cache volume in each successive `RUN` command,
// and configures the `Converter` to persist the cache in the image at the end of the target.
func (c *Converter) Cache(ctx context.Context, mountTarget string, opts commandflag.CacheOpts) error {
	err := c.checkAllowed(cacheCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	key := cacheKey(c.target)
	cacheID := path.Join("/run/cache", key, path.Clean(mountTarget))
	if c.ftrs.GlobalCache && opts.ID != "" {
		cacheID = opts.ID
	}

	var shareMode llb.CacheMountSharingMode
	switch opts.Sharing {
	case "shared":
		shareMode = llb.CacheMountShared
	case "private":
		shareMode = llb.CacheMountPrivate
	case "locked", "":
		shareMode = llb.CacheMountLocked
	default:
		return errors.Errorf("invalid cache sharing mode %q", opts.Sharing)
	}

	if _, exists := c.persistentCacheDirs[mountTarget]; !exists {
		var mountOpts []llb.MountOption
		mountOpts = append(mountOpts, llb.AsPersistentCacheDir(cacheID, shareMode))
		mountOpts = append(mountOpts, llb.SourcePath("/cache"))
		var mountMode int
		if opts.Mode == "" {
			opts.Mode = "0644"
		}
		mountMode, err = ParseMode(opts.Mode)
		if err != nil {
			return errors.Errorf("failed to parse mount mode %s", opts.Mode)
		}
		persisted := true // Without new --cache-persist-option we use old behaviour which is persisted
		if c.ftrs.CachePersistOption {
			persisted = opts.Persist
		} else if opts.Persist {
			return errors.Errorf("the --persist flag is only available when VERSION --cache-persist-option is enabled")
		}
		c.persistentCacheDirs[mountTarget] = states.CacheMount{
			Persisted: persisted,
			RunOption: pllb.AddMount(mountTarget, pllb.Scratch().File(pllb.Mkdir("/cache", os.FileMode(mountMode))), mountOpts...),
		}
	}
	return nil
}

// Host handles a `HOST` command in a Target.
func (c *Converter) Host(ctx context.Context, hostname string, ip net.IP) error {
	err := c.checkAllowed(hostCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.mts.Final.MainState = c.mts.Final.MainState.AddExtraHost(hostname, ip)
	return nil
}

// Project handles a "PROJECT" command in base target.
func (c *Converter) Project(ctx context.Context, org, project string) error {
	err := c.checkAllowed(projectCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.varCollection.SetOrg(org)
	c.varCollection.SetProject(project)
	c.opt.ProjectAdder.AddProject(org, project)
	return nil
}

// Pipeline handles a "PIPELINE" command.
func (c *Converter) Pipeline(ctx context.Context) error {
	err := c.checkAllowed(pipelineCmd)
	if err != nil {
		return err
	}
	c.isPipeline = true
	c.mts.Final.RanFromLike = true
	return nil
}

func (c *Converter) ExpandWildcard(ctx context.Context, fullTargetName string, cmd spec.Command) ([]spec.Command, error) {

	parsedTarget, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		return nil, err
	}

	if strings.Contains(fullTargetName, "**") {
		return nil, errors.New("globstar (**) pattern not yet supported")
	}

	var target domain.Target
	if c.target.IsRemote() {
		target = c.target
	} else {
		target = parsedTarget
	}

	matches, err := c.opt.Resolver.ExpandWildcard(ctx, c.opt.GwClient, c.platr, target, parsedTarget.GetLocalPath())
	if err != nil {
		return nil, err
	}

	children := []spec.Command{}
	for _, match := range matches {
		childTargetName := fmt.Sprintf("./%s+%s", match, parsedTarget.GetName())

		childTarget, err := domain.ParseTarget(childTargetName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse target %q", childTargetName)
		}

		data, _, _, err := c.ResolveReference(ctx, childTarget)
		if err != nil {
			notExist := buildcontext.ErrEarthfileNotExist{}
			if errors.As(err, &notExist) {
				continue
			}
			return nil, errors.Wrapf(err, "unable to resolve target %q", childTargetName)
		}

		var found bool
		for _, target := range data.Earthfile.Targets {
			if target.Name == childTarget.GetName() {
				found = true
				break
			}
		}

		if !found {
			continue
		}

		cloned := cmd.Clone()
		cloned.Args[0] = childTargetName
		children = append(children, cloned)
	}

	if len(children) == 0 {
		return nil, errors.Wrapf(err, "no matching targets found for pattern %q", parsedTarget.GetLocalPath())
	}

	return children, nil
}

// ResolveReference resolves a reference's build context given the current state: relativity to the Earthfile, imports etc.
func (c *Converter) ResolveReference(ctx context.Context, ref domain.Reference) (bc *buildcontext.Data, allowPrivileged, allowPrivilegedSet bool, err error) {
	derefed, allowPrivileged, allowPrivilegedSet, err := c.varCollection.Imports().Deref(ref)
	if err != nil {
		return nil, false, false, err
	}
	refToResolve, err := c.joinRefs(derefed)
	if err != nil {
		return nil, false, false, err
	}
	bc, err = c.opt.Resolver.Resolve(ctx, c.opt.GwClient, c.platr, refToResolve)
	if err != nil {
		return nil, false, false, err
	}
	return bc, allowPrivileged, allowPrivilegedSet, nil
}

// EnterScopeDo introduces a new variable scope. Globals and imports are fetched from baseTarget.
func (c *Converter) EnterScopeDo(ctx context.Context, command domain.Command, baseTarget domain.Target, allowPrivileged, passArgs bool, scopeName string, buildArgs []string) error {
	topArgs := buildArgs
	if c.ftrs.ArgScopeSet {
		tmpScope, err := variables.ParseArgs(buildArgs, nil, nil)
		if err != nil {
			return err
		}
		topArgs = variables.CombineScopes(tmpScope, c.varCollection.TopOverriding()).BuildArgs()
	}

	baseMts, err := c.buildTarget(ctx, baseTarget.String(), c.platr.Current(), allowPrivileged, passArgs, topArgs, true, enterScopeDoCmd, "")
	if err != nil {
		return err
	}

	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}
	overriding, err := variables.ParseArgs(buildArgs, pncvf, c.varCollection)
	if err != nil {
		return err
	}
	if passArgs {
		overriding = variables.CombineScopesInactive(overriding, c.varCollection.Overriding(), c.varCollection.Args(), c.varCollection.Globals())
	}
	c.varCollection.EnterFrame(
		scopeName, command, overriding, baseMts.Final.VarCollection.Globals(),
		baseMts.Final.GlobalImports)
	return nil
}

// ExitScope exits the most recent variable scope.
func (c *Converter) ExitScope(ctx context.Context) error {
	c.varCollection.ExitFrame()
	return nil
}

// StackString string returns the current command stack string.
func (c *Converter) StackString() string {
	return c.varCollection.StackString()
}

// FinalizeStates returns the LLB states.
func (c *Converter) FinalizeStates(ctx context.Context) (*states.MultiTarget, error) {
	c.markFakeDeps()

	if !c.varCollection.IsStackAtBase() {
		// Should never happen.
		return nil, errors.New("internal error: stack not at base in FinalizeStates")
	}

	// Persists any cache directories created by using a `CACHE` command
	c.mts.Final.MainState = c.persistCache(c.mts.Final.MainState)

	c.mts.Final.PlatformResolver = c.platr
	c.mts.Final.VarCollection = c.varCollection
	c.mts.Final.GlobalImports = c.varCollection.Imports().Global()
	if c.opt.DoSaves {
		c.mts.Final.SetDoSaves()
	}
	if c.opt.DoPushes {
		c.mts.Final.SetDoPushes()
	}

	if c.ftrs.WaitBlock {
		c.waitBlock().AddItem(newStateWaitItem(&c.mts.Final.MainState, c))
	}
	close(c.mts.Final.Done())

	// Force execution asynchronously, and then mark the logbusTarget as finished.
	// This ensures that the execution actually took place, for timing purposes.
	c.opt.ErrorGroup.Go(func() error {
		rel, err := c.opt.Parallelism.Acquire(ctx, 1)
		if err != nil {
			return errors.Wrapf(err, "acquiring parallelism semaphore for %s", c.mts.FinalTarget().String())
		}
		defer rel()
		if c.ftrs.ExecAfterParallel {
			err = c.forceExecution(ctx, c.mts.Final.MainState, c.mts.Final.PlatformResolver)
			if err != nil {
				return errors.Wrapf(err, "async force execution for %s", c.mts.FinalTarget().String())
			}
		}
		now := time.Now()
		st := logstream.RunStatus_RUN_STATUS_SUCCESS
		c.logbusTarget.SetEnd(now, st, c.platr.Current().String())
		return nil
	})
	return c.mts, nil
}

// RecordTargetFailure records a failure in a target.
func (c *Converter) RecordTargetFailure(ctx context.Context, err error) {
	var st logstream.RunStatus
	switch {
	case errors.Is(err, context.Canceled) || status.Code(errors.Cause(err)) == codes.Canceled:
		st = logstream.RunStatus_RUN_STATUS_CANCELED
	default:
		st = logstream.RunStatus_RUN_STATUS_FAILURE
	}
	now := time.Now()
	c.logbusTarget.SetEnd(now, st, c.platr.Current().String())
}

var errShellOutNotPermitted = errors.New("shell-out not permitted")

// ExpandArgs expands args in the provided word.
func (c *Converter) ExpandArgs(ctx context.Context, runOpts ConvertRunOpts, word string, allowShellOut bool) (string, error) {
	if !c.opt.Features.ShellOutAnywhere {
		return c.varCollection.ExpandOld(word), nil
	}
	return c.varCollection.Expand(word, func(cmd string) (string, error) {
		if !allowShellOut {
			return "", errShellOutNotPermitted
		}
		runOpts.Args = []string{cmd}
		return c.RunCommand(ctx, "internal-expand-args", runOpts)
	})
}

func (c *Converter) absolutizeTarget(fullTargetName string, allowPrivileged bool) (domain.Target, domain.Target, bool, error) {
	relTarget, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		return domain.Target{}, domain.Target{}, false, errors.Wrapf(err, "earthly target parse %s", fullTargetName)
	}

	derefedTarget, allowPrivilegedImport, isImport, err := c.varCollection.Imports().Deref(relTarget)
	if err != nil {
		return domain.Target{}, domain.Target{}, false, err
	}

	if isImport {
		allowPrivileged = allowPrivileged && allowPrivilegedImport
	}

	targetRef, err := c.joinRefs(derefedTarget)
	if err != nil {
		return domain.Target{}, domain.Target{}, false, errors.Wrap(err, "join targets")
	}

	return targetRef.(domain.Target), relTarget, allowPrivileged, nil
}

func (c *Converter) checkAutoSkip(ctx context.Context, fullTargetName string, allowPrivileged, passArgs bool, buildArgs []string) (bool, func(), error) {
	console := c.opt.Console.WithPrefix("auto-skip")

	nopFn := func() {}

	if !c.opt.Features.BuildAutoSkip {
		return false, nopFn, nil
	}

	if c.opt.BuildkitSkipper == nil {
		console.Warnf("BUILD --auto-skip option disabled due to client initialization failure")
		return false, nopFn, nil
	}

	if c.opt.NoAutoSkip {
		console.VerbosePrintf("BUILD --auto-skip ignored due to --no-auto-skip flag")
		return false, nopFn, nil
	}

	target, relTarget, _, err := c.absolutizeTarget(fullTargetName, allowPrivileged)
	if err != nil {
		return false, nil, err
	}

	overriding, _, err := c.prepOverridingVars(ctx, relTarget, passArgs, buildArgs)
	if err != nil {
		return false, nil, err
	}

	targetHash, _, err := inputgraph.HashTarget(ctx, inputgraph.HashOpt{
		Target:         target,
		Console:        c.opt.Console,
		CI:             c.opt.IsCI,
		BuiltinArgs:    c.opt.BuiltinArgs,
		OverridingVars: overriding,
	})
	if err != nil {
		return false, nil, errors.Wrapf(err, "auto-skip is unable to calculate hash for %s", target)
	}

	orgName := c.varCollection.Org()

	exists, err := c.opt.BuildkitSkipper.Exists(ctx, orgName, targetHash)
	if err != nil {
		console.Warnf("Unable to check if target %s (hash %x) has already been run: %s", target.String(), targetHash, err.Error())
		return false, nopFn, nil
	}

	if exists {
		console.Printf("Target %s (hash %x) has already been run. Skipping.", target.String(), targetHash)
		return true, nil, nil
	}

	return exists, func() {
		err := c.opt.BuildkitSkipper.Add(ctx, orgName, target.StringCanonical(), targetHash)
		if err != nil {
			console.Warnf("Failed to add target %s (hash %x) to the auto-skip DB.", target.String(), targetHash)
		}
	}, nil
}

func (c *Converter) prepOverridingVars(ctx context.Context, relTarget domain.Target, passArgs bool, buildArgs []string) (*variables.Scope, bool, error) {
	var buildArgFunc variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		buildArgFunc = c.processNonConstantBuildArgFunc(ctx)
	}

	overriding, err := variables.ParseArgs(buildArgs, buildArgFunc, c.varCollection)
	if err != nil {
		return nil, false, errors.Wrap(err, "parse build args")
	}

	// Don't allow transitive overriding variables to cross project boundaries (unless --pass-args is used).
	propagateBuildArgs := !relTarget.IsExternal()
	if passArgs {
		overriding = variables.CombineScopes(overriding, c.varCollection.Overriding(), c.varCollection.Args(), c.varCollection.Globals())
	} else if propagateBuildArgs {
		overriding = variables.CombineScopes(overriding, c.varCollection.Overriding())
	}

	return overriding, propagateBuildArgs, nil
}

func (c *Converter) prepBuildTarget(
	ctx context.Context,
	fullTargetName string,
	platform platutil.Platform,
	allowPrivileged, passArgs bool,
	buildArgs []string,
	isDangling bool,
	cmdT cmdType,
	parentCmdID string,
) (domain.Target, ConvertOpt, bool, error) {
	target, relTarget, allowPrivileged, err := c.absolutizeTarget(fullTargetName, allowPrivileged)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, err
	}

	overriding, propagateBuildArgs, err := c.prepOverridingVars(ctx, relTarget, passArgs, buildArgs)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, err
	}

	// Recursion.
	opt := c.opt
	opt.OverridingVars = overriding
	opt.GlobalImports = nil
	opt.parentDepSub = c.mts.Final.NewDependencySubscription()
	opt.PlatformResolver = c.platr.SubResolver(platform)
	opt.HasDangling = isDangling
	opt.AllowPrivileged = allowPrivileged
	opt.parentTargetID = c.mts.Final.ID
	opt.parentCommandID = parentCmdID

	if cmdT == buildCmd {
		// only BUILD commands get propagated
		opt.waitBlock = c.waitBlock()
	} else {
		// FROM/COPY commands will return a llb state, which will cause a wait to occur
		// if the wait block was passed here, calling SetDoSaves would get propagated
		opt.waitBlock = nil
	}

	if c.opt.Features.ReferencedSaveOnly {
		// DoSaves should only be potentially turned-off when the ReferencedSaveOnly feature is flipped
		opt.DoSaves = (cmdT == buildCmd && c.opt.DoSaves && !c.opt.OnlyFinalTargetImages)
		opt.DoPushes = (cmdT == buildCmd && c.opt.DoPushes)
		opt.ForceSaveImage = false
	} else {
		opt.DoSaves = c.opt.DoSaves && !target.IsRemote()   // legacy mode only saves artifacts from local targets
		opt.DoPushes = c.opt.DoPushes && !target.IsRemote() // legacy mode only saves artifacts from local targets
	}
	return target, opt, propagateBuildArgs, nil
}

func (c *Converter) buildTarget(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged, passArgs bool, buildArgs []string, isDangling bool, cmdT cmdType, parentCmdID string) (*states.MultiTarget, error) {
	target, opt, propagateBuildArgs, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, passArgs, buildArgs, isDangling, cmdT, parentCmdID)
	if err != nil {
		return nil, err
	}
	mts, err := Earthfile2LLB(ctx, target, opt, false)
	if err != nil {
		return nil, errors.Wrapf(err, "earthfile2llb for %s", fullTargetName)
	}
	c.directDeps = append(c.directDeps, mts.Final)
	if propagateBuildArgs {
		// Propagate build arg inputs upwards (a child target depending on a build arg means
		// that the parent also depends on that build arg).
		for _, bai := range mts.Final.TargetInput().BuildArgs {
			// Check if the build arg has been overridden. If it has, it can no longer be an input
			// directly, so skip it.
			_, found := opt.OverridingVars.Get(bai.Name)
			if found {
				continue
			}
			c.mts.Final.AddBuildArgInput(bai)
		}
		if cmdT == fromCmd {
			// Propagate globals.
			globals := mts.Final.VarCollection.Globals()
			for _, k := range globals.Sorted(variables.WithActive()) {
				_, alreadyActive := c.varCollection.Get(k, variables.WithActive())
				if alreadyActive {
					// Globals don't override any variables in current scope.
					continue
				}
				v, _ := globals.Get(k, variables.WithActive())
				// Look for the default arg value in the built target's TargetInput.
				defaultArgValue := ""
				for _, childBai := range mts.Final.TargetInput().BuildArgs {
					if childBai.Name == k {
						defaultArgValue = childBai.DefaultValue
						break
					}
				}
				c.mts.Final.AddBuildArgInput(
					dedup.BuildArgInput{
						Name:          k,
						DefaultValue:  defaultArgValue,
						ConstantValue: v,
					})
			}
			c.varCollection.SetGlobals(globals)
			c.varCollection.Imports().SetGlobal(mts.Final.GlobalImports)
			c.varCollection.SetProject(mts.Final.VarCollection.Project())
			c.varCollection.SetOrg(mts.Final.VarCollection.Org())
		}
	}

	return mts, nil
}

func getDebuggerSecretKey(saveFilesSettings []debuggercommon.SaveFilesSettings) string {
	h := sha1.New()
	b := make([]byte, 8)

	addToHash := func(path string) {
		h.Write([]byte(path))
		inode := inodeutil.GetInodeBestEffort(path)
		binary.LittleEndian.PutUint64(b, inode)
		h.Write(b)
	}

	for _, saveFile := range saveFilesSettings {
		addToHash(saveFile.Src)
		addToHash(saveFile.Dst)
	}
	return hex.EncodeToString(h.Sum(nil))

}

func (c *Converter) internalRun(ctx context.Context, opts ConvertRunOpts) (pllb.State, error) {
	isInteractive := (opts.Interactive || opts.InteractiveKeep)
	if !c.opt.AllowInteractive && isInteractive {
		return pllb.State{}, errors.New("interactive options are not allowed, when --strict is specified or otherwise implied")
	}
	if opts.Locally {
		if len(opts.Secrets) != 0 {
			return pllb.State{}, errors.New("secrets not yet supported with LOCALLY") // TODO
		}
		if len(opts.Mounts) != 0 {
			return pllb.State{}, errors.New("mounts not supported with LOCALLY")
		}
		if opts.WithSSH {
			return pllb.State{}, errors.New("--ssh not supported with LOCALLY")
		}
		if opts.Privileged {
			return pllb.State{}, errors.New("--privileged not supported with LOCALLY")
		}
		if isInteractive {
			return pllb.State{}, errors.New("interactive mode not supported with LOCALLY")
		}
		if opts.Push {
			return pllb.State{}, errors.New("--push not supported with LOCALLY")
		}
		if opts.Transient {
			return pllb.State{}, errors.New("Transient run not supported with LOCALLY")
		}
		if opts.NoNetwork {
			return pllb.State{}, errors.New("--network=none is not supported with LOCALLY")
		}
	}
	if opts.shellWrap == nil {
		opts.shellWrap = withShellAndEnvVars
	}

	finalArgs := opts.Args[:]
	if opts.WithEntrypoint {
		if len(finalArgs) == 0 {
			// No args provided. Use the image's CMD.
			args := make([]string, len(c.mts.Final.MainImage.Config.Cmd))
			copy(args, c.mts.Final.MainImage.Config.Cmd)
		}
		finalArgs = append(c.mts.Final.MainImage.Config.Entrypoint, finalArgs...)
		opts.WithShell = false // Don't use shell when --entrypoint is passed.
	}

	runOpts := opts.extraRunOpts[:]
	if opts.Privileged {
		runOpts = append(runOpts, llb.Security(llb.SecurityModeInsecure))
	}
	mountRunOpts, err := c.parseMounts(opts.Mounts)
	if err != nil {
		return pllb.State{}, errors.Wrap(err, "parse mounts")
	}
	if opts.NoNetwork {
		runOpts = append(runOpts, llb.Network(llb.NetModeNone))
	}

	runOpts = append(runOpts, mountRunOpts...)
	commandStr := fmt.Sprintf(
		"%s %s%s%s%s%s%s%s%s",
		opts.CommandName, // e.g. "RUN", "IF", "FOR", "ARG"
		strIf(opts.Privileged, "--privileged "),
		strIf(opts.Push, "--push "),
		strIf(opts.WithSSH, "--ssh "),
		strIf(opts.NoCache, "--no-cache "),
		strIf(opts.NoNetwork, "--network=none "),
		strIf(opts.Interactive, "--interactive "),
		strIf(opts.InteractiveKeep, "--interactive-keep "),
		strings.Join(opts.Args, " "))

	prefix, _, err := c.newVertexMeta(ctx, opts.Locally, isInteractive, false, opts.Secrets)
	if err != nil {
		return pllb.State{}, err
	}
	runOpts = append(runOpts, llb.WithCustomNamef("%s%s", prefix, commandStr))

	var extraEnvVars []string

	// Build args.
	for _, buildArgName := range c.varCollection.SortedVariables(variables.WithActive()) {
		ba, _ := c.varCollection.Get(buildArgName, variables.WithActive())
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=%s", buildArgName, shellescape.Quote(ba)))
	}
	// Secrets.
	for _, secretKeyValue := range opts.Secrets {
		secretName, envVar, err := c.parseSecretFlag(secretKeyValue)
		if err != nil {
			return pllb.State{}, err
		}
		if secretName != "" {
			secretPath := path.Join("/run/secrets", secretName)
			secretOpts := []llb.SecretOption{
				llb.SecretID(c.secretID(secretName)),
				// TODO: Perhaps this should just default to the current user automatically from
				//       buildkit side. Then we wouldn't need to open this up to everyone.
				llb.SecretFileOpt(0, 0, 0444),
			}
			runOpts = append(runOpts, llb.AddSecret(secretPath, secretOpts...))
			// TODO: The use of cat here might not be portable.
			extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"$(cat %s)\"", envVar, secretPath))
		}
	}
	if !opts.Locally {
		// Debugger.
		err := c.opt.LLBCaps.Supports(solverpb.CapExecMountSock)
		if err != nil {
			switch errors.Cause(err).(type) {
			case *apicaps.CapError:
				if c.opt.InteractiveDebuggerEnabled || isInteractive {
					return pllb.State{}, errors.Wrap(err, "interactive debugger requires a newer version of buildkit")
				}
			default:
				c.opt.Console.Warnf("failed to check LLBCaps for CapExecMountSock: %v", err) // keep going
			}
		} else {
			runOpts = append(runOpts, llb.SocketTarget("earthly_interactive", debuggercommon.DebuggerDefaultSocketPath, 0666, 0, 0))
			runOpts = append(runOpts, llb.SocketTarget("earthly_save_file", debuggercommon.DefaultSaveFileSocketPath, 0666, 0, 0))
		}

		localPathAbs, err := filepath.Abs(c.target.LocalPath)
		if err != nil {
			return pllb.State{}, errors.Wrapf(err, "unable to determine absolute path of %s", c.target.LocalPath)
		}
		saveFiles := []debuggercommon.SaveFilesSettings{}
		for _, interactiveSaveFile := range opts.InteractiveSaveFiles {

			canSave, err := c.canSave(ctx, interactiveSaveFile.Dst)
			if err != nil {
				return pllb.State{}, err
			}
			if !canSave {
				return pllb.State{}, fmt.Errorf("unable to save to %s; path must be located under %s", interactiveSaveFile.Dst, c.target.LocalPath)
			}
			dst := path.Join(localPathAbs, interactiveSaveFile.Dst)
			c.opt.LocalArtifactWhiteList.Add(dst)
			// The receiveFile handler will only be given the relative, so needs to whitelisted as well.
			// This is needed when e.g. the user specifies a destination of "out/", which is then rewritten
			// to "out/file".  If the user specified a full file path, e.g. "out/file", then this is redundant
			// since the waitBlock will add that same value.  This makes it explicit.
			c.opt.LocalArtifactWhiteList.Add(interactiveSaveFile.Dst)

			saveFiles = append(saveFiles, debuggercommon.SaveFilesSettings{
				Src:      interactiveSaveFile.Src,
				Dst:      interactiveSaveFile.Dst,
				IfExists: interactiveSaveFile.IfExists,
			})
		}

		debuggerSettingsSecretsKey := getDebuggerSecretKey(saveFiles)
		debuggerSettings := debuggercommon.DebuggerSettings{
			DebugLevelLogging: c.opt.InteractiveDebuggerDebugLevelLogging,
			Enabled:           c.opt.InteractiveDebuggerEnabled,
			SocketPath:        debuggercommon.DebuggerDefaultSocketPath,
			Term:              os.Getenv("TERM"),
			SaveFiles:         saveFiles,
		}
		debuggerSettingsData, err := json.Marshal(&debuggerSettings)
		if err != nil {
			return pllb.State{}, errors.Wrap(err, "debugger settings json marshal")
		}
		err = c.opt.InternalSecretStore.SetSecret(ctx, c.secretID(debuggerSettingsSecretsKey), debuggerSettingsData)
		if err != nil {
			return pllb.State{}, errors.Wrap(err, "InternalSecretStore.SetSecret")
		}

		secretOpts := []llb.SecretOption{
			llb.SecretID(c.secretID(debuggerSettingsSecretsKey)),
			llb.SecretFileOpt(0, 0, 0444),
		}
		debuggerSecretMount := llb.AddSecret(
			fmt.Sprintf("/run/secrets/%s", debuggercommon.DebuggerSettingsSecretsKey), secretOpts...)
		debuggerMount := pllb.AddMount(debuggerPath, pllb.Scratch(),
			llb.HostBind(), llb.SourcePath("/usr/bin/earth_debugger"))
		runOpts = append(runOpts, debuggerSecretMount, debuggerMount)
		if opts.WithSSH {
			runOpts = append(runOpts, llb.AddSSHSocket())
		}
	}
	// Shell and debugger wrap.
	prependDebugger := !opts.Locally
	finalArgs = opts.shellWrap(finalArgs, extraEnvVars, opts.WithShell, prependDebugger, isInteractive)
	if opts.Locally {
		// buildkit-hack in order to run locally, we prepend the command with a magic UUID.
		finalArgs = append(
			[]string{localhost.RunOnLocalHostMagicStr},
			finalArgs...)
	}

	if c.ftrs.WaitBlock && opts.Push {
		// The WAIT / END feature treats a --push as syntactic sugar for
		// IF [ "$EARTHLY_PUSH" = "true" ]
		//    RUN --no-cache ...
		// END
		if !c.opt.DoPushes {
			// quick return when EARTHLY_PUSH != true
			return c.mts.Final.MainState, nil
		}

		// convert "push" to "no-cache"
		opts.Push = false
		opts.NoCache = true
	}

	runOpts = append(runOpts, llb.Args(finalArgs))
	if opts.NoCache || opts.Locally || opts.Push || isInteractive {
		runOpts = append(runOpts, llb.IgnoreCache)
	}

	if opts.Push {
		if !c.mts.Final.RunPush.HasState {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.Final.RunPush.State = c.mts.Final.MainState
			c.mts.Final.RunPush.HasState = true
		}
	}

	var state pllb.State
	if opts.Push {
		state = c.mts.Final.RunPush.State
	} else {
		state = c.mts.Final.MainState
	}
	if opts.statePrep != nil {
		state, err = opts.statePrep(ctx, state)
		if err != nil {
			return pllb.State{}, err
		}
	}
	if isInteractive {
		c.mts.Final.RanInteractive = true

		switch {
		case opts.Interactive:
			is := states.InteractiveSession{
				CommandStr:  commandStr,
				Initialized: true,
				Kind:        states.SessionEphemeral,
			}

			if opts.Push {
				is.State = state.Run(runOpts...).Root()
				c.mts.Final.RunPush.InteractiveSession = is
				return c.mts.Final.RunPush.State, nil

			}
			is.State = state.Run(runOpts...).Root()
			c.mts.Final.InteractiveSession = is
			return c.mts.Final.MainState, nil

		case opts.InteractiveKeep:
			c.mts.Final.InteractiveSession = states.InteractiveSession{
				CommandStr:  commandStr,
				Initialized: true,
				Kind:        states.SessionKeep,
			}
		}
	}

	if opts.Push {
		// Don't run on MainState. We want push-flagged commands to be executed only
		// *after* the build. Save this for later.
		c.mts.Final.RunPush.State = state.Run(runOpts...).Root()
		c.mts.Final.RunPush.CommandStrs = append(c.mts.Final.RunPush.CommandStrs, commandStr)
		return c.mts.Final.RunPush.State, nil
	} else if opts.Transient {
		transientState := state.Run(runOpts...).Root()
		return transientState, nil
	} else {
		c.mts.Final.MainState = state.Run(runOpts...).Root()

		if opts.Locally {
			err = c.forceExecution(ctx, c.mts.Final.MainState, c.platr)
			if err != nil {
				return pllb.State{}, err
			}
		}

		return c.mts.Final.MainState, nil
	}
}

// secretID returns query parameter style string that contains the secret
// version, name, org, and project. The version value informs the secret
// providers how to handle the secret name and whether to use the new
// project-based secrets endpoints.
func (c *Converter) secretID(name string) string {
	v := url.Values{}
	v.Set("name", name)
	if c.ftrs.UseProjectSecrets {
		v.Set("v", "1")
		v.Set("org", c.varCollection.Org())
		v.Set("project", c.varCollection.Project())
	} else {
		v.Set("v", "0")
	}
	return v.Encode()
}

func (c *Converter) parseSecretFlag(secretKeyValue string) (secretID string, envVar string, err error) {
	if secretKeyValue == "" {
		// If empty string, don't use (used for optional secrets).
		// TODO: This should be an actual secret (with an empty value),
		//       so that the cache works correctly.
		return "", "", nil
	}
	parts := strings.SplitN(secretKeyValue, "=", 2)

	// validate environment name is correct
	defer func() {
		if err != nil {
			return
		}
		if envVar != "" && !shell.IsValidEnvVarName(envVar) {
			err = fmt.Errorf("invalid secret environment name: %s", envVar)
			secretID = ""
			envVar = ""
			return
		}
	}()

	if len(parts) == 1 {
		return parts[0], parts[0], nil
	}

	secretID = parts[1]

	if secretID == "" {
		// If empty string, don't use (used for optional secrets).
		// TODO: This should be an actual secret (with an empty value),
		//       so that the cache works correctly.
		return "", "", nil
	}

	if c.ftrs.UseProjectSecrets {
		if strings.HasPrefix(secretID, "+secrets/") {
			secretID = strings.TrimPrefix(secretID, "+secrets/")
			c.opt.Console.Printf("Deprecation: the '+secrets/' prefix is not required and support for it will be removed in an upcoming release")
		}
		return secretID, parts[0], nil
	}

	if strings.HasPrefix(secretID, "+secrets/") {
		secretID = strings.TrimPrefix(secretID, "+secrets/")
		return secretID, parts[0], nil
	}

	return "", "", errors.Errorf("secret definition %s not supported. Format must be either <env-var>=+secrets/<secret-id> or <secret-id>", secretKeyValue)
}

func (c *Converter) forceExecution(ctx context.Context, state pllb.State, platr *platutil.Resolver) error {
	if state.Output() == nil {
		// Scratch - no need to execute.
		return nil
	}
	ref, err := llbutil.StateToRef(
		ctx, c.opt.GwClient, state, c.opt.NoCache,
		platr, c.opt.CacheImports.AsSlice())
	if err != nil {
		return errors.Wrap(err, "force execution state to ref")
	}
	if ref == nil {
		return nil
	}
	// We're not really interested in reading the dir - we just
	// want to un-lazy the ref so that the commands have executed.
	_, err = ref.ReadDir(ctx, gwclient.ReadDirRequest{Path: "/"})
	if err != nil {
		return errors.Wrap(err, "unlazy force execution")
	}
	return nil
}

func (c *Converter) readArtifact(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact) ([]byte, error) {
	if mts.Final.ArtifactsState.Output() == nil {
		// ArtifactsState is scratch - no artifact has been copied.
		return nil, errors.Errorf("artifact %s not found; no SAVE ARTIFACT command was issued in %s", artifact.String(), artifact.Target.String())
	}
	ref, err := llbutil.StateToRef(
		ctx, c.opt.GwClient, mts.Final.ArtifactsState, c.opt.NoCache,
		mts.Final.PlatformResolver, c.opt.CacheImports.AsSlice())
	if err != nil {
		return nil, errors.Wrap(err, "state to ref solve artifact")
	}
	artDt, err := ref.ReadFile(ctx, gwclient.ReadRequest{
		Filename: artifact.Artifact,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "read artifact %s", artifact.String())
	}
	return artDt, nil
}

func (c *Converter) internalFromClassical(ctx context.Context, imageName string, platform platutil.Platform, opts ...llb.ImageOption) (pllb.State, *image.Image, *variables.Scope, error) {
	llbPlatform := c.platr.ToLLBPlatform(platform)
	if imageName == "scratch" {
		// FROM scratch
		img := image.NewImage()
		img.OS = llbPlatform.OS
		img.Architecture = llbPlatform.Architecture
		return pllb.Scratch().Platform(llbPlatform), img, nil, nil
	}
	sourceRef, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", imageName)
	}
	baseImageName := reference.TagNameOnly(sourceRef).String()
	logName := fmt.Sprintf(
		"%sLoad metadata %s %s",
		c.imageVertexPrefix(imageName, platform), imageName, platforms.Format(llbPlatform))
	ref, dgst, dt, err := c.opt.MetaResolver.ResolveImageConfig(
		ctx, baseImageName,
		llb.ResolveImageConfigOpt{
			Platform:    &llbPlatform,
			ResolveMode: c.opt.ImageResolveMode.String(),
			LogName:     logName,
		})
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "resolve image config for %s", imageName)
	}
	sourceRef, err = reference.ParseNormalizedNamed(ref)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", ref)
	}
	var img image.Image
	err = json.Unmarshal(dt, &img)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "unmarshal image config for %s", imageName)
	}
	if dgst != "" {
		sourceRef, err = reference.WithDigest(sourceRef, dgst)
		if err != nil {
			return pllb.State{}, nil, nil, errors.Wrapf(err, "reference add digest %v for %s", dgst, imageName)
		}
	}
	allOpts := append(opts, llb.Platform(c.platr.ToLLBPlatform(platform)), c.opt.ImageResolveMode)
	state := pllb.Image(sourceRef.String(), allOpts...)
	state, img2, envVars := c.applyFromImage(state, &img)
	return state, img2, envVars, nil
}

func (c *Converter) checkOldPlatformIncompatibility(platform platutil.Platform) error {
	if c.ftrs.NewPlatform {
		return nil
	}
	if c.platr.Default() == platutil.DefaultPlatform || platform == platutil.DefaultPlatform {
		return nil
	}
	if !c.platr.PlatformEquals(c.platr.Default(), platform) {
		return errors.Errorf(
			"platform contradiction: \"%s\" vs \"%s\"",
			platform.String(), c.platr.Default().String())
	}
	return nil
}

func (c *Converter) applyFromImage(state pllb.State, img *image.Image) (pllb.State, *image.Image, *variables.Scope) {
	// Reset variables.
	ev := variables.ParseEnvVars(img.Config.Env)
	for _, name := range ev.Sorted(variables.WithActive()) {
		v, _ := ev.Get(name, variables.WithActive())
		state = state.AddEnv(name, v)
	}
	// Init config maps if not already initialized.
	if img.Config.ExposedPorts == nil {
		img.Config.ExposedPorts = make(map[string]struct{})
	}
	if img.Config.Labels == nil {
		img.Config.Labels = make(map[string]string)
	}
	if img.Config.Volumes == nil {
		img.Config.Volumes = make(map[string]struct{})
	}
	if img.Config.WorkingDir != "" {
		state = state.Dir(img.Config.WorkingDir)
	}
	if img.Config.User != "" {
		state = state.User(img.Config.User)
	}
	// No need to apply entrypoint, cmd, volumes and others.
	// The fact that they exist in the image configuration is enough.
	// TODO: Apply any other settings? Shell?
	return state, img, ev
}

func (c *Converter) nonSaveCommand() {
	if c.ranSave {
		c.mts.Final.HasDangling = true
	}
}

func (c *Converter) processNonConstantBuildArgFunc(ctx context.Context) variables.ProcessNonConstantVariableFunc {
	return func(name string, expression string) (string, int, error) {
		opts := ConvertRunOpts{
			CommandName: fmt.Sprintf("ARG %s = RUN", name),
			Args:        strings.Split(expression, " "),
			Transient:   true,
			WithShell:   true,
		}
		output, err := c.RunExpression(ctx, name, opts)
		if err != nil {
			return "", 0, err
		}
		return output, 0, nil
	}
}

func (c *Converter) newLogbusCommand(ctx context.Context, name string) (string, *logbus.Command, error) {

	cmdID := fmt.Sprintf("%s/%d", c.mts.Final.ID, c.newCmdID())

	var gitURL, gitHash, fileRelToRepo string
	if c.gitMeta != nil {
		gitURL = c.gitMeta.RemoteURL
		gitHash = c.gitMeta.Hash
		fileRelToRepo = path.Join(c.gitMeta.RelDir, "Earthfile")
	}

	platform := c.platr.Materialize(c.platr.Current())

	cmd, err := c.opt.Logbus.Run().NewCommand(
		cmdID,
		name,
		c.mts.Final.ID,
		c.mts.Final.Target.String(),
		platform.String(),
		false, // cached
		false, // local
		false, // interactive
		SourceLocationFromContext(ctx),
		gitURL,
		gitHash,
		fileRelToRepo,
	)
	if err != nil {
		return "", nil, err
	}

	cmd.SetStart(time.Now())

	return cmdID, cmd, nil
}

func (c *Converter) newVertexMeta(ctx context.Context, local, interactive, internal bool, secrets []string) (string, string, error) {
	activeOverriding := make(map[string]string)
	for _, arg := range c.varCollection.SortedOverridingVariables() {
		v, ok := c.varCollection.Get(arg, variables.WithActive())
		if ok {
			activeOverriding[arg] = v
		}
	}

	platform := c.platr.Materialize(c.platr.Current())
	platformStr := platform.String()
	isNativePlatform := c.platr.PlatformEquals(platform, platutil.NativePlatform)

	var gitURL, gitHash, fileRelToRepo string
	if c.gitMeta != nil {
		gitURL = c.gitMeta.RemoteURL
		gitHash = c.gitMeta.Hash
		fileRelToRepo = path.Join(c.gitMeta.RelDir, "Earthfile")
	}

	var (
		srcLoc = SourceLocationFromContext(ctx)
		cmdID  = fmt.Sprintf("%s/%d", c.mts.Final.ID, c.newCmdID())
		name   = "" // Name is initially empty. It will be set by SolverMonitor in most cases.
	)

	_, err := c.opt.Logbus.Run().NewCommand(
		cmdID,
		name,
		c.mts.Final.ID,
		c.mts.Final.Target.String(),
		platformStr,
		false, // cached
		local,
		interactive,
		srcLoc,
		gitURL,
		gitHash,
		fileRelToRepo,
	)
	if err != nil {
		return "", "", err
	}

	vm := &vertexmeta.VertexMeta{
		SourceLocation:      srcLoc,
		RepoGitURL:          gitURL,
		RepoGitHash:         gitHash,
		RepoFileRelToRepo:   fileRelToRepo,
		CommandID:           cmdID,
		TargetID:            c.mts.Final.ID,
		TargetName:          c.mts.Final.Target.String(),
		CanonicalTargetName: c.mts.Final.Target.StringCanonical(),
		Platform:            platformStr,
		NonDefaultPlatform:  !isNativePlatform,
		Local:               local,
		Interactive:         interactive,
		OverridingArgs:      activeOverriding,
		Secrets:             secrets,
		Internal:            internal,
		Runner:              c.opt.Runner,
	}

	return vm.ToVertexPrefix(), cmdID, nil
}

func (c *Converter) imageVertexPrefix(id string, platform platutil.Platform) string {
	platform = c.platr.Materialize(platform)
	isNativePlatform := c.platr.PlatformEquals(platform, platutil.NativePlatform)
	vm := &vertexmeta.VertexMeta{
		TargetName:         id,
		Platform:           platform.String(),
		NonDefaultPlatform: !isNativePlatform,
	}
	return vm.ToVertexPrefix()
}

func (c *Converter) vertexMetaWithURL(url string) string {
	return fmt.Sprintf("[%s(%s)] ", c.mts.Final.Target.String(), url)
}

func (c *Converter) markFakeDeps() {
	if !c.opt.UseFakeDep {
		return
	}
	for _, dep := range c.directDeps {
		select {
		case <-dep.Done():
		default:
			panic("mark fake dep but dep not done")
		}
		if dep.HasDangling {
			c.mts.Final.MainState = llbutil.WithDependency(
				c.mts.Final.MainState, dep.MainState, c.mts.Final.Target.String(), dep.Target.String(),
				c.platr)
		}
	}
	// Clear the direct deps so we don't do this again.
	c.directDeps = nil
}

func (c *Converter) copyOwner(keepOwn bool, chown string) string {
	own := c.mts.Final.MainImage.Config.User
	if own == "" {
		own = "root:root"
	}
	if keepOwn {
		own = ""
	}
	if chown != "" {
		own = chown
	}
	return own
}

func (c *Converter) setPlatform(platform platutil.Platform) platutil.Platform {
	newPlatform := c.platr.UpdatePlatform(platform)
	c.varCollection.SetPlatform(c.platr)
	return newPlatform
}

func (c *Converter) joinRefs(relRef domain.Reference) (domain.Reference, error) {
	return domain.JoinReferences(c.varCollection.AbsRef(), relRef)
}

func (c *Converter) checkAllowed(command cmdType) error {
	if command == projectCmd && !c.ftrs.UseProjectSecrets {
		return errors.New("--use-project-secrets must be enabled in order to use PROJECT")
	}

	if (command == pipelineCmd || command == triggerCmd) && !c.ftrs.UsePipelines {
		return errors.New("--use-pipelines must be enabled in order to use PIPELINE or TRIGGER")
	}

	if c.mts.Final.RanInteractive && !(command == saveImageCmd || command == saveArtifactCmd) {
		return errors.New("If present, a single --interactive command must be the last command in a target")
	}

	if command == pipelineCmd && c.isPipeline {
		return errors.New("only 1 PIPELINE command is allowed")
	}

	if !c.mts.Final.RanFromLike {
		switch command {
		case fromCmd, fromDockerfileCmd, locallyCmd, buildCmd, argCmd, letCmd, setCmd, importCmd, projectCmd, pipelineCmd:
			return nil
		default:
			return errors.New("the first command has to be FROM, FROM DOCKERFILE, LOCALLY, ARG, BUILD or IMPORT")
		}
	}

	switch command {
	case setCmd, letCmd:
		if !c.ftrs.ArgScopeSet {
			return errors.New("--arg-scope-and-set must be enabled in order to use LET and SET")
		}
	default:
	}

	return nil
}

// persistCache makes temporary cache directories permanent by writing their contents
// from the cached directory to the persistent image layers at the same directory.
// This only has an effect when the Target contains at least one `CACHE /my/directory` command.
func (c *Converter) persistCache(srcState pllb.State) pllb.State {
	dest := srcState
	// User may have multiple CACHE commands in a single target
	for dir, state := range c.persistentCacheDirs {
		if state.Persisted {
			// Copy the contents of the user's cache directory to the temporary backup.
			// It's important to use DockerfileCopy here, since traditional llb.Copy()
			// doesn't support adding mounts via RunOptions.
			runOpts := []llb.RunOption{state.RunOption, llb.WithCustomName("persist cache directory: " + dir)}
			dest = llbutil.CopyWithRunOptions(
				dest,
				dir, // cache dir from external mount
				dir, // cache dir on dest state (same location but without the mount)
				c.platr,
				runOpts...,
			)
		}
	}

	return dest
}

func (c *Converter) newCmdID() int {
	cmdID := c.nextCmdID
	c.nextCmdID++
	return cmdID
}

func joinWrap(a []string, before string, sep string, after string) string {
	if len(a) > 0 {
		return fmt.Sprintf("%s%s%s", before, strings.Join(a, sep), after)
	}
	return ""
}

func strIf(condition bool, str string) string {
	if condition {
		return str
	}
	return ""
}
