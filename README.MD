# Introduction

This package will allow you to deal with panics in the following ways:

- Call a function when a panic is called.
- Send the information from the panic over a channel.
- Cancel contexts with a provided context.CancelFunc.
- Implement the panic_handler.Task interface by defining its DoPanicTask method.
- All these things in a struct, which can execute everything stated above, only some of these things, or only one of them.


## Explanation

In writing one of my programs, I wanted a way to handle panics in the same fassion across the application. Using a defer statement to do this seemed fine, but calling an included function within that specific program seemed somewhat limiting to me, so I created this package, just in case anyone else might want to do something similar. I expanded it since then to work with channels, context cancel functions, implemented a Task interface, and a struct which can deal with everything.

[Full documentation](https://pkg.go.dev/github.com/tech10/panic_handler) is available, but here are a couple of different things.

The panic_handler.Info struct has variables that will be populated with the precise value of the call to recover, along with its string and byte formatted copies. The stack trace is also available as bytes and a string. A complete format of the string value to recover, then the stack trace on another line, is available in both string and bytes using appropriate functions.

With all this information, you should be able to create a defer statement anywhere you like, which can then catch panics that you will then be able to do anything you like with, such as sending yourself an email with their values, etc. Typically, you should create the defer statement in one of two places, either when you begin execution of a new goroutine, or at a location where you may not want the entire program to crash.

The panic_handler.HandlerFunc is only called when a panic is caught. If a panic is not caught, this function isn't called. If you initialize the panic_handler.Handle function with a nil value, the program will print the string formatted value and stack trace, then exit the program with a default exit code of 111. You can change this exit code with the panic_handler.ExitCode variable.

The same exit code will be used if the function you define to catch a panic, causes a panic itself. This will be the same for the Task interface, also. The complete stack trace will be printed for both panics, along with the string value of the reason behind the panic. There may be redundent information in the stack traces, as in my testing, the stack trace from the second panic included everything from the first. Both are printed for now, however, to provide a bit more clarity on where individual panics occurred. This may change in the future, and as it won't be a breaking API change, only an information update, a patch version update will satisfy this, should it be changed.

When sending the data of a panic over a channel with the HandleWithChan function, the data will only be sent if a panic is caught.

When using the panic_handler.HandleWithContextCancel function, you cannot declare the context.CancelFunc value as nil. The program will warn you of an uncatchable panic, then will terminate. You can, however, make the panic_handler.HandleFunc a nil value, in which case, only the context.CancelFunc will be called. I wouldn't recommend this, though, as you may want to do something more with your panics.

The panic_handler.Task interface defines one method, DoPanicTask(*panic_handler.Info)

The panic_handler.Capture struct can catch panics and pass them to a function, task interface, channel, and context cancel function, in that specific order and one at a time. A nil value on all such values will cause the program to exit with the stacktrace of the panic that was essentially uncaught. Catching panics with an uninitialized struct is not going to do anything for you, so the library won't let you do this. There is only one partial exception to this rule. If you catch a panic by using an already canceled contexts CancelFunc, the panic will be caught silently. This is bad practice, don't do this, and please, do check your programs.

This package makes no attempt to be concurrent safe within itself, due to the fact that multiple panics could be caught concurrently at one time, each one performing a specific task once they are caught. It should be up to the user to properly test their program for race conditions and make it concurrent safe.

All tests should pass race testing.


# Examples

## Function

```go
package main

import (
	"fmt"
	"github.com/tech10/panic_handler"
)

func main() {
	fmt.Println("This will test catching panics.")
	defer panic_handler.Handle(func(i *panic_handler.Info) {
		fmt.Println("This function has been called because a panic has been caught. Here is the reason for this panic.")
		fmt.Println(i.PanicString)
		fmt.Println("Here is the stack trace.")
		fmt.Println(i.StackString)
		fmt.Println("Goodbye!")
	})
	panic("This is a test panic.")
}
```


## Channel

```go
package main

import (
	"fmt"
	"github.com/tech10/panic_handler"
)

func main() {
	fmt.Println("This will test catching panics over a channel.")
	c := make(chan *panic_handler.Info)
	go func() {
		defer panic_handler.HandleWithChan(c)
		panic("Testing channels.")
	}()
	i := <-c
	fmt.Println("This program has continued execution because a panic has been caught. Here is the reason for this panic.")
	fmt.Println(i.PanicString)
	fmt.Println("Here is the stack trace.")
	fmt.Println(i.StackString)
	fmt.Println("Goodbye!")
}
```


## Context cancelation

```go
package main

import (
	"context"
	"fmt"
	"github.com/tech10/panic_handler"
	"sync"
)

func main() {
	fmt.Println("Catching a panic and canceling contexts.")
	var wg sync.WaitGroup
	mctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-mctx.Done()
		fmt.Println("Parent context canceled.")
	}()
	var cctx context.Context
	cctx, _ = context.WithCancel(mctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-cctx.Done()
		fmt.Println("Child context canceled.")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer panic_handler.HandleWithContextCancel(cancel, func(i *Info) {
			fmt.Println("Panic caught.\n", i.String())
		})
		panic("Testing context cancelations.")
	}()
	wg.Wait()
	fmt.Println("All contexts have been canceled. Goodbye.")
}
```


## Task interface

```go
package main

import (
	"fmt"
	"github.com/tech10/panic_handler"
	"sync"
)

type CP struct {
	sync.WaitGroup
}

func (cp *CP) DoPanicTask(i *panic_handler.Info) {
	defer cp.Done()
	fmt.Println("Panic captured.")
	fmt.Println(i.String())
}

func main() {
	fmt.Println("Capture panics with the Task interface.")
	c := &CP{}
	c.Add(2)
	go func() {
		defer c.Done()
		defer panic_handler.HandleTask(c)
		panic("Testing the Task interface with this panic.")
	}()
	c.Wait()
	fmt.Println("Panic caught, task complete. Goodbye.")
}
```


## Capture struct

This is the only example that will have an uncaught panic, to demonstrate that an uninitialized panic_handler.Capture struct can't be used.

```go
package main

import (
	"context"
	"fmt"
	"github.com/tech10/panic_handler"
	"sync"
)

type captT struct {
	f func()
}

func (c *captT) DoPanicTask(i *panic_handler.Info) {
	c.f()
}

func main() {
	fmt.Println("Testing capture struct. Watch as a panic is caught and passed down the stack of a function, interface, channel, and context cancelation.")
	var l sync.Mutex
	var wg sync.WaitGroup
	ctx, ctxc := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		l.Lock()
		defer l.Unlock()
		fmt.Println("Context canceled.")
	}()
	c := panic_handler.New()
	c.CC = ctxc
	c.C = make(chan *panic_handler.Info)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-c.C
		l.Lock()
		defer l.Unlock()
		fmt.Println("Channel received data.")
	}()
	wg.Add(1)
	c.F = func(i *panic_handler.Info) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
		fmt.Println("Function called.")
	}
	wg.Add(1)
	c.T = &captT{
		f: func() {
			defer wg.Done()
			l.Lock()
			defer l.Unlock()
			fmt.Println("Task interface executed.")
		},
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.Catch()
		panic("Demonstrating capture struct.")
	}()
	wg.Wait()
	fmt.Println("All the work is done.")
	fmt.Println("Let's take the task and channel out of the work, now. Notice that the context will not cancel again, since it has already been canceled.")
	c.T = nil
	c.C = nil
	fmt.Println("Catch another panic.")
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer c.Catch()
		panic("Removing work.")
	}()
	wg.Wait()
	fmt.Println("Let's only cancel a context. Don't do this with an already canceled context like this, it won't work out well for you.")
	c.F = nil
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.Catch()
		panic("A silent panic, you really don't want this in production.")
	}()
	wg.Wait()
	fmt.Println("And as you can see, a Capture struct with all values at nil won't capture panics, crash the program instead, as we are expecting to do something in the struct.")
	c.CC = nil
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.Catch()
		panic("Empty the struct.")
	}()
	wg.Wait()
	fmt.Println("If we got this far, something is wrong. Have a good day!")
}
```


# Contributions

Open an issue or a pull request with any contributions. Remember to properly format your code with gofmt. If creating any new features, create tests for them as well.
