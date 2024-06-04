package variables

import (
	"fmt"
	"sort"
	"strings"

	"github.com/earthly/earthly/util/types/variable"
	"github.com/earthly/earthly/variables/reserved"
)

func NewStringVariable(s string) variable.Value {
	return variable.Value{Str: s}
}

// Scope represents a variable scope.
type Scope struct {
	variables map[string]variable.Value
	// activeVariables are variables that are active right now as we have passed the point of
	// their declaration.
	activeVariables map[string]bool
}

// NewScope creates a new variable scope.
func NewScope() *Scope {
	return &Scope{
		variables:       make(map[string]variable.Value),
		activeVariables: make(map[string]bool),
	}
}

// DebugString returns a string that can be printed while debugging
func (s *Scope) DebugString() string {
	var sb strings.Builder
	for _, k := range s.SortedNames() {
		v := s.variables[k]
		sb.WriteString(fmt.Sprintf("%s=%s", k, v))
		if s.activeVariables[k] {
			sb.WriteString(" (active)")
		} else {
			sb.WriteString(" (inactive)")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Clone returns a copy of the scope.
func (s *Scope) Clone() *Scope {
	ret := NewScope()
	for k, v := range s.variables {
		ret.variables[k] = v
	}
	for k := range s.activeVariables {
		ret.activeVariables[k] = true
	}
	return ret
}

// Get gets a variable by name.
func (s *Scope) Get(name string, opts ...ScopeOpt) (variable.Value, bool) {
	opt := applyOpts(opts...)
	v, ok := s.variables[name]
	if !ok {
		return variable.Value{}, false
	}
	if opt.active && !s.activeVariables[name] {
		return variable.Value{}, false
	}
	return v, true
}

// Add sets a variable to a value within this scope. It returns true if the
// value was set.
func (s *Scope) Add(name string, value variable.Value, opts ...ScopeOpt) bool {
	opt := applyOpts(opts...)
	_, existed := s.variables[name]
	if opt.noOverride && existed {
		return false
	}
	s.variables[name] = value
	if opt.active {
		s.activeVariables[name] = true
	}
	return true
}

// Remove removes a variable from the scope.
func (s *Scope) Remove(name string) {
	delete(s.variables, name)
	delete(s.activeVariables, name)
}

// Map returns a name->value variable map of variables in this scope.
func (s *Scope) Map(opts ...ScopeOpt) map[string]variable.Value {
	opt := applyOpts(opts...)
	m := map[string]variable.Value{}
	for k, v := range s.variables {
		if opt.active && !s.activeVariables[k] {
			continue
		}
		m[k] = v
	}
	return m
}

func (s *Scope) MapWithStringValues(opts ...ScopeOpt) map[string]string {
	m := map[string]string{}
	for k, v := range s.Map(opts...) {
		m[k] = v.String()
	}
	return m
}

// Keys returns a sorted list of variable names in this Scope.
func (s *Scope) SortedNames(opts ...ScopeOpt) []string {
	opt := applyOpts(opts...)
	var sorted []string
	for k := range s.variables {
		if opt.active && !s.activeVariables[k] {
			continue
		}
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	return sorted
}

// BuildArgs returns s as a slice of build args, as they would have been passed
// in originally at the CLI or in a BUILD command.
func (s *Scope) BuildArgs(opts ...ScopeOpt) []variable.KeyValue {
	vars := s.SortedNames(opts...)
	var args []variable.KeyValue
	for _, v := range vars {
		val, _ := s.Get(v)
		args = append(args, variable.KeyValue{
			Key:   v,
			Value: &val,
		})
	}
	return args
}

// RemoveReservedArgsFromScope returns a new scope ommits any builtin arguments
func RemoveReservedArgsFromScope(scope *Scope) *Scope {
	s := NewScope()
	for k, v := range scope.variables {
		if reserved.IsBuiltIn(k) {
			continue
		}
		s.variables[k] = v
		if scope.activeVariables[k] {
			s.activeVariables[k] = true
		}
	}
	return s
}

// CombineScopes combines all the variables across all scopes, with the
// following precedence:
//
// 1. Active variables
// 2. Inactive variables
// 3. All other things equal, left-most scopes have precedence
func CombineScopes(scopes ...*Scope) *Scope {
	s := NewScope()
	precedence := [][]ScopeOpt{
		{WithActive()},
		nil,
	}
	for _, opts := range precedence {
		addOpts := append(opts, NoOverride())
		for _, scope := range scopes {
			for k, v := range scope.Map(opts...) {
				s.Add(k, v, addOpts...)
			}
		}
	}
	return s
}

// CombineScopesInactive combines all scopes (leaving all variables inactive), with left-most scope having precedence,
func CombineScopesInactive(scopes ...*Scope) *Scope {
	s := NewScope()
	for _, scope := range scopes {
		for k, v := range scope.Map() {
			s.Add(k, v, NoOverride())
		}
	}
	return s
}

type scopeOpts struct {
	active     bool
	noOverride bool
}

func applyOpts(opts ...ScopeOpt) scopeOpts {
	var opt scopeOpts
	for _, o := range opts {
		opt = o(opt)
	}
	return opt
}

// ScopeOpt is an option function for setting flags when adding to or getting
// from a Scope.
type ScopeOpt func(scopeOpts) scopeOpts

// WithActive is a ScopeOpt for active variables. When passed to Add, it sets
// the variable to active; when passed to Get or Map, it causes them to only
// return active variables.
func WithActive() ScopeOpt {
	return func(o scopeOpts) scopeOpts {
		o.active = true
		return o
	}
}

// NoOverride only applies to Add. When passed to Add, NoOverride will prevent
// Add from overriding an existing value.
//
// This will also prevent Add from applying other opts to the existing variable,
// so if the caller wishes to set options on the existing value, they should
// update the value on a false return from Add.
func NoOverride() ScopeOpt {
	return func(o scopeOpts) scopeOpts {
		o.noOverride = true
		return o
	}
}
