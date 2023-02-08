package analytics

import "sync"

type ProjectTracker struct {
	earthfileOrg       string
	earthfileProject   string
	commandLineOrg     string
	commandLineProject string
	mutex              sync.Locker
}

func NewProjectTracker() *ProjectTracker {
	return &ProjectTracker{
		mutex: &sync.Mutex{},
	}
}
func (pt *ProjectTracker) AddEarthfileProject(org, project string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.earthfileOrg = org
	pt.earthfileProject = project
}

func (pt *ProjectTracker) AddCommandLineProject(org, project string) {
	pt.commandLineOrg = org
	pt.commandLineProject = project
}

func (pt *ProjectTracker) ProjectDetails() (string, string) {
	if pt.earthfileOrg != "" && pt.earthfileProject != "" {
		return pt.earthfileOrg, pt.earthfileProject
	}
	org := pt.commandLineOrg
	project := pt.commandLineProject
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	if pt.earthfileOrg != "" {
		org = pt.earthfileOrg
	}
	if pt.earthfileProject != "" {
		project = pt.earthfileProject
	}
	return org, project
}
