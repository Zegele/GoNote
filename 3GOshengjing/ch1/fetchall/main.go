// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string) //创建了一个传递string类型参数的channel，

	for _, url := range os.Args[1:] {
		go fetch(url, ch) //start a goroutine
		//让fetch函数在goroutine中异步执行http.Get方法。
	}
	//对每一个命令行参数，我们都用go这个关键字来创建一个goroutine。
	//并且让函数咋这个goroutine异步执行http.Get方法。（联系下面fetch函数的参数查看）

	for range os.Args[1:] {
		fmt.Println(<-ch) //receive from channel ch 还能这样用？？？
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) //send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	//io.Copy会把响应的Body内容拷贝到ioutil.Discard输出流中（译注：可以把这个变量看做一个垃圾桶，可以向里面写一些不需要的数据）
	//但是因为我们需要这个方法返回的自己数（也就是需要知道内容的大小），但是有不想要传入的内容。

	resp.Body.Close() //don't leak resources

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
		//每当请求返回内容时，fetch函数都会往ch这个channel里写如一个字符串，
		//由main函数里的第二个for循环来处理并打印channel里的这个字符串。
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

//当一个goroutine尝试在一个channel上做send或者receive操作时，
//这个goroutine会阻塞在调用处，直到另一个goroutine往这个channel里写入、或接收值，
//这样两个goroutine才会继续执行channel操作之后的逻辑。

//在这个例子中，每一个fetch函数在执行时都会往channel里发送一个值（ch<-expression）,
//主函数负责接收这些值（<-ch）。
//这个程序中我们用main函数来接收所有fetch函数传回的字符串，可以避免在goroutine异步执行还没有完成时main函数提前退出。
//

/*命令
go build ...fetchall
./fetchall.exe http://golang.org http://gopl.io http://godoc.org
后面就会打印出获取的时间，大小，以及对应的网址。
*/
