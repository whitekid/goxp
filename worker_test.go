package goxp

import (
	"context"
	"iter"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/whitekid/goxp/errors"
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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := DoWithWorker(ctx, tt.args.workers, func(ctx context.Context, i int) error {
				for x := range ch {
					atomic.AddInt32(&sum, x)
				}
				return nil
			})
			require.NoError(t, err)

			require.Equal(t, int32(tt.want), sum)
		})
	}
}

func TestDoWithWorkerCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	t1 := time.Now()
	DoWithWorker(ctx, 0, func(ctx context.Context, i int) error {
		after := time.NewTimer(time.Hour)

		select {
		case <-ctx.Done():
			break
		case <-after.C:
			require.Fail(t, "must canceled by context")
		}

		return nil
	})

	require.Truef(t, time.Now().Before(t1.Add(time.Second)), "work should done in %s, it takes %s", time.Second, time.Since(t1))
}

func TestEvery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callCount := int32(0)
	Every(ctx, 100*time.Second, true, func(ctx context.Context) {
		atomic.AddInt32(&callCount, 1)
	})
	require.Greater(t, callCount, int32(0))
}

func TestAfter(t *testing.T) {
	type args struct {
		ret error
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{}, false},
		{"err", args{errors.New("error")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := After(ctx, 100*time.Millisecond, func(ctx context.Context) error { return tt.args.ret })
			require.Truef(t, (err != nil) == tt.wantErr, `After() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
		})
	}
}

func TestAfterContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := After(ctx, 200*time.Millisecond, func(ctx context.Context) error { return nil })
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestAsync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	it := Async(ctx, func(ctx context.Context) int {
		time.Sleep(time.Second)
		return 7
	})
	next, stop := iter.Pull(it)
	defer stop()

	got, ok := next()
	require.True(t, ok)
	require.Equal(t, 7, got)
}

func TestAsync2(t *testing.T) {
	it := Async2(func() (int, time.Time) {
		time.Sleep(time.Second)
		return 7, time.Now()
	})
	next, stop := iter.Pull2(it)
	defer stop()
	v1, v2, ok := next()
	require.True(t, ok)

	require.Equal(t, 7, v1)
	require.True(t, v2.Before(time.Now()))
}
