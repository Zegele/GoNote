//Dedup prints only one instance of each line; deplicates are removed.

package main

import (
	"bufio"
	"fmt"
	"os"
)

//!+
func main() {
	seen := make(map[string]bool) // a set of strings
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "exit" {
			os.Exit(0) //表示退出。
		}
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
		fmt.Println(seen)
	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}

//!-

//go build
// ./dedup.exe
// a 输入a
// a 返回a 再输入a就不会返回a了。也就是输入相同的无效。只有输入不同的才会加入到map中。
