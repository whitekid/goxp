package utils

import "sync"

// DoWithWorker iteraate chan and run do() with n workers
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
