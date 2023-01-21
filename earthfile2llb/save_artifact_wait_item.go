package earthfile2llb

import (
	"sync"

	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/waitutil"
)

type saveArtifactLocalWaitItem struct {
	c           *Converter
	saveLocal   states.SaveLocal
	localExport bool
	mu          sync.Mutex
}

// SetDoPush has no effect, but exists to satisfy interface
func (salwi *saveArtifactLocalWaitItem) SetDoPush() {
}

func (salwi *saveArtifactLocalWaitItem) SetDoSave() {
	salwi.mu.Lock()
	defer salwi.mu.Unlock()
	salwi.localExport = true
}

func newSaveArtifactLocal(state states.SaveLocal, c *Converter, localExport bool) waitutil.WaitItem {
	return &saveArtifactLocalWaitItem{
		c:           c,
		saveLocal:   state,
		localExport: localExport,
	}
}
