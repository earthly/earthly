package proj_test

import (
	"bytes"
	"context"
	"io/fs"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/earthly/earthly/util/proj"
	"github.com/pkg/errors"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

const (
	timeout     = time.Second
	mockTimeout = 5 * time.Second
)

func TestGolang(t *testing.T) {
	type testCtx struct {
		*testing.T
		ctx    context.Context
		expect expect.Expectation
		golang *proj.Golang

		fs   *mockFS
		exec *mockExecer

		cancel func()
	}

	o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) testCtx {
		fs := newMockFS(t, mockTimeout)
		exec := newMockExecer(t, mockTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		return testCtx{
			T:      t,
			ctx:    ctx,
			expect: expect.New(t),
			golang: proj.NewGolang(fs, exec),
			fs:     fs,
			exec:   exec,
			cancel: cancel,
		}
	})
	defer o.Run()

	o.AfterEach(func(t testCtx) {
		t.cancel()
	})

	o.Spec("Type", func(t testCtx) {
		t.expect(t.golang.Type(t.ctx)).To(equal("go"))
	})

	o.Group("ForDir", func() {
		o.Spec("it skips projects without a go.mod", func(t testCtx) {
			pers.Return(t.fs.StatOutput, nil, fs.ErrNotExist)
			_, err := t.golang.ForDir(t.ctx, ".")
			t.expect(t.fs).To(haveMethodExecuted("Stat", withArgs("go.mod")))
			t.expect(err).To(beErr(proj.ErrSkip))
		})

		o.Spec("it errors if reading go.mod fails", func(t testCtx) {
			pers.Return(t.fs.StatOutput, nil, errors.New("boom"))
			_, err := t.golang.ForDir(t.ctx, ".")
			t.expect(t.fs).To(haveMethodExecuted("Stat", withArgs("go.mod")))
			t.expect(err).To(haveOccurred())
			t.expect(err).To(not(beErr(proj.ErrSkip)))
		})

		o.Spec("it errors if 'go' is not available", func(t testCtx) {
			pers.Return(t.fs.StatOutput, nil, nil)
			cmd := newMockCmd(t, mockTimeout)
			pers.Return(t.exec.CommandOutput, cmd)
			const projDir = "some/path/to/a/project"
			stdout := bytes.NewBufferString(projDir)
			pers.Return(cmd.RunOutput, stdout, nil, fs.ErrNotExist)
			_, err := t.golang.ForDir(t.ctx, ".")
			t.expect(t.fs).To(haveMethodExecuted("Stat", withArgs("go.mod")))
			t.expect(t.exec).To(haveMethodExecuted("Command", withArgs("go", "list", "-f", "{{.Dir}}")))
			t.expect(cmd).To(haveMethodExecuted("Run"))
			t.expect(err).To(haveOccurred())
			t.expect(err).To(not(beErr(proj.ErrSkip)))
		})
	})
}
