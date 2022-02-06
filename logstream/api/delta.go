package api

import (
	"time"

	"github.com/pkg/errors"
)

// Delta represents a set of changes to be applied atomically.
type Delta struct {
	// Version is the format version of the manifest.
	Version int `json:"version,omitempty"`

	DeltaManifests []*DeltaManifest `json:"delta_manifests,omitempty"`
	DeltaLogs      []*DeltaLog      `json:"delta_logs,omitempty"`
}

// DeltaManifest represents a change to a manifest.
type DeltaManifest struct {
	// OrderID is the ordering ID of the manifest.
	OrderID int64 `json:"order_id"`

	// Reset is the snapshot to reset the manifest to (overwrites everything).
	// If Reset is set, then the rest of the fields must be zero/unset.
	Reset *Manifest `json:"snapshot,omitempty"`

	// StartedAt is the start time of the build.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the build.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Status is the status of the build.
	Status Status `json:"status,omitempty"`
	// Targets is a map of target manifests. The key of the map is the sts ID.
	Targets map[string]*DeltaTargetManifest `json:"targets,omitempty"`
	// FailedTarget is the sts ID of the target that failed.
	FailedTarget string `json:"failed_target,omitempty"`
	// FailedSummary is an exerpt of the failing command's output.
	FailedSummary string `json:"summary,omitempty"`
}

// DeltaTargetManifest is the manifest of a target within a given build.
type DeltaTargetManifest struct {
	// Name is the name of the target (e.g. +something or ./path/to/something+else).
	Name string `json:"name,omitempty"`
	// OverrideArgs is the override args used for invoking this target, as a list of "key=value" strings.
	OverrideArgs []string `json:"override_args,omitempty"`
	// Platform is the override platform of the target.
	Platform string `json:"platform,omitempty"`
	// Status is the status of the target.
	Status Status `json:"status,omitempty"`
	// StartedAt is the start time of the target.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the target.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Commands is a map of command manifests. The key of the map is the order number,
	// starting from 0.
	Commands map[int]*DeltaCommandManifest `json:"commands,omitempty"`
}

// DeltaCommandManifest is the manifest of a command within a given target.
type DeltaCommandManifest struct {
	// Name is the name of the command, with all of its args
	// (e.g. "RUN hello world").
	Name string `json:"command,omitempty"`
	// Status is the status of the command.
	Status Status `json:"status,omitempty"`
	// Cached represents whether the command was previously cached.
	Cached *bool `json:"cached,omitempty"`
	// StartedAt is the start time of the command.
	StartedAt *time.Time `json:"started_at,omitempty"`
	// FinishedAt is the finish time of the command.
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	// Progress is an integer from 0 to 100 representing the % progress of the command.
	// Note that not all commands have a valid updating progress value.
	Progress *int `json:"progress,omitempty"`
}

// StartOrderID returns the data ordering ID of the delta start.
func (dm *DeltaManifest) StartOrderID() int64 {
	return dm.OrderID
}

// EndOrderID returns the data ordering ID of the delta end.
func (dm *DeltaManifest) EndOrderID() int64 {
	return dm.OrderID + 1
}

// DeltaLog represents log lines that should be appeneded to the log output for
// a given target.
type DeltaLog struct {
	// SeekIndex is the byte number where the data of the log belongs to.
	SeekIndex int64 `json:"seek_index"`
	// TargetID is the sts ID of the target.
	TargetID string `json:"target_id"`

	// Data is the data to append to the log output.
	Data []byte
}

// StartOrderID returns the data ordering ID of the delta start.
func (dl *DeltaLog) StartOrderID() int64 {
	return dl.SeekIndex
}

// EndOrderID returns the data ordering ID of the delta end.
func (dl *DeltaLog) EndOrderID() int64 {
	return dl.SeekIndex + int64(len(dl.Data))
}

// SimplifyDeltas takes a delta and
// * merges associated delta logs
// * merges delta manifests into a single, larger delta manifest
// It is assumed that the deltas contained within are ordered correctly.
func SimplifyDeltas(delta Delta, manifestOrderID int64) (Delta, int64) {
	if delta.Version != VersionNumber {
		panic(errors.Errorf("unsupported delta version %d", delta.Version))
	}
	var dms []*DeltaManifest
	nextManifestOrderID := manifestOrderID
	if len(delta.DeltaManifests) > 0 {
		dms = []*DeltaManifest{mergeDeltaManifests(delta.DeltaManifests, manifestOrderID)}
		nextManifestOrderID = manifestOrderID + 1
	}
	return Delta{
		Version:        VersionNumber,
		DeltaManifests: dms,
		DeltaLogs:      mergeDeltaLogs(delta.DeltaLogs),
	}, nextManifestOrderID
}

func mergeDeltaLogs(dls []*DeltaLog) []*DeltaLog {
	targets := make(map[string]*DeltaLog)
	for _, dl := range dls {
		tdl, found := targets[dl.TargetID]
		if !found {
			tdl = &DeltaLog{
				TargetID:  dl.TargetID,
				SeekIndex: dl.SeekIndex,
			}
			targets[dl.TargetID] = tdl
		}
		tdl.Data = append(tdl.Data, dl.Data...)
	}
	ret := make([]*DeltaLog, 0, len(targets))
	for _, tdl := range targets {
		ret = append(ret, tdl)
	}
	return ret
}

func mergeDeltaManifests(dms []*DeltaManifest, manifestOrderID int64) *DeltaManifest {
	ret := &DeltaManifest{
		OrderID: manifestOrderID,
	}
	for _, dm := range dms {
		if dm.Reset != nil {
			ret = &DeltaManifest{
				OrderID: manifestOrderID,
				Reset:   dm.Reset,
			}
			continue
		}
		if dm.StartedAt != nil {
			ret.StartedAt = dm.StartedAt
		}
		if dm.FinishedAt != nil {
			ret.FinishedAt = dm.FinishedAt
		}
		if dm.Status != "" {
			ret.Status = dm.Status
		}
		if dm.FailedTarget != "" {
			ret.FailedTarget = dm.FailedTarget
		}
		if dm.FailedSummary != "" {
			ret.FailedSummary = dm.FailedSummary
		}
		if dm.Targets != nil {
			if ret.Targets == nil {
				ret.Targets = make(map[string]*DeltaTargetManifest)
			}
			for targetID, dt := range dm.Targets {
				retDt, found := ret.Targets[targetID]
				if !found {
					retDt = &DeltaTargetManifest{}
					ret.Targets[targetID] = retDt
				}
				if dt.Name != "" {
					retDt.Name = dt.Name
				}
				if dt.OverrideArgs != nil {
					retDt.OverrideArgs = append(retDt.OverrideArgs, dt.OverrideArgs...)
				}
				if dt.Platform != "" {
					retDt.Platform = dt.Platform
				}
				if dt.Status != "" {
					retDt.Status = dt.Status
				}
				if dt.StartedAt != nil {
					retDt.StartedAt = dt.StartedAt
				}
				if dt.FinishedAt != nil {
					retDt.FinishedAt = dt.FinishedAt
				}
				for execOrder, dc := range dt.Commands {
					retDc, found := retDt.Commands[execOrder]
					if !found {
						retDc = &DeltaCommandManifest{}
						retDt.Commands[execOrder] = retDc
					}
					if dc.Name != "" {
						retDc.Name = dc.Name
					}
					if dc.Status != "" {
						retDc.Status = dc.Status
					}
					if dc.Cached != nil {
						retDc.Cached = dc.Cached
					}
					if dc.StartedAt != nil {
						retDc.StartedAt = dc.StartedAt
					}
					if dc.FinishedAt != nil {
						retDc.FinishedAt = dc.FinishedAt
					}
					if dc.Progress != nil {
						retDc.Progress = dc.Progress
					}
				}
			}
		}
	}
	return ret
}
