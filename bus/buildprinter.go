package bus

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
)

// BuildPrinter is a build log printer.
type BuildPrinter struct {
	b     *Bus
	mu    sync.Mutex
	tps   map[string]*TargetPrinter
	ended bool

	gpMu sync.Mutex
	gps  map[string]*genericPrinter
}

func newBuildPrinter(b *Bus) *BuildPrinter {
	return &BuildPrinter{
		b:   b,
		tps: make(map[string]*TargetPrinter),
		gps: make(map[string]*genericPrinter),
	}
}

// genericPrinter is a generic printer for build output unrelated to a specific target.
type genericPrinter struct {
	bp       *BuildPrinter
	category string
	mu       sync.Mutex
	size     int64
}

// Write writes the given bytes to the generic printer.
func (gp *genericPrinter) Write(dt []byte) (int, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	seekIndex := gp.size
	gp.size += int64(len(dt))
	gp.bp.b.RawDelta(&logstream.Delta{
		DeltaLogs: []*logstream.DeltaLog{
			{
				TargetId:  fmt.Sprintf("_generic:%s", gp.category),
				SeekIndex: seekIndex,
				DeltaLogOneof: &logstream.DeltaLog_Data{
					Data: dt,
				},
			},
		},
	})
	return len(dt), nil
}

// GenericPrinter returns a generic printer for build output unrelated to a specific target.
func (bp *BuildPrinter) GenericPrinter(category string) io.Writer {
	bp.gpMu.Lock()
	defer bp.gpMu.Unlock()
	gp, ok := bp.gps[category]
	if ok {
		return gp
	}
	gp = &genericPrinter{
		bp:       bp,
		category: category,
	}
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
				Name: targetID, // TODO: Remove?
				// ShortName:    shortTargetName, // TODO
				// CanonicalName: canonicalTargetName, // TODO
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
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:    logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS,
		StartedAt: start.Unix(),
	})
}

// SetFatalError sets a fatal error for the build.
func (bp *BuildPrinter) SetFatalError(end time.Time, failedTargetID string, failedCommandIndex int32, output []byte, errString string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
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
	bp.mu.Lock()
	defer bp.mu.Unlock()
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
