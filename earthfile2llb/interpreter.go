package earthfile2llb

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Interpreter interprets Earthly AST's into calls to the converter.
type Interpreter struct {
	converter *Converter

	isBase          bool
	isWith          bool
	pushOnlyAllowed bool
	local           bool

	withDocker    *WithDockerOpt
	withDockerRan bool
}

func newInterpreter(c *Converter) *Interpreter {
	return &Interpreter{
		converter: c,
	}
}

// Run interprets the commands in the given Earthfile AST, for a specific target.
func (i *Interpreter) Run(ctx context.Context, ef spec.Earthfile, target string) error {
	if target == "base" {
		i.isBase = true
		err := i.handleBlock(ctx, ef.BaseRecipe)
		i.isBase = false
		return err
	}
	for _, t := range ef.Targets {
		if t.Name == target {
			return i.handleTarget(ctx, t)
		}
	}
	return Errorf(ef.SourceLocation, "target %s not found", target)
}

func (i *Interpreter) handleTarget(ctx context.Context, t spec.Target) error {
	// Apply implicit FROM +base
	err := i.converter.From(ctx, "+base", nil, nil)
	if err != nil {
		return WrapError(err, t.SourceLocation, "apply FROM")
	}
	return i.handleBlock(ctx, t.Recipe)
}

func (i *Interpreter) handleBlock(ctx context.Context, b spec.Block) error {
	for _, stmt := range b {
		err := i.handleStatement(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) handleStatement(ctx context.Context, stmt spec.Statement) error {
	if stmt.Command != nil {
		return i.handleCommand(ctx, *stmt.Command)
	} else if stmt.With != nil {
		return i.handleWith(ctx, *stmt.With)
	} else if stmt.If != nil {
		return i.handleIf(ctx, *stmt.If)
	} else {
		return Errorf(stmt.SourceLocation, "unexpected statement type")
	}
}

func (i *Interpreter) handleCommand(ctx context.Context, cmd spec.Command) error {
	if i.isWith {
		switch cmd.Name {
		case "DOCKER":
			return i.handleWithDocker(ctx, cmd)
		default:
			return Errorf(cmd.SourceLocation, "unexpected WITH command %s", cmd.Name)
		}
	}

	switch cmd.Name {
	case "FROM":
		return i.handleFrom(ctx, cmd)
	case "RUN":
		return i.handleRun(ctx, cmd)
	case "FROM DOCKERFILE":
		return i.handleFromDockerfile(ctx, cmd)
	case "LOCALLY":
		return i.handleLocally(ctx, cmd)
	case "COPY":
		return i.handleCopy(ctx, cmd)
	case "SAVE ARTIFACT":
		return i.handleSaveArtifact(ctx, cmd)
	case "SAVE IMAGE":
		return i.handleSaveImage(ctx, cmd)
	case "BUILD":
		return i.handleBuild(ctx, cmd)
	case "WORKDIR":
		return i.handleWorkdir(ctx, cmd)
	case "USER":
		return i.handleUser(ctx, cmd)
	case "CMD":
		return i.handleCmd(ctx, cmd)
	case "ENTRYPOINT":
		return i.handleEntrypoint(ctx, cmd)
	case "EXPOSE":
		return i.handleExpose(ctx, cmd)
	case "VOLUME":
		return i.handleVolume(ctx, cmd)
	case "ENV":
		return i.handleEnv(ctx, cmd)
	case "ARG":
		return i.handleArg(ctx, cmd)
	case "LABEL":
		return i.handleLabel(ctx, cmd)
	case "GIT CLONE":
		return i.handleGitClone(ctx, cmd)
	case "HEALTHCHECK":
		return i.handleHealthcheck(ctx, cmd)
	case "ADD":
		return i.handleAdd(ctx, cmd)
	case "STOPSIGNAL":
		return i.handleStopsignal(ctx, cmd)
	case "ONBUILD":
		return i.handleOnbuild(ctx, cmd)
	case "SHELL":
		return i.handleShell(ctx, cmd)
	default:
		return Errorf(cmd.SourceLocation, "unexpected command %s", cmd.Name)
	}
}

func (i *Interpreter) handleWith(ctx context.Context, with spec.WithStatement) error {
	i.isWith = true
	err := i.handleCommand(ctx, with.Command)
	i.isWith = false
	if err != nil {
		return err
	}
	if i.withDocker != nil && len(with.Body) > 1 {
		return Errorf(with.SourceLocation, "only one command is allowed in WITH DOCKER and it has to be RUN")
	}
	err = i.handleBlock(ctx, with.Body)
	if err != nil {
		return err
	}
	i.withDocker = nil
	if !i.withDockerRan {
		return Errorf(with.SourceLocation, "no RUN command found in WITH DOCKER")
	}
	i.withDockerRan = false
	return nil
}

func (i *Interpreter) handleIf(ctx context.Context, ifStmt spec.IfStatement) error {
	if i.pushOnlyAllowed {
		return Errorf(ifStmt.SourceLocation, "no non-push commands allowed after a --push")
	}
	isZero, err := i.handleIfExpression(ctx, ifStmt.Expression, ifStmt.ExecMode, ifStmt.SourceLocation)
	if err != nil {
		return err
	}

	if isZero {
		return i.handleBlock(ctx, ifStmt.IfBody)
	}
	for _, elseIf := range ifStmt.ElseIf {
		isZero, err = i.handleIfExpression(ctx, elseIf.Expression, elseIf.ExecMode, elseIf.SourceLocation)
		if err != nil {
			return err
		}
		if isZero {
			return i.handleBlock(ctx, elseIf.Body)
		}
	}
	if ifStmt.ElseBody != nil {
		return i.handleBlock(ctx, *ifStmt.ElseBody)
	}
	return nil
}

func (i *Interpreter) handleIfExpression(ctx context.Context, expression []string, execMode bool, sl *spec.SourceLocation) (bool, error) {
	if len(expression) < 1 {
		return false, Errorf(sl, "not enough arguments for IF")
	}

	fs := flag.NewFlagSet("IF", flag.ContinueOnError)
	privileged := fs.Bool("privileged", false, "Enable privileged mode")
	withSSH := fs.Bool("ssh", false, "Make available the SSH agent of the host")
	noCache := fs.Bool("no-cache", false, "Always execute this specific condition, ignoring cache")
	secrets := new(StringSliceFlag)
	fs.Var(secrets, "secret", "Make available a secret")
	mounts := new(StringSliceFlag)
	fs.Var(mounts, "mount", "Mount a file or directory")
	err := fs.Parse(expression)
	if err != nil {
		return false, WrapError(err, sl, "invalid RUN arguments %v", expression)
	}
	withShell := !execMode

	for index, s := range secrets.Args {
		secrets.Args[index] = i.expandArgs(s, true)
	}
	for index, m := range mounts.Args {
		mounts.Args[index] = i.expandArgs(m, false)
	}
	// Note: Not expanding args for the expression itself, as that will be take care of by the shell.

	var exitCode int
	commandName := "IF"
	if i.local {
		if len(mounts.Args) > 0 {
			return false, Errorf(sl, "mounts are not supported in combination with the LOCALLY directive")
		}
		if *withSSH {
			return false, Errorf(sl, "the --ssh flag has no effect when used with the  LOCALLY directive")
		}
		if *privileged {
			return false, Errorf(sl, "the --privileged flag has no effect when used with the LOCALLY directive")
		}
		if *noCache {
			return false, Errorf(sl, "the --no-cache flag has no effect when used with the LOCALLY directive")
		}

		// TODO these should be supported, but haven't yet been implemented
		if len(secrets.Args) > 0 {
			return false, Errorf(sl, "secrets need to be implemented for the LOCALLY directive")
		}

		exitCode, err = i.converter.RunLocalExitCode(ctx, commandName, fs.Args())
		if err != nil {
			return false, WrapError(err, sl, "apply RUN")
		}
	} else {
		exitCode, err = i.converter.RunExitCode(
			ctx, commandName, fs.Args(), mounts.Args, secrets.Args, *privileged,
			withShell, *withSSH, *noCache)
		if err != nil {
			return false, WrapError(err, sl, "apply IF")
		}
	}
	return (exitCode == 0), nil
}

// Commands -------------------------------------------------------------------

func (i *Interpreter) handleFrom(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	fs := flag.NewFlagSet("FROM", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	platformStr := fs.String("platform", "", "The platform to use")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid FROM arguments")
	}
	if fs.NArg() != 1 {
		if fs.NArg() == 3 && fs.Arg(1) == "AS" {
			return Errorf(cmd.SourceLocation, "AS not supported, use earthly targets instead")
		}
		return Errorf(cmd.SourceLocation, "invalid number of arguments for FROM: %s", cmd.Args)
	}
	imageName := i.expandArgs(fs.Arg(0), true)
	*platformStr = i.expandArgs(*platformStr, false)
	platform, err := llbutil.ParsePlatform(*platformStr)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "parse platform %s", *platformStr)
	}
	for index, ba := range buildArgs.Args {
		buildArgs.Args[index] = i.expandArgs(ba, true)
	}

	i.local = false
	err = i.converter.From(ctx, imageName, platform, buildArgs.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "apply FROM %s", imageName)
	}
	return nil
}

func (i *Interpreter) handleRun(ctx context.Context, cmd spec.Command) error {
	if len(cmd.Args) < 1 {
		return Errorf(cmd.SourceLocation, "not enough arguments for RUN")
	}

	fs := flag.NewFlagSet("RUN", flag.ContinueOnError)
	pushFlag := fs.Bool(
		"push", false,
		"Execute this command only if the build succeeds and also if earthly is invoked in push mode")
	privileged := fs.Bool("privileged", false, "Enable privileged mode")
	withEntrypoint := fs.Bool(
		"entrypoint", false,
		"Include the entrypoint of the image when running the command")
	withDocker := fs.Bool("with-docker", false, "Deprecated")
	withSSH := fs.Bool("ssh", false, "Make available the SSH agent of the host")
	noCache := fs.Bool("no-cache", false, "Always run this specific item, ignoring cache")
	interactive := new(StringSliceFlag)
	fs.Var(interactive, "interactive", "Run this command with an interactive session")
	secrets := new(StringSliceFlag)
	fs.Var(secrets, "secret", "Make available a secret")
	mounts := new(StringSliceFlag)
	fs.Var(mounts, "mount", "Mount a file or directory")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid RUN arguments %v", cmd.Args)
	}
	withShell := !cmd.ExecMode
	if *withDocker {
		*privileged = true
	}
	if !*pushFlag && i.pushOnlyAllowed {
		return Errorf(cmd.SourceLocation, "no non-push commands allowed after a --push")
	}
	// TODO: In the bracket case, should flags be outside of the brackets?

	for index, s := range secrets.Args {
		secrets.Args[index] = i.expandArgs(s, true)
	}
	for index, m := range mounts.Args {
		mounts.Args[index] = i.expandArgs(m, false)
	}
	for index, ia := range interactive.Args {
		interactive.Args[index] = i.expandArgs(ia, false)
	}
	// Note: Not expanding args for the run itself, as that will be take care of by the shell.

	if i.local {
		if len(mounts.Args) > 0 {
			return Errorf(cmd.SourceLocation, "mounts are not supported in combination with the LOCALLY directive")
		}
		if *withSSH {
			return Errorf(cmd.SourceLocation, "the --ssh flag has no effect when used with the  LOCALLY directive")
		}
		if *privileged {
			return Errorf(cmd.SourceLocation, "the --privileged flag has no effect when used with the LOCALLY directive")
		}
		if *noCache {
			return Errorf(cmd.SourceLocation, "the --no-cache flag has no effect when used with the LOCALLY directive")
		}

		// TODO these should be supported, but haven't yet been implemented
		if len(secrets.Args) > 0 {
			return Errorf(cmd.SourceLocation, "secrets need to be implemented for the LOCALLY directive")
		}

		if i.withDocker != nil {
			return Errorf(cmd.SourceLocation, "the WITH DOCKER directive is not (yet) supported with the LOCALLY directive")
		}

		err = i.converter.RunLocal(ctx, fs.Args(), *pushFlag)
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "apply RUN")
		}
		return nil
	}

	if i.withDocker == nil {
		err = i.converter.Run(
			ctx, fs.Args(), mounts.Args, secrets.Args, *privileged, *withEntrypoint, *withDocker,
			withShell, *pushFlag, *withSSH, *noCache, interactive.Args)
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "apply RUN")
		}
		if *pushFlag {
			i.pushOnlyAllowed = true
		}
	} else {
		if *pushFlag {
			return Errorf(cmd.SourceLocation, "RUN --push not allowed in WITH DOCKER")
		}
		if i.withDockerRan {
			return Errorf(cmd.SourceLocation, "only one RUN command allowed in WITH DOCKER")
		}
		i.withDockerRan = true
		i.withDocker.Mounts = mounts.Args
		i.withDocker.Secrets = secrets.Args
		i.withDocker.WithShell = withShell
		i.withDocker.WithEntrypoint = *withEntrypoint
		i.withDocker.NoCache = *noCache
		i.withDocker.Interactive = interactive.Args
		err = i.converter.WithDockerRun(ctx, fs.Args(), *i.withDocker)
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "with docker run")
		}
	}
	return nil
}

func (i *Interpreter) handleFromDockerfile(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	fs := flag.NewFlagSet("FROM DOCKERFILE", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target and also to the Dockerfile build")
	platformStr := fs.String("platform", "", "The platform to use")
	dfTarget := fs.String("target", "", "The Dockerfile target to inherit from")
	dfPath := fs.String("f", "", "Not supported")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid FROM DOCKERFILE arguments %v", cmd.Args)
	}
	if fs.NArg() != 1 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for FROM DOCKERFILE")
	}
	path := i.expandArgs(fs.Arg(0), false)
	_, parseErr := domain.ParseArtifact(path)
	if parseErr != nil {
		// Treat as context path, not artifact path.
		path = i.expandArgs(fs.Arg(0), false)
	}
	for index, ba := range buildArgs.Args {
		buildArgs.Args[index] = i.expandArgs(ba, true)
	}
	*platformStr = i.expandArgs(*platformStr, false)
	platform, err := llbutil.ParsePlatform(*platformStr)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "parse platform %s", *platformStr)
	}
	*dfPath = i.expandArgs(*dfPath, false)
	*dfTarget = i.expandArgs(*dfTarget, false)
	i.local = false
	err = i.converter.FromDockerfile(ctx, path, *dfPath, *dfTarget, platform, buildArgs.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "from dockerfile")
	}
	return nil
}

func (i *Interpreter) handleLocally(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}

	i.local = true
	err := i.converter.Locally(ctx, nil)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "apply LOCALLY")
	}
	return nil
}

func (i *Interpreter) handleCopy(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	fs := flag.NewFlagSet("COPY", flag.ContinueOnError)
	from := fs.String("from", "", "Not supported")
	isDirCopy := fs.Bool("dir", false, "Copy entire directories, not just the contents")
	chown := fs.String("chown", "", "Apply a specific group and/or owner to the copied files and directories")
	keepTs := fs.Bool("keep-ts", false, "Keep created time file timestamps")
	keepOwn := fs.Bool("keep-own", false, "Keep owner info")
	ifExists := fs.Bool("if-exists", false, "Do not fail if the artifact does not exist")
	platformStr := fs.String("platform", "", "The platform to use")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid COPY arguments %v", cmd.Args)
	}
	if fs.NArg() < 2 {
		return Errorf(cmd.SourceLocation, "not enough COPY arguments %v", cmd.Args)
	}
	if *from != "" {
		return Errorf(cmd.SourceLocation, "COPY --from not implemented. Use COPY artifacts form instead")
	}
	srcs := fs.Args()[:fs.NArg()-1]
	dest := i.expandArgs(fs.Arg(fs.NArg()-1), false)
	for index, ba := range buildArgs.Args {
		buildArgs.Args[index] = i.expandArgs(ba, true)
	}
	*chown = i.expandArgs(*chown, false)
	*platformStr = i.expandArgs(*platformStr, false)
	platform, err := llbutil.ParsePlatform(*platformStr)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "parse platform %s", *platformStr)
	}
	allClassical := true
	allArtifacts := true
	for index, src := range srcs {
		// If it parses as an artifact, treat as artifact.
		artifactSrc, parseErr := domain.ParseArtifact(i.expandArgs(src, true))
		if parseErr == nil {
			srcs[index] = artifactSrc.String()
			allClassical = false
		} else {
			srcs[index] = i.expandArgs(src, false)
			allArtifacts = false
		}
	}
	if !allClassical && !allArtifacts {
		return Errorf(cmd.SourceLocation, "combining artifacts and build context arguments in a single COPY command is not allowed: %v", srcs)
	}
	if allArtifacts {
		if dest == "" || dest == "." || len(srcs) > 1 {
			dest += string(filepath.Separator)
		}
		for _, src := range srcs {
			if i.local {
				err = i.converter.CopyArtifactLocal(ctx, src, dest, platform, buildArgs.Args, *isDirCopy)
				if err != nil {
					return WrapError(err, cmd.SourceLocation, "copy artifact locally")
				}
			} else {
				err = i.converter.CopyArtifact(ctx, src, dest, platform, buildArgs.Args, *isDirCopy, *keepTs, *keepOwn, *chown, *ifExists)
				if err != nil {
					return WrapError(err, cmd.SourceLocation, "copy artifact")
				}
			}
		}
	} else {
		if len(buildArgs.Args) != 0 {
			return Errorf(cmd.SourceLocation, "build args not supported for non +artifact arguments case %v", cmd.Args)
		}
		if i.local {
			return Errorf(cmd.SourceLocation, "unhandled locally artifact copy when allArtifacts is false")
		}

		i.converter.CopyClassical(ctx, srcs, dest, *isDirCopy, *keepTs, *keepOwn, *chown)
	}
	return nil
}

func (i *Interpreter) handleSaveArtifact(ctx context.Context, cmd spec.Command) error {
	fs := flag.NewFlagSet("SAVE ARTIFACT", flag.ContinueOnError)
	keepTs := fs.Bool("keep-ts", false, "Keep created time file timestamps")
	keepOwn := fs.Bool("keep-own", false, "Keep owner info")
	ifExists := fs.Bool("if-exists", false, "Do not fail if the artifact does not exist")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid SAVE arguments %v", cmd.Args)
	}

	if fs.NArg() == 0 {
		return Errorf(cmd.SourceLocation, "no arguments provided to the SAVE ARTIFACT command")
	}
	if fs.NArg() > 5 {
		return Errorf(cmd.SourceLocation, "too many arguments provided to the SAVE ARTIFACT command: %v", cmd.Args)
	}
	saveAsLocalTo := ""
	saveTo := "./"
	if fs.NArg() >= 4 {
		if strings.Join(fs.Args()[fs.NArg()-3:fs.NArg()-1], " ") == "AS LOCAL" {
			saveAsLocalTo = fs.Args()[fs.NArg()-1]
			if fs.NArg() == 5 {
				saveTo = fs.Args()[1]
			}
		} else {
			return Errorf(cmd.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Args)
		}
	} else if fs.NArg() == 2 {
		saveTo = fs.Args()[1]
	} else if fs.NArg() == 3 {
		return Errorf(cmd.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Args)
	}

	saveFrom := i.expandArgs(fs.Args()[0], false)
	saveTo = i.expandArgs(saveTo, false)
	saveAsLocalTo = i.expandArgs(saveAsLocalTo, false)

	if i.local {
		if saveAsLocalTo != "" {
			return Errorf(cmd.SourceLocation, "SAVE ARTIFACT AS LOCAL is not implemented under LOCALLY targets")
		}
		err = i.converter.SaveArtifactFromLocal(ctx, saveFrom, saveTo, *keepTs, *ifExists, "")
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
		}
		return nil
	}

	err = i.converter.SaveArtifact(ctx, saveFrom, saveTo, saveAsLocalTo, *keepTs, *keepOwn, *ifExists, i.pushOnlyAllowed)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
	}
	return nil
}

func (i *Interpreter) handleSaveImage(ctx context.Context, cmd spec.Command) error {
	fs := flag.NewFlagSet("SAVE IMAGE", flag.ContinueOnError)
	pushFlag := fs.Bool(
		"push", false,
		"Push the image to the remote registry provided that the build succeeds and also that earthly is invoked in push mode")
	cacheHint := fs.Bool(
		"cache-hint", false,
		"Instruct Earthly that the current target shuold be saved entirely as part of the remote cache")
	insecure := fs.Bool(
		"insecure", false,
		"Use unencrypted connection for the push")
	cacheFrom := new(StringSliceFlag)
	fs.Var(cacheFrom, "cache-from", "Declare additional cache import as a Docker tag")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid SAVE IMAGE arguments %v", cmd.Args)
	}
	for index, cf := range cacheFrom.Args {
		cacheFrom.Args[index] = i.expandArgs(cf, false)
	}
	if *pushFlag && fs.NArg() == 0 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for SAVE IMAGE --push: %v", cmd.Args)
	}

	imageNames := fs.Args()
	for index, img := range imageNames {
		imageNames[index] = i.expandArgs(img, false)
	}
	if len(imageNames) == 0 && !*cacheHint && len(cacheFrom.Args) == 0 {
		fmt.Fprintf(os.Stderr, "Deprecation: using SAVE IMAGE with no arguments is no longer necessary and can be safely removed\n")
		return nil
	}
	err = i.converter.SaveImage(ctx, imageNames, *pushFlag, *insecure, *cacheHint, cacheFrom.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "save image")
	}
	if *pushFlag {
		i.pushOnlyAllowed = true
	}
	return nil
}

func (i *Interpreter) handleBuild(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	fs := flag.NewFlagSet("BUILD", flag.ContinueOnError)
	platformsStr := new(StringSliceFlag)
	fs.Var(platformsStr, "platform", "The platform to build")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid BUILD arguments %v", cmd.Args)
	}
	if fs.NArg() != 1 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for BUILD: %s", cmd.Args)
	}
	fullTargetName := i.expandArgs(fs.Arg(0), true)
	platformsSlice := make([]*specs.Platform, 0, len(platformsStr.Args))
	for index, p := range platformsStr.Args {
		platformsStr.Args[index] = i.expandArgs(p, false)
		platform, err := llbutil.ParsePlatform(p)
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	for index, arg := range buildArgs.Args {
		buildArgs.Args[index] = i.expandArgs(arg, true)
	}
	if len(platformsSlice) == 0 {
		platformsSlice = []*specs.Platform{nil}
	}

	crossProductBuildArgs, err := buildArgMatrix(buildArgs.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "build arg matrix")
	}

	for _, bas := range crossProductBuildArgs {
		for _, platform := range platformsSlice {
			err = i.converter.Build(ctx, fullTargetName, platform, bas)
			if err != nil {
				return WrapError(err, cmd.SourceLocation, "apply BUILD %s", fullTargetName)
			}
		}
	}
	return nil
}

func (i *Interpreter) handleWorkdir(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) != 1 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for WORKDIR: %v", cmd.Args)
	}
	workdirPath := i.expandArgs(cmd.Args[0], false)
	i.converter.Workdir(ctx, workdirPath)
	return nil
}

func (i *Interpreter) handleUser(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) != 1 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for USER: %v", cmd.Args)
	}
	user := i.expandArgs(cmd.Args[0], false)
	i.converter.User(ctx, user)
	return nil
}

func (i *Interpreter) handleCmd(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	withShell := !cmd.ExecMode
	cmdArgs := cmd.Args
	if !withShell {
		for index, arg := range cmdArgs {
			cmdArgs[index] = i.expandArgs(arg, false)
		}
	}
	i.converter.Cmd(ctx, cmdArgs, withShell)
	return nil
}

func (i *Interpreter) handleEntrypoint(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	withShell := !cmd.ExecMode
	entArgs := cmd.Args
	if !withShell {
		for index, arg := range entArgs {
			entArgs[index] = i.expandArgs(arg, false)
		}
	}
	i.converter.Entrypoint(ctx, entArgs, withShell)
	return nil
}

func (i *Interpreter) handleExpose(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) == 0 {
		return Errorf(cmd.SourceLocation, "no arguments provided to the EXPOSE command")
	}
	ports := cmd.Args
	for index, port := range ports {
		ports[index] = i.expandArgs(port, false)
	}
	i.converter.Expose(ctx, ports)
	return nil
}

func (i *Interpreter) handleVolume(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) == 0 {
		return Errorf(cmd.SourceLocation, "no arguments provided to the VOLUME command")
	}
	volumes := cmd.Args
	for index, volume := range volumes {
		volumes[index] = i.expandArgs(volume, false)
	}
	i.converter.Volume(ctx, volumes)
	return nil
}

func (i *Interpreter) handleEnv(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	var key, value string
	switch len(cmd.Args) {
	case 3:
		if cmd.Args[1] != "=" {
			return Errorf(cmd.SourceLocation, "invalid syntax")
		}
		value = i.expandArgs(cmd.Args[2], false)
		fallthrough
	case 1:
		key = cmd.Args[0] // Note: Not expanding args for key.
	default:
		return Errorf(cmd.SourceLocation, "invalid syntax")
	}
	i.converter.Env(ctx, key, value)
	return nil
}

func (i *Interpreter) handleArg(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	var key, value string
	switch len(cmd.Args) {
	case 3:
		if cmd.Args[1] != "=" {
			return Errorf(cmd.SourceLocation, "invalid syntax")
		}
		value = i.expandArgs(cmd.Args[2], true)
		fallthrough
	case 1:
		key = cmd.Args[0] // Note: Not expanding args for key.
	default:
		return Errorf(cmd.SourceLocation, "invalid syntax")
	}
	// Args declared in the base target are global.
	global := i.isBase
	err := i.converter.Arg(ctx, key, value, global)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "apply ARG")
	}
	return nil
}

func (i *Interpreter) handleLabel(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	labels := make(map[string]string)
	var key string
	nextEqual := false
	nextKey := true
	for _, arg := range cmd.Args {
		if nextKey {
			key = i.expandArgs(arg, false)
			nextEqual = true
			nextKey = false
		} else if nextEqual {
			if arg != "=" {
				return Errorf(cmd.SourceLocation, "syntax error")
			}
			nextEqual = false
		} else {
			value := i.expandArgs(arg, false)
			labels[key] = value
			nextKey = true
		}
	}
	if nextKey != true {
		return Errorf(cmd.SourceLocation, "syntax error")
	}
	if len(labels) == 0 {
		return Errorf(cmd.SourceLocation, "no labels provided in LABEL command")
	}
	i.converter.Label(ctx, labels)
	return nil
}

func (i *Interpreter) handleGitClone(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	fs := flag.NewFlagSet("GIT CLONE", flag.ContinueOnError)
	branch := fs.String("branch", "", "The git ref to use when cloning")
	keepTs := fs.Bool("keep-ts", false, "Keep created time file timestamps")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid GIT CLONE arguments %v", cmd.Args)
	}
	if fs.NArg() != 2 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for GIT CLONE: %s", cmd.Args)
	}
	gitURL := i.expandArgs(fs.Arg(0), false)
	gitCloneDest := i.expandArgs(fs.Arg(1), false)
	*branch = i.expandArgs(*branch, false)
	err = i.converter.GitClone(ctx, gitURL, *branch, gitCloneDest, *keepTs)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "git clone")
	}
	return nil
}

func (i *Interpreter) handleHealthcheck(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
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
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid HEALTHCHECK arguments %v", cmd.Args)
	}
	if fs.NArg() == 0 {
		return Errorf(cmd.SourceLocation, "invalid number of arguments for HEALTHCHECK: %s", cmd.Args)
	}
	isNone := false
	var cmdArgs []string
	switch fs.Arg(0) {
	case "NONE":
		if fs.NArg() != 1 {
			return Errorf(cmd.SourceLocation, "invalid arguments for HEALTHCHECK: %s", cmd.Args)
		}
		isNone = true
	case "CMD":
		if fs.NArg() == 1 {
			return Errorf(cmd.SourceLocation, "invalid number of arguments for HEALTHCHECK CMD: %s", cmd.Args)
		}
		cmdArgs = fs.Args()[1:]
	default:
		if strings.HasPrefix(fs.Arg(0), "[") {
			return Errorf(cmd.SourceLocation, "exec form not yet supported for HEALTHCHECK CMD: %s", cmd.Args)
		}
		return Errorf(cmd.SourceLocation, "invalid arguments for HEALTHCHECK: %s", cmd.Args)
	}
	for index, arg := range cmdArgs {
		cmdArgs[index] = i.expandArgs(arg, false)
	}
	i.converter.Healthcheck(ctx, isNone, cmdArgs, *interval, *timeout, *startPeriod, *retries)
	return nil
}

func (i *Interpreter) handleWithDocker(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return pushOnlyErr(cmd.SourceLocation)
	}
	if i.withDocker != nil {
		return Errorf(cmd.SourceLocation, "cannot use WITH DOCKER within WITH DOCKER")
	}

	fs := flag.NewFlagSet("WITH DOCKER", flag.ContinueOnError)
	composeFiles := new(StringSliceFlag)
	fs.Var(composeFiles, "compose", "A compose file used to bring up services from")
	composeServices := new(StringSliceFlag)
	fs.Var(composeServices, "service", "A compose service to bring up")
	loads := new(StringSliceFlag)
	fs.Var(loads, "load", "An image produced by Earthly which is loaded as a Docker image")
	platformStr := fs.String("platform", "", "The platform to use")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "A build arg override passed on to a referenced Earthly target")
	pulls := new(StringSliceFlag)
	fs.Var(pulls, "pull", "An image which is pulled and made available in the docker cache")
	err := fs.Parse(cmd.Args)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "invalid WITH DOCKER arguments %v", cmd.Args)
	}
	if len(fs.Args()) != 0 {
		return Errorf(cmd.SourceLocation, "invalid WITH DOCKER arguments %v", fs.Args())
	}

	*platformStr = i.expandArgs(*platformStr, false)
	platform, err := llbutil.ParsePlatform(*platformStr)
	if err != nil {
		return WrapError(err, cmd.SourceLocation, "parse platform %s", *platformStr)
	}
	for index, cf := range composeFiles.Args {
		composeFiles.Args[index] = i.expandArgs(cf, false)
	}
	for index, cs := range composeServices.Args {
		composeServices.Args[index] = i.expandArgs(cs, false)
	}
	for index, load := range loads.Args {
		loads.Args[index] = i.expandArgs(load, true)
	}
	for index, ba := range buildArgs.Args {
		buildArgs.Args[index] = i.expandArgs(ba, true)
	}
	for index, p := range pulls.Args {
		pulls.Args[index] = i.expandArgs(p, false)
	}

	i.withDocker = &WithDockerOpt{
		ComposeFiles:    composeFiles.Args,
		ComposeServices: composeServices.Args,
	}
	for _, pullStr := range pulls.Args {
		i.withDocker.Pulls = append(i.withDocker.Pulls, DockerPullOpt{
			ImageName: pullStr,
			Platform:  platform,
		})
	}
	for _, loadStr := range loads.Args {
		loadImg, loadTarget, err := parseLoad(loadStr)
		if err != nil {
			return WrapError(err, cmd.SourceLocation, "parse load")
		}
		i.withDocker.Loads = append(i.withDocker.Loads, DockerLoadOpt{
			Target:    loadTarget,
			ImageName: loadImg,
			Platform:  platform,
			BuildArgs: buildArgs.Args,
		})
	}
	return nil
}

func (i *Interpreter) handleAdd(ctx context.Context, cmd spec.Command) error {
	return Errorf(cmd.SourceLocation, "command ADD not yet supported")
}

func (i *Interpreter) handleStopsignal(ctx context.Context, cmd spec.Command) error {
	return Errorf(cmd.SourceLocation, "command STOPSIGNAL not yet supported")
}

func (i *Interpreter) handleOnbuild(ctx context.Context, cmd spec.Command) error {
	return Errorf(cmd.SourceLocation, "command ONBUILD not supported")
}

func (i *Interpreter) handleShell(ctx context.Context, cmd spec.Command) error {
	return Errorf(cmd.SourceLocation, "command SHELL not yet supported")
}

// ----------------------------------------------------------------------------

func (i *Interpreter) expandArgs(word string, keepPlusEscape bool) string {
	ret := i.converter.ExpandArgs(escapeSlashPlus(word))
	if keepPlusEscape {
		return ret
	}
	return unescapeSlashPlus(ret)
}

func escapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "\\\\+")
}

func unescapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "+")
}

func pushOnlyErr(sl *spec.SourceLocation) error {
	return Errorf(sl, "no non-push commands allowed after a --push")
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

type argGroup struct {
	key    string
	values []*string
}

func buildArgMatrix(args []string) ([][]string, error) {
	groupedArgs := make([]argGroup, 0, len(args))
	for _, arg := range args {
		k, v, err := parseKeyValue(arg)
		if err != nil {
			return nil, err
		}

		found := false
		for i, g := range groupedArgs {
			if g.key == k {
				groupedArgs[i].values = append(groupedArgs[i].values, v)
				found = true
				break
			}
		}
		if !found {
			groupedArgs = append(groupedArgs, argGroup{
				key:    k,
				values: []*string{v},
			})
		}
	}
	return crossProduct(groupedArgs, nil), nil
}

func crossProduct(ga []argGroup, prefix []string) [][]string {
	if len(ga) == 0 {
		return [][]string{prefix}
	}
	var ret [][]string
	for _, v := range ga[0].values {
		newPrefix := prefix[:]
		var kv string
		if v == nil {
			kv = ga[0].key
		} else {
			kv = fmt.Sprintf("%s=%s", ga[0].key, *v)
		}
		newPrefix = append(newPrefix, kv)

		cp := crossProduct(ga[1:], newPrefix)
		ret = append(ret, cp...)
	}
	return ret
}

func parseKeyValue(arg string) (string, *string, error) {
	var name string
	splitArg := strings.SplitN(arg, "=", 2)
	if len(splitArg) < 1 {
		return "", nil, fmt.Errorf("invalid build arg %s", splitArg)
	}
	name = splitArg[0]
	var value *string
	if len(splitArg) == 2 {
		value = &splitArg[1]
	}
	return name, value, nil
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
