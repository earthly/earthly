package logstreamer

import (
	"context"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	"github.com/hashicorp/go-multierror"
)

type Orchestrator struct {
	buildID         string
	bus             LogBus
	c               CloudClient
	initialManifest *logstream.RunManifest

	doneMU sync.Mutex
	mu     sync.Mutex
	errMu  sync.Mutex

	doneCH chan struct{}
	subCH  chan struct{}

	errors      []error
	startOnce   sync.Once
	retries     int
	streamer    *LogStreamer
	deltas      *deltasIter
	verbose     bool
	deltaBuffer int
}

type LOpt func(*Orchestrator) *Orchestrator

// WithVerbose sets the verbose option on the Orchestrator
func WithVerbose(verbose bool) func(orchestrator *Orchestrator) *Orchestrator {
	return func(orchestrator *Orchestrator) *Orchestrator {
		orchestrator.verbose = verbose
		return orchestrator
	}
}

func WithDeltaBuffer(buffer int) func(orchestrator *Orchestrator) *Orchestrator {
	return func(orchestrator *Orchestrator) *Orchestrator {
		orchestrator.deltaBuffer = buffer
		return orchestrator
	}
}

func NewOrchestrator(bus LogBus, c CloudClient, initialManifest *logstream.RunManifest, opts ...LOpt) *Orchestrator {
	ls := &Orchestrator{
		buildID:         initialManifest.GetBuildId(),
		bus:             bus,
		c:               c,
		initialManifest: initialManifest,
		retries:         10,
		deltaBuffer:     DefaultBufferSize,
		doneCH:          make(chan struct{}),
		subCH:           nil, // nil on purpose - only use when subscribing
	}
	for _, o := range opts {
		ls = o(ls)
	}
	return ls
}

// Start will restart streaming to the cloud retrying up the retry count
// Callers should listen to Done to be notified when the streaming contract completes
// Start may only be called once
func (l *Orchestrator) Start(ctx context.Context) {
	l.startOnce.Do(func() {
		go func() {
			for i := 0; i < l.retries; i++ {
				l.subscribe()
				shouldRetry, err := l.streamer.Stream(ctx)
				l.addError(err)
				l.markDone()
				if !shouldRetry {
					return
				}
				l.restart()
			}
		}()
	})
}

func (l *Orchestrator) addError(err error) {
	if err == nil {
		return
	}
	l.errMu.Lock()
	defer l.errMu.Unlock()
	l.errors = append(l.errors, err)
}

func (l *Orchestrator) getError() error {
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
func (l *Orchestrator) Close() (manifestsWritten int32, logsWritten int32, _ error) {
	l.mu.Lock()
	defer l.mu.Unlock()

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

func (l *Orchestrator) WriteToDeltaIter(delta *logstream.Delta) {
	l.mu.Lock()
	defer l.mu.Unlock()
	d := l.deltas
	d.Write(delta)
}

// Done returns a channel that is closed once the Logstreamer has finished
// communicating with the server.
// Callers should listen to the Done channel before exiting
func (l *Orchestrator) Done() chan struct{} {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	return l.doneCH
}

func (l *Orchestrator) markDone() {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	close(l.doneCH)
}

func (l *Orchestrator) restart() {
	l.doneMU.Lock()
	defer l.doneMU.Unlock()
	l.doneCH = make(chan struct{})
}

func (l *Orchestrator) subscribe() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.deltas != nil {
		// wait for previous subscriber to finish being added (if there is one)
		if l.subCH != nil {
			<-l.subCH
		}
		l.bus.RemoveSubscriber(l.deltas)
		l.deltas.close()
	}
	if l.streamer != nil {
		l.streamer.Close()
	}
	l.deltas = newDeltasIter(l.deltaBuffer, l.initialManifest, l.verbose)
	l.streamer = New(l.c, l.buildID, l.deltas)
	l.subCH = make(chan struct{})
	go func(subCH chan struct{}) {
		defer close(subCH)
		l.bus.AddSubscriber(l.deltas)
	}(l.subCH)
}
