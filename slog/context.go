package slog

import "context"

type contextKey string

const loggerContextKey = contextKey("logger")

// WithLogger returns a new context with a logger added to it.
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

// GetLogger returns a logger associated with this context.
func GetLogger(ctx context.Context) Logger {
	v := ctx.Value(loggerContextKey)
	if v == nil {
		return Logger{}
	}
	return v.(Logger) // Note that this panics if not a real Logger.
}

// With adds logging metadata to the logger within the context.
func With(ctx context.Context, key string, value interface{}) context.Context {
	return WithLogger(ctx, GetLogger(ctx).With(key, value))
}
