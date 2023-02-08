package analytics

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

type testMutex struct {
	lockCalled   bool
	unlockCalled bool
	callback     func()
}

func (tm *testMutex) Lock() {
	tm.lockCalled = true
	if tm.callback != nil {
		tm.callback()
	}
}

func (tm *testMutex) Unlock() {
	tm.unlockCalled = true
}

// TestProjectTracker_AddEarthfileProject ensures the org and project are set correctly and that a lock is being acquired
func TestProjectTracker_AddEarthfileProject(t *testing.T) {
	tm := &testMutex{}
	pt := &ProjectTracker{
		mutex: tm,
	}
	org := "some org"
	project := "some project"
	pt.AddEarthfileProject(org, project)
	True(t, tm.lockCalled)
	True(t, tm.unlockCalled)
	Equal(t, org, pt.earthfileOrg)
	Equal(t, project, pt.earthfileProject)
}

func TestProjectTracker_AddCommandLineProject(t *testing.T) {
	tm := &testMutex{}
	pt := &ProjectTracker{
		mutex: tm,
	}
	org := "some org"
	project := "some project"
	pt.AddCommandLineProject(org, project)
	False(t, tm.lockCalled)
	False(t, tm.unlockCalled)
	Equal(t, org, pt.commandLineOrg)
	Equal(t, project, pt.commandLineProject)
}

func TestProjectTracker_ProjectDetails(t *testing.T) {
	earthfileOrg := "earthfile org"
	earthfileProject := "earthfile org"
	commandLineOrg := "command line org"
	commandLineProject := "command line org"

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
		"use org from commandline after lock when it's initially not set": {
			earthfileOrg:            "",
			earthfileProject:        earthfileProject,
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             commandLineOrg,
			expectedProject:         earthfileProject,
		},
		"use project from commandline after lock when it's initially not set": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         commandLineProject,
		},
		"use org from earthfile and project from commandline after lock when initially the former is set and the latter is not": {
			earthfileOrg:            earthfileOrg,
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     earthfileOrg,
			updatedEarthfileProject: "",
			expectedOrg:             earthfileOrg,
			expectedProject:         commandLineProject,
		},
		"use org from commandline and project from earthfile after lock when initially the former is not set and the latter is": {
			earthfileOrg:            "",
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: earthfileProject,
			expectedOrg:             commandLineOrg,
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
		"use org & project from commandline after lock when initially they are both not set": {
			earthfileOrg:            "",
			earthfileProject:        "",
			expectLockUnlock:        true,
			updatedEarthfileOrg:     "",
			updatedEarthfileProject: "",
			expectedOrg:             commandLineOrg,
			expectedProject:         commandLineProject,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			tm := &testMutex{}
			pt := &ProjectTracker{
				earthfileOrg:       tc.earthfileOrg,
				earthfileProject:   tc.earthfileProject,
				commandLineOrg:     commandLineOrg,
				commandLineProject: commandLineProject,
				mutex:              tm,
			}

			tm.callback = func() {
				if tc.updatedEarthfileOrg != "" {
					pt.earthfileOrg = tc.updatedEarthfileOrg
				}
				if tc.updatedEarthfileProject != "" {
					pt.earthfileProject = tc.updatedEarthfileProject
				}
			}

			org, project := pt.ProjectDetails()
			Equal(t, tc.expectLockUnlock, tm.lockCalled)
			Equal(t, tc.expectLockUnlock, tm.unlockCalled)
			Equal(t, tc.expectedOrg, org)
			Equal(t, tc.expectedProject, project)
		})
	}
}
