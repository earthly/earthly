package logstreamer

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/earthly/cloud-api/logstream"
	"github.com/hashicorp/go-multierror"
)

type LogstreamOrchestrator struct {
	buildID         string
	bus             LogBus
	c               CloudClient
	initialManifest *logstream.RunManifest

	doneMU sync.Mutex
	mu     sync.Mutex
	doneCH chan struct{}

	errMu       sync.Mutex
	errors      []error
	started     atomic.Bool
	retries     int
	streamer    *LogStreamer
	deltas      *deltasIter
	verbose     bool
	deltaBuffer int
}

type LOpt func(*LogstreamOrchestrator) *LogstreamOrchestrator

// WithVerbose sets the verbose option on the LogstreamOrchestrator
func WithVerbose(verbose bool) func(orchestrator *LogstreamOrchestrator) *LogstreamOrchestrator {
	return func(orchestrator *LogstreamOrchestrator) *LogstreamOrchestrator {
		orchestrator.verbose = verbose
		return orchestrator
	}
}

func WithDeltaBuffer(buffer int) func(orchestrator *LogstreamOrchestrator) *LogstreamOrchestrator {
	return func(orchestrator *LogstreamOrchestrator) *LogstreamOrchestrator {
		orchestrator.deltaBuffer = buffer
		return orchestrator
	}
}

func NewLogstreamOrchestrator(bus LogBus, c CloudClient, initialManifest *logstream.RunManifest, opts ...LOpt) *LogstreamOrchestrator {
	ls := &LogstreamOrchestrator{
		buildID:         initialManifest.GetBuildId(),
		bus:             bus,
		c:               c,
		initialManifest: initialManifest,
		retries:         10,
		deltaBuffer:     DefaultBufferSize,
	}
	for _, o := range opts {
		ls = o(ls)
	}
	return ls
}

// StartLogstreamer will start streaming to the cloud retrying up the retry count
// Callers should listen to Done to be notified when the streaming contract completes
// StartLogstreamer may only be called once
func (l *LogstreamOrchestrator) StartLogstreamer(ctx context.Context) {
	if l.started.Swap(true) {
		// Can only start once
		return
	}
	go func() {
		for i := 0; i < l.retries; i++ {
			l.start()
			l.CloseLastLogstreamer()

			l.mu.Lock()
			l.deltas = newDeltasIter(l.deltaBuffer, l.initialManifest, l.verbose)
			l.streamer = New(l.c, l.buildID, l.deltas)
			l.mu.Unlock()

			go l.bus.AddSubscriber(l.deltas)
			shouldRetry, err := l.streamer.Stream(ctx)
			l.addError(err)
			l.markDone()
			if !shouldRetry {
				return
			}
		}
	}()
}

// CloseLastLogstreamer Will close any previous logstreamer / deltas.
func (l *LogstreamOrchestrator) CloseLastLogstreamer() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.deltas != nil {
		l.bus.RemoveSubscriber(l.deltas)
		l.deltas.close()
	}
	if l.streamer != nil {
		l.streamer.Close()
	}
}

func (l *LogstreamOrchestrator) addError(err error) {
	if err == nil {
		return
	}
	l.errMu.Lock()
	defer l.errMu.Unlock()
	l.errors = append(l.errors, err)
}

func (l *LogstreamOrchestrator) getError() error {
	l.errMu.Lock()
	defer l.errMu.Unlock()
	var retErr error
	for _, err := range l.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}

// Close will mark the deltas and streamer as closed and prevent further writes to the cloud
// Callers should listen to Done() to know when it is safe to exit.
func (l *LogstreamOrchestrator) Close() (int32, int32, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var logsWritten int32
	var manifestsWritten int32
	if l.deltas != nil {
		l.bus.RemoveSubscriber(l.deltas)
		manifestsWritten, logsWritten = l.deltas.close()
	}
	if l.streamer != nil {
		l.streamer.Close()
		return manifestsWritten, logsWritten, l.getError()
	}
	return manifestsWritten, logsWritten, nil
}

func (l *LogstreamOrchestrator) WriteToDeltaIter(delta *logstream.Delta) {
	l.mu.Lock()
	d := l.deltas
	defer l.mu.Unlock()
	d.Write(delta)
}

// Done returns a channel that is closed once the Logstreamer has finished
// communicating wit the server.
// Callers should listen to the Done channel before exiting
func (l *LogstreamOrchestrator) Done() chan struct{} {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	return l.doneCH
}

func (l *LogstreamOrchestrator) markDone() {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	close(l.doneCH)
}

func (l *LogstreamOrchestrator) start() {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	l.doneCH = make(chan struct{})
}
