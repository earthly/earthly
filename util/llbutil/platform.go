package llbutil

import (
	"context"
	"runtime"

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/moby/buildkit/client"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
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

// FromLLBPlatform returns a platform from a containerd platform.
func FromLLBPlatform(p specs.Platform) Platform {
	return Platform{p: &p}
}

// ParsePlatform parses a given platform string. Empty string is a valid selection:
// it means "the default platform".
func ParsePlatform(str string, allowNativeAndUser bool) (Platform, error) {
	switch str {
	case "native":
		if !allowNativeAndUser {
			return Platform{}, errors.New("platform \"native\" is not allowed in this version")
		}
		return NativePlatform, nil
	case "user":
		if !allowNativeAndUser {
			return Platform{}, errors.New("platform \"user\" is not allowed in this version")
		}
		return UserPlatform, nil
	case "":
		return DefaultPlatform, nil
	default:
		p, err := platforms.Parse(str)
		if err != nil {
			return Platform{}, err
		}
		return Platform{p: &p}, nil
	}
}

// ToLLBPlatform returns the containerd platform.
func (p Platform) ToLLBPlatform(nativePlatform specs.Platform) specs.Platform {
	ret := p.Resolve(nativePlatform).p
	return *ret
}

// Resolve retnurns the specific platform to use for the build. It resolves
// platforms such as "", "native" or "user" to an actual value.
func (p Platform) Resolve(nativePlatform specs.Platform) Platform {
	if p.p != nil {
		return p
	}
	if p.user {
		dp := userPlatform()
		return Platform{p: &dp}
	}
	// p.native or none set (default)
	return Platform{p: &nativePlatform}
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

func userPlatform() specs.Platform {
	p := platforms.DefaultSpec()
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		// Use linux so that this works with Docker Desktop app.
		p.OS = "linux"
	}
	return platforms.Normalize(p)
}

// ScratchWithPlatform is the scratch state with the default platform readily set.
func ScratchWithPlatform(nativePlatform specs.Platform) pllb.State {
	return pllb.Scratch().Platform(nativePlatform)
}

// PlatformEquals compares whether two platforms equate to the same platform.
// A "native" platform can still equate to a "user" platform, if they resolve
// to the same platform.
func PlatformEquals(p1 Platform, p2 Platform, nativePlatform specs.Platform) bool {
	if p1.native && p2.native {
		return true
	}
	if p1.user && p2.user {
		return true
	}
	if p1.p == p2.p {
		return true
	}
	p1 = p1.Resolve(nativePlatform)
	p2 = p2.Resolve(nativePlatform)
	if p1.p == p2.p {
		return true
	}
	return p1.p.OS == p2.p.OS &&
		p1.p.Architecture == p2.p.Architecture &&
		p1.p.Variant == p2.p.Variant
}

// GetNativePlatform returns the native platform for a given gwClient.
func GetNativePlatform(gwClient gwclient.Client) (specs.Platform, error) {
	ws := gwClient.BuildOpts().Workers
	if len(ws) == 0 {
		return specs.Platform{}, errors.New("no worker found via gwclient")
	}
	nps := ws[0].Platforms
	if len(nps) == 0 {
		return specs.Platform{}, errors.New("no platform found for worker via gwclient")
	}
	return platforms.Normalize(nps[0]), nil
}

// GetNativePlatformViaBkClient returns the native platform for a given buildkit client.
func GetNativePlatformViaBkClient(ctx context.Context, bkClient *client.Client) (specs.Platform, error) {
	ws, err := bkClient.ListWorkers(ctx)
	if err != nil {
		return specs.Platform{}, errors.Wrap(err, "failed to list workers")
	}
	if len(ws) == 0 {
		return specs.Platform{}, errors.New("no worker found via bkClient")
	}
	nps := ws[0].Platforms
	if len(nps) == 0 {
		return specs.Platform{}, errors.New("no platform found for worker via bkClient")
	}
	return platforms.Normalize(nps[0]), nil
}
