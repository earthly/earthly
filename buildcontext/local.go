package buildcontext

import (
	"context"
	"path/filepath"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/llbfactory"
	"github.com/earthly/earthly/util/syncutil/synccache"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type localResolver struct {
	gitMetaCache         *synccache.SyncCache // local path -> *gitutil.GitMetadata
	sessionID            string
	console              conslogging.ConsoleLogger
	featureFlagOverrides string
}

func (lr *localResolver) resolveLocal(ctx context.Context, ref domain.Reference) (*Data, error) {
	analytics.Count("localResolver.resolveLocal", "local-reference")
	if ref.IsRemote() {
		return nil, errors.Errorf("unexpected remote target %s", ref.String())
	}

	metadataValue, err := lr.gitMetaCache.Do(ctx, ref.GetLocalPath(), func(ctx context.Context, _ interface{}) (interface{}, error) {
		metadata, err := gitutil.Metadata(ctx, ref.GetLocalPath())
		if err != nil {
			if errors.Is(err, gitutil.ErrNoGitBinary) ||
				errors.Is(err, gitutil.ErrNotAGitDir) ||
				errors.Is(err, gitutil.ErrCouldNotDetectRemote) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitHash) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitBranch) {
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

	buildFilePath, err := detectBuildFile(ref, filepath.FromSlash(ref.GetLocalPath()))
	if err != nil {
		return nil, err
	}

	// Resolver is responsible for parsing the Earthfile later, however, we need to pre-emptively parse the
	// version only so we can check the Earthfile version and feature gates. This is so the local resolver knows
	// which files to exclude based on certain feature flags like --no-implicit-ignore.
	version, err := ast.ParseVersion(buildFilePath, false)
	if err != nil {
		return nil, err
	}

	ftrs, err := features.GetFeatures(version)
	if err != nil {
		return nil, err
	}

	err = features.ApplyFlagOverrides(ftrs, lr.featureFlagOverrides)
	if err != nil {
		return nil, err
	}

	var buildContextFactory llbfactory.Factory
	if _, isTarget := ref.(domain.Target); isTarget {
		excludes, err := readExcludes(ref.GetLocalPath(), ftrs.NoImplicitIgnore)
		if err != nil {
			return nil, err
		}
		buildContextFactory = llbfactory.Local(
			ref.GetLocalPath(),
			llb.ExcludePatterns(excludes),
			llb.SessionID(lr.sessionID),
			llb.Platform(llbutil.DefaultPlatform()),
			llb.WithCustomNamef("[context %s] local context %s", ref.GetLocalPath(), ref.GetLocalPath()),
		)
	} else {
		// Commands don't come with a build context.
	}

	return &Data{
		BuildFilePath:       buildFilePath,
		BuildContextFactory: buildContextFactory,
		GitMetadata:         metadata,
	}, nil
}
