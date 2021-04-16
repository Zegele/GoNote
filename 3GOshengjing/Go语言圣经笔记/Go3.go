第三章 基础数据类型

1. go语言将数据类型分为4类：基础类型、复合类型、引用类型和接口类型。
基础类型包括：数字、字符串和布尔型。（本章）

复合数据类型：数组（4.1）和结构体（4.2）
是通过组合简单类型，来表达更加复杂的数据结构。

引用类型包括：指针（2.3.2）、切片（4.2）、字典（4.3）、函数（5）、通道（8）
虽然数据钟来很多，但他们都是对程序中一个变量或状态的简介引用。
这意味着对任一引用类型数据的修改都会影响所有该引用的拷贝。

接口类型：第 7 章介绍


3.1 整型

2. 数值类型包括整形、浮点数和复数。每种类型都决定了对应的大小范围和是否支持正负符号。（P84）

int8 int16 int32 int64 有符号整形数
uint8 uint16 uint32 uint64 无符号整形数

Unicode字符 rune 类型是和 int32 等价的类型，通常用于表示一个Unicode码点。
byte 也是 uint8 类型的等价类型，byte 类型一般用于强调数值是一个原始的数据而不是一个小的整数。

uintptr 一种无符号的整数类型，没有制定具体的bit大小，但是足以容纳指针。
uintptr 只有在底层编程才需要。 特别是go函数库或操作系统接口交互的地方。

int8 -128--127（负的2的n-1次方到，2的n-1次方-1）， uint8 0--255（0到2的n次方-1）

3. 算数运算、逻辑运算、比较运算的二元运算符，按照优先级递减的顺序排列。（p84）

*(乘)  /（除）  %（余）  <<  >>  &  &^

+（加）  -（减）  |  ^  

==（等于） !=（不等于）  <（小于）  <=（小于等于）  >（大于）  >=（大于等于）

&&

||

4. + - * / 可以适用于整数、浮点数、和复数 (p85)
取模运算符 % 仅用于整数间的运算。

5. %
GO语言中驱魔运算的符号（正负），和被取模数的符号总是一致的。(p85)
如-5%3 和-5%-3结果都是-2

6. / (p85)
除法运算符的行为则依赖于操作数是否全为整数，
如 5.0 / 4.0结果是1.25， 5/4结果是1
因为整数除法会想着0方向截断余数。

7. 溢出 (p85)没懂怎么是这么个结果
如果一个算数运算的结果，不管是有符合或者是无符号的，如果需要更多的bit位才能正确表示的话，
就说明计算结果是溢出了。超出的高位bit位部分将被丢弃。//不理解这句
如果原始的数值是有符号类型，而且最左边的bit位是1的话，那么最终结果可能是负的。
例如：
var u uint8 = 255
fmt.Println(u, u+1, u*u) // 255 0 1

var i int8 =127
fmt.Println(i, i+1,i*i) // 127 -128 1  


8. 两个相同类型的整数类型可以使用下面的二元比较运算符进行比较。
比较表达式的结果是布尔类型。
==（等于） !=（不等于）  <（小于）  <=（小于等于）  >（大于）  >=（大于等于）

9. 布尔型，数字类型和字符串等类型都是可比较的，两个相同类型的值可用==和！=进行比较。
可比较就可以根据比较结果排序。如：整数、浮点数和字符串可以根据比较结果排序。

10. + -
对于整数，+x是0+x的简写，-x是0-x的简写；
对于浮点数和复数，+x就是x，-x是x的负数。

11. Go语言还提供了以下的bit位操作运算符，前面4个操作运算符并不区分是有符号还是无符号数：（p85）
&  位运算 AND（和）（二进制下，两个数，每一位，同时都是1，则返回1，否则返回0）
|  位运算 OR（或）（二进制下，两个数，每一位，任意一位有1，则返回1，否则返回0）
^  位运算 XOR （ ）（两个数，每一位，只有一位是1，则返回1，否则返回0）
&^ 位清空 （AND NOT） （a &^ b 的意思是：将b中为1的位对应与a的位清0，a其他的位不变。例如：0110 &^ 1011 = 0100）
<< 左移
>> 右移

12. ^（p86）
作为二元运算符时是按位异或（XOR），当用作一元运算符时表示按位取反，就是，他返回一个每一个bit位都取反的数。


13. <<左移 >>右移 （p86）
代码演示 uint8 类型值的8个独立的bit位。

var x uint8 = 1<<1 | 1<<5 // x= 1左移1位 或 1左移5为 = 10 或 100000 = 100010
var y uint8 = 1<<1 | 1<<2

fmt.Printf("%08b\n",x) // “00100010”
// %b是Printf函数的参数，打印二进制格式的数字。
// 其中%08b中08表示打印至少8个字符宽度，不足的前缀部分用0填充。
fmt.Printf("%08b\n",y) // "00000110"

fmt.Printf("%08b\n", x&y) // "00000010", the intersection{1}
fmt.Printf("%08b\n", x|y) // "00100110", the union{1,2,5}
fmt.Printf("%08b\n", x^y) // "00100100", the symmetric difference{2,5}
fmt.Printf("%08b\n", x&^y) // "00100000", the difference{5}

for i := uint(0); i<8; i++{ // uint(0)意思是将0强转成uint类型。
	if x&(1<<i) != 0 { //membership test （检查在哪个位被占，也就是哪个位为1。）
		fmt.Println(i) //"1", "5"
	}
}
/*
额外的例子：
上面的循环是测试第几位被设置（占用，是1）
2. 设置指定位：
var a int8 = 8 //00001000
a = a | (1<<2)
fmt.Printf("%08b\n",a) // 00001100,设置了a的第3位。

3. 将第n位的值不设置
var a int8 = 13  //1101
a = a &^ (1<<2)
fmt.Printf("%04b\n",a)  //1001 将第3位的值不设置。
*/

fmt.Printf("%08b\n", x<<1) // "01000100", the set {2,6}
fmt.Printf("%08b\n", x>>1) // "00010001", the set {0,4}


(6.5节给出了一个可以远大于一个字节的整数集的实现)


14. <<左移 >>右移 之二（p86）
在 x<<n 和 x>>n 移位运算中，决定了移位操作bit数部分必须是无符号数；
被操作的x数可以是有符号或无符号数。 
算数上，一个 x<<n 左移运算等价于 x乘以2的n次方；一个 x>>n 右移运算，等价于x除以2的n次方。

左移运算用0填充右边空缺的bit位，无符号数的右移运算也是用0填充左边空缺的bit位，
但是，有符号数的右移运算会用符号位的值填充左边空缺的bit位。
因为这个原因，最好用无符号运算，这样你可以将整数完全当作一个bit位模式处理。

15. 选择无符号数，还是有符号数？（p87）
即使数值本身不可能出现复数，我们还是倾向于使用有符号的int类型。
如：内置len函数返回一个有符号的int
medals ：= []string{"gold", "silver", "bronze"}
for i := len(medals) - 1; i >= 0; i-- {
	fmt.Println(medals[i]) // "bronze", "silver", "gold"
}
如果 len 函数返回一个无符号数。。。（查看p87）


无符号数往往只有在位运算或其他特殊的运算场景才会使用，
就像bit集合、分析二进制文件格式或者是哈希和加密操作等。
他们通常并不用于仅仅是表达非负数的场合。

16. 算术和逻辑运算的二元操作中必须是相同的类型。(p87)
var apples int32 = 1
var oranges int16 = 2
var compote int = apples + oranges // compile error
// invalid operation: apples + oranges (mismatched types int32 and int16)

这种类型不匹配的问题可以有几种不同的方法修复，最常见方法是将他们都显式转型为一个常见类型：
如：
var compote = int(apples) + int(oranges) //将apples 和 oranges都转化为int型，然后进行运算。

如2.5节所述，对于每种类型T，如果转化允许的话，类型转换操作 T(x) 将x转换为T类型。

17. 转化类型可能改变数值，或丢失精度（p88）
许多整形数之间的转换并不会改变数值；他们只是告诉编译器如何解释这个值。
但是对于将一个大尺寸的整数类型转为一个小尺寸的整数类型，或者是将一个浮点数转为整数，可能会改变数值或丢失精度。
f := 3.131 // a float64
i := int(f)
fmt.Println(f, i) // "3.141 3"
f = 1.99
fmt.Println(int(f)) // "1"

浮点数到整数的转换将丢失任何小数部分，然后想数轴零方向截断。
你应该避免对可能会超出目标类型表示范围的数值类型转换，因为截断的行为可能依赖于具体的实现：
f := 1e100 // a float64
i := int(f) // 结果依赖于具体实现


18. 任何大小的整数字面值都可以用以0开始的八进制格式书写，例如0666；
或用以0x或0X开头的十六进制格式书写，例如0xdeadbeef 。十六进制数字可以用大写或小鞋字母。
如今八进制数据通常用于POSIX操作系统上的文件访问权限标志，十六进制数字则更强调数字值的bit位模式。

当使用fmt包打印一个数值时，我们可以用%d（十进制整数）、 %o（八进制）和%x（16进制）参数控制输出的进制格式，就像下面的例子：
o := 0666 //八进制表示 
fmt.Printf("%d %[1]o %#[1]o\n", o) // "438 666 0666"
x := int64(0xdeadbeef) //16进制表示
fmt.Printf("%d %[1]x %#[1]x %#[1]X\n ") // 3735928559 deadbeef oxdeadbeef oXDEADBEEF


19. 使用fmt的技巧 一（p88）
1 通常Printf格式化字符串包含多个 % 参数时将会包含对应相同数量的额外操作数，
但是%之后的[1]副词告诉Printf函数再次使用第一个操作数。
2 % 后的 # 副词告诉Printf在用%o、 %x 或 %X输出是产生0、0x或0X前缀。

20. 字符面值通过一对单引号直接包含对应字符。
最简单的例子是ASCII中类似'a'写法的字符面值，但是我们也可以通过转义的数值来表示任意的Unicode码点对应的字符。

字符使用%c参数打印，或者是用%q参数打印带单引号的字符;
ascii := 'a'
unicode := '国'
newline := '\n'
fmt.Printf("%d %[1]c %[1]q\n", ascii) // "97 a 'a'"
fmt.Printf("%d %[1]c %[1]q\n", unicode) // "22269 国 '国'"
fmt.Printf("%d %[1]q\n", newline) // "10 '\n'"
//这个符号对应的整数叫字符面值。unicode包含汉子等字符。ascii包含字母等。


3.2 浮点数
21. Go语言提供两种精度的浮点数：float32 float64 （p90）

22. 常量math.MaxFloat32表示float32的最大数值，大约 3.4e38;
常量math.MaxFloat64常量大约是 1.8e308。
它们分别能表示的最小值近似为 1.4e-45 和 4.9e-324

一个float32类似的浮点数可以提供大约6个十进制数的精度，而float64则可以提供约15个十进制数的精度。
通常应该优先使用float64类型，因为float32类型的累计计算误差很容易扩散，
并且float32能精确表示的正整数并不是很大（注：float32的有效bit位只有23个，其他的bit位用于指数和符号；当整数大于23bit能表达的范围时，float32的表示将出现误差）：
var f float32 = 16777216 //1<<24
fmt.Println(f == f+1) //"true"!


23. 浮点数的字面值可以直接写小数部分，如下：
const e = 2.71828 // (approximately)
小数点前面或后面的数字都可能被省略。很小或很大的数最好用科学计数法书写，通过e或E来指定指数部分：
const Avogadro = 6.02214129e23 // 阿伏伽德罗常数
const Planck = 6.62606957e-34 // 普朗克常数

24. 用Printf函数的%g参数打印浮点数，将采用更紧凑的表示形式打印，并提供足够的精度，但是对应表格的数据，使用 %e (带指数)或 %f 的形式打印可能更合适。
所有的这三种打印形式都可以指定打印的宽度和控制打印精度。
for x := 0; x < 8; x++{
	fmt.Printf("x = %d e^x = %8.3f\n", x, math.Exp(float64(x)))
}
代码打印e的幂，打印精度是小数点后三个小数，精度8个字符宽度;
x = 0 e^x = 1.000
x = 1 e^x = 2.718
...
x = 7 e^x = 1096.633 //8个字符宽度

25. math包（p91）
math包中除了提供大量常用的数学函数外，还提供了IEEE754浮点数标准中定义的特殊值的创建和测试：
正无穷大和负无穷大，分别用于表示太大溢出的数则和除0的结果；
还有NaN非数，一般用于表示无效的书法操作结果 0/0 或 Sqrt(-1) //Sqrt开方
var z float64
fmt.Println(z, -z, 1/z, -1/z, z/z) // "0 -0 +Inf(正无穷大) -Inf（负无穷大） NaN（非数）"

26. math.IsNaN（p91）
math.IsNaN用于测试一个数是否为非数NaN，math.NaN则返回非数对应的值。
虽然可以用math.NaN来表示一个非法的结果，但是测试一个结果是否为非数NaN则是充满风险的，
因为NaN和任何数都是不相等的（注：在浮点数中，NaN、正无穷大和负无穷大都不是唯一的，每个都有非常多种的bit模式表示）：
nan := math.NaN()
fmt.Println(nan) // "NaN"
fmt.Println(nan == nan, nan < nan, nan > nan) //"false false false"

27. 如果一个函数返回的浮点数结果可能失败，最好的做法是用单独的标志报告失败，如：
func compute() (value float64, ok bool){
	// ...
	if failed{
		return 0, false
	}
	return result, true
}

28. 可缩放矢量图形（SVG）格式输出，SVG是一个用于矢量线绘制的XML标准。(p91)

3.3 复数

29. Go语言提供了两种精度的复数类型：complex64 和 complex128， （p95）
分别对应 float32 和 float64两种浮点数精度。
内置的complex函数用于构建复数，内建的real和imag函数分别返回复数的实部和虚部：
var x complex128 = complex(1, 2) // 1+2i
var y complex128 = complex(3, 4) // 3+4i
fmt.Println(x*y) //"(-5+10i)" (1+2i)*(3+4i)
fmt.Println(real(x*y)) //"-5"
fmt.Println(imag(x*y)) //"10"

30. 如果一个浮点数面值或一个十进制整数面值后面跟着一个i，如 3.141592i 或 2i，它将构成一个复数的虚部，复数的实部是0：
fmt.Println(1i * 1i ) //"-1+0i" , i^2 = -1

31. 命名复数可以简化
如：
x := 1 + 2i

32. 复数可以== 和 != 进行比较。（p95）
只有两个复数的实部和虚部相等的时候他们才是相等的。
（浮点数的相等比较是危险的，需要特别小心处理精度问题） 

33. math/cmplx包提供了复数处理的许多函数，例如求复数的平方根函数和求幂函数。
fmt.Println(complx.Sqrt(-1)) // "(0+1i)"


3.4 布尔型（p98）

34. 一个布尔型的值只有两种：true 和 false。

35. if和for语句的条件部分都是布尔类型的值，并且== 和 <等比较错做也会产生布尔型的值。

36. 一元操作符 ! 对应逻辑非操作，因此 !true 的值为 false。
简洁的布尔表达式 x == true

37. 布尔值可以和&&（and）和||（or）操作符结合，并且可能会有短路行为：
如果运算符左边值已经可以确定整个布尔表达式的值，那么运算符右边的值将不在被求值，
因此下面的表达式总数安全的：
s != "" && s[0] =="x" //不再被求值是什么意思？？？没理解。
//s就是空的，所以 s != ""为false，左边可以确定了，就不用管右边的了？？？
其中s[0]操作如果应用于空字符串将会导致panic异常。

38. 因为&&的优先级比||高。
（助记：&&对应乘，||对应加，乘法优先级高于加法）
下面形式的布尔表达式是不需要加小括号的：
if 'a'<=c && c <= 'z' || //'a' 单引号是ASCII码
	'A'<=c && c <= 'Z' ||
	'0'<=c && c <= '9'{
	//...ASCII letter or digit....
}

39. 布尔值不会隐式转换为数字值0或1，反之亦然。
必须使用一个显式的if语句辅助转换：
i:=0
if b{//b是bool型，如果b为true，i=1，如果不是。。。
	i=1
}
//这是干了个啥？转换啥了？

如果需要经常做类似的转换，包装成一个函数会更方便：
//btoi return 1 if b is true and 0 if false.
func btoi (b bool)int {
	if b{
		return 1
	}
	return 0
}

40. 数字到布尔型的逆转换则非常简单，
不过为了保持对称，我们也可以包装一个函数：（p99）
// itob reports whether i is non-zero.
func itob(i int) bool{return i != 0}//只要不是0，bool值就是true


3.5 字符串（p100）

41. 字符串
一个字符串是一个不可改变的字节序列。字符串可以包含任意的数据，包括byte值0，
但是通常是用来包含人类可读的文本。
文本字符串通常被解释为采用UTF8编码的Unicode码点（rune）序列。

42. 内置的len函数可以返回一个字符串中的字节数目（不是rune字符数目），
索引操作 s[i] 返回第i个字节的字节值，i必须满足 0<=i<len(s)条件约束。
s := "hello, world"
fmt.Println(len(s)) //12
fmt.Println(s[0], s[7]) // 104 119 ('h' and 'w')
如果试图访问超出字符串索引范围的字节将会导致panic异常：
c := s[len(s)] // panic: index out of range

43. 第i个字节并不一定是字符串的第i个字符，因为对于非ASCII字符的UTF8编码会有两个或多个字节。
（没懂？如中文，一个中文（具体看是什么码）是两个字节 或3-4个字节？第1个中文和第一个字节不是同一个东西。是这个意思吧？）
子字符串操作s[i:j]基于原始的s字符串的第i个字节开始到第j个自己（并不包含j本身）生成一个新字符串。
生成的新字符串将包含j-i个字节。
fmt.Prinln(s[0:5]) // "hello"
同样，如果索引超出字符串范围或者j小于i的话将导致panic异常。
不管i还是j都可能被忽略，当他们被忽略是将采用0最为开始位置，采用 len(s)作为结束的位置。
fmt.Println(s[:5]) // "hello"
fmt.Println(s[7:]) // "world"
fmt.Println(s[:])  // "hello world"

44. 其中+操作符将两个字符串链接构成一个新字符串：
fmt.Println("goodbye"+s[5:]) // "goodbye, world"

45. 字符串可以用== 和 < 进行比较；比较通过逐个字节比较完成的，因此比较的结果是字符串自然编码的顺序。

46. 字符串的值是不可变的：一个字符串包含的字节序列永远不会被改变，当然我们也可以给一个字符串变量分配一个新字符串值。
s := "left foot"
t := s
s +=", right foot"
这并不会导致原始的字符串值被改变，但是变量s将因为+=语句持有一个新的字符串值，
但是t依然是包含原先的字符串值。
fmt.Println(s) // "left foot, right foot"
fmt.Println(t) // "left foot"

47. 因为字符串是不可修改的，因此尝试修改字符串内部数据的操作也是被禁止的：
s[0]="L" //compile error: cannot assign to s[0].
不变性意味如果两个字符串共享相同的底层数据的话也是安全的，这使得复制任何长度的字符串代价是低廉的。
同样，一个字符串s和对应的子字符串切片s[7:]的操作也可以安全地共享相同的内存，因此字符串切片操作代价也是低廉的。
在这两种情况下都没有必要分配新的内存。

3.5.1 字符串面值（p101）
48. 字符串值也可以用字符串面值方式编写，只要将一系列字节序列包含在双引号即可：
"Hello, 世界"

因为Go语言源文件总是用UTF-8编码，并且Go语言的文本字符串也以UTF8编码的方式处理，
因此我们可以将Unicode码点也写到字符串面值中。
（Unicode字符的编码方式一般有三种：UTF-8，UTF-16，UTF-32）

在一个双引号包含的字符串面值中，可以用以反斜杠 \ 开头的转译序列插入任意的数据。
下面的换行、回车和制表符等是常见的ASCII控制代码的转义方式：

\a 响铃 
\b 退格
\f 换页
\n 换行
\r 回车
\t 制表符
\v 垂直制表符
\' 单引号（只用在 '\'' 形式的rune符号面值中 ）
\" 双引号（只用在 "..." 形式的字符串面值中）"
\\ 反斜杠

49. 可以通过十六进制或八进制转义在字符串面值包含任意的字节。 
十六进制的转义形式是： \xhh,	其中两个h表示十六进制数字（大小写都可以）。
八进制转义形式是： \ooo, 包含三个八进制的o数字（0到7），但是不能超过 \377(对应一个字节的范围，十进制为255)
每一个单一的字节表达一个特定的值。

50. 原生的字符串面值（p102）
原生字符串面值形式是 `...` ，使用反引号，代替双引号。
在原生的字符串面值中，没哟转义操作；
全部内容为字面意思，包含退格和换行，因此一个而程序中的原生字符串面值可能跨越多行。
（在原生字符串面值内部是无法直接写字符的，可以用八进制或十六进制转义或+"```"链接字符串 常量完成）
唯一的特殊处理是会删除回车以保证在所有平台上的值都是一样的，包括那些把回车也放入文本文件的系统（windows系统会把回车和换行一起放入文本文件中）。

原生字符串面值用于编写正则表达式会很方便，因为正则表达式往往会包含很多反斜杠。原生字符串面值同时被广泛用于HTML模板、JSON面值、命令行提示信息以及那些需要扩展到多行的场景。

const GoUsage = `Go is a tool for managing Go source code.

Usage:
	go command [arguments]
...`

3.5.2 Unicode(p102)

51. ASCII字符集：美国信息标准代码。
使用7bit来表示128个字符：包含英文字母的大小写、数字、各种标点符号和设置控制符。

Unicode（ http://unicode.org ）， 
它收集了世界上所有的符号系统，包括重音符号和其他变音符号，制表符和回车符，还有很多神秘的符号，每个符号都分配一个唯一的Unicode码点，
Unicode码点对应Go语言中发rune整数类型。
（rune是int32等价类型）

3.5.3 UTF-8

52. UTF-8是一个将Unicode码点编码为字节序列的变长编码。
UTF8编码使用1到4个字节来表示每个Unicode码点，
ASCII部分字符只使用1个字节，常用字符部分使用2或3个字节表示。
每个符号编码后第一个字节的高端bit位用与表示总共有多少个字节。

如果第一个字节的高端bit为0，则表示对应7bit的ASCII字符，ASCII字符每个字符依然是一个字节，和传统的ASCII编码兼容。
如果第一个字节的高端bit是110，则说明需要2个自己；
后续的每个高端都以10开头。
更大的Unicode码点也是采用类似的策略处理。
0xxxxxxx rune 0-127 (ASCII)
110xxxxx 10xxxxxx 128-2047 (values <128 unused)
1110xxxx 10xxxxxx 10xxxxxx 2048-65535 (values <2048 unused)
11110xxx 10xxxxxx 10xxxxxx 10xxxxxx 65536-0x10ffff (other values unused)

53. utf-8编码的好处（详细p103）

54. GO语言中utf-8的优势（p104）
除了优秀，unicode包提供了诸多处理rune字符相关功能的函数（比如区分字母和数组，或者是字母的大写和小写转换等），
unicode/utf8包则提供了用于rune字符序列的UTF8编码和解码功能。

55. 很多Unicode字符很难直接从键盘输入，并且还有很多字符有着相似的结构；
甚至是不可见的字符。
Go语言字符串面值中的Unicode转义字符让我们可以通过Unicode码点输入特殊的字符。
有两种形式：、\uhhhh对应16bit的码点值，\Uhhhhhhhh对应32bit的码点值，其中h是一个十六进制数字；
一般很少需要使用32bit的形式。
下面的字符串面值都表示相同的值：
"世界"
"\xe4\xb8\x96\xe7\x95\x8c"
"\u4e16\u754c"
"\U00004e16\U0000754c"
上面三个转义序列都为第一个字符串提供代替写法，但是他们的值都是相同的。
Unicode转义也可以使用在 rune字符 （'...'）中。下面三个私服是等价的：
'世' '\u4e16' '\U00004e16'

56. 对于小于256码点值可以写在一个十六进制转义字节中，
例如：'\x41'对应字符'A'，
但是对于更大的码点则必须使用\u或\U转义形式。
因此，'\xe4\xb8\x96'并不是一个合法的rune字符（？？？），
虽然这三个字节对应一个有效的UTF8编码的码点。
(仔细琢磨这里说的)

57. 得益于UTF8编码优良的设计，诸多字符串操作都不需要解码操作。我们可以不用解码直接测试一个字符串是否是另一个字符串的前缀：
前缀测试：
func HasPrefic(s, prefix string) bool{
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

后缀测试：
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

包含子串测试:
func Contains(s, substr string) bool{
	for i := 0; i < len(s); i++ {
		if HasPrefix(s[i:], substr) {
			return true
		}
	}
	return false
}

58. （p105）字符串包含13个字节，以UTF8形式编码，但是只对应9个Unicode字符：
import"unicode/utf8"

s := "Hello, 世界"
fmt.Println(len(s)) // "13" 包含了13个字节
fmt.Println(utf8.RuneCountInString(s)) //"9" 对应9个Unicode字符

为了处理这些真实的字符，我们需要一个UTF8解码器。
unicode/utf8包提供了该功能，我们可以这样使用：
for i := 0; i < len(s); {
	r, size := utf8.DecodeRuneInString(s[i:])//s[i:]这里的i代表索引的位置。
	fmt.Printf("%d\t%c\n", i, r)//%c 字符（rune） (Unicode码点)
	i += size
}
每一次调用DecodeRuneInString函数都返回一个r和长度，r对应字符本身，长度对应r采用UTF8编码后的编码字节数。
长度可以用于更新第i个字符在字符串中的字节索引位置。


59. Go 语言的range循环在处理字符串的时候，会自动隐式解码UTF8字符串。(p106)
for i, r := range "Hello, 世界"{
	fmt.Printf("%d\t%q\t%d\n", i, r, r)// /t是制表符（见Go1:13（p30））
}

60. 统计字符串 字符 的数目。（p106）
我们可以使用一个简单的循环来统计字符串中字符的数目，像这样：
n := 0
for _, _ = range s{
	n++
}
像其他形式的循环那样，我们也可以忽略不需要的变量：
n := 0
for range s {
	n++
}

或者我们可以直接调用utf8.RuneCountInString(s)函数。

61. UTF8字符串作为交换格式是非常方便的，但是在程序内部采用rune序列可能更方便。
因为rune大小一致，支持数组索引和方便切割。
string 接受到 []rune的类型转换，可以将一个UTF8编码的字符串解码为Unicode字符串序列：

// "program" in Japanese katakana
s := "プログラム"
fmt.Printf("% x\n", s) // e3 83 97 e3 83 ad e3 82 b0 e3 83 a9 e3 83 a0
// % x用于在每个十六进制数字前插入一个空格。
r := []rune(s)
fmt.Printf("%x\n", r) // "[30d7 30ed 30b0 30e9 30e0]"

如果是将一个[]rune类型的Unicode字符slice或数组转为string，则对他们进行UTF8编码：

fmt.Println(string(r)) //"プログラム"

62. 将一个整数转为字符串意思是生成以只包换对饮gUnicode码点字符的UTF8字符串：(p107)
fmt.Println(string(65)) // "A", not "65"
fmt.Println(string(0x4eac)) // "京"

如果对应码点的字符是无效的，则用'\uFFFD'无效字符作为替换：
fmt.Println(string(1234567)) // "�"

3.5.4 字符串和Byte切片（p107）

63. 标准库中有四个包对字符串处理尤为重要：bytes\ strings\ strconv\ uncode包。
strings包提供了许多如字符串的查询、替换、比较、截断、拆分和合并等功能。

bytes包也提供了很多类似功能的函数，但是针对和字符串有着相同结构的[]byte类型。
因为字符串是只读的，因此逐步构建字符串会导致很多分配和复制。
在这种情况下，使用bytes.Buffer类型将会更有效，稍后我们将展示。

strconv包提供了布尔型、整行数、浮点数和对应字符串的相互转换，还提供了双引号转义相关的转换。

unicode包提供了IsDigit、IsLetter、IsUpper和IsLower等类似功能，用于给字符分类。
每个函数有一个单一的rune类型的参数，然后返回一个布尔值。
而像ToUpper和ToLower之类的转换函数将用于rune字符的带下写转换。
所有的这些函数都是遵循Unicode标准定义的字母、数字等分类规范。

strings包也有类似的函数，他
们是ToUpper和Tolower，将原始字符串的每个字符都做响应的转换，然后返回新的字符串。

64. strings.LastIndex (p108)
下面例子的basename函数灵感于Unix shell的同名工具。
在我们现实的版本中，basename(s)将看来像是系统路径的前缀删除，同时将看似文件类型的后缀名部分删除：

fmt.Println(basename("a/b/c.go")) // "c"
fmt.Println(basename("c.d.go")) // "c.d"
fmt.Prinrln(basename("abc")) // "abc"

65. path和path/filepath 包提供了关于文件路劲名更一般的函数操作。
使用斜杠分隔路劲可以在任何操作系统上工作。
斜杠本身不应该用于文件名，但是在其他一些领域可能会用于文件名。
例如URL路径组件。相比之下，path/filepath包则使用操作系统本身的路径规则，
例如POSIX系统使用/foo/bar，而microsoft Windows使用c:\foo\bar等。

函数的功能是将一个表示整值的字符串，每隔三个字符插入一个逗号分隔符，例如“12345”处理后成为“12，345”。
这个版本只适用于整数类型；


66. 一个字符串是包含的只读字节数组，一旦创建，是不可变的。
相比之下，一个字节slice的元素则可以自由地修改。

67. 字符串和字节slice之间可以相互转换：
s := "abc"
b := []byte(s) //将s字符串转化为byte类型的切片
s2 := string(b)

68. 从概念上讲，一个[]byte(s)转换是分配了一个新的字节数组用于保存字符串数据的拷贝，
然后引用这个底层的字节数组。
编译器的优化可以避免在一些场景下分配和复制字符串数据，
但总的来说需要确保在变量b被修改的情况下，原始的s字符串也不会改变。
讲一个字节slice转到字符串的 string(b)操作则是构造一个字符串拷贝，以确保s2字符传是只读的。

为了避免转换中不必要的内存分配，bytes包和strings同时提供了许多使用函数。
下面是strings包中的六个函数：
func Contains(s, substr string) bool
func Count(s, sep string) int
func Fields(s string) []string
func HasPrefix(s, prefix string) bool
func Index(s, sep string) int
func Join(a []string, sep string) string

bytes包中也对应的六个函数：
func Contains(b, subslice []byte) bool
func Count(s, sep []byte) int
func Fields(s []byte) [][]byte
func HasPrefix(s, prefix []byte) bool
func Index(s, sep []byte) int
func Join(s [][]byte, sep []byte) []byte

它们之间唯一的区别是字符串类型参数被替换成了字节slice类型的参数。

69. bytes包还提供了Buffer类型用于字节slice的缓存。(p110)
一个Buffer开始是空的，但是随着string、byte或[]byte等类型数据的写入可以动态增长，
一个bytes.Buffer变量并不需要处理化，因为零值也是有效的。
（有示例）

70. 当向bytes.Buffer添加任意字符的UTF8编码时，最好使用bytes.Buffer的WriteRune方法，
但是WirteByte方法对于写入类似'['和']'等ASCII字符则会更加有效。

71. bytes.Buffer类型有着很多使用的功能，我们在第七章讨论接口是会涉及到，（p111）
我们将看看如何将他用作一个I/O的输入和输出对象，例如当做Fprintf的io.Writer输出对象，
或者当做io.Reader类型的输入源对象。

3.5.5 字符串和数字的转换

72. 处理字符串、字符、字节之间的转换，字符串和数值之间的转换也比较常见。
由strconv包提供这类转换功能。
将一个整数转换为字符串：
一种方法是用fmt.Sprintf返回一个格式化的字符串；
另一个方法是用strconv.Itoa ("整数到ASCII")：
x := 123
y := fmt.Sprintf("%d", x)
//Sprintf:将数据按固定格式输出到某个缓存中，可以这么理解，其实它功能和printf差不读，只是printf输出到了屏幕上，而sprintf只是将同样的内容装到了某个缓存中。
fmt.Println(y, strconv.Itoa(x))// "123 123"

73. Formatlnt 和 FormatUint函数可以用不同的进制来格式化数字：
fmt.Println(strconv.FormatInt(int64(x),2)) //"1110011"

74. fmt.Printf函数的%b %d %o %x等参数提供功能往往比strconv包的Format函数方便很多，
特别是在需要包含附加额外信息的时候：
s := fmt.Sprintf("x=%b", x) // "x=1111011"

75. 如果要将一个字符串解析为整数，可以使用strconv包的Atoi或Parselnt函数，
还有解析无符号整数的ParseUint函数;
x, err := strconv.Atoi("123") // x is an int
y, err := strconv.ParseInt("123", 10, 64) //base 10, up to 64 bits(10进制，最多64位？)
ParsenInt函数的第三个参数是用于指定整型数的大小；例如16表示int16,0则表示int。
在任何情况下，返回的结果y勇士int64，你可以通过强制类型转换将他转为更小的整数类型。


76. 有时候也会使用fmt.Scanf来解析输入的字符串和数字，特别是当字符串和数字混合在一行的时候，
它可以灵活处理不完整或不规则的输入。

3.6 常量
77. 常量表达式的值在编译期计算，而不是在运行期。
每种常量的潜在类型都是基础类型：boolean（布尔型）、string或数字

一个常量的声明语句定义了常量的名字、和变量的声明语法类似，常量的值不可修改。
const pi = 3.14159265359 //这是近似值，使用math.Pi是更好的近似值。
可以批量声明多个常量:
const(
	e = 2.718281828459045235360287471
	pi = 3.1515926535897932384626433832
)

常量的运算都可以在编译期完成，这样减少运行时的工作，也方便其他编译优化，也有助于提前发现操作数是常数时的错误。

78. 常量间的所有算术运算、逻辑运算和比较运算的结果也是常量，(p113)
对常量的类型转换操作或以下函数调用都是返回常量结果：
len \ cap \ real \ imag \ complex 和 unsafe.Sizeof（13.1）

79. 因为常量的值是在编译期就确定的，因此常量可以是构成类型的一部分，例如用于指定数组类型的长度：
const IPv4Len = 4

//parseIPv4 parses an IPv4 address (d.d.d.d).(解析一个ipv4地址)
func parseIPv4(s string) IP {
	var p [IPv4Len]byte
	//...
}

80. 一个常量的声明也可以包含一个类型和一个值，
但是如果没有显式指明类型，那么将从右边的表达式推断类型。
time.Duration是一个命名类型，底层类型是int64，
time.Minute是对应类型的常量。
下面声明的两个常量都是time.Duration类型，
可以通过%T参数打印类型信息：
const noDelay time.Duration = 0 //time.Duration是一个命名类型，底层类型是int64
const timeout = 5 * time.Minute
fmt.Printf("%T %[1]v\n", noDelay) // "time.Duration 0"
fmt.Printf("%T %[1]v\n", timeout) // "time.Duration 5m0s"
fmt.Printf("%T %[1]v\n", time.Minute) // "time.Duration 1m0s"
//%[1]v 中的[1]表示再次使用第一个操作数

81. 如果是批量声明的常量，除了第一个外其他的常量右边的初始化表达式都可以省略，
如果省略初始化表达式则表示使用前面常量的初始化表达式写法，对应的常量类型也一样的。
const (
	a = 1
	b
	c = 2
	d
)
fmt.Println(a, b, c, d) // "1 1 2 2"

3.6.1 iota常量生成器（p114）
82. 常量声明可以使用iota常量生成器初始化，它用于生成一组以相似规则初始化的常量，
但是不用每行都写一遍初始化表达式。
在一个const声明语句中，在第一个声明的常量所在的行，
iota将会被置为0，然后在每一个有常量声明的行加1.package main


83. 来自time包的例子（p114）
它首先定义了一个Weekday命名类型，然后为一周的每天定义了一个常量，
从周日0开始。在其他编程语言中，这种类型一般被成为枚举类型。
type Weekday int

const (
	Sunday Weekday = iota //0
	Monday //1
	Tuesday //2
	Wednesday //3
	Thursday //4
	Friday //5
	Saturday //6
)

84. 我们也可以在复杂的常量表达式中使用iota， 下面是来自net包的例子，
用于给一个无符号整数的最低5bit的每个bit指定一个名字：
byte Flags uint

const（
	FlagUp Flags = 1 << iota // is up
	FlagBroadcast // supports broadcast access capability
	FlagLoopback // is a loopback interface
	FlagPointToPoint // belongs to a point-to-point link
	FlagMulticast // supports multicast access capability
）
随着iota的递增，每个常量对应表达式1<<iota，是连续的2的幂，分别对应一个bit位置。
是用这个常量可以用于测试、设置或清楚对应的bit位的值：

85. 下面是一个更复杂的例子，每个常量都是1024的幂
const (
	_ = 1 << (10 * iota)
	KiB // 1024  2^10
	MiB // 1048576 2^20
	GiB // 1073741824 2^30
	TiB // 1099511627776 2^40
	PiB // 1125899906842624 2^50
	EiB // 1152921504606846976 2^60
	ZiB // 1180591620717411303424 2^70
	YiB // 1208925819614629174706176 2^80
)

86. 不过iota常量生成规则也有其局限性。
例如，它不能用于产生1000的幂（KB，MB等），因为Go语言并没有计算幂的运算符。

3.6.2 无类型常量 
87.1 许多常量并没有一个明确的基础类型。（var f float64 = 212 其中212就是将212这个没有基础类型的无类型常量，赋值为float64类型。）
87.2. 六种未明确类型的常量类型，(p116)
分别是：无类型的布尔型、无类型的整形、无类型的字符、
无类型的浮点数、无类型的复数、无类型的字符串。

通过延迟明确常量的具体类型，无类型的常量不仅可以提供更高的运算精度，而且可以直接用于更多的表达式而不需要显式的类型转换。
上面例子中的ZiB和YiB的值已经超出任何Go语言中整数类型能表达的范围，
但是它们依然是合法的常量，而且可以想下面常量表达式依然有效。
（注：YiB/ZiB是在编译期计算出来的，并且结果常量是1024，是Go语言int变量能有效表示的）：
fmt.Println(YiB/ZiB) //"1024"

88. math.Pi (p116)
另一个例子，math.Pi无类型的浮点数常量，可以直接用于任意需要浮点数或复数的地方：
var x float32 = math.Pi
var y float64 = math.Pi
var z complex128 = math.Pi

如果math.Pi被确定为特定类型，比如float64，那么结果精度可能会不一样，
同时对于需要float32或complex28类型值的地方则会强制需要一个明确的类型转换：
const Pi64 float64 = math.Pi

var x float32 = float32(Pi64)
var y float64 = Pi64
var z complex128 = complex128(Pi64)

89. 对于常量面值，不同的写法可能会对应不同的类型。
例如0、0.0、0i和'\u0000'虽然有着相同的常量值，
但是他们分别对应无类型的整数、无类型的浮点数、无类型的复数和无类型的字符等不同的常量类型。
同样，true和false也是无类型的布尔类型，
字符串面值常量是无类型的字符串类型。

90. 前面说过除法运算符/会根据操作数的类型生成对应类型的结果。
因此，不同写法的常量除法表达式可能对应不同的结果：
var f float64 = 212
fmt.Println((f - 32) *5 / 9) // "100"; (f-32)*5 is a float
fmt.Println(5 / 9 * (f - 32)) // "0"; 5/9 is an Untyped integer, 0
fmt.Println(5.0 / 9.0 * (f - 32)) // "100"; 5.0/9.0 is an untyped float

91. 只有常量可以是无类型的。
当一个无类型的常量被赋值给一个变量的时候，就像上面的第一行语句，或者是像其余三个语句中右边表达是中含有明确类型的值，
无类型的常量将会被隐式转换为对应的类型，如果转换合法的话。
var f float64 = 3 + 0i //untyped complex -> float64
f = 2 // untype integer -> float64
f = 1e123 // untype floating-point -> float64
f = 'a' // untype rune -> float64
上面的语句相当于:
var f float64 = float64(3+0i)
f = float64(2)
f = float64(1e123)
f = float64('a')

92. 无论是隐式或显示转换，将一种类型转换为另一种类型都要求目标可以表示原始值。
对于浮点数和复数，可能会有舍入处理。
const(
	deadbeef = 0xdeadbeef // untyped int with value 3735928559
	a = uint32(deadbeef) // uint32 with value 3735928559
	b = float32(deadbeef) // float32 with value 3735928576 (rounded up四舍五入？)
	c = float64(deadbeef) // float64 with value 3735928559(exact)
	d = int32(deadbeef) //compile error: constant overflows int32(溢出)
	e = float64(1e309) //compile error: constant overflows float64
	f = uint(-1) // compile error: constant underflows uint(常量下溢)
)

93. 对于一个没有显示类型的变量声明语法（包括短变量声明语法），无类型的常量会被隐式转为默认的变量类型，就像下面的例子:（p117）
i := 0 // untyped integer; implicit int(0) 
r := '\000' // untype rune; implicit rune('\000')
f := 0.0 // untyped floating-point; implicit float64(0.0)
c := 0i // untype complex; implicit complex128(0i)

94.（p118） 默认类型是规则的：无类型的整数常量默认转换为int，对于不确定的内存大小，但是浮点数和复数常量则默认转换为float64 和complex128。
Go语言本身并没有不确定内存大小的浮点数和复数类型，而且如果不知道浮点数类型的话将很难写出正确的数值算法。

如果要给变量一个不同的类型，我们必须显示地将无类型的常量转化为所需的类型，或给声明的变量指定明确的类型，像下面的例子：
var i = int8(0)
var i int8 = 0

95. 当尝试将这些无类型的常量转为一个接口值时（见第7章），这些默认类型将显得尤为重要，（p118）
因为要靠他们明确接口对应的动态类型。
fmt.Printf("%T\n", 0) // "int"
fmt.Printf("%T\n", 0.0) // "float64"
fmt.Printf("%T\n", 0i) // "complex128"
fmt.Printf("%T\n", '\000') // "int32" (rune)

