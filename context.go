package goxp

import (
	"context"
	"time"
)

// IsContextDone return true if context is done
func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func WithTimeout(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(ctx)
}
