package panic_handler

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

// Handle panics in a function and cancel the provided context.CancelFunc.
// Call it like this.
// panic_handler.HandleWithContextCancel(context.CancelFunc, panic_handler.HandlerFunc)
func HandleWithContextCancel(cancel context.CancelFunc, c HandlerFunc) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	caller(i, c, ExitCode)
	if cancel == nil {
		fmt.Fprintf(os.Stderr, "Nil context CancelFunc provided, uncatchable panic.\n%s", i.String())
		os.Exit(ExitCode)
	}
	cancel()
}
