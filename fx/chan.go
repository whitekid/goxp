package fx

import (
	"context"
	"sync"
)

// IterChan iter chan and run fx(), until context is done
func IterChan[T any](ctx context.Context, ch <-chan T, fx func(T)) {
exit:
	for {
		select {
		case <-ctx.Done():
			break exit

		case v, ok := <-ch:
			if !ok {
				break exit
			}

			fx(v)
		}
	}
}

// CloseChan close chan when context is done
func CloseChan[T any](ctx context.Context, ch chan T) {
	<-ctx.Done()
	close(ch)
}

func FanIn[T any](ctx context.Context, chans ...chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		ch := ch
		go func() {
			IterChan(ctx, ch, func(v T) { out <- v })
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Async[T any](f func() T) <-chan T {
	ch := make(chan T)
	go func() { ch <- f(); close(ch) }()
	return ch
}

func Async2[U1, U2 any](f func() (U1, U2)) <-chan Tuple2[U1, U2] {
	ch := make(chan Tuple2[U1, U2])
	go func() { ch <- T2(f()); close(ch) }()
	return ch
}
