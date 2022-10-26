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
		m.CreatedAt = m2.GetCreatedAt()
		m.StartedAt = m2.GetStartedAt()
		m.EndedAt = m2.GetEndedAt()
		m.Status = m2.GetStatus()
		m.MainTarget = m2.GetMainTarget()
		m.Failure = m2.GetFailure()
		m.UserId = m2.GetUserId()
		m.OrgId = m2.GetOrgId()
		m.ProjectId = m2.GetProjectId()
		newTargets := make(map[string]*pb.TargetManifest)
		for targetID, t2 := range m2.GetTargets() {
			newTargets[targetID] = &pb.TargetManifest{
				Name:          t2.GetName(),
				CanonicalName: t2.GetCanonicalName(),
				OverrideArgs:  append([]string{}, t2.GetOverrideArgs()...),
				Platform:      t2.GetPlatform(),
				Status:        t2.GetStatus(),
				StartedAt:     t2.GetStartedAt(),
				EndedAt:       t2.GetEndedAt(),
				Commands:      make([]*pb.CommandManifest, len(t2.GetCommands())),
			}
			for index, c2 := range t2.GetCommands() {
				newTargets[targetID].Commands[index] = &pb.CommandManifest{
					Name:           c2.GetName(),
					Status:         c2.GetStatus(),
					IsCached:       c2.GetIsCached(),
					IsPush:         c2.GetIsPush(),
					IsLocal:        c2.GetIsLocal(),
					StartedAt:      c2.GetStartedAt(),
					EndedAt:        c2.GetEndedAt(),
					HasProgress:    c2.GetHasProgress(),
					Progress:       c2.GetProgress(),
					ErrorMessage:   c2.GetErrorMessage(),
					SourceLocation: c2.GetSourceLocation(),
				}
			}
		}
		m.Targets = newTargets
	case *pb.DeltaManifest_Fields:
		f := dm.GetFields()
		if f.GetStartedAt().IsValid() {
			m.StartedAt = f.GetStartedAt()
		}
		if f.GetEndedAt().IsValid() {
			m.EndedAt = f.GetEndedAt()
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
			if t2.GetPlatform() != "" {
				t.Platform = t2.GetPlatform()
			}
			if t2.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
				t.Status = t2.GetStatus()
			}
			if t2.GetStartedAt().IsValid() {
				t.StartedAt = t2.GetStartedAt()
			}
			if t2.GetEndedAt().IsValid() {
				t.EndedAt = t2.GetEndedAt()
			}
			for index, c2 := range t2.GetCommands() {
				if index >= int32(len(t.Commands)) {
					for i := int32(len(t.Commands)); i <= index; i++ {
						t.Commands = append(t.Commands, &pb.CommandManifest{})
					}
				}
				c := t.Commands[index]
				if c2.GetName() != "" {
					c.Name = c2.GetName()
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
				if c2.GetStartedAt().IsValid() {
					c.StartedAt = c2.GetStartedAt()
				}
				if c2.GetEndedAt().IsValid() {
					c.EndedAt = c2.GetEndedAt()
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
	}
	return m, nil
}
