package variables

import (
	"os"
	"strings"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/types/variable"
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

// ParseArgs2 parses args passed as --build-arg to an Earthly command, such as BUILD or FROM.
func ParseArgs2(args []variable.KeyValue, pncvf ProcessNonConstantVariableFunc, current *Collection, currentTarget domain.Reference) (*Scope, error) {
	ret := NewScope()
	for _, arg := range args {
		arg, err := parseArg2(arg, pncvf, current, currentTarget)
		if err != nil {
			return nil, errors.Wrapf(err, "parse build arg %s", arg)
		}
		ret.Add(arg.Key, *arg.Value) // parseArg2 always sets the Value
	}
	return ret, nil
}

func parseArg2(arg variable.KeyValue, pncvf ProcessNonConstantVariableFunc, current *Collection, currentTarget domain.Reference) (variable.KeyValue, error) {
	//fmt.Printf("in parseArg2 with %+v\n", arg)
	var name string
	name = arg.Key
	if arg.Value != nil {
		if reserved.IsBuiltIn(name) {
			return variable.KeyValue{}, errors.Errorf("value cannot be specified for built-in build arg %s", name)
		}
		//if !strings.Contains(arg.Value.Str, "$") {
		//	// keep existing value
		//	return arg, nil
		//}
		//val := arg.Value.String(currentTarget)
		//expandedValue, err := parseArgValue(name, val, pncvf)
		//if err != nil {
		//	return variable.KeyValue{}, err
		//}
		//return name, v, nil
	}
	//v, ok := current.GetValue(name, WithActive())
	//if !ok {
	//	return variable.KeyValue{}, errors.Errorf("value not specified for build arg %s and no value can be inferred", name)
	//}
	//arg.Value = &v
	return arg, nil
}

// ParseArgs parses args passed as --build-arg to an Earthly command, such as BUILD or FROM.
//func ParseArgs(args []string, pncvf ProcessNonConstantVariableFunc, current *Collection) (*Scope, error) {
//	ret := NewScope()
//	for _, arg := range args {
//		name, variable, err := parseArg(arg, pncvf, current)
//		if err != nil {
//			return nil, errors.Wrapf(err, "parse build arg %s", arg)
//		}
//		ret.Add(name, NewStringVariable(variable))
//	}
//	return ret, nil
//}

//func parseArg(arg string, pncvf ProcessNonConstantVariableFunc, current *Collection) (string, string, error) {
//	var name string
//	splitArg := strings.SplitN(arg, "=", 2)
//	if len(splitArg) < 1 {
//		return "", "", errors.Errorf("invalid build arg %s", splitArg)
//	}
//	name = splitArg[0]
//	value := ""
//	hasValue := false
//	if len(splitArg) == 2 {
//		value = splitArg[1]
//		hasValue = true
//	}
//	if hasValue {
//		if reserved.IsBuiltIn(name) {
//			return "", "", errors.Errorf("value cannot be specified for built-in build arg %s", name)
//		}
//		v, err := parseArgValue(name, value, pncvf)
//		if err != nil {
//			return "", "", err
//		}
//		return name, v, nil
//	}
//	v, ok := current.GetValue(name, WithActive())
//	if !ok {
//		return "", "", errors.Errorf("value not specified for build arg %s and no value can be inferred", name)
//	}
//	return name, v.String(), nil
//}

func parseArgValue2(name string, value variable.Value, pncvf ProcessNonConstantVariableFunc, currentTarget domain.Reference) (variable.Value, error) {
	if pncvf == nil {
		return value, nil
	}
	if strings.HasPrefix(value.Str, "$(") {
		// Variable build arg - resolve value.
		var err error
		value.Str, _, err = pncvf(name, value.String(currentTarget))
		if err != nil {
			return variable.Value{}, err
		}
		value.Type = variable.TypeUnknown // unsure what to do here
		value.ComeFrom = domain.Target{}  // clear it out? TODO: should we set something else saying it must be a string in this case? or can it still work.... maybe it can.
	}
	return value, nil
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
