package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/server/snapshot"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Build is a server for a given build. The server is able to both
// ingest and serve deltas. It also performs occasional snapshotting,
// to pack the deltas more efficiently.
type Build struct {
	manifest         api.Manifest
	manifestStream   *DeltaStream
	manifestOrderID  int64
	targets          map[string]*target
	subs             map[string]*subscriber
	snapshotter      snapshot.Snapshotter
	latestSnapshotID string

	mu sync.RWMutex
}

type target struct {
	stream         *DeltaStream
	logData        []byte
	startSeekIndex int64
	endSeekIndex   int64

	mu sync.RWMutex
}

// SubscriberOpt holds options for a subscriber.
type SubscriberOpt struct {
	Manifest        bool
	ManifestOrderID int64
	AllTargets      bool
	TargetSeeks     map[string]int64
	Emit            func(delta api.Delta) error
}

type subscriber struct {
	manifest   bool
	allTargets bool
	targets    map[string]bool
	emit       func(delta api.Delta) error

	mu sync.RWMutex
}

// Build returns a new Build.
func NewBuild(snapshotter snapshot.Snapshotter) *Build {
	b := &Build{
		manifest: api.Manifest{
			Version: api.VersionNumber,
		},
		manifestStream: nil, // set below
		targets:        make(map[string]*target),
		subs:           make(map[string]*subscriber),
		snapshotter:    snapshotter,
	}
	b.manifestStream = NewManifestDeltaStream(b.handleDeltaManifest, nil)
	return b
}

// ReceiveDelta takes in a delta and applies it to the server's build model,
// taking into account ordering (via DeltaStreams).
func (b *Build) ReceiveDelta(delta api.Delta) error {
	for _, dm := range delta.DeltaManifests {
		err := b.manifestStream.Receive(dm)
		if err != nil {
			return err
		}
	}
	for _, dl := range delta.DeltaLogs {
		err := b.receiveDeltaLog(dl)
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe creates a new subscriber for the given build.
func (b *Build) Subscribe(opt *SubscriberOpt) (string, error) {
	subID := uuid.New().String()
	b.mu.Lock()
	defer b.mu.Unlock()
	sub := &subscriber{
		manifest:   opt.Manifest,
		allTargets: opt.AllTargets,
		targets:    make(map[string]bool),
		emit:       opt.Emit,
	}
	for stsID := range opt.TargetSeeks {
		sub.targets[stsID] = true
	}

	b.subs[subID] = sub
	sub.mu.Lock()
	defer sub.mu.Unlock()

	initialDelta := api.Delta{}
	if sub.manifest && opt.ManifestOrderID < b.manifestOrderID {
		dm := &api.DeltaManifest{
			OrderID: b.manifestOrderID,
			Reset: &api.Manifest{
				Version: api.VersionNumber,
			},
		}
		dm.Reset.CopyFrom(&b.manifest)
		initialDelta.DeltaManifests = append(initialDelta.DeltaManifests, dm)
	}

	if sub.allTargets {
		for stsID, t := range b.targets {
			seek := int64(0)
			subSeek, found := opt.TargetSeeks[stsID]
			if found {
				seek = subSeek
			}
			dl := b.getTargetLogsSeek(stsID, t, seek)
			initialDelta.DeltaLogs = append(initialDelta.DeltaLogs, dl)
		}
	} else {
		for stsID, subSeek := range opt.TargetSeeks {
			t, found := b.targets[stsID]
			if !found {
				// Subscriber knows of a target we don't (knows the future).
				continue
			}
			dl := b.getTargetLogsSeek(stsID, t, subSeek)
			initialDelta.DeltaLogs = append(initialDelta.DeltaLogs, dl)
		}
	}
	err := sub.emit(initialDelta)
	if err != nil {
		return "", err
	}
	return subID, nil
}

func (b *Build) getTargetLogsSeek(targetID string, t *target, seek int64) *api.DeltaLog {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if seek < 0 {
		seek = t.endSeekIndex + seek
	}
	if seek < t.startSeekIndex {
		seek = t.startSeekIndex
	}
	return &api.DeltaLog{
		TargetID:  targetID,
		SeekIndex: seek,
		Data:      t.logData[seek-t.startSeekIndex:],
	}
}

// Resubscribe causes an existing subscriber reset their orderID for a number
// of streams. This is useful in the case of a client reconnecting.
func (b *Build) Resubscribe(subID string, opt *SubscriberOpt) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	sub, found := b.subs[subID]
	if !found {
		return errors.Errorf("subscriber %s not found", subID)
	}
	sub.mu.Lock()
	defer sub.mu.Unlock()
	// TODO
	panic("not implemented")
}

// Unsubscribe removes a subscriber.
func (b *Build) Unsubscribe(subID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs, subID)
}

// Snapshot creates a snapshot of the current build state and stores it to
// persistent storage.
func (b *Build) Snapshot(ctx context.Context) error {
	s := b.makeSnapshot()
	err := b.snapshotter.Write(ctx, s)
	if err != nil {
		return err
	}
	err = b.dropPreSnapshotData(s)
	if err != nil {
		return err
	}
	return nil
}

func (b *Build) LatestSnapshotID(ctx context.Context) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.latestSnapshotID, nil
}

func (b *Build) makeSnapshot() *snapshot.Snapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()
	s := &snapshot.Snapshot{
		SnapshotID:      fmt.Sprintf("snap-%d", b.manifestOrderID),
		Manifest:        new(api.Manifest),
		ManifestOrderID: b.manifestOrderID,
		Logs:            make(map[string]*snapshot.Target),
	}
	s.Manifest.CopyFrom(&b.manifest)
	for stsID, t := range b.targets {
		t.mu.RLock()
		mt, found := s.Manifest.Targets[stsID]
		if !found {
			mt = new(api.TargetManifest)
			s.Manifest.Targets[stsID] = mt
		}
		mt.Size = new(int64)
		*mt.Size = t.endSeekIndex
		s.Logs[stsID] = &snapshot.Target{
			LogData: append([]byte{}, t.logData...),
		}
		t.mu.RUnlock()
	}
	return s
}

func (b *Build) dropPreSnapshotData(s *snapshot.Snapshot) error {
	// TODO: We don't drop the logs yet, because they are needed in other
	//       places (e.g. snapshot logic).
	b.mu.Lock()
	defer b.mu.Unlock()
	b.latestSnapshotID = s.SnapshotID
	return nil
}

func (b *Build) receiveDeltaLog(dl *api.DeltaLog) error {
	b.mu.Lock()
	t, found := b.targets[dl.TargetID]
	if !found {
		t = new(target)
		t.stream = NewLogDeltaStream(func(dl *api.DeltaLog) error {
			return b.handleDeltaLog(t, dl)
		}, nil)
		b.targets[dl.TargetID] = t
	}
	b.mu.Unlock()
	return t.stream.Receive(dl)
}

func (b *Build) handleDeltaManifest(dm *api.DeltaManifest) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.manifestOrderID = dm.OrderID
	err := b.manifest.ApplyDelta(dm)
	if err != nil {
		return err
	}
	relevantSubs := make(map[string]*subscriber)
	for subID, sub := range b.subs {
		sub.mu.RLock()
		if !sub.manifest {
			sub.mu.RUnlock()
			continue
		}
		relevantSubs[subID] = sub
		sub.mu.RUnlock()
	}
	for _, sub := range relevantSubs {
		err := sub.emit(api.Delta{
			Version:        api.VersionNumber,
			DeltaManifests: []*api.DeltaManifest{dm},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Build) handleDeltaLog(t *target, dl *api.DeltaLog) error {
	b.mu.RLock()
	relevantSubs := make(map[string]*subscriber)
	for subID, sub := range b.subs {
		sub.mu.RLock()
		_, found := sub.targets[dl.TargetID]
		if !sub.allTargets && !found {
			sub.mu.RUnlock()
			continue
		}
		relevantSubs[subID] = sub
		sub.mu.RUnlock()
	}
	b.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()
	t.endSeekIndex = dl.EndOrderID()
	t.logData = append(t.logData, dl.Data...)
	for _, sub := range relevantSubs {
		err := sub.emit(api.Delta{
			Version:   api.VersionNumber,
			DeltaLogs: []*api.DeltaLog{dl},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
