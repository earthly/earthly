package buildcontext

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/syncutil/synccache"

	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// DockerfileMetaTarget is a target name prefix which signals the resolver that the build file is a
// dockerfile. The DockerfileMetaTarget is really not a valid Earthly target otherwise.
const DockerfileMetaTarget = "@dockerfile:"

// Data represents a resolved target's build context data.
type Data struct {
	// The parsed Earthfile AST.
	Earthfile spec.Earthfile
	// BuildFilePath is the local path where the Earthfile or Dockerfile can be found.
	BuildFilePath string
	// BuildContext is the state to use for the build.
	BuildContextFactory llbfactory.Factory
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

	parseCache *synccache.SyncCache // local path -> AST
	console    conslogging.ConsoleLogger
}

// NewResolver returns a new NewResolver.
func NewResolver(sessionID string, cleanCollection *cleanup.Collection, gitLookup *GitLookup, console conslogging.ConsoleLogger) *Resolver {
	return &Resolver{
		gr: &gitResolver{
			cleanCollection: cleanCollection,
			projectCache:    synccache.New(),
			buildFileCache:  synccache.New(),
			gitLookup:       gitLookup,
		},
		lr: &localResolver{
			gitMetaCache: synccache.New(),
			sessionID:    sessionID,
			console:      console,
		},
		parseCache: synccache.New(),
		console:    console,
	}
}

// Resolve returns resolved context data for a given Earthly reference. If the reference is a target,
// then the context will include a build context and possibly additional local directories.
func (r *Resolver) Resolve(ctx context.Context, gwClient gwclient.Client, ref domain.Reference) (*Data, error) {
	if ref.IsUnresolvedImportReference() {
		return nil, errors.Errorf("cannot resolve non-dereferenced import ref %s", ref.String())
	}
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
	if !strings.HasPrefix(ref.GetName(), DockerfileMetaTarget) {
		d.Earthfile, err = r.parseEarthfile(ctx, d.BuildFilePath)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (r *Resolver) parseEarthfile(ctx context.Context, path string) (spec.Earthfile, error) {
	path = filepath.Clean(path)
	efValue, err := r.parseCache.Do(ctx, path, func(ctx context.Context, k interface{}) (interface{}, error) {
		return ast.Parse(ctx, k.(string), true)
	})
	if err != nil {
		return spec.Earthfile{}, err
	}
	ef := efValue.(spec.Earthfile)
	return ef, nil
}
