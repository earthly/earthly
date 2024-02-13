package platutil

import (
	"context"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// GetNativePlatform returns the native platform for a given gwClient.
func GetNativePlatform(gwClient gwclient.Client) (specs.Platform, error) {
	ws := gwClient.BuildOpts().Workers
	if len(ws) == 0 {
		return specs.Platform{}, errors.New("no worker found via gwclient")
	}
	nps := ws[0].Platforms
	if len(nps) == 0 {
		return specs.Platform{}, errors.New("no platform found for worker via gwclient")
	}
	return platforms.Normalize(nps[0]), nil
}

// GetNativePlatformViaBkClient returns the native platform for a given buildkit client.
func GetNativePlatformViaBkClient(ctx context.Context, bkClient *client.Client) (specs.Platform, error) {
	ws, err := bkClient.ListWorkers(ctx)
	if err != nil {
		return specs.Platform{}, errors.Wrap(err, "failed to list workers")
	}
	if len(ws) == 0 {
		return specs.Platform{}, errors.New("no worker found via bkClient")
	}
	nps := ws[0].Platforms
	if len(nps) == 0 {
		return specs.Platform{}, errors.New("no platform found for worker via bkClient")
	}
	return platforms.Normalize(nps[0]), nil
}
