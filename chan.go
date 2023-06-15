package goxp

import (
	"context"
	"sync"

	"github.com/whitekid/iter"
)

// IterChan iter chan and run fx(), until context is done
func IterChan[T any](ctx context.Context, ch <-chan T, fn func(T)) {
	it := iter.C(ch)
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		if IsContextDone(ctx) {
			break
		}

		fn(v)
	}
}

// CloseChan close chan when context is done
func CloseChan[T any](ctx context.Context, ch chan T) {
	<-ctx.Done()
	close(ch)
}

// FadeIn collect chan to single chan
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		ch := ch
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FadeOut distribute chan
func FadeOut[T any](ctx context.Context, ch <-chan T, size int) []<-chan T {
	chans := make([]chan T, 0, size)
	for i := 0; i < size; i++ {
		chans = append(chans, make(chan T))
	}

	go func() {
		i := 0
	exit:
		for v := range ch {
			select {
			case <-ctx.Done():
				break exit
			default:
				chans[i] <- v

				i++
				i %= size
			}
		}

		for i := 0; i < size; i++ {
			close(chans[i])
		}
	}()

	r := make([]<-chan T, size)
	for i := 0; i < size; i++ {
		r[i] = chans[i]
	}

	return r

}

// Async run func and returns with channel
func Async[T any](fn func() T) <-chan T {
	ch := make(chan T)
	go func() {
		ch <- fn()
		close(ch)
	}()
	return ch
}

// Async2 run func and returns with channel
func Async2[U1, U2 any](fn func() (U1, U2)) <-chan *Tuple2[U1, U2] {
	ch := make(chan *Tuple2[U1, U2])
	go func() {
		ch <- T2(fn())
		close(ch)
	}()
	return ch
}
