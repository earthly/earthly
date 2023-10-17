package inputgraph

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/stretchr/testify/require"
)

func TestHashTargetWithDocker(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	org, project, hash, err := HashTarget(ctx, target, cons)
	r.NoError(err)
	r.Equal("earthly-technologies", org)
	r.Equal("core", project)

	hex := fmt.Sprintf("%x", hash)
	r.Equal("9d2903bc18c99831f4a299090abaf94d25d89321", hex)
}

func TestHashTargetWithDockerNoAlias(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load-no-alias",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	org, project, hash, err := HashTarget(ctx, target, cons)
	r.NoError(err)
	r.Equal("earthly-technologies", org)
	r.Equal("core", project)

	hex := fmt.Sprintf("%x", hash)
	r.Equal("d73e37689c7780cbff2cba2de1a23141618b7b14", hex)
}

func TestHashTargetWithDockerArgs(t *testing.T) {
	r := require.New(t)
	target := domain.Target{
		LocalPath: "./testdata/with-docker",
		Target:    "with-docker-load-args",
	}

	ctx := context.Background()
	cons := conslogging.New(os.Stderr, &sync.Mutex{}, conslogging.NoColor, 0, conslogging.Info)

	_, _, _, err := HashTarget(ctx, target, cons)
	r.Error(err)
}
