package states

import (
	"context"
	"sync"

	"github.com/moby/buildkit/client/llb"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/waitutil"
	"github.com/earthly/earthly/variables"
	"github.com/google/uuid"
)

// MultiTarget holds LLB states representing multiple earthly targets,
// in the order in which they should be built.
type MultiTarget struct {
	// Visited represents the previously visited states, grouped by target
	// name. Duplicate targets are possible if same target is called with different
	// build args.
	Visited VisitedCollection
	// Final is the main target to be built.
	Final *SingleTarget
}

// CacheMount holds run options needed to cache mounts, and some extra options.
type CacheMount struct {
	// Persisted should the cache be persisted to image.
	Persisted bool
	// RunOption Run options
	RunOption llb.RunOption
}

// FinalTarget returns the final target of the states.
func (mts *MultiTarget) FinalTarget() domain.Target {
	return mts.Final.Target
}

// All returns all SingleTarget contained within.
func (mts *MultiTarget) All() []*SingleTarget {
	return mts.Visited.All()
}

// SingleTarget holds LLB states representing an earthly target.
type SingleTarget struct {
	// ID is a random unique string.
	ID                     string
	Target                 domain.Target
	PlatformResolver       *platutil.Resolver
	MainImage              *image.Image
	MainState              pllb.State
	ArtifactsState         pllb.State
	SeparateArtifactsState []pllb.State
	SaveLocals             []SaveLocal
	SaveImages             []SaveImage
	VarCollection          *variables.Collection
	RunPush                RunPush
	InteractiveSession     InteractiveSession
	GlobalImports          map[string]domain.ImportTrackerVal
	// HasDangling represents whether the target has dangling instructions -
	// ie if there are any non-SAVE commands after the first SAVE command,
	// or if the target is invoked via BUILD command (not COPY nor FROM).
	HasDangling bool
	// RanFromLike represents whether we have encountered a FROM-like command
	// (eg FROM, FROM DOCKERFILE, LOCALLY).
	RanFromLike bool
	// RanInteractive represents whether we have encountered an --interactive command.
	RanInteractive bool

	Sboms []string

	// doSavesMu is a mutex for doSave.
	doSavesMu sync.Mutex
	// doSaves indicates whether the SaveImages and the SaveLocals should actually be saved
	doSaves bool

	// doPushes indicates whether the SaveImages should actually be pushed
	doPushes bool

	// WaitBlocks contains the caller's waitblock plus any additional waitblocks defined in the target
	WaitBlocks []waitutil.WaitBlock

	// WaitItems contains all wait items which are created by the target
	// it exists for tracking items in the target vs a caller's wait block that is shared between multiple targets
	WaitItems []waitutil.WaitItem

	// doneCh is a channel that is closed when the sts is complete.
	doneCh chan struct{}

	tiMu        sync.Mutex
	targetInput dedup.TargetInput

	depMu sync.Mutex
	// dependentIDs are the sts IDs of the transitive dependants of this target.
	dependentIDs map[string]bool
	// outgoingNewSubscriptions is a list of channels to update when new dependentIDs are added.
	outgoingNewSubscriptions []chan string
	incomingNewSubscriptions chan string
}

func newSingleTarget(ctx context.Context, target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, parentDepSub chan string) (*SingleTarget, error) {
	targetStr := target.StringCanonical()
	sts := &SingleTarget{
		ID:               uuid.New().String(),
		Target:           target,
		PlatformResolver: nil, // Will be set in converter's FinalizeStates.
		targetInput: dedup.TargetInput{
			TargetCanonical: targetStr,
			Platform:        platr.Materialize(platr.Current()).String(),
			AllowPrivileged: allowPrivileged,
		},
		MainState:                platr.Scratch(),
		MainImage:                image.NewImage(),
		ArtifactsState:           platr.Scratch(),
		dependentIDs:             make(map[string]bool),
		doneCh:                   make(chan struct{}),
		incomingNewSubscriptions: make(chan string, 1024),
	}
	sts.addOverridingVarsAsBuildArgInputs(overridingVars)
	if parentDepSub == nil {
		// New simplified algorithm.
		return sts, nil
	}
	// Consume all items from the parent subscription before returning control.
OuterLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case newID := <-parentDepSub:
			sts.AddDependentIDs(map[string]bool{newID: true})
		default:
			break OuterLoop
		}
	}
	// Keep monitoring async.
	sts.MonitorDependencySubscription(ctx, parentDepSub)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case newID := <-sts.incomingNewSubscriptions:
				sts.AddDependentIDs(map[string]bool{newID: true})
			}

		}
	}()
	return sts, nil
}

// GetDoSaves returns whether the SaveImages and the SaveLocals should be
// actually saved.
func (sts *SingleTarget) GetDoSaves() bool {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	return sts.doSaves
}

// SetDoSaves sets the DoSaves flag.
func (sts *SingleTarget) SetDoSaves() {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	sts.doSaves = true
	for _, wi := range sts.WaitItems {
		wi.SetDoSave()
	}
	for _, wb := range sts.WaitBlocks {
		wb.SetDoSaves()
	}
}

// GetDoPushes returns whether the SAVE IMAGE --push or RUN --push commands
// should be executed
func (sts *SingleTarget) GetDoPushes() bool {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	return sts.doPushes
}

// SetDoPushes sets the doPushes flag.
func (sts *SingleTarget) SetDoPushes() {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	sts.doPushes = true
	for _, wi := range sts.WaitItems {
		wi.SetDoPush()
	}
	for _, wb := range sts.WaitBlocks {
		wb.SetDoPushes()
	}
}

// AddWaitBlock adds a wait block to the state
func (sts *SingleTarget) AddWaitBlock(waitBlock waitutil.WaitBlock) {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	sts.WaitBlocks = append(sts.WaitBlocks, waitBlock)
}

// Wait performs a Wait on all wait blocks
func (sts *SingleTarget) Wait(ctx context.Context) error {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	for i := len(sts.WaitBlocks) - 1; i >= 0; i-- {
		err := sts.WaitBlocks[i].Wait(ctx, sts.doPushes, sts.doSaves)
		if err != nil {
			return err
		}
	}
	return nil
}

// AttachTopLevelWaitItems adds pre-created wait items to a new waitblock
func (sts *SingleTarget) AttachTopLevelWaitItems(ctx context.Context, waitBlock waitutil.WaitBlock) {
	sts.doSavesMu.Lock()
	defer sts.doSavesMu.Unlock()
	for _, item := range sts.WaitItems {
		waitBlock.AddItem(item)
	}
}

// TargetInput returns the target input in a concurrent-safe way.
func (sts *SingleTarget) TargetInput() dedup.TargetInput {
	sts.tiMu.Lock()
	defer sts.tiMu.Unlock()
	return sts.targetInput
}

// AddBuildArgInput adds a bai to the sts's target input.
func (sts *SingleTarget) AddBuildArgInput(bai dedup.BuildArgInput) {
	sts.tiMu.Lock()
	defer sts.tiMu.Unlock()
	sts.targetInput = sts.targetInput.WithBuildArgInput(bai)
}

// LastSaveImage returns the last save image available (if any).
func (sts *SingleTarget) LastSaveImage() SaveImage {
	if len(sts.SaveImages) == 0 {
		// Use main state / image if no save image exists.
		return SaveImage{
			State: sts.MainState,
			Image: sts.MainImage,
		}
	}
	return sts.SaveImages[len(sts.SaveImages)-1]
}

// AddDependentIDs adds additional IDs that depend on this sts.
func (sts *SingleTarget) AddDependentIDs(dependentIDs map[string]bool) {
	select {
	case <-sts.doneCh:
		// We don't really care anymore if we're done.
		return
	default:
	}
	sts.depMu.Lock()
	defer sts.depMu.Unlock()
	for ID := range dependentIDs {
		sts.dependentIDs[ID] = true
	}
	for _, sub := range sts.outgoingNewSubscriptions {
		for ID := range dependentIDs {
			sub <- ID
		}
	}
}

// MonitorDependencySubscription monitors for new dependencies.
func (sts *SingleTarget) MonitorDependencySubscription(ctx context.Context, inCh chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ID := <-inCh:
				sts.incomingNewSubscriptions <- ID
			}
		}
	}()
}

// NewDependencySubscription adds additional IDs that depend on this sts.
func (sts *SingleTarget) NewDependencySubscription() chan string {
	sts.depMu.Lock()
	defer sts.depMu.Unlock()
	ch := make(chan string, 1024) // size is an arbitrary maximum cycle length
	sts.outgoingNewSubscriptions = append(sts.outgoingNewSubscriptions, ch)
	// Send everything we have so far.
	ch <- sts.ID // send our ID
	for depID := range sts.dependentIDs {
		ch <- depID
	}
	return ch
}

// Done returns a channel that is closed when the sts is complete.
func (sts *SingleTarget) Done() chan struct{} {
	return sts.doneCh
}

func (sts *SingleTarget) addOverridingVarsAsBuildArgInputs(overridingVars *variables.Scope) {
	sts.tiMu.Lock()
	defer sts.tiMu.Unlock()
	for _, key := range overridingVars.Sorted() {
		ovVar, _ := overridingVars.Get(key)
		sts.targetInput = sts.targetInput.WithBuildArgInput(
			dedup.BuildArgInput{ConstantValue: ovVar, Name: key})
	}
}

// SaveLocal is an artifact path to be saved to local disk.
type SaveLocal struct {
	// DestPath is the local dest path to copy the artifact to.
	DestPath string
	// ArtifactPath is the relative path within the artifacts image.
	ArtifactPath string
	// Index is the index number of the "save as local" command encountered. Starts as 0.
	Index int
	// IfExists allows the artifact to be optional.
	IfExists bool
}

// SaveImage is a docker image to be saved.
type SaveImage struct {
	State        pllb.State
	Image        *image.Image
	DockerTag    string
	Push         bool
	InsecurePush bool
	// CacheHint instructs Earthly to save a separate ref for this image, even if no tag is
	// provided.
	CacheHint           bool
	HasPushDependencies bool
	// ForceSave indicates whether the image should be force-saved and (possibly pushed).
	ForceSave bool
	// CheckDuplicate indicates whether to check if the image name shows up
	// multiple times during output.
	CheckDuplicate bool
	// NoManifestList indicates that the image should not include a manifest
	// list (usually used for multi-platform setups). This means that the image
	// can only be a single-platform image.
	NoManifestList bool

	Platform    platutil.Platform
	HasPlatform bool // true when the --platform value was set (either on cli, or via FROM --platform=..., or BUILD --platform=...)

	SkipBuilder bool // for use with WAIT/END
}

// RunPush is a series of RUN --push commands to be run after the build has been deemed as
// successful, along with artifacts to save and images to push
type RunPush struct {
	CommandStrs        []string
	State              pllb.State
	SaveLocals         []SaveLocal
	SaveImages         []SaveImage
	InteractiveSession InteractiveSession
	HasState           bool
}

// InteractiveSessionKind represents what kind of interactive session has been encountered.
type InteractiveSessionKind string

const (
	// SessionKeep is a session where the data *persists* in the image when it exits.
	SessionKeep InteractiveSessionKind = "keep"
	// SessionEphemeral is a session where the data *does not persist* in the image when it exits.
	SessionEphemeral InteractiveSessionKind = "ephemeral"
)

// InteractiveSession holds the relevant data for running an interactive session when
// it is not desired to save the resulting changes into an image.
type InteractiveSession struct {
	CommandStr  string
	State       pllb.State
	Initialized bool
	Kind        InteractiveSessionKind
}
