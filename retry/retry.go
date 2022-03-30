package retry

import (
	"context"
	"time"

	"github.com/whitekid/goex"
	"github.com/whitekid/goex/log"
)

// Retry return new default retrier
func Retry() Retrier {
	return &retrier{
		limit:          5,
		initialBackoff: time.Millisecond * 100,
		backoffRatio:   1.3,
	}
}

// Retrier is retrier interface
type Retrier interface {
	// total try limit
	Limit(limit int) Retrier

	// backoff settings when fails
	Backoff(initial time.Duration, backoffRatio float64) Retrier

	// run the func
	Do(ctx context.Context, fn func() error) error
}

type retrier struct {
	limit          int
	initialBackoff time.Duration
	backoffRatio   float64
}

func (r *retrier) Limit(limit int) Retrier {
	r.limit = limit
	return r
}

func (r *retrier) Backoff(initial time.Duration, ratio float64) Retrier {
	r.initialBackoff = initial
	r.backoffRatio = ratio
	return r
}

func (r *retrier) Do(ctx context.Context, fn func() error) (err error) {
	backoff := r.initialBackoff

	for i := 0; i < r.limit; i++ {
		if goex.IsContextDone(ctx) {
			return ctx.Err()
		}

		if err = fn(); err != nil {
			log.Infof("try %d failed, retry in %s", i+1, backoff)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * r.backoffRatio)
			continue
		}

		return nil
	}

	return err
}
