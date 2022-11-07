package logbus

import (
	"fmt"
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// Generic is a generic writer for build output unrelated to a specific target.
type Generic struct {
	run      *Run
	category string
}

func newGeneric(run *Run, category string) *Generic {
	return &Generic{
		run:      run,
		category: category,
	}
}

// Write writes the given bytes to the generic printer.
func (g *Generic) Write(dt []byte) (int, error) {
	return g.WriteWithTimestamp(dt, time.Now())
}

// WriteWithTimestamp writes the given bytes to the generic printer.
func (g *Generic) WriteWithTimestamp(dt []byte, ts time.Time) (int, error) {
	g.run.b.WriteRawLog(&logstream.DeltaLog{
		CommandId:          fmt.Sprintf("_generic:%s", g.category),
		TimestampUnixNanos: uint64(ts.UnixNano()),
		Data:               dt,
	})
	return len(dt), nil
}
