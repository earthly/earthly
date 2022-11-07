package llbutil

import (
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/moby/buildkit/client/llb"
)

// WithDependency creates a fake dependency between two states.
func WithDependency(state pllb.State, depState pllb.State, stateStr, depStr string, platr *platutil.Resolver, opts ...llb.RunOption) pllb.State {
	// TODO: Is there a better way to mark two states as depending on each other?
	vm := &vertexmeta.VertexMeta{
		Internal: true,
	}
	if state.Output() == nil {
		// TODO: A buildkit bug causes a hang if we try to do this operation
		// on top of a scratch state, and then try to export this as remote
		// cache.
		state = state.File(pllb.Mkdir("/tmp", 0755),
			llb.WithCustomNamef("%s(fakecopy) init scratch", vm.ToVertexPrefix()))
	}
	if depState.Output() == nil {
		// depState is Scratch.
		return state
	}

	// Copy a wildcard that could never exist.
	// (And allow for the wildcard to match nothing).
	interm := platr.Scratch()
	interm = interm.File(pllb.Copy(
		depState, "/fake-745cb405-fbfb-4ea7-83b0-a85c26b4aff0-*", "/tmp/",
		&llb.CopyInfo{
			CreateDestPath:      true,
			AllowWildcard:       true,
			AllowEmptyWildcard:  true,
			CopyDirContentsOnly: true,
		}),
		llb.WithCustomNamef(
			"%s(fakecopy1) %s depends on %s",
			vm.ToVertexPrefix(), stateStr, depStr),
	)

	// Do this again. The extra step is needed to prevent the need for BuildKit
	// to re-hash the input in certain cases (can be slow if depState is large).
	return state.File(pllb.Copy(
		interm, "/fake-5fa01e05-ca9e-45c9-8721-05b9183a2914-*", "/tmp/",
		&llb.CopyInfo{
			CreateDestPath:      true,
			AllowWildcard:       true,
			AllowEmptyWildcard:  true,
			CopyDirContentsOnly: true,
		}),
		llb.WithCustomNamef(
			"%s(fakecopy2) %s depends on %s",
			vm.ToVertexPrefix(), stateStr, depStr),
	)
}
