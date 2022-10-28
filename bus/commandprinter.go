package bus

import (
	"sync"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/cloud-api/logstream"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const tailErrorBufferSizeBytes = 80 * 1024 // About as much as 1024 lines of 80 chars each.

// CommandPrinter is a build log printer for a command.
type CommandPrinter struct {
	b        *Bus
	tp       *TargetPrinter
	targetID string
	index    int32
	cached   bool
	push     bool
	local    bool

	tailOutput *circbuf.Buffer

	mu          sync.Mutex
	started     bool
	hasProgress bool
}

func newCommandPrinter(b *Bus, tp *TargetPrinter, targetID string, index int32, cached bool, push bool, local bool) *CommandPrinter {
	to, err := circbuf.NewBuffer(tailErrorBufferSizeBytes)
	if err != nil {
		panic(errors.Wrap(err, "failed to create tail buffer"))
	}
	return &CommandPrinter{
		b:          b,
		tp:         tp,
		targetID:   targetID,
		index:      index,
		cached:     cached,
		push:       push,
		local:      local,
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
				TargetId:     cp.targetID,
				CommandIndex: cp.index,
				Stream:       stream,
				Timestamp:    timestamppb.New(ts),
				Log:          dt,
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

// Index returns the index of the command.
func (cp *CommandPrinter) Index() int32 {
	return cp.index
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
		StartedAt: timestamppb.New(start),
		Status:    logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
	})
	cp.tp.maybeSetStart(start)
}

// SetProgress sets the progress of the command.
func (cp *CommandPrinter) SetProgress(progress int32) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if !cp.hasProgress {
		cp.commandDelta(&logstream.DeltaCommandManifest{
			HasHasProgress: true,
			HasProgress:    true,
		})
	}
	cp.hasProgress = true
	cp.commandDelta(&logstream.DeltaCommandManifest{
		Progress: progress,
	})
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
		Status:       status,
		ErrorMessage: errorStr,
		EndedAt:      timestamppb.New(end),
	})
	cp.tp.setEnd(end, status)
}

func (cp *CommandPrinter) commandDelta(dcm *logstream.DeltaCommandManifest) {
	cp.b.RawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: &logstream.DeltaManifest_FieldsDelta{
						Targets: map[string]*logstream.DeltaTargetManifest{
							cp.targetID: {
								Commands: map[int32]*logstream.DeltaCommandManifest{
									cp.index: dcm,
								},
							},
						},
					},
				},
			},
		},
	})
}
