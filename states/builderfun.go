package states

import (
	"context"
)

// DockerTarImageSolver can create a Docker image and make it available as a tar
// file.
type DockerTarImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string, printOutput bool) error
}

// ImageSolverResult contains data and channels that allow one to act on images
// during and after they are built.
type ImageSolverResult struct {
	ImageName   string
	ResultChan  chan string
	ErrChan     chan error
	ReleaseFunc func()
}

// DockerImageSolver can create a Docker image for the WITH DOCKER command using
// the embedded BuildKit registry.
type DockerImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string) (*ImageSolverResult, error)
}
