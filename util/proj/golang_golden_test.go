package proj_test

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/earthly/earthly/util/proj"
	"github.com/fatih/color"
	"github.com/kylelemons/godebug/diff"
)

const (
	goOut_base  = "./testdata/golang_base.out"
	goOut_named = "./testdata/golang_named.out"

	version = "VERSION --arg-scope-and-set 0.7\n\n"
)

var (
	update = flag.Bool("update", false, "Update the testdata for golden tests")
)

func goldenFile(t *testing.T, path string) []byte {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("got error opening golden file %q: %v", path, err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("got error reading golden file %q: %v", path, err)
	}
	return b
}

func saveGoldenFile(t *testing.T, path string, b []byte) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("got error creating golden file %q: %v", path, err)
	}
	for len(b) > 0 {
		n, err := f.Write(b)
		if err != nil {
			t.Fatalf("got error writing golden file %q: %v", path, err)
		}
		b = b[n:]
	}
}

func colorDiff(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		switch line[0] {
		case '-':
			lines[i] = color.RedString(line)
		case '+':
			lines[i] = color.GreenString(line)
		default:
		}
	}
	return strings.Join(lines, "\n")
}

func matchGolden(t *testing.T, actualBytes []byte, path string) {
	goldenBytes := goldenFile(t, path)
	golden := string(goldenBytes)
	actual := string(actualBytes)
	if golden != actual {
		t.Fatalf("output did not match golden file. diff:\n\n%v", colorDiff(diff.Diff(golden, actual)))
	}
}

func TestGolang_Targets_Base(t *testing.T) {
	buf := bytes.NewBufferString(version)
	g := proj.NewGolang(proj.StdFS(), proj.StdExecer())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tgts, err := g.Targets(ctx)
	if err != nil {
		t.Fatalf("failed to load golang targets: %v", err)
	}
	for i, tgt := range tgts {
		tgt.SetPrefix(ctx, "")
		if i > 0 {
			buf.WriteString("\n")
		}
		err := tgt.Format(ctx, buf, "    ", 0)
		if err != nil {
			t.Fatalf("failed to format code: %v", err)
		}
	}
	if *update {
		saveGoldenFile(t, goOut_base, buf.Bytes())
	}
	matchGolden(t, buf.Bytes(), goOut_base)
}

func TestGolang_Targets_Named(t *testing.T) {
	buf := bytes.NewBufferString(version)
	g := proj.NewGolang(proj.StdFS(), proj.StdExecer())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pfx := g.Type(ctx)

	tgts, err := g.Targets(ctx)
	if err != nil {
		t.Fatalf("failed to load golang targets: %v", err)
	}
	for i, tgt := range tgts {
		tgt.SetPrefix(ctx, pfx)
		if i > 0 {
			buf.WriteString("\n")
		}
		err := tgt.Format(ctx, buf, "    ", 0)
		if err != nil {
			t.Fatalf("failed to format code: %v", err)
		}
	}
	if *update {
		saveGoldenFile(t, goOut_named, buf.Bytes())
	}
	matchGolden(t, buf.Bytes(), goOut_named)
}
