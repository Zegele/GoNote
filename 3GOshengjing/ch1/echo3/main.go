package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], "*"))
	//strings.Join是将n个字符合并，后面的“*”值在合并的中间加入这个符号，当然也可以是空的。
}
