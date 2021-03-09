package buildcontext

import (
	"context"
	"path/filepath"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
)

// DockerfileMetaTarget is a target name which signals the resolver that the build file is a
// dockerfile. The DockerfileMetaTarget is really not a valid Earthly target otherwise.
const DockerfileMetaTarget = "@dockerfile"

// Data represents a resolved target's build context data.
type Data struct {
	// The parsed Earthfile AST.
	Earthfile spec.Earthfile
	// BuildFilePath is the local path where the Earthfile or Dockerfile can be found.
	BuildFilePath string
	// BuildContext is the state to use for the build.
	BuildContext llb.State
	// GitMetadata contains git metadata information.
	GitMetadata *gitutil.GitMetadata
	// Target is the earthly reference.
	Ref domain.Reference
	// LocalDirs is the local dirs map to be passed as part of the buildkit solve.
	LocalDirs map[string]string
}

// Resolver is a build context resolver.
type Resolver struct {
	gr *gitResolver
	lr *localResolver

	parseCache map[string]spec.Earthfile // local path -> AST
}

// NewResolver returns a new NewResolver.
func NewResolver(sessionID string, cleanCollection *cleanup.Collection, gitLookup *GitLookup) *Resolver {
	return &Resolver{
		gr: &gitResolver{
			cleanCollection: cleanCollection,
			projectCache:    make(map[string]*resolvedGitProject),
			gitLookup:       gitLookup,
		},
		lr: &localResolver{
			gitMetaCache: make(map[string]*gitutil.GitMetadata),
			sessionID:    sessionID,
		},
		parseCache: make(map[string]spec.Earthfile),
	}
}

// Resolve returns resolved context data for a given Earthly reference. If the reference is a target,
// then the context will include a build context and possibly additional local directories.
func (r *Resolver) Resolve(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (*Data, error) {
	var d *Data
	var err error
	localDirs := make(map[string]string)
	if ref.IsRemote() {
		// Remote.
		d, err = r.gr.resolveEarthProject(ctx, gwClient, ref)
		if err != nil {
			return nil, err
		}
	} else {
		// Local.
		if _, isTarget := ref.(domain.Target); isTarget {
			localDirs[ref.GetLocalPath()] = ref.GetLocalPath()
		}
		d, err = r.lr.resolveLocal(ctx, ref)
		if err != nil {
			return nil, err
		}
	}
	d.Ref = gitutil.ReferenceWithGitMeta(ref, d.GitMetadata)
	d.LocalDirs = localDirs
	if ref.GetName() != DockerfileMetaTarget {
		d.Earthfile, err = r.parseEarthfile(ctx, d.BuildFilePath)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (r *Resolver) parseEarthfile(ctx context.Context, path string) (spec.Earthfile, error) {
	path = filepath.Clean(path)
	ef, found := r.parseCache[path]
	if found {
		return ef, nil
	}

	ef, err := ast.Parse(ctx, path, true)
	if err != nil {
		return spec.Earthfile{}, err
	}
	r.parseCache[path] = ef
	return ef, nil
}
