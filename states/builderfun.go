package states

import (
	"context"

	"github.com/earthly/earthly/domain"
)

// DockerBuilderFun is a function able to build a target into a docker tar file.
type DockerBuilderFun = func(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string) error

// ArtifactBuilderFun is a function able to build an artifact and output it locally.
type ArtifactBuilderFun = func(ctx context.Context, mts *MultiTarget, artifact domain.Artifact, outFile string) error
