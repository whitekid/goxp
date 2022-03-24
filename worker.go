package utils

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/whitekid/go-utils/log"
)

// DoWithWorker iterate chan and run do() with n workers
// if works <=0 then worker set to runtime.NumCPU()
func DoWithWorker(workers int, gen func(), do func(i int)) {
	var wg sync.WaitGroup
	wg.Add(workers)

	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	if gen != nil {
		go gen()
	}

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			do(i)
		}(i)
	}

	wg.Wait()
}

type EveryFunc func() error

// Every execute fn() in every interval
func Every(ctx context.Context, interval time.Duration, fn EveryFunc) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				if err := fn(); err != nil {
					log.Errorf("task failed with %+v", err)
				}
			}
		}
	}()
}
