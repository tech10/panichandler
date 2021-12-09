package panichandler

import (
	"sync"
	"testing"
)

var called = false

var l sync.Mutex

var h HandlerFunc = func(i *Info) {
	l.Lock()
	defer l.Unlock()
	called = true
}

func Test_panic_uncaught(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Handle(h)
	}()
	wg.Wait()
	l.Lock()
	defer l.Unlock()
	if called {
		called = false
		t.Fatal("A panic was never caught here, but the function to catch them has been called.")
	}
	t.Log("No panic was caught.")
}

func Test_panic_caught(t *testing.T) {
	var wg sync.WaitGroup
	l.Lock()
	called = false
	l.Unlock()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Handle(h)
		panic("testing")
	}()
	wg.Wait()
	l.Lock()
	defer l.Unlock()
	if !called {
		t.Fatal("A panic was not caught here, but it should have been.")
	}
	called = false
	t.Log("A panic was caught.")
}

func Test_panic_value(t *testing.T) {
	var wg sync.WaitGroup
	l.Lock()
	value := "This is a test panic value."
	pstr := ""
	fstr := ""
	l.Unlock()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Handle(func(i *Info) {
			l.Lock()
			pstr = i.PanicString
			fstr = i.String()
			l.Unlock()
		})
		panic(value)
	}()
	wg.Wait()
	l.Lock()
	defer l.Unlock()
	if value != pstr {
		t.Fatal("The following should have been equal, and are not:\n\"" + value + "\", \"" + pstr + "\"")
	}
	t.Log("Panic caught with the following value: \"" + pstr)
	t.Log("The complete formatted string of the panic value and stack trace:\n" + fstr)
}
