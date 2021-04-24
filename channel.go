package panic_handler

import (
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
