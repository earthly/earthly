package conslogging

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prettyPrefix(t *testing.T) {
	testCases := []struct {
		name          string
		prefixPadding int
		prefix        string
		expected      string
	}{
		{
			name:          "does not truncate if prefixPadding is NoPadding",
			prefixPadding: NoPadding,
			prefix:        "github.com/earthly/earthly:80524f0d82a353b3444e83f056207e15f4d5447c+hello-world",
			expected:      "github.com/earthly/earthly:80524f0d82a353b3444e83f056207e15f4d5447c+hello-world",
		},
		{
			name:          "shortens git SHA if prefix is present",
			prefixPadding: DefaultPadding,
			prefix:        "github.com/earthly/earthly:80524f0d82a353b3444e83f056207e15f4d5447c+hello-world",
			expected:      "g/e/earthly:80524f0+hello-world",
		},
		{
			name:          "keeps branch name",
			prefixPadding: DefaultPadding,
			prefix:        "github.com/earthly/earthly:some-feature-branch+hello-world",
			expected:      "g/e/earthly:some-feature-branch+hello-world",
		},
		{
			name:          "keeps branch name closely resembling sha",
			prefixPadding: DefaultPadding,
			prefix:        "/e/hello-world:feedfacecafe",
			expected:      "/e/hello-world:feedfacecafe",
		},
		{
			name:          "keeps branch name containing special characters /-_",
			prefixPadding: DefaultPadding,
			prefix:        "github.com/earthly/earthly:-_/ryan_-/branch-names/-_in-here+hello-world",
			expected:      "g/e/earthly:-_/ryan_-/branch-names/-_in-here+hello-world",
		},
		{
			name:          "simple target with no path or github info",
			prefixPadding: DefaultPadding,
			prefix:        "+run",
			expected:      strings.Repeat(" ", DefaultPadding-4) + "+run",
		},
		{
			name:          "simple target with path",
			prefixPadding: DefaultPadding,
			prefix:        "github.com/earthly/earthly+run",
			expected:      "github.com/earthly/earthly+run",
		},
		{
			name:          "does not add padding if prefix longer than prefixPadding",
			prefixPadding: 3,
			prefix:        "+run",
			expected:      "+run",
		},
		{
			name:          "negative padding results in no change",
			prefixPadding: -10,
			prefix:        "+run",
			expected:      "+run",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, prettyPrefix(tc.prefixPadding, tc.prefix))
		})
	}

}
