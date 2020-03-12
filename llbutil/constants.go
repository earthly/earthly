package llbutil

import "github.com/containerd/containerd/platforms"

var (
	// TargetPlatform is the platform used for the build and for the resulting output images.
	// This will be configurable in the future.
	TargetPlatform = platforms.MustParse("linux/amd64")
)
