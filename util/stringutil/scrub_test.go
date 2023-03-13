package stringutil

import (
	"testing"
)

func TestScrub(t *testing.T) {
	s := ScrubCredentials("https://user:password@github.com/org/repo.git")
	Equal(t, "https://user:***@github.com/org/repo.git", s)
}

func TestScrubMissingProtocol(t *testing.T) {
	s := ScrubCredentials("user:password@github.com/org/repo.git")
	Equal(t, "user:***@github.com/org/repo.git", s)
}
