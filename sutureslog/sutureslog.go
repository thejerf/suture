/*
Package sutureslog implements a slog-based handler for suture events.

Note this is a separate module from suture to avoid forcing suture
itself to require a Go version that has log/slog. That dependency is
isolated to just this module. That means you must run a separate

	go get github.com/thejerf/suture/sutureslog

Passing this as a logger for a Supervisor looks like:

	// have some *slog.Logger called logger
	supervisor := suture.New(
	    "my supervisor name",
	    suture.Spec{
	        EventHook: sutureslog.Handler{Logger: logger}.MustHook(),
	    },
	)
*/
package sutureslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/thejerf/suture/v4"
)

const (
	// LogStopTimeout is logged for suture.EventBackoff.
	LogStopTimeout = "stopping service has timed out"

	// LogResume is logged for suture.EventResume
	LogResume = "supervisor has resumed restarting services"

	// LogBackoff is logged for suture.EventBackoff
	LogBackoff = "supervisor has entered backoff state"

	// LogServicePanic is logged for suture.EventServicePanic
	LogServicePanic = "service under supervision has panicked"

	// LogServiceTerminate is logged for
	// suture.EventServiceTerminate
	LogServiceTerminate = "service has terminated by returning"
)

var ErrNoHandlerSpecified = errors.New("no slog handler specified")

// A Handler handles incoming events and manages logging them to slog.
//
// Once a Supervisor starts running with a Hook from this handler, the
// values must not be changed.
//
// If the slog.Leveler is left nil, a default level will be used. A
// minimal definition includes only a Handler.
//
// For EventStopTimeout and EventServicePanic, this default level will
// be slog.LevelError.
//
// For EventServiceTerminate, the default level is slog.LevelInfo.
//
// For the rest it will be slog.LevelDebug.
//
// For Services passed as logging values, if the service implements
// slog.LogValuer, that will be passed to slog. Otherwise, if it
// implements fmt.Stringer, it will be stringified using
// that. Otherwise, it will be passed as an Any value.
//
// The log messages that will be emitted by each event correspond to
// the constants defined by this package; this can be used by
// something like [github.com/thejerf/slogassert] to validate proper
// functioning of your code. See sutureslog_test.go in this package
// for examples.
type Handler struct {
	StopTimeoutLevel      slog.Leveler
	ServicePanicLevel     slog.Leveler
	ServiceTerminateLevel slog.Leveler
	BackoffLevel          slog.Leveler
	ResumeLevel           slog.Leveler

	Logger *slog.Logger
}

func valFromService(s suture.Service) slog.Value {
	logValuer, isLogValuer := s.(slog.LogValuer)
	if isLogValuer {
		return logValuer.LogValue()
	}

	stringer, isStringer := s.(fmt.Stringer)
	if isStringer {
		return slog.StringValue(stringer.String())
	}

	return slog.AnyValue(s)
}

func defaultLevel(l1 slog.Leveler, l2 slog.Leveler) slog.Leveler {
	if l1 == nil {
		return l2
	}
	return l1
}

// Hook returns an event hook suitable for passing into a suture.Spec.
//
// The only possible error is ErrNoHandlerSpecified. If you statically
// guarantee that can't happen you can elide this error.
func (h *Handler) Hook() (func(suture.Event), error) {
	if h.Logger == nil {
		return nil, ErrNoHandlerSpecified
	}

	return func(event suture.Event) {
		switch evt := event.(type) {
		case suture.EventBackoff:
			h.Logger.LogAttrs(
				context.Background(),
				defaultLevel(h.BackoffLevel, slog.LevelDebug).Level(),
				LogBackoff,
				slog.String("supervisor", evt.Supervisor.Name))

		case suture.EventResume:
			h.Logger.LogAttrs(
				context.Background(),
				defaultLevel(h.ResumeLevel, slog.LevelDebug).Level(),
				LogResume,
				slog.String("supervisor", evt.Supervisor.Name))

		case suture.EventServicePanic:
			h.Logger.LogAttrs(
				context.Background(),
				defaultLevel(h.ServicePanicLevel, slog.LevelError).Level(),
				LogServicePanic,
				slog.String("supervisor", evt.Supervisor.Name),
				slog.Any("service", valFromService(evt.Service)),
				slog.String("service_name", evt.ServiceName),
				slog.Float64("current_failures", evt.CurrentFailures),
				slog.Float64("failure_threshold", evt.FailureThreshold),
				slog.Bool("restarting", evt.Restarting),
				slog.String("panic_msg", evt.PanicMsg),
				slog.String("stacktrace", evt.Stacktrace),
			)

		case suture.EventServiceTerminate:
			h.Logger.LogAttrs(
				context.Background(),
				defaultLevel(h.ServiceTerminateLevel, slog.LevelInfo).Level(),
				LogServiceTerminate,
				slog.String("supervisor", evt.Supervisor.Name),
				slog.Any("service", valFromService(evt.Service)),
				slog.String("service_name", evt.ServiceName),
				slog.Float64("current_failures", evt.CurrentFailures),
				slog.Float64("failure_threshold", evt.FailureThreshold),
				slog.Bool("restarting", evt.Restarting),
				slog.Any("error", evt.Err),
			)

		case suture.EventStopTimeout:
			h.Logger.LogAttrs(
				context.Background(),
				defaultLevel(h.StopTimeoutLevel, slog.LevelError).Level(),
				LogStopTimeout,
				slog.String("supervisor", evt.Supervisor.Name),
				slog.Any("service", valFromService(evt.Service)),
				slog.String("service_name", evt.ServiceName),
			)

		default:
			// while the Event interface is not closed,
			// suture should still NEVER emit an event
			// that isn't in the set of events it
			// understands; it is effectively a sum type.
			panic(fmt.Sprintf("unknown event passed to handler, file an issue at github.com/thejerf/suture that event type %T is not handled in sutureslog", event))
		}
	}, nil
}

// MustHook will return a valid event hook function or panic.
func (h *Handler) MustHook() func(suture.Event) {
	hook, err := h.Hook()
	if err != nil {
		panic("panic in getting sutureslog hook: no slog handler given")
	}
	return hook
}
