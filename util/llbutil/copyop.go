package llbutil

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// CopyOp is a simplified llb copy operation.
func CopyOp(ctx context.Context, srcState pllb.State, srcs []string, destState pllb.State, dest string, allowWildcard, isDir, keepTs bool, chown string, chmod *fs.FileMode, ifExists, symlinkNoFollow, merge bool, opts ...llb.ConstraintsOpt) (pllb.State, error) {
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
		if ifExists && len(src) != 0 {
			// Strip ./ and / prefixes as to make paths relative to top-level.
			src = strings.TrimPrefix(src, "./")
			src = strings.TrimPrefix(src, "/")
			// HACK: For COPY --if-exists, we can use a glob expression (e.g., '[f]oo') to
			// prevent errors caused by non-existing files. This approach also works
			// with additional wildcards (e.g., '[f]oo/*').
			src = fmt.Sprintf("[%s]%s", string(src[0]), string(src[1:]))
		}
		copyOpts := append([]llb.CopyOption{
			&llb.CopyInfo{
				Mode:                chmod,
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
		return destState, nil
	}
	if merge && chown == "" {
		cwd, err := destState.GetDir(ctx)
		if err != nil {
			return pllb.State{}, err
		}
		return pllb.Merge([]pllb.State{destState, pllb.Scratch().Dir(cwd).File(fa)}, opts...).Dir(cwd), nil
	}
	return destState.File(fa, opts...), nil
}

// CopyWithRunOptions copies from `src` to `dest` and returns the result in a separate LLB State.
// This operation is similar llb.Copy, however, it can apply llb.RunOptions (such as a mount)
// Internally, the operation runs on the internal COPY image used by Dockerfile.
func CopyWithRunOptions(srcState pllb.State, src, dest string, platr *platutil.Resolver, opts ...llb.RunOption) pllb.State {
	// Docker's internal image for running COPY.
	// Ref: https://github.com/moby/buildkit/blob/v0.9.3/frontend/dockerfile/dockerfile2llb/convert.go#L40
	const copyImg = "docker/dockerfile-copy:v0.1.9@sha256:e8f159d3f00786604b93c675ee2783f8dc194bb565e61ca5788f6a6e9d304061"
	// Use the native platform instead of the target platform.
	imgOpts := []llb.ImageOption{llb.MarkImageInternal, llb.Platform(platr.LLBNative())}

	// The following executes the `copy` command, which is a custom executable
	// contained in the Dockerfile COPY image above. The following .Run()
	// operation executes in a state constructed from that Dockerfile COPY image,
	// with the Earthly user's state mounted at /dest on that image.
	opts = append(opts, []llb.RunOption{
		llb.ReadonlyRootFS(),
		llb.Shlexf("copy %s /dest/%s", src, dest)}...)
	copyState := pllb.Image(copyImg, imgOpts...)
	run := copyState.Run(opts...)
	destState := run.AddMount("/dest", srcState)
	destState = destState.Platform(platr.ToLLBPlatform(platr.Current()))
	return destState
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
