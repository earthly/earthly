package llbutil

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// CopyOp is a simplified llb copy operation.
func CopyOp(srcState llb.State, srcs []string, destState llb.State, dest string, allowWildcard bool, isDir bool, chown string, opts ...llb.ConstraintsOpt) llb.State {
	destAdjusted := dest
	if dest == "." || dest == "" || strings.HasSuffix(dest, string(filepath.Separator)) {
		destAdjusted += string(filepath.Separator)
	}
	var baseCopyOpts []llb.CopyOption
	if chown != "" {
		baseCopyOpts = append(baseCopyOpts, llb.WithUser(chown))
	}
	var fa *llb.FileAction
	for _, src := range srcs {
		copyOpts := append([]llb.CopyOption{
			&llb.CopyInfo{
				FollowSymlinks:      true,
				CopyDirContentsOnly: !isDir,
				AttemptUnpack:       false,
				CreateDestPath:      true,
				AllowWildcard:       allowWildcard,
				AllowEmptyWildcard:  false,
			},
		}, baseCopyOpts...)
		if fa == nil {
			fa = llb.Copy(srcState, src, destAdjusted, copyOpts...)
		} else {
			fa = fa.Copy(srcState, src, destAdjusted, copyOpts...)
		}
	}
	if fa == nil {
		return destState
	}
	return destState.File(fa, opts...)
}

// Abs pre-pends the working dir to the given path, if the
// path is relative.
func Abs(ctx context.Context, s llb.State, p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	dir, err := s.GetDir(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get dir")
	}
	return filepath.Join(dir, p), nil
}
