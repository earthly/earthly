package gatewaycrafter

import (
	"fmt"
	"sort"
	"sync"

	"github.com/earthly/earthly/util/dockerutil"
)

// ExportCoordinator is a thread-safe data-store used for coordinating the export
// of images, and artifacts (e.g. OnPull, OnImage, and Artifact summaries)
type ExportCoordinator struct {
	m                     sync.Mutex
	imageEntries          map[string]imageEntry
	localOutputSummary    []LocalOutputSummaryEntry
	artifactOutputSummary []ArtifactOutputSummaryEntry
	pushedImageSummary    []PushedImageSummaryEntry
	imgIndex              int
}

type imageEntry struct {
	manifest   *dockerutil.Manifest
	localImage string
}

// LocalOutputSummaryEntry contains a summary of output images
type LocalOutputSummaryEntry struct {
	Target    string
	DockerTag string
	Salt      string
}

// PushedImageSummaryEntry contains a summary of images which were pushed
type PushedImageSummaryEntry struct {
	Target    string
	DockerTag string
	Salt      string
	Pushed    bool
}

// ArtifactOutputSummaryEntry contains a summary of output artifacts
type ArtifactOutputSummaryEntry struct {
	Target string
	Path   string
	Salt   string
}

// NewExportCoordinator returns a new ExportCoordinator
func NewExportCoordinator() *ExportCoordinator {
	return &ExportCoordinator{
		imageEntries: map[string]imageEntry{},
	}
}

// GetImage fetches an existing entry from the data-store or returns false if none exists
func (ec *ExportCoordinator) GetImage(k string) (*dockerutil.Manifest, string, bool) {
	ec.m.Lock()
	defer ec.m.Unlock()
	v, ok := ec.imageEntries[k]
	return v.manifest, v.localImage, ok
}

// AddImage creates a new entry for the value under sessionID/<v'>-<uuid>
// Where v' is v without special chars
func (ec *ExportCoordinator) AddImage(sessionID, localImage string, manifest *dockerutil.Manifest) string {
	ec.m.Lock()
	defer ec.m.Unlock()
	k := fmt.Sprintf("sess-%s/pullping:img-%d", sessionID, ec.imgIndex)
	ec.imgIndex++
	ec.imageEntries[k] = imageEntry{
		localImage: localImage,
		manifest:   manifest,
	}
	return k
}

// AddArtifactSummary adds an entry of a local target and docker tag, which is used to output a summary text at the end of earthly execution
func (ec *ExportCoordinator) AddArtifactSummary(target, path, salt string) {
	ec.m.Lock()
	defer ec.m.Unlock()
	ec.artifactOutputSummary = append(ec.artifactOutputSummary, ArtifactOutputSummaryEntry{
		Target: target,
		Path:   path,
		Salt:   salt,
	})
}

// GetArtifactSummary returns a list of artifact summary entries, sorted by target name
func (ec *ExportCoordinator) GetArtifactSummary() []ArtifactOutputSummaryEntry {
	entries := []ArtifactOutputSummaryEntry{}

	ec.m.Lock()
	for _, x := range ec.artifactOutputSummary {
		entries = append(entries, x)
	}
	ec.m.Unlock()

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Target < entries[j].Target
	})
	return entries
}

// AddLocalOutputSummary adds an entry of a local target and docker tag, which is used to output a summary text at the end of earthly execution
func (ec *ExportCoordinator) AddLocalOutputSummary(target, dockerTag, salt string) {
	ec.m.Lock()
	defer ec.m.Unlock()
	ec.localOutputSummary = append(ec.localOutputSummary, LocalOutputSummaryEntry{
		Target:    target,
		DockerTag: dockerTag,
		Salt:      salt,
	})
}

// GetLocalOutputSummary returns a list of output summary entries, sorted by target name
func (ec *ExportCoordinator) GetLocalOutputSummary() []LocalOutputSummaryEntry {
	entries := []LocalOutputSummaryEntry{}

	ec.m.Lock()
	for _, x := range ec.localOutputSummary {
		entries = append(entries, x)
	}
	ec.m.Unlock()

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Target < entries[j].Target
	})
	return entries
}

// AddPushedImageSummary adds an entry of a pushed images, which is used to output a summary text at the end of earthly execution
func (ec *ExportCoordinator) AddPushedImageSummary(target, dockerTag, salt string, pushed bool) {
	ec.m.Lock()
	defer ec.m.Unlock()
	ec.pushedImageSummary = append(ec.pushedImageSummary, PushedImageSummaryEntry{
		Target:    target,
		DockerTag: dockerTag,
		Salt:      salt,
		Pushed:    pushed,
	})
}

// GetPushedImageSummary returns a list of pushed image summary entries, sorted by target name
func (ec *ExportCoordinator) GetPushedImageSummary() []PushedImageSummaryEntry {
	entries := []PushedImageSummaryEntry{}

	ec.m.Lock()
	for _, x := range ec.pushedImageSummary {
		entries = append(entries, x)
	}
	ec.m.Unlock()

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Target < entries[j].Target
	})
	return entries
}
