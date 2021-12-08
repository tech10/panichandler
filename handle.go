// Handle panics in a simple manner.
package panic_handler

import (
	"fmt"
	"os"
	"runtime/debug"
)

// ExitCode used if one isn't provided.
// This can be set with panic_handler.ExitCode = code.
var ExitCode int = 111

func caller(i *Info, c HandlerFunc, e int) bool {
	if c == nil {
		return false
	}
	defer nestedPanic(i, e)
	c(i)
	return true
}

// Catch a panic within the function designed to run upon receiving a panic.
// This will crash the program after printing out all stack traces.
func nestedPanic(i *Info, e int) {
	iN := newInfo(recover(), debug.Stack())
	if iN == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "WARNING!!!\nA panic within a panic catching function has been detected, this is a severe bug. Never fear, all stack traces are below.\nOriginally caught panic:\n%s\nPanic caused while catching original panic:\n%s\n", i.String(), iN.String())
	os.Exit(e)
}

// Handle panics. Call this in a defer statement, like this.
// panic_handler.Handle(panic_handler.HandlerFunc).
func Handle(c HandlerFunc) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if caller(i, c, ExitCode) {
		return
	}
	fmt.Fprintln(os.Stderr, i.String())
	os.Exit(ExitCode)
}
