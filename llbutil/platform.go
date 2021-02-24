package llbutil

import (
	"runtime"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// ParsePlatform parses a given platform string. Empty string is a valid selection:
// it means "the default platform".
func ParsePlatform(str string) (*specs.Platform, error) {
	if str == "" {
		return nil, nil
	}
	p, err := platforms.Parse(str)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

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

// PlatformEquals compares whether two platform pointers equate to the same platform.
// If any of the pointers is nil, then the default platform is assumed for it.
func PlatformEquals(p1 *specs.Platform, p2 *specs.Platform) bool {
	if p1 == p2 {
		// Quick way out.
		return true
	}
	pp1 := PlatformWithDefault(p1)
	pp2 := PlatformWithDefault(p2)
	return pp1.OS == pp2.OS &&
		pp1.Architecture == pp2.Architecture &&
		pp1.Variant == pp2.Variant
}

// PlatformWithDefaultToString turns a platform pointer into a string representation.
func PlatformWithDefaultToString(p *specs.Platform) string {
	return platforms.Format(PlatformWithDefault(p))
}

// PlatformToString turns a platform pointer into a string representation.
func PlatformToString(p *specs.Platform) string {
	if p == nil {
		return ""
	}
	return platforms.Format(*p)
}

// ResolvePlatform returns the non-nil platform provided. If both are nil, nil is returned.
// If both are non-nil, override is returned.
func ResolvePlatform(base *specs.Platform, override *specs.Platform) *specs.Platform {
	if override == nil {
		return base
	}
	return override
}

// PlatformWithDefault returns the same platform provided if not nil, or the default
// platform otherwise.
func PlatformWithDefault(p *specs.Platform) specs.Platform {
	if p != nil {
		return *p
	}
	return DefaultPlatform()
}
