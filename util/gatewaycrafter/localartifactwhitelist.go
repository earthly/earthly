package gatewaycrafter

import (
	"sync"
)

// LocalArtifactWhiteList is a set of paths which have been seen in a SAVE ARTIFACT ... AS LOCAL command
type LocalArtifactWhiteList struct {
	m     sync.Mutex
	paths map[string]struct{}
}

// NewLocalArtifactWhiteList returns a new LocalArtifactWhiteList
func NewLocalArtifactWhiteList() *LocalArtifactWhiteList {
	return &LocalArtifactWhiteList{
		paths: map[string]struct{}{},
	}
}

// Exists returns true if the path exists in the set
func (l *LocalArtifactWhiteList) Exists(k string) bool {
	l.m.Lock()
	defer l.m.Unlock()
	_, exists := l.paths[k]
	return exists
}

// Add adds the path to the set of paths
func (l *LocalArtifactWhiteList) Add(path string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.paths[path] = struct{}{}
}

// AsList returns a copy of the set as a list
func (l *LocalArtifactWhiteList) AsList() []string {
	l.m.Lock()
	defer l.m.Unlock()
	paths := []string{}
	for path := range l.paths {
		paths = append(paths, path)
	}
	return paths
}
