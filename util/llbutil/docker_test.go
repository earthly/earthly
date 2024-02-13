package llbutil

import (
	"testing"
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
		ans := DockerTagSafe(tt.tag)
		if ans != tt.safe {
			t.Errorf("got %s, want %s", ans, tt.safe)
		}
	}
}
