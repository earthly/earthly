package domain

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Artifact is a earth artifact identifier.
type Artifact struct {
	Target   Target
	Artifact string
}

// String returns a string representation of the Artifact.
func (ea Artifact) String() string {
	return fmt.Sprintf("%s%s", ea.Target.String(), path.Join("/", ea.Artifact))
}

// StringCanonical returns a string representation of the Artifact.
func (ea Artifact) StringCanonical() string {
	return fmt.Sprintf("%s%s", ea.Target.StringCanonical(), path.Join("/", ea.Artifact))
}

// ParseArtifact parses a string representation of a Artifact.
func ParseArtifact(artifactName string) (Artifact, error) {
	parts := strings.SplitN(artifactName, "+", 2)
	if len(parts) != 2 {
		return Artifact{}, fmt.Errorf("Invalid artifact name %s", artifactName)
	}
	partsSlash := strings.SplitN(parts[1], "/", 2)
	if len(partsSlash) != 2 {
		return Artifact{}, fmt.Errorf("Invalid artifact name %s", artifactName)
	}
	earthTargetName := fmt.Sprintf("%s+%s", parts[0], partsSlash[0])
	target, err := ParseTarget(earthTargetName)
	if err != nil {
		return Artifact{}, errors.Wrapf(err, "Invalid artifact name %s", artifactName)
	}
	artifactPath := fmt.Sprintf("/%s", partsSlash[1])
	return Artifact{
		Target:   target,
		Artifact: artifactPath,
	}, nil
}
