package llbutil

import (
	"github.com/moby/buildkit/client/llb"
)

// WithDependency creates a fake dependency between two states.
func WithDependency(state llb.State, depState llb.State, stateStr, depStr string, opts ...llb.RunOption) llb.State {
	// TODO: Is there a better way to mark two states as depending on each other?
	if depState.Output() == nil {
		// depState is Scratch.
		return state
	}

	// Copy a wildcard that could never exist.
	// (And allow for the wildcard to match nothing).
	interm := ScratchWithPlatform()
	interm = interm.File(llb.Copy(
		depState, "/fake-745cb405-fbfb-4ea7-83b0-a85c26b4aff0-*", "/tmp/",
		&llb.CopyInfo{
			CreateDestPath:      true,
			AllowWildcard:       true,
			AllowEmptyWildcard:  true,
			CopyDirContentsOnly: true,
		}), llb.WithCustomNamef("[internal] (fakecopy1) %s depends on %s", stateStr, depStr))

	// Do this again. The extra step is needed to prevent the need for BuildKit
	// to re-hash the input in certain cases (can be slow if depState is large).
	return state.File(llb.Copy(
		depState, "/fake-5fa01e05-ca9e-45c9-8721-05b9183a2914-*", "/tmp/",
		&llb.CopyInfo{
			CreateDestPath:      true,
			AllowWildcard:       true,
			AllowEmptyWildcard:  true,
			CopyDirContentsOnly: true,
		}), llb.WithCustomNamef("[internal] (fakecopy2) %s depends on %s", stateStr, depStr))
}
