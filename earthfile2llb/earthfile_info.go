package earthfile2llb

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/ast/commandflag"
	"github.com/earthly/earthly/ast/spec"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/platutil"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

// These are functions that are used for getting information about an Earthfile,
// most notably for `earthly doc` and `earthly ls` output.

// GetTargets returns a list of targets from an Earthfile.
// Note that the passed in domain.Target's target name is ignored (only the reference to the Earthfile is used)
func GetTargets(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, target domain.Target) ([]string, error) {
	platr := platutil.NewResolver(platutil.GetUserPlatform())
	bc, err := resolver.Resolve(ctx, gwClient, platr, target)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve build context for target %s", target.String())
	}
	targets := make([]string, 0, len(bc.Earthfile.Targets))
	for _, target := range bc.Earthfile.Targets {
		targets = append(targets, target.Name)
	}
	return targets, nil
}

// GetTargetArgs returns a list of build arguments for a specified target
func GetTargetArgs(ctx context.Context, resolver *buildcontext.Resolver, gwClient gwclient.Client, target domain.Target) ([]string, error) {
	platr := platutil.NewResolver(platutil.GetUserPlatform())
	bc, err := resolver.Resolve(ctx, gwClient, platr, target)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve build context for target %s", target.String())
	}
	var t *spec.Target
	for _, tt := range bc.Earthfile.Targets {
		if tt.Name == target.Target {
			t = &tt
			break
		}
	}
	if t == nil {
		return nil, fmt.Errorf("failed to find %s", target.String())
	}
	var args []string
	for _, stmt := range t.Recipe {
		if stmt.Command != nil && stmt.Command.Name == "ARG" {
			isBase := t.Name == "base"
			// since Arg opts are ignored (and feature flags are not available) we set explicitGlobalArgFlag as false
			explicitGlobal := false
			_, argName, _, err := parseArgArgs(ctx, *stmt.Command, isBase, explicitGlobal)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse ARG arguments %v", stmt.Command.Args)
			}
			args = append(args, argName)
		}
	}
	return args, nil

}

// ArgName returns the parsed name of an ARG command, the default value (if
// any), and the state of the --required and --global flags.
func ArgName(ctx context.Context, cmd spec.Command, isBase bool, explicitGlobal bool) (_ string, _ *string, isRequired bool, isGlobal bool, _ error) {
	if cmd.Name != "ARG" {
		return "", nil, false, false, errors.Errorf("ArgName was called with non-arg command type '%v'", cmd.Name)
	}
	opts, argName, dflt, err := parseArgArgs(ctx, cmd, isBase, explicitGlobal)
	if err != nil {
		return "", nil, false, false, errors.Wrapf(err, "could not parse opts for ARG [%v]", cmd)
	}
	return argName, dflt, opts.Required, opts.Global, nil
}

// ArtifactName returns the parsed name of a SAVE ARTIFACT command and its local
// name (if any).
func ArtifactName(ctx context.Context, cmd spec.Command) (string, *string, error) {
	from, to, asLocal, ok := parseSaveArtifactArgs(cmd.Args)
	if !ok {
		return "", nil, errors.Errorf("could not parse opts for SAVE TARGET [%v]", cmd)
	}
	if to == "./" {
		to = from
	}
	if asLocal == "" {
		return to, nil, nil
	}
	return to, &asLocal, nil
}

// ImageNames returns the parsed names of a SAVE IMAGE command.
func ImageNames(ctx context.Context, cmd spec.Command) ([]string, error) {
	var opts commandflag.SaveImageOpts
	args, err := parseArgs("SAVE IMAGE", &opts, getArgsCopy(cmd))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid SAVE IMAGE arguments %v", cmd.Args)
	}
	return args, nil
}
