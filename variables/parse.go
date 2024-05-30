package variables

import (
	"os"
	"strings"

	"github.com/earthly/earthly/variables/reserved"

	"github.com/pkg/errors"
)

// ProcessNonConstantVariableFunc is a function which takes in an expression and
// turns it into a state, target intput and arg index.
type ProcessNonConstantVariableFunc func(name string, expression string) (value string, argIndex int, err error)

// ParseCommandLineArgs parses a slice of old build args
// (the ones passed via --build-arg) and returns a new scope.
func ParseCommandLineArgs(args []string) (*Scope, error) {
	ret := NewScope()
	for _, arg := range args {
		splitArg := strings.SplitN(arg, "=", 2)
		if len(splitArg) < 1 {
			return nil, errors.Errorf("invalid build arg %s", splitArg)
		}
		key := splitArg[0]
		value := ""
		hasValue := false
		if len(splitArg) == 2 {
			value = splitArg[1]
			hasValue = true
		}
		if reserved.IsBuiltIn(key) {
			return nil, errors.Errorf("built-in arg %s cannot be passed on the command line", key)
		}
		if !hasValue {
			var found bool
			value, found = os.LookupEnv(key)
			if !found {
				return nil, errors.Errorf("env var %s not set", key)
			}
		}
		ret.Add(key, NewStringVariable(value))
	}
	return ret, nil
}

// ParseArgs parses args passed as --build-arg to an Earthly command, such as BUILD or FROM.
func ParseArgs(args []string, pncvf ProcessNonConstantVariableFunc, current *Collection) (*Scope, error) {
	ret := NewScope()
	for _, arg := range args {
		name, variable, err := parseArg(arg, pncvf, current)
		if err != nil {
			return nil, errors.Wrapf(err, "parse build arg %s", arg)
		}
		ret.Add(name, NewStringVariable(variable))
	}
	return ret, nil
}

func parseArg(arg string, pncvf ProcessNonConstantVariableFunc, current *Collection) (string, string, error) {
	var name string
	splitArg := strings.SplitN(arg, "=", 2)
	if len(splitArg) < 1 {
		return "", "", errors.Errorf("invalid build arg %s", splitArg)
	}
	name = splitArg[0]
	value := ""
	hasValue := false
	if len(splitArg) == 2 {
		value = splitArg[1]
		hasValue = true
	}
	if hasValue {
		if reserved.IsBuiltIn(name) {
			return "", "", errors.Errorf("value cannot be specified for built-in build arg %s", name)
		}
		v, err := parseArgValue(name, value, pncvf)
		if err != nil {
			return "", "", err
		}
		return name, v, nil
	}
	v, ok := current.Get(name, WithActive())
	if !ok {
		return "", "", errors.Errorf("value not specified for build arg %s and no value can be inferred", name)
	}
	return name, v, nil
}

func parseArgValue(name string, value string, pncvf ProcessNonConstantVariableFunc) (string, error) {
	if pncvf == nil {
		return value, nil
	}
	if strings.HasPrefix(value, "$(") {
		// Variable build arg - resolve value.
		var err error
		value, _, err = pncvf(name, value)
		if err != nil {
			return "", err
		}
	}
	return value, nil
}

// ParseEnvVars parses env vars from a slice of strings of the form "key=value".
func ParseEnvVars(envVars []string) *Scope {
	ret := NewScope()
	for _, envVar := range envVars {
		k, v, _ := ParseKeyValue(envVar)
		ret.Add(k, NewStringVariable(v), WithActive())
	}
	return ret
}
