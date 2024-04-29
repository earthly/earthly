package earthfile2llb

import "sync"

var (
	sbomMu sync.Mutex
	sboms  []SBOMEntry
)

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

func GetSBOMs() []SBOMEntry {
	sbomMu.Lock()
	defer sbomMu.Unlock()

	cpy := make([]SBOMEntry, len(sboms))
	copy(cpy, sboms)

	return cpy
}
