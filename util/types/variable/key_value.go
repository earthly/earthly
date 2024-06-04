package variable

import (
	"fmt"
	"strings"
)

type KeyValue struct {
	Key   string
	Value *Value
}

type KeyValueSlice []KeyValue

func (kvs KeyValueSlice) DebugString() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, x := range kvs {
		if i > 0 {
			sb.WriteString(", ")
		}
		if x.Value == nil {
			sb.WriteString(fmt.Sprintf("%s: <undefined>", x.Key))
		} else {
			sb.WriteString(fmt.Sprintf("%s: %s (from %s)", x.Key, x.Value.Str, x.Value.ComeFrom))
		}
	}
	sb.WriteString("]")
	return sb.String()
}
