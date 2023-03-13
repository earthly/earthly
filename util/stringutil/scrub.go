package stringutil

import "regexp"

var scrubRegexp = regexp.MustCompile(`^(([a-zA-Z]+://)?([^:]+)):([^@]+)@`)

// ScrubCredentials removes credentials from a string
func ScrubCredentials(s string) string {
	return string(scrubRegexp.ReplaceAll([]byte(s), []byte("$1:***@")))
}
