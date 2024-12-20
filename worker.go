package goxp

import (
	"context"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

// DoWithWorker iterate chan and run do() with n workers
// if works <=0 then worker set to runtime.NumCPU()
func DoWithWorker(ctx context.Context, workers int, do func(ctx context.Context, i int) error) error {
	eg, _ := errgroup.WithContext(ctx)

	workers = Ternary(workers <= 0, runtime.NumCPU(), workers)
	eg.SetLimit(workers)

	for i := 0; i < workers; i++ {
		eg.Go(func() error { return do(ctx, i) })
	}

	return eg.Wait()
}

// Every execute fn() in every time interval, return when context is done.
//
// if you want run scheduled task like cron spec. please see github.com/robfig/cron
func Every(ctx context.Context, interval time.Duration, initialRun bool, fn func(ctx context.Context)) error {
	if initialRun {
		fn(ctx)

		if IsContextDone(ctx) {
			return ctx.Err()
		}
	}

	for {
		after := time.NewTimer(interval)

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-after.C:
			fn(ctx)
		}
	}
}

// After run func after duration
func After(ctx context.Context, duration time.Duration, fn func(ctx context.Context) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(duration):
		return fn(ctx)
	}
}

// Async run func in background and returns with iter.Seq
func Async[T any](ctx context.Context, fn func(ctx context.Context) T) <-chan T {
	ch := make(chan T)
	go func() {
		ch <- fn(ctx)
		close(ch)
	}()

	return ch
}

// Async2 run func in background and returns with iter.Seq
func Async2[U1, U2 any](fn func() (U1, U2)) <-chan *Tuple2[U1, U2] {
	ch := make(chan *Tuple2[U1, U2])
	go func() {
		ch <- T2(fn())
		close(ch)
	}()

	return ch
}
