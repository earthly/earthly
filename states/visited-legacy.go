package states

import (
	"context"
	"sync"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/variables"
	"github.com/pkg/errors"
)

// legacyVisitedCollection is a collection of visited targets.
type legacyVisitedCollection struct {
	mu      sync.Mutex
	visited map[string][]*SingleTarget // targetStr -> sts list
	// Same collection as above, but as a list, to make the ordering consistent.
	visitedList []*SingleTarget
}

// NewLegacyVisitedCollection returns a collection of visited targets.
func NewLegacyVisitedCollection() VisitedCollection {
	return &legacyVisitedCollection{
		visited: make(map[string][]*SingleTarget),
	}
}

// All returns all visited items.
func (vc *legacyVisitedCollection) All() []*SingleTarget {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return append([]*SingleTarget{}, vc.visitedList...)
}

// Add adds a target to the collection, if it hasn't yet been visited. The returned sts is
// either the previously visited one or a brand new one.
func (vc *legacyVisitedCollection) Add(ctx context.Context, target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, parentDepSub chan string) (*SingleTarget, bool, error) {
	dependents, err := vc.waitAllDoneAndLock(ctx, target, parentDepSub)
	if err != nil {
		return nil, false, err
	}
	// Put the deps back into the channel for the new sts to consume.
	for depID := range dependents {
		parentDepSub <- depID
	}
	defer vc.mu.Unlock()
	for _, sts := range vc.visited[target.StringCanonical()] {
		same, err := compareTargetInputs(target, platr, allowPrivileged, overridingVars, sts.TargetInput())
		if err != nil {
			return nil, false, err
		}
		if same {
			// Existing sts.
			if dependents[sts.ID] {
				// Infinite recursion. The previously visited sts is a dependent of us.
				return nil, false, errors.Errorf(
					"infinite recursion detected for target %s", target.String())
			}
			// If it's not a dependent, then it *has* to be done at this point.
			// Sanity check.
			select {
			case <-sts.Done():
			default:
				panic("same sts but not done")
			}
			// Subscribe that sts to the dependencies of our parent.
			sts.MonitorDependencySubscription(ctx, parentDepSub)
			return sts, true, nil
		}
	}
	// None are the same. Create new sts.
	sts, err := newSingleTarget(ctx, target, platr, allowPrivileged, overridingVars, parentDepSub)
	if err != nil {
		return nil, false, err
	}
	targetStr := target.StringCanonical()
	vc.visited[targetStr] = append(vc.visited[targetStr], sts)
	vc.visitedList = append(vc.visitedList, sts)
	return sts, false, nil
}

// waitAllDoneAndLock acquires mu at a point when all sts are done for a particular
// target, allowing for comparisons across the board while the lock is held.
func (vc *legacyVisitedCollection) waitAllDoneAndLock(ctx context.Context, target domain.Target, parentDepSub chan string) (map[string]bool, error) {
	// Build up dependents from parentDepSub. The list needs to be complete when returning
	// from this function for proper infinite loop detection.
	dependents := make(map[string]bool)
	// wait all done & lock loop
	prevLenList := 0
	for {
		vc.mu.Lock()
		list := append([]*SingleTarget{}, vc.visited[target.StringCanonical()]...)
		if prevLenList == len(list) {
			// The list we have now is the same we just checked if it's done or waiting on us.
			// We are finished.
			return dependents, nil // no unlocking on purpose
		}
		prevLenList = len(list)
		vc.mu.Unlock()
		// Wait for sts's to be done outside of the mu lock.
	stsLoop:
		for _, sts := range list {
			if dependents[sts.ID] {
				// No need to wait if it's a dependent, because the sts is waiting on us.
				// It's safe to perform comparison if they are waiting on us.
				continue
			}
			for {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case newID := <-parentDepSub:
					dependents[newID] = true
					if newID == sts.ID {
						// Just the one we were waiting for. It seems that it is
						// now waiting on us.
						continue stsLoop
					}
				case <-sts.Done():
					continue stsLoop
				}
			}
		}
	}
}

// compareTargetInputs compares two targets and their inputs to check if they are the same.
func compareTargetInputs(target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, other dedup.TargetInput) (bool, error) {
	if target.StringCanonical() != other.TargetCanonical {
		return false, nil
	}
	if allowPrivileged != other.AllowPrivileged {
		return false, nil
	}
	stsPlat, err := platr.ParseAllowNativeAndUser(other.Platform)
	if err != nil {
		return false, err
	}
	if !platr.PlatformEquals(platr.Current(), stsPlat) {
		return false, nil
	}
	for _, bai := range other.BuildArgs {
		variable, found := overridingVars.Get(bai.Name)
		if found {
			baiVariable := dedup.BuildArgInput{
				Name:          bai.Name,
				DefaultValue:  bai.DefaultValue,
				ConstantValue: variable.Str,
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
