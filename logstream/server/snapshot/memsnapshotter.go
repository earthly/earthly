package snapshot

import (
	"context"
	"sync"

	"github.com/earthly/earthly/logstream/api"
	"github.com/pkg/errors"
)

// MemSnapshotter is a snapshotter that stores all data in memory. This is meant
// for testing / mocking only.
type MemSnapshotter struct {
	snapshots map[string]*Snapshot

	mu sync.Mutex
}

// NewMemSnapshotter returns a new MemSnapshotter.
func NewMemSnapshotter() *MemSnapshotter {
	return &MemSnapshotter{
		snapshots: make(map[string]*Snapshot),
	}
}

func (ms *MemSnapshotter) Write(ctx context.Context, snapshot *Snapshot) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.snapshots[snapshot.SnapshotID] = snapshot
	return nil
}

func (ms *MemSnapshotter) ReadManifest(ctx context.Context, snapshotID string) (*api.Manifest, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	snapshot, found := ms.snapshots[snapshotID]
	if !found {
		return nil, errors.Errorf("snapshot %s not found", snapshotID)
	}
	return snapshot.Manifest, nil
}

func (ms *MemSnapshotter) ReadLogFragment(ctx context.Context, snapshotID string, targetID string, startSeekIndex int64, endSeekIndex int64) ([]byte, int64, int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	snapshot, found := ms.snapshots[snapshotID]
	if !found {
		return nil, 0, 0, errors.Errorf("snapshot %s not found", snapshotID)
	}
	target, found := snapshot.Logs[targetID]
	if !found {
		return nil, 0, 0, errors.Errorf("target %s not found", targetID)
	}
	if startSeekIndex < 0 {
		startSeekIndex = int64(len(target.LogData)) + startSeekIndex
	}
	if startSeekIndex > int64(len(target.LogData)) {
		return nil, 0, 0, errors.Errorf("start seek index %d out of bounds", startSeekIndex)
	}
	if endSeekIndex <= 0 {
		endSeekIndex = int64(len(target.LogData)) + endSeekIndex
	}
	if endSeekIndex > int64(len(target.LogData)) {
		return nil, 0, 0, errors.Errorf("end seek index %d out of bounds", endSeekIndex)
	}
	return target.LogData[startSeekIndex:endSeekIndex], startSeekIndex, endSeekIndex, nil
}
