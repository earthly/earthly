package logbus

import (
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
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
	run       *Run
	createdAt time.Time

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
		run:       nil, // set below
		createdAt: time.Now(),
	}
	b.run = newRun(b)
	return b
}

// CreatedAt returns the time the bus was created.
func (b *Bus) CreatedAt() time.Time {
	return b.createdAt
}

// NowUnixNanos returns the current time in unix nanoseconds, ensuring
// monotonically increasing time.
func (b *Bus) NowUnixNanos() uint64 {
	return b.TsUnixNanos(time.Now())
}

// TsUnixNanos returns a given timestamp in unix nanoseconds, ensuring
// monotonically increasing time.
func (b *Bus) TsUnixNanos(t2 time.Time) uint64 {
	// The following is necessary to ensure that the time is monotonically increasing.
	// t2.UnixNano() sometimes strips the monotonic clock reading (e.g. time.Now().UnixNano()).
	// Sub maintains the monotonic clock reading https://pkg.go.dev/time#hdr-Monotonic_Clocks.
	deltaT := t2.Sub(b.createdAt).Nanoseconds()
	return uint64(b.CreatedAt().UnixNano()) + uint64(deltaT)
}

// AddSubscriber adds a subscriber to the bus. A subscriber receives both the
// raw and formatted deltas.
func (b *Bus) AddSubscriber(sub Subscriber) {
	b.AddRawSubscriber(sub)
	b.AddFormattedSubscriber(sub)
}

// RemoveSubscriber removes a subscriber from the bus.
func (b *Bus) RemoveSubscriber(sub Subscriber) {
	b.RemoveRawSubscriber(sub)
	b.RemoveFormattedSubscriber(sub)
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

// RemoveRawSubscriber removes a raw subscriber from the bus.
func (b *Bus) RemoveRawSubscriber(sub Subscriber) {
	b.rawMu.Lock()
	defer b.rawMu.Unlock()
	for i, s := range b.rawSubs {
		if s == sub {
			b.rawSubs = append(b.rawSubs[:i], b.rawSubs[i+1:]...)
			return
		}
	}
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

// RemoveFormattedSubscriber removes a formatted subscriber from the bus.
func (b *Bus) RemoveFormattedSubscriber(sub Subscriber) {
	b.formattedMu.Lock()
	defer b.formattedMu.Unlock()
	for i, s := range b.formattedSubs {
		if s == sub {
			b.formattedSubs = append(b.formattedSubs[:i], b.formattedSubs[i+1:]...)
			return
		}
	}
}

// Run returns the underlying run.
func (b *Bus) Run() *Run {
	return b.run
}

// WriteDeltaManifest write a raw delta log to the bus.
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
func (b *Bus) FormattedWriter(targetID, commandID string) *FormattedWriter {
	return NewFormattedWriter(b, targetID, commandID)
}
