package resolver

import (
	"context"
	"errors"
)

type contextKey struct{}

func WithResolver(ctx context.Context, r *Resolver) context.Context {
	return context.WithValue(ctx, contextKey{}, r)
}

func FromContext(ctx context.Context) (*Resolver, error) {
	d, ok := ctx.Value(contextKey{}).(*Resolver)
	if !ok {
		return d, errors.New("Resolver not present in the context")
	}
	return d, nil
}
