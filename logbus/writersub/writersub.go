package writersub

import (
	"io"
	"sync"

	"github.com/earthly/cloud-api/logstream"
)

// WriterSub is a bus subscriber that can print formatted logs to a writer.
type WriterSub struct {
	w io.Writer

	mu     sync.Mutex
	errors []error
}

// New creates a new WriterSub.
func New(w io.Writer) *WriterSub {
	return &WriterSub{
		w: w,
	}
}

// Write writes the given delta to the writer, if it is a formatted log delta.
func (ws *WriterSub) Write(delta *logstream.Delta) {
	switch d := delta.DeltaTypeOneof.(type) {
	case *logstream.Delta_DeltaFormattedLog:
		_, err := ws.w.Write(d.DeltaFormattedLog.Data)
		if err != nil {
			ws.mu.Lock()
			ws.errors = append(ws.errors, err)
			ws.mu.Unlock()
		}
	default:
	}
}

// Errors returns any errors that occurred while writing to the writer.
func (ws *WriterSub) Errors() []error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return ws.errors
}
