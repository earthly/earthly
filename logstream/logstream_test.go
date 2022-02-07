package logstream

import (
	"context"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/chanutil"
	"github.com/earthly/earthly/logstream/consumer"
	"github.com/earthly/earthly/logstream/producer"
	"github.com/earthly/earthly/logstream/server"
	"github.com/earthly/earthly/logstream/server/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	targetID := "target-1"
	otherTargetID := "target-2"
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	w2 := p.GetWriter(otherTargetID)
	w2.Write([]byte("SHOULD NOT BE PASSED ALONG"))
	pCh := p.DeltaCh()

	splitCh := chanutil.Splitter(ctx, pCh, 2, 1000)
	pChLocal := splitCh[0]
	pChRemote := splitCh[1]
	rLocal := consumer.NewTargetConsumerReader(ctx, targetID, pChLocal)
	doneStdoutRead := make(chan struct{})
	go func() {
		testStdout, err := ioutil.ReadAll(rLocal)
		assert.NoError(t, err)
		assert.Equal(t, "hello world and the entire universe!", string(testStdout))
		close(doneStdoutRead)
	}()

	batchedRemoteCh := chanutil.IntervalBatcher(ctx, 25*time.Millisecond, pChRemote)

	memSnapshotter := snapshot.NewMemSnapshotter()
	sb := server.NewBuild(memSnapshotter)
	go func() {
		for delta := range batchedRemoteCh {
			err := sb.ReceiveDelta(delta)
			assert.NoError(t, err)
		}
	}()
	time.Sleep(1 * time.Second)
	w.Write([]byte(" and the"))
	time.Sleep(1 * time.Second)

	err := sb.Snapshot(ctx)
	assert.NoError(t, err)
	snapshotID, err := sb.LatestSnapshotID(ctx)
	assert.NoError(t, err)
	dt, _, _, err := memSnapshotter.ReadLogFragment(ctx, snapshotID, targetID, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, "hello world and the", string(dt))

	uiDeltas := make(chan api.Delta, 1000)
	uiDoneCh := make(chan struct{})
	var uiLogData []byte
	go func() {
		uiReader := consumer.NewTargetConsumerReader(ctx, targetID, uiDeltas)
		tailOutput, err := ioutil.ReadAll(uiReader)
		assert.NoError(t, err)
		uiLogData = append(uiLogData, tailOutput...)
		assert.Equal(t, "the entire universe!", string(uiLogData))
		close(uiDoneCh)
	}()
	var onDataMu sync.Mutex
	onData := func(delta api.Delta) error {
		for _, dm := range delta.DeltaManifests {
			if (dm.Reset != nil && dm.Reset.Status == api.StatusSuccess) ||
				dm.Status == api.StatusSuccess {
				go func() {
					time.Sleep(1 * time.Second)
					onDataMu.Lock()
					defer onDataMu.Unlock()
					close(uiDeltas)
				}()
			}
		}
		onDataMu.Lock()
		defer onDataMu.Unlock()
		uiDeltas <- delta
		return nil
	}

	w.Write([]byte(" entire"))
	time.Sleep(time.Second)

	// A mock UI subscriber gets the last N=3 bytes of the snapshot and
	// subscribes from that point onwards to get the rest.
	snapshotID, err = sb.LatestSnapshotID(ctx)
	assert.NoError(t, err)
	var startSeek, endSeek int64
	uiLogData, startSeek, endSeek, err = memSnapshotter.ReadLogFragment(ctx, snapshotID, targetID, -3, 0)
	assert.NoError(t, err)
	assert.Equal(t, "the", string(uiLogData))
	assert.Equal(t, int64(len("hello world and ")), startSeek)
	assert.Equal(t, int64(len("hello world and the")), endSeek)
	_, err = sb.Subscribe(&server.SubscriberOpt{
		Manifest: true,
		TargetSeeks: map[string]int64{
			targetID: endSeek,
		},
		Emit: onData,
	})
	assert.NoError(t, err)

	w.Write([]byte(" universe!"))
	p.SetBuildMetadata(&api.DeltaManifest{
		Status: api.StatusSuccess,
	})
	err = p.Close()
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	err = sb.Snapshot(ctx)
	assert.NoError(t, err)
	snapshotID, err = sb.LatestSnapshotID(ctx)
	assert.NoError(t, err)
	m, err := memSnapshotter.ReadManifest(ctx, snapshotID)
	assert.NoError(t, err)
	assert.Equal(t, api.StatusSuccess, m.Status)
	dt, _, _, err = memSnapshotter.ReadLogFragment(ctx, snapshotID, targetID, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, "hello world and the entire universe!", string(dt))

	<-doneStdoutRead
	<-uiDoneCh
}
