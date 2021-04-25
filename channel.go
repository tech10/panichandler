package panic_handler

import (
	"fmt"
	"os"
	"runtime/debug"
)

// Send the *Info struct to a channel rather than a function.
// Call it like this.
// panic_handler.HandleWithChan(chan *panic_handler.Info)
func HandleWithChan(c chan<- *Info) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	if c == nil {
		fmt.Fprintf(os.Stderr, "WARNING!!!\nThe HandleWithChan function cannot have a nil channel.\nPanic reason and stack trace:\n%s\n", i.String())
		os.Exit(ExitCode)
	}
	channelSend(i, c)
}

func channelSend(i *Info, c chan<- *Info) {
	if c == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	c <- i
}
