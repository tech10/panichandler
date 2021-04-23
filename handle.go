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

	// Bytes of the returned panic interface.
	PanicBytes []byte

	// The direct interface panic was provided when called, either by the Go runtime or by the user.
	PanicInterface interface{}

	// String of the returned panic interface.
	PanicString string

	// The stack trace as taken from debug.Stack()
	StackBytes []byte

	// The stack as taken by debug.Stack() converted to a string.
	StackString string
}

// Populate the Info struct with all values.
func (i *Info) populate(r interface{}, d []byte) {
	i.StackBytes = d
	i.StackString = fmt.Sprintf("%s", d)
	i.PanicInterface = r
	pstr := fmt.Sprintf("%s", r)
	i.PanicString = pstr
	i.PanicBytes = []byte(pstr)
}

// Handle panics. Call this in a defer statement, like this.
// panic_handler.Handle(HandlerFunc)
func Handle(c HandlerFunc) {
	r := recover()
	if r == nil {
		return
	}
	d := debug.Stack()
	i := &Info{}
	i.populate(r, d)
	if c != nil {
		c(i)
		return
	}
	fmt.Fprintln(os.Stderr, i.String())
	os.Exit(ExitCode)
}

// Returns a string formatted output of the panic and stack trace.
func (i *Info) String() string {
	return i.PanicString + "\n" + i.StackString
}

// Returns a byte formatted output of the panic and stack trace.
func (i *Info) Bytes() []byte {
	return []byte(i.PanicString + "\n" + i.StackString)
}

// Handle panics with this function.
type HandlerFunc func(*Info)
