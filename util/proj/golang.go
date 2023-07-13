package proj

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/ast/hint"
	"github.com/pkg/errors"
)

const (
	goMod = "go.mod"
	goSum = "go.sum"

	goBase = `
LET go_version = 1.20
LET distro = alpine3.18

FROM golang:${go_version}-${distro}
WORKDIR /go-workdir`

	goDepsBlock = `
    # Copying only go.mod and go.sum means that the cache for this
    # target will only be busted when go.mod/go.sum change. This
    # means that we can cache the results of 'go mod download'.
    COPY go.mod .
    # Projects with no external dependencies do not have a go.sum.
    COPY --if-exists go.sum .
    RUN go mod download`

	goDepsFromTgt = `
go-deps:
    FROM +%s
` + goDepsBlock

	goDepsFromBase = `
go-deps:` + goDepsBlock

	goTestBase = `
go-test-base:
    FROM +go-deps

    # gcc and g++ are required for -race.
    RUN apk add --update gcc g++

    # This copies the whole project. If you want better caching, try
    # limiting this to _just_ files required by your go tests.
    COPY . .`

	goTestRace = `
# go-test-race runs 'go test -race'.
go-test-race:
    FROM +go-test-base

    # package sets the package that tests will run against.
    ARG package = ./...

    RUN go test -race "$package"`

	goTestIntegration = `
# go-test-integration runs 'go test -tags integration'.
go-test-integration:
    FROM +go-test-base

    # package sets the package that tests will run against.
    ARG package = ./...

    RUN go test -tags integration "$package"`

	goTest = `
# go-test runs all go test targets
go-test:
    BUILD +go-test-race
    BUILD +go-test-integration`

	goProjBase = `
go-proj-base:
    FROM +go-deps

    # This copies the whole project. If you want better caching, try
    # limiting this to _just_ files required by your go project.
    COPY . .`

	goTidy = `
# go-mod-tidy runs 'go mod tidy', saving go.mod and go.sum locally.
go-mod-tidy:
    FROM +go-proj-base

    RUN go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT --if-exists go.sum AS LOCAL go.sum`

	goBuild = `
# go-build runs 'go build ./...', saving artifacts locally.
go-build:
    FROM +go-proj-base

    ENV GOBIN = "/tmp/build"
    RUN go install ./...

    # outputDir sets the directory that build artifacts will be saved to.
    ARG outputDir = "./build"

    FOR bin IN $(ls -1 "/tmp/build")
        SAVE ARTIFACT "/tmp/build/${bin}" AS LOCAL "${outputDir}/${bin}"
    END`
)

// Golang is used to auto-generate Earthfiles for go projects.
type Golang struct {
	root   string
	fs     FS
	execer Execer
}

// NewGolang returns a new Golang.
func NewGolang(fs FS, execer Execer) *Golang {
	return &Golang{
		fs:     fs,
		execer: execer,
	}
}

// ForDir returns a Project for the given directory. It returns ErrSkip if the
// directory does not contain a go project.
func (g *Golang) ForDir(ctx context.Context, dir string) (Project, error) {
	_, err := g.fs.Stat(filepath.Join(dir, goMod))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.Wrap(ErrSkip, "no go.mod found")
	}
	out, _, err := g.execer.Command("go", "list", "-f", "{{.Dir}}").Run(ctx)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, hint.Wrap(errors.Wrap(err, "go.mod and go.sum exist, but go is not installed"),
			"go must be installed for 'go list' so that earthly can read information about your go project",
		)
	}
	rootBytes, err := io.ReadAll(out)
	if err != nil {
		return nil, errors.Wrap(err, "could not read go project root directory")
	}
	root, err := filepath.Abs(strings.TrimSpace(string(rootBytes)))
	if err != nil {
		return nil, errors.Wrapf(err, "could not get absolute path for directory %q", string(rootBytes))
	}
	return &Golang{
		root:   root,
		fs:     g.fs,
		execer: g.execer,
	}, nil
}

// Root returns the root path of this Golang project.
func (g *Golang) Root(context.Context) string {
	return g.root
}

// BaseBlock returns the block of commands that need to be in the base target
// for go targets.
func (g *Golang) BaseBlock(context.Context) (Formatter, error) {
	return strFormatter(goBase), nil
}

// Targets returns the targets that should be used for this Golang project.
func (g *Golang) Targets(_ context.Context, baseName string) ([]Formatter, error) {
	goDeps := goDepsFromBase
	if baseName != "" {
		goDeps = fmt.Sprintf(goDepsFromTgt, baseName)
	}
	return []Formatter{
		strFormatter(goDeps),
		strFormatter(goTestBase),
		strFormatter(goTestRace),
		strFormatter(goTestIntegration),
		strFormatter(goTest),
		strFormatter(goProjBase),
		strFormatter(goTidy),
		strFormatter(goBuild),
	}, nil
}

type strFormatter string

func (f strFormatter) Format(_ context.Context, w io.Writer, indent string, level int) error {
	s := strings.TrimSpace(string(f))
	if level > 0 {
		fullIndent := strings.Repeat(indent, level)
		lines := strings.Split(s, "\n")
		for i, l := range lines {
			if len(l) == 0 {
				continue
			}
			lines[i] = fullIndent + l
		}
		s = strings.Join(lines, "\n")
	}
	s += "\n"
	b := []byte(s)
	for len(b) > 0 {
		n, err := w.Write(b)
		if err != nil {
			return errors.Wrap(err, "errored writing code to writer")
		}
		b = b[n:]
	}
	return nil
}
