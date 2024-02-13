package inputgraph

import (
	"context"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/buildkitskipper/hasher"
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

	// Bypass further analysis for remote targets as there's nothing to do
	// beyond hashing the full target name.
	if t := opt.Target; t.IsRemote() {
		if supportedRemoteTarget(t) {
			h := hasher.New()
			h.HashString(t.StringCanonical())
			return h.GetHash(), Stats{}, nil
		}
		return nil, Stats{}, errInvalidRemoteTarget
	}

	// Continue processing local targets (which may include remote transitive targets).
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
