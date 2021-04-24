package panic_handler

import (
	"fmt"
	"os"
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
	if t == nil {
		fmt.Fprintf(os.Stderr, "WARNING!!!\nThe HandleTask function cannot have a nil pointer.\n%s", i.String())
		os.Exit(ExitCode)
	}
	taskRun(i, t, ExitCode)
}

func taskRun(i *Info, t Task, e int) {
	defer nestedPanic(i, e)
	t.DoPanicTask(i)
}
