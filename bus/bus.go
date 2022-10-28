package bus

import (
	"context"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
)

const deltaBufferSize = 10240

// Bus is a build log data bus.
// It listens for events from BuildKit and forwards them to the console and
// possibly also to a remote logstream service.
type Bus struct {
	ch chan *logstream.Delta
	bp *BuildPrinter
	sm *SolverMonitor

	mu      sync.Mutex
	history []*logstream.Delta
	subs    []chan *logstream.Delta
	closed  bool
}

// New creates a new Bus.
func New(ctx context.Context) *Bus {
	b := &Bus{
		ch: make(chan *logstream.Delta, deltaBufferSize),
		bp: nil, // set below
		sm: nil, // set below
	}
	b.bp = newBuildPrinter(b)
	b.sm = newSolverMonitor(b)
	go func() {
		<-ctx.Done()
		b.Close()
	}()
	go b.messageLoop(ctx)
	return b
}

// AddSubscriber adds a subscriber to the bus.
func (b *Bus) AddSubscriber() chan *logstream.Delta {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan *logstream.Delta, deltaBufferSize)
	b.subs = append(b.subs, ch)
	// Send history to the new subscriber. The lock ensures
	// that no new messages are added while we are sending.
	for _, delta := range b.history {
		ch <- delta
	}
	return ch
}

// Printer returns the underlying BuildPrinter.
func (b *Bus) Printer() *BuildPrinter {
	return b.bp
}

// SolverMonitor returns the underlying SolverMonitor.
func (b *Bus) SolverMonitor() *SolverMonitor {
	return b.sm
}

// RawDelta sends a raw delta on the bus.
func (b *Bus) RawDelta(delta *logstream.Delta) {
	b.ch <- delta
}

func (b *Bus) messageLoop(ctx context.Context) {
	returnedCh := make(chan struct{})
	defer close(returnedCh)
	go func() {
		<-ctx.Done()
		// Delay closing to allow any pending messages
		// to be processed.
		select {
		case <-returnedCh:
		case <-time.After(5 * time.Second):
		}
		b.Close()
	}()
	for delta := range b.ch {
		// if delta.GetDeltaManifest() != nil {
		// 	fmt.Printf("@#@# Delta manifest: %v\n", delta.GetDeltaManifest())
		// }
		// if delta.GetDeltaLog() != nil {
		// 	fmt.Printf("@#@#@# Delta log %s\n", string(delta.GetDeltaLog().GetLog()))
		// }
		b.mu.Lock()
		b.history = append(b.history, delta)
		var subs = append([]chan *logstream.Delta{}, b.subs...)
		b.mu.Unlock()
		for _, sub := range subs {
			sub <- delta
		}
	}
	b.mu.Lock()
	for _, sub := range b.subs {
		close(sub)
	}
	b.mu.Unlock()
}

// Close closes the bus.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	close(b.ch)
}
