// Findlins3 crawls the web, starting with the URLs on the command line.
package main

import (
	"GoNote/3GOshengjing/ch5/links"
	"fmt"
	"log"
	"os"
)

//!+breadthFirst
//breadthFirst calls f for each item in the worlist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//!-crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}

//!-main
