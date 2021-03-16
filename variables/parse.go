package variables

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ProcessNonConstantVariableFunc is a function which takes in an expression and
// turns it into a state, target intput and arg index.
type ProcessNonConstantVariableFunc func(name string, expression string) (value string, argIndex int, err error)

// ParseCommandLineArgs parses a slice of constant build args and returns a new scope.
func ParseCommandLineArgs(args []string, dotEnvMap map[string]string) (*Scope, error) {
	ret := NewScope()
	for k, v := range dotEnvMap {
		ret.AddInactive(k, Var{
			Value: v,
			Type:  StringType,
		})
	}
	for _, arg := range args {
		splitArg := strings.SplitN(arg, "=", 2)
		if len(splitArg) < 1 {
			return nil, fmt.Errorf("invalid build arg %s", splitArg)
		}
		key := splitArg[0]
		value := ""
		hasValue := false
		if len(splitArg) == 2 {
			value = splitArg[1]
			hasValue = true
		}
		if !hasValue {
			var found bool
			value, found = os.LookupEnv(key)
			if !found {
				return nil, fmt.Errorf("env var %s not set", key)
			}
		}
		ret.AddInactive(key, Var{
			Value: value,
			Type:  StringType,
		})
	}
	return ret, nil
}

// ParseArgs parses args passed to an Earthly command, such as BUILD or FROM.
func ParseArgs(args []string, pncvf ProcessNonConstantVariableFunc, current *Collection2) (*Scope, error) {
	ret := NewScope()
	for _, arg := range args {
		name, variable, err := parseArg(arg, pncvf, current)
		if err != nil {
			return nil, errors.Wrapf(err, "parse build arg %s", arg)
		}
		ret.AddInactive(name, variable)
	}
	return ret, nil
}

func parseArg(arg string, pncvf ProcessNonConstantVariableFunc, current *Collection2) (string, Var, error) {
	var name string
	splitArg := strings.SplitN(arg, "=", 2)
	if len(splitArg) < 1 {
		return "", Var{}, fmt.Errorf("invalid build arg %s", splitArg)
	}
	name = splitArg[0]
	value := ""
	hasValue := false
	if len(splitArg) == 2 {
		value = splitArg[1]
		hasValue = true
	}
	if hasValue {
		v, err := parseArgValue(name, StringType, value, pncvf)
		if err != nil {
			return "", Var{}, err
		}
		return name, v, nil
	}
	v, ok := current.GetActive(name)
	if !ok {
		return "", Var{}, errors.Errorf("value not specified for build arg %s and no value can be inferred", name)
	}
	return name, v, nil
}

func parseArgValue(name string, varType Type, value string, pncvf ProcessNonConstantVariableFunc) (Var, error) {
	if strings.HasPrefix(value, "$(") {
		// Variable build arg - resolve value.
		var err error
		value, _, err = pncvf(name, value)
		if err != nil {
			return Var{}, err
		}
	}
	err := ValidateArgType(varType, value)
	if err != nil {
		return Var{}, err
	}

	return Var{
		Value: value,
		Type:  varType,
	}, nil
}

// ParseEnvVars parses env vars from a slice of strings of the form "key=value".
func ParseEnvVars(envVars []string) *Scope {
	ret := NewScope()
	for _, envVar := range envVars {
		k, v, _ := ParseKeyValue(envVar)
		ret.AddActive(k, Var{Type: StringType, Value: v})
	}
	return ret
}
