package goxp

import (
	"context"
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
func Every(ctx context.Context, interval time.Duration, initialRun bool, fn func() error, errC chan<- error) {
	firstInterval, origInterval := interval, interval
	if initialRun {
		firstInterval = 0
	}
	firstRun := true

	for {
		if firstRun {
			interval = firstInterval
			firstRun = false
		} else {
			interval = origInterval
		}

		after := time.NewTimer(interval)

		select {
		case <-ctx.Done():
			if !after.Stop() {
				go func() { <-after.C }()
			}
			return

		case <-after.C:
			if err := fn(); err != nil {
				errC <- err
			}
		}
	}
}

// After run func after duration
func After(ctx context.Context, duration time.Duration, fn func() error) error {
	after := time.NewTimer(duration)

	select {
	case <-ctx.Done():
		if !after.Stop() {
			go func() { <-after.C }()
		}
		return ctx.Err()

	case <-after.C:
		return fn()
	}
}
