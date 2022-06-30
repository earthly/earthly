package earthfile2llb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/variables"
	"github.com/earthly/earthly/variables/reserved"

	"github.com/alessio/shellescape"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session/localhost"
	solverpb "github.com/moby/buildkit/solver/pb"
	"github.com/pkg/errors"
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
	persistentCacheDirs map[string]llb.RunOption // maps path->mount
	varCollection       *variables.Collection
	ranSave             bool
	cmdSet              bool
	ftrs                *features.Features
	localWorkingDir     string
	containerFrontend   containerutil.ContainerFrontend
	waitBlockStack      []*waitBlock
}

// NewConverter constructs a new converter for a given earthly target.
func NewConverter(ctx context.Context, target domain.Target, bc *buildcontext.Data, sts *states.SingleTarget, opt ConvertOpt) (*Converter, error) {
	opt.BuildContextProvider.AddDirs(bc.LocalDirs)
	sts.HasDangling = opt.HasDangling
	mts := &states.MultiTarget{
		Final:   sts,
		Visited: opt.Visited,
	}
	sts.AddOverridingVarsAsBuildArgInputs(opt.OverridingVars)
	newCollOpt := variables.NewCollectionOpt{
		Console:          opt.Console,
		Target:           target,
		PlatformResolver: opt.PlatformResolver,
		GitMeta:          bc.GitMetadata,
		BuiltinArgs:      opt.BuiltinArgs,
		OverridingVars:   opt.OverridingVars,
		GlobalImports:    opt.GlobalImports,
		Features:         opt.Features,
	}
	return &Converter{
		target:              target,
		gitMeta:             bc.GitMetadata,
		platr:               opt.PlatformResolver,
		opt:                 opt,
		mts:                 mts,
		buildContextFactory: bc.BuildContextFactory,
		cacheContext:        pllb.Scratch(),
		persistentCacheDirs: make(map[string]llb.RunOption),
		varCollection:       variables.NewCollection(newCollOpt),
		ftrs:                bc.Features,
		localWorkingDir:     filepath.Dir(bc.BuildFilePath),
		containerFrontend:   opt.ContainerFrontend,
		waitBlockStack:      []*waitBlock{opt.waitBlock},
	}, nil
}

// From applies the earthly FROM command.
func (c *Converter) From(ctx context.Context, imageName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	err := c.checkAllowed(fromCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	if len(c.persistentCacheDirs) > 0 {
		c.persistentCacheDirs = make(map[string]llb.RunOption)
	}
	c.cmdSet = false
	err = c.checkOldPlatformIncompatibility(platform)
	if err != nil {
		return err
	}
	platform = c.setPlatform(platform)
	if strings.Contains(imageName, "+") {
		// Target-based FROM.
		return c.fromTarget(ctx, imageName, platform, allowPrivileged, buildArgs)
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
	state, img, envVars, err := c.internalFromClassical(
		ctx, imageName, platform,
		llb.WithCustomNamef("%sFROM %s", c.vertexPrefix(local, false, internal), imageName))
	if err != nil {
		return err
	}
	c.mts.Final.MainState = state
	c.mts.Final.MainImage = img
	c.mts.Final.RanFromLike = true
	c.varCollection.ResetEnvVars(envVars)
	return nil
}

func (c *Converter) fromTarget(ctx context.Context, targetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}
	mts, err := c.buildTarget(ctx, depTarget.String(), platform, allowPrivileged, buildArgs, false, fromCmd)
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
func (c *Converter) FromDockerfile(ctx context.Context, contextPath string, dfPath string, dfTarget string, platform platutil.Platform, buildArgs []string) error {
	err := c.checkAllowed(fromDockerfileCmd)
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
	var dfData []byte
	if dfPath != "" {
		dfArtifact, parseErr := domain.ParseArtifact(dfPath)
		if parseErr == nil {
			// The Dockerfile is from a target's artifact.
			mts, err := c.buildTarget(ctx, dfArtifact.Target.String(), platform, false, buildArgs, false, fromDockerfileCmd)
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
		// The build context is from a target's artifact.
		// TODO: The build args are used for both the artifact and the Dockerfile. This could be
		//       confusing to the user.
		mts, err := c.buildTarget(ctx, contextArtifact.Target.String(), platform, false, buildArgs, false, fromDockerfileCmd)
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
		BuildContextFactory = llbfactory.PreconstructedState(llbutil.CopyOp(
			mts.Final.ArtifactsState, []string{contextArtifact.Artifact},
			c.platr.Scratch(), "/", true, true, false, "", nil, false, false,
			c.ftrs.UseCopyLink,
			llb.WithCustomNamef(
				"%sFROM DOCKERFILE (copy build context from) %s%s",
				c.vertexPrefix(false, false, true),
				joinWrap(buildArgs, "(", " ", ") "), contextArtifact.String())))
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
	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}
	overriding, err := variables.ParseArgs(buildArgs, pncvf, c.varCollection)
	if err != nil {
		return err
	}
	caps := solverpb.Caps.CapSet(solverpb.Caps.All())
	bcRawState, done := BuildContextFactory.Construct().RawState()
	state, dfImg, err := dockerfile2llb.Dockerfile2LLB(ctx, dfData, dockerfile2llb.ConvertOpt{
		BuildContext:     &bcRawState,
		ContextLocalName: c.mts.FinalTarget().String(),
		MetaResolver:     c.opt.MetaResolver,
		ImageResolveMode: c.opt.ImageResolveMode,
		Target:           dfTarget,
		TargetPlatform:   &plat,
		LLBCaps:          &caps,
		BuildArgs:        overriding.AllValueMap(),
		Excludes:         nil, // TODO: Need to process this correctly.
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

	// reset WORKDIR to current directory where Earthfile is
	c.mts.Final.MainState = c.mts.Final.MainState.Dir(workingDir)
	c.mts.Final.MainImage.Config.WorkingDir = workingDir
	c.setPlatform(platutil.UserPlatform)
	return nil
}

// CopyArtifactLocal applies the earthly COPY artifact command which are invoked under a LOCALLY target.
func (c *Converter) CopyArtifactLocal(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDir bool) error {
	err := c.checkAllowed(copyCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, buildArgs, false, copyCmd)
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
			c.vertexPrefix(false, false, false),
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
func (c *Converter) CopyArtifact(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDir bool, keepTs bool, keepOwn bool, chown string, chmod *fs.FileMode, ifExists, symlinkNoFollow bool) error {
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
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, buildArgs, false, copyCmd)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", artifact.Target.String())
	}
	if artifact.Target.IsLocalInternal() {
		artifact.Target.LocalPath = c.mts.Final.Target.LocalPath
	}
	// Grab the artifacts state in the dep states, after we've built it.
	relevantDepState := mts.Final
	// Copy.
	c.mts.Final.MainState = llbutil.CopyOp(
		relevantDepState.ArtifactsState, []string{artifact.Artifact},
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), chmod, ifExists, symlinkNoFollow,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s%s%s %s",
			c.vertexPrefix(false, false, false),
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest))
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
	c.mts.Final.MainState = llbutil.CopyOp(
		srcState,
		srcs,
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), chmod, ifExists, false,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s %s",
			c.vertexPrefix(false, false, false),
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			strings.Join(srcs, " "),
			dest))
	return nil
}

// ConvertRunOpts represents a set of options needed for the RUN command.
type ConvertRunOpts struct {
	CommandName     string
	Args            []string
	Locally         bool
	Mounts          []string
	Secrets         []string
	WithEntrypoint  bool
	WithShell       bool
	Privileged      bool
	Push            bool
	Transient       bool
	WithSSH         bool
	NoCache         bool
	Interactive     bool
	InteractiveKeep bool

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

	for _, cache := range c.persistentCacheDirs {
		opts.extraRunOpts = append(opts.extraRunOpts, cache)
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
		opts.statePrep = func(ctx context.Context, state pllb.State) (pllb.State, error) {
			return state.File(
				pllb.Mkdir("/run", 0755, llb.WithParents(true)),
				llb.WithCustomNamef(
					"%smkdir %s",
					c.vertexPrefix(false, false, true), "/run"),
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
			c.platr, c.opt.CacheImports.AsMap())
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
	return c.runCommand(ctx, expressionName, true, opts)
}

// RunCommand runs a command and returns its output. The run is transient - any state created
// is not used in subsequent commands.
func (c *Converter) RunCommand(ctx context.Context, commandName string, opts ConvertRunOpts) (string, error) {
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
		outputFile = path.Join(srcBuildArgDir, outputFileName)
		opts.statePrep = func(ctx context.Context, state pllb.State) (pllb.State, error) {
			return state.File(
				pllb.Mkdir(srcBuildArgDir, 0777, llb.WithParents(true)), // Mkdir is performed as root even when USER is set; we must use 0777
				llb.WithCustomNamef(
					"%smkdir %s",
					c.vertexPrefix(false, false, true), srcBuildArgDir),
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
			c.platr, c.opt.CacheImports.AsMap())
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
func (c *Converter) SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string, keepTs bool, keepOwn bool, ifExists, symlinkNoFollow, force bool, isPush bool) error {
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

	c.mts.Final.ArtifactsState = llbutil.CopyOp(
		pcState, []string{saveFrom}, c.mts.Final.ArtifactsState,
		saveToAdjusted, true, true, keepTs, own, nil, ifExists, symlinkNoFollow,
		c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sSAVE ARTIFACT %s%s%s %s",
			c.vertexPrefix(false, false, false),
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			saveFrom,
			artifact.String()))
	if saveAsLocalTo != "" {
		separateArtifactsState := c.platr.Scratch()
		if isPush {
			pushState := c.persistCache(c.mts.Final.RunPush.State)
			separateArtifactsState = llbutil.CopyOp(
				pushState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", nil, ifExists, symlinkNoFollow,
				c.ftrs.UseCopyLink,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					c.vertexPrefix(false, false, false),
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
		} else {
			separateArtifactsState = llbutil.CopyOp(
				pcState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", nil, ifExists, symlinkNoFollow,
				c.ftrs.UseCopyLink,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					c.vertexPrefix(false, false, false),
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
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
		if isPush {
			c.mts.Final.RunPush.SaveLocals = append(c.mts.Final.RunPush.SaveLocals, saveLocal)
		} else {
			c.mts.Final.SaveLocals = append(c.mts.Final.SaveLocals, saveLocal)
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
	opts := []llb.RunOption{
		llb.Args([]string{localhost.CopyFileMagicStr, saveFrom, saveTo}),
		llb.IgnoreCache,
		llb.WithCustomNamef(
			"%sCopyFileMagicStr %s %s",
			c.vertexPrefix(true, false, true), saveFrom, saveTo),
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
	c.mts.Final.ArtifactsState = llbutil.CopyOp(
		c.mts.Final.MainState, []string{absSaveTo}, c.mts.Final.ArtifactsState,
		absSaveTo, true, true, keepTs, own, nil, ifExists, false,
		c.ftrs.UseCopyLink,
	)
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
	c.opt.Console.Warnf("WAIT/END code is experimental and may be incomplete")
	c.waitBlockStack = append(c.waitBlockStack, newWaitBlock())
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
		c.waitBlock().addState(&c.mts.Final.MainState, c)
	}

	i := n - 1
	waitBlock := c.waitBlockStack[i]
	c.waitBlockStack = c.waitBlockStack[:i]

	return waitBlock.wait(ctx)
}

// SaveImage applies the earthly SAVE IMAGE command.
func (c *Converter) SaveImage(ctx context.Context, imageNames []string, pushImages bool, insecurePush bool, cacheHint bool, cacheFrom []string, noManifestList bool) error {
	err := c.checkAllowed(saveImageCmd)
	if err != nil {
		return err
	}
	if noManifestList && !c.ftrs.UseNoManifestList {
		return fmt.Errorf("SAVE IMAGE --no-manifest-list is not supported in this version")
	}
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
					Push:                pushImages,
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
				Push:                pushImages,
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
				shouldPush := pushImages && si.DockerTag != "" && c.opt.DoSaves

				//TODO opt.DoSaves isn't enough here, the original builder code did
				// shouldExport := (!opt.NoOutput && opt.OnlyArtifact == nil && !(opt.OnlyFinalTargetImages && sts != mts.Final) && saveImage.DockerTag != "")
				// however, not all of those are available here right now.
				shouldExportLocally := si.DockerTag != "" && c.opt.DoSaves

				c.waitBlock().addSaveImage(si, c, shouldPush, shouldExportLocally)
			} else {
				c.mts.Final.SaveImages = append(c.mts.Final.SaveImages, si)
			}

		}

		if pushImages && imageName != "" && c.opt.UseInlineCache {
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
func (c *Converter) Build(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	err := c.checkAllowed(buildCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	_, err = c.buildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, true, buildCmd)
	return err
}

type afterParallelFunc func(context.Context, *states.MultiTarget) error

// BuildAsync applies the earthly BUILD command asynchronously.
func (c *Converter) BuildAsync(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, cmdT cmdType, apf afterParallelFunc, sem semutil.Semaphore) error {
	target, opt, _, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, true, cmdT)
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
		if c.ftrs.ExecAfterParallel && mts != nil && mts.Final != nil {
			err = c.forceExecution(ctx, mts.Final.MainState, mts.Final.PlatformResolver)
			if err != nil {
				return errors.Wrapf(err, "async force execution for %s", fullTargetName)
			}
		}
		if apf != nil {
			err = apf(ctx, mts)
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
		opts := []llb.ConstraintsOpt{
			llb.WithCustomNamef("%sWORKDIR %s", c.vertexPrefix(false, false, false), workdirPath),
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
func (c *Converter) Arg(ctx context.Context, argKey string, defaultArgValue string, opts argOpts) error {
	err := c.checkAllowed(argCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()

	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}

	effective, effectiveDefault, err := c.varCollection.DeclareArg(argKey, defaultArgValue, opts.Global, pncvf)
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
func (c *Converter) GitClone(ctx context.Context, gitURL string, branch string, dest string, keepTs bool) error {
	err := c.checkAllowed(gitCloneCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	gitURLScrubbed := stringutil.ScrubCredentials(gitURL)
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef(
			"%sGIT CLONE (--branch %s) %s", c.vertexPrefixWithURL(gitURLScrubbed), branch, gitURLScrubbed),
		llb.KeepGitDir(),
	}
	gitState := pllb.Git(gitURL, branch, gitOpts...)
	c.mts.Final.MainState = llbutil.CopyOp(
		gitState, []string{"."}, c.mts.Final.MainState, dest, false, false, keepTs,
		c.mts.Final.MainImage.Config.User, nil, false, false, c.ftrs.UseCopyLink,
		llb.WithCustomNamef(
			"%sCOPY GIT CLONE (--branch %s) %s TO %s", c.vertexPrefix(false, false, false),
			branch, gitURLScrubbed, dest))
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
func (c *Converter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int) error {
	err := c.checkAllowed(healthcheckCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	hc := &dockerfile2llb.HealthConfig{}
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
func (c *Converter) Cache(ctx context.Context, mountTarget string) error {
	err := c.checkAllowed(cacheCmd)
	if err != nil {
		return err
	}
	c.nonSaveCommand()

	key, err := cacheKeyTargetInput(c.targetInputActiveOnly())
	if err != nil {
		return err
	}
	mountID := path.Clean(mountTarget)
	cachePath := path.Join("/run/cache", key, mountID)

	if _, exists := c.persistentCacheDirs[mountTarget]; !exists {
		c.persistentCacheDirs[mountTarget] = pllb.AddMount(mountTarget, pllb.Scratch(),
			llb.AsPersistentCacheDir(cachePath, llb.CacheMountShared))
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

// EnterScopeDo introduces a new variable scope. Gloabls and imports are fetched from baseTarget.
func (c *Converter) EnterScopeDo(ctx context.Context, command domain.Command, baseTarget domain.Target, allowPrivileged bool, scopeName string, buildArgs []string) error {
	baseMts, err := c.buildTarget(ctx, baseTarget.String(), c.platr.Current(), allowPrivileged, buildArgs, true, enterScopeDoCmd)
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

	if c.ftrs.WaitBlock {
		c.waitBlock().addState(&c.mts.Final.MainState, c)
	}

	close(c.mts.Final.Done())
	return c.mts, nil
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

func (c *Converter) prepBuildTarget(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDangling bool, cmdT cmdType) (domain.Target, ConvertOpt, bool, error) {
	relTarget, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, errors.Wrapf(err, "earthly target parse %s", fullTargetName)
	}
	derefedTarget, allowPrivilegedImport, isImport, err := c.varCollection.Imports().Deref(relTarget)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, err
	}
	if isImport {
		allowPrivileged = allowPrivileged && allowPrivilegedImport
	}
	targetRef, err := c.joinRefs(derefedTarget)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, errors.Wrap(err, "join targets")
	}
	target := targetRef.(domain.Target)

	var pncvf variables.ProcessNonConstantVariableFunc
	if !c.opt.Features.ShellOutAnywhere {
		pncvf = c.processNonConstantBuildArgFunc(ctx)
	}

	overriding, err := variables.ParseArgs(buildArgs, pncvf, c.varCollection)
	if err != nil {
		return domain.Target{}, ConvertOpt{}, false, errors.Wrap(err, "parse build args")
	}
	// Don't allow transitive overriding variables to cross project boundaries.
	propagateBuildArgs := !relTarget.IsExternal()
	if propagateBuildArgs {
		overriding = variables.CombineScopes(overriding, c.varCollection.Overriding())
	}

	// Recursion.
	opt := c.opt
	opt.OverridingVars = overriding
	opt.GlobalImports = nil
	opt.parentDepSub = c.mts.Final.NewDependencySubscription()
	opt.PlatformResolver = c.platr.SubResolver(platform)
	opt.HasDangling = isDangling
	opt.AllowPrivileged = allowPrivileged
	opt.waitBlock = c.waitBlock()
	if c.opt.Features.ReferencedSaveOnly {
		// DoSaves should only be potentially turned-off when the ReferencedSaveOnly feature is flipped
		opt.DoSaves = (cmdT == buildCmd && c.opt.DoSaves)
		opt.ForceSaveImage = false
	} else {
		opt.DoSaves = c.opt.DoSaves && !target.IsRemote() // legacy mode only saves artifacts from local targets
	}
	return target, opt, propagateBuildArgs, nil
}

func (c *Converter) buildTarget(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDangling bool, cmdT cmdType) (*states.MultiTarget, error) {
	target, opt, propagateBuildArgs, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, isDangling, cmdT)
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
			_, found := opt.OverridingVars.GetAny(bai.Name)
			if found {
				continue
			}
			c.mts.Final.AddBuildArgInput(bai)
		}
		if cmdT == fromCmd {
			// Propagate globals.
			globals := mts.Final.VarCollection.Globals()
			for _, k := range globals.SortedActive() {
				_, alreadyActive := c.varCollection.GetActive(k)
				if alreadyActive {
					// Globals don't override any variables in current scope.
					continue
				}
				v, _ := globals.GetActive(k)
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
		}
	}

	return mts, nil
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
	}
	if opts.shellWrap == nil {
		opts.shellWrap = withShellAndEnvVars
	}

	if c.ftrs.WaitBlock && opts.Push {
		return pllb.State{}, errors.New("RUN --push is not currently supported with --wait-block, you must rewrite it as IF ... RUN --no-cache END") // TODO this will be done automatically
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
	mountRunOpts, err := parseMounts(opts.Mounts, c.mts.Final.Target, c.targetInputActiveOnly(), c.cacheContext, c.platr)
	if err != nil {
		return pllb.State{}, errors.Wrap(err, "parse mounts")
	}
	runOpts = append(runOpts, mountRunOpts...)
	commandStr := fmt.Sprintf(
		"%s %s%s%s%s%s%s",
		opts.CommandName, // e.g. "RUN", "IF", "FOR", "ARG"
		strIf(opts.Privileged, "--privileged "),
		strIf(opts.Push, "--push "),
		strIf(opts.NoCache, "--no-cache "),
		strIf(opts.Interactive, "--interactive "),
		strIf(opts.InteractiveKeep, "--interactive-keep "),
		strings.Join(opts.Args, " "))
	runOpts = append(runOpts, llb.WithCustomNamef("%s%s", c.vertexPrefix(opts.Locally, isInteractive, false), commandStr))

	var extraEnvVars []string
	// Secrets.
	for _, secretKeyValue := range opts.Secrets {
		secretID, envVar, err := c.parseSecretFlag(secretKeyValue)
		if err != nil {
			return pllb.State{}, err
		}
		if secretID != "" {
			secretPath := path.Join("/run/secrets", secretID)
			secretOpts := []llb.SecretOption{
				llb.SecretID(secretID),
				// TODO: Perhaps this should just default to the current user automatically from
				//       buildkit side. Then we wouldn't need to open this up to everyone.
				llb.SecretFileOpt(0, 0, 0444),
			}
			runOpts = append(runOpts, llb.AddSecret(secretPath, secretOpts...))
			// TODO: The use of cat here might not be portable.
			extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"$(cat %s)\"", envVar, secretPath))
		}
	}
	// Build args.
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _ := c.varCollection.GetActive(buildArgName)
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=%s", buildArgName, shellescape.Quote(ba)))
	}
	if !opts.Locally {
		// Debugger.
		secretOpts := []llb.SecretOption{
			llb.SecretID(common.DebuggerSettingsSecretsKey),
			llb.SecretFileOpt(0, 0, 0444),
		}
		debuggerSecretMount := llb.AddSecret(
			fmt.Sprintf("/run/secrets/%s", common.DebuggerSettingsSecretsKey), secretOpts...)
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

func (c *Converter) parseSecretFlag(secretKeyValue string) (secretID string, envVar string, err error) {
	parts := strings.SplitN(secretKeyValue, "=", 2)
	if len(parts) == 2 {
		if strings.HasPrefix(parts[1], "+secrets/") {
			secretID := strings.TrimPrefix(parts[1], "+secrets/")
			return secretID, parts[0], nil
		} else if parts[1] == "" {
			// If empty string, don't use (used for optional secrets).
			// TODO: This should be an actual secret (with an empty value),
			//       so that the cache works correctly.
			return "", "", nil
		} else {
			return "", "", errors.Errorf("secret definition %s not supported. Format must be either <env-var>=+secrets/<secret-id> or <secret-id>", secretKeyValue)
		}
	} else if len(parts) == 1 {
		if secretKeyValue == "" {
			// If empty string, don't use (used for optional secrets).
			// TODO: This should be an actual secret (with an empty value),
			//       so that the cache works correctly.
			return "", "", nil
		}
		return parts[0], parts[0], nil
	} else {
		return "", "", errors.Errorf("secret definition %s not supported. Format must be either <env-var>=+secrets/<secret-id> or <secret-id>", secretKeyValue)
	}
}

func (c *Converter) forceExecutionWithSemaphore(ctx context.Context, state pllb.State, platr *platutil.Resolver) error {
	sem := c.opt.Parallelism
	rel, err := sem.Acquire(ctx, 1)
	if err != nil {
		return errors.Wrapf(err, "acquiring parallelism semaphore during forceExecutionWithSemaphore for %s", c.target.String())
	}
	defer rel()
	return c.forceExecution(ctx, state, platr)
}

func (c *Converter) forceExecution(ctx context.Context, state pllb.State, platr *platutil.Resolver) error {
	if state.Output() == nil {
		// Scratch - no need to execute.
		return nil
	}
	ref, err := llbutil.StateToRef(
		ctx, c.opt.GwClient, state, c.opt.NoCache,
		platr, c.opt.CacheImports.AsMap())
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
		mts.Final.PlatformResolver, c.opt.CacheImports.AsMap())
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
	ref, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", imageName)
	}
	baseImageName := reference.TagNameOnly(ref).String()
	logName := fmt.Sprintf(
		"%sLoad metadata %s",
		c.imageVertexPrefix(imageName, platform), platforms.Format(llbPlatform))
	dgst, dt, err := c.opt.MetaResolver.ResolveImageConfig(
		ctx, baseImageName,
		llb.ResolveImageConfigOpt{
			Platform:    &llbPlatform,
			ResolveMode: c.opt.ImageResolveMode.String(),
			LogName:     logName,
		})
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "resolve image config for %s", imageName)
	}
	var img image.Image
	err = json.Unmarshal(dt, &img)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "unmarshal image config for %s", imageName)
	}
	if dgst != "" {
		ref, err = reference.WithDigest(ref, dgst)
		if err != nil {
			return pllb.State{}, nil, nil, errors.Wrapf(err, "reference add digest %v for %s", dgst, imageName)
		}
	}
	allOpts := append(opts, llb.Platform(c.platr.ToLLBPlatform(platform)), c.opt.ImageResolveMode)
	state := pllb.Image(ref.String(), allOpts...)
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
	for _, name := range ev.SortedActive() {
		v, _ := ev.GetActive(name)
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

func (c *Converter) vertexPrefix(local bool, interactive bool, internal bool) string {
	activeOverriding := make(map[string]string)
	for _, arg := range c.varCollection.SortedOverridingVariables() {
		v, ok := c.varCollection.GetActive(arg)
		if ok {
			activeOverriding[arg] = v
		}
	}
	platform := c.platr.Materialize(c.platr.Current())
	platformStr := platform.String()
	isNativePlatform := c.platr.PlatformEquals(platform, platutil.NativePlatform)
	vm := &outmon.VertexMeta{
		TargetID:           c.mts.Final.ID,
		TargetName:         c.mts.Final.Target.String(),
		Platform:           platformStr,
		NonDefaultPlatform: !isNativePlatform,
		Local:              local,
		Interactive:        interactive,
		OverridingArgs:     activeOverriding,
		Internal:           internal,
	}
	return vm.ToVertexPrefix()
}

func (c *Converter) imageVertexPrefix(id string, platform platutil.Platform) string {
	platform = c.platr.Materialize(platform)
	isNativePlatform := c.platr.PlatformEquals(platform, platutil.NativePlatform)
	vm := &outmon.VertexMeta{
		TargetName:         id,
		Platform:           platform.String(),
		NonDefaultPlatform: !isNativePlatform,
	}
	return vm.ToVertexPrefix()
}

func (c *Converter) vertexPrefixWithURL(url string) string {
	return fmt.Sprintf("[%s(%s) %s] ", c.mts.Final.Target.String(), url, url)
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
	if c.mts.Final.RanInteractive && command != saveImageCmd {
		return errors.New("If present, a single --interactive command must be the last command in a target")
	}

	if c.mts.Final.RanFromLike {
		return nil
	}

	switch command {
	case fromCmd, fromDockerfileCmd, locallyCmd, buildCmd, argCmd, importCmd:
		return nil
	default:
		return errors.New("the first command has to be FROM, FROM DOCKERFILE, LOCALLY, ARG, BUILD or IMPORT")
	}
}

func (c *Converter) targetInputActiveOnly() dedup.TargetInput {
	activeBuildArgs := make(map[string]bool)
	for _, k := range c.varCollection.SortedActiveVariables() {
		activeBuildArgs[k] = true
	}
	return c.mts.Final.TargetInput().WithFilterBuildArgs(activeBuildArgs)
}

// persistCache makes temporary cache directories permanent by writing their contents
// from the cached directory to the persistent image layers at the same directory.
// This only has an effect when the Target contains at least one `CACHE /my/directory` command.
func (c *Converter) persistCache(srcState pllb.State) pllb.State {
	dest := srcState
	// User may have multiple CACHE commands in a single target
	for dir, cache := range c.persistentCacheDirs {
		// Copy the contents of the user's cache directory to the temporary backup.
		// It's important to use DockerfileCopy here, since traditional llb.Copy()
		// doesn't support adding mounts via RunOptions.
		runOpts := []llb.RunOption{cache, llb.WithCustomName("persist cache directory")}
		dest = llbutil.CopyWithRunOptions(
			dest,
			dir, // cache dir from external mount
			dir, // cache dir on dest state (same location but without the mount)
			c.platr,
			runOpts...,
		)
	}

	return dest
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
