package earthfile2llb

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/shell"
	"github.com/earthly/earthly/variables"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var errCannotAsync = errors.New("cannot run async operation")

// Interpreter interprets Earthly AST's into calls to the converter.
type Interpreter struct {
	converter *Converter

	target domain.Target

	isBase          bool
	isWith          bool
	pushOnlyAllowed bool
	local           bool
	allowPrivileged bool

	stack string

	withDocker    *WithDockerOpt
	withDockerRan bool

	parallelConversion bool
	console            conslogging.ConsoleLogger
	gitLookup          *buildcontext.GitLookup
}

func newInterpreter(c *Converter, t domain.Target, allowPrivileged, parallelConversion bool, console conslogging.ConsoleLogger, gitLookup *buildcontext.GitLookup) *Interpreter {
	return &Interpreter{
		converter:          c,
		target:             t,
		stack:              c.StackString(),
		allowPrivileged:    allowPrivileged,
		parallelConversion: parallelConversion,
		console:            console,
		gitLookup:          gitLookup,
	}
}

// Run interprets the commands in the given Earthfile AST, for a specific target.
func (i *Interpreter) Run(ctx context.Context, ef spec.Earthfile) (err error) {
	if i.target.Target == "base" {
		i.isBase = true
		err := i.handleBlock(ctx, ef.BaseRecipe)
		i.isBase = false
		return err
	}
	for _, t := range ef.Targets {
		if t.Name == i.target.Target {
			return i.handleTarget(ctx, t)
		}
	}
	return i.errorf(ef.SourceLocation, "target %s not found", i.target.Target)
}

func (i *Interpreter) handleTarget(ctx context.Context, t spec.Target) error {
	// Apply implicit FROM +base
	err := i.converter.From(ctx, "+base", platutil.DefaultPlatform, i.allowPrivileged, nil)
	if err != nil {
		return i.wrapError(err, t.SourceLocation, "apply FROM")
	}
	return i.handleBlock(ctx, t.Recipe)
}

func (i *Interpreter) handleBlock(ctx context.Context, b spec.Block) error {
	prevWasArg := true // not exactly true, but makes the logic easier
	for index, stmt := range b {
		if i.parallelConversion && prevWasArg {
			err := i.handleBlockParallel(ctx, b, index)
			if err != nil {
				return err
			}
		}
		err := i.handleStatement(ctx, stmt)
		if err != nil {
			return err
		}
		prevWasArg = (stmt.Command != nil && stmt.Command.Name == "ARG")
	}
	return nil
}

func (i *Interpreter) handleBlockParallel(ctx context.Context, b spec.Block, startIndex int) error {
	if i.local {
		// Don't do any preemptive execution for LOCALLY targets.
		return nil
	}
	// Look ahead of the execution and fire off asynchronous builds for mentioned targets,
	// as long as they don't have variable args $(...).
	for index := startIndex; index < len(b); index++ {
		stmt := b[index]
		if stmt.Command != nil {
			switch stmt.Command.Name {
			case "ARG", "LOCALLY", "FROM", "FROM DOCKERFILE":
				// Cannot do any further parallel builds - these commands need to be
				// executed to ensure that they don't impact the outcome. As such,
				// commands following these cannot be executed preemptively.
				return nil
			case "BUILD":
				err := i.handleBuild(ctx, *stmt.Command, true)
				if err != nil {
					if errors.Is(err, errCannotAsync) {
						continue // no biggie
					}
					return err
				}
			case "COPY":
				// TODO
			}
		} else if stmt.With != nil {
			switch stmt.With.Command.Name {
			case "DOCKER":
				// TODO
			}
		} else if stmt.If != nil || stmt.For != nil || stmt.Wait != nil {
			// Cannot do any further parallel builds - these commands need to be
			// executed to ensure that they don't impact the outcome. As such,
			// commands following these cannot be executed preemptively.
			return nil
		} else {
			return i.errorf(stmt.SourceLocation, "unexpected statement type")
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
	} else if stmt.For != nil {
		return i.handleFor(ctx, *stmt.For)
	} else if stmt.Wait != nil {
		return i.handleWait(ctx, *stmt.Wait)
	} else {
		return i.errorf(stmt.SourceLocation, "unexpected statement type")
	}
}

func (i *Interpreter) handleCommand(ctx context.Context, cmd spec.Command) (err error) {
	// The AST should not be modified by any operation. This is a consistency check.
	argsCopy := getArgsCopy(cmd)
	defer func() {
		if err != nil {
			return
		}
		if len(argsCopy) != len(cmd.Args) {
			err = i.errorf(cmd.SourceLocation, "internal error: args were modified in command handling")
			return
		}
		for index, arg := range cmd.Args {
			if arg != argsCopy[index] {
				err = i.errorf(cmd.SourceLocation, "internal error: args were modified in command handling")
				return
			}
		}
	}()

	analytics.Count("cmd", cmd.Name)

	if i.isWith {
		switch cmd.Name {
		case "DOCKER":
			return i.handleWithDocker(ctx, cmd)
		default:
			return i.errorf(cmd.SourceLocation, "unexpected WITH command %s", cmd.Name)
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
		return i.handleBuild(ctx, cmd, false)
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
	case "COMMAND":
		return i.handleUserCommand(ctx, cmd)
	case "DO":
		return i.handleDo(ctx, cmd)
	case "IMPORT":
		return i.handleImport(ctx, cmd)
	case "CACHE":
		return i.handleCache(ctx, cmd)
	case "HOST":
		return i.handleHost(ctx, cmd)
	case "PROJECT":
		return i.handleProject(ctx, cmd)
	default:
		return i.errorf(cmd.SourceLocation, "unexpected command %s", cmd.Name)
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
		return i.errorf(with.SourceLocation, "only one command is allowed in WITH DOCKER and it has to be RUN")
	}
	err = i.handleBlock(ctx, with.Body)
	if err != nil {
		return err
	}
	i.withDocker = nil
	if !i.withDockerRan {
		return i.errorf(with.SourceLocation, "no RUN command found in WITH DOCKER")
	}
	i.withDockerRan = false
	return nil
}

func (i *Interpreter) handleIf(ctx context.Context, ifStmt spec.IfStatement) error {
	if i.pushOnlyAllowed {
		return i.errorf(ifStmt.SourceLocation, "no non-push commands allowed after a --push")
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
		return false, i.errorf(sl, "not enough arguments for IF")
	}
	opts := ifOpts{}
	args, err := parseArgs("IF", &opts, expression)
	if err != nil {
		return false, i.wrapError(err, sl, "invalid IF arguments %v", expression)
	}
	withShell := !execMode

	for index, s := range opts.Secrets {
		expanded, err := i.expandArgs(ctx, s, true, false)
		if err != nil {
			return false, i.wrapError(err, sl, "failed to expand IF secret %v", s)
		}
		opts.Secrets[index] = expanded
	}
	for index, m := range opts.Mounts {
		expanded, err := i.expandArgs(ctx, m, false, false)
		if err != nil {
			return false, i.wrapError(err, sl, "failed to expand IF mount %v", m)
		}
		opts.Mounts[index] = expanded
	}
	// Note: Not expanding args for the expression itself, as that will be take care of by the shell.

	var exitCode int
	runOpts := ConvertRunOpts{
		CommandName: "IF",
		Args:        args,
		Locally:     i.local,
		Mounts:      opts.Mounts,
		Secrets:     opts.Secrets,
		WithShell:   withShell,
		Privileged:  opts.Privileged,
		WithSSH:     opts.WithSSH,
		NoCache:     opts.NoCache,
		Transient:   !i.local,
	}
	exitCode, err = i.converter.RunExitCode(ctx, runOpts)
	if err != nil {
		return false, i.wrapError(err, sl, "apply IF")
	}
	return (exitCode == 0), nil
}

func (i *Interpreter) handleFor(ctx context.Context, forStmt spec.ForStatement) error {
	if !i.converter.ftrs.ForIn {
		return i.errorf(forStmt.SourceLocation, "the FOR command is not supported in this version")
	}
	variable, instances, err := i.handleForArgs(ctx, forStmt.Args, forStmt.SourceLocation)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		err = i.converter.SetArg(ctx, variable, instance)
		if err != nil {
			return i.wrapError(err, forStmt.SourceLocation, "set %s=%s", variable, instance)
		}
		err = i.handleBlock(ctx, forStmt.Body)
		if err != nil {
			return err
		}
		err = i.converter.UnsetArg(ctx, variable)
		if err != nil {
			return i.wrapError(err, forStmt.SourceLocation, "unset %s", variable)
		}
	}
	return nil
}

func (i *Interpreter) handleForArgs(ctx context.Context, forArgs []string, sl *spec.SourceLocation) (string, []string, error) {
	opts := forOpts{
		Separators: "\n\t ",
	}
	args, err := parseArgs("FOR", &opts, forArgs)
	if err != nil {
		return "", nil, i.wrapError(err, sl, "invalid FOR arguments %v", forArgs)
	}
	if len(args) < 3 {
		return "", nil, i.errorf(sl, "not enough arguments for FOR")
	}
	if args[1] != "IN" {
		return "", nil, i.errorf(sl, "expected IN, got %s", args[1])
	}
	variable := args[0]
	runOpts := ConvertRunOpts{
		CommandName: "FOR",
		Args:        args[2:],
		Locally:     i.local,
		Mounts:      opts.Mounts,
		Secrets:     opts.Secrets,
		WithShell:   true,
		Privileged:  opts.Privileged,
		WithSSH:     opts.WithSSH,
		NoCache:     opts.NoCache,
		Transient:   !i.local,
	}
	output, err := i.converter.RunExpression(ctx, variable, runOpts)
	if err != nil {
		return "", nil, i.wrapError(err, sl, "apply FOR ... IN")
	}
	instances := strings.FieldsFunc(output, func(r rune) bool {
		return strings.ContainsRune(opts.Separators, r)
	})
	return variable, instances, nil
}

func (i *Interpreter) handleWait(ctx context.Context, waitStmt spec.WaitStatement) error {
	if !i.converter.ftrs.WaitBlock {
		return i.errorf(waitStmt.SourceLocation, "the WAIT command is not supported in this version")
	}

	if len(waitStmt.Args) != 0 {
		return i.errorf(waitStmt.SourceLocation, "WAIT does not accept any options")
	}

	err := i.converter.PushWaitBlock(ctx)
	if err != nil {
		return err
	}

	err = i.handleBlock(ctx, waitStmt.Body)
	if err != nil {
		return err
	}

	return i.converter.PopWaitBlock(ctx)
}

// Commands -------------------------------------------------------------------

func (i *Interpreter) handleFrom(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := fromOpts{}
	args, err := parseArgs("FROM", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid FROM arguments %v", cmd.Args)
	}
	if len(args) != 1 {
		if len(args) == 3 && args[1] == "AS" {
			return i.errorf(cmd.SourceLocation, "AS not supported, use earthly targets instead")
		}
		if len(args) < 1 {
			return i.errorf(cmd.SourceLocation, "invalid number of arguments for FROM: %s", cmd.Args)
		}
	}
	imageName, err := i.expandArgs(ctx, args[0], true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand image name for FROM: %s", args[0])
	}
	expandedPlatform, err := i.expandArgs(ctx, opts.Platform, false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand platform for FROM: %s", opts.Platform)
	}
	platform, err := i.converter.platr.Parse(expandedPlatform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", expandedPlatform)
	}
	expandedBuildArgs, err := i.expandArgsSlice(ctx, opts.BuildArgs, true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand build args for FROM: %v", opts.BuildArgs)
	}
	expandedFlagArgs, err := i.expandArgsSlice(ctx, args[1:], true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand flag args for FROM: %v", args[1:])
	}
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}
	expandedBuildArgs = append(parsedFlagArgs, expandedBuildArgs...)

	allowPrivileged, err := i.getAllowPrivilegedTarget(imageName, opts.AllowPrivileged)
	if err != nil {
		return err
	}

	i.local = false
	err = i.converter.From(ctx, imageName, platform, allowPrivileged, expandedBuildArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply FROM %s", imageName)
	}
	return nil
}

func (i *Interpreter) getAllowPrivilegedTarget(targetName string, allowPrivileged bool) (bool, error) {
	if !strings.Contains(targetName, "+") {
		return false, nil
	}
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return false, errors.Wrapf(err, "parse target name %s", targetName)
	}

	return i.getAllowPrivileged(depTarget, allowPrivileged)
}

func (i *Interpreter) getAllowPrivileged(depTarget domain.Target, allowPrivileged bool) (bool, error) {
	if depTarget.IsRemote() {
		return i.allowPrivileged && allowPrivileged, nil
	}
	if allowPrivileged {
		i.console.Printf("the --allow-privileged flag has no effect when referencing a local target\n")
	}
	return i.allowPrivileged, nil
}

func (i *Interpreter) getAllowPrivilegedArtifact(artifactName string, allowPrivileged bool) (bool, error) {
	artifact, err := domain.ParseArtifact(artifactName)
	if err != nil {
		return false, errors.Wrapf(err, "parse artifact name %s", artifactName)
	}

	return i.getAllowPrivileged(artifact.Target, allowPrivileged)
}

func (i *Interpreter) flagValModifierFuncWithContext(ctx context.Context) func(string, *flags.Option, *string) (*string, error) {
	return func(flagName string, flagOpt *flags.Option, flagVal *string) (*string, error) {
		if flagOpt.IsBool() && flagVal != nil {
			newFlag, err := i.expandArgs(ctx, *flagVal, false, false)
			if err != nil {
				return nil, err
			}
			return &newFlag, nil
		}
		return flagVal, nil
	}
}

func (i *Interpreter) handleRun(ctx context.Context, cmd spec.Command) error {
	if len(cmd.Args) < 1 {
		return i.errorf(cmd.SourceLocation, "not enough arguments for RUN")
	}
	opts := runOpts{}
	args, err := parseArgsWithValueModifier("RUN", &opts, getArgsCopy(cmd), i.flagValModifierFuncWithContext(ctx))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid RUN arguments %v", cmd.Args)
	}
	withShell := !cmd.ExecMode
	if opts.WithDocker {
		opts.Privileged = true
	}
	if !opts.Push && i.pushOnlyAllowed {
		return i.errorf(cmd.SourceLocation, "no non-push commands allowed after a --push")
	}
	// TODO: In the bracket case, should flags be outside of the brackets?

	for index, s := range opts.Secrets {
		expanded, err := i.expandArgs(ctx, s, true, false)
		if err != nil {
			return i.errorf(cmd.SourceLocation, "failed to expand secrets arg in RUN: %s", s)
		}
		opts.Secrets[index] = expanded
	}
	for index, m := range opts.Mounts {
		expanded, err := i.expandArgs(ctx, m, false, false)
		if err != nil {
			return i.errorf(cmd.SourceLocation, "failed to expand mount arg in RUN: %s", m)
		}
		opts.Mounts[index] = expanded
	}
	// Note: Not expanding args for the run itself, as that will be take care of by the shell.

	if opts.Privileged && !i.allowPrivileged {
		return i.errorf(cmd.SourceLocation, "Permission denied: unwilling to run privileged command; did you reference a remote Earthfile without the --allow-privileged flag?")
	}

	if i.withDocker == nil {
		if opts.WithDocker {
			return i.errorf(cmd.SourceLocation, "--with-docker is obsolete. Please use WITH DOCKER ... RUN ... END instead")
		}
		opts := ConvertRunOpts{
			CommandName:     cmd.Name,
			Args:            args,
			Locally:         i.local,
			Mounts:          opts.Mounts,
			Secrets:         opts.Secrets,
			WithShell:       withShell,
			WithEntrypoint:  opts.WithEntrypoint,
			Privileged:      opts.Privileged,
			Push:            opts.Push,
			WithSSH:         opts.WithSSH,
			NoCache:         opts.NoCache,
			Interactive:     opts.Interactive,
			InteractiveKeep: opts.InteractiveKeep,
		}
		err = i.converter.Run(ctx, opts)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply RUN")
		}
		if opts.Push {
			i.pushOnlyAllowed = true
		}
	} else {
		if i.withDockerRan {
			return i.errorf(cmd.SourceLocation, "only one RUN command allowed in WITH DOCKER")
		}
		if opts.Push {
			return i.errorf(cmd.SourceLocation, "RUN --push not allowed in WITH DOCKER")
		}
		i.withDocker.Mounts = opts.Mounts
		i.withDocker.Secrets = opts.Secrets
		i.withDocker.WithShell = withShell
		i.withDocker.WithEntrypoint = opts.WithEntrypoint
		i.withDocker.WithSSH = opts.WithSSH
		i.withDocker.NoCache = opts.NoCache
		i.withDocker.Interactive = opts.Interactive
		i.withDocker.interactiveKeep = opts.InteractiveKeep
		// TODO: Could this be allowed in the future, if dynamic build args
		//       are expanded ahead of time?
		allowParallel := true
		for _, l := range i.withDocker.Loads {
			if !isSafeAsyncBuildArgsKVStyle(l.BuildArgs) {
				allowParallel = false
				break
			}
		}

		if i.local {
			err = i.converter.WithDockerRunLocal(ctx, args, *i.withDocker, allowParallel)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "with docker run")
			}
		} else {
			err = i.converter.WithDockerRun(ctx, args, *i.withDocker, allowParallel)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "with docker run")
			}
		}
		i.withDockerRan = true
	}
	return nil
}

func (i *Interpreter) handleFromDockerfile(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := fromDockerfileOpts{}
	args, err := parseArgs("FROM DOCKERFILE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid FROM DOCKERFILE arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for FROM DOCKERFILE")
	}
	path, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand FROM DOCKERFILE path arg %s", args[0])
	}
	_, parseErr := domain.ParseArtifact(path)
	if parseErr != nil {
		// Treat as context path, not artifact path.
		path, err = i.expandArgs(ctx, args[0], false, false)
		if err != nil {
			return i.errorf(cmd.SourceLocation, "failed to expand FROM DOCKERFILE path arg %s", args[0])
		}
	}
	expandedBuildArgs, err := i.expandArgsSlice(ctx, opts.BuildArgs, true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand FROM DOCKERFILE build args %s", opts.BuildArgs)
	}
	expandedFlagArgs, err := i.expandArgsSlice(ctx, args[1:], true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand FROM DOCKERFILE flag args %s", args[1:])
	}
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}
	expandedBuildArgs = append(parsedFlagArgs, expandedBuildArgs...)
	expandedPlatform, err := i.expandArgs(ctx, opts.Platform, false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand FROM DOCKERFILE platform %s", opts.Platform)
	}
	platform, err := i.converter.platr.Parse(expandedPlatform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", expandedPlatform)
	}
	expandedPath, err := i.expandArgs(ctx, opts.Path, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand path %s", opts.Path)
	}
	expandedTarget, err := i.expandArgs(ctx, opts.Target, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand target %s", opts.Target)
	}
	i.local = false
	err = i.converter.FromDockerfile(ctx, path, expandedPath, expandedTarget, platform, expandedBuildArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "from dockerfile")
	}
	return nil
}

func (i *Interpreter) handleLocally(ctx context.Context, cmd spec.Command) error {
	if !i.allowPrivileged {
		return i.errorf(cmd.SourceLocation, "Permission denied: unwilling to allow locally directive from remote Earthfile; did you reference a remote Earthfile without the --allow-privileged flag?")
	}

	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}

	i.local = true
	err := i.converter.Locally(ctx)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply LOCALLY")
	}
	return nil
}

func (i *Interpreter) handleCopy(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := copyOpts{}
	args, err := parseArgs("COPY", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid COPY arguments %v", cmd.Args)
	}
	if len(args) < 2 {
		return i.errorf(cmd.SourceLocation, "not enough COPY arguments %v", cmd.Args)
	}
	if opts.From != "" {
		return i.errorf(cmd.SourceLocation, "COPY --from not implemented. Use COPY artifacts form instead")
	}
	srcs := args[:len(args)-1]
	srcFlagArgs := make([][]string, len(srcs))
	dest, err := i.expandArgs(ctx, args[len(args)-1], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY args %v", args[len(args)-1])
	}
	expandedBuildArgs, err := i.expandArgsSlice(ctx, opts.BuildArgs, true, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY buildargs %v", opts.BuildArgs)
	}
	expandedChown, err := i.expandArgs(ctx, opts.Chown, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY chown: %v", opts.Chown)
	}
	var fileModeParsed *os.FileMode
	if opts.Chmod != "" {
		expandedMode, err := i.expandArgs(ctx, opts.Chmod, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY chmod: %v", opts.Platform)
		}
		mask, err := strconv.ParseUint(expandedMode, 8, 32)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to parse COPY chmod: %v", opts.Platform)
		}
		mode := os.FileMode(uint32(mask))
		fileModeParsed = &mode
	}
	expandedPlatform, err := i.expandArgs(ctx, opts.Platform, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY platform: %v", opts.Platform)
	}
	platform, err := i.converter.platr.Parse(expandedPlatform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", expandedPlatform)
	}

	allClassical := true
	allArtifacts := true
	for index, src := range srcs {
		var artifactSrc domain.Artifact
		var parseErr error
		if isInParansForm(src) {
			// COPY (<src> <flag-args>) ...
			artifactStr, extraArgs, err := parseParans(src)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "parse parans %s", src)
			}
			expandedArtifact, err := i.expandArgs(ctx, artifactStr, true, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY artifact %s", artifactStr)
			}
			artifactSrc, parseErr = domain.ParseArtifact(expandedArtifact)
			if parseErr != nil {
				// Must parse in the parans case.
				return i.wrapError(err, cmd.SourceLocation, "parse artifact")
			}
			srcFlagArgs[index] = extraArgs
		} else {
			expandedSrc, err := i.expandArgs(ctx, src, true, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY src %s", src)
			}
			artifactSrc, parseErr = domain.ParseArtifact(expandedSrc)
		}
		// If it parses as an artifact, treat as artifact.
		if parseErr == nil {
			srcs[index] = artifactSrc.String()
			allClassical = false
		} else {
			expandedSrc, err := i.expandArgs(ctx, src, false, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY src %s", src)
			}
			srcs[index] = expandedSrc
			allArtifacts = false
		}
	}
	if !allClassical && !allArtifacts {
		return i.errorf(cmd.SourceLocation, "combining artifacts and build context arguments in a single COPY command is not allowed: %v", srcs)
	}
	if allArtifacts {
		if dest == "" || dest == "." || len(srcs) > 1 {
			dest += string("/") // TODO needs to be the containers platform, not the earthly hosts platform. For now, this is always Linux.
		}
		for index, src := range srcs {
			allowPrivileged, err := i.getAllowPrivilegedArtifact(src, opts.AllowPrivileged)
			if err != nil {
				return err
			}

			expandedFlagArgs, err := i.expandArgsSlice(ctx, srcFlagArgs[index], true, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY flag %s", srcFlagArgs[index])
			}
			parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "parse flag args")
			}
			srcBuildArgs := append(parsedFlagArgs, expandedBuildArgs...)

			if i.local {
				err = i.converter.CopyArtifactLocal(ctx, src, dest, platform, allowPrivileged, srcBuildArgs, opts.IsDirCopy)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "copy artifact locally")
				}
			} else {
				err = i.converter.CopyArtifact(ctx, src, dest, platform, allowPrivileged, srcBuildArgs, opts.IsDirCopy, opts.KeepTs, opts.KeepOwn, expandedChown, fileModeParsed, opts.IfExists, opts.SymlinkNoFollow)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "copy artifact")
				}
			}
		}
	} else {
		if len(expandedBuildArgs) != 0 {
			return i.errorf(cmd.SourceLocation, "build args not supported for non +artifact arguments case %v", cmd.Args)
		}
		if i.local {
			return i.errorf(cmd.SourceLocation, "unhandled locally artifact copy when allArtifacts is false")
		}

		err = i.converter.CopyClassical(ctx, srcs, dest, opts.IsDirCopy, opts.KeepTs, opts.KeepOwn, expandedChown, fileModeParsed, opts.IfExists)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "copy classical")
		}
	}
	return nil
}

func (i *Interpreter) handleSaveArtifact(ctx context.Context, cmd spec.Command) error {
	opts := saveArtifactOpts{}
	args, err := parseArgs("SAVE ARTIFACT", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid SAVE ARTIFACT arguments %v", cmd.Args)
	}

	if len(args) == 0 {
		return i.errorf(cmd.SourceLocation, "no arguments provided to the SAVE ARTIFACT command")
	}
	if len(args) > 5 {
		return i.errorf(cmd.SourceLocation, "too many arguments provided to the SAVE ARTIFACT command: %v", cmd.Args)
	}
	saveAsLocalTo := ""
	saveTo := "./"
	if len(args) >= 4 {
		if strings.Join(args[len(args)-3:len(args)-1], " ") == "AS LOCAL" {
			saveAsLocalTo = args[len(args)-1]
			if len(args) == 5 {
				saveTo = args[1]
			}
		} else {
			return i.errorf(cmd.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Args)
		}
	} else if len(args) == 2 {
		saveTo = args[1]
	} else if len(args) == 3 {
		return i.errorf(cmd.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Args)
	}

	saveFrom, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand SAVE ARTIFACT src: %s", args[0])
	}
	expandedSaveTo, err := i.expandArgs(ctx, saveTo, false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand SAVE ARTIFACT dst: %s", saveTo)
	}
	expandedSaveAsLocalTo, err := i.expandArgs(ctx, saveAsLocalTo, false, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "failed to expand SAVE ARTIFACT local dst: %s", saveAsLocalTo)
	}

	if i.local {
		if expandedSaveAsLocalTo != "" {
			return i.errorf(cmd.SourceLocation, "SAVE ARTIFACT AS LOCAL is not implemented under LOCALLY targets")
		}
		err = i.converter.SaveArtifactFromLocal(ctx, saveFrom, expandedSaveTo, opts.KeepTs, opts.IfExists, "")
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
		}
		return nil
	}

	err = i.converter.SaveArtifact(ctx, saveFrom, expandedSaveTo, expandedSaveAsLocalTo, opts.KeepTs, opts.KeepOwn, opts.IfExists, opts.SymlinkNoFollow, opts.Force, i.pushOnlyAllowed)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
	}
	return nil
}

func (i *Interpreter) handleSaveImage(ctx context.Context, cmd spec.Command) error {
	opts := saveImageOpts{}
	args, err := parseArgs("SAVE IMAGE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid SAVE IMAGE arguments %v", cmd.Args)
	}
	for index, cf := range opts.CacheFrom {
		expandedCacheFrom, err := i.expandArgs(ctx, cf, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand SAVE IMAGE cache-from: %s", cf)
		}
		opts.CacheFrom[index] = expandedCacheFrom
	}
	if opts.Push && len(args) == 0 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for SAVE IMAGE --push: %v", cmd.Args)
	}

	imageNames := args
	for index, img := range imageNames {
		expandedImageName, err := i.expandArgs(ctx, img, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand SAVE IMAGE img: %s", img)
		}
		imageNames[index] = expandedImageName
	}
	if len(imageNames) == 0 && !opts.CacheHint && len(opts.CacheFrom) == 0 {
		fmt.Fprintf(os.Stderr, "Deprecation: using SAVE IMAGE with no arguments is no longer necessary and can be safely removed\n")
		return nil
	}
	err = i.converter.SaveImage(ctx, imageNames, opts.Push, opts.Insecure, opts.CacheHint, opts.CacheFrom, opts.NoManifestList)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "save image")
	}
	if opts.Push {
		i.pushOnlyAllowed = true
	}
	return nil
}

func (i *Interpreter) handleBuild(ctx context.Context, cmd spec.Command, async bool) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := buildOpts{}
	args, err := parseArgs("BUILD", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid BUILD arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for BUILD: %s", cmd.Args)
	}
	fullTargetName, err := i.expandArgs(ctx, args[0], true, async)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand BUILD target %s", args[0])
	}
	platformsSlice := make([]platutil.Platform, 0, len(opts.Platforms))
	for index, p := range opts.Platforms {
		expandedPlatform, err := i.expandArgs(ctx, p, false, async)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand BUILD platform %s", p)
		}
		opts.Platforms[index] = expandedPlatform
		platform, err := i.converter.platr.Parse(p)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	if async && !(isSafeAsyncBuildArgsKVStyle(opts.BuildArgs) && isSafeAsyncBuildArgs(args[1:])) {
		return errCannotAsync
	}
	if i.local && !(isSafeAsyncBuildArgsKVStyle(opts.BuildArgs) && isSafeAsyncBuildArgs(args[1:])) {
		return i.errorf(cmd.SourceLocation, "BUILD args do not currently support shelling-out in combination with LOCALLY")
	}
	expandedBuildArgs, err := i.expandArgsSlice(ctx, opts.BuildArgs, true, async)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand BUILD args %v", opts.BuildArgs)
	}
	expandedFlagArgs, err := i.expandArgsSlice(ctx, args[1:], true, async)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand BUILD flags %v", args[1:])
	}
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}
	expandedBuildArgs = append(parsedFlagArgs, expandedBuildArgs...)
	if len(platformsSlice) == 0 {
		platformsSlice = []platutil.Platform{platutil.DefaultPlatform}
	}

	crossProductBuildArgs, err := buildArgMatrix(expandedBuildArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "build arg matrix")
	}

	allowPrivileged, err := i.getAllowPrivilegedTarget(fullTargetName, opts.AllowPrivileged)
	if err != nil {
		return err
	}

	for _, bas := range crossProductBuildArgs {
		for _, platform := range platformsSlice {
			if async {
				err = i.converter.BuildAsync(ctx, fullTargetName, platform, allowPrivileged, bas, buildCmd, nil, nil)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "apply BUILD %s", fullTargetName)
				}
			} else {
				err = i.converter.Build(ctx, fullTargetName, platform, allowPrivileged, bas)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "apply BUILD %s", fullTargetName)
				}
			}
		}
	}
	return nil
}

func (i *Interpreter) handleWorkdir(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) != 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for WORKDIR: %v", cmd.Args)
	}
	workdirPath, err := i.expandArgs(ctx, cmd.Args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand WORKDIR path %s", cmd.Args[0])
	}
	err = i.converter.Workdir(ctx, workdirPath)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply WORKDIR")
	}
	return nil
}

func (i *Interpreter) handleUser(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) != 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for USER: %v", cmd.Args)
	}
	user, err := i.expandArgs(ctx, cmd.Args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand USER %s", cmd.Args[0])
	}
	err = i.converter.User(ctx, user)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply USER")
	}
	return nil
}

func (i *Interpreter) handleCmd(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	withShell := !cmd.ExecMode
	cmdArgs := getArgsCopy(cmd)
	if withShell {
		for index, arg := range cmdArgs {
			expandedCmd, err := i.expandArgs(ctx, arg, false, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand CMD %s", arg)
			}
			cmdArgs[index] = expandedCmd
		}
	}
	err := i.converter.Cmd(ctx, cmdArgs, withShell)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply CMD")
	}
	return nil
}

func (i *Interpreter) handleEntrypoint(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	withShell := !cmd.ExecMode
	entArgs := getArgsCopy(cmd)
	if withShell {
		for index, arg := range entArgs {
			expandedEntrypoint, err := i.expandArgs(ctx, arg, false, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand ENTRYPOINT %s", arg)
			}
			entArgs[index] = expandedEntrypoint
		}
	}
	err := i.converter.Entrypoint(ctx, entArgs, withShell)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply ENTRYPOINT")
	}
	return nil
}

func (i *Interpreter) handleExpose(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) == 0 {
		return i.errorf(cmd.SourceLocation, "no arguments provided to the EXPOSE command")
	}
	ports := getArgsCopy(cmd)
	for index, port := range ports {
		expandedPort, err := i.expandArgs(ctx, port, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand EXPOSE %s", port)
		}
		ports[index] = expandedPort
	}
	err := i.converter.Expose(ctx, ports)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply EXPOSE")
	}
	return nil
}

func (i *Interpreter) handleVolume(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	if len(cmd.Args) == 0 {
		return i.errorf(cmd.SourceLocation, "no arguments provided to the VOLUME command")
	}
	volumes := getArgsCopy(cmd)
	for index, volume := range volumes {
		expandedVolume, err := i.expandArgs(ctx, volume, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand VOLUME %s", volume)
		}
		volumes[index] = expandedVolume
	}
	err := i.converter.Volume(ctx, volumes)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply VOLUME")
	}
	return nil
}

func (i *Interpreter) handleEnv(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	var err error
	var key, value string
	switch len(cmd.Args) {
	case 3:
		if cmd.Args[1] != "=" {
			return i.errorf(cmd.SourceLocation, "invalid syntax")
		}
		value, err = i.expandArgs(ctx, cmd.Args[2], false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand ENV %s", cmd.Args[2])
		}
		fallthrough
	case 1:
		key = cmd.Args[0] // Note: Not expanding args for key.
	default:
		return i.errorf(cmd.SourceLocation, "invalid syntax")
	}
	err = i.converter.Env(ctx, key, value)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply ENV")
	}
	return nil
}

var errInvalidSyntax = errors.New("invalid syntax")
var errRequiredArgHasDefault = errors.New("required ARG cannot have a default value")
var errGlobalArgNotInBase = errors.New("global ARG can only be set in the base target")

// parseArgArgs parses the ARG command's arguments
// and returns the argOpts, key, value (or nil if missing), or error
func parseArgArgs(ctx context.Context, cmd spec.Command, isBaseTarget bool, explicitGlobalFeature bool) (argOpts, string, *string, error) {
	opts := argOpts{}
	args, err := parseArgs("ARG", &opts, getArgsCopy(cmd))
	if err != nil {
		return argOpts{}, "", nil, err
	}
	if opts.Global {
		// since the global flag is part of the struct, we need to manually return parsing error if it's used while the feature flag is off
		if !explicitGlobalFeature {
			return argOpts{}, "", nil, errors.New("unknown flag --global")
		}
		// global flag can only bet set on base targets
		if !isBaseTarget {
			return argOpts{}, "", nil, errGlobalArgNotInBase
		}
	} else if !explicitGlobalFeature {
		// if the feature flag is off, all base target args are considered global
		opts.Global = isBaseTarget
	}
	switch len(args) {
	case 3:
		if args[1] != "=" {
			return argOpts{}, "", nil, errInvalidSyntax
		}
		if opts.Required {
			return argOpts{}, "", nil, errRequiredArgHasDefault
		}
		return opts, args[0], &args[2], nil
	case 1:
		return opts, args[0], nil, nil
	default:
		return argOpts{}, "", nil, errInvalidSyntax
	}
}

func (i *Interpreter) handleArg(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts, key, valueOrNil, err := parseArgArgs(ctx, cmd, i.isBase, i.converter.ftrs.ExplicitGlobal)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid ARG arguments %v", cmd.Args)
	}

	var value string
	if valueOrNil != nil {
		value, err = i.expandArgs(ctx, *valueOrNil, true, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand ARG %s", *valueOrNil)
		}
	}

	err = i.converter.Arg(ctx, key, value, opts)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply ARG")
	}
	return nil
}

func (i *Interpreter) handleLabel(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	labels := make(map[string]string)
	var err error
	var key string
	nextEqual := false
	nextKey := true
	for _, arg := range cmd.Args {
		if nextKey {
			key, err = i.expandArgs(ctx, arg, false, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand LABEL key %s", arg)
			}
			nextEqual = true
			nextKey = false
		} else if nextEqual {
			if arg != "=" {
				return i.errorf(cmd.SourceLocation, "syntax error")
			}
			nextEqual = false
		} else {
			value, err := i.expandArgs(ctx, arg, false, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand LABEL value %s", arg)
			}
			labels[key] = value
			nextKey = true
		}
	}
	if !nextKey {
		return i.errorf(cmd.SourceLocation, "syntax error")
	}
	if len(labels) == 0 {
		return i.errorf(cmd.SourceLocation, "no labels provided in LABEL command")
	}
	err = i.converter.Label(ctx, labels)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply LABEL")
	}
	return nil
}

func (i *Interpreter) handleGitClone(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := gitCloneOpts{}
	args, err := parseArgs("GIT CLONE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid GIT CLONE arguments %v", cmd.Args)
	}
	if len(args) != 2 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for GIT CLONE: %s", cmd.Args)
	}
	gitURL, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand GIT CLONE url: %s", args[0])
	}

	gitCloneDest, err := i.expandArgs(ctx, args[1], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand GIT CLONE dest: %s", args[1])
	}
	gitBranch, err := i.expandArgs(ctx, opts.Branch, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand GIT CLONE dest: %s", opts.Branch)
	}

	convertedGitURL, _, err := i.gitLookup.ConvertCloneURL(gitURL)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to use %v with configured earthly credentials from ~/.earthly/config.yml", cmd.Args)
	}

	err = i.converter.GitClone(ctx, convertedGitURL, gitBranch, gitCloneDest, opts.KeepTs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "git clone")
	}
	return nil
}

func (i *Interpreter) handleHealthcheck(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := healthCheckOpts{}
	args, err := parseArgs("HEALTHCHECK", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid HEALTHCHECK arguments %v", cmd.Args)
	}
	if len(args) == 0 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for HEALTHCHECK: %s", cmd.Args)
	}
	isNone := false
	var cmdArgs []string
	switch args[0] {
	case "NONE":
		if len(args) != 1 {
			return i.errorf(cmd.SourceLocation, "invalid arguments for HEALTHCHECK: %s", cmd.Args)
		}
		isNone = true
	case "CMD":
		if len(args) == 1 {
			return i.errorf(cmd.SourceLocation, "invalid number of arguments for HEALTHCHECK CMD: %s", cmd.Args)
		}
		cmdArgs = args[1:]
	default:
		if strings.HasPrefix(args[0], "[") {
			return i.errorf(cmd.SourceLocation, "exec form not yet supported for HEALTHCHECK CMD: %s", cmd.Args)
		}
		return i.errorf(cmd.SourceLocation, "invalid arguments for HEALTHCHECK: %s", cmd.Args)
	}
	for index, arg := range cmdArgs {
		expandedArg, err := i.expandArgs(ctx, arg, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand HEALTHCHECK arguments %s", arg)
		}
		cmdArgs[index] = expandedArg
	}
	err = i.converter.Healthcheck(ctx, isNone, cmdArgs, opts.Interval, opts.Timeout, opts.StartPeriod, opts.Retries)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply HEALTHCHECK")
	}
	return nil
}

func (i *Interpreter) handleWithDocker(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	if i.withDocker != nil {
		return i.errorf(cmd.SourceLocation, "cannot use WITH DOCKER within WITH DOCKER")
	}
	opts := withDockerOpts{}
	args, err := parseArgs("WITH DOCKER", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid WITH DOCKER arguments %v", cmd.Args)
	}
	if len(args) != 0 {
		return i.errorf(cmd.SourceLocation, "invalid WITH DOCKER arguments %v", args)
	}
	expandedPlatform, err := i.expandArgs(ctx, opts.Platform, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER platform %s", opts.Platform)
	}
	platform, err := i.converter.platr.Parse(expandedPlatform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", expandedPlatform)
	}
	for index, cf := range opts.ComposeFiles {
		expandedComposeFile, err := i.expandArgs(ctx, cf, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER compose: %s", cf)
		}
		opts.ComposeFiles[index] = expandedComposeFile
	}
	for index, cs := range opts.ComposeServices {
		expandedComposeService, err := i.expandArgs(ctx, cs, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER compose service: %s", cs)
		}
		opts.ComposeServices[index] = expandedComposeService
	}
	for index, load := range opts.Loads {
		expandedLoad, err := i.expandArgs(ctx, load, true, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER load: %s", load)
		}
		opts.Loads[index] = expandedLoad
	}
	expandedBuildArgs, err := i.expandArgsSlice(ctx, opts.BuildArgs, true, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER build args %v", opts.BuildArgs)
	}
	for index, p := range opts.Pulls {
		expandedPull, err := i.expandArgs(ctx, p, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER pull: %s", p)
		}
		opts.Pulls[index] = expandedPull
	}

	i.withDocker = &WithDockerOpt{
		ComposeFiles:    opts.ComposeFiles,
		ComposeServices: opts.ComposeServices,
	}
	for _, pullStr := range opts.Pulls {
		i.withDocker.Pulls = append(i.withDocker.Pulls, DockerPullOpt{
			ImageName: pullStr,
			Platform:  platform,
		})
	}
	for _, loadStr := range opts.Loads {
		loadImg, loadTarget, flagArgs, err := parseLoad(loadStr)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "parse load")
		}
		expandedFlagArgs, err := i.expandArgsSlice(ctx, flagArgs, true, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand WITH DOCKER load flag: %s", flagArgs)
		}
		parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "parse flag args")
		}
		loadBuildArgs := append(parsedFlagArgs, expandedBuildArgs...)

		allowPrivileged, err := i.getAllowPrivilegedTarget(loadTarget, opts.AllowPrivileged)
		if err != nil {
			return err
		}

		i.withDocker.Loads = append(i.withDocker.Loads, DockerLoadOpt{
			Target:          loadTarget,
			ImageName:       loadImg,
			Platform:        platform,
			BuildArgs:       loadBuildArgs,
			AllowPrivileged: allowPrivileged,
		})
	}
	return nil
}

func (i *Interpreter) handleAdd(ctx context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command ADD not yet supported")
}

func (i *Interpreter) handleStopsignal(ctx context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command STOPSIGNAL not yet supported")
}

func (i *Interpreter) handleOnbuild(ctx context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command ONBUILD not supported")
}

func (i *Interpreter) handleShell(ctx context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command SHELL not yet supported")
}

func (i *Interpreter) handleUserCommand(ctx context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command COMMAND not allowed in a target definition")
}

func (i *Interpreter) handleDo(ctx context.Context, cmd spec.Command) error {
	opts := doOpts{}
	args, err := parseArgs("DO", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid DO arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for DO: %s", args)
	}

	expandedFlagArgs, err := i.expandArgsSlice(ctx, args[1:], true, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand DO flags %v", args[1:])
	}
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}

	ucName, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand user command %v", args[0])
	}
	relCommand, err := domain.ParseCommand(ucName)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to parse user command reference %s", ucName)
	}

	allowPrivileged := i.allowPrivileged
	if relCommand.IsRemote() {
		allowPrivileged = i.allowPrivileged && opts.AllowPrivileged
	} else if opts.AllowPrivileged {
		i.console.Printf("the --allow-privileged flag has no effect when referencing a local target\n")
	}

	bc, resolvedAllowPrivileged, resolvedAllowPrivilegedSet, err := i.converter.ResolveReference(ctx, relCommand)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to resolve user command %s", ucName)
	}
	command := bc.Ref.(domain.Command)
	if resolvedAllowPrivilegedSet {
		allowPrivileged = allowPrivileged && resolvedAllowPrivileged
	}

	for _, uc := range bc.Earthfile.UserCommands {
		if uc.Name == command.Command {
			return i.handleDoUserCommand(ctx, command, relCommand, uc, cmd, parsedFlagArgs, allowPrivileged)
		}
	}
	return i.errorf(cmd.SourceLocation, "user command %s not found", ucName)
}

func (i *Interpreter) handleImport(ctx context.Context, cmd spec.Command) error {
	opts := importOpts{}
	args, err := parseArgs("IMPORT", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid IMPORT arguments %v", cmd.Args)
	}

	if len(args) != 1 && len(args) != 3 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for IMPORT: %s", args)
	}
	if len(args) == 3 && args[1] != "AS" {
		return i.errorf(cmd.SourceLocation, "invalid arguments for IMPORT: %s", args)
	}
	importStr, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand IMPORT %s", args[0])
	}
	var as string
	if len(args) == 3 {
		as, err = i.expandArgs(ctx, args[2], false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand IMPORT as: %s", as)
		}
	}
	isGlobal := (i.target.Target == "base")
	err = i.converter.Import(ctx, importStr, as, isGlobal, i.allowPrivileged, opts.AllowPrivileged)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply IMPORT")
	}
	return nil
}

func (i *Interpreter) handleProject(ctx context.Context, cmd spec.Command) error {
	projectVal, err := i.expandArgs(ctx, cmd.Args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand PROJECT %s", cmd.Args[0])
	}

	if !strings.Contains(projectVal, "/") {
		return i.errorf(cmd.SourceLocation, "format for PROJECT statement should be: <organization>/<project>")
	}

	parts := strings.Split(projectVal, "/")
	if len(parts) != 2 {
		return i.errorf(cmd.SourceLocation, "unexpected format for PROJECT statement")
	}

	i.converter.Project(ctx, parts[0], parts[1])

	return nil
}

func (i *Interpreter) handleCache(ctx context.Context, cmd spec.Command) error {
	if !i.converter.ftrs.UseCacheCommand {
		return i.errorf(cmd.SourceLocation, "the CACHE command is not supported in this version")
	}
	if len(cmd.Args) != 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for CACHE: %s", cmd.Args)
	}
	if i.local {
		return i.errorf(cmd.SourceLocation, "CACHE command not supported with LOCALLY")
	}
	dir, err := i.expandArgs(ctx, cmd.Args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand CACHE %s", cmd.Args[0])
	}
	if !path.IsAbs(dir) {
		dir = path.Clean(path.Join("/", i.converter.mts.Final.MainImage.Config.WorkingDir, dir))
	}
	if err := i.converter.Cache(ctx, dir); err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply CACHE")
	}
	return nil
}

func (i *Interpreter) handleHost(ctx context.Context, cmd spec.Command) error {
	if !i.converter.ftrs.UseHostCommand {
		return i.errorf(cmd.SourceLocation, "the HOST command is not supported in this version")
	}
	if len(cmd.Args) != 2 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for HOST: %s", cmd.Args)
	}
	if i.local {
		return i.errorf(cmd.SourceLocation, "HOST command not supported with LOCALLY")
	}

	host := cmd.Args[0]
	ipStr := cmd.Args[1]

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return i.errorf(cmd.SourceLocation, "invalid HOST ip %s", ipStr)
	}

	if err := i.converter.Host(ctx, host, ip); err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply HOST")
	}
	return nil
}

// ----------------------------------------------------------------------------

func (i *Interpreter) handleDoUserCommand(ctx context.Context, command domain.Command, relCommand domain.Command, uc spec.UserCommand, do spec.Command, buildArgs []string, allowPrivileged bool) error {
	if allowPrivileged && !i.allowPrivileged {
		return i.errorf(uc.SourceLocation, "invalid privileged in COMMAND") // this shouldn't happen, but check just in case
	}
	if len(uc.Recipe) == 0 || uc.Recipe[0].Command == nil || uc.Recipe[0].Command.Name != "COMMAND" {
		return i.errorf(uc.SourceLocation, "command recipes must start with COMMAND")
	}
	if len(uc.Recipe[0].Command.Args) > 0 {
		return i.errorf(uc.Recipe[0].SourceLocation, "COMMAND takes no arguments")
	}
	scopeName := fmt.Sprintf(
		"%s (%s line %d:%d)",
		command.StringCanonical(), do.SourceLocation.File, do.SourceLocation.StartLine, do.SourceLocation.StartColumn)
	err := i.converter.EnterScopeDo(ctx, command, baseTarget(relCommand), allowPrivileged, scopeName, buildArgs)
	if err != nil {
		return i.wrapError(err, uc.SourceLocation, "enter scope")
	}
	prevAllowPrivileged := i.allowPrivileged
	i.allowPrivileged = allowPrivileged
	i.stack = i.converter.StackString()
	err = i.handleBlock(ctx, uc.Recipe[1:])
	if err != nil {
		return err
	}
	err = i.converter.ExitScope(ctx)
	if err != nil {
		return i.wrapError(err, uc.SourceLocation, "exit scope")
	}
	i.stack = i.converter.StackString()
	i.allowPrivileged = prevAllowPrivileged
	return nil
}

// ----------------------------------------------------------------------------

func (i *Interpreter) expandArgsSlice(ctx context.Context, words []string, keepPlusEscape, async bool) ([]string, error) {
	ret := make([]string, 0, len(words))
	for _, word := range words {
		expanded, err := i.expandArgs(ctx, word, keepPlusEscape, async)
		if err != nil {
			return nil, err
		}
		ret = append(ret, expanded)
	}
	return ret, nil
}

func (i *Interpreter) errorf(sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	return Errorf(sl, i.stack, format, args...)
}

func (i *Interpreter) wrapError(cause error, sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	return WrapError(cause, sl, i.stack, format, args...)
}

func (i *Interpreter) pushOnlyErr(sl *spec.SourceLocation) error {
	return i.errorf(sl, "no non-push commands allowed after a --push")
}

func (i *Interpreter) expandArgs(ctx context.Context, word string, keepPlusEscape, async bool) (string, error) {
	runOpts := ConvertRunOpts{
		CommandName: "expandargs",
		Args:        nil, // this gets replaced whenver a shell-out is encountered
		Locally:     i.local,
		Transient:   !i.local,
		WithShell:   true,
	}

	ret, err := i.converter.ExpandArgs(ctx, runOpts, escapeSlashPlus(word), !async)
	if err != nil {
		if async && errors.Is(err, errShellOutNotPermitted) {
			return "", errCannotAsync
		}
		return "", err
	}
	if keepPlusEscape {
		return ret, nil
	}
	return unescapeSlashPlus(ret), nil
}

func escapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "\\\\+")
}

func unescapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "+")
}

func parseLoad(loadStr string) (image string, target string, extraArgs []string, err error) {
	words := strings.SplitN(loadStr, " ", 2)
	if len(words) == 0 {
		return "", "", nil, nil
	}
	firstWord := words[0]
	splitFirstWord := strings.SplitN(firstWord, "=", 2)
	if len(splitFirstWord) < 2 {
		// <target-name>
		// (will infer image name from SAVE IMAGE of that target)
		image = ""
		target = loadStr
	} else {
		// <image-name>=<target-name>
		image = splitFirstWord[0]
		if len(words) == 1 {
			target = splitFirstWord[1]
		} else {
			words[0] = splitFirstWord[1]
			target = strings.Join(words, " ")
		}
	}
	if isInParansForm(target) {
		target, extraArgs, err = parseParans(target)
		if err != nil {
			return "", "", nil, err
		}
	}
	return image, target, extraArgs, nil
}

func getArgsCopy(cmd spec.Command) []string {
	argsCopy := make([]string, len(cmd.Args))
	copy(argsCopy, cmd.Args)
	return argsCopy
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
		return "", nil, errors.Errorf("invalid build arg %s", splitArg)
	}
	name = splitArg[0]
	var value *string
	if len(splitArg) == 2 {
		value = &splitArg[1]
	}
	return name, value, nil
}

func isInParansForm(str string) bool {
	return (strings.HasPrefix(str, "\"(") && strings.HasSuffix(str, "\")")) ||
		(strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")"))
}

// parseParans turns "(+target --flag=something)" into "+target" and []string{"--flag=something"},
// or "\"(+target --flag=something)\"" into "+target" and []string{"--flag=something"}
func parseParans(str string) (string, []string, error) {
	if !isInParansForm(str) {
		return "", nil, errors.New("parans atom not in ( ... )")
	}
	if strings.HasPrefix(str, "\"(") {
		str = str[2 : len(str)-2] // remove \"( and )\"
	} else {
		str = str[1 : len(str)-1] // remove ( and )
	}
	var parts []string
	var part []rune
	nextEscaped := false
	inQuotes := false
	for _, char := range str {
		switch char {
		case '"':
			if !nextEscaped {
				inQuotes = !inQuotes
			}
			nextEscaped = false
		case '\\':
			nextEscaped = true
		case ' ', '\t', '\n':
			if !inQuotes && !nextEscaped {
				if len(part) > 0 {
					parts = append(parts, string(part))
					part = []rune{}
					nextEscaped = false
					continue
				} else {
					nextEscaped = false
					continue
				}
			}
			nextEscaped = false
		default:
			nextEscaped = false
		}
		part = append(part, char)
	}
	if nextEscaped {
		return "", nil, errors.New("unterminated escape sequence")
	}
	if inQuotes {
		return "", nil, errors.New("no ending quotes")
	}
	if len(part) > 0 {
		parts = append(parts, string(part))
	}

	if len(parts) < 1 {
		return "", nil, errors.New("invalid empty parans")
	}
	return parts[0], parts[1:], nil
}

// requiresShellOutOrCmdInvalid returns true if
// cmd requires shelling out via $(...), or if the cmd is invalid.
// This function is best-effort, and returns false on errors, errors
// will be handled during the synchronous portion of earthfile2llb handling.
func requiresShellOutOrCmdInvalid(s string) bool {
	var required bool
	shlex := shell.NewLex('\\')
	shlex.ShellOut = func(cmd string) (string, error) {
		required = true
		return "", nil
	}
	_, err := shlex.ProcessWordWithMap(s, map[string]string{})
	return required || err != nil
}

// isSafeAsyncBuildArgsKVStyle is used for "key=value" style buildargs
func isSafeAsyncBuildArgsKVStyle(args []string) bool {
	for _, arg := range args {
		_, v, _ := variables.ParseKeyValue(arg)
		if requiresShellOutOrCmdInvalid(v) {
			return false
		}
	}
	return true
}

// isSafeAsyncBuildArgs is used for "BUILD +target --key=value" style buildargs
func isSafeAsyncBuildArgs(args []string) bool {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			return false // malformed build arg
		}
		_, v, _ := variables.ParseKeyValue(arg[2:])
		if requiresShellOutOrCmdInvalid(v) {
			return false
		}
	}
	return true
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

func baseTarget(ref domain.Reference) domain.Target {
	return domain.Target{
		GitURL:    ref.GetGitURL(),
		Tag:       ref.GetTag(),
		ImportRef: ref.GetImportRef(),
		LocalPath: ref.GetLocalPath(),
		Target:    "base",
	}
}

func parseArgs(cmdName string, opts interface{}, args []string) ([]string, error) {
	processed := processParansAndQuotes(args)
	return flagutil.ParseArgs(cmdName, opts, processed)
}

func parseArgsWithValueModifier(cmdName string, opts interface{}, args []string, argumentModFunc flagutil.ArgumentModFunc) ([]string, error) {
	processed := processParansAndQuotes(args)
	return flagutil.ParseArgsWithValueModifier(cmdName, opts, processed, argumentModFunc)
}

// processParansAndQuotes takes in a slice of strings, and rearranges the slices
// depending on quotes and paranthesis.
// For example "hello ", "wor(", "ld)" becomes "hello ", "wor( ld)".
func processParansAndQuotes(args []string) []string {
	curQuote := rune(0)
	allowedQuotes := map[rune]rune{
		'"':  '"',
		'\'': '\'',
		'(':  ')',
	}
	ret := make([]string, 0, len(args))
	var newArg []rune
	for _, arg := range args {
		for _, char := range arg {
			newArg = append(newArg, char)
			if curQuote == 0 {
				_, isQuote := allowedQuotes[char]
				if isQuote {
					curQuote = char
				}
			} else {
				if char == allowedQuotes[curQuote] {
					curQuote = rune(0)
				}
			}
		}
		if curQuote == 0 {
			ret = append(ret, string(newArg))
			newArg = []rune{}
		} else {
			// Unterminated quote - join up two args into one.
			// Add a space between joined-up args.
			newArg = append(newArg, ' ')
		}
	}
	if curQuote != 0 {
		// Unterminated quote case.
		newArg = newArg[:len(newArg)-1] // remove last space
		ret = append(ret, string(newArg))
	}

	return ret
}
