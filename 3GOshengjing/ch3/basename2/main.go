// Basename2 reads file names from stdin and prints the base name of each one.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(basename(input.Text()))
	}
	//Note: ignoring potential errors from input.Err()
}

// basename removes directory components and a trailing .suffix.
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c

func basename(s string) string {
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	//slash ：= strings.LastIndex（s, "/"）该函数是在s中找排在最后的"/"字符在第几位，返回一个整数。如果没有该字符，则返回-1。
	fmt.Println(slash)
	s = s[slash+1:]
	fmt.Println(s)
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}

//go build
//./ .. .exe
//输入测试
