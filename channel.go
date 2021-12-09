package panichandler

import (
	"fmt"
	"os"
	"runtime/debug"
)

// Send the *Info struct to a channel rather than a function.
// Call it like this.
// panichandler.HandleWithChan(chan *panichandler.Info).
func HandleWithChan(c chan<- *Info) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if c == nil {
		fmt.Fprintf(os.Stderr, "WARNING!!!\nThe HandleWithChan function cannot have a nil channel.\nPanic reason and stack trace:\n%s\n", i.String())
		os.Exit(ExitCode)
	}
	channelSend(i, c, ExitCode)
}

func channelSend(i *Info, c chan<- *Info, e int) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r == nil {
			return
		}
		fmt.Fprintf(os.Stderr, "WARNING!!!\nYour program has sent the panic information to a closed channel. Panic information and stack trace:\n%s\n", i.String())
		os.Exit(e)
	}()
	c <- i
}
