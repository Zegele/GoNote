//  Nonempty is an example of an in-place slice algorithm（算法）.在切片空间内部的算法。

//!+nonempty
package main

import "fmt"

func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" { //当循环到 ""空值时，结束if循环了，这时的i依旧=1，所以当range到第3个元素时，将第3个元素添加到了string[1]的位置。所以没有空值了。
			strings[i] = s
			i++
		}
	}

	return strings[:i]
}

//!-nonempty

func main() {
	//!+main
	data := []string{"one", "", "three"}
	fmt.Printf("%v\n", nonempty2(data))
	fmt.Printf("%q\n", nonempty(data)) // '["one" "three"]'
	fmt.Printf("%q\n", data)           // '["one" "three" "three"]'

	//!-main
}

// +alt
func nonempty2(strings []string) []string {
	out := strings[:0] //zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
