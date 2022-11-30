package logbus

import (
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// Target is a delta generator for a target.
type Target struct {
	b        *Bus
	targetID string
}

func newTarget(b *Bus, targetID string) *Target {
	return &Target{
		b:        b,
		targetID: targetID,
	}
}

// SetStart sets the start time of the target.
func (target *Target) SetStart(start time.Time) {
	target.targetDelta(&logstream.DeltaTargetManifest{
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAtUnixNanos: target.b.TsUnixNanos(start),
	})
}

// SetEnd sets the end time of the target.
func (target *Target) SetEnd(end time.Time, status logstream.RunStatus, finalPlatform string) {
	target.targetDelta(&logstream.DeltaTargetManifest{
		Status:           status,
		EndedAtUnixNanos: target.b.TsUnixNanos(end),
		FinalPlatform:    finalPlatform,
	})
}

func (target *Target) targetDelta(dtm *logstream.DeltaTargetManifest) {
	target.b.WriteDeltaManifest(&logstream.DeltaManifest{
		DeltaManifestOneof: &logstream.DeltaManifest_Fields{
			Fields: &logstream.DeltaManifest_FieldsDelta{
				Targets: map[string]*logstream.DeltaTargetManifest{
					target.targetID: dtm,
				},
			},
		},
	})
}
