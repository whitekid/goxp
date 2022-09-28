package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	ctxDone, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx     context.Context
		limit   int
		initial time.Duration
		ratio   float64
		fn      func() error
	}

	tests := [...]struct {
		name      string
		args      args
		wantErr   bool
		wantTries int
	}{
		{"default", args{nil, 3, time.Millisecond * 100, 1.3, func() error { return nil }}, false, 1},
		{"error", args{nil, 3, time.Millisecond * 100, 1.3, func() error { return errors.New("fail") }}, true, 3},
		{"stop", args{nil, 3, time.Millisecond * 100, 1.3, func() error { return ErrStop }}, true, 1},
		{"done context", args{ctxDone, 3, time.Millisecond * 100, 1.3, func() error { return nil }}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if tt.args.ctx != nil {
				ctx = tt.args.ctx
			}

			tries := 0
			err := New().Limit(tt.args.limit).Backoff(tt.args.initial, tt.args.ratio).
				Do(ctx, func() error {
					tries++
					return tt.args.fn()
				})
			if (err != nil) != tt.wantErr {
				require.Fail(t, "retry failed", "error: %v, want: %v", err, tt.wantErr)
			}

			require.Equal(t, tt.wantTries, tries)
		})
	}
}
