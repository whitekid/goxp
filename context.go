package goxp

import "context"

// IsContextDone return true if context is done
// Depreciated: use ValidContext()
func IsContextDone(ctx context.Context) bool {
	return !ValidContext(ctx)
}

// ValidContext return true if context is valid(not canceled or deadline exceed)
func ValidContext(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return true
	}
}
