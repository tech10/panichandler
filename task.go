package panic_handler

import (
	"runtime/debug"
)

// Interface for defining your own handler, perhaps within a struct.
type Task interface {
	DoPanicTask(*Info)
}

// Handle panics within the Task interface.
// Call it like this.
// panic_handler.HandleTask(panic_handler.Task)
func HandleTask(t Task) {
	i := newInfo(recover(), debug.Stack())
	if i == nil {
		return
	}
	defer nestedPanic(i)
	t.DoPanicTask(i)
}
