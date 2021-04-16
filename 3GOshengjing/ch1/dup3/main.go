// Dup3 prints the count and text of lines that
// appear more than once in the named input files.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename) //ReadFile函数返回一个字节切片（byte slice），赋值给data
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n") { //把字节切片转化为string类型，才能被string.Split使用。
			counts[line]++
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

// go build
// ./dup3.exe a.txt b.txt c.txt
// 即可看到结果
