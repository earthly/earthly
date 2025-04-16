package variable

import (
	"fmt"
	"path/filepath"
	"strings"

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

func (v *Value) String(currentTarget domain.Reference) string {
	switch v.Type {
	case TypeUnknown, TypeString:
		return v.Str
	case TypeArg:
		rel := getRelative(currentTarget, v.ComeFrom)
		if rel == "" {
			return v.Str
		}
		if strings.HasPrefix(v.Str, "+") {
			return rel + v.Str
		}
		str := strings.TrimPrefix(v.Str, "./")
		return rel + "/" + str
	case TypePath:
		return fmt.Sprintf("TYPE_PATH: %s\n", v.Str)
	default:
		panic(fmt.Sprintf("Value corrupt; unknown type %d", v.Type))
	}
}

func getRelative(currentTarget, comeFrom domain.Reference) string {
	rel, err := filepath.Rel(currentTarget.GetLocalPath(), comeFrom.GetLocalPath())
	if err != nil {
		panic(err)
	}
	if rel == "" || strings.HasPrefix(rel, "..") {
		return rel
	}
	return "./" + rel
}
