package panic_handler

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

// Capture panics in something you can pass around the program in an easy to use variable.
type Capture struct {
	F        HandlerFunc        // Function to execute to handle panics, this is called first.
	T        Task               // Interface to run the DoPanicTask method on, this is called second.
	C        chan *Info         // Channel to pass panic information to, this is done third.
	CC       context.CancelFunc // Context cancelation function, this is called last.
	ExitCode int                // Status to exit with if a panic occurrs that crashes the program, and isn't caught by anything else.
}

func NewCapture() *Capture {
	return &Capture{
		ExitCode: ExitCode,
	}
}

func (c *Capture) Catch() {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if c.F == nil && c.T == nil && c.C == nil && c.CC == nil {
		fmt.Fprintf(os.Stderr, "Uninitialized Capture struct used, invalid operation.\n%s", i.String())
		os.Exit(c.ExitCode)
	}
	caller(i, c.F, c.ExitCode)
	taskRun(i, c.T, c.ExitCode)
	channelSend(i, c.C)
	if c.CC != nil {
		c.CC()
	}
}
