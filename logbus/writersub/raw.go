package writersub

import (
	"io"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// RawWriterSub is a bus subscriber that can print formatted logs to a writer.
type RawWriterSub struct {
	w    io.Writer
	json bool

	mu     sync.Mutex
	errors []error
}

// New creates a new WriterSub.
func NewRaw(w io.Writer, json bool) *RawWriterSub {
	return &RawWriterSub{
		w:    w,
		json: json,
	}
}

// Write writes the given delta to the writer, if it is a formatted log delta.
func (rws *RawWriterSub) Write(delta *logstream.Delta) {
	rws.mu.Lock()
	defer rws.mu.Unlock()
	var dt []byte
	var err error
	if rws.json {
		dt, err = protojson.Marshal(delta)
		dt = append(dt, '\n')
	} else {
		dt, err = proto.Marshal(delta)
	}
	if err != nil {
		rws.errors = append(rws.errors, err)
		return
	}
	_, err = rws.w.Write(dt)
	if err != nil {
		rws.errors = append(rws.errors, err)
		return
	}
}

// Errors returns any errors that occurred while writing to the writer.
func (rws *RawWriterSub) Errors() []error {
	rws.mu.Lock()
	defer rws.mu.Unlock()
	return rws.errors
}
