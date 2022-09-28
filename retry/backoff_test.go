package retry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/log"
)

func TestBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type args struct {
		backoff backoffer
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"zero", args{newZeroBackoff(ctx)}},
		{"fixed", args{newFixedBackoff(ctx, 100*time.Millisecond)}},
		{"ratio", args{newRatioBackoff(ctx, 100*time.Millisecond, time.Second, 0.1)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := time.Now()
			for i := 0; i < 5; i++ {
				select {
				case <-ctx.Done():
					require.Fail(t, "deadline exceed")
				case <-tt.args.backoff.NextBackoff():
					now := time.Now()
					log.Debugf("backoff: %s", now.Sub(prev))
					prev = now
				}
			}
		})
	}
}
