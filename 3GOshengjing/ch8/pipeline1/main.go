// Pipeline1 demonstrates an infinite 3-stage pipeline.
package main

import "fmt"

//!+
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ { //这里中间空的，可以没有条件，也就是没有“如果xxx”,直接跟“就YYY这样”。
			naturals <- x
		}
	}()

	//Squarer
	go func() {
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}
