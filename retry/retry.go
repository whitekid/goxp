package retry

import (
	"context"
	"time"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
)

// Retry return new default retrier
func Retry() Interface {
	return &retrier{
		limit:          5,
		initialBackoff: time.Millisecond * 100,
		backoffRatio:   1.3,
	}
}

// Interface is retrier interface
type Interface interface {
	// total try limit
	Limit(limit int) Interface

	// backoff settings when fails
	Backoff(initial time.Duration, backoffRatio float64) Interface

	// run the func
	Do(ctx context.Context, fn func() error) error
}

// Depreciated: please use Interface
type Retrier = Interface

type retrier struct {
	limit          int
	initialBackoff time.Duration
	backoffRatio   float64
}

func (r *retrier) Limit(limit int) Interface {
	r.limit = limit
	return r
}

func (r *retrier) Backoff(initial time.Duration, ratio float64) Interface {
	r.initialBackoff = initial
	r.backoffRatio = ratio
	return r
}

func (r *retrier) Do(ctx context.Context, fn func() error) (err error) {
	backoff := NewRatioBackoff(ctx, r.initialBackoff, 0, r.backoffRatio)

	for i := 0; i < r.limit; i++ {
		if goxp.IsContextDone(ctx) {
			return ctx.Err()
		}

		if err = fn(); err == nil {
			return nil
		}

		log.Infof("try %d failed, retry in %s", i+1, backoff.CurrentBackoff())

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-backoff.NextBackoff():
			continue
		}
	}

	return err
}
