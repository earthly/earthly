package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// BuildPrinter is a build log printer.
type BuildPrinter struct {
	b     *Bus
	tps   map[string]*TargetPrinter
	mu    sync.Mutex
	ended bool
}

func newBuildPrinter(b *Bus) *BuildPrinter {
	return &BuildPrinter{
		b:   b,
		tps: make(map[string]*TargetPrinter),
	}
}

// GetTargetPrinter creates a new target printer.
func (bp *BuildPrinter) GetTargetPrinter(targetID string, overrideArgs []string, platform string) *TargetPrinter {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	tp, ok := bp.tps[targetID]
	if ok {
		return tp
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Targets: map[string]*logstream.DeltaTargetManifest{
			targetID: {
				Name:         targetID, // TODO: Can this be a human-readable name instead?
				Status:       logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
				OverrideArgs: overrideArgs,
				Platform:     platform,
				StartedAt:    time.Now().Unix(),
			},
		},
	})
	tp = &TargetPrinter{
		b:        bp.b,
		targetID: targetID,
		platform: platform,
	}
	bp.tps[targetID] = tp
	return tp
}

// SetStart sets the start time of the build.
func (bp *BuildPrinter) SetStart(start time.Time) {
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:    logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
		StartedAt: start.Unix(),
	})
}

// SetFatalError sets a fatal error for the build.
func (bp *BuildPrinter) SetFatalError(end time.Time, failedTargetID string, failedCommandIndex int32, output []byte, errString string) {
	if bp.ended {
		return
	}
	bp.ended = true
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:       logstream.BuildStatus_BUILD_STATUS_FAILURE,
		FinishedAt:   end.Unix(),
		FailedTarget: failedTargetID,
		// FailedCommand: failedCommandIndex, // TODO
		FailedSummary: errString,
		// FailedOutput:  output, // TODO
	})
}

// SetEnd sets the end time of the build.
func (bp *BuildPrinter) SetEnd(end time.Time, success bool, canceled bool, failureSummary string) {
	if bp.ended {
		return
	}
	bp.ended = true
	var status logstream.BuildStatus
	switch {
	case canceled:
		status = logstream.BuildStatus_BUILD_STATUS_CANCELLED
	case success:
		status = logstream.BuildStatus_BUILD_STATUS_SUCCESS
	default:
		status = logstream.BuildStatus_BUILD_STATUS_FAILURE
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:        status,
		FinishedAt:    end.Unix(),
		FailedSummary: failureSummary,
	})
}

func (bp *BuildPrinter) buildDelta(fd *logstream.DeltaManifest_FieldsDelta) {
	bp.b.RawDelta(&logstream.Delta{
		DeltaManifests: []*logstream.DeltaManifest{
			{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: fd,
				},
			},
		},
	})
}
