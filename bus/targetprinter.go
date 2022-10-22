package bus

import (
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// TargetPrinter is a build log printer for a target.
type TargetPrinter struct {
	b            *Bus
	targetID     string
	platform     string
	commandIndex int32
	started      bool
}

// NewCommandPrinter creates a new command printer.
func (tp *TargetPrinter) NewCommandPrinter(command string, cached bool, push bool, local bool) *CommandPrinter {
	// TODO: Add command source location.
	index := tp.commandIndex
	tp.commandIndex++
	tp.targetDelta(&logstream.DeltaTargetManifest{
		Commands: map[int32]*logstream.DeltaCommandManifest{
			index: {
				Name:      command,
				Status:    logstream.BuildStatus_BUILD_STATUS_NOT_STARTED,
				HasCached: true,
				IsCached:  cached,
				HasPush:   true,
				IsPush:    push,
				HasLocal:  true,
				IsLocal:   local,
			},
		},
	})
	return &CommandPrinter{
		b:        tp.b,
		tp:       tp,
		targetID: tp.targetID,
		index:    index,
		cached:   cached,
		push:     push,
		local:    local,
	}
}

func (tp *TargetPrinter) maybeSetStart(start time.Time) {
	if tp.started {
		tp.targetDelta(&logstream.DeltaTargetManifest{
			Status: logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
		})
		return
	}
	tp.started = true
	tp.targetDelta(&logstream.DeltaTargetManifest{
		Status:    logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
		StartedAt: start.Unix(),
	})
}

func (tp *TargetPrinter) setEnd(end time.Time, status logstream.BuildStatus, errorStr string) {
	tp.targetDelta(&logstream.DeltaTargetManifest{
		Status: status,
		// Error:      errorStr, // TODO
		FinishedAt: end.Unix(),
	})
}

func (tp *TargetPrinter) targetDelta(dtm *logstream.DeltaTargetManifest) {
	tp.b.RawDelta(&logstream.Delta{
		DeltaManifests: []*logstream.DeltaManifest{
			{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: &logstream.DeltaManifest_FieldsDelta{
						Targets: map[string]*logstream.DeltaTargetManifest{
							tp.targetID: dtm,
						},
					},
				},
			},
		},
	})
}
