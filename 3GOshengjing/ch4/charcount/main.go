// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)    //rune等价int32。 counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int //count of lengths of UTF-8 encodings 这是个数组
	invalid := 0                    //count of invalid UTF-8 characters
	in := bufio.NewReader(os.Stdin)
	//if in == "exit" {
	//	os.Exit(0)
	//}
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 { //unicode.ReplacementChar表示无效字符，并且无效字符编码长度是1
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}

}

//去好好看看怎样使用标准输入
// go build
// ./charcount.exe < a.txt
// 然后会打印出数据
