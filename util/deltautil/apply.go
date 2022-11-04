package deltautil

import (
	"fmt"

	pb "github.com/earthly/cloud-api/logstream"
)

// Version is the version of the deltautil package.
const Version = 2

// ApplyDeltaManifest takes a delta and applies it to the given manifest,
// then returns it. This will mutate the originally passed-in
// manifest.
func ApplyDeltaManifest(m *pb.RunManifest, d *pb.Delta) (*pb.RunManifest, error) {
	if m.GetVersion() != 0 && m.GetVersion() != Version {
		return nil, fmt.Errorf("unsupported manifest version %d", m.GetVersion())
	}
	if d.GetVersion() != 0 && d.GetVersion() != Version {
		return nil, fmt.Errorf("unsupported delta version %d", d.GetVersion())
	}
	var dm *pb.DeltaManifest
	switch d.GetDeltaTypeOneof().(type) {
	case *pb.Delta_DeltaManifest:
		dm = d.GetDeltaManifest()
	default:
		// No action needed if this is not a manifest delta.
		return m, nil
	}
	switch dm.GetDeltaManifestOneof().(type) {
	case *pb.DeltaManifest_ResetAll:
		m2 := dm.GetResetAll()
		if m2.GetVersion() != Version {
			return nil, fmt.Errorf("unsupported manifest version %d", m2.GetVersion())
		}
		m.Version = m2.GetVersion()
		m.CreatedAtUnixNanos = m2.GetCreatedAtUnixNanos()
		m.StartedAtUnixNanos = m2.GetStartedAtUnixNanos()
		m.EndedAtUnixNanos = m2.GetEndedAtUnixNanos()
		m.Status = m2.GetStatus()
		m.MainTarget = m2.GetMainTarget()
		m.Failure = m2.GetFailure()
		m.OrgName = m2.GetOrgName()
		m.ProjectName = m2.GetProjectName()
		m.UserId = m2.GetUserId()
		m.OrgId = m2.GetOrgId()
		m.ProjectId = m2.GetProjectId()
		m.Targets = make(map[string]*pb.TargetManifest)
		for targetID, t2 := range m2.GetTargets() {
			m.Targets[targetID] = &pb.TargetManifest{
				Name:               t2.GetName(),
				CanonicalName:      t2.GetCanonicalName(),
				OverrideArgs:       append([]string{}, t2.GetOverrideArgs()...),
				InitialPlatform:    t2.GetInitialPlatform(),
				FinalPlatform:      t2.GetFinalPlatform(),
				Status:             t2.GetStatus(),
				StartedAtUnixNanos: t2.GetStartedAtUnixNanos(),
				EndedAtUnixNanos:   t2.GetEndedAtUnixNanos(),
			}
		}
		m.Commands = make(map[string]*pb.CommandManifest)
		for commandID, c2 := range m2.GetCommands() {
			m.Commands[commandID] = &pb.CommandManifest{
				Name:               c2.GetName(),
				TargetId:           c2.GetTargetId(),
				Platform:           c2.GetPlatform(),
				Status:             c2.GetStatus(),
				IsCached:           c2.GetIsCached(),
				IsPush:             c2.GetIsPush(),
				IsLocal:            c2.GetIsLocal(),
				StartedAtUnixNanos: c2.GetStartedAtUnixNanos(),
				EndedAtUnixNanos:   c2.GetEndedAtUnixNanos(),
				HasProgress:        c2.GetHasProgress(),
				Progress:           c2.GetProgress(),
				ErrorMessage:       c2.GetErrorMessage(),
				SourceLocation:     c2.GetSourceLocation(),
			}
		}
	case *pb.DeltaManifest_Fields:
		f := dm.GetFields()
		if f.GetStartedAtUnixNanos() != 0 {
			m.StartedAtUnixNanos = f.GetStartedAtUnixNanos()
		}
		if f.GetEndedAtUnixNanos() != 0 {
			m.EndedAtUnixNanos = f.GetEndedAtUnixNanos()
		}
		if f.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
			m.Status = f.GetStatus()
		}
		if f.GetHasFailure() {
			m.Failure = f.GetFailure()
		}
		for targetID, t2 := range f.GetTargets() {
			if m.Targets == nil {
				m.Targets = make(map[string]*pb.TargetManifest)
			}
			t, ok := m.Targets[targetID]
			if !ok {
				t = &pb.TargetManifest{}
				m.Targets[targetID] = t
			}

			if t2.GetName() != "" {
				t.Name = t2.GetName()
			}
			if t2.GetCanonicalName() != "" {
				t.CanonicalName = t2.GetCanonicalName()
			}
			if len(t2.GetOverrideArgs()) > 0 {
				t.OverrideArgs = append([]string{}, t2.GetOverrideArgs()...)
			}
			if t2.GetInitialPlatform() != "" {
				t.InitialPlatform = t2.GetInitialPlatform()
			}
			if t2.GetFinalPlatform() != "" {
				t.FinalPlatform = t2.GetFinalPlatform()
			}
			if t2.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
				t.Status = t2.GetStatus()
			}
			if t2.GetStartedAtUnixNanos() != 0 {
				t.StartedAtUnixNanos = t2.GetStartedAtUnixNanos()
			}
			if t2.GetEndedAtUnixNanos() != 0 {
				t.EndedAtUnixNanos = t2.GetEndedAtUnixNanos()
			}
		}
		for commandID, c2 := range f.GetCommands() {
			if m.Commands == nil {
				m.Commands = make(map[string]*pb.CommandManifest)
			}
			c, ok := m.Commands[commandID]
			if !ok {
				c = &pb.CommandManifest{}
				m.Commands[commandID] = c
			}

			if c2.GetName() != "" {
				c.Name = c2.GetName()
			}
			if c2.GetTargetId() != "" {
				c.TargetId = c2.GetTargetId()
			}
			if c2.GetPlatform() != "" {
				c.Platform = c2.GetPlatform()
			}
			if c2.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
				c.Status = c2.GetStatus()
			}
			if c2.GetHasCached() {
				c.IsCached = c2.GetIsCached()
			}
			if c2.GetHasPush() {
				c.IsPush = c2.GetIsPush()
			}
			if c2.GetHasLocal() {
				c.IsLocal = c2.GetIsLocal()
			}
			if c2.GetStartedAtUnixNanos() != 0 {
				c.StartedAtUnixNanos = c2.GetStartedAtUnixNanos()
			}
			if c2.GetEndedAtUnixNanos() != 0 {
				c.EndedAtUnixNanos = c2.GetEndedAtUnixNanos()
			}
			if c2.GetHasHasProgress() {
				c.HasProgress = c2.GetHasProgress()
			}
			if c2.GetHasProgress() {
				c.Progress = c2.GetProgress()
			}
			if c2.GetErrorMessage() != "" {
				c.ErrorMessage = c2.GetErrorMessage()
			}
			if c2.GetHasSourceLocation() {
				c.SourceLocation = c2.GetSourceLocation()
			}
		}
	}
	return m, nil
}
