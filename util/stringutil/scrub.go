package stringutil

import (
	"regexp"
	"strings"
)

var scrubRegexp = regexp.MustCompile(`^(([a-zA-Z]+://)?([^:]+)):([^@]+)@`)

// ScrubCredentials removes credentials from a string
func ScrubCredentials(s string) string {
	return string(scrubRegexp.ReplaceAll([]byte(s), []byte("$1:xxxxx@")))
}

// ScrubCredentialsAll scrubs all credentials from a longer piece of text.
func ScrubCredentialsAll(s string) string {
	parts := strings.Split(s, " ")
	ret := []string{}
	for _, part := range parts {
		if strings.Contains(part, "@") {
			ret = append(ret, (ScrubCredentials(part)))
		} else {
			ret = append(ret, part)
		}
	}
	return strings.Join(ret, " ")
}
