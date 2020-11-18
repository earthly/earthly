package buildcontext

import (
	"testing"
)

func TestParseGitRemoteURL(t *testing.T) {
	//TODO FIX THIS var tests = []struct {
	//TODO FIX THIS 	gitURL          string
	//TODO FIX THIS 	expectedVendor  string
	//TODO FIX THIS 	expectedProject string
	//TODO FIX THIS 	valid           bool
	//TODO FIX THIS }{
	//TODO FIX THIS 	{
	//TODO FIX THIS 		"github.com:user/repo",
	//TODO FIX THIS 		"github.com",
	//TODO FIX THIS 		"user/repo",
	//TODO FIX THIS 		true,
	//TODO FIX THIS 	},
	//TODO FIX THIS 	{
	//TODO FIX THIS 		"git@github.com:user/repo.git",
	//TODO FIX THIS 		"github.com",
	//TODO FIX THIS 		"user/repo",
	//TODO FIX THIS 		true,
	//TODO FIX THIS 	},
	//TODO FIX THIS 	{
	//TODO FIX THIS 		"git@gitlab.com:user/repo.git",
	//TODO FIX THIS 		"gitlab.com",
	//TODO FIX THIS 		"user/repo",
	//TODO FIX THIS 		true,
	//TODO FIX THIS 	},
	//TODO FIX THIS 	{
	//TODO FIX THIS 		"ssh://git@github.com/earthly/earthly.git",
	//TODO FIX THIS 		"github.com",
	//TODO FIX THIS 		"earthly/earthly",
	//TODO FIX THIS 		true,
	//TODO FIX THIS 	},
	//TODO FIX THIS 	{
	//TODO FIX THIS 		"https://git@github.com/earthly/earthly.git",
	//TODO FIX THIS 		"github.com",
	//TODO FIX THIS 		"earthly/earthly",
	//TODO FIX THIS 		true,
	//TODO FIX THIS 	},
	//TODO FIX THIS }
	//TODO FIX THIS for _, test := range tests {
	//TODO FIX THIS 	vendor, project, err := parseGitRemoteURL(test.gitURL)
	//TODO FIX THIS 	if !test.valid {
	//TODO FIX THIS 		if err == nil {
	//TODO FIX THIS 			t.Errorf("expected error did not occur")
	//TODO FIX THIS 		}
	//TODO FIX THIS 		continue
	//TODO FIX THIS 	}
	//TODO FIX THIS 	if err != nil {
	//TODO FIX THIS 		t.Errorf("got err: %v", err)
	//TODO FIX THIS 	}
	//TODO FIX THIS 	if vendor != test.expectedVendor {
	//TODO FIX THIS 		t.Errorf("want: %v; got %v", test.expectedVendor, vendor)
	//TODO FIX THIS 	}
	//TODO FIX THIS 	if project != test.expectedProject {
	//TODO FIX THIS 		t.Errorf("want: %v; got %v", test.expectedProject, project)
	//TODO FIX THIS 	}
	//TODO FIX THIS }
}
