package suture

import (
	"testing"
)

type testFooBar struct {
	doneC chan struct{}
}

func (tfb *testFooBar) Foo() {
	<-tfb.doneC
}

func (tfb *testFooBar) Bar() {
	close(tfb.doneC)
}

func TestGenerator(t *testing.T) {
	fb := &testFooBar{
		doneC: make(chan struct{}),
	}

	service := GenerateService("test service", fb.Foo, fb.Bar, false)

	supervisor := NewSimple("test supervisor")
	token := supervisor.Add(service)

	go supervisor.Serve()

	err := supervisor.Remove(token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
