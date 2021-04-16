// Basename1 read file names from stdin and prints the base name of each one.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(basename(input.Text()))
	}
	// NOTE: ignoring potential errors from input.Err()
}

// basename removes directory components and a .suffix.
// e.g. a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c

func basename(s string) string {
	// Discard last "/" and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break //会结束循环
		}
	}

	//fmt.Println(s)
	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	//fmt.Println(s)
	return s
}

//go build
//./ .. .exe
//输入测试
