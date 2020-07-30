package suture

import (
	"context"
	"fmt"
	"testing"
)

const (
	JobLimit = 2
)

type IncrementorJob struct {
	current int
	next    chan int
}

func (i *IncrementorJob) Serve(ctx context.Context) error {
	for {
		select {
		case i.next <- i.current + 1:
			i.current++
			if i.current >= JobLimit {
				fmt.Println("Stopping the service")
				return ErrComplete
			}
		}
	}
}

func TestCompleteJob(t *testing.T) {
	supervisor := NewSimple("Supervisor")
	service := &IncrementorJob{0, make(chan int)}
	supervisor.Add(service)

	supervisor.ServeBackground()

	fmt.Println("Got:", <-service.next)
	fmt.Println("Got:", <-service.next)

	supervisor.Stop()

	// Output:
	// Got: 1
	// Got: 2
	// Stopping the service
}
