package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/ast/spec"
)

// TargetPrinter is a build log printer for a target.
type TargetPrinter struct {
	b        *Bus
	targetID string
	platform string

	mu      sync.Mutex
	started bool
	cps     []*CommandPrinter
}

// NextCommandPrinter creates a new command printer.
func (tp *TargetPrinter) NextCommandPrinter(command string, cached bool, push bool, local bool, sourceLocation *spec.SourceLocation) (int32, *CommandPrinter) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	index := int32(len(tp.cps))
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
				// SourceLocation: sourceLocation, // TODO
			},
		},
	})
	cp := &CommandPrinter{
		b:        tp.b,
		tp:       tp,
		targetID: tp.targetID,
		index:    index,
		cached:   cached,
		push:     push,
		local:    local,
	}
	tp.cps = append(tp.cps, cp)
	return int32(len(tp.cps)), cp
}

// CommandPrinter returns a command printer for a given index.
func (tp *TargetPrinter) CommandPrinter(index int32) *CommandPrinter {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.cps[index]
}

func (tp *TargetPrinter) maybeSetStart(start time.Time) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
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
	tp.mu.Lock()
	defer tp.mu.Unlock()
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
