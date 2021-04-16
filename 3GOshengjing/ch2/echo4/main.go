//P61
//Echo4 prints its command-line arguments Echo4打印其命令行參數

package main

import (
	"flag"
	"fmt"
	"strings"
)

var n = flag.Bool("n", false, "omit trailing newline省略尾隨換行符") //
var sep = flag.String("s", " ", "separator")

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))
	if !*n { //如果是非n
		fmt.Println() //因为是非n，所以结尾要换行的话，所以就要打印一个空行。
	}
}

/*
$ go build .../echo4
$ ./echo4.exe a bc def
a bc def

$ ./echo4.exe -s / a bc def //这是把分隔符设置成了斜杠
a/bc/def

$ ./echo4.exe -n a bc def //省略行尾的换行符
a bc def$ //$在尾部表示省略了尾部的换行符的意思？

$ ./echo4.exe -help
Usage of ./echo4:
	-n  omit trailing newline
	-s string
		separator (default " ") //分隔符（默认 “ ”），分隔符默认是空格。
*/
