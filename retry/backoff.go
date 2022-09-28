package retry

import (
	"context"
	"time"
)

type backoffer interface {
	// return current backoff durations
	CurrentBackoff() time.Duration

	// when next backoff comes data will be here
	NextBackoff() <-chan struct{}
}

func newZeroBackoff(ctx context.Context) backoffer { return newFixedBackoff(ctx, 0) }

func newFixedBackoff(ctx context.Context, interval time.Duration) backoffer {
	return newRatioBackoff(ctx, interval, interval, 0)
}

const defaultMaxBackoff = 60 * time.Second

func newRatioBackoff(ctx context.Context, initial time.Duration, maxBackoff time.Duration, ratio float64) backoffer {
	ch := make(chan struct{})

	go func() {
		<-ctx.Done()
		close(ch)
	}()

	if maxBackoff == 0 {
		maxBackoff = defaultMaxBackoff
	}

	return &backoffImpl{
		ch:         ch,
		initial:    initial,
		maxBackoff: maxBackoff,
		ratio:      ratio,
		backoff:    initial,
	}
}

type backoffImpl struct {
	ch         chan struct{}
	initial    time.Duration
	maxBackoff time.Duration
	ratio      float64

	backoff time.Duration
}

var _ backoffer = (*backoffImpl)(nil)

func (back *backoffImpl) CurrentBackoff() time.Duration { return back.backoff }

func (back *backoffImpl) NextBackoff() <-chan struct{} {
	go func() {
		<-time.After(back.next())

		back.ch <- struct{}{}
	}()

	return back.ch
}

// calculate next backoff
func (back *backoffImpl) next() time.Duration {
	current := back.backoff

	backoff := back.backoff + time.Duration(float64(back.backoff)*back.ratio)
	if backoff > back.maxBackoff {
		backoff = current
	} else {
		back.backoff = backoff
	}

	return current
}
