package services

import (
	"context"

	"golang.org/x/sync/errgroup"
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
		return ErrNoRegisteredService
	}

	eg, ctx := errgroup.WithContext(ctx)

	// run sub service
	for _, service := range s.services {
		eg.Go(func() error {
			if err := service.Serve(ctx); err != nil {
				return err
			}

			return nil
		})
	}

	return eg.Wait()
}
