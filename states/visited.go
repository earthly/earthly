package states

import (
	"context"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/variables"
)

// VisitedCollection represents a collection of visited targets.
type VisitedCollection interface {
	// All returns all visited items.
	All() []*SingleTarget
	// Add adds a target to the collection, if it hasn't yet been visited. The returned sts is
	// either the previously visited one or a brand new one.
	Add(ctx context.Context, target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, parentDepSub chan string) (*SingleTarget, bool, error)
}
