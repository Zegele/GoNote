// Bytecounter demonstates an implementation of io.Writer that counts bytes.

package main

import "fmt"

//!+bytecounter

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

//!-bytecounter

func main() {
	//!+main
	var c ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c) // "5", = len("hello")
	fmt.Println(&c)
	fmt.Println(*&c) // * &c （*后面跟指针类型，才是指指针所值向地址的值）

	c = 0 // reset the counter
	a := &c
	fmt.Println(*a)
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // "12", = len("hello, Dolly")
	fmt.Println(&c)
	//!-main
}
