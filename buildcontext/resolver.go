package buildcontext

import (
	"context"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
)

// Data represents a resolved target's build context data.
type Data struct {
	// EarthfilePath is the local path where the build.earth file can be found.
	EarthfilePath string
	// BuildContext is the state to use for the build.
	BuildContext llb.State
	// GitMetadata contains git metadata information.
	GitMetadata *GitMetadata
	// Target is the earthly target.
	Target domain.Target
	// LocalDirs is the local dirs map to be passed as part of the buildkit solve.
	LocalDirs map[string]string
}

// Resolver is a build context resolver.
type Resolver struct {
	gr *gitResolver
	lr *localResolver
}

// NewResolver returns a new NewResolver.
func NewResolver(bkClient *client.Client, console conslogging.ConsoleLogger, sessionID string) *Resolver {
	return &Resolver{
		gr: &gitResolver{
			bkClient:     bkClient,
			console:      console,
			projectCache: make(map[string]*resolvedGitProject),
		},
		lr: &localResolver{
			gitMetaCache: make(map[string]*GitMetadata),
			sessionID:    sessionID,
		},
	}
}

// Resolve returns resolved build context data.
func (r *Resolver) Resolve(ctx context.Context, target domain.Target) (*Data, error) {
	localDirs := make(map[string]string)
	if target.IsRemote() {
		// Remote.
		d, err := r.gr.resolveEarthProject(ctx, target)
		if err != nil {
			return nil, err
		}

		d.Target = TargetWithGitMeta(target, d.GitMetadata)
		d.LocalDirs = localDirs
		return d, nil
	}

	// Local.
	localDirs[target.LocalPath] = target.LocalPath
	d, err := r.lr.resolveLocal(ctx, target)
	if err != nil {
		return nil, err
	}
	d.Target = TargetWithGitMeta(target, d.GitMetadata)
	d.LocalDirs = localDirs
	return d, nil
}

// Close closes the resolver, freeing up any internal resources.
func (r *Resolver) Close() error {
	return r.gr.close()
}
