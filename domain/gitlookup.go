package domain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type gitMatcher struct {
	pattern  string
	user     string
	suffix   string
	protocol string
}

// GitLookup looksup gits
type GitLookup struct {
	matchers []gitMatcher
}

// NewGitLookup creates new lookuper
func NewGitLookup() *GitLookup {
	matchers := []gitMatcher{
		{
			pattern:  "github.com/[^/]+/[^/]+",
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
		},
		{
			pattern:  "gitlab.com/[^/]+/[^/]+",
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
		},
		{
			pattern:  "bitbucket.com/[^/]+/[^/]+",
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
		},
		{
			pattern:  "192.168.0.116/my/test/path/[^/]+",
			user:     "alex",
			suffix:   ".git",
			protocol: "ssh",
		},
	}

	gl := &GitLookup{
		matchers: matchers,
	}
	return gl
}

// GlobalGitLookup allows for converting git urls
var GlobalGitLookup = NewGitLookup()

// ErrNoMatch occurs when no git matcher is found
var ErrNoMatch = fmt.Errorf("no git match found")

// DisableSSH changes all git matchers from ssh to https
func (gl *GitLookup) DisableSSH() {
	for i, m := range gl.matchers {
		if m.protocol == "ssh" {
			gl.matchers[i].protocol = "https"
		}
	}
}

func (gl *GitLookup) getGitMatcher(path string) (string, *gitMatcher, error) {
	for _, m := range gl.matchers {
		r, err := regexp.Compile(m.pattern)
		if err != nil {
			return "", nil, errors.Wrapf(err, "failed to compile regex %s", m.pattern)
		}
		match := r.FindString(path)
		if match != "" {
			return match, &m, nil
		}
	}
	return "", nil, ErrNoMatch
}

// SplitGitTarget splits a git repo target into base repo and relative repo path; for example:
//   "github.com/earthly/earthly"             ---> ("github.com/earthly/earthly", "")
//   "github.com/earthly/earthly/examples"    ---> ("github.com/earthly/earthly", "examples")
//   "github.com/earthly/earthly/examples/go" ---> ("github.com/earthly/earthly", "examples/go")
func (gl *GitLookup) SplitGitTarget(path string) (string, string, error) {
	match, _, err := gl.getGitMatcher(path)
	if err != nil {
		return "", "", err
	}
	n := len(match) + 1
	subPath := ""
	if len(path) > n {
		subPath = path[n:]
	}
	return match, subPath, nil
}

// GetCloneURL returns a string
func (gl *GitLookup) GetCloneURL(path string) (string, error) {
	match, m, err := gl.getGitMatcher(path)
	if err != nil {
		return "", err
	}
	var gitURL string
	switch m.protocol {
	case "ssh":
		gitURL = m.user + "@" + strings.Replace(match, "/", ":", 1) + m.suffix
	case "http":
		gitURL = "http://" + match + m.suffix
	case "https":
		gitURL = "https://" + match + m.suffix
	default:
		return "", fmt.Errorf("unsupported protocol: %s", m.protocol)
	}
	return gitURL, nil
}
