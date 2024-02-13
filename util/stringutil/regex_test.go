package stringutil

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestNamedGroupMatches(t *testing.T) {
	tests := map[string]struct {
		s             string
		re            *regexp.Regexp
		expectedMap   map[string][]string
		expectedSlice []string
	}{
		"empty map when no matches": {
			s:             "123",
			re:            regexp.MustCompile(`[a-z]+`),
			expectedMap:   map[string][]string{},
			expectedSlice: []string{},
		},
		"no match when no name": {
			s:             "123",
			re:            regexp.MustCompile(`[0-9]+`),
			expectedMap:   map[string][]string{},
			expectedSlice: []string{},
		},
		"can return a single match": {
			s:             "123",
			re:            regexp.MustCompile(`(?P<number>[0-9]+)`),
			expectedMap:   map[string][]string{"number": {"123"}},
			expectedSlice: []string{"number"},
		},
		"can return multiple matches for a single key": {
			s:             "123",
			re:            regexp.MustCompile(`(?P<number>[0-9])`),
			expectedMap:   map[string][]string{"number": {"1", "2", "3"}},
			expectedSlice: []string{"number"},
		},
		"can return multiple matches for multiple keys": {
			s:             "!1a2b3c!@",
			re:            regexp.MustCompile(`(?P<number>[0-9])(?P<letter>[a-z])(?P<notfound>(?:foo)?)`),
			expectedMap:   map[string][]string{"number": {"1", "2", "3"}, "letter": {"a", "b", "c"}},
			expectedSlice: []string{"number", "letter"},
		},
	}
	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mapRes, sliceRes := NamedGroupMatches(tc.s, tc.re)
			assert.Equal(t, tc.expectedMap, mapRes)
			assert.Equal(t, tc.expectedSlice, sliceRes)
			assert.Equal(t, len(mapRes), len(sliceRes))
		})
	}
}
