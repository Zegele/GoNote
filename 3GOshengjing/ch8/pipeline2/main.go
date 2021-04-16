// Pipeline2 demonstrates a finite 3-stage pipeline.

package main

import "fmt"

//!+
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 10; x++ {
			naturals <- x //channel内如果没有接收，会阻塞，所以channel内只会有一个值。
		}
		close(naturals)
	}()

	// Squarer
	go func() {
		for x := range naturals { //这个range是怎么进行的？因为channel内只有一个值。
			squares <- x * x
		}
		close(squares)
	}()

	// Printer( in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}
}
