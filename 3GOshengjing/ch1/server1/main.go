// Server1 is a minimal "echo" server
//比如用户访问的是http://localhost:8000/hello ，那么响应是URL.Path = "hello"。
//这里涉及的标准库后面一定要搞懂。
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler) //each request calls handler
	log.Fatel(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q \n", r.URL.Path)
	//%q 带双引号的字符串"abc"或带单引号的字符'c'
}

//相关命令
//1. 后台运行这个服务器程序
// 在后台运行这个服务程序 在运行命令的末尾加上一个 & 符号。
// $ go run src/.../main.go &
// windows下可以在另外一个命令行窗口去运行这个程序。
//
//2. 通过命令行发送客户端请求：
//$ go build .../fetch
//$ ./fetch http://localhost:8000
//URL.Path = "/"
//$ ./fetch http://localhost:8000/help
//URL.Path = "/help"
