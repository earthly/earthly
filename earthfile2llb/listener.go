package earthfile2llb

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/parser"
	"github.com/pkg/errors"
)

var _ parser.EarthParserListener = &listener{}

type listener struct {
	*parser.BaseEarthParserListener
	converter *Converter
	ctx       context.Context

	executeTarget   string
	currentTarget   string
	targetFound     bool
	pushOnlyAllowed bool

	envArgKey   string
	envArgValue string
	labelKeys   []string
	labelValues []string

	withDocker    *WithDockerOpt
	withDockerRan bool

	execMode  bool
	stmtWords []string

	err error
}

func newListener(ctx context.Context, converter *Converter, executeTarget string) *listener {
	return &listener{
		ctx:           ctx,
		converter:     converter,
		executeTarget: executeTarget,
		currentTarget: "base",
		targetFound:   (executeTarget == "base"),
	}
}

func (l *listener) Err() error {
	if l.err != nil {
		return l.err
	}
	if !l.targetFound {
		return fmt.Errorf("target %s not defined", l.executeTarget)
	}
	return nil
}

func (l *listener) EnterTargetHeader(c *parser.TargetHeaderContext) {
	l.currentTarget = strings.TrimSuffix(c.GetText(), ":")
	if l.currentTarget == l.executeTarget {
		if l.targetFound {
			l.err = fmt.Errorf("target %s is declared twice", l.currentTarget)
			return
		}
		l.targetFound = true
	}
	if l.shouldSkip() {
		return
	}
	if l.currentTarget == "base" || l.currentTarget == "secrets" {
		l.err = errors.New("target name cannot be \"base\" or \"secrets\"")
		return
	}
	// Apply implicit FROM +base
	err := l.converter.From(l.ctx, "+base", nil)
	if err != nil {
		l.err = errors.Wrap(err, "apply implicit FROM +base")
		return
	}
}

func (l *listener) EnterStmts(c *parser.StmtsContext) {
	if l.shouldSkip() {
		return
	}
	l.pushOnlyAllowed = false
}

func (l *listener) ExitStmts(c *parser.StmtsContext) {
	if l.shouldSkip() {
		return
	}
	if l.withDocker != nil {
		l.err = errors.New("no matching END found for WITH DOCKER")
		return
	}
}

//
// Commands.

func (l *listener) EnterStmt(c *parser.StmtContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = nil
	l.envArgKey = ""
	l.envArgValue = ""
	l.labelKeys = nil
	l.labelValues = nil
	l.execMode = false
}

func (l *listener) ExitFromStmt(c *parser.FromStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("FROM", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid FROM arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 1 {
		if fs.NArg() == 3 && fs.Arg(1) == "AS" {
			l.err = errors.New("AS not supported, use earthly targets instead")
		} else {
			l.err = fmt.Errorf("invalid number of arguments for FROM: %s", l.stmtWords)
		}
		return
	}
	imageName := l.expandArgs(fs.Arg(0), true)
	for i, ba := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(ba, true)
	}
	err = l.converter.From(l.ctx, imageName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply FROM %s", imageName)
		return
	}
}

func (l *listener) ExitFromDockerfileStmt(c *parser.FromDockerfileStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("FROM DOCKERFILE", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target and also to the Dockerfile build")
	dfTarget := fs.String("target", "", "The Dockerfile target to inherit from")
	dfPath := fs.String("f", "", "Not supported")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid FROM DOCKERFILE arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 1 {
		l.err = errors.New("invalid number of arguments for FROM DOCKERFILE")
		return
	}
	path := l.expandArgs(fs.Arg(0), false)
	_, parseErr := domain.ParseArtifact(path)
	if parseErr != nil {
		// Treat as context path, not artifact path.
		path = l.expandArgs(fs.Arg(0), false)
	}
	for i, ba := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(ba, true)
	}
	*dfPath = l.expandArgs(*dfPath, false)
	*dfTarget = l.expandArgs(*dfTarget, false)
	err = l.converter.FromDockerfile(l.ctx, path, *dfPath, *dfTarget, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrap(err, "from dockerfile")
		return
	}
}

func (l *listener) ExitCopyStmt(c *parser.CopyStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("COPY", flag.ContinueOnError)
	from := fs.String("from", "", "Not supported")
	isDirCopy := fs.Bool("dir", false, "Copy entire directories, not just the contents")
	chown := fs.String("chown", "", "Apply a specific group and/or owner to the copied files and directories")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid COPY arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() < 2 {
		l.err = fmt.Errorf("not enough COPY arguments %v", l.stmtWords)
		return
	}
	if *from != "" {
		l.err = errors.New("COPY --from not implemented. Use COPY artifacts form instead")
		return
	}
	srcs := fs.Args()[:fs.NArg()-1]
	dest := l.expandArgs(fs.Arg(fs.NArg()-1), false)
	for i, ba := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(ba, true)
	}
	*chown = l.expandArgs(*chown, false)
	allClassical := true
	allArtifacts := true
	for i, src := range srcs {
		// If it parses as an artifact, treat as artifact.
		artifactSrc, parseErr := domain.ParseArtifact(l.expandArgs(src, true))
		if parseErr == nil {
			srcs[i] = artifactSrc.String()
			allClassical = false
		} else {
			srcs[i] = l.expandArgs(src, false)
			allArtifacts = false
		}
	}
	if !allClassical && !allArtifacts {
		l.err = fmt.Errorf("combining artifacts and build context arguments in a single COPY command is not allowed: %v", srcs)
		return
	}
	if allArtifacts {
		for _, src := range srcs {
			err = l.converter.CopyArtifact(l.ctx, src, dest, buildArgs.Args, *isDirCopy, *chown)
			if err != nil {
				l.err = errors.Wrapf(err, "copy artifact")
				return
			}
		}
	} else {
		if len(buildArgs.Args) != 0 {
			l.err = fmt.Errorf("build args not supported for non +artifact arguments case %v", l.stmtWords)
			return
		}
		l.converter.CopyClassical(l.ctx, srcs, dest, *isDirCopy, *chown)
	}
}

func (l *listener) ExitRunStmt(c *parser.RunStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) < 1 {
		l.err = errors.New("not enough arguments for RUN")
		return
	}

	fs := flag.NewFlagSet("RUN", flag.ContinueOnError)
	pushFlag := fs.Bool(
		"push", false,
		"Execute this command only if the build succeeds and also if earth is invoked in push mode")
	privileged := fs.Bool("privileged", false, "Enable privileged mode")
	withEntrypoint := fs.Bool(
		"entrypoint", false,
		"Include the entrypoint of the image when running the command")
	withDocker := fs.Bool("with-docker", false, "Deprecated")
	withSSH := fs.Bool("ssh", false, "Make available the SSH agent of the host")
	secrets := new(StringSliceFlag)
	fs.Var(secrets, "secret", "Make available a secret")
	mounts := new(StringSliceFlag)
	fs.Var(mounts, "mount", "Mount a file or directory")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid RUN arguments %v", l.stmtWords)
		return
	}
	withShell := !l.execMode
	if *withDocker {
		*privileged = true
	}
	if !*pushFlag && l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	// TODO: In the bracket case, should flags be outside of the brackets?

	for i, s := range secrets.Args {
		secrets.Args[i] = l.expandArgs(s, true)
	}
	for i, m := range mounts.Args {
		mounts.Args[i] = l.expandArgs(m, false)
	}
	// Note: Not expanding args for the run itself, as that will be take care of by the shell.

	if l.withDocker == nil {
		err = l.converter.Run(
			l.ctx, fs.Args(), mounts.Args, secrets.Args, *privileged, *withEntrypoint, *withDocker,
			withShell, *pushFlag, *withSSH)
		if err != nil {
			l.err = errors.Wrap(err, "run")
			return
		}
		if *pushFlag {
			l.pushOnlyAllowed = true
		}
	} else {
		if *pushFlag {
			l.err = fmt.Errorf("RUN --push not allowed in WITH DOCKER")
			return
		}
		if l.withDockerRan {
			l.err = fmt.Errorf("Only one RUN command allowed in WITH DOCKER")
			return
		}
		l.withDockerRan = true
		l.withDocker.Mounts = mounts.Args
		l.withDocker.Secrets = secrets.Args
		l.withDocker.WithShell = withShell
		l.withDocker.WithEntrypoint = *withEntrypoint
		err = l.converter.WithDockerRun(l.ctx, fs.Args(), *l.withDocker)
		if err != nil {
			l.err = errors.Wrap(err, "with docker run")
			return
		}
	}
}

func (l *listener) ExitSaveArtifact(c *parser.SaveArtifactContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the SAVE ARTIFACT command")
		return
	}
	if len(l.stmtWords) > 5 {
		l.err = fmt.Errorf("too many arguments provided to the SAVE ARTIFACT command: %v", l.stmtWords)
		return
	}
	saveAsLocalTo := ""
	saveTo := "./"
	if len(l.stmtWords) >= 4 {
		if strings.Join(l.stmtWords[len(l.stmtWords)-3:len(l.stmtWords)-1], " ") == "AS LOCAL" {
			saveAsLocalTo = l.stmtWords[len(l.stmtWords)-1]
			if len(l.stmtWords) == 5 {
				saveTo = l.stmtWords[1]
			}
		} else {
			l.err = fmt.Errorf("invalid arguments for SAVE ARTIFACT command: %v", l.stmtWords)
			return
		}
	} else if len(l.stmtWords) == 2 {
		saveTo = l.stmtWords[1]
	} else if len(l.stmtWords) == 3 {
		l.err = fmt.Errorf("invalid arguments for SAVE ARTIFACT command: %v", l.stmtWords)
		return
	}

	saveFrom := l.expandArgs(l.stmtWords[0], false)
	saveTo = l.expandArgs(saveTo, false)
	saveAsLocalTo = l.expandArgs(saveAsLocalTo, false)
	err := l.converter.SaveArtifact(l.ctx, saveFrom, saveTo, saveAsLocalTo)
	if err != nil {
		l.err = errors.Wrap(err, "apply SAVE ARTIFACT")
		return
	}
}

func (l *listener) ExitSaveImage(c *parser.SaveImageContext) {
	if l.shouldSkip() {
		return
	}

	fs := flag.NewFlagSet("SAVE IMAGE", flag.ContinueOnError)
	pushFlag := fs.Bool(
		"push", false,
		"Push the image to the remote registry provided that the build succeeds and also that earth is invoked in push mode")
	cacheHint := fs.Bool(
		"cache-hint", false,
		"Instruct Earthly that the current target shuold be saved entirely as part of the remote cache")
	cacheFrom := new(StringSliceFlag)
	fs.Var(cacheFrom, "cache-from", "Declare additional cache import as a Docker tag")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid SAVE IMAGE arguments %v", l.stmtWords)
		return
	}
	for i, cf := range cacheFrom.Args {
		cacheFrom.Args[i] = l.expandArgs(cf, false)
	}
	if !*pushFlag && l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if *pushFlag && fs.NArg() == 0 {
		l.err = fmt.Errorf("invalid number of arguments for SAVE IMAGE --push: %v", l.stmtWords)
		return
	}

	imageNames := fs.Args()
	for i, img := range imageNames {
		imageNames[i] = l.expandArgs(img, false)
	}
	if len(imageNames) == 0 && !*cacheHint && len(cacheFrom.Args) == 0 {
		fmt.Printf("Deprecation: using SAVE IMAGE with no arguments is no longer necessary and can be safely removed\n")
		return
	}
	err = l.converter.SaveImage(l.ctx, imageNames, *pushFlag, *cacheHint, cacheFrom.Args)
	if err != nil {
		l.err = errors.Wrap(err, "save image")
		return
	}
	if *pushFlag {
		l.pushOnlyAllowed = true
	}
}

func (l *listener) ExitBuildStmt(c *parser.BuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("BUILD", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid BUILD arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 1 {
		l.err = fmt.Errorf("invalid number of arguments for BUILD: %s", l.stmtWords)
		return
	}
	fullTargetName := l.expandArgs(fs.Arg(0), true)
	for i, arg := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(arg, true)
	}
	err = l.converter.Build(l.ctx, fullTargetName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply BUILD %s", fullTargetName)
		return
	}
}

func (l *listener) ExitWorkdirStmt(c *parser.WorkdirStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for WORKDIR: %v", l.stmtWords)
		return
	}
	workdirPath := l.expandArgs(l.stmtWords[0], false)
	l.converter.Workdir(l.ctx, workdirPath)
}

func (l *listener) ExitUserStmt(c *parser.UserStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for USER: %v", l.stmtWords)
		return
	}
	user := l.expandArgs(l.stmtWords[0], false)
	l.converter.User(l.ctx, user)
}

func (l *listener) ExitCmdStmt(c *parser.CmdStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	withShell := !l.execMode
	cmdArgs := l.stmtWords
	if !withShell {
		for i, arg := range cmdArgs {
			cmdArgs[i] = l.expandArgs(arg, false)
		}
	}
	l.converter.Cmd(l.ctx, cmdArgs, withShell)
}

func (l *listener) ExitEntrypointStmt(c *parser.EntrypointStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	withShell := !l.execMode
	entArgs := l.stmtWords
	if !withShell {
		for i, arg := range entArgs {
			entArgs[i] = l.expandArgs(arg, false)
		}
	}
	l.converter.Entrypoint(l.ctx, entArgs, withShell)
}

func (l *listener) ExitExposeStmt(c *parser.ExposeStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the EXPOSE command")
		return
	}
	ports := l.stmtWords
	for i, port := range ports {
		ports[i] = l.expandArgs(port, false)
	}
	l.converter.Expose(l.ctx, ports)
}

func (l *listener) ExitVolumeStmt(c *parser.VolumeStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the VOLUME command")
		return
	}
	volumes := l.stmtWords
	for i, volume := range volumes {
		volumes[i] = l.expandArgs(volume, false)
	}
	l.converter.Volume(l.ctx, volumes)
}

func (l *listener) ExitEnvStmt(c *parser.EnvStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	key := l.envArgKey // Note: Not expanding args for key.
	value := l.expandArgs(l.envArgValue, false)
	l.converter.Env(l.ctx, key, value)
}

func (l *listener) ExitArgStmt(c *parser.ArgStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	key := l.envArgKey // Note: Not expanding args for key.
	value := l.expandArgs(l.envArgValue, true)
	// Args declared in the base target are global.
	global := (l.currentTarget == "base")
	l.converter.Arg(l.ctx, key, value, global)
}

func (l *listener) ExitLabelStmt(c *parser.LabelStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.labelKeys) == 0 {
		l.err = fmt.Errorf("no labels provided in LABEL command: %s", c.GetText())
		return
	}
	if len(l.labelKeys) != len(l.labelValues) {
		l.err = fmt.Errorf("label keys and values do not match: %s", c.GetText())
		return
	}
	labels := make(map[string]string)
	for i := range l.labelKeys {
		labels[l.expandArgs(l.labelKeys[i], false)] = l.expandArgs(l.labelValues[i], false)
	}
	l.converter.Label(l.ctx, labels)
}

func (l *listener) ExitGitCloneStmt(c *parser.GitCloneStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("GIT CLONE", flag.ContinueOnError)
	branch := fs.String("branch", "", "The git ref to use when cloning")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid GIT CLONE arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 2 {
		l.err = fmt.Errorf("invalid number of arguments for GIT CLONE: %s", l.stmtWords)
		return
	}
	gitURL := l.expandArgs(fs.Arg(0), false)
	gitCloneDest := l.expandArgs(fs.Arg(1), false)
	*branch = l.expandArgs(*branch, false)
	err = l.converter.GitClone(l.ctx, gitURL, *branch, gitCloneDest)
	if err != nil {
		l.err = errors.Wrap(err, "git clone")
		return
	}
}

func (l *listener) ExitDockerLoadStmt(c *parser.DockerLoadStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("DOCKER LOAD", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid DOCKER LOAD arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 2 {
		l.err = fmt.Errorf("invalid number of arguments for DOCKER LOAD: %s", l.stmtWords)
		return
	}
	fullTargetName := l.expandArgs(fs.Arg(0), true)
	imageName := l.expandArgs(fs.Arg(1), false)
	for i, arg := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(arg, true)
	}
	if l.withDocker == nil {
		err = l.converter.DockerLoadOld(l.ctx, fullTargetName, imageName, buildArgs.Args)
		if err != nil {
			l.err = errors.Wrap(err, "docker load")
			return
		}
	} else {
		if l.withDockerRan {
			l.err = fmt.Errorf("cannot DOCKER LOAD after the RUN command in a WITH DOCKER clause")
			return
		}
		fmt.Printf("Warning: DOCKER LOAD is deprecated. Please use WITH DOCKER --load %s=%s instead\n", imageName, fullTargetName)
		l.withDocker.Loads = append(l.withDocker.Loads, DockerLoadOpt{
			Target:    fullTargetName,
			ImageName: imageName,
			BuildArgs: buildArgs.Args,
		})
	}
}

func (l *listener) ExitDockerPullStmt(c *parser.DockerPullStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for DOCKER PULL: %s", l.stmtWords)
		return
	}
	imageName := l.expandArgs(l.stmtWords[0], false)
	if l.withDocker == nil {
		err := l.converter.DockerPullOld(l.ctx, imageName)
		if err != nil {
			l.err = errors.Wrap(err, "docker pull")
			return
		}
	} else {
		if l.withDockerRan {
			l.err = fmt.Errorf("cannot DOCKER PULL after the RUN command in a WITH DOCKER clause")
			return
		}
		fmt.Printf("Warning: DOCKER PULL is deprecated. Please use WITH DOCKER --pull %s instead\n", imageName)
		l.withDocker.Pulls = append(l.withDocker.Pulls, imageName)
	}
}

func (l *listener) ExitHealthcheckStmt(c *parser.HealthcheckStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("HEALTHCHECK", flag.ContinueOnError)
	interval := fs.Duration(
		"interval", 30*time.Second,
		"The interval between healthchecks")
	timeout := fs.Duration(
		"timeout", 30*time.Second,
		"The timeout before the command is considered failed")
	startPeriod := fs.Duration(
		"start-period", 0,
		"An initialization time period in which failures are not counted towards the maximum number of retries")
	retries := fs.Int(
		"retries", 3,
		"The number of retries before a container is considered unhealthy")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid HEALTHCHECK arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() == 0 {
		l.err = fmt.Errorf("invalid number of arguments for HEALTHCHECK: %s", l.stmtWords)
		return
	}
	isNone := false
	var cmdArgs []string
	switch fs.Arg(0) {
	case "NONE":
		if fs.NArg() != 1 {
			l.err = fmt.Errorf("invalid arguments for HEALTHCHECK: %s", l.stmtWords)
			return
		}
		isNone = true
	case "CMD":
		if fs.NArg() == 1 {
			l.err = fmt.Errorf("invalid number of arguments for HEALTHCHECK CMD: %s", l.stmtWords)
			return
		}
		cmdArgs = fs.Args()[1:]
	default:
		if strings.HasPrefix(fs.Arg(0), "[") {
			l.err = fmt.Errorf("exec form not yet supported for HEALTHCHECK CMD: %s", l.stmtWords)
			return
		}
		l.err = fmt.Errorf("invalid arguments for HEALTHCHECK: %s", l.stmtWords)
		return
	}
	for i, arg := range cmdArgs {
		cmdArgs[i] = l.expandArgs(arg, false)
	}
	l.converter.Healthcheck(l.ctx, isNone, cmdArgs, *interval, *timeout, *startPeriod, *retries)
}

func (l *listener) ExitWithDockerStmt(c *parser.WithDockerStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if l.withDocker != nil {
		l.err = fmt.Errorf("cannot use WITH DOCKER within WITH DOCKER")
		return
	}

	fs := flag.NewFlagSet("WITH DOCKER", flag.ContinueOnError)
	composeFiles := new(StringSliceFlag)
	fs.Var(composeFiles, "compose", "A compose file used to bring up services from")
	composeServices := new(StringSliceFlag)
	fs.Var(composeServices, "service", "A compose service to bring up")
	loads := new(StringSliceFlag)
	fs.Var(loads, "load", "An image produced by Earthly which is loaded as a Docker image")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	pulls := new(StringSliceFlag)
	fs.Var(pulls, "pull", "An image which is pulled and made available in the docker cache")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid WITH DOCKER arguments %v", l.stmtWords)
		return
	}
	if len(fs.Args()) != 0 {
		l.err = fmt.Errorf("invalid WITH DOCKER arguments %v", fs.Args())
		return
	}

	for i, cf := range composeFiles.Args {
		composeFiles.Args[i] = l.expandArgs(cf, false)
	}
	for i, cs := range composeServices.Args {
		composeServices.Args[i] = l.expandArgs(cs, false)
	}
	for i, load := range loads.Args {
		loads.Args[i] = l.expandArgs(load, true)
	}
	for i, ba := range buildArgs.Args {
		buildArgs.Args[i] = l.expandArgs(ba, true)
	}
	for i, p := range pulls.Args {
		pulls.Args[i] = l.expandArgs(p, false)
	}

	l.withDocker = &WithDockerOpt{
		Pulls:           pulls.Args,
		ComposeFiles:    composeFiles.Args,
		ComposeServices: composeServices.Args,
	}
	for _, loadStr := range loads.Args {
		loadImg, loadTarget, err := parseLoad(loadStr)
		if err != nil {
			l.err = err
			return
		}
		l.withDocker.Loads = append(l.withDocker.Loads, DockerLoadOpt{
			Target:    loadTarget,
			ImageName: loadImg,
			BuildArgs: buildArgs.Args,
		})
	}
}

func (l *listener) ExitEndStmt(c *parser.EndStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) != 0 {
		l.err = fmt.Errorf("END does not take any arguments: %s", c.GetText())
		return
	}
	if l.withDocker == nil {
		l.err = fmt.Errorf("END can only be used to end a WITH DOCKER clause")
		return
	}
	if !l.withDockerRan {
		l.err = fmt.Errorf("No RUN command found in WITH DOCKER")
		return
	}
	l.withDocker = nil
	l.withDockerRan = false
}

func (l *listener) ExitAddStmt(c *parser.AddStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command ADD not yet supported")
}

func (l *listener) ExitStopsignalStmt(c *parser.StopsignalStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command STOPSIGNAL not yet supported")
}

func (l *listener) ExitOnbuildStmt(c *parser.OnbuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command ONBUILD not supported")
}

func (l *listener) ExitShellStmt(c *parser.ShellStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command SHELL not yet supported")
}

func (l *listener) ExitGenericCommandStmt(c *parser.GenericCommandStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Invalid command %s", c.GetText())
}

//
// Variables.

func (l *listener) EnterEnvArgKey(c *parser.EnvArgKeyContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgKey = c.GetText()
	err := checkEnvVarName(l.envArgKey)
	if err != nil {
		l.err = err
		return
	}
}

func (l *listener) EnterEnvArgValue(c *parser.EnvArgValueContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgValue = c.GetText()
}

func (l *listener) EnterLabelKey(c *parser.LabelKeyContext) {
	if l.shouldSkip() {
		return
	}
	l.labelKeys = append(l.labelKeys, c.GetText())
}

func (l *listener) EnterLabelValue(c *parser.LabelValueContext) {
	if l.shouldSkip() {
		return
	}
	l.labelValues = append(l.labelValues, c.GetText())
}

func (l *listener) ExitStmtWordsMaybeJSON(c *parser.StmtWordsMaybeJSONContext) {
	if l.shouldSkip() {
		return
	}
	// Try to parse as JSON. If parse works, override the already collected stmtWords.
	var words []string
	err := json.Unmarshal([]byte(c.GetText()), &words)
	if err == nil {
		l.stmtWords = words
		l.execMode = true
	}
}

func (l *listener) EnterStmtWord(c *parser.StmtWordContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = append(l.stmtWords, replaceEscape(c.GetText()))
}

func (l *listener) shouldSkip() bool {
	return l.err != nil || l.currentTarget != l.executeTarget
}

func (l *listener) expandArgs(word string, keepPlusEscape bool) string {
	ret := l.converter.ExpandArgs(escapeSlashPlus(word))
	if keepPlusEscape {
		return ret
	}
	return unescapeSlashPlus(ret)
}

// StringSliceFlag is a flag backed by a string slice.
type StringSliceFlag struct {
	Args []string
}

// String returns a string representation of the flag.
func (ssf *StringSliceFlag) String() string {
	if ssf == nil {
		return ""
	}
	return strings.Join(ssf.Args, ",")
}

// Set adds a flag value to the string slice.
func (ssf *StringSliceFlag) Set(arg string) error {
	ssf.Args = append(ssf.Args, arg)
	return nil
}

var envVarNameRegexp = regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")

func checkEnvVarName(str string) error {
	itMatch := envVarNameRegexp.MatchString(str)
	if !itMatch {
		return fmt.Errorf("invalid env key definition %s", str)
	}
	return nil
}

var lineContinuationRegexp = regexp.MustCompile("\\\\(\\n|(\\r\\n))[\\t ]*")

func replaceEscape(str string) string {
	return lineContinuationRegexp.ReplaceAllString(str, "")
}

func parseLoad(loadStr string) (string, string, error) {
	splitLoad := strings.SplitN(loadStr, "=", 2)
	if len(splitLoad) < 2 {
		// --load <target-name>
		// (will infer image name from SAVE IMAGE of that target)
		return "", loadStr, nil
	}
	// --load <image-name>=<target-name>
	return splitLoad[0], splitLoad[1], nil
}

func escapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "\\\\+")
}

func unescapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "+")
}
