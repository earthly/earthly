package states

import (
	"context"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerBuilderFun is a function able to build a target into a docker tar file.
type DockerBuilderFun = func(ctx context.Context, mts *MultiTarget, nativePlatform specs.Platform, dockerTag string, outFile string) error
