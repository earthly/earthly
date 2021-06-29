// Package pllb is a wrapper around llb, which makes it compatible with concurrent
// code. The standard BuildKit llb package does not allow llb.State to be used
// from different goroutines.
package pllb

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// gmu is a global lock used for any interaction with the llb package.
var gmu sync.Mutex

// hacks are used to support changing opts before copying
var hackName map[string]string
var hackOpts map[string][]llb.LocalOption

func init() {
	hackName = map[string]string{}
	hackOpts = map[string][]llb.LocalOption{}
}

// State is a wrapper around llb.State.
type State struct {
	st llb.State
	id string
}

// FromRawState creates a wrapper around a raw llb.State.
func FromRawState(st llb.State) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: st}
}

// Scratch is a wrapper around llb.Scratch.
func Scratch() State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: llb.Scratch()}
}

// Local is a wrapper around llb.Local.
func Local(name string, opts ...llb.LocalOption) State {
	gmu.Lock()
	defer gmu.Unlock()
	s := State{st: llb.Local(name, opts...)}
	id := uuid.NewString()
	hackName[id] = name
	hackOpts[id] = opts
	s.id = id
	return s
}

// Image is a wrapper around llb.Image.
func Image(ref string, opts ...llb.ImageOption) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: llb.Image(ref, opts...)}
}

// Git is a wrapper around llb.Git.
func Git(remote, ref string, opts ...llb.GitOption) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: llb.Git(remote, ref, opts...)}
}

// RawState returns the wrapped llb.State, but requires an unlock from the caller.
func (s State) RawState() (llb.State, func()) {
	gmu.Lock()
	return s.st, gmu.Unlock
}

// UnsafeUnwrap returns the underlying llb.State without locking.
func (s State) UnsafeUnwrap() llb.CopyInput {
	return s.st
}

// Output is a wrapper around llb.Output.
func (s State) Output() llb.Output {
	gmu.Lock()
	defer gmu.Unlock()
	return s.st.Output()
}

// SetMarshalDefaults is a wrapper around llb.SetMarshalDefaults.
func (s State) SetMarshalDefaults(co ...llb.ConstraintsOpt) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.SetMarshalDefaults(co...)}
}

// Marshal is a wrapper around llb.Marshal.
func (s State) Marshal(ctx context.Context, co ...llb.ConstraintsOpt) (*llb.Definition, error) {
	gmu.Lock()
	defer gmu.Unlock()
	return s.st.Marshal(ctx, co...)
}

// Run is a wrapper around llb.Run.
func (s State) Run(ro ...llb.RunOption) ExecState {
	gmu.Lock()
	defer gmu.Unlock()
	return ExecState{est: s.st.Run(ro...)}
}

// File is a wrapper around llb.File.
func (s State) File(a *FileAction, opts ...llb.ConstraintsOpt) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.File(a.fia, opts...)}
}

// AddEnv is a wrapper around llb.AddEnv.
func (s State) AddEnv(key, value string) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.AddEnv(key, value)}
}

// Dir is a wrapper around llb.Dir.
func (s State) Dir(str string) State {
	gmu.Lock()
	defer gmu.Unlock()
	ss := State{st: s.st.Dir(str)}
	id := uuid.NewString()
	hackName[id] = str
	hackOpts[id] = nil
	ss.id = id
	return ss
}

// GetDir is a wrapper around llb.GetDir.
func (s State) GetDir(ctx context.Context) (string, error) {
	gmu.Lock()
	defer gmu.Unlock()
	return s.st.GetDir(ctx)
}

// WithInclude needs work
func (s State) WithInclude(incl []string) State {
	gmu.Lock()
	defer gmu.Unlock()

	if s.id == "" {
		fmt.Fprintf(os.Stderr, "internal warning: state uuid was never set for %v\n", s)
		return s
	}

	name, ok := hackName[s.id]
	if !ok {
		panic(fmt.Sprintf("internal error: hackName has no entry for %q", name))
	}
	opts, ok := hackOpts[s.id]
	if !ok {
		panic(fmt.Sprintf("internal error: hackOpts has no entry for %q", name))
	}

	newOpts := []llb.LocalOption{}

	for _, o := range opts {
		newOpts = append(newOpts, o)
	}
	newOpts = append(newOpts, llb.IncludePatterns(incl))

	return State{st: llb.Local(name, newOpts...), id: s.id}
}

// User is a wrapper around llb.User.
func (s State) User(v string) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.User(v)}
}

// Platform is a wrapper around llb.Platform.
func (s State) Platform(p specs.Platform) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.Platform(p)}
}

// ExecState is a wrapper around llb.ExecState.
type ExecState struct {
	est llb.ExecState
}

// AddMount is a wrapper around llb.AddMount.
func (e ExecState) AddMount(target string, source State, opt ...llb.MountOption) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: e.est.AddMount(target, source.st, opt...)}
}

// Root is a wrapper around llb.Root.
func (e ExecState) Root() State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: e.est.Root()}
}

// AddMount is a wrapper around llb.AddMount.
func AddMount(dest string, mountState State, opts ...llb.MountOption) llb.RunOption {
	gmu.Lock()
	defer gmu.Unlock()
	return llb.AddMount(dest, mountState.st, opts...)
}

// FileAction is a wrapper around llb.FileAction.
type FileAction struct {
	fia *llb.FileAction
}

// Mkdir is a wrapper around llb.Mkdir.
func (fa *FileAction) Mkdir(p string, m os.FileMode, opt ...llb.MkdirOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: fa.fia.Mkdir(p, m, opt...)}
}

// Mkfile is a wrapper around llb.Mkfile.
func (fa *FileAction) Mkfile(p string, m os.FileMode, dt []byte, opt ...llb.MkfileOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: fa.fia.Mkfile(p, m, dt, opt...)}
}

// Rm is a wrapper around llb.Rm.
func (fa *FileAction) Rm(p string, opt ...llb.RmOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: fa.fia.Rm(p, opt...)}
}

// Copy is a wrapper around llb.Copy.
func (fa *FileAction) Copy(input CopyInput, src, dest string, opt ...llb.CopyOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: fa.fia.Copy(input.UnsafeUnwrap(), src, dest, opt...)}
}

// Mkdir is a wrapper around llb.Mkdir.
func Mkdir(p string, m os.FileMode, opt ...llb.MkdirOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: llb.Mkdir(p, m, opt...)}
}

// Mkfile is a wrapper around llb.Mkfile.
func Mkfile(p string, m os.FileMode, dt []byte, opts ...llb.MkfileOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: llb.Mkfile(p, m, dt, opts...)}
}

// Copy is a wrapper around llb.Copy.
func Copy(input CopyInput, src, dest string, opts ...llb.CopyOption) *FileAction {
	gmu.Lock()
	defer gmu.Unlock()
	return &FileAction{fia: llb.Copy(input.UnsafeUnwrap(), src, dest, opts...)}
}

// CopyInput is a mirror of llb.CopyInput.
type CopyInput interface {
	UnsafeUnwrap() llb.CopyInput
}
