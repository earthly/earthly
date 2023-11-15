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

type ProtoEnum interface {
	fmt.Stringer
}

// EnumToStringFunc takes an ProtoEnum and returns a string
type EnumToStringFunc func(item ProtoEnum) string

// Title takes an enum and returns its string value in title mode
func Title(e ProtoEnum) string {
	return pretty(titleCaser, e)
}

// Lower takes an enum and returns its string value in lower case mode
func Lower(e ProtoEnum) string {
	return pretty(lowerCaser, e)
}

// EnumToStringArray takes an array of enum values and returns an array of their string values after applying EnumToStringFunc on each item
func EnumToStringArray[T ProtoEnum](items []T, f EnumToStringFunc) []string {
	strs := make([]string, 0, len(items))
	for _, item := range items {
		strs = append(strs, f(item))
	}
	return strs
}

func pretty(caser caser, e ProtoEnum) string {
	val := reflect.TypeOf(e).Name()
	idx := strings.Index(val, "_")
	val = val[idx+1:]
	prefix := fmt.Sprintf("%s_", strcase.ToScreamingSnake(val))
	return caser.String(strings.ReplaceAll(strings.TrimPrefix(e.String(), prefix), "_", " "))
}
