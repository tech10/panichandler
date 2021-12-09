package panichandler

import (
	"fmt"
	"strings"
)

// Handle panics with this function.
type HandlerFunc func(*Info)

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
	i.StackString = string(d)
	i.PanicInterface = r
	pstr := fmt.Sprintf("%s", r)
	i.PanicString = pstr
	i.PanicBytes = []byte(pstr)
	return i
}

// Returns a string formatted output of the panic and stack trace.
func (i *Info) String() string {
	return i.PanicString + "\n" + strings.TrimSpace(i.StackString)
}

// Returns a byte formatted output of the panic and stack trace.
func (i *Info) Bytes() []byte {
	return []byte(i.String())
}
