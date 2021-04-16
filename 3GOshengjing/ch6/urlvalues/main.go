// The urlvalues command demonstrates a map type with methods.

package main

import (
	"fmt"
	"net/url"
)

/*
// ！+values
package url

// Values maps a string key to a list of values.
type Values map[string][]string
// Get returns the first value associated with the given key.
// or "" if there are none.
func (v Values) Get(key string) string{
	if vs := v[key]; len(vs) > 0{
		return vs[0]
	}
	return ""
}

// Add adds the value to key.
// It appends to any existing values associated with key.
func (v Values) Add(key, value string){
	v[key] = append(v[key], value)
}
//!-values
*/

func main() {
	//!+main
	m := url.Values{"lang": {"en"}} // direct construction
	n := &m
	fmt.Println(m)
	fmt.Println(n)
	fmt.Println(*n)
	fmt.Println(&n)
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Println(m.Get("lang")) //"en"
	fmt.Println(m.Get("q"))    // ""
	fmt.Println(m.Get("item")) // "1" (first value)
	fmt.Println(m["item"])     // "[1 2]" (direct map access)
	b := &m
	fmt.Println(m)
	fmt.Println(&b)
	fmt.Println(m["item"][0])

	m = nil
	v := &m
	fmt.Println(m)
	fmt.Println(&v)
	fmt.Println(m.Get("item")) // ""
	m.Add("item", "3")         // panic: assignment to entry in nil map

}
