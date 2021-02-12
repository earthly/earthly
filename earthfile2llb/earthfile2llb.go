package earthfile2llb

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/variables"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// ConvertOpt holds conversion parameters needed for conversion.
type ConvertOpt struct {
	// GwClient is the BuildKit gateway client.
	GwClient gwclient.Client
	// Resolver is the build context resolver.
	Resolver *buildcontext.Resolver
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
	// VarCollection is a collection of build args used for overriding args in the build.
	VarCollection *variables.Collection
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
	// EnableAst is an internal feature flag for the new Earthly AST.
	EnableAst bool
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
		same := (sts.TargetInput.Platform == llbutil.PlatformToString(opt.Platform))
		if same {
			for _, bai := range sts.TargetInput.BuildArgs {
				variable, _, found := opt.VarCollection.Get(bai.Name)
				if found {
					if !variable.BuildArgInput(bai.Name, bai.DefaultValue).Equals(bai) {
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
				return nil, fmt.Errorf(
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
	converter, err := NewConverter(ctx, bc.Target, bc, opt)
	if err != nil {
		return nil, err
	}
	// TODO: Use a parser cache.
	ef, err := ast.Parse(ctx, bc.BuildFilePath, true)
	if err != nil {
		return nil, err
	}
	interpreter := newInterpreter(converter)
	err = interpreter.Run(ctx, ef, target.Target)
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
