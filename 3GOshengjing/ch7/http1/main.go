// Http1 is a rudimentary e-commerce server.

package main

import (
	"fmt"
	"log"
	"net/http"
)

//!+main
func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) } //$50.00  %.2f 是表示两位小数的浮点数。

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

//!-main

/*
//!+handler
package http

type Handler interface{
	ServeHTTP(w ResponseWriter, r *Request)
}

func ListenAndServe(address string, h Handler) error
//!-handler
*/

// $ go build .../http1
// ./http1 &
// &是什么意思？
