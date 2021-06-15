package stringutil

import (
	"regexp"
	"strings"
)

var scrubRegexp = regexp.MustCompile(`(//[^:]+):([^@]+)@`)

// ScrubCredentials removes credentials from a string
func ScrubCredentials(s string) string {
	if strings.Contains(s, ":xxxxx@") {
		// buildkit already does this under redact_credentials.go; don't double-scrub
		return s
	}
	return string(scrubRegexp.ReplaceAll([]byte(s), []byte("$1:***@")))
}
