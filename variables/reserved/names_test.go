package reserved

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

// getConsts parses names.go and returns all constants.
// Unfortunately reflect is unable to work on the package level.
func getConsts() (map[string]string, error) {
	_, filename, _, _ := runtime.Caller(0)
	namesPath := path.Join(path.Dir(filename), "names.go")
	namesSrc, err := os.ReadFile(namesPath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", namesSrc, 0)
	if err != nil {
		return nil, err
	}

	consts := map[string]string{}
	for name, obj := range f.Scope.Objects {
		if obj.Kind == ast.Con {
			val := obj.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
			parsedVal, ok := trimDoubleQuotes(val)
			if !ok {
				return nil, fmt.Errorf("failed to parse %s", val)
			}
			consts[name] = parsedVal
		}
	}
	return consts, nil
}

// trimDoubleQuotes takes a string such as "abc" and returns abc
// this is a hack due to the ast ast.BasicLit returning the literal go source code.
// There is probably a better way to do this, but Alex couldn't quickly figure it out.
// This function should not be used outside of testing the very-specific use-case,
// it is very fragile, and is only designed to work with basic strings that are only ever
// defined on a single line in the go source code (which works fine for builtin ARG names).
func trimDoubleQuotes(s string) (string, bool) {
	n := len(s)
	if n < 2 {
		return "", false
	}
	if s[0] != '"' || s[n-1] != '"' {
		return "", false
	}
	return s[1:(n - 1)], true
}

func isUpper(s string) bool {
	return strings.ToUpper(s) == s
}

// TestAllConstsInMap tests that this package only defines constants
// that represent reserved ARG names, and that they exist in the map that IsBuiltIn
// relies on.
func TestAllConstsInMap(t *testing.T) {
	consts, err := getConsts()
	Nil(t, err)

	constsValues := map[string]struct{}{}

	// tests all consts return true for IsBuiltIn (i.e. all consts are in args map)
	for _, v := range consts {
		t.Run(v, func(t *testing.T) {
			True(t, IsBuiltIn(v))
			True(t, isUpper(v))
		})
		constsValues[v] = struct{}{}
	}

	// tests all values in args are defined as consts
	for k := range args {
		t.Run(k, func(t *testing.T) {
			_, exists := constsValues[k]
			True(t, exists)
		})
	}
}
