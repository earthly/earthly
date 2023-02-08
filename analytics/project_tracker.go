package analytics

import "sync"

type ProjectTracker struct {
	earthfileOrg     string
	earthfileProject string
	cliOrg           string
	cliProject       string
	mutex            sync.Locker
}

var projectTracker = ProjectTracker{
	mutex: &sync.Mutex{},
}

func (pt *ProjectTracker) AddEarthfileProject(org, project string) {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()
	pt.earthfileOrg = org
	pt.earthfileProject = project
}

func (pt *ProjectTracker) AddCLIProject(org, project string) {
	pt.cliOrg = org
	pt.cliProject = project
}

func (pt *ProjectTracker) ProjectDetails() (string, string) {
	if pt.earthfileOrg != "" && pt.earthfileProject != "" {
		return pt.earthfileOrg, pt.earthfileProject
	}
	org := pt.cliOrg
	project := pt.cliProject
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
