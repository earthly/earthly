package bk

import (
	"context"

	"github.com/earthly/earthly/util/buildkitskipper"
	"github.com/pkg/errors"
)

type BuildkitSkipper interface {
	Add(ctx context.Context, key []byte) error
	Exists(ctx context.Context, key []byte) (bool, error)
}

// NewBuildkitSkipper returns a local buildkitskipper when localSkipDB is specified, or alternatively a cloud-based skipper
func NewBuildkitSkipper(localSkipDB, orgName, projectName, pipelineName string, cloudClient buildkitskipper.ASKVClient) (BuildkitSkipper, error) {
	if localSkipDB != "" {
		skipDB, err := buildkitskipper.NewLocal(localSkipDB)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open buildkit skipper database %s", localSkipDB)
		}
		return skipDB, nil
	}

	skipDB, err := buildkitskipper.NewCloud(orgName, projectName, pipelineName, cloudClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cloud-based buildkit skipper database")
	}
	return skipDB, nil
}
