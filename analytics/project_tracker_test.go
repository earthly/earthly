package analytics

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

type mockMutex struct {
	lockCalled   bool
	unlockCalled bool
	callback     func()
}

func (mm *mockMutex) Lock() {
	mm.lockCalled = true
	if mm.callback != nil {
		mm.callback()
	}
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
	False(t, mm.lockCalled)
	False(t, mm.unlockCalled)
	Equal(t, org, pt.cliOrg)
	Equal(t, project, pt.cliProject)
}

func TestProjectTracker_ProjectDetails(t *testing.T) {
	earthfileOrg := "earthfile org"
	earthfileProject := "earthfile org"
	cliOrg := "cli org"
	cliProject := "cli org"

	testCases := map[string]struct {
		earthfileOrg            string
		earthfileProject        string
		expectLockUnlock        bool
		updatedEarthfileOrg     string
		updatedEarthfileProject string
		expectedOrg             string
		expectedProject         string
	}{
		"use details from earthfile when both set": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        earthfileProject,
			expectLockUnlock:        false,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         earthfileProject,
		},
		"use org from earthfile after lock when it's initially not set": {
			earthfileOrg:            "",
			earthfileProject:        earthfileProject,
			expectLockUnlock:        true,
			updatedEarthfileOrg:     earthfileOrg,
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         earthfileProject,
		},
		"use project from earthfile after lock when it's initially not set": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: earthfileProject,
			expectedOrg:             earthfileOrg,
			expectedProject:         earthfileProject,
		},
		"use org from cli after lock when it's initially not set": {
			earthfileOrg:            "",
			earthfileProject:        earthfileProject,
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             cliOrg,
			expectedProject:         earthfileProject,
		},
		"use project from cli after lock when it's initially not set": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         cliProject,
		},
		"use org from earthfile and project from cli after lock when initially the former is set and the latter is not": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     earthfileOrg,
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         cliProject,
		},
		"use org from cli and project from earthfile after lock when initially the former is not set and the latter is": {
			earthfileOrg:            "",
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: earthfileProject,
			expectedOrg:             cliOrg,
			expectedProject:         earthfileProject,
		},
		"use org & project from earthfile after lock when initially they are both not set": {
			earthfileOrg:            "",
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     earthfileOrg,
			updatedEarthfileProject: earthfileProject,
			expectedOrg:             earthfileOrg,
			expectedProject:         earthfileProject,
		},
		"use org & project from cli after lock when initially they are both not set": {
			earthfileOrg:            "",
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             cliOrg,
			expectedProject:         cliProject,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			mm := &mockMutex{}
			pt := &ProjectTracker{
				earthfileOrg:     tc.earthfileOrg,
				earthfileProject: tc.earthfileProject,
				cliOrg:           cliOrg,
				cliProject:       cliProject,
				mutex:            mm,
			}

			mm.callback = func() {
				if tc.updatedEarthfileOrg != "" {
					pt.earthfileOrg = tc.updatedEarthfileOrg
				}
				if tc.updatedEarthfileProject != "" {
					pt.earthfileProject = tc.updatedEarthfileProject
				}
			}

			org, project := pt.ProjectDetails()
			Equal(t, tc.expectLockUnlock, mm.lockCalled)
			Equal(t, tc.expectLockUnlock, mm.unlockCalled)
			Equal(t, tc.expectedOrg, org)
			Equal(t, tc.expectedProject, project)
		})
	}
}
