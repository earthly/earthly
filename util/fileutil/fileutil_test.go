package fileutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobDirs(t *testing.T) {
	tests := []struct {
		pattern string
		results []string
	}{
		{
			pattern: "testdata/globdirs/*",
			results: []string{"testdata/globdirs/bar", "testdata/globdirs/baz", "testdata/globdirs/foo"},
		},
		{
			pattern: "testdata/globdirs/b*",
			results: []string{"testdata/globdirs/bar", "testdata/globdirs/baz"},
		},
		{
			pattern: "testdata/globdirs/file.txt",
			results: nil,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			results, err := GlobDirs(test.pattern)
			r := require.New(t)
			r.Equal(test.results, results)
			r.NoError(err)
		})
	}
}
