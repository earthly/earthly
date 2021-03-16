package variables

import (
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	dfShell "github.com/moby/buildkit/frontend/dockerfile/shell"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Collection is a collection of variable scopes used within a single target.
type Collection struct {
	// Always inactive scopes. These scopes only influence newly declared
	// args. They do not otherwise participate when args are expanded.
	overriding *Scope
	builtin    *Scope

	// Always active scopes. These scopes influence the value of args directly.
	argsStack []*Scope
	envs      *Scope
	globals   *Scope

	// A scope containing all scopes above, combined.
	effectiveCache *Scope
}

// NewCollection creates a new Collection to be used in the context of a target.
func NewCollection(target domain.Target, platform specs.Platform, gitMeta *gitutil.GitMetadata, overridingVars *Scope) *Collection {
	return &Collection{
		overriding: overridingVars,
		builtin:    BuiltinArgs(target, platform, gitMeta),
		argsStack:  []*Scope{NewScope()},
		envs:       NewScope(),
		globals:    NewScope(),
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

// EnvVars returns a copy of the env vars.
func (c *Collection) EnvVars() *Scope {
	return c.envs.Clone()
}

// Globals returns a copy of the globals.
func (c *Collection) Globals() *Scope {
	return c.globals.Clone()
}

// SetGlobals sets the global variables.
func (c *Collection) SetGlobals(globals *Scope) {
	c.globals = globals
	c.effectiveCache = nil
}

// Overriding returns a copy of the overriding args.
func (c *Collection) Overriding() *Scope {
	return c.overriding.Clone()
}

// SetOverriding sets the overriding args.
func (c *Collection) SetOverriding(overriding *Scope) {
	c.overriding = overriding
	c.effectiveCache = nil
}

// SetPlatform sets the platform, updating the builtin args.
func (c *Collection) SetPlatform(platform specs.Platform) {
	SetPlatformArgs(c.builtin, platform)
	c.effectiveCache = nil
}

// GetActive returns an active variable by name.
func (c *Collection) GetActive(name string) (Var, bool) {
	return c.effective().GetActive(name)
}

// SortedActiveVariables returns the active variable names in a sorted slice.
func (c *Collection) SortedActiveVariables() []string {
	return c.effective().SortedActive()
}

// SortedOverridingVariables returns the overriding variable names in a sorted slice.
func (c *Collection) SortedOverridingVariables() []string {
	return c.overriding.SortedActive()
}

// Expand expands variables within the given word.
func (c *Collection) Expand(word string) string {
	shlex := dfShell.NewLex('\\')
	varMap := c.effective().ActiveValueMap()
	ret, err := shlex.ProcessWordWithMap(word, varMap)
	if err != nil {
		// No effect if there is an error.
		return word
	}
	return ret
}

// DeclareArg declares an arg. The effective value may be
// different than the default, if the variable has been overridden.
func (c *Collection) DeclareArg(name string, varType Type, defaultValue string, global bool, pncvf ProcessNonConstantVariableFunc) (Var, error) {
	ef := c.effective()
	var finalValue Var
	existing, found := ef.GetAny(name)
	if found {
		finalValue = existing
	} else {
		v, err := parseArgValue(name, varType, defaultValue, pncvf)
		if err != nil {
			return Var{}, err
		}
		finalValue = v
	}
	err := ValidateArgType(varType, finalValue.Value)
	if err != nil {
		return Var{}, err
	}
	c.args().AddActive(name, finalValue)
	if global {
		c.globals.AddActive(name, finalValue)
	}
	c.effectiveCache = nil
	return Var{
		Type:  varType,
		Value: finalValue.Value,
	}, nil
}

// DeclareEnv declares an env var.
func (c *Collection) DeclareEnv(name string, value string) {
	c.envs.AddActive(name, Var{
		Value: value,
		Type:  StringType,
	})
	c.effectiveCache = nil
}

func (c *Collection) args() *Scope {
	return c.argsStack[len(c.argsStack)-1]
}

func (c *Collection) pushArgsStack() {
	c.argsStack = append(c.argsStack, NewScope())
	c.effectiveCache = nil
}

func (c *Collection) popArgsStack() {
	if len(c.argsStack) == 0 {
		panic("trying to pop an empty argsStack")
	}
	c.argsStack = c.argsStack[:(len(c.argsStack) - 1)]
	c.effectiveCache = nil
}

// effective returns the variables as a single combined scope.
func (c *Collection) effective() *Scope {
	if c.effectiveCache == nil {
		if len(c.argsStack) == 1 {
			// Not in a UDC.
			c.effectiveCache = CombineScopes(c.overriding, c.builtin, c.args(), c.envs, c.globals)
		} else {
			// Within a UDC.
			c.effectiveCache = CombineScopes(c.builtin, c.args(), c.envs)
		}
	}
	return c.effectiveCache
}
