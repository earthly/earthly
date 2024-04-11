package domain

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

var _ Reference = Command{}

const commandNamePattern = "^[A-Z][A-Z0-9._]*$"

var commandNameRegex = regexp.MustCompile(commandNamePattern)

// Command is an earthly command identifier.
type Command struct {
	// Remote representation.
	GitURL string `json:"gitUrl"` // e.g. "github.com/earthly/earthly/examples/tutorial/go/part3"
	Tag    string `json:"tag"`    // e.g. "main"
	// Local representation. E.g. in "./some/path+something" this is "./some/path".
	LocalPath string `json:"localPath"`
	// Import representation. E.g. in "foo+bar" this is "foo".
	ImportRef string `json:"importRef"`

	// Command name. E.g. in "+SOMETHING" this is "SOMETHING".
	Command string `json:"command"`
}

// GetGitURL returns the GitURL portion of the command.
func (ec Command) GetGitURL() string {
	return ec.GitURL
}

// GetTag returns the Tag portion of the command.
func (ec Command) GetTag() string {
	return ec.Tag
}

// GetLocalPath returns the Path portion of the command.
func (ec Command) GetLocalPath() string {
	return ec.LocalPath
}

// GetImportRef returns the ImportRef portion of the command.
func (ec Command) GetImportRef() string {
	return ec.ImportRef
}

// GetName returns the Name portion of the command.
func (ec Command) GetName() string {
	return ec.Command
}

// IsExternal returns whether the command is external to the current project.
func (ec Command) IsExternal() bool {
	return ec.IsRemote() || ec.IsLocalExternal() || ec.IsImportReference()
}

// IsLocalInternal returns whether the command is in the same Earthfile.
func (ec Command) IsLocalInternal() bool {
	return ec.LocalPath == "."
}

// IsLocalExternal returns whether the command is a local, but external command.
func (ec Command) IsLocalExternal() bool {
	return ec.LocalPath != "." && ec.LocalPath != ""
}

// IsRemote returns whether the command is remote.
func (ec Command) IsRemote() bool {
	return ec.GitURL != "" && !ec.IsLocalInternal() && !ec.IsLocalExternal()
}

// IsImportReference returns whether the target is a reference to an import.
func (ec Command) IsImportReference() bool {
	return ec.ImportRef != ""
}

// IsUnresolvedImportReference returns whether the command is an import reference that has
// no remote or local information set.
func (ec Command) IsUnresolvedImportReference() bool {
	return ec.IsImportReference() && !ec.IsRemote() && !ec.IsLocalExternal() && !ec.IsLocalInternal()
}

// DebugString returns a string that can be printed out for debugging purposes
func (ec Command) DebugString() string {
	return fmt.Sprintf("gitURL: %q; tag: %q; LocalPath: %q; ImportRef: %q; Command: %q", ec.GitURL, ec.Tag, ec.LocalPath, ec.ImportRef, ec.Command)
}

// String returns a string representation of the command.
func (ec Command) String() string {
	return referenceString(ec)
}

// StringCanonical returns a string representation of the command, in canonical form.
func (ec Command) StringCanonical() string {
	return referenceStringCanonical(ec)
}

// ProjectCanonical returns a string representation of the project of the command, in canonical form.
func (ec Command) ProjectCanonical() string {
	return referenceProjectCanonical(ec)
}

// ParseCommand parses a string into a Command.
func ParseCommand(fullCommandName string) (Command, error) {
	gitURL, tag, localPath, importRef, command, err := parseCommon(fullCommandName)
	if err != nil {
		return Command{}, err
	}
	ok := commandNameRegex.MatchString(command)
	if !ok {
		return Command{}, errors.Errorf("command name %s does not match %s", command, commandNamePattern)
	}
	return Command{
		GitURL:    gitURL,
		Tag:       tag,
		LocalPath: localPath,
		ImportRef: importRef,
		Command:   command,
	}, nil
}
