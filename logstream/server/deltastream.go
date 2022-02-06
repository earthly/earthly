package server

import (
	"github.com/earthly/earthly/logstream/api"
	"github.com/pkg/errors"
)

type applyFun func(delta orderedDelta) error
type applyDeltaManifestFun func(dm *api.DeltaManifest) error
type applyDeltaLogFun func(dl *api.DeltaLog) error

type orderedDelta interface {
	StartOrderID() int64
	EndOrderID() int64
}

// DeltaStream is a processor that ensures that deltas are applied in the
// correct order.
type DeltaStream struct {
	nextOrderID  int64
	apply        applyFun
	onOutOfOrder func(int64)
}

// NewManifestDeltaStream returns a new DeltaStream that applies DeltaManifests
// taking care of deduplication and reordering.
func NewManifestDeltaStream(apply applyDeltaManifestFun, onOutOfOrder func(int64)) *DeltaStream {
	return &DeltaStream{
		apply: func(delta orderedDelta) error {
			d, ok := delta.(*api.DeltaManifest)
			if !ok {
				return errors.Errorf("unexpected delta type %T", delta)
			}
			return apply(d)
		},
		onOutOfOrder: onOutOfOrder,
	}
}

// NewLogDeltaStream returns a new DeltaStream that applies DeltaLogs
// taking care of deduplication and reordering.
func NewLogDeltaStream(apply applyDeltaLogFun, onOutOfOrder func(int64)) *DeltaStream {
	return &DeltaStream{
		apply: func(delta orderedDelta) error {
			d, ok := delta.(*api.DeltaLog)
			if !ok {
				return errors.Errorf("unexpected delta type %T", delta)
			}
			return apply(d)
		},
		onOutOfOrder: onOutOfOrder,
	}
}

// Receive handles an incoming delta.
func (ds *DeltaStream) Receive(delta orderedDelta) error {
	if delta.StartOrderID() == delta.EndOrderID() {
		return nil
	}
	if delta.StartOrderID() > delta.EndOrderID() {
		return errors.Errorf("invalid delta: startOrderID %d > endOrderID %d", delta.StartOrderID(), delta.EndOrderID())
	}
	if delta.EndOrderID() <= ds.nextOrderID {
		// Old (duplicate) delta. Ignore.
		if ds.onOutOfOrder != nil {
			ds.onOutOfOrder(ds.nextOrderID)
		}
		return nil
	} else if delta.StartOrderID() < ds.nextOrderID {
		// Partial overlap. Not yet implemented. Should not normally happen.
		return errors.Errorf("delta has partial overlap")
	} else if delta.StartOrderID() == ds.nextOrderID {
		ds.nextOrderID = delta.EndOrderID()
		err := ds.apply(delta)
		if err != nil {
			return err
		}
	}
	// Future delta. Ignore.
	if ds.onOutOfOrder != nil {
		ds.onOutOfOrder(ds.nextOrderID)
	}
	return nil
}
