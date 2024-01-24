package containerutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const containerListText = `7a3c4741b86e,earthly-darwin-proxy-T58UvV,Up 5 hours,alpine/socat:1.7.4.4,2024-01-23 12:53:32 -0800 PST
d8183461827c,earthly-dev-buildkitd,Up 5 hours,earthly/buildkitd:dev-main,2024-01-23 12:50:02 -0800 PST
3084cac7996e,earthly-buildkitd,Up 6 hours,earthly/buildkitd:prerelease,2024-01-23 12:31:06 -0800 PST`

func Test_parseContainerList(t *testing.T) {
	ret, err := parseContainerList(containerListText)
	r := require.New(t)
	r.NoError(err)
	r.Len(ret, 3)
	r.Equal("7a3c4741b86e", ret[0].ID)
	r.Equal("earthly-darwin-proxy-T58UvV", ret[0].Name)
	r.Equal("Up 5 hours", ret[0].Status)
	r.Equal("alpine/socat:1.7.4.4", ret[0].Image)
	r.Equal(int64(1706043212), ret[0].Created.Unix())
}

func Test_parseContainerList_whitespace(t *testing.T) {
	ret, err := parseContainerList(containerListText + "\n\n")
	r := require.New(t)
	r.NoError(err)
	r.Len(ret, 3)
}

func Test_parseContainerList_empty(t *testing.T) {
	ret, err := parseContainerList("\n\n")
	r := require.New(t)
	r.NoError(err)
	r.Len(ret, 0)
}
