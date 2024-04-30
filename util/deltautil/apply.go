package deltautil

import (
	"fmt"
	"slices"

	pb "github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/proto"
)

// Version is the version of the deltautil package.
const Version = 2

// ApplyDelta applies the manifest changes specified by d to the manifest, m. It
// does so by modifying the manifest directly.
func ApplyDelta(m *pb.RunManifest, d *pb.Delta) error {
	err := checkVersion(m, d)
	if err != nil {
		return err
	}
	dm := d.GetDeltaManifest()
	if dm == nil {
		return nil
	}
	switch dm.GetDeltaManifestOneof().(type) {
	case *pb.DeltaManifest_ResetAll:
		*m = *dm.GetResetAll()
	case *pb.DeltaManifest_Fields:
		setManifestFields(dm, m)
	}
	return nil
}

// WithDeltaManifest takes a delta and a manifest and returns the result of
// applying the delta to the manifest. The original passed-in manifest is not
// changed. If the delta would not have any effect on the manifest, the original
// manifest is returned.
func WithDeltaManifest(m *pb.RunManifest, d *pb.Delta) (*pb.RunManifest, error) {
	err := checkVersion(m, d)
	if err != nil {
		return nil, err
	}
	dm := d.GetDeltaManifest()
	if dm == nil {
		// No action needed if this is not a manifest delta.
		return m, nil
	}
	var ret *pb.RunManifest
	switch dm.GetDeltaManifestOneof().(type) {
	case *pb.DeltaManifest_ResetAll:
		ret = proto.Clone(dm.GetResetAll()).(*pb.RunManifest)
	case *pb.DeltaManifest_Fields:
		ret = proto.Clone(m).(*pb.RunManifest)
		setManifestFields(dm, ret)
	}
	return ret, nil
}

func checkVersion(m *pb.RunManifest, d *pb.Delta) error {
	if m.GetVersion() != 0 && m.GetVersion() != Version {
		return fmt.Errorf("unsupported manifest version %d", m.GetVersion())
	}
	if d.GetVersion() != 0 && d.GetVersion() != Version {
		return fmt.Errorf("unsupported delta version %d", d.GetVersion())
	}
	return nil
}

func setManifestFields(dm *pb.DeltaManifest, ret *pb.RunManifest) {
	f := dm.GetFields()
	if f.GetStartedAtUnixNanos() != 0 {
		ret.StartedAtUnixNanos = f.GetStartedAtUnixNanos()
	}
	if f.GetEndedAtUnixNanos() != 0 {
		ret.EndedAtUnixNanos = f.GetEndedAtUnixNanos()
	}
	if f.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
		ret.Status = f.GetStatus()
	}
	if f.GetHasFailure() {
		ret.Failure = f.GetFailure()
	}
	if f.GetMainTargetId() != "" {
		ret.MainTargetId = f.GetMainTargetId()
	}
	if f.GetIsCi() {
		ret.IsCi = true
	}
	if f.GetGitAuthor() != "" {
		ret.GitAuthor = f.GetGitAuthor()
	}
	if f.GetGitConfigEmail() != "" {
		ret.GitConfigEmail = f.GetGitConfigEmail()
	}
	for targetID, t2 := range f.GetTargets() {
		t := ensureTargetExists(ret, targetID)

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
		if t2.GetRunner() != "" {
			t.Runner = t2.GetRunner()
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
		if t2.GetGitUrl() != "" {
			t.GitUrl = t2.GetGitUrl()
		}
		if t2.GetTag() != "" {
			t.Tag = t2.GetTag()
		}
		if t2.GetImportRef() != "" {
			t.ImportRef = t2.GetImportRef()
		}
		if t2.GetLocalPath() != "" {
			t.LocalPath = t2.GetLocalPath()
		}
		if len(t2.GetDependsOn()) > 0 {
			t.DependsOn = append(t.DependsOn, t2.GetDependsOn()...)
		}
	}
	for commandID, c2 := range f.GetCommands() {
		c := ensureCommandExists(ret, commandID)

		if c2.GetName() != "" {
			c.Name = c2.GetName()
		}
		if c2.GetTargetId() != "" {
			c.TargetId = c2.GetTargetId()
			ensureTargetExists(ret, c2.GetTargetId())
		}
		if c2.GetCategory() != "" {
			c.Category = c2.GetCategory()
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
		if c2.GetHasLocal() {
			c.IsLocal = c2.GetIsLocal()
		}
		if c2.GetHasInteractive() {
			c.IsInteractive = c2.GetIsInteractive()
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
		if len(c2.GetDependsOn()) > 0 {
			c.DependsOn = append(c.DependsOn, c2.GetDependsOn()...)
		}

		for _, image := range c2.GetImages() {
			i := ensureImageExists(ret, image.GetDigest())
			i.ImageName = image.GetImageName()
			if image.GetTag() != "" {
				i.Tags = append(i.GetTags(), image.GetTag())
				slices.Sort(i.Tags)
				i.Tags = slices.Compact(i.Tags)
			}
			i.Platforms = append(i.GetPlatforms(), image.GetPlatform())
			slices.Sort(i.Platforms)
			i.Platforms = slices.Compact(i.Platforms)
			i.Vulnerabilities = image.GetVulnerabilities()
			//sboms := GetSBOMSbyImageName(i.ImageName)
			//i.Spdx = append(i.Spdx, sboms...)
			//fmt.Printf("Added %d SPDXs to %s\n", len(sboms), i.ImageName)
		}
	}
}

func ensureTargetExists(r *pb.RunManifest, targetID string) *pb.TargetManifest {
	if r.Targets == nil {
		r.Targets = make(map[string]*pb.TargetManifest)
	}
	t, ok := r.Targets[targetID]
	if !ok {
		t = &pb.TargetManifest{}
		r.Targets[targetID] = t
	}
	return t
}

func ensureCommandExists(r *pb.RunManifest, commandID string) *pb.CommandManifest {
	if r.Commands == nil {
		r.Commands = make(map[string]*pb.CommandManifest)
	}
	c, ok := r.Commands[commandID]
	if !ok {
		c = &pb.CommandManifest{}
		r.Commands[commandID] = c
	}
	return c
}

func ensureImageExists(r *pb.RunManifest, digest string) *pb.ImageManifest {
	if r.Images == nil {
		r.Images = make(map[string]*pb.ImageManifest)
	}
	i, ok := r.Images[digest]
	if !ok {
		i = &pb.ImageManifest{
			Digest: digest,
		}
		r.Images[digest] = i
	}
	return i
}
