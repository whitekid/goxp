package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	type args struct {
		limit   int
		initial time.Duration
		backoff float64
		fn      func() error
	}

	tests := [...]struct {
		name    string
		args    args
		wantErr bool
		tries   int
	}{
		{"default", args{3, time.Millisecond * 100, 1.3, func() error { return nil }}, false, 1},
		{"error", args{3, time.Millisecond * 100, 1.3, func() error { return errors.New("fail") }}, true, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tries := 0

			err := Retry().Limit(tt.args.limit).Backoff(tt.args.initial, tt.args.backoff).
				Do(ctx, func() error {
					tries++
					return tt.args.fn()
				})
			if (err != nil) != tt.wantErr {
				require.Fail(t, "retry failed", "error: %v, want: %v", err, tt.wantErr)
			}

			require.Equal(t, tt.tries, tries)
		})
	}
}
