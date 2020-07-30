package suture

import (
	"context"
)

type DeprecatedService interface {
	Serve()
	Stop()
}

func AsNewService(service DeprecatedService) Service {
	return &serviceShim{service: service}
}

type serviceShim struct {
	service DeprecatedService
}

func (s *serviceShim) Serve(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.service.Serve()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			s.service.Stop()
		case <-done:
			return ctx.Err()
		}
	}
}
