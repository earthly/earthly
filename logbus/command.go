package logbus

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/cloud-api/logstream"
	"github.com/pkg/errors"
)

const tailErrorBufferSizeBytes = 80 * 1024 // About as much as 1024 lines of 80 chars each.

// Command is a build log writer for a command.
type Command struct {
	b         *Bus
	commandID string
	targetID  string

	tailOutput *circbuf.Buffer

	mu           sync.Mutex
	started      atomic.Bool
	lastProgress atomic.Int32
}

func newCommand(b *Bus, commandID string, targetID string) *Command {
	to, err := circbuf.NewBuffer(tailErrorBufferSizeBytes)
	if err != nil {
		panic(errors.Wrap(err, "failed to create tail buffer"))
	}
	return &Command{
		b:          b,
		commandID:  commandID,
		targetID:   targetID,
		tailOutput: to,
	}
}

// Write prints a byte slice with a timestamp.
func (c *Command) Write(dt []byte, ts time.Time, stream int32) (int, error) {
	c.mu.Lock()
	_, err := c.tailOutput.Write(dt)
	c.mu.Unlock()
	if err != nil {
		return 0, errors.Wrap(err, "write to tail output")
	}
	c.b.WriteRawLog(&logstream.DeltaLog{
		TargetId:           c.targetID,
		CommandId:          c.commandID,
		Stream:             stream,
		TimestampUnixNanos: c.b.TsUnixNanos(ts),
		Data:               dt,
	})
	return len(dt), nil
}

// TailOutput returns the tail of the output.
func (c *Command) TailOutput() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.tailOutput.Bytes()
}

// SetStart sets the start time of the command.
func (c *Command) SetStart(start time.Time) {
	if c.started.Load() {
		return
	}
	c.started.Store(true)
	c.commandDelta(&logstream.DeltaCommandManifest{
		StartedAtUnixNanos: c.b.TsUnixNanos(start),
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
	})
}

// SetProgress sets the progress of the command.
func (c *Command) SetProgress(progress int32) {
	if c.lastProgress.Load() == progress {
		return
	}
	c.commandDelta(&logstream.DeltaCommandManifest{
		HasHasProgress: true,
		HasProgress:    true,
		Progress:       progress,
	})
	c.lastProgress.Store(progress)
}

// SetCached sets the cached status of the command.
func (c *Command) SetCached(cached bool) {
	c.commandDelta(&logstream.DeltaCommandManifest{
		HasCached: true,
		IsCached:  cached,
	})
}

// SetEnd sets the end time of the command.
func (c *Command) SetEnd(end time.Time, status logstream.RunStatus, errorStr string) {
	c.commandDelta(&logstream.DeltaCommandManifest{
		Status:           status,
		ErrorMessage:     errorStr,
		EndedAtUnixNanos: c.b.TsUnixNanos(end),
	})
}

func (c *Command) commandDelta(dcm *logstream.DeltaCommandManifest) {
	c.b.WriteDeltaManifest(&logstream.DeltaManifest{
		DeltaManifestOneof: &logstream.DeltaManifest_Fields{
			Fields: &logstream.DeltaManifest_FieldsDelta{
				Commands: map[string]*logstream.DeltaCommandManifest{
					c.commandID: dcm,
				},
			},
		},
	})
}
