package gatewaycrafter

import "sync"

type onPull struct {
	m           sync.Mutex
	localImages map[string]string
}

var (
	OnPullInst *onPull
)

func init() {
	OnPullInst = newOnPull()
}

func newOnPull() *onPull {
	return &onPull{
		localImages: map[string]string{},
	}
}

func (op *onPull) Get(k string) (string, bool) {
	op.m.Lock()
	defer op.m.Unlock()
	v, ok := op.localImages[k]
	return v, ok
}

func (op *onPull) Set(k, v string) {
	op.m.Lock()
	defer op.m.Unlock()
	op.localImages[k] = v
}
