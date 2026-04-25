package config

import (
	"context"
	"errors"
)

type contextKey struct{}

func WithDns(ctx context.Context, d *Dns) context.Context {
	return context.WithValue(ctx, contextKey{}, d)
}

func DnsFromContext(ctx context.Context) (*Dns, error) {
	d, ok := ctx.Value(contextKey{}).(*Dns)
	if !ok {
		return d, errors.New("Logger not present in the context")
	}
	return d, nil
}
