package ast_test

import (
	"strings"
	"testing"

	"github.com/earthly/earthly/ast"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	namedReader := namedStringReader{strings.NewReader("VERSION 0.6")}
	ver, err := ast.ParseVersionOpts(ast.FromReader(&namedReader))
	r := require.New(t)
	r.NoError(err)
	r.Len(ver.Args, 1)
	r.Equal("0.6", ver.Args[0])
	r.Nil(ver.SourceLocation)
}
