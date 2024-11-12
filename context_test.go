package goxp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithTimeout(t *testing.T) {
	err := WithTimeout(context.Background(), time.Second, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Minute):
			return nil
		}
	})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
