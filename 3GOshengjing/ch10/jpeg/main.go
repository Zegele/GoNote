//!+main

// The jpeg command read a PNG image from the standard input
// and writes it as a JPEG image to the standard output.
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // register PNG decoder //如果没有这一句，后面的结果会不一样
	"io"
	"os"
)

func main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

//!-main
/*
//!+with
$ go build ...ch3/mandlbort
$ go build .../ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
//!-with

//包引用中如果没有 	_ "image/png" 这一条，就是下面这个结果。
//!+without
$ go build ...ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format
//!-without
*/
