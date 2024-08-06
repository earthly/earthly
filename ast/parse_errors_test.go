package ast_test

import (
	"context"
	"strings"
	"testing"

	"github.com/earthly/earthly/ast"
	"github.com/stretchr/testify/require"
)

func TestParserErrors(t *testing.T) {

	tests := []struct {
		name         string
		earthfile    string
		expectedHint string
	}{
		{
			name: "missing newline token",
			earthfile: `
VERSION 0.7

test:
    FROM alpine
    IF $foo END
`,
			expectedHint: `
  Hints:
  - I couldn't find a pattern that completes the current statement - check your quote pairs, paren pairs, and newlines
  - I parsed 'END' as a word, but it looks like it should be a keyword - is it on the wrong line?
`,
		},
		{
			name: "key-value with missing EQUALS",
			earthfile: `
VERSION 0.7

test:
    FROM alpine
    LABEL a
`,
			expectedHint: `
  Hint: I got lost looking for '=' - did you define a key/value pair without a value?
`,
		},
		{
			name: "unrecognized keyword",
			earthfile: `
VERSION 0.7

test:
	RIN apk --update add build-base cmake bash
`,
			expectedHint: `
Hint: 'RIN ' is not a recognized keyword.
`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			namedReader := namedStringReader{strings.NewReader(test.earthfile)}
			_, err := ast.ParseOpts(context.Background(), ast.FromReader(&namedReader))
			r := require.New(t)
			r.Error(err)
			r.ErrorContains(err, strings.TrimSpace(test.expectedHint))
		})
	}
}
