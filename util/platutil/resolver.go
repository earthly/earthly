package platutil

import (
	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/util/llbutil/pllb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Resolver is a platform resolver.
type Resolver struct {
	AllowNativeAndUser bool // this will be set via feature flag.

	currentPlatform Platform
	defaultPlatform Platform
	userPlatform    specs.Platform
	nativePlatform  specs.Platform
}

// NewResolver returns a new platform resolver.
func NewResolver(nativePlatform specs.Platform) *Resolver {
	return &Resolver{
		currentPlatform: DefaultPlatform,
		defaultPlatform: DefaultPlatform,
		nativePlatform:  nativePlatform,
		userPlatform:    GetUserPlatform(),
	}
}

// SubResolver returns a copy of this resolver, but with the current and default
// platform overriden.
func (r *Resolver) SubResolver(newPlatform Platform) *Resolver {
	if newPlatform == DefaultPlatform {
		newPlatform = r.defaultPlatform
	}
	return &Resolver{
		currentPlatform:    newPlatform,
		defaultPlatform:    newPlatform,
		userPlatform:       r.userPlatform,
		nativePlatform:     r.nativePlatform,
		AllowNativeAndUser: r.AllowNativeAndUser,
	}
}

// SubPlatform returns a platform that defaults to the current platform of the
// resolver (as opposed to the default platform of the resolver).
func (r *Resolver) SubPlatform(in Platform) Platform {
	if in == DefaultPlatform {
		return r.currentPlatform
	}
	return in
}

// UpdatePlatform set the current platform of this resolver and returns
// the effective platform set, resolving the default accordingly.
func (r *Resolver) UpdatePlatform(newPlatform Platform) Platform {
	if newPlatform == DefaultPlatform {
		newPlatform = r.defaultPlatform
	}
	r.currentPlatform = newPlatform
	return newPlatform
}

// Current returns the current platform.
func (r *Resolver) Current() Platform {
	return r.currentPlatform
}

// Default returns the default platform.
func (r *Resolver) Default() Platform {
	return r.defaultPlatform
}

// LLBNative returns the native platform.
func (r *Resolver) LLBNative() specs.Platform {
	return r.nativePlatform
}

// LLBUser returns the user platform.
func (r *Resolver) LLBUser() specs.Platform {
	return r.userPlatform
}

// Parse parses a given platform string. Empty string is a valid selection:
// it means "the default platform".
func (r *Resolver) Parse(str string) (Platform, error) {
	if r.AllowNativeAndUser {
		return r.ParseAllowNativeAndUser(str)
	}
	if str == "" {
		return DefaultPlatform, nil
	}
	p, err := platforms.Parse(str)
	if err != nil {
		return Platform{}, err
	}
	p = platforms.Normalize(p)
	return Platform{p: &p}, nil
}

// ParseAllowNativeAndUser parses a given platform string. Empty string is a valid selection:
// it means "the default platform". This variant forces allowing "native" and "user".
func (r *Resolver) ParseAllowNativeAndUser(str string) (Platform, error) {
	switch str {
	case "native":
		return NativePlatform, nil
	case "user":
		return UserPlatform, nil
	case "":
		return DefaultPlatform, nil
	default:
		p, err := platforms.Parse(str)
		if err != nil {
			return Platform{}, err
		}
		p = platforms.Normalize(p)
		return Platform{p: &p}, nil
	}
}

// Materialize turns a platform into a concrete platform
// (resolves "user" / "native" / "") to an actual value.
func (r *Resolver) Materialize(in Platform) Platform {
	var out specs.Platform
	switch {
	case in.p != nil:
		out = *in.p
	case in.user:
		out = r.userPlatform
	default: // in.native or none (default)
		out = r.nativePlatform
	}
	out = platforms.Normalize(out)
	return Platform{p: &out}
}

// Scratch is the scratch state with the native platform readily set.
func (r *Resolver) Scratch() pllb.State {
	return pllb.Scratch().Platform(r.nativePlatform)
}

// PlatformEquals compares two platforms if the equate to the same platform.
// A "native" platform can still equate to a "user" platform, if they
// materialize to the same platform.
func (r *Resolver) PlatformEquals(p1, p2 Platform) bool {
	p1 = r.Materialize(p1)
	p2 = r.Materialize(p2)
	return p1.p.OS == p2.p.OS &&
		p1.p.Architecture == p2.p.Architecture &&
		p1.p.Variant == p2.p.Variant
}

// ToLLBPlatform returns the containerd platform.
func (r *Resolver) ToLLBPlatform(in Platform) specs.Platform {
	return *r.Materialize(in).p
}

// GetUserPlatform returns the user platform.
func GetUserPlatform() specs.Platform {
	return platforms.Normalize(platforms.DefaultSpec())
}
