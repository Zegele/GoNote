// The sum program demonstrates a variadic function.

package main

import "fmt"

//!+
func sum(vals ...int) int { //vals被看作是类型为[]int的切片（vals被看作是类型为int数组的切片）
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

//!-

func main() {
	//!+main
	fmt.Println(sum())           // "0"
	fmt.Println(sum(3))          // "3"
	fmt.Println(sum(1, 2, 4, 3)) // "10"
	//!-main

	//!+slice
	values := []int{1, 2, 3, 4}
	fmt.Println(sum(values...)) // "10"
	//!-slice

}
