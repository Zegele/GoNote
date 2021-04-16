// The thumbnail command produces thumbnails of JPEG files
// those names are provided on each line of the standard input.

// The "+build ignore" tag excludes this file from the
// thumbnail package, but it can be compiled as a command and run like
// this:

// Run with:
// $ go run $GOPATH/src/...thumbnail/main.go
// foo.jpeg
// ^D (linux系统的命令？)
//

package main

import (
	"GoNote/3GOshengjing/ch8/thumbnail"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		thumb, err := thumbnail.ImageFile(input.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(thumb)
	}
	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}
