//go:build hasgitdirectory
// +build hasgitdirectory

package analytics

import (
	"testing"
)

// TestGetRepoHash tests the git repo hashing never changes
// in order to ensure our analytics stay consistent
func TestGetRepoHash(t *testing.T) {
	hash := hashString(getLocalRepo())
	Equal(t, hash, "5d892560f423223dec22b9b03e11d3aa3775871a80962326ac80543401843749")
}
