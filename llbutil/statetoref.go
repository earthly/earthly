package llbutil

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// StateToRef takes an LLB state, solves it using gateway and returns the ref.
func StateToRef(ctx context.Context, gwClient gwclient.Client, state llb.State, cacheImports []string) (gwclient.Reference, error) {
	var coes []gwclient.CacheOptionsEntry
	for _, ci := range cacheImports {
		coe := gwclient.CacheOptionsEntry{
			Type:  "registry",
			Attrs: map[string]string{"ref": ci},
		}
		coes = append(coes, coe)
	}
	def, err := state.Marshal(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "marshal state")
	}
	r, err := gwClient.Solve(ctx, gwclient.SolveRequest{
		Definition:   def.ToPB(),
		CacheImports: coes,
	})
	if err != nil {
		return nil, errors.Wrap(err, "solve state")
	}
	ref, err := r.SingleRef()
	if err != nil {
		return nil, errors.Wrap(err, "single ref")
	}
	return ref, nil
}
