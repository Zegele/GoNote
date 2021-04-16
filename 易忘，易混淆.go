1.等价于什么？用途

1.1 rune
Unicode字符 rune 类型是和 int32 等价的类型，通常用于表示一个Unicode码点。
文本字符串通常被解释为采用UTF8编码的Unicode码点（rune）序列。
%c 字符（rune） (Unicode码点)

我们可以将Unicode码点也写到字符串面值中:

\'      单引号（只用在 '\'' 形式的rune符号面值中 ）
\"      双引号（只用在 "..." 形式的字符串面值中）"
\r     	回车


Unicode（ http://unicode.org ）， 
它收集了世界上所有的符号系统，包括重音符号和其他变音符号，制表符和回车符，还有很多神秘的符号，每个符号都分配一个唯一的Unicode码点，
Unicode码点对应Go语言中发rune整数类型。

更大的Unicode码点也是采用类似的策略处理。
0xxxxxxx rune 0-127 (ASCII)


1.2 []byte
byte 也是 uint8 类型的等价类型，byte 类型一般用于强调数值是一个原始的数据,而不是一个小的整数。

字符串和字节slice之间可以相互转换：
s := "abc"
b := []byte(s) //将s字符串转化为byte类型的切片。
s2 := string(b) //将byte类型的切片转化为字符串类型。

bytes包中也对应的六个函数：
func Contains(b, subslice []byte) bool
func Count(s, sep []byte) int
func Fields(s []byte) [][]byte
func HasPrefix(s, prefix []byte) bool
func Index(s, sep []byte) int
func Join(s [][]byte, sep []byte) []byte

 bytes包还提供了Buffer类型用于字节slice的缓存。(p110)


1.3 go 的 [] rune 和 [] byte 区别
https://blog.csdn.net/zhizhengguan/article/details/104538840?utm_medium=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-BlogCommendFromMachineLearnPai2-1.nonecase
first := "社区"
fmt.Println([]rune(first))//[31038 21306] 一个对应一个汉字
fmt.Println([]byte(first))//[231 164 190 229 140 186] 三个对应一个汉字，所以一个汉字是三个字节

s := "截取中文"
//试试这样能不能截取?
res := []rune(s)
fmt.Println(string(res[:2])) // 截取


s := "截取中文"
//试试这样能不能截取?
fmt.Println(s[:3]) // 截 
底层会默认将中文转化成 []byte， 而不是 []rune。

2.结构体怎么赋值
type Point struct {
	X, Y float64
}

var P Point
P = Point{X: 10.0, Y: 10.0}

var Q Point
Q = Point{3, 4}

Z := Point{0, 0}




3. 指针
一个指针的值是另一个变量的地址。一个指针对应变量在内存中的储存位置（内存地址）。

如果用“var x int”声明语句声明一个x变量，那么&x表达式（取x变量的内存地址）将产生一个指向该整数变量的指针，
指针对应的数据类型是*int，
指针被称之为“指向int类型的指针”。
如果指针名字为p，那么可以说“p指针指向变量x”，
或者说“p指针保存了x变量的内存地址”。
同时*p表达式对应p指针指向的变量的值。
一般*p表达式读取指针指向的变量的值，这里为int类型的值。
同时因为*p对应一个变量，多以该表达式也可以出现在赋值语句的左边，表示更新指针所指向的变量的值。
x := 1
p := &x //p,of type *int, points to x 
//p是指针, *int类型（p、&x都是*int类型）；“p指针指向变量x”或“p指针保存了x变量的内存地址”
fmt.Println(*p) // "1" 取指针指向的变量的值
fmt.Println(p) //会打印出保存的内存地址？？？（这一条是自己加的）
*p = 2 //equivalent to x = 2 通过指针改变了变量的值。
fmt.Println(x) // "2"

c := 0
fmt.Println(&c) // 内存地址
fmt.Println(*&c) // "0" * &c （*后面跟指针类型，才是指指针所值向地址的值）


4. map
4.1 
ages := make(map[string]int)

4.2 
ages := map[string]int{
	"alice": 31,
	"charlie": 43,
}

这相当于:
ages:= make(map[string]int)
ages["alice"] = 31
ages["charlie"] =34

4.3 给map赋值
ages["alice"] = 32
fmt.Println(ages["alice"]) // "32"

4.5 删除map中的元素（p135）
使用内置的delete函数可以删除元素：
delete(ages, "alice") // remove element ages["slice"]

4.6 以上所有操作是安全的，即使这些元素不在map中也没有关系。
如果一个查找失败将返回value类型对应的零值，例如，即使map中不存在“bob”下面的代码也可以正常工作，因为ages["bob"]失败时将返回0。
ages["bob"] = ages["bob"]+1 
//等号右边的bob检索不到返回0值，但是加1后，总数是1。然后赋值给ages["bob"]，所以创建了bob这个key以及也有对应的value。

4.7. 而且 x += y 和 x++等简短赋值语法也可以用在map上，所以上面的代码可以改成;
ages["bob"] += 1

更简短的写法：
ages["bob"]++

4.8. 但是map中的元素并不是一个变量，因此我们不能对map的元素进行取指操作：
_ = &ages["bob"] // compile error: cannot take address of map element.
禁止对map元素取址的原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而可能导致之前的地址无效。

5. contunue(p31)
continue语句直接跳到for循环的下个迭代开始执行。

6. Go语言还提供了以下的bit位操作运算符，前面4个操作运算符并不区分是有符号还是无符号数：（p85）
&  位运算 AND（和）（二进制下，两个数，每一位，同时都是1，则返回1，否则返回0）
|  位运算 OR（或）（二进制下，两个数，每一位，任意一位有1，则返回1，否则返回0）
^  位运算 XOR （ ）（两个数，每一位，只有一位是1，则返回1，否则返回0）
&^ 位清空 （AND NOT） （a &^ b 的意思是：将b中为1的位对应与a的位清0，a其他的位不变。例如：0110 &^ 1011 = 0100）
<< 左移
>> 右移

7. 怎样调用包


8. break和continue（p50）
break和continue语句会改变控制流。
break会中断当前的循环，并开始执行循环之后的内容。
continue会跳过当前循环，并开始执行下一次循环。

