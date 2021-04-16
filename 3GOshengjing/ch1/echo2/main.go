package main

import (
	"fmt"
	"os"
)

func main() {
	s, sep := "", "" //定义s和sep是两个字符串类型，并且是空值
	for _, arg := range os.Args[1:] {
		s += sep + arg //这个时候sep是空值
		//相当于s = s + sep +arg 计算顺序也是这样的
		sep = " XXX" //这个时候sep是“XXX”
	}
	fmt.Println(s)
}
