package sutureslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/thejerf/slogassert"
	"github.com/thejerf/suture/v4"
)

// This provides a service that is a slog.LogValuer.
type ServiceSlogValue struct{}

func (ssv ServiceSlogValue) Serve(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (ssv ServiceSlogValue) LogValue() slog.Value {
	return slog.StringValue("ServiceSlogValue")
}

// This provides a service that is a fmt.Stringer.
type ServiceStringer struct{}

func (ss ServiceStringer) Serve(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (ss ServiceStringer) String() string {
	return "ServiceStringer"
}

// This provides a service that has neither special case.
type NormalService struct {
	Value string
}

func (ns NormalService) Serve(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func TestSutureSlog(t *testing.T) {
	handler := slogassert.New(t, slog.LevelDebug, nil)
	logger := slog.New(handler)

	sutureHandler := &Handler{Logger: logger}
	hook := sutureHandler.MustHook()

	supervisorName := "testsupervisor"

	// this is just to have a valid reference, we never run it. We
	// don't need to test that supervisors actually emit these
	// events, suture covers that, we just need to test that they
	// are handled in the documented manner.
	supervisor := suture.NewSimple(supervisorName)

	hook(suture.EventBackoff{
		Supervisor:     supervisor,
		SupervisorName: supervisorName,
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogBackoff,
		Level:         slog.LevelDebug,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor": supervisorName,
		},
	})
	handler.AssertEmpty()

	hook(suture.EventResume{
		Supervisor:     supervisor,
		SupervisorName: supervisorName,
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogResume,
		Level:         slog.LevelDebug,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor": supervisorName,
		},
	})
	handler.AssertEmpty()

	// Test a normal service for panic
	hook(suture.EventServicePanic{
		Supervisor:       supervisor,
		SupervisorName:   supervisorName,
		Service:          NormalService{"Hello"},
		ServiceName:      "normal_service",
		CurrentFailures:  1.0,
		FailureThreshold: 5.0,
		Restarting:       false,
		PanicMsg:         "oh no",
		Stacktrace:       "A Stack Trace",
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogServicePanic,
		Level:         slog.LevelError,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor":        supervisorName,
			"service":           NormalService{"Hello"},
			"service_name":      "normal_service",
			"current_failures":  float64(1.0),
			"failure_threshold": float64(5.0),
			"restarting":        false,
			"panic_msg":         "oh no",
			"stacktrace":        "A Stack Trace",
		},
	})

	// Test a fmt.Stringer service
	hook(suture.EventServicePanic{
		Supervisor:       supervisor,
		SupervisorName:   supervisorName,
		Service:          ServiceStringer{},
		ServiceName:      "normal_service",
		CurrentFailures:  1.0,
		FailureThreshold: 5.0,
		Restarting:       false,
		PanicMsg:         "oh no",
		Stacktrace:       "A Stack Trace",
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogServicePanic,
		Level:         slog.LevelError,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor":        supervisorName,
			"service":           "ServiceStringer",
			"service_name":      "normal_service",
			"current_failures":  float64(1.0),
			"failure_threshold": float64(5.0),
			"restarting":        false,
			"panic_msg":         "oh no",
			"stacktrace":        "A Stack Trace",
		},
	})

	// Test a slog.LogValuer service
	hook(suture.EventServicePanic{
		Supervisor:       supervisor,
		SupervisorName:   supervisorName,
		Service:          ServiceSlogValue{},
		ServiceName:      "normal_service",
		CurrentFailures:  1.0,
		FailureThreshold: 5.0,
		Restarting:       false,
		PanicMsg:         "oh no",
		Stacktrace:       "A Stack Trace",
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogServicePanic,
		Level:         slog.LevelError,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor":        supervisorName,
			"service":           "ServiceSlogValue",
			"service_name":      "normal_service",
			"current_failures":  float64(1.0),
			"failure_threshold": float64(5.0),
			"restarting":        false,
			"panic_msg":         "oh no",
			"stacktrace":        "A Stack Trace",
		},
	})

	hook(suture.EventServiceTerminate{
		Supervisor:       supervisor,
		SupervisorName:   supervisorName,
		Service:          ServiceSlogValue{},
		ServiceName:      "normal_service",
		CurrentFailures:  1.0,
		FailureThreshold: 5.0,
		Restarting:       true,
		Err:              errors.New("error"),
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogServiceTerminate,
		Level:         slog.LevelInfo,
		AllAttrsMatch: false,
		Attrs: map[string]any{
			"supervisor":        supervisorName,
			"service":           "ServiceSlogValue",
			"service_name":      "normal_service",
			"current_failures":  float64(1.0),
			"failure_threshold": float64(5.0),
			"restarting":        false,
			"error":             errors.New("error"),
		},
	})

	hook(suture.EventStopTimeout{
		Supervisor:     supervisor,
		SupervisorName: supervisorName,
		Service:        ServiceSlogValue{},
		ServiceName:    "normal_service",
	})
	handler.AssertPrecise(slogassert.LogMessageMatch{
		Message:       LogStopTimeout,
		Level:         slog.LevelError,
		AllAttrsMatch: true,
		Attrs: map[string]any{
			"supervisor":   supervisorName,
			"service":      "ServiceSlogValue",
			"service_name": "normal_service",
		},
	})
}

func TestMissingLogger(t *testing.T) {
	_, err := (&Handler{}).Hook()
	if err != ErrNoHandlerSpecified {
		fmt.Println(err)
		t.Fatal("missing logger not detected")
	}
}

func TestFailedMustMissingLogger(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("did not panic on missing logger")
		}
	}()

	(&Handler{}).MustHook()
}
