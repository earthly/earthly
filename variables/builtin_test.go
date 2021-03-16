package variables

import (
	"testing"

	. "github.com/stretchr/testify/assert"
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
	}

	for _, tt := range tests {
		ans := getProjectName(tt.tag)
		Equal(t, tt.safe, ans)
	}
}
