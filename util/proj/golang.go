package proj

import (
	"context"
	"text/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/ast/hint"
	"github.com/pkg/errors"
)

const (
	goMod      = "go.mod"
	goSum      = "go.sum"
	goCache    = "/.go-cache"
	goModCache = "/.go-mod-cache"

	goBase = `
{{- $indent := and .Prefix .Indent}}{{/* if .Prefix is empty string, empty string; otherwise .Indent */}}
{{- if .Prefix }}{{.Prefix}}base:
{{ end -}}
{{$indent}}LET go_version = 1.20
{{$indent}}LET distro = alpine3.18

{{$indent}}FROM golang:${go_version}-${distro}
{{$indent}}WORKDIR /go-workdir`

	goDeps = `
{{.Prefix}}deps:
    {{- if .Prefix }}
    FROM +{{.Prefix}}base{{"\n"}}
    {{- end }}
    # These cache dirs will be used in later test and build targets
    # to persist cached go packages.
    #
    # NOTE: cache only gets persisted on successful builds. A test
    # failure will prevent the go cache from being persisted.
    ENV GOCACHE = "` + goCache + `"
    ENV GOMODCACHE = "` + goModCache + `"

    # Copying only go.mod and go.sum means that the cache for this
    # target will only be busted when go.mod/go.sum change. This
    # means that we can cache the results of 'go mod download'.
    COPY go.mod .
    # Projects with no external dependencies do not have a go.sum.
    COPY --if-exists go.sum .
    RUN go mod download`

	goTestBase = `
{{.Prefix}}test-base:
    FROM +{{.Prefix}}deps

    # gcc and g++ are required for -race.
    RUN apk add --update gcc g++

    # This copies the whole project. If you want better caching, try
    # limiting this to _just_ files required by your go tests.
    COPY . .`

	goTestRace = `
# {{.Prefix}}test-race runs 'go test -race'.
{{.Prefix}}test-race:
    FROM +{{.Prefix}}test-base

    CACHE --sharing shared "$GOCACHE"
    CACHE --sharing shared "$GOMODCACHE"

    # package sets the package that tests will run against.
    ARG package = ./...

    RUN go test -race "$package"`

	goTestIntegration = `
# {{.Prefix}}test-integration runs 'go test -tags integration'.
{{.Prefix}}test-integration:
    FROM +{{.Prefix}}test-base

    CACHE --sharing shared "$GOCACHE"
    CACHE --sharing shared "$GOMODCACHE"

    # package sets the package that tests will run against.
    ARG package = ./...

    RUN go test -tags integration "$package"`

	goTest = `
# {{.Prefix}}test runs all go test targets
{{.Prefix}}test:
    BUILD +{{.Prefix}}test-race
    BUILD +{{.Prefix}}test-integration`

	goProjBase = `
{{.Prefix}}proj-base:
    FROM +{{.Prefix}}deps

    # This copies the whole project. If you want better caching, try
    # limiting this to _just_ files required by your go project.
    COPY . .`

	goTidy = `
# {{.Prefix}}mod-tidy runs 'go mod tidy', saving go.mod and go.sum locally.
{{.Prefix}}mod-tidy:
    FROM +{{.Prefix}}proj-base

    RUN go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT --if-exists go.sum AS LOCAL go.sum`

	goBuild = `
# {{.Prefix}}build runs 'go build ./...', saving artifacts locally.
{{.Prefix}}build:
    FROM +{{.Prefix}}proj-base

    CACHE --sharing shared "$GOCACHE"
    CACHE --sharing shared "$GOMODCACHE"

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

// Type returns 'go'
func (g *Golang) Type(context.Context) string {
	return "go"
}

// ForDir returns a Project for the given directory. It returns ErrSkip if the
// directory does not contain a go project.
func (g *Golang) ForDir(ctx context.Context, dir string) (Project, error) {
	_, err := fs.Stat(g.fs, filepath.Join(dir, goMod))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.Wrap(ErrSkip, "no go.mod found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "error reading go.mod")
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

// Targets returns the targets that should be used for this Golang project.
func (g *Golang) Targets(_ context.Context) ([]Target, error) {
	return []Target{
		&targetFormatter{template: goBase},
		&targetFormatter{template: goDeps},
		&targetFormatter{template: goTestBase},
		&targetFormatter{template: goTestRace},
		&targetFormatter{template: goTestIntegration},
		&targetFormatter{template: goTest},
		&targetFormatter{template: goProjBase},
		&targetFormatter{template: goTidy},
		&targetFormatter{template: goBuild},
	}, nil
}

type targetFormatter struct {
	prefix   string
	template string
}

func (f *targetFormatter) SetPrefix(_ context.Context, pfx string) {
	if pfx == "" {
		f.prefix = ""
		return
	}
	if !strings.HasSuffix(pfx, "-") {
		pfx += "-"
	}
	f.prefix = pfx
}

func (f *targetFormatter) Format(_ context.Context, w io.Writer, indent string, level int) error {
	t := strings.TrimSpace(f.template) + "\n"
	tmpl, err := template.New("").Parse(t)
	if err != nil {
		return errors.Wrap(err, "golang: failed to parse target template")
	}
	type tmplCtx struct {
		Prefix string
		Indent string
	}
	err = tmpl.Execute(w, tmplCtx{Prefix: f.prefix, Indent: indent})
	if err != nil {
		return errors.Wrap(err, "golang: failed to execute target template")
	}
	return nil
}
