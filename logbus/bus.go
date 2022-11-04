package logbus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/google/uuid"
)

// Subscriber is an object that can receive deltas.
type Subscriber interface {
	Write(*logstream.Delta)
}

// Bus is a build log data bus.
// It listens for raw deltas via WriteDeltaManifest and WriteRawLog, passes
// them on to raw subscribers, who can then write formatted deltas back to the
// bus via WriteFormattedLog. The formatted deltas are then passed on to
// formatted subscribers.
type Bus struct {
	run *Run

	rawMu      sync.Mutex
	rawSubs    []Subscriber
	rawHistory []*logstream.Delta

	formattedMu      sync.Mutex
	formattedSubs    []Subscriber
	formattedHistory []*logstream.Delta
}

// New creates a new Bus.
func New() *Bus {
	b := &Bus{
		run: nil, // set below
	}
	b.run = newRun(b)
	// TODO (vladaionescu): This should be issued somewhere else
	//                      (after we've parsed org and project names).
	b.WriteDeltaManifest(&logstream.DeltaManifest{
		DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{
			ResetAll: &logstream.RunManifest{
				BuildId:            uuid.NewString(),
				Version:            deltautil.Version,
				CreatedAtUnixNanos: uint64(time.Now().UnixNano()),
				OrgName:            "TODO",
				ProjectName:        "TODO",
			},
		},
	})
	return b
}

// AddSubscriber adds a subscriber to the bus. A subscriber receives both the
// raw and formatted deltas.
func (b *Bus) AddSubscriber(sub Subscriber) {
	b.AddRawSubscriber(sub)
	b.AddFormattedSubscriber(sub)
}

// AddRawSubscriber adds a raw subscriber to the bus. A raw subscriber only
// receives the raw deltas: DeltaManifest and DeltaLog.
func (b *Bus) AddRawSubscriber(sub Subscriber) {
	b.rawMu.Lock()
	defer b.rawMu.Unlock()
	for _, delta := range b.rawHistory {
		sub.Write(delta)
	}
	b.rawSubs = append(b.rawSubs, sub)
}

// AddFormattedSubscriber adds a formatted subscriber to the bus. A formatted
// subscriber receives only the formatted deltas: DeltaFormattedLog.
func (b *Bus) AddFormattedSubscriber(sub Subscriber) {
	b.formattedMu.Lock()
	defer b.formattedMu.Unlock()
	b.formattedSubs = append(b.formattedSubs, sub)
	for _, delta := range b.formattedHistory {
		sub.Write(delta)
	}
}

// Run returns the underlying run.
func (b *Bus) Run() *Run {
	return b.run
}

// WriteLog write a raw delta log to the bus.
func (b *Bus) WriteDeltaManifest(dm *logstream.DeltaManifest) {
	delta := &logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: dm,
		},
	}
	b.rawMu.Lock()
	defer b.rawMu.Unlock()
	b.rawHistory = append(b.rawHistory, delta)
	for _, sub := range b.rawSubs {
		sub.Write(delta)
	}
}

// WriteRawLog write a raw delta log to the bus.
func (b *Bus) WriteRawLog(dl *logstream.DeltaLog) {
	delta := &logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaLog{
			DeltaLog: dl,
		},
	}
	b.rawMu.Lock()
	defer b.rawMu.Unlock()
	b.rawHistory = append(b.rawHistory, delta)
	for _, sub := range b.rawSubs {
		sub.Write(delta)
	}
}

// WriteFormattedLog writes a formatted delta to the bus.
func (b *Bus) WriteFormattedLog(dfl *logstream.DeltaFormattedLog) {
	delta := &logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaFormattedLog{
			DeltaFormattedLog: dfl,
		},
	}
	b.formattedMu.Lock()
	defer b.formattedMu.Unlock()
	b.formattedHistory = append(b.formattedHistory, delta)
	for _, sub := range b.formattedSubs {
		sub.Write(delta)
	}
}

// FormattedWriter returns a writer that writes formatted deltas to the bus.
func (b *Bus) FormattedWriter(targetID string) *FormattedWriter {
	return NewFormattedWriter(b, targetID)
}
