package goxp

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoWithWorker(t *testing.T) {
	type args struct {
		workers int
		sumTo   int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{"default", args{0, 10000}, 49995000},
		{"default", args{4, 1000}, 499500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := int32(0)

			ch := make(chan int32)

			go func() {
				defer close(ch)
				for i := 0; i < tt.args.sumTo; i++ {
					ch <- int32(i)
				}
			}()

			DoWithWorker(tt.args.workers, func(i int) {
				for x := range ch {
					atomic.AddInt32(&sum, x)
				}
			})

			require.Equal(t, int32(tt.want), sum)
		})
	}
}

func TestDoWithWorkerCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	t1 := time.Now()
	DoWithWorker(0, func(i int) {
		select {
		case <-ctx.Done():
			break
		case <-time.After(time.Hour):
			require.Fail(t, "must canceled by context")
		}
	})

	require.Truef(t, time.Now().Before(t1.Add(time.Second)), "work should done in %s, it takes %s", time.Second, time.Since(t1))
}
