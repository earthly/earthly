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
	// Short circuit when no basic auth credentials present.
	if !strings.Contains(s, "@") {
		return s
	}
	parts := strings.Split(s, " ")
	ret := []string{}
	for _, part := range parts {
		if strings.Contains(part, "@") {
			ret = append(ret, ScrubCredentials(part))
		} else {
			ret = append(ret, part)
		}
	}
	return strings.Join(ret, " ")
}

// ScrubANSICodes removes ANSI escape codes from a string.
func ScrubANSICodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleanedString := re.ReplaceAllString(input, "")
	return cleanedString
}
