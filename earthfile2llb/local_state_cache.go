package earthfile2llb

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/earthly/earthly/util/inodeutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/llbutil/pllb"
)

// LocalStateCache provides caching of local States
type LocalStateCache struct {
	mu    sync.Mutex
	cache map[string]pllb.State
}

// NewSharedLocalStateCache creates a new local state cache
func NewSharedLocalStateCache() *LocalStateCache {
	return &LocalStateCache{
		cache: map[string]pllb.State{},
	}
}

// getOrConstruct returns a cached pllb.State with the same shared key, or creates a new one
// if it doesn't exist.
func (lsc *LocalStateCache) getOrConstruct(factory llbfactory.Factory) pllb.State {
	localFactory, ok := factory.(*llbfactory.LocalFactory)
	if !ok {
		return factory.Construct()
	}

	lsc.mu.Lock()
	defer lsc.mu.Unlock()

	key := localFactory.GetSharedKey()

	if st, ok := lsc.cache[key]; ok {
		return st
	}

	st := factory.Construct()
	lsc.cache[key] = st
	return st
}

func getSharedKeyHintFromInclude(name string, incl []string) string {
	h := sha1.New()
	b := make([]byte, 8)

	addToHash := func(path string) {
		h.Write([]byte(path))
		inode := inodeutil.GetInodeBestEffort(path)
		binary.LittleEndian.PutUint64(b, inode)
		h.Write(b)
	}

	addToHash(name)
	for _, path := range incl {
		addToHash(path)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func createIncludePatterns(incl []string) []string {
	incl2 := []string{}
	for _, inc := range incl {
		if inc == "." {
			inc = "./*"
		} else if strings.HasSuffix(inc, "/.") {
			inc = inc[:len(inc)-1] + "*"
		}
		inc = quoteMeta(inc)
		incl2 = append(incl2, inc)
	}
	return incl2
}

func addIncludePathAndSharedKeyHint(factory llbfactory.Factory, src []string) llbfactory.Factory {
	localFactory, ok := factory.(*llbfactory.LocalFactory)
	if !ok {
		return factory
	}

	includePatterns := createIncludePatterns(src)
	sharedKey := getSharedKeyHintFromInclude(localFactory.GetName(), includePatterns)

	return localFactory.
		WithInclude(includePatterns).
		WithSharedKeyHint(sharedKey)
}
