package consumer

import (
	"context"
	"io"
	"sync"

	"github.com/earthly/earthly/logstream/api"
	"github.com/pkg/errors"
)

type logReader struct {
	targetID string
	buffer   chan []byte
	c        *Consumer
}

func (lr *logReader) Read(p []byte) (int, error) {
	data, ok := <-lr.buffer
	if !ok {
		lr.c.mu.Lock()
		defer lr.c.mu.Unlock()
		if lr.c.err != nil {
			return 0, lr.c.err
		}
		return 0, io.EOF
	}
	n := copy(p, data)
	return n, nil
}

// Consumer is a consumer of deltas from a delta channel.
type Consumer struct {
	deltaCh           chan api.Delta
	doneCh            chan struct{}
	onManifestChanges func(api.Manifest, api.Delta)

	mu       sync.Mutex
	manifest api.Manifest
	readers  map[string]*logReader
	err      error
}

// NewConsumer creates a new Consumer.
func NewConsumer(ctx context.Context, deltaCh chan api.Delta, onManifestChanges func(api.Manifest, api.Delta)) *Consumer {
	c := &Consumer{
		manifest: api.Manifest{
			Version: api.VersionNumber,
		},
		readers:           make(map[string]*logReader),
		deltaCh:           deltaCh,
		doneCh:            make(chan struct{}),
		onManifestChanges: onManifestChanges,
	}
	go func() {
		for {
			select {
			case delta, ok := <-deltaCh:
				if !ok {
					for _, r := range c.readers {
						close(r.buffer)
					}
					close(c.doneCh)
					return
				}
				err := c.processDelta(delta)
				if err != nil {
					c.mu.Lock()
					c.err = err
					c.mu.Unlock()
					for _, r := range c.readers {
						close(r.buffer)
					}
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return c
}

// GetReader returns a log reader for a given targetID.
func (c *Consumer) GetReader(targetID string) io.Reader {
	return c.getReader(targetID)
}

// GetManifest returns the current manifest.
func (c *Consumer) GetManifest() api.Manifest {
	c.mu.Lock()
	defer c.mu.Unlock()
	ret := api.Manifest{}
	ret.CopyFrom(&c.manifest)
	return ret
}

func (c *Consumer) getReader(targetID string) *logReader {
	c.mu.Lock()
	defer c.mu.Unlock()
	r, found := c.readers[targetID]
	if !found {
		r = &logReader{
			c:        c,
			targetID: targetID,
			buffer:   make(chan []byte, 1000),
		}
		c.readers[targetID] = r
	}
	return r
}

func (c *Consumer) processDelta(delta api.Delta) error {
	if delta.Version != api.VersionNumber {
		return errors.Errorf("unsupported delta version: %d", delta.Version)
	}
	if len(delta.DeltaManifests) > 0 {
		c.mu.Lock()
		for _, dm := range delta.DeltaManifests {
			err := c.manifest.ApplyDelta(dm)
			if err != nil {
				c.mu.Unlock()
				return err
			}
		}
		m := api.Manifest{}
		err := m.CopyFrom(&c.manifest)
		if err != nil {
			c.mu.Unlock()
			return err
		}
		c.mu.Unlock()
		if c.onManifestChanges != nil {
			c.onManifestChanges(m, delta)
		}
	}
	for _, dl := range delta.DeltaLogs {
		lr := c.getReader(dl.TargetID)
		lr.buffer <- dl.Data
	}
	return nil
}

// Done return a channel that is closed when the consumer is finished.
func (c *Consumer) Done() chan struct{} {
	return c.doneCh
}
