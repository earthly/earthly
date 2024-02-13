package platutil

import (
	"github.com/containerd/containerd/platforms"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	// DefaultPlatform returns the default platform object.
	DefaultPlatform = Platform{}
	// NativePlatform returns the native platform.
	NativePlatform = Platform{native: true}
	// UserPlatform returns the user platform.
	UserPlatform = Platform{user: true}
)

// Platform is a platform set to either the user's platform, the native platform where the
// build executes, or another specific platform.
type Platform struct {
	user   bool
	native bool
	p      *specs.Platform
}

// IsPlatformDefined returns true when the platform was explicitly set
func IsPlatformDefined(p Platform) bool {
	return p != DefaultPlatform
}

// FromLLBPlatform returns a platform from a containerd platform.
func FromLLBPlatform(p specs.Platform) Platform {
	p = platforms.Normalize(p)
	return Platform{p: &p}
}

// String returns the string representation of the platform.
func (p Platform) String() string {
	if p.p != nil {
		return platforms.Format(*p.p)
	}
	if p.native {
		return "native"
	}
	if p.user {
		return "user"
	}
	return ""
}
