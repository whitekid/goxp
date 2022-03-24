package utils

import (
	"context"
	"sync"
	"time"

	"github.com/whitekid/go-utils/log"
)

// DoWithWorker iterate chan and run do() with n workers
func DoWithWorker(workers int, gen func(), do func(i int)) {
	var wg sync.WaitGroup
	wg.Add(workers)

	go gen()
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
