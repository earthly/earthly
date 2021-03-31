package variables

import "sort"

// Scope represents a variable scope.
type Scope struct {
	variables map[string]string
	// activeVariables are variables that are active right now as we have passed the point of
	// their declaration.
	activeVariables map[string]bool
}

// NewScope creates a new variable scope.
func NewScope() *Scope {
	return &Scope{
		variables:       make(map[string]string),
		activeVariables: make(map[string]bool),
	}
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

// GetAny returns a variable by name, even if it is not active.
func (s *Scope) GetAny(name string) (string, bool) {
	variable, found := s.variables[name]
	return variable, found
}

// GetActive returns an active variable by name.
func (s *Scope) GetActive(name string) (string, bool) {
	variable, found := s.variables[name]
	active := false
	if found {
		active = s.activeVariables[name]
	}
	if !active {
		variable = ""
	}
	return variable, active
}

// AddInactive adds an inactive variable in the collection.
func (s *Scope) AddInactive(name string, variable string) {
	s.variables[name] = variable
}

// AddActive adds and activates a variable in the collection.
func (s *Scope) AddActive(name string, variable string) {
	s.activeVariables[name] = true
	s.variables[name] = variable
}

// ActiveValueMap returns a map of the values of the active variables.
func (s *Scope) ActiveValueMap() map[string]string {
	ret := make(map[string]string)
	for name := range s.activeVariables {
		ret[name] = s.variables[name]
	}
	return ret
}

// AllValueMap returns a map of the values of all the variables.
func (s *Scope) AllValueMap() map[string]string {
	ret := make(map[string]string)
	for name, value := range s.variables {
		ret[name] = value
	}
	return ret
}

// SortedActive returns the active variable names in a sorted slice.
func (s *Scope) SortedActive() []string {
	varNames := make([]string, 0, len(s.activeVariables))
	for varName := range s.activeVariables {
		varNames = append(varNames, varName)
	}
	sort.Strings(varNames)
	return varNames
}

// SortedAny returns the variable names in a sorted slice.
func (s *Scope) SortedAny() []string {
	varNames := make([]string, 0, len(s.variables))
	for varName := range s.variables {
		varNames = append(varNames, varName)
	}
	sort.Strings(varNames)
	return varNames
}

// CombineScopes combines all the variables across all scopes, with
// left precedence.
func CombineScopes(scopes ...*Scope) *Scope {
	allActive := make(map[string]bool)
	for _, scope := range scopes {
		for name := range scope.activeVariables {
			allActive[name] = true
		}
	}
	allInactive := make(map[string]bool)
	for _, scope := range scopes {
		for name := range scope.variables {
			if !allActive[name] {
				allInactive[name] = true
			}
		}
	}

	s := NewScope()
AllActiveLoop:
	for name := range allActive {
		for _, scope := range scopes {
			variable, active := scope.GetActive(name)
			if active {
				s.AddActive(name, variable)
				continue AllActiveLoop
			}
		}
	}
AllInactiveLoop:
	for name := range allInactive {
		for _, scope := range scopes {
			variable, found := scope.GetAny(name)
			if found {
				s.AddInactive(name, variable)
				continue AllInactiveLoop
			}
		}
	}
	return s
}
