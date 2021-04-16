//Server2 is a minimal "echo" and counter server.

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler) //URL后面是任意的，都会调用handler这个函数。一般就是默认的处理函数
	//如果你的请求pattern（模式，“对象”的意思？）是以/结尾，那么所有以该url为前缀的url都会被这条规则匹配。

	http.HandleFunc("/count", counter)                    //URL后是/count，调用counter这个函数。
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) //这是起一个服务器？
}

//handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock() //这啥意思？？？
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

//counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

//这些代码的背后，服务器每一次接收请求处理时都会另起一个goroutine，这样服务器就可以同一时间处理多个请求。
//然后在并发情况下，假如真的有两个请求同一时刻去更新count，那么这个值可能并不会被正确地增加；
//这个程序可能会引发一个严重的bug：竞态条件（参见9.1）
//为了避免这个问题，我们必须保证每次修改变量的最多只能有一个goroutine，这也就是代码里的mu.Lock()和mu.Unlock()调用将修改count的所有行为包在中间的目的。
//第九章中我们会进一步讲解共享变量。
