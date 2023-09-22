package stringutil

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestNamedGroupMatches(t *testing.T) {
	tests := map[string]struct {
		s        string
		re       *regexp.Regexp
		expected map[string][]string
	}{
		"empty map when no matches": {
			s:        "123",
			re:       regexp.MustCompile(`[a-z]+`),
			expected: map[string][]string{},
		},
		"no match when no name": {
			s:        "123",
			re:       regexp.MustCompile(`[0-9]+`),
			expected: map[string][]string{},
		},
		"can return a single match": {
			s:        "123",
			re:       regexp.MustCompile(`(?P<number>[0-9]+)`),
			expected: map[string][]string{"number": {"123"}},
		},
		"can return multiple matches for a single key": {
			s:        "123",
			re:       regexp.MustCompile(`(?P<number>[0-9])`),
			expected: map[string][]string{"number": {"1", "2", "3"}},
		},
		"can return multiple matches for multiple keys": {
			s:        "!1a2b3c!@",
			re:       regexp.MustCompile(`(?P<number>[0-9])(?P<letter>[a-z])`),
			expected: map[string][]string{"number": {"1", "2", "3"}, "letter": {"a", "b", "c"}},
		},
	}
	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := NamedGroupMatches(tc.s, tc.re)
			assert.Equal(t, tc.expected, res)
		})
	}
}
