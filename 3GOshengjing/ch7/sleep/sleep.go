// The sleep program sleeps for a specified period of time.

package mian

import (
	"flag"
	"fmt"
	"time"
)

//!+sleep
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

//!-sleep

// $ go build ...sleep
// $ ./sleep
// Sleeping for 1s...
