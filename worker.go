package goxp

import (
	"context"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

// DoWithWorker executes the provided function with a pool of worker goroutines.
// It creates 'workers' number of goroutines and calls do() for each worker with its index (0 to workers-1).
//
// If workers <= 0, it defaults to runtime.NumCPU().
// The function blocks until all workers complete or context is cancelled.
// Returns the first error encountered from any worker.
//
// Example:
//
//	err := DoWithWorker(ctx, 4, func(ctx context.Context, workerID int) error {
//	    fmt.Printf("Worker %d processing\n", workerID)
//	    return processData(ctx)
//	})
func DoWithWorker(ctx context.Context, workers int, do func(ctx context.Context, i int) error) error {
	eg, ctx := errgroup.WithContext(ctx)

	workers = Ternary(workers <= 0, runtime.NumCPU(), workers)
	eg.SetLimit(workers)

	for i := range workers {
		eg.Go(func() error { return do(ctx, i) })
	}

	return eg.Wait()
}

// Every executes fn() at regular time intervals until context is cancelled.
// If initialRun is true, fn() is called immediately before starting the ticker.
//
// The function returns when context is done, returning ctx.Err().
// For cron-style scheduling, see github.com/robfig/cron.
//
// Example:
//
//	err := Every(ctx, 5*time.Second, true, func(ctx context.Context) {
//	    fmt.Println("Running periodic task")
//	    updateMetrics(ctx)
//	})
func Every(ctx context.Context, interval time.Duration, initialRun bool, fn func(ctx context.Context)) error {
	if initialRun {
		fn(ctx)

		if IsContextDone(ctx) {
			return ctx.Err()
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			fn(ctx)
		}
	}
}

// After executes fn() once after the specified duration.
// The function respects context cancellation and returns ctx.Err() if cancelled before duration elapses.
//
// Example:
//
//	err := After(ctx, 10*time.Second, func(ctx context.Context) error {
//	    fmt.Println("Delayed execution")
//	    return processDelayedTask(ctx)
//	})
func After(ctx context.Context, duration time.Duration, fn func(ctx context.Context) error) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		return fn(ctx)
	}
}

// Async executes fn() in a goroutine and returns a channel that receives the result.
// The channel is closed after the result is sent.
//
// Example:
//
//	resultCh := Async(ctx, func(ctx context.Context) string {
//	    return fetchData(ctx)
//	})
//	result := <-resultCh
func Async[T any](ctx context.Context, fn func(ctx context.Context) T) <-chan T {
	ch := make(chan T, 1)
	go func() {
		ch <- fn(ctx)
		close(ch)
	}()

	return ch
}

// Async2 run func in background and returns with channel
// The goroutine respects context cancellation to prevent leaks
func Async2[U1, U2 any](ctx context.Context, fn func(ctx context.Context) (U1, U2)) <-chan *Tuple2[U1, U2] {
	ch := make(chan *Tuple2[U1, U2], 1)
	go func() {
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- T2(fn(ctx)):
		}
	}()

	return ch
}
