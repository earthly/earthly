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
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/synccache"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	buildkitgitutil "github.com/moby/buildkit/util/gitutil"
	"github.com/pkg/errors"
)

// DockerfileMetaTarget is a target name prefix which signals the resolver that the build file is a
// dockerfile. The DockerfileMetaTarget is really not a valid Earthly target otherwise.
const DockerfileMetaTarget = "@dockerfile:"

// Data represents a resolved target's build context data.
type Data struct {
	// The parsed Earthfile AST.
	Earthfile spec.Earthfile
	// EarthlyOrgName is the org that the target belongs to.
	EarthlyOrgName string
	// EarthlyProjectName is the project that the target belongs to.
	EarthlyProjectName string
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
	// Features holds the feature state for the build context
	Features *features.Features
}

// Resolver is a build context resolver.
type Resolver struct {
	gr *gitResolver
	lr *localResolver

	parseCache *synccache.SyncCache // local path -> AST
	console    conslogging.ConsoleLogger

	featureFlagOverrides string
}

// NewResolver returns a new NewResolver.
func NewResolver(cleanCollection *cleanup.Collection, gitLookup *GitLookup, console conslogging.ConsoleLogger, featureFlagOverrides, gitBranchOverride, gitLFSInclude string, gitLogLevel buildkitgitutil.GitLogLevel, gitImage string) *Resolver {
	return &Resolver{
		gr: &gitResolver{
			gitBranchOverride: gitBranchOverride,
			gitImage:          gitImage,
			lfsInclude:        gitLFSInclude,
			logLevel:          gitLogLevel,
			cleanCollection:   cleanCollection,
			projectCache:      synccache.New(),
			buildFileCache:    synccache.New(),
			gitLookup:         gitLookup,
			console:           console,
		},
		lr: &localResolver{
			buildFileCache:    synccache.New(),
			gitMetaCache:      synccache.New(),
			gitBranchOverride: gitBranchOverride,
			console:           console,
		},
		parseCache:           synccache.New(),
		console:              console,
		featureFlagOverrides: featureFlagOverrides,
	}
}

// ExpandWildcard will expand a wildcard BUILD target in a local path or remote
// Git repository. The pattern is the path relative to the target path and
// should be in the form 'my/path/*'
func (r *Resolver) ExpandWildcard(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, rootTarget, target domain.Target, pattern string) ([]string, error) {

	if target.IsRemote() {
		matches, err := r.gr.expandWildcard(ctx, gwClient, platr, target, pattern)
		if err != nil {
			return nil, errors.Wrap(err, "failed to expand remote BUILD target path")
		}
		return matches, nil
	}

	abs, err := filepath.Abs(rootTarget.LocalPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not determine absolute path for %q", rootTarget.LocalPath)
	}

	parent := filepath.Join(abs, target.GetLocalPath())

	matches, err := fileutil.GlobDirs(parent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to expand BUILD target path")
	}

	return matches, nil
}

// Resolve returns resolved context data for a given Earthly reference. If the reference is a target,
// then the context will include a build context and possibly additional local directories.
func (r *Resolver) Resolve(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, ref domain.Reference) (*Data, error) {
	if ref.IsUnresolvedImportReference() {
		return nil, errors.Errorf("cannot resolve non-dereferenced import ref %s", ref.String())
	}
	var d *Data
	var err error
	localDirs := make(map[string]string)
	if ref.IsRemote() {
		// Remote.
		d, err = r.gr.resolveEarthProject(ctx, gwClient, platr, ref, r.featureFlagOverrides)
		if err != nil {
			return nil, err
		}
	} else {
		// Local.
		if _, isTarget := ref.(domain.Target); isTarget {
			localDirs[ref.GetLocalPath()] = ref.GetLocalPath()
		}

		d, err = r.lr.resolveLocal(ctx, gwClient, platr, ref, r.featureFlagOverrides)
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
		org, project, err := extractOrgAndProjectName(d.Earthfile)
		if err != nil {
			return nil, err
		}
		d.EarthlyOrgName = org
		d.EarthlyProjectName = project
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

func extractOrgAndProjectName(ef spec.Earthfile) (string, string, error) {
	for _, cmd := range ef.BaseRecipe {
		if cmd.Command == nil {
			continue
		}
		if cmd.Command.Name != "PROJECT" {
			continue
		}
		if len(cmd.Command.Args) != 1 {
			return "", "", errors.Errorf("invalid PROJECT command")
		}
		orgProj := cmd.Command.Args[0]
		parts := strings.SplitN(orgProj, "/", 2)
		if len(parts) != 2 {
			return "", "", errors.Errorf("invalid PROJECT command")
		}
		return parts[0], parts[1], nil
	}
	return "", "", nil
}
