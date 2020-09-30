package autocomplete

import (
	"fmt"
	"os/user"
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagCompletion(t *testing.T) {
	matches, err := GetPotentials("earth --fl", 10, []string{"flag", "fleet", "fig"}, nil)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag ", "--fleet "}, matches)
}

func TestCommandCompletion(t *testing.T) {
	matches, err := GetPotentials("earth pru", 9, nil, []string{"prune", "foo"})
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune "}, matches)
}

func TestPathCompletion(t *testing.T) {
	matches, err := GetPotentials("earth .", 7, nil, nil)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"./", "../"}, matches)
}

func Test_getPotentialPaths(t *testing.T) {
	return
	u, err := user.Current()
	require.NoError(t, err)

	t.Run("user list", func(t *testing.T) {
		userList, err := getPotentialPaths("~")
		require.NoError(t, err)
		require.True(t, len(userList) > 0)
	})

	t.Run("user path", func(t *testing.T) {
		res, err := getPotentialPaths(fmt.Sprintf("~%s", u.Username))
		require.NoError(t, err)
		require.True(t, len(res) > 0)
	})
	t.Run("directory with path /", func(t *testing.T) {
		res, err := getPotentialPaths("/")
		require.NoError(t, err)
		require.True(t, len(res) > 0)
		t.Log(res)
	})
}
