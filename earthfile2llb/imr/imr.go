package imr

import (
	"context"
	"os"
	"sync"

	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/docker/cli/cli/config"
	"github.com/docker/docker/pkg/locker"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/contentutil"
	"github.com/moby/buildkit/util/imageutil"
	digest "github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

var defaultImageMetaResolver llb.ImageMetaResolver
var defaultImageMetaResolverOnce sync.Once

var WithDefault = imageOptionFunc(func(ii *llb.ImageInfo) {
	llb.WithMetaResolver(Default()).SetImageOption(ii)
})

type imageMetaResolverOpts struct {
	platform *specs.Platform
}

type ImageMetaResolverOpt func(o *imageMetaResolverOpts)

func WithDefaultPlatform(p *specs.Platform) ImageMetaResolverOpt {
	return func(o *imageMetaResolverOpts) {
		o.platform = p
	}
}

func New(ctx context.Context, with ...ImageMetaResolverOpt) llb.ImageMetaResolver {
	r := docker.NewResolver(docker.ResolverOptions{
		Authorizer: docker.NewDockerAuthorizer(
			docker.WithAuthCreds(makeCredentialsFun()),
		),
	})
	var opts imageMetaResolverOpts
	for _, f := range with {
		f(&opts)
	}
	return &imageMetaResolver{
		resolver: r,
		platform: opts.platform,
		buffer:   contentutil.NewBuffer(),
		cache:    map[string]resolveResult{},
		locker:   locker.New(),
	}
}

func Default() llb.ImageMetaResolver {
	defaultImageMetaResolverOnce.Do(func() {
		defaultImageMetaResolver = New(context.Background())
	})
	return defaultImageMetaResolver
}

type imageMetaResolver struct {
	resolver remotes.Resolver
	buffer   contentutil.Buffer
	platform *specs.Platform
	locker   *locker.Locker
	cache    map[string]resolveResult
}

type resolveResult struct {
	config []byte
	dgst   digest.Digest
}

func (imr *imageMetaResolver) ResolveImageConfig(ctx context.Context, ref string, opt llb.ResolveImageConfigOpt) (digest.Digest, []byte, error) {
	imr.locker.Lock(ref)
	defer imr.locker.Unlock(ref)

	platform := opt.Platform
	if platform == nil {
		platform = imr.platform
	}

	k := imr.key(ref, platform)

	if res, ok := imr.cache[k]; ok {
		return res.dgst, res.config, nil
	}

	dgst, config, err := imageutil.Config(ctx, ref, imr.resolver, imr.buffer, nil, platform)
	if err != nil {
		return "", nil, err
	}

	imr.cache[k] = resolveResult{dgst: dgst, config: config}
	return dgst, config, nil
}

func (imr *imageMetaResolver) key(ref string, platform *specs.Platform) string {
	if platform != nil {
		ref += platforms.Format(*platform)
	}
	return ref
}

type imageOptionFunc func(*llb.ImageInfo)

func (fn imageOptionFunc) SetImageOption(ii *llb.ImageInfo) {
	fn(ii)
}

func makeCredentialsFun() func(host string) (string, string, error) {
	var mu sync.Mutex
	conf := config.LoadDefaultConfigFile(os.Stderr)
	return func(host string) (string, string, error) {
		mu.Lock()
		defer mu.Unlock()
		if host == "registry-1.docker.io" {
			host = "https://index.docker.io/v1/"
		}
		ac, err := conf.GetAuthConfig(host)
		if err != nil {
			return "", "", err
		}
		var secret string
		if ac.IdentityToken != "" {
			secret = ac.IdentityToken
		} else {
			secret = ac.Password
		}
		return ac.Username, secret, nil
	}
}
