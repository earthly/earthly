package domain

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Target is a earth target identifier.
type Target struct {
	GitURL  string // "github.com/earthly/earthly"
	GitPath string // "examples/go"
	Tag     string // "main"

	// Local representation.
	LocalPath string `json:"localPath"`

	// Target name.
	Target string `json:"target"`
}

// IsExternal returns whether the target is external to the current project.
func (et Target) IsExternal() bool {
	return et.IsRemote() || et.IsLocalExternal()
}

// IsLocalInternal returns whether the target is a local.
func (et Target) IsLocalInternal() bool {
	return et.LocalPath == "."
}

// IsLocalExternal returns whether the target is a local, but external target.
func (et Target) IsLocalExternal() bool {
	return et.LocalPath != "." && et.LocalPath != ""
}

// IsRemote returns whether the target is remote.
func (et Target) IsRemote() bool {
	return !et.IsLocalExternal() && !et.IsLocalInternal()
}

// DebugString returns a string that can be printed out for debugging purposes
func (et Target) DebugString() string {
	return fmt.Sprintf("gitURL: %q; gitPath: %q; tag: %q; LocalPath: %q; Target: %q", et.GitURL, et.GitPath, et.Tag, et.LocalPath, et.Target)
}

// String returns a string representation of the Target.
func (et Target) String() string {
	if et.IsLocalExternal() {
		return fmt.Sprintf("%s+%s", escapePlus(et.LocalPath), et.Target)
	}
	if et.IsRemote() {
		s := et.GitURL
		if et.GitPath != "" {
			s += "/" + escapePlus(et.GitPath)
		}
		if et.Tag != "" {
			s += ":" + escapePlus(et.Tag)
		}
		s += "+" + escapePlus(et.Target)
		return s
	}
	// Local internal.
	return fmt.Sprintf("+%s", et.Target)
}

// StringCanonical returns a string representation of the Target, in canonical form.
func (et Target) StringCanonical() string {
	if et.GitURL != "" {
		s := et.GitURL
		if et.GitPath != "" {
			s += "/" + escapePlus(et.GitPath)
		}
		if et.Tag != "" {
			s += ":" + escapePlus(et.Tag)
		}
		s += "+" + escapePlus(et.Target)
		return s
	}
	return et.String()
}

// ProjectCanonical returns a string representation of the project of the target, in canonical form.
func (et Target) ProjectCanonical() string {
	if et.GitURL != "" {
		s := escapePlus(et.GitURL)
		if et.GitPath != "" {
			s += "/" + escapePlus(et.GitPath)
		}
		if et.Tag != "" {
			s += ":" + escapePlus(et.Tag)
		}
		return s
	}
	if et.LocalPath == "." {
		return ""
	}
	return escapePlus(path.Base(et.LocalPath))
}

// ParseTarget parses a string into a Target.
func ParseTarget(fullTargetName string) (Target, error) {
	partsPlus, err := splitUnescapePlus(fullTargetName)
	if err != nil {
		return Target{}, err
	}
	if len(partsPlus) != 2 {
		return Target{}, fmt.Errorf("Invalid target ref %s", fullTargetName)
	}
	if partsPlus[0] == "" {
		// Local target.
		return Target{
			LocalPath: ".",
			Target:    partsPlus[1],
		}, nil
	} else if strings.HasPrefix(partsPlus[0], ".") ||
		strings.HasPrefix(partsPlus[0], "/") {
		// Local external target.
		localPath := partsPlus[0]
		if path.IsAbs(localPath) {
			localPath = path.Clean(localPath)
		} else {
			localPath = path.Clean(localPath)
			if !strings.HasPrefix(localPath, ".") {
				localPath = fmt.Sprintf("./%s", localPath)
			}
		}
		return Target{
			LocalPath: localPath,
			Target:    partsPlus[1],
		}, nil
	} else {
		// Remote target.
		tag := ""
		partsColon := strings.SplitN(partsPlus[0], ":", 2)
		if len(partsColon) == 2 {
			tag = partsColon[1]
		}

		gitURL, gitPath, err := GlobalGitLookup.SplitGitTarget(partsColon[0])
		if err != nil {
			return Target{}, err
		}
		return Target{
			GitURL:  gitURL,
			GitPath: gitPath,
			Tag:     tag,
			Target:  partsPlus[1],
		}, nil
	}
}

// JoinTargets returns the result of interpreting target2 as relative to target1.
func JoinTargets(target1 Target, target2 Target) (Target, error) {
	ret := target2
	if target1.IsRemote() {
		// target1 is remote. Turn relative targets into remote targets.
		if !ret.IsRemote() {
			ret.GitURL = target1.GitURL
			ret.GitPath = target1.GitPath
			ret.Tag = target1.Tag
			if ret.IsLocalExternal() {
				if path.IsAbs(ret.LocalPath) {
					return Target{}, fmt.Errorf(
						"Absolute path %s not supported as reference in external target context", ret.LocalPath)
				}

				ret.GitPath = path.Join(target1.GitPath, ret.LocalPath)
				ret.LocalPath = ""
			} else if ret.IsLocalInternal() {
				ret.LocalPath = ""
			}
		}
	} else {
		if ret.IsLocalExternal() {
			if path.IsAbs(ret.LocalPath) {
				ret.LocalPath = path.Clean(ret.LocalPath)
			} else {
				ret.LocalPath = path.Join(target1.LocalPath, ret.LocalPath)
				if !strings.HasPrefix(ret.LocalPath, ".") {
					ret.LocalPath = fmt.Sprintf("./%s", ret.LocalPath)
				}
			}
		} else if ret.IsLocalInternal() {
			ret.LocalPath = target1.LocalPath
		}
	}
	return ret, nil
}

// splitUnescapePlus performs a split on "+", but it accounts for escaping as "\+".
func splitUnescapePlus(str string) ([]string, error) {
	escape := false
	ret := make([]string, 0, 2)
	word := make([]rune, 0, len(str))
	for _, c := range str {
		if escape {
			word = append(word, c)
			escape = false
			continue
		}

		switch c {
		case '\\':
			escape = true
		case '+':
			ret = append(ret, string(word))
			word = word[:0]
		default:
			word = append(word, c)
		}
	}
	if escape {
		return nil, errors.Errorf("cannot split by +: unterminated escape sequence at the end of %s", str)
	}
	if len(word) > 0 {
		ret = append(ret, string(word))
	}
	return ret, nil
}

func escapePlus(str string) string {
	return strings.ReplaceAll(str, "+", "\\+")
}

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
