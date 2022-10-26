package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CommandPrinter is a build log printer for a command.
type CommandPrinter struct {
	b        *Bus
	tp       *TargetPrinter
	targetID string
	index    int32
	cached   bool
	push     bool
	local    bool

	mu          sync.Mutex
	started     bool
	hasProgress bool
}

// Write prints a byte slice with a timestamp.
func (cp *CommandPrinter) Write(dt []byte, ts time.Time, stream int32) (int, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
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
		// TODO: Tail output.
		EndedAt: timestamppb.New(end),
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
