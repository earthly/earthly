package conslogging

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	// urlPath is used to find all parts of a url
	// for example github.com/earthly/my-repo => p1="github.com/", p2="earthly/" (note: "my-repo" is intentionally not captured as p3)
	urlPath = regexp.MustCompile(`(.*?/)`)
	// githubRegExp Matches :2dd88e53f2e59e96ec1f9215f24a3981e5565edf+ in a prefix.
	// 	Prefix containing hash may resemble: g/e/hello-world:2dd88e53f2e59e96ec1f9215f24a3981e5565edf+base
	//	Prefix must be exactly 40 characters
	githubRegExp = regexp.MustCompile(`:[a-f0-9]{40}\+`)
	// gitURL matches the url appearing in parentheses for example:
	// +my-target(https://github/earthly/earthly)
	gitURL = regexp.MustCompile(`\(.+?\)`)
	// urlPrefix is used to captured url protocol, for example "https://" in "https://github.com/earthly/earthly"
	urlPrefix = regexp.MustCompile("^.+?//")
	// targetURL is used to capture any target path - relative (./my-dir+my-target), absolute (/abs/my-dir+my-target) or remote (github.com/my-org-my-repo+my-target)
	// the url my include an optional branch name or commit sha, e.g. github.com/my-org/my-repo:my-branch+my-target
	targetURL = regexp.MustCompile(`^.+?(:|\+)`)
	formatter = NewPrefixFormatter(truncateSha, truncateGITURL, truncateTargetURL)
)

type prefixFormatter struct {
	formatOpts []formatOpt
	cache      sync.Map
}

type formatOpt func(str string, padding int, curLen int) string

func truncateSha(str string, _, _ int) string {
	return githubRegExp.ReplaceAllStringFunc(str, func(s string) string {
		return s[:8] + "+"
	})
}

func truncateURL(str string, padding int, curLen int) string {
	return urlPath.ReplaceAllStringFunc(str, func(part string) string {
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
	return gitURL.ReplaceAllStringFunc(str, func(s string) string {
		s = s[1 : len(s)-1]
		urlProtocol := urlPrefix.FindString(s)
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
	return targetURL.ReplaceAllStringFunc(str, func(s string) string {
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

func padStr(s string, padding int) string {
	formatString := fmt.Sprintf("%%%vv", padding)
	return fmt.Sprintf(formatString, s)
}

func normalize(s string) string {
	isLocalDirPrefix := strings.HasPrefix(s, "./")
	s = filepath.Clean(s)
	if isLocalDirPrefix {
		s = fmt.Sprintf("./%s", s)
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

func (pb *prefixFormatter) Format(prefix string, padding int) string {
	if padding <= NoPadding {
		return prefix
	}
	key := pb.getKey(prefix, padding)
	if cachedPrefix, ok := pb.cache.Load(key); ok {
		return cachedPrefix.(string)
	}
	curLen := len(prefix)
	for _, formatOpt := range pb.formatOpts {
		prefix = formatOpt(prefix, padding, curLen)
		curLen = len(prefix)
	}
	modifiedPrefix := padStr(prefix, padding)
	pb.cache.Store(key, modifiedPrefix)
	return modifiedPrefix
}
