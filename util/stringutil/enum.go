package stringutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)
var lowerCaser = cases.Lower(language.English)

type caser interface {
	String(s string) string
}

type Enum interface {
	String() string
}

type EnumToStringFunc func(item Enum) string

func Title(e Enum) string {
	return pretty(titleCaser, e)
}

func Lower(e Enum) string {
	return pretty(lowerCaser, e)
}

func EnumToStringArray[T Enum](items []T, f EnumToStringFunc) []string {
	strs := make([]string, 0, len(items))
	for _, item := range items {
		strs = append(strs, f(item))
	}
	return strs
}

func pretty(caser caser, e Enum) string {
	val := reflect.TypeOf(e).Name()
	idx := strings.Index(val, "_")
	val = val[idx+1:]
	prefix := fmt.Sprintf("%s_", strcase.ToScreamingSnake(val))
	return caser.String(strings.ReplaceAll(strings.TrimPrefix(e.String(), prefix), "_", " "))
}
