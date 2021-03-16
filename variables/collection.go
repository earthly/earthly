package variables

import (
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	dfShell "github.com/moby/buildkit/frontend/dockerfile/shell"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Collection2 is a collection of variable scopes used within a single target.
type Collection2 struct {
	// Always inactive scopes.
	overriding *Scope
	builtin    *Scope

	// Always active scopes.
	argsStack []*Scope
	envs      *Scope
	globals   *Scope

	// A scope containing all scopes above, combined.
	effectiveCache *Scope
}

// NewCollection2 creates a new Collection2 to be used in the context of a target.
func NewCollection2(target domain.Target, platform specs.Platform, gitMeta *gitutil.GitMetadata, overridingVars *Scope) *Collection2 {
	return &Collection2{
		overriding: overridingVars,
		builtin:    BuiltinArgs(target, platform, gitMeta),
		argsStack:  []*Scope{NewScope()},
		envs:       NewScope(),
		globals:    NewScope(),
	}
}

// ResetEnvVars resets the collection's env vars.
func (c *Collection2) ResetEnvVars(envs *Scope) {
	if envs == nil {
		envs = NewScope()
	}
	c.envs = envs
	c.effectiveCache = nil
}

// EnvVars returns a copy of the env vars.
func (c *Collection2) EnvVars() *Scope {
	return c.envs.Clone()
}

// Globals returns a copy of the globals.
func (c *Collection2) Globals() *Scope {
	return c.globals.Clone()
}

// SetGlobals sets the global variables.
func (c *Collection2) SetGlobals(globals *Scope) {
	c.globals = globals
	c.effectiveCache = nil
}

// Overriding returns a copy of the overriding args.
func (c *Collection2) Overriding() *Scope {
	return c.overriding.Clone()
}

// SetOverriding sets the overriding args.
func (c *Collection2) SetOverriding(overriding *Scope) {
	c.overriding = overriding
	c.effectiveCache = nil
}

// SetPlatform sets the platform, updating the builtin args.
func (c *Collection2) SetPlatform(platform specs.Platform) {
	SetPlatformArgs(c.builtin, platform)
	c.effectiveCache = nil
}

// GetActive returns an active variable by name.
func (c *Collection2) GetActive(name string) (Var, bool) {
	return c.effective().GetActive(name)
}

// SortedActiveVariables returns the active variable names in a sorted slice.
func (c *Collection2) SortedActiveVariables() []string {
	return c.effective().SortedActive()
}

// SortedOverridingVariables returns the overriding variable names in a sorted slice.
func (c *Collection2) SortedOverridingVariables() []string {
	return c.overriding.SortedActive()
}

// Expand expands variables within the given word.
func (c *Collection2) Expand(word string) string {
	ef := c.effective()
	shlex := dfShell.NewLex('\\')
	varMap := make(map[string]string)
	for varName := range ef.AllActive() {
		variable, _ := ef.GetActive(varName)
		varMap[varName] = variable.Value
	}
	ret, err := shlex.ProcessWordWithMap(word, varMap)
	if err != nil {
		// No effect if there is an error.
		return word
	}
	return ret
}

// DeclareArg declares an arg. The effective value may be
// different than the default, if the variable has been overridden.
func (c *Collection2) DeclareArg(name string, varType Type, defaultValue string, global bool, pncvf ProcessNonConstantVariableFunc) (Var, error) {
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
func (c *Collection2) DeclareEnv(name string, value string) {
	c.envs.AddActive(name, Var{
		Value: value,
		Type:  StringType,
	})
}

func (c *Collection2) args() *Scope {
	return c.argsStack[len(c.argsStack)-1]
}

func (c *Collection2) pushArgsStack() {
	c.argsStack = append(c.argsStack, NewScope())
	c.effectiveCache = nil
}

func (c *Collection2) popArgsStack() {
	if len(c.argsStack) == 0 {
		panic("trying to pop an empty argsStack")
	}
	c.argsStack = c.argsStack[:(len(c.argsStack) - 1)]
	c.effectiveCache = nil
}

// effective returns the variables as a single combined scope.
func (c *Collection2) effective() *Scope {
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
