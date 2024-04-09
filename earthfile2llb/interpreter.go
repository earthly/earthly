package earthfile2llb

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast/command"
	"github.com/earthly/earthly/ast/commandflag"
	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/internal/version"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/shell"
	"github.com/earthly/earthly/variables"

	"github.com/docker/go-connections/nat"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

const maxCommandRenameWarnings = 3

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

	withDocker    *WithDockerOpt
	withDockerRan bool

	parallelConversion bool
	console            conslogging.ConsoleLogger
	gitLookup          *buildcontext.GitLookup

	interactiveSaveFiles []debuggercommon.SaveFilesSettings
}

func newInterpreter(c *Converter, t domain.Target, allowPrivileged, parallelConversion bool, console conslogging.ConsoleLogger, gitLookup *buildcontext.GitLookup) *Interpreter {
	return &Interpreter{
		converter:          c,
		target:             t,
		allowPrivileged:    allowPrivileged,
		parallelConversion: parallelConversion,
		console:            console,
		gitLookup:          gitLookup,
	}
}

// Run interprets the commands in the given Earthfile AST, for a specific target.
func (i *Interpreter) Run(ctx context.Context, ef spec.Earthfile) (retErr error) {
	defer func() {
		if retErr != nil {
			i.converter.RecordTargetFailure(ctx, retErr)
		}
	}()
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

func (i *Interpreter) isPipelineTarget(_ context.Context, t spec.Target) bool {
	for _, stmt := range t.Recipe {
		if stmt.Command != nil && stmt.Command.Name == "PIPELINE" {
			return true
		}
	}
	return false
}

func (i *Interpreter) handleTarget(ctx context.Context, t spec.Target) error {
	ctx = ContextWithSourceLocation(ctx, t.SourceLocation)
	// Apply implicit FROM +base
	err := i.converter.From(ctx, "+base", platutil.DefaultPlatform, i.allowPrivileged, false, nil)
	if err != nil {
		return i.wrapError(err, t.SourceLocation, "apply FROM")
	}

	if i.isPipelineTarget(ctx, t) {
		return i.handlePipelineBlock(ctx, t.Name, t.Recipe)
	}

	return i.handleBlock(ctx, t.Recipe)
}

func (i *Interpreter) handleBlock(ctx context.Context, b spec.Block) error {
	prevWasArgLike := true // not exactly true, but makes the logic easier
	for index, stmt := range b {
		if i.parallelConversion && prevWasArgLike {
			err := i.handleBlockParallel(ctx, b, index)
			if err != nil {
				return err
			}
		}
		err := i.handleStatement(ctx, stmt)
		if err != nil {
			return err
		}
		prevWasArgLike = i.isArgLike(stmt.Command)

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
			case command.Arg, command.Locally, command.From, command.FromDockerfile:
				// Cannot do any further parallel builds - these commands need to be
				// executed to ensure that they don't impact the outcome. As such,
				// commands following these cannot be executed preemptively.
				return nil
			case command.Let, command.Set:
				if i.converter.ftrs.LetSetBlockParallel {
					// treat LET/SET the same as ARG if the feature flag is on,
					// otherwise fallthrough to handle the build
					return nil
				}
				fallthrough
			case command.Build:
				err := i.handleBuild(ctx, *stmt.Command, true)
				if err != nil {
					if errors.Is(err, errCannotAsync) {
						continue // no biggie
					}
					return err
				}
			case command.Copy:
				// TODO
			}
		} else if stmt.With != nil {
			switch stmt.With.Command.Name {
			case command.Docker:
				// TODO
			}
		} else if stmt.If != nil || stmt.For != nil || stmt.Wait != nil || stmt.Try != nil {
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
	ctx = ContextWithSourceLocation(ctx, stmt.SourceLocation)
	if stmt.Command != nil {
		return i.handleCommand(ctx, *stmt.Command)
	}
	if stmt.With != nil {
		return i.handleWith(ctx, *stmt.With)
	}
	if stmt.If != nil {
		return i.handleIf(ctx, *stmt.If)
	}
	if stmt.For != nil {
		return i.handleFor(ctx, *stmt.For)
	}
	if stmt.Wait != nil {
		return i.handleWait(ctx, *stmt.Wait)
	}
	if stmt.Try != nil {
		return i.handleTry(ctx, *stmt.Try)
	}
	return i.errorf(stmt.SourceLocation, "unexpected statement type")
}

func (i *Interpreter) handleCommand(ctx context.Context, cmd spec.Command) (err error) {
	// The AST should not be modified by any operation. This is a consistency check.
	argsCopy := flagutil.GetArgsCopy(cmd)
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

	ctx = ContextWithSourceLocation(ctx, cmd.SourceLocation)
	analytics.Count("cmd", cmd.Name)

	if i.isWith {
		switch cmd.Name {
		case command.Docker:
			return i.handleWithDocker(ctx, cmd)
		default:
			return i.errorf(cmd.SourceLocation, "unexpected WITH command %s", cmd.Name)
		}
	}

	switch cmd.Name {
	case command.From:
		return i.handleFrom(ctx, cmd)
	case command.Run:
		return i.handleRun(ctx, cmd)
	case command.FromDockerfile:
		return i.handleFromDockerfile(ctx, cmd)
	case command.Locally:
		return i.handleLocally(ctx, cmd)
	case command.Copy:
		return i.handleCopy(ctx, cmd)
	case command.SaveArtifact:
		return i.handleSaveArtifact(ctx, cmd)
	case command.SaveImage:
		return i.handleSaveImage(ctx, cmd)
	case command.Build:
		return i.handleBuild(ctx, cmd, false)
	case command.Workdir:
		return i.handleWorkdir(ctx, cmd)
	case command.User:
		return i.handleUser(ctx, cmd)
	case command.Cmd:
		return i.handleCmd(ctx, cmd)
	case command.Entrypoint:
		return i.handleEntrypoint(ctx, cmd)
	case command.Expose:
		return i.handleExpose(ctx, cmd)
	case command.Volume:
		return i.handleVolume(ctx, cmd)
	case command.Env:
		return i.handleEnv(ctx, cmd)
	case command.Arg:
		return i.handleArg(ctx, cmd)
	case command.Let:
		return i.handleLet(ctx, cmd)
	case command.Set:
		return i.handleSet(ctx, cmd)
	case command.Label:
		return i.handleLabel(ctx, cmd)
	case command.GitClone:
		return i.handleGitClone(ctx, cmd)
	case command.HealthCheck:
		return i.handleHealthcheck(ctx, cmd)
	case command.Add:
		return i.handleAdd(ctx, cmd)
	case command.StopSignal:
		return i.handleStopsignal(ctx, cmd)
	case command.OnBuild:
		return i.handleOnbuild(ctx, cmd)
	case command.Shell:
		return i.handleShell(ctx, cmd)
	case command.Command:
		return i.handleUserCommand(ctx, cmd)
	case command.Function:
		return i.handleFunction(ctx, cmd)
	case command.Do:
		return i.handleDo(ctx, cmd)
	case command.Import:
		return i.handleImport(ctx, cmd)
	case command.Cache:
		return i.handleCache(ctx, cmd)
	case command.Host:
		return i.handleHost(ctx, cmd)
	case command.Project:
		return i.handleProject(ctx, cmd)
	case command.Trigger:
		return i.handleTrigger(ctx, cmd)
	default:
		return i.errorf(cmd.SourceLocation, "unexpected command %q", cmd.Name)
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
	opts := commandflag.IfOpts{}
	args, err := flagutil.ParseArgsCleaned("IF", &opts, expression)
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
	opts := commandflag.NewForOpts()
	args, err := flagutil.ParseArgsCleaned("FOR", &opts, forArgs)
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

	if !i.converter.ftrs.ReferencedSaveOnly {
		return i.errorf(waitStmt.SourceLocation, "the WAIT command requires the --referenced-save-only feature")
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

func (i *Interpreter) handleTry(ctx context.Context, tryStmt spec.TryStatement) error {
	if !i.converter.ftrs.TryFinally {
		return i.errorf(tryStmt.SourceLocation, "the TRY/CATCH/FINALLY commands are not supported in this version")
	}

	if len(tryStmt.TryBody) == 0 {
		return i.errorf(tryStmt.SourceLocation, "TRY body is empty")
	}
	if len(tryStmt.TryBody) != 1 {
		return i.errorf(tryStmt.SourceLocation, "TRY body can (currently) only contain a single command")
	}

	if tryStmt.CatchBody != nil {
		return i.errorf(tryStmt.SourceLocation, "TRY/FINALLY doesn't (currently) support CATCH statements")
	}

	isRun := tryStmt.TryBody[0].Command != nil && tryStmt.TryBody[0].Command.Name == "RUN"
	isDocker := tryStmt.TryBody[0].With != nil && tryStmt.TryBody[0].With.Command.Name == "DOCKER"

	if isDocker {
		if len(tryStmt.TryBody[0].With.Body) != 1 {
			return i.errorf(tryStmt.SourceLocation, "TRY body can (currently) only contain a single command")
		}
	} else if !isRun {
		return i.errorf(tryStmt.SourceLocation, "TRY body must (currently) be a RUN command (or a RUN inside a WITH DOCKER)")
	}

	saveArtifacts := []debuggercommon.SaveFilesSettings{}
	if tryStmt.FinallyBody != nil {
		for _, cmd := range *tryStmt.FinallyBody {
			if cmd.Command == nil || cmd.Command.Name != "SAVE ARTIFACT" {
				return i.errorf(tryStmt.SourceLocation, "CATCH/FINALLY body only (currently) supports SAVE ARTIFACT ... AS LOCAL commands; got %s", cmd.Command.Name)
			}
			opts := commandflag.SaveArtifactOpts{}
			args, err := flagutil.ParseArgsCleaned("SAVE ARTIFACT", &opts, flagutil.GetArgsCopy(*cmd.Command))
			if err != nil {
				return i.wrapError(err, cmd.Command.SourceLocation, "invalid SAVE ARTIFACT arguments %v", cmd.Command.Args)
			}
			if opts.KeepTs || opts.KeepOwn || opts.SymlinkNoFollow || opts.Force {
				return i.wrapError(err, cmd.Command.SourceLocation, "only the SAVE ARTIFACT --if-exists option is allowed in a TRY/FINALLY block: %v", cmd.Command.Args)
			}
			saveFrom, _, saveAsLocalTo, ok := parseSaveArtifactArgs(args)
			if !ok {
				return i.errorf(cmd.Command.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Command.Args)
			}

			if strings.Contains(saveFrom, "*") {
				return i.errorf(cmd.Command.SourceLocation, "TRY/CATCH/FINALLY does not (currently) support wildcard SAVE ARTIFACT paths")
			}
			if saveAsLocalTo == "" {
				return i.errorf(cmd.Command.SourceLocation, "missing local name for SAVE ARTIFACT within TRY/CATCH/FINALLY")
			}
			if strings.Contains(saveAsLocalTo, "$") {
				return i.errorf(cmd.Command.SourceLocation, "TRY/CATCH/FINALLY does not (currently) support expanding args for SAVE ARTIFACT paths")
			}
			destIsDir := strings.HasSuffix(saveAsLocalTo, "/") || saveAsLocalTo == "."
			if destIsDir {
				saveAsLocalTo = path.Join(saveAsLocalTo, path.Base(saveFrom))
			}
			saveArtifacts = append(saveArtifacts, debuggercommon.SaveFilesSettings{
				Src:      saveFrom,
				Dst:      saveAsLocalTo,
				IfExists: opts.IfExists,
			})
		}
	}

	i.interactiveSaveFiles = saveArtifacts

	// process TRY body (i.e. perform the single RUN
	err := i.handleStatement(ctx, tryStmt.TryBody[0])
	if err != nil {
		return err
	}

	// process the FINALLY body (which will only happen when the try RUN succeeds, on failure
	// the SAVE ARTIFACTS are handled by the embedded debugger that was run under the try)
	if tryStmt.FinallyBody != nil {
		for _, cmd := range *tryStmt.FinallyBody {
			err := i.handleStatement(ctx, cmd)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Commands -------------------------------------------------------------------

func (i *Interpreter) handleFrom(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := commandflag.FromOpts{}
	args, err := flagutil.ParseArgsCleaned("FROM", &opts, flagutil.GetArgsCopy(cmd))
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

	if !i.converter.ftrs.PassArgs && opts.PassArgs {
		return i.errorf(cmd.SourceLocation, "the FROM --pass-args flag must be enabled with the VERSION --pass-args feature flag.")
	}

	i.local = false // FIXME https://github.com/earthly/earthly/issues/2044
	err = i.converter.From(ctx, imageName, platform, allowPrivileged, opts.PassArgs, expandedBuildArgs)
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
	opts := commandflag.RunOpts{}
	args, err := flagutil.ParseArgsWithValueModifierCleaned("RUN", &opts, flagutil.GetArgsCopy(cmd), i.flagValModifierFuncWithContext(ctx))
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

	var noNetwork bool
	if opts.Network != "" {
		if !i.converter.ftrs.NoNetwork {
			return i.errorf(cmd.SourceLocation, "the RUN --network=none flag must be enabled with the VERSION --no-network feature flag.")
		}
		if opts.Network != "none" {
			return i.errorf(cmd.SourceLocation, "invalid network value %s; only \"none\" is currently supported", opts.Network)
		}
		noNetwork = true
	}

	if opts.WithAWS && !i.converter.opt.Features.RunWithAWS {
		return i.errorf(cmd.SourceLocation, "RUN --aws requires the --run-with-aws feature flag")
	}

	if i.withDocker == nil {
		if opts.WithDocker {
			return i.errorf(cmd.SourceLocation, "--with-docker is obsolete. Please use WITH DOCKER ... RUN ... END instead")
		}
		opts := ConvertRunOpts{
			CommandName:          cmd.Name,
			Args:                 args,
			Locally:              i.local,
			Mounts:               opts.Mounts,
			Secrets:              opts.Secrets,
			WithShell:            withShell,
			WithEntrypoint:       opts.WithEntrypoint,
			Privileged:           opts.Privileged,
			NoNetwork:            noNetwork,
			Push:                 opts.Push,
			WithSSH:              opts.WithSSH,
			NoCache:              opts.NoCache,
			Interactive:          opts.Interactive,
			InteractiveKeep:      opts.InteractiveKeep,
			InteractiveSaveFiles: i.interactiveSaveFiles,
			WithAWSCredentials:   opts.WithAWS,
		}
		err = i.converter.Run(ctx, opts)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply RUN")
		}
		if opts.Push && !i.converter.ftrs.WaitBlock {
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
	opts := commandflag.FromDockerfileOpts{}
	args, err := flagutil.ParseArgsCleaned("FROM DOCKERFILE", &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid FROM DOCKERFILE arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for FROM DOCKERFILE")
	}

	if !i.converter.ftrs.AllowPrivilegedFromDockerfile && opts.AllowPrivileged {
		return i.errorf(cmd.SourceLocation, "the FROM DOCKERFILE --allow-privileged flag must be enabled with the VERSION --allow-privileged-from-dockerfile feature flag.")
	}
	allowPrivileged := opts.AllowPrivileged && i.allowPrivileged

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
	err = i.converter.FromDockerfile(ctx, path, expandedPath, expandedTarget, platform, allowPrivileged, expandedBuildArgs)
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
	opts := commandflag.CopyOpts{}
	args, err := flagutil.ParseArgsCleaned("COPY", &opts, flagutil.GetArgsCopy(cmd))
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
		if flagutil.IsInParamsForm(src) {
			// COPY (<src> <flag-args>) ...
			artifactStr, extraArgs, err := flagutil.ParseParams(src)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "parse params %s", src)
			}
			expandedArtifact, err := i.expandArgs(ctx, artifactStr, true, false)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to expand COPY artifact %s", artifactStr)
			}
			artifactSrc, parseErr = domain.ParseArtifact(expandedArtifact)
			if parseErr != nil {
				// Must parse in the params case.
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

			if i.converter.opt.LocalArtifactWhiteList.Exists(expandedSrc) {
				return i.errorf(cmd.SourceLocation, "unable to copy file %s, which has is outputted elsewhere by SAVE ARTIFACT AS LOCAL", expandedSrc)
			}

			srcs[index] = expandedSrc
			allArtifacts = false
		}
	}
	if !allClassical && !allArtifacts {
		return i.errorf(cmd.SourceLocation, "combining artifacts and build context arguments in a single COPY command is not allowed: %v", srcs)
	}
	if slices.ContainsFunc(strings.Split(dest, "/"), func(s string) bool {
		return s == "~" || strings.HasPrefix(s, "~")
	}) {
		i.console.Warnf(`destination path %q contains a "~" which does not expand to a home directory`, dest)
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

			if !i.converter.ftrs.PassArgs && opts.PassArgs {
				return i.errorf(cmd.SourceLocation, "the COPY --pass-args flag must be enabled with the VERSION --pass-args feature flag.")
			}

			if i.local {
				err = i.converter.CopyArtifactLocal(ctx, src, dest, platform, allowPrivileged, opts.PassArgs, srcBuildArgs, opts.IsDirCopy)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "copy artifact locally")
				}
			} else {
				err = i.converter.CopyArtifact(ctx, src, dest, platform, allowPrivileged, opts.PassArgs, srcBuildArgs, opts.IsDirCopy, opts.KeepTs, opts.KeepOwn, expandedChown, fileModeParsed, opts.IfExists, opts.SymlinkNoFollow)
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

func parseSaveArtifactArgs(args []string) (from, to, asLocal string, _ bool) {
	saveAsLocalTo := ""
	saveTo := "./"
	if len(args) >= 4 {
		if strings.Join(args[len(args)-3:len(args)-1], " ") == "AS LOCAL" {
			saveAsLocalTo = args[len(args)-1]
			if len(args) == 5 {
				saveTo = args[1]
			}
		} else {
			return "", "", "", false
		}
	} else if len(args) == 2 {
		saveTo = args[1]
	} else if len(args) == 0 || len(args) == 3 {
		return "", "", "", false
	}
	saveFrom := args[0]
	return saveFrom, saveTo, saveAsLocalTo, true
}

func (i *Interpreter) handleSaveArtifact(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.SaveArtifactOpts{}
	args, err := flagutil.ParseArgsCleaned("SAVE ARTIFACT", &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid SAVE ARTIFACT arguments %v", cmd.Args)
	}

	if len(args) == 0 {
		return i.errorf(cmd.SourceLocation, "no arguments provided to the SAVE ARTIFACT command")
	}
	if len(args) > 5 {
		return i.errorf(cmd.SourceLocation, "too many arguments provided to the SAVE ARTIFACT command: %v", cmd.Args)
	}
	saveFrom, saveTo, saveAsLocalTo, ok := parseSaveArtifactArgs(args)
	if !ok {
		return i.errorf(cmd.SourceLocation, "invalid arguments for SAVE ARTIFACT command: %v", cmd.Args)
	}

	saveFrom, err = i.expandArgs(ctx, saveFrom, false, false)
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

	if i.converter.ftrs.SaveArtifactKeepOwn {
		if opts.KeepOwn {
			fmt.Fprintf(os.Stderr, "Deprecation: SAVE ARTIFACT --keep-own is now applied by default, setting it no longer has any effect\n")
		}
		opts.KeepOwn = true
	}

	err = i.converter.SaveArtifact(ctx, saveFrom, expandedSaveTo, expandedSaveAsLocalTo, opts.KeepTs, opts.KeepOwn, opts.IfExists, opts.SymlinkNoFollow, opts.Force, i.pushOnlyAllowed)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
	}
	return nil
}

func (i *Interpreter) handleSaveImage(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.SaveImageOpts{}
	args, err := flagutil.ParseArgsCleaned("SAVE IMAGE", &opts, flagutil.GetArgsCopy(cmd))
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

	labels := map[string]string{
		"dev.earthly.version":  version.Version,
		"dev.earthly.git-sha":  version.GitSha,
		"dev.earthly.built-by": version.BuiltBy,
	}
	err = i.converter.Label(ctx, labels)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to create dev.earthly.* labels during SAVE IMAGE")
	}

	err = i.converter.SaveImage(ctx, imageNames, opts.Push, opts.Insecure, opts.CacheHint, opts.CacheFrom, opts.NoManifestList)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "save image")
	}
	if opts.Push && !i.converter.ftrs.WaitBlock {
		i.pushOnlyAllowed = true
	}
	return nil
}

func (i *Interpreter) handleBuild(ctx context.Context, cmd spec.Command, async bool) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := commandflag.BuildOpts{}
	args, err := flagutil.ParseArgsCleaned("BUILD", &opts, flagutil.GetArgsCopy(cmd))
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
	// Expand wildcards into a set of BUILD spec.Command's, one for each discovered Earthfile.
	if strings.Contains(fullTargetName, "*") {
		if !i.converter.ftrs.WildcardBuilds {
			return i.errorf(cmd.SourceLocation, "wildcard BUILD commands are not enabled")
		}
		return i.handleWildcardBuilds(ctx, fullTargetName, cmd, async)
	}
	platformsSlice := make([]platutil.Platform, 0, len(opts.Platforms))
	for index, p := range opts.Platforms {
		expandedPlatform, err := i.expandArgs(ctx, p, false, async)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand BUILD platform %s", p)
		}
		opts.Platforms[index] = expandedPlatform
		platform, err := i.converter.platr.Parse(expandedPlatform)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	asyncSafeArgs := isSafeAsyncBuildArgsKVStyle(opts.BuildArgs) && isSafeAsyncBuildArgs(args[1:])
	if async && (!asyncSafeArgs || opts.AutoSkip) {
		return errCannotAsync
	}
	if i.local && !asyncSafeArgs {
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

	crossProductBuildArgs, err := flagutil.BuildArgMatrix(expandedBuildArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "build arg matrix")
	}

	allowPrivileged, err := i.getAllowPrivilegedTarget(fullTargetName, opts.AllowPrivileged)
	if err != nil {
		return err
	}

	if !i.converter.ftrs.PassArgs && opts.PassArgs {
		return i.errorf(cmd.SourceLocation, "the BUILD --pass-args flag must be enabled with the VERSION --pass-args feature flag.")
	}

	if !i.converter.ftrs.BuildAutoSkip && opts.AutoSkip {
		return i.errorf(cmd.SourceLocation, "the BUILD --auto-skip flag must be enabled with the VERSION --build-auto-skip feature flag.")
	}

	for _, buildArgs := range crossProductBuildArgs {
		saveHashFn := func() {}
		if opts.AutoSkip {
			skip, fn, err := i.converter.checkAutoSkip(ctx, fullTargetName, allowPrivileged, opts.PassArgs, buildArgs)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "failed to determine whether target can be skipped")
			}
			if skip {
				continue
			}
			saveHashFn = fn
		}
		for _, platform := range platformsSlice {
			if async {
				err := i.converter.BuildAsync(ctx, fullTargetName, platform, allowPrivileged, opts.PassArgs, buildArgs, buildCmd, nil, nil)
				if err != nil {
					return i.wrapError(err, cmd.SourceLocation, "apply BUILD %s", fullTargetName)
				}
				continue
			}
			err := i.converter.Build(ctx, fullTargetName, platform, allowPrivileged, opts.PassArgs, buildArgs)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "apply BUILD %s", fullTargetName)
			}
			saveHashFn()
		}
	}
	return nil
}

func (i *Interpreter) handleWildcardBuilds(ctx context.Context, fullTargetName string, cmd spec.Command, async bool) error {

	children, err := i.converter.ExpandWildcard(ctx, fullTargetName, cmd)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand wildcard BUILD %q", fullTargetName)
	}

	for _, child := range children {
		if err := i.handleBuild(ctx, child, async); err != nil {
			return err
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
	cmdArgs := flagutil.GetArgsCopy(cmd)
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
	entArgs := flagutil.GetArgsCopy(cmd)
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
	ports := flagutil.GetArgsCopy(cmd)
	for index, port := range ports {
		expandedPort, err := i.expandArgs(ctx, port, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand EXPOSE %s", port)
		}

		ports[index] = expandedPort
	}

	// Dockerfile syntax allows defining host bindings; however, they are ignored when generating the image
	// see: https://github.com/earthly/buildkit/blob/dad0cead57a2d92d43e44c9212153ffe53d9ebc9/frontend/dockerfile/dockerfile2llb/convert.go#L1207
	ps, _, err := nat.ParsePortSpecs(ports)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to parse EXPOSE command")
	}
	parsedPorts := []string{}
	for p := range ps {
		parsedPorts = append(parsedPorts, string(p))
	}

	err = i.converter.Expose(ctx, parsedPorts)
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
	volumes := flagutil.GetArgsCopy(cmd)
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

func (i *Interpreter) handleArg(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts, key, valueOrNil, err := flagutil.ParseArgArgs(ctx, cmd, i.isBase, i.converter.ftrs.ExplicitGlobal)
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

func (i *Interpreter) handleLet(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	var opts commandflag.LetOpts
	argsCpy := flagutil.GetArgsCopy(cmd)
	args, err := flagutil.ParseArgsCleaned("LET", &opts, argsCpy)
	if err != nil {
		return errors.Wrap(err, "failed to parse LET args")
	}
	if len(args) != 3 || args[1] != "=" {
		return hint.Wrap(flagutil.ErrInvalidSyntax, "LET requires a variable assignment, e.g. LET foo = bar")
	}

	key := args[0]
	baseVal := args[2]
	val, err := i.expandArgs(ctx, baseVal, true, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand LET value %q", baseVal)
	}

	err = i.converter.Let(ctx, key, val)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply LET")
	}
	return nil
}

func parseSetArgs(ctx context.Context, cmd spec.Command) (name, value string, _ error) {
	var opts commandflag.SetOpts
	argsCpy := flagutil.GetArgsCopy(cmd)
	args, err := flagutil.ParseArgsCleaned("SET", &opts, argsCpy)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to parse SET args")
	}
	if len(args) != 3 {
		return "", "", flagutil.ErrInvalidSyntax
	}
	if args[1] != "=" {
		return "", "", flagutil.ErrInvalidSyntax
	}
	return args[0], args[2], nil
}

func (i *Interpreter) handleSet(ctx context.Context, cmd spec.Command) error {
	if !i.converter.ftrs.ArgScopeSet {
		return errors.New("unknown command SET")
	}
	key, value, err := parseSetArgs(ctx, cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to parse SET arguments")
	}
	newVal, err := i.expandArgs(ctx, value, true, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand SET %s", value)
	}
	return i.converter.UpdateArg(ctx, key, newVal, i.isBase)
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
			if strings.HasPrefix(key, "dev.earthly.") {
				return i.wrapError(err, cmd.SourceLocation, "LABEL keys starting with \"dev.earthly.\" are reserved")
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
	opts := commandflag.GitCloneOpts{}
	args, err := flagutil.ParseArgsCleaned("GIT CLONE", &opts, flagutil.GetArgsCopy(cmd))
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

	convertedGitURL, _, sshCommand, err := i.gitLookup.ConvertCloneURL(gitURL)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to use %v with configured earthly credentials from ~/.earthly/config.yml", cmd.Args)
	}

	err = i.converter.GitClone(ctx, convertedGitURL, sshCommand, gitBranch, gitCloneDest, opts.KeepTs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "git clone")
	}
	return nil
}

func (i *Interpreter) handleHealthcheck(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := commandflag.HealthCheckOpts{}
	args, err := flagutil.ParseArgsCleaned("HEALTHCHECK", &opts, flagutil.GetArgsCopy(cmd))
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
	err = i.converter.Healthcheck(ctx, isNone, cmdArgs, opts.Interval, opts.Timeout, opts.StartPeriod, opts.Retries, opts.StartInterval)
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
	opts := commandflag.WithDockerOpts{}
	args, err := flagutil.ParseArgsCleaned("WITH DOCKER", &opts, flagutil.GetArgsCopy(cmd))
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
		ComposeFiles:          opts.ComposeFiles,
		ComposeServices:       opts.ComposeServices,
		TryCatchSaveArtifacts: i.interactiveSaveFiles,
	}
	for _, pullStr := range opts.Pulls {
		i.withDocker.Pulls = append(i.withDocker.Pulls, DockerPullOpt{
			ImageName: pullStr,
			Platform:  platform,
		})
	}
	for _, loadStr := range opts.Loads {
		loadImg, loadTarget, flagArgs, err := flagutil.ParseLoad(loadStr)
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

		if !i.converter.ftrs.PassArgs && opts.PassArgs {
			return i.errorf(cmd.SourceLocation, "the WITH DOCKER --pass-args flag must be enabled with the VERSION --pass-args feature flag.")
		}

		i.withDocker.Loads = append(i.withDocker.Loads, DockerLoadOpt{
			Target:          loadTarget,
			ImageName:       loadImg,
			Platform:        platform,
			BuildArgs:       loadBuildArgs,
			AllowPrivileged: allowPrivileged,
			PassArgs:        opts.PassArgs,
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

func (i *Interpreter) handleUserCommand(_ context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command COMMAND not allowed in a target definition")
}

func (i *Interpreter) handleFunction(_ context.Context, cmd spec.Command) error {
	return i.errorf(cmd.SourceLocation, "command FUNCTION not allowed in a target definition")
}

func (i *Interpreter) handleDo(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.DoOpts{}
	args, err := flagutil.ParseArgsCleaned("DO", &opts, flagutil.GetArgsCopy(cmd))
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

	if !i.converter.ftrs.PassArgs && opts.PassArgs {
		return i.errorf(cmd.SourceLocation, "the DO --pass-args flag must be enabled with the VERSION --pass-args feature flag.")
	}

	for _, uc := range bc.Earthfile.Functions {
		if uc.Name == command.Command {
			sourceFilePath := bc.Ref.ProjectCanonical() + "/Earthfile"
			return i.handleDoFunction(ctx, command, relCommand, uc, cmd, parsedFlagArgs, allowPrivileged, opts.PassArgs, sourceFilePath, bc.Features.UseFunctionKeyword)
		}
	}
	return i.errorf(cmd.SourceLocation, "user command %s not found", ucName)
}

func (i *Interpreter) handleImport(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.ImportOpts{}
	args, err := flagutil.ParseArgsCleaned("IMPORT", &opts, flagutil.GetArgsCopy(cmd))
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
	// Note: Expanding args for PROJECT is not allowed. The value needs to be
	// lifted straight from the AST.
	projectVal := cmd.Args[0]
	parts := strings.Split(projectVal, "/")
	if len(parts) != 2 {
		return i.errorf(cmd.SourceLocation, "unexpected format for PROJECT statement, should be: <organization>/<project>")
	}

	err := i.converter.Project(ctx, parts[0], parts[1])
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to process PROJECT")
	}

	return nil
}

func (i *Interpreter) handlePipelineBlock(ctx context.Context, name string, block spec.Block) error {
	if len(block) == 0 {
		return errors.New("pipeline targets require sub-commands")
	}

	if block[0].Command == nil || block[0].Command.Name != "PIPELINE" {
		return i.errorf(block[0].Command.SourceLocation, "PIPELINE must be the first command in a pipeline target")
	}

	for _, stmt := range block {
		if stmt.Command == nil {
			return errors.New("pipeline targets do not support IF, WITH, FOR, or WAIT commands")
		}
		cmd := *stmt.Command
		ctx = ContextWithSourceLocation(ctx, cmd.SourceLocation)
		var err error
		switch cmd.Name {
		case command.Pipeline:
			err = i.handlePipeline(ctx, cmd)
		case command.Trigger:
			err = i.handleTrigger(ctx, cmd)
		case command.Arg:
			err = i.handleArg(ctx, cmd)
		case command.Build:
			err = i.handleBuild(ctx, cmd, false)
		default:
			return i.errorf(cmd.SourceLocation, "pipeline targets only support PIPELINE, TRIGGER, ARG, and BUILD commands")
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) handlePipeline(ctx context.Context, cmd spec.Command) error {

	if len(cmd.Args) > 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of PIPELINE arguments")
	}

	var opts commandflag.PipelineOpts
	_, err := flagutil.ParseArgsCleaned("PIPELINE", &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid PIPELINE arguments")
	}

	return i.converter.Pipeline(ctx)
}

func (i *Interpreter) handleTrigger(ctx context.Context, cmd spec.Command) error {

	if len(cmd.Args) < 1 {
		return i.errorf(cmd.SourceLocation, "TRIGGER requires at least 1 argument")
	}

	switch cmd.Args[0] {
	case "manual":
		if len(cmd.Args) != 1 {
			return i.errorf(cmd.SourceLocation, "invalid argument")
		}
	case "pr", "push":
		if len(cmd.Args) != 2 {
			return i.errorf(cmd.SourceLocation, "'pr' and 'push' triggers require a branch name")
		}
	default:
		return i.errorf(cmd.SourceLocation, "valid triggers include: 'manual', 'pr', or 'push'")
	}

	return nil
}

func (i *Interpreter) handleCache(ctx context.Context, cmd spec.Command) error {
	if !i.converter.ftrs.UseCacheCommand {
		return i.errorf(cmd.SourceLocation, "the CACHE command is not supported in this version")
	}
	opts := commandflag.CacheOpts{}
	args, err := flagutil.ParseArgsCleaned("CACHE", &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid CACHE arguments %v", cmd.Args)
	}
	if len(args) != 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for CACHE: %s", args)
	}
	if i.local {
		return i.errorf(cmd.SourceLocation, "CACHE command not supported with LOCALLY")
	}
	dir, err := i.expandArgs(ctx, args[0], false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand CACHE directory %s", args[0])
	}
	expandedMode, err := i.expandArgs(ctx, opts.Mode, false, false)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "failed to expand CACHE mode %s", opts.Mode)
	} else {
		opts.Mode = expandedMode
	}
	if !path.IsAbs(dir) {
		dir = path.Clean(path.Join("/", i.converter.mts.Final.MainImage.Config.WorkingDir, dir))
	}
	if opts.ID != "" {
		opts.ID, err = i.expandArgs(ctx, opts.ID, false, false)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "failed to expand CACHE id %s", opts.ID)
		}
	}
	if err := i.converter.Cache(ctx, dir, opts); err != nil {
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
	host, err := i.expandArgs(ctx, cmd.Args[0], true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand host name for HOST: %s", cmd.Args)
	}
	ipStr, err := i.expandArgs(ctx, cmd.Args[1], true, false)
	if err != nil {
		return i.errorf(cmd.SourceLocation, "unable to expand IP addr for HOST: %s", cmd.Args)
	}
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

func (i *Interpreter) handleDoFunction(ctx context.Context, command domain.Command, relCommand domain.Command, uc spec.Function, do spec.Command, buildArgs []string, allowPrivileged, passArgs bool, sourceLocationFile string, useFunctionCmd bool) error {
	cmdName := "FUNCTION"
	if !useFunctionCmd {
		cmdName = "COMMAND"
	}
	if allowPrivileged && !i.allowPrivileged {
		return i.errorf(uc.SourceLocation, "invalid privileged in %s", cmdName) // this shouldn't happen, but check just in case
	}
	if len(uc.Recipe) == 0 || uc.Recipe[0].Command == nil || uc.Recipe[0].Command.Name != cmdName {
		return i.errorf(uc.SourceLocation, "%s recipes must start with %s", strings.ToLower(cmdName), cmdName)
	}
	if !useFunctionCmd && len(i.converter.opt.FilesWithCommandRenameWarning) < maxCommandRenameWarnings && !i.converter.opt.FilesWithCommandRenameWarning[sourceLocationFile] {
		i.console.Printf(
			`Note that the COMMAND keyword will be replaced by FUNCTION starting with VERSION 0.8.
To start using the FUNCTION keyword now (experimental) please use VERSION --use-function-keyword 0.7 in %s. Note that switching now may cause breakages for your colleagues if they are using older Earthly versions.
`, sourceLocationFile)
		i.converter.opt.FilesWithCommandRenameWarning[sourceLocationFile] = true
	}
	if len(uc.Recipe[0].Command.Args) > 0 {
		return i.errorf(uc.Recipe[0].SourceLocation, "%s takes no arguments", cmdName)
	}
	scopeName := fmt.Sprintf(
		"%s (%s line %d:%d)",
		command.StringCanonical(), do.SourceLocation.File, do.SourceLocation.StartLine, do.SourceLocation.StartColumn)
	err := i.converter.EnterScopeDo(ctx, command, baseTarget(relCommand), allowPrivileged, passArgs, scopeName, buildArgs)
	if err != nil {
		return i.wrapError(err, uc.SourceLocation, "enter scope")
	}
	prevAllowPrivileged := i.allowPrivileged
	i.allowPrivileged = allowPrivileged
	err = i.handleBlock(ctx, uc.Recipe[1:])
	if err != nil {
		return err
	}
	err = i.converter.ExitScope(ctx)
	if err != nil {
		return i.wrapError(err, uc.SourceLocation, "exit scope")
	}
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

func (i *Interpreter) stack() string {
	return i.converter.varCollection.StackString()
}

func (i *Interpreter) errorf(sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	targetID := i.converter.mts.Final.ID
	return Errorf(sl, targetID, i.stack(), format, args...)
}

func (i *Interpreter) wrapError(cause error, sl *spec.SourceLocation, format string, args ...interface{}) *InterpreterError {
	targetID := i.converter.mts.Final.ID
	return WrapError(cause, sl, targetID, i.stack(), format, args...)
}

func (i *Interpreter) pushOnlyErr(sl *spec.SourceLocation) error {
	return i.errorf(sl, "no non-push commands allowed after a --push")
}

func (i *Interpreter) expandArgs(ctx context.Context, word string, keepPlusEscape, async bool) (string, error) {
	runOpts := ConvertRunOpts{
		CommandName: "expandargs",
		Args:        nil, // this gets replaced whenever a shell-out is encountered
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

// isArgLike returns true if the command is ARG/LET/SET
func (i *Interpreter) isArgLike(cmd *spec.Command) bool {
	if cmd == nil {
		return false
	}

	switch cmd.Name {
	case command.Let, command.Set:
		return i.converter.ftrs.LetSetBlockParallel
	case command.Arg:
		return true
	}

	return false
}

func escapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "\\\\+")
}

func unescapeSlashPlus(str string) string {
	// TODO: This is not entirely correct in a string like "\\\\+".
	return strings.ReplaceAll(str, "\\+", "+")
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
	_, err := shlex.ProcessWordWithMap(s, map[string]string{}, variables.ShellOutEnvs)
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
