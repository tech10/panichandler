package panic_handler

import (
	"sync"
	"testing"
	"time"
)

var called bool = false

var l sync.Mutex

var h HandlerFunc = func(i *Info) {
	l.Lock()
	defer l.Unlock()
	called = true
}

func Test_panic_uncaught(t *testing.T) {
	go func() {
		defer Handle(h)
	}()
	time.Sleep(time.Millisecond * 10)
	l.Lock()
	if called {
		called = false
		l.Unlock()
		t.Fatal("A panic was never caught here, but the function to catch them has been called.")
	}
	l.Unlock()
	t.Log("No panic was caught.")
}

func Test_panic_caught(t *testing.T) {
	l.Lock()
	called = false
	l.Unlock()
	go func() {
		defer Handle(h)
		panic("testing")
	}()
	time.Sleep(time.Millisecond * 10)
	l.Lock()
	if !called {
		l.Unlock()
		t.Fatal("A panic was not caught here, but it should have been.")
	}
	called = false
	l.Unlock()
	t.Log("A panic was caught.")
}

func Test_panic_value(t *testing.T) {
	l.Lock()
	value := "This is a test panic value."
	pstr := ""
	sstr := ""
	fstr := ""
	l.Unlock()
	go func() {
		defer Handle(func(i *Info) {
			l.Lock()
			pstr = i.PanicString
			sstr = i.StackString
			fstr = i.String()
			l.Unlock()
		})
		panic(value)
	}()
	time.Sleep(time.Millisecond * 10)
	l.Lock()
	if value != pstr {
		l.Unlock()
		t.Fatal("The following should have been equal, and are not:\n\"" + value + "\", \"" + pstr + "\"")
	}
	l.Unlock()
	t.Log("Panic caught with the following value: \"" + pstr)
	t.Log("The following stack trace was retrieved:\n\"" + sstr + "\"")
	t.Log("The complete formatted string of the panic value and runtime:\n\"" + fstr + "\"")
}
