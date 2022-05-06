package states

import (
	"context"

	"github.com/earthly/earthly/util/platutil"
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

// ImageSolverResults contains data and channels that allow one to act on images
// during and after they are built.
type ImageSolverResults struct {
	ResultChan  chan string
	ErrChan     chan error
	ReleaseFunc func()
}

// ImageDef includes the information required to build an image in BuildKit.
type ImageDef struct {
	MTS       *MultiTarget
	ImageName string
	Platform  platutil.Platform
}

// MultiImageSolver can create a Docker image for the WITH DOCKER command using
// the embedded BuildKit registry.
type MultiImageSolver interface {
	SolveImages(ctx context.Context, defs []*ImageDef) (*ImageSolverResults, error)
}
