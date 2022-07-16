package llbutil

import (
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/moby/buildkit/client/llb"
)

// FlattenTimestamp flattens the timestamps of all new files created or
// modified in upper since lower.
func FlattenTimestamp(lower pllb.State, upper pllb.State, opts ...llb.ConstraintsOpt) pllb.State {
	if upper.Output() == nil {
		// Quick way out for scratch.
		return upper
	}
	diff := pllb.Diff(lower, upper, opts...)
	t := defaultTs()
	fa := pllb.Copy(diff, "/", "/", llb.WithCreatedTime(*t))
	return pllb.Merge([]pllb.State{upper, pllb.Scratch().File(fa, opts...)}, opts...)
}
