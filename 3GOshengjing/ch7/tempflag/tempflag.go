// Tempflag prints the value of its -temp (temerature) flag.
package main

import (
	"GoNote/3GOshengjing/ch7/tempconv"
	"flag"
	"fmt"
)

//!+
var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

// temp就是CelsiusFlag函数返回值的类型 fmt.Printf("%T", temp) //*tempconv.Celsius（tempconv包的Celsius类型的指针类型）

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

/*
$ go build ...tempflag
$ ./tempflag
20°C

$ ./tempflag -temp -18°C
-18°C

$ ./tempflag -temp -212°F
100°C

$ ./tempflag -temp -273.15K
invallid value "273.15K" for flag -temp: invalid temperature "273.15K"

Usage of ./tempflag:
	-temp value
		the temperature (default 20°C)

$ ./tempflag -help
Usage of ./tempflag:
	-temp  value
		the temperature (default 20°C)
*/
