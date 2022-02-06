package producer

import (
	"io"

	"github.com/earthly/earthly/logstream/api"
)

// Producer is a producer of deltas.
type Producer struct {
	targets map[string]*target
	deltaCh chan api.Delta
}

type target struct {
	w *logWriter
}

// NewProducer returns a new producer.
func NewProducer() *Producer {
	return &Producer{
		targets: make(map[string]*target),
		deltaCh: make(chan api.Delta, 1000),
	}
}

// GetWriter returns a writer for a given target.
func (p *Producer) GetWriter(targetID string) io.Writer {
	return p.getTarget(targetID).w
}

// SetBuildMetadata sets the build metadata for the given target.
func (p *Producer) SetBuildMetadata(dm *api.DeltaManifest) {
	p.deltaCh <- api.Delta{
		Version:        api.VersionNumber,
		DeltaManifests: []*api.DeltaManifest{dm},
	}
}

// SetTargetMetadata sets the target metadata for the given target.
func (p *Producer) SetTargetMetadata(targetID string, dtm *api.DeltaTargetManifest) {
	p.deltaCh <- api.Delta{
		Version: api.VersionNumber,
		DeltaManifests: []*api.DeltaManifest{
			{
				Targets: map[string]*api.DeltaTargetManifest{
					targetID: dtm,
				},
			},
		},
	}
}

// SetCommandMetadata sets the command metadata for the given target.
func (p *Producer) SetCommandMetadata(targetID string, execOrder int, dcm *api.DeltaCommandManifest) {
	p.deltaCh <- api.Delta{
		Version: api.VersionNumber,
		DeltaManifests: []*api.DeltaManifest{
			{
				Targets: map[string]*api.DeltaTargetManifest{
					targetID: {
						Commands: map[int]*api.DeltaCommandManifest{
							execOrder: dcm,
						},
					},
				},
			},
		},
	}
}

func (p *Producer) getTarget(targetID string) *target {
	t, found := p.targets[targetID]
	if !found {
		t = &target{
			w: &logWriter{
				deltaCh:  p.deltaCh,
				targetID: targetID,
			},
		}
		p.targets[targetID] = t
	}
	return t
}

// DeltaCh returns the channel used to emit deltas to.
func (p *Producer) DeltaCh() chan api.Delta {
	return p.deltaCh
}

// Close closes the producer, signaling that there will be no more deltas.
func (p *Producer) Close() error {
	close(p.deltaCh)
	return nil
}
