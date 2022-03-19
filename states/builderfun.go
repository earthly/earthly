package states

import (
	"context"
)

// DockerBuilderFun is a function able to build a target into a docker tar file.
type DockerBuilderFun = func(ctx context.Context, mts *MultiTarget, dockerTag string, outFile string, printOutput bool) error
