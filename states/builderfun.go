package states

import (
	"context"

	"github.com/earthly/earthly/util/platutil"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerTarImageSolver can create a Docker image and make it available as a tar
// file.
type DockerTarImageSolver interface {
	SolveImage(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string, printOutput bool) error
}

// ImageSolverResults contains data and channels that allow one to act on images
// during and after they are built.
type ImageSolverResults struct {
	ResultChan  chan *ImageResult
	ErrChan     chan error
	ReleaseFunc func()
}

// ImageResult contains data about an image that was built.
type ImageResult struct {
	IntermediateImageName    string
	FinalImageName           string
	FinalImageNameWithDigest string
	ImageDigest              string
	ConfigDigest             string
	ImageDescriptor          *ocispecs.Descriptor
	Annotations              map[string]string
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
