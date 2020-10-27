package buildcontext

import (
	"context"

	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
)

// DockerfileMetaTarget is a target name which signals the resolver that the build file is a
// dockerfile. The DockerfileMetaTarget is really not a valid Earthly target otherwise.
const DockerfileMetaTarget = "@dockerfile"

// Data represents a resolved target's build context data.
type Data struct {
	// BuildFilePath is the local path where the Earthfile or Dockerfile can be found.
	BuildFilePath string
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
func NewResolver(sessionID string, cleanCollection *cleanup.Collection) *Resolver {
	return &Resolver{
		gr: &gitResolver{
			cleanCollection: cleanCollection,
			projectCache:    make(map[string]*resolvedGitProject),
		},
		lr: &localResolver{
			gitMetaCache: make(map[string]*GitMetadata),
			sessionID:    sessionID,
		},
	}
}

// Resolve returns resolved build context data.
func (r *Resolver) Resolve(ctx context.Context, gwClient gwclient.Client, target domain.Target) (*Data, error) {
	localDirs := make(map[string]string)
	if target.IsRemote() {
		// Remote.
		d, err := r.gr.resolveEarthProject(ctx, gwClient, target)
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
