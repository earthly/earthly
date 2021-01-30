package variables

import (
	"github.com/earthly/earthly/states/dedup"
)

// Variable is an object representing a build arg or an env var.
type Variable struct {
	isEnvVar bool
	value    string
}

// NewConstant creates a new constant build arg.
func NewConstant(value string) Variable {
	return Variable{
		value: value,
	}
}

// NewConstantEnvVar cretes a new constant env var.
func NewConstantEnvVar(value string) Variable {
	return Variable{
		isEnvVar: true,
		value:    value,
	}
}

// IsEnvVar returns whether the variable is an env var.
func (v Variable) IsEnvVar() bool {
	return v.isEnvVar
}

// ConstantValue returns the value of the constant build arg.
func (v Variable) ConstantValue() string {
	return v.value
}

// BuildArgInput returns the BuildArgInput for this variable.
func (v Variable) BuildArgInput(name string, defaultValue string) dedup.BuildArgInput {
	return dedup.BuildArgInput{
		Name:          name,
		ConstantValue: v.value,
		DefaultValue:  defaultValue,
	}
}
