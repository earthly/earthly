package autocomplete

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestFlagCompletion(t *testing.T) {
	matches, err := GetPotentials("--fl", []string{"flag", "fleet", "fig"}, nil)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag", "--fleet"}, matches)
}
func TestCommandCompletion(t *testing.T) {
	matches, err := GetPotentials("pru", nil, []string{"prune", "foo"})
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune"}, matches)
}
func TestPathCompletion(t *testing.T) {
	matches, err := GetPotentials(".", nil, nil)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"./", "../"}, matches)
}
