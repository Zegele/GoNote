// dup1 prints the text of each line that appears more than once in the standard input, preceded by its count.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	//NOTE: ignoring potential errors from input.Err()
	//难道是没有这个错误处理就运行不下去？
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

//去好好看看怎样使用标准输入
// go build
// ./dup1.exe < a.txt
// 然后会打印出数据

/*
func main() {
	Counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text() //input.Text()
		if line == "exit" {
			os.Exit(0)
		}

		Counts[input.Text()]++
		fmt.Println(Counts)
		//上面一句等于下面的表述：
		//line := input.Text()
		//counts[line] = counts[line] + 1
		for line, n := range Counts {

			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line)
			}
		}
	}
}

//为啥把for line, n这整段代码放入上面就行。
//独立放出来就不行，独立放出来是把内容打印到哪里了???
*/
