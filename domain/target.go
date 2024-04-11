package domain

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var _ Reference = Target{}

const targetNamePattern = "^[a-z][a-zA-Z0-9.\\-]*$"

var targetNameRegex = regexp.MustCompile(targetNamePattern)

// Target is an earthly target identifier.
type Target struct {
	// Remote representation.
	GitURL string `json:"gitUrl"` // e.g. "github.com/earthly/earthly/examples/tutorial/go/part3"
	Tag    string `json:"tag"`    // e.g. "main"
	// Local representation. E.g. in "./some/path+something" this is "./some/path".
	LocalPath string `json:"localPath"`
	// Import representation. E.g. in "foo+bar" this is "foo".
	ImportRef string `json:"importRef"`

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

// GetImportRef returns the ImportRef portion of the target.
func (et Target) GetImportRef() string {
	return et.ImportRef
}

// GetName returns the Name portion of the target.
func (et Target) GetName() string {
	return et.Target
}

// IsExternal returns whether the target is external to the current project.
func (et Target) IsExternal() bool {
	return et.IsRemote() || et.IsLocalExternal() || et.IsImportReference()
}

// IsLocalInternal returns whether the target is in the same Earthfile.
func (et Target) IsLocalInternal() bool {
	return et.LocalPath == "."
}

// IsLocalExternal returns whether the target is a local, but external target.
func (et Target) IsLocalExternal() bool {
	return et.LocalPath != "." && et.LocalPath != ""
}

// IsRemote returns whether the target is remote.
func (et Target) IsRemote() bool {
	return et.GitURL != "" && !et.IsLocalInternal() && !et.IsLocalExternal()
}

// IsImportReference returns whether the target is a reference to an import.
func (et Target) IsImportReference() bool {
	return et.ImportRef != ""
}

// IsUnresolvedImportReference returns whether the target is an import reference that has
// no remote or local information set.
func (et Target) IsUnresolvedImportReference() bool {
	return et.IsImportReference() && !et.IsRemote() && !et.IsLocalExternal() && !et.IsLocalInternal()
}

// DebugString returns a string that can be printed out for debugging purposes
func (et Target) DebugString() string {
	return fmt.Sprintf("GitURL: %q; Tag: %q; LocalPath: %q; ImportRef: %q; Target: %q", et.GitURL, et.Tag, et.LocalPath, et.ImportRef, et.Target)
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
	gitURL, tag, localPath, importRef, target, err := parseCommon(fullTargetName)
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
		ImportRef: importRef,
		Target:    target,
	}, nil
}
