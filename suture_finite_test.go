package suture

import (
	"fmt"
	"testing"
)

const (
	JOB_LIMIT = 2
)

type IncrementorJob struct {
	current int
	next    chan int
	stop    chan bool
}

func (i *IncrementorJob) Stop() {
	fmt.Println("Stopping the service")
	i.stop <- true
}

func (i *IncrementorJob) Serve() {
	for {
		select {
		case i.next <- i.current + 1:
			i.current++
			if i.current >= JOB_LIMIT {
				return
			}
		case <-i.stop:
			// We sync here just to guarantee the output of "Stopping the service",
			// so this passes the test reliably.
			// Most services would simply "return" here.
			i.stop <- true
			return
		}
	}
}

func (i *IncrementorJob) Complete() bool {
	// fmt.Println("IncrementorJob exited as Complete()")
	return i.current >= JOB_LIMIT
}

func TestFiniteJob(t *testing.T) {
	supervisor := NewSimple("Supervisor")
	service := &IncrementorJob{0, make(chan int), make(chan bool)}
	supervisor.Add(service)

	supervisor.ServeBackground()

	fmt.Println("Got:", <-service.next)
	fmt.Println("Got:", <-service.next)

	<-service.stop

	fmt.Println("IncrementorJob exited as Complete()")

	supervisor.Stop()

	// Output:
	// Got: 1
	// Got: 2
	// Stopping the service
}
