package domain

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Target is a earth target identifier.
type Target struct {
	GitURL string // e.g. "github.com/earthly/earthly/examples/go"
	Tag    string // e.g. "main"

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
	return fmt.Sprintf("gitURL: %q; tag: %q; LocalPath: %q; Target: %q", et.GitURL, et.Tag, et.LocalPath, et.Target)
}

// String returns a string representation of the Target.
func (et Target) String() string {
	if et.IsLocalExternal() {
		return fmt.Sprintf("%s+%s", escapePlus(et.LocalPath), et.Target)
	}
	if et.IsRemote() {
		s := escapePlus(et.GitURL)
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
		s := escapePlus(et.GitURL)
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

		return Target{
			GitURL: partsColon[0],
			Tag:    tag,
			Target: partsPlus[1],
		}, nil
	}
}

// JoinTargets returns the result of interpreting target2 as relative to target1.
func JoinTargets(target1 Target, target2 Target) (Target, error) {
	ret := target2
	if target1.IsRemote() {
		// target1 is remote. Turn relative targets into remote targets.
		if !ret.IsRemote() {
			ret.Tag = target1.Tag
			if ret.IsLocalExternal() {
				if path.IsAbs(ret.LocalPath) {
					return Target{}, fmt.Errorf(
						"Absolute path %s not supported as reference in external target context", ret.LocalPath)
				}

				ret.GitURL = path.Join(target1.GitURL, ret.LocalPath)
				ret.LocalPath = ""
			} else if ret.IsLocalInternal() {
				ret.GitURL = target1.GitURL
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
