package chanutil

import (
	"context"

	"github.com/earthly/earthly/logstream/api"
)

// Splitter is a channel splitter that takes an input channel and duplicates
// any incoming messages to nDup number of output channels. The outputs are closed
// if the context is cancelled.
func Splitter(ctx context.Context, in chan api.Delta, nDup int, outSize int) []chan api.Delta {
	outs := make([]chan api.Delta, nDup)
	for i := 0; i < nDup; i++ {
		outs[i] = make(chan api.Delta, outSize)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, out := range outs {
					close(out)
				}
			case delta, ok := <-in:
				if !ok {
					for _, out := range outs {
						close(out)
					}
					return
				}
				for _, out := range outs {
					out <- delta
				}
			}
		}
	}()
	return outs
}
