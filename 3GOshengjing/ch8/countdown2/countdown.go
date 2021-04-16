// Countdown implements the countdown for a recket launch.

package main

import (
	"fmt"
	"os"
	"time"
)

//!+
func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
		//struct{} ：表示struct类型
		//struct{}{} :表示struct类型的值，该值也是空。
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown. Press return to abour.")
	select {
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off !")
}
