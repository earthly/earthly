package earthfile2llb

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/variables"

	"github.com/alessio/shellescape"
	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session/localhost"
	solverpb "github.com/moby/buildkit/solver/pb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

const exitCodeFile = "/run/exit_code"

// Converter turns earthly commands to buildkit LLB representation.
type Converter struct {
	gitMeta             *gitutil.GitMetadata
	opt                 ConvertOpt
	mts                 *states.MultiTarget
	directDeps          []*states.SingleTarget
	buildContextFactory llbfactory.Factory
	cacheContext        pllb.State
	varCollection       *variables.Collection
	ranSave             bool
	cmdSet              bool
	ftrs                *features.Features
}

// NewConverter constructs a new converter for a given earthly target.
func NewConverter(ctx context.Context, target domain.Target, bc *buildcontext.Data, sts *states.SingleTarget, opt ConvertOpt, ftrs *features.Features) (*Converter, error) {
	for k, v := range bc.LocalDirs {
		sts.LocalDirs[k] = v
	}
	sts.HasDangling = opt.HasDangling
	mts := &states.MultiTarget{
		Final:   sts,
		Visited: opt.Visited,
	}
	sts.AddOverridingVarsAsBuildArgInputs(opt.OverridingVars)
	vc := variables.NewCollection(opt.Console,
		target, llbutil.PlatformWithDefault(opt.Platform), bc.GitMetadata, opt.OverridingVars,
		opt.GlobalImports)
	return &Converter{
		gitMeta:             bc.GitMetadata,
		opt:                 opt,
		mts:                 mts,
		buildContextFactory: bc.BuildContextFactory,
		cacheContext:        pllb.Scratch(),
		varCollection:       vc,
		ftrs:                ftrs,
	}, nil
}

// From applies the earthly FROM command.
func (c *Converter) From(ctx context.Context, imageName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string) error {
	err := c.checkAllowed("FROM")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	c.cmdSet = false
	platform, err = llbutil.ResolvePlatform(platform, c.opt.Platform)
	if err != nil {
		return err
	}
	c.setPlatform(platform)
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

func (c *Converter) fromClassical(ctx context.Context, imageName string, platform *specs.Platform, local bool) error {
	var prefix string
	if local {
		// local mode uses a fake image containing /bin/true
		// we want to prefix this as internal so it doesn't show up in the output
		prefix = "[internal] "
	} else {
		prefix = c.vertexPrefix(false, false)
	}
	plat := llbutil.PlatformWithDefault(platform)
	state, img, envVars, err := c.internalFromClassical(
		ctx, imageName, plat,
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

func (c *Converter) fromTarget(ctx context.Context, targetName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string) error {
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}
	mts, err := c.buildTarget(ctx, depTarget.String(), platform, allowPrivileged, buildArgs, false, "FROM")
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
	for dirKey, dirValue := range relevantDepState.LocalDirs {
		c.mts.Final.LocalDirs[dirKey] = dirValue
	}
	c.varCollection.ResetEnvVars(mts.Final.VarCollection.EnvVars())
	c.mts.Final.MainImage = saveImage.Image.Clone()
	c.mts.Final.RanFromLike = mts.Final.RanFromLike
	c.mts.Final.RanInteractive = mts.Final.RanInteractive
	c.setPlatform(mts.Final.Platform)
	return nil
}

// FromDockerfile applies the earthly FROM DOCKERFILE command.
func (c *Converter) FromDockerfile(ctx context.Context, contextPath string, dfPath string, dfTarget string, platform *specs.Platform, buildArgs []string) error {
	err := c.checkAllowed("FROM DOCKERFILE")
	if err != nil {
		return err
	}
	platform, err = llbutil.ResolvePlatform(platform, c.opt.Platform)
	if err != nil {
		return err
	}
	c.setPlatform(platform)
	plat := llbutil.PlatformWithDefault(platform)
	c.nonSaveCommand()
	var dfData []byte
	if dfPath != "" {
		dfArtifact, parseErr := domain.ParseArtifact(dfPath)
		if parseErr == nil {
			// The Dockerfile is from a target's artifact.
			mts, err := c.buildTarget(ctx, dfArtifact.Target.String(), platform, false, buildArgs, false, "FROM DOCKERFILE")
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
			data, err := c.opt.Resolver.Resolve(ctx, c.opt.GwClient, dockerfileMetaTarget)
			if err != nil {
				return errors.Wrap(err, "resolve build context for dockerfile")
			}
			for ldk, ld := range data.LocalDirs {
				c.mts.Final.LocalDirs[ldk] = ld
			}
			dfData, err = ioutil.ReadFile(data.BuildFilePath)
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
		mts, err := c.buildTarget(ctx, contextArtifact.Target.String(), platform, false, buildArgs, false, "FROM DOCKERFILE")
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
			llbutil.ScratchWithPlatform(), "/", true, true, false, "", false, false,
			llb.WithCustomNamef(
				"[internal] FROM DOCKERFILE (copy build context from) %s%s",
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
			Target:    fmt.Sprintf("%sDockerfile", buildcontext.DockerfileMetaTarget),
			LocalPath: path.Join(contextPath),
		}
		dockerfileMetaTargetRef, err := c.joinRefs(dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "join targets")
		}
		dockerfileMetaTarget = dockerfileMetaTargetRef.(domain.Target)
		data, err := c.opt.Resolver.Resolve(ctx, c.opt.GwClient, dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "resolve build context for dockerfile")
		}
		for ldk, ld := range data.LocalDirs {
			c.mts.Final.LocalDirs[ldk] = ld
		}
		if dfPath == "" {
			// Imply dockerfile as being ./Dockerfile in the root of the build context.
			dfData, err = ioutil.ReadFile(data.BuildFilePath)
			if err != nil {
				return errors.Wrapf(err, "read file %s", data.BuildFilePath)
			}
		}
		BuildContextFactory = data.BuildContextFactory
	}
	overriding, err := variables.ParseArgs(
		buildArgs, c.processNonConstantBuildArgFunc(ctx), c.varCollection)
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
func (c *Converter) Locally(ctx context.Context, workdirPath string, platform *specs.Platform) error {
	err := c.checkAllowed("LOCALLY")
	if err != nil {
		return err
	}
	if !c.opt.AllowLocally {
		return errors.New("LOCALLY cannot be used when --strict is specified or otherwise implied")
	}
	if !path.IsAbs(workdirPath) {
		return errors.New("workdirPath must be absolute")
	}

	err = c.fromClassical(ctx, "scratch", platform, true)
	if err != nil {
		return err
	}

	// reset WORKDIR to current directory where Earthfile is
	c.mts.Final.MainState = c.mts.Final.MainState.Dir(workdirPath)
	c.mts.Final.MainImage.Config.WorkingDir = workdirPath
	return nil
}

// CopyArtifactLocal applies the earthly COPY artifact command which are invoked under a LOCALLY target.
func (c *Converter) CopyArtifactLocal(ctx context.Context, artifactName string, dest string, platform *specs.Platform, allowPrivileged bool, buildArgs []string, isDir bool) error {
	err := c.checkAllowed("COPY")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, buildArgs, false, "COPY")
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
			c.vertexPrefix(false, false),
			strIf(isDir, "--dir "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest),
	}

	c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()
	return nil
}

// CopyArtifact applies the earthly COPY artifact command.
func (c *Converter) CopyArtifact(ctx context.Context, artifactName string, dest string, platform *specs.Platform, allowPrivileged bool, buildArgs []string, isDir bool, keepTs bool, keepOwn bool, chown string, ifExists, symlinkNoFollow bool) error {
	err := c.checkAllowed("COPY")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, allowPrivileged, buildArgs, false, "COPY")
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
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), ifExists, symlinkNoFollow,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s%s%s %s",
			c.vertexPrefix(false, false),
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest))
	return nil
}

// CopyClassical applies the earthly COPY command, with classical args.
func (c *Converter) CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool, keepTs bool, keepOwn bool, chown string) error {
	err := c.checkAllowed("COPY")
	if err != nil {
		return err
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
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), false, false,
		llb.WithCustomNamef(
			"%sCOPY %s%s %s",
			c.vertexPrefix(false, false),
			strIf(isDir, "--dir "),
			strings.Join(srcs, " "),
			dest))
	return nil
}

// RunLocal applies a RUN statement locally rather than in a container
func (c *Converter) RunLocal(ctx context.Context, args []string, pushFlag bool) error {
	err := c.checkAllowed("RUN")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	runStr := fmt.Sprintf("RUN %s%s", strIf(pushFlag, "--push "), strings.Join(args, " "))

	// Build args get propagated into env.
	extraEnvVars := []string{}
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _ := c.varCollection.GetActive(buildArgName)
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"%s\"", buildArgName, shellescape.Quote(ba)))
	}

	// buildkit-hack in order to run locally, we prepend the command with a UUID
	finalArgs := append([]string{localhost.RunOnLocalHostMagicStr}, withShellAndEnvVars(args, extraEnvVars, true, false, false)...)
	opts := []llb.RunOption{
		llb.Args(finalArgs),
		llb.IgnoreCache,
		llb.WithCustomNamef("%s%s", c.vertexPrefix(true, false), runStr),
	}

	if pushFlag {
		if !c.mts.Final.RunPush.HasState {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.Final.RunPush.State = c.mts.Final.MainState
			c.mts.Final.RunPush.HasState = true
		}
		c.mts.Final.RunPush.State = c.mts.Final.RunPush.State.Run(opts...).Root()
		c.mts.Final.RunPush.CommandStrs = append(
			c.mts.Final.RunPush.CommandStrs, runStr)
	} else {
		c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()
	}
	return nil
}

// RunExitCode executes a run for the purpose of determining the exit code of the command. This can be used in conditionals.
func (c *Converter) RunExitCode(ctx context.Context, commandName string, args, mounts, secretKeyValues []string, privileged, isWithShell, withSSH, noCache bool) (int, error) {
	err := c.checkAllowed("RUN")
	if err != nil {
		return 0, err
	}
	c.nonSaveCommand()
	if !isWithShell {
		return 0, errors.New("non-shell mode not yet supported")
	}

	// The exit code will be placed in /run. We need that dir to have been created.
	c.mts.Final.MainState = c.mts.Final.MainState.File(
		pllb.Mkdir("/run", 0755, llb.WithParents(true)),
		llb.WithCustomNamef("[internal] mkdir %s", "/run"))

	// Perform execution, but append the command with the right shell incantation that
	// causes it to output the exit code to a file. This is done via the shellWrap.
	var opts []llb.RunOption
	mountRunOpts, err := parseMounts(mounts, c.mts.Final.Target, c.targetInputActiveOnly(), c.cacheContext)
	if err != nil {
		return 0, errors.Wrap(err, "parse mounts")
	}
	opts = append(opts, mountRunOpts...)
	if privileged {
		opts = append(opts, llb.Security(llb.SecurityModeInsecure))
	}
	runStr := fmt.Sprintf(
		"%s %s%s%s",
		commandName, // eg IF
		strIf(privileged, "--privileged "),
		strIf(noCache, "--no-cache "),
		strings.Join(args, " "))
	shellWrap := withShellAndEnvVarsExitCode(exitCodeFile)
	opts = append(opts, llb.WithCustomNamef("%s%s", c.vertexPrefix(false, false), runStr))
	state, err := c.internalRun(
		ctx, args, secretKeyValues, isWithShell, shellWrap, false,
		true, withSSH, noCache, false, false, runStr, opts...)
	if err != nil {
		return 0, err
	}
	ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, state, c.opt.Platform, c.opt.CacheImports.AsMap())
	if err != nil {
		return 0, errors.Wrap(err, "run exit code state to ref")
	}
	codeDt, err := ref.ReadFile(ctx, gwclient.ReadRequest{
		Filename: exitCodeFile,
	})
	if err != nil {
		return 0, errors.Wrap(err, "read exit code")
	}
	exitCode, err := strconv.ParseInt(string(bytes.TrimSpace(codeDt)), 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse exit code as int")
	}
	return int(exitCode), err
}

// RunLocalExitCode runs a command locally rather than in a container and returns its exit code.
func (c *Converter) RunLocalExitCode(ctx context.Context, commandName string, args []string) (int, error) {
	err := c.checkAllowed("RUN")
	if err != nil {
		return 0, err
	}
	c.nonSaveCommand()
	runStr := fmt.Sprintf("%s %s", commandName, strings.Join(args, " "))
	// Build args get propagated into env.
	extraEnvVars := []string{}
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _ := c.varCollection.GetActive(buildArgName)
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"%s\"", buildArgName, shellescape.Quote(ba)))
	}

	exitCodeDir, err := ioutil.TempDir(os.TempDir(), "earthlyexitcode")
	if err != nil {
		return 0, errors.Wrap(err, "create temp dir")
	}
	exitCodeFile := filepath.Join(exitCodeDir, "/exit_code")
	c.opt.CleanCollection.Add(func() error {
		return os.RemoveAll(exitCodeDir)
	})

	// buildkit-hack in order to run locally, we prepend the command with a UUID
	finalArgs := append(
		[]string{localhost.RunOnLocalHostMagicStr},
		withShellAndEnvVarsExitCode(exitCodeFile)(args, extraEnvVars, true, false, false)...,
	)
	opts := []llb.RunOption{
		llb.Args(finalArgs),
		llb.IgnoreCache,
		llb.WithCustomNamef("%s%s", c.vertexPrefix(true, false), runStr),
	}
	c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()

	ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, c.mts.Final.MainState, c.opt.Platform, c.opt.CacheImports.AsMap())
	if err != nil {
		return 0, errors.Wrap(err, "run exit code state to ref")
	}
	// Cause the execution to complete. We're not really interested in reading the dir - we just
	// want to un-lazy the ref so that the local commands have executed.
	_, err = ref.ReadDir(ctx, gwclient.ReadDirRequest{Path: "/"})
	if err != nil {
		return 0, errors.Wrap(err, "unlazy locally")
	}

	codeDt, err := ioutil.ReadFile(exitCodeFile)
	if err != nil {
		return 0, errors.Wrap(err, "read exit code file")
	}
	exitCode, err := strconv.ParseInt(string(bytes.TrimSpace(codeDt)), 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "parse exit code as int")
	}
	return int(exitCode), err
}

// Run applies the earthly RUN command.
func (c *Converter) Run(ctx context.Context, args, mounts, secretKeyValues []string, privileged, withEntrypoint, withDocker, isWithShell, pushFlag, withSSH, noCache, interactive, interactiveKeep bool) error {
	err := c.checkAllowed("RUN")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	if withDocker {
		return errors.New("RUN --with-docker is obsolete. Please use WITH DOCKER ... RUN ... END instead")
	}

	var opts []llb.RunOption
	mountRunOpts, err := parseMounts(mounts, c.mts.Final.Target, c.targetInputActiveOnly(), c.cacheContext)
	if err != nil {
		return errors.Wrap(err, "parse mounts")
	}
	opts = append(opts, mountRunOpts...)

	finalArgs := args
	if withEntrypoint {
		if len(args) == 0 {
			// No args provided. Use the image's CMD.
			args := make([]string, len(c.mts.Final.MainImage.Config.Cmd))
			copy(args, c.mts.Final.MainImage.Config.Cmd)
		}
		finalArgs = append(c.mts.Final.MainImage.Config.Entrypoint, args...)
		isWithShell = false // Don't use shell when --entrypoint is passed.
	}
	if privileged {
		opts = append(opts, llb.Security(llb.SecurityModeInsecure))
	}
	runStr := fmt.Sprintf(
		"RUN %s%s%s%s%s%s%s%s",
		strIf(privileged, "--privileged "),
		strIf(withDocker, "--with-docker "),
		strIf(withEntrypoint, "--entrypoint "),
		strIf(pushFlag, "--push "),
		strIf(noCache, "--no-cache "),
		strIf(interactive, "--interactive "),
		strIf(interactiveKeep, "--interactive-keep "),
		strings.Join(finalArgs, " "))
	shellWrap := withShellAndEnvVars
	opts = append(opts, llb.WithCustomNamef("%s%s", c.vertexPrefix(false, interactive || interactiveKeep), runStr))
	_, err = c.internalRun(
		ctx, finalArgs, secretKeyValues, isWithShell, shellWrap, pushFlag,
		false, withSSH, noCache, interactive, interactiveKeep, runStr, opts...)
	return err
}

// SaveArtifact applies the earthly SAVE ARTIFACT command.
func (c *Converter) SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string, keepTs bool, keepOwn bool, ifExists, symlinkNoFollow bool, isPush bool) error {
	err := c.checkAllowed("SAVE ARTIFACT")
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
	c.mts.Final.ArtifactsState = llbutil.CopyOp(
		c.mts.Final.MainState, []string{saveFrom}, c.mts.Final.ArtifactsState,
		saveToAdjusted, true, true, keepTs, own, ifExists, symlinkNoFollow,
		llb.WithCustomNamef(
			"%sSAVE ARTIFACT %s%s%s %s",
			c.vertexPrefix(false, false),
			strIf(ifExists, "--if-exists "),
			strIf(symlinkNoFollow, "--symlink-no-follow "),
			saveFrom,
			artifact.String()))
	if saveAsLocalTo != "" && c.opt.DoSaves {
		separateArtifactsState := llbutil.ScratchWithPlatform()
		if isPush {
			separateArtifactsState = llbutil.CopyOp(
				c.mts.Final.RunPush.State, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", ifExists, symlinkNoFollow,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					c.vertexPrefix(false, false),
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
		} else {
			separateArtifactsState = llbutil.CopyOp(
				c.mts.Final.MainState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", ifExists, symlinkNoFollow,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s%s %s AS LOCAL %s",
					c.vertexPrefix(false, false),
					strIf(ifExists, "--if-exists "),
					strIf(symlinkNoFollow, "--symlink-no-follow "),
					saveFrom,
					artifact.String(),
					saveAsLocalTo))
		}
		c.mts.Final.SeparateArtifactsState = append(c.mts.Final.SeparateArtifactsState, separateArtifactsState)

		saveLocal := states.SaveLocal{
			DestPath:     saveAsLocalTo,
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

// SaveArtifactFromLocal saves a local file into the ArtifactsState
func (c *Converter) SaveArtifactFromLocal(ctx context.Context, saveFrom, saveTo string, keepTs, keepOwn bool, chown string) error {
	err := c.checkAllowed("SAVE ARTIFACT")
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
		llb.WithCustomNamef("[internal] CopyFileMagicStr %s %s", saveFrom, saveTo),
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
		absSaveTo, true, true, keepTs, own, ifExists, false,
	)
	c.ranSave = true
	c.markFakeDeps()
	return nil
}

// SaveImage applies the earthly SAVE IMAGE command.
func (c *Converter) SaveImage(ctx context.Context, imageNames []string, pushImages bool, insecurePush bool, cacheHint bool, cacheFrom []string) error {
	err := c.checkAllowed("SAVE IMAGE")
	if err != nil {
		return err
	}
	for _, cf := range cacheFrom {
		c.opt.CacheImports.Add(cf)
	}
	if !c.opt.DoSaves {
		imageNames = []string{}
	}
	justCacheHint := false
	if len(imageNames) == 0 && cacheHint {
		imageNames = []string{""}
		justCacheHint = true
	}
	for _, imageName := range imageNames {
		if c.mts.Final.RunPush.HasState {
			// SAVE IMAGE --push when it comes before any RUN --push should be treated as if they are in the main state,
			// since thats their only dependency. It will still be marked as a push.
			c.mts.Final.RunPush.SaveImages = append(c.mts.Final.RunPush.SaveImages,
				states.SaveImage{
					State:               c.mts.Final.RunPush.State,
					Image:               c.mts.Final.MainImage.Clone(), // We can get away with this because no Image details can vary in a --push. This should be fixed before then.
					DockerTag:           imageName,
					Push:                pushImages,
					InsecurePush:        insecurePush,
					CacheHint:           cacheHint,
					HasPushDependencies: true,
				})
		} else {
			c.mts.Final.SaveImages = append(c.mts.Final.SaveImages,
				states.SaveImage{
					State:               c.mts.Final.MainState,
					Image:               c.mts.Final.MainImage.Clone(),
					DockerTag:           imageName,
					Push:                pushImages,
					InsecurePush:        insecurePush,
					CacheHint:           cacheHint,
					HasPushDependencies: false,
				})
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
func (c *Converter) Build(ctx context.Context, fullTargetName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string) error {
	err := c.checkAllowed("BUILD")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	_, err = c.buildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, true, "BUILD")
	return err
}

// BuildAsync applies the earthly BUILD command asynchronously.
func (c *Converter) BuildAsync(ctx context.Context, fullTargetName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string, cmdName string) chan error {
	errChan := make(chan error, 1)
	target, opt, _, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, false, cmdName)
	if err != nil {
		errChan <- err
		return errChan
	}
	go func() {
		err := c.opt.Parallelism.Acquire(ctx, 1)
		if err != nil {
			errChan <- errors.Wrapf(err, "acquiring parallelism semaphore for %s", fullTargetName)
			return
		}
		defer c.opt.Parallelism.Release(1)
		_, err = Earthfile2LLB(ctx, target, opt)
		if err != nil {
			errChan <- errors.Wrapf(err, "async earthfile2llb for %s", fullTargetName)
			return
		}
		errChan <- nil
	}()
	return errChan
}

// Workdir applies the WORKDIR command.
func (c *Converter) Workdir(ctx context.Context, workdirPath string) error {
	err := c.checkAllowed("WORKDIR")
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
			llb.WithCustomNamef("%sWORKDIR %s", c.vertexPrefix(false, false), workdirPath),
		}
		c.mts.Final.MainState = c.mts.Final.MainState.File(
			pllb.Mkdir(workdirAbs, 0755, mkdirOpts...), opts...)
	}
	return nil
}

// User applies the USER command.
func (c *Converter) User(ctx context.Context, user string) error {
	err := c.checkAllowed("USER")
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
	err := c.checkAllowed("CMD")
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
	err := c.checkAllowed("ENTRYPOINT")
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
	err := c.checkAllowed("EXPOSE")
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
	err := c.checkAllowed("VOLUME")
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
	err := c.checkAllowed("ENV")
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
func (c *Converter) Arg(ctx context.Context, argKey string, defaultArgValue string, global bool) error {
	err := c.checkAllowed("ARG")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	effective, err := c.varCollection.DeclareArg(argKey, defaultArgValue, global, c.processNonConstantBuildArgFunc((ctx)))
	if err != nil {
		return err
	}
	c.mts.Final.AddBuildArgInput(dedup.BuildArgInput{
		Name:          argKey,
		DefaultValue:  defaultArgValue,
		ConstantValue: effective,
	})
	return nil
}

// Label applies the LABEL command.
func (c *Converter) Label(ctx context.Context, labels map[string]string) error {
	err := c.checkAllowed("LABEL")
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
	err := c.checkAllowed("GIT CLONE")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef(
			"%sGIT CLONE (--branch %s) %s", c.vertexPrefixWithURL(gitURL), branch, gitURL),
		llb.KeepGitDir(),
	}
	gitState := pllb.Git(gitURL, branch, gitOpts...)
	c.mts.Final.MainState = llbutil.CopyOp(
		gitState, []string{"."}, c.mts.Final.MainState, dest, false, false, keepTs,
		c.mts.Final.MainImage.Config.User, false, false,
		llb.WithCustomNamef(
			"%sCOPY GIT CLONE (--branch %s) %s TO %s", c.vertexPrefix(false, false),
			branch, gitURL, dest))
	return nil
}

// WithDockerRun applies an entire WITH DOCKER ... RUN ... END clause.
func (c *Converter) WithDockerRun(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := c.checkAllowed("RUN")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	wdr := &withDockerRun{
		c: c,
	}
	return wdr.Run(ctx, args, opt)
}

// WithDockerRunLocal applies an entire WITH DOCKER ... RUN ... END clause.
func (c *Converter) WithDockerRunLocal(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := c.checkAllowed("RUN")
	if err != nil {
		return err
	}
	c.nonSaveCommand()
	wdrl := &withDockerRunLocal{
		c: c,
	}
	return wdrl.Run(ctx, args, opt)
}

// Healthcheck applies the HEALTHCHECK command.
func (c *Converter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int) error {
	err := c.checkAllowed("HEALTHCHECK")
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
	err := c.checkAllowed("IMPORT")
	if err != nil {
		return err
	}
	return c.varCollection.Imports().Add(importStr, as, isGlobal, currentlyPrivileged, allowPrivilegedFlag)
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
	bc, err = c.opt.Resolver.Resolve(ctx, c.opt.GwClient, refToResolve)
	if err != nil {
		return nil, false, false, err
	}
	return bc, allowPrivileged, allowPrivilegedSet, nil
}

// EnterScope introduces a new variable scope. Gloabls and imports are fetched from baseTarget.
func (c *Converter) EnterScope(ctx context.Context, command domain.Command, baseTarget domain.Target, allowPrivileged bool, scopeName string, buildArgs []string) error {
	baseMts, err := c.buildTarget(ctx, baseTarget.String(), c.mts.Final.Platform, allowPrivileged, buildArgs, true, "ENTER-SCOPE")
	if err != nil {
		return err
	}

	overriding, err := variables.ParseArgs(
		buildArgs, c.processNonConstantBuildArgFunc(ctx), c.varCollection)
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
	c.opt.BuildContextProvider.AddDirs(c.mts.Final.LocalDirs)
	c.mts.Final.VarCollection = c.varCollection
	c.mts.Final.GlobalImports = c.varCollection.Imports().Global()
	close(c.mts.Final.Done())
	return c.mts, nil
}

// ExpandArgs expands args in the provided word.
func (c *Converter) ExpandArgs(word string) string {
	return c.varCollection.Expand(word)
}

func (c *Converter) prepBuildTarget(ctx context.Context, fullTargetName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string, isDangling bool, cmdName string) (domain.Target, ConvertOpt, bool, error) {
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

	overriding, err := variables.ParseArgs(buildArgs, c.processNonConstantBuildArgFunc(ctx), c.varCollection)
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
	opt.Platform, err = llbutil.ResolvePlatform(platform, c.opt.Platform)
	opt.HasDangling = isDangling
	opt.AllowPrivileged = allowPrivileged
	if c.opt.Features.ReferencedSaveOnly {
		// DoSaves should only be potentially turned-off when the ReferencedSaveOnly feature is flipped
		opt.DoSaves = (cmdName == "BUILD" && c.opt.DoSaves)
	}
	if err != nil {
		// Contradiction allowed. You can BUILD another target with different platform.
		opt.Platform = platform
	}
	return target, opt, propagateBuildArgs, nil
}

func (c *Converter) buildTarget(ctx context.Context, fullTargetName string, platform *specs.Platform, allowPrivileged bool, buildArgs []string, isDangling bool, cmdName string) (*states.MultiTarget, error) {
	target, opt, propagateBuildArgs, err := c.prepBuildTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs, isDangling, cmdName)
	if err != nil {
		return nil, err
	}
	mts, err := Earthfile2LLB(ctx, target, opt)
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
		if cmdName == "FROM" {
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

func (c *Converter) internalRun(ctx context.Context, args, secretKeyValues []string, isWithShell bool, shellWrap shellWrapFun, pushFlag, transient, withSSH, noCache, interactive, interactiveKeep bool, commandStr string, opts ...llb.RunOption) (pllb.State, error) {
	isInteractive := (interactive || interactiveKeep)
	if !c.opt.AllowInteractive && isInteractive {
		// This check is here because other places also call here to evaluate RUN-like statements. We catch all potential interactives here.
		return pllb.State{}, errors.New("--interactive options are not allowed, when --strict is specified or otherwise implied")
	}

	finalOpts := opts
	var extraEnvVars []string
	// Secrets.
	for _, secretKeyValue := range secretKeyValues {
		parts := strings.SplitN(secretKeyValue, "=", 2)
		if len(parts) != 2 {
			return pllb.State{}, errors.Errorf("invalid secret definition %s", secretKeyValue)
		}
		if strings.HasPrefix(parts[1], "+secrets/") {
			envVar := parts[0]
			secretID := strings.TrimPrefix(parts[1], "+secrets/")
			secretPath := path.Join("/run/secrets", secretID)
			secretOpts := []llb.SecretOption{
				llb.SecretID(secretID),
				// TODO: Perhaps this should just default to the current user automatically from
				//       buildkit side. Then we wouldn't need to open this up to everyone.
				llb.SecretFileOpt(0, 0, 0444),
			}
			finalOpts = append(finalOpts, llb.AddSecret(secretPath, secretOpts...))
			// TODO: The use of cat here might not be portable.
			extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"$(cat %s)\"", envVar, secretPath))
		} else if parts[1] == "" {
			// If empty string, don't use (used for optional secrets).
			// TODO: This should be an actual secret (with an empty value),
			//       so that the cache works correctly.
		} else {
			return pllb.State{}, errors.Errorf("secret definition %s not supported. Must start with +secrets/ or be an empty string", secretKeyValue)
		}
	}
	// Build args.
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _ := c.varCollection.GetActive(buildArgName)
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=%s", buildArgName, shellescape.Quote(ba)))
	}
	// Debugger.
	secretOpts := []llb.SecretOption{
		llb.SecretID(common.DebuggerSettingsSecretsKey),
		llb.SecretFileOpt(0, 0, 0444),
	}
	debuggerSecretMount := llb.AddSecret(
		fmt.Sprintf("/run/secrets/%s", common.DebuggerSettingsSecretsKey), secretOpts...)
	debuggerMount := pllb.AddMount(debuggerPath, pllb.Scratch(),
		llb.HostBind(), llb.SourcePath("/usr/bin/earth_debugger"))
	finalOpts = append(finalOpts, debuggerSecretMount, debuggerMount)
	if withSSH {
		finalOpts = append(finalOpts, llb.AddSSHSocket())
	}
	// Shell and debugger wrap.
	finalArgs := shellWrap(args, extraEnvVars, isWithShell, true, isInteractive)
	finalOpts = append(finalOpts, llb.Args(finalArgs))
	if noCache {
		finalOpts = append(finalOpts, llb.IgnoreCache)
	}

	if pushFlag {
		// For push-flagged commands, make sure they run every time - don't use cache.
		finalOpts = append(finalOpts, llb.IgnoreCache)
		if !c.mts.Final.RunPush.HasState {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.Final.RunPush.State = c.mts.Final.MainState
			c.mts.Final.RunPush.HasState = true
		}
	}

	if isInteractive {
		finalOpts = append(finalOpts, llb.IgnoreCache)
		c.mts.Final.RanInteractive = true

		switch {
		case interactive:
			is := states.InteractiveSession{
				CommandStr:  commandStr,
				Initialized: true,
				Kind:        states.SessionEphemeral,
			}

			if pushFlag {
				is.State = c.mts.Final.RunPush.State.Run(finalOpts...).Root()
				c.mts.Final.RunPush.InteractiveSession = is
				return c.mts.Final.RunPush.State, nil

			}
			is.State = c.mts.Final.MainState.Run(finalOpts...).Root()
			c.mts.Final.InteractiveSession = is
			return c.mts.Final.MainState, nil

		case interactiveKeep:
			c.mts.Final.InteractiveSession = states.InteractiveSession{
				CommandStr:  commandStr,
				Initialized: true,
				Kind:        states.SessionKeep,
			}
		}
	}

	if pushFlag {
		// Don't run on MainState. We want push-flagged commands to be executed only
		// *after* the build. Save this for later.
		c.mts.Final.RunPush.State = c.mts.Final.RunPush.State.Run(finalOpts...).Root()
		c.mts.Final.RunPush.CommandStrs = append(c.mts.Final.RunPush.CommandStrs, commandStr)
		return c.mts.Final.RunPush.State, nil
	} else if transient {
		transientState := c.mts.Final.MainState.Run(finalOpts...).Root()
		return transientState, nil
	} else {
		c.mts.Final.MainState = c.mts.Final.MainState.Run(finalOpts...).Root()
		return c.mts.Final.MainState, nil
	}
}

func (c *Converter) readArtifact(ctx context.Context, mts *states.MultiTarget, artifact domain.Artifact) ([]byte, error) {
	if mts.Final.ArtifactsState.Output() == nil {
		// ArtifactsState is scratch - no artifact has been copied.
		return nil, errors.Errorf("artifact %s not found; no SAVE ARTIFACT command was issued in %s", artifact.String(), artifact.Target.String())
	}
	ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, mts.Final.ArtifactsState, mts.Final.Platform, c.opt.CacheImports.AsMap())
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

func (c *Converter) internalFromClassical(ctx context.Context, imageName string, platform specs.Platform, opts ...llb.ImageOption) (pllb.State, *image.Image, *variables.Scope, error) {
	if imageName == "scratch" {
		// FROM scratch
		img := image.NewImage()
		img.OS = platform.OS
		img.Architecture = platform.Architecture
		return pllb.Scratch().Platform(platform), img, nil, nil
	}
	ref, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return pllb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", imageName)
	}
	baseImageName := reference.TagNameOnly(ref).String()
	logName := fmt.Sprintf(
		"%sLoad metadata %s",
		c.imageVertexPrefix(imageName), llbutil.PlatformToString(&platform))
	dgst, dt, err := c.opt.MetaResolver.ResolveImageConfig(
		ctx, baseImageName,
		llb.ResolveImageConfigOpt{
			Platform:    &platform,
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
	allOpts := append(opts, llb.Platform(platform), c.opt.ImageResolveMode)
	state := pllb.Image(ref.String(), allOpts...)
	state, img2, envVars := c.applyFromImage(state, &img)
	return state, img2, envVars, nil
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
		// Run the expression on the main state, but don't alter main with the resulting state.
		srcBuildArgDir := "/run/buildargs"
		srcBuildArgPath := path.Join(srcBuildArgDir, name)
		c.mts.Final.MainState = c.mts.Final.MainState.File(
			pllb.Mkdir(srcBuildArgDir, 0755, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] mkdir %s", srcBuildArgDir))
		args := strings.Split(fmt.Sprintf("echo \"%s\" >%s", expression, srcBuildArgPath), " ")
		transient := true
		state, err := c.internalRun(
			ctx, args, []string{}, true, withShellAndEnvVars, false, transient, false, false, false, false, expression,
			llb.WithCustomNamef("%sRUN %s", c.vertexPrefix(false, false), expression))
		if err != nil {
			return "", 0, errors.Wrapf(err, "run %v", expression)
		}
		ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, state, c.opt.Platform, c.opt.CacheImports.AsMap())
		if err != nil {
			return "", 0, errors.Wrapf(err, "build arg state to ref")
		}
		value, err := ref.ReadFile(ctx, gwclient.ReadRequest{Filename: srcBuildArgPath})
		if err != nil {
			return "", 0, errors.Wrapf(err, "non constant build arg read request")
		}
		// echo adds a trailing \n.
		value = bytes.TrimSuffix(value, []byte("\n"))
		return string(value), 0, nil
	}
}

var base64True = base64.StdEncoding.EncodeToString([]byte("true"))

func (c *Converter) vertexPrefix(local bool, interactive bool) string {
	overriding := c.varCollection.SortedOverridingVariables()
	varStrBuilder := make([]string, 0, len(overriding)+1)
	if c.mts.Final.Platform != nil {
		b64Platform := base64.StdEncoding.EncodeToString(
			[]byte(llbutil.PlatformToString(c.opt.Platform)))
		varStrBuilder = append(varStrBuilder, fmt.Sprintf("@platform=%s", b64Platform))
	}
	if local {
		varStrBuilder = append(varStrBuilder, fmt.Sprintf("@local=%s", base64True))
	}
	if interactive {
		varStrBuilder = append(varStrBuilder, fmt.Sprintf("@interactive=%s", base64True))
	}
	for _, key := range overriding {
		variable, isActive := c.varCollection.GetActive(key)
		if !isActive {
			continue
		}
		b64Value := base64.StdEncoding.EncodeToString([]byte(variable))
		varStrBuilder = append(varStrBuilder, fmt.Sprintf("%s=%s", key, b64Value))
	}
	var varStr string
	if len(varStrBuilder) > 0 {
		varStr = fmt.Sprintf("(%s)", strings.Join(varStrBuilder, " "))
	}
	return fmt.Sprintf("[%s%s %s] ", c.mts.Final.Target.String(), varStr, c.mts.Final.ID)
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
				c.mts.Final.MainState, dep.MainState, c.mts.Final.Target.String(), dep.Target.String())
		}
	}
	// Clear the direct deps so we don't do this again.
	c.directDeps = nil
}

func (c *Converter) imageVertexPrefix(id string) string {
	h := fnv.New32a()
	h.Write([]byte(id))
	return fmt.Sprintf("[%s %d] ", id, h.Sum32())
}

func (c *Converter) vertexPrefixWithURL(url string) string {
	return fmt.Sprintf("[%s(%s) %s] ", c.mts.Final.Target.String(), url, url)
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

func (c *Converter) setPlatform(platform *specs.Platform) {
	c.opt.Platform = platform
	c.mts.Final.SetPlatform(platform)
	c.varCollection.SetPlatform(llbutil.PlatformWithDefault(platform))
}

func (c *Converter) joinRefs(relRef domain.Reference) (domain.Reference, error) {
	return domain.JoinReferences(c.varCollection.AbsRef(), relRef)
}

func (c *Converter) checkAllowed(command string) error {
	if c.mts.Final.RanInteractive && command != "SAVE IMAGE" {
		return errors.New("If present, a single --interactive command must be the last command in a target")
	}

	if c.mts.Final.RanFromLike {
		return nil
	}

	switch command {
	case "FROM", "FROM DOCKERFILE", "LOCALLY", "BUILD", "ARG", "IMPORT":
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
