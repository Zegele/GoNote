// Dup2 prints the count and text of lines that appear more than once
// in the input. It reads from stdin or from a list of named fileds.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int) //创建了一个map
	files := os.Args[1:]           //files 从命令的第一个元素取起。
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg) //os.Open返回两个值，第一个值是被打开的文件（文件类型是：*os.File）
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue //直接跳到for循环的下个迭代开始执行。不是向下执行，而是继续去for循环。
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}

	//NOTE: ignoring potential errors from input.Err()
}
