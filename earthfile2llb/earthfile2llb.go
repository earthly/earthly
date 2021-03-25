package earthfile2llb

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/variables"
)

// ConvertOpt holds conversion parameters.
type ConvertOpt struct {
	// GwClient is the BuildKit gateway client.
	GwClient gwclient.Client
	// Resolver is the build context resolver.
	Resolver *buildcontext.Resolver
	// GlobalImports is a map of imports used to dereference import ref targets, commands, etc.
	GlobalImports map[string]string
	// The resolve mode for referenced images (force pull or prefer local).
	ImageResolveMode llb.ResolveMode
	// DockerBuilderFun is a fun that can be used to execute an image build. This
	// is used as part of operations like DOCKER LOAD and DOCKER PULL, where
	// a tar image is needed in the middle of a build.
	DockerBuilderFun states.DockerBuilderFun
	// CleanCollection is a collection of cleanup functions.
	CleanCollection *cleanup.Collection
	// Visited is a collection of target states which have been converted to LLB.
	// This is used for deduplication and infinite cycle detection.
	Visited *states.VisitedCollection
	// Platform is the target platform of the build.
	Platform *specs.Platform
	// OverridingVars is a collection of build args used for overriding args in the build.
	OverridingVars *variables.Scope
	// A cache for image solves. (maybe dockerTag +) depTargetInputHash -> context containing image.tar.
	SolveCache *states.SolveCache
	// BuildContextProvider is the provider used for local build context files.
	BuildContextProvider *provider.BuildContextProvider
	// MetaResolver is the image meta resolver to use for resolving image metadata.
	MetaResolver llb.ImageMetaResolver
	// CacheImports is a set of docker tags that can be used to import cache. Note that this
	// set is modified by the converter if InlineCache is enabled.
	CacheImports map[string]bool
	// UseInlineCache enables the inline caching feature (use any SAVE IMAGE --push declaration as
	// cache import).
	UseInlineCache bool
	// UseFakeDep is an internal feature flag for fake dep.
	UseFakeDep bool
	// AllowLocally is an internal feature flag for controlling if LOCALLY directives can be used.
	AllowLocally bool
	// AllowInteractive is an internal feature flag for controlling if interactive sessions can be initiated.
	AllowInteractive bool
}

// Earthfile2LLB parses a earthfile and executes the statements for a given target.
func Earthfile2LLB(ctx context.Context, target domain.Target, opt ConvertOpt) (mts *states.MultiTarget, err error) {
	if opt.SolveCache == nil {
		opt.SolveCache = states.NewSolveCache()
	}
	if opt.Visited == nil {
		opt.Visited = states.NewVisitedCollection()
	}
	if opt.MetaResolver == nil {
		opt.MetaResolver = opt.GwClient
	}
	// Check if we have previously converted this target, with the same build args.
	targetStr := target.String()
	for _, sts := range opt.Visited.Visited[targetStr] {
		stsPlat, err := llbutil.ParsePlatform(sts.TargetInput.Platform)
		if err != nil {
			return nil, err
		}
		same := llbutil.PlatformEquals(stsPlat, opt.Platform)
		if same {
			for _, bai := range sts.TargetInput.BuildArgs {
				variable, found := opt.OverridingVars.GetAny(bai.Name)
				if found {
					baiVariable := dedup.BuildArgInput{
						Name:          bai.Name,
						DefaultValue:  bai.DefaultValue,
						ConstantValue: variable,
					}
					if !baiVariable.Equals(bai) {
						same = false
						break
					}
				} else {
					if !bai.IsDefaultValue() {
						same = false
						break
					}
				}
			}
		}
		if same {
			if sts.Ongoing {
				return nil, errors.Errorf(
					"infinite recursion detected for target %s", targetStr)
			}
			// Use the already built states.
			return &states.MultiTarget{
				Final:   sts,
				Visited: opt.Visited,
			}, nil
		}
	}
	// Resolve build context.
	bc, err := opt.Resolver.Resolve(ctx, opt.GwClient, target)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve build context for target %s", target.String())
	}
	converter, err := NewConverter(ctx, bc.Ref.(domain.Target), bc, opt)
	if err != nil {
		return nil, err
	}
	interpreter := newInterpreter(converter, target)
	err = interpreter.Run(ctx, bc.Earthfile)
	if err != nil {
		return nil, err
	}
	return converter.FinalizeStates(ctx)
}

// GetTargets returns a list of targets from an Earthfile.
func GetTargets(filename string) ([]string, error) {
	ef, err := ast.Parse(context.TODO(), filename, false)
	if err != nil {
		return nil, err
	}
	targets := make([]string, 0, len(ef.Targets))
	for _, target := range ef.Targets {
		targets = append(targets, target.Name)
	}
	return targets, nil
}
