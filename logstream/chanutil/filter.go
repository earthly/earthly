package chanutil

import (
	"context"

	"github.com/earthly/earthly/logstream/api"
)

// Map takes a channel and applies a function to each element of the channel.
func Map(ctx context.Context, in chan api.Delta, mapFun func(api.Delta) api.Delta) chan api.Delta {
	out := make(chan api.Delta)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case delta, ok := <-in:
				if !ok {
					close(out)
					return
				}
				delta2 := mapFun(delta)
				if len(delta2.DeltaLogs) > 0 || len(delta2.DeltaManifests) > 0 {
					out <- delta2
				}
			}
		}
	}()
	return out
}

// Filter is a delta channel filter that only permits deltas of a certain
// type.
func Filter(ctx context.Context, in chan api.Delta, manifest bool, allTargets bool, targetsAllow map[string]bool) chan api.Delta {
	filterFun := func(delta api.Delta) api.Delta {
		d2 := api.Delta{
			Version: api.VersionNumber,
		}
		if len(delta.DeltaManifests) > 0 {
			if manifest {
				d2.DeltaManifests = append(d2.DeltaManifests, delta.DeltaManifests...)
			}
		}
		if len(delta.DeltaLogs) > 0 {
			if allTargets {
				d2.DeltaLogs = append(d2.DeltaLogs, delta.DeltaLogs...)
			} else {
				if targetsAllow != nil {
					for _, dl := range delta.DeltaLogs {
						if targetsAllow[dl.TargetID] {
							d2.DeltaLogs = append(d2.DeltaLogs, dl)
						}
					}
				}
			}
		}
		return d2
	}
	return Map(ctx, in, filterFun)
}
