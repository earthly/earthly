package states

import (
	"context"
)

// DockerImageSolver can create a Docker image for the WITH DOCKER command.
type DockerImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string, printOutput bool) (chan string, func(), error)
}
