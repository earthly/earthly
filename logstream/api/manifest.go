package api

import (
	"time"

	"github.com/pkg/errors"
)

// VersionNumber is the currently supported manifest version.
const VersionNumber = 2

// Status represents the status of a build, of a target, or of a command.
type Status string

const (
	// StatusNotStarted is the not started status.
	StatusNotStarted Status = "not_started"
	// StatusInProgress is the in progress status.
	StatusInProgress Status = "in_progress"
	// StatusSuccess is the success status.
	StatusSuccess Status = "success"
	// StatusFailure is the failure status.
	StatusFailure Status = "failure"
	// StatusCancelled is the cancelled status.
	StatusCancelled Status = "cancelled"
)

// Manifest is the metadata associated with a build.
type Manifest struct {
	// Version is the format version of the manifest.
	Version int `json:"version"`
	// CreatedAt is the creation time of the build.
	CreatedAt time.Time `json:"created_at"`
	// StartedAt is the start time of the build.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the build.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Status is the status of the build.
	Status Status `json:"status,omitempty"`
	// Targets is a map of target manifests. The key of the map is the sts ID.
	Targets map[string]*TargetManifest `json:"targets"`
	// MainTarget is the main (final) target of the build.
	MainTarget string `json:"main_target"`
	// FailedTarget is the sts ID of the target that failed.
	FailedTarget string `json:"failed_target,omitempty"`
	// FailedSummary is an exerpt of the failing command's output.
	FailedSummary string `json:"summary,omitempty"`
}

// TargetManifest is the manifest of a target within a given build.
type TargetManifest struct {
	// Name is the name of the target (e.g. +something or ./path/to/something+else).
	Name string `json:"name"`
	// OverrideArgs is the override args used for invoking this target, as a list of "key=value" strings.
	OverrideArgs []string `json:"override_args"`
	// Platform is the override platform of the target.
	Platform string `json:"platform"`
	// Status is the status of the target.
	Status Status `json:"status"`
	// StartedAt is the start time of the target.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the target.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Size is the size of the logs of the target in bytes.
	Size *int64 `json:"size"`
	// Commands is a an list of commands ordered by their execution order.
	Commands []*CommandManifest `json:"commands"`
}

// CommandManifest is the manifest of a command within a given target.
type CommandManifest struct {
	// Name is the name of the command, with all of its args
	// (e.g. "RUN hello world").
	Name string `json:"command"`
	// Status is the status of the command.
	Status Status `json:"status"`
	// Cached represents whether the command was previously cached.
	Cached bool `json:"cached"`
	// StartedAt is the start time of the command.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the command.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Progress is an integer from 0 to 100 representing the % progress of the command.
	// Note that not all commands have a valid updating progress value.
	Progress *int `json:"progress,omitempty"`
}

// Apply delta takes a delta manifest and applies it to the current manifest.
func (m *Manifest) ApplyDelta(dm *DeltaManifest) error {
	if dm.Reset != nil {
		err := m.CopyFrom(dm.Reset)
		if err != nil {
			return err
		}
		return nil
	}
	if dm.StartedAt != nil {
		m.StartedAt = dm.StartedAt
	}
	if dm.FinishedAt != nil {
		m.FinishedAt = dm.FinishedAt
	}
	if dm.Status != "" {
		m.Status = dm.Status
	}
	if dm.FailedTarget != "" {
		m.FailedTarget = dm.FailedTarget
	}
	if dm.FailedSummary != "" {
		m.FailedSummary = dm.FailedSummary
	}
	if dm.Targets != nil {
		if m.Targets == nil {
			m.Targets = make(map[string]*TargetManifest)
		}
		for stsID, dt := range dm.Targets {
			t, found := m.Targets[stsID]
			if !found {
				t = new(TargetManifest)
				m.Targets[stsID] = t
			}
			if dt.Name != "" {
				t.Name = dt.Name
			}
			if dt.OverrideArgs != nil {
				t.OverrideArgs = append([]string{}, dt.OverrideArgs...)
			}
			if dt.Platform != "" {
				t.Platform = dt.Platform
			}
			if dt.Status != "" {
				t.Status = dt.Status
			}
			if dt.StartedAt != nil {
				t.StartedAt = dt.StartedAt
			}
			if dt.FinishedAt != nil {
				t.FinishedAt = dt.FinishedAt
			}
			if dt.Commands != nil {
				if t.Commands == nil {
					t.Commands = make([]*CommandManifest, 0, len(dt.Commands))
				}
				for iCmd, dc := range dt.Commands {
					for i := len(t.Commands); i <= iCmd; i++ {
						t.Commands = append(t.Commands, new(CommandManifest))
					}
					c := t.Commands[iCmd]
					if dc.Name != "" {
						c.Name = dc.Name
					}
					if dc.Status != "" {
						c.Status = dc.Status
					}
					if dc.Cached != nil {
						c.Cached = *dc.Cached
					}
					if dc.StartedAt != nil {
						c.StartedAt = dc.StartedAt
					}
					if dc.FinishedAt != nil {
						c.FinishedAt = dc.FinishedAt
					}
					if dc.Progress != nil {
						c.Progress = dc.Progress
					}
				}
			}
		}
	}
	return nil
}

func (m *Manifest) CopyFrom(m2 *Manifest) error {
	if m2.Version != VersionNumber {
		return errors.Errorf("manifest version %d is not supported", m2.Version)
	}
	*m = *m2
	m.Targets = make(map[string]*TargetManifest)
	for stsID, t2 := range m2.Targets {
		t := new(TargetManifest)
		m.Targets[stsID] = t
		*t = *t2
		t.Commands = append([]*CommandManifest{}, t2.Commands...)
	}
	return nil
}
