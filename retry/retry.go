package retry

import (
	"context"
	"errors"
	"time"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
)

// New return new default retrier
func New() Interface {
	return &retrierImpl{
		limit:          defaultLimit,
		initialBackoff: defaultInitialBackoff,
		backoffRatio:   defaultBackoffRatio,
	}
}

const (
	defaultLimit          = 5
	defaultInitialBackoff = time.Millisecond * 100
	defaultBackoffRatio   = 1.3
)

// Interface is retrier interface
type Interface interface {
	// total try limit
	Limit(limit int) Interface

	// backoff settings when fails
	Backoff(initial time.Duration, backoffRatio float64) Interface

	// run the func with retry
	// exit when ctx is done or returns ErrRetryStop
	Do(ctx context.Context, fn func() error) error
}

var ErrStop = errors.New("stop retry")

type retrierImpl struct {
	limit          int
	initialBackoff time.Duration
	backoffRatio   float64
}

var _ Interface = (*retrierImpl)(nil)

func (r *retrierImpl) Limit(limit int) Interface {
	r.limit = limit
	return r
}

func (r *retrierImpl) Backoff(initial time.Duration, ratio float64) Interface {
	r.initialBackoff = initial
	r.backoffRatio = ratio
	return r
}

func (r *retrierImpl) Do(ctx context.Context, fn func() error) (err error) {
	backoff := newRatioBackoff(ctx, r.initialBackoff, 0, r.backoffRatio)

	for i := 0; i < r.limit; i++ {
		if goxp.IsContextDone(ctx) {
			return ctx.Err()
		}

		if err = fn(); err == nil {
			return nil
		}

		if errors.Is(err, ErrStop) {
			return err
		}

		log.Infof("try %d failed, retry in %s: err=%v", i+1, backoff.CurrentBackoff(), err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-backoff.NextBackoff():
			continue
		}
	}

	return err
}
