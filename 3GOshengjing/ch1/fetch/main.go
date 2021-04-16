package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil { //看例子中的错误处理
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1) //用于退出
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() //关闭resp的Body流
		if err != nil {   //看例子中的错误处理
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

/*命令（git或终端中输入）

$ go build ...ch1/fetch
$ ./fetch.exe http://gopl.io
//会打印网页代码

*/
