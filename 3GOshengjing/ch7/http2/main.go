// Http2 is an e-commerce server with /list and /price endpoints.

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

//!+handler
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound) //404
			fmt.Fprintf(w, "no such item: %q\n", item)
			return // 这里为什么放个return？？必须放的么？
			// return 的作用是 这个检查应该在向w写入任何值前完成？（p256）
		}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound) //404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}

//!-handler
