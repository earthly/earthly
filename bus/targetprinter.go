package bus

import (
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// TargetPrinter is a build log printer for a target.
type TargetPrinter struct {
	b        *Bus
	targetID string
}

func newTargetPrinter(b *Bus, targetID string) *TargetPrinter {
	return &TargetPrinter{
		b:        b,
		targetID: targetID,
	}
}

func (tp *TargetPrinter) SetStart(start time.Time) {
	tp.targetDelta(&logstream.DeltaTargetManifest{
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAtUnixNanos: uint64(start.UnixNano()),
	})
}

// SetEnd sets the end time of the target.
func (tp *TargetPrinter) SetEnd(end time.Time, status logstream.RunStatus, finalPlatform string) {
	tp.targetDelta(&logstream.DeltaTargetManifest{
		Status:           status,
		EndedAtUnixNanos: uint64(end.UnixNano()),
		FinalPlatform:    finalPlatform,
	})
}

func (tp *TargetPrinter) targetDelta(dtm *logstream.DeltaTargetManifest) {
	tp.b.SendRawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
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
