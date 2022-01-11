package llbutil

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// CopyOp is a simplified llb copy operation.
func CopyOp(srcState pllb.State, srcs []string, destState pllb.State, dest string, allowWildcard bool, isDir bool, keepTs bool, chown string, ifExists, symlinkNoFollow bool, opts ...llb.ConstraintsOpt) pllb.State {
	destAdjusted := dest
	if dest == "." || dest == "" || len(srcs) > 1 {
		destAdjusted += string("/") // TODO: needs to be the containers platform, not the earthly hosts platform. For now, this is always Linux.
	}
	var baseCopyOpts []llb.CopyOption
	if chown != "" {
		baseCopyOpts = append(baseCopyOpts, llb.WithUser(chown))
	}
	var fa *pllb.FileAction
	if !keepTs {
		baseCopyOpts = append(baseCopyOpts, llb.WithCreatedTime(*defaultTs()))
	}
	for _, src := range srcs {
		if ifExists {
			// If the copy came in as optional (ifExists), then we need to trigger the
			// underlying wildcard matching and allow empty wildcards. The matching uses
			// the filepath.Match syntax, so by simply creating a wildcard where the
			// first letter needs to match the current first letter gets us the single
			// match; and no error if it is missing.

			//Normalize path by dropping './'
			src = strings.TrimPrefix(src, "./")
			src = fmt.Sprintf("[%s]%s", string(src[0]), string(src[1:]))
		}
		copyOpts := append([]llb.CopyOption{
			&llb.CopyInfo{
				FollowSymlinks:      !symlinkNoFollow,
				CopyDirContentsOnly: !isDir,
				AttemptUnpack:       false,
				CreateDestPath:      true,
				AllowWildcard:       allowWildcard,
				AllowEmptyWildcard:  ifExists,
			},
		}, baseCopyOpts...)
		if fa == nil {
			fa = pllb.Copy(srcState, src, destAdjusted, copyOpts...)
		} else {
			fa = fa.Copy(srcState, src, destAdjusted, copyOpts...)
		}
	}
	if fa == nil {
		return destState
	}
	return destState.File(fa, opts...)
}

// Abs prepends the working dir to the given path, if the
// path is relative.
func Abs(ctx context.Context, s pllb.State, p string) (string, error) {
	if path.IsAbs(p) {
		return p, nil
	}
	dir, err := s.GetDir(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get dir")
	}
	return path.Join(dir, p), nil
}

var defaultTsValue time.Time
var defaultTsParse sync.Once

func defaultTs() *time.Time {
	defaultTsParse.Do(func() {
		var err error
		defaultTsValue, err = time.Parse(time.RFC3339, "2020-04-16T12:00:00+00:00")
		if err != nil {
			panic(err)
		}
	})
	return &defaultTsValue
}
