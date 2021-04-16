package main

import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		fmt.Println(len(os.Args))
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)

	for n, v := range os.Args { //os.Args本身就是有两个值的
		fmt.Printf("num:%d, val:%v\n", n, v)
	}
}

//运行：go run main.go aaa bbb ccc
//os.Args的意思是 命令和字符都切片化，
//第一个元素是命令，也就是go run ...
//后面的元素就是命令行参数。（例子中是3个参数）
