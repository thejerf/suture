package suture

import (
	"fmt"
)

type (
	serviceID uint32

	// ServiceToken is an opaque identifier that can be used to terminate a
	// service that has been Add()ed to a Supervisor.
	ServiceToken struct {
		id uint64
	}

	/*Service is the interface that describes a service to a Supervisor.

	Serve Method

	The Serve method is called by a Supervisor to start the service.
	The service should execute within the goroutine that this is
	called in. If this function either returns or panics, the Supervisor
	will call it again.

	A Serve method SHOULD do as much cleanup of the state as possible,
	to prevent any corruption in the previous state from crashing the
	service again.

	Stop Method

	This method is used by the supervisor to stop the service. Calling this
	directly on a Service given to a Supervisor will simply result in the
	Service being restarted; use the Supervisor's .Remove(ServiceToken) method
	to stop a service. A supervisor will call .Stop() only once. Thus, it may
	be as destructive as it likes to get the service to stop.

	Once Stop has been called on a Service, the Service SHOULD NOT be
	reused in any other supervisor! Because of the impossibility of
	guaranteeing that the service has actually stopped in Go, you can't
	prove that you won't be starting two goroutines using the exact
	same memory to store state, causing completely unpredictable behavior.

	Stop should not return until the service has actually stopped.
	"Stopped" here is defined as "the service will stop servicing any
	further requests in the future". For instance, a common implementation
	is to receive a message on a dedicated "stop" channel and immediately
	returning. Once the stop command has been processed, the service is
	stopped.

	Another common Stop implementation is to forcibly close an open socket
	or other resource, which will cause detectable errors to manifest in the
	service code. Bear in mind that to perfectly correctly use this
	approach requires a bit more work to handle the chance of a Stop
	command coming in before the resource has been created.

	If a service does not Stop within the supervisor's timeout duration, a log
	entry will be made with a descriptive string to that effect. This does
	not guarantee that the service is hung; it may still get around to being
	properly stopped in the future. Until the service is fully stopped,
	both the service and the spawned goroutine trying to stop it will be
	"leaked".

	Stringer Interface

	It is not mandatory to implement the fmt.Stringer interface on your
	service, but if your Service does happen to implement that, the log
	messages that describe your service will use that when naming the
	service. Otherwise, you'll see the GoString of your service object,
	obtained via fmt.Sprintf("%#v", service).

	If you implement the Stringer interface, it must be callable from any
	goroutine, and must not be dependent on the state of your Service, i.e.,
	it must still successfully execute and return something even if your
	Service is completely crashed.

	*/
	Service interface {
		Serve()
		Stop()
	}

	listServices struct {
		c chan []Service
	}

	removeService struct {
		id serviceID
	}

	serviceFailed struct {
		id         serviceID
		err        interface{}
		stacktrace []byte
	}

	serviceEnded struct {
		id serviceID
	}

	// added by the Add() method
	addService struct {
		service  Service
		response chan serviceID
	}

	serviceTerminated struct {
		id serviceID
	}
)

// serviceName returns the name of a given Service
func serviceName(service Service) string {
	// Return service.String() if service is a fmt.Stringer
	if stringer, ok := service.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%#v", service)
}

func (ls listServices) isSupervisorMessage()  {}
func (rs removeService) isSupervisorMessage() {}
func (sf serviceFailed) isSupervisorMessage() {}

func (s serviceEnded) isSupervisorMessage() {}

func (as addService) isSupervisorMessage() {}

func (st serviceTerminated) isSupervisorMessage() {}
