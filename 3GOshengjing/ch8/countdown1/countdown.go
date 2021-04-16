// Countdown implements the coutdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second) //time.Tick函数返回一个channel。
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick //把这个值丢弃了？也就是把channel清空了？
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
