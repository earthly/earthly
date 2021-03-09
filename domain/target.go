package domain

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var _ Reference = &Target{}

const targetNamePattern = "^[a-z][a-zA-Z0-9.\\-]*$"

var targetNameRegex = regexp.MustCompile(targetNamePattern)

// Target is an earthly target identifier.
type Target struct {
	GitURL string // e.g. "github.com/earthly/earthly/examples/go"
	Tag    string // e.g. "main"

	// Local representation. E.g. in "./some/path+something" this is "./some/path".
	LocalPath string `json:"localPath"`

	// Target name. E.g. in "+something" this is "something".
	Target string `json:"target"`
}

// GetGitURL returns the GitURL portion of the command.
func (et Target) GetGitURL() string {
	return et.GitURL
}

// GetTag returns the Tag portion of the target.
func (et Target) GetTag() string {
	return et.Tag
}

// GetLocalPath returns the Path portion of the target.
func (et Target) GetLocalPath() string {
	return et.LocalPath
}

// GetName returns the Name portion of the target.
func (et Target) GetName() string {
	return et.Target
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
	return referenceString(et)
}

// StringCanonical returns a string representation of the Target, in canonical form.
func (et Target) StringCanonical() string {
	return referenceStringCanonical(et)
}

// ProjectCanonical returns a string representation of the project of the target, in canonical form.
func (et Target) ProjectCanonical() string {
	return referenceProjectCanonical(et)
}

// ParseTarget parses a string into a Target.
func ParseTarget(fullTargetName string) (Target, error) {
	gitURL, tag, localPath, target, err := parseCommon(fullTargetName)
	if err != nil {
		return Target{}, err
	}
	ok := targetNameRegex.MatchString(target)
	if !ok {
		return Target{}, errors.Errorf("target name %s does not match %s", target, targetNamePattern)
	}
	return Target{
		GitURL:    gitURL,
		Tag:       tag,
		LocalPath: localPath,
		Target:    target,
	}, nil
}
