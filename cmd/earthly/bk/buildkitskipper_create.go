package bk

import (
	"context"

	"github.com/earthly/earthly/util/buildkitskipper"
	"github.com/pkg/errors"
)

// BuildkitSkipper adds new auto-skip hashes to the backing datastore & allows
// us to check for their existence.
type BuildkitSkipper interface {
	Add(ctx context.Context, target string, key []byte) error
	Exists(ctx context.Context, key []byte) (bool, error)
}

// NewBuildkitSkipper returns a local buildkitskipper when localSkipDB is specified
func NewBuildkitSkipper(localSkipDB string) (BuildkitSkipper, error) {
	if localSkipDB == "" {
		return nil, nil // will disable autoskipper.
	}

	skipDB, err := buildkitskipper.NewLocal(localSkipDB)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open buildkit skipper database %s", localSkipDB)
	}

	return skipDB, nil
}
