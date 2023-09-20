package ast_test

import (
	"context"
	"io"
	"testing"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/earthly/earthly/ast"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

func TestParse(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
		reader *mockNamedReader
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		return testCtx{
			t:      t,
			expect: expect.New(t),
			reader: newMockNamedReader(t, timeout),
		}
	})
	defer o.Run()

	o.Spec("it parses SET commands", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.7

foo:
    LET foo = bar
    SET foo = baz
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(1))
		foo := f.Targets[0]
		tt.expect(foo.Recipe).To(haveLen(2))
		set := foo.Recipe[1]
		tt.expect(set.Command).To(not(beNil()))
		tt.expect(set.Command.Name).To(equal("SET"))
		tt.expect(set.Command.Args).To(equal([]string{"foo", "=", "baz"}))
	})

	o.Spec("it parses LET commands", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.7

LET foo = bar

foo:
    LET bacon = eggs
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.BaseRecipe).To(haveLen(1))
		global := f.BaseRecipe[0]
		tt.expect(global.Command).To(not(beNil()))
		tt.expect(global.Command.Name).To(equal("LET"))
		tt.expect(global.Command.Args).To(equal([]string{"foo", "=", "bar"}))

		tt.expect(f.Targets).To(haveLen(1))
		foo := f.Targets[0]
		tt.expect(foo.Recipe).To(haveLen(1))
		let := foo.Recipe[0]
		tt.expect(let.Command).To(not(beNil()))
		tt.expect(let.Command.Name).To(equal("LET"))
		tt.expect(let.Command.Args).To(equal([]string{"bacon", "=", "eggs"}))
	})

	o.Spec("it safely ignores comments outside of documentation", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
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

    # Even if they don't have a trailing newline.`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(3))
		foo := f.Targets[2]
		tt.expect(foo.Name).To(equal("foo"))
		tt.expect(foo.Docs).To(equal(""))
	})

	o.Spec("targets with leading whitespace cause errors", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

  foo:
    RUN echo foo
`))
		_, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(haveOccurred())
		tt.expect(err.Error()).To(containSubstring("no viable alternative at input '  '"))
	})

	o.Spec("it parses a basic target", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo foo
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Version.Args).To(haveLen(1))
		tt.expect(f.Version.Args[0]).To(equal("0.6"))

		tt.expect(f.Targets).To(haveLen(1))
		tgt := f.Targets[0]
		tt.expect(tgt.Name).To(equal("foo"))
		tt.expect(tgt.Recipe).To(haveLen(1))
		rcp := tgt.Recipe[0]
		tt.expect(rcp.Command).To(not(beNil()))
		tt.expect(rcp.Command.Name).To(equal("RUN"))
		tt.expect(rcp.Command.Args).To(equal([]string{"echo", "foo"}))
	})

	o.Spec("nested quotes inside of shellouts do not break parent quotes", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo "$(echo "foo     bar")"
    ENV FOO="$(echo "foo     bar")"
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(1))
		tgt := f.Targets[0]
		tt.expect(tgt.Name).To(equal("foo"))
		tt.expect(tgt.Recipe).To(haveLen(2))

		run := tgt.Recipe[0]
		tt.expect(run.Command).To(not(beNil()))
		tt.expect(run.Command.Name).To(equal("RUN"))
		tt.expect(run.Command.Args).To(equal([]string{"echo", `"$(echo "foo     bar")"`}))

		env := tgt.Recipe[1]
		tt.expect(env.Command).To(not(beNil()))
		tt.expect(env.Command.Name).To(equal("ENV"))
		tt.expect(env.Command.Args).To(equal([]string{"FOO", "=", `"$(echo "foo     bar")"`}))
	})

	o.Spec("nested shellouts inside of shellouts do not break parent shellouts", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo $(echo $(echo -n foo) $(echo -n bar))
    ENV FOO=$(echo $(echo -n foo) $(echo -n bar))
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(1))
		tgt := f.Targets[0]
		tt.expect(tgt.Name).To(equal("foo"))
		tt.expect(tgt.Recipe).To(haveLen(2))

		run := tgt.Recipe[0]
		tt.expect(run.Command).To(not(beNil()))
		tt.expect(run.Command.Name).To(equal("RUN"))
		tt.expect(run.Command.Args).To(equal([]string{"echo", "$(echo $(echo -n foo) $(echo -n bar))"}))

		env := tgt.Recipe[1]
		tt.expect(env.Command).To(not(beNil()))
		tt.expect(env.Command.Name).To(equal("ENV"))
		tt.expect(env.Command.Args).To(equal([]string{"FOO", "=", "$(echo $(echo -n foo) $(echo -n bar))"}))
	})

	o.Spec("nested parens inside of quotes do not break parent shellouts", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    ARG foo = "$(echo "()")"
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(1))
		tgt := f.Targets[0]
		tt.expect(tgt.Name).To(equal("foo"))
		tt.expect(tgt.Recipe).To(haveLen(1))

		run := tgt.Recipe[0]
		tt.expect(run.Command).To(not(beNil()))
		tt.expect(run.Command.Name).To(equal("ARG"))
		tt.expect(run.Command.Args).To(equal([]string{"foo", "=", `"$(echo "()")"`}))
	})

	o.Spec("ENV and ARG values retain inner whitespace", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    ARG foo = $ ( foo )
    ENV foo = $ ( foo )
`))
		f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))

		tt.expect(f.Targets).To(haveLen(1))
		tgt := f.Targets[0]
		tt.expect(tgt.Name).To(equal("foo"))
		tt.expect(tgt.Recipe).To(haveLen(2))

		arg := tgt.Recipe[0]
		tt.expect(arg.Command).To(not(beNil()))
		tt.expect(arg.Command.Name).To(equal("ARG"))
		tt.expect(arg.Command.Args).To(equal([]string{"foo", "=", "$ ( foo )"}))

		env := tgt.Recipe[1]
		tt.expect(env.Command).To(not(beNil()))
		tt.expect(env.Command.Name).To(equal("ENV"))
		tt.expect(env.Command.Args).To(equal([]string{"foo", "=", "$ ( foo )"}))
	})

	o.Spec("it successfully parses unindented comments mid-recipe", func(tt testCtx) {
		mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.7

foo:
    RUN some_command
# Comment regarding something
    SAVE ARTIFACT /stuff
`))
		_, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))
	})

	o.Group("target docs", func() {
		o.Spec("it parses target documentation", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

# foo echoes 'foo'
foo:
    RUN echo foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal("foo echoes 'foo'\n"))
		})

		o.Spec("it respects leading whitespace in documentation", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
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
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal(`foo outputs formatted JSON

Sample output:

    $ earthly +foo --json='{"a":"b","c":"d"}'
    {
        "a": "b",
        "c": "d"
    }
`))
		})

		o.Spec("it parses documentation on later targets", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

bar:
    RUN echo bar

# foo echoes 'foo'
foo:
    RUN echo foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(2))
			tgt := f.Targets[1]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal("foo echoes 'foo'\n"))
		})

		o.Spec("it parses multiline documentation", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

# foo echoes 'foo'
#
# and that's all.
foo:
    RUN echo foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal("foo echoes 'foo'\n\nand that's all.\n"))
		})

		o.Spec("it does not parse comments with empty lines after them as documentation", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

# foo echoes 'foo'

foo:
    RUN echo foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal(""))
		})

		o.Spec("it does not check the comment against the target name", func(tt testCtx) {
			// It felt cleaner to check the doc comment's first word against the
			// target's name at a higher level where we can display hints to the
			// user about why the comments are not considered documentation.
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

# echoes 'foo'
foo:
    RUN echo foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Name).To(equal("foo"))
			tt.expect(tgt.Docs).To(equal("echoes 'foo'\n"))
		})

		o.Spec("it skips comments that have different indentation", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo foo
    # this is a trailing comment in foo
# bar is a documented target
bar:
    RUN echo bar
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(2))
			tgt := f.Targets[1]
			tt.expect(tgt.Name).To(equal("bar"))
			tt.expect(tgt.Docs).To(equal("bar is a documented target\n"))
		})

		o.Spec("it does not treat comments in otherwise-empty targets as documentation for the next target", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.7


foo:
    # bar is not a documentation line

bar:
    RUN echo bar
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(2))
			tgt := f.Targets[1]
			tt.expect(tgt.Name).To(equal("bar"))
			tt.expect(tgt.Docs).To(equal(""))
		})

		o.Spec("it parses documentation on ARGs", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    # foo is the argument that will be echoed
    ARG foo = bar
    RUN echo $foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Recipe).To(haveLen(2))
			arg := tgt.Recipe[0]
			tt.expect(arg.Command).To(not(beNil()))
			tt.expect(arg.Command.Name).To(equal("ARG"))
			tt.expect(arg.Command.Docs).To(equal("foo is the argument that will be echoed\n"))
		})

		o.Spec("it parses multiline documentation on global ARGs", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.7
FROM alpine:3.18

# globalArg is a documented global arg
# with multiple lines.
ARG --global globalArg
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.BaseRecipe).To(haveLen(2))
			arg := f.BaseRecipe[1]
			tt.expect(arg.Command).To(not(beNil()))
			tt.expect(arg.Command.Name).To(equal("ARG"))
			tt.expect(arg.Command.Docs).To(equal("globalArg is a documented global arg\nwith multiple lines.\n"))
		})

		o.Spec("it parses documentation on SAVE ARTIFACT", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo foo > bar.txt
    # bar.txt will contain the output of this target
    SAVE ARTIFACT bar.txt
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Recipe).To(haveLen(2))
			arg := tgt.Recipe[1]
			tt.expect(arg.Command).To(not(beNil()))
			tt.expect(arg.Command.Name).To(equal("SAVE ARTIFACT"))
			tt.expect(arg.Command.Docs).To(equal("bar.txt will contain the output of this target\n"))
		})

		o.Spec("it parses documentation on SAVE IMAGE", func(tt testCtx) {
			mockEarthfile(tt.t, tt.reader, []byte(`
VERSION 0.6

foo:
    RUN echo foo > bar.txt
    # foo is an image that contains a bar.txt file
    SAVE IMAGE foo
`))
			f, err := ast.ParseOpts(context.Background(), ast.FromReader(tt.reader))
			tt.expect(err).To(not(haveOccurred()))

			tt.expect(f.Targets).To(haveLen(1))
			tgt := f.Targets[0]
			tt.expect(tgt.Recipe).To(haveLen(2))
			arg := tgt.Recipe[1]
			tt.expect(arg.Command).To(not(beNil()))
			tt.expect(arg.Command.Name).To(equal("SAVE IMAGE"))
			tt.expect(arg.Command.Docs).To(equal("foo is an image that contains a bar.txt file\n"))
		})
	})
}

// mockEarthfile mocks out an Earthfile for testing purposes.
func mockEarthfile(t *testing.T, reader *mockNamedReader, earthfileBody []byte) {
	t.Helper()

	pers.ConsistentlyReturn(t, reader.NameOutput, "Earthfile")
	handleMockFile(t, reader, earthfileBody)
}

// handleMockFile helps us perform slightly more black-box testing by handling a
// mockNamedReader as if it were a file-like io.ReadSeeker. This way, we don't
// need to know in the test how many times the file is seeked back to zero and
// re-read.
//
// This cannot handle non-zero seeks and will fail if it receives a non-zero
// seek call.
func handleMockFile(t *testing.T, r *mockNamedReader, body []byte) {
	t.Helper()

	idx := 0
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() {
		for {
			select {
			case <-r.ReadCalled:
				buff := <-r.ReadInput.Buff
				cpyEnd := idx + len(buff)
				if cpyEnd > len(body) {
					cpy := body[idx:]
					copy(buff, cpy)
					idx = len(body)
					pers.Return(r.ReadOutput, len(cpy), io.EOF)
					continue
				}
				copy(buff, body[idx:cpyEnd])
				idx = cpyEnd
				pers.Return(r.ReadOutput, len(buff), nil)
			case <-r.SeekCalled:
				offset := <-r.SeekInput.Offset
				whence := <-r.SeekInput.Whence
				if offset != 0 || whence != 0 {
					t.Fatalf("ast: handleMockFile cannot handle non-zero offset or whence values in calls to Seek(); got offset=%d, whence=%d", offset, whence)
				}
				idx = 0
				pers.Return(r.SeekOutput, 0, nil)
			case <-ctx.Done():
				return
			}
		}
	}()
}
