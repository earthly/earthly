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
	mu sync.Mutex

	visited map[string][]*SingleTarget
	// Same collection as above, but as a list, to keep the ordering consistent.
	visitedList []*SingleTarget
}

// NewVisitedCollection returns a collection of visited targets.
func NewVisitedCollection() *VisitedCollection {
	return &VisitedCollection{
		visited: make(map[string][]*SingleTarget),
	}
}

// Add adds a target to the collection, if it hasn't yet been visited. The returned sts is
// either the previously visited one or a brand new one.
func (vc *VisitedCollection) Add(target domain.Target, platform *specs.Platform, overridingVars *variables.Scope) (*SingleTarget, bool, error) {
	vc.waitAllDoneAndLock(target)
	defer vc.mu.Unlock()
	targetStr := target.String()
	for _, sts := range vc.visited[target.String()] {
		same, err := CompareTargetInputs(target, platform, overridingVars, sts.TargetInput)
		if err != nil {
			return nil, false, err
		}
		if same {
			return sts, true, nil
		}
	}
	sts := &SingleTarget{
		Target: target,
		TargetInput: dedup.TargetInput{
			TargetCanonical: target.StringCanonical(),
			Platform:        llbutil.PlatformWithDefaultToString(platform),
		},
		MainState:      llbutil.ScratchWithPlatform(),
		MainImage:      image.NewImage(),
		ArtifactsState: llbutil.ScratchWithPlatform(),
		Salt:           fmt.Sprintf("%d", rand.Int()),
		Done:           make(chan struct{}),
	}
	vc.visited[targetStr] = append(vc.visited[targetStr], sts)
	vc.visitedList = append(vc.visitedList, sts)
	return sts, false, nil
}

// waitAllDoneAndLock acquires mu at a point when all sts are done for a particular
// target, allowing for comparisons across the board while the lock is held.
func (vc *VisitedCollection) waitAllDoneAndLock(target domain.Target) {
	lenList := 0
	for {
		vc.mu.Lock()
		list := append([]*SingleTarget{}, vc.visited[target.String()]...)
		if lenList == len(list) {
			return // no unlocking on purpose
		}
		lenList = len(list)
		vc.mu.Unlock()
		// Wait for sts's to be done outside of the mu lock.
		for _, sts := range list {
			<-sts.Done
		}
	}
}

// All returns the list of visited members as a slice.
func (vc *VisitedCollection) All() []*SingleTarget {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return append([]*SingleTarget{}, vc.visitedList...)
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
