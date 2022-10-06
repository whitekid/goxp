package goxp

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/whitekid/goxp/fx"
)

// DoWithWorker iterate chan and run do() with n workers
// if works <=0 then worker set to runtime.NumCPU()
func DoWithWorker(workers int, do func(i int)) {
	var wg sync.WaitGroup

	workers = fx.Ternary(workers <= 0, runtime.NumCPU(), workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			do(i)
		}(i)
	}

	wg.Wait()
}

// Every execute fn() in every time interval
//
// if you want run scheduled task like cron spec. please see github.com/robfig/cron
func Every(ctx context.Context, interval time.Duration, fn func() error, errCh chan<- error) {
exit:
	for {
		select {
		case <-ctx.Done():
			break exit

		case <-time.After(interval):
			if err := fn(); err != nil {
				if errCh != nil {
					errCh <- err
				}
			}
		}
	}
}

// After run func after duration
func After(ctx context.Context, duration time.Duration, fn func()) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(duration):
		fn()
	}
}
