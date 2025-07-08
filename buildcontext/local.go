package buildcontext

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/synccache"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

type localResolver struct {
	gitMetaCache      *synccache.SyncCache // local path -> *gitutil.GitMetadata
	gitBranchOverride string
	buildFileCache    *synccache.SyncCache
	console           conslogging.ConsoleLogger
}

func (lr *localResolver) resolveLocal(ctx context.Context, gwClient gwclient.Client, platr *platutil.Resolver, ref domain.Reference, featureFlagOverrides string) (*Data, error) {
	if ref.IsRemote() {
		return nil, errors.Errorf("unexpected remote target %s", ref.String())
	}

	metadataValue, err := lr.gitMetaCache.Do(ctx, ref.GetLocalPath(), func(ctx context.Context, _ interface{}) (interface{}, error) {
		metadata, err := gitutil.Metadata(ctx, ref.GetLocalPath(), lr.gitBranchOverride)
		if err != nil {
			if errors.Is(err, gitutil.ErrNoGitBinary) ||
				errors.Is(err, gitutil.ErrNotAGitDir) ||
				errors.Is(err, gitutil.ErrCouldNotDetectRemote) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitHash) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitShortHash) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitBranch) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitTags) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitRefs) {
				// Keep going anyway. Either not a git dir, or git not installed, or
				// remote not detected.
				if errors.Is(err, gitutil.ErrNoGitBinary) {
					lr.console.Warnf("Warning: %s\n", err.Error())
				}
			} else {
				return nil, err
			}
		}
		return metadata, nil
	})
	if err != nil {
		return nil, err
	}
	metadata := metadataValue.(*gitutil.GitMetadata)

	localPath := filepath.FromSlash(ref.GetLocalPath())
	key := localPath
	isDockerfile := strings.HasPrefix(ref.GetName(), DockerfileMetaTarget)
	if isDockerfile {
		// Different key for dockerfiles to include the dockerfile name itself.
		key = ref.String()
	}
	buildFileValue, err := lr.buildFileCache.Do(ctx, key, func(ctx context.Context, _ interface{}) (interface{}, error) {
		buildFilePath, err := detectBuildFile(ref, localPath)
		if err != nil {
			return nil, err
		}
		var ftrs *features.Features
		if isDockerfile {
			ftrs = new(features.Features)
		} else {
			ftrs, err = parseFeatures(buildFilePath, featureFlagOverrides, ref.GetLocalPath(), lr.console)
			if err != nil {
				return nil, err
			}
		}
		return &buildFile{
			path: buildFilePath,
			ftrs: ftrs,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	bf := buildFileValue.(*buildFile)

	var buildContextFactory llbfactory.Factory
	// guard against auto-complete code's GetTargetArgs() func which passes in a nil gwClient (but doesn't actually invoke buildkit)
	if gwClient != nil {
		if _, isTarget := ref.(domain.Target); isTarget {
			noImplicitIgnore := bf.ftrs != nil && bf.ftrs.NoImplicitIgnore

			useDockerIgnore := isDockerfile
			ftrs := features.FromContext(ctx)
			if ftrs != nil {
				useDockerIgnore = useDockerIgnore && ftrs.UseDockerIgnore
			}

			excludes, err := readExcludes(ref.GetLocalPath(), noImplicitIgnore, useDockerIgnore)
			if err != nil {
				return nil, err
			}
			buildContextFactory = llbfactory.Local(
				ref.GetLocalPath(),
				llb.ExcludePatterns(excludes),
				llb.Platform(platr.LLBNative()),
				llb.WithCustomNamef("[context %s] local context %s", ref.GetLocalPath(), ref.GetLocalPath()),
			)
		}
		// Else not needed: Commands don't come with a build context.
	}

	return &Data{
		BuildFilePath:       bf.path,
		BuildContextFactory: buildContextFactory,
		GitMetadata:         metadata,
		Features:            bf.ftrs,
	}, nil
}
