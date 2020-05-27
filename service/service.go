package service

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/whitekid/go-utils/logging"
)

var log = logging.Named("service")

// Interface ...
type Interface interface {
	// Run service and block until stop
	Serve(ctx context.Context, args ...string) error
}

// AddrGetter for get listen address
type AddrGetter interface {
	GetAddr() *net.TCPAddr
}

// SetupSignal return context done when get system terminmaton signal
//
// returned context will be done when os.Interrupt, syscall.SIGTERM are invoked
// and call os.Exit() to exit program
func SetupSignal(ctx context.Context) context.Context {
	c := make(chan os.Signal, 2)
	defer close(c)

	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	termCtx, cancel := context.WithCancel(ctx)
	go func() {
		sig := <-c

		log.Debugf("got signal: %s", sig)

		cancel()
		os.Exit(1)
	}()

	return termCtx
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

// Multi 여러 sub service들을 돌릴 수 있는 서비스
type Multi struct {
	services []Interface
	ctx      context.Context
}

// NewMulti create Multi
func NewMulti(services ...Interface) *Multi {
	return &Multi{
		services: services,
	}
}

// Serve runs sub services
func (s *Multi) Serve(ctx context.Context, args ...string) error {
	if len(s.services) == 0 {
		return errors.New("No registered services")
	}

	errorC := make(chan error)
	defer close(errorC)

	// run sub service
	for _, service := range s.services {
		go func(service Interface) {
			if err := service.Serve(ctx); err != nil {
				errorC <- err
			}
		}(service)
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errorC:
		log.Errorf("Error %s", err)
		return err
	}
}
