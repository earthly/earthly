package stringutil

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestScrub(t *testing.T) {
	s := ScrubCredentials("https://user:password@github.com/org/repo.git")
	Equal(t, "https://user:***@github.com/org/repo.git", s)
}
