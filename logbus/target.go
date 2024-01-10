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
func (t *Target) SetStart(start time.Time) {
	t.targetDelta(&logstream.DeltaTargetManifest{
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAtUnixNanos: t.b.TsUnixNanos(start),
	})
}

// SetEnd sets the end time of the target.
func (t *Target) SetEnd(end time.Time, status logstream.RunStatus, finalPlatform string) {
	t.targetDelta(&logstream.DeltaTargetManifest{
		Status:           status,
		EndedAtUnixNanos: t.b.TsUnixNanos(end),
		FinalPlatform:    finalPlatform,
	})
}

// AddDependsOn creates a delta that will be used to merge the specified target
// ID into the current target's list of targets on which it depends.
func (t *Target) AddDependsOn(targetID string) {
	t.targetDelta(&logstream.DeltaTargetManifest{
		DependsOn: []string{targetID},
	})
}

func (t *Target) targetDelta(dtm *logstream.DeltaTargetManifest) {
	t.b.WriteDeltaManifest(&logstream.DeltaManifest{
		DeltaManifestOneof: &logstream.DeltaManifest_Fields{
			Fields: &logstream.DeltaManifest_FieldsDelta{
				Targets: map[string]*logstream.DeltaTargetManifest{
					t.targetID: dtm,
				},
			},
		},
	})
}
