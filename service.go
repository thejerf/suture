package suture

import (
	"context"
	"errors"
)

/*
Service is the interface that describes a service to a Supervisor.

Serve Method

The Serve method is called by a Supervisor to start the service.
The service should execute within the goroutine that this is
called in, that is, it should not spawn a "worker" goroutine.
If this function either returns error or panics, the Supervisor
will call it again, unless it returns special error sentinel values

A Serve method SHOULD do as much cleanup of the state as possible,
to prevent any corruption in the previous state from crashing the
service again.

The Supervisor will pass a ctx down to the service. That context
can not be independently used to cancel the service, as the supervisor
will not notice the context is cancelled. Use the Remove method on the
Supervisor to remove it properly, or have the service return
ErrDoNotRestart.

The error returned by the service, if any, will be part of the log
message generated for it. There are two distinguished errors a
Service can return: ErrDoNotRestart indicates that the service should
not be restarted and removed from the supervisor entirely.

ErrTerminateTree indicates that the entire supervision tree should be
restarted.

In Go 1.13 and greater, this is checked via errors.Is, so the error
can be further wrapped with whatever additional info you like. Prior
to Go 1.13, it will be checked via directly equality check, so the
distinguished errors can not be wrapped.

Once the service has been instructed to stop, the Service SHOULD NOT be
reused in any other supervisor! Because of the impossibility of
guaranteeing that the service has actually stopped in Go, you can't
prove that you won't be starting two goroutines using the exact
same memory to store state, causing completely unpredictable behavior.

Serve should not return until the service has actually stopped.
"Stopped" here is defined as "the service will stop servicing any
further requests in the future". Any mandatory cleanup related to
the Service should also have been performed.

If a service does not stop within the supervisor's timeout duration, a log
entry will be made with a descriptive string to that effect. This does
not guarantee that the service is hung; it may still get around to being
properly stopped in the future. Until the service is fully stopped,
both the service and the spawned goroutine trying to stop it will be
"leaked".

Stringer Interface

When a Service is added to a Supervisor, the Supervisor will create a
string representation of that service used for logging.

If you implement the fmt.Stringer interface, that will be used.

If you do not implement the fmt.Stringer interface, a default
fmt.Sprintf("%#v") will be used.

*/
type Service interface {
	Serve(ctx context.Context) error
}

// ErrDoNotRestart can be returned by a service to voluntarily not
// be restarted.
var ErrDoNotRestart = errors.New("service should not be restarted")

// ErrTerminateTree can can be returned by a service to terminate the
// entire supervision tree above it as well.
var ErrTerminateTree = errors.New("tree should be terminated")
