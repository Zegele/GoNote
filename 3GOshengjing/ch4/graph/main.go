// Graph shows how to use a map of maps to represent a directed graph.
package main

import "fmt"

//!+
var graph = make(map[string]map[string]bool) //map中嵌套了map，map的key是string类型，value是map类型。

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil { //如果edges是个空值，这里edges是空map
		//fmt.Printf("%T\t%[1]v\n", edges)
		//map[string]bool	map[]
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to] // graph对应的[from]的元素，然后元素中[to]对应的bool值。
	//也就是返回值是true或false

}

func main() {
	addEdge("a", "b")              // 得到—— map[a:map[b:true]]
	addEdge("c", "d")              //  map[a:map[b:true] c:map[d:true]]
	addEdge("a", "d")              //  map[a:map[b:true d:true] c:map[d:true]]
	addEdge("d", "a")              //  map[a:map[b:true d:true] c:map[d:true] d:map[a:true]]
	fmt.Println(hasEdge("a", "b")) // true
	fmt.Println(hasEdge("c", "d")) // true
	fmt.Println(hasEdge("a", "d")) // true
	fmt.Println(hasEdge("d", "a")) // true
	fmt.Println(hasEdge("x", "b")) // false
	fmt.Println(hasEdge("c", "d")) // true
	fmt.Println(hasEdge("x", "d")) // false
	fmt.Println(hasEdge("d", "x")) // false
	fmt.Println(hasEdge("a", "c")) // false 0值也是false
}
