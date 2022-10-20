package outmon

import (
	"time"

	"github.com/earthly/earthly/ast/spec"
)

// Status is the status of a command.
type Status string

const (
	// StatusSuccess is the status of a command that succeeded.
	StatusSuccess Status = "success"
	// StatusFailure is the status of a command that failed.
	StatusFailure Status = "failure"
	// StatusCanceled is the status of a command that was canceled.
	StatusCanceled Status = "canceled"
)

// BuildAnalytics is a struct that contains the analytics for a build.
type BuildAnalytics struct {
	FinalTarget string
	StartTime   time.Time
	EndTime     time.Time
	Status      Status

	// The following fields are only set if the build failed.
	FailedTarget  string
	FailedSummary string

	// Runner information.
	SatelliteName string // empty if ran locally

	// Target details.
	Targets map[string]*TargetAnalytics
}

func (ba *BuildAnalytics) getTarget(target string) *TargetAnalytics {
	if ba.Targets == nil {
		ba.Targets = map[string]*TargetAnalytics{}
	}
	if _, ok := ba.Targets[target]; !ok {
		ba.Targets[target] = &TargetAnalytics{}
	}
	return ba.Targets[target]
}

// TargetAnalytics is a struct that contains the analytics for a target.
type TargetAnalytics struct {
	StartTime time.Time
	EndTime   time.Time
	// ProcessingDuration represents the real processing time. This can differ
	// than EndTime - StartTime if the build scheduler did not schedule some
	// intermediary commands immediately.
	ProcessingDuration time.Duration
	Status             Status
	Commands           []*CommandAnalytics
}

// CommandAnalytics is a struct that contains the analytics for a command.
type CommandAnalytics struct {
	Command        string
	SourceLocation *spec.SourceLocation
	StartTime      time.Time
	EndTime        time.Time
	Status         Status
	IsCached       bool
	IsLocal        bool
}
