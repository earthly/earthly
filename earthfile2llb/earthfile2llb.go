package earthfile2llb

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/variables"
)

// ConvertOpt holds conversion parameters.
type ConvertOpt struct {
	// GwClient is the BuildKit gateway client.
	GwClient gwclient.Client
	// Resolver is the build context resolver.
	Resolver *buildcontext.Resolver
	// GlobalImports is a map of imports used to dereference import ref targets, commands, etc.
	GlobalImports map[string]domain.ImportTrackerVal
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
	CacheImports *states.CacheImports
	// UseInlineCache enables the inline caching feature (use any SAVE IMAGE --push declaration as
	// cache import).
	UseInlineCache bool
	// UseFakeDep is an internal feature flag for fake dep.
	UseFakeDep bool
	// AllowLocally is an internal feature flag for controlling if LOCALLY directives can be used.
	AllowLocally bool
	// AllowInteractive is an internal feature flag for controlling if interactive sessions can be initiated.
	AllowInteractive bool
	// HasDangling represents whether the target has dangling instructions -
	// ie if there are any non-SAVE commands after the first SAVE command,
	// or if the target is invoked via BUILD command (not COPY nor FROM).
	HasDangling bool
	// Console is for logging
	Console conslogging.ConsoleLogger
	// AllowPrivileged is used to allow (or prevent) any "RUN --privileged" or RUNs under a LOCALLY target to be executed,
	// when set to false, it prevents other referenced remote targets from requesting elevated privileges
	AllowPrivileged bool
	// DoSaves is used to control when SAVE ARTIFACT AS LOCAL calls will actually output the artifacts locally
	// this is to differentiate between calling a target that saves an artifact directly vs using a FROM which indirectly
	// calls a target which saves an artifact as a side effect.
	DoSaves bool
	// ForceSaveImage is used to force all SAVE IMAGE commands are executed regardless of if they are
	// for a local or remote target; this is to support the legacy behaviour that was first introduced in earthly (up to 0.5)
	// When this is set to false, SAVE IMAGE commands are only executed when DoSaves is true.
	ForceSaveImage bool
	// Gitlookup is used to attach credentials to GIT CLONE operations
	GitLookup *buildcontext.GitLookup
	// LocalStateCache provides a cache for local pllb.States
	LocalStateCache *LocalStateCache

	// Features is the set of enabled features
	Features *features.Features

	// ParallelConversion is a feature flag enabling the parallel conversion algorithm.
	ParallelConversion bool
	// Parallelism is a semaphore controlling the maximum parallelism.
	Parallelism *semaphore.Weighted

	// parentDepSub is a channel informing of any new dependencies from the parent.
	parentDepSub chan string // chan of sts IDs.

	// FeatureFlagOverride is used to override feature flags that are defined in specific Earthfiles
	FeatureFlagOverrides string
}

// Earthfile2LLB parses a earthfile and executes the statements for a given target.
func Earthfile2LLB(ctx context.Context, target domain.Target, opt ConvertOpt, initialCall bool) (mts *states.MultiTarget, err error) {
	if opt.SolveCache == nil {
		opt.SolveCache = states.NewSolveCache()
	}
	if opt.Visited == nil {
		opt.Visited = states.NewVisitedCollection()
	}
	if opt.MetaResolver == nil {
		opt.MetaResolver = NewCachedMetaResolver(opt.GwClient)
	}
	// Resolve build context.
	bc, err := opt.Resolver.Resolve(ctx, opt.GwClient, target)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve build context for target %s", target.String())
	}

	ftrs, err := features.GetFeatures(bc.Earthfile.Version)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve feature set for version %v for target %s", bc.Earthfile.Version.Args, target.String())
	}
	err = features.ApplyFlagOverrides(ftrs, opt.FeatureFlagOverrides)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to apply version feature overrides")
	}
	opt.Features = ftrs
	if initialCall {
		// It's not possible to know if we should DoSaves until after we have parsed the target's VERSION features.
		if ftrs.ReferencedSaveOnly {
			opt.DoSaves = true
		} else {
			if !target.IsRemote() {
				opt.DoSaves = true // legacy mode only saves artifacts that are locally referenced
			}
			opt.ForceSaveImage = true // legacy mode always saves images regardless of locally or remotely referenced
		}
	}

	targetWithMetadata := bc.Ref.(domain.Target)
	sts, found, err := opt.Visited.Add(ctx, targetWithMetadata, opt.Platform, opt.AllowPrivileged, opt.OverridingVars, opt.parentDepSub)
	if err != nil {
		return nil, err
	}
	if found {
		// This target has already been done.
		return &states.MultiTarget{
			Final:   sts,
			Visited: opt.Visited,
		}, nil
	}
	converter, err := NewConverter(ctx, targetWithMetadata, bc, sts, opt, ftrs)
	if err != nil {
		return nil, err
	}
	interpreter := newInterpreter(converter, targetWithMetadata, opt.AllowPrivileged, opt.ParallelConversion, opt.Parallelism, opt.Console, opt.GitLookup)
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
