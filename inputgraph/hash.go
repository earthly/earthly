package inputgraph

import (
	"context"
	"errors"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/variables"
)

// HashOpt contains all of the options available to the hasher.
type HashOpt struct {
	Target          domain.Target
	Console         conslogging.ConsoleLogger
	CI              bool
	Push            bool
	BuiltinArgs     variables.DefaultArgs
	OverridingVars  *variables.Scope
	EarthlyCIRunner bool
}

// HashTarget produces a hash from an Earthly target.
func HashTarget(ctx context.Context, opt HashOpt) (org, project string, hash []byte, err error) {
	l := newLoader(ctx, opt)

	org, project, err = l.findProject(ctx)
	if err != nil {
		return "", "", nil, err
	}

	err = l.load(ctx)
	if err != nil {
		return "", "", nil, errors.Join(ErrUnableToDetermineHash, err)
	}

	return org, project, l.hasher.GetHash(), nil
}
