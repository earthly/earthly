package earthfile2llb

import (
	"context"

	"github.com/earthly/earthly/ast/spec"
)

type contextKey string

var (
	contextKeySourceLocation contextKey = "sourceLocation"
)

// ContextWithSourceLocation returns a new context with the given source location.
func ContextWithSourceLocation(ctx context.Context, sl *spec.SourceLocation) context.Context {
	if sl == nil {
		return ctx
	}
	return context.WithValue(ctx, contextKeySourceLocation, sl)
}

// SourceLocationFromContext returns the source location from the given context.
func SourceLocationFromContext(ctx context.Context) *spec.SourceLocation {
	sl, _ := ctx.Value(contextKeySourceLocation).(*spec.SourceLocation)
	return sl
}
