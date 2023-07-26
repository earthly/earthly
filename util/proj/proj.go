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

type stdFS struct{}

func (stdFS) Open(path string) (fs.File, error) {
	return os.Open(path)
}

func (stdFS) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

// StdFS returns a standard library FS for use in most use cases of this
// package.
func StdFS() FS {
	return stdFS{}
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

// Target is a type that can write a formatted target with a given indent string
// and indentation level
type Target interface {
	// SetPrefix sets a prefix to prepend to this target's name.
	SetPrefix(context.Context, string)

	// Format writes out the target with the given indentation string and level.
	Format(ctx context.Context, w io.Writer, indent string, level int) error
}

// ProjectType represents a type of project (typically a language).
type ProjectType interface {
	// ForDir returns a Project for the named directory. It should return
	// ErrSkip if there is no project matching this ProjectType at the requested
	// dir.
	ForDir(ctx context.Context, dir string) (Project, error)
}

// Project is a type that can generate Earthfile code for a given project.
type Project interface {
	// Root returns the root directory for this Project.
	Root(context.Context) string

	// Type returns a unique name for this project type. It will be used for
	// conflict avoidance (i.e. making sure we don't have two go modules loaded)
	// and as a prefix for targets in multi-project-type Earthfiles.
	Type(context.Context) string

	// Targets returns a list of targets for this Project.
	Targets(ctx context.Context) ([]Target, error)
}

// All returns all available project types for the given dir.
func All(ctx context.Context, dir string) ([]Project, error) {
	known := []ProjectType{
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
