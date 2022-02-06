package consumer

import (
	"context"
	"io"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/chanutil"
)

// NewTargetConsumerReader consumes deltas and returns a reader for a specific
// target.
func NewTargetConsumerReader(ctx context.Context, targetID string, deltaCh chan api.Delta) io.Reader {
	filteredCh := chanutil.Filter(ctx, deltaCh, false, false, map[string]bool{targetID: true})
	c := NewConsumer(ctx, filteredCh, nil)
	return c.GetReader(targetID)
}
