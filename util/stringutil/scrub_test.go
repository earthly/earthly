package stringutil

import (
	"testing"
)

func TestScrub(t *testing.T) {
	s := ScrubCredentials("https://user:password@github.com/org/repo.git")
	Equal(t, "https://user:xxxxx@github.com/org/repo.git", s)
}

func TestScrubMissingProtocol(t *testing.T) {
	s := ScrubCredentials("user:password@github.com/org/repo.git")
	Equal(t, "user:xxxxx@github.com/org/repo.git", s)
}

func TestScrubInline(t *testing.T) {
	s := ScrubCredentialsAll("Here is a URL: https://user:password@github.com/org/repo.git")
	Equal(t, "Here is a URL: https://user:xxxxx@github.com/org/repo.git", s)
}

func TestANSICodes(t *testing.T) {
	s := ScrubANSICodes("\033[0;32mCommand succeeded.\033[0m")
	Equal(t, "Command succeeded.", s)
}
