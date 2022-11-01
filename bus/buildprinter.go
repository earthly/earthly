package bus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/ast/spec"
)

// BuildPrinter is a build log printer.
type BuildPrinter struct {
	b     *Bus
	mu    sync.Mutex
	tps   map[string]*TargetPrinter
	cps   map[string]*CommandPrinter
	ended bool

	gpMu sync.Mutex
	gps  map[string]*GenericPrinter
}

func newBuildPrinter(b *Bus) *BuildPrinter {
	return &BuildPrinter{
		b:   b,
		tps: make(map[string]*TargetPrinter),
		cps: make(map[string]*CommandPrinter),
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

// NewTargetPrinter creates a new target printer.
func (bp *BuildPrinter) NewTargetPrinter(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string) *TargetPrinter {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	tp, ok := bp.tps[targetID]
	if ok {
		return tp
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Targets: map[string]*logstream.DeltaTargetManifest{
			targetID: {
				Name:            shortTargetName,
				CanonicalName:   canonicalTargetName,
				OverrideArgs:    overrideArgs,
				InitialPlatform: initialPlatform,
			},
		},
	})
	tp = newTargetPrinter(bp.b, targetID)
	bp.tps[targetID] = tp
	return tp
}

// TargetPrinter returns the target printer for the given target ID.
func (bp *BuildPrinter) TargetPrinter(targetID string) (*TargetPrinter, bool) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	tp, ok := bp.tps[targetID]
	return tp, ok
}

// NewCommandPrinter creates a new command printer.
func (bp *BuildPrinter) NewCommandPrinter(commandID string, command string, targetID string, platform string, cached bool, push bool, local bool, sourceLocation *spec.SourceLocation, repoURL, repoHash, fileRelToRepo string) *CommandPrinter {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	cp, ok := bp.cps[commandID]
	if ok {
		return cp
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Commands: map[string]*logstream.DeltaCommandManifest{
			commandID: {
				Name:           command,
				TargetId:       targetID,
				Platform:       platform,
				IsCached:       cached,
				IsPush:         push,
				IsLocal:        local,
				SourceLocation: sourceLocationToProto(repoURL, repoHash, fileRelToRepo, sourceLocation),
			},
		},
	})
	cp = newCommandPrinter(bp.b, commandID, targetID)
	bp.cps[commandID] = cp
	return cp
}

// CommandPrinter returns the command printer for the given command ID.
func (bp *BuildPrinter) CommandPrinter(commandID string) (*CommandPrinter, bool) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	cp, ok := bp.cps[commandID]
	return cp, ok
}

// SetStart sets the start time of the build.
func (bp *BuildPrinter) SetStart(start time.Time) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAtUnixNanos: uint64(start.UnixNano()),
	})
}

// SetFatalError sets a fatal error for the build.
func (bp *BuildPrinter) SetFatalError(end time.Time, targetID string, commandID string, failureType logstream.FailureType, errString string) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.ended {
		return
	}
	bp.ended = true
	var tailOutput []byte
	if commandID != "" {
		cp, ok := bp.CommandPrinter(commandID)
		if !ok {
			panic("command printer not found")
		}
		tailOutput = cp.TailOutput()
	}
	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:           logstream.RunStatus_RUN_STATUS_FAILURE,
		EndedAtUnixNanos: uint64(end.UnixNano()),
		HasFailure:       true,
		Failure: &logstream.Failure{
			Type:         failureType,
			TargetId:     targetID,
			CommandId:    commandID,
			Output:       tailOutput,
			ErrorMessage: errString,
		},
	})
}

// SetEnd sets the end time of the build.
func (bp *BuildPrinter) SetEnd(end time.Time, success bool, canceled bool, failureOutput []byte, errorMessage string) {
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
			Output:       failureOutput,
			ErrorMessage: errorMessage,
		}
	}

	bp.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:           status,
		EndedAtUnixNanos: uint64(end.UnixNano()),
		Failure:          f,
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

func sourceLocationToProto(repoURL, repoHash, fileRelToRepo string, sl *spec.SourceLocation) *logstream.SourceLocation {
	if sl == nil {
		return nil
	}
	file := fileRelToRepo
	if fileRelToRepo == "" && repoURL == "" {
		file = sl.File
	}
	return &logstream.SourceLocation{
		RepositoryUrl:  repoURL,
		RepositoryHash: repoHash,
		File:           file,
		StartLine:      int32(sl.StartLine),
		StartColumn:    int32(sl.StartColumn),
		EndLine:        int32(sl.EndLine),
		EndColumn:      int32(sl.EndColumn),
	}
}
