package variables

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/shell"
	"github.com/pkg/errors"

	dfShell "github.com/moby/buildkit/frontend/dockerfile/shell"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	ErrRedeclared   = errors.New("ARG was declared twice in the same target")
	ErrArgNotFound  = errors.New("no matching ARG found in this scope")
	ErrInvalidScope = errors.New("this action is not allowed in this scope")

	ShellOutEnvs = map[string]struct{}{
		"HOME": {},
		"PATH": {},
	}
)

type stackFrame struct {
	frameName string
	// absRef is the ref any other ref in this frame would be relative to.
	absRef  domain.Reference
	imports *domain.ImportTracker

	// Always inactive scopes. These scopes only influence newly declared
	// args. They do not otherwise participate when args are expanded.
	overriding *Scope

	// Always active scopes. These scopes influence the value of args directly.
	args    *Scope
	globals *Scope

	// Explicitly defined scopes. These are declared using commands to alter
	// existing values and will always be active, and even override the
	// overriding scopes.
	assignedArgs    *Scope
	assignedGlobals *Scope
}

// Collection is a collection of variable scopes used within a single target.
type Collection struct {
	// These scopes are always present, regardless of the stack position.
	builtin *Scope // inactive
	envs    *Scope // active

	errorOnRedeclare bool
	shelloutAnywhere bool

	project string
	org     string

	stack []*stackFrame

	// A scope containing all scopes above, combined.
	effectiveCache *Scope

	console conslogging.ConsoleLogger
}

// NewCollectionOpt contains supported arguments which
// the `NewCollection` function may accept.
type NewCollectionOpt struct {
	Console          conslogging.ConsoleLogger
	Target           domain.Target
	Push             bool
	CI               bool
	PlatformResolver *platutil.Resolver
	NativePlatform   specs.Platform
	GitMeta          *gitutil.GitMetadata
	BuiltinArgs      DefaultArgs
	OverridingVars   *Scope
	AssignedVars     *Scope
	Features         *features.Features
	GlobalImports    map[string]domain.ImportTrackerVal
}

// NewCollection creates a new Collection to be used in the context of a target.
func NewCollection(opts NewCollectionOpt) *Collection {
	target := opts.Target
	console := opts.Console
	return &Collection{
		builtin:          BuiltinArgs(target, opts.PlatformResolver, opts.GitMeta, opts.BuiltinArgs, opts.Features, opts.Push, opts.CI),
		envs:             NewScope(),
		errorOnRedeclare: opts.Features.ArgScopeSet,
		shelloutAnywhere: opts.Features.ShellOutAnywhere,
		stack: []*stackFrame{{
			frameName:       target.StringCanonical(),
			absRef:          target,
			imports:         domain.NewImportTracker(console, opts.GlobalImports),
			overriding:      opts.OverridingVars,
			args:            NewScope(),
			globals:         NewScope(),
			assignedArgs:    NewScope(),
			assignedGlobals: NewScope(),
		}},
		console: console,
	}
}

// ResetEnvVars resets the collection's env vars.
func (c *Collection) ResetEnvVars(envs *Scope) {
	if envs == nil {
		envs = NewScope()
	}
	c.envs = envs
	c.effectiveCache = nil
}

// SetOrg sets the organization name.
func (c *Collection) SetOrg(org string) {
	c.org = org
}

// Org returns the organization name.
func (c *Collection) Org() string {
	return c.org
}

// SetProject sets the project name.
func (c *Collection) SetProject(project string) {
	c.project = project
}

// Project returns the project name.
func (c *Collection) Project() string {
	return c.project
}

// EnvVars returns a copy of the env vars.
func (c *Collection) EnvVars() *Scope {
	return c.envs.Clone()
}

// Globals returns a copy of the globals.
func (c *Collection) Globals() *Scope {
	return c.globals().Clone()
}

// SetGlobals sets the global variables.
func (c *Collection) SetGlobals(globals *Scope) {
	c.frame().globals = globals
	c.effectiveCache = nil
}

// AssignedGlobals returns a copy of the assigned globals (i.e. globals that
// have been reassigned using SET).
func (c *Collection) AssignedGlobals() *Scope {
	return c.assignedGlobals().Clone()
}

// SetAssignedGlobals sets the assigned global variables.
func (c *Collection) SetAssignedGlobals(assigned *Scope) {
	c.frame().assignedGlobals = assigned
	c.effectiveCache = nil
}

// Overriding returns a copy of the overriding args.
func (c *Collection) Overriding() *Scope {
	return c.overriding().Clone()
}

// SetOverriding sets the overriding args.
func (c *Collection) SetOverriding(overriding *Scope) {
	c.frame().overriding = overriding
	c.effectiveCache = nil
}

// SetPlatform sets the platform, updating the builtin args.
func (c *Collection) SetPlatform(platr *platutil.Resolver) {
	SetPlatformArgs(c.builtin, platr)
	c.effectiveCache = nil
}

// SetLocally sets the locally flag, updating the builtin args.
func (c *Collection) SetLocally(locally bool) {
	SetLocally(c.builtin, locally)
	c.effectiveCache = nil
}

// Get returns a variable by name.
func (c *Collection) Get(name string, opts ...ScopeOpt) (string, bool) {
	return c.effective().Get(name, opts...)
}

// SortedVariables returns the current variable names in a sorted slice.
func (c *Collection) SortedVariables(opts ...ScopeOpt) []string {
	return c.effective().Sorted(opts...)
}

// SortedOverridingVariables returns the overriding variable names in a sorted slice.
func (c *Collection) SortedOverridingVariables() []string {
	return c.overriding().Sorted()
}

// ExpandOld expands variables within the given word, it does not perform shelling-out.
// it will eventually be removed when the ShellOutAnywhere feature is fully-adopted
func (c *Collection) ExpandOld(word string) string {
	shlex := dfShell.NewLex('\\')
	varMap := c.effective().Map(WithActive())
	ret, err := shlex.ProcessWordWithMap(word, varMap)
	if err != nil {
		// No effect if there is an error.
		return word
	}
	return ret
}

// Expand expands variables within the given word.
func (c *Collection) Expand(word string, shellOut shell.EvalShellOutFn) (string, error) {
	shlex := shell.NewLex('\\')
	shlex.ShellOut = shellOut
	varMap := c.effective().Map(WithActive())
	return shlex.ProcessWordWithMap(word, varMap, ShellOutEnvs)
}

func (c *Collection) overridingOrDefault(name string, defaultValue string, pncvf ProcessNonConstantVariableFunc) (string, error) {
	if v, ok := c.overriding().Get(name); ok {
		return v, nil
	}
	return parseArgValue(name, defaultValue, pncvf)
}

func (c *Collection) declareOldArg(name string, defaultValue string, global bool, pncvf ProcessNonConstantVariableFunc) (string, string, error) {
	ef := c.effective()
	finalDefaultValue := defaultValue
	var finalValue string
	existing, found := ef.Get(name)
	if found {
		finalValue = existing
	} else {
		v, err := parseArgValue(name, defaultValue, pncvf)
		if err != nil {
			return "", "", err
		}
		finalValue = v
		finalDefaultValue = v
	}
	opts := []ScopeOpt{WithActive()}
	c.args().Add(name, finalValue, opts...)
	if global {
		c.globals().Add(name, finalValue, opts...)
	}
	c.effectiveCache = nil
	return finalValue, finalDefaultValue, nil
}

// DeclareArg declares an arg. The effective value may be
// different than the default, if the variable has been overridden.
func (c *Collection) DeclareArg(name string, defaultValue string, global bool, pncvf ProcessNonConstantVariableFunc) (string, string, error) {
	if !c.errorOnRedeclare {
		return c.declareOldArg(name, defaultValue, global, pncvf)
	}
	if !c.shelloutAnywhere {
		return "", "", errors.New("the --arg-scope-and-set feature flag requires --shell-out-anywhere")
	}

	v, err := c.overridingOrDefault(name, defaultValue, pncvf)
	if err != nil {
		return "", "", err
	}

	c.effectiveCache = nil
	opts := []ScopeOpt{WithActive(), NoOverride()}
	if global {
		if _, ok := c.args().Get(name); ok {
			return "", "", errors.Wrapf(ErrRedeclared, "could not override non-global ARG '%[1]v' with global ARG [hint: '%[1]v' was already declared as a non-global ARG in this scope - did you mean to add '--global' to the original declaration?]", name)
		}
		ok := c.globals().Add(name, v, opts...)
		if !ok {
			return "", "", errors.Wrapf(ErrRedeclared, "global ARG '%v' redeclared [hint: use SET to reassign an existing ARG]", name)
		}
		return v, v, nil
	}
	ok := c.args().Add(name, v, opts...)
	if !ok {
		return "", "", errors.Wrapf(ErrRedeclared, "ARG '%v' redeclared [hint: use SET to reassign an existing ARG]", name)
	}
	return v, defaultValue, nil
}

// SetArg sets the value of an arg.
func (c *Collection) SetArg(name string, value string) {
	c.args().Add(name, value, WithActive())
	c.effectiveCache = nil
}

// UnsetArg removes an arg if it exists.
func (c *Collection) UnsetArg(name string) {
	c.args().Remove(name)
	c.effectiveCache = nil
}

// DeclareEnv declares an env var.
func (c *Collection) DeclareEnv(name string, value string) {
	c.envs.Add(name, value, WithActive())
	c.effectiveCache = nil
}

// UpdateArg updates the value of an existing ARG. It will override the value of
// the arg, regardless of where the value was previously defined.
//
// It returns ErrArgNotFound if the variable was not found.
func (c *Collection) UpdateArg(name, value string, pncvf ProcessNonConstantVariableFunc, isBase bool) (retErr error) {
	defer func() {
		if retErr == nil {
			c.effectiveCache = nil
		}
	}()
	if _, ok := c.effective().Get(name, WithActive()); !ok {
		return errors.Wrapf(ErrArgNotFound, "could not SET undeclared variable '%[1]v' [hint: '%[1]v' needs to be declared with 'ARG %[1]v' before it can be used with SET]", name)
	}
	v, err := parseArgValue(name, value, pncvf)
	if err != nil {
		return errors.Wrap(err, "failed to parse SET arg value")
	}
	if _, ok := c.args().Get(name, WithActive()); ok {
		c.assignedArgs().Add(name, v, WithActive())
		return nil
	}
	if _, ok := c.globals().Get(name, WithActive()); ok {
		if !isBase {
			return errors.Wrapf(ErrInvalidScope, "could not SET global variable '%[1]v' outside the base target [hint: you can declare a new non-global ARG to be used with SET using 'ARG myArg=$%[1]v']", name)
		}
		c.assignedGlobals().Add(name, v, WithActive())
		return nil
	}
	return errors.New("variable %v was found in the effective cache, but not in args or globals - this should not happen, please report this incident in github")
}

// Imports returns the imports tracker of the current frame.
func (c *Collection) Imports() *domain.ImportTracker {
	return c.frame().imports
}

// EnterFrame creates a new stack frame.
func (c *Collection) EnterFrame(frameName string, absRef domain.Reference, assignedGlobals *Scope, overriding *Scope, globals *Scope, globalImports map[string]domain.ImportTrackerVal) {
	c.stack = append(c.stack, &stackFrame{
		frameName:       frameName,
		absRef:          absRef,
		imports:         domain.NewImportTracker(c.console, globalImports),
		overriding:      overriding,
		assignedGlobals: assignedGlobals,
		globals:         globals,
		assignedArgs:    NewScope(),
		args:            NewScope(),
	})
	c.effectiveCache = nil
}

// ExitFrame exits the latest stack frame.
func (c *Collection) ExitFrame() {
	if len(c.stack) == 0 {
		panic("trying to pop an empty argsStack")
	}
	c.stack = c.stack[:(len(c.stack) - 1)]
	c.effectiveCache = nil
}

// AbsRef returns a ref that any other reference should be relative to as part of the stack frame.
func (c *Collection) AbsRef() domain.Reference {
	return c.frame().absRef
}

// IsStackAtBase returns whether the stack has size 1.
func (c *Collection) IsStackAtBase() bool {
	return len(c.stack) == 1
}

// StackString returns the stack as a string.
func (c *Collection) StackString() string {
	builder := make([]string, 0, len(c.stack))
	for i := len(c.stack) - 1; i >= 0; i-- {
		activeNames := c.stack[i].args.Sorted(WithActive())
		row := make([]string, 0, len(activeNames)+1)
		row = append(row, c.stack[i].frameName)
		for _, k := range activeNames {
			v, _ := c.stack[i].overriding.Get(k)
			row = append(row, fmt.Sprintf("--%s=%s", k, v))
		}
		builder = append(builder, strings.Join(row, " "))
	}
	return strings.Join(builder, "\ncalled from\t")
}

func (c *Collection) frame() *stackFrame {
	return c.stack[len(c.stack)-1]
}

func (c *Collection) args() *Scope {
	return c.frame().args
}

func (c *Collection) globals() *Scope {
	return c.frame().globals
}

func (c *Collection) assignedArgs() *Scope {
	return c.frame().assignedArgs
}

func (c *Collection) assignedGlobals() *Scope {
	return c.frame().assignedGlobals
}

func (c *Collection) overriding() *Scope {
	return c.frame().overriding
}

// effective returns the variables as a single combined scope.
func (c *Collection) effective() *Scope {
	if c.effectiveCache == nil {
		ag := c.assignedGlobals().Clone()
		for k := range c.args().Map() {
			// c.assignedGlobals() override c.overriding(), but not c.args() -
			// so this effectively creates a difference (in set terms) between
			// c.assignedGlobals() and c.args() to avoid overriding c.args().
			ag.Remove(k)
		}
		c.effectiveCache = CombineScopes(c.assignedArgs(), ag, c.overriding(), c.builtin, c.args(), c.envs, c.globals())
	}
	return c.effectiveCache
}
