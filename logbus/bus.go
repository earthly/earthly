package logbus

import (
	"sync"

	"github.com/earthly/cloud-api/logstream"
)

// SubscriberFun is a function that is called for each delta.
type SubscriberFun func(delta *logstream.Delta)

// Bus is a build log data bus.
// It listens for raw deltas via WriteDeltaManifest and WriteRawLog, passes
// them on to raw subscribers, who can then write formatted deltas back to the
// bus via WriteFormattedLog. The formatted deltas are then passed on to
// formatted subscribers.
type Bus struct {
	run *Run

	rawMu      sync.Mutex
	rawSubs    []SubscriberFun
	rawHistory []*logstream.Delta

	formattedMu      sync.Mutex
	formattedSubs    []SubscriberFun
	formattedHistory []*logstream.Delta
}

// New creates a new Bus.
func New() *Bus {
	b := &Bus{
		run: nil, // set below
	}
	b.run = newRun(b)
	return b
}

// AddSubscriber adds a subscriber to the bus. A subscriber receives both the
// raw and formatted deltas.
func (b *Bus) AddSubscriber(sub SubscriberFun) {
	b.AddRawSubscriber(sub)
	b.AddFormattedSubscriber(sub)
}

// AddRawSubscriber adds a raw subscriber to the bus. A raw subscriber only
// receives the raw deltas: DeltaManifest and DeltaLog.
func (b *Bus) AddRawSubscriber(sub SubscriberFun) {
	b.rawMu.Lock()
	defer b.rawMu.Unlock()
	for _, delta := range b.rawHistory {
		sub(delta)
	}
	b.rawSubs = append(b.rawSubs, sub)
}

// AddFormattedSubscriber adds a formatted subscriber to the bus. A formatted
// subscriber receives only the formatted deltas: DeltaFormattedLog.
func (b *Bus) AddFormattedSubscriber(sub SubscriberFun) {
	b.formattedMu.Lock()
	defer b.formattedMu.Unlock()
	b.formattedSubs = append(b.formattedSubs, sub)
	for _, delta := range b.formattedHistory {
		sub(delta)
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
		sub(delta)
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
		sub(delta)
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
		sub(delta)
	}
}

// FormattedWriter returns a writer that writes formatted deltas to the bus.
func (b *Bus) FormattedWriter(targetID string) *FormattedWriter {
	return NewFormattedWriter(b, targetID)
}
