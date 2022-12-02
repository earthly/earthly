package earthfile2llb

import (
	"sync"

	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/waitutil"
)

type saveImageWaitItem struct {
	c  *Converter
	si states.SaveImage

	push        bool
	localExport bool

	mu sync.Mutex
}

func newSaveImage(si states.SaveImage, c *Converter, push, localExport bool) waitutil.WaitItem {
	return &saveImageWaitItem{
		c:           c,
		si:          si,
		push:        push,
		localExport: localExport,
	}
}

func (siwi *saveImageWaitItem) SetDoSave() {
	siwi.mu.Lock()
	defer siwi.mu.Unlock()
	siwi.localExport = true
}

func (siwi *saveImageWaitItem) SetDoPush() {
	siwi.mu.Lock()
	defer siwi.mu.Unlock()
	siwi.push = true
}
