package ast_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/ast/spec"
)

type namedStringReader struct {
	*strings.Reader
}

func (n *namedStringReader) Name() string {
	return "Earthfile"
}

var _ ast.NamedReader = &namedStringReader{}

func TestParse(t *testing.T) {

	tests := []struct {
		note      string
		earthfile string
		check     func(*require.Assertions, spec.Earthfile, error)
	}{
		{
			note: "it parses SET commands",
			earthfile: `
VERSION 0.7

foo:
    LET foo = bar
    SET foo = baz
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				foo := s.Targets[0]
				r.Len(foo.Recipe, 2)
				set := foo.Recipe[1]
				r.NotNil(set.Command)
				r.Equal("SET", set.Command.Name)
				r.Equal([]string{"foo", "=", "baz"}, set.Command.Args)
			},
		},
		{
			note: "it parses LET commands",
			earthfile: `
VERSION 0.7

LET foo = bar

foo:
    LET bacon = eggs
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.BaseRecipe, 1)
				global := s.BaseRecipe[0]
				r.NotNil(global.Command)
				r.Equal("LET", global.Command.Name)
				r.Equal([]string{"foo", "=", "bar"}, global.Command.Args)
				r.Len(s.Targets, 1)
				foo := s.Targets[0]
				r.Len(foo.Recipe, 1)
				let := foo.Recipe[0]
				r.NotNil(let.Command)
				r.Equal("LET", let.Command.Name)
				r.Equal([]string{"bacon", "=", "eggs"}, let.Command.Args)
			},
		},
		{
			note: "it safely ignores comments outside of documentation",
			earthfile: `
# this is an early comment.

# VERSION does not get documentation.
VERSION 0.6 # Trailing comments do not cause parsing errors at the top level
WORKDIR /tmp

# a comment before an IF or a FOR does not cause parser errors
IF foo
    RUN echo foo
END

bar:

baz:
    # comments in an otherwise empty target should be
    # ignored.

# foo - Comments between targets should not be parsed as
# documentation, even if they start with the target's name.

foo: # inline  comments do not consume newlines
    # RUN does not get documentation.
    RUN echo foo

    ARG foo=bar # inline comments should also be ignored.

    # Lonely comment blocks in
    # targets should be ignored.

    # Even if they don't have a trailing newline.`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 3)
				foo := s.Targets[2]
				r.Equal("foo", foo.Name)
				r.Equal("", foo.Docs)
			},
		},
		{
			note: "targets with leading whitespace cause errors",
			earthfile: `
VERSION 0.6

  foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.Error(err)
				r.ErrorContains(err, "no viable alternative at input '  '")
			},
		},
		{
			note: "it parses a basic target",
			earthfile: `
VERSION 0.6

foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Version.Args, 1)
				r.Equal("0.6", s.Version.Args[0])
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Len(target.Recipe, 1)
				recipe := target.Recipe[0]
				r.NotNil(recipe.Command)
				r.Equal("RUN", recipe.Command.Name)
				r.Equal([]string{"echo", "foo"}, recipe.Command.Args)
			},
		},
		{
			note: "nested quotes inside of shellouts do not break parent quotes",
			earthfile: `
VERSION 0.6

foo:
    RUN echo "$(echo "foo     bar")"
    ENV FOO="$(echo "foo     bar")"
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Len(target.Recipe, 2)
				run := target.Recipe[0]
				r.NotNil(run.Command)
				r.Equal("RUN", run.Command.Name)
				r.Equal([]string{"echo", `"$(echo "foo     bar")"`}, run.Command.Args)
				env := target.Recipe[1]
				r.Equal("ENV", env.Command.Name)
				r.Equal([]string{"FOO", "=", `"$(echo "foo     bar")"`}, env.Command.Args)
			},
		},
		{
			note: "multi key value pairs in ENV ",
			earthfile: `
VERSION 0.6

env:
  FROM alpine
  ENV GOLANG=1.22.2 \
    GO_VERSION=1.22.2 \
	GOOS=linux \
	GOARCH=amd64 \
    GO_DOWNLOAD_SHA256=5901c52b7a78002aeff14a21f93e0f064f74ce1360fce51c6ee68cd471216a17
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("env", target.Name)
				r.Len(target.Recipe, 1)
				env := target.Recipe[1]
				r.Equal("ENV", env.Command.Name)
				r.Equal([]string{"GOLANG", "=", `1.22.2 \
    GO_VERSION=1.22.2 \
	GOOS=linux \
	GOARCH=amd64 \
    GO_DOWNLOAD_SHA256=5901c52b7a78002aeff14a21f93e0f064f74ce1360fce51c6ee68cd471216a17
`}, env.Command.Args)
			},
		},
		{
			note: "nested shellouts inside of shellouts do not break parent shellouts",
			earthfile: `
VERSION 0.6

foo:
    RUN echo $(echo $(echo -n foo) $(echo -n bar))
    ENV FOO=$(echo $(echo -n foo) $(echo -n bar))
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Len(target.Recipe, 2)
				run := target.Recipe[0]
				r.NotNil(run.Command)
				r.Equal("RUN", run.Command.Name)
				r.Equal([]string{"echo", "$(echo $(echo -n foo) $(echo -n bar))"}, run.Command.Args)
				env := target.Recipe[1]
				r.Equal("ENV", env.Command.Name)
				r.Equal([]string{"FOO", "=", "$(echo $(echo -n foo) $(echo -n bar))"}, env.Command.Args)
			},
		},
		{
			note: "nested parens inside of quotes do not break parent shellouts",
			earthfile: `
VERSION 0.6

foo:
    ARG foo = "$(echo "()")"
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Len(target.Recipe, 1)
				run := target.Recipe[0]
				r.NotNil(run.Command)
				r.Equal("ARG", run.Command.Name)
				r.Equal([]string{"foo", "=", `"$(echo "()")"`}, run.Command.Args)
			},
		},
		{
			note: "ENV and ARG values retain inner whitespace",
			earthfile: `
VERSION 0.6

foo:
    ARG foo = $ ( foo )
    ENV foo = $ ( foo )
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Len(target.Recipe, 2)
				arg := target.Recipe[0]
				r.NotNil(arg.Command)
				r.Equal("ARG", arg.Command.Name)
				r.Equal([]string{"foo", "=", "$ ( foo )"}, arg.Command.Args)
				env := target.Recipe[1]
				r.Equal("ENV", env.Command.Name)
				r.Equal([]string{"foo", "=", "$ ( foo )"}, env.Command.Args)
			},
		},
		{
			note: "it successfully parses unindented comments mid-recipe",
			earthfile: `
VERSION 0.7

foo:
    RUN some_command
# Comment regarding something
    SAVE ARTIFACT /stuff
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
			},
		},
		{
			note: "it parses target documentation",
			earthfile: `
VERSION 0.6

# foo echoes 'foo'
foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Equal("foo echoes 'foo'\n", target.Docs)
			},
		},
		{
			note: "it respects leading whitespace in documentation",
			earthfile: `
VERSION 0.7

# foo outputs formatted JSON
#
# Sample output:
#
#     $ earthly +foo --json='{"a":"b","c":"d"}'
#     {
#         "a": "b",
#         "c": "d"
#     }
foo:
    ARG json
    RUN echo $json | jq .
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Equal(`foo outputs formatted JSON

Sample output:

    $ earthly +foo --json='{"a":"b","c":"d"}'
    {
        "a": "b",
        "c": "d"
    }
`, target.Docs)
			},
		},
		{
			note: "it parses documentation on later targets",
			earthfile: `
VERSION 0.6

bar:
    RUN echo bar

# foo echoes 'foo'
foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 2)
				target := s.Targets[1]
				r.Equal("foo", target.Name)
				r.Equal("foo echoes 'foo'\n", target.Docs)
			},
		},
		{
			note: "it parses multiline documentation",
			earthfile: `
VERSION 0.6

# foo echoes 'foo'
#
# and that's all.
foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Equal("foo echoes 'foo'\n\nand that's all.\n", target.Docs)
			},
		},
		{
			note: "it does not parse comments with empty lines after them as documentation",
			earthfile: `
VERSION 0.6

# foo echoes 'foo'

foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Equal("", target.Docs)
			},
		},
		{
			note: "it does not check the comment against the target name",
			earthfile: `
VERSION 0.6

# echoes 'foo'
foo:
    RUN echo foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Equal("foo", target.Name)
				r.Equal("echoes 'foo'\n", target.Docs)
			},
		},
		{
			note: "it skips comments that have different indentation",
			earthfile: `
VERSION 0.6

foo:
    RUN echo foo
    # this is a trailing comment in foo
# bar is a documented target
bar:
    RUN echo bar
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 2)
				target := s.Targets[1]
				r.Equal("bar", target.Name)
				r.Equal("bar is a documented target\n", target.Docs)
			},
		},
		{
			note: "it does not treat comments in otherwise-empty targets as documentation for the next target",
			earthfile: `
VERSION 0.7


foo:
    # bar is not a documentation line

bar:
    RUN echo bar
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 2)
				target := s.Targets[1]
				r.Equal("bar", target.Name)
				r.Equal("", target.Docs)
			},
		},
		{
			note: "it parses documentation on ARGs",
			earthfile: `
VERSION 0.6

foo:
    # foo is the argument that will be echoed
    ARG foo = bar
    RUN echo $foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				arg := target.Recipe[0]
				r.NotNil(arg.Command)
				r.Equal("ARG", arg.Command.Name)
				r.Equal("foo is the argument that will be echoed\n", arg.Command.Docs)
			},
		},
		{
			note: "it parses multiline documentation on global ARGs",
			earthfile: `
VERSION 0.7
FROM alpine:3.18

# globalArg is a documented global arg
# with multiple lines.
ARG --global globalArg
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.BaseRecipe, 2)
				arg := s.BaseRecipe[1]
				r.NotNil(arg.Command)
				r.Equal("ARG", arg.Command.Name)
				r.Equal("globalArg is a documented global arg\nwith multiple lines.\n", arg.Command.Docs)
			},
		},
		{
			note: "it parses documentation on SAVE ARTIFACT",
			earthfile: `
VERSION 0.6

foo:
    RUN echo foo > bar.txt
    # bar.txt will contain the output of this target
    SAVE ARTIFACT bar.txt
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				arg := target.Recipe[1]
				r.NotNil(arg.Command)
				r.Equal("SAVE ARTIFACT", arg.Command.Name)
				r.Equal("bar.txt will contain the output of this target\n", arg.Command.Docs)
			},
		},
		{
			note: "it parses documentation on SAVE IMAGE",
			earthfile: `
VERSION 0.6

foo:
    RUN echo foo > bar.txt
    # foo is an image that contains a bar.txt file
    SAVE IMAGE foo
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				arg := target.Recipe[1]
				r.NotNil(arg.Command)
				r.Equal("SAVE IMAGE", arg.Command.Name)
				r.Equal("foo is an image that contains a bar.txt file\n", arg.Command.Docs)
			},
		},
		{
			note: "complex character sequences in single quotes",
			earthfile: `VERSION 0.8

target:
  RUN find . -type f -iname '*.md' | xargs -n 1 sed -i 's/{[^}]*}//g'
  RUN find . -type f -iname '*.md' | xargs vale --config /etc/vale/vale.ini --output line --minAlertLevel error
`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				r.Len(s.Targets[0].Recipe, 2)
			},
		},
		{
			note: "regression test for single-quoted #",
			earthfile: `VERSION 0.8

test:
    FROM debian:9
    RUN set -x \
     && sed -i \
            -e 's, universe multiverse, universe # multiverse,' \
            /etc/apt/sources.list
    SAVE IMAGE --push blah`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 3)
				// Confirm that the single-quoted string is intact
				r.Contains(target.Recipe[1].Command.Args, `'s, universe multiverse, universe # multiverse,'`)
			},
		},
		{
			note: "regression test for escaped # in $()",
			earthfile: `VERSION 0.8

thebug:
    FROM alpine
    ARG myarg=$(echo "a#b#c" | cut -f2 -d\#)
    RUN touch /some-file
    RUN echo "myarg is \"$myarg\""
    RUN test -f /some-file`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 5)
				// Confirm that the escaped expression is intact
				r.Contains(target.Recipe[1].Command.Args, `$(echo "a#b#c" | cut -f2 -d\#)`)
			},
		},
		{
			note: "regression test for single-quoted string in shell expression",
			earthfile: `VERSION 0.8
FROM alpine

arg-plain:
    ARG val=$(echo run | tr -d '"')
    RUN echo $val`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				// Confirm that the single-quoted string is intact
				r.Contains(target.Recipe[0].Command.Args, `$(echo run | tr -d '"')`)
			},
		},
		{
			note: "regression test for single-quoted string in RUN",
			earthfile: `VERSION 0.8
FROM alpine

run-plain:
    RUN echo run | tr -d '"'`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 1)
			},
		},
		{
			note: "regression test for escaped double-quoted strings in shell expression",
			earthfile: `VERSION 0.8
FROM alpine

arg-esc:
    ARG val=$(echo single | tr -d "\"")
    RUN echo $val`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				// Confirm that the single-quoted string is intact
				r.Contains(target.Recipe[0].Command.Args, `$(echo single | tr -d "\"")`)
			},
		},
		{
			note: "regression test for escaped \\ & double-quotes in shell expression",
			earthfile: `VERSION 0.8
FROM alpine

arg-esc:
    ARG val=$(echo single | tr -d "\\\"")
    RUN echo $val`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				// Confirm that the single-quoted string is intact
				r.Contains(target.Recipe[0].Command.Args, `$(echo single | tr -d "\\\"")`)
			},
		},
		{
			note: "regression test for single-quoted commands",
			earthfile: `VERSION 0.8

FROM alpine

test:
  RUN 'echo "message'
  RUN 'echo "message"'`,
			check: func(r *require.Assertions, s spec.Earthfile, err error) {
				r.NoError(err)
				r.Len(s.Targets, 1)
				target := s.Targets[0]
				r.Len(target.Recipe, 2)
				// Confirm that the single-quoted strings are intact
				r.Contains(target.Recipe[0].Command.Args, `'echo "message'`)
				r.Contains(target.Recipe[1].Command.Args, `'echo "message"'`)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.note, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			r := namedStringReader{strings.NewReader(test.earthfile)}
			s, err := ast.ParseOpts(ctx, ast.FromReader(&r))
			test.check(require.New(t), s, err)
		})
	}
}
