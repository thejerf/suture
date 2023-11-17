package suture

import (
	"context"
	"testing"
)

var _ Service = (*ServiceWithSecret)(nil)

type ServiceWithSecret struct {
	secret string
}

func (s *ServiceWithSecret) Serve(context.Context) error {
	return nil
}

func TestEventStopTimeout_String(t *testing.T) {
	tests := map[string]struct {
		Supervisor *Supervisor
		Service    Service
		want       string
	}{
		"service doesn't implement fmt.Stringer": {
			Supervisor: &Supervisor{Name: "test"},
			Service:    &ServiceWithSecret{secret: "MySecret"},
			want:       "test: Service '*suture.ServiceWithSecret' failed to terminate in a timely manner",
		},
		"service implements fmt.Stringer": {
			Supervisor: &Supervisor{Name: "test"},
			Service:    NewService("named"),
			want:       "test: Service 'named' failed to terminate in a timely manner",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := EventStopTimeout{
				Supervisor: tt.Supervisor,
				Service:    tt.Service,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("EventStopTimeout.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
