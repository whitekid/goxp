package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
)

type simpleService struct {
	started bool
}

func newSampleService() Interface {
	return &simpleService{}
}

func (s *simpleService) Serve(ctx context.Context) error {
	s.started = true

	go goxp.Every(ctx, time.Second, func() error {
		if goxp.IsContextDone(ctx) {
			return nil
		}

		log.Infof("Now: %s", time.Now().UTC().Format(time.RFC3339))
		return nil
	}, nil)

	<-ctx.Done()
	return nil
}

func TestSingle(t *testing.T) {
	svc := newSampleService()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	svc.Serve(ctx)

	require.True(t, svc.(*simpleService).started)
}

func TestMulti(t *testing.T) {
	services := []Interface{newSampleService(), newSampleService()}
	m := NewMulti(services...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.Nil(t, m.Serve(ctx))
	for _, svc := range services {
		require.True(t, svc.(*simpleService).started)
	}
}

func TestCascadeCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	something := func(ctx context.Context) {
		log.Info("entering something")
		defer log.Info("exit somthing")

		<-ctx.Done()
	}

	for i := 0; i < 3; i++ {
		go something(ctx)
	}

	cancel()

	time.Sleep(time.Second)
}
