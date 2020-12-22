package variables

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestDockerTagSafe(t *testing.T) {
	var tests = []struct {
		tag  string
		safe string
	}{
		{"pull/123", "pull_123"},
		{"", "latest"},
		{"-asdf", "_asdf"},
		{"a-asdf", "a-asdf"},
		{"0123", "0123"},
		{"/a/b/c/d", "_a_b_c_d"},
		{"asdf:aa", "asdf_aa"},
		{"verylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongSHOULDTRUNCATE", "verylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylong"},
		{"a", "a"},
		{".", "_"},
		{"v1.2.3", "v1.2.3"},
		{"john/my-branch", "john_my-branch"},
	}

	for _, tt := range tests {
		ans := dockerTagSafe(tt.tag)
		Equal(t, tt.safe, ans)
	}
}

func TestGetProjectName(t *testing.T) {
	var tests = []struct {
		tag  string
		safe string
	}{
		{"http://github.com/earthly/earthly", "earthly/earthly"},
		{"http://gitlab.com/earthly/earthly", "earthly/earthly"},
		{"https://github.com/earthly/earthly", "earthly/earthly"},
		{"https://user@github.com/earthly/earthly", "earthly/earthly"},
		{"https://user:password@github.com/earthly/earthly", "earthly/earthly"},
		{"git@github.com:earthly/earthly", "earthly/earthly"},
		{"git@bitbucket.com:earthly/earthly", "earthly/earthly"},
		{"ssh://git@github.com/earthly/earthly", "earthly/earthly"},
		{"ssh://git@github.com:22/earthly/earthly", "earthly/earthly"},
	}

	for _, tt := range tests {
		ans := getProjectName(tt.tag)
		Equal(t, tt.safe, ans)
	}
}
