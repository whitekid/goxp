package services

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/whitekid/goxp/errors"
	"github.com/whitekid/goxp/log"
)

type MultiService interface {
	Interface

	Append(Interface)
	Len() int
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

func (s *multiServiceImpl) Append(srv Interface) { s.services = append(s.services, srv) }
func (s *multiServiceImpl) Len() int             { return len(s.services) }

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
				if !errors.Is(err, context.DeadlineExceeded) {
					log.Errorf("error: %+v", err)
				}

				return err
			}

			return nil
		})
	}

	return eg.Wait()
}
