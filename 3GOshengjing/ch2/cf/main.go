package main

import (
	"GoNote/3GOshengjing/ch2/tempconv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	for _, arg := range os.Args[1:] { //for（等关键字）后面手动空格，这种格式化不了。。。
		//Args查询参照Go1.go文件。
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Fahrenheit(t)
		c := tempconv.Celsius(t)
		fmt.Printf("%s = %s, %s = %s\n",
			f, tempconv.FToC(f), c, tempconv.CToF(c))

	}
}

//go build 路径/cf(src之后的路径)
// ./cf.exe 100 45或./cf 100 45 (运行cf文件，并且计算后面的参数)
