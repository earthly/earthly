package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedactSecretsFromArgs(t *testing.T) {
	for _, testCase := range []struct {
		args     []string
		expected []string
	}{
		{
			args:     []string{"earthly", "--secret", "foo=bar"},
			expected: []string{"earthly", "--secret", "foo=XXXXX"},
		},
		{
			args:     []string{"earthly", "--secret", "foo=bar", "--ci"},
			expected: []string{"earthly", "--secret", "foo=XXXXX", "--ci"},
		},
		{
			args:     []string{"earthly", "--secret", "foo", "--ci"},
			expected: []string{"earthly", "--secret", "foo", "--ci"},
		},
		{
			args:     []string{"earthly", "-s", "foo=bar"},
			expected: []string{"earthly", "-s", "foo=XXXXX"},
		},
		{
			args:     []string{"earthly", "-s", "foo=bar", "--ci"},
			expected: []string{"earthly", "-s", "foo=XXXXX", "--ci"},
		},
		{
			args:     []string{"earthly", "-s", "foo", "--ci"},
			expected: []string{"earthly", "-s", "foo", "--ci"},
		},
	} {
		actual := redactSecretsFromArgs(testCase.args)
		require.ElementsMatch(t, testCase.expected, actual)
	}
}
