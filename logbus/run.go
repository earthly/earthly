package logbus

import (
	"errors"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/domain"
)

// GenericDefault is the internal name used to identify messages unrelated to a specific target or command.
const GenericDefault = "_generic:default"

// Run is a run logstream delta generator for a run.
type Run struct {
	b             *Bus
	mu            sync.Mutex
	targets       map[string]*Target
	commands      map[string]*Command
	ended         bool
	hasMainTarget bool

	generic *Generic
}

func newRun(b *Bus) *Run {
	run := &Run{
		b:             b,
		targets:       make(map[string]*Target),
		commands:      make(map[string]*Command),
		generic:       nil, // set below
		hasMainTarget: false,
	}
	run.generic = newGeneric(run)
	return run
}

// Generic returns a generic writer for build output unrelated to a specific target.
func (run *Run) Generic() *Generic {
	return run.generic
}

// NewTarget creates a new target printer.
func (run *Run) NewTarget(targetID string, target domain.Target, overrideArgs []string, initialPlatform string, runner string) (*Target, error) {
	run.mu.Lock()
	defer run.mu.Unlock()
	mainTargetID := ""
	if !run.hasMainTarget {
		// The first target is deemed as the main target.
		run.hasMainTarget = true
		mainTargetID = targetID
	}
	_, ok := run.targets[targetID]
	if ok {
		return nil, errors.New("target printer already exists")
	}
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		MainTargetId: mainTargetID,
		Targets: map[string]*logstream.DeltaTargetManifest{
			targetID: {
				Name:            target.String(), // Includes "+" prefix (e.g., "+target-name").
				CanonicalName:   target.StringCanonical(),
				GitUrl:          target.GetGitURL(),
				LocalPath:       target.GetLocalPath(),
				Tag:             target.GetTag(),
				ImportRef:       target.GetImportRef(),
				OverrideArgs:    overrideArgs,
				InitialPlatform: initialPlatform,
				Runner:          runner,
			},
		},
	})
	t := newTarget(run.b, targetID)
	run.targets[targetID] = t
	return t, nil
}

// Target returns the target printer for the given target ID.
func (run *Run) Target(targetID string) (*Target, bool) {
	run.mu.Lock()
	defer run.mu.Unlock()
	target, ok := run.targets[targetID]
	return target, ok
}

// NewCommand creates a new command printer.
func (run *Run) NewCommand(commandID, command string, targetID, category, platform string, cached, local, interactive bool, sourceLocation *spec.SourceLocation, repoURL, repoHash, fileRelToRepo string) (*Command, error) {
	run.mu.Lock()
	defer run.mu.Unlock()
	_, ok := run.commands[commandID]
	if ok {
		return nil, errors.New("command printer already exists")
	}
	sl := sourceLocationToProto(repoURL, repoHash, fileRelToRepo, sourceLocation)
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Commands: map[string]*logstream.DeltaCommandManifest{
			commandID: {
				Name:              command,
				TargetId:          targetID,
				Category:          category,
				Platform:          platform,
				HasCached:         true,
				IsCached:          cached,
				HasLocal:          true,
				IsLocal:           local,
				HasInteractive:    true,
				IsInteractive:     interactive,
				HasSourceLocation: (sl != nil),
				SourceLocation:    sl,
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
func (run *Run) SetFatalError(end time.Time, targetID string, commandID string, failureType logstream.FailureType, helpMsg, errorMsg string) {
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
			HelpMessage:  helpMsg,
			ErrorMessage: errorMsg,
		},
	})
}

// SetGenericFatalError sets a fatal error for the build with an empty target id and a command id indicating not to prefix the error with target info.
func (run *Run) SetGenericFatalError(end time.Time, failureType logstream.FailureType, helpMsg, errorMsg string) {
	run.SetFatalError(end, "", GenericDefault, failureType, helpMsg, errorMsg)
}

// SetEnd sets the end time and status of the build.
func (run *Run) SetEnd(end time.Time, status logstream.RunStatus) {
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.ended {
		return
	}
	run.ended = true
	run.buildDelta(&logstream.DeltaManifest_FieldsDelta{
		Status:           status,
		EndedAtUnixNanos: run.b.TsUnixNanos(end),
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
