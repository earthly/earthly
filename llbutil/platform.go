package llbutil

import (
	"runtime"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// DefaultPlatform returns the default platform to use if none is specified.
func DefaultPlatform() specs.Platform {
	p := platforms.DefaultSpec()
	if runtime.GOOS == "darwin" {
		// Use linux so that this works with Docker Desktop app.
		p.OS = "linux"
	}
	return platforms.Normalize(p)
}

// ScratchWithPlatform is the scratch state with the default platform readily set.
func ScratchWithPlatform() llb.State {
	return llb.Scratch().Platform(DefaultPlatform())
}

// PlatformToString turns a platform pointer into a string representation.
func PlatformToString(p *specs.Platform) string {
	if p == nil {
		return ""
	}
	return platforms.Format(*p)
}
