package bus

import (
	"sync"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/cloud-api/logstream"
	"github.com/pkg/errors"
)

const tailErrorBufferSizeBytes = 80 * 1024 // About as much as 1024 lines of 80 chars each.

// CommandPrinter is a build log printer for a command.
type CommandPrinter struct {
	b         *Bus
	commandID string
	targetID  string

	tailOutput *circbuf.Buffer

	mu           sync.Mutex
	started      bool
	lastProgress int32
}

func newCommandPrinter(b *Bus, commandID string, targetID string) *CommandPrinter {
	to, err := circbuf.NewBuffer(tailErrorBufferSizeBytes)
	if err != nil {
		panic(errors.Wrap(err, "failed to create tail buffer"))
	}
	return &CommandPrinter{
		b:          b,
		commandID:  commandID,
		targetID:   targetID,
		tailOutput: to,
	}
}

// Write prints a byte slice with a timestamp.
func (cp *CommandPrinter) Write(dt []byte, ts time.Time, stream int32) (int, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	_, err := cp.tailOutput.Write(dt)
	if err != nil {
		return 0, errors.Wrap(err, "write to tail output")
	}
	cp.b.RawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaLog{
			DeltaLog: &logstream.DeltaLog{
				TargetId:           cp.targetID,
				CommandId:          cp.commandID,
				Stream:             stream,
				TimestampUnixNanos: uint64(ts.UnixNano()),
				Data:               dt,
			},
		},
	})
	return len(dt), nil
}

// TailOutput returns the tail of the output.
func (cp *CommandPrinter) TailOutput() []byte {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return cp.tailOutput.Bytes()
}

// SetStart sets the start time of the command.
func (cp *CommandPrinter) SetStart(start time.Time) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.started {
		return
	}
	cp.started = true
	cp.commandDelta(&logstream.DeltaCommandManifest{
		StartedAtUnixNanos: uint64(start.UnixNano()),
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
	})
}

// SetProgress sets the progress of the command.
func (cp *CommandPrinter) SetProgress(progress int32) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.lastProgress == progress {
		return
	}
	cp.commandDelta(&logstream.DeltaCommandManifest{
		HasHasProgress: true,
		HasProgress:    true,
		Progress:       progress,
	})
	cp.lastProgress = progress
}

// SetCached sets the cached status of the command.
func (cp *CommandPrinter) SetCached(cached bool) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.commandDelta(&logstream.DeltaCommandManifest{
		HasCached: true,
		IsCached:  cached,
	})
}

// SetEnd sets the end time of the command.
func (cp *CommandPrinter) SetEnd(end time.Time, success bool, canceled bool, errorStr string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	var status logstream.RunStatus
	switch {
	case canceled:
		status = logstream.RunStatus_RUN_STATUS_CANCELED
	case success:
		status = logstream.RunStatus_RUN_STATUS_SUCCESS
	default:
		status = logstream.RunStatus_RUN_STATUS_FAILURE
	}
	cp.commandDelta(&logstream.DeltaCommandManifest{
		Status:           status,
		ErrorMessage:     errorStr,
		EndedAtUnixNanos: uint64(end.UnixNano()),
	})
}

func (cp *CommandPrinter) commandDelta(dcm *logstream.DeltaCommandManifest) {
	cp.b.RawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: &logstream.DeltaManifest_FieldsDelta{
						Commands: map[string]*logstream.DeltaCommandManifest{
							cp.commandID: dcm,
						},
					},
				},
			},
		},
	})
}
