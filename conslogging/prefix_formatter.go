package conslogging

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/earthly/earthly/util/stringutil"
)

var (
	// urlPathPartRegex is used to find all parts of a url
	// for example github.com/earthly/my-repo => p1="github.com/", p2="earthly/" (note: "my-repo" is intentionally not captured as p3)
	urlPathPartRegex = regexp.MustCompile(`(.*?/)`)
	// githubRegex Matches :2dd88e53f2e59e96ec1f9215f24a3981e5565edf+ in a prefix.
	// 	Prefix containing hash may resemble: g/e/hello-world:2dd88e53f2e59e96ec1f9215f24a3981e5565edf+base
	//	Prefix must be exactly 40 characters
	githubRegex = regexp.MustCompile(`:[a-f0-9]{40}\+`)
	// gitURLRegex matches the url appearing in parentheses for example:
	// +my-target(https://github/earthly/earthly)
	gitURLRegex = regexp.MustCompile(`\(.+?\)`)
	// urlPrefixRegex is used to captured url protocol, for example "https://" in "https://github.com/earthly/earthly"
	urlPrefixRegex = regexp.MustCompile("^.+?//")
	// targetURLRegex is used to capture any target path - relative (./my-dir+my-target), absolute (/abs/my-dir+my-target) or remote (github.com/my-org-my-repo+my-target)
	// the url my include an optional branch name or commit sha, e.g. github.com/my-org/my-repo:my-branch+my-target
	targetURLRegex       = regexp.MustCompile(`^.+?(:|\+)`)
	gitURLWithCredsRegex = regexp.MustCompile(`(?P<protocol>.+?)://(?P<user>.+?):(?P<password>.+?)@(?P<repoURL>.+?).git#(?P<ref>.+?)$`)
	formatter            = NewPrefixFormatter(truncateURLWithCreds, truncateSha, truncateGITURL, truncateTargetURL)
)

type prefixFormatter struct {
	formatOpts []formatOpt
	cache      sync.Map
}

type formatOpt func(str string, padding int, curLen int) string

func truncateURLWithCreds(str string, padding int, curLen int) string {
	namedMatches, namedGroups := stringutil.NamedGroupMatches(str, gitURLWithCredsRegex)
	if len(namedMatches) != 5 {
		// no match for the regex, return original string
		return str
	}
	matches := make([]string, 0, len(namedMatches))
	for _, name := range namedGroups {
		if len(namedMatches[name]) == 0 {
			//something was wrong with the regex, return original string
			return str
		}
		if curLen <= padding || namedMatches[name][0] == "" {
			// no need to keep truncating the url parts
			matches = append(matches, namedMatches[name][0])
		} else if name == "repoURL" {
			truncatedURL := truncateURL(namedMatches[name][0], padding, curLen)
			matches = append(matches, truncatedURL)
			curLen -= len(namedMatches[name][0]) - len(truncatedURL)
		} else {
			matches = append(matches, string(namedMatches[name][0][0]))
			curLen -= len(namedMatches[name][0]) - 1
		}
	}
	seps := []string{"://", ":", "@", "#", ""}
	var sb strings.Builder
	for i := range matches {
		sb.WriteString(matches[i])
		sb.WriteString(seps[i])
	}
	return sb.String()
}

func truncateSha(str string, _, _ int) string {
	return githubRegex.ReplaceAllStringFunc(str, func(s string) string {
		return s[:8] + "+"
	})
}

func truncateURL(str string, padding int, curLen int) string {
	return urlPathPartRegex.ReplaceAllStringFunc(str, func(part string) string {
		if curLen <= padding || len(part) <= 1 || part == ".." {
			return part
		}
		if strings.HasSuffix(part, "/") {
			curLen -= len(part) - 1
			return fmt.Sprintf("%c%c", part[0], '/')
		}
		curLen -= len(part) - 1
		return string(part[0])
	})
}

func truncateGITURL(str string, padding int, curLen int) string {
	return gitURLRegex.ReplaceAllStringFunc(str, func(s string) string {
		s = s[1 : len(s)-1]
		urlProtocol := urlPrefixRegex.FindString(s)
		s = strings.TrimPrefix(s, urlProtocol)
		l1 := len(s)
		s = normalize(s)
		charsRemoved := l1 - len(s)
		curLen -= charsRemoved
		if curLen <= padding {
			return fmt.Sprintf("(%s)", s)
		}
		return fmt.Sprintf("(%s%s)", urlProtocol, truncateURL(s, padding, curLen))
	})
}

func truncateTargetURL(str string, padding int, curLen int) string {
	return targetURLRegex.ReplaceAllStringFunc(str, func(s string) string {
		suffixChar := s[len(s)-1]
		s = s[:len(s)-1]
		l1 := len(s)
		s = normalize(s)
		charsRemoved := l1 - len(s)
		curLen -= charsRemoved
		if curLen <= padding {
			return fmt.Sprintf("%s%c", s, suffixChar)
		}
		return fmt.Sprintf("%s%c", truncateURL(s, padding, curLen), suffixChar)
	})
}

func normalize(s string) string {
	isLocalDirPrefix := strings.HasPrefix(s, "./")
	s = filepath.Clean(s)
	if isLocalDirPrefix {
		return fmt.Sprintf("./%s", s)
	}
	return s
}

func NewPrefixFormatter(formatOpt ...formatOpt) *prefixFormatter {
	return &prefixFormatter{
		formatOpts: formatOpt,
	}
}

func (pb *prefixFormatter) getKey(prefix string, padding int) string {
	return fmt.Sprintf("%s-%d", prefix, padding)
}

func (pb *prefixFormatter) Format(prefix string, padding int) (modifiedPrefix string) {
	if padding <= NoPadding {
		return prefix
	}
	key := pb.getKey(prefix, padding)
	if cachedPrefix, ok := pb.cache.Load(key); ok {
		return cachedPrefix.(string)
	}
	defer func() {
		modifiedPrefix = fmt.Sprintf("%*s", padding, prefix)
		pb.cache.Store(key, modifiedPrefix)
	}()
	curLen := len(prefix)
	for _, formatOpt := range pb.formatOpts {
		if curLen <= padding {
			return
		}
		prefix = formatOpt(prefix, padding, curLen)
		curLen = len(prefix)
	}
	return
}
