package chanutil

import (
	"context"
	"time"

	"github.com/earthly/earthly/logstream/api"
)

// IntervalBatcher consumes deltas from an input channel, merges them over time,
// and sends the merged deltas to an output channel at regular time intervals.
func IntervalBatcher(ctx context.Context, interval time.Duration, in chan api.Delta) chan api.Delta {
	out := make(chan api.Delta)
	go func() {
		batch := api.Delta{
			Version: api.VersionNumber,
		}
		var hasData bool
		ticker := time.NewTicker(interval)
		nextManifestOrderID := int64(0)
		for {
			select {
			case <-ctx.Done():
				if hasData {
					batch, _ = api.SimplifyDeltas(batch, nextManifestOrderID)
					out <- batch
				}
				close(out)
				ticker.Stop()
				return
			case <-ticker.C:
				if hasData {
					batch, nextManifestOrderID = api.SimplifyDeltas(batch, nextManifestOrderID)
					out <- batch
					batch = api.Delta{
						Version: api.VersionNumber,
					}
					hasData = false
				}
				ticker.Stop()
			case delta, ok := <-in:
				if !ok {
					if hasData {
						batch, _ = api.SimplifyDeltas(batch, nextManifestOrderID)
						out <- batch
					}
					close(out)
					ticker.Stop()
					return
				}
				if !hasData {
					ticker.Reset(interval)
				}
				hasData = true
				batch.DeltaManifests = append(batch.DeltaManifests, delta.DeltaManifests...)
				batch.DeltaLogs = append(batch.DeltaLogs, delta.DeltaLogs...)
			}
		}
	}()
	return out
}
