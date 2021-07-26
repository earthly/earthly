package earthfile2llb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/variables"

	flags "github.com/jessevdk/go-flags"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
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
	parallelism        *semaphore.Weighted
	parallelErrChan    chan error
	console            conslogging.ConsoleLogger
	gitLookup          *buildcontext.GitLookup
}

func newInterpreter(c *Converter, t domain.Target, allowPrivileged, parallelConversion bool, parallelism *semaphore.Weighted, console conslogging.ConsoleLogger, gitLookup *buildcontext.GitLookup) *Interpreter {
	return &Interpreter{
		converter:          c,
		target:             t,
		stack:              c.StackString(),
		allowPrivileged:    allowPrivileged,
		parallelism:        parallelism,
		parallelConversion: parallelConversion,
		parallelErrChan:    make(chan error),
		console:            console,
		gitLookup:          gitLookup,
	}
}

// Run interprets the commands in the given Earthfile AST, for a specific target.
func (i *Interpreter) Run(ctx context.Context, ef spec.Earthfile) (err error) {
	done := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer close(done)
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
	})
	eg.Go(func() error {
		select {
		case parallelErr := <-i.parallelErrChan:
			return parallelErr
		case <-done:
			return nil
		case <-ctx.Done():
			return nil
		}
	})
	return eg.Wait()
}

func (i *Interpreter) handleTarget(ctx context.Context, t spec.Target) error {
	// Apply implicit FROM +base
	err := i.converter.From(ctx, "+base", nil, i.allowPrivileged, nil)
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
	// Look ahead of the execution and fire off asynchronous builds for mentioned targets,
	// as long as they don't have variable args $(...).
	for index := startIndex; index < len(b); index++ {
		stmt := b[index]
		if stmt.Command != nil {
			switch stmt.Command.Name {
			case "ARG":
				// Cannot do any further parallel builds - args may change the outcome.
				return nil
			case "BUILD":
				err := i.handleBuild(ctx, *stmt.Command, true)
				if err != nil {
					if errors.Is(err, errCannotAsync) {
						continue // no biggie
					}
					return err
				}
			case "FROM":
				// TODO
			case "COPY":
				// TODO
			case "FROM DOCKERFILE":
				// TODO
			}
		} else if stmt.With != nil {
			switch stmt.With.Command.Name {
			case "DOCKER":
				// TODO
			}
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
	args, err := flagutil.ParseArgs("IF", &opts, expression)
	if err != nil {
		return false, i.wrapError(err, sl, "invalid IF arguments %v", expression)
	}
	withShell := !execMode

	for index, s := range opts.Secrets {
		opts.Secrets[index] = i.expandArgs(s, true)
	}
	for index, m := range opts.Mounts {
		opts.Mounts[index] = i.expandArgs(m, false)
	}
	// Note: Not expanding args for the expression itself, as that will be take care of by the shell.

	var exitCode int
	commandName := "IF"
	if i.local {
		if len(opts.Mounts) > 0 {
			return false, i.errorf(sl, "mounts are not supported in combination with the LOCALLY directive")
		}
		if opts.WithSSH {
			return false, i.errorf(sl, "the --ssh flag has no effect when used with the  LOCALLY directive")
		}
		if opts.Privileged {
			return false, i.errorf(sl, "the --privileged flag has no effect when used with the LOCALLY directive")
		}
		if opts.NoCache {
			return false, i.errorf(sl, "the --no-cache flag has no effect when used with the LOCALLY directive")
		}

		// TODO these should be supported, but haven't yet been implemented
		if len(opts.Secrets) > 0 {
			return false, i.errorf(sl, "secrets need to be implemented for the LOCALLY directive")
		}

		exitCode, err = i.converter.RunLocalExitCode(ctx, commandName, args)
		if err != nil {
			return false, i.wrapError(err, sl, "apply RUN")
		}
	} else {
		exitCode, err = i.converter.RunExitCode(
			ctx, commandName, args, opts.Mounts, opts.Secrets, opts.Privileged,
			withShell, opts.WithSSH, opts.NoCache)
		if err != nil {
			return false, i.wrapError(err, sl, "apply IF")
		}
	}
	return (exitCode == 0), nil
}

// Commands -------------------------------------------------------------------

func (i *Interpreter) handleFrom(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := fromOpts{}
	args, err := flagutil.ParseArgs("FROM", &opts, getArgsCopy(cmd))
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
	imageName := i.expandArgs(args[0], true)
	opts.Platform = i.expandArgs(opts.Platform, false)
	platform, err := llbutil.ParsePlatform(opts.Platform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", opts.Platform)
	}
	expandedBuildArgs := i.expandArgsSlice(opts.BuildArgs, true)
	expandedFlagArgs := i.expandArgsSlice(args[1:], true)
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

func (i *Interpreter) flagValModifier(flagName string, flagOpt *flags.Option, flagVal *string) *string {
	if flagOpt.IsBool() && flagVal != nil {
		newFlag := i.expandArgs(*flagVal, false)
		return &newFlag
	}
	return flagVal
}

func (i *Interpreter) handleRun(ctx context.Context, cmd spec.Command) error {
	if len(cmd.Args) < 1 {
		return i.errorf(cmd.SourceLocation, "not enough arguments for RUN")
	}
	opts := runOpts{}
	args, err := flagutil.ParseArgsWithValueModifier("RUN", &opts, getArgsCopy(cmd), i.flagValModifier)
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
		opts.Secrets[index] = i.expandArgs(s, true)
	}
	for index, m := range opts.Mounts {
		opts.Mounts[index] = i.expandArgs(m, false)
	}
	// Note: Not expanding args for the run itself, as that will be take care of by the shell.

	if i.local {
		if len(opts.Mounts) > 0 {
			return i.errorf(cmd.SourceLocation, "mounts are not supported in combination with the LOCALLY directive")
		}
		if opts.WithSSH {
			return i.errorf(cmd.SourceLocation, "the --ssh flag has no effect when used with the  LOCALLY directive")
		}
		if opts.Privileged {
			return i.errorf(cmd.SourceLocation, "the --privileged flag has no effect when used with the LOCALLY directive")
		}
		if opts.NoCache {
			return i.errorf(cmd.SourceLocation, "the --no-cache flag has no effect when used with the LOCALLY directive")
		}
		if opts.Interactive {
			// I mean its literally just your terminal but with extra steps. No reason to support this?
			return i.errorf(cmd.SourceLocation, "the --interactive flag is not supported in combination with the LOCALLY directive")
		}
		if opts.InteractiveKeep {
			// I mean its literally just your terminal but with extra steps. No reason to support this?
			return i.errorf(cmd.SourceLocation, "the --interactive-keep flag is not supported in combination with the LOCALLY directive")
		}

		// TODO these should be supported, but haven't yet been implemented
		if len(opts.Secrets) > 0 {
			return i.errorf(cmd.SourceLocation, "secrets need to be implemented for the LOCALLY directive")
		}

		if i.withDocker != nil {
			if opts.Push {
				return i.errorf(cmd.SourceLocation, "RUN --push not allowed in WITH DOCKER")
			}
			if i.withDockerRan {
				return i.errorf(cmd.SourceLocation, "only one RUN command allowed in WITH DOCKER")
			}
			i.withDockerRan = true
			err = i.converter.WithDockerRunLocal(ctx, args, *i.withDocker)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "with docker run")
			}
			return nil
		}

		err = i.converter.RunLocal(ctx, args, opts.Push)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply RUN")
		}
		return nil
	}

	if opts.Privileged && !i.allowPrivileged {
		return i.errorf(cmd.SourceLocation, "Permission denied: unwilling to run privileged command; did you reference a remote Earthfile without the --allow-privileged flag?")
	}

	if i.withDocker == nil {
		err = i.converter.Run(
			ctx, args, opts.Mounts, opts.Secrets, opts.Privileged, opts.WithEntrypoint, opts.WithDocker,
			withShell, opts.Push, opts.WithSSH, opts.NoCache, opts.Interactive, opts.InteractiveKeep)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply RUN")
		}
		if opts.Push {
			i.pushOnlyAllowed = true
		}
	} else {
		if opts.Push {
			return i.errorf(cmd.SourceLocation, "RUN --push not allowed in WITH DOCKER")
		}
		if i.withDockerRan {
			return i.errorf(cmd.SourceLocation, "only one RUN command allowed in WITH DOCKER")
		}
		i.withDockerRan = true
		i.withDocker.Mounts = opts.Mounts
		i.withDocker.Secrets = opts.Secrets
		i.withDocker.WithShell = withShell
		i.withDocker.WithEntrypoint = opts.WithEntrypoint
		i.withDocker.NoCache = opts.NoCache
		i.withDocker.Interactive = opts.Interactive
		i.withDocker.interactiveKeep = opts.InteractiveKeep
		err = i.converter.WithDockerRun(ctx, args, *i.withDocker)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "with docker run")
		}
	}
	return nil
}

func (i *Interpreter) handleFromDockerfile(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	opts := fromDockerfileOpts{}
	args, err := flagutil.ParseArgs("FROM DOCKERFILE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid FROM DOCKERFILE arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for FROM DOCKERFILE")
	}
	path := i.expandArgs(args[0], false)
	_, parseErr := domain.ParseArtifact(path)
	if parseErr != nil {
		// Treat as context path, not artifact path.
		path = i.expandArgs(args[0], false)
	}
	expandedBuildArgs := i.expandArgsSlice(opts.BuildArgs, true)
	expandedFlagArgs := i.expandArgsSlice(args[1:], true)
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}
	expandedBuildArgs = append(parsedFlagArgs, expandedBuildArgs...)
	opts.Platform = i.expandArgs(opts.Platform, false)
	platform, err := llbutil.ParsePlatform(opts.Platform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", opts.Platform)
	}
	opts.Path = i.expandArgs(opts.Path, false)
	opts.Target = i.expandArgs(opts.Target, false)
	i.local = false
	err = i.converter.FromDockerfile(ctx, path, opts.Path, opts.Target, platform, expandedBuildArgs)
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

	workingDir, err := filepath.Abs(filepath.Dir(cmd.SourceLocation.File))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to get abs path in LOCALLY")
	}

	i.local = true
	err = i.converter.Locally(ctx, workingDir, nil)
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
	args, err := flagutil.ParseArgs("COPY", &opts, getArgsCopy(cmd))
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
	dest := i.expandArgs(args[len(args)-1], false)
	expandedBuildArgs := i.expandArgsSlice(opts.BuildArgs, true)
	opts.Chown = i.expandArgs(opts.Chown, false)
	opts.Platform = i.expandArgs(opts.Platform, false)
	platform, err := llbutil.ParsePlatform(opts.Platform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", opts.Platform)
	}

	allClassical := true
	allArtifacts := true
	for index, src := range srcs {
		var artifactSrc domain.Artifact
		var parseErr error
		if strings.HasPrefix(src, "(") && strings.HasSuffix(src, ")") {
			// COPY (<src> <flag-args>) ...
			artifactStr, extraArgs, err := parseParans(src)
			if err != nil {
				return i.wrapError(err, cmd.SourceLocation, "parse parans %s", src)
			}
			artifactSrc, parseErr = domain.ParseArtifact(i.expandArgs(artifactStr, true))
			if parseErr != nil {
				// Must parse in the parans case.
				return i.wrapError(err, cmd.SourceLocation, "parse artifact")
			}
			srcFlagArgs[index] = extraArgs
		} else {
			artifactSrc, parseErr = domain.ParseArtifact(i.expandArgs(src, true))
		}
		// If it parses as an artifact, treat as artifact.
		if parseErr == nil {
			srcs[index] = artifactSrc.String()
			allClassical = false
		} else {
			srcs[index] = i.expandArgs(src, false)
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

			expandedFlagArgs := i.expandArgsSlice(srcFlagArgs[index], true)
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
				err = i.converter.CopyArtifact(ctx, src, dest, platform, allowPrivileged, srcBuildArgs, opts.IsDirCopy, opts.KeepTs, opts.KeepOwn, opts.Chown, opts.IfExists, opts.SymlinkNoFollow)
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

		err = i.converter.CopyClassical(ctx, srcs, dest, opts.IsDirCopy, opts.KeepTs, opts.KeepOwn, opts.Chown)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "copy classical")
		}
	}
	return nil
}

func (i *Interpreter) handleSaveArtifact(ctx context.Context, cmd spec.Command) error {
	opts := saveArtifactOpts{}
	args, err := flagutil.ParseArgs("SAVE ARTIFACT", &opts, getArgsCopy(cmd))
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

	saveFrom := i.expandArgs(args[0], false)
	saveTo = i.expandArgs(saveTo, false)
	saveAsLocalTo = i.expandArgs(saveAsLocalTo, false)

	if i.local {
		if saveAsLocalTo != "" {
			return i.errorf(cmd.SourceLocation, "SAVE ARTIFACT AS LOCAL is not implemented under LOCALLY targets")
		}
		err = i.converter.SaveArtifactFromLocal(ctx, saveFrom, saveTo, opts.KeepTs, opts.IfExists, "")
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
		}
		return nil
	}

	err = i.converter.SaveArtifact(ctx, saveFrom, saveTo, saveAsLocalTo, opts.KeepTs, opts.KeepOwn, opts.IfExists, opts.SymlinkNoFollow, i.pushOnlyAllowed)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply SAVE ARTIFACT")
	}
	return nil
}

func (i *Interpreter) handleSaveImage(ctx context.Context, cmd spec.Command) error {
	opts := saveImageOpts{}
	args, err := flagutil.ParseArgs("SAVE IMAGE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid SAVE IMAGE arguments %v", cmd.Args)
	}
	for index, cf := range opts.CacheFrom {
		opts.CacheFrom[index] = i.expandArgs(cf, false)
	}
	if opts.Push && len(args) == 0 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for SAVE IMAGE --push: %v", cmd.Args)
	}

	imageNames := args
	for index, img := range imageNames {
		imageNames[index] = i.expandArgs(img, false)
	}
	if len(imageNames) == 0 && !opts.CacheHint && len(opts.CacheFrom) == 0 {
		fmt.Fprintf(os.Stderr, "Deprecation: using SAVE IMAGE with no arguments is no longer necessary and can be safely removed\n")
		return nil
	}
	err = i.converter.SaveImage(ctx, imageNames, opts.Push, opts.Insecure, opts.CacheHint, opts.CacheFrom)
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
	args, err := flagutil.ParseArgs("BUILD", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid BUILD arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for BUILD: %s", cmd.Args)
	}
	fullTargetName := i.expandArgs(args[0], true)
	platformsSlice := make([]*specs.Platform, 0, len(opts.Platforms))
	for index, p := range opts.Platforms {
		opts.Platforms[index] = i.expandArgs(p, false)
		platform, err := llbutil.ParsePlatform(p)
		if err != nil {
			return i.wrapError(err, cmd.SourceLocation, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	if async && !isSafeAsyncBuildArgs(opts.BuildArgs) {
		return errCannotAsync
	}
	expandedBuildArgs := i.expandArgsSlice(opts.BuildArgs, true)
	expandedFlagArgs := i.expandArgsSlice(args[1:], true)
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}
	expandedBuildArgs = append(parsedFlagArgs, expandedBuildArgs...)
	if len(platformsSlice) == 0 {
		platformsSlice = []*specs.Platform{nil}
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
				errChan := i.converter.BuildAsync(ctx, fullTargetName, platform, allowPrivileged, bas, buildCmd)
				i.monitorErrChan(ctx, errChan)
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
	workdirPath := i.expandArgs(cmd.Args[0], false)
	err := i.converter.Workdir(ctx, workdirPath)
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
	user := i.expandArgs(cmd.Args[0], false)
	err := i.converter.User(ctx, user)
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
	if !withShell {
		for index, arg := range cmdArgs {
			cmdArgs[index] = i.expandArgs(arg, false)
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
	if !withShell {
		for index, arg := range entArgs {
			entArgs[index] = i.expandArgs(arg, false)
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
		ports[index] = i.expandArgs(port, false)
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
		volumes[index] = i.expandArgs(volume, false)
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
	var key, value string
	switch len(cmd.Args) {
	case 3:
		if cmd.Args[1] != "=" {
			return i.errorf(cmd.SourceLocation, "invalid syntax")
		}
		value = i.expandArgs(cmd.Args[2], false)
		fallthrough
	case 1:
		key = cmd.Args[0] // Note: Not expanding args for key.
	default:
		return i.errorf(cmd.SourceLocation, "invalid syntax")
	}
	err := i.converter.Env(ctx, key, value)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply ENV")
	}
	return nil
}

func (i *Interpreter) handleArg(ctx context.Context, cmd spec.Command) error {
	if i.pushOnlyAllowed {
		return i.pushOnlyErr(cmd.SourceLocation)
	}
	var key, value string
	switch len(cmd.Args) {
	case 3:
		if cmd.Args[1] != "=" {
			return i.errorf(cmd.SourceLocation, "invalid syntax")
		}
		value = i.expandArgs(cmd.Args[2], true)
		fallthrough
	case 1:
		key = cmd.Args[0] // Note: Not expanding args for key.
	default:
		return i.errorf(cmd.SourceLocation, "invalid syntax")
	}
	// Args declared in the base target are global.
	global := i.isBase
	err := i.converter.Arg(ctx, key, value, global)
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
				return i.errorf(cmd.SourceLocation, "syntax error")
			}
			nextEqual = false
		} else {
			value := i.expandArgs(arg, false)
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
	err := i.converter.Label(ctx, labels)
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
	args, err := flagutil.ParseArgs("GIT CLONE", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid GIT CLONE arguments %v", cmd.Args)
	}
	if len(args) != 2 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for GIT CLONE: %s", cmd.Args)
	}
	gitURL := i.expandArgs(args[0], false)
	gitCloneDest := i.expandArgs(args[1], false)
	opts.Branch = i.expandArgs(opts.Branch, false)

	convertedGitURL, _, err := i.gitLookup.ConvertCloneURL(gitURL)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "unable to use %v with configured earthly credentials from ~/.earthly/config.yml", cmd.Args)
	}

	err = i.converter.GitClone(ctx, convertedGitURL, opts.Branch, gitCloneDest, opts.KeepTs)
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
	args, err := flagutil.ParseArgs("HEALTHCHECK", &opts, getArgsCopy(cmd))
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
		cmdArgs[index] = i.expandArgs(arg, false)
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
	args, err := flagutil.ParseArgs("WITH DOCKER", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid WITH DOCKER arguments %v", cmd.Args)
	}
	if len(args) != 0 {
		return i.errorf(cmd.SourceLocation, "invalid WITH DOCKER arguments %v", args)
	}
	opts.Platform = i.expandArgs(opts.Platform, false)
	platform, err := llbutil.ParsePlatform(opts.Platform)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse platform %s", opts.Platform)
	}
	for index, cf := range opts.ComposeFiles {
		opts.ComposeFiles[index] = i.expandArgs(cf, false)
	}
	for index, cs := range opts.ComposeServices {
		opts.ComposeServices[index] = i.expandArgs(cs, false)
	}
	for index, load := range opts.Loads {
		opts.Loads[index] = i.expandArgs(load, true)
	}
	expandedBuildArgs := i.expandArgsSlice(opts.BuildArgs, true)
	for index, p := range opts.Pulls {
		opts.Pulls[index] = i.expandArgs(p, false)
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
		expandedFlagArgs := i.expandArgsSlice(flagArgs, true)
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
	args, err := flagutil.ParseArgs("DO", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid DO arguments %v", cmd.Args)
	}
	if len(args) < 1 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for DO: %s", args)
	}

	expandedFlagArgs := i.expandArgsSlice(args[1:], true)
	parsedFlagArgs, err := variables.ParseFlagArgs(expandedFlagArgs)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "parse flag args")
	}

	ucName := i.expandArgs(args[0], false)
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
	args, err := flagutil.ParseArgs("IMPORT", &opts, getArgsCopy(cmd))
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "invalid IMPORT arguments %v", cmd.Args)
	}

	if len(args) != 1 && len(args) != 3 {
		return i.errorf(cmd.SourceLocation, "invalid number of arguments for IMPORT: %s", args)
	}
	if len(args) == 3 && args[1] != "AS" {
		return i.errorf(cmd.SourceLocation, "invalid arguments for IMPORT: %s", args)
	}
	importStr := i.expandArgs(args[0], false)
	var as string
	if len(args) == 3 {
		as = i.expandArgs(args[2], false)
	}
	isGlobal := (i.target.Target == "base")
	err = i.converter.Import(ctx, importStr, as, isGlobal, i.allowPrivileged, opts.AllowPrivileged)
	if err != nil {
		return i.wrapError(err, cmd.SourceLocation, "apply IMPORT")
	}
	return nil
}

// ----------------------------------------------------------------------------

func (i *Interpreter) handleDoUserCommand(ctx context.Context, command domain.Command, relCommand domain.Command, uc spec.UserCommand, do spec.Command, buildArgs []string, allowPrivileged bool) error {
	if allowPrivileged && !i.allowPrivileged {
		return i.errorf(uc.SourceLocation, "invalid privileged in COMMAND") // this shouldn't happen, but check just in case
	}
	if len(uc.Recipe) == 0 || uc.Recipe[0].Command.Name != "COMMAND" {
		return i.errorf(uc.SourceLocation, "command recipes must start with COMMAND")
	}
	if len(uc.Recipe[0].Command.Args) > 0 {
		return i.errorf(uc.Recipe[0].SourceLocation, "COMMAND takes no arguments")
	}
	scopeName := fmt.Sprintf(
		"%s (%s line %d:%d)",
		command.StringCanonical(), do.SourceLocation.File, do.SourceLocation.StartLine, do.SourceLocation.StartColumn)
	err := i.converter.EnterScope(ctx, command, baseTarget(relCommand), allowPrivileged, scopeName, buildArgs)
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

func (i *Interpreter) expandArgsSlice(words []string, keepPlusEscape bool) []string {
	ret := make([]string, 0, len(words))
	for _, word := range words {
		ret = append(ret, i.expandArgs(word, keepPlusEscape))
	}
	return ret
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

func (i *Interpreter) expandArgs(word string, keepPlusEscape bool) string {
	ret := i.converter.ExpandArgs(escapeSlashPlus(word))
	if keepPlusEscape {
		return ret
	}
	return unescapeSlashPlus(ret)
}

func (i *Interpreter) monitorErrChan(ctx context.Context, errChan chan error) {
	go func() {
		select {
		case err := <-errChan:
			if err != nil && !errors.Is(err, context.Canceled) {
				i.parallelErrChan <- err
			}
		case <-ctx.Done():
		}
	}()
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
	splitLoad := strings.SplitN(loadStr, "=", 2)
	if len(splitLoad) < 2 {
		// --load <target-name>
		// (will infer image name from SAVE IMAGE of that target)
		image = ""
		target = loadStr
	} else {
		// --load <image-name>=<target-name>
		image = splitLoad[0]
		target = splitLoad[1]
	}
	if strings.HasPrefix(target, "(") && strings.HasSuffix(target, ")") {
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

// parseParans turns "(+target --flag=something)" into "+target" and []string{"--flag=something"}.
func parseParans(str string) (string, []string, error) {
	if !strings.HasPrefix(str, "(") || !strings.HasSuffix(str, ")") {
		return "", nil, errors.New("parans atom not in ( ... )")
	}
	str = str[1 : len(str)-1] // remove ( and )
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

func isSafeAsyncBuildArgs(args []string) bool {
	for _, arg := range args {
		_, v, _ := variables.ParseKeyValue(arg)
		if strings.HasPrefix(v, "$(") {
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
