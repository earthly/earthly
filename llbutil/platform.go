package llbutil

import (
	"runtime"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
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

// ResolvePlatform returns the non-nil platform provided. If both are nil, nil is returned.
// If both are non-nil, they are compared to ensure they are identical. If they are not,
// an error is returned.
func ResolvePlatform(p1 *specs.Platform, p2 *specs.Platform) (*specs.Platform, error) {
	if p1 == nil {
		return p2, nil
	}
	if p2 == nil {
		return p1, nil
	}
	plat1 := platforms.Normalize(*p1)
	plat2 := platforms.Normalize(*p2)
	if plat1.OS != plat2.OS {
		return nil, errors.Errorf(
			"platform contradiction %s vs %s",
			platforms.Format(*p1), platforms.Format(*p2))
	}
	if plat1.Architecture != plat2.Architecture {
		return nil, errors.Errorf(
			"platform contradiction %s vs %s",
			platforms.Format(*p1), platforms.Format(*p2))
	}
	if plat1.Variant != plat2.Variant {
		return nil, errors.Errorf(
			"platform contradiction %s vs %s",
			platforms.Format(*p1), platforms.Format(*p2))
	}
	return p1, nil
}

// PlatformWithDefault returns the same platform provided if not nil, or the default
// platform otherwise.
func PlatformWithDefault(p *specs.Platform) specs.Platform {
	if p != nil {
		return *p
	}
	return DefaultPlatform()
}
