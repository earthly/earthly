package states

import (
	"context"
	"sync"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/variables"
)

// visitedUpfrontHashCollection is a collection of visited targets.
type visitedUpfrontHashCollection struct {
	mu      sync.Mutex
	visited map[string]*SingleTarget // targetInputHash -> sts
	// visitedList is the same collection as above, but as a list,
	// to make the ordering consistent.
	visitedList []*SingleTarget
}

// NewVisitedUpfrontHashCollection returns a collection of visited targets.
func NewVisitedUpfrontHashCollection() VisitedCollection {
	return &visitedUpfrontHashCollection{
		visited: make(map[string]*SingleTarget),
	}
}

// All returns all visited items.
func (vc *visitedUpfrontHashCollection) All() []*SingleTarget {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	return append([]*SingleTarget{}, vc.visitedList...)
}

// Add adds a target to the collection, if it hasn't yet been visited. The returned sts is
// either the previously visited one or a brand new one.
func (vc *visitedUpfrontHashCollection) Add(ctx context.Context, target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, parentDepSub chan string) (*SingleTarget, bool, error) {
	// Constructing a new sts early to be able to compute its target input hash.
	newSts, err := newSingleTarget(ctx, target, platr, allowPrivileged, overridingVars, nil)
	if err != nil {
		return nil, false, err
	}
	newKey, err := newSts.targetInput.Hash()
	if err != nil {
		return nil, false, err
	}
	vc.mu.Lock()
	defer vc.mu.Unlock()
	sts, found := vc.visited[newKey]
	if found {
		return sts, true, nil
	}
	vc.visited[newKey] = newSts
	vc.visitedList = append(vc.visitedList, newSts)
	return newSts, false, nil
}
