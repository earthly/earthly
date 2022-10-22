package bus

import (
	"github.com/earthly/cloud-api/logstream"
)

// Printer is a build log printer.
type Printer struct {
	b        *Bus
	targetID string
	cached   bool
	push     bool
	local    bool
}

// NewPrinter creates a new Printer.
func NewPrinter(b *Bus) Printer {
	return Printer{
		b: b,
	}
}

func (p Printer) clone() Printer {
	return Printer{
		b:        p.b,
		targetID: p.targetID,
		cached:   p.cached,
		push:     p.push,
		local:    p.local,
	}
}

// PrintBytes prints a byte slice.
func (p Printer) PrintBytes(data []byte) {
	p.b.RawDelta(&logstream.Delta{
		DeltaLogs: []*logstream.DeltaLog{
			{
				TargetId: p.targetID,
				DeltaLogOneof: &logstream.DeltaLog_Data{
					Data: data,
				},
			},
		},
	})
}

// WithTargetID sets the target for the printer.
func (p Printer) WithTargetID(targetID string) Printer {
	p2 := p.clone()
	p2.targetID = targetID
	return p2
}

// AddCommand adds a new command and sets it for the printer.
func (p Printer) AddCommand(command string, cached bool, push bool, local bool) {
	p.b.RawDelta(&logstream.Delta{
		DeltaManifests: []*logstream.DeltaManifest{
			{
				DeltaManifestOneof: &logstream.DeltaManifest_Fields{
					Fields: &logstream.DeltaManifest_FieldsDelta{
						Targets: map[string]*logstream.DeltaTargetManifest{
							p.targetID: {
								Commands: map[int32]*logstream.DeltaCommandManifest{
									0: {
										Name: command,
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

// WithLocal sets the local flag for the printer.
func (p Printer) WithLocal(local bool) Printer {
	p2 := p.clone()
	p2.local = local
	return p2
}

// WithCached sets the cached flag for the printer.
func (p Printer) WithCached(cached bool) Printer {
	p2 := p.clone()
	p2.cached = cached
	return p2
}
