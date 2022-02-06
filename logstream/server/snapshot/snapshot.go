package snapshot

import (
	"context"

	"github.com/earthly/earthly/logstream/api"
)

// Snapshot is a snapshot of a build's logs in a moment in time.
type Snapshot struct {
	// SnapshotID is a unique identifier for this snapshot, unique within a given build.
	SnapshotID string
	// Manifest is the manifest of the build.
	Manifest *api.Manifest
	// ManifestOrderID is the ordering ID of the manifest.
	ManifestOrderID int64
	// Logs is a map of target ID to logs.
	Logs map[string]*Target // targetID -> target
}

// Target contains the logs for a given target.
type Target struct {
	// LogData contains the logs for the target.
	LogData []byte
}

// EndSeekIndex returns the end seek index for the target.
func (t *Target) EndSeekIndex() int64 {
	return int64(len(t.LogData))
}

// Snapshotter is an object that can read and write snapshots.
type Snapshotter interface {
	Write(ctx context.Context, snapshot *Snapshot) error
	ReadManifest(ctx context.Context, snapshotID string) (*api.Manifest, error)
	ReadLogFragment(ctx context.Context, snapshotID string, targetID string, startSeekIndex int64, endSeekIndex int64) ([]byte, error)
}
