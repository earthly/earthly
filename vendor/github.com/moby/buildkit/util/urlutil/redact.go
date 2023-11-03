package urlutil

import (
	"net/url"
	"regexp"
	"strings"
)

const mask = "xxxxx"

// RedactCredentials takes a URL and redacts username and password from it.
// e.g. "https://user:password@host.tld/path.git" will be changed to
// "https://xxxxx:xxxxx@host.tld/path.git".
func RedactCredentials(s string) string {
	ru, err := url.Parse(s)
	if err != nil {
		return s // string is not a URL, just return it
	}
	var (
		hasUsername bool
		hasPassword bool
	)
	if ru.User != nil {
		hasUsername = len(ru.User.Username()) > 0
		_, hasPassword = ru.User.Password()
	}
	if hasUsername && hasPassword {
		ru.User = url.UserPassword(mask, mask)
	} else if hasUsername {
		ru.User = url.User(mask)
	} else if hasPassword {
		ru.User = url.UserPassword(ru.User.Username(), mask)
	}
	return ru.String()
}

var urlRegexp = regexp.MustCompile(`(https?://\S+)`)

// RedactAllCredentials is earthly-specific
func RedactAllCredentials(s string) string {
	var sb strings.Builder
	matches := urlRegexp.FindAllStringIndex(s, -1)
	i := 0
	for _, m := range matches {
		sb.WriteString(s[i:m[0]])
		sb.WriteString(RedactCredentials(s[m[0]:m[1]]))
		i = m[1]
	}
	sb.WriteString(s[i:len(s)])
	return sb.String()
}
