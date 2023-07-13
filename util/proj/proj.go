//go:generate hel --output helheim_mocks_test.go

// Package proj contains types and functions for managing a project's
// Earthfile(s).
package proj

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// FS represents the type that proj types need to inspect files in the
// filesystem in order to detect what should be in the generated Earthfile(s).
type FS interface {
	Open(name string) (fs.File, error)
	Stat(name string) (fs.FileInfo, error)
}

// StdFS returns a standard library FS for use in most use cases of this
// package.
func StdFS() FS {
	return os.DirFS(".").(FS)
}

// Cmd represents the type that proj types will use to execute commands.
type Cmd interface {
	// Run is used to run the command. If the command needs to be canceled (or
	// time out), then ctx should be used to do so.
	Run(ctx context.Context) (stdout, stderr io.Reader, _ error)
}

// Execer represents a type that can return commands that may be executed in the
// OS.
type Execer interface {
	Command(name string, args ...string) Cmd
}

type stdCmd struct {
	name string
	args []string
}

func (c stdCmd) Run(ctx context.Context) (stdout, stderr io.Reader, _ error) {
	cmd := exec.CommandContext(ctx, c.name, c.args...)
	var outw, errw bytes.Buffer
	cmd.Stdout = &outw
	cmd.Stderr = &errw
	if err := cmd.Run(); err != nil {
		return &outw, &errw, err
	}
	return &outw, &errw, nil
}

type stdExecer struct{}

func (e stdExecer) Command(name string, args ...string) Cmd {
	return stdCmd{name: name, args: args}
}

// StdExecer returns a standard library Execer (mostly wrapping the exec
// package) for use in most use cases of this package.
func StdExecer() Execer {
	return stdExecer{}
}

// ErrSkip is an error that means that the project should skip this generator.
var ErrSkip = errors.New("proj: this project is not a supported type")

// Formatter is a type that can return formatted code with a given indent string
// and indentation level
type Formatter interface {
	Format(ctx context.Context, w io.Writer, indent string, level int) error
}

// Project is a type that can generate Earthfile code for a given project.
type Project interface {
	ForDir(ctx context.Context, dir string) (Project, error)
	Root(context.Context) string
	Targets(context.Context) ([]Formatter, error)
}

// All returns all available project types for the given dir.
func All(ctx context.Context, dir string) ([]Project, error) {
	known := []Project{
		NewGolang(StdFS(), StdExecer()),
	}
	var active []Project
	for _, proj := range known {
		forDir, err := proj.ForDir(ctx, dir)
		if errors.Is(err, ErrSkip) {
			continue
		}
		if err != nil {
			return nil, errors.Wrapf(err, "checking for project type %T failed", proj)
		}
		active = append(active, forDir)
	}
	return active, nil
}
