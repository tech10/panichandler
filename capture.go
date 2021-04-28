package panic_handler

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

// Capture panics in something you can pass around the program in an easy to use struct.
// The execution order for catching panics is the following:
// Your own defined function, a task, a channel, and a context.CancelFunc.
// Any of these values can be omitted, but all of them can't be omitted.
// If you try and catch panics without filling at least one value,
// and a panic is recovered from,
// your program will have the panic passed along to standard error,
// and will exit immediately.
type Capture struct {
	F        HandlerFunc        // Function to execute, handling panics, this is called first.
	T        Task               // Interface to run the DoPanicTask method on, this is called second.
	C        chan *Info         // Channel to pass panic information to, this is done third.
	CC       context.CancelFunc // Context cancelation function, this is called last.
	ExitCode int                // Status to exit with if a panic occurrs that crashes the program, and isn't caught by anything else.
}

// Create a new instance of the Capture struct, uninitialized.
// The only variable set here is its ExitCode,
// which is set to the default ExitCode variable.
func New() *Capture {
	return &Capture{
		ExitCode: ExitCode,
	}
}

func (c *Capture) catcher(r interface{}) {
	i := newInfo(r, debug.Stack())
	if i == nil {
		return
	}
	if c.F == nil && c.T == nil && c.C == nil && c.CC == nil {
		fmt.Fprintf(os.Stderr, "Uninitialized Capture struct used, invalid operation.\n%s\n", i.String())
		os.Exit(c.ExitCode)
	}
	if c.F != nil {
		caller(i, c.F, c.ExitCode)
	}
	if c.T != nil {
		taskRun(i, c.T, c.ExitCode)
	}
	if c.C != nil {
		channelSend(i, c.C, c.ExitCode)
	}
	if c.CC != nil {
		c.CC()
	}
}

// Catch panics, call this in a defer statement.
// Only do this if you have initialized the Capture struct
// with a function, Task interface,
// channel, or context.CancelFunc
func (c *Capture) Catch() {
	c.catcher(recover())
}

// Get a context that will be returned to you, and canceled upon a panic.
// If you have already set a context cancelation function, this will override it.
// The context returned should be used on anything you want to cancel,
// and will not be derived from any parent contexts.
// This is a context that is designed to be canceled
// upon a panic, though you could call the Capture.CC function yourself.
// You may need to do this if a panic is not caught.
func (c *Capture) GetContext() context.Context {
	var ctx context.Context
	ctx, c.CC = getContext()
	return ctx
}

// Always cancel a context, if it is available.
func (c *Capture) CatchAndCancelContext() {
	c.catcher(recover())
	if c.CC != nil {
		c.CC()
	}
}
