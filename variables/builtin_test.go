package variables

import (
	"testing"
)

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
		{"http://github.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"http://gitlab.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"https://github.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"https://user@github.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"https://user:password@github.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"git@github.com:earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"git@bitbucket.com:earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"ssh://git@github.com/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
		{"ssh://git@github.com:22/earthly/earthly/subdir/anothersubdir", "earthly/earthly/subdir/anothersubdir"},
	}

	for _, tt := range tests {
		ans := getProjectName(tt.tag)
		Equal(t, tt.safe, ans)
	}
}
