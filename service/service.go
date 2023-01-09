package service

import (
	"context"
	"errors"
	"net"

	"github.com/whitekid/goxp/log"
)

// Interface the service interface
type Interface interface {
	// Run service and block until stop
	Serve(ctx context.Context) error
}

// AddrGetter for get listen address
type AddrGetter interface {
	GetAddr() *net.TCPAddr
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
func (s *Multi) Serve(ctx context.Context) error {
	if len(s.services) == 0 {
		return errors.New("No registered services")
	}

	errorC := make(chan error)
	defer close(errorC)

	// run sub service
	for _, service := range s.services {
		service := service
		go func() {
			if err := service.Serve(ctx); err != nil {
				errorC <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errorC:
		log.Errorf("Error %s", err)
		return err
	}
}
