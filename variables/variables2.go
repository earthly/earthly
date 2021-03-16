package variables

import (
	"github.com/earthly/earthly/states/dedup"
	"github.com/pkg/errors"
)

// Type is the variable type (string, bool etc).
type Type int

const (
	// StringType is the type representing string.
	StringType Type = iota
	// BoolType is the type representing bool.
	BoolType
)

// Var represents a variable
// (build arg, env var, command line arg, builtin arg, etc).
type Var struct {
	// Type is the type of the variable.
	Type Type
	// Value is the value of the variable.
	Value string
}

// BuildArgInput returns the BuildArgInput for this variable.
func (v Var) BuildArgInput(name string, defaultValue string) dedup.BuildArgInput {
	return dedup.BuildArgInput{
		Name:          name,
		ConstantValue: v.Value,
		DefaultValue:  defaultValue,
	}
}

// ValidateArgType verifies that the value conforms to the type specified in varType.
func ValidateArgType(varType Type, value string) error {
	switch varType {
	case StringType:
		return nil
	case BoolType:
		switch value {
		case "true", "false":
			return nil
		default:
			return errors.Errorf("invalid value %s for arg type bool", value)
		}
	default:
		return errors.Errorf("unsupported arg type %v", varType)
	}
}
