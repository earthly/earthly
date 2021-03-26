package states

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/variables"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// VisitedCollection is a collection of visited targets.
type VisitedCollection struct {
	mu      sync.Mutex
	visited map[string][]*SingleTarget // targetStr -> sts list
	// Same collection as above, but as a list, to make the ordering consistent.
	visitedList []*SingleTarget
}

// NewVisitedCollection returns a collection of visited targets.
func NewVisitedCollection() *VisitedCollection {
	return &VisitedCollection{
		visited: make(map[string][]*SingleTarget),
	}
}

// Add adds a target to the collection.
func (vc *VisitedCollection) Add(target domain.Target, platform *specs.Platform) *SingleTarget {
	targetStr := target.StringCanonical()
	sts := &SingleTarget{
		Target:   target,
		Platform: platform,
		targetInput: dedup.TargetInput{
			TargetCanonical: targetStr,
			Platform:        llbutil.PlatformWithDefaultToString(platform),
		},
		MainState:      llbutil.ScratchWithPlatform(),
		MainImage:      image.NewImage(),
		ArtifactsState: llbutil.ScratchWithPlatform(),
		Ongoing:        true,
		LocalDirs:      make(map[string]string),
		Salt:           fmt.Sprintf("%d", rand.Int()),
	}
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.visited[targetStr] = append(vc.visited[targetStr], sts)
	vc.visitedList = append(vc.visitedList, sts)
	return sts
}

// All returns all visited items.
func (vc *VisitedCollection) All() []*SingleTarget {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return append([]*SingleTarget{}, vc.visitedList...)
}

// AllTarget returns all visited items for a given target.
func (vc *VisitedCollection) AllTarget(target domain.Target) []*SingleTarget {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	targetStr := target.StringCanonical()
	return append([]*SingleTarget{}, vc.visited[targetStr]...)
}

// CompareTargetInputs compares two targets and their inputs to check if they are the same.
func CompareTargetInputs(target domain.Target, platform *specs.Platform, overridingVars *variables.Scope, other dedup.TargetInput) (bool, error) {
	if target.StringCanonical() != other.TargetCanonical {
		return false, nil
	}
	stsPlat, err := llbutil.ParsePlatform(other.Platform)
	if err != nil {
		return false, err
	}
	if !llbutil.PlatformEquals(stsPlat, platform) {
		return false, nil
	}
	for _, bai := range other.BuildArgs {
		variable, found := overridingVars.GetAny(bai.Name)
		if found {
			baiVariable := dedup.BuildArgInput{
				Name:          bai.Name,
				DefaultValue:  bai.DefaultValue,
				ConstantValue: variable,
			}
			if !baiVariable.Equals(bai) {
				return false, nil
			}
		} else {
			if !bai.IsDefaultValue() {
				return false, nil
			}
		}
	}
	return true, nil
}
