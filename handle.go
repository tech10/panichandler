// Handle panics in a simple manner.
package panic_handler

import (
	"fmt"
	"os"
	"runtime/debug"
)

// Exit code used if one isn't provided.
// This can be set with panic_handler.ExitCode = code
var ExitCode int = 111

// Contains all information about a panic, formatted in various ways.
type Info struct {
	PanicBytes     []byte      // Bytes of the returned panic interface.
	PanicInterface interface{} // The direct interface panic was provided when called, either by the Go runtime or by the user.
	PanicString    string      // String of the returned panic interface.
	StackBytes     []byte      // The stack trace as taken from debug.Stack()
	StackString    string      // The stack as taken by debug.Stack() converted to a string.
}

// Return the Info struct with all values.
func newInfo(r interface{}, d []byte) *Info {
	if r == nil {
		return nil
	}
	i := &Info{}
	i.StackBytes = d
	i.StackString = fmt.Sprintf("%s", d)
	i.PanicInterface = r
	pstr := fmt.Sprintf("%s", r)
	i.PanicString = pstr
	i.PanicBytes = []byte(pstr)
	return i
}

// Handle panics. Call this in a defer statement, like this.
// panic_handler.Handle(HandlerFunc)
func Handle(c HandlerFunc) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if caller(i, c) {
		return
	}
	fmt.Fprintln(os.Stderr, i.String())
	os.Exit(ExitCode)
}

func caller(i *Info, c HandlerFunc) bool {
	if c == nil {
		return false
	}
	defer nestedPanic(i)
	c(i)
	return true
}

// Catch a panic within the function designed to run upon receiving a panic.
// This will crash the program after printing out all stack traces.
func nestedPanic(i *Info) {
	i_n := newInfo(recover(), debug.Stack())
	if i_n == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "WARNING!!!\nA panic within a panic catching function has been detected, this is a severe bug. Never fear, all stack traces are below.\nOriginally caught panic:\n%s\nPanic caused while catching original panic:\n%s\n", i.String(), i_n.String())
	os.Exit(ExitCode)
}

// Returns a string formatted output of the panic and stack trace.
func (i *Info) String() string {
	return i.PanicString + "\n" + i.StackString
}

// Returns a byte formatted output of the panic and stack trace.
func (i *Info) Bytes() []byte {
	return []byte(i.String())
}

// Handle panics with this function.
type HandlerFunc func(*Info)
