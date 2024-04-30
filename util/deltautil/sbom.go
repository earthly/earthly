package deltautil

import "sync"

var (
	sbomMu      sync.Mutex
	sboms       []SBOMEntry
	byImageName map[string][]string
)

func init() {
	byImageName = map[string][]string{}
}

type SBOMEntry struct {
	Target string
	SBOM   string
}

func AddSBOM(target, s string) {
	sbomMu.Lock()
	defer sbomMu.Unlock()
	sboms = append(sboms, SBOMEntry{
		Target: target,
		SBOM:   s,
	})
}

func AddSBOMToSaveImage(imageName string, sboms []string) {
	sbomMu.Lock()
	defer sbomMu.Unlock()
	byImageName[imageName] = sboms
}

func GetSBOMSbyImageName(imageName string) []string {
	sbomMu.Lock()
	defer sbomMu.Unlock()

	sboms, ok := byImageName[imageName]
	if !ok {
		return []string{}
	}
	cpy := make([]string, len(sboms))
	copy(cpy, sboms)
	return cpy
}

func GetSBOMs() []SBOMEntry {
	sbomMu.Lock()
	defer sbomMu.Unlock()

	cpy := make([]SBOMEntry, len(sboms))
	copy(cpy, sboms)

	return cpy
}
