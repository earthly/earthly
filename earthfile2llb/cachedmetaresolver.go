package earthfile2llb

import (
	"context"

	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/syncutil/synccache"
	"github.com/moby/buildkit/client/llb"
	"github.com/opencontainers/go-digest"
)

var _ llb.ImageMetaResolver = &CachedMetaResolver{}

type cachedMetaResolverKey struct {
	ref      string
	platform string
}

type cachedMetaResolverEntry struct {
	dgst   digest.Digest
	config []byte
}

// CachedMetaResolver is an image meta resolver with a local cache.
type CachedMetaResolver struct {
	metaResolver llb.ImageMetaResolver
	cache        *synccache.SyncCache // cachedMetaResolverKey -> cachedMetaResolverEntry
}

// NewCachedMetaResolver creates a new cached meta resolver based on an underlying meta resolver
// which needs to be provided.
func NewCachedMetaResolver(metaResolver llb.ImageMetaResolver) *CachedMetaResolver {
	return &CachedMetaResolver{
		metaResolver: metaResolver,
		cache:        synccache.New(),
	}
}

// ResolveImageConfig implements llb.ImageMetaResolver.ResolveImageConfig.
func (cmr *CachedMetaResolver) ResolveImageConfig(ctx context.Context, ref string, opt llb.ResolveImageConfigOpt) (digest.Digest, []byte, error) {
	key := cachedMetaResolverKey{
		ref:      ref,
		platform: llbutil.PlatformToString(opt.Platform),
	}
	value, err := cmr.cache.Do(ctx, key, func(ctx context.Context, _ interface{}) (interface{}, error) {
		dgst, config, err := cmr.metaResolver.ResolveImageConfig(ctx, ref, opt)
		if err != nil {
			return nil, err
		}
		return cachedMetaResolverEntry{
			dgst:   dgst,
			config: config,
		}, nil
	})
	if err != nil {
		return "", nil, err
	}
	entry := value.(cachedMetaResolverEntry)
	return entry.dgst, entry.config, nil
}
