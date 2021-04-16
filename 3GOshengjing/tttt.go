package main

import (
	"fmt"
)

var (
	firstName, lastName, s string
	i                      int
	f                      float32
	input                  = "56.12 / 5212 / Go"
	format                 = "%f / %d / %s"
)

func main() {
	fmt.Println("Please input your full name: ")
	fmt.Scanln(&firstName, &lastName)
	//fmt.Scanf("%s %s", &firstName, &lastName)
	fmt.Printf("Hi %s %s!\n", firstName, lastName)
	fmt.Sscanf(input, format, &f, &i, &s)
	fmt.Println("From the string we read: ", f, i, s)
}

/*
func main() {
	var buffer [512]byte

	n, err := os.Stdin.Read(buffer[:])
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Println("count:", n, ", msg:", string(buffer[:]))
}
*/

/*
func main(){
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("please input your name: ")
	input, err := inputReader.ReadString('\n')
	// 对unix：使用“\n”作为界定符，而window使用“\r\n”为界定符

if err!= nil{
	fmt.Println("There ware errors reading, exiting program.")
	return
}

fmt.Printf("Your name is %s", input)
}

//参考网址： https://my.oschina.net/u/1590519/blog/336118
*/
