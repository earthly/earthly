package earthfile2llb

import (
	"sync"

	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/waitutil"
)

type saveImageWaitItem struct {
	c  *Converter
	si states.SaveImage

	allowPush   bool
	doPush      bool
	localExport bool

	mu sync.Mutex
}

func newSaveImage(si states.SaveImage, c *Converter, allowPush, localExport bool) waitutil.WaitItem {
	return &saveImageWaitItem{
		c:           c,
		si:          si,
		allowPush:   allowPush,
		localExport: localExport,
	}
}

func (siwi *saveImageWaitItem) SetDoSave() {
	siwi.mu.Lock()
	defer siwi.mu.Unlock()
	if siwi.si.DockerTag != "" {
		siwi.localExport = true
	}
}

func (siwi *saveImageWaitItem) SetDoPush() {
	siwi.mu.Lock()
	defer siwi.mu.Unlock()
	if siwi.si.DockerTag != "" {
		siwi.doPush = siwi.allowPush
	}
}
