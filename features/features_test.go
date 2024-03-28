package features

import (
	"testing"
)

func TestMustParseVersion(t *testing.T) {
	testCases := []struct {
		version  string
		expected []int
	}{
		{
			version:  "0.5",
			expected: []int{0, 5},
		},
		{
			version:  "0.67",
			expected: []int{0, 67},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.version, func(t *testing.T) {
			major, minor := mustParseVersion(tc.version)
			Equal(t, tc.expected[0], major)
			Equal(t, tc.expected[1], minor)
		})
	}
}
