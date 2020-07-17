package llbutil

import "github.com/containerd/containerd/platforms"

var (
	// TargetPlatform is the platform used for the build and for the resulting output images.
	// This will be configurable in the future.
	TargetPlatform = platforms.MustParse("linux/amd64")

	// DefaultPathEnv is the default PATH to use.
	DefaultPathEnv = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)
