package goxp

import (
	"context"
	"iter"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

// DoWithWorker iterate chan and run do() with n workers
// if works <=0 then worker set to runtime.NumCPU()
func DoWithWorker(ctx context.Context, workers int, do func(i int) error) error {
	eg, _ := errgroup.WithContext(ctx)

	workers = Ternary(workers <= 0, runtime.NumCPU(), workers)
	eg.SetLimit(workers)

	for i := 0; i < workers; i++ {
		i := i
		eg.Go(func() error { return do(i) })
	}

	return eg.Wait()
}

// Every execute fn() in every time interval
//
// if you want run scheduled task like cron spec. please see github.com/robfig/cron
func Every(ctx context.Context, interval time.Duration, initialRun bool, fn func()) error {
	if initialRun {
		fn()
	}

	for {
		after := time.NewTimer(interval)

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-after.C:
			fn()
		}
	}
}

// After run func after duration
func After(ctx context.Context, duration time.Duration, fn func() error) error {
	after := time.NewTimer(duration)

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-after.C:
		return fn()
	}
}

// Async run func and returns with channel
func Async[T any](fn func() T) iter.Seq[T] {
	ch := make(chan T)
	go func() {
		ch <- fn()
		close(ch)
	}()

	return func(yield func(T) bool) {
		for c := range ch {
			if !yield(c) {
				return
			}
		}
	}
}

// Async2 run func and returns with channel
func Async2[U1, U2 any](fn func() (U1, U2)) iter.Seq2[U1, U2] {
	ch := make(chan *Tuple2[U1, U2])
	go func() {
		ch <- T2(fn())
		close(ch)
	}()

	return func(yield func(U1, U2) bool) {
		for v := range ch {
			if !yield(v.V1, v.V2) {
				return
			}
		}
	}
}
