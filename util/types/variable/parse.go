package variable

import (
	"strings"

	"github.com/earthly/earthly/domain"
)

func ParseKeyValue(arg string, comefrom domain.Target) KeyValue {
	var name string
	splitArg := strings.SplitN(arg, "=", 2)
	if len(splitArg) < 1 {
		panic("bad")
	}
	name = splitArg[0]
	var value *Value
	if len(splitArg) == 2 {
		value = &Value{
			Str:      splitArg[1],
			ComeFrom: comefrom,
		}
	}
	return KeyValue{
		Key:   name,
		Value: value,
	}
}
