第七章 接口 
(p231-288)

1. 接口类型是对其它类型行为的抽象和概括；
因为接口类型不会和特定的实现细节绑定在一起，
通过这种抽象的方式我们可以让我们的函数更加灵活和更具有适应能力。

2. 很多面向对象的语言都有相似的接口概念，但go语言中接口类型的独特之处在于它是满足隐式实现的。
也就是说，我们没有必要对于给定的具体类型定义所有满足的接口类型；
简单拥有一些必须的方法就足够了。
这种设计可以让你创建一个新的接口类型满足已经存在的具体类型却不会去改变这些类型的定义；
当我们使用的类型来自于不受我们控制的包时，这种设计尤其有用。

在本章，我们会开始看到接口类型和值的一些基本技巧。
顺着这种方式我们将学习几个来自标准库的重要接口。
很多Go程序中都尽可能多的去使用标准库中的接口。
最后，我们会在（7.10）看到类型断言的知识，
在（7.13）看到类型开关的使用并且学到他们是怎样让不同的类型的概括成为可能。

7.1 接口约定
（p232）

3. 目前为止，我们看到的类型都是具体的类型。
一个具体的类型可以准确的描述它所代表的值并且展示出对类型本身的一些操作方式就想数字类型的算术操作，
切片类型的索引、附加和取范围操作。
具体的类型还可以通过他的 方法 提供额外的行为操作。
总的来说，当你拿到一个具体的类型时你就知道他的本身是什么和你可以用它来做什么。

在Go语言中还存在着另外一种类型：接口类型。
接口类型是一种抽象的类型。
它不会暴露出它所代表的对象的内部值的结构和这个对象支持的基础操作的集合；
他们只会展示出它们自己的方法。
也就是说当你看到一个接口类型的值时，你不知道它是什么，
唯一知道的就是可以通过他的方法来做什么。

在本书中，我们一直使用两个相似的函数来进行字符串的格式化：fmt.Printf它会把结果写到标准输出
和fmt.Sprintf它会把结果以字符串的形式返回。
得益于使用接口，我们不必可悲的因为返回结果在使用方式上的一些浅显不同就必须把格式化这个最困难的过程复制一份。
实际上，这两个函数都使用了另一个函数fmt.Fprintf来进行封装。
fmt.Fprintf这个函数对它的计算结果会被怎么使用是完全不知道的。

package fmt

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error){
	return Fprintf(os.Stdout, format, args...) //os.Stdout是*os.File类型
}
func Sprintf(format string, args ...interface{}) string{
	var buf bytes.Buffer
	Fprintf(&buf, format, args...) //第一个参数&buf是一个指向可以写入字节的内存缓冲区
	return buf.String()

}

Fprintf的前缀F表示文件（file）也表明格式化输出结果应该被写入第一个参数提供的文件中。
在Printf函数中的第一个参数os.Stdout是*os.File类型；
在Sprintf函数中的第一个参数&buf是一个指向可以写入字节的内存缓冲区，然而它并不是一个文件类型尽管它在某种意义上和文件类型相似。

即使Fprintf函数中的第一个参数也不是一个文件类型。
它是io.Writer类型这是一个接口类型定义如下：

package io

// Writer is the interface that wraps the basic Write method.
type Writer interface{
	// Write writes len(p) bytes from p to the underlying data stream.
	// It returns the number of bytes written from p (0 <= n <= len(p))
	// and any error encountered that caused the write to stop early.
	// Write must return a non-nil error if it returns n < len(p).
	// Write must not modify the slice data, even temporarily.
	//
	// Implementtations must not retain p.
	Write(p []byte) (n int, err error)
}

io.Writer类型定义了函数Fprintf和这个函数调用者之间的约定。
一方面这个约定需要调用者提供具体类型的值就像*os.File和*bytes.Buffer，
这些类型都有一个特定签名和行为的Write的函数。
另一方面这个约定保证了Fprintf接收任何满足io.Writer接口的值都可以工作。
Fprintf函数可能没有假定写入的是一个文件或是一段内存，而是写入一个可以调用Write函数的值。

因为fmt.Fprintf函数没有对具体操作的值做任何假设而是仅仅通过io.Writer接口的约定来保证行为，
所以第一个参数可以安全地传入一个任何具体类型的值只需要满足io.Writer接口。
一个类型可以自由的使用另一个满足相同接口的类型来进行替换被称为可替换性（LSP里氏替换）。
这是一个面向对象的特征。

让我们通过一个新的类型阿里进行校验，下面*ByteCounter类型里的Write方法，仅仅在丢失写向它的字节前统计他们的长度。
（在这个+=赋值语句中，让 len(p)的类型和*c的类匹配的转换是必须的。）

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}



因为*ByteCounter满足io.Writer的约定，我们可以把它传入Fprintf函数中；
Fprintf函数执行字符串格式化的过程不会去关注ByteCounter正确的累加结果的长度。
var c ByteCounter
c.Write([]byte("hello"))
fmt.Println(c) // "5", = len("hello")
c = 0 // reset the counter
var name = "Dolly"
fmt.Fprintf(&c, "hello, %s", name) //把格式化之后的东西写入&c
fmt.Println(c) // "12", = len("hello, Dolly")

除了io.Writer这个接口类型，还有另一个对fmt包很重要的接口类型。
Fprintf和Fprintln函数向类型提供了一种控制他们值输出的途径。
在2.5节中，我们为Celsius类型提供了一个String方法以便于可以打印成这样"100°C"，在6.5节中我们给*IntSet添加一个String方法，这样集合可以用传统的符号来进行表示就像"{1 2 3}"。
给一个类型定义String方法，可以让它满足最广泛使用之一的接口类型fmt.Stringer：
package fmt

// The String method is used to print values passed
// as an operand to any format that accepts a string
// or to an unformatted printer such as Print.
type Stringer interface{
	String() string
}

我们会在7.10 节解释fmt包怎么发现哪些值是满足这个接口类型的。


7.2 接口类型
4. 接口类型具体描述了一系列方法的集合，一个实现了这些方法的具体类型是这个接口类型的示例。

io.Writer类型是用的最广泛的接口之一，因为它提供了所有的类型写入bytes的抽象，
包括文件类型，内存缓冲区，网络链接，HTTP客户端，压缩工具，哈希等。
io包中定义了很多其他有用的接口类型。
Reader可以代表任意可以读取bytes的类型，
Closer可以是任意可以关闭的值，
例如一个文件或是网络链接。（到现在你可能注意到了很多Go语言中单方法接口的命名习惯）

package io
type Reader interface {
	Read(p []byte)(n int, err error)
}
type Close interface{
	Close() error
}

再往下看，我们发现有些新的接口类型通过组合已经有的接口来定义。
下面是两个例子：
type ReadWriter interface{
	Reader
	Writer
}

type ReadWriteCloser interface{
	Reader
	Writer
	Closer
}

5. 接口内嵌
上面用到的语法和结构内嵌相似，我们可以用这种方式以一个简写命名另一个接口。
而不用声明它所有的方法。
这种方式本称为接口内嵌。
尽管略失简洁，我们可以像下面这样，不使用内嵌来声明io.Writer接口。

type ReadWriter interface{
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

//只使用Reader接口中的Read函数（方法）。

或者甚至使用混合的风格：

type ReadWriter interface{
	Read(p []byte) (n int, err error)
	Writer
}

上面3中定义方式都是一样的效果。
方法的顺序变化也没有影响，唯一重要的就是这个集合里面的方法。


7.3 现实接口的条件（p237）
6. 一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。
// 反过来说如果一个接口接受你这种类型，则就可以使用该接口。
例如，*os.File类型实现了io.Reader, Writer, Closer, 和ReadWriter接口。
*bytes.Buffer实现了Reader，Writer， 和ReadWriter这些接口，但是他没有实现Closer接口，
因为它不具有Close方法。
Go的程序员经常会简要的把一个具体的类型描述。
举个例子，*bytes.Buffer 是 io.Writer ; *os.Files 是io.ReadWriter。

7 接口指定的规则非常简单：表达一个类型属于某个接口只要这个类型实现这个接口。
所以：
var w io.Writer // io.Writer是接口类型
w = os.Stdout // OK: *os.File has Write method
w = new(bytes.Buffer) // OK: *bytes.Buffer has Write method
w = time.Second // compile error: time.Duration lacks Write method

var rwc io.ReadWriteCloser // io.ReadWriteCloser是接口类型
rwc = os.Stdout // OK: *os.File has Read, Write,  Close methods.
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method.

这个规则甚至适用于等式右边本身也是一个接口类型
w = rwc // OK: io.ReadWriteCloser has Write method
rwc = w // compile error: io.Writer lacks Close method

因为ReadWriter和ReadWriteCloser包含所有Writer的方法，
所以任何实现了ReadWrite和ReadWriteCloser的类型必定也实现了Writer接口。

8. 在进一步学习前，必须先解释表示一个类型持有一个方法当中的细节。
回想在6.2 章中，对于每一个命名过的具体类型T；它一些方法的接受者是类型T本身，
然而另一些则是一个T的指针。
还记得在T指针类型的参数上调用一个T的方法是合法的，只要这个参数是一个变量；
编译器隐式的获取了他的地址
但这仅仅是一个语法糖：T类型的值不拥有所有*T指针的方法，那这样它就可能只实现更少的接口。

举个例子可能会更清晰一点。
在6.5章，IntSet类型的String方法的接受者是一个指针类型，所以我们不能在一个不能寻址的IntSet值上调用这个方法：
type IntSet struct { /*...*/}
func (*IntSet)String()string
var _ = IntSet{}.String() // compile error: String requires *IntSet receiver

但是我们可以在一个IntSet值上调用这个方法：

var s IntSet
var _ = s.String() // OK: s is a variable and &s has a String method

然而，由于只有*IntSet类型有String方法，所有也只有*IntSet类型实现了fmt.Stringer接口：
var _ fmt.Stringer = &s // OK
var _ fmt.Stringer = s // compile error: IntSet lacks String method
//这个时候该是什么类型就得是什么类型。不像上一章可以隐式识别为另一种类型。

12.8章包含了一个打印出任意值的所有方法的程序，
然后可以使用godoc-analysis=typetool(10.7.4)展示每个类型的方法和具体类型和接口之间的关系。

9. 就像信封封装和隐藏新建一样，接口类型封装和隐藏具体类型和它的值。
即使具体类型有其他的方法也只有接口类型暴露出来的方法会被调用到：

os.Stdouot.Write([]byte("hello")) // OK: *os.File has Write method
os.Stdout.Close() // OK: *os.File has Close method

var w io.Writer
w = os.Stdout
w.Write([]byte("hello")) // OK: io.Writer has Write method
w.Close() // compile error: io.Writer lacks Close method

//暴露出来才有效， w.Close() 与 os.Stdout.Close()

一个有更多方法的接口类型，如io.ReadWriter， 和少一些方法的解扣子类型，如io.Reader，
进行对比，更多方法的接口类型会告诉我们更多关于他的值持有的信息，并且对实现他的类型要求更加严格。
那么关于interface{}类型，他没有任何方法，请讲出哪些具体的类型实现了他？

这看上去好像没有用，但实际上interface{}被成为空接口类型，是不可或缺的。
因为空接口类型对实现它的类型没有要求，所以我们可以将任意一个值赋给空接口类型。

var any interface{}

any = ture
any = 12.34
any = "hello"
any = map[string]int{"one":1}
any = new(bytes.Buffer)

尽管不是很明显，从本书最早的例子中我们就已经在使用空接口类型。
它允许像fmt.Println或者5.7章中的errorf函数接收任何的参数。

对于创建的一个interface{}值持有一个boolean， float， string， map， pointer， 或者任意其他的类型；
我们当然不能直接对它持有的值做操作，
因为interface{}没有任何方法。
我们会在7.10章中学到一种用类型断言来获取interface{}中值的方法。

因为接口实现只依赖于判断的两个类型的方法，所以没有必要定义一个具体类型和它实现的接口之间的关系。
也就是说，尝试文档化和断言这种关系几乎没有用，所以并没有通过程序强制定义。
下面的定义在编译期断言一个*bytes.Buffer的值实现了io.Writer接口类型：

// *bytes.Buffer must satisy io.Writer
var w io.Writer = new(bytes.Buffer)

因为任意bytes.Buffer的值，甚至包括nil通过 (bytes.Buffer)(nil)进行显示的转换都实现了这个接口，
所以我们不必分配一个新的变量。
并且因为我们绝不会引用变量w，我们可以使用空标识符来进行代替。
总的看，这些变化可以让我们得到一个更朴素的版本：

// *bytes.Buffer must satisfy io.Writer
var _ io.Writer = (*bytes.Buffer)(nil)

非空的接口类型比如io.Writer经常被指针类型实现，尤其当一个或多个接口方法像Write方法那样隐式的给接收者带来变化的时候。
一个结构体的指针是非常常见的承载方法的类型。

但是并不意味着只有指针类型满足接口类型，甚至连一些有设置方法的接口类型也可能会被Go语言中其他的引用类型实现。
我们已经看过slice类型的方法（geometry.Path， 6.1）和map类型的方法（url.Values, 6.2.1）后面还会看到函数类型的方法的例子（http.HandlerFunc, 7.7）。
甚至基本的类型也可能会实现一些接口；
就如我们在7.4 章中看到的time.Duration 类型实现了fmt.Stringer接口。

一个具体的类型可能实现了很多不相关的接口。
考虑在一个组织出售数字文化船票如音乐，电影和书籍的程序中可能定义了下列的具体类型：

Album
Book
Movie
Magazine
Podcast
TVEposode
Track

我们可以把每个抽象的特点用接口表示。
一些特性对于所有的这些文化产品都是共通的，例如标题，创作日期和作者列表。

type Artifact interface{
	Title() string
	Creators() []string
	Created() time.Time
}

其他的一些特性只对特定类型的文化产品才有。
和文字排版特性相关的只有books和magazines，还有只有movies和TV剧和屏幕分辨率相关。

type Text interface{
	Pages() int
	Words() int
	PageSize() int
}

type Audio interface{
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string //e.g., "MP3, "WAV"
}

type Video interface{
	Stream() (io.ReadClose, error)
	RunningTime() time.Duration
	Format() string // e.g., "MP4", "WMV"
	Reslution() (x, y int)
}

这些接口不止是一种有用的方式来分组相关的具体类型和表示他们之间的共同特点。
我们后面可能会发现其他的分组。
举例，如果我们发现我们需要以同样的方式处理Audio和Video，我们可以定义一个Streamer接口来代表他们之间相同的部分而不必对已经存在的类型做改变。
type Streamer interface{
	Stream() (io.ReadCloser, error) //这是 Video和Audio接口中共同有的方法
	RunningTime() time.Duration //这是 Video和Audio接口中共同有的方法
	Format() string //这是 Video和Audio接口中共同有的方法
}

每一个具体类型的组基于他们相同的行为可以表示成一个接口类型。
不像基于类的语言，他们一个类实现的接口集合需要进行显式的定义，
在Go语言中我们可以在需要的时候定义一个新的抽象或者特定特定的组，
而不需要修改具体类型的定义。
当具体的类型来自不同的作者时这种方式会特别有用。
当然也确实没有必要在具体的类型中指出这些共性。


7.4 flag.Value接口（p241）
10. 在本节，我们会学到另一个标准的接口类型flag.Value是怎么帮助命令行标记定义新的符号的。
思考下面这个会休眠特定时间的程序：

var period = flag.Duration("period", 1*time.Second, "sleep period")

func main(){
	flag.Parse()
	fmt.Printf("Sleepping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

在它休眠前它会打印出休眠的时间周期。
fmt包调用time.Duration的String方法打印这个时间周期是以用户友好的注解方式，而不是一个纳秒数字：

// $ go build ...sleep
// $ ./sleep
// Sleeping for 1s... 

默认情况下，休眠周期是一秒，但是可以通过 -period这个命令标记来控制。
flag.Duration函数创建一个time.Duration类型的标记变量并且允许用户通过多种用户友好的方式来设置这个变量的大小，
这种方式还包括和String方法相同的符号排版形式。
这种对称设计使得用户交互良好。

$ ./sleep -period 50ms
Sleeping for 50sm...
$ ./sleep -period 2m30s
Sleeping for 2m30s...
$ ./sleep -period 1.5h
Sleeping for 1h30m0s...
$ ./sleep -period "1 day"
invalid value "1 day" for flag -period: time: invalid duration 1 day

12. 因为时间周期标记值非常有用，所以这个特性被构建到了flag包中；但是我们为我们自己的数据类型定义新的标记符号是简单容易的。
我们只需要定义一个实现flag.Value接口的类型，如下：
package flag

// Value is the interface to the value stored in a flag.
type Value interface{
	String() string
	Set(string) error
}

String方法格式化标记的值用在命令行帮组消息中；这样每一个flag.Value也是一个fmt.Stringer。
Set方法解析它的字符串参数并且更新标记变量的值。
实际上，Set方法和String是两个相反的操作，
所以做好的办法就是对他们使用相同的注解方式。

让我们定义个允许通过摄氏度或者华氏度变换的形式制定温度的celsiusFlag类型。
注意celsiusFlag内嵌了一个Celsius类型（2.5），因此不用实现本身就已经有String方法了。
为了实现flag.Value，我们只需要定义Set方法。

示例代码

13. 调用fmt.Sscanf函数，从s中解析一个浮点数（value），一个字符串（unit）。
虽然通常必须检查Sscanf的错误返回，但是在这个例子中我们不需要因为如果有错误发生，就没有switch case会匹配到。

14. 下面的CelsiusFlag函数将所有逻辑都封装在一起。
他返回一个内嵌在celsiusFlag变量f中的Celsius指针给调用者。
Celsius字段是一个会通过Set方法在标记处理的过程中更新的变量。
调用Var方法将标记加入应用的命令行标记集合中，有异常复杂命令接口的全局变量flag.CommandLine.Programs可能有几个这个类型的变量。
调用Var方法将一个celsiusFlag参数赋值给一个flag.Value参数，导致编译器去检查celsiusFlag是否有必须的方法。


7.5 接口值 （p244）
15. 概念上讲一个接口的值、接口值，有两个部分组成，一个具体的类型和那个类型的值。
他们被成为接口的动态类型和动态值。
由于像Go语言这种景泰类型的语言，类型是编译期的概念；
因此一个类型不是一个值。
在我们的概念模型中，一些提供每个类型信息的值被称为类型描述符，比如类型的名称和方法。
在一个接口值中，类型部分代表与之相关类型的描述符。

16. 下面4个语句中，变量w得到了3个不同的值。（开始和最后的值是相同的）

var w io.Writer //io.Writer是一个接口类型
w = os.Stdout // os.Stdout也是个接口类型？
w = new(bytes.Buffer)
w = nil

让我们进一步观察在每一个语句后的w变量的值和动态行为。
第一个语句定义了变量w:

var w io.Writer

在Go语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。
对于一个接口的零值就是它的类型和值的部分都是nil。
		w
type	nil
value	nil

一个接口值基于它的动态类型被描述为空或非空，所以这是一个空的接口值。
你可以通过使用 w == nil 或 w != nil 来判读接口值是否为空。
调用一个空接口值上的任意方法都会产生panic:
w.Write([]byte("hello")) // panic: nil pointer dereference

第二个语句将*os.File类型的值赋给变量w:
w = os.Stdout

这个赋值过程调用了一个具体类型到接口类型的隐式转换，这和显式的使用io.Writer(os.Stdout)是等价的。
这类转换不管是显式的还是隐式的，都会刻画出操作到的类型和值。
这个接口值的动态类型被设为*os.Stdout指针的类型描述符，
他的动态值持有os.Stdout的拷贝；这是一个代表处理标准输出的os.File类型变量的指针。
		w
type	*os.File		os.File
value	--------------->fd int = 1(stdout)

调用一个包含*os.File类型指针的接口值的Write方法，使得 (*os.File).Write方法被调用。
这个调用输出"hello"
w.Write([]byte("hello")) // "hello"

通常在编译期，我们不知道接口值的动态类型是什么，
所以一个接口上的调用必须使用动态分配。
因为不是直接进行调用，所以编译期必须把代码生成在类型描述符发方法Write上，
然后间接调用那个地址。
这样调用的接收者是一个接口动态值的拷贝，os.Stdout。
效果和下面这个直接调用一样：
os.Stdout.Write([]byte("hello")) // "hello"

第三个语句给接口值附了一个*bytes.Buffer类型的值。
w = new(bytes.Buffer)

现在动态类型是*bytes.Buffer，并且动态值是一个指向新分配的缓冲区的指针。

		w						bytes.Buffer
type	*bytes.Buffer		--->data []byte
value	--------------------|	...

Write方法的调用也使用了和之前一样的机制：
w.Write([]byte("hello")) // writes "hello" to the bytes.Buffers

这次类型描述符是*bytes.Buffer，所以调用了 (*bytes.Buffer).Write方法，
并且接收者是该缓冲区的地址。
这个调用把字符串"hello"添加到缓冲区中。

最后，第四个语句将nil赋值给了接口值：
w = nil

这个重置将它所有的部分都设为nil，把变量w回复到和它之前定义时相同的状态图，在7.1中可以看到。

17. 一个接口值可以持有任意大的动态值。
例如，表示时间实例的time.Time类型，这个类型有几个对外不公开的字段。
我们从他上面创建一个接口值。

var x interface{} = time.Now()

结果可能和图7.4 相似。
从概念上讲，不论接口值多大，动态值总是可以容下它。
（这只是一个概念上的模型；具体实现可能会非常不同）
		x
type	time.Time
value	sec:	xxxxxxx
		nsec:	xxxxxx
		loc:	"UTC"

18. 接口值可以使用== 和 != 来进行比较。
两个接口值相等仅当他们都是nil值或者他们的动态值类型相同兵器动态值也根据这个动态类型的==操作相等。
因为接口值是可比较的，所以他们可以用在map的键或者switch语句的操作数。

然而， 如果两个接口值的动态类型相同，但是这个动态类型是不可比较的（比如切片），
将他们进行比较就会失败并且panic：
var x interface{} = []int{1, 2, 3}
fmt.Println(x == x) //panic： comparing uncomparable type []int

考虑到这点，接口类型是非常与众不同的。
其他类型要么是安全的可比较类型（如基本类型和指针）
要么是完全不可比较的类型（如切片，映射类型，和函数）
但是在比较接口值或者包含了接口值的聚合类型时，我们必须要意识到潜在的panic。
同样的风险也存在于使用接口作为map的键或switch的操作数。
只能比较你非常确定它们的动态值是可比较类型的接口值。

19. 当我们处理错误或者调试的过程中，得知接口值的动态类型是非常有帮助的。
所以我们使用fmt包的%T动作：

var w io.Writer
fmt.Printf("%T\n", w) // "<nil>"
w = os.Stdout
fmt.Printf("%T\n", w) // "*os.File"
w = new(bytes.Buffer)
fmt.Printf("%T\n", w) // "*bytes.Buffer"

在fmt包内部，使用反射来获取接口动态类型的名称。
我们会在第12章学到反射相关的知识。

7.5.1（p247）
20. 警告： 一个包含nil指针的接口不是nil接口
一个不包含任何值的nil接口值和一个刚好包含nil指针的接口值是不同的。
这个细微区别产生了一个容易绊倒每个Go程序员的陷阱。
思考下面的程序。
当debug变量设置为true时，main函数会将f函数的输出收集到一个bytes.Buffer类型中。

const debug = true

func main(){
	var buf *bytes.Buffer
	if debug{
		buf = new(bytes.Buffer) // enable collection of output
	}
	f(buf) // NOTE: subtly incorrect!
	if debug{
		// ...use buf...
	}
}

// If out is non-nil, output will be written to it..
func f(out io.Writer){
	// ...do something...
	if out != nil{
		out.Write([]byte("done!\n"))
	}
}

我们可能会预计当把变量debug设置为false时可以进制对输出的收集，
但是实际上在out.Write方法调用是程序发生了panic:
if out != nil{
	out.Write([]byte("done!\n")) // panic: nil pointer dereference.
}

当main函数调用函数f时，它给f函数的out参数附了一个*bytes.Buffer的空指针，
所以out的动态值是nil。
然而，他的动态类型是*bytes.Buffer，意思就是out变量是一个包含空指针值的非空接口，
所以防御性检查out!=nil的结果依然是true。
		w
type	*bytes.Buffer  //类型不是空的，
value	nil	//值是空的

动态分配机制依然决定 (*btyes.Buffer).Write的方法会被调用，但是这次的接收者的值是nil。
对于一些如*os.File的类型，nil是一个有效的接收者（6.2.1），但是*bytes.Buffer类型不在这些类型中。
这个方法会被调用，但是当它尝试去获取缓冲区时会发生panic。

问题在于尽管一个nil的*bytes.Buffer指针有实现这个接口的方法，
他也不满足这个接口具体的行为上的要求。
特别是这个调用违反了(*bytes.Buffer).Write方法的接收者非空的隐含先觉条件，
所以将nil指针赋值给这个接口是错误的。
解决方案就是将main函数中的变量bug的类型改为io.Writer,
因此可以避免一开始就将一个不完全的值赋值给这个接口。

var buf io.Writer
if debug{
	buf = new(bytes.Buffer) // enable collection of output
}
f(buf) // OK

现在我们已经把接口值的技巧都讲完了，
让我们来看更多的一些在Go标准库中重要的接口类型。
在下面的三章中，我们会看到接口类型是怎样用在排序，web服务，错误处理中的。

7.6 sort.Interface接口（p249）

排序操作和字符串格式化一样是很多程序经常使用的操作。
尽管一个最短的快排序只要15行就可以搞定，但是一个健壮的实现需要更多的代码，并且我们不希望每次我们需要的时候都重写或者拷贝这些代码。

幸运的是，sort包内置的提供了根据一些排序函数来对任何序列排序的功能。
他的设计非常独到。在很多语言中，排序算法都是和序列数据类型关联，同时排序函数和具体类型元素关联。
相比之下，Go语言的sort.Sort函数不会对具体的序列和他的元素做任何假设。
相反，他使用了一个接口类型sort.Interface来指定通用的排序算法和可能被排序到的序列类型之间的约定。
这个接口的实现由序列的具体表示和它希望排序的元素决定，序列的表示经常是一个切片。

21. 一个内置的排序算法需要知道3个东西：序列的长度，表示两个元素比较的结果，一种交换两个元素的方式；这就是sort.Interface的三个方式：

package sort

type Interface interface{
	Len() int
	Less(i, j int) bool // i, j are indices of sequence elements. i、 j是序列元素的指数
	Swap(i, j int) 
}

为了对序列进行排序，我们需要定义一个实现了这三个方法的类型，然后对这个类型的一个实例应用sort.Sort函数。
思考对一个字符串切片进行排序，这可能是最简单的例子了。
下面是这个新的类型StringSlice和他的Len，Less和Swap方法

type StringSlice []string
func (p StringSlice) Len() int {return len(p)}
func (p StringSlice) Less(i, j int) bool {return p[i]<p[i]}
func (p StringSlice) Swap(i, j int) {p[i], p[j] = p[j], p[i]}

现在我们可以通过下面这样将一个切片转化为一个StringSlice类型来进行排序：
sort.Sort(StringSlice(names)) // 把name转变为StringSlice类型，然后调用sort包的Sort函数

这个转换得到一个相同长度，容量，和基于names数组的切片值；
并且这个切片值的类型有三个排序需要的方法。

对字符串切片的排序是很常用的需要，所以sort包提供了StringSlice类型，
也提供了Strings函数能让上面这些调用简化成sort.Strings(names)。

这里用到的技术很容易适用到其他排序序列中，例如我们可以忽略大些或者含有特殊的字符。
（本书使用Go程序对索引词和页码进行排序也用到了这个技术，对罗马数字做了额外逻辑处理）
对于更复杂的排序，我们使用相同个的方法，但是会用更复杂的数据结构和更复杂的实现sort.Interface的方法。

我们会运行上面的例子来对一个表格中的音乐播放列表进行排序。
每个track都是单独的一个行，每一列都是这个track的属性像艺术家，标题，和运行时间。
想象一个图形用户界面来呈现这个表格，并且点击一个属性的顶部会使这个列表按照这个属性进行排序；
再一次点击相同属性的顶部会进行逆向排序。
让我们看下每个点击会发生什么相应。

下面的变量tracks包好了一个播放列表。
（One of authorshipapologizes for the other author’s musical tastes.）
每个元素都不是Track本身而是指向他的指针。
尽管我们在下面的代码中直接存储Tracks也可以工作，sort函数会交换很多对元素，
所以如果每个元素都是指针会更快而不是全部Track类型，指针是一个机器字码长度而Track类型可能是八个或更多。

示例代码

22. printTracks函数将播放列表打印成一个列表。
一个图形化的展示可能会更好点，但是这个小程序用text/tabwrite包来生成一个列是整齐对其和隔开的表格，像下面展示的这样。
注意到*tabwrite.Writer是满足io.Writer接口的。
他会收集每一片写向它的数据；
他的Flush方法会格式化表格并且将它写向os.Stdout （标准输出）。

func printTracks(track []*Track){
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0) //这什么意思？
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}

为了能按照Artist字段对播放列表进行排序，我们会像对StringSlice那样定义一个新的带有必须Len，Less和Swap方法的切片类型。

type byArtist []*Track
func (x byArtist) Len() int { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

为了调用通用的排序程序，我们必须现将tracks转换为新的byArtist类型，
他定义了具体的排序：

sort.Sort(byArtist(tracks))

在按照artist对这个切片进行排序后，printTrack的输出如下：

示例结果

如果用户第二次请求“按照artist排序”，我们会对tracks进行逆向排序。
然而我们不需要定义一个有颠倒Less方法的新类型byReverseAritst，
因为sort包中提供了Reverse函数将排序顺序转换成逆序。

sort.Sort(sort.Reverse(byArtist(tracks)))

在按照artist对这个切片进行逆向排序后，printTrack的输出如下：

示例结果（逆向的排序）

sort.Reverse函数值得进行更近一步的学习因为它使用了（6.3）章中的组合，
这是一个重要的思路。
sort包含定义了一个不公开的struct类型reverse，它嵌入了一个sort.Interface。
reverse的Less方法调用了内嵌的sort.Interface值的Less方法，但是通过交换索引的方式使排序结果变成逆序。

package sort

type reverse struct{Interface} // that is, sort.Interface
func (r reverse) Less(i, j int) bool {return r.Interface.Less(j,i)}
func Reverse(data Interface) Interface{return reverse{data}}

reverse的另外两个方法Len和Swap隐式地由原有内嵌的sort.Interface提供。
因为reverse是一个不公开的类型，所以导出函数Reverse函数返回一个包含原有sort.Interface值的reverse类型实例。

为了可以按照不同的列进行排序，我们必须定义一个新的类型例如byYear：
type byYear []*Track
func (x byYear) Len() int {return len(x)}
func (x byYear) Less(i, j int) bool {return x[i].Year< x[j].Year}
func (x byYear) Swap(i, j int) {x[i], x[j] = x[j], x[i]}

在使用sort.Sort(byYear(tracks))按照年对tracks进行排序后，printTrack展示了一个按时间先后顺序的列表：

结果示例：

对于我们需要的每个切片元素类型和每个排序函数，我们需要定义一个新的sort.Interface实现。
如你所见，Len和Swap方法对于所有的切片类型都有相同的定义。
下个例子，具体的类型customSort会将一个切片和函数结合，是我们只需要写比较函数就可以定义一个新的排序。
顺便说下，实现了sort.Interface的具体类型不一定是切片类型；
customSort是一个结构体类型。

type customSort struct{
	t []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int
func (x customSort) Less(i, j int) bool {return x.less(x.t[i], x.t[j])}
func (x customSort) Swap(i, j int) {x.t[i],x.t[j] = x.t[j], x.t[i]}

让我们定义一个多层的排序函数，他主要的排序键是标题，第二个键是年，第三个键是运行时间Length。
下面是该排序的调用，其中这个排序使用了匿名排序函数:

sort.Sort(customSort{tracks, func(x, y *Track)bool{  // 这个函数就是less函数
	if x.Title !=y.Title{
		return x.Title < y.Title
	}
	if x.Year != y.Year{
		return x.Year < y.Year
	}
	if x.Length != y.Length{
		return x.Length < y.Length
	}
	return false
}})

这下面是排序的结果。
注意到两个标题是"Go"是track按照标题排序是相同的顺序，但是在按照year排序上更久的那个track优先。

结果示例

尽管对长度为n的序列排序需要0（n log n）次比较操作，检查一个序列是否已经有序至少需要n-1次比较。
sort包中的IsSorted函数帮我们做这样的检查。
像sort.Sort一样，它也使用sort.Interface对这个序列和它的排序函数进行抽象，
但是它从不会调用Swap方法：这段代码示范了IntsAreSorted和Ints函数和IntSlice类型的使用：

values := []int{3,1,4,1}
fmt.Println(sort.IntsAreSorted(values)) // "false"
sort.Ints(values)
fmt.Println(values) //"[1 1 3 4]"
fmt.Println(sort.IntsAreSorted(value)) // "true"
sort.Sort(sort.Reverse(sort.IntSlice(values)))
fmt.Pritnln(values) // "[4 3 1 1]"
fmt.Println(sort.IntsAreSorted(values)) // "false"

为了使用方便，sort包为[]int，[]string 和[]float64 的正常排序提供了特定版本的函数和类型。
对于其他类型，例如[]int64或者[]uint，尽管路径也很简单，还是依赖我们自己实现。

7.7 http.Handler接口（p255）
在第一章中，我们粗略的了解了怎么用net/http包去实现网络客户端（1.5）和服务器（1.7）。
在这个小节中，我们会对那些基于http.Handler接口的服务器API做更进一步的学习：
net/http

package http

type Handler interface{
	ServeHTTP(w ResponseWriter,, r *Request)
}

func ListenAndServe(address string, h Handler) error

ListenAndServe函数需要一个例如“localhost:8000”的服务器地址，和一个所有请求都可以分派的Handler接口实例。
他会一直运行，知道这个服务器因为一个错误而失败（或者启动失败），他的返回值一定是一个非空的错误。

想象一个电子商务网站，为了销售它的数据库将它物品的价格映射成美元。
下面这个程序可能是想到最简单的实现了。
他将库存清单模型化为一个命名为database的map类型，我们给这个类型一个ServeHttp方法，这样它可以满足http.Handler接口。
这个handler会遍历整个map并输出物品信息。

示例代码

如果我们启动这个服务，

$ go .../http1
$ ./http1 & // &是什么意思？

然后用1.5节中的获取程序（如果你更喜欢可以使用web浏览器）来链接服务器，我们得到下面的输出：
$ go build .../fetch
$ ./fetch http://localhost:8000
shoes: $50.00
socks: $5.00

目前为止，这个服务器不考虑URL只能为每个请求列出它全部的库存清单。(ServeHTTP方法中没有req参数的原因？)
更真实的服务器会定义多个不同的URL，每一个都会触发一个不同的行为。
让我们使用/list来调用已经存在的这个行为并且增加另一个、price调用表明单个货品的价格，
像这样/price?item=socks来制定一个请求参数。

示例代码

现在handler基于URL的路径部分（rep.URL.Path)来决定执行什么逻辑。
如果这个handler不能识别这个路径，他会通过调用w.WriteHeader(http.StatusNotFound)返回客户端一个HTTP错误；
这个检查应该在向w写入任何值前完成。
（顺便提一下，http.ResponseWriter是另一个接口。它在io.Writer上增加了发送HTTP相应的方法）
等效的，我们可以使用使用的http.Error函数：
msg := fmt.Sprintf("no such page: %s\n", req.URL)
http.Error(w, msg, http.StatusNotFound) // 404

/price 的case会调用URL的Query方法来将HTTP请求参数解析为一个map，
或者更准确地说一个net/url包中url.Value （6.2.1）类型的多重映射。
然后找到第一个item参数并查找它的价格。
如果这个货品没有找到会返回一个错误。

这里是一个和新服务器会话的例子：（没有说怎么关闭服务器，所以说新服务器？）

$ go build gopl.io/ch7/http2
$ go build gopl.io/ch1/fetch
$ ./http2 &
$ ./fetch http://localhost:8000/list
shoes: $50.00
socks: $5.00
$ ./fetch http://localhost:8000/price?item=socks
$5.00
$ ./fetch http://localhost:8000/price?item=shoes
$50.00
$ ./fetch http://localhost:8000/price?item=hat
no such item: "hat"
$ ./fetch http://localhost:8000/help
no such page: /help

显然我们可以继续向ServeHTTP方法中添加case，但在一个实际的应用中，
将每个case中的逻辑定义到一个分开的方法或函数中会很使用。
此外，相近的URL可能需要相似的逻辑；
例如几个图片文件可能有形如/images/ *.png的URL。
因为这些原因，net/http包提供了一个请求多路器ServeMux来简化URL和handlers的联系。
一个ServeMux将一批http.Handler聚集到一个单一的http.Handler中。
再一次，我们可以看到满足同一接口的不同类型是可替换的：
web服务器将请求指派给任意的http.Handler而不需要考虑它后面的具体类型。

对于更复杂的应用，一些ServeMux可以通过组合来处理更加错综复杂的路由需求。
Go语言目前没哟一个权威的web框架，就像Ruby语言有Rails和python有Django。
这并不是说这样的框架不存在，而是Go语言标准库中的构建模块就已经非常灵活以至于这些框架都是不必要的。
此外，尽管在一个项目早期使用框架是非常方便的，但是他们带来额外的复杂度会使长期的维护更加困难。

在下面的程序中，我们床架一个ServeMux并且使用它将URL和相应处理/list和/price操作的handler联系起来，
这些操作逻辑都已被分到不同的方法中。
然后我们在调用ListenAndServe函数中使用ServeMux最为主要handler。

示例代码

让我们关注这两个注册到handlers上的调用。
第一个db.list是一个方法值（6.4），它是下面这个类型的值。

func(w http.ResponseWriter, req *http.Reqwuest)

也就是说db.list的调用会援引一个接收者是db的database.list方法。
所以db.list是一个实现了handler类似行为的函数，
但是因为它没有方法，
所以他不满足http.Handler接口并且不能直接传给mux.Handle。
//没懂。。。？？？

语句http.HandlerFunc(db.list)是一个转换而非一个函数调用，
因为http.HandlerFunc是一个类型。
他有如下定义：
net/http
package http

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request){
	f(w, r)
}

HandlerFunc显示了在Go语言接口机制中一些不同寻常的特点。
这是一个有实现了接口http.Handler方法的函数类型。
ServeHTTP方法的行为调用了它本身的函数。
因此，HandlerFunc是一个让函数值满足一个接口的适配器，这里函数和这个接口仅有的方法有相同的函数签名。
实际上，这个技巧让一个单一的类型例如database以多种方式满足http.Handler接口：
一种通过它的list方法，一种通过他的price方法等等。

因为handler通过这种方法注册非常普遍，ServeMux有一个方便的HandleFunc方法，
它帮我们简化handler注册代码成这样：

http3a
mux.HandleFunc("/list", db.list)
mux.HandleFunc("/price", db.list)

从上面的代码很容易看出应该怎么都贱一个程序，他有两个不同的web服务器鉴定不同的端口的，
并且定义不同的URL将它们指派到不同的handler。
我们只要构建另外一个ServeMux并且再调用一次ListenAndServe （可能并行的）。
但是在大多数程序中，一个web服务器就足够了。
此外，在一个应用程序的多个文件中定义HTTP handler也是非常经典的，
如果它们必须全部都显示的注册到这个应用的ServeMux实例上会比较麻烦。

所以为了方便，net/http包提供了一个全局的ServeMux实例DefaultServerMux和
包级别的htp.Handle和http.HandleFunc函数。
现在，为了使用DefaultServeMux作为服务器的主handler，
我们不需要将他传给ListenAndServe函数；
nil值就可以工作。

然后服务器的主函数可以简化成：

示例代码

最后，一个重要的提示：就像我们在1.7节中提到的，web服务器在一个新的协程中调用每一个handler，
所以当handler获取其他协程或者这个handler本身的其他请求也可以访问的变量时一定要使用预防措施比如锁机制。
我们后面的两章中讲到并发相关的知识。

7.8 error接口（p261）
从本书的开始，我们就已经创建和使用过神秘的预定义error类型，
而且没有解释它究竟是什么。
实际上他就是interface类型，这个类型有一个返回错误信息的单一方法：
type error interface{
	Error()string
}

创建一个error最简单的方法就是调用errors.New函数，他会根据传入的错误信息，返回一个新的error。
整个errors包仅只有4行：
package errors
func New(text string) error {return &errorString{text}}

type errorString struct {text string}

func (e *errorString) Error() string {return e.text}

承载errorString的类型是一个结构体而非一个字符串，这是为了保护它表示的错误避免粗心（或有意）的更新。
并且因为是指针类型*errorString满足error接口而非errorString类型，
所以每个New函数的调用都分配了一个独特的和其他错误不相同的实例。
我们也不想要重要的error例如EOF和一个刚好有相同错误消息的error比较后相等。

fmt.Println(errors.New("EOF") == errors.New("EOF")) // "false"

调用errors.Now函数是非常稀少的，因为有一个方便的封装函数fmt.Errorf，
他还会处理字符串格式化。
我们曾多次在第5章中用到他。

package fmt

import "errors"

func Errorf(format string, args ...interface{})error{
	return errors.New(Sprintf(format, args...))
}

虽然*errorString可能是最简单的错误模型，但远非只有它一个。
例如，syscall包提供了Go语言底层系统调用API。
在多个平台上，他定义一个实现error接口的数字类型Errno，
并且在Unix平台上，Errno的Error方法会从一个字符串表中查找错误信息，如下面展示的这样：

package syscall

type Errno uintptr // operating system error code

var errors = [...]string{
	1:	"operation no permitted", // EPERM
	2:	"no such file or directory", // ENOENT
	3:	"no such process", // ESRCH
	// ...
}

func (e Errno) Error() string{
	if 0<=int(e) && int(e) < len(errors){
		return errors[e]
	}
	return fmt.Sprintf("errno %d", e)
}

下面的语句创建了一个持有Errno值为2的接口值，表示POSIX ENOENT状况：

var err error = syscall.Errno(2)
fmt.Println(err.Error()) // "no such file or directory"
fmt.Println(err) // "no such file or directory"

err的值图形化的呈现：
		err
type	syscall.Errno
value	2

Errno是一个系统调用错误的高效表示方式，它通过一个有限的集合进行描述，并且他满足标准的错误接口。
我们会在7.11节了解到其他满足这个接口的类型。

7.9 示例：表达式求值（p263）
在本节中，我们会构建一个简单算术表达式的求值器。
我们将使用一个接口Expr来表示Go语言中任意的表达式。
现在这个接口不需要有方法，但是我们后面会为它增加一些。

// An Expr is arithmetic expression.
type Expr interface{}

我们的表达式语言有浮点数符号（小数点）；二元操作符+，-，*，和/；
一元操作符-x和+x；
调用 pow(x,y),sin(s),和 sqrt(x)的函数；
例如x和pi的变量；
当然也有括号和标准的优先级运算符。
所有的值都是float64类型。
这下面是一些表达式的例子：

sqrt(A / pi)
pow(x, 3) + pow(y, 3)
(F - 32) *5 / 9

下面的五个具体类型表示了具体的表达式类型。
Var类型表示对一个变量的引用。
（我们很快会知道为什么它可以被输出）
literal类型表示一个浮点型常量。
unary和binary类型表示有一到两个运算对象的运算符表达式，这些操作数可以是任意的Expr类型。
call类型表示对函数的调用；
我们限制他的字段只能是pow，sin或者sprt。

// A Var identifies a variable, e.g., x.
type Var string
// A literal is a numeric constant, e.g., 3.141.
type literal float64
// A unary represents a unary operator expression, e.g., -x.
type unary struct {
	op rune // one of '+', '-'
	x Expr
}
// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
	op rune // one of '+', '-', '*', '/'
	x,? , y Expr
}
// A call represents a function call expression, e.g., sin(x).
type call struct {
	fn string // one of "pow", "sin", "sqrt"
	args []Expr
}


为了计算一个包含变量的表达式，我们需要一个environment变量将变量的名字映射成对应的值：

type Env map[Var]float64

我们也需要每个表达式去定义一个Eval方法，这个方法会根据给定的environment变量返回表达式的值。
因为每个表达式都必须提供这个方法，我们将他加入到Expr接口中。
这个包只会对外公开Expr，Env，和Var类型。
调用方不需要获取其他的表达式类型就可以使用这个求值器。

type Expr interface{
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
}

下面给大家展示一个具体的Eval方法。
Var类型的这个方法对一个environment变量进行查找，
如果这个变量没有在environment中定义过这个方法,会返回一个零值，
literal 类型的这个方法简单的返回他真实的值。

func (v Var) Eval (env Env) float64{
	return env[v]
}

func (l literal) Eval(_ Env) float64{
	return float64(l) //把l转化成float64类型
}

unary 和 binary 的Eval方法会递归的计算他的运算对象，然后将运算符op作用到他们上。
我们不讲被零或无穷数除作为一个错误，因为他们都会产生一个固定的结果无限。
最后，call的这个方法会计算对于pow，sin，或者sqrt函数的参数值，然后调用对应在math包中的函数。

func (u unary) Eval(env Env) float64{
	switch u.op{
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env)float64{
	switch b.op{
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval (env Env) float64{
	switch c.fn{
	case "pow"://次方运算
		//这里为什么是双引号
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

一些方法会失败。
例如一个call表达式可能未知的函数或者错误的参数个数。
用一个无效的运算符如 ! 或者 < 去构建一个unary或者binary表达式也是可能发生的
（尽管下面提到的Parse函数不会这样做）
这些错误会让Eval方法panic。
其他的错误，像计算一个没有在environment变量中出现过的Var，
只会让Eval方法返回一个错误的结果。
所有的这些错误都可以通过在计算前检查Expr来发现。
这是我们接下来要讲的Check方法的工作，但是让我们先测试Eval方法。

下面的TestEval函数是对evaluator的一个测试。
它使用了我们会在第11章讲解的testing包，但是现在知道调用t.Errof会报告一个错误就足够了。
这个函数循环遍历一个表格中的输入，
这个表格中定义了三个表达式和针对每个表达式不同的环境变量。
第一个表达式根据给定圆的面积Ａ计算他的半径，
第二个表达式通过两个变量x和y计算两个立方体的体积之和，
第三个表达式将华氏温度F转化成摄氏度。

示例代码

对于表格中的每一条记录，这个测试会解析它的表达式然后在环境变量中计算他，输出结果。
这里我们没有空间来展示Parse函数，但是如果你使用go get下载这个包你就可以看到这个函数。
go test命名会运行一个包的测试用例：

$ go test -v ...eval

这个-v标识可以让我们看到测试用例打印的输出；
正常情况下像这个一样成功的测试用例会阻止打印结果的输出。
这里是测试用例fmt.Printf语句的输出：
sqrt(A / pi)
	map[A:87616 pi:3.141592653589793] => 167

pow(x, 3) + pow(y, 3)
	map[x:12 y:1] => 1729
	map[x:9 y:10] => 1729

5 / 9 * (F - 32)
	map[F:-40] => -40
	map[F:32] => 0
	map[F:212] => 100

幸运的是目前为止所有的输入都是适合的格式，但是我们的运气不可能一直都有。
甚至在解释型语言中，为了静态错误检查语法是非常常见的；
静态错误就是不用运行程序就可以检测出来的错误。
通过将静态检查和动态的部分分开，
我们可以快速的检查错误并且对于多次检查只执行一次而不是每次表达式计算的时候都进行检查。

让我们往Expr接口中增加另一个方法。Check方法在一个表达式语义树检查出静态错误。
我们马上会说明他的vars参数。

type Expr interface{
	Eval(env Env)float64
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]bool)error
}

具体的Check方法展示在下面。
literal和Var类型的计算不可能失败，
所以这些类型的Check方法会返回一个nil值。
对于unary和binary的Cheak方法会首先检查操作符是否有效，
然后递归的检查运算单元。
相似地对于call的这个方法首先检查调用的函数是否已知并且有没有正确个数的参数，
然后递归的检查每一个参数。

示例代码

我们在两个组中有选择地列出有问题的输入和他们得出的错误。
Parse函数（这里没有出现）会报出一个语法错误和Check函数会报出语义错误。

示例结果

Check方法的参数是一个Var类型的集合，这个集合聚集从表达式中找到的变量名。
为了保证成功的计算，这些变量的每一个都必须出现在环境变量中。
从逻辑上讲，这个集合就是调用Check方法返回的结果，
但是因为这个方法是递归调用的，所以对于Check方法填充结果到一个作为参数传入的集合中会更加的方便。
调用方在初始调用时必须提供一个空的集合。

在第3.2节中，我们绘制了一个在编译器才确定的函数 f(x, y)。
现在我们可以解析，检查和计算在字符串中的表达式，我们可以构建一个在运行时从客户端接收表达式的web应用并且他会绘制这个函数的表示的曲面。
我们可以使用集合vars来检查表达式是否是一个只有两个变量，
x和y的函数————实际上是3个， 
因为我们为了方便会提供半径大小r。
并且我们会在计算前使用Check方法拒绝有格式问题的表达式，
这样我们就不会在下面函数的40000个计算过程（100*100个栅格，每一个有4个角）重复这些检查。

这个ParseAndCheck函数混合了解析和检查步骤的过程：

示例代码

为了编写这个web应用，所有我们需要做的就是下面这个plot函数，
这个函数有和http.HandlerFunc相似的签名：

示例代码
根据相应公式，生成相应图片

这个plot函数解析和检查在HTTP请求中指定的表达式，
并且用它来创建一个两个变量的匿名函数。
这个匿名函数和来自原来surface-plotting程序中的固定函数f有相同的签名，
但是它计算一个用户提供的表达式。
环境变量中定义了x, y和半径r。
最后plot调用surface函数，他就是ch3/surface中的主要函数，修改后他可以接收plot中的函数和输出io.Writer作为参数，
而不是使用固定的函数f和os.Stdout。
显示了通过程序产生的3个曲面。

7.10 类型断言（p273）
类型断言是一个使用在接口值上的操作。
语法上他看起来像x.(T),被成为断言类型，这里x表示一个接口的类型和T表示一个类型。
一个类型断言检查他操作对象的动态类型是否和断言的类型匹配。

这里有两种可能。
第一种，如果断言的类型T是一个具体类型，
然后类型断言检查x的动态类型是否和T相同。
如果这个检查成功了，类型断言的结果是x的动态值，
当然他的类型是T。
换句话说，具体类型的类型断言从他的操作对象中获得具体的值。
如果检查失败，接下来这个操作会抛出panic。
例如：
var w io.Writer // 这个是个接口
w = os.Stdout //os.Stdout是接口值
f := w.(*os.File) // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer

第二种，如果相反断言的类型T是一个接口类型，
然后类型断言检查是否x的动态类型满足T。
如果这个检查成功了，动态值没有获取到；
这个结果仍然是一个有相同类型和值部分的接口值，但是结果有类型T。
换句话说，对一个接口类型的类型断言改变了类型的表述方式，
改变了可以获取的方法集合（通常更大），但是他保护了接口值内部的动态类型和值的部分。

在下面的第一个类型断言后，w和rw都持有os.Stdout因此他们每个有一个动态类型*os.File,
但是变量w是一个io.Writer类型只对外公开出文件的Write方法，然而rw变量也只公开他的Read方法。

var w io.Writer// 这个是个接口
w = os.Stdout//os.Stdout是接口值
rw :=w.(io.ReadWriter) // success: *os.File has both Read and Write
// io.ReadWriter也是个接口
w = new(ByetCounter)
rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method

如果断言操作的对象是一个nil接口值，那么不论被断言的类型是什么这个类型断言都会失败。
我们几乎不需要对一个更少限制性的接口类型（更少的方法集合）做断言，
因为他表现的就像赋值操作一样，除了对于nil接口值的情况。
w = rw // io.ReadWriter is assignable to io.Writer
w = rw.(io.Writer) // fails only if rw == nil
// w = nil.(io.Writer) ??这个意思？

经常地我们对一个接口值的动态类型是不确定的，并且我们更愿意去检验它是不是一些特定的类型。
如果类型断言出现在一个预期有两个结果的操作中，例如如下的定义，
这个操作不会在失败的时候发生panic，但是代替地返回一个额外的第二个结果，
这个结果是一个标识成功的布尔值：
var w io.Writer = os.Stdout
f, ok := w.(os.File) // success: ok, f == os.Stdout
b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil

第二个结果常规地赋值给一个命名为ok的变量。
如果这个操作失败了，那么ok就是false值，
第一个结果等于被断言类型的零值，在这个例子中就是一个nil的*bytes.Buffer类型。

这个ok结果经常立即用于决定程序下面做什么。
if语句的扩展格式让这个编的很简洁：
if f, ok := w.(*os,File); ok{
	// ...use f...
}

当类型断言的操作对象是一个变量，你有时会看见原来的变量名重用而不是生命吧一个新的本地变量，
这个重用的变量会覆盖原来的值，如下面这样：
if w, ok := w.(*os.File); ok{
	// ...use w ...
}


7.11 基于类型断言区别错误类型（p275）
思考在os包中文件操作返回的错误集合。
I/O可以因为任何数量的原因失败，
但是有三种经常的错误必须进行不同的处理：
1.文件已经存在（对于创建操作）
2.找不到文件（对于读取操作）
3.权限拒绝。
os包中提供了这三个帮助函数来对给定的错误值表示的失败进行分类：
package os

func IsExist(err error) bool
func IsNotWxist(err error)bool
func IsPermission(err error)bool

对这些判断的一个缺乏经验的实现可能会去检查错误消息是否包含了特定的子字符串，

func IsNotExist(err error)bool{
	//NOTE: not robust!
	return strings.Contains(err.Error(), "file does not exist")
}

但是处理I/O错误的逻辑可能一个和另一个平台非常的不同，
所以这种方案并不健壮并且对相同的失败可能会报出各种不同的错误消息。
在测试的过程中，通过检查错误消息的子字符串来保证特定的函数一起网的方式失败是非常有用的，
但对于线上的代码是不够的。

一个更可靠的方式是一个专门的类型来描述结构化的错误。
os包中定义了一个PathError类型来描述在文件路径操作中涉及到的失败，像Open或者Delete操作，
并且定义了一个叫LinkError的变体来描述涉及到两个文件路径的操作，
像Symlink和Rename。
这下面是os.PathError:
package os

// PathError records an error and the operation and file path that caused it.
type PathError struct{
	Op string
	Path string
	Err error
}

func (e *PathError)Error()string{
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

大多数调用方都不知道PathError并且通过调用错误本身的Error方法来统一处理所有的错误。
尽管PathError的Error方法简单地把这些字段链接起来生成错误消息，
PathError的结构保护了内部的错误组建。
调用方需要使用类型断言来检测错误的具体类型一遍将一种失败和另一种区分开；
具体的类型比字符串可以提供更多的细节。

_, err := os.Open("/no/such/file")
fmt.Println(err) // "open /no/such/file: No such file or directory"
fmt.Printf("%#v\n", err)
// Output:
// &os.PathError{Op:"open", Path:"/no/such/file", Err:0x2}

这就是三个帮助函数是怎么工作的。
例如下面展示的IsNotExist，
它会报出是否一个错误和syscall.ENOENT （7.8）或者和有名的错误os.ErrNotExist相等（可以在5.4.2中找到io.EOF）;
或者是一个*PathError，它内部的错误是syscall.ENOENT和os.ErrNotExist其中之一。

import (
	"errors"
	"syscall"
)

var ErrNotExist = errors.New("file does not exist")

// IsNotExist returns a boolean indicating whether the error is know to 
// report that a file or directory does not exist. It is satisfied by
// ErrNotExist as well as some syscall errors.
func IsNotExist(err error)bool{
	if pe, ok:=err.(*PathError);ok{
		err = pe.Err
	}
	return err == syscall.ENOENT || err == ErrNotExist
}

下面这里是他的实际使用：
_, err := os.Open("/no/such/file")
fmt.Println(os.IsNotExist(err)) // "true"

如果错误消息结合成一个更大的字符串，当然PathError的结构就不再为人所知，
例如通过一个对fmt.Errorf函数的调用。
区别错误通常必须在失败操作后，错误传回调用者前进行。

7.12 通过类型断言询问行为（p277）
下面这段逻辑和net/http包中web服务器负责写入HTTP头字段
（例如："Content-type:text/html"）的部分相似。
io.Writer接口类型的变量w代表HTTP响应；
写入他的字节最终被发送到某个人的web浏览器上。
func writeHeader(w io.Writer, contentType string)error{
	if _,  err := w.Write([]byte("Content-Type: ")); err != nil{
		//[]byte("string")是将内容转化成[]byte类型。
		return err
	}
	if _, err := w.Write([]byte(contentType)), err != nil{
		return err
	}
	// ...
}

因为Write方法需要传入一个byte切片，而我们希望写入的值是一个字符串，
所以我们需要使用[]byte(...)进行转换。
这个转换分配内存并且做一个拷贝，
但是这个拷贝在转换后几乎立马就被丢弃掉。
让我们假装这是一个web服务器的核心部分并且我们的性能分析表示这个内存分配使服务器的速度变慢。
这里我们可以避免掉内存分配么？

这个io.Writer接口告诉我们关于w持有的具体类型的唯一东西：
就是可以向它写入字节切片。
如果我们回顾net/http包中的内幕，我们知道在这个程序中的w变量持有的动态类型也有一个允许字符串高效写入的WriteString方法；
这个方法会避免去分配一个临时的拷贝。
（这可能像在黑夜中射击一样，但是许多满足io.Writer接口的重要类型同时也有WriteString方法，包括*bytes.Buffer，*bufio.Writer。）

我们不能对任意io.Writer类型的变量w，假设他也拥有WriteString方法。
但是我们可以定义一个只有这个方法的新街口并且使用类型断言来检测是否w的动态类型满足这个新接口。

// writeString writes s to w.
// If w has a WriteString method, it is invoked instead of w.Write.
func writeString(w io.Writer, s string)(n int, err error){
	type stringWriter interface{
		WriteString(string)(n int, err error)
	}
	if sw, ok:=w.(stringWriter);ok{
		return sw.WriteString(s) // avoid a copy
	}
	return w.Write([]byte(s)) // allocate temporary copy
}

func writeHeader(w io.Writer, contentType string) error{
	if _, err := writeString(w, "Content-Type: "); err != nil{
		return err
	}
	if _, err := writeString(w, contentType); err != nil{
		return err
	}
	// ...
}

为了避免重复定义，我们将这个检查移入到一个使用工具函数writeString中，
但是他太有用了，以至于标准库将他作为io.WriteString函数提供。
这是向一个io.Writer接口写入字符串的推荐方法。

这个例子的神奇支出在于没有定义了WriteString方法的标准接口和没有指定他是一个需要行为的标准接口。
而且一个具体类型只会通过他的方法决定它是否满足stringWriter接口，
而不是任何他和这个接口类型表明的关系。
他的意思就是上面的技术依赖于一个假设；
这个假设就是，如果一个类型满足下面的这个接口，
然后 WriteString(s)方法就必须和 Write([]byte(s))有相同的效果。

interface{
	io.Writer
	WriteString(s string)(n int, err error)
}

尽管io.WriteString记录了他的假设，但是调用他的函数极少有可能会去记录他们也做了同样的假设。
定义一个特定类型的方法隐式地获取了对特定行为的协约。
对于Go语言的新手，特别是那些来自有强类语言使用背景的新手，可能会发现他缺乏显示的意图令人感到混乱，
但是在实战的过程中这几乎不是一个问题。
除了空接口interface{}，接口类型很少意外巧合地被实现。

上面的WriteString函数使用一个类型断言来知道一个普遍接口类型的值是否满足一个更加具体的接口类型；
并且如果满足，他会使用这个更具体接口的行为。
这个技术可以被很好的使用不论这个被询问的接口是一个标准的如io.ReadWriter或者用户定义的如stringWriter。

这也是fmt.Fprintf函数怎么从其他所有值中区分满足error或者fmt.Stringer接口的值。
在fmt.Fprintf内部，有一个将单个操作对象转换成一个字符串的步骤，像下面这样：
package fmt

func formatOneValue(x interface{}) string{
	if err, ok := x.(error); ok{
		return err.Error()
	}
	if str, ok := x.(Stringer); ok{
		return str.String()
	}
	// ... all other types ...
}

如果x满足这两个接口类型中的一个，具体满足的接口决定对值的格式化方式。
如果都不满足，默认的case或多或少会同时地使用反射来处理所有的其他类型；
我们可以在第12章知道具体是怎么实现的。

再一次，他假设任何有String方法的类型满足fmt.Stringer中约定的行为，
这个行为会返回一个适合打印的字符串。

7.13 类型开关（p280）
接口被以两种不同的方式使用。
在第一个方式中，以io.Reader, io.Writer,fmt.Stringer,sort.Interface, http.Handler, 和error为典型，
一个接口的方法表达了实现这个接口的具体类型间的相似性，
但是隐藏了代表的细节和这些具体类型本身的操作。
重点在于方法上，而不是具体的类型上。

第二个方式利用一个接口值可以持有各种具体类型值的能力，并且将这个接口认为是这些类型的union（联合）。
类型断言用来动态地区别这些类型并且对每一种情况都 不一样。
在这个方式中，重点在于具体的类型满足这个接口，
而不是在于接口的方法（如果它确实有一些的话），并且没有任何的信息隐藏。
我们将以这种方式使用的接口描述为discriminated unions（可辨识联合）。

如果你熟悉面向对象编程，你可能会将这两种方式当作是subtype polymorphism（子类型多态）
和ad hoc polymorphism （非参数多态），但是你不需要去记住这些术语。
对于本章剩下的部分，我们将会呈现一些第二种方式的例子。

和其他那些语言一样，Go语言查询一个SQL数据库的API会干净地将查询中固定的部分和变化的分布分开。
一个调用的例子可能看起来像这样：

import "database/sql"

func listTracks(db sql.DB, artist string, minYear int) {
	result, err := db.Exec(
		"SELECT * FROM tracks WHERE artist = ? AND ? <= year AND year <= ?",
		artist, minYear, maxYear)
		// ...
	
}

Exec方法使用SQL字面量替换在查询字符串中的每个'?'；
SQL字面量表示相应参数的值，他有可能是一个布尔值，一个数字，一个字符串，或者nil空值。
用这种方式构造查询可以帮助避免SQL注入攻击；
这种攻击就是对手可以通过利用输入内容中不正确的引文来控制查询语句。
在Exec函数内部，我们可能会找到像下面这样的一个函数，他会将每一个参数值转换成他的SQL字面量符号。

func sqlQote(x interface{})string{
	if x == nil{
		return "NULL"
	}else if _, ok := x.(int); ok{
		return fmt.Sprintf("%d", x)
	}else if _, ok := x.(uint); ok{
		return fmt.Sprintf("%d", x)
	}else if b, ok := x.(bool); ok{
		if b {
			return "TRUE"
		}
		return "FALSE"
	}else if s, ok := x.(string); ok{
		return sqlQuoteString(s) // (not shown)
	}else{
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}

switch语句可以简化if-else链，如果这个if-else链对一连串值做相等测试、
一个相似的type switch（类型开关）可以简化类型断言的if-else链。

在它最简单的形式中，一个类型开关像普通的switch语句一样，
他的运算对象是x.(type)————他使用了关键词字面量type————并且每个case有一或多个类型。
一个类型开关基于这个接口值的动态类型使一个多路分支有效。
这个nil的case和if x == nil 匹配，并且这个default的case和如果其他case都不匹配的情况匹配。
一个对sqlQuote的类型开关可能会有这些case：
switch x.(type){
	case nil: // ...
	case int, uint: // ...
	case bool: // ...
	case string: // ...
	default: // ...
}

和（1.8）中的普通switch语句一样，每一个case会被顺序的进行考虑，
并且当一个匹配找到时，这个case中的内容会被执行。
当一个或多个case类型是接口时，case的顺序就会变得很重要，因为可能会有两个case同时匹配的情况。
default case相对其他case的位置是无所谓的。
他不会允许落空发生。

注意到在原来的函数中，对于bool和string情况的逻辑需要通过类型断言访问提取的值。
因为这个做法很典型，类型开关语句有一个扩展的形式，
他可以将提取的值绑定到一个在每个case范围内的新变量。

switch x := x.(type){ /*... */ }

这里我们已经将新的变量也命名为x；
和类型断言一样，重用变量名是很常见的。
和一个switch语句相似地，一个类型开关隐式的创建了一个语言块，
因此新变量x的定义不会和外面块中的x变量冲突。每一个case也会隐式的创建一个单独的语言块。

使用类型开关的扩展形式来重写sqlQuote函数会让这个函数更加的清晰：

func sqlQuote(x interface{}) string{
	switch x := x.(type){
	case nil:
		return "NULL"
	case int, uint:
		return fmt.Sprintf("%d", x) // x has type interface{} here.
	case bool:
		if x{
			return "TRUE"
		}
		return "FALSE"
	case string:
		return sqlQuoteString(x) // (not shown)
	default:
		panic(fmt.Sprintf("unexpected type %T: %v", x, x))
	}
}

在这个版本的函数中，在每个单一类型的case内部，变量x和这个case的类型相同。
例如，变量x在bool的case中是bool类型，和string的case中是string类型。
在所有其他的情况中，变量x是switch运算对象的类型（接口）；
在这个例子中运算对象是一个interface{}。
当多个case需要相同的操作时，比如int和uint的情况，类型开关可以很容易的合并这些情况。

尽管sqlQuote接受一个任意类型的参数，但是这个函数只会在他的参数匹配类型开关中的一个case时运行到结束；
其他情况的他会panic出“unexpected type”消息。
虽然x的类型是interface{}，但是我们把它认为是一个int，uint，bool，string,和nil值的discriminated union（可识别联合）。


7.14 示例： 基于标记的XML解码 (p283)
第4.5章节展示了如何使用encoding/json包中的Marshal和Unmarshal函数来将JSON文档转换成Go语言的数据结构。
encoding/xml包提供了一个相似的API。
当我们想构造一个文档树的表示时使用encoding/xml包会很方便，
但是对于很多程序并不是必须的。
encoding/xml包也提供了一个更底层的基于标记的API用于XML解码。
在基于标记的样式中，解析器消费输入和产生一个标记流；
四个主要的标记类型——StartElement， EndElement，CharData，和Comment————
每一个都是encoding/xml包中的具体类型。
每一个对 (*xml.Decoder).Token的调用都返回一个标记。
这里显示的是和这个API相关的部分：

encoding/xml

package xml

type Name struct{
	Local string // e.g., "Title" or "id"
}

type Attr struct { // e.g., name= "value"
	Name Name
	Value string
}

// A Token includes StartElement, EndElement, CharData,
// and Comment, plus a few esoteric types (not shown).
type Token interface{}
type StartElement struct{ // e.g., <name>
	Name Name
	Attr []Attr
}
type EndElement struct {Name Name} // e.g., </name>
type CharData []byte // e.g., <p>CharData</p>
type Comment []byte // e.g., <!-- Comment -->

type Decoder strunt{ /*...*/}
func NewDecoder(io.Reader)*Decoder
func (*Decoder) Token()(Token, error) // returns next Token in sequence

这个没有方法的Token接口也是一个可识别联合的例子。
传统的接口如io.Reader的目的是隐藏满足他的具体类型的细节，
这样就可以创造出新的实现；
在这个实现中每个具体类型都被统一地对待。
相反，满足可识别联合的具体类型的集合被设计确定和暴露，
而不是隐藏。
可识别的联合类型几乎没有方法；
操作他们的函数使用一个类型开关的case集合来进行表述；
这个case集合中每一个case这个有不同的逻辑。
下面的xmlselect程序获取和打印在一个XML文档树中确定的元素下找到的文本。
使用上面的API，它可以在输入上一次完成他的工作而从来不要具体化这个文档树。

示例代码

每次main函数中的循环遇到一个StartElement时，它把这个元素的名称压到一个栈里；
并且每次遇到EndElement时，它将名称从这个栈中推出。
这个API保证了StartElement和EndElement的序列可以被完全的匹配，甚至咋一个糟糕的文档格式中。
注释会被忽略。
当xmlselect遇到一个CharData时，
只有当栈中有序地包含不所有通过命令行参数传入的元素名称时他才会输出相应的文本。

下面的命令打印出任意出现在两层div元素下的h2元素的文本。
他的输入是XML的说明文档，并且它自己就是XML文档格式的。

示例结果

7.15 一些建议 （p288）
当设计一个新的包时，新的Go程序员总是通过创建一个接口的集合开始和后面定义满足他们的具体类型。
这种方式的结果就是有很多的接口，他们中的每一个仅只有一个实现。
不要在这么做了。
这种接口是不必要的抽象；
他们也有一个运行时损耗。
你可以使用导出机制（6.6）来限制一个类型的方法或一个结构体的字段是否在包外可见。
接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要。

当一个接口只被一个单一的具体类型实现时有一个例外，就是由于它的依赖，
这个具体类型不能和这个接口存在在一个相同的包中。
这种情况下，一个接口是解耦这两个包的一个好方式。

因为在Go语言中只有当两个或更多的类型实现一个接口时才使用接口，
他们必定会从任意特定的实现细节中抽象出来。
结果就是有更少和更简单方法（经常和io.Writer或fmt.Stringer一样只有一个）的更小的接口。
当新的类型出现时，小的接口更容易满足。
对于接口设计的一个好的标准就是ask only for what you need (只考虑你需要的东西)

我们完成了对methods和接口的学习过程。
Go语言良好的支持面向对象风格的编程，但不是说你仅仅只能使用他。
不是任何事物都需要被当做成一个对象；
独立的函数有他们自己的用处，未封装的数据类型也是这样。
同时观察到这两个，在本书的前五章的例子中没有调用超过两打方法
像input.Scan，与之相反的是普遍的函数调用如fmt.Printf。
