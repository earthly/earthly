package states

import (
	"errors"

	"github.com/earthly/earthly/llbutil"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/variables"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// MultiTarget holds LLB states representing multiple earthly targets,
// in the order in which they should be built.
type MultiTarget struct {
	// Visited represents the previously visited states, grouped by target
	// name. Duplicate targets are possible if same target is called with different
	// build args.
	Visited *VisitedCollection
	// Final is the main target to be built.
	Final *SingleTarget
}

// FinalTarget returns the final target of the states.
func (mts *MultiTarget) FinalTarget() domain.Target {
	return mts.Final.Target
}

// All returns all SingleTarget contained within.
func (mts *MultiTarget) All() []*SingleTarget {
	return mts.Visited.VisitedList
}

type Phase int

const (
	PhaseMain Phase = iota
	PhasePush
	PhasePostPush
)

const (
	phaseShadow = -(iota + 1)
)

// SingleTarget holds LLB states representing an earthly target.
type SingleTarget struct {
	Target                 domain.Target
	Platform               *specs.Platform
	TargetInput            dedup.TargetInput
	CurrentPhase           Phase
	States                 map[Phase]llb.State
	Images                 map[Phase]*image.Image
	ArtifactsState         llb.State
	SeparateArtifactsState []llb.State
	SaveLocals             map[Phase][]SaveLocal
	SaveImages             []SaveImage //phase
	VarCollection          *variables.Collection
	LocalDirs              map[string]string
	Ongoing                bool
	Salt                   string
	// HasDangling represents whether the target has dangling instructions -
	// ie if there are any non-SAVE commands after the first SAVE command,
	// or if the target is invoked via BUILD command (not COPY nor FROM).
	HasDangling bool
}

// LastSaveImage returns the last save image available (if any).
func (sts *SingleTarget) LastSaveImage() SaveImage {
	if len(sts.SaveImages) == 0 {
		// Use main state / image if no save image exists.
		return SaveImage{
			State: sts.StateForPhase(sts.CurrentPhase),
			Image: sts.ImageForPhase(sts.CurrentPhase),
		}
	}
	return sts.SaveImages[len(sts.SaveImages)-1]
}

func (sts *SingleTarget) StateForPhase(p Phase) llb.State {
	return sts.States[p]
}

func (sts *SingleTarget) ImageForPhase(p Phase) *image.Image {
	return sts.Images[p]
}

func (sts *SingleTarget) NextPhase(state llb.State) error {
	nextState := Phase(int(sts.CurrentPhase + 1))
	if nextState > PhasePostPush {
		return errors.New("already in end phase")
	}

	sts.States[nextState] = state
	sts.Images[nextState] = sts.Images[sts.CurrentPhase].Clone()

	sts.CurrentPhase = nextState

	return nil
}

func (sts *SingleTarget) HasPhase(p Phase) bool {
	return p <= sts.CurrentPhase
}

func (sts *SingleTarget) CurrentState() llb.State {
	return sts.StateForPhase(sts.CurrentPhase)
}

func (sts *SingleTarget) CurrentImage() *image.Image {
	return sts.ImageForPhase(sts.CurrentPhase)
}

func (sts *SingleTarget) SetCurrentState(f llbutil.StatesAdapter) {
	sts.States[sts.CurrentPhase] = f(sts.States[sts.CurrentPhase])
	sts.States[phaseShadow] = f(sts.States[phaseShadow])
}

func (sts *SingleTarget) SetCurrentStateWithAltShadow(real, shadow llb.State) {
	sts.States[sts.CurrentPhase] = real
	sts.States[phaseShadow] = shadow
}

func (sts *SingleTarget) ShadowState() llb.State {
	return sts.States[phaseShadow]
}

func (sts *SingleTarget) SetCurrentImage(image *image.Image) {
	sts.Images[sts.CurrentPhase] = image
}

func (sts *SingleTarget) CurrentSaveLocal() []SaveLocal {
	return sts.SaveLocals[sts.CurrentPhase]
}

func (sts *SingleTarget) SetCurrentSaveLocal(locals []SaveLocal) {
	sts.SaveLocals[sts.CurrentPhase] = locals
}

func (sts *SingleTarget) SaveLocalsForPhase(p Phase) []SaveLocal {
	return sts.SaveLocals[p]
}

func DirectAssign(new llb.State) llbutil.StatesAdapter {
	return func(old llb.State) llb.State {
		return new
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
	State        llb.State
	Image        *image.Image
	DockerTag    string
	Push         bool
	InsecurePush bool
	// CacheHint instructs Earthly to save a separate ref for this image, even if no tag is
	// provided.
	CacheHint bool
}
