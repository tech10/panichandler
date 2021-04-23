package panic_handler

import (
	"testing"
	"time"
)

var called bool = false

var h HandlerFunc = func(i *Info) {
	called = true
}

func Test_panic_uncaught(t *testing.T) {
	go func() {
		defer Handle(h)
	}()
	time.Sleep(time.Millisecond * 10)
	if called {
		t.Fatal("A panic was never caught here, but the function to catch them has been called.")
	}
	t.Log("No panic was caught.")
}

func Test_panic_caught(t *testing.T) {
	called = false
	go func() {
		defer Handle(h)
		panic("testing")
	}()
	time.Sleep(time.Millisecond * 10)
	if !called {
		t.Fatal("A panic was not caught here, but it should have been.")
	}
	t.Log("A panic was caught.")
}

func Test_panic_value(t *testing.T) {
	value := "This is a test panic value."
	pstr := ""
	sstr := ""
	fstr := ""
	go func() {
		defer Handle(func(i *Info) {
			pstr = i.PanicString
			sstr = i.StackString
			fstr = i.String()
		})
		panic(value)
	}()
	time.Sleep(time.Millisecond * 10)
	if value != pstr {
		t.Fatal("The following should have been equal, and are not:\n\"" + value + "\", \"" + pstr + "\"")
	}
	t.Log("Panic caught with the following value: \"" + pstr)
	t.Log("The following stack trace was retrieved:\n\"" + sstr + "\"")
	t.Log("The complete formatted string of the panic value and runtime:\n\"" + fstr + "\"")
}
