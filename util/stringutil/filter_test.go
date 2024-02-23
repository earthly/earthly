package stringutil_test

import (
	"testing"

	. "github.com/earthly/earthly/util/stringutil"

	"github.com/stretchr/testify/require"
)

func TestRemoveFromFromArgs(t *testing.T) {
	for _, testCase := range []struct {
		args     []string
		remove   []string
		expected []string
	}{
		{
			args:     []string{"a", "b", "c"},
			remove:   []string{"a"},
			expected: []string{"b", "c"},
		},
		{
			args:     []string{"a", "b", "c"},
			remove:   []string{"b", "c"},
			expected: []string{"a"},
		},
		{
			args:     []string{"a", "b", "c"},
			remove:   []string{},
			expected: []string{"a", "b", "c"},
		},
	} {
		actual := FilterElementsFromList(testCase.args, testCase.remove...)
		require.ElementsMatch(t, testCase.expected, actual)
	}
}
