package logbus

import (
	"github.com/earthly/cloud-api/logstream"
)

// FormattedWriter is a writer that produces DeltaFormattedLog messages.
type FormattedWriter struct {
	bus       *Bus
	targetID  string
	commandID string
}

// NewFormattedWriter creates a new FormattedWriter.
func NewFormattedWriter(bus *Bus, targetID, commandID string) *FormattedWriter {
	return &FormattedWriter{
		bus:       bus,
		targetID:  targetID,
		commandID: commandID,
	}
}

// Write writes the given bytes to the writer.
func (w *FormattedWriter) Write(dt []byte) (int, error) {
	// TODO (vladaionescu): Can the timestamp be passed along straight
	// 						from buildkit?
	now := w.bus.NowUnixNanos()
	w.bus.WriteFormattedLog(&logstream.DeltaFormattedLog{
		TargetId:           w.targetID,
		CommandId:          w.commandID,
		TimestampUnixNanos: now,
		Data:               dt,
	})
	w.bus.WriteFormattedLog(&logstream.DeltaFormattedLog{
		TargetId:           "_full",
		CommandId:          w.commandID,
		TimestampUnixNanos: now,
		Data:               dt,
	})
	return len(dt), nil
}
