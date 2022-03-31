package retry

import (
	"context"
	"time"
)

type Backoff interface {
	// return current backoff durations
	CurrentBackoff() time.Duration

	// when next backoff comes data will be here
	NextBackoff() <-chan *time.Time
}

func NewZeroBackoff(ctx context.Context) Backoff { return NewFixedBackoff(ctx, 0) }

func NewFixedBackoff(ctx context.Context, interval time.Duration) Backoff {
	return NewRatioBackoff(ctx, interval, interval, 0)
}

const MaxBackoff = 60 * time.Second

func NewRatioBackoff(ctx context.Context, initial time.Duration, maxBackoff time.Duration, ratio float64) Backoff {
	ch := make(chan *time.Time)

	go func() {
		<-ctx.Done()
		close(ch)
	}()

	if maxBackoff == 0 {
		maxBackoff = MaxBackoff
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
	ch         chan *time.Time
	initial    time.Duration
	maxBackoff time.Duration
	ratio      float64

	backoff time.Duration
}

func (back *backoffImpl) CurrentBackoff() time.Duration { return back.backoff }

func (back *backoffImpl) NextBackoff() <-chan *time.Time {
	go func() {
		<-time.After(back.next())

		now := time.Now()
		back.ch <- &now
	}()

	return back.ch
}

func (back *backoffImpl) next() time.Duration {
	current := back.backoff

	// calculate next backoff
	backoff := back.backoff + time.Duration(float64(back.backoff)*back.ratio)
	if backoff > back.maxBackoff {
		backoff = current
	} else {
		back.backoff = backoff
	}

	return current
}
