package inputgraph

import (
	"context"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/variables"
)

// HashOpt contains all of the options available to the hasher.
type HashOpt struct {
	Target         domain.Target
	Console        conslogging.ConsoleLogger
	CI             bool
	BuiltinArgs    variables.DefaultArgs
	OverridingVars *variables.Scope
}

// HashTarget produces a hash from an Earthly target.
func HashTarget(ctx context.Context, opt HashOpt) ([]byte, Stats, error) {
	l := newLoader(ctx, opt)

	b, err := l.load(ctx)
	if err != nil {
		return nil, Stats{}, err
	}

	stats := Stats{}
	if l.stats != nil {
		stats = *l.stats
	}

	return b, stats, nil
}
