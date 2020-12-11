package llbutil

import (
	"github.com/moby/buildkit/client/llb"
)

const fakeDepImg = "busybox:1.31.1"

// WithDependency creates a fake dependency between two states.
func WithDependency(state llb.State, depState llb.State, stateStr, depStr string, opts ...llb.RunOption) llb.State {
	// TODO: Is there a better way to mark two states as depending on each other?
	if depState.Output() == nil {
		// depState is Scratch.
		return state
	}

	// Copy a file that is known to exist in (almost) all images.
	// (The copy is necessary to prevent situation where BuildKit needs to re-hash the entire
	// depState).
	interm := llb.Scratch().Platform(TargetPlatform)
	interm = CopyOp(
		depState, []string{"/bin/sh"}, interm, "/bin/sh", false, true, "",
		llb.WithCustomNamef("[internal] (copy) %s depends on %s", stateStr, depStr))

	// Execute a command which doesn't do anything (but it creates a new layer, which casues it
	// to depend on depState).
	runOpts := []llb.RunOption{
		llb.Args([]string{"/bin/sh", "-c", "true"}),
		llb.Dir("/"),
		llb.ReadonlyRootFS(),
		llb.AddMount("/fake-dep", interm, llb.Readonly),
		llb.WithCustomNamef("[internal] (run) %s depends on %s", stateStr, depStr),
	}
	runOpts = append(runOpts, opts...)
	opImg := llb.Image(
		fakeDepImg, llb.MarkImageInternal, llb.Platform(TargetPlatform),
		llb.WithCustomNamef("[internal] helper image for fake dep operations"))
	return opImg.Run(runOpts...).AddMount("/fake", state)
}
