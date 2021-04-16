// p55
//Boiling Prints the boiling point of water.

package main //说明该文件属于哪个包。这个就属于main包

import "fmt"

//导入依赖的包

//然后是包一级的类型、变量、常量、函数的声明语句。
//包一级的声明语句的顺序无关紧要。
//（但函数内部的名字则必须先声明之后才能使用）
const boilingF = 212.0

func main() {
	var f = boilingF
	var c = (f - 32) * 5 / 9
	fmt.Printf("boiling point= %g°F or %g°C\n", f, c)
	//%g表示浮点数。（格式化输出，可参考p29）

	//Output:
	//boiling point = 212°F or 100°C
}
