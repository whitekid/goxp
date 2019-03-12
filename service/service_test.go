package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/go-utils/log"
)

type sampleService struct {
	started bool
}

func newSampleService() Interface {
	return &sampleService{}
}

func (s *sampleService) Serve(ctx context.Context, args ...string) error {
	s.started = true

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		log.Infof("Now: %s", time.Now().UTC().String())
		time.Sleep(time.Second)
	}
}

func TestSingle(t *testing.T) {
	svc := newSampleService()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	svc.Serve(ctx)

	require.True(t, svc.(*sampleService).started)
}

func TestMulti(t *testing.T) {
	services := []Interface{newSampleService(), newSampleService()}
	m := NewMulti(services...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.Nil(t, m.Serve(ctx))
	for _, svc := range services {
		require.True(t, svc.(*sampleService).started)
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
