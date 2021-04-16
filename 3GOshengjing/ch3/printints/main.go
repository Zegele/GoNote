// intsToString is like fmt.Sprint(values) but adds commas.
package main

import (
	"bytes"
	"fmt"
)

func intsToString(values []int) string { //[]int 是切片类型
	var buf bytes.Buffer //bytes.Buffer Buffer类型用于字节slice的缓存，开始Buffer是空的
	fmt.Println(buf)
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v) //将数据v，按照%d这个格式，输入到&buf这个内存地址中。
		//Fprintf就是将数据格式化，并输入到指定的文件中。
		//Printf：将数据 按照定义的格式固定输出到屏幕；
		//Sprintf:将数据按固定格式输出到某个缓存中，可以这么理解，其实它功能和printf差不读，只是printf输出到了屏幕上，而sprintf只是将同样的内容装到了某个缓存中。
		//Fprintf:同上，该函数将要输出的数据，按照定义的格式，输出到了文件中，比如可以输入到一个txt中。
	}
	buf.WriteByte(']')
	fmt.Println(buf)    // 	{[91 49 44 32 50 44 32 51 93] 0 0} 44代表 ，(逗号) ,32代表 (空格)
	return buf.String() //将buf转成字符串

}
func main() {
	fmt.Println(intsToString([]int{1, 2, 3})) //"[1, 2, 3]"
}
