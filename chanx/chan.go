package chanx

import "context"

// Iter iter chan and run fx(), until context is done or fn returns error
func Iter[T any](ctx context.Context, ch <-chan T, fn func(T) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-ch:
			if !ok {
				return nil
			}

			if err := fn(v); err != nil {
				return err
			}
		}
	}
}

// Close close chan when context is done
func Close[T any](ctx context.Context, ch chan T) {
	<-ctx.Done()
	close(ch)
}
