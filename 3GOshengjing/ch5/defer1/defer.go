// Defer1 demonstrates a deferred call being invoded during a panic.

package main

import "fmt"

//!+f
func main() {
	f(3)
}

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

// !-f
