package gatewaycrafter

import (
	"fmt"
	"sync"

	"github.com/earthly/earthly/util/dockerutil"
	"github.com/earthly/earthly/util/stringutil"
)

// ExportCoordinator is a thread-save data-store used for coordinating the export
// of images, and artifacts (e.g. OnPull, OnImage, and Artifact summaries)
type ExportCoordinator struct {
	m            sync.Mutex
	imageEntries map[string]imageEntry
}

type imageEntry struct {
	manifest   *dockerutil.Manifest
	localImage string
}

// NewExportCoordinator returns a new ExportCoordinator
func NewExportCoordinator() *ExportCoordinator {
	return &ExportCoordinator{
		imageEntries: map[string]imageEntry{},
	}
}

// GetImage fetches an existing entry from the data-store or returns false if none exists
func (ppm *ExportCoordinator) GetImage(k string) (*dockerutil.Manifest, string, bool) {
	ppm.m.Lock()
	defer ppm.m.Unlock()
	v, ok := ppm.imageEntries[k]
	return v.manifest, v.localImage, ok
}

// AddImage creates a new entry for the value under sessionID/<v'>-<uuid>
// Where v' is v without special chars
func (ppm *ExportCoordinator) AddImage(sessionID, localImage string, manifest *dockerutil.Manifest) string {
	k := fmt.Sprintf("sess-%s/pullping:%s-%s", sessionID, stringutil.AlphanumericOnly(localImage), stringutil.RandomAlphanumeric(32))
	ppm.m.Lock()
	defer ppm.m.Unlock()
	ppm.imageEntries[k] = imageEntry{
		localImage: localImage,
		manifest:   manifest,
	}
	return k
}
