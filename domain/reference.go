package domain

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Reference is a target or a command reference.
type Reference interface {
	// GetGitURL is the git URL part of the reference. E.g. "github.com/earthly/earthly/examples/tutorial/go/part3"
	GetGitURL() string
	// GetTag is the git tag of the reference. E.g. "main"
	GetTag() string
	// GetLocalPath is the local path representation of the reference. E.g. in "./some/path+something" this is "./some/path".
	GetLocalPath() string
	// GetImportRef is the import identifier. E.g. in "foo+bar" this is "foo".
	GetImportRef() string
	// GetName is the target name or the command name of the reference. E.g. in "+something" this is "something".
	GetName() string

	// IsExternal returns whether the ref is external to the current project.
	IsExternal() bool

	// IsLocalInternal returns whether the ref is a local.
	IsLocalInternal() bool

	// IsLocalExternal returns whether the ref is a local, but external target.
	IsLocalExternal() bool

	// IsRemote returns whether the ref is remote.
	IsRemote() bool

	// IsImportReference returns whether the ref is a reference to an import.
	IsImportReference() bool
	// IsUnresolvedImportReference returns whether the ref is an import reference that has
	// no remote or local information set.
	IsUnresolvedImportReference() bool

	// DebugString returns a string that can be printed out for debugging purposes.
	DebugString() string

	// String returns a string representation of the ref.
	String() string

	// StringCanonical returns a string representation of the ref, in canonical form.
	StringCanonical() string
	// ProjectCanonical returns a string representation of the project of the ref, in canonical form.
	ProjectCanonical() string
}

// JoinReferences returns the result of interpreting r2 as relative to r1. The result will always have
// the same type as r2.
func JoinReferences(r1 Reference, r2 Reference) (Reference, error) {
	if r1.IsUnresolvedImportReference() || r2.IsUnresolvedImportReference() {
		return nil, errors.New("unresolved import references cannot be joined")
	}
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
				if !(strings.HasPrefix(localPath, ".") || strings.HasPrefix(localPath, "/")) {
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
			ImportRef: r2.GetImportRef(),
			Target:    name,
		}, nil
	case Command:
		return Command{
			GitURL:    gitURL,
			Tag:       tag,
			LocalPath: localPath,
			ImportRef: r2.GetImportRef(),
			Command:   name,
		}, nil
	default:
		return nil, errors.Errorf("joining references not supported for type %T", r2)
	}
}

func referenceString(r Reference) string {
	if r.IsImportReference() {
		return fmt.Sprintf("%s+%s", escapePlus(r.GetImportRef()), r.GetName())
	}
	if r.IsLocalExternal() {
		return fmt.Sprintf("%s+%s", escapePlus(r.GetLocalPath()), r.GetName())
	}
	if r.IsRemote() {
		s := escapePlus(r.GetGitURL())
		if r.GetTag() != "" {
			s += ":" + escapePlus(r.GetTag())
		}
		s += "+" + r.GetName()
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
		s += "+" + r.GetName()
		return s
	}
	if r.GetLocalPath() == "." {
		return fmt.Sprintf("+%s", r.GetName())
	}
	if r.GetLocalPath() == "" && r.GetImportRef() != "" {
		return fmt.Sprintf("%s+%s", escapePlus(r.GetImportRef()), r.GetName())
	}
	// Local external.
	return fmt.Sprintf("%s+%s", escapePlus(r.GetLocalPath()), r.GetName())
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
	if r.GetLocalPath() == "" && r.GetImportRef() != "" {
		return escapePlus(r.GetImportRef())
	}
	return escapePlus(r.GetLocalPath())
}

func parseCommon(fullName string) (gitURL string, tag string, localPath string, importRef string, name string, err error) {
	partsPlus, err := splitUnescapePlus(fullName)
	if err != nil {
		return "", "", "", "", "", err
	}
	if len(partsPlus) != 2 {
		return "", "", "", "", "", errors.Errorf("invalid target ref %s", fullName)
	}
	if partsPlus[0] == "" {
		// Local target.
		return "", "", ".", "", partsPlus[1], nil
	} else if strings.HasPrefix(partsPlus[0], ".") || filepath.IsAbs(partsPlus[0]) {
		// Local external target.
		localPath := partsPlus[0]
		if filepath.IsAbs(localPath) {
			localPath = path.Clean(localPath)
		} else {
			localPath = path.Clean(localPath)
			if !strings.HasPrefix(localPath, ".") {
				localPath = fmt.Sprintf("./%s", localPath)
			}
		}
		return "", "", localPath, "", partsPlus[1], nil
	}

	if strings.ContainsAny(partsPlus[0], "/:") {
		// Remote target.
		partsColon := strings.SplitN(partsPlus[0], ":", 2)
		if len(partsColon) == 2 {
			tag = partsColon[1]
		}

		return partsColon[0], tag, "", "", partsPlus[1], nil
	}

	// Import reference.
	return "", "", "", partsPlus[0], partsPlus[1], nil
}

// splitUnescapePlus performs a split on "+" to return the target and path separately (i.e. always an array of 2).
// The function accounts for escaping a "\+" in a path before a target (which also begins with "+").
// For example, a path to target that contains "+" like `/my/some\+dir+my-target` will be returned as
// [ "/my/some\+dir", "my-target" ]. Special care is given for other cases where backslash might be used,
// such as an escaped whitespace "\ ", or escaped back-slash "\\" or a Windows path with backslashes.
func splitUnescapePlus(str string) ([]string, error) {
	escape := false
	ret := make([]string, 0, 2)
	word := make([]rune, 0, len(str))
	for _, c := range str {
		if escape {
			if c != '+' {
				word = append(word, '\\')
			}
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
