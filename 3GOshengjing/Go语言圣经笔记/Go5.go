第五章 函数
函数可以让我们将一个语句序列打包为一个单元，然后可以从程序中其他地方多次调用。
函数的机制可以让我们将一个大的工作分解为小的任务 ，这样的小任务可以让不同程序员在不同时间、不同地方独立完成。
一个函数同时对用户隐藏了其实现细节。

一个网络爬虫，给我们足够的机会去探索递归函数，匿名函数，错误处理和函数其他的很多特性。

5.1 函数声明（p167）

1. 函数声明包括函数名，形式参数列表，返回值列表（可省略），以及函数体。(p167)
func name(parameter-list)(result-list){
	body
}

2. 形式参数列表描述了函数的参数名以及参数类型。
这些参数作为局部变量，其值由参数调用者提供。

3. 返回值列表描述了函数返回值的变量名以及类型。
如果函数返回一个 无名变量 或者 没有返回值，返回值列表的括号是可以省略的。
如果一个函数声明不包括返回值列表，那么函数体执行完毕后，不会返回任何值。

func hypot(x,y float64) float64{
	return math.Sqrt(x*x+y*y)
}
fmt.Println(hypot(3,4)) //"5"

x, y是形参名，3和4是调用时的传入的实数，函数返回了一个float64类型的值。
返回值也可以像形式参数一样被命名。
在这种情况下，每个返回值被声明成一个局部变量，并根据该返回值的类型，将其初始化为0。

如果一个函数函数在声明时，包含返回值列表，该函数必须以return语句结尾，
除非函数明显无法运行到结尾处。
例如函数在结尾时调用了panic异常或函数中存在无限循环。

4. 正如hypot一样，如果一组形参或返回值有相同类型，饿哦们不必每个形参都写出参数类型：
func f(i,j,k int, s, t string){/*...*/}
func f(i int, j int, k int, s string, t string){/*...*/}

5. 给出4种方法声明拥有2个int型参数和1个int型返回值的函数.blank identifier可以强调某个参数未被使用。
func add(x int, y int) int {return x + y}
func sub(x, y int)(z int) {z = x-y; return}
func first(x int, _ int) int {return x }
func zero(int, int) int {return 0}

fmt.Printf("%T\n", add) // "func(int, int) int"
fmt.Printf("%T\n", sub)// "func(int, int) int"
fmt.Printf("%T\n", first)// "func(int, int) int"
fmt.Printf("%T\n", zero)// "func(int, int) int"

6. 函数的类型被称为函数的标识符。
如果两个函数形式参数列表和返回值列表中的变量类型一一对应，那么两个函数被认为有相同的类型和标识符。
形参和返回值的变量名不影响函数标识符也不影响它们是否可以以省略参数类型的形式表示。

7. 每一次函数调用都必须按照声明顺序为所有参数提供实参。
在函数调用时，Go语言没有默认参数值，也没有任何方法可以通过参数名制定形参，
因此形参和返回值的变量名对于函数调用者而言没有意义。

8. 在函数体中，函数的形参作为局部变量，被初始化为调用者提供的值。
函数的形参和有名返回值作为函数最外层的局部变量，被存储在相同的词法块中。

9. 实参通过值的方式传递，因此函数的形参是实参的拷贝。
对应惭进行修改不会影响实参。
但是，如果实参包括引用类型，如指针，slice、map、function、channel等类型，
实参可能会由于函数的间接引用被修改。

10. 你可能会偶尔遇到没有函数体的函数声明，这表示该函数不是以Go实现的。（p168）
这样的声明定义了函数标识符。
package math
func Sin(x float64) float //implemented in assembly language 已被汇编语言实现了。

5.2 递归 (p169)
11. 函数可以是递归的，这意味着函数可以直接或间接的调用自身。
递归是一种强有力的技术，处理递归的数据结构。

golang.org/x/net/html扩展包没有加入到标准库的原因：1.仍在开发中；2.扩展包提供的功能很少被使用。

12. findlist1 (p169)
下列代码使用了非标准包golang.org/x/net/html，解析HTML。
golang.org/x/... 目录下存储了一些由Go团队设计、维护、对网络编程、
国际化文件处理、移动平台、图像处理、加密解密、开发者工具提供支持的扩展包。 

例子中调用golang.org/x/net/html的部分api如下所示。
html.Parse函数读入一组bytes.解析后，返回html.node类型的HTML页面树状结构根节点。
HTML拥有很多类型的节点如text（文本），commnet（注释）类型，在下面的例子中，我们值关注<name key='value'>形式的节点。

main函数解析HTML标准输入，通过递归函数visit获得links（链接），并打印出这些links：
</i>...findlinks1</i>

visit函数遍历HTML的节点树，从每一个anchor元素的href属性获得link，将这些links存入字符串数组中，并返回这个字符串数组。

为了遍历节点n的所有后代节点，每次遇到n的子节点时，visit递归的调用自身。
这些子节点存放在FirstChild列表中。


13. outline (p171)
在outline函数中，我们通过递归的方式遍历整个HTML节点树，并输出树的结构。
在outline内部，没遇到一个HTML元素标签，就将其入栈，并输出。

有一点值得注意： outline有入栈操作，但没有相对应的出栈操作。
当outline调用哦自身时，被调用者接收的是stack的拷贝。
被调用者的入栈操作，修改的是stack的拷贝，而不是调用者的stack，因对当函数返回时，调用者的stack并未被修改。

正如上面实验所见，大部分HTML页面只需要几层递归就能被处理，
但任然有些页面需要深层次的递归。

Go语言使用可变栈，栈的大小按需增加（初始时很小）。
这使得我们使用递归时不必考虑溢出和安全问题。


5.3 多返回值（p173）
14. 在Go语言中，一个函数可以返回多个值。
许多标准库中的函数返回2个值，一个是期望得到的返回值，另一个是函数出错是的错误信息。

findlinks2可以自己发起HTTP请求，这样我们就不必再运行fetch。
因为HTTP请求和解析操作可能会失败，因此findlinks2声明了2个返回值：
链接列表和错误信息。

HTML的解析器可以处理HTML页面的错误节点，构造出HTML页面结构，
所以解析HTML很少失败。
这意味着如果findlinks函数失败了，很可能是由于I/O（什么东西？）的错误导致的。
(I/O的错误是什么？？？)

在findlinks函数中，有4处return语句，每一处return都返回了一组值。
前三处return，将http和html包中的错误信息传递给findlinks的调用者。
第一处return直接返回错误信息，其他两处通过fmt.Errorf(7.8)输出详细的错误信息。
如果findlinks成功结束，最后的return语句将一组解析获得的链接返回给用户。

在findlinks函数中，我们必须确保resp.Body被关闭，释放网络资源。
虽然Go的垃圾回收机制会回收不被使用的内存，但是这不包含操作系统层面的资源。
比如打开的文件、网络链接。
因此，我们必须显式的释放这些资源。

调用多返回函数时，返回给调用者的是一组值，调用者必须显式的将这些值分配给变量：
links, err := findLinks(url)

如果某个值不被使用，可以将其分配给blank identifier：
links, _ := findLinks(url) // errors ingored

一个函数内部可以将另一个有多返回值的函数作为返回值，下面的例子展示了与findLinks有相同功能的函数，
两者的区别在于下面的例子先输出参数：
func findLinksLog(url string) ([]string, error){
	log.Printf("findLinks %s", url)
	return findLinks(url)
}

当你调用接受多参数的函数时，可以将一个返回多参数的函数作为该函数的参数。
虽然这很少出现在现实生产代码中，但这个特性在debug（调试）时很方便，
我们只需要一条语句就可以输出所有的返回值。
下面的代码是等价的：
log.Println(findLinks(url))

links, err := findLinks(url)
log.Println(links, err)

15. 准确的变量名可以传达函数返回值的含义。（p174）
尤其在返回值的类型都相同时，就像下面这样：
func Size(rect image.Rectangle) (width, height int)
func Split(path string) (dir, file string)
func HourMinSec(t time.Time)(hour, minute, second int)

虽然良好的命名很重要，但你也不必为每一个返回值都去一个适当的名字。
如，惯例，函数的最后一个bool类型的返回值表示函数示范运行成功，
error类型的返回值代表函数的错误信息，对于这些类似的惯例，我们不必思考合适的命名，他们都无需解释。

如果一个函数将所有的返回值都显示的变量名，那么该函数的return语句可以省略操作数。
这称之为 bare return。
func CountWordsAndimages(url string) (words, images int, err error){
	resp, err := http.Get(url)
	if err != nil{
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil{
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}
func countWordsAndImages(n *html.Node)(words, images int){/*...*/}

按照返回值列表的顺序，返回所有的返回值，在上面的例子中，每一个 return 语句等价于：
return words, images, err

16. 当一个函数有多处 return 语句以及许多返回值时，bare return 可以减少代码的重复，（p175）
但是使得代码难以被理解。
不宜过度使用 bare return。
举例：
如果你没有仔细审查代码，很难发现前2处return等价于return 0, 0, err 
(GO会将返回值words和images在函数体的开始处， 根据他们的类型，将其初始化为0值。)
最后一处return等价于 return words, image, nil。

5.4 错误(p176)
17. panic异常是来自被调函数的信号，表示发生了某个已知的bug。
一个良好的程序永远不应该发生panic异常。

18. 在Go的初五处理中，错误是软件包API和应用程序用户界面的一个重要组成部分，
程序运行失败仅被认为是几个预期的结果之一。

对于那些将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，
通常是最后一个，来传递错误信息。
如果导致失败的原因只有一个，额外的返回值可以是一个布尔值，
通常被命名为ok。
如：
cache.Lookup失败的唯一原因是key不存在，那么代码可以按照下面的方式组织：
value, ok := cache.Lookup(key)
if !ok{//如果ok不存在
	// ...cache[key] does not exist
}

19. 通常，导致失败的原因不止一种，尤其是对I/O操作而言，用户u需要了解耕读偶读错误信息。
因此，额外的返回值不再是简单的布尔类型，而是error类型。

内置的error是接口类型。
我们将在第7章了解接口类型的含义，以及它对错误处理的影响。
现在我们只需要明白error类型可能是nil或者non-nil。
nil意味着函数运行成功，non-nil表示表示失败。
对于non-nil的error类型，我们可以通过调用error的Error函数或者输出函数获得字符串类型的错误信息。
fmt.Println(err)
fmt.Println("%v", err)

20. 通常，当函数返回non-nil的error时，其他的返回值是未定义的（underfined），(p177)
这些未定义的返回值应该被忽略。 然而，有少部分函数在发送错误时，任然会返回一些有用的返回值。
比如，当读取文件发生错误时，Read函数会返回可以读取的字节数以及错误信息。
对于这种情况，正确的处理方式应该是先处理这些不完整的数据，再处理错误。
因此对函数的返回值要有清晰的说明，以便于其他人使用。

在Go中，函数运行失败时会返回错误信息，这些错误信息被认为是一种预期的值而非异常（exception），
这使得Go有别于那些将函数运行失败看作是异常的语言。
虽然Go有各种异常机制，但这些机制仅被使用在处理那些未被遇到到的错误，即bug，而不是那些在健康程序中应该被避免的错误。
对于Go的异常机制我们将在5.9介绍。

Go这样涉及的原因是由于对于某个应该在控制流程中处理的错误而言，将这个错误以异常的形式抛出会混乱对错误的描述，这通常会导致一些糟糕的后果。
当某个程序错误被当做异常处理后，这个错误会将队栈根据信息返回给终端用户，这些信息复杂且无用，无法帮助定位错误。

正因此，GO使用控制流机制（如if和return）处理异常，这使得编码人员能更多的关注错误处理。

5.4.1 错误处理策略
21. 当一次函数调用返回错误时，调用者有应该选择何时的方式处理错误。
根据情况的不同，有很多处理方式，让我们来看看常用的五中方式。

21.1 首先，也是最常用的方式是传播错误。这意味着函数中某个子程序的失败，会变成该函数的失败。
下面，我们以5.3节的findLink函数作为例子。
如果findLinks对http.Get的调用失败，findLinks会直接将这个HTTP错误返回给调用者：

resp, err := http.Get(url)
if err != nil{
	return nil, err
}

当对html.Parse的调用失败时，findLinks不会直接返回html.Parse的错误，
因为缺少两条重要信息：
1、错误发生在解析器；
2、url已经被解析。
这些信息有助于错误的处理，
findLinks会构造新的错误信息返回给调用者:
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil{
	return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
}

22. fmt.Errorf函数使用fmt.Sprintf格式化错误信息并返回。
我们使用该函数前缀添加额外的上下文信息到原始错误信息。
当错误最终由main函数处理时，错误信息应提供清晰的从原因到后果的因果链，就像美国宇航局事故调查时做的那样：
genesis: crashed: no parachute: G-switch failed: bad relay orientation

23. 由于错误信息经常是以链式组合在一起的，所以错误信息中应避免大写和换行符。
最终的错误信息可能很长，我们可以通过类似grep的工具处理错误信息（grep是一种文本搜索工具）
编写错误信息时，我们要确保错误信息对问题细节的描述是详尽的。
尤其是要注意错误信息表达的一致性，即相同的函数或同包内的同一组函数返回的错误在构成和处理方式上是相似的。

以OS包为例，OS包确保文件操作（如os.Open\Read\Write\Close）
返回的每个错误的描述不仅仅包含错误的原因（如无权限，文件目录不存在）也包含文件名，
这样调用者在构造新的错误信息时无需再添加这些信息。

一般而言，被调函数 f(x)会将调用信息和参数信息作为发生错误时的上下文放在错误信息中并返回给调用者，
调用者需要添加一些错误信息中不包含的信息，比如添加url到html.Parse返回中的错误中。

21.2 偶然性的错误
如果错误的发生是偶然性的，或由不可预知的问题导致的。
一个明智的选择是重新尝试失败的操作。
在重试时，我们需要限制重试的时间间隔或重试的次数，防止无限制的重试。

21.3 如果错误发生后，程序无法继续运行，我们就可以采用第三章策略：
输出错误信息并结束程序。
需要注意的是，这种策略只应在main中执行。
对库函数而言，应仅向上传播错误，除非该错误意味着程序内部包含不一致性，即遇到了bug，才能在库函数中结束程序。
// (In function main.)
if err := WaitForServer(url); err != nil{
	fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
	os.Exit(1)
}

调用log.Fatalf可以更简洁的代码达到与上文相同的效果。log中的所有函数，都默认会在错误信息之前输出时间。
if err := WaitForServer(url); err != nil{
	log.Fatalf("Site is down: %v\n", err)
}

长时间运行的服务器常采用默认的时间格式，而交互式工具很少采用包含如此多信息的格式。
2020/03/22 15:04:05 Site is down: no such domain:
bad.gopl.io

我们可以设置log的前缀信息屏蔽时间信息，一般而言，前缀信息会被设置成命令名。
log.SetPrefix("wait: ")
log.SetFlags(0)

21.4 第四种策略：有时，我们只需要输出错误信息就足够了，
不需要中断程序的运行。我们可以通过log包提供函数
if err := Ping(); err != nil{
	log.Printf("ping failed: %v; networking disabled", err)
}

24. log包中的所有函数会为没有换行符的字符串增加换行符。（p179）

21.5 第五种，我们可以直接忽略掉错误。
dir, err := ioutil.TempDir("", "scratch")
if err != nil{
	return fmt.Errorf("failed to creat temp dir: %v", err)
}
// ...use temp dir...
os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodically

尽管os.RemoveAll会失败，但上面的例子并没有做错误处理。
这是因为操作系统会定期的清理临时目录。
正因如此，虽然程序没有处理错误，但程序的逻辑不会因此受到影响。

25. 处理错误的态度：
我们应该在每次函数调用后，都养成考虑错误处理的习惯，
当你决定忽略某个错误时，你应该在清晰的记录下你的意图。

26. 在Go中，错误处理有一套独特的编码风格。（p180）
检查某个子函数是否失败后，我们通常将处理失败的逻辑代码放在处理成功的代码之前。
如果某个错误会导致函数返回，那么成功时的逻辑代码不应放在else语句块中，而应该直接放在函数体中。

27. Go中大部分函数的代码结构几乎相同，首先是一系列的初始检查，防止错误发生，之后是函数的实际逻辑。

5.4.2 文件结尾错误（EOF）(p180)
28. 函数经常会返回多种错误，这对终端用户来说可能会很有趣，
但对程序而言，这使得情况变得复杂。很多时候，程序必须根据错误类型，作出不同的相应。
让我们考虑这样一个例子：从文件中读取n个字节。
如果n等于文件的长度，读取过程的任何错误都表示失败。
如果n小于文件的长度，调用者会重复的读取固定大小的数据直到文件结束。
这会导致调用者必须分别处理由问价结束引起的各种错误。
基于这样的原因，io包保证任何由文件结束引起的读取失败都返回用一个错误——io.EOF，该错误在io包中定义：
package io
import "errors"

// EOF is the error returned by Read when no more input is available.
var EOF = errors.New("EOF")

调用者只需通过简单的比较，就可以检测出这个错误。
下面的例子展示了如何从标准输入中读取字符，以及判断文件结束。

in := bufio.NewReader(os.Stdin)
for{
	r, _, err := in.ReadRune()
	if err == io.EOF{
		break // finished reading
	}
	if err != nil{
		return fmt.Errorf("read failed:%v", err)
	}
	// ...use r...
}

因为文件结束这种错误不需要更多的描述，所以io.EOF有固定的错误信息————“EOF”。
对于其他错误，我们可能需要在错误信息中描述错误的类型和数量，这是的我们不能像io.EOF一样采用固定的错误信息。
在7.11节中，我们会提出更系统的方法区分某些固定的错误值。

5.5 函数值（p182）
29. 在Go中，函数被看作第一类值（first-class values）：函数像其他值一样，拥有类型，
可以被赋值给其他变量，传递给函数，从函数返回。
对函数值（function value）的调用类似函数调用。
如下
func square(n int) int {return n*n}
func negative(n int) int {return -n}
func product(m, n int)int {return m*n}

f := square
fmt.Println(f(3)) // "9" f(3) 就叫函数值

f = negative
fmt.Println(f(3)) //"-3"
fmt.Printf("%T\n", f)//"func(int) int"

f = product // compile error: can't assign func(int,int) int to func(int) int
//类型不对

30. 函数类型的零值是nil。
调用值为nil的函数值会引起panic错误。

var f func(int)int
f(3) //此处f的值为nil,会引起panic错误。

31. 函数值可以与nil比较：
var f func(int) int
if f != nil {
	f(3)
}

但是函数值之间是不可比较的，也不能用函数值作为map的key。

32. 函数值使得我们不仅可以通过数据来参数化函数，亦可通过行为。(p182)
标准库中包含许多这样的例子。
下面的代码展示了如何使用这个技巧。
strings.Map对字符串的每个字符调用add1函数，并将每个add1函数的返回值组成一个新的字符串返回给调用者。

func add1(r rune) rune {return r + 1}//rune等价int32,unicode码点
fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
fmt.Println(strings.Map(add1, "VMS")) // "WNT"
fmt.Println(strings.Map(add1, "Admix")) //"Benjy"

5.2 节的findLinks函数使用了辅助函数visit，遍历和操作了HTML页面的所有节点。
使用函数值，我们可以将遍历结点的逻辑和操作节点的逻辑分离，使得我们可以复用遍历的逻辑，
从而对结点进行不同的操作。

示例代码。。。

该函数接收2个函数作为参数，分别在子结点被访问前后访问后调用。
这样的设计给调用者更大的灵活性。
举个例子，现在我们有startElement和endElement两个函数用于输出HTML元素的开始标签和结束标签<b>...</b>:

var depth int
func startElement(n *html.Node){
	if n.Type == html.ElementNode{
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data) // n.Data是标签的意思？？？
		depth++
	}
}

func endElement(n *html.Node){
	if n.Type == html.ElementNode{
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}


33. 上面的代码利用fmt.Printf的一个小技巧控制输出的缩进。
%*s中的*会在字符串之前填充一些空格。在例子中，每次出书会先填充depth*2数量的空格，再输出""，最后在出书HTML标签。
如果我们像下面这样调用forEachNode：
forEachNode(doc, startElement, endElement)
与之前的outline程序相比，我们得到了更加详细的页面结构：
$ go build .../outline2
$ ./outline2 http://gopl.io

<html>
	<head>
		<meta>
		</meta>
		<title>
		</title>
		<style>
		</style>
	</head>
	<body>
		<table>
			<tbody>
				<tr>
					<td>
						<a>
							<img>
							</img>
...


5.6 匿名函数（p185）

34. 拥有函数名的函数只能在包级语法块中被声明，通过函数字面量（function literal），
我们可绕过这一限制，在任何表达式中表示一个函数值。
函数字面量的语法和函数声明相似，区别在于func关键字后没有函数名。
函数字面量是一种表达式，它的值被称为匿名函数（anonymous function)。

函数字面量允许我们在使用函数时，再定义它。
通过这种技巧，我们可以改写之前对strings.Map的调用：
strings.Map(func(r rune) rune{return r+1}, "HAL-9000")
更为重要的是，通过这种方式定义的函数可以访问完整的词法环境（lexical environment），
这意味着在函数中定义的内部函数可以引用该函数的变量，如下例所示。

示例代码

函数squares返回另一个类型为 func()int的函数。
对squares的一次调用会生成一个局部变量x并返回一个匿名函数。
每次调用匿名函数时，该函数都会先使x的值加1，在返回x的平方。
第二次调用squares时，会生成第二个x变量，并返回一个新的匿名函数。
新匿名桉树操作的是第二个x变量。

35. squares的例子证明，函数值不仅仅是遗传代码，还记录了状态。
在squares中定义的匿名内部函数可以访问和更新squares中的局部变量，这意味着匿名函数和squares中，
存在变量引用。
这就是函数值属于引用类型和函数值不可比较的原因。
Go使用闭包（closures）技术实现函数值，Go程序员也把函数值叫做闭包。

通过这个例子，我们看到变量的生命周期不由它的作用域决定：squares返回后，变量x仍然隐式的存在于f中。

36. 我们讨论一个有点学术性的例子，考虑这样一个问题：
给定一些计算机课程，每个课程都有前置课程，只有完成了前置课程才可以开始当前课程的学习；
我们的目标是选择出一组课程，这组课程必须确保按顺序学习时，能全部被完成。
每个课程的前置课程如下：

示例代码

这类问题被称作拓扑排序。从概念上说，前置条件可以过程有向图。
图中的顶点表示课程，边表示课程键的依赖关系。
显然，图中应该无环，这也就是说从某点出发的边，最终不会回到该点。
下面的代码用深度优先搜索了整张图，获得了符合要求的可能序列。

示例代码

当匿名韩式需要被递归调用时，我们必须首先声明一个变量（在上面的例子中，我们首先声明了visitAll）
，再将匿名函数赋值给这个变量。
如果不分成两部，函数字面量无法与visitAll绑定，我们也无法递归调用该匿名函数。

vistiAll := func(items []string){
	//...
	visitAll(m[item]) // compile error: undefined: visitAll
	//...
}

在topsort中，首先对prereqs中的key排序，再调用visitAll。
因为prereqs映射的是切片而不是更复杂的map，所以数据的遍历次序是固定的，这意味着你每次运行topsort得到的输出都是一样的。
topsort的输出结果如下：（略）

37. 让我们回到findLinks这个例子。
我们将代码移动到了links包下，将函数重命名为Extract，
在第8章我们会再次用到这个函数。
新的匿名函数被引入，用于替换原来了的visit函数。
该匿名函数负责将新链接添加到切片中。
在Extract中，使用forEachNode遍历HTML页面，
由于Extract只需要在遍历节点前操作节点，所以forEachNode的post参数被传入nil。

示例代码

上面的代码对之前的版本做了改进，现在links中存储的不是href属性的原始值，
而是通过resp.Request.URL解析后的值。解析后，这些连接以绝对路径的形式存在，
可以直接被http.Get访问。
网页抓取的核心问题就是如何遍历图。
在topoSort的例子中，已经展示了深度优先遍历，在网页抓取中，我们会展示如何用广度优先遍历图。
在第8章，我们会介绍如何将深度优先和广度优先结合使用。

38. 下面的函数实现了广度优先算法。（p190）
调用者需要输入一个初始的待访问列表和一个函数f。
待访问列表中的每个元素被定义为string类型。
广度优先算法会为每个元素调用一次f。
每次f执行完毕后，会返回一组待访问元素。
这些元素会被加入到待访问列表中。
当待访问列表中的所有元素都被访问后，hreadthFirst函数运行结束。
为了避免同一个元素被访问两次，代码中维护了一个map。

示例代码

就像我们在章节3解释的那样，append的参数“f(item)...”，会将f返回的一组元素一个一个添加到worklist中。

在我们网页抓取器中，元素的类型是url。crawl函数会将URL输出，提取其中的新链接，
并将这些新链接返回。我们会将crawl作为参数传递给breadthFirst。

示例代码

为了使抓取器开始运行，我们用命令行输入的参数作为初始的待访问url。

让我们从https://golang.org 开始，下面是程序的输出结果：
$ go build ...findlinks3
$ ./findlinks3 https://golang.org
https://golang.org/
https://golang.org/doc/
https://golang.org/pkg/
...

当所有发现的链接都已经被访问或电脑的内存耗尽时，程序运行结束。

5.6.1 警告：捕获迭代变量（p191）
39.本节，将介绍Go词法作用域的一个陷阱。
请务必仔细的阅读，弄清楚发生问题的原因。
即使是经验丰富的程序员也会在这个问题上犯错误。

考虑这样一个问题：你被要求首先创建一些目录，再将目录删除。
在下面的例子中我们用函数值来完成删除操作。
下面的示例代码需要引入os包。
为了使代码简单，我们忽略了所有的异常处理。

var rmdirs []func()
for _, d := range tempDirs(){
	dir := d // NOTE: necessary!
	os.MkdirAll(dir, 0755) // creates parent directories too
	rmdirs = append(rmdirs, func(){
		os.RemoveAll(dir)
	})
}
// ...do some work...
for _, rmdir := range rmdirs{
	rmdir() // clean up 
}

你可能会感到困惑，为什么要在循环体中用循环变量d赋值一个新的局部变量，
而不是像下面的代码一样直接使用循环变量dir。
需要注意，下面的代码是错误的。

var rmdirs []func()
for _, dir := range tempDirs(){
	os.MkdirAll(dir, 0755)
	rmdirs = append(rmdirs, func(){
		os.RemoveAll(dir)// NOTE: incorrect!
	})
}

问题的原因在于循环变量的作用域。
在上面的程序中，for循环语句引入了新的词法域，
循环变量dir在这个词法块中被声明。
在该循环中生成的所有函数值都共享相同的循环变量。
需要注意，函数值中记录的是循环变量的内存地址，而不是循环变量有一时刻的值。
以dir为例，后续的迭代会不断更新dir的值，当删除操作执行时，for循环已完成，
dir中存储的值等于最后一次迭代的值。
这意味着，每次对os.RemoveAll的调用欧删除的都是相同的目录。

通常，为了解这个问题，我们会引入一个与循环变量同名的局部变量，作为循环变量的副本。
比如下面的变量dir，虽然这看起来很奇怪，但却很有用。
for _, dir := range tempDirs(){
	dir := dir // declares inner dir, initialized to outer dir 
	// ...
}

这个问题不仅存在基于range的循环，在下面的例子中，对循环变量i的使用也存下同样的问题：

var rmdirs []func()
dirs := tempDirs()
for i := 0; i<len(dirs); i++{
	os.MkdirAll(dirs[i], 0755) //OK
	rmdirs = append(rmdirs, func(){
		os.RemoveAll(dirs[i]) // NOTE: incorrect!
	})
}

如果你使用Go语言（第八章）或者defer语句（5.8）会经常遇到此类问题。
这不是go或deger本身导致的，
而是因为他们都会等待循环结束后，再执行函数值。


5.7 可变参数（p194）

40. 参数数量可变的函数称为可变参数函数。
典型的例子就是fmt.Printf和类似函数。
Printf首先接收一个必备参数，之后接收任意个数的后续参数。

func sum(val ...int) int{
	total := 0
	for _, val := range vals{
		total += val
	}
	return total
}

sum函数返回任意个int类型参数的和。
在函数体中，vals被看作是类型为[]int的切片。
sum可以接收任意数量的int型参数。
fmt.Println(sum(1,2,3,4))// "10"

如果原始参数已经是切片类型，我们该如何传递给sum？
只需在最后一个参数后加上省略符。
values := []int{1,2,3,4}
fmt.Println(sum(values...))

虽然在可变参数函数内部，...int型参数的行为看起来很像切片类型，
但实际上，可变参数函数和以切片作为参数的函数是不同的。

func f (...int){}
func g ([]int){}
fmt.Printf("%T\n", f) // "func(...int)" //可变参数函数
fmt.Printf("%T\n", g) // "func([]int)" //参数是切片类型的函数

41. 可变参数函数经常被用于格式化字符串。
下面的errorf函数构造了一个以行号开头的，经过格式化的错误信息。
函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串。
func errorf(linenum int, format string, args ...interface{}){
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprinln(os.Stderr)
/*
在通常情况下，UNIX每个程序在开始运行的时刻，都会有3个已经打开的stream. 分别用来输入，输出，打印诊断和错误信息。通常他们会被连接到用户终端. 但也可以改变到其它文件或设备。
Linux内核启动的时候默认打开的这三个I/O设备文件：标准输入文件stdin，标准输出文件stdout，标准错误输出文件stderr，分别得到文件描述符 0, 1, 2。
stdin是标准输入，stdout是标准输出，stderr是标准错误输出。大多数的命令行程序从stdin输入，输出到stdout或stderr。
*/
}
linenum, name := 12, "count"
errorf(linenum, "undefined: %s", name) // "Line 12: underfined: count"

interface{}表示函数的最后一个参数可以接收任意类型，第7章详细介绍。


5.8 Deferred函数（p196）
在findLinks的例子中，我们用http.Get的输出作为html.Parse的输入。
只有url的内容的确是HTML格式的，html.Parse才可以正常工作，但实际上，url指向的内容很丰富，可能是图片，纯文本或是其他。
将这些格式的内容传递给html.Parse，会产生不良后果。

下面的例子获取HTML页面并输出页面的标题。
title函数会检查服务器返回的Content-Type字段，如果发现页面不是HTML，将终止函数运行，返回错误。

随着函数变得复杂，需要处理的徐哦呜也变多，维护清理逻辑变得越来越困难。
而Go语言独有的defer机制可以让事情变得简单。

42. 你只需要在调用普通函数或方法前加上关键字defer，就完成了defer所需要的语法。
当defer语句被执行时，跟在defer后面的函数会被延迟执行。直到包含该语句的函数执行完毕时，
defer后的函数才会被执行，不论包含defer语句的函数是通过return正常结束，还是由于panic导致的异常结束。
你可以在一个函数中执行多条defer语句，他们的执行顺序与声明顺序相反。


43. defer语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。
通过defer机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。
释放资源的defer应该直接跟在资源的语句后。
在下面的代码中，一条defer语句代替了之前的所有resp.Body.Close

示例代码

处理其他资源时，也可以采用defer机制，
比如对文件的操作：io/ioutil
package ioutil
func ReadFile(filename string)([]byte, error){
	f, err := os.Open(filename)
	if err != nil{
		return nil, err
	}
	defer f.Close()
	return ReadAll(f)
}

或是处理互斥锁（9.2 章）
var mu sync.Mutex
var m = make(map[string]int)
func lookup(key string)int{
	mu.Lock()
	defer mu.Unlock()
	return m[key]
}

44. 调试复杂程序时，defer机制也常被用于记录合适进入和退出函数。
下列中的bigSlowOperation函数，直接调用trace记录函数的被调情况。
bigSlowOperation被调时，trace会返回一个函数值，该函数值会在bigSlowOperation退出是被调用。
通过这种方式，我们可以值通过一条语句控制函数的入口和所有的出口，甚至可以记录函数的运行时间，如例子中的start。
需要之一一点：不要忘记defer语句后的圆括号，否则本该在进入时执行的操作会在退出时执行，而本该在退出时执行的，将永远不会被执行。

示例代码。

45. 我们知道，defer 语句中的函数会在return语句更新返回值变量后再执行，又因为在函数中定义的匿名函数可以访问函数包括返回值变量在内的所有变量，所以，对匿名函数采用defer机制，可以使其观察函数的返回值。
以double函数为例：
func double(x int) int {
	return x + x
}
我们只需要首先命名double的返回值，再增加defer语句，我们就可以在double每次被调用时，输出参数以及返回值。
func double (x int)(result int){
	defer func(){fmt.Printf("double(%d) = %d\n", x, result)	}() //别忘了小括号，这是defer函数
	return x+x
}
_ = double(4)
// Output:
// "double(4) = 8"

可能double函数过于简单，看不出这个小技巧的作用，但对于有许多return语句的函数而言，这个技巧很有用。

46. 被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值：
func triple(x int)(result int){
	defer func(){result += x}()
	return double(x)	
}
fmt.Println(triple(4)) // "12"


47. 在循环体中的defer语句需要特别注意，因为只有在函数执行完毕后，这些被延迟的函数才会执行。
下面的代码会导致系统的文件描述符耗尽，因为在所有文件都被处理之前，没有文件会关闭。

for _, filename := range filenames{
	f, err := os.Open(filename)
	if err != nil{
		return err
	}
	defer f.Close() // NOTE: risky; could run out of file
	descriptors
	// ...process f...
}

一种解决方法是将循环体中的defer语句移至另外一个函数。在每次循环时，调用这个函数。
for _, filename := range filename{
	if err := doFile(filename); err != nil{
		return err
	}
}

func doFile(filename string)error{
	f, err := os.Open(filename)
	if err != nil{
		return err
	}
	defer f.Close()
	// ...process f...
}

48. 下面的代码是fetch（1.5 节）的改进版，我们将http响应信息写入本地文件而不是从标准输出流输出。
我们通过path.Base提出url路径的最后一段作为文件名。

示例代码

对resp.Body.Close延迟调用我们已经见过了，再次不做解释。
上例中，通过os.Create打开文件进行写入，在关闭文件时，我们没有对f.close采用defer机制，
因为这会产生一些微妙的错误。
许多文件系统，尤其是NFS，写入文件时发生的错误会被延迟到文件关闭时反馈。
如果没有检查文件关闭时的反馈信息，可能会导致数据丢失，而我们还误以为写入操作成功。
如果io.Copy和f.close都失败了，我们倾向于将io.Copy的错误信息反馈给调用者，因为它先于f.close发生，更有可能接近问题的本质。

5.9 Panic 异常（p202）

49. Go的类型系统会在编译时捕获很多错误，但有些错误只能在运行时检查，如数组访问越界、空指针引用等。
这些运行时错误会引起panic异常。

一般而言，当panic异常发生时，程序会中断运行，并立即执行在该goroutine（可以先理解成线程，在第8章详细介绍）中被延迟的函数（defer机制）。
随后，程序崩溃并输出日志信息。
日志信息包括panic value和函数调用的堆栈跟踪信息。
panic value通常是某种错误信息。
对于每个goroutine， 日志信息中都会有与之相对的，发生panic时的函数调用堆栈跟踪信息。
通常，我们不需要再次运行程序去定位问题，日志信息已经提供了足够的诊断依据。
因此，在我们填写问题报告时，一般会将panic异常和日志信息一并记录。

不是所有的panic异常都来自运行时，直接调用内置的panic函数也会引发panic异常；
panic函数接收任何值作为参数。
当某些不该发生的场景发生时，我们就应该调用panic。
比如，当程序到达了某条逻辑上不可能到达的路径：
switch s := suit(drawCard()); s{
case "Spades":
case "Hearts":
case "Diamonds":
case "Clubs":
	panic(fmt.Sprintf("invalid suit %q",s))
}

50. 断言函数必须满足的前置条件是明智的做法，但这很容易被滥用。
除非你嫩更提供更多的错误信息，或者能更快速的发现错误，否则不需要是用断言，编译器在运行时会帮你检查代码。

func Reset(x *Buffer){
	if x == nil{
		panic("x is nil") // unnecessary!
	}
	x.elements = nil
} 

虽然Go的panic机制类似于其他语言的异常，但panic的适用场景有一些不同。
由于panic会引起程序的崩溃，因此panic一般用于严重错误，如程序内部的逻辑不一致。
勤奋的程序员认为任何崩溃都表面代码中存在漏洞，所以对于大部分漏洞，我们应该适用Go提供的错误机制，
而不是panic，尽量避免程序的崩溃。
在健壮的程序中，任何可以预料到的错误，如不正确的输入、错误的配置或是失败的i/o操作都应该被优雅的处理。
最好的处理方式，就是适用Go的错误机制。

考虑regexp.Compile函数，该函数将正则表达式编译成有效的可匹配格式。当输入的正则表达式不合法时，该函数会返回一个错误。
当调用者明确的知道正确的输入不会引起函数错误时，要求调用这检查这个错误是不必要和累赘的。
我们应该假设函数的输入一直合法，就如前面的断言一样：当调用者输入了不应该出现的输入时，触发panic异常。

在程序源码中，大多数正则表达式是字符串字面值（string literals），
因此regexp包提供了包装函数regexp.MustCompile检查输入的合法性。
package regexp
func Compile(expr string)(*Regexp, error){/*...*/}
func MustCompile(expr string)*Regexp{
	re, err := Compile(expr)
	if err != nil{
		panic(err)
	}
	return re
}

包装函数使得调用者可以便捷的用一个编译后的正则表达式为包级别的变量赋值：
var httpSchemeRE = regexp.MustCompile(`^https?:`) // "http:" or "https:"
显然， MustCompile不能接收不合法的输入。函数名中的Must前缀是一种针对此类函数的命名约定，比如与template.Must(4.6节)

func main(){
	f(3)
}
func f(x int){
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)//先defer后执行
	f(x - 1)
}

我们在下一节将看到，如何使程序从panic异常中恢复，阻止程序的崩溃。

为了方便诊断问题，runtime包允许程序员输出堆栈信息。
在下面的例子中，我们通过在main函数中延迟调用printStack输出堆栈信息。

func main(){
	defer printStack()
	f(3)
}
func printStack(){
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

51. 将panic机制类比其他语言异常机制的渎职可能会惊讶，runtime.Stack为何能输出已经被释放函数的信息？
在Go的panic机制中，延迟函数的调用在释放堆栈信息之前。


5.10 Recover捕获异常（p206）

通常来说，不应该对panic异常做任何处理，
但有时，也许我们可以从异常中回复，至少我们可以在程序崩溃前，做一些操作。
举个例子，当web服务器遇到不可预料的严重问题时，在崩溃前应该将所有的连接关闭；
如果不做任何处理，会使得客户端一直处于等待状态。
如果web服务器还在开发阶段，服务器甚至可以将异常信息反馈到客户端，帮助调试。

52. 如果在deferred函数中调用了内置函数recover，并且定义该defer语句的函数发生了panic异常，recover会是程序从panic中恢复，并返回panic value。
导致panic异常的函数不会继续运行，但能正常返回。
在未发生panic时调用recover，recover会返回nil。

让我们以语言解析器为例，说明recover的使用场景。
考虑到语言解析器的复杂性，即使某个语言解析器目前工作正常，也无法肯定他没有漏洞。
因此，当某个异常出现时，我们不会选择让解析器崩溃，而是会将panic异常当做普通的解析错误，并附加额外信息体总用户报告此错误。

func Parse(input string)(s *Syntax, err error){
	defer func(){
		if p := recover(); p != nil{
			err = fmt.Errorf("internal error: %v", p)
		}
	}()
	// ...parser...
}

deferred函数帮助Parse从panic中恢复。
在deferred函数内部，panic value被附加到错误信息中；并用err变量接收错误信息，返回给调用者。
我们也可以通过调用runtime.Stack往错误信息中添加完整的堆栈调用信息。

53. 不加区分的恢复所有的panic异常，不是可取的方法；
因此在panic之后，无法保证包级变量的状态任然和我们预期一致。
比如，对数据结构的一次重要更新没有被完整完成、文件或者网络连接没有被关闭、获得的锁没有被释放。
此外，如果写日志时产生的panic被不加区分的恢复，可能会导致漏洞被忽略。

虽然把对panic的处理都集中在一个包下，有助于主简化对复杂和不可以预料问题的处理，
但作为被广泛遵守的规范，你不应该试图去恢复其他包引起的panic。
公有的API应该将函数的运行失败作为error返回，而不是panic。
同样的，你也不应该恢复一个由他人开发的函数引起的panic，比如说调用者传入的回调函数，因为你无法确保这样做是安全的。

有时我们很难完全遵循规范，举个例子，net/http包中提供了一个web服务器，将收到的请求分发给用户提供的处理函数。
很显然，我们不能因为某个处理函数引发的panic异常，杀掉整个进程；
web服务器遇到处理函数导致的panic时会调用recover，输出堆栈信息，继续运行。

这样的做法在实践中很便捷，但也会引起资源泄露，或是因为recover操作，导致其他问题。

基于以上原因，安全的做法是有选择性的recover。
换句话说，只恢复应该被恢复的panic异常，
此外，这些异常所占的比例应该尽可能的低。
为了标识某个panic是否应该被恢复，我们可以将panic value设置成特殊类型。
在recover时对panic value进行检查，如果发现panic value是特殊类型，就将这个panic作为error处理，如果不是，则按照正常的panic进行处理。
（下面的例子中，我们会看到这种方式）

下面的自理是title函数的变形，如果HTML页面包含多个<title>，该函数会给调用者返回一个错误（error）。
在sole Title内部处理时，如果检测到有多个<title>，会调用panic，组织函数继续递归，并将特殊类型bailout作为panic的参数。

示例代码

在上例中，deferred函数调用recover，并检查panic value。
当panic value是bailout{}类型时，deferred函数生成一个error返回给调用者。
当panic value是其他non-nil值时，表示发生了位置的panic异常，deferred函数将调用panic函数并将当前的panic value作为参数传入；
此时，等同于recover没有做任何操作。
（注意，在例子中，对可预期的错误采用了panic，这违反了之前的建议，我们在此只是想向读者演示这种机制。）

有些情况下，我们无法恢复。某些指明错误会导致Go在运行时终止程序，如内存不足。

