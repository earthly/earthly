package buildkitskipper

import "context"

type BuildkitSkipper interface {
	Add(ctx context.Context, key []byte) error
	Exists(ctx context.Context, key []byte) (bool, error)
}
