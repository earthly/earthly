package variables

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/shell"
	"github.com/earthly/earthly/util/types/variable"
	"github.com/pkg/errors"

	dfShell "github.com/moby/buildkit/frontend/dockerfile/shell"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	ErrRedeclared   = errors.New("this variable was declared twice in the same target")
	ErrVarNotFound  = errors.New("no matching variable found in this scope")
	ErrInvalidScope = errors.New("this action is not allowed in this scope")
	ErrSetArg       = errors.New("ARG values cannot be reassigned")

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

	// Explicitly defined scopes. These are declared as non-argument variables
	// and will always be active, and even override the overriding scopes.
	vars *Scope
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
	EarthlyCIRunner  bool
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
	if opts.OverridingVars == nil {
		opts.OverridingVars = NewScope()
	}
	return &Collection{
		builtin:          BuiltinArgs(target, opts.PlatformResolver, opts.GitMeta, opts.BuiltinArgs, opts.Features, opts.Push, opts.CI, opts.EarthlyCIRunner),
		envs:             NewScope(),
		errorOnRedeclare: opts.Features.ArgScopeSet,
		shelloutAnywhere: opts.Features.ShellOutAnywhere,
		stack: []*stackFrame{{
			frameName:  target.StringCanonical(),
			absRef:     target,
			imports:    domain.NewImportTracker(console, opts.GlobalImports),
			overriding: opts.OverridingVars,
			args:       NewScope(),
			globals:    NewScope(),
			vars:       NewScope(),
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

// Args returns a copy of the args.
func (c *Collection) Args() *Scope {
	return c.args().Clone()
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

// TopOverriding returns a copy of the top-level overriding args, for use in
// commands that may need to re-parse the base target but have drastically
// different variable scopes.
func (c *Collection) TopOverriding() *Scope {
	if len(c.stack) == 0 {
		return NewScope()
	}
	return c.stack[0].overriding.Clone()
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

// Get returns a variable by name. // TODO rename to GetValueAsString
//func (c *Collection) Get(name string, opts ...ScopeOpt) (string, bool) {
//	v, ok := c.effective().Get(name, opts...)
//	return v.String(), ok
//}

// Get returns a variable by name. // TODO rename back to Get
func (c *Collection) GetValue(name string, opts ...ScopeOpt) (variable.Value, bool) {
	v, ok := c.effective().Get(name, opts...)
	return v, ok
}

// SortedVariables returns the current variable names in a sorted slice.
func (c *Collection) SortedVariables(opts ...ScopeOpt) []string {
	return c.effective().SortedNames(opts...)
}

// SortedOverridingVariables returns the overriding variable names in a sorted slice.
func (c *Collection) SortedOverridingVariables() []string {
	return c.overriding().SortedNames()
}

// ExpandOld expands variables within the given word, it does not perform shelling-out.
// it will eventually be removed when the ShellOutAnywhere feature is fully-adopted
func (c *Collection) ExpandOld(word string) string {
	shlex := dfShell.NewLex('\\')
	varMap := c.effective().MapWithStringValues(c.AbsRef(), WithActive())
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
	varMap := c.effective().MapWithStringValues(c.AbsRef(), WithActive())
	return shlex.ProcessWordWithMap(word, varMap, ShellOutEnvs)
}

func (c *Collection) overridingOrDefault(name string, defaultValue variable.Value, pncvf ProcessNonConstantVariableFunc) (variable.Value, error) {
	if v, ok := c.overriding().Get(name); ok {
		v.Type = defaultValue.Type
		return v, nil
	}
	if v, ok := c.builtin.Get(name); ok {
		v.Type = defaultValue.Type
		return v, nil
	}
	return parseArgValue2(name, defaultValue, pncvf, c.AbsRef())
}

func (c *Collection) declareOldArg(name string, defaultValue variable.Value, global bool, pncvf ProcessNonConstantVariableFunc) (variable.Value, variable.Value, error) {
	ef := c.effective()
	finalDefaultValue := defaultValue
	var finalValue variable.Value
	existing, found := ef.Get(name)
	if found {
		finalValue = existing
	} else {
		v, err := parseArgValue2(name, defaultValue, pncvf, c.AbsRef())
		if err != nil {
			return variable.Value{}, variable.Value{}, err
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

type declarePrefs struct {
	val    variable.Value
	global bool
	arg    bool
	pncvf  ProcessNonConstantVariableFunc
}

// DeclareOpt is an option function for declaring variables.
type DeclareOpt func(declarePrefs) declarePrefs

// WithValue is an option function for setting a variable's value. For ARGs,
// this is only the default value - it can be overridden when calling a target
// at the CLI.
func WithValue(val variable.Value) DeclareOpt {
	return func(o declarePrefs) declarePrefs {
		o.val = val
		return o
	}
}

// AsGlobal is an option function to declare a global variable.
func AsGlobal() DeclareOpt {
	return func(o declarePrefs) declarePrefs {
		o.global = true
		return o
	}
}

// AsArg is an option function to declare an argument.
func AsArg() DeclareOpt {
	return func(o declarePrefs) declarePrefs {
		o.arg = true
		return o
	}
}

// WithPNCVFunc is an option function to apply a ProcessNonConstantVariableFunc
// to ARGs. This supports deprecated functionality and is never used in
// Earthfiles with `VERSION 0.7` and higher.
func WithPNCVFunc(f ProcessNonConstantVariableFunc) DeclareOpt {
	return func(o declarePrefs) declarePrefs {
		o.pncvf = f
		return o
	}
}

// DeclareVar declares a variable. The effective value may be
// different than the default, if the variable has been overridden.
func (c *Collection) DeclareVar(name string, opts ...DeclareOpt) (variable.Value, variable.Value, error) {
	var prefs declarePrefs
	for _, o := range opts {
		prefs = o(prefs)
	}
	if !c.errorOnRedeclare {
		if !prefs.arg {
			return variable.Value{}, variable.Value{}, errors.New("LET requires the --arg-scope-and-set feature")
		}
		//fmt.Printf("adding old %s -> %+v\n", name, prefs.val)
		return c.declareOldArg(name, prefs.val, prefs.global, prefs.pncvf)
	}
	if !c.shelloutAnywhere {
		return variable.Value{}, variable.Value{}, errors.New("the --arg-scope-and-set feature flag requires --shell-out-anywhere")
	}

	c.effectiveCache = nil
	scope := []ScopeOpt{WithActive(), NoOverride()}

	if !prefs.arg {
		//fmt.Printf("adding %s -> %+v\n", name, prefs.val)
		ok := c.vars().Add(name, prefs.val, scope...)
		if !ok {
			return variable.Value{}, variable.Value{}, hint.Wrapf(ErrRedeclared, "if you want to change the value of '%[1]v', use 'SET %[1]v = %[2]q'", name, prefs.val)
		}
		return prefs.val, prefs.val, nil
	}

	if _, ok := c.vars().Get(name, WithActive()); ok {
		return variable.Value{}, variable.Value{}, hint.Wrapf(ErrRedeclared, "'%v' was already declared with LET and cannot be redeclared as an ARG", name)
	}

	//fmt.Printf("overridingOrDefault %s -> %+v\n", name, prefs.val)
	v, err := c.overridingOrDefault(name, prefs.val, prefs.pncvf)
	if err != nil {
		return variable.Value{}, variable.Value{}, err
	}
	v.Type = prefs.val.Type

	if prefs.global {
		if _, ok := c.args().Get(name); ok {
			baseErr := errors.Wrap(ErrRedeclared, "could not override non-global ARG with global ARG")
			return variable.Value{}, variable.Value{}, hint.Wrapf(baseErr, "'%[1]v' was already declared as a non-global ARG in this scope - did you mean to add '--global' to the original declaration?", name)
		}
		ok := c.globals().Add(name, v, scope...)
		if !ok {
			return variable.Value{}, variable.Value{}, hint.Wrapf(ErrRedeclared, "if you want to change the value of '%[1]v', redeclare it as a non-argument variable with 'LET %[1]v = %[2]q'", name, prefs.val)
		}
		return v, v, nil
	}
	//fmt.Printf("adding new %s -> %+v\n", name, v)
	ok := c.args().Add(name, v, scope...)
	if !ok {
		return variable.Value{}, variable.Value{}, hint.Wrapf(ErrRedeclared, "if you want to change the value of '%[1]v', redeclare it as a non-argument variable with 'LET %[1]v = %[2]q'", name, prefs.val)
	}
	return v, prefs.val, nil
}

// SetArg sets the value of an arg.
func (c *Collection) SetArg(name string, value variable.Value) {
	c.args().Add(name, value, WithActive())
	c.effectiveCache = nil
}

// UnsetArg removes an arg if it exists.
func (c *Collection) UnsetArg(name string) {
	c.args().Remove(name)
	c.effectiveCache = nil
}

// DeclareEnv declares an env var.
func (c *Collection) DeclareEnv(name string, value variable.Value) {
	c.envs.Add(name, value, WithActive())
	c.effectiveCache = nil
}

// UpdateVar updates the value of an existing variable. It will override the
// value of the variable, regardless of where the value was previously defined.
//
// It returns ErrVarNotFound if the variable was not found.
func (c *Collection) UpdateVar(name string, value variable.Value, pncvf ProcessNonConstantVariableFunc) (retErr error) {
	defer func() {
		if retErr == nil {
			c.effectiveCache = nil
		}
	}()
	if _, ok := c.effective().Get(name, WithActive()); !ok {
		return hint.Wrapf(ErrVarNotFound, "'%[1]v' needs to be declared with 'LET %[1]v = someValue' before it can be used with SET", name)
	}
	if _, ok := c.vars().Get(name, WithActive()); !ok {
		return hint.Wrapf(ErrSetArg, "'%[1]v' is an ARG and cannot be used with SET - try declaring 'LET %[1]v = $%[1]v' first", name)
	}
	v, err := parseArgValue2(name, value, pncvf, c.AbsRef())
	if err != nil {
		return errors.Wrap(err, "failed to parse SET value")
	}
	c.vars().Add(name, v, WithActive())
	return nil
}

// Imports returns the imports tracker of the current frame.
func (c *Collection) Imports() *domain.ImportTracker {
	return c.frame().imports
}

// EnterFrame creates a new stack frame.
func (c *Collection) EnterFrame(frameName string, absRef domain.Reference, overriding *Scope, globals *Scope, globalImports map[string]domain.ImportTrackerVal) {
	c.stack = append(c.stack, &stackFrame{
		frameName:  frameName,
		absRef:     absRef,
		imports:    domain.NewImportTracker(c.console, globalImports),
		overriding: overriding,
		globals:    globals,
		vars:       NewScope(),
		args:       NewScope(),
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
		activeNames := c.stack[i].args.SortedNames(WithActive())
		row := make([]string, 0, len(activeNames)+1)
		row = append(row, c.stack[i].frameName)
		for _, k := range activeNames {
			v, _ := c.stack[i].overriding.Get(k)
			fmt.Printf("key: %s\n", k)
			fmt.Printf("======= %+v\n", v)
			fmt.Printf("======= %s\n", v)
			fmt.Printf("======= %s\n", v.String(c.AbsRef()))
			row = append(row, fmt.Sprintf("--%s=%s", k, v.String(c.AbsRef())))
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

func (c *Collection) vars() *Scope {
	return c.frame().vars
}

func (c *Collection) overriding() *Scope {
	return c.frame().overriding
}

// effective returns the variables as a single combined scope.
func (c *Collection) effective() *Scope {
	if c.effectiveCache == nil {
		c.effectiveCache = CombineScopes(c.vars(), c.overriding(), c.builtin, c.args(), c.envs, c.globals())
	}
	return c.effectiveCache
}
