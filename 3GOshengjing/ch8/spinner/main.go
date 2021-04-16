// Spinner displays an animatio while computing the 45th Fibonacci number.

package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	//spinner函数和fib(45)函数在同时进行。
	//main函数结束后，goroutine才结束。
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
	//\r 回车
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r) //%c 字符（rune） (Unicode码点)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
