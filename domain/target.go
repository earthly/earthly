package domain

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

// Target is a earth target identifier.
type Target struct {
	// Remote and canonical representation.
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

// String returns a string representation of the Target.
func (et Target) String() string {
	if et.IsLocalExternal() {
		return fmt.Sprintf("%s+%s", et.LocalPath, et.Target)
	}
	if et.IsRemote() {
		s := et.GitURL
		if et.GitPath != "" {
			s += "/" + et.GitPath
		}
		if et.Tag != "" {
			s += ":" + et.Tag
		}
		s += "+" + et.Target
		return s
	}
	// Local internal.
	return fmt.Sprintf("+%s", et.Target)
}

// StringCanonical returns a string representation of the Target, in canonical form.
func (et Target) StringCanonical() string {
	if et.GitURL != "" {
		tag := fmt.Sprintf(":%s", et.Tag)
		if et.Tag == "" {
			tag = ""
		}
		return fmt.Sprintf("%s/%s%s+%s", et.GitURL, et.GitPath, tag, et.Target)
	}
	return et.String()
}

// ProjectCanonical returns a string representation of the project of the target, in canonical form.
func (et Target) ProjectCanonical() string {
	if et.GitURL != "" {
		tag := fmt.Sprintf(":%s", et.Tag)
		if et.Tag == "" {
			tag = ""
		}
		return fmt.Sprintf("%s/%s%s", et.GitURL, et.GitPath, tag)
	}
	if et.LocalPath == "." {
		return ""
	}
	return path.Base(et.LocalPath)
}

type gitMatcher struct {
	pattern string
	user    string
	suffix  string
}

// GitLookup looksup gits
type GitLookup struct {
	matchers []gitMatcher
}

// NewGitLookup creates new lookuper
func NewGitLookup() *GitLookup {
	matchers := []gitMatcher{
		{
			pattern: "github.com/[^/]+/[^/]+",
			user:    "git",
			suffix:  ".git",
		},
		{
			pattern: "gitlab.com/[^/]+/[^/]+",
			user:    "git",
			suffix:  ".git",
		},
		{
			pattern: "bitbucket.com/[^/]+/[^/]+",
			user:    "git",
			suffix:  ".git",
		},
		{
			pattern: "192.168.0.116/my/test/path/[^/]+",
			user:    "alex",
			suffix:  ".git",
		},
	}

	gl := &GitLookup{
		matchers: matchers,
	}
	return gl
}

// TODO needs fixing
var TODO = NewGitLookup()

// ParseGitURLandPath returns git path in the form user@host:path/to/repo.git, and any subdir
func (gl *GitLookup) ParseGitURLandPath(path string) (string, string, error) {
	fmt.Printf("ParseGitURLandPath(%q)\n", path)
	for _, m := range gl.matchers {
		r, err := regexp.Compile(m.pattern)
		if err != nil {
			panic(err)
		}
		match := r.FindString(path)
		if match != "" {
			n := len(match) + 1
			subPath := ""
			if len(path) > n {
				subPath = path[n:]
			}
			if strings.HasSuffix(match, "/") {
				panic("bad")
			}
			if strings.HasPrefix(subPath, "/") {
				fmt.Println(subPath)
				panic("bad1")
			}
			fmt.Printf("parsed %q into %q and %q\n", path, match, subPath)
			return match, subPath, nil
		}
	}
	fmt.Printf("failed to parse %q\n", path)
	return "", "", nil
}

// GetCloneURL returns a string
func (gl *GitLookup) GetCloneURL(path string) (string, error) {
	for _, m := range gl.matchers {
		r, err := regexp.Compile(m.pattern)
		if err != nil {
			panic(err)
		}
		match := r.FindString(path)
		if match == "" {
			continue
		}
		subPath := path[len(match):]
		if strings.HasSuffix(match, "/") {
			panic("bad2")
		}
		if subPath != "" {
			panic("bad3")
		}

		parts := strings.SplitN(match, "/", 2)
		if len(parts) != 2 {
			panic("bad4")
		}

		s := fmt.Sprintf("%s@%s:%s", m.user, parts[0], parts[1])
		return s, nil
	}
	fmt.Printf("failed to parse %q\n", path)
	return "", nil
}

// ParseTarget parses a string into a Target.
func ParseTarget(fullTargetName string) (Target, error) {
	partsPlus := strings.SplitN(fullTargetName, "+", 2)
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

		gitURL, gitPath, err := TODO.ParseGitURLandPath(partsColon[0])
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
