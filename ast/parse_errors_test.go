package ast_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/earthly/earthly/ast"
	"github.com/poy/onpar/v2"
	"github.com/poy/onpar/v2/expect"
)

type namedReader struct {
	*bytes.Reader

	name string
}

func (b namedReader) Name() string {
	return b.name
}

func TestParserErrors(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		return testCtx{
			t:      t,
			expect: expect.New(t),
		}
	})

	for _, tt := range []struct {
		name         string
		body         string
		expectedHint string
	}{
		{
			name: "missing newline token",
			body: `
VERSION 0.7

test:
    FROM alpine
    IF $foo END
`,
			expectedHint: `
Hints:
  - I couldn't find a pattern that completes the current statement - check your quote pairs, paren pairs, and newlines
  - I parsed 'END' as a word, but it looks like it should be a keyword - is it on the wrong line?`,
		},
		{
			name: "key-value with missing EQUALS",
			body: `
VERSION 0.7

test:
    FROM alpine
    LABEL a
`,
			expectedHint: `
Hints:
  - I got lost looking for '=' - did you define a key/value pair without a value?`,
		},
	} {
		tt := tt
		o.Spec(tt.name, func(tc testCtx) {
			b := namedReader{
				Reader: bytes.NewReader([]byte(tt.body)),
				name:   strings.Replace(tt.name, " ", "_", -1) + ".earth",
			}
			_, err := ast.ParseOpts(context.Background(), ast.FromReader(b))
			tc.expect(err).To(haveOccurred())
			tc.expect(err.Error()).To(endWith(tt.expectedHint))
		})
	}
}
