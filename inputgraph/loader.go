package inputgraph

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/earthly/earthly/ast/command"
	"github.com/earthly/earthly/ast/commandflag"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/buildkitskipper/hasher"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/variables"
	"github.com/pkg/errors"
)

var ErrRemoteNotSupported = fmt.Errorf("remote targets not supported")

type loader struct {
	conslog         conslogging.ConsoleLogger
	target          domain.Target
	visited         map[string]struct{}
	hasher          *hasher.Hasher
	varCollection   *variables.Collection
	features        *features.Features
	isBaseTarget    bool
	ci              bool
	builtinArgs     variables.DefaultArgs
	overridingVars  *variables.Scope
	earthlyCIRunner bool
}

func newLoader(ctx context.Context, opt HashOpt) *loader {
	// Other important values are set by load().
	return &loader{
		conslog:         opt.Console,
		target:          opt.Target,
		visited:         map[string]struct{}{},
		hasher:          hasher.New(),
		isBaseTarget:    opt.Target.Target == "base",
		ci:              opt.CI,
		builtinArgs:     opt.BuiltinArgs,
		overridingVars:  opt.OverridingVars,
		earthlyCIRunner: opt.EarthlyCIRunner,
	}
}

func (l *loader) handleFrom(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.FromOpts{}
	args, err := flagutil.ParseArgsCleaned(command.From, &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return err
	}

	fromTarget := args[0]
	if !strings.Contains(fromTarget, "+") {
		return nil
	}

	return l.loadTargetFromString(ctx, fromTarget, args[1:], false)
}

func (l *loader) handleBuild(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.BuildOpts{}
	args, err := flagutil.ParseArgsCleaned(command.Build, &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing BUILD arg")
	}

	argCombos, err := flagutil.BuildArgMatrix(args)
	if err != nil {
		return errors.Wrap(err, "failed to compute arg matrix")
	}

	for _, args := range argCombos {
		err := l.loadTargetFromString(ctx, args[0], args[1:], opts.PassArgs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) handleCopy(ctx context.Context, cmd spec.Command) error {
	opts := commandflag.CopyOpts{}
	args, err := flagutil.ParseArgsCleaned(command.Copy, &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return err
	}

	if opts.From != "" {
		return errors.New("COPY --from is not supported")
	}

	if len(args) < 2 {
		return errors.New("COPY must include a source and destination")
	}

	srcs := args[:len(args)-1]
	for _, src := range srcs {
		if err := l.handleCopySrc(ctx, src, opts.IsDirCopy); err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) handleCopySrc(ctx context.Context, src string, isDir bool) error {

	if strings.Contains(src, "$") {
		return errors.Errorf("dynamic COPY source %q cannot be resolved", src)
	}

	var (
		classical   bool
		artifactSrc domain.Artifact
		extraArgs   []string
		err         error
	)

	// Complex form with args: (+target --arg=1)
	if flagutil.IsInParamsForm(src) {
		var artifactName string
		classical = false
		artifactName, extraArgs, err = flagutil.ParseParams(src)
		if err != nil {
			return errors.Wrap(err, "failed to parse COPY params")
		}
		artifactSrc, err = domain.ParseArtifact(artifactName)
		if err != nil {
			return errors.Wrap(err, "failed to parse artifact")
		}
	} else { // Simpler form: '+target/artifact' or 'file/path'
		artifactSrc, err = domain.ParseArtifact(src)
		if err != nil {
			classical = true
		}
	}

	// COPY classical (not from another target)
	if classical {
		path := filepath.Join(l.target.GetLocalPath(), src)
		files, err := l.expandCopyFiles(path)
		if err != nil {
			return err
		}
		sort.Strings(files)
		for _, file := range files {
			if err := l.hasher.HashFile(ctx, file); err != nil {
				return errors.Wrapf(err, "failed to hash file %s", path)
			}
		}
		return nil
	}

	// Remote targets aren't supported.
	if artifactSrc.Target.IsRemote() {
		return errors.New("unable to handle remote target")
	}

	targetName := artifactSrc.Target.LocalPath + "+" + artifactSrc.Target.Target
	if err := l.loadTargetFromString(ctx, targetName, extraArgs, false); err != nil {
		return err
	}

	return nil
}

// expandCopyFiles expands a single COPY source into a slice containing all
// nested files. The file names will then be used in our hash.
func (l *loader) expandCopyFiles(src string) ([]string, error) {
	if strings.Contains(src, "**") {
		return nil, errors.New("globstar (**) not supported")
	}

	if strings.Contains(src, "*") {
		matches, err := filepath.Glob(src)
		if err != nil {
			return nil, errors.Wrap(err, "unable to expand glob pattern")
		}
		return l.expandDirs(matches...)
	}

	stat, err := os.Stat(src)
	if err != nil {
		return nil, errors.Wrap(err, "failed to stat file")
	}

	if stat.IsDir() {
		return l.expandDirs(src)
	}

	return []string{src}, nil
}

// expandDirs takes a list of paths (directories and files) and recursively
// expands all directories into a list of nested files. The final list will not
// contain directories.
func (l *loader) expandDirs(dirs ...string) ([]string, error) {
	ret := []string{}
	for _, dir := range dirs {
		stat, err := os.Stat(dir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to stat file")
		}
		if stat.IsDir() {
			entries, err := os.ReadDir(dir)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read dir")
			}
			for _, entry := range entries {
				next := filepath.Join(dir, entry.Name())
				if entry.IsDir() {
					found, err := l.expandDirs(next)
					if err != nil {
						return nil, err
					}
					ret = append(ret, found...)
				} else {
					ret = append(ret, next)
				}
			}
		} else {
			ret = append(ret, dir)
		}
	}
	return uniqStrs(ret), nil
}

func (l *loader) expandArgs(ctx context.Context, args []string) ([]string, error) {
	ret := []string{}
	for _, arg := range args {
		expanded, err := l.varCollection.Expand(arg, func(cmd string) (string, error) {
			return arg, nil // Return the original expression so it can be referenced later.
		})
		if err != nil {
			return nil, err
		}
		ret = append(ret, expanded)
	}
	return ret, nil
}

func (l *loader) handleCommand(ctx context.Context, cmd spec.Command) error {
	// All commands are expanded and hashed at a minimum.
	var err error
	cmd.Args, err = l.expandArgs(ctx, cmd.Args)
	if err != nil {
		return err
	}
	l.hashCommand(cmd)

	// Some commands require more processing.
	switch cmd.Name {
	case command.From:
		return l.handleFrom(ctx, cmd)
	case command.Build:
		return l.handleBuild(ctx, cmd)
	case command.Copy:
		return l.handleCopy(ctx, cmd)
	case command.Arg:
		return l.handleArg(ctx, cmd)
	default:
		return nil
	}
}

func (l *loader) handleArg(ctx context.Context, cmd spec.Command) error {
	opts, key, valueOrNil, err := flagutil.ParseArgArgs(ctx, cmd, l.isBaseTarget, l.features.ExplicitGlobal)
	if err != nil {
		return errors.Wrap(err, "failed to parse ARG args")
	}

	declOpts := []variables.DeclareOpt{
		variables.AsArg(),
	}

	if valueOrNil != nil {
		declOpts = append(declOpts, variables.WithValue(*valueOrNil))
	}

	if opts.Global {
		declOpts = append(declOpts, variables.AsGlobal())
	}

	_, _, err = l.varCollection.DeclareVar(key, declOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to declare variable")
	}

	return nil
}

func (l *loader) handleWith(ctx context.Context, with spec.WithStatement) error {
	if with.Command.Name != command.Docker {
		return errors.New("expected WITH DOCKER")
	}
	err := l.handleWithDocker(ctx, with.Command)
	if err != nil {
		return err
	}
	return l.loadBlock(ctx, with.Body)
}

func (l *loader) handleWithDocker(ctx context.Context, cmd spec.Command) error {
	// Special case since handleWithDocker doesn't get called from handleCommand.
	var err error
	cmd.Args, err = l.expandArgs(ctx, cmd.Args)
	if err != nil {
		return err
	}
	l.hashCommand(cmd)
	opts := commandflag.WithDockerOpts{}
	_, err = flagutil.ParseArgsCleaned("WITH DOCKER", &opts, flagutil.GetArgsCopy(cmd))
	if err != nil {
		return errors.New("failed to parse WITH DOCKER flags")
	}
	for _, load := range opts.Loads {
		_, target, extraArgs, err := earthfile2llb.ParseLoad(load)
		if err != nil {
			return errors.Wrap(err, "failed to parse --load value")
		}
		err = l.loadTargetFromString(ctx, target, extraArgs, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) handleIf(ctx context.Context, ifStmt spec.IfStatement) error {
	l.hashIfStatement(ifStmt)
	if err := l.loadBlock(ctx, ifStmt.IfBody); err != nil {
		return err
	}
	if ifStmt.ElseBody != nil {
		if err := l.loadBlock(ctx, *ifStmt.ElseBody); err != nil {
			return err
		}
	}
	for _, elseIf := range ifStmt.ElseIf {
		l.hashElseIf(elseIf)
		if err := l.loadBlock(ctx, elseIf.Body); err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) handleFor(ctx context.Context, forStmt spec.ForStatement) error {
	l.hashForStatement(forStmt)

	opts := commandflag.NewForOpts()

	args, err := flagutil.ParseArgsCleaned("FOR", &opts, forStmt.Args)
	if err != nil {
		return errors.Wrap(err, "failed to parse FOR args")
	}

	expandedArgs, err := l.expandArgs(ctx, args)
	if err != nil {
		return err
	}

	name := expandedArgs[0]
	vals := flattenForArgs(expandedArgs, opts.Separators)

	l.hasher.HashInt(len(vals))

	for _, val := range vals {
		l.hasher.HashString(fmt.Sprintf("FOR %s=%s", name, val))
		l.varCollection.SetArg(name, val)
		err := l.loadBlock(ctx, forStmt.Body)
		if err != nil {
			return err
		}
		l.varCollection.UnsetArg(name)
	}

	return nil
}

func flattenForArgs(args []string, seps string) []string {
	var ret []string
	for _, arg := range args {
		if strings.ContainsAny(arg, seps) {
			found := strings.FieldsFunc(arg, func(r rune) bool {
				return strings.ContainsRune(seps, r)
			})
			ret = append(ret, found...)
		} else {
			ret = append(ret, arg)
		}
	}
	return ret
}

func (l *loader) handleWait(ctx context.Context, waitStmt spec.WaitStatement) error {
	l.hashWaitStatement(waitStmt)
	return l.handleStatements(ctx, waitStmt.Body)
}

func (l *loader) handleTry(ctx context.Context, tryStmt spec.TryStatement) error {
	l.hashTryStatement(tryStmt)
	if err := l.handleStatements(ctx, tryStmt.TryBody); err != nil {
		return err
	}
	if tryStmt.CatchBody != nil {
		if err := l.handleStatements(ctx, *tryStmt.CatchBody); err != nil {
			return err
		}
	}
	if tryStmt.FinallyBody != nil {
		if err := l.handleStatements(ctx, *tryStmt.FinallyBody); err != nil {
			return err
		}
	}
	return nil
}

func (l *loader) handleStatements(ctx context.Context, stmts []spec.Statement) error {
	l.hasher.HashInt(len(stmts))
	for _, stmt := range stmts {
		if err := l.handleStatement(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
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
	return errors.New("unexpected statement type")
}

func (l *loader) loadBlock(ctx context.Context, b spec.Block) error {
	return l.handleStatements(ctx, b)
}

func (l *loader) forTarget(ctx context.Context, target domain.Target, args []string, passArgs bool) (*loader, error) {
	fullTargetName := target.String()

	visited := copyVisited(l.visited)
	visited[fullTargetName] = struct{}{}

	flagArgs, err := variables.ParseFlagArgs(args)
	if err != nil {
		return nil, err
	}

	overriding, err := variables.ParseCommandLineArgs(flagArgs)
	if err != nil {
		return nil, err
	}

	if passArgs {
		overriding = variables.CombineScopes(overriding, l.overridingVars)
	}

	return &loader{
		conslog:         l.conslog,
		target:          target,
		visited:         visited,
		hasher:          l.hasher,
		isBaseTarget:    target.Target == "base",
		ci:              l.ci,
		builtinArgs:     l.builtinArgs,
		overridingVars:  overriding,
		earthlyCIRunner: l.earthlyCIRunner,
	}, nil
}

func (l *loader) loadTargetFromString(ctx context.Context, targetName string, args []string, passArgs bool) error {
	// If the target name contains a variable that hasn't been expanded, we
	// won't be able to explore the rest of the graph and generate a valid hash.
	if strings.Contains(targetName, "$") {
		return errors.Errorf("dynamic target %q cannot be resolved", targetName)
	}

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
		// Prevent infinite loops; the converter does a better job since it also
		// looks at args and if conditions.
		return errors.Errorf("circular dependency detected; %s already called", fullTargetName)
	}

	targetLoader, err := l.forTarget(ctx, target, args, passArgs)
	if err != nil {
		return errors.Wrapf(err, "failed to create loader for target %q", targetName)
	}

	return targetLoader.load(ctx)
}

func (l *loader) findProject(ctx context.Context) (org, project string, err error) {
	if l.target.IsRemote() {
		return "", "", ErrRemoteNotSupported
	}

	resolver := buildcontext.NewResolver(nil, nil, l.conslog, "", "", "", 0, "")

	buildCtx, err := resolver.Resolve(ctx, nil, nil, l.target)
	if err != nil {
		return "", "", err
	}

	ef := buildCtx.Earthfile
	if ef.Version != nil {
		l.hashVersion(*ef.Version)
	}

	for _, stmt := range ef.BaseRecipe {
		if stmt.Command != nil && stmt.Command.Name == command.Project {
			args := stmt.Command.Args
			if len(args) != 1 {
				return "", "", errors.New("failed to parse PROJECT command")
			}
			parts := strings.Split(args[0], "/")
			if len(parts) != 2 {
				return "", "", errors.New("failed to parse PROJECT command")
			}
			return parts[0], parts[1], nil
		}
	}

	return "", "", errors.New("PROJECT command missing")
}

func (l *loader) load(ctx context.Context) error {
	if l.target.IsRemote() {
		return ErrRemoteNotSupported
	}

	resolver := buildcontext.NewResolver(nil, nil, l.conslog, "", "", "", 0, "")

	buildCtx, err := resolver.Resolve(ctx, nil, nil, l.target)
	if err != nil {
		return err
	}

	l.features = buildCtx.Features

	collOpt := variables.NewCollectionOpt{
		Console:         l.conslog,
		Target:          l.target,
		CI:              l.ci,
		BuiltinArgs:     l.builtinArgs,
		OverridingVars:  l.overridingVars,
		EarthlyCIRunner: l.earthlyCIRunner,
		GitMeta:         buildCtx.GitMetadata,
		Features:        l.features,
	}
	l.varCollection = variables.NewCollection(collOpt)

	ef := buildCtx.Earthfile
	if ef.Version != nil {
		l.hashVersion(*ef.Version)
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
