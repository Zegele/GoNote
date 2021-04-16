package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			os.Exit(0) //表示退出。
		}
		fmt.Println("line=", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取输入错误!", err)
	}
}
