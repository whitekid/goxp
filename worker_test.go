package goxp

import (
	"context"
	"sync"
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

	got := <-Async(ctx, func(ctx context.Context) int {
		time.Sleep(time.Second)
		return 7
	})

	require.Equal(t, 7, got)
}

func TestAsync2(t *testing.T) {
	ctx := context.Background()
	got := <-Async2(ctx, func(ctx context.Context) (int, time.Time) {
		time.Sleep(time.Second)
		return 7, time.Now()
	})

	require.Equal(t, 7, got.V1)
	require.True(t, got.V2.Before(time.Now()))
}

// Test edge cases and error conditions
func TestDoWithWorkerErrorPropagation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test error propagation
	expectedErr := errors.New("worker error")
	err := DoWithWorker(ctx, 2, func(ctx context.Context, i int) error {
		if i == 0 {
			return expectedErr
		}
		// This should be canceled when the other worker returns an error
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Hour):
			return nil
		}
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "worker error")
}

func TestDoWithWorkerIndexCapture(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerCount := 5
	indices := make([]int, workerCount)
	var mu sync.Mutex

	err := DoWithWorker(ctx, workerCount, func(ctx context.Context, i int) error {
		mu.Lock()
		indices[i] = i // This tests the closure variable capture fix
		mu.Unlock()
		return nil
	})

	require.NoError(t, err)

	// Verify each worker got the correct index
	for i := 0; i < workerCount; i++ {
		require.Equal(t, i, indices[i], "Worker %d should have index %d", i, i)
	}
}

func TestEveryPanicRecovery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	callCount := int32(0)

	// Every function doesn't recover from panics, so test that it panics
	require.Panics(t, func() {
		Every(ctx, 50*time.Millisecond, true, func(ctx context.Context) {
			atomic.AddInt32(&callCount, 1)
			if atomic.LoadInt32(&callCount) == 1 {
				// First call panics
				panic("test panic")
			}
		})
	})

	// Should have called at least once before panicking
	require.Greater(t, atomic.LoadInt32(&callCount), int32(0))
}

func TestEveryContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	callCount := int32(0)
	done := make(chan struct{})

	go func() {
		defer close(done)
		Every(ctx, 10*time.Millisecond, false, func(ctx context.Context) {
			atomic.AddInt32(&callCount, 1)
		})
	}()

	time.Sleep(50 * time.Millisecond) // Let it run a few times
	cancel()                          // Cancel context

	select {
	case <-done:
		// Should exit quickly after cancellation
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Every did not respect context cancellation")
	}

	finalCount := atomic.LoadInt32(&callCount)
	require.Greater(t, finalCount, int32(0))
	require.Less(t, finalCount, int32(20)) // Should not run too many times
}

func TestAfterConcurrentCalls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const numCalls = 10
	results := make(chan int, numCalls)

	// Start multiple After calls concurrently
	for i := 0; i < numCalls; i++ {
		i := i
		go func() {
			err := After(ctx, 10*time.Millisecond, func(ctx context.Context) error {
				results <- i
				return nil
			})
			require.NoError(t, err)
		}()
	}

	// Collect all results
	received := make(map[int]bool)
	for i := 0; i < numCalls; i++ {
		select {
		case result := <-results:
			received[result] = true
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for After results")
		}
	}

	// Verify all calls completed
	require.Len(t, received, numCalls)
}

func TestAsyncContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := Async(ctx, func(ctx context.Context) int {
		select {
		case <-ctx.Done():
			return -1 // Signal that context was canceled
		case <-time.After(time.Hour):
			return 42
		}
	})

	// Cancel context immediately
	cancel()

	select {
	case result := <-ch:
		require.Equal(t, -1, result) // Should receive cancellation signal
	case <-time.After(time.Second):
		t.Fatal("Async did not respect context cancellation")
	}
}

func TestAsyncChannelBuffering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test that channel is buffered (won't block sender)
	ch := Async(ctx, func(ctx context.Context) int {
		return 42
	})

	// Don't read from channel immediately
	time.Sleep(100 * time.Millisecond)

	// Should still be able to read the result
	select {
	case result := <-ch:
		require.Equal(t, 42, result)
	case <-time.After(time.Second):
		t.Fatal("Channel should be buffered")
	}
}
