package service

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/whitekid/go-utils/logging"
)

var (
	log = logging.New()
)

// Service ...
type Service interface {
	// Run service and block until stop
	Serve(ctx context.Context, args ...string) error
}

// SetupSignal setup termination signal context
func SetupSignal(ctx context.Context) context.Context {
	c := make(chan os.Signal, 2)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	ret, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
		os.Exit(1)
	}()

	return ret
}

// Helper is helper for services
type Helper struct {
}

// Done return true if context.Done
func (s *Helper) Done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// MultiService 여러 sub service들을 돌릴 수 있는 서비스
type MultiService struct {
	services []Service
	ctx      context.Context
}

// Multi create MultiService
func Multi(services ...Service) *MultiService {
	return &MultiService{
		services: services,
	}
}

// Serve runs sub services
func (s *MultiService) Serve(ctx context.Context, args ...string) error {
	if len(s.services) == 0 {
		return errors.New("No registered services")
	}

	wg := sync.WaitGroup{}
	ctx1, cancel := context.WithCancel(ctx)

	var startError error

	// run sub service
	for _, service := range s.services {
		wg.Add(1)
		go func(service Service) {
			defer wg.Done()
			if err := service.Serve(ctx1); err != nil {
				log.Errorf("Error %s", err)
				startError = err

				cancel()
			}
		}(service)
	}

	// wait context canceling
	go func() {
		select {
		case <-ctx.Done():
			log.Debug("get cancel")
			cancel()
		}
	}()

	wg.Wait()

	return startError
}
