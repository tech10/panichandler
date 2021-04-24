package panic_handler

import (
	"runtime/debug"
)

// Send the *Info struct to a channel rather than a function.
func HandleWithChan(c chan<- *Info) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	defer func() {
		recover()
	}()
	c <- i
}
