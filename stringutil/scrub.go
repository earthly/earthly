package stringutil

import "regexp"

var scrubRegexp = regexp.MustCompile(`:([^@]+)@`)

// ScrubCredentials removes credentials from a string
func ScrubCredentials(s string) string {
	return string(scrubRegexp.ReplaceAll([]byte(s), []byte("***")))
}
