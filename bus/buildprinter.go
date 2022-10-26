package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// BuildPrinter is a build log printer.
type BuildPrinter struct {
	b     *Bus
	mu    sync.Mutex
	tps   map[string]*TargetPrinter
	ended bool

	gpMu sync.Mutex
	gps  map[string]*GenericPrinter
}

func newBuildPrinter(b *Bus) *BuildPrinter {
	return &BuildPrinter{
		b:   b,
		tps: make(map[string]*TargetPrinter),
		gps: make(map[string]*GenericPrinter),
	}
}

// GenericPrinter returns a generic printer for build output unrelated to a specific target.
func (bp *BuildPrinter) GenericPrinter(category string) *GenericPrinter {
	bp.gpMu.Lock()
	defer bp.gpMu.Unlock()
	gp, ok := bp.gps[category]
	if ok {
		return gp
	}
	gp = newGenericPrinter(bp, category)
	bp.gps[category] = gp
	return gp
}

// TargetPrinter creates a new target printer.
func (bp *BuildPrinter) TargetPrinter(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, platform string) *TargetPrinter {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	tp, ok := bp.tps[targetID]
	if ok {
		return tp
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Targets: map[string]*logstream.DeltaTargetManifest{
			targetID: {
				Name:          shortTargetName,
				CanonicalName: canonicalTargetName,
				Status:        logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
				OverrideArgs:  overrideArgs,
				Platform:      platform,
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
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:    logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAt: timestamppb.New(start),
	})
}

// SetFatalError sets a fatal error for the build.
func (bp *BuildPrinter) SetFatalError(end time.Time, targetID string, hasCommandIndex bool, commandIndex int32, output []byte, errString string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.ended {
		return
	}
	bp.ended = true
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:     logstream.RunStatus_RUN_STATUS_FAILURE,
		EndedAt:    timestamppb.New(end),
		HasFailure: true,
		Failure: &logstream.Failure{
			TargetId:        targetID,
			HasCommandIndex: hasCommandIndex,
			CommandIndex:    commandIndex,
			Output:          output,
			Summary:         errString,
		},
	})
}

// SetEnd sets the end time of the build.
func (bp *BuildPrinter) SetEnd(end time.Time, success bool, canceled bool, failureOutput []byte, failureSummary string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.ended {
		return
	}
	bp.ended = true
	var status logstream.RunStatus
	var f *logstream.Failure
	switch {
	case canceled:
		status = logstream.RunStatus_RUN_STATUS_CANCELED
	case success:
		status = logstream.RunStatus_RUN_STATUS_SUCCESS
	default:
		status = logstream.RunStatus_RUN_STATUS_FAILURE
		f = &logstream.Failure{
			Output:  failureOutput,
			Summary: failureSummary,
		}
	}

	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:  status,
		EndedAt: timestamppb.New(end),
		Failure: f,
	})
}

func (bp *BuildPrinter) buildDelta(fd *logstream.DeltaManifest_FieldsDelta) {
	bp.b.RawDelta(&logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: fd,
				},
			},
		},
	})
}
