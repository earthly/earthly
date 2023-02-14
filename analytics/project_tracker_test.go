package analytics

import (
	"testing"
)

type mockMutex struct {
	lockCalled   bool
	unlockCalled bool
}

func (mm *mockMutex) Lock() {
	mm.lockCalled = true
}

func (mm *mockMutex) Unlock() {
	mm.unlockCalled = true
}

// TestProjectTracker_AddEarthfileProject ensures the org and project are set correctly and that a lock is being acquired
func TestProjectTracker_AddEarthfileProject(t *testing.T) {
	mm := &mockMutex{}
	pt := &ProjectTracker{
		mutex: mm,
	}
	org := "some org"
	project := "some project"
	pt.AddEarthfileProject(org, project)
	True(t, mm.lockCalled)
	True(t, mm.unlockCalled)
	Equal(t, org, pt.earthfileOrg)
	Equal(t, project, pt.earthfileProject)
}

func TestProjectTracker_AddCLIProject(t *testing.T) {
	mm := &mockMutex{}
	pt := &ProjectTracker{
		mutex: mm,
	}
	org := "some org"
	project := "some project"
	pt.AddCLIProject(org, project)
	True(t, mm.lockCalled)
	True(t, mm.unlockCalled)
	Equal(t, org, pt.cliOrg)
	Equal(t, project, pt.cliProject)
}

func TestProjectTracker_ProjectDetails(t *testing.T) {
	earthfileOrg := "earthfile org"
	earthfileProject := "earthfile org"
	cliOrg := "cli org"
	cliProject := "cli org"

	testCases := map[string]struct {
		earthfileOrg     string
		earthfileProject string
		expectedOrg      string
		expectedProject  string
	}{
		"use details from earthfile when both org and project are set": {
			earthfileOrg:     earthfileOrg,
			earthfileProject: earthfileProject,
			expectedOrg:      earthfileOrg,
			expectedProject:  earthfileProject,
		},
		"use org from cli when it's initially not set in earthfile": {
			earthfileOrg:     "",
			earthfileProject: earthfileProject,
			expectedOrg:      cliOrg,
			expectedProject:  earthfileProject,
		},
		"use project from cli when it's not set in earthfile": {
			earthfileOrg:     earthfileOrg,
			earthfileProject: "",
			expectedOrg:      earthfileOrg,
			expectedProject:  cliProject,
		},
		"use org & project from cli after when they are both not set in earthfile": {
			earthfileOrg:     "",
			earthfileProject: "",
			expectedOrg:      cliOrg,
			expectedProject:  cliProject,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mm := &mockMutex{}
			pt := &ProjectTracker{
				earthfileOrg:     tc.earthfileOrg,
				earthfileProject: tc.earthfileProject,
				cliOrg:           cliOrg,
				cliProject:       cliProject,
				mutex:            mm,
			}

			org, project := pt.ProjectDetails()
			True(t, mm.lockCalled)
			True(t, mm.unlockCalled)
			Equal(t, tc.expectedOrg, org)
			Equal(t, tc.expectedProject, project)
		})
	}
}
