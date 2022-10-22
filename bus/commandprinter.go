package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
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
	size        int64
}

// Write prints a byte slice.
func (cp *CommandPrinter) Write(dt []byte) (int, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	seekIndex := cp.size
	cp.size += int64(len(dt))
	cp.b.RawDelta(&logstream.Delta{
		DeltaLogs: []*logstream.DeltaLog{
			{
				TargetId:  cp.targetID,
				SeekIndex: seekIndex,
				// TODO: Add command index?
				DeltaLogOneof: &logstream.DeltaLog_Data{
					Data: dt,
				},
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
		StartedAt: start.Unix(),
		Status:    logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
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
	var status logstream.BuildStatus
	switch {
	case canceled:
		status = logstream.BuildStatus_BUILD_STATUS_CANCELLED
	case success:
		status = logstream.BuildStatus_BUILD_STATUS_SUCCESS
	default:
		status = logstream.BuildStatus_BUILD_STATUS_FAILURE
	}
	cp.commandDelta(&logstream.DeltaCommandManifest{
		Status: status,
		// Error:      errorStr, // TODO
		FinishedAt: end.Unix(),
	})
	cp.tp.setEnd(end, status, errorStr)
}

func (cp *CommandPrinter) commandDelta(dcm *logstream.DeltaCommandManifest) {
	cp.b.RawDelta(&logstream.Delta{
		DeltaManifests: []*logstream.DeltaManifest{
			{
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
