package states

import (
	"context"
)

// DockerTarImageSolver can create a Docker image and make it available as a tar
// file.
type DockerTarImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string, printOutput bool) error
}

// DockerImageSolver can create a Docker image for the WITH DOCKER command using
// the embedded BuildKit registry.
type DockerImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string) (chan string, func(), chan error, error)
}
