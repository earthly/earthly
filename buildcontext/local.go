package buildcontext

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/moby/buildkit/client/llb"
	"github.com/vladaionescu/earthly/domain"
	"github.com/vladaionescu/earthly/llbutil"
	"github.com/vladaionescu/earthly/logging"
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
				errors.Is(err, ErrCouldNotDetectRemote) {
				// Keep going anyway. Either not a git dir, or git not installed, or
				// remote not detected.
				logging.GetLogger(ctx).Warning(err.Error())
				if errors.Is(err, ErrNoGitBinary) ||
					errors.Is(err, ErrCouldNotDetectRemote) {
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

	return &Data{
		EarthfilePath: filepath.Join(filepath.FromSlash(target.LocalPath), "build.earth"),
		BuildContext: llb.Local(
			target.LocalPath,
			llb.SharedKeyHint(target.LocalPath),
			llb.ExcludePatterns(excludes),
			llb.SessionID(lr.sessionID),
			llb.Platform(llbutil.TargetPlatform),
			llb.WithCustomNamef("[context] local context %s", target.LocalPath),
		),
		GitMetadata: metadata,
	}, nil
}
