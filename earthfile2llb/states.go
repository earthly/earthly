package earthfile2llb

import (
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/dedup"
	"github.com/earthly/earthly/earthfile2llb/image"
	"github.com/moby/buildkit/client/llb"
)

// MultiTargetStates holds LLB states representing multiple earth targets,
// in the order in which they should be built.
type MultiTargetStates struct {
	// VisitedStates represents the previously visited states, grouped by target
	// name. Duplicate targets are possible if same target is called with different
	// build args.
	VisitedStates map[string][]*SingleTargetStates
	// FinalStates is the main target to be built.
	FinalStates *SingleTargetStates
}

// FinalTarget returns the final target of the states.
func (mts *MultiTargetStates) FinalTarget() domain.Target {
	return mts.FinalStates.Target
}

// AllStates returns all SingleTargetStates contained within.
func (mts *MultiTargetStates) AllStates() []*SingleTargetStates {
	var ret []*SingleTargetStates
	for _, stss := range mts.VisitedStates {
		ret = append(ret, stss...)
	}
	return ret
}

// SingleTargetStates holds LLB states representing a earth target.
type SingleTargetStates struct {
	Target                 domain.Target
	TargetInput            dedup.TargetInput
	SideEffectsImage       *image.Image
	SideEffectsState       llb.State
	ArtifactsState         llb.State
	SeparateArtifactsState []llb.State
	SaveLocals             []SaveLocal
	SaveImages             []SaveImage
	RunPush                RunPush
	LocalDirs              map[string]string
	Ongoing                bool
}

// LastSaveImage returns the last save image available (if any).
func (sts *SingleTargetStates) LastSaveImage() (SaveImage, bool) {
	if len(sts.SaveImages) == 0 {
		return SaveImage{}, false
	}
	return sts.SaveImages[len(sts.SaveImages)-1], true
}

// SaveLocal is an artifact path to be saved to local disk.
type SaveLocal struct {
	// DestPath is the local dest path to copy the artifact to.
	DestPath string
	// ArtifactPath is the relative path within the artifacts image.
	ArtifactPath string
	// Index is the index number of the "save as local" command encountered. Starts as 0.
	Index int
}

// SaveImage is a docker image to be saved.
type SaveImage struct {
	State     llb.State
	Image     *image.Image
	DockerTag string
	Push      bool
}

// RunPush is a series of RUN --push commands to be run after the build has been deemed as
// successful.
type RunPush struct {
	Initialized bool
	CommandStrs []string
	State       llb.State
}
