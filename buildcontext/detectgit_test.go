package buildcontext

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestParseGitRemoteURL(t *testing.T) {
	var tests = []struct {
		gitURL         string
		expectedGitURL string
		valid          bool
	}{
		{
			"github.com:user/repo",
			"github.com/user/repo",
			true,
		},
		{
			"git@github.com:user/repo.git",
			"github.com/user/repo",
			true,
		},
		{
			"git@gitlab.com:user/repo.git",
			"gitlab.com/user/repo",
			true,
		},
		{
			"ssh://git@github.com/earthly/earthly.git",
			"github.com/earthly/earthly",
			true,
		},
		{
			"https://git@github.com/earthly/earthly.git",
			"github.com/earthly/earthly",
			true,
		},
	}
	for _, test := range tests {
		gitURL, err := parseGitRemoteURL(test.gitURL)
		if !test.valid {
			if err == nil {
				t.Errorf("expected error did not occur")
			}
			continue
		}
		NoError(t, err, "parseGitRemoteURL failed")
		Equal(t, test.expectedGitURL, gitURL)
	}
}
