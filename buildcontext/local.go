package buildcontext

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/moby/buildkit/client/llb"
)

type localResolver struct {
	gitMetaCache map[string]*GitMetadata
	sessionID    string
}

func (lr *localResolver) resolveLocal(ctx context.Context, target domain.Target) (*Data, error) {
	if target.IsRemote() {
		return nil, fmt.Errorf("Unexpected remote target %s", target.String())
	}
	excludes, err := readExcludes(target.LocalPath)
	if err != nil {
		return nil, err
	}

	metadata, found := lr.gitMetaCache[target.LocalPath]
	if !found {
		metadata, err = Metadata(ctx, target.LocalPath)
		if err != nil {
			if errors.Is(err, ErrNoGitBinary) ||
				errors.Is(err, ErrNotAGitDir) ||
				errors.Is(err, ErrCouldNotDetectRemote) ||
				errors.Is(err, ErrCouldNotDetectGitHash) ||
				errors.Is(err, ErrCouldNotDetectGitBranch) {
				// Keep going anyway. Either not a git dir, or git not installed, or
				// remote not detected.
				if errors.Is(err, ErrNoGitBinary) {
					// TODO: Log this properly in the console.
					fmt.Printf("Warning: %s\n", err.Error())
				}
			} else {
				return nil, err
			}
		}
		// Note that this could be nil in some cases.
		lr.gitMetaCache[target.LocalPath] = metadata
	}

	buildFilePath, err := detectBuildFile(target, filepath.FromSlash(target.LocalPath))
	if err != nil {
		return nil, err
	}
	return &Data{
		BuildFilePath: buildFilePath,
		BuildContext: llb.Local(
			target.LocalPath,
			llb.SharedKeyHint(target.LocalPath),
			llb.ExcludePatterns(excludes),
			llb.SessionID(lr.sessionID),
			llb.Platform(llbutil.TargetPlatform),
			llb.WithCustomNamef("[context %s] local context %s", target.LocalPath, target.LocalPath),
		),
		GitMetadata: metadata,
	}, nil
}
