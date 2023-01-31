package logbus

import (
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
	// TODO (vladaionescu): Can the timestamp be passed along straight
	// 						from buildkit?
	now := fw.bus.NowUnixNanos()
	fw.bus.WriteFormattedLog(&logstream.DeltaFormattedLog{
		TargetId:           fw.targetID,
		TimestampUnixNanos: now,
		Data:               dt,
	})
	fw.bus.WriteFormattedLog(&logstream.DeltaFormattedLog{
		TargetId:           "_full",
		TimestampUnixNanos: now,
		Data:               dt,
	})
	return len(dt), nil
}
