package llbutil

import (
	"github.com/moby/buildkit/client/llb"
)

const fakeDepImg = "busybox:1.31.1"

// WithDependency creates a fake dependency between two states.
func WithDependency(state llb.State, depState llb.State, opts ...llb.RunOption) llb.State {
	// TODO: Is there a better way to mark two states as depending on each other?
	if depState.Output() == nil {
		// depState is Scratch.
		return state
	}
	runOpts := []llb.RunOption{
		llb.Args([]string{"/bin/sh", "-c", "true"}),
		llb.Dir("/"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/fake-dep", depState, llb.Readonly),
	}
	runOpts = append(runOpts, opts...)
	opImg := llb.Image(
		fakeDepImg, llb.MarkImageInternal, llb.Platform(TargetPlatform),
		llb.ResolveModePreferLocal,
		llb.WithCustomNamef("[internal] mark dependency"))
	return opImg.Run(runOpts...).AddMount("/fake", state)
}
