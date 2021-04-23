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
	r := recover()
	if r == nil {
		return
	}
	i := newInfo(r, debug.Stack())
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
	c(i)
	return true
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
