第八章 Goroutines和Channels
（p289-338）

1. 并发
并发程序指同时进行多个任务的程序，随着硬件的发展，
并发程序变得越来越重要。
Web服务器会一次处理成千上万的请求。
平板电脑和手机app在渲染用户画面同时还会后台执行各种计算任务和网络请求。
即使是传统的批量处理问题——读取数据，计算，写输出——现在也会用并发来隐藏掉I/O的操作延迟以充分利用现代计算机设备的多个核心。
计算机的性能每年都在以非线性的速度增长。

Go语言中的并发程序可以用两种手段来实现。
本章讲解goroutine和channel，
其支持“顺序通信进程”（communicating sequential processes）或简称为CSP。
CSP是一种现代的并发编程模型，在这种编程模型中，值会在不同的运行实例（goroutine）中传递，
尽管大多数情况下任然是被限制在单一实例（gorountine）中。
第9章覆盖更为传统的并发模型：多线程共享内存，如果你在其他的主流语言中写过并发程序的话
可能会更熟悉一些。
第9章也会深入介绍一些并发程序带来的风险和陷阱。

尽管Go对并发的支持是众多强力特性之一，但跟踪调试并发程序还是很困难，
在线性程序中形成的直觉往往还会使我们误入歧途。
如果这是读者第一次接触并发，推荐稍微多花一些时间来思考这两个章节中的样例。

8.1 Goroutines（p290）
2. 在Go语言中，每一个并发的执行单元叫作一个goroutine。
设想这里的一个程序有两个函数，一个函数做计算，另一个输出结果，
假设两个函数没有相互之间的调用关系。
一个线性的程序会先调用其中的一个函数，然后再调用另一个。
如果程序中包含多个goroutine，对两个函数的调用则可能发生在同一时刻。
马上就会看到这样的一个程序。

如果你使用过操作系统或者其他语言提供的线程，那么你可以简单地把goroutine类比作一个线程，
这样你就可以写出一些正确的程序了。
goroutine和线程的本质区别会在9.8节中讲。

当一个程序启动时，其主函数即在一个单独的goroutine中运行，
我们叫它main goroutine。
新的goroutine会用go语句来创建。
在语法上，go语句是一个普通的函数或方法调用前加上关键字go。
go语句会使其语句中的函数在一个新创建的goroutine中运行。
而go语句本身会迅速地完成。

f() // call f(); wait for it to return
go f() // create a new goroutine that calls f(); don't wait

下面的例子，main goroutine将计算斐波那契数列的第45个元素值。
由于计算函数使用低效的递归，所以会运行相当长时间，
在此期间我们想让用户看到一个可见的标识来表明程序依然在正常运行，
所以来做一个动画的小图标：

示例代码

动画显示了几秒之后，fib（45）的调用成功地返回，并且打印结果：
Fibonacci(45) = 1134903170

然后主函数返回。
主函数返回时，所有的goroutine都会被直接打断，程序退出。
除了从主函数退出或者直接终止程序之外，没有其他的编程方法能够让一个goroutine来打断另一个的执行，
但是之后可以看到一种方式来实现这个目的，
通过goroutine之间的通信来让一个goroutine请求其他的goroutine，
并被请求的goroutine自行结束执行。

留意一下这里的两个独立的单元是如何进行组合的，
spinning和斐波那契的计算。
分别在独立的函数中，但两个函数会同时执行。

8.2 示例：并发的Clock服务（p292）
网络编程是并发大显身手的一个领域，由于服务器是最典型的需要同时处理很多链接的程序，
这些链接一般来自彼此独立的客户端。
在本小节中，我们会讲解go语言的net包，
这个包提供编写一个网络客户端或者服务器程序的基本组件，
无论两者间通信是使用TCP，UDP或者Unix domain sockets。
在第一章中我们已经使用过的net/http里的方法，也算是net包的一部分。
我们的第一个例子是一个顺序执行的时钟服务器，他会每隔一秒钟将当前时间写到客户端：

示例代码

Listen函数创建了一个net.Listener的对象，
这个对象会监听一个网络端口上到来的连接，
在这个例子里我们用的是TCP的localhost:8000端口。
listener对象的Accept方法会直接阻塞，
直到一个新的连接被创建，然后会返回一个net.Conn对象表示这个连接。

handleConn函数会处理一个完整的客户端连接。
在一个for死循环中，将当前的时候用time.Now()函数得到，
然后写到客户端。
由于net.Conn实现了io.Writer接口，我们可以直接向其写入内容。
这个死循环会一直执行，直到写入失败。
最可能的原因是客户端主动断开连接。
这种情况下handleConn函数会用defer调用关闭服务器侧的连接，然后返回到主函数，
继续等待下一个连接请求。

time.Time.Format方法提供了一种格式化日期和时间信息的方式。
他的参数是一个格式化模板标识如何来格式化时间，而这个格式化模板限定为Mon Jan 2 03:04:05PM 2006 UTC-0700。
有八个部分（周几，月份，日期，几时，几分，几秒，年 时区）
可以以任意的形式来咋和前面这个模板；
出现在模板中的部分会作为参考来对时间格式进行输出。
在上面的例子中我们只用到了小时、分钟和秒。
time包里定义了很多标准时间格式，比如time.RFC1123。
在进行格式化的逆向操作time.Parse时，也会用到同样的策略。
（这是Go语言和其他语言相比比较奇葩的一个地方。你需要技术格式化字符串是1月2日下午3点4分5秒零6年UTC-0700，
而不像其他语言那样Y-m-d H:i:s一样，当然了这里可以用1234567的方式来记忆，倒也不麻烦）

3.涉及一个netcat工具（nc命令）
为了连接例子里的服务器，我们需要一个客户端程序，比如netcat这个工具（nc命令），这个工具可以用来执行网络连接操作。

示例结果（没有nc工具，我没有演示这个结果，后续补上）

客户端将服务器发来的时间显示了出来，我们用Control+C来终端客户端的执行，
在Unix系统上，你会看到^C这样的响应。
如果你的系统没有装nc这个工具，你可以用telnet来实现同样的效果，或者也可以用我们下面的这个go写的简单的telnet程序，
用net.Dial就可以简单地创建一个TCP连接：

示例代码/ch8/netcat1

这个程序会从连接中读取数据，并将督导的内容刚写到标准输出中，
直到遇到end of file的条件或者发生错误。
mustCopy这个函数我们在本节的几个例子中都会用到。
让我们同时运行两个客户端来进行一个测试，
这里可以开两个终端窗口，
下面左边的是其中的一个的输出，右边的是另一个的输出：

$ go build gopl.io/ch8/netcat1
$ ./netcat1
13:58:54 		$ ./netcat1
13:58:55
13:58:56
^C
				13:58:57
				13:58:58
				13:58:59
				^C
$ killall clock1

killall命令是Unix命令行工具，可以用给定的进程名来杀掉所有名字匹配的进程。

第二个客户端必须等待第一个客户端完成工作，
这样服务端才能继续向后执行；
因为我们这里的服务器程序同一时间只能处理一个客户端链接。
我们这里对服务端程序做一点小改动，时期支持并发：
在handleConn函数调用的地方增加go关键字，
让每一次handleConn的调用都进入一个独立的goroutine。

示例代码ch8/clock2

现在多个客户端可以同时同时接收到时间了：

$ go build gopl.io/ch8/clock2
$ ./clock2 &
$ go build gopl.io/ch8/netcat1
$ ./netcat1
14:02:54 		$ ./netcat1
14:02:55 		14:02:55
14:02:56 		14:02:56
14:02:57 		^C
14:02:58
14:02:59 		$ ./netcat1
14:03:00 		14:03:00
14:03:01 		14:03:01
^C 				14:03:02
				^C
$ killall clock2

8.3 示例：并发的Echo服务（p298）
clock服务器每一个链接都会起一个goroutine。
在本节中，我们会创建一个echo服务器，
这个服务在每个链接中会有多个goroutine。
大多数echo服务仅仅会返回他们读取到的内容，
就像下面这个简单的handleConn函数所做的一样：

func handleConn(c net.Conn){
	io.Copy(c, c) //NOTE: ignoring errors
	c.Close()
}

一个更有意思的echo服务应该模拟一个实际的echo的“回响”，
并且一开始要用大写HELLO来表示“声音很大”，之后经过一小段延迟返回一个有所缓和的Hello，
然后一个全小写字母的hello表示声音渐渐变小直至消失，
像下面这个版本的handleConn：

示例代码 ch8/reverb1

我们需要升级我们的客户端程序，这样它就可以发送终端的输入到服务器，
并把服务端的返回输出到终端上，这使我们有了使用并发的另一个好机会：

示例代码 ch8/netcat2

当main goroutine 从标准输入流中读取内部并将其发送给服务器时，
另一个goroutine会读取并打印服务端的响应。
当main goroutine碰到输入终止时，例如，用户在终端中按了Control-D(^D)，
在windows上是Control-Z，这时程序就会被终止，
尽管其他goroutine中还有进行中的任务。
（在8.4.1中引入了channels后我们会明白如何让程序等待两边都结束）。

下面这个会话中，客户端的输入是左对齐的，服务端的响应会用缩进来区别显示。
客户端会向服务器“喊话三次”：

$ go build gopl.io/ch8/reverb1
$ ./reverb1 &
$ go build gopl.io/ch8/netcat2
$ ./netcat2
Hello?
	HELLO?
	Hello?
	hello?
Is there anybody there?
	IS THERE ANYBODY THERE?
Yooo-hooo!
	Is there anybody there?
	is there anybody there?
	YOOO-HOOO!
	Yooo-hooo!
	yooo-hooo!
^D
$ killall reverb1


注意客户端的第三次shout在前一个shout处理完之前一直没有被处理，
这貌似看起来不是特别“现实”。
真实世界里的回响应该是由三次shout的回声组合而成的。
为了模拟真实世界的回响，我们需要更多的goroutine来做这件事情。
这样我们就再一次地需要go这个关键词了，
这次我们用它来调用echo：

示例代码/ch8/reverb2

go后跟的函数的参数会在go语句自身执行时被求值；
因此input.Text()会在main goroutine中被求值。
现在回响是并发并且会按时间来覆盖掉其他响应了：

$ go build gopl.io/ch8/reverb2
$ ./reverb2 &
$ ./netcat2
Is there anybody there?
	IS THERE ANYBODY THERE?
Yooo-hooo!
	Is there anybody there?
	YOOO-HOOO!
	is there anybody there?
	Yooo-hooo!
	yooo-hooo!
^D
$ killall reverb2

//注意观察和上面的示例结果有不一样的地方。

让服务使用并发不只是处理多个客户端的请求，
甚至在处理单个连接时也可能会用到，
就像我们上面的两个go关键词的用法。
然而在我们使用go关键词的同时，需要慎重地考虑net.Conn中的方法在并发地调用时是否安全，
事实上对于大多数类型来说也确实不安全。
我们会在下一章中详细地探讨并发安全性。

8.4 Channels （p301）
如果说goroutine是Go语言程序的并发体的话，那么channels是他们之间的通信机制。
一个channels是一个通信机制，它可以让一个goroutine通过他给另一个goroutine发送值信息。
每个channel都有一个特殊的类型，也就是channels可发送数据的类型。
一个可以发送int类型数据的channel一般写为chan int。

使用内置的make函数，我们可以创建一个channel：

ch := make(chan int) // ch has type 'chan int'

和map类型， channel也一个对应make创建的底层数据结构的引用哦。
当我们复制一个channel或用于函数参数传递时，
我们只是拷贝了一个channel引用，因此调用者可被调用者将引用同一个channel对象。
和其他的引用类型一样，channel的零值也是nil。

两个相同类型的channel可以使用 == 运算符比较。
如果两个channel引用的是相同的对象，那么比较的结果为真。
一个channel也可以和nil进行比较。

一个channel有发送和接受两个主要操作，都是通信行为。
一个发送语句将一个值从一个goroutine通过channel发送到另一个执行接受操作的goroutine。
发送和接收两个操作都是用 <- 运算符。
在发送语句中，<-运算符分割channel和要发送的值。
在接收语句中，<-运算符写在channel对象之前。
一个不使用接收结果的接收操作也是合法的。

ch <- x // a send statement 发送语句
x = <-ch // a receive expression in an assignment statement赋值语句中的接收表达式
<-ch // a receive statement; result is discarded.结果被丢弃。

Channel还支持close操作，用于关闭channel，
随后对基于该channel的任何发送操作都将导致panic异常。
对一个而已经被close过的channel执行接收操作依然可以接收到之前已经成功发送的数据；
如果channel中已经没有数据的话将产生一个零值的数据。

使用内置的close函数就可以关闭一个channel：
close(ch)

以最简单方式调用make函数创建的是一个无缓存的channel，
但是我们也可以制定第二个整形参数，
对应channel的容量。
如果channel的容量大于零，那么该channel就是带缓存的channel。

ch = make(chan int) // unbuffered channel
ch = make(chan int, 0) // unbuffered channel
ch = make(chan int, 3) // buffered channel with capacity 3

我们将先讨论无缓存的channel，然后在8.4.4节讨论带缓存的channel。

8.4.1不带缓存的channels（p302）
一个基于无缓存Channels的发送操作将导致发送者goroutine阻塞，
直到另一个goroutine在相同的Channels上执行接收操作，当发送的值通过Channels成功传输之后，
两个goroutine可以继续执行后面的语句。
反之，如果接收操作先发生，那么接收者goroutine也将阻塞，
直到有另一个goroutine在相同的Channels上执行发送操作。

基于无缓存Channels的发送和接收操作将导致两个goroutine做一次同步操作。
因为这个原因，无缓存channels有时候也被成为同步Channels。
当通过一个无缓冲Channels发送数据时，接收者收到数据发生在唤醒发送者goroutine之前。
（译注： happens before， 这是Go语言并发内存模型的一个关键术语）

在讨论并发编程时，当我们说x事件在y事件之前发生（happens before），
我们并不是说x事件在时间上比y时间更早；
我们要表达的意思是要保证在此之前的事件都已经完成了，
例如在此之前的更新某些变量的操作应完成，你可以放心依赖这些已经完成的事件了。

当我们说x事件既不是在y事件之前发生也不是在y事件之后发生，我们就说x事件和y事件是并发的。
这并不是意味着x事件和y事件就一定是同时发生的，我们只是不能确定这两个事件发生的先后顺序。
在下一章中我们将看到，当两个goroutine并发访问了相同的变量时，
我们有必要保证某些事件的执行顺序，以避免出现某些并发问题。

在8.3节的客户端程序，它在主goroutine中（就是执行main函数的goroutine）将标注输入复制到server，
因此当客户端程序关闭标准输入时，后台goroutine可能依然在工作。
我们需要让主goroutine等待后台goroutine完成工作后在退出，
我们使用了一个channel来同步两个goroutine：

示例代码ch8/netcat3

当用户关闭了标准输入，主goroutine中的mustCopy函数调用将返回，
然后调用conn.Close()关闭读和写方向的网络连接。
关闭网络链接中的写方向的链接将导致server程序收到一个文件（end-of-file）结束的信号。
关闭网络链接中读的方向的链接将导致后台goroutine的io.Copy函数调用返回一个“read from closed connection”（从关闭的链接读）类似的错误，
因此我们临时溢出了错误日志语句；
在练习8.3将会提供一个更好的解决方案。
（需要注意的是go语言调用了一个函数字面量，这go语言中启动goroutine常用的形式。）

在后台goroutine返回之前，它先打印一个日志信息，
然后向done对应的channel发送一个值。
主goroutine在退出前先等待从done对应的channel接收一个值。
因此，总是可以在程序退出前正确输出“done”消息。

基于channels发送消息有两个重要方面。
首先每个消息都有一个值，但是有时候通讯的事实和发生的时刻也同样重要。
当我们更希望强调通讯发生的时刻，我们将它称为消息事件。
有些消息事件并不携带额外信息，它仅仅是用作两个goroutine之间的同步，这时候我们可以用struct{}
空结构体作为channels元素的类型，
虽然也可以使用bool或int类型实现同样的功能，done <- 1 语句也比 done <- struct{} 更短。

8.4.2 串联的Channels（Pipeline）

Channels也可以用于将多个goroutine链接在一起，
一个Channels的输出作为下一个Channels的输入。
这种串联的Channels就是所谓的管道（pipeline）。
下面的程序用两个channels将三个goroutine串联起来，

第一个goroutine是一个计数器，用于生成0、1、2....形式的整数序列，
然后通过channel将该整数序列发送给第二个goroutine；
第二个goroutine是一个求平方的程序，
对收到的每个整数求平方，
然后将平方后的结果通过第二个channel发送给第三个goroutine；
第三个goroutine是一个打印程序，打印收到的每个整数。
为了报出例子清晰，我们有意选择了非常简单的函数，
当然三个goroutine的计算很简单，在现实中确实没有必要为如此简单的运算构建三个goroutine。

示例代码 ch8/pipeline1

如您所料，上面的程序将生成 0、1、4、9...形式的无穷数列。
像这样的串联Channels的管道（Pipelines）可以用在需要长时间运行的服务中，
每个长时间运行的goroutine可能会包含一个死循环，
在不同goroutine的死循环内部使用串联的Channels来通信。
但是，如果我们希望通过Channels只发送有限的数列该如何处理呢？

如果发送者知道，没有更多的值需要发送到channel的话，
那么让接收者也能即使知道没有多余的值可接收将是有用的，
因为接收者可以停止不必要的接收等待。
这可以通过内置的close函数来关闭channel实现：

close(naturals)

当一个channel被关闭后，再向该channel发送数据将导致panic异常。
当一个被关闭的channel中已经发送的数据都被成功接收后，
后续的接收操作将不再阻塞，他们会立即返回一个零值。
关闭上面例子中的natruals变量对应的channel并不能终止循环，
它依然会收到一个永无休止的零值序列，然后将它们发送给打印者goroutine。

没有办法直接测试一个channel是否被关闭，但是接收操作有一个变体形式：
它多接收一个结果，多接收的第二个结果是一个布尔值ok,
true 表示成功从channels接收到值，
false channel已经并关闭，并且里面没有值可接收。
使用这个特性，我们可以修改squarer函数中循环代码，当naturals对应的channel被关闭并没有值可接收时跳出循环，
并且也关闭squares对应的channel。

// Squarer
go func(){
	for{
		x, ok := <-naturals
		if !ok{
			break // channel was closed and drained
		}
		squares <- x*x
	}
	close(squares)
}()

因为上面的语法是笨拙的，而且这种处理模式很场景，
因此Go语言的range循环可直接在channel上面迭代。
使用range循环是上面处理模式的简洁语法，它依次从channel接收数据，
当channel被关闭并且没有值可接收时跳出循环。

在下面的改进中，我们的计算器goroutine只生成100个含数字的序列，
然后关闭naturals对应的channel，
这将导致计算平方数的squarer对应的goroutine可以正常终止循环并关闭squares对应的channel。
（在一个更复杂的程序中，可以通过defer语句关闭对应的channel）
最后，主goroutine也可以正常终止循环并退出程序。

代码示例 ch8/pipeline2

其实你并不需要关闭每一个channel。
只要当需要告诉接收者goroutine，
所有的数据已经全部发送时才需要关闭channel。
不管一个channel是否被关闭，
当它没有被引用时将会被Go语言的垃圾自动回收器回收。
（不要将关闭一个打开文件的操作和关闭一个channel操作混淆。
对于每个打开的文件，都需要在不使用的时候调用对应的Close方法来关闭文件。）

试图重复关闭一个channel将导致panic异常，
试图关闭一个nil值的channel也将导致panic异常。
关闭一个channels还会触发一个广播机制，我们将在8.9节讨论。

8.4.3 单方向的channel（p306）
随着程序的增长，人们习惯于将大的函数拆分为小的函数。
我们前面的例子中使用了三个goroutine，然后用两个channel链接他们，
他们都是main函数的局部变量。
将三个goroutine拆分为以下三个函数是自然的想法：

func counter(out chan int)
func squarer(out, in chan int)
func printer(in chan int)

其中squarer计算平方的函数在两个串联Channels的中间，
因此拥有两个channel类型的参数，
一个用于输入一个用于输出。
每个channel都用有相同的类型，但是他们的使用方式相反：
一个只用于接收，另一个只用于发送。
参数的名字in和out已经明确表示了这个意图，
但是并无法保证squarer函数向一个in参数对应的channel发送数据，
或者从一个out参数对应的channel接收数据。

这种场景是典型的。
当一个channel作为一个函数参数时，他一般总是被专门用于只发送或者值接收。

为了表明这种意图并防止被滥用，
Go语言的类型系统提供了单方向的channel类型，
分别用于只发送或只接收的channel。
类型 chan<- int 表示一个只发送int的channel，
只能发送不能接收。
相反，类型 <-chan int 表示一个只接收int的channel，
只能接收不能发送。
（箭头 <- 和关键字chan的相对位置表明了channel的方向）
这种限制将在编译期检测。

因为关闭操作只用于断言不再向channel发送新的数据，
所以只有在发送者所在的goroutine才会调用close函数，
因此对一个只接收的channel调用close将是一个编译错误。

这是改进的版本，这一次参数使用了单方向channel类型：

示例代码 ch8/pipeline3

调用 counter(naturals) 将导致chan int类型的naturals
隐式地转换为 chan<- int类型只发送型的channel。
调用 printer(squares)也会导致相似的隐式转换，
这一次是转换为 <-chan int类型，只接收型的channel。
任何双向channel向单向channel变量的赋值操作都将导致该隐式转换。
这里并没有反向转换的语法：也就是不能将一个类似 chan<- int 类型的单向型的channel转换为
chan int 类型的双向型channel。

8.4.4 带缓存的channel（p308）
带缓存的Channel内部持有一个元素队列。
队列的最大容量是在调用make函数创建channel时通过第二个参数指定的。
下面的语句创建了一个可以持有三个字符串元素的带缓存channel。
ch = make(chan string, 3) 

向缓存channel的发送操作就是向内部缓存队列的尾部插入元素，
接收操作则是从队列的头部删除元素。
如果内部缓存队列是满的，那么发送操作将阻塞直到因另一个goroutine执行接收操作而释放了新的队列空间。
相反，如果channel是空的，接收操作将阻塞知道有另一个goroutine执行发送操作而向队列插入元素。

我们可以在无阻塞的情况下连续向新创建的channel发送三个值：
ch <- "A"
ch <- "B"
ch <- "C"

此刻，channel的内部缓存队列将是满的，如果有第四个发送操作将发生阻塞。
ch -------> "A" "B" "C"

如果我们接收一个值，
fmt.Println(<-ch)  // "A"

那么channel的缓存队列将不是满的，也不是空的，
因此对该channel执行的发送或接收操作都不会发生阻塞。
通过这种方式，channel的缓存队列解耦了接收和发送的goroutine。

ch ------> "B" "C" " "

在某些特殊情况下，程序可能需要知道channel内部缓存的容量，
可以用内置的cap函数获取:
fmt.Println(cap(ch)) // "3"

同样，对于内置的len函数，如果传入的是channel，那么将返回channel内部缓存队列中有效元素的个数。
因为在并发程序中该信息会随着接收操作而失效，但是它对某些故障诊断和性能优化会有帮助。
fmt.Println(len(ch)) // "2"

在继续执行两次接收操作后channel内部的缓存队列将又变成空的，如果有第四个接收操作将发生阻塞：
fmt.Println(<-ch) // "B"
fmt.Println(<-ch) // "C"

在这个例子中，发送和接收操作都发生在同一个goroutine中，
但是在真实的程序中它们一般由不同的goroutine执行。
Go语言新手有时候会将一个带缓存的channel当作用一个goroutine中的队列使用，
虽然语法看似简单，但实际上这是一个错误。
Channel和goroutine的调度器机制是紧密相连的，
一个发送操作————或许是整个程序————可能会永远阻塞。
如果你只是需要一个简单的队列，使用slice就可以了。

下面的例子展示了一个使用了带缓存channel的应用。
他并发地向三个镜像站点发出请求，
三个镜像站点分散在不同的地理位置。
他们分别将收到的响应发送到带缓存channel，
最后接收者只接收第一个收到的响应，也就是最快的那个响应。
因此mirroredQuery函数可能在另外两个响应慢的镜像站点响应之前就返回了结果。
（顺便说下，多个goroutines并发地向同一个channel发送数据，或从同一个channel接收数据都是常见的用法）
func mirroredQuery() string{
	responses := make(chan string, 3)
	go func (){responses <- request("asia.gopl.io")}()
	go func (){responses <- request("europe.gopl.io")}()
	go func (){responses <- request("americas.gopl.io")}()
	return <-responses // return the quickest response
}

func request(hostname string)(response string){/*...*/}

如果我们使用了无缓存的channel，那么两个慢的goroutines将会因为没有人接收而被永远卡住。
这种情况，成为goroutine泄露，浙江是也给bug。
和垃圾变量不同，泄露的goroutines并不会被自动回收，因此确保每个不再需要的goroutine能正常退出是重要的。

关于无缓存或者带缓存channel之间的选择，或者是带缓存channel的容量大小的选择，
都可能影响程序的正确性。
无缓存channel更强地保证了每个发送操作与响应的同步接收操作；
但是对于带缓存channel，这些操作是解耦的。
同样，即使我们知道将要发送到一个channel的信息的数量上限，
创建一个对应容量大小带缓存channel也是不现实的，
因为这要求在执行任何接收操作之前缓存已经发送完所有的值。
如果未能分配足够的缓存将导致程序死锁。

channel的缓存也可能影响程序的性能。
想象一家蛋糕店有三个厨师，一个烘焙，一个上糖衣，还有一个将每个蛋糕递到下一个厨师在的生产线。
在狭小的厨房空间环境，每个厨师在完成蛋糕后必须等待下一个厨师已经准备好接受它；
这类似于在一个无缓存的channel上进行沟通。

如果在每个厨师之间有一个放置一个蛋糕的额外空间，那么每个厨师就可以将一个完成的蛋糕临时放在那里而马上进入下一个蛋糕的制作；
这类似于将channel的缓存队列的容量设置为1。
只要每个厨师的平均工作效率相近，那么其中大部分的传输工作将是迅速的，
个体之间细小的效率差异将在交接过程中弥补。
如果厨师之间有更大的额外空间————也就是更大容量的缓存队列————将可以在不停止生产线的前提下消除更大的效率波动，
例如一个厨师可以短暂地休息，然后再加快赶上进度而不影响其他人。

另一方面，如果生成线的前期阶段一直快于后续阶段，那么他们之间的缓存在大部分时间都将是满的。
相反，如果后续阶段币前期阶段更快，那么他们之间的缓存在大部分时间都将是空的。
对于这类场景，额外的缓存并没有带来任何好处。

生产线的隐喻对于理解channel和goroutine的工作机制是很有帮助的。
例如，如果第二阶段是需要精心制作的复杂操作，一个厨师可能无法跟上第一个厨师的进度，或者无法满足第三阶段厨师的需求。
要解决这个问题，我们可以雇佣另一个厨师来帮助完成第二阶段的工作，
他执行相同的任务但是独立工作。
这类似基于相同的channel创建另一个独立的goroutine。

我们没有太多的空间展示全部细节，但是ch8/cake包模拟了这个蛋糕店，
可以通过不同的参数调整。
他还对上面提到的集中场景提供对应的基准测试。（11.4）



8.5 并发的循环（p312）

本节中，我们会探索一些用来在并行时循环迭代的常见并发模型。
我们会探究从全尺寸图片生成一些缩略图的问题。
ch8/thumbnail包提供了ImageFile函数来帮助我们拉伸图片。
我们不会说明这个函数的实现，只需要从gopl.io下载他。

package thumbnail

func ImageFile(infile string)(string, error)

下面的程序会循环迭代一些图片文件名，并为每一张图片生成一个缩略图：
// makeThumbnails makes thumbnauls of the specified files.
func makeThumbnails(filenames []string){
	for _, f := range filenames{
		if _, err := thumbnail.ImageFile(f); err != nil{
			log.Println(err)
		}
	}
}

显然我们处理文件的顺序无关紧要，因为每一个图片的拉伸操作和其他图片的处理操作都是彼此独立的。
像这种子问题都是完全彼此独立的问题被叫做易并行问题（embarrassiongly parallel，直译的话更像是尴尬并行）。
易并行问题是最容易被实现成并行的一类问题（废话），
并且是最能够享受并发带来的好处，能够随着并行的规模线性地扩展。

下面让我们并行地执行这些操作，从而将文件IO的延迟隐藏掉，
并用上多核cpu的计算能力来拉伸图像。
我们的第一个并发程序只是使用了go关键字。
这里我们先忽略掉错误，之后在进行处理。

//NOTE: incorrect!
func makeThumbnail2(filenames []string){
	for _, f := range filenemes{
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
	}
}

这个版本运行的是在有点太快了，实际上，
由于它比最早的版本使用的时间要短得多，
即使当文件名的slice中只包含有一个元素。
这就有点奇怪了，如果程序没有并发执行的话，那为什么一个并发的版本还是要快呢？
答案其实是makeThumbnail在它还没有完成工作之前就已经返回了。
他启动了所有的goroutine，每一个文件名对应一个，但没有等待他们一直到执行完毕。

没有什么直接的办法能够等待goroutine完成，
但是我们可以改变goroutine里的代码让其能够将完成情况报告给外部的goroutine知晓，
使用的方式是向一个共享的channel中发送事件。
因为我们已经知道内部的goroutine只有len（filenames），
所以外部的goroutine只需要在返回之前对这些事件计数。

// makeThumbnail3 makes thumbnails of the specified files in paralled.
func makeThumbnail3(filenames []string){
	ch := make(chan struct{})
	for _, f := range filenames{
		go func(f string){
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f) // 这个(f)是什么意思？
	}
	//Wait for goroutines to complete.
	for range filenames{
		<-ch
	}
}

注意我们将f的值作为一个显式的变量传给了函数，
而不是在循环的闭包中声明：
for _, f := range filenames{
	go func(){
		thumbnail.ImageFile(f) // NOTE: incorrect!
	}()
}

回忆一下之前在5.6.1中，匿名函数中的循环变量快照问题。
上面这个单独的变量f是被所有的匿名函数值所共享，
且会被连续的循环迭代所更新的。
当新的goroutine开始执行字面函数时，for循环可能已经更新了f并且开始了另一轮的迭代或者（更有可能的）已经结束了整个循环，
所以当这些goroutine开始读取f的值时，他们所看到的值已经是slice的最后一个元素了。
显示地添加这个参数，我们能够确保使用的f是当go语句执行时的“当前”那个f。

如果我们想要从每一个worker goroutine往主goroutine中返回值时该怎么办？
当我们调用thumbnail.ImageFile创建文件失败的时候，他会返回一个错误。
下一个版本的makeThumbnails会返回其在做拉伸操作时接收到的第一个错误：

// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error of any step failed.
func makeThumbnail4(filenames []string)error{
	errors := make(chan error)

	for _, f := range filenames{
		go func(f string){
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames{
		if err := <-errors; err != nil{
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil
}

这个程序有一个微妙的bug。
当它遇到第一个非nil的error时会直接将error返回到调用方，
使得没有一个goroutine去排空errors channel。
这样剩下的worker goroutine在向这个channel中发送值时，
都会永远地阻塞下去，并且永远都不会退出。
这种情况叫做goroutine泄露（8.4.4），可能会导致整个程序卡住或者跑出out of memory（内存不足）的错误。

最简单的解决办法就是用一个具有合适大小的buffered channel，
这样这些worker goroutine向channel中发送测向（什么意思？？？）时就不会被阻塞。
（一个可选的解决办法是创建一个另外的goroutine，当main goroutine返回第一个错误的同时去排空channel）

下一个版本的的makeThumbnails使用了一个buffered channel来返回生成的图片文件的名字，
附带生成时的错误。

// makeThumbnails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order, 
// or an error if any step failed.

func makeThumbnails5(filenames []string)(thumbfiles []string, err error){
	type item struct{
		thumbfile string
		err error
	}

	ch := make(chan item, len(filenames))
	for _, f := range filenames{
		go func (f string){
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames{
		it := <-ch
		if it.err != nil{
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

我们最后一个版本的makeThumbnails返回了新文件的大小总计数（bytes）。
和前面的版本都不一样的一点是我们在这个版本里没有把文件名放在slice里，
而是通过一个string的channel传过来，所以我们无法对循环的次数进行预测。

为了知道最后一个goroutine什么时候结束（最后一个结束并不一定是最后一个开始），
我们需要一个递增的计数器，在每一个goroutine启动时加一，在goroutine退出时减一。
这需要一种特殊的计数器，这个计数器需要在多个goroutine操作时做到安全并且提供在其减为零之前一直等待的一种方法。
这种技术类型被成为sync.WaitGroup，下面的代码就用到了这种方法：

// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.

func makeThumbnails6(filenames <-chan string) int64{
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for f := range filenames{
		wg.Add(1)
		//worker
		go func(f string){
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil{
				log.Println(err)
				return
			}
			info, _ :=os.Stat(thumb)  // OK to ingore error
			sizes <- info.Size()
		}(f)
	}


	//closer
	go func(){
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes{
		total += size
	}
	return total
}

注意Add和Done方法的不对策。
Add是为计数器加一，必须在worker goroutine开始之前调用，
而不是在goroutine中；
否则的话我们没办法确定Add是在“closer” goroutine调用Wait之前被调用。
并且Add还有一个参数，但Done却没有任何参数；
其实他和 Add(-1)是等价的。
我们使用defer来确保计数器即使是在出错的情况下依然能够正确地被减掉。
上面的程序代码结构是当我们使用并发循环，
但又不知道迭代次数时很通常而且很地道的写法。

sizes channel携带了每一个文件的大小到main goroutine，
在main goroutine中使用了range loop来计算总和。
观察以下我们是怎样创建一个closer goroutine，
并让其等待worker们在关闭掉sizes channel之前退出的。
两步操作：wait和close，必须是基于sizes的循环的并发。
考虑一下另一种方案：如果等待操作被放在了mian goroutine中，
在循环之前，这样的话就永远都不会结束了，
如果在循环之后，那么又变成了不可达的部分，
因为没有任何东西去关闭这个channel，这个循环就永远不会终止。

插图8.5


8.6 示例：并发的Web爬虫（p318）
在5.6节，我们做了一个简单的web爬虫，用bfs（广度优先）算法来抓取整个网站。
在本节中，我们会让这个爬虫并行化，这样每一个彼此独立的抓取命令可以并行进行IO，
最大化利用网络资源。
crawl函数和ch5/findlinks3中的是一样的。

示例代码ch8/crawl1

主函数和5.6节中的breadthFirst（深度优先）类似。
像之前的一样，一个worklist是一个记录了需要处理的元素的队列，
每一个元素都是一个需要抓取的URL列表，
不过这一次我们用channel代替slice来做这个队列。
每一个对crawl的调用都会在他们自己的goroutine中进行，
并且会把他们抓到的链接发送回worklist。

示例代码ch8/crawl1

注意这里的crawl所在的goroutine会将link作为一个显式的参数传入，
来避免“循环变量快照”的问题（5.6.1中有讲解）。
另外注意这里将命令行参数传入worklist也是在一个另外的goroutine中进行的，
这是为了避免在main goroutine和crawler goroutine中同时向另一个goroutine通过channel发送内容时发生死锁
（因为另一边的接收操作还没有准备好）。
当然，这里我们也可以用buffered channel来解决问题，这里不再赘述。

现在爬虫可以高并发地运行起来，并且可以产生一大坨的URL了，
不过还是会有两问题。
一个问题是在运行一段时间后可能会出现在log的错误信息里：

示例结果

最初的错误信息是一个让人莫名其妙的DNS查找失败，
即使这个域名是完全可靠的。
而随后的错误信息揭示了原因：这个程序一次性创建了太多网络链接，
超过了每一个进程的打开文件数限制，继而导致了在调用net.Dial像DNS查找失败这样的问题。

这个程序实在是他妈并行了。
无穷无尽地并行化并不是什么好事，因为不管怎么说，
你的系统总是会有一些限制因素，比如CPU核心数会限制你的计算负载，
比如你的硬盘转轴和磁头数限制了你的本地磁盘IO操作频率，
比如你的网络带宽限制了你的下载速度上限，
或者是你的一个web服务器容量上限等等。
为了解决这个问题，我们可以限制并发程序所使用的资源来使之适应自己的运行环境。
对于我们的例子来说，最简单的方法就是限制对links.Extract在同一时间最多不会有超过n次调用，
这里的n是fd的limit-20，一般情况下。
这个一个夜店里限制客人数目是一个道理，只有当客人离开时，
才会允许新的客人进入店内。

我们可以用一个有容量限制的buffered channel来控制并发，
这类似于操作系统里的计数信号量概念。
从概念上讲，channel里的n个空槽代表n个可以处理内容的token（通行证），
从channel里接收一个值会释放其中的一个token，并且生成一个新的空槽位。
这样保证了在没有接收介入时最多有n个发送操作。
（这里可能我们拿channel里填充的槽来做token更直观一些，不过还是这样吧）。
由于channel里的元素类型并不重要，我们用一个零值的struct{}来作为其元素。

让我们重写crawl函数，将对links.Extract的调用操作用获取、释放token的操作包裹起来，
来确保同一时间对其只有20个调用。
信号量数量和岂能操作的IO资源数量应保持接近。

示例代码ch8/crawl2

第二个问题是这个程序永远都不会终止，即使它已近爬到了所有初始链接衍生出的链接。
（当然，除非你慎重地选择了合适的初始化URL或者已经实现了练习8.6中的深度限制，
你应该还没有意识到这个问题）。
为了使这个程序能够终止，我们需要在worklist为空或者没有crawl的goroutine在运行是退出主循环。

示例代码ch8/crawl2

这个版本中，计算器n对worklist的发送操作数量进行了限制。
每一次我们发现有元素需要被发送到worklist时，我们都会对n进行++操作，
在向worklist中发送初始的命令行参数之前，
我们也进行过一次++操作。
这里的操作++是在每启动一个crawler的goroutine之前。
主循环会在n减为0时终止，这时候说明没活可干了。

现在这个并发爬虫会比5.6节中的深度优先所有板块快上20倍，
而且不会出什么错，并且在其完成任务时也会正确地终止。

下面的程序是避免过度并发的另一种思路。
这个版本使用了原来的crawl函数，但没有使用计数信号量，取而代之用了20个长活的crawler goroutine，
这样来保证最多20个HTTP请求在并发。

示例代码ch8/crawl3

所有的爬虫goroutine现在都是被同一个channel-unseenLinks喂饱的了。
主goroutine负责拆分它从worklist里拿到的元素，
然后把没有抓过的经由unseenLinks channel发动给一个爬虫的goroutine。

seen这个map被限定在main goroutine中；
也就是说这个map只能在main goroutine中进行访问。
类似于其他的信息隐藏方式，这样的约束可以让我们从一定程度上保证程序的正确性。
例如，内部变量不能够在函数外部被访问到；
变量（2.3.4）在没有被转义的情况下是无法在函数外部访问的；
一个对象的封装字段无法被该对象的方法意外的方法访问到。
在所有的情况下，信息隐藏都可以帮助我们约束我们的程序，使其不发生意料之外的情况。

crawl函数爬到的链接在一个专有的goroutine中被发送到worklist中拉力避免死锁。
为了节省空间，这个例子的终止问题我们先不详细阐述了。

译注：拓展阅读 Handling 1 Million Requests per Minute with Go。

8.7 基于select的多路复用（p323）
下面的程序会进行火箭发射的倒计时。
time.Tick函数返回一个channel，
程序会周期性地像一个节拍器一样向这个channel发送事件。
每一个事件的值是一个时间戳，
不过更有意思的是其传送方式。

示例代码ch8/countdown1

现在我们让这个程序支持在倒计时中，用户按下return键时直接中断发射流程。
首先，我们启动一个goroutine，
这个goroutine会尝试从标准输入中调入一个单独的byte，
并且如果成功了，会向名为abort的channel发送一个值。

示例代码ch8/countdown2

现在每一次技术循环的迭代都需要等待两个channel中的其中一个返回事件了：
ticker channel当一切正常时，或者异常时返回的abort事件。
我们无法做到从每一个channel中接收信息，
如果我们这么做的话，如果第一个channel中没有事件发过来那么程序就会立刻被阻塞，
这样我们就无法收到第二个channel中发过来的事件。
这时候我们需要多路复用（multiplex）这些操作了，
为了能够多路复用，我们使用了select语句。

select{
case <- ch1:
	// ...
case x := <-ch2:
	// ...use x ...
case ch3 <- y:
	// ...
default:
	// ...
}

上面是select语句的一般形式。
和switch语句稍微有点相似，
也会有几个case和最后的default选择支。
每一个case代表一个通信操作（在某个channel上进行发送或者接收），
并且会包含一些语句组成的一个语句块。
一个接收表达式可能只包含接收表达式自身（译注：不把接收到的值赋值给变量什么的）
就像上面的第一个case，
或者包含在一个简短的变量声明中，
像第二个case里一样；
第二种形式让你能够用引用接收到的值。

select会等待case中又能够执行的case时去执行。
当条件满足时，select才会去通信并执行case之后的语句；
这时候其他同行是不会执行的。
一个没有任何case的select语句写作select{},会永远地等待下去。

让我们回到我们的火箭发射程序。
time.After函数会立即返回一个channel，
并起一个新的goroutine在经过特定的时间后向该channel发送一个独立的值。
下面的select语句会一直等待到两个事件中的一个到达，
无论是abort事件或者一个10秒经过的事件。
如果10秒经过了还没有abort事件进入，那么火箭就会发射。

示例代码ch8/countdown2

下面这个例子更微妙。
ch这个channel的buffer大小是1，所以会交替的为空或为满，
所以只有一个case可以进行下去，无论i是奇数或偶数，它都会打印0 2 4 6 8。

ch := make(chan int, 1)
for i := 0; i<10; i++{
	select{
		//因为select每次只能执行可执行的。这两个case中只能执行1个。
	case x := <-ch:
		fmt.Println(x) // "0" "2" "4" "6" "8"
	case ch <- i:
	}
}

如果多个case同时就绪时，select会随机地选择一个执行，
这样来保证每一个channel都有平等的被select的机会。
增加前一个例子的buffer大小会使其输出变得不确定，
因为当buffer既不未满也不为空时，select语句的执行情况就像是抛硬币的行为一样是随机的。

下面让我们的发生程序打印倒计时。
这里的select语句会使每次循环迭代等待一秒来执行退出操作。

示例代码ch8/countdown3

tiem.Tick函数表现得好像他创建了一个在循环中调用time.Sleep的goroutine，
每次被唤醒时发送一个事件。
当countdown函数返回时，他会停止从tick中接收事件，
但是ticker这个goroutine还依然存活，
继续徒劳地尝试从channel中发送值，然而这时候已经没有其他的goroutine会从该channel中接收值了
这被称为goroutine泄露（8.4.4）

Tick函数挺方便，但是只有当程序整个生命周期都需要这个时间时，
我们使用他才比较合适。
否则的话，我们应该使用下面这种模式：

ticker := time.NewTicker(1 * time.Second)
<-ticker.C // receiver from the ticker's channel
ticker.Stop() // cause the ticker's goroutine to terminate

有时候我们希望能够从channel中发送或接收值，并避免因为发送或者接收导致的阻塞，
尤其是当channel没有准备写好或者读时。
select语句就可以实现这样的功能。
select会有一个default来设置当其他的操作都不能够马上被处理时程序需要执行哪些逻辑。

下面的select语句会在abort channel中有值时，从其中接收值；
无值时什么都不做。
这是一个非阻塞的接收操作；反复地做这样的操作叫做“轮询channel”
select{
case <-abort:
	fmt.Printf("Launch aborted!\n")
	return
default:
	// do nothing
}

channel的零值是nil。
也许会让你觉得比较奇怪，nil的channel有时候也是有一些用处的。
因为对一个nil的channel发送和接收操作会永远阻塞，
在select语句中操作nil的channel永远都不会被select到。

这使得我们可以用nil来激活或者禁用case，来达成处理其他输入或输出时间时超时和取消的逻辑。
我们会在下一节中看到一个例子。

8.8 示例：并发的字典遍历 （p327）
在本小节中，我们会创建一个程序来生成指定目录的硬盘使用情况报告，
这个程序和Unix里的du工具比较相似。
大多数工作用下面这个walkDir函数来完成，
这个函数使用dierents函数来枚举一个目录下的所有入口。

示例代码ch8/du1

ioutil.ReadDir函数会返回一个os.FileInfo类型的slice，
os.FileInfo类型也是os.Stat这个函数的返回值。
对每一个子目录而言，walkDir会递归地调用其自身，
并且会对每一个文件也递归调用。
walkDir函数会向fileSizes这个channel发送一条消息。
这条消息包含了文件的字节大小。

下面的主函数，用了两个goroutine。
后台的goroutine调用walkDir来遍历命令行给出的每一个路径，
并最终关闭fileSizes这个channel。
主goroutine会对其从channel中接收到的文件大小进行累加，并输出其和。

示例代码ch8/du1

这个程序会在打印其结果之前卡住很长时间。

$ go build .../ch8/du1
$ ./du1 $HOME /usr /bin /etc //某一个文件夹
213201 files 62.7 GB

如果在运行的时候能够让我们知道处理进度的话想必更好。
但是，如果简单地把printDiskUsage函数调用移动到循环里会到导致其打印出成百上千个输出。

下面这个du的变种会间歇打印内容，
不过只有在调用时提供了-v的flag才会显示程序进度信息。
在roots目录上循环的后台goroutine在这里保持不变。
主goroutine现在使用了计时器来每500ms生成事件，
然后用select语句来等待文件大小的消息来更新总大小数据，
或者一个计时器的事件来打印当前的总大小数据。
如果-v的flag在运行时没有传入的话，
tick这个channel会保持为nil，
这样在select里的case也就相当于被禁用了。

示例代码ch8/du2

由于我们的程序不再使用range循环，
第一个select的case必须显式地判断fileSizes的channel是不是已经被关闭了，
这里可以用到channel接收的二值形式。
如果channel已经被关闭了的话，程序会直接退出循环。
这里的break语句用到了标签break，
这样可以同时终结select和for两个循环；
如果没有用标签就break的话只会退出内层的select循环，
而外层的for循环会使之进入下一轮select循环。

现在程序会悠闲地为我们打印更新流：

$ go build .../ch8/du2
$ ./du2 -v $HOME /usr /bin /etc
28608 files 8.3 GB
54147 files 10.3 GB
93591 files 15.1 GB
127169 files 52.9 GB
175931 files 62.2 GB
213201 files 62.7 GB

然而这个程序还是会花上很长时间才会结束。
无法对walkDir做并行化处理没有什么别的原因，
无非是因为磁盘系统并行限制。
下面这个第三个版本的du，会对每一个walkDir的调用创建一个新的goroutine。
他使用sync.WaitGroup（8.5）来对仍旧活跃的walkDir调用进行计数，
另一个goroutine会在计数器减为零的时候将fileSizes这个channel关闭。

示例代码ch8/du3

由于这个程序在高峰期会创建成百上千的goroutine，
我们需要修改dirents函数，用于技术信号量来阻止它同时打开太多的文件，
就像我们在8.7节中的并发爬虫一样：

// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo{
	sema <- struct{}{} // acquire token
	defer func(){ <- sema }() // release token
	//...
}

这个版本币之前那个快了好几倍，尽管其具体效率还是和你的运行环境，机器配置相关。


8.9 并发的退出（p332）
有时候我们需要通知goroutine停止它正在干的事情，
比如一个正在执行计算的web服务，
然而他的客户端已经断开了和服务端的连接。

Go语言并没有提供在一个goroutine终止另一个goroutine的方法，
由于这样会导致goroutine之间的共享变量落在未定义的状态上。
在8.7节中rocket launch程序中，我们往名字叫abort的channel里发送了一个简单的值，
在countdown的goroutine中会把这个值理解为自己的退出信号。
但是如果我们想要退出两个或者任意多个goroutine怎么办呢？

一种可能的手段是向abort的channel里发送和goroutine数目一样多的时间来退出他们。
如果这些goroutine中已经有一些自己退出了，
那么会导致我们的channel里的时间数比goroutine还多，
这样导致我们的发送直接被阻塞。
另一方面，如果这些goroutine可能会无法接收到退出消息。
一般情况下我们是很难知道在某一个时刻具体有多少个goroutine在运行着的。
另外，当一个goroutine从abort channel中接收到一个值的时候，他会消费掉这个值，
这样其他的goroutine就没法看到这条信息。
为了能够达到我们退出goroutine的目的，
我们需要更靠谱的策略，
来通过一个channel把消息广播出去，这样goroutine们能够看到这条事件消息，
并且在事件完成之后，可以知道这件事已经发生过了。

回忆一下我们关闭了一个channel并且被消费掉了所有已发送的值，
操作channel之后的代码可以理解被执行，并且会产生零值。
我们可以将这个机制扩展一下，来作为我们的广播机制：
不要向channel发送值，而是用关闭一个channel来进行广播。

只要一些小修改，我们就可以把退出逻辑加入到前一节的du程序。
首先，我们创建一个退出的channel，这个channel不会向其中发送任何值，
但其所在的闭包内要写明程序需要退出。
我们同时还定义了一个工具函数，cancelled，
这个函数在被调用的时候会轮询退出状态。

示例代码ch8/du4
var done = make(chan struct{})

func cancelled() bool{
	select {
	case <- done:
		return true
	default:
		return false
	}
}

下面我们创建一个从标准输入流中读取内容的goroutine，
这是一个比较典型的连接到终端的程序。
每当有输入被读到（比如用户按了回车键），
这个goroutine就会把取消消息通过关闭done的channel广播出去。

// Cancel traversal when input is detected.（当检测到输入时，取消遍历）
go func(){
	os.Stdin.Read(make([]byte, 1)) // read a single byte
	close(done)
}()

现在我们需要使我们的goroutine来对取消进行响应。
在main goroutine中，我们添加了select的第三个case语句，
尝试从done channel中接收内容。
如果这个case被满足的话，在select到的时候即会返回，
但在结束之前我们需要把fileSizes channel中的内容“排”空，
在channel被关闭之前，舍弃掉所有值。
这样可以保证对walkDir的调用不要被向fileSizes发送信息阻塞住，
可以正确地完成。

for{
	select{
	case <- done:
		// Drain fileSizes to allow existing goroutine to finish.
		for range fileSizes{
			// Do nothing.
		}
		return
	case size, ok:= <-fileSizes:
		//...
	}
}

walkDir这个goroutine一启动就会轮询取消状态，
如果取消状态被设置的话会直接返回，
并且不做额外的事情。
这样我们将所有在取消事件之后创建的goroutine改变为无操作。

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64){
	defer n.Done()
	if cancelled(){
		return
	}
	for _, entry := range dirents(dir){
		// ...
	}
}

在walkDir函数的循环中我们对取消状态进行轮询可以带来明显的益处，
可以避免在取消事件发生时还去创建goroutine。
取消本身是有一些代价的；
想要快速的响应需要对程序逻辑进行侵入式的修改。
确保在取消发生之后不要有代价太大的操作，可能会需要修改你代码里的很多地方，
但是在一些重要的地方去检查取消时间也确实能带来很大的好处。

对这个程序的一个简单的性能分析可以揭示瓶颈在dirents函数中获取一个信号量。
下面的select可以让这种操作可以被取消，
并且可以将取消时的延迟从几百毫秒降低到几十毫秒。

func dirents(dir string) []os.FileInfo{
	select{
	case sema<-struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func(){<-sema}() // release token
}

现在当取消发生时，所有后台的goroutine都会迅速停止并且主函数会返回。
当然，当主函数返回时，一个程序会退出，而我们又无法在主函数退出的时候确认其已经释放了所有的资源。
（因为程序都退出了，你的代码都没法执行了）。
这里有一个方便的窍门我们可以一用：
取代掉直接从主函数返回，我们调用一个panic，然后runtime会把每一个goroutine的栈dump下来。
如果main goroutine是唯一一个剩下的goroutine的话，
他会清理掉自己的一切资源。
但是如果还有其他goroutine没有退出，他们可能没办法别正确地取消掉，也有可能被取消但是取消操作会花时间；
所以这里的一个调研还是很有必要的。
我们用panic来获取到足够的信息来验证我们上面的判断，看看最终到底是什么样的情况。

8.10 示例：聊天服务（p335）
我们用一个聊天服务器来终结本章的内容，这个程序可以让一些用户通过服务器向其他所有用户广播文本消息。
这个程序中有四种goroutine。
main和broadcaster各自是一个goroutine实例，
每一个客户端的连接都会有一个handleConn和clientWriter的goroutine。
broadcaster是select用法的不错的样例，因为它需要处理三种不同类型的消息。

下面演示的main goroutine的工作，是listen和accept（译注：网络编程里的概念）从客户端过来的连接。
每一个连接，程序都会建立一个新的handleConn的goroutine，
就像我们在本章开头的并发的echo服务器里所做的那样。

代码示例ch8/chat
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

然后是broadcaster的goroutine。
他的内容变量clients会记录当前建立连接的客户端集合。
其记录的内容每一个客户端的消息发出channel的“资格”信息。

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
	select {
	case msg := <-messages:
		// Broadcast incoming message to all
		// clients' outgoing message channels.
		for cli := range clients {
			cli <- msg
		}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

broadcaster监听来自全局的entering和leaving的channel来获知客户端的到来和离开事件。
当其接收到其中的一个事件时，会更新clients集合，
当该事件是离开行为时，他会关闭客户端的消息发出channel。
broadcaster也会监听全局的消息channel，所有的客户端都会向这个channel中发送消息。
当broadcaster接收到什么消息时，就会将其广播至所有连接到服务端的客户端。

现在让我们看看每一个客户端的goroutine。
handleConn函数回味他的客户端创建一个消息发出channel，
并通过entering channel来通知客户端的到来。
然后它会读取客户端发来的每一行文本，
并通过全局的消息channel来将这些文本发送出去，
并为每条消息带上发送者的前缀来标明消息身份。
当客户端发送完毕后，handleConn会通过leaving这个channel来通知客户端的离开关闭连接。

示例代码ch8/chat

另外，handleConn为每一个客户端创建了一个clientWriter的goroutine来接收向客户端发出消息channel中发送的广播消息，
并将他们写入到客户端的网络连接。
客户端的读取方循环会在broadcaster接收到leaving通知并关闭了channel后终止。

下面演示的是当服务器有两个活动的客户端连接，并且在两个窗口中运行的情况，使用netcat来聊天：

$ go build gopl.io/ch8/chat
$ go build gopl.io/ch8/netcat3
$ ./chat &
$ ./netcat3
You are 127.0.0.1:64208 		$ ./netcat3
127.0.0.1:64211 has arrived 	You are 127.0.0.1:64211
Hi!
127.0.0.1:64208: Hi!
								127.0.0.1:64208: Hi!
								Hi yourself.
127.0.0.1:64211: Hi yourself. 	127.0.0.1:64211: Hi yourself.
^C
								127.0.0.1:64208 has left
$ ./netcat3
You are 127.0.0.1:64216 		127.0.0.1:64216 has arrived
								Welcome.
127.0.0.1:64211: Welcome. 		127.0.0.1:64211: Welcome.
								^C
127.0.0.1:64211 has left”

当与n个客户端保持聊天session时，
这个程序会有2n+2个并发的goroutine，
然而这个程序却并不需要显式的锁（9.2）。
clients这个map被限制在了一个独立的goroutine中，
broadcaster，所以它不能被并发地访问。
多个goroutine共享的变量只有这个channel和net.Conn的实例，
两个东西都是并发安全的。
我们会在下一章中更多地解决约束，
并发安全以及goroutine中共享变量的含义。

