package domain

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Reference is a target or a command reference.
type Reference interface {
	// GetGitURL is the git URL part of the reference. E.g. "github.com/earthly/earthly/examples/go"
	GetGitURL() string
	// GetTag is the git tag of the reference. E.g. "main"
	GetTag() string
	// GetLocalPath is the local path representation of the reference. E.g. in "./some/path+something" this is "./some/path".
	GetLocalPath() string
	// GetName is the target name or the command name of the reference. E.g. in "+something" this is "something".
	GetName() string

	// IsExternal returns whether the target is external to the current project.
	IsExternal() bool

	// IsLocalInternal returns whether the target is a local.
	IsLocalInternal() bool

	// IsLocalExternal returns whether the target is a local, but external target.
	IsLocalExternal() bool

	// IsRemote returns whether the target is remote.
	IsRemote() bool

	// DebugString returns a string that can be printed out for debugging purposes.
	DebugString() string

	// String returns a string representation of the Target.
	String() string

	// StringCanonical returns a string representation of the Target, in canonical form.
	StringCanonical() string
	// ProjectCanonical returns a string representation of the project of the target, in canonical form.
	ProjectCanonical() string
}

// JoinReferences returns the result of interpreting r2 as relative to r1. The result will always have
// the same type as r2.
func JoinReferences(r1 Reference, r2 Reference) (Reference, error) {
	gitURL := r2.GetGitURL()
	tag := r2.GetTag()
	localPath := r2.GetLocalPath()
	name := r2.GetName()
	if r1.IsRemote() {
		// r1 is remote. Turn relative targets into remote targets.
		if !r2.IsRemote() {
			tag = r1.GetTag()
			if r2.IsLocalExternal() {
				if path.IsAbs(r2.GetLocalPath()) {
					return Target{}, errors.Errorf(
						"absolute path %s not supported as reference in external target context", r2.GetLocalPath())
				}

				gitURL = path.Join(r1.GetGitURL(), localPath)
				localPath = ""
			} else if r2.IsLocalInternal() {
				gitURL = r1.GetGitURL()
				localPath = ""
			}
		}
	} else {
		if r2.IsLocalExternal() {
			if path.IsAbs(localPath) {
				localPath = path.Clean(localPath)
			} else {
				localPath = path.Join(r1.GetLocalPath(), localPath)
				if !strings.HasPrefix(localPath, ".") {
					localPath = fmt.Sprintf("./%s", localPath)
				}
			}
		} else if r2.IsLocalInternal() {
			localPath = r1.GetLocalPath()
		}
	}
	switch r2.(type) {
	case Target:
		return Target{
			GitURL:    gitURL,
			Tag:       tag,
			LocalPath: localPath,
			Target:    name,
		}, nil
	case Command:
		return Command{
			GitURL:    gitURL,
			Tag:       tag,
			LocalPath: localPath,
			Command:   name,
		}, nil
	default:
		return nil, errors.New("joining references not supported for this type")
	}
}

func referenceString(r Reference) string {
	if r.IsLocalExternal() {
		return fmt.Sprintf("%s+%s", escapePlus(r.GetLocalPath()), r.GetName())
	}
	if r.IsRemote() {
		s := escapePlus(r.GetGitURL())
		if r.GetTag() != "" {
			s += ":" + escapePlus(r.GetTag())
		}
		s += "+" + escapePlus(r.GetName())
		return s
	}
	// Local internal.
	return fmt.Sprintf("+%s", r.GetName())
}

func referenceStringCanonical(r Reference) string {
	if r.GetGitURL() != "" {
		s := escapePlus(r.GetGitURL())
		if r.GetTag() != "" {
			s += ":" + escapePlus(r.GetTag())
		}
		s += "+" + escapePlus(r.GetName())
		return s
	}
	return r.String()
}

func referenceProjectCanonical(r Reference) string {
	if r.GetGitURL() != "" {
		s := escapePlus(r.GetGitURL())
		if r.GetTag() != "" {
			s += ":" + escapePlus(r.GetTag())
		}
		return s
	}
	if r.GetLocalPath() == "." {
		return ""
	}
	return escapePlus(path.Base(r.GetLocalPath()))
}

func parseCommon(fullName string) (gitURL string, tag string, localPath string, name string, err error) {
	partsPlus, err := splitUnescapePlus(fullName)
	if err != nil {
		return "", "", "", "", err
	}
	if len(partsPlus) != 2 {
		return "", "", "", "", fmt.Errorf("invalid target ref %s", fullName)
	}
	if partsPlus[0] == "" {
		// Local target.
		return "", "", ".", partsPlus[1], nil
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
		return "", "", localPath, partsPlus[1], nil
	}

	// Remote target.
	partsColon := strings.SplitN(partsPlus[0], ":", 2)
	if len(partsColon) == 2 {
		tag = partsColon[1]
	}

	return partsColon[0], tag, "", partsPlus[1], nil
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
