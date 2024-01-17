package ast

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	namedReader := namedStringReader{strings.NewReader("VERSION 0.6")}
	ver, err := ParseVersionOpts(FromReader(&namedReader))
	r := require.New(t)
	r.NoError(err)
	r.Len(ver.Args, 1)
	r.Equal("0.6", ver.Args[0])
	r.Nil(ver.SourceLocation)
}
