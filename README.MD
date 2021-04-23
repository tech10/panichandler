# Introduction

This package will allow you to deal with panics in an extremely simple fassion. No channels, no contexts, no process forking, only a function call when a panic is caught.


## Explanation

In writing one of my programs, I wanted a way to handle panics in the same fassion across the application. Using a defer statement to do this seemed fine, but calling an included function within that specific program seemed somewhat limiting to me, so I created this package, just in case anyone else might want to do something similar.

Full documentation is available with go doc, but here are a couple of different things.

The panic_handler.Info struct has variables that will be populated with the precise value of the call to recover, along with its string and byte formatted copies. The stack trace is also available as bytes and a string. A complete format of the string value to recover, then the stack trace on another line, is available in both string and bytes using appropriate functions.

With all this information, you should be able to create a defer statement anywhere you like, which can then catch panics that you will then be able to do anything you like with, such as sending yourself an email with their values, etc. Typically, you should create the defer statement in one of two places, either when you begin execution of a new goroutine, or at a location where you may not want the entire program to crash.

The panic_handler.HandlerFunc is only called when a panic is caught. If a panic is not caught, this function isn't called. If you initialize the panic_handler.Handle function with a nil value, the program will print the string formatted value and stack trace, then exit the program with a default exit code of 111. You can change this exit code with the panic_handler.ExitCode variable.

This package makes no attempt to be concurrent safe within itself, due to the fact that multiple panics could be caught concurrently at one time, each one performing a specific task once they are caught. It should be up to the user to properly test their program for race conditions and make it concurrent safe.

# Example

```go
package main

import (
	"fmt"
	"github.com/tech10/panic_handler"
)

func main() {
	fmt.Println("This will test catching panics.")
	defer panic_handler.Handle(func(i *panic_handler.Info) {
		fmt.Println("This function has been called because a panic has been caught. Here is the value of this panic.")
		fmt.Println(i.PanicString)
		fmt.Println("Here is the stack trace.")
		fmt.Println(i.StackString)
		fmt.Println("Goodbye!")
	})
	panic("This is a test panic.")
}
```

# Contributions

Open an issue or a pull request with any contributions. Remember to properly format your code with gofmt. If creating any new features, create tests for them as well.
