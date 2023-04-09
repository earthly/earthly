package ast_test

import (
	"io"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/earthly/earthly/ast"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

// timeout is used to catch any hung tests that are not caught by the deadlock
// detector.
const timeout = time.Second

func TestParseVersion(t *testing.T) {
	type testCtx struct {
		expect expect.Expectation
		reader *mockNamedReader
	}

	o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) testCtx {
		return testCtx{
			expect: expect.New(t),
			reader: newMockNamedReader(t, timeout),
		}
	})
	defer o.Run()

	o.Spec("it parses a basic version", func(tt testCtx) {
		go func() {
			const version = "VERSION 0.6"
			var buf []byte
			tt.expect(tt.reader).To(haveMethodExecuted("Read", storeArgs(&buf), within(timeout)))
			copy(buf, []byte(version))
			pers.Return(tt.reader.ReadOutput, len(version), io.EOF)
		}()
		ver, err := ast.ParseVersionOpts(ast.FromReader(tt.reader))
		tt.expect(err).To(not(haveOccurred()))
		tt.expect(ver.Args).To(haveLen(1))
		tt.expect(ver.Args[0]).To(equal("0.6"))
		tt.expect(ver.SourceLocation).To(beNil())
	})
}
