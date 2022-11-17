package logbus

import (
	"errors"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/ast/spec"
)

// Run is a run logstream delta generator for a run.
type Run struct {
	b        *Bus
	mu       sync.Mutex
	targets  map[string]*Target
	commands map[string]*Command
	ended    bool

	generic *Generic
}

func newRun(b *Bus) *Run {
	run := &Run{
		b:        b,
		targets:  make(map[string]*Target),
		commands: make(map[string]*Command),
		generic:  nil, // set below
	}
	run.generic = newGeneric(run)
	return run
}

// Generic returns a generic writer for build output unrelated to a specific target.
func (run *Run) Generic() *Generic {
	return run.generic
}

// NewTarget creates a new target printer.
func (run *Run) NewTarget(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string) (*Target, error) {
	run.mu.Lock()
	defer run.mu.Unlock()
	_, ok := run.targets[targetID]
	if ok {
		return nil, errors.New("target printer already exists")
	}
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Targets: map[string]*logstream.DeltaTargetManifest{
			targetID: {
				Name:            shortTargetName,
				CanonicalName:   canonicalTargetName,
				OverrideArgs:    overrideArgs,
				InitialPlatform: initialPlatform,
			},
		},
	})
	target := newTarget(run.b, targetID)
	run.targets[targetID] = target
	return target, nil
}

// Target returns the target printer for the given target ID.
func (run *Run) Target(targetID string) (*Target, bool) {
	run.mu.Lock()
	defer run.mu.Unlock()
	target, ok := run.targets[targetID]
	return target, ok
}

// NewCommand creates a new command printer.
func (run *Run) NewCommand(commandID string, command string, targetID string, platform string, cached bool, push bool, local bool, sourceLocation *spec.SourceLocation, repoURL, repoHash, fileRelToRepo string) (*Command, error) {
	run.mu.Lock()
	defer run.mu.Unlock()
	_, ok := run.commands[commandID]
	if ok {
		return nil, errors.New("command printer already exists")
	}
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
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
	cp := newCommand(run.b, commandID, targetID)
	run.commands[commandID] = cp
	return cp, nil
}

// Command returns the command printer for the given command ID.
func (run *Run) Command(commandID string) (*Command, bool) {
	run.mu.Lock()
	defer run.mu.Unlock()
	cp, ok := run.commands[commandID]
	return cp, ok
}

// SetStart sets the start time of the build.
func (run *Run) SetStart(start time.Time) {
	run.mu.Lock()
	defer run.mu.Unlock()
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:             logstream.RunStatus_RUN_STATUS_IN_PROGRESS,
		StartedAtUnixNanos: run.b.TsUnixNanos(start),
	})
}

// SetFatalError sets a fatal error for the build.
func (run *Run) SetFatalError(end time.Time, targetID string, commandID string, failureType logstream.FailureType, errString string) {
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.ended {
		return
	}
	run.ended = true
	var tailOutput []byte
	if commandID != "" {
		cp, ok := run.commands[commandID]
		if ok {
			tailOutput = cp.TailOutput()
		}
	}
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:           logstream.RunStatus_RUN_STATUS_FAILURE,
		EndedAtUnixNanos: run.b.TsUnixNanos(end),
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
func (run *Run) SetEnd(end time.Time, success bool, canceled bool, failureOutput []byte, errorMessage string) {
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.ended {
		return
	}
	run.ended = true
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

	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:           status,
		EndedAtUnixNanos: run.b.TsUnixNanos(end),
		Failure:          f,
	})
}

func (run *Run) buildDelta(fd *logstream.DeltaManifest_FieldsDelta) {
	run.b.WriteDeltaManifest(&logstream.DeltaManifest{
		DeltaManifestOneof: &logstream.DeltaManifest_Fields{
			Fields: fd,
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
