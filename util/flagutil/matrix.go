package flagutil

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/util/types/variable"
	"github.com/pkg/errors"
)

type argGroup2 struct {
	key    string
	values []*variable.Value
}

// BuildArgMatrix2 builds a 2-dimensional slice of arguments that contains all
// combinations
func BuildArgMatrix2(args []variable.KeyValue) ([][]variable.KeyValue, error) {
	groupedArgs := make([]argGroup2, 0, len(args))
	for _, arg := range args {
		//k, v, err := parseKeyValue(arg)
		//if err != nil {
		//	return nil, err
		//}

		found := false
		for i, g := range groupedArgs {
			if g.key == arg.Key {
				groupedArgs[i].values = append(groupedArgs[i].values, arg.Value)
				found = true
				break
			}
		}
		if !found {
			groupedArgs = append(groupedArgs, argGroup{
				key:    k,
				values: []*string{v},
			})
		}
	}
	return crossProduct(groupedArgs, nil), nil
}

type argGroup struct {
	key    string
	values []*string
}

// BuildArgMatrix builds a 2-dimensional slice of arguments that contains all
// combinations
func BuildArgMatrix(args []string) ([][]string, error) {
	groupedArgs := make([]argGroup, 0, len(args))
	for _, arg := range args {
		k, v, err := parseKeyValue(arg)
		if err != nil {
			return nil, err
		}
		fmt.Printf("parsed %s into %s -> %s\n", arg, k, v)

		found := false
		for i, g := range groupedArgs {
			if g.key == k {
				groupedArgs[i].values = append(groupedArgs[i].values, v)
				found = true
				break
			}
		}
		if !found {
			groupedArgs = append(groupedArgs, argGroup{
				key:    k,
				values: []*string{v},
			})
		}
	}
	return crossProduct(groupedArgs, nil), nil
}

func crossProduct(ga []argGroup, prefix []string) [][]string {
	if len(ga) == 0 {
		return [][]string{prefix}
	}
	var ret [][]string
	for _, v := range ga[0].values {
		newPrefix := prefix[:]
		var kv string
		if v == nil {
			kv = ga[0].key
		} else {
			kv = fmt.Sprintf("%s=%s", ga[0].key, *v)
		}
		newPrefix = append(newPrefix, kv)

		cp := crossProduct(ga[1:], newPrefix)
		ret = append(ret, cp...)
	}
	return ret
}

func parseKeyValue(arg string) (string, *string, error) {
	var name string
	splitArg := strings.SplitN(arg, "=", 2)
	if len(splitArg) < 1 {
		return "", nil, errors.Errorf("invalid build arg %s", splitArg)
	}
	name = splitArg[0]
	var value *string
	if len(splitArg) == 2 {
		value = &splitArg[1]
	}
	return name, value, nil
}
