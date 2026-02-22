package logger

import (
	"context"
	"errors"
)

type contextKey struct{}

func WithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

func FromContext(ctx context.Context) (*Logger, error) {
	l, ok := ctx.Value(contextKey{}).(*Logger)
	if !ok {
		return l, errors.New("Logger not present in the context")
	}
	return l, nil
}
