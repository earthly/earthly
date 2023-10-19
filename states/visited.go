package states

import (
	"context"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/variables"
)

type VisitedCollection interface {
	All() []*SingleTarget
	Add(ctx context.Context, target domain.Target, platr *platutil.Resolver, allowPrivileged bool, overridingVars *variables.Scope, parentDepSub chan string) (*SingleTarget, bool, error)
}
