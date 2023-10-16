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
	r.Equal("b644ef18314c4d28c98a26ff21c593ac25a21c86", hex)
}

func TestHashTargetWithDockerImplicit(t *testing.T) {
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
	r.Equal("16f68555b6894a3e757505546ffaccdcc964cc07", hex)
}
