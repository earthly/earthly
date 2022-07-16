// Package pllb is a wrapper around llb, which makes it compatible with concurrent
// code. The standard BuildKit llb package does not allow llb.State to be used
// from different goroutines.
package pllb

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// gmu is a global lock used for any interaction with the llb package.
var gmu sync.Mutex

// State is a wrapper around llb.State.
type State struct {
	st llb.State
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
	return State{st: llb.Local(name, opts...)}
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

// Merge is a wrapper around llb.Merge.
func Merge(sts []State, opts ...llb.ConstraintsOpt) State {
	sts2 := make([]llb.State, len(sts))
	for i, st := range sts {
		sts2[i] = st.st
	}
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: llb.Merge(sts2, opts...)}
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
	return State{st: s.st.Dir(str)}
}

// GetDir is a wrapper around llb.GetDir.
func (s State) GetDir(ctx context.Context) (string, error) {
	gmu.Lock()
	defer gmu.Unlock()
	return s.st.GetDir(ctx)
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

// AddExtraHost is a wrapper around llb.AddExtraHost
func (s State) AddExtraHost(hostname string, ip net.IP) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: s.st.AddExtraHost(hostname, ip)}
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

// Diff is a wrapper around llb.Diff.
func Diff(lower State, upper State, opts ...llb.ConstraintsOpt) State {
	gmu.Lock()
	defer gmu.Unlock()
	return State{st: llb.Diff(lower.st, upper.st, opts...)}
}
