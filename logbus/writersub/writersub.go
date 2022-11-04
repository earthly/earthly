package writersub

import (
	"io"
	"sync"

	"github.com/earthly/cloud-api/logstream"
)

// WriterSub is a bus subscriber that can print formatted logs to a writer.
type WriterSub struct {
	w              io.Writer
	targetIDFilter string

	mu     sync.Mutex
	errors []error
}

// New creates a new WriterSub.
func New(w io.Writer, targetIDFilter string) *WriterSub {
	return &WriterSub{
		w:              w,
		targetIDFilter: targetIDFilter,
	}
}

// Write writes the given delta to the writer, if it is a formatted log delta.
func (ws *WriterSub) Write(delta *logstream.Delta) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	switch d := delta.DeltaTypeOneof.(type) {
	case *logstream.Delta_DeltaFormattedLog:
		if ws.targetIDFilter != "" && d.DeltaFormattedLog.TargetId != ws.targetIDFilter {
			return
		}
		_, err := ws.w.Write(d.DeltaFormattedLog.Data)
		if err != nil {
			ws.errors = append(ws.errors, err)
			return
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
