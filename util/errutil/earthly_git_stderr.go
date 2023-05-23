package errutil

import (
	"encoding/base64"
	"regexp"
	"strings"
)

const EarthlyGitStdErrMagicString = "EARTHLY_GIT_STDERR"

var gitStdErrRegexp = regexp.MustCompile(`EARTHLY_GIT_STDERR: ([A-Za-z0-9+/]*={0,2}):`)

// ExtractEarthlyGitStdErr scans an error string for a base64 payload that contains the git stderr
// and returns the extracted stderr and a shorter error string which does not include the base64 payload.
// if no payload was extracted, then false is returned
func ExtractEarthlyGitStdErr(errStr string) (extracted, shorterErr string, ok bool) {
	shorterErr = gitStdErrRegexp.ReplaceAllString(errStr, "")
	matches := gitStdErrRegexp.FindStringSubmatch(errStr)
	if len(matches) == 2 {
		if stderr, err := base64.StdEncoding.DecodeString(matches[1]); err == nil {
			return strings.TrimSpace(string(stderr)), shorterErr, true
		}
	}
	return "", "", false
}
