package inputgraph

import (
	"context"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/variables"
)

// HashOpt contains all of the options available to the hasher.
type HashOpt struct {
	Target           domain.Target
	Console          conslogging.ConsoleLogger
	CI               bool
	BuiltinArgs      variables.DefaultArgs
	OverridingVars   *variables.Scope
	EarthlyCIRunner  bool
	SkipProjectCheck bool
}

// HashTarget produces a hash from an Earthly target.
func HashTarget(ctx context.Context, opt HashOpt) (org string, hash []byte, err error) {
	l := newLoader(ctx, opt)

	if !opt.SkipProjectCheck {
		org = l.findOrg(ctx)
	}

	err = l.load(ctx)
	if err != nil {
		return "", nil, err
	}

	return org, l.hasher.GetHash(), nil
}
