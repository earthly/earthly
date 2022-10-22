package deltautil

import (
	"fmt"

	pb "github.com/earthly/cloud-api/logstream"
)

// ApplyDeltaManifest takes a delta manifest and applies it to the given manifest,
// then returns it. This will mutate the originally passed-in
// manifest.
func ApplyDeltaManifest(m *pb.BuildManifest, d *pb.DeltaManifest) (*pb.BuildManifest, error) {
	if m.GetVersion() != 0 && m.GetVersion() != Version {
		return nil, fmt.Errorf("unsupported manifest version %d", m.GetVersion())
	}
	m.SnapshotOrderId = d.GetOrderId()
	switch d.GetDeltaManifestOneof().(type) {
	case *pb.DeltaManifest_Eof:
		if d.GetEof() {
			m.SnapshotEof = true
		}
		return m, nil
	case *pb.DeltaManifest_ResetAll:
		m2 := d.GetResetAll()
		if m2.GetVersion() != Version {
			return nil, fmt.Errorf("unsupported manifest version %d", m2.GetVersion())
		}
		m.Version = m2.GetVersion()
		m.CreatedAt = m2.GetCreatedAt()
		m.StartedAt = m2.GetStartedAt()
		m.FinishedAt = m2.GetFinishedAt()
		m.Status = m2.GetStatus()
		m.MainTarget = m2.GetMainTarget()
		m.FailedTarget = m2.GetFailedTarget()
		m.FailedSummary = m2.GetFailedSummary()
		newTargets := make(map[string]*pb.TargetManifest)
		for targetID, t2 := range m2.GetTargets() {
			ssize := int64(0)
			seof := false
			oldT, found := m.GetTargets()[targetID]
			if found {
				ssize = oldT.GetSnapshotSize()
				seof = oldT.GetSnapshotEof()
			}
			newTargets[targetID] = &pb.TargetManifest{
				Name:         t2.GetName(),
				OverrideArgs: append([]string{}, t2.GetOverrideArgs()...),
				Platform:     t2.GetPlatform(),
				Status:       t2.GetStatus(),
				StartedAt:    t2.GetStartedAt(),
				FinishedAt:   t2.GetFinishedAt(),
				Commands:     make([]*pb.CommandManifest, len(t2.GetCommands())),

				SnapshotSize: ssize,
				SnapshotEof:  seof,
			}
			for index, c2 := range t2.GetCommands() {
				newTargets[targetID].Commands[index] = &pb.CommandManifest{
					Name:        c2.GetName(),
					Status:      c2.GetStatus(),
					IsCached:    c2.GetIsCached(),
					IsPush:      c2.GetIsPush(),
					IsLocal:     c2.GetIsLocal(),
					StartedAt:   c2.GetStartedAt(),
					FinishedAt:  c2.GetFinishedAt(),
					HasProgress: c2.GetHasProgress(),
					Progress:    c2.GetProgress(),
				}
			}
		}
		m.Targets = newTargets
	case *pb.DeltaManifest_Fields:
		f := d.GetFields()
		if f.GetStartedAt() != 0 {
			m.StartedAt = f.GetStartedAt()
		}
		if f.GetFinishedAt() != 0 {
			m.FinishedAt = f.GetFinishedAt()
		}
		if f.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
			m.Status = f.GetStatus()
		}
		if f.GetFailedTarget() != "" {
			m.FailedTarget = f.GetFailedTarget()
		}
		if f.GetFailedSummary() != "" {
			m.FailedSummary = f.GetFailedSummary()
		}
		if f.GetTargets() != nil {
			if m.Targets == nil {
				m.Targets = make(map[string]*pb.TargetManifest)
			}
			for targetID, dt := range f.GetTargets() {
				t, found := m.Targets[targetID]
				if !found {
					t = new(pb.TargetManifest)
					m.Targets[targetID] = t
				}
				if dt.GetName() != "" {
					t.Name = dt.GetName()
				}
				if dt.GetOverrideArgs() != nil {
					t.OverrideArgs = append([]string{}, dt.GetOverrideArgs()...)
				}
				if dt.GetPlatform() != "" {
					t.Platform = dt.GetPlatform()
				}
				if dt.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
					t.Status = dt.GetStatus()
				}
				if dt.GetStartedAt() != 0 {
					t.StartedAt = dt.GetStartedAt()
				}
				if dt.GetFinishedAt() != 0 {
					t.FinishedAt = dt.GetFinishedAt()
				}
				if dt.GetCommands() != nil {
					for index, dcm := range dt.GetCommands() {
						for i := int32(len(t.Commands)); i <= index; i++ {
							t.Commands = append(t.Commands, new(pb.CommandManifest))
						}
						cm := t.Commands[index]
						if dcm.GetName() != "" {
							cm.Name = dcm.GetName()
						}
						if dcm.GetStatus() != pb.BuildStatus_BUILD_STATUS_UNKNOWN {
							cm.Status = dcm.GetStatus()
						}
						if dcm.GetHasCached() {
							cm.IsCached = dcm.GetIsCached()
						}
						if dcm.GetHasPush() {
							cm.IsPush = dcm.GetIsPush()
						}
						if dcm.GetHasLocal() {
							cm.IsLocal = dcm.GetIsLocal()
						}
						if dcm.GetStartedAt() != 0 {
							cm.StartedAt = dcm.GetStartedAt()
						}
						if dcm.GetFinishedAt() != 0 {
							cm.FinishedAt = dcm.GetFinishedAt()
						}
						if dcm.GetHasHasProgress() {
							cm.HasProgress = dcm.GetHasProgress()
						}
						if dcm.GetProgress() != 0 {
							cm.Progress = dcm.GetProgress()
						}
					}
				}

			}
		}
	}
	return m, nil
}

// ApplyDelta takes a delta and applies it to the given manifest, then returns
// it. This will mutate the originally passed-in manifest.
func ApplyDelta(m *pb.BuildManifest, d *pb.Delta) (*pb.BuildManifest, error) {
	if m.Targets == nil {
		m.Targets = make(map[string]*pb.TargetManifest)
	}
	for _, dm := range d.GetDeltaManifests() {
		var err error
		m, err = ApplyDeltaManifest(m, dm)
		if err != nil {
			return nil, err
		}
	}
	for _, dl := range d.GetDeltaLogs() {
		t, found := m.GetTargets()[dl.GetTargetId()]
		if !found {
			t = &pb.TargetManifest{}
			m.Targets[dl.GetTargetId()] = t
		}
		switch dl.GetDeltaLogOneof().(type) {
		case *pb.DeltaLog_Eof:
			if dl.GetEof() {
				t.SnapshotEof = true
			}
		case *pb.DeltaLog_Data:
			t.SnapshotSize += int64(len(dl.GetData()))
		}
	}
	return m, nil
}
