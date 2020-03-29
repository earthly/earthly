package llbutil

import (
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

// CopyOp is a simplified llb copy operation.
func CopyOp(srcState llb.State, srcs []string, destState llb.State, dest string, allowWildcard bool, isDir bool, opts ...llb.ConstraintsOpt) llb.State {
	destAdjusted := dest
	if dest == "." || dest == "" || strings.HasSuffix(dest, string(filepath.Separator)) {
		destAdjusted += string(filepath.Separator)
	}
	var fa *llb.FileAction
	for _, src := range srcs {
		copyOpts := []llb.CopyOption{
			&llb.CopyInfo{
				FollowSymlinks:      true,
				CopyDirContentsOnly: !isDir,
				AttemptUnpack:       false,
				CreateDestPath:      true,
				AllowWildcard:       allowWildcard,
				AllowEmptyWildcard:  false,
			},
		}
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
func Abs(s llb.State, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(s.GetDir(), p)
}
