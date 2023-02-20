package suture

import "fmt"

type Incrementor struct {
	current int
	next    chan int
	stop    chan struct{}
}

func (i *Incrementor) Stop() {
	fmt.Println("Stopping the service")
	close(i.stop)
}

func (i *Incrementor) Serve() {
	for {
		select {
		case i.next <- i.current:
			i.current++
		case <-i.stop:
			return
		}
	}
}

func ExampleNew_simple() {
	supervisor := NewSimple("Supervisor")
	service := &Incrementor{0, make(chan int), make(chan struct{})}
	supervisor.Add(service)

	supervisor.ServeBackground()

	fmt.Println("Got:", <-service.next)
	fmt.Println("Got:", <-service.next)
	supervisor.Stop()

	// We sync here just to guarantee the output of "Stopping the service"
	<-service.stop

	// Output:
	// Got: 0
	// Got: 1
	// Stopping the service
}
