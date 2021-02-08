package earthfile2llb

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/alessio/shellescape"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/variables"

	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session/localhost"
	solverpb "github.com/moby/buildkit/solver/pb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Converter turns earthly commands to buildkit LLB representation.
type Converter struct {
	gitMeta          *gitutil.GitMetadata
	opt              ConvertOpt
	mts              *states.MultiTarget
	directDeps       []*states.SingleTarget
	directDepIndices []int
	buildContext     llb.State
	cacheContext     llb.State
	varCollection    *variables.Collection
	nextArgIndex     int
	ranSave          bool
}

// NewConverter constructs a new converter for a given earthly target.
func NewConverter(ctx context.Context, target domain.Target, bc *buildcontext.Data, opt ConvertOpt) (*Converter, error) {
	sts := &states.SingleTarget{
		Target:   target,
		Platform: opt.Platform,
		TargetInput: dedup.TargetInput{
			TargetCanonical: target.StringCanonical(),
			Platform:        llbutil.PlatformToString(opt.Platform),
		},
		MainState:      llbutil.ScratchWithPlatform(),
		MainImage:      image.NewImage(),
		ArtifactsState: llbutil.ScratchWithPlatform(),
		LocalDirs:      bc.LocalDirs,
		Ongoing:        true,
		Salt:           fmt.Sprintf("%d", rand.Int()),
	}
	mts := &states.MultiTarget{
		Final:   sts,
		Visited: opt.Visited,
	}
	for _, key := range opt.VarCollection.SortedOverridingVariables() {
		ovVar, _, _ := opt.VarCollection.Get(key)
		sts.TargetInput = sts.TargetInput.WithBuildArgInput(ovVar.BuildArgInput(key, ""))
	}
	targetStr := target.String()
	opt.Visited.Add(targetStr, sts)
	return &Converter{
		gitMeta:      bc.GitMetadata,
		opt:          opt,
		mts:          mts,
		buildContext: bc.BuildContext,
		cacheContext: makeCacheContext(target),
		varCollection: opt.VarCollection.WithBuiltinBuildArgs(
			target, llbutil.PlatformWithDefault(opt.Platform), bc.GitMetadata),
	}, nil
}

// From applies the earthly FROM command.
func (c *Converter) From(ctx context.Context, imageName string, platform *specs.Platform, buildArgs []string) error {
	c.nonSaveCommand()
	platform, err := llbutil.ResolvePlatform(platform, c.opt.Platform)
	if err != nil {
		return err
	}
	c.setPlatform(platform)
	if strings.Contains(imageName, "+") {
		// Target-based FROM.
		return c.fromTarget(ctx, imageName, platform, buildArgs)
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
		prefix = c.vertexPrefix(false)
	}
	plat := llbutil.PlatformWithDefault(platform)
	state, img, newVariables, err := c.internalFromClassical(
		ctx, imageName, plat,
		llb.WithCustomNamef("%sFROM %s", prefix, imageName))
	if err != nil {
		return err
	}
	c.mts.Final.MainState = state
	c.mts.Final.MainImage = img
	c.varCollection = newVariables
	return nil
}

func (c *Converter) fromTarget(ctx context.Context, targetName string, platform *specs.Platform, buildArgs []string) error {
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}
	mts, err := c.buildTarget(ctx, depTarget.String(), platform, buildArgs, false)
	if err != nil {
		return errors.Wrapf(err, "apply build %s", depTarget.String())
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
	for _, kv := range saveImage.Image.Config.Env {
		k, v, _ := variables.ParseKeyValue(kv)
		c.varCollection.AddActive(k, variables.NewConstantEnvVar(v))
	}
	c.mts.Final.MainImage = saveImage.Image.Clone()
	return nil
}

// FromDockerfile applies the earthly FROM DOCKERFILE command.
func (c *Converter) FromDockerfile(ctx context.Context, contextPath string, dfPath string, dfTarget string, platform *specs.Platform, buildArgs []string) error {
	platform, err := llbutil.ResolvePlatform(platform, c.opt.Platform)
	if err != nil {
		return err
	}
	c.setPlatform(platform)
	plat := llbutil.PlatformWithDefault(platform)
	c.nonSaveCommand()
	if dfPath != "" {
		// TODO: It's not yet very clear what -f should do. Should it be referencing a Dockerfile
		//       from the build context or the build environment?
		//       Build environment is likely better as it gives maximum flexibility to do
		//       anything.
		return errors.New("FROM DOCKERFILE -f not yet supported")
	}
	var buildContext llb.State
	var dfData []byte
	contextArtifact, parseErr := domain.ParseArtifact(contextPath)
	if parseErr == nil {
		// The Dockerfile and build context are from a target's artifact.
		// TODO: The build args are used for both the artifact and the Dockerfile. This could be
		//       confusing to the user.
		mts, err := c.buildTarget(ctx, contextArtifact.Target.String(), platform, buildArgs, false)
		if err != nil {
			return err
		}
		dfArtifact := contextArtifact
		dfArtifact.Artifact = path.Join(dfArtifact.Artifact, "Dockerfile")
		dfData, err = c.readArtifact(ctx, mts, dfArtifact)
		if err != nil {
			return err
		}
		buildContext = llbutil.ScratchWithPlatform()
		buildContext = llbutil.CopyOp(
			mts.Final.ArtifactsState, []string{contextArtifact.Artifact},
			buildContext, "/", true, true, false, "", false,
			llb.WithCustomNamef(
				"[internal] FROM DOCKERFILE (copy build context from) %s%s",
				joinWrap(buildArgs, "(", " ", ") "), contextArtifact.String()))
	} else {
		// The Dockerfile and build context are from the host.
		if contextPath != "." &&
			!strings.HasPrefix(contextPath, "./") &&
			!strings.HasPrefix(contextPath, "../") &&
			!strings.HasPrefix(contextPath, "/") {
			contextPath = fmt.Sprintf("./%s", contextPath)
		}
		dockerfileMetaTarget := domain.Target{
			Target:    buildcontext.DockerfileMetaTarget,
			LocalPath: contextPath,
		}
		dockerfileMetaTarget, err := domain.JoinTargets(c.mts.FinalTarget(), dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "join targets")
		}
		data, err := c.opt.Resolver.Resolve(ctx, c.opt.GwClient, dockerfileMetaTarget)
		if err != nil {
			return errors.Wrap(err, "resolve build context for dockerfile")
		}
		for ldk, ld := range data.LocalDirs {
			c.mts.Final.LocalDirs[ldk] = ld
		}
		dfPath = data.BuildFilePath
		dfData, err = ioutil.ReadFile(dfPath)
		if err != nil {
			return errors.Wrapf(err, "read file %s", dfPath)
		}
		buildContext = data.BuildContext
	}
	newVarCollection, _, err := c.varCollection.WithParseBuildArgs(
		buildArgs, c.processNonConstantBuildArgFunc(ctx), false)
	if err != nil {
		return err
	}
	caps := solverpb.Caps.CapSet(solverpb.Caps.All())
	state, dfImg, err := dockerfile2llb.Dockerfile2LLB(ctx, dfData, dockerfile2llb.ConvertOpt{
		BuildContext:     &buildContext,
		ContextLocalName: c.mts.FinalTarget().String(),
		MetaResolver:     c.opt.MetaResolver,
		ImageResolveMode: c.opt.ImageResolveMode,
		Target:           dfTarget,
		TargetPlatform:   &plat,
		LLBCaps:          &caps,
		BuildArgs:        newVarCollection.AsMap(),
		Excludes:         nil, // TODO: Need to process this correctly.
	})
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
	state2, img2, newVarCollection := c.applyFromImage(*state, &img)
	c.mts.Final.MainState = state2
	c.mts.Final.MainImage = img2
	c.varCollection = newVarCollection
	return nil
}

// Locally applies the earthly Locally command.
func (c *Converter) Locally(ctx context.Context, platform *specs.Platform) error {
	return c.fromClassical(ctx, "scratch", platform, true)
}

// CopyArtifact applies the earthly COPY artifact command.
func (c *Converter) CopyArtifact(ctx context.Context, artifactName string, dest string, platform *specs.Platform, buildArgs []string, isDir bool, keepTs bool, keepOwn bool, chown string, ifExists bool) error {
	c.nonSaveCommand()
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return errors.Wrapf(err, "parse artifact name %s", artifactName)
	}
	mts, err := c.buildTarget(ctx, artifact.Target.String(), platform, buildArgs, false)
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
		c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), ifExists,
		llb.WithCustomNamef(
			"%sCOPY %s%s%s%s %s",
			c.vertexPrefix(false),
			strIf(isDir, "--dir "),
			strIf(ifExists, "--if-exists "),
			joinWrap(buildArgs, "(", " ", ") "),
			artifact.String(),
			dest))
	return nil
}

// CopyClassical applies the earthly COPY command, with classical args.
func (c *Converter) CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool, keepTs bool, keepOwn bool, chown string) {
	c.nonSaveCommand()
	c.mts.Final.MainState = llbutil.CopyOp(
		c.buildContext, srcs, c.mts.Final.MainState, dest, true, isDir, keepTs, c.copyOwner(keepOwn, chown), false,
		llb.WithCustomNamef(
			"%sCOPY %s%s %s",
			c.vertexPrefix(false),
			strIf(isDir, "--dir "),
			strings.Join(srcs, " "),
			dest))
}

// RunLocal applies a RUN statement locally rather than in a container
func (c *Converter) RunLocal(ctx context.Context, args []string, pushFlag bool) error {
	runStr := fmt.Sprintf("RUN %s%s", strIf(pushFlag, "--push "), strings.Join(args, " "))

	// Build args get propagated into env.
	extraEnvVars := []string{}
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _, _ := c.varCollection.Get(buildArgName)
		if ba.IsEnvVar() {
			continue
		}
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=\"%s\"", buildArgName, ba.ConstantValue()))
	}

	// buildkit-hack in order to run locally, we prepend the command with a UUID
	finalArgs := append([]string{localhost.RunOnLocalHostMagicStr}, withShellAndEnvVars(args, extraEnvVars, true, false)...)
	opts := []llb.RunOption{
		llb.Args(finalArgs),
		llb.IgnoreCache,
		llb.WithCustomNamef("%s%s", c.vertexPrefix(true), runStr),
	}

	if pushFlag {
		if !c.mts.Final.RunPush.Initialized {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.Final.RunPush.State = c.mts.Final.MainState
			c.mts.Final.RunPush.Initialized = true
		}
		c.mts.Final.RunPush.State = c.mts.Final.RunPush.State.Run(opts...).Root()
		c.mts.Final.RunPush.CommandStrs = append(
			c.mts.Final.RunPush.CommandStrs, runStr)
	} else {
		c.mts.Final.MainState = c.mts.Final.MainState.Run(opts...).Root()
	}
	return nil
}

// Run applies the earthly RUN command.
func (c *Converter) Run(ctx context.Context, args, mounts, secretKeyValues []string, privileged, withEntrypoint, withDocker, isWithShell, pushFlag, withSSH, noCache bool) error {
	c.nonSaveCommand()
	if withDocker {
		return errors.New("RUN --with-docker is obsolete. Please use WITH DOCKER ... RUN ... END instead")
	}
	var opts []llb.RunOption
	mountRunOpts, err := parseMounts(mounts, c.mts.Final.Target, c.mts.Final.TargetInput, c.cacheContext)
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
		"RUN %s%s%s%s%s%s",
		strIf(privileged, "--privileged "),
		strIf(withDocker, "--with-docker "),
		strIf(withEntrypoint, "--entrypoint "),
		strIf(pushFlag, "--push "),
		strIf(noCache, "--no-cache "),
		strings.Join(finalArgs, " "))
	shellWrap := withShellAndEnvVars
	opts = append(opts, llb.WithCustomNamef("%s%s", c.vertexPrefix(false), runStr))
	_, err = c.internalRun(
		ctx, finalArgs, secretKeyValues, isWithShell, shellWrap, pushFlag,
		false, withSSH, noCache, runStr, opts...)
	return err
}

// SaveArtifact applies the earthly SAVE ARTIFACT command.
func (c *Converter) SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string, keepTs bool, keepOwn bool, ifExists bool, isPush bool) error {
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
		saveToAdjusted, true, true, keepTs, own, ifExists,
		llb.WithCustomNamef(
			"%sSAVE ARTIFACT %s%s %s", c.vertexPrefix(false), strIf(ifExists, "--if-exists "), saveFrom, artifact.String()))
	if saveAsLocalTo != "" {
		separateArtifactsState := llbutil.ScratchWithPlatform()
		if isPush {
			separateArtifactsState = llbutil.CopyOp(
				c.mts.Final.RunPush.State, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", ifExists,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s %s AS LOCAL %s",
					c.vertexPrefix(false), strIf(ifExists, "--if-exists "), saveFrom, artifact.String(), saveAsLocalTo))
		} else {
			separateArtifactsState = llbutil.CopyOp(
				c.mts.Final.MainState, []string{saveFrom}, separateArtifactsState,
				saveToAdjusted, true, true, keepTs, "root:root", ifExists,
				llb.WithCustomNamef(
					"%sSAVE ARTIFACT %s%s %s AS LOCAL %s",
					c.vertexPrefix(false), strIf(ifExists, "--if-exists "), saveFrom, artifact.String(), saveAsLocalTo))
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
	return c.SaveArtifact(ctx, saveTo, saveTo, "", keepTs, keepOwn, false, false)
}

// SaveImage applies the earthly SAVE IMAGE command.
func (c *Converter) SaveImage(ctx context.Context, imageNames []string, pushImages bool, insecurePush bool, cacheHint bool, cacheFrom []string) error {
	for _, cf := range cacheFrom {
		c.opt.CacheImports[cf] = true
	}
	justCacheHint := false
	if len(imageNames) == 0 && cacheHint {
		imageNames = []string{""}
		justCacheHint = true
	}
	for _, imageName := range imageNames {
		c.mts.Final.SaveImages = append(c.mts.Final.SaveImages, states.SaveImage{
			State:        c.mts.Final.MainState,
			Image:        c.mts.Final.MainImage.Clone(),
			DockerTag:    imageName,
			Push:         pushImages,
			InsecurePush: insecurePush,
			CacheHint:    cacheHint,
		})
		if pushImages && imageName != "" && c.opt.UseInlineCache {
			// Use this image tag as cache import too.
			c.opt.CacheImports[imageName] = true
		}
	}
	if !justCacheHint {
		c.ranSave = true
		c.markFakeDeps()
	}
	return nil
}

// Build applies the earthly BUILD command.
func (c *Converter) Build(ctx context.Context, fullTargetName string, platform *specs.Platform, buildArgs []string) error {
	c.nonSaveCommand()
	_, err := c.buildTarget(ctx, fullTargetName, platform, buildArgs, true)
	return err
}

// Workdir applies the WORKDIR command.
func (c *Converter) Workdir(ctx context.Context, workdirPath string) {
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
			llb.WithCustomNamef("%sWORKDIR %s", c.vertexPrefix(false), workdirPath),
		}
		c.mts.Final.MainState = c.mts.Final.MainState.File(
			llb.Mkdir(workdirAbs, 0755, mkdirOpts...), opts...)
	}
}

// User applies the USER command.
func (c *Converter) User(ctx context.Context, user string) {
	c.nonSaveCommand()
	c.mts.Final.MainState = c.mts.Final.MainState.User(user)
	c.mts.Final.MainImage.Config.User = user
}

// Cmd applies the CMD command.
func (c *Converter) Cmd(ctx context.Context, cmdArgs []string, isWithShell bool) {
	c.nonSaveCommand()
	c.mts.Final.MainImage.Config.Cmd = withShell(cmdArgs, isWithShell)
}

// Entrypoint applies the ENTRYPOINT command.
func (c *Converter) Entrypoint(ctx context.Context, entrypointArgs []string, isWithShell bool) {
	c.nonSaveCommand()
	c.mts.Final.MainImage.Config.Entrypoint = withShell(entrypointArgs, isWithShell)
}

// Expose applies the EXPOSE command.
func (c *Converter) Expose(ctx context.Context, ports []string) {
	c.nonSaveCommand()
	for _, port := range ports {
		c.mts.Final.MainImage.Config.ExposedPorts[port] = struct{}{}
	}
}

// Volume applies the VOLUME command.
func (c *Converter) Volume(ctx context.Context, volumes []string) {
	c.nonSaveCommand()
	for _, volume := range volumes {
		c.mts.Final.MainImage.Config.Volumes[volume] = struct{}{}
	}
}

// Env applies the ENV command.
func (c *Converter) Env(ctx context.Context, envKey string, envValue string) {
	c.nonSaveCommand()
	c.varCollection.AddActive(envKey, variables.NewConstantEnvVar(envValue))
	c.mts.Final.MainState = c.mts.Final.MainState.AddEnv(envKey, envValue)
	c.mts.Final.MainImage.Config.Env = variables.AddEnv(
		c.mts.Final.MainImage.Config.Env, envKey, envValue)
}

// Arg applies the ARG command.
func (c *Converter) Arg(ctx context.Context, argKey string, defaultArgValue string, global bool) error {
	c.nonSaveCommand()
	effective, err := c.varCollection.DeclareActive(argKey, defaultArgValue, global, c.processNonConstantBuildArgFunc((ctx)))
	if err != nil {
		return err
	}
	c.mts.Final.TargetInput = c.mts.Final.TargetInput.WithBuildArgInput(
		effective.BuildArgInput(argKey, effective.ConstantValue()))
	return nil
}

// Label applies the LABEL command.
func (c *Converter) Label(ctx context.Context, labels map[string]string) {
	c.nonSaveCommand()
	for key, value := range labels {
		c.mts.Final.MainImage.Config.Labels[key] = value
	}
}

// GitClone applies the GIT CLONE command.
func (c *Converter) GitClone(ctx context.Context, gitURL string, branch string, dest string, keepTs bool) error {
	c.nonSaveCommand()
	gitOpts := []llb.GitOption{
		llb.WithCustomNamef(
			"%sGIT CLONE (--branch %s) %s", c.vertexPrefixWithURL(gitURL), branch, gitURL),
		llb.KeepGitDir(),
	}
	gitState := llb.Git(gitURL, branch, gitOpts...)
	c.mts.Final.MainState = llbutil.CopyOp(
		gitState, []string{"."}, c.mts.Final.MainState, dest, false, false, keepTs,
		c.mts.Final.MainImage.Config.User, false,
		llb.WithCustomNamef(
			"%sCOPY GIT CLONE (--branch %s) %s TO %s", c.vertexPrefix(false),
			branch, gitURL, dest))
	return nil
}

// WithDockerRun applies an entire WITH DOCKER ... RUN ... END clause.
func (c *Converter) WithDockerRun(ctx context.Context, args []string, opt WithDockerOpt) error {
	c.nonSaveCommand()
	wdr := &withDockerRun{
		c: c,
	}
	return wdr.Run(ctx, args, opt)
}

// Healthcheck applies the HEALTHCHECK command.
func (c *Converter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int) {
	c.nonSaveCommand()
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
	c.mts.Final.MainImage.Config.Healthcheck = hc
}

// FinalizeStates returns the LLB states.
func (c *Converter) FinalizeStates(ctx context.Context) (*states.MultiTarget, error) {
	c.markFakeDeps()

	c.opt.BuildContextProvider.AddDirs(c.mts.Final.LocalDirs)
	c.mts.Final.VarCollection = c.varCollection
	c.mts.Final.Ongoing = false
	return c.mts, nil
}

// ExpandArgs expands args in the provided word.
func (c *Converter) ExpandArgs(word string) string {
	return c.varCollection.Expand(word)
}

func (c *Converter) buildTarget(ctx context.Context, fullTargetName string, platform *specs.Platform, buildArgs []string, isDangling bool) (*states.MultiTarget, error) {
	relTarget, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		return nil, errors.Wrapf(err, "earthly target parse %s", fullTargetName)
	}
	target, err := domain.JoinTargets(c.mts.Final.Target, relTarget)
	if err != nil {
		return nil, errors.Wrap(err, "join targets")
	}
	// Don't allow transitive overriding variables to cross project boundaries.
	propagateBuildArgs := !relTarget.IsExternal()
	var newVars map[string]bool
	newVarCollection, newVars, err := c.varCollection.WithParseBuildArgs(
		buildArgs, c.processNonConstantBuildArgFunc(ctx), propagateBuildArgs)
	if err != nil {
		return nil, errors.Wrap(err, "parse build args")
	}
	// Recursion.
	opt := c.opt
	opt.Visited = c.mts.Visited
	opt.VarCollection = newVarCollection
	opt.Platform, err = llbutil.ResolvePlatform(platform, c.opt.Platform)
	if err != nil {
		// Contradiction allowed. You can BUILD another target with different platform.
		opt.Platform = platform
	}
	mts, err := Earthfile2LLB(ctx, target, opt)
	if err != nil {
		return nil, errors.Wrapf(err, "earthfile2llb for %s", fullTargetName)
	}
	if isDangling {
		mts.Final.HasDangling = true
	}
	c.directDeps = append(c.directDeps, mts.Final)
	if propagateBuildArgs {
		// Propagate build arg inputs upwards (a child target depending on a build arg means
		// that the parent also depends on that build arg).
		for _, bai := range mts.Final.TargetInput.BuildArgs {
			// Check if the build arg has been overridden. If it has, it can no longer be an input
			// directly, so skip it.
			_, found := newVars[bai.Name]
			if found {
				continue
			}
			c.mts.Final.TargetInput = c.mts.Final.TargetInput.WithBuildArgInput(bai)
		}
		// Propagate globals.
		globals := mts.Final.VarCollection.WithOnlyGlobals()
		for _, k := range globals.SortedActiveVariables() {
			_, alreadyActive, _ := c.varCollection.Get(k)
			if alreadyActive {
				// Globals don't override any variables in current scope.
				continue
			}
			v, _, _ := globals.Get(k)
			c.varCollection.AddActive(k, v)
			c.mts.Final.TargetInput = c.mts.Final.TargetInput.WithBuildArgInput(
				v.BuildArgInput(k, "")) // TODO: Set correct default value for bai.
		}
	}
	return mts, nil
}

func (c *Converter) internalRun(ctx context.Context, args, secretKeyValues []string, isWithShell bool, shellWrap shellWrapFun, pushFlag, transient, withSSH, noCache bool, commandStr string, opts ...llb.RunOption) (llb.State, error) {
	finalOpts := opts
	var extraEnvVars []string
	// Secrets.
	for _, secretKeyValue := range secretKeyValues {
		parts := strings.SplitN(secretKeyValue, "=", 2)
		if len(parts) != 2 {
			return llb.State{}, fmt.Errorf("invalid secret definition %s", secretKeyValue)
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
			return llb.State{}, fmt.Errorf("secret definition %s not supported. Must start with +secrets/ or be an empty string", secretKeyValue)
		}
	}
	// Build args.
	for _, buildArgName := range c.varCollection.SortedActiveVariables() {
		ba, _, _ := c.varCollection.Get(buildArgName)
		if ba.IsEnvVar() {
			continue
		}
		extraEnvVars = append(extraEnvVars, fmt.Sprintf("%s=%s", buildArgName, shellescape.Quote(ba.ConstantValue())))
	}
	// Debugger.
	secretOpts := []llb.SecretOption{
		llb.SecretID(common.DebuggerSettingsSecretsKey),
		llb.SecretFileOpt(0, 0, 0444),
	}
	debuggerSecretMount := llb.AddSecret(
		fmt.Sprintf("/run/secrets/%s", common.DebuggerSettingsSecretsKey), secretOpts...)
	debuggerMount := llb.AddMount(debuggerPath, llb.Scratch(),
		llb.HostBind(), llb.SourcePath("/usr/bin/earth_debugger"))
	runEarthlyMount := llb.AddMount("/run/earthly", llb.Scratch(),
		llb.HostBind(), llb.SourcePath("/run/earthly"))
	finalOpts = append(finalOpts, debuggerSecretMount, debuggerMount, runEarthlyMount)
	if withSSH {
		finalOpts = append(finalOpts, llb.AddSSHSocket())
	}
	// Shell and debugger wrap.
	finalArgs := shellWrap(args, extraEnvVars, isWithShell, true)
	finalOpts = append(finalOpts, llb.Args(finalArgs))
	if noCache {
		finalOpts = append(finalOpts, llb.IgnoreCache)
	}

	if pushFlag {
		// For push-flagged commands, make sure they run every time - don't use cache.
		finalOpts = append(finalOpts, llb.IgnoreCache)
		if !c.mts.Final.RunPush.Initialized {
			// If this is the first push-flagged command, initialize the state with the latest
			// side-effects state.
			c.mts.Final.RunPush.State = c.mts.Final.MainState
			c.mts.Final.RunPush.Initialized = true
		}
		// Don't run on MainState. We want push-flagged commands to be executed only
		// *after* the build. Save this for later.
		c.mts.Final.RunPush.State = c.mts.Final.RunPush.State.Run(finalOpts...).Root()
		c.mts.Final.RunPush.CommandStrs = append(
			c.mts.Final.RunPush.CommandStrs, commandStr)
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
	ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, mts.Final.ArtifactsState, mts.Final.Platform, c.opt.CacheImports)
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

func (c *Converter) internalFromClassical(ctx context.Context, imageName string, platform specs.Platform, opts ...llb.ImageOption) (llb.State, *image.Image, *variables.Collection, error) {
	if imageName == "scratch" {
		// FROM scratch
		img := image.NewImage()
		img.OS = platform.OS
		img.Architecture = platform.Architecture
		return llb.Scratch().Platform(platform), img,
			c.varCollection.WithResetEnvVars(), nil
	}
	ref, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return llb.State{}, nil, nil, errors.Wrapf(err, "parse normalized named %s", imageName)
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
	allOpts := append(opts, llb.Platform(platform), c.opt.ImageResolveMode)
	state := llb.Image(ref.String(), allOpts...)
	state, img2, newVarCollection := c.applyFromImage(state, &img)
	return state, img2, newVarCollection, nil
}

func (c *Converter) applyFromImage(state llb.State, img *image.Image) (llb.State, *image.Image, *variables.Collection) {
	// Reset variables.
	newVarCollection := c.varCollection.WithResetEnvVars()
	for _, envVar := range img.Config.Env {
		k, v, _ := variables.ParseKeyValue(envVar)
		newVarCollection.AddActive(k, variables.NewConstantEnvVar(v))
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
	return state, img, newVarCollection
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
			llb.Mkdir(srcBuildArgDir, 0755, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] mkdir %s", srcBuildArgDir))
		args := strings.Split(fmt.Sprintf("echo \"%s\" >%s", expression, srcBuildArgPath), " ")
		transient := true
		state, err := c.internalRun(
			ctx, args, []string{}, true, withShellAndEnvVars, false, transient, false, false, expression,
			llb.WithCustomNamef("%sRUN %s", c.vertexPrefix(false), expression))
		if err != nil {
			return "", 0, errors.Wrapf(err, "run %v", expression)
		}
		ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, state, c.opt.Platform, c.opt.CacheImports)
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

func (c *Converter) vertexPrefix(local bool) string {
	overriding := c.varCollection.SortedOverridingVariables()
	varStrBuilder := make([]string, 0, len(overriding)+1)
	if c.opt.Platform != nil {
		varStrBuilder = append(
			varStrBuilder,
			fmt.Sprintf("platform=%s", llbutil.PlatformToString(c.opt.Platform)))
	}
	for _, key := range overriding {
		variable, _, _ := c.varCollection.Get(key)
		if variable.IsEnvVar() {
			continue
		}
		varStrBuilder = append(varStrBuilder, fmt.Sprintf("%s=%s", key, variable.ConstantValue()))
	}
	var varStr string
	if len(varStrBuilder) > 0 {
		b64VarStr := base64.StdEncoding.EncodeToString([]byte(strings.Join(varStrBuilder, " ")))
		varStr = fmt.Sprintf("(%s)", b64VarStr)
	}
	return fmt.Sprintf("[%s%s%s %s] ", c.mts.Final.Target.String(), varStr, strIf(local, " *local*"), c.mts.Final.Salt)
}

func (c *Converter) markFakeDeps() {
	if !c.opt.UseFakeDep {
		return
	}
	for _, dep := range c.directDeps {
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
	c.mts.Final.Platform = platform
	c.varCollection.SetPlatformArgs(llbutil.PlatformWithDefault(platform))
}

func makeCacheContext(target domain.Target) llb.State {
	sessionID := cacheKey(target)
	opts := []llb.LocalOption{
		llb.SharedKeyHint(target.ProjectCanonical()),
		llb.SessionID(sessionID),
		llb.Platform(llbutil.DefaultPlatform()),
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
