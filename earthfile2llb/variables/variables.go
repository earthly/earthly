package variables

import (
	"github.com/earthly/earthly/earthfile2llb/dedup"
	"github.com/moby/buildkit/client/llb"
)

// Variable is an object representing a build arg or an env var.
type Variable struct {
	isConstant        bool
	isEnvVar          bool
	value             string
	state             llb.State
	variableFromInput dedup.VariableFromInput
}

// NewConstant creates a new constant build arg.
func NewConstant(value string) Variable {
	return Variable{
		isConstant: true,
		value:      value,
	}
}

// NewConstantEnvVar cretes a new constant env var.
func NewConstantEnvVar(value string) Variable {
	return Variable{
		isConstant: true,
		isEnvVar:   true,
		value:      value,
	}
}

// NewVariable creates a new variable build arg.
func NewVariable(state llb.State, targetInput dedup.TargetInput, argIndex int) Variable {
	return Variable{
		isConstant: false,
		state:      state,
		variableFromInput: dedup.VariableFromInput{
			TargetInput: targetInput,
			Index:       argIndex,
		},
	}
}

// IsConstant returns whether this build arg is constant.
func (v Variable) IsConstant() bool {
	return v.isConstant
}

// IsEnvVar returns whether the variable is an env var.
func (v Variable) IsEnvVar() bool {
	return v.isEnvVar
}

// ConstantValue returns the value of the constant build arg.
func (v Variable) ConstantValue() string {
	return v.value
}

// VariableState returns the state that holds a file containing the expression
// result of the build arg.
func (v Variable) VariableState() llb.State {
	return v.state
}

// BuildArgInput returns the BuildArgInput for this variable.
func (v Variable) BuildArgInput(name string, defaultValue string) dedup.BuildArgInput {
	if v.isConstant {
		return dedup.BuildArgInput{
			Name:          name,
			IsConstant:    true,
			ConstantValue: v.value,
			DefaultValue:  defaultValue,
		}
	}
	return dedup.BuildArgInput{
		Name:              name,
		IsConstant:        false,
		VariableFromInput: v.variableFromInput,
		DefaultValue:      defaultValue,
	}
}
