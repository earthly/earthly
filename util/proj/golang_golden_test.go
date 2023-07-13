package proj_test

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/earthly/earthly/util/proj"
)

const (
	goOut = "./testdata/golang.out"

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

func TestGolang_Targets(t *testing.T) {
	buf := bytes.NewBufferString(version)
	g := proj.NewGolang(proj.StdFS(), proj.StdExecer())
	tgts, err := g.Targets(context.TODO())
	if err != nil {
		t.Fatalf("failed to load golang targets: %v", err)
	}
	for _, tgt := range tgts {
		err := tgt.Format(context.TODO(), buf, "    ", 0)
		if err != nil {
			t.Fatalf("failed to format code: %v", err)
		}
		buf.WriteString("\n")
	}
	if *update {
		saveGoldenFile(t, goOut, buf.Bytes())
	}
	golden := goldenFile(t, goOut)
	if string(golden) != buf.String() {
		t.Fatalf("output did not match golden file. output:\n\n%v", buf.String())
	}
}
