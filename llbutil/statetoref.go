package llbutil

import (
	"context"
	"sort"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// StateToRef takes an LLB state, solves it using gateway and returns the ref.
func StateToRef(ctx context.Context, gwClient gwclient.Client, state llb.State, cacheImports map[string]bool) (gwclient.Reference, error) {
	cacheImportsSlice := make([]string, 0, len(cacheImports))
	for ci := range cacheImports {
		cacheImportsSlice = append(cacheImportsSlice, ci)
	}
	sort.Strings(cacheImportsSlice)
	var coes []gwclient.CacheOptionsEntry
	for _, ci := range cacheImportsSlice {
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
