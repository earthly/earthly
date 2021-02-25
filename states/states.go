package states

import (
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

// SingleTarget holds LLB states representing an earthly target.
type SingleTarget struct {
	Target                 domain.Target
	Platform               *specs.Platform
	TargetInput            dedup.TargetInput
	MainImage              *image.Image
	MainState              llb.State
	ArtifactsState         llb.State
	SeparateArtifactsState []llb.State
	SaveLocals             []SaveLocal
	SaveImages             []SaveImage
	VarCollection          *variables.Collection
	RunPush                RunPush
	LocalDirs              map[string]string
	Ongoing                bool
	Salt                   string
	// HasDangling represents whether the target has dangling instructions -
	// ie if there are any non-SAVE commands after the first SAVE command,
	// or if the target is invoked via BUILD command (not COPY nor FROM).
	HasDangling bool
	// RanFromLike represents whether we have encountered a FROM-like command
	// (eg FROM, FROM DOCKERFILE, LOCALLY).
	RanFromLike bool
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
	CacheHint           bool
	HasPushDependencies bool
}

// RunPush is a series of RUN --push commands to be run after the build has been deemed as
// successful, along with artifacts to save and images to push
type RunPush struct {
	CommandStrs []string
	State       llb.State
	SaveLocals  []SaveLocal
	SaveImages  []SaveImage
	HasState    bool
}
