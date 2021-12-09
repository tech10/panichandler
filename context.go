package panichandler

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

// Handle panics in a function and cancel the provided context.CancelFunc.
// Call it like this.
// panichandler.HandleWithContextCancel(context.CancelFunc, panichandler.HandlerFunc).
func HandleWithContextCancel(cancel context.CancelFunc, c HandlerFunc) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if cancel == nil {
		fmt.Fprintf(os.Stderr, "WARNING!!!\nNil context CancelFunc provided, uncatchable panic.\n%s\n", i.String())
		os.Exit(ExitCode)
	}
	caller(i, c, ExitCode)
	cancel()
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
