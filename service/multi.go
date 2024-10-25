package service

import (
	"context"
	"errors"

	"github.com/whitekid/goxp/log"
)

type MultiService interface {
	Interface

	Append(Interface)
}

// multiServiceImpl 여러 sub service들을 돌릴 수 있는 서비스
type multiServiceImpl struct {
	services []Interface
}

// NewMulti create Multi
func NewMulti(services ...Interface) MultiService {
	return &multiServiceImpl{
		services: services,
	}
}

func (s *multiServiceImpl) Append(srv Interface) {
	s.services = append(s.services, srv)
}

// Serve runs sub services
func (s *multiServiceImpl) Serve(ctx context.Context) error {
	if len(s.services) == 0 {
		return errors.New("no registered services")
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
