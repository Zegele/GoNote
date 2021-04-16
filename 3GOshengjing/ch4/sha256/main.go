// The sha256 command computes the SHA256 hash (an array) of a string.

package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1) // %x是十六进制
	//fmt.Println(len(c1))
	//fmt.Println(c1)
	//fmt.Println(c2)
	c11 := &c1
	fmt.Printf("%T\n", c11)
	fmt.Println(c11)
	zero(c11)
	var zzz byte                   //byte是字节，zzz是字节类型
	fmt.Printf("%v\n%[1]T\n", zzz) // "0 uint8"字节类型本身就是uint8类型的。byte和uint8是等价的。
	//Output:
	//2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	//4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	//false
	//[32]unit8 这个uint8表示byte的数值范围，也就是byte的范围是uint8
	//对比 [32]byte 这个整体表示这是32位的字节类型的数组。字节类型。

}

func zero(ptr *[32]byte) { // ptr是[32]byte指针类型
	for i := range ptr {
		ptr[i] = 0
	}
	fmt.Println(*ptr)
	fmt.Println(ptr)
}
