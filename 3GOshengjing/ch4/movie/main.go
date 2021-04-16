// Movie prints Movies as JSON

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title  string
	Year   int      `json:"released"`
	Color  bool     `json:"color, omitempty"`
	Actors []string //数组类型？
}

var movies = []Movie{ //定义了movies是个Movie结构体类型的切片
	{Title: "Casabance", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}

func main() {
	{
		//!+Marchal
		data, err := json.Marshal(movies)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)

		}
		fmt.Printf("%s\n", data)
		//!-Marshal
	}
	{
		//!+MarshalIndent
		data, err := json.MarshalIndent(movies, "", "   ")
		if err != nil {
			log.Fatalf("JSON marshling failed: %s", err)

		}
		fmt.Printf("%s\n", data)
		//!-MarshalIndent

		//!+ Unmarshal
		var titles []struct {
			Title  string
			Year   int
			Color  bool
			Actors []string
		}
		//fmt.Println(data)
		//fmt.Printf("%s\n", data)
		if err := json.Unmarshal(data, &titles); err != nil { // 这个data就是上面命名过的，是JSON格式的文件
			log.Fatalf("JSON unmarshaling failed: %s", err)
		}
		fmt.Println(titles) // "[{Casablance} {Cool Hand Luke} {Bullitt}]"
		//!-Unmarshal
	}
}
