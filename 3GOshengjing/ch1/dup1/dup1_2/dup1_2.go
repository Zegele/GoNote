package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var i = 0
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
		i += 1
		if i > 5 {
			break
		}

	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
