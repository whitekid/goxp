package fx

import "context"

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
