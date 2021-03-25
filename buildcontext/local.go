package buildcontext

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/gitutil"
	"github.com/earthly/earthly/llbutil"

	"github.com/moby/buildkit/client/llb"
)

type localResolver struct {
	gitMetaCache map[string]*gitutil.GitMetadata
	sessionID    string
}

func (lr *localResolver) resolveLocal(ctx context.Context, ref domain.Reference) (*Data, error) {
	if ref.IsRemote() {
		return nil, errors.Errorf("unexpected remote target %s", ref.String())
	}

	metadata, found := lr.gitMetaCache[ref.GetLocalPath()]
	if !found {
		var err error
		metadata, err = gitutil.Metadata(ctx, ref.GetLocalPath())
		if err != nil {
			if errors.Is(err, gitutil.ErrNoGitBinary) ||
				errors.Is(err, gitutil.ErrNotAGitDir) ||
				errors.Is(err, gitutil.ErrCouldNotDetectRemote) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitHash) ||
				errors.Is(err, gitutil.ErrCouldNotDetectGitBranch) {
				// Keep going anyway. Either not a git dir, or git not installed, or
				// remote not detected.
				if errors.Is(err, gitutil.ErrNoGitBinary) {
					// TODO: Log this properly in the console.
					fmt.Printf("Warning: %s\n", err.Error())
				}
			} else {
				return nil, err
			}
		}
		// Note that this could be nil in some cases.
		lr.gitMetaCache[ref.GetLocalPath()] = metadata
	}

	buildFilePath, err := detectBuildFile(ref, filepath.FromSlash(ref.GetLocalPath()))
	if err != nil {
		return nil, err
	}

	var buildContext llb.State
	if _, isTarget := ref.(domain.Target); isTarget {
		excludes, err := readExcludes(ref.GetLocalPath())
		if err != nil {
			return nil, err
		}
		buildContext = llb.Local(
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
		BuildFilePath: buildFilePath,
		BuildContext:  buildContext,
		GitMetadata:   metadata,
	}, nil
}
