package panicHandler

import (
	"context"
	"sync"
	"testing"
)

func Test_HandleWithContextCancel(t *testing.T) {
	var wg sync.WaitGroup
	mctx, mcancel := context.WithCancel(context.Background())
	c1, cancel1 := context.WithCancel(mctx)
	c2, cancel2 := context.WithCancel(c1)
	wg.Add(2)
	go func() {
		<-c1.Done()
		wg.Done()
	}()
	go func() {
		<-c2.Done()
		wg.Done()
	}()
	go func() {
		defer HandleWithContextCancel(mcancel, nil)
		panic("Testing context cancelation.")
	}()
	wg.Wait()
	t.Log("Context cancelation successful.")
	mcancel()
	cancel1()
	cancel2()
}
