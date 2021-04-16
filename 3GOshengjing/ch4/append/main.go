// Append illustrates the behavior of the built-in append function.

package main

import "fmt"

func appendslice(x []int, y ...int) []int { // ...表示接收变长的slice为参数
	//换成 y []int 则main程序能运行。这个...int到底是什么东西？？？
	var z []int
	zlen := len(x) + len(y)
	if zlen <= cap(x) {
		//There is room to expand the slice.
		z = x[:zlen]
	} else {
		//There is insufficient space.
		//Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // 给z中，复制x
	}
	copy(z[len(x):], y) // 给z的[len(x):],赋值y
	return z
}

// !+ append
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		// There is room to grow. Extend the slice.
		z = x[:zlen]
	} else {
		//There is insufficient space. 没有足够的空间。Allocate a new array.分配了一个新的数组
		//Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) //a built-in function（内置函数）; see text
	}
	z[len(x)] = y
	return z
}

// !- append

//!+growth
func main() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)

		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
		x = y
		fmt.Println(appendslice(x, y))
	}

}

// !-growth

/*
// !+ output

*/
