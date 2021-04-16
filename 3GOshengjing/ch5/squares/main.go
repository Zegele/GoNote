// The squares program demonstrates a function value with state.

package main

import "fmt"

//!+
// squares returns a function that returns
// the next square number each time it is called

func squares() func() int {
	var x int
	return func() int {
		x++ //匿名函数可以引用squares函数中的变量
		//这也是匿名函数可以访问完整的词法环境的意思
		return x * x
	}
}

func main() {
	f := squares()
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	fmt.Println(f()) // "16"
}

//!-
