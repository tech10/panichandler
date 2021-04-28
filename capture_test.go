package panic_handler

import (
	"context"
	"sync"
	"testing"
)

type captT struct {
	f func()
}

func (c *captT) DoPanicTask(i *Info) {
	c.f()
}

func Test_Capture(t *testing.T) {
	var l sync.Mutex
	var wg sync.WaitGroup
	logstr := ""
	ctx, ctxc := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		<-ctx.Done()
		l.Lock()
		defer l.Unlock()
		logstr += "Context canceled."
		wg.Done()
	}()
	c := New()
	c.CC = ctxc
	c.C = make(chan *Info)
	wg.Add(1)
	go func() {
		<-c.C
		l.Lock()
		defer l.Unlock()
		logstr += "Channel received data.\n"
		wg.Done()
	}()
	wg.Add(1)
	c.F = func(i *Info) {
		l.Lock()
		defer l.Unlock()
		logstr += "Function called.\n"
		wg.Done()
	}
	wg.Add(1)
	c.T = &captT{
		f: func() {
			l.Lock()
			defer l.Unlock()
			logstr += "Task executed.\n"
			wg.Done()
		},
	}
	go func() {
		defer c.Catch()
		panic("testing capture struct")
	}()
	wg.Wait()
	l.Lock()
	defer l.Unlock()
	if logstr == "" {
		t.Fatal("Nothing was executed, but should have been.")
	}
	t.Log(logstr)
}

func Test_CaptureGetContext(t *testing.T) {
	c := New()
	ctx := c.GetContext()
	go func() {
		defer c.Catch()
		panic("Testing get context.")
	}()
	<-ctx.Done()
	t.Log("Context canceled successfully.")
}

func Test_CaptureCatchAndCancelContext(t *testing.T) {
	c := New()
	ctx := c.GetContext()
	go func() {
		defer c.CatchAndCancelContext()
	}()
	<-ctx.Done()
	t.Log("Context canceled successfully, and no panic was caused.")
}
