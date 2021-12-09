package panichandler

import (
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
