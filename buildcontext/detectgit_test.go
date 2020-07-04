package buildcontext

import (
	"testing"
)

func TestParseGitRemoteURL(t *testing.T) {
	var tests = []struct {
		gitURL          string
		expectedVendor  string
		expectedProject string
		valid           bool
	}{
		{
			"git@github.com:user/repo.git",
			"github.com",
			"user/repo",
			true,
		},
		{
			"git@gitlab.com:user/repo.git",
			"gitlab.com",
			"user/repo",
			true,
		},
		{
			"ssh://git@github.com/earthly/earthly.git",
			"github.com",
			"earthly/earthly",
			true,
		},
		{
			"https://git@github.com/earthly/earthly.git",
			"github.com",
			"earthly/earthly",
			true,
		},
		{
			"http://google.com/hello/world",
			"",
			"",
			false,
		},
	}
	for _, test := range tests {
		vendor, project, err := parseGitRemoteURL(test.gitURL)
		if !test.valid {
			if err == nil {
				t.Errorf("expected error did not occur")
			}
			continue
		}
		if err != nil {
			t.Errorf("got err: %v", err)
		}
		if vendor != test.expectedVendor {
			t.Errorf("want: %v; got %v", test.expectedVendor, vendor)
		}
		if project != test.expectedProject {
			t.Errorf("want: %v; got %v", test.expectedProject, project)
		}
	}
}
