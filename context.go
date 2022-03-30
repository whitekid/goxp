package goxp

import "context"

// IsContextDone return true if context is done
func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
