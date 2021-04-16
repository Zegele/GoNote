// The trace program uses defer to add entry/exit diagnostics to a function.

package main

import (
	"fmt"
	"log"
	"time"
)

//!+main
func bigSlowOperation() int {
	defer trace("bigSlowOperation")() //  don't forget the extra parentheses
	// ...lots of work...
	time.Sleep(10 * time.Second) // simulate slow operation by sleeping
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}

//!-mian

func main() {
	fmt.Println(bigSlowOperation())
}

/*
!+output
$ go build .../trace
$ ./trace
2020/...    enter bigSlowOperation
2020/...    exit bigSlowOperation (10.xxx s)
!-output
*/
