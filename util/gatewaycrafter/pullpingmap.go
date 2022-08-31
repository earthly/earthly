package gatewaycrafter

import (
	"fmt"
	"sync"

	"github.com/earthly/earthly/util/dockerutil"
	"github.com/earthly/earthly/util/stringutil"
)

// PullPingMap is a thread-save map used for coordinating pullpings
type PullPingMap struct {
	m       sync.Mutex
	entries map[string]pullPingMapEntry
}

type pullPingMapEntry struct {
	manifest   *dockerutil.Manifest
	localImage string
}

// NewPullPingMap returns a new PullPingMap
func NewPullPingMap() *PullPingMap {
	return &PullPingMap{
		entries: map[string]pullPingMapEntry{},
	}
}

// Get fetches an existing entry from the map or returns false if none exists
func (ppm *PullPingMap) Get(k string) (*dockerutil.Manifest, string, bool) {
	ppm.m.Lock()
	defer ppm.m.Unlock()
	v, ok := ppm.entries[k]
	return v.manifest, v.localImage, ok
}

// Insert creates a new entry for the value under sessionID/<v'>-<uuid>
// Where v' is v without special chars
func (ppm *PullPingMap) Insert(sessionID, localImage string, manifest *dockerutil.Manifest) string {
	k := fmt.Sprintf("sess-%s/pullping:%s-%s", sessionID, stringutil.AlphanumericOnly(localImage), stringutil.RandomAlphanumeric(32))
	ppm.m.Lock()
	defer ppm.m.Unlock()
	ppm.entries[k] = pullPingMapEntry{
		localImage: localImage,
		manifest:   manifest,
	}
	return k
}
