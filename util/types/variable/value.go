package variable

import (
	"fmt"

	"github.com/earthly/earthly/domain"
)

type ValueType int

const (
	TypeUnknown ValueType = iota
	TypeArg
	TypePath
	TypeString
)

func (vt ValueType) String() string {
	switch vt {
	case TypeUnknown:
		return "unknown"
	case TypeArg:
		return "arg"
	case TypePath:
		return "path"
	case TypeString:
		return "string"
	default:
		return "corrupt"
	}
}

type Value struct {
	Str      string
	ComeFrom domain.Target
	Type     ValueType
}

func (v *Value) String() string {
	//fmt.Printf("type is %d\n", v.Type)
	switch v.Type {
	case TypeUnknown, TypeString:
		return v.Str
	case TypeArg:
		return fmt.Sprintf("TYPE_ARG: %s\n", v.Str)
	case TypePath:
		return fmt.Sprintf("TYPE_PATH: %s\n", v.Str)
	default:
		panic(fmt.Sprintf("Value corrupt; unknown type %d", v.Type))
	}
}
