package server

import (
	"context"
	"testing"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/server/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestBuildSnapshot(t *testing.T) {
	ctx := context.Background()
	ms := snapshot.NewMemSnapshotter()
	b := NewBuild(ms)
	err := b.ReceiveDelta(api.Delta{
		DeltaManifests: []*api.DeltaManifest{
			{
				Reset: &api.Manifest{
					Version: api.VersionNumber,
				},
				OrderID: 0,
			},
		},
	})
	assert.NoError(t, err)
	targetID := "target1"
	name := "+some-target"
	b.ReceiveDelta(api.Delta{
		DeltaManifests: []*api.DeltaManifest{
			{
				Targets: map[string]*api.DeltaTargetManifest{
					targetID: {
						Name: name,
					},
				},
				OrderID: 1,
			},
		},
	})
	data := [][]byte{
		[]byte("hello"),
		[]byte(" world"),
		[]byte(" and the entire"),
		[]byte(" universe!"),
	}
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: 0,
				Data:      data[0],
			},
		},
	})
	assert.NoError(t, err)
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])),
				Data:      data[1],
			},
		},
	})
	assert.NoError(t, err)
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])) + int64(len(data[1])),
				Data:      data[2],
			},
		},
	})
	assert.NoError(t, err)
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])) + int64(len(data[1])) + int64(len(data[2])),
				Data:      data[3],
			},
		},
	})
	assert.NoError(t, err)
	err = b.Snapshot(ctx)
	assert.NoError(t, err)

	snapshotID, err := b.LatestSnapshotID(ctx)
	assert.NoError(t, err)
	m, err := ms.ReadManifest(ctx, snapshotID)
	assert.NoError(t, err)
	assert.Equal(t, name, m.Targets[targetID].Name)

	dt, _, _, err := ms.ReadLogFragment(ctx, snapshotID, targetID, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, "hello world and the entire universe!", string(dt))
}

func TestBuildSubscriber(t *testing.T) {
	ms := snapshot.NewMemSnapshotter()
	b := NewBuild(ms)
	err := b.ReceiveDelta(api.Delta{
		DeltaManifests: []*api.DeltaManifest{
			{
				Reset: &api.Manifest{
					Version: api.VersionNumber,
				},
				OrderID: 0,
			},
		},
	})
	assert.NoError(t, err)
	targetID := "target1"
	name := "+some-target"
	b.ReceiveDelta(api.Delta{
		DeltaManifests: []*api.DeltaManifest{
			{
				Targets: map[string]*api.DeltaTargetManifest{
					targetID: {
						Name: name,
					},
				},
				OrderID: 1,
			},
		},
	})

	recvM0 := &api.Manifest{}
	recvLog0 := []byte{}
	emitFun0 := func(d api.Delta) error {
		for _, dm := range d.DeltaManifests {
			assert.NotNil(t, dm.Reset)
			recvM0.CopyFrom(dm.Reset)
		}
		for _, dl := range d.DeltaLogs {
			assert.Equal(t, targetID, dl.TargetID)
			recvLog0 = append(recvLog0, dl.Data...)
		}
		return nil
	}
	_, err = b.Subscribe(&SubscriberOpt{
		Manifest:   true,
		AllTargets: true,
		Emit:       emitFun0,
	})
	assert.NoError(t, err)
	assert.Equal(t, name, recvM0.Targets[targetID].Name)
	assert.Equal(t, []byte(""), recvLog0)

	data := [][]byte{
		[]byte("hello"),
		[]byte(" world"),
		[]byte(" and the entire"),
		[]byte(" universe!"),
	}
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: 0,
				Data:      data[0],
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello"), recvLog0)
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])),
				Data:      data[1],
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), recvLog0)

	recvM1 := &api.Manifest{}
	recvLog1 := []byte{}
	emitFun1 := func(d api.Delta) error {
		for _, dm := range d.DeltaManifests {
			assert.NotNil(t, dm.Reset)
			recvM1.CopyFrom(dm.Reset)
		}
		for _, dl := range d.DeltaLogs {
			assert.Equal(t, targetID, dl.TargetID)
			recvLog1 = append(recvLog1, dl.Data...)
		}
		return nil
	}
	_, err = b.Subscribe(&SubscriberOpt{
		Manifest:    true,
		TargetSeeks: map[string]int64{targetID: 0},
		Emit:        emitFun1,
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), recvLog1)
	assert.Equal(t, name, recvM1.Targets[targetID].Name)

	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])) + int64(len(data[1])),
				Data:      data[2],
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world and the entire"), recvLog0)
	assert.Equal(t, []byte("hello world and the entire"), recvLog1)
	err = b.ReceiveDelta(api.Delta{
		DeltaLogs: []*api.DeltaLog{
			{
				TargetID:  targetID,
				SeekIndex: int64(len(data[0])) + int64(len(data[1])) + int64(len(data[2])),
				Data:      data[3],
			},
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, "hello world and the entire universe!", string(recvLog0))
	assert.Equal(t, "hello world and the entire universe!", string(recvLog1))
}
