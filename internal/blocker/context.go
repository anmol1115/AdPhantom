package blocker

import (
	"context"
	"errors"
)

type contextKey struct{}

func WithFilterList(ctx context.Context, fl *FilterList) context.Context {
	return context.WithValue(ctx, contextKey{}, fl)
}

func FromContext(ctx context.Context) (*FilterList, error) {
	l, ok := ctx.Value(contextKey{}).(*FilterList)
	if !ok {
		return l, errors.New("Filter List not present in the context")
	}
	return l, nil
}
