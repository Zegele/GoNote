第一章 入门

1. go run (p19)
这个命令编译一个或多个以.go结尾的源文件，链接（源文件涉及的）库文件，并运行最终生成的可执行文件。

2. go build(p19)
go build 后跟文件。如：go build aaa.go 
//这条命令路径必须在该文件夹内，并且在该路径生成一个aaa.exe文件
go build 后跟文件夹。 如：go build ssss
//这条命令可以在任何路径下生成一个ssss.exe文件。但引用的ssss文件夹必须在GOPATH的src文件夹下，或在GOROOT的src文件夹下。
在命令终端
$ ./ssss.exe则可运行该.exe文件（exe可运行文件） 

3. os包（p23）
os包以跨平台的方式，提供了一些与操作系统交互的函数和变量。
程序的命令行参数可从os包的Args变量获取；
//怎样操作呢？演示一下？
os包外部使用os.Args访问该变量。
//

4. os.Args的第一个元素，os.Args[0], 命令本身的名字。其它的元素则是程序启动时传给它的参数。(p23)

5. s += sep + os.Args[i]  等价于：s = s + sep + os.Args[i]  (p24)

6. for循环（p25）
for initialization; condition; post{
    //zero or more statements 
}
//
initialization语句是可选的，在循环开始前执行。
initalization如果存在，必须是一条简单语句（simple statement），即，短变量声明、自增语句、赋值语句或函数调用。 
condition 是一个布尔表达式（boolean expression），其值在每次循环迭代开始时计算。如果为 true 则执行循环体语句。 
post 语句在循环体执行结束后执行，之后再次对conditon 求值。 condition 值为 false 时，循环结束。
for循环的这三个部分每个都可以省略。

7. range (p26)
range 产生一对值，索引和在该索引处的元素值。
可以将忽略的值用 空标识符（blank identifier）,即 _ (就是下划线)。

for _, arg := range os.Args[1:]{
	s += sep + arg
	sep = " "
}

8. s := "" //短变量声明，最简洁，但只能用在函数内部。
   var s string //依赖于字符串的默认初始化零值机制，被初始化为"" (string的空值)。
   var s = "" //用的很少，除非同时声明多个变量。
   var s string = "" //不常用，类型冗余。

9. if语句 (p28)
if语句中条件两边不加括号，但主体部分要加。if语句的else部分是可选的，在if的条件为false是执行。
if n>1{
	fmt.Printf("%d\t%s\n", n, line)
	} else{
		fmt.Println("NIUBI")
	}

10. map(p28)
map储存了键/值（key/value）的集合，对集合元素，提供常数时间的存、取或测试操作。
键 可以是任意类型，只要其能用 == 运算符比较就行，最常见的例子是字符串；
值 可以是任意类型。
内置函数make创建空map，此外还有别的作用（4.3节讨论）
counts := make(map[string]int)
键是字符串，值是整数。

11. bufio包（p29） 
它使处理输入和输出方便又高效。Scanner类型是该包最有用的特性之一，它读取输入并将其拆成行或单词；通常是处理行形式的输入最简单的方法。

input := bufio.NewScanner(os.Stdin) //os.Stdin表示标准输入
//使用短变量声明创建bufio.Scanner类型的变量input。

该变量从程序的标准输入中读取内容。
每次调用input.Scanner,即读入下一行，并移除行末的换行符；
读取的内容可以调用input.Text()得到。
Scan函数在读到一行时返回true，在无输入时返回false。

12. printf函数（p29）
fmt.Printf 函数对一些表达式产生格式化输出。该函数的首个参数是个格式字符串，制定后续参数被如何格式化。
各个参数的格式取决于“转换字符”（conversion character），形式为百分号后跟一个字母。如%d表示以十进制形式打印一个整型操作数，%s表示把字符串型操作数的值打出来（展开）。

更多如下表：
%d 十进制整数
%x, %o, %b 十六进制，八进制，二进制整数。
%f, %g, %e 浮点数： 3.141593 3.141592653589793 3.141593e+00
%t 布尔：true或false
%c 字符（rune） (Unicode码点)
%s 字符串
%q 带双引号的字符串"abc"或带单引号的字符'c'
%v 变量的自然形式（natural format）
%T 变量的类型
%% 字面上的百分号标志（无操作数）
%p 打印内存地址//千峰教育加的

13. 不可见字符的转译字符（escap sequences）（p30）
\t 制表符
\n 换行符

14. 很多程序要么从 标准输入 中读取数据，要么从一系列具名文件中读取数据。

15. os.Open 函数返回两个值。(p31)
f, err := os.Open(arg)
if err != nil{
	fmt.Fprintf(os.Stderr, "dup2:%v\n", err)
	continue
}

返回的第一个值是被打开的文件（*os.File），其后被Scanner读取。
返回的第二个值是内置error类型的值。如果err等于内置值nil，那么文件被成功打开。读取文件，知道文件结束，然后调用close关闭该文件，并释放占用的所有资源。
如果err的值不是nil，说明打开文件时出错了。这种情况下，错误值描述了所遇到的问题。

16. contunue(p31)
continue语句直接跳到for循环的下个迭代开始执行。

17. 参数传递（p32）
以map为例，map作为参数传递给某函数时，该函数接收这个引用的一份拷贝，被调用函数对map底层数据结构的任何修改，调用这函数都可以通过持有的map引用看到。
（意思是说，当map的底层数据被修改，调用这个map时（已修改的）看到的是被修改后的map。）

18. ReadFile函数（p32）
来自io/ioutil包，其读取指定文件的全部内容。

19. strings.Split函数把字符串分割成子串的切片。（p32）
strings.Split的作用与strings.Join相反。
strings.Join 大概是合并字符串吧。

20. 引用（import）一个包路径包含有多个单词的包（package）时，（p36）
如image/color，通常我们只需要最后那个单词表示这个就可以。所以color.White时，这个变量指向的是image/color包里的变量。

21. struct内部的变量可以以一个点（.）来进行访问。（p36）

22. http.Get（net/http包）函数是创建HTTP请求的函数，如果获得过程没有出错，那么会在resp这个结构体中得到访问的请求结果。
resp的Body字段包括一个可读的服务器响应流。（p38）

23. b, err := ioutil.ReadAll(resp.Body)  (p38)
ioutil.ReadAll函数从response中读取到全部内容,将其结果保存在变量b中。

24. resp.Body.Close 关闭resp的Body流，防止资源泄露。

25. os.Exit(1) 
无论实例代码中出现哪种失败原有，我们的程序都用了os.Exit函数来终止进程，并且返回一个status错误码，其值为1。(p39)

26. goroutine是一种函数的并发执行方式，而channel是用来在goroutine之间进行参数传递。

27. 有关goroutine和channel最出的例子和知识，请联系实例代码一同对照、查看、理解。（p42）

28. Web服务（p43）
Web程序服务中，标准库里的方法已经帮我们完成了大量工作。
main函数将所有发送到 /路径下的请求和handler函数关联起来，/开头的请求其实就是所有发送到当前站点的请求，服务监听8000端口。
发送到这个服务的“请求”是一个http.Request类型的对象，这个对象中包含了请求中的一系列相关字段，其中就包括我们需要的URL。
当请求到达服务器时，这个请求会被传给handler函数来处理，这个函数会将/hello这个路径从请求的URL中解析出来，然后把其发送到响应中，这里我们用的是标准输出流的fmt.Fprintf。

29. 调用url(p45)
请求的url不同会调用不同的函数：对“/count”这个url的请求会调用到count这个函数，“/”url会调用默认的处理函数。

29.1. 后台运行服务程序
如果你的操作系统是Mac OS X或者Linux，那么在运行命令的末尾加上一个&符号，即可让程序简单地跑在后台。
windows下可以在另外一个命令行窗口运行这个程序、
$ go run .../main.go &
然后可以通过明亮来发送客户端请求了。

30. mu.Lock()和mu.Unlock()
竞态条件：在并发情况下，如果有两个请求同一时刻去更新count，那么这个值可能并不会被正确地增加；这个程序可能会引发严重的bug。
本例子中为了避免这个bug，保证每次修改变量的最多只能有一个goroutine，这就是代码里的mu.Lock()和mu.Unlock()调用将修改count的所有行为，抱在lock和unlcok中间的目的。
第九章，我们会进一步讲解共享变量。
(p45)

31. if语句相关(p46)
GO允许一个简单的语句结果作为循环的变量声明出现在if语句的最前面，这一点对徐哦呜处理很有用。
这样让代码更加简单，并且可以限制err这个变量的作用域（只限定在这里使用），这么做是很不错的。

如:
if err := r.ParseForm(); err != nil{
	log.Print(err)
}
// if 简单语句的结果； 判断语句{
// 	 执行
// }

和这样是一样的：
err：=r.ParseForm()//这样导致作用域扩大了，万一后面再有个err调用的，可能会引发错误。
if err != nil{
	log.Print(err)
}

32. 有关接口（p46）
在这些程序中，我们看到了很多不同的类型被输出到标准输出流中。
比如前面的fetch程序，把HTTP的响应数据拷贝到了os.Stdout
lissajous程序里我们输出的是一个文件。
fetchall程序则完全忽略了HTTP的响应Body，只是计算了下响应Body的大小，这个程序中把响应Body拷贝到了ioutil.Discard.package main
在本节的web服务器程序中则是用fmt.Fprintf直接写到了http.ResponseWriter中。

尽管三种具体的实现流程并不太一样，他们都实现一个共同的接口，即当它们被调用需要一个标准流输出时都可以满足。
这个接口叫作io.Writer，在7.1节中会详细讨论。

33.switch 多路选择（p49）
	
	switch coinflip(){
	case "heads"://注意冒号
		heads++
	case "tails":
		tails++
	default：
		fmt.Println("landed on edge!")
	}
符合一个case后，则break，也就是退出switch循环。

34.switch中的fallthrough语句
让相邻几个case都执行同一个逻辑。
也就是所有的case都执行一遍逻辑，然后才退出。

35.switch还可以不带操作对象。（不带操作对象时默认用哦true值代替，然后将每个case的表达式和true值进行比较）
可以直接罗列多种条件，像其他语言里面的多个if else一样，下面是一个例子。
func Signum(x int) int {
	switch {//这里switch后没有操作对象
	case x > 0:
		return +1
	default:
		return 0
	case x <0:
		return -1
	}
}
这种形式叫做无tag switch(tagless switch);这和switch true是等价的。

像for和if控制语句一样，
switch也可以紧跟一个简短的变量声明，一个自增表达式、赋值语句，或者一个函数调用。

36. break和continue（p50）
break和continue语句会改变控制流。
break会中断当前的循环，并开始执行循环之后的内容。
continue会跳过当前循环，并开始执行下一次循环。

这两个语句可以控制for、switch、select语句。

循环往往会有很多层，
如果我们想跳过的是更外层的循环的话，我们可以在响应的位置加上label（标签），这样break或continue就可以根据我们的想法来break或continue任意循环。

37. 命名类型（p50）
命名类型：类型声明使得我们可以很方便地给一个特殊类型一个名字。
因为struct类型声明通常非常地长，所以我们总要给这种struct取一个名字。

如：
type Point struct{
	X，Y int
}
var p Point
//后面就可以用p或Point使用这个struct。
方便区别。

38. 指针（p50）
指针：是一种直接存储了变量的内存地址的数据类型。
Go语言中，指针是可见的内存地址 ，
&操作符可以返回一个变量的内存地址，
并且*操作符可以获取指针指向的变量内容。
但是在Go语言中没有指针运算，也就是不能像C语言里可以对指针进行加或减操作。

39. 方法（p50）
方法是和命名类型关联的一类函数。
Go语言里比较特殊的是方法可以被关联到任意一种命名类型。（关注第六章）

40. 接口（p50）
接口是一种抽象类型，这种类型更可以让我们以同样的方式来处理不同的固有类型，不用关心它们的具体实现，
而只需要关注它们提供的方法。（关注第七章）

41 包（p50）
Go语言提供了一些很好用的package,并且这些package是可以扩展的。
Go语言编程大多数情况下就是用已有的package来写我们自己的代码。
可以在https://golang.org/pkg 和 https://godoc.org中找到标准库和社区写的package、


疑问：
什么是标注输入？以及怎么运行dup1，dup2这两个教程？

