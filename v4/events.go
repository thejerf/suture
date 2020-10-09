package suture

import (
	"fmt"
)

type (
	EventType int

	Event interface {
		fmt.Stringer
		Type() EventType
	}

	EventHook func(Event)
)

const (
	EventTypeStopTimeout EventType = iota
	EventTypeServicePanic
	EventTypeServiceTerminate
	EventTypeBackoff
	EventTypeResume
)

type ServiceFailureEvent struct {
	Supervisor       string
	Service          string
	CurrentFailures  float64
	FailureThreshold float64
	Restarting       bool
}

type StopTimeoutEvent struct {
	Supervisor string
	Service    string
}

func (e StopTimeoutEvent) Type() EventType {
	return EventTypeStopTimeout
}

func (e StopTimeoutEvent) String() string {
	return fmt.Sprintf(
		"%s: Service %s failed to terminate in a timely manner",
		e.Supervisor,
		e.Service,
	)
}

type ServicePanicEvent struct {
	Supervisor       string
	Service          string
	CurrentFailures  float64
	FailureThreshold float64
	Restarting       bool
	PanicMsg         string
	Stacktrace       string
}

func (e ServiceFailureEvent) String() string {
	return fmt.Sprintf(
		"%s: Failed service '%s' (%f failures of %f), restarting: %#v",
		e.Supervisor,
		e.Service,
		e.CurrentFailures,
		e.FailureThreshold,
		e.Restarting,
	)
}

func (e ServicePanicEvent) Type() EventType {
	return EventTypeServicePanic
}

func (e ServicePanicEvent) String() string {
	return fmt.Sprintf(
		"%s, panic: %s, stacktrace: %s",
		serviceFailureString(e.Supervisor, e.Service, e.CurrentFailures, e.FailureThreshold, e.Restarting),
		e.PanicMsg,
		string(e.Stacktrace),
	)
}

type ServiceTerminateEvent struct {
	Supervisor       string
	Service          string
	CurrentFailures  float64
	FailureThreshold float64
	Restarting       bool
	Err              error
}

func (e ServiceTerminateEvent) Type() EventType {
	return EventTypeServiceTerminate
}

func (e ServiceTerminateEvent) String() string {
	return fmt.Sprintf(
		"%s, error: %s",
		serviceFailureString(e.Supervisor, e.Service, e.CurrentFailures, e.FailureThreshold, e.Restarting),
		e.Err)
}

func serviceFailureString(supervisor, service string, currentFailures, failureThreshold float64, restarting bool) string {
	return fmt.Sprintf(
		"%s: Failed service '%s' (%f failures of %f), restarting: %#v",
		supervisor,
		service,
		currentFailures,
		failureThreshold,
		restarting,
	)
}

type BackoffEvent struct {
	Supervisor string
}

func (e BackoffEvent) Type() EventType {
	return EventTypeBackoff
}

func (e BackoffEvent) String() string {
	return fmt.Sprintf("%s: Entering the backoff state.", e.Supervisor)
}

type ResumeEvent struct {
	Supervisor string
}

func (e ResumeEvent) Type() EventType {
	return EventTypeResume
}

func (e ResumeEvent) String() string {
	return fmt.Sprintf("%s: Exiting backoff state.", e.Supervisor)
}
