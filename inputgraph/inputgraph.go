package inputgraph

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/ast/command"
	"github.com/earthly/earthly/ast/commandflag"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/buildkitskipper/hasher"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/variables"
	"github.com/pkg/errors"
)

var (
	ErrRemoteNotSupported    = fmt.Errorf("remote targets not supported")
	ErrUnableToDetermineHash = fmt.Errorf("unable to determine hash")
)

func argsContainsStr(args []string, substr string) bool {
	for _, s := range args {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func requiresCrossProduct(args []string) bool {
	seen := map[string]struct{}{}
	for _, s := range args {
		k := strings.SplitN(s, "=", 2)[0]
		if _, found := seen[k]; found {
			return true
		}
		seen[k] = struct{}{}
	}
	return false
}

func getArgsCopy(cmd spec.Command) []string {
	argsCopy := make([]string, len(cmd.Args))
	copy(argsCopy, cmd.Args)
	return argsCopy
}

func parseArgs(cmdName string, opts interface{}, args []string) ([]string, error) {
	processed := stringutil.ProcessParamsAndQuotes(args)
	return flagutil.ParseArgs(cmdName, opts, processed)
}

func (l *loader) handleFrom(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.FromOpts{}
	args, err := parseArgs(command.From, &opts, getArgsCopy(cmd))
	if err != nil {
		return err
	}
	if argsContainsStr(args, "$") {
		return errors.Wrap(ErrUnableToDetermineHash, "unable to handle arg in FROM")
	}
	fromTarget := args[0]
	if !strings.Contains(fromTarget, "+") {
		return nil
	}
	return l.loadTargetFromString(ctx, fromTarget)
}

func (l *loader) handleBuild(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.BuildOpts{}
	args, err := parseArgs(command.Build, &opts, getArgsCopy(cmd))
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return errors.Wrap(ErrUnableToDetermineHash, "missing BUILD arg")
	}
	targetName := args[0]
	if strings.Contains(targetName, "$") {
		return errors.Wrap(ErrUnableToDetermineHash, "unable to handle arg in BUILD")
	}
	if requiresCrossProduct(args) {
		return errors.Wrap(ErrUnableToDetermineHash, "unable to cross-product in BUILD")
	}
	return l.loadTargetFromString(ctx, targetName)
}

func (l *loader) handleCopy(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.CopyOpts{}
	args, err := parseArgs(command.Copy, &opts, getArgsCopy(cmd))
	if err != nil {
		return err
	}
	if argsContainsStr(args, "$") || argsContainsStr(args, "*") || len(args) < 2 || opts.From != "" { // TODO handle globbing (e.g. *), ignored for now
		return errors.Wrap(ErrUnableToDetermineHash, "unable to handle COPY with arg or wildcard")
	}
	srcs := args[:len(args)-1]
	for _, src := range srcs {
		artifactSrc, parseErr := domain.ParseArtifact(src)
		if parseErr != nil {
			// COPY classical
			path := filepath.Join(l.target.GetLocalPath(), src)
			err := l.hasher.HashFile(ctx, path)
			if err != nil {
				return errors.Wrapf(ErrUnableToDetermineHash, "failed to hash file %s: %s", path, err)
			}
			continue
		}
		// COPY from a different target
		if artifactSrc.Target.IsRemote() {
			return errors.Wrap(ErrUnableToDetermineHash, "unable to handle remote target")
		}
		targetName := artifactSrc.Target.LocalPath + "+" + artifactSrc.Target.Target

		err := l.loadTargetFromString(ctx, targetName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) handlePipeline(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.PipelineOpts{}
	_, err := parseArgs(command.Copy, &opts, getArgsCopy(cmd))
	if err != nil {
		return err
	}
	l.isPipeline = !opts.NoPipelineCache
	return nil
}

func (l *loader) handleCommand(ctx context.Context, cmd spec.Command) error {
	l.hasher.HashCommand(cmd)
	switch cmd.Name {
	case command.From:
		return l.handleFrom(ctx, cmd)
	case command.Build:
		return l.handleBuild(ctx, cmd)
	case command.Copy:
		return l.handleCopy(ctx, cmd)
	case command.Pipeline:
		return l.handlePipeline(ctx, cmd)
	default:
		return nil
	}
}

func (l *loader) handleWith(ctx context.Context, with spec.WithStatement) error {
	if with.Command.Name != command.Docker {
		return errors.Wrap(ErrUnableToDetermineHash, "expected WITH DOCKER")
	}
	err := l.handleWithDocker(ctx, with.Command)
	if err != nil {
		return err
	}
	return l.loadBlock(ctx, with.Body)
}

func (l *loader) handleWithDocker(ctx context.Context, cmd spec.Command) error {
	l.hasher.HashCommand(cmd) // special case since handleWithDocker doesn't get called from handleCommand
	opts := commandflag.WithDockerOpts{}
	_, err := parseArgs("WITH DOCKER", &opts, getArgsCopy(cmd))
	if err != nil {
		return errors.Wrap(ErrUnableToDetermineHash, "failed to parse WITH DOCKER flags")
	}
	for _, load := range opts.Loads {
		if strings.Contains(load, "$") {
			return errors.Wrap(ErrUnableToDetermineHash, "unable to handle arg in WITH DOCKER --load")
		}
		_, v, _ := variables.ParseKeyValue(load)
		if v == "" {
			return errors.Wrap(ErrUnableToDetermineHash, "unable to handle WITH DOCKER --load with implicit image name (hint: specify the image name rather than relying on the target's SAVE IMAGE command)")
		}
		err := l.loadTargetFromString(ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) handleIf(ctx context.Context, ifStmt spec.IfStatement) error {
	return errors.Wrap(ErrUnableToDetermineHash, "if not supported")
}

func (l *loader) handleFor(ctx context.Context, forStmt spec.ForStatement) error {
	return errors.Wrap(ErrUnableToDetermineHash, "for not supported")
}

func (l *loader) handleWait(ctx context.Context, waitStmt spec.WaitStatement) error {
	return errors.Wrap(ErrUnableToDetermineHash, "wait not supported")
}

func (l *loader) handleTry(ctx context.Context, tryStmt spec.TryStatement) error {
	return errors.Wrap(ErrUnableToDetermineHash, "try not supported")
}

func (l *loader) handleStatement(ctx context.Context, stmt spec.Statement) error {
	if stmt.Command != nil {
		return l.handleCommand(ctx, *stmt.Command)
	}
	if stmt.With != nil {
		return l.handleWith(ctx, *stmt.With)
	}
	if stmt.If != nil {
		return l.handleIf(ctx, *stmt.If)
	}
	if stmt.For != nil {
		return l.handleFor(ctx, *stmt.For)
	}
	if stmt.Wait != nil {
		return l.handleWait(ctx, *stmt.Wait)
	}
	if stmt.Try != nil {
		return l.handleTry(ctx, *stmt.Try)
	}
	return errors.Wrap(ErrUnableToDetermineHash, "unexpected statement type")
}

func (l *loader) loadBlock(ctx context.Context, b spec.Block) error {
	for _, stmt := range b {
		err := l.handleStatement(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyVisited(m map[string]struct{}) map[string]struct{} {
	m2 := map[string]struct{}{}
	for k := range m {
		m2[k] = struct{}{}
	}
	return m2
}

func (l *loader) loadTargetFromString(ctx context.Context, targetName string) error {
	relTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target name %s", targetName)
	}
	targetRef, err := domain.JoinReferences(l.target, relTarget)
	if err != nil {
		return errors.Wrapf(err, "failed to join %s and %s", l.target, relTarget)
	}
	target := targetRef.(domain.Target)
	fullTargetName := target.String()
	if fullTargetName == "" {
		return fmt.Errorf("missing target string")
	}
	if _, exists := l.visited[fullTargetName]; exists {
		// prevent infinite loops; the converter does a better job since it also looks at args and if conditions
		return errors.Wrapf(ErrUnableToDetermineHash, "circular dependency detected; %s already called", fullTargetName)
	}
	visited := copyVisited(l.visited)
	visited[fullTargetName] = struct{}{}
	loaderInst := &loader{
		conslog: l.conslog,
		target:  target,
		visited: visited,
		hasher:  l.hasher,
	}
	return loaderInst.load(ctx)
}

func (l *loader) findProject(ctx context.Context) (org, project string, err error) {
	if l.target.IsRemote() {
		return "", "", ErrRemoteNotSupported
	}
	resolver := buildcontext.NewResolver(nil, nil, l.conslog, "", "", "", 0, "")
	bc, err := resolver.Resolve(ctx, nil, nil, l.target)
	if err != nil {
		return "", "", err
	}
	ef := bc.Earthfile

	if ef.Version != nil {
		l.hasher.HashVersion(*ef.Version)
	}

	for _, stmt := range ef.BaseRecipe {
		if stmt.Command != nil && stmt.Command.Name == command.Project {
			args := stmt.Command.Args
			if len(args) != 1 {
				return "", "", errors.Wrapf(ErrUnableToDetermineHash, "failed to parse PROJECT command")
			}
			parts := strings.Split(args[0], "/")
			if len(parts) != 2 {
				return "", "", errors.Wrapf(ErrUnableToDetermineHash, "failed to parse PROJECT command")
			}
			return parts[0], parts[1], nil
		}
	}
	return "", "", errors.Wrapf(ErrUnableToDetermineHash, "PROJECT command missing")
}

func (l *loader) load(ctx context.Context) error {
	if l.target.IsRemote() {
		return ErrRemoteNotSupported
	}
	resolver := buildcontext.NewResolver(nil, nil, l.conslog, "", "", "", 0, "")
	bc, err := resolver.Resolve(ctx, nil, nil, l.target)
	if err != nil {
		return err
	}
	ef := bc.Earthfile

	if ef.Version != nil {
		l.hasher.HashVersion(*ef.Version)
	}

	if l.target.Target == "base" {
		return l.loadBlock(ctx, ef.BaseRecipe)
	}
	for _, t := range ef.Targets {
		if t.Name == l.target.Target {
			return l.loadBlock(ctx, t.Recipe)
		}
	}
	return fmt.Errorf("target %s not found", l.target.Target)
}

type loader struct {
	conslog    conslogging.ConsoleLogger
	target     domain.Target
	visited    map[string]struct{}
	hasher     *hasher.Hasher
	isPipeline bool
}

func HashTarget(ctx context.Context, target domain.Target, conslog conslogging.ConsoleLogger) (org, project string, hash []byte, err error) {
	loaderInst := &loader{
		conslog: conslog,
		target:  target,
		visited: map[string]struct{}{},
		hasher:  hasher.New(),
	}
	org, project, err = loaderInst.findProject(ctx)
	if err != nil {
		return "", "", nil, err
	}

	err = loaderInst.load(ctx)
	if err != nil {
		return "", "", nil, err
	}
	if !loaderInst.isPipeline {
		return "", "", nil, errors.Wrap(ErrUnableToDetermineHash, "target is not a pipeline")
	}

	return org, project, loaderInst.hasher.GetHash(), nil
}
