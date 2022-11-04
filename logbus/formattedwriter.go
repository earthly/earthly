package logbus

import (
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// FormattedWriter is a writer that produces DeltaFormattedLog messages.
type FormattedWriter struct {
	bus      *Bus
	targetID string
}

// NewFormattedWriter creates a new FormattedWriter.
func NewFormattedWriter(bus *Bus, targetID string) *FormattedWriter {
	return &FormattedWriter{
		bus:      bus,
		targetID: targetID,
	}
}

// Write writes the given bytes to the writer.
func (fw *FormattedWriter) Write(dt []byte) (int, error) {
	fw.bus.WriteFormattedLog(&logstream.DeltaFormattedLog{
		TargetId: fw.targetID,
		// TODO (vladaionescu): Can the timestamp be passed along straight
		// 						from buildkit?
		TimestampUnixNanos: uint64(time.Now().UnixNano()),
		Data:               dt,
	})
	return len(dt), nil
}
