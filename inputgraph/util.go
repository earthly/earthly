package inputgraph

import (
	"context"
	"strings"

	"github.com/earthly/earthly/ast/command"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/pkg/errors"
)

func ParseProjectCommand(ctx context.Context, target domain.Target, console conslogging.ConsoleLogger) (string, string, error) {
	if target.IsRemote() {
		return "", "", errCannotLoadRemoteTarget
	}

	resolver := buildcontext.NewResolver(nil, nil, console, "", "", "", 0, "")

	buildCtx, err := resolver.Resolve(ctx, nil, nil, target)
	if err != nil {
		return "", "", err
	}

	ef := buildCtx.Earthfile

	for _, stmt := range ef.BaseRecipe {
		if stmt.Command != nil && stmt.Command.Name == command.Project {
			args := stmt.Command.Args
			if len(args) != 1 {
				return "", "", errors.New("failed to parse PROJECT command")
			}
			parts := strings.Split(args[0], "/")
			if len(parts) != 2 {
				return "", "", errors.New("failed to parse PROJECT command")
			}
			return parts[0], parts[1], nil
		}
	}

	return "", "", errors.New("PROJECT command is required for remote storage")
}

func copyVisited(m map[string]struct{}) map[string]struct{} {
	m2 := map[string]struct{}{}
	for k := range m {
		m2[k] = struct{}{}
	}
	return m2
}

func uniqStrs(all []string) []string {
	m := map[string]struct{}{}
	for _, v := range all {
		m[v] = struct{}{}
	}
	ret := []string{}
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}
