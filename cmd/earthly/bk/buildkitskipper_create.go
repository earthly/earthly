package bk

import (
	"context"

	"github.com/earthly/earthly/util/buildkitskipper"
	"github.com/pkg/errors"
)

// BuildkitSkipper adds new auto-skip hashes to the backing datastore & allows
// us to check for their existence.
type BuildkitSkipper interface {
	Add(ctx context.Context, org, target string, key []byte) error
	Exists(ctx context.Context, org string, key []byte) (bool, error)
}

// NewBuildkitSkipper returns a local buildkitskipper when localSkipDB is specified, or alternatively a cloud-based skipper
func NewBuildkitSkipper(localSkipDB string, cloudClient buildkitskipper.ASKVClient) (BuildkitSkipper, error) {
	if localSkipDB != "" {
		skipDB, err := buildkitskipper.NewLocal(localSkipDB)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open buildkit skipper database %s", localSkipDB)
		}
		return skipDB, nil
	}

	skipDB, err := buildkitskipper.NewCloud(cloudClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cloud-based buildkit skipper database")
	}

	return skipDB, nil
}
