package chanx

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// FadeIn collect chan to single chan
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)

	eg, ctx := errgroup.WithContext(ctx)
	for _, ch := range chans {
		eg.Go(func() error {
			return Iter(ctx, ch, func(ctx context.Context, v T) error {
				out <- v
				return nil
			})
		})
	}

	go func() {
		defer close(out)
		eg.Wait()
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
