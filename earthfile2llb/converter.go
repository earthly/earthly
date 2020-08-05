package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/dedup"
	"github.com/earthly/earthly/earthfile2llb/image"
	"github.com/earthly/earthly/earthfile2llb/imr"
	"github.com/earthly/earthly/earthfile2llb/variables"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/llbutil/llbgit"
	"github.com/earthly/earthly/logging"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/pkg/errors"
)

// DockerBuilderFun is a function able to build a target into a docker tar file.
type DockerBuilderFun = func(ctx context.Context, mts *MultiTargetStates, dockerTag string, outFile string) error

var _ commandInterpreter = (*Converter)(nil)

// Converter turns earth commands to buildkit LLB representation.
type Converter struct {
	gitMeta              *buildcontext.GitMetadata
	resolver             *buildcontext.Resolver
	mts                  *MultiTargetStates
	directDeps           []*SingleTargetStates
	directDepIndices     []int
	buildContext         llb.State
	cacheContext         llb.State
	varCollection        *variables.Collection
	dockerBuilderFun     DockerBuilderFun
	cleanCollection      *cleanup.Collection
	nextArgIndex         int
	interactiveDebugging bool
	debuggerImage        string
	remoteConsoleAddr    string
}

// NewConverter constructs a new converter for a given earth target.
func NewConverter(ctx context.Context, target domain.Target, resolver *buildcontext.Resolver, dockerBuilderFun DockerBuilderFun, cleanCollection *cleanup.Collection, bc *buildcontext.Data, visitedStates map[string][]*SingleTargetStates, varCollection *variables.Collection, interactiveDebugging bool, debuggerImage, remoteConsoleAddr string) (*Converter, error) {
	sts := &SingleTargetStates{
		Target: target,
		TargetInput: dedup.TargetInput{
			TargetCanonical: target.StringCanonical(),
		},
		SideEffectsState: llb.Scratch().Platform(llbutil.TargetPlatform),
		SideEffectsImage: image.NewImage(),
		ArtifactsState:   llb.Scratch().Platform(llbutil.TargetPlatform),
		LocalDirs:        bc.LocalDirs,
		Ongoing:          true,
		Salt:             fmt.Sprintf("%d", rand.Int()),
	}
	mts := &MultiTargetStates{
		FinalStates:   sts,
		VisitedStates: visitedStates,
	}
	targetStr := target.String()
	visitedStates[targetStr] = append(visitedStates[targetStr], sts)
	return &Converter{
		gitMeta:              bc.GitMetadata,
		resolver:             resolver,
		mts:                  mts,
		buildContext:         bc.BuildContext,
		cacheContext:         makeCacheContext(target),
		varCollection:        varCollection.WithBuiltinBuildArgs(target, bc.GitMetadata),
		dockerBuilderFun:     dockerBuilderFun,
		cleanCollection:      cleanCollection,
		interactiveDebugging: interactiveDebugging,
		debuggerImage:        debuggerImage,
		remoteConsoleAddr:    remoteConsoleAddr,
	}, nil
}

// From applies the earth FROM command.
func (c *Converter) From(ctx context.Context, imageName string, buildArgs []string) error {
	imageName = c.expandArgs(imageName)
	for i := range buildArgs {
		buildArgs[i] = c.expandArgs(buildArgs[i])
	}
	if strings.Contains(imageName, "+") {
		// Target-based FROM.
		return c.fromTarget(ctx, imageName, buildArgs)
	}

	// Docker image based FROM.
	if len(buildArgs) != 0 {
		return errors.New("--build-arg not supported in non-target FROM")
	}
	return c.fromClassical(ctx, imageName)
}

func (c *Converter) fromClassical(ctx context.Context, imageName string) error {
	state, img, newVariables, err := c.internalFromClassical(
		ctx, imageName,
		llb.WithCustomNamef("%sFROM %s", c.vertexPrefix(), imageName))
	if err != nil {
		return err
	}
	c.mts.FinalStates.SideEffectsState = state
	c.mts.FinalStates.SideEffectsImage = img
	c.varCollection = newVariables
	return nil
}

func (c *Converter) fromTarget(ctx context.Context, targetName string, buildArgs []string) error {
	logger := logging.GetLogger(ctx).With("from-target", targetName).With("build-args", buildArgs)
	logger.Info("Applying FROM target")
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}
	mts, err := c.Build(ctx, depTarget.String(), buildArgs)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", depTarget.String())
	}
	if depTarget.IsLocalInternal() {
		depTarget.LocalPath = c.mts.FinalStates.Target.LocalPath
	}
	// Look for the built state in the dep states, after we've built it.
	relevantDepState := mts.FinalStates
	saveImage, ok := relevantDepState.LastSaveImage()
	if !ok {
		return fmt.Errorf(
			"FROM statement: referenced target %s does not contain a SAVE IMAGE statement",
			depTarget.String())
	}

	// Pass on dep state over to this state.
	c.mts.FinalStates.SideEffectsState = saveImage.State
	for dirKey, dirValue := range relevantDepState.LocalDirs {
		c.mts.FinalStates.LocalDirs[dirKey] = dirValue
	}
	c.mts.FinalStates.SideEffectsImage = saveImage.Image.Clone()
	return nil
}

// CopyArtifact applies the earth COPY artifact command.
func (c *Converter) CopyArtifact(ctx context.Context, artifactName string, dest string, buildArgs []string, isDir bool) error {
	artifactName = c.expandArgs(artifactName)
	dest = c.expandArgs(dest)
	for i := range buildArgs {
		buildArgs[i] = c.expandArgs(buildArgs[i])
	}
	logging.GetLogger(ctx).
		With("srcArtifact", artifactName).
		With("dest", dest).
		With("build-args", buildArgs).
		With("dir", isDir).
		Info("Applying COPY (artifact)")
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	mts, err := c.Build(ctx, artifact.Target.String(), buildArgs)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", artifact.Target.String())
	}
	if artifact.Target.IsLocalInternal() {
		artifact.Target.LocalPath = c.mts.FinalStates.Target.LocalPath
	}
	// Grab the artifacts state in the dep states, after we've built it.
	relevantDepState := mts.FinalStates
	// Copy.
	c.mts.FinalStates.SideEffectsState = llbutil.CopyOp(
		relevantDepState.ArtifactsState, []string{artifact.Artifact},
		c.mts.FinalStates.SideEffectsState, dest, true, isDir,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s %s",
			c.vertexPrefix(),
			strIf(isDir, "--dir "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest))
	return nil
}

// CopyClassical applies the earth COPY command, with classical args.
func (c *Converter) CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool) {
	dest = c.expandArgs(dest)
	for i := range srcs {
		srcs[i] = c.expandArgs(srcs[i])
	}
	logging.GetLogger(ctx).
		With("srcs", srcs).
		With("dest", dest).
		With("dir", isDir).
		Info("Applying COPY (classical)")
	c.mts.FinalStates.SideEffectsState = llbutil.CopyOp(
		c.buildContext, srcs, c.mts.FinalStates.SideEffectsState, dest, true, isDir,
		llb.WithCustomNamef(
			"%sCOPY %s%s %s",
			c.vertexPrefix(),
			strIf(isDir, "--dir "),
			strings.Join(srcs, " "),
			dest))
}

// Run applies the earth RUN command.
func (c *Converter) Run(ctx context.Context, args []string, mounts []string, secretKeyValues []string, privileged bool, withEntrypoint bool, withDocker bool, isWithShell bool, pushFlag bool) error {
	// TODO: This does not work, because it strips away some quotes, which are valuable to the shell.
	//       In any case, this is probably working as intended as is.
	// if !isWithShell {
	// 	for i := range args {
	// 		args[i] = c.expandArgs(args[i])
	// 	}
	// }
	for i := range mounts {
		mounts[i] = c.expandArgs(mounts[i])
	}
	logging.GetLogger(ctx).
		With("args", args).
		With("mounts", mounts).
		With("secrets", secretKeyValues).
		With("privileged", privileged).
		With("withEntrypoint", withEntrypoint).
		With("withDocker", withDocker).
		With("push", pushFlag).
		Info("Applying RUN")
	var opts []llb.RunOption
	mountRunOpts, err := parseMounts(mounts, c.mts.FinalStates.Target, c.mts.FinalStates.TargetInput, c.cacheContext)
	if err != nil {
		return errors.Wrap(err, "parse mounts")
	}
	opts = append(opts, mountRunOpts...)

	finalArgs := args
	if withEntrypoint {
		if len(args) == 0 {
			// No args provided. Use the image's CMD.
			args := make([]string, len(c.mts.FinalStates.SideEffectsImage.Config.Cmd))
			copy(args, c.mts.FinalStates.SideEffectsImage.Config.Cmd)
		}
		finalArgs = append(c.mts.FinalStates.SideEffectsImage.Config.Entrypoint, args...)
		isWithShell = false // Don't use shell when --entrypoint is passed.
	}
	if privileged {
		opts = append(opts, llb.Security(llb.SecurityModeInsecure))
	}
	runStr := fmt.Sprintf(
		"RUN %s%s%s%s%s",
		strIf(privileged, "--privileged "),
		strIf(withDocker, "--with-docker "),
		strIf(withEntrypoint, "--entrypoint "),
		strIf(pushFlag, "--push "),
		strings.Join(finalArgs, " "))
	opts = append(opts, llb.WithCustomNamef("%s%s", c.vertexPrefix(), runStr))
	return c.internalRun(ctx, finalArgs, secretKeyValues, withDocker, isWithShell, pushFlag, runStr, opts...)
}

// SaveArtifact applies the earth SAVE ARTIFACT command.
func (c *Converter) SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string) error {
	saveFrom = c.expandArgs(saveFrom)
	saveTo = c.expandArgs(saveTo)
	saveAsLocalTo = c.expandArgs(saveAsLocalTo)
	logging.GetLogger(ctx).
		With("saveFrom", saveFrom).
		With("saveTo", saveTo).
		With("saveAsLocalTo", saveAsLocalTo).
		Info("Applying SAVE ARTIFACT")
	saveToAdjusted := saveTo
	if saveTo == "" || saveTo == "." || strings.HasSuffix(saveTo, "/") {
		absSaveFrom, err := llbutil.Abs(ctx, c.mts.FinalStates.SideEffectsState, saveFrom)
		if err != nil {
			return err
		}
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
		Target:   c.mts.FinalStates.Target,
		Artifact: artifactPath,
	}
	c.mts.FinalStates.ArtifactsState = llbutil.CopyOp(
		c.mts.FinalStates.SideEffectsState, []string{saveFrom}, c.mts.FinalStates.ArtifactsState,
		saveToAdjusted, true, true,
		llb.WithCustomNamef("%sSAVE ARTIFACT %s %s", c.vertexPrefix(), saveFrom, artifact.String()))
	if saveAsLocalTo != "" {
		separateArtifactsState := llb.Scratch().Platform(llbutil.TargetPlatform)
		separateArtifactsState = llbutil.CopyOp(
			c.mts.FinalStates.SideEffectsState, []string{saveFrom}, separateArtifactsState,
			saveToAdjusted, true, false,
			llb.WithCustomNamef("%sSAVE ARTIFACT %s %s", c.vertexPrefix(), saveFrom, artifact.String()))
		c.mts.FinalStates.SeparateArtifactsState = append(c.mts.FinalStates.SeparateArtifactsState, separateArtifactsState)
		c.mts.FinalStates.SaveLocals = append(c.mts.FinalStates.SaveLocals, SaveLocal{
			DestPath:     saveAsLocalTo,
			ArtifactPath: artifactPath,
			Index:        len(c.mts.FinalStates.SeparateArtifactsState) - 1,
		})
	}
	return nil
}

// SaveImage applies the earth SAVE IMAGE command.
func (c *Converter) SaveImage(ctx context.Context, imageNames []string, pushImages bool) {
	for i := range imageNames {
		imageNames[i] = c.expandArgs(imageNames[i])
	}
	logging.GetLogger(ctx).With("image", imageNames).With("push", pushImages).Info("Applying SAVE IMAGE")
	if len(imageNames) == 0 {
		// Use an empty image name if none provided. This will not be exported
		// as docker image, but will allow for importing / referencing within
		// earthfiles.
		imageNames = []string{""}
	}
	for _, imageName := range imageNames {
		c.mts.FinalStates.SaveImages = append(c.mts.FinalStates.SaveImages, SaveImage{
			State:     c.mts.FinalStates.SideEffectsState,
			Image:     c.mts.FinalStates.SideEffectsImage.Clone(),
			DockerTag: imageName,
			Push:      pushImages,
		})
	}
}

// Build applies the earth BUILD command.
func (c *Converter) Build(ctx context.Context, fullTargetName string, buildArgs []string) (*MultiTargetStates, error) {
	fullTargetName = c.expandArgs(fullTargetName)
	for i := range buildArgs {
		buildArgs[i] = c.expandArgs(buildArgs[i])
	}
	logging.GetLogger(ctx).
		With("full-target-name", fullTargetName).
		With("build-args", buildArgs).
		Info("Applying BUILD")

	target, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		return nil, errors.Wrapf(err, "earth target parse %s", fullTargetName)
	}

	if c.mts.FinalStates.Target.IsRemote() {
		// Current target is remotee. Turn relative targets into remote
		// targets.
		if !target.IsRemote() {
			target.Registry = c.mts.FinalStates.Target.Registry
			target.ProjectPath = c.mts.FinalStates.Target.ProjectPath
			target.Tag = c.mts.FinalStates.Target.Tag
			if target.IsLocalExternal() {
				if path.IsAbs(target.LocalPath) {
					return nil, fmt.Errorf(
						"Absolute path %s not supported as reference in external target context", target.LocalPath)
				}
				target.ProjectPath = path.Join(
					c.mts.FinalStates.Target.ProjectPath, target.LocalPath)
				target.LocalPath = ""
			} else if target.IsLocalInternal() {
				target.LocalPath = ""
			}
		}
	} else {
		if target.IsLocalExternal() {
			if path.IsAbs(target.LocalPath) {
				target.LocalPath = path.Clean(target.LocalPath)
			} else {
				target.LocalPath = path.Join(c.mts.FinalStates.Target.LocalPath, target.LocalPath)
				if !strings.HasPrefix(target.LocalPath, ".") {
					target.LocalPath = fmt.Sprintf("./%s", target.LocalPath)
				}
			}
		} else if target.IsLocalInternal() {
			target.LocalPath = c.mts.FinalStates.Target.LocalPath
		}
	}
	newVarCollection, err := c.varCollection.WithParseBuildArgs(
		buildArgs, c.processNonConstantBuildArgFunc(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "parse build args")
	}
	// Recursion.
	mts, err := Earthfile2LLB(
		ctx, target, c.resolver, c.dockerBuilderFun, c.cleanCollection,
		c.mts.VisitedStates, newVarCollection, c.interactiveDebugging, c.debuggerImage, c.remoteConsoleAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "earthfile2llb for %s", fullTargetName)
	}
	c.directDeps = append(c.directDeps, mts.FinalStates)
	return mts, nil
}

// Workdir applies the WORKDIR command.
func (c *Converter) Workdir(ctx context.Context, workdirPath string) {
	workdirPath = c.expandArgs(workdirPath)
	logging.GetLogger(ctx).With("workdir", workdirPath).Info("Applying WORKDIR")
	c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.Dir(workdirPath)
	workdirAbs := workdirPath
	if !path.IsAbs(workdirAbs) {
		workdirAbs = path.Join("/", c.mts.FinalStates.SideEffectsImage.Config.WorkingDir, workdirAbs)
	}
	c.mts.FinalStates.SideEffectsImage.Config.WorkingDir = workdirAbs
	if workdirAbs != "/" {
		// Mkdir.
		mkdirOpts := []llb.MkdirOption{
			llb.WithParents(true),
		}
		if c.mts.FinalStates.SideEffectsImage.Config.User != "" {
			mkdirOpts = append(mkdirOpts, llb.WithUser(c.mts.FinalStates.SideEffectsImage.Config.User))
		}
		opts := []llb.ConstraintsOpt{
			llb.WithCustomNamef("%sWORKDIR %s", c.vertexPrefix(), workdirPath),
		}
		c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.File(
			llb.Mkdir(workdirAbs, 0755, mkdirOpts...), opts...)
	}
}

// User applies the USER command.
func (c *Converter) User(ctx context.Context, user string) {
	user = c.expandArgs(user)
	logging.GetLogger(ctx).With("user", user).Info("Applying USER")
	c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.User(user)
	c.mts.FinalStates.SideEffectsImage.Config.User = user
}

// Cmd applies the CMD command.
func (c *Converter) Cmd(ctx context.Context, cmdArgs []string, isWithShell bool) {
	if !isWithShell {
		for i := range cmdArgs {
			cmdArgs[i] = c.expandArgs(cmdArgs[i])
		}
	}
	logging.GetLogger(ctx).With("cmd", cmdArgs).Info("Applying CMD")
	c.mts.FinalStates.SideEffectsImage.Config.Cmd = withShell(cmdArgs, isWithShell)
}

// Entrypoint applies the ENTRYPOINT command.
func (c *Converter) Entrypoint(ctx context.Context, entrypointArgs []string, isWithShell bool) {
	if !isWithShell {
		for i := range entrypointArgs {
			entrypointArgs[i] = c.expandArgs(entrypointArgs[i])
		}
	}
	logging.GetLogger(ctx).With("entrypoint", entrypointArgs).Info("Applying ENTRYPOINT")
	c.mts.FinalStates.SideEffectsImage.Config.Entrypoint = withShell(entrypointArgs, isWithShell)
}

// Expose applies the EXPOSE command.
func (c *Converter) Expose(ctx context.Context, ports []string) {
	for i := range ports {
		ports[i] = c.expandArgs(ports[i])
	}
	logging.GetLogger(ctx).With("ports", ports).Info("Applying EXPOSE")
	for _, port := range ports {
		c.mts.FinalStates.SideEffectsImage.Config.ExposedPorts[port] = struct{}{}
	}
}

// Volume applies the VOLUME command.
func (c *Converter) Volume(ctx context.Context, volumes []string) {
	for i := range volumes {
		volumes[i] = c.expandArgs(volumes[i])
	}
	logging.GetLogger(ctx).With("volumes", volumes).Info("Applying VOLUME")
	for _, volume := range volumes {
		c.mts.FinalStates.SideEffectsImage.Config.Volumes[volume] = struct{}{}
	}
}

// Env applies the ENV command.
func (c *Converter) Env(ctx context.Context, envKey string, envValue string) {
	envValue = c.expandArgs(envValue)
	logging.GetLogger(ctx).With("env-key", envKey).With("env-value", envValue).Info("Applying ENV")
	c.varCollection.AddActive(envKey, variables.NewConstantEnvVar(envValue), true)
	c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.AddEnv(envKey, envValue)
	c.mts.FinalStates.SideEffectsImage.Config.Env = variables.AddEnv(
		c.mts.FinalStates.SideEffectsImage.Config.Env, envKey, envValue)
}

// Arg applies the ARG command.
func (c *Converter) Arg(ctx context.Context, argKey string, defaultArgValue string) {
	defaultArgValue = c.expandArgs(defaultArgValue)
	logging.GetLogger(ctx).With("arg-key", argKey).With("arg-value", defaultArgValue).Info("Applying ARG")
	effective := c.varCollection.AddActive(argKey, variables.NewConstant(defaultArgValue), false)
	c.mts.FinalStates.TargetInput.BuildArgs = append(
		c.mts.FinalStates.TargetInput.BuildArgs,
		effective.BuildArgInput(argKey, defaultArgValue))
}

// Label applies the LABEL command.
func (c *Converter) Label(ctx context.Context, labels map[string]string) {
	labels2 := make(map[string]string)
	for key, value := range labels {
		key2 := c.expandArgs(key)
		value2 := c.expandArgs(value)
		labels2[key2] = value2
	}
	logging.GetLogger(ctx).With("labels", labels2).Info("Applying LABEL")
	for key, value := range labels2 {
		c.mts.FinalStates.SideEffectsImage.Config.Labels[key] = value
	}
}

// GitClone applies the GIT CLONE command.
func (c *Converter) GitClone(ctx context.Context, gitURL string, branch string, dest string) error {
	gitURL = c.expandArgs(gitURL)
	branch = c.expandArgs(branch)
	dest = c.expandArgs(dest)
	logging.GetLogger(ctx).With("git-url", gitURL).With("branch", branch).Info("Applying GIT CLONE")
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef(
			"%sGIT CLONE (--branch %s) %s", c.vertexPrefixWithURL(gitURL), branch, gitURL),
		llb.KeepGitDir(),
	}
	gitState := llbgit.Git(gitURL, branch, gitOpts...)
	c.mts.FinalStates.SideEffectsState = llbutil.CopyOp(
		gitState, []string{"."}, c.mts.FinalStates.SideEffectsState, dest, false, false,
		llb.WithCustomNamef(
			"%sCOPY GIT CLONE (--branch %s) %s TO %s", c.vertexPrefix(),
			branch, gitURL, dest))
	return nil
}

// DockerLoad applies the DOCKER LOAD command.
func (c *Converter) DockerLoad(ctx context.Context, targetName string, dockerTag string, buildArgs []string) error {
	targetName = c.expandArgs(targetName)
	dockerTag = c.expandArgs(dockerTag)
	for i := range buildArgs {
		buildArgs[i] = c.expandArgs(buildArgs[i])
	}
	logging.GetLogger(ctx).With("target-name", targetName).With("dockerTag", dockerTag).Info("Applying DOCKER LOAD")
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target %s", targetName)
	}
	mts, err := c.Build(ctx, depTarget.String(), buildArgs)
	if err != nil {
		return err
	}
	err = c.solveAndLoad(
		ctx, mts, depTarget.String(), dockerTag,
		llb.WithCustomNamef(
			"%sDOCKER LOAD %s %s", c.vertexPrefix(), depTarget.String(), dockerTag))
	if err != nil {
		return err
	}
	return nil
}

// DockerPull applies the DOCKER PULL command.
func (c *Converter) DockerPull(ctx context.Context, dockerTag string) error {
	dockerTag = c.expandArgs(dockerTag)
	logging.GetLogger(ctx).With("dockerTag", dockerTag).Info("Applying DOCKER PULL")
	state, image, _, err := c.internalFromClassical(
		ctx, dockerTag,
		llb.WithCustomNamef("%sDOCKER PULL %s", c.vertexPrefix(), dockerTag),
	)
	if err != nil {
		return err
	}
	mts := &MultiTargetStates{
		FinalStates: &SingleTargetStates{
			SideEffectsState: state,
			SideEffectsImage: image,
			SaveImages: []SaveImage{
				{
					State:     state,
					Image:     image,
					DockerTag: dockerTag,
				},
			},
		},
	}
	err = c.solveAndLoad(
		ctx, mts, dockerTag, dockerTag,
		llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", c.vertexPrefix(), dockerTag))
	if err != nil {
		return err
	}
	return nil
}

// Healthcheck applies the HEALTHCHECK command.
func (c *Converter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int) {
	for index := range cmdArgs {
		cmdArgs[index] = c.expandArgs(cmdArgs[index])
	}
	logging.GetLogger(ctx).
		With("isNone", isNone).
		With("cmdArgs", cmdArgs).
		With("interval", interval).
		With("timeout", timeout).
		With("startPeriod", startPeriod).
		With("retries", retries).
		Info("Applying HEALTHCHECK")
	hc := &dockerfile2llb.HealthConfig{}
	if isNone {
		hc.Test = []string{"NONE"}
	} else {
		// TODO: Should support also CMD without shell (exec form).
		//       See https://github.com/moby/buildkit/blob/master/frontend/dockerfile/dockerfile2llb/image.go#L18
		hc.Test = append([]string{"CMD-SHELL", strings.Join(cmdArgs, " ")})
		hc.Interval = interval
		hc.Timeout = timeout
		hc.StartPeriod = startPeriod
		hc.Retries = retries
	}
	c.mts.FinalStates.SideEffectsImage.Config.Healthcheck = hc
}

// FinalizeStates returns the LLB states.
func (c *Converter) FinalizeStates() *MultiTargetStates {
	// Create an artificial bond to depStates so that side-effects of deps are built automatically.
	for _, depStates := range c.directDeps {
		c.mts.FinalStates.SideEffectsState = withDependency(
			c.mts.FinalStates.SideEffectsState,
			c.mts.FinalStates.Target,
			depStates.SideEffectsState,
			depStates.Target)
	}

	c.mts.FinalStates.Ongoing = false
	return c.mts
}

func (c *Converter) internalRun(ctx context.Context, args []string, secretKeyValues []string, withDocker bool, isWithShell bool, pushFlag bool, commandStr string, opts ...llb.RunOption) error {
	finalOpts := opts
	var extraEnvVars []string
	for _, secretKeyValue := range secretKeyValues {
		parts := strings.SplitN(secretKeyValue, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Invalid secret definition %s", secretKeyValue)
		}
		if !strings.HasPrefix(parts[1], "+secrets/") {
			return fmt.Errorf("Secret definition %s not supported. Must start with +secrets/", secretKeyValue)
		}
		envVar := parts[0]
		secretID := strings.TrimPrefix(parts[1], "+secrets/")
		secretPath := path.Join("/run/secrets", secretID)
		secretOpts := []llb.SecretOption{
			llb.SecretID(secretID),
		}
		finalOpts = append(finalOpts, llb.AddSecret(secretPath, secretOpts...))
		// TODO: The use of cat here might not be portable.
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"$(cat %s)\"", envVar, secretPath))
	}
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _, _ := c.varCollection.Get(buildArgName)
		if ba.IsEnvVar() {
			continue
		}
		if ba.IsConstant() {
			extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"%s\"", buildArgName, ba.ConstantValue()))
		} else {
			buildArgPath := path.Join("/run/buildargs", buildArgName)
			finalOpts = append(finalOpts, llb.AddMount(buildArgPath, ba.VariableState(), llb.SourcePath(buildArgPath)))
			// TODO: The use of cat here might not be portable.
			extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"$(cat %s)\"", buildArgName, buildArgPath))
		}
	}

	finalOpts = append(finalOpts, llb.AddMount("/usr/bin/earth_debugger", llb.Image(c.debuggerImage, llb.MarkImageInternal), llb.SourcePath("/earth_debugger"), llb.Readonly))
	secretOpts := []llb.SecretOption{
		llb.SecretID("earthly_debugger_settings"),
	}
	finalOpts = append(finalOpts, llb.AddSecret("/run/secrets/earthly_debugger_settings", secretOpts...))

	var finalArgs []string
	if withDocker {
		finalArgs = withDockerdWrap(args, extraEnvVars, isWithShell)
	} else {
		finalArgs = withShellAndEnvVars(args, extraEnvVars, isWithShell)
	}

	finalArgs = append([]string{"earth_debugger"}, finalArgs...)

	finalOpts = append(finalOpts, llb.Args(finalArgs))
	if pushFlag {
		// For push-flagged commands, make sure they run every time - don't use cache.
		finalOpts = append(finalOpts, llb.IgnoreCache)
		if !c.mts.FinalStates.RunPush.Initialized {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.FinalStates.RunPush.State = c.mts.FinalStates.SideEffectsState
			c.mts.FinalStates.RunPush.Initialized = true
		}
		// Don't run on SideEffectsState. We want push-flagged commands to be executed only
		// *after* the build. Save this for later.
		c.mts.FinalStates.RunPush.State = c.mts.FinalStates.RunPush.State.Run(finalOpts...).Root()
		c.mts.FinalStates.RunPush.CommandStrs = append(
			c.mts.FinalStates.RunPush.CommandStrs, commandStr)
	} else {
		c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.Run(finalOpts...).Root()
	}
	return nil
}

func (c *Converter) solveAndLoad(ctx context.Context, mts *MultiTargetStates, opName string, dockerTag string, opts ...llb.RunOption) error {
	// Use a builder to create docker .tar file, mount it via a local build context,
	// then docker load it within the current side effects state.
	outDir, err := ioutil.TempDir("/tmp", "earthly-docker-load")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for docker load")
	}
	c.cleanCollection.Add(func() error {
		return os.RemoveAll(outDir)
	})
	outFile := path.Join(outDir, "image.tar")
	// TODO: This ends up printing some repetitive output, as it builds
	//       the dep twice (even though it's cached the second time).
	err = c.dockerBuilderFun(ctx, mts, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "build target %s for docker load", opName)
	}
	dockerImageID, err := dockertar.GetID(outFile)
	if err != nil {
		return errors.Wrap(err, "inspect docker tar after build")
	}
	// Use the docker image ID + dockerTag as sessionID. This will cause
	// buildkit to use cache when these are the same as before (eg a docker image
	// that is identical as before).
	sessionIDKey := fmt.Sprintf("%s-%s", dockerTag, dockerImageID)
	sha256SessionIDKey := sha256.Sum256([]byte(sessionIDKey))
	sessionID := hex.EncodeToString(sha256SessionIDKey[:])
	// Add the tar to the local context.
	tarContext := llb.Local(
		opName,
		llb.SharedKeyHint(opName),
		llb.SessionID(sessionID),
		llb.Platform(llbutil.TargetPlatform),
		llb.WithCustomNamef("[internal] docker tar context %s %s", opName, sessionID),
	)
	c.mts.FinalStates.LocalDirs[opName] = outDir

	c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.File(
		llb.Mkdir("/var/lib/docker", 0755, llb.WithParents(true)),
		llb.WithCustomNamef("[internal] mkdir /var/lib/docker"),
	)
	loadOpts := []llb.RunOption{
		llb.Args(
			withDockerdWrap(
				[]string{"docker", "load", "</src/image.tar"}, []string{}, true)),
		llb.AddMount("/src", tarContext, llb.Readonly),
		llb.Dir("/src"),
		llb.Security(llb.SecurityModeInsecure),
	}
	loadOpts = append(loadOpts, opts...)
	loadOp := c.mts.FinalStates.SideEffectsState.Run(loadOpts...)
	c.mts.FinalStates.SideEffectsState = loadOp.AddMount(
		"/var/lib/docker", c.mts.FinalStates.SideEffectsState,
		llb.SourcePath("/var/lib/docker"))
	return nil
}

func (c *Converter) internalFromClassical(ctx context.Context, imageName string, opts ...llb.ImageOption) (llb.State, *image.Image, *variables.Collection, error) {
	logging.GetLogger(ctx).With("image", imageName).Info("Applying FROM")
	if imageName == "scratch" {
		// FROM scratch
		return llb.Scratch().Platform(llbutil.TargetPlatform), image.NewImage(),
			variables.NewCollection(), nil
	}
	ref, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return llb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", imageName)
	}
	baseImageName := reference.TagNameOnly(ref).String()
	metaResolver := imr.Default()
	dgst, dt, err := metaResolver.ResolveImageConfig(
		ctx, baseImageName,
		llb.ResolveImageConfigOpt{
			Platform:    &llbutil.TargetPlatform,
			ResolveMode: llb.ResolveModePreferLocal.String(),
		})
	if err != nil {
		return llb.State{}, nil, nil, errors.Wrapf(err, "resolve image config for %s", imageName)
	}
	var img image.Image
	err = json.Unmarshal(dt, &img)
	if err != nil {
		return llb.State{}, nil, nil, errors.Wrapf(err, "unmarshal image config for %s", imageName)
	}
	if dgst != "" {
		ref, err = reference.WithDigest(ref, dgst)
		if err != nil {
			return llb.State{}, nil, nil, errors.Wrapf(err, "reference add digest %v for %s", dgst, imageName)
		}
	}
	allOpts := append(opts, llb.Platform(llbutil.TargetPlatform))
	state := llb.Image(ref.String(), allOpts...)
	// Reset variables.
	newVarCollection := c.varCollection.WithResetEnvVars()
	for _, envVar := range img.Config.Env {
		k, v := variables.ParseKeyValue(envVar)
		newVarCollection.AddActive(k, variables.NewConstantEnvVar(v), true)
		state = state.AddEnv(k, v)
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
	return state, &img, newVarCollection, nil
}

func (c *Converter) expandArgs(word string) string {
	return c.varCollection.Expand(word)
}

func (c *Converter) processNonConstantBuildArgFunc(ctx context.Context) variables.ProcessNonConstantVariableFunc {
	return func(name string, expression string) (llb.State, dedup.TargetInput, int, error) {
		// Run the expression on the side effects state.
		srcBuildArgDir := "/run/buildargs-src"
		srcBuildArgPath := path.Join(srcBuildArgDir, name)
		c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.File(
			llb.Mkdir(srcBuildArgDir, 0755, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] mkdir %s", srcBuildArgDir))
		buildArgPath := path.Join("/run/buildargs", name)
		args := strings.Split(fmt.Sprintf("echo \"%s\" >%s", expression, srcBuildArgPath), " ")
		err := c.internalRun(
			ctx, args, []string{}, false, true, false, expression,
			llb.WithCustomNamef("%sRUN %s", c.vertexPrefix(), expression))
		if err != nil {
			return llb.State{}, dedup.TargetInput{}, 0, errors.Wrapf(err, "run %v", expression)
		}
		// Copy the result of the expression into a separate, isolated state.
		buildArgState := llb.Scratch().Platform(llbutil.TargetPlatform)
		buildArgState = llbutil.CopyOp(
			c.mts.FinalStates.SideEffectsState, []string{srcBuildArgPath}, buildArgState, buildArgPath, false, false,
			llb.WithCustomNamef("[internal] copy buildarg %s", name))
		// Store the state with the expression result for later use.
		argIndex := c.nextArgIndex
		c.nextArgIndex++
		// Remove intermediary file from side effects state.
		c.mts.FinalStates.SideEffectsState = c.mts.FinalStates.SideEffectsState.File(
			llb.Rm(srcBuildArgPath, llb.WithAllowNotFound(true)),
			llb.WithCustomNamef("[internal] rm %s", srcBuildArgPath))

		return buildArgState, c.mts.FinalStates.TargetInput, argIndex, nil
	}
}

func (c *Converter) vertexPrefix() string {
	return fmt.Sprintf("[%s %s] ", c.mts.FinalStates.Target.String(), c.mts.FinalStates.Salt)
}

func (c *Converter) vertexPrefixWithURL(url string) string {
	return fmt.Sprintf("[%s(%s) %s] ", c.mts.FinalStates.Target.String(), url, url)
}

func withDependency(state llb.State, target domain.Target, depState llb.State, depTarget domain.Target) llb.State {
	return llbutil.WithDependency(
		state, depState,
		llb.WithCustomNamef(
			"[internal] create artificial dependency: %s depends on %s",
			target.String(), depTarget.String()))
}

func makeCacheContext(target domain.Target) llb.State {
	sessionID := cacheKey(target)
	opts := []llb.LocalOption{
		llb.SharedKeyHint(target.ProjectCanonical()),
		llb.SessionID(sessionID),
		llb.Platform(llbutil.TargetPlatform),
		llb.WithCustomNamef("[internal] cache context %s", target.ProjectCanonical()),
	}
	return llb.Local("earthly-cache", opts...)
}

func cacheKey(target domain.Target) string {
	// Use the canonical target, but wihout the tag for cache matching.
	targetCopy := target
	targetCopy.Tag = ""
	digest := sha256.Sum256([]byte(targetCopy.StringCanonical()))
	return hex.EncodeToString(digest[:])
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
