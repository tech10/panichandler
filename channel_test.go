package panic_handler

import (
	"sync"
	"testing"
)

func Test_HandleWithChan(t *testing.T) {
	c := make(chan *Info)
	go func() {
		defer HandleWithChan(c)
		panic("Test sending a panic value over a channel.")
	}()
	i, open := <-c
	if !open {
		t.Fatal("The channel is closed, and should not have been.")
	}
	if i == nil {
		t.Fatal("The caught panic value of *Info is nil, and should not have been.")
	}
	t.Log("Retrieved value of panic:\n" + i.PanicString)
}

func Test_HandleWithChanClosed(t *testing.T) {
	var cl sync.Mutex
	cl.Lock()
	c := make(chan *Info)
	close(c)
	go func() {
		defer cl.Unlock()
		defer HandleWithChan(c)
		panic("Test sending a panic value over a closed channel.")
	}()
	cl.Lock()
	defer cl.Unlock()
	i, open := <-c
	if open {
		t.Fatal("The channel is open, and should not have been.")
	}
	if i != nil {
		t.Fatal("The caught panic value of *Info is not nil, and should have been. Full value and stack trace:\n", i.String())
	}
	t.Log("A panic was caught, but data failed to send over a closed channel. The program has silently recovered.")
}
