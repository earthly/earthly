package buildcontext

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type gitMatcher struct {
	name     string
	re       *regexp.Regexp
	user     string
	suffix   string
	protocol string
	password string
	keyScan  string
}

// GitLookup looksup gits
type GitLookup struct {
	matchers []*gitMatcher
	catchAll *gitMatcher
}

// NewGitLookup creates new lookuper
func NewGitLookup() *GitLookup {
	matchers := []*gitMatcher{
		{
			name:     "github.com",
			re:       regexp.MustCompile("github.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
			keyScan:  "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			name:     "gitlab.com",
			re:       regexp.MustCompile("gitlab.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
			keyScan:  "gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9",
		},
		{
			name:     "bitbucket.com",
			re:       regexp.MustCompile("bitbucket.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
			keyScan:  "bitbucket.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==",
		},
	}

	gl := &GitLookup{
		matchers: matchers,
		catchAll: &gitMatcher{
			name:     "",
			re:       regexp.MustCompile("[^/]+/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "ssh",
		},
	}
	return gl
}

// ErrNoMatch occurs when no git matcher is found
var ErrNoMatch = fmt.Errorf("no git match found")

// DisableSSH changes all git matchers from ssh to https
func (gl *GitLookup) DisableSSH() {
	for i, m := range gl.matchers {
		if m.protocol == "ssh" {
			gl.matchers[i].protocol = "https"
		}
	}
	if gl.catchAll.protocol == "ssh" {
		gl.catchAll.protocol = "https"
	}
}

// AddMatcher adds a new matcher for looking up git repos
func (gl *GitLookup) AddMatcher(name, pattern, user, password, suffix, protocol, keyScan string) error {
	if protocol == "http" && password != "" {
		return fmt.Errorf("using a password with http for %s is insecure", name)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "failed to compile regex %s", pattern)
	}
	switch protocol {
	case "http", "https", "ssh":
		break
	default:
		return fmt.Errorf("unsupported git protocol %q", protocol)
	}

	gm := &gitMatcher{
		name:     name,
		re:       re,
		user:     user,
		password: password,
		suffix:   suffix,
		protocol: protocol,
		keyScan:  keyScan,
	}

	// update existing entry
	for i, m := range gl.matchers {
		if m.name == name {
			if gm.keyScan == "" {
				gm.keyScan = m.keyScan
			}
			gl.matchers[i] = gm
			return nil
		}
	}

	// add new entry
	gl.matchers = append(gl.matchers, gm)
	return nil
}

func (gl *GitLookup) getGitMatcher(path string) (string, *gitMatcher, error) {
	if len(gl.matchers) == 0 {
		panic("no matchers")
	}
	for _, m := range gl.matchers {
		match := m.re.FindString(path)
		if match != "" {
			return match, m, nil
		}
	}

	match := gl.catchAll.re.FindString(path)
	if match != "" {
		return match, gl.catchAll, nil
	}

	return "", nil, ErrNoMatch
}

// GetCloneURL returns the repo to clone, and a path relative to the repo
//   "github.com/earthly/earthly"             ---> ("git@github.com/earthly/earthly.git", "")
//   "github.com/earthly/earthly/examples"    ---> ("git@github.com/earthly/earthly.git", "examples")
//   "github.com/earthly/earthly/examples/go" ---> ("git@github.com/earthly/earthly.git", "examples/go")
// Additionally a ssh keyscan might be returned (or an empty string indicating none was configured)
func (gl *GitLookup) GetCloneURL(path string) (string, string, string, error) {
	match, m, err := gl.getGitMatcher(path)
	if err != nil {
		return "", "", "", err
	}

	n := len(match) + 1
	subPath := ""
	if len(path) > n {
		subPath = path[n:]
	}

	var gitURL, keyScan string
	switch m.protocol {
	case "ssh":
		gitURL = m.user + "@" + strings.Replace(match, "/", ":", 1) + m.suffix
		keyScan = m.keyScan
	case "http", "https":
		var userAndPass string
		if m.user != "" && m.password != "" {
			userAndPass = url.QueryEscape(m.user) + ":" + url.QueryEscape(m.password) + "@"
		}
		gitURL = m.protocol + "://" + userAndPass + match + m.suffix
	default:
		return "", "", "", fmt.Errorf("unsupported protocol: %s", m.protocol)
	}
	return gitURL, subPath, keyScan, nil
}
