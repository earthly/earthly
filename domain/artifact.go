package domain

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Artifact is an earthly artifact identifier.
type Artifact struct {
	Target   Target
	Artifact string
}

// Clone returns a copy of the Artifact
func (a Artifact) Clone() Artifact {
	newArtifact := a
	return newArtifact
}

// String returns a string representation of the Artifact.
func (ea Artifact) String() string {
	return fmt.Sprintf("%s%s", ea.Target.String(), path.Join("/", escapePlus(ea.Artifact)))
}

// StringCanonical returns a string representation of the Artifact.
func (ea Artifact) StringCanonical() string {
	return fmt.Sprintf("%s%s", ea.Target.StringCanonical(), path.Join("/", escapePlus(ea.Artifact)))
}

// ParseArtifact parses a string representation of an Artifact.
func ParseArtifact(artifactName string) (Artifact, error) {
	parts, err := splitUnescapePlus(artifactName)
	if err != nil {
		return Artifact{}, err
	}
	if len(parts) != 2 {
		return Artifact{}, errors.Errorf("invalid artifact name %s", artifactName)
	}
	partsSlash := strings.SplitN(parts[1], "/", 2)
	if len(partsSlash) != 2 {
		return Artifact{}, errors.Errorf("invalid artifact name %s", artifactName)
	}
	earthTargetName := fmt.Sprintf("%s+%s", escapePlus(parts[0]), partsSlash[0])
	target, err := ParseTarget(earthTargetName)
	if err != nil {
		return Artifact{}, errors.Wrapf(err, "invalid artifact name %s", artifactName)
	}
	artifactPath := fmt.Sprintf("/%s", partsSlash[1])
	return Artifact{
		Target:   target,
		Artifact: artifactPath,
	}, nil
}
