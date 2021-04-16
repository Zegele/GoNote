第九章 基于共享变量的并发 
（P339-372）

前一章我们介绍了一些使用goroutine和channel这样直接而自然的方式来实现并发的方法。
然而这样做我们实际上屏蔽掉了在写并发代码时必须处理的一些重要而且细微的问题。

在本章中，我们会细致地了解并发机制。
尤其是在多goroutine之间的共享变量，
并发问题的分析手段，
以及解决这些问题的基本模式。
最后我们会解释goroutine和操作系统线程之间的技术上的一些区别。

9.1 竞争条件（p340）
在一个线性（就是说只有一个goroutine的）的程序中，
程序的执行顺序只由程序的逻辑来决定。
例如，我们有一段语句序列，第一个在第二个之前（废话），以此类推。
在有两个或更多goroutine的程序中，
每一个goroutine内的语句也是按照既定的顺序去执行的，
x是在y之前还是之后，还是同时发生是没法判断的。
当我们没有办法自信地确认一个事件是在另一个事件的前面或者后面发生的话，
就说明x和y这两个事件是并发的。

考虑下，一个函数在线性程序中可以正确地工作。
如果在并发的情况下，这个函数依然可以正确地工作的话，
那么我们就说这个函数是并发安全的，
鬓发安全的函数不需要额外的同步工作。
我们可以把这个概念概括为一个特定类型的一些方法和操作函数，
如果这个类型是并发安全的话，那么所有他的访问方法和操作就都是并发安全的。

在一个程序中有非并发安全的类型的情况下，我们依然可以使这个程序并发安全。
确实，并发安全的类型是例外，而不是规则，
所以只有当文档中明确地说明了其是并发安全的情况下，
你才可以并发地去访问它。
我们会避免并发访问大多数的类型，
无论是将变量局限在单一的一个goroutine内，
还是用互斥条件维持更高级别的不变性都是为了这个目的。
我们会在本章中说明这些术语。

相反，导出包级别的函数一般情况下都是并发安全的。
由于package级的变量没法被限制在单一的goroutine，
所以修改这些变量“必须”使用互斥条件。

一个函数在并发调用时没法工作的原因太多了，
比如死锁（deadlock），活锁（livelock）和饿死（resource starvation）。
我们没有空去讨论所有问题，这里我们只聚焦在竞争条件上。

竞争条件指的是程序在多个goroutine交叉执行操作时，
没有给出正确的结果。
竞争条件是很恶劣的一种场景，
因为这种问题会一直潜伏在你的程序里，
然后在非常少见的时候蹦出来，
或许只是会在很大的负载时才会发生，
又或许是会在使用了某一个编译器、某一种平台，
或某一种架构的时候才会出现。
这些使得竞争条件带来的问题非常难以复现，
而且难以分析诊断。

传统上经常用经济损失来为竞争条件做比喻，所以我们来看一个简单的银行账户程序。

// Package bank implements a bank with only one account.
package bank

var balance int
func Deposit(amount int) {balance = balance + amount}
func Balance()int {return balance} //balance余额

(当然我们也可以把Deposit存款函数写成balance+=amount，这种形式也是等价的，
不过长一些的形式解释起来更方便一些。)

对于这个具体的程序而言，我们可以瞅一眼各种存款和余额的顺序调用，
都能给出正确的结果。
也就是说，Balance函数会给出之前的所有存入的额度之和。
然而，当我们并发地而不是顺序地调用这些函数的话，
Balance就再也没办法保证结果正确了。
考虑一下下面的两个goroutine，
其代表了一个银行联合账户的两笔交易：

// Alice:
go func(){
	bank.Deposit(200) // A1
	fmt.Println("=", bank.Balance()) // A2
}()

// Bob:
go bank.Deposit(100) // B

Alice存了$200，然后检查她的余额，同时Bob存了$100。
因为A1和A2是和B并发执行的，我们没法预测他们发生的先后顺序。
直观地来看的话，我们会认为其执行顺序只有三种可能性：“Alice先”，
“Bob先”以及“Alice/Bob/Alice”交错执行。
下面的表格会展示经过每一步骤后balance变量的值。
引号里的字符串表示余额单。

Alice first | Bob first | Alice/Bob/Alice
0				0			0
A1 200			B  100		A1 200
A2 "=200"		A1 300		B  300
B  300			A2 "=300" 	A2 "=300"

所有情况下最终的余额都是$300。
唯一的变数是Alice的余额单是否包含了Bob交易，
不过无论怎么着客户都不会在意。

但是事实上面的直觉推断是错误的。
第四种可能结果是事实存在的，
这种情况下Bob的存款会在Alice存款操作中间，
在余额被读到（balance + amount）之后，
在余额被更新之前（balance = ...），
这样会导致Bob的交易丢失。
而这是因为Alice的存款操作A1实际上是两个操作的一个序列，
读取然后写；可以称之为A1r和A1w。
下面是交叉时产生的问题：
Data	race
0		
A1r		0			... = balance + amount
B 		100
A1w		200			balance = ...
A2	  	"= 200"

在A1r之后，balance + amount会被计算为200，
所以这是A1w会写入的值，
并不受其他存款操作的干扰。
最终的余额是$200。
银行的账户上的资产比Bob实际的资产多了$100。
（译注：因为丢失了Bob的存款操作，所以其实是说Bob的钱丢了）

这个程序包含了一个特定的竞争条件，叫作数据竞争。
无论任何时候，只要有两个goroutine并发访问同一变量，
且至少其中的一个是写操作的时候就会发生数据竞争。

如果数据竞争的对象是一个比一个机器字
（译注：32位机器上一个字=4个字节）
更大的类型时。
事情就变得更麻烦了，
比如interface，string或者slice类型都是如此。
下面的代码会并发地更新两个不同长度的slice：

var x []int
go func(){x = make([]int, 10)}()
go func(){x = make([]int, 1000000)}()
x[999999] = 1 // NOTE: undefined behavior; memory corruption possible!

最后一个语句中的x的值是未定义的；
其可能是nil，或也可能是一个长度为10的slice，
也可能是一个长度为1,000,000 的slice。
但是回忆以下slice的三个组成部分：
指针（pointer）、长度（length)和容量（capacity）。
如果指针是从第一个make调用来，
而长度从第二个make来，x就变成了一个混合体，一个自称长度为1000000，
但实际上内部只有10个元素的slice。
这样导致的结果是存储999999元素的位置会碰撞一个遥远的内存位置，
这种情况下难以对值进行预测，
而且定位和debug也会变成噩梦。
这种语义雷区被称为未定义行为，
对C程序员来说应该很熟悉；
幸运的是在Go语言里造成的麻烦要比C里小得多。

尽管并发程序的概念让我们知道并发并不是简单的语句交叉执行。
我们将会在9.4节中看到，
数据竞争可能会有奇怪的结果。
许多程序员，甚至一些非常聪明的人也还是会偶尔提出一些理由来允许数据竞争。
比如：“互斥条件代价太高”，“这个逻辑只是用来做logging”，
“我不介意丢失一些消息”等等。
因为在他们的编译器或者平台上很少遇到问题，
可能给了他们错误的信心。
一个好的经验法则是根本就没有什么所谓的良性数据竞争。
所以我们一定要避免数据竞争，那么在我们的程序中要如何做到呢？

我们来重复一下数据竞争的定义，
因为是在太重要了：数据竞争会在两个以上的goroutine并发访问相同的变量，
且至少其中一个为写操作时发生。
根据上述定义，有三种方式可以避免数据竞争：

第一种方法是不要去写变量。
考虑以下下面的map，
会被“懒填充”，
也就是说在每个key被第一次请求到的时候才回去填值。
如果Icon是被顺序调用的话，这个程序会工作很正常，
但如果Icon被并发调用，那么对于这个map来说就会存在数据竞争。

var icons = make(map[string]image.Image)
func loadIcon(name string) image.Image

// NOTE: not concurrency-safe!
func Icon (name string) image.Image{
	icon, ok := icons[name]
	if !ok {
		icon = loadIcon(name)
		icons[name] = icon
	}
	return icon
}

反之，如果我们在创建goroutine之前的初始化阶段，
就初始化了map中的所有条目并且再也不去修改他们，
那么任意数量的goroutine并发访问icon都是安全的，
因为每一个goroutine都只是读取而已。

var icons = map [string]image.Image{
	"spades.png": loadIcon("spades.png"),
	"hearts.png": lpadIcon("hearts.png"),
	"diamonds.png": loadIcon("diamonds.png"),
	"clubs.png":loadIcon("clubs.png"),
}

// Concurrency-safe.
func Icon(name string) image.Image{return icons[name]}

上面的例子里icons变量在包初始化阶段就已经被赋值了，
包的初始化是在程序main函数开始执行之前就完成了的。
只要初始化完成了，icons就再也不会修改的或者不变量是本来就并发安全的，
这种变量不需要进行同步。
不过显然我们没法用这种方法，因为update操作是必要的操作，
尤其对于银行账户来说。

第二种避免数据竞争的方法是，避免从多个goroutine访问变量。
这也是前一章中大多数程序所采用的方法。
例如前面的并发web爬虫（8.6）的main goroutine是唯一一个能够访问seen map的goroutine，
而聊天服务器（8.10）中的broadcaster goroutine是唯一一个能够访问clients mao的goroutine。
这些变量都被限定在了一个单独的goroutine中。

由于其他的goroutine不能直接访问变量，他们只能使用一个channel来发送给制定的goroutine请求来查询更新变量。
这也就是Go的口头禅“不要使用共享数据来通信；使用通信来共享数据”。
// 不要用共享数据来通信，而要用通道（channel）来共享数据，这样通信。
一个提供对一个指定的变量通过channel来请求的goroutine就做这个变量的监控（monitor）goroutine。
例如broadcaster goroutine会监控（monitor）clients map的全部访问。

下面是一个重写了银行的例子，这个例子中balance变量被限制在了monitor goroutine中，
名为teller：

示例代码ch9/bank1

即使当一个变量无法在其整个生命周期内被绑定到一个独立的goroutine，
绑定依然是并发问题的一个解决方案。
例如在一条流水线上的goroutine之间共享变量是很普遍的行为，
在这两者间会通过channel来传输地址信息。
如果流水线的每一个阶段都能够避免在将变量传送到下一阶段时再去访问它，
那么对这个变量的所有访问就是线性的。
其效果是变量会被绑定到流水线的一个阶段，传送万之后被绑定到下一个，
以此类推。
这种规则有事呗成为串行绑定。

下面的例子中，Cakes会被严格地顺序访问，显示baker goroutine， 然后是icer goroutine：

type Cake struct{ state string}

func baker(cooked chan<- *Cake){ //cooked 是指针类型的channel，且该channel是输出类型的channel（数据输出给channel）
	for {
		cake := new(Cake)
		cake.state = "cooked"
		cooked <- cake // baker never touches this cake again
	}
}

func icer(iced chan<- *Cake, cooked <-chan *Cake){
	for cake := range cooked{
		cake.state = "iced"
		iced <- cake // icer never touches this cake again
	}
}

第三种避免数据竞争的方法是允许很多goroutine去访问变量，
但是在同一个时刻最多只有一个goroutine在访问。
这种方式被称为“互斥”，在下一节来讨论这个主题。

9.2 sync.Mutex互斥锁（p346）
在8.6节中，我们使用了一个buffered channel作为一个技术信号量，
来保证最多只有20个goroutine会同时执行HTTP请求。
同理，我们可以用一个容量只有1的channel来保证最多只有一个goroutine在同一时刻访问一个共享变量。
一个只能为1和0的信号叫做二元信号量（binary semaphore）。

示例代码ch9/bank2

这个种互斥很实用，而且被sync包里的Mutex类型直接支持。
它的Lock方法能过获取到token（这里叫锁），并且Unlock方法会释放这个token：


示例代码ch9/bank3

每次一个goroutine访问bank变量时（这里只有balance余额变量），
它都会调用mutex的Lock方法来获取一个互斥锁。
如果其他的goroutine已经获得了这个锁的话，
这个操作会被阻塞知道其他goroutine调用了Unlock使该锁便会可用状态。
mutex会保护共享变量。
惯例来说，被mutex所保护的变量是在mutex变量声明之后立刻声明的。
如果你的做法和惯例不符，确保在文档里对你的做法进行说明。

在Lock和Unlock之间的代码段中的内容goroutine可以随便读取或者修改，
这个代码段叫做临界区。
goroutine在结束后释放锁是必要的，
无论以哪条路径通过函数都需要释放，
即使是在错误路径中，也要记得释放。 

上面的bank程序例证了一种通用的并发模式。
一系列的导出程序封装了一个或多个变量，
那么访问这些变量唯一的方式就是通过这些函数来做（或者方法，对于一个对象的变量来说）。
每一个函数在一开始就获取互斥锁并在最后释放锁，
从而保证共享变量不会被并发访问。
这种函数、互斥锁和变量的编排叫做监控（monitor）
（这种老式单词的monitor是受“monitor goroutine”的术语启发而来的。
两种用法都是一个代理人保证变量被顺序访问。）

由于在村换和查询余额函数中的临界区代码这么短--只有一行，没有分支调用--
在代码最后去调用Unlock就显得更为直截了当。
在更复杂的临界区的应用中，尤其是必须要尽早处理错误并返回的情况下，
就很难去（靠人）判断对Lock和Unlock的调用是在所有路径都能够严格配对的了。
Go语言里的defer简直就是这种情况下的救星：
我们用defer来调用Unlock，临界区会隐式地延伸到函数作用域的最后，
这样我们就从“总要记得在函数返回之后或者发生错误返回时要记得调用一次Unlock”
这种状态中获得了解放。
Go会自动帮我们完成这些事情。

func Balance()int{
	mu.Lock()
	defer mu.Unlock()
	return balance
}

上面的例子里Unlock会在return语句读取完balance的值之后执行，
所以Balance函数是并发安全的。
这带来的另一点好处是，我们再也不需要一个本地变量b了。

此外，一个deferred Unlock函数即使在临界区发生panic时依然会执行，
这对于用recover（5.10）来恢复的程序来说是很重要的。
defer调用只会比显式地调用Unlock成本高那么一点点，
不过却在很大程度上保证了代码的整洁性。
大多数情况下对于并发程序来说，
代码的整洁性比过度的优化更重要。
如果可能的话尽量使用defer来将临界去扩展到函数的结束。

考虑一下下面的Withdraw函数。
成功的时候，他会正确地减掉余额并返回true。
但如果银行记录资金对交易来说不足，
那么取款就会回复余额，并返回false。

// NOTE: not atomic!
func Withdraw(amount int)bool{
	Deposit(-amount)
	if Balance()<0{
		Deposit(amount)
		return false // insufficient funds
	}
	return true
}

函数终于给出了正确的结果，但是还有一点讨厌的副作用。
当过多的取款操作同时执行时，
balance可能会瞬时被减到0以下。
这可能会引起一个并发的取款被不合逻辑地拒绝。
所以如果Bob尝试买一辆sports car时，
Alice可能就没办法为她的早咖啡付款了。
这里的问题是取款不是一个原子操作：
他包含了三个步骤，每一步都需要去获取并释放互斥锁，
但任何一次锁都不会锁上整个取款流程。

理想情况下，取款应该只在整个操作中获得一次互斥锁。
下面这样的尝试是错误的：

// NOTE: incorrect!
func Withdraw(amount int) bool{
	mu.Lock()
	defer mu.Unlock()
	Deposit(-amount)
	if Balance() <0 {
		Deposit(amount)
		return false // insufficient funds
	}
	return true
}

上面这个例子中，Deposit会调用mu.Lock()第二次去获取互斥锁，
但因为mutex已经锁上了，
而无法被重入（译注：go里没有重入锁，关于重入锁的概念，请参考java）--
也就是说没法对一个已经锁上的mutex来再次上锁--
这会导致程序死锁，
没法继续执行下去，
Withdraw会永远阻塞下去。

关于Go的互斥量不能重入这一点我们有很充分的理由。
互斥量的目的是为了确保共享变量在程序执行时的关键点上能够保证不变性。
不变性的其中之一是“没有goroutine访问共享变量”。
但实际上对于mutex保护的变量来说，不变性还包括其他方面。
当一个goroutine获得了一个互斥锁时，
他会断定这种不变性能够被保持。
其获取并保持锁期间，可能会去更新共享变量，
这样不变性只是短暂地被破坏。
然而档期解放锁之后，他必须保证不变性已经恢复原样。
尽管一个可以重入的mutex也可以保证没有其他的goroutine在访问共享变量，
但这种方式没法保证这些变量额外的不变性。

一个通用的解决方案是将一个函数分离为多个函数，比如我们把Deposit分离成两个：
一个不导出的函数deposit，这个函数假设锁总是被保持并去做实际的操作，
另一个是导出的函数（Deposit），这个函数会调用deposit，
但在调用前会先去获取锁。
同理我们可以将Withdraw也表示成这种形式：

func Withdraw(amount int) bool{
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance<0{
		deposit(amount)
		return false
	} // insufficient funds
	return true
}

func Deposit(amount int){
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance()int{
	mu.Lock()
	defer mu.Unlock()
	return balance
}

// This function requires that the lock be held.
func deposit(amount int){
	balance += amount
}

当然，这里的存款deposit函数很小实际上取款withdraw函数不需要理会对它的调用，
尽管如此，这里的表达还是表明了规则。

封装（6.6），用限制一个程序中的意外交互的方式，
可以使我们获得数据结果的不变性。
因为某种原因，封装还帮我们获得了并发的不变性。
当你使用mutex时，确保mutex和其保护的变量没有被导出
（在go里也就是小写，且不要被大写字母开头的函数访问了），
无论这些变量是包级的变量还是一个struct的字段。

9.3 sync.RWMutex 读写锁（p351）
在100刀的存款消失时不做记录多少还是会让我们有些恐慌，
Bob写了一个程序，每秒运行几百次来检查他的银行余额。
他会在家，在工作中，甚至会在他的手机上来运行这个程序。
银行注意到这些陡增的流量使得存款和取款有了延时，
因为所有的余额查询请求是顺序执行的，
这样会互斥地获得锁，并且会暂时组织其他的goroutine运行。

由于Balance函数只需要读取变量的状态，
所以我们同时让多个Balance调用并发运行事实上是安全的，
只要在运行的时候没有存款或者取款操作就行。
在这种场景下我们需要一种特殊类型的锁，其允许多个只读操作并行执行，
但写操作会完全互斥。
这种锁叫做“多读单写”锁（multiple readers, single writer lock），
Go语言提供的这样的锁是sync.RWMutex:

var mu sync.RWMutex
var balance int
func Balance() int {
	mu.RLock() // readers lock
	defer mu.RUnlock()
	return balance
}

Balance函数现在调用了RLock和RUnlock方法来获取和释放一个读取或者共享锁。
Deposit函数没有变化，会调用mu.Lock和mu.Unlock方法来获取和释放一个写或互斥锁。

在这次修改后，Bob的余额查询请求就可以彼此并行地执行并且会很快地完成了。
锁在更多的时间范围可用，并且存款请求也能够及时地被响应了。

RLock只能在临界区共享变量没有任何写入操作时可用。
一般来说，我们不应该假设逻辑上的只读函数/方法也不会去更新某一些变量。
比如一个方法功能是访问一个变量，
但它也有可能会同时去给一个内部的计数器+1
（译注：可能是记录这个方法的访问次数啥的）
或者去更新缓存--使即时的调用能够更快。
如果有疑惑的话，请使用互斥锁。

RWMutex只有当获得锁的大部分goroutine都是读操作，
而锁在竞争条件下，也就是说，
goroutine们必须等待才能获取到锁的时候，
RWMutex才是最能带来好处的。
RWMutex需要更复杂的内部记录，所以会让它币一般无竞争锁的mutex慢一些。

9.4 内存同步（p352）
你可能比较纠结为什么Balance方法需要用到互斥条件，
无论是基于channel还是基于互斥量。
毕竟和存款不一样，它只由一个简单的操作组成，
所以不会碰到其他goroutine在其执行“中”执行其他的逻辑的风险。
这里使用mutex有两方面考虑。
第一Balance不会在其他操作比如Withdraw“中间”执行。
第二（更重要）的是“同步”不仅仅是一堆goroutine执行顺序的问题；
同样也会涉及到内存的问题。

在现代计算机中可能会有一堆处理器，每一个都会有其本地缓存（local cache）。
为了效率，对内存的写入一般会在每一个处理器中缓存，并在必要时一起flush到主存。
这种情况下这些数据可能会以与当初goroutine写入顺序不同的顺序被提交到主存。
像channel通信或者互斥量操作这样的原语会使处理器将其聚集的写入flush并commit，
这样goroutine在某个时间点上的执行结果才能被其他处理器上运行的goroutine得到。

考虑以下代码片段的可能输出：
var x, y int
go func(){
	x = 1 // A1
	fmt.Print("y:",y, " ") // A2
}()
go func(){
	y = 1 // B1
	fmt.Print("x:", x, " ") // B2
}()

因为两个goroutine是并发执行，并且访问共享变量时也没有互斥，
会有数据竞争，所以程序的运行结果没法预测的话也请不要惊讶。
我们可能希望它能够打印出下面这四种结果中的一种，
相当于几种不同的交错执行时的情况：
y:0 x:1
x:0 y:1
x:1 y:1
y:1 x:1

第四行可以被解释为执行顺序A1,B1,A2,B2或者B1,A1,A2,B2的执行结果。
然而实际的运行时还是有些情况让我们有点惊讶：
x:0 y:0
y:0 x:0

但是根据所使用的编译器，CPU，或者其他很多影响因子，
这两种情况也是有可能发送的。
那么这两种情况要怎么解释呢？

在一个独立的goroutine中，每一个语句的执行顺序是可以被保证的；
也就是说goroutine是顺序连贯的。
但是在不使用channel，且不使用mutex这样的显式同步操作时，
我们就没法保证事件在不同的goroutine中看到的执行顺序是一致的了。
尽管goroutine A中一定需要观察到x=1执行成功之后才会去读取y，
但它没法确保自己观察得到goroutine B中对y的写入，
所以A还可能会打印出y的一个旧版的值。

尽管去理解并发的一种尝试是去将其运行理解为不同goroutine语句的交错执行，
但看看上面的例子，这已经不是现代编译器和cpu的工作方式了。
因为赋值和打印指向不同的变量，
编译器可能会断定两条语句的顺序不会影响执行结果，
并且会交换两个语句的执行顺序。
如果两个goroutine在不同的CPU上执行，每一个核心有自己的缓存，
这样一个goroutine的写入对于其他goroutine的Print，
在主存同步之前就是不可见的了。

所有并发的问题都可以用一致的、简单的既定模式来规避。
所以可能的话，将变量限定在goroutine内部；
如果是多个goroutine都需要访问的变量，使用互斥条件来访问。

9.5 sync.Once初始化（p354）
如果初始化成本比较大的话，那么将初始化延迟到需要的时候再去做就是一个比较好的选择。
如果在程序启动的时候就去做类似的初始化的话，会增加程序的启动时间，
并且因为执行的时候可能也并不需要这些变量，
所以实际上有一些浪费。
让我们在本章早一些时候看到的icons变量：

var icons map[string]image.Image

这个版本的Icon用到了懒初始化（lazy initialization）

func loadIcons(){
	icons = map[string]image.Image{
		"spades.png": loadIcon("spades.png"),
		"hearts.png": loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png": loadIcon("clubs.png"),
	}
}

// NOTE: not concurrency-safe!
func Icon(name string) image.Image{
	if icons == nil{
		loadIcons() // one-time initialization
	}
	return icons[name]
}

如果一个变量只被一个单独的goroutine所访问的话，我们可以使用上面的这种模板，
但这种模板在Icon被并发调用时并不安全。
就像前面银行的那个Deposit（存款）函数一样，
Icon函数也是由多个步骤组成的：
首先测试icons是否为空，然后load这些icons，
之后将icons更新为一个非空的值。
直觉会告诉我们最差的情况是loadIcons函数被多次访问会带来数据竞争。
当第一个goroutine在忙着loading这些icons的时候，另一个goroutine进入了Icon函数，
发现变量是nil，然后也会调用loadIcons函数。

不过这种直觉是错误的。
（我们希望现在你从开始能够构建自己对并发的直觉，也就是说对并发的直觉总是不能被信任的！）
回忆一下9.4节。
因为缺少显示的同步，编译器和cpu是可以随意地去更改访问内部的指令顺序，
以任意方式，只要保证每一个goroutine自己的执行顺序一致。
其中一种可能loadIcons的语句重排是下面这样。
他会在填写icons变量的值之前先用一个空map来初始化icons变量。

func loadIcons(){
	icons = make(map[string]image.Image)
	icons["spades.png"] = loadIcon("spades.png")
	icons["hearts.png"] = loadIcon("hearts.png")
	icons["diamonds.png"] = loadIcon("diamonds.png")
	icons["clubs.png"] = loadIcon("clubs.png")
}

因此一个goroutine在检查icons是非空时，
也并不能就假设这个变量的初始化流程已经走完了。
（译注：可能只是塞了个空map，里面的值还没填完，也就是说填值的语句都没执行完呢）。

最简单且正确的保证所有goroutine能够观察到loadcons效果的方式，
是用一个mutex来同步检查。

var mu sync.Mutex // guards icons
var icons map[string]image.Image

// Concurrency-safe.
func Icon(name string) image.Image{
	mu.Lock()
	defer mu.Unlock()
	if icons == nil{
		loadIcons()
	}
	return icons[name]
}

然而使用互斥访问icons的代价就是没有办法对该变量进行并发访问，
即使变量已经被初始化完毕，且再也不会进行变动。
这里我们可以引入一个允许多读的锁：

var mu sync.RWMutex // guards icons
var icons map[string]image.Image
// Concurrency-safe.
func Icon(name string) image.Image{
	mu.RLock()
	if icons != nil{
		icon := icons[name]
		mu.RUnlock()
		return icon
	}
	mu.RUnlock()

	// acquire an exclusive lock
	mu.Lock()
	if icons == nil{ // NOTE: must recheck for nil
		loadIcons()
	}
	icon := icons[name]
	mu.Unlock()
	return icon
}

上面的代码有两个临界区。
goroutine首先会获取一个写锁，查询map，然后释放锁。
如果条目被找到了（一般情况下），
那么会直接返回。
如果没有找到，那goroutine会获取一个写锁。
不释放共享锁的话，也没有任何办法来将一个共享锁升级为一个互斥锁，
所以我们必须重新检查icons变量是否为nil，
以防止在执行这一段代码的时候，
icons变量已经被其他goroutine初始化过了。

上面的模板使我们的程序能够更好的并发，
但是有一点太复杂且容易出错。
幸运的是，sync包为我们提供了一个专门的方案来解决这种一次性初始化的问题：sync.Once
概念上来讲，一次性的初始化需要一个互斥量mutex和一个boolean变量来记录初始化是不是已经完成了；
互斥量用来保护boolean变量和客户端数据结构。
Do这个唯一的方法需要接收初始化函数作为其参数。
让我们用sync.Once来简化前面的Icon函数吧：

var loadIconsOnce sync.Once
var icons map[string]image.Image
// Concurrency-safe.
func Icon(name string) image.Image{
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}

每一次对Do（loadIcons）的调用都会锁定mutex，
并会检查boolean变量。
在第一次调用时，变量的值是false，
Do会调用loadIcons并会将boolean设置为true。
随后的调用什么都不会做，但是mutex同步会保证loadIcons对内存（这里其实就是指icons变量）
产生的效果能够对所有goroutine可见。用这种方式来使用sync.Once的话，
我们能够避免在变量被构建完成之前和其它goroutine共享该变量。

9.6 竞争条件检测（p358）
即使我们小心到不能再小心，但在并发程序中犯错还是太容易了。
幸运的是，Go的runtime和工具链为我们装备了一个复杂但好用的动态分析工具，
竞争检查器（the race detector）。

只要在go build，go run或者go test命令后面加上-race的flag，
就会使编译器创建一个你的应用的“修改”版或者一个附带了能够记录所有运行期对共享变量访问工具的test，
并且会记录下每一个读或者写共享变量的goroutine的身份信息。
另外，修改版的程序会记录下所有的同步事件，比如go语句，channel操作，
以及对 (*sync.Mutex).Lock，(*sync.WaitGroup).Wait等等的调用。
（完整的同步事件集合是在The Go Memory Model文档中说明，该文档是和语言文档放在一起的。译注：https://golang.org/ref/mem ）

竞争检查器会检查这些事件，会寻找在哪一个goroutine中出现了这样的case，
例如其读或者写了一个共享变量，这个共享变量是被另一个goroutine在没有进行干涉同步操作便直接写入的。
这种情况也就表明了是对一个共享变量的并发访问，即数据竞争。
这个公布会打印一份报告，内容包含变量身份，
读取和写入的goroutine中活跃的函数的调用栈。
这些信息在定位问题时通常很有用。
9.7节中会有一个竞争检查器的示例样例。

竞争检查器会报告所有的已经发生的数据竞争。
然而，它只能检测到运行时的竞争条件；
并不能证明之后不会发生数据竞争。
所以为了使结果尽量正确，请保证你的测试并发地覆盖到了你的包。

由于需要额外的记录，因此构建时加了竞争检测的程序跑起来会慢一些，且需要更大的内存，
即时是这样，这些代价对于很多生成环境的工作来说还是可以接受的。
对于一些偶发的竞争条件来说，让竞争检查器来干活可以节省无数日夜的debugging。
（译注：多少服务端C和C++程序员为此尽折腰）

9.7 示例：并发的非阻塞缓存（p359）
本节中我们会做一个无阻塞的缓存，这种工具可以帮助我们来解决显示世界中并发程序出现，
但没有现成的库可以解决的问题。
这个问题叫作缓存（memoizing）函数
（译注：Memoization的定义：memoization一词是Donald Michie根据拉丁语memorandum杜撰的一个词。
响应的动词，过去分词、ing形式有memoiz、memoized、memoizing），
也就是说，我们需要缓存函数的返回结果，这样在对函数进行调用的时候，
我们就只需要一次计算，之后只要返回计算的结果就可以了。
我们的解决方案会使并发安全且会避免对整个缓存加锁，
而导致所有操作都去争一个锁的设计。

我们将使用下面的httpGetBody函数作为我们需要缓存的函数的一个样例。
这个函数会去进行HTTP GET请求并且获取http响应body。
对这个函数的调用本身开销是比较大的，
所以我们尽量尽量避免在不必要的时候反复调用。

func httpGetBody(url string)(interface{}, error){
	resp, err := http.Get(url)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

最后一行稍微隐藏了一些细节。ReadAll会返回两个结果，一个[]byte数组和一个错误，
不过这两个对象可以被赋值给httpGetBody的返回声明里的interface{}和error；类型，
所以我们也就 可以这样返回结果并且不需要额外的工作了。
我们在httpGetBody中选用这种返回类型视为了使其可以与缓存匹配。

下面是我们要设计的cache的第一个“草稿”：

示例代码ch9/memo1

Memo实例会记录需要缓存的函数f（类型为Func），
以及缓存内容（里面是一个string到result映射的map）。
每一个result都是简单的函数返回的值对儿--一个值和一个错误值。
继续下去我们会展示一些Memo的变种，不过所有的例子都会遵循这些上面的这些方面。

下面是一个使用Memo的例子。
对于流入的URL的每一个元素我们都会调用Get，
并打印调用延时以及返回的数据大小的log：

m := memo.New(httpGetBody)
for url := range incomingURLs(){
	start := time.Now()
	value, err := m.Get(url)
	if err != nil{
		log.Print(err)
	}
	fmt.Printf("%s, %s， %d bytes\n",
		url, time.Since(start), len(value.([]byte)))
}

我们可以使用测试包（第11章的主题）来系统地鉴定缓存的效果。
从下面的测试输出，我们可以看到URL流包含了一些重复的情况，
尽管我们第一次对每一个URL的 (*Memo).Get的调用都会花上几百毫秒，
但第二次就只需要花1毫秒就可以返回完整的数据了。

测试示例代码

这个测试是顺序地去做所有的调用的。

由于这种彼此独立的HTTP请求可以很好地并发，
我们可以把这个测试改成并发形式。
可以使用sync.WaitGroup来等待所有的请求都完成之后再返回。

m := memo.New(httpGetBody)
var n sync.WaitGroup
for url := range incomingURLs(){
	n.Add(1)
	go func (url string)(){
		start := time.Now()
		value, err := m.Get(url)
		if err != nil{
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
			n.Done()
	}(url)
}
n.Wait()

这次测试跑起来更快了，然而不幸的是貌似这个测试不是每次都能够正常工作。
我们注意到有一些意料之外的cache miss（缓存未命中），
或者命中了缓存但却返回了错误的值，
或者甚至会直接崩溃。

但更糟糕的是，有时候这个程序还是能正确的运行（译注：也就是最让人崩溃的偶发bug），
所以我们甚至可能都不会意识到这个程序有bug。
但是我们可以使用-race这个flag来运行程序，
竞争检测器（9.6）会打印像下面这样的报告：

$ go test -run=TestConcurrent -race -v gopl.io/ch9/memo1
=== RUN TestConcurrent
...
WARNING: DATA RACE
Write by goroutine 36:
runtime.mapassign1()
~/go/src/runtime/hashmap.go:411 +0x0
gopl.io/ch9/memo1.(*Memo).Get()
~/gobook2/src/gopl.io/ch9/memo1/memo.go:32 +0x205
...
Previous write by goroutine 35:
runtime.mapassign1()
~/go/src/runtime/hashmap.go:411 +0x0
gopl.io/ch9/memo1.(*Memo).Get()
~/gobook2/src/gopl.io/ch9/memo1/memo.go:32 +0x205
...
Found 1 data race(s)
FAIL gopl.io/ch9/memo1 2.393s

memo.go的32行出现了两次，
说明有两个goroutine在没有同步干涉的情况下更新了cache map。
这表明Get不是并发安全的，存在数据竞争。

28 func (memo *Memo) Get(key string) (interface{}, error) {
29 res, ok := memo.cache(key)
30 if !ok {
31 res.value, res.err = memo.f(key)
32 memo.cache[key] = res
33 }
34 return res.value, res.err
35 }

最简单的使cache并发安全的方式是使用基于监控的同步。
只要给Memo加上一个mutex，
在Get的一开始获取互斥锁，
return的时候释放锁，就可以让cache的操作发生在临界区内了：

示例代码ch9/memo2

测试依然并发进行，但这会竞争检查器“沉默”了。
不幸的是对于Memo的这一点改变使我们完全丧失了并发的性能优点。
每次对f的调用期间都会持有锁，
Get将本来可以并行运行的I/O操作串行化了。
我们本章的目的是完成一个无锁缓存，
而不是现在这样的将所有请求串行化的函数的缓存。

下一个Get的实现，调用Get的goroutine会两次获取锁：
查找阶段获取一次，如果查找没有返回任何内容，
那么进入更新阶段会再次获取。
在这两次获取锁的中间阶段，
其他goroutine可以随意使用cache。

示例代码ch9/memo3

这些修改使性能再次得到了提升，
但有一些URL被获取了两次。
这种情况咋两个以上的goroutine同一时刻调用Get来请求同样的URL时会发生。
多个goroutine一起查询cache，
发现没有值，
然后一起调用f这个慢不拉叽的函数。
在得到结果后，也都会去更新map。
其中一个获得的结果会覆盖掉另一个的结果。

理想情况下是应该避免掉多余的工作的。
而这种“避免”工作一般被成为duplicate suppression（重复抑制/避免）。
下面版本的Memo每一个map元素都是指向一个条目的指针。
每一个条目包含对函数f调用结果的内容缓存。
与之前不同的是这次entry还包含了一个叫ready的channel。
在条目的结果被设置之后，这个channel就会被关闭，
以向其他goroutine广播（8.9）去读取该条目的结果是安全的了。

示例代码ch9/memo4

现在Get函数包括下面这些步骤了：
获取互斥锁来保护共享变量cache map， 查询map中是否存在指定条目，
如果没有找到那么分配空间插入一个新条目，释放互斥锁。
如果存在条目的话，且其值没有写入完成
（也就是有其他的goroutine在调用f这个慢函数）时，
goroutine必须等待值ready之后才能读到条目的结果。
而想知道是否ready的话，
可以直接从ready channel中读取，
由于这个读取操作在channel关闭之前一直是阻塞。

如果没有条目的话，需要向map中插入一个没有ready的条目，
当前正在调用的goroutine就需要负责调用慢函数、更新条目，
以及向其他所有goroutine关闭条目已经ready可读的消息了。

条目中的e.res.value和e.res.err变量是在多个goroutine之间共享的。
创建条目的goroutine同时也会设置条目的值，
其他goroutine在收到“ready”的广播消息之后立刻回去读取条目的值。
尽管会被多个goroutine同时访问，但却并不需要互斥锁。
ready channel的关闭一定会发生在其他goroutine接收到广播事件之前，
因此第一个goroutine对这些变量的写操作是一定发生在这些读操作之前的。
不会发生数据竞争。

这样并发、不重复、无阻塞的cache就完成了。

上面这样Memo的实现使用了一个互斥量来保护多个goroutine调用Get时的共享map变量。
不妨把这种设计和前面提到的把map变量限制在一个单独的monitor goroutine的方案做一些对比，
后再在调用Get时需要发消息。

Func、result和entry的声明和之前保持一致：

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct{
	value interface{}
	err error
}

type entry struct{
	res result
	ready chan struct{} // closed when res is ready
}

然而Memo类型现在包含了一个叫做requests的channel，
Get的调用方法这个channel来和monitor goroutine来通信。
requests channel中的元素类型是request。
Get的调用方法会把这个结构中的两组key都填充好，
实际上用这两个变量来对函数进行缓存的。
另一个叫response的channel会被拿来发送响应结果。
这个channel只会传回一个单独的值。

示例代码ch9/memo5

上面的Get方法，会创建一个response channel，
把它放进request结构中，然后发送给monitor goroutine，
然后马上又会接受到它。

cache变量被限制在了monitor goroutine(*Memo).server中， 
下面会看到。
monitor会在循环中一直读取请求，
直到request channel被Close方法关闭。
每一个请求都会去查询cache，
如果没有找到条目的话，那么就会创建/插入一个新的条目。

示例代码ch9/memo5

和基于互斥量的版本类似，第一个对某个key的请求需要负责去调用函数f，
并传入这个key，将结果存在条目里，
并关闭ready channel来广播条目的ready消息。
使用 (*entry).call来完成上述工作。

紧接着对同一个key的请求会发现map中已经有了存在的条目，
然后会等待结果变为ready，
并将结果从response发送给客户端的goroutine。
上述工作是用 (*entry).deliver来完成的。
对call和deliver方法的调用必须在自己的goroutine中进行，
以确保monitor goroutine不会因此而被阻塞住，
而没法处理新的请求。

这个例子说明我们可以用上锁，通信（channel？）来建立并发程序都是可行的。
（这个翻译。。。）

9.8 Goroutine和线程（P370）
在上一章中我们说goroutine和操作系统的线程区别可以先忽略。
尽管两者的却别实际上只是一个量的区别，
但量变会引起质变的道理同样适用于goroutine和线程。
现在正是我们来区分开两者的最佳时机。

9.8.1 动态栈
每一个OS（操作系统）线程都有一个固定大小的内存块（一般会是2MB）来做栈，
这个栈会用来存储当前正在被调用或挂起（指在调用其它函数时）的函数的内部变量。
这个固定大小的栈同时很大又很小。
因为2MB的栈对于一个小小的goroutine来说是很大的内存浪费，
比如对于我们用到的，一个只是用来WaitGroup之后关闭channel的goroutine来说。
而对于go程序来说，同时创建成百上千个goroutine是非常普遍的，
如果每一个goroutine都需要这么大的栈的话，
那这么多的goroutine就不太可能了。
除去大小的问题之外，固定大小的栈对于更复杂或者更深层次的递归函数调用来说显然是不够的。
修改固定的大小可以提升空间的利用率允许创建更多的线程，
并且可以允许更深的递归调用，
不过这两者是没法同时兼备的。

相反，一个goroutine会以一个很小的栈开始其生命周期，一般只需要2KB。
一个goroutine的栈，和操作系统线程一样，会保存其活跃或挂起的函数调用的本地变量，
但是和OS线程不太一样的是一个goroutine的栈大小并不是固定的；
栈的大小会根据需要动态地伸缩。
而goroutine的栈的最大值有1GB，比传统的固定大小的线程要大得多，
尽管一般情况下，大多goroutine都不需要这么大的栈。

9.8.2 Goroutine调度
OS线程会被操作系统内核调度。
每几毫秒，一个硬件计时器会中断处理器，
这会调用一个叫做scheduler的内核函数。
这个函数会挂起当前执行的线程并保存内存中他的寄存器内容，
检查线程列表并决定下一次哪个线程可以被运行，
并从内存中恢复giant线程的寄存器信息，
然后恢复执行该线程的现成并开始执行线程。
因为操作系统线程是被内核所调度，所以从一个线程向另一个“移动”需要完整的上下文切换，
也就是说，保存一个用户线程的状态到内存，恢复另一个线程的到寄存器，
然后更新调度器的数据结构。
这几步操作很慢，因为其局部性很差需要几次内存访问，并且会增加运行的cpu周期。

Go的运行时包含了其自己的调度器，
这个调度器使用了一些技术手段，
比如m:n调度，因为其会在n个操作系统线程上多工（调度）m个goroutine。
Go调度器的工作和内核的调度是相似的，
但是这个调度器只关注单独的Go程序中的goroutine（译注：按程序独立）。
和操作系统的线程调度不同的是，
Go调度器并不是用一个硬件定时器而是被Go语言“建筑”本身进行调度的。
例如当一个goroutine调用了time.Sleep或者被channel调用或者mutex操作阻塞时，
调度器会使其进入休眠并开始执行另一个goroutine直到实际到了再去唤醒第一个goroutine。
因为这种调度方式不需要进入内核的上下文，
所以重新调度一个goroutine比调度一个线程代价要低得多。

9.8.3 GOMAXPROCS（Go最大进程）

Go的调度器使用了一个叫做GOMAXPROCS的变量来决定会有多少个操作系统的线程同时执行Go的代码。
其默认的值是运行机器上的CPU的核心数，
所以在一个有8个核心的机器上时，
调度器一次会在8个OS线程上去调度GO代码。
（GOMAXPROCES是前面说的m:n调度中的n）
在休眠中的或者在通信中被阻塞的goroutine是不需要一个对应的线程来做调度的。
在I/O中或系统调用中或调用非Go语言函数时，是需要一个对应的操作系统线程的，
但是GOMAXPROCS并不需要将这集中情况计数在内。

你可以用GOMAXPROCS的环境变量吕显示地控制这个参数，
或者也可以在运行时用runtime.GOMAXPROCS函数来修改它。
我们在下面的小程序中会看到GOMAXPROCS的效果，
这个程序会无限打印0和1.package main

for{
	go fmt.Print(0)
	fmt.Print(1)
}

$ GOMAXPROCS = 1 go run backer-cliche.go
111111111111111111110000000000000000000011111...

$ GOMAXPROCS=2 go run hacker-cliché.go
010101010101010101011001100101011010010100110...

在第一个执行时，最多同时只能有一个goroutine被执行。
厨师情况下只有main goroutine被执行，所以会打印很多1，
过了一段时间后，Go调度器会将其置为休眠，并唤醒另一个goroutine，
这时候就开始打印多个0了，
在打印的时候，goroutine是被调度到操作系统线程上的。
在第二次执行时，我们使用了两个操作系统线程，
所以两个goroutine的调度是受很多因子影响的，
而runtime也是在不断地发展演进的，
所以这里的你实际得到的结果可能会因为版本的不同而与我们运行的结果有所不同。

9.8.4 Goroutine没有ID号（p372）
在大多数支持多线程的操作系统和程序语言中，
当前的线程都有一个独特的身份（id），
并且这个身份信息可以以一个普通值的形式被很容易地获取到，
典型的可以是一个integer或者指针值。
这种情况下我们做一个抽象化的thread-local storage（线程本地存储，
多线程编程中不希望其他线程访问的内容）
就很容易，
只需要以线程的id作为key的一个map就可以解决问题，
每一个线程以其id就能从中获取到值，且和其他线程互不冲突。

goroutine没有可以被程序员获取到的身份（id）的概念。
这一点是设计上故意而为之，
由于thread-local storage总是会被滥用。
比如说，一个web server是用一种支持tls的语言实现的，
而非常普遍的是很多函数会去寻找HTTP请求的信息，
这代表它们就是去其存储层（这个存储层可能是tls）查找的。
这就像是那些过分依赖全局变量的程序一样，会导致一种非健康的“距离外行为”，
在这种行为下，一个函数的行为可能不是由其自己内部的便令锁决定，
而是由其所运行在的线程所决定。
因此，如果线程本身的身份会改变————比如一些worker线程之类的————
那么函数的行为就会变得神秘莫测。

Go鼓励更为简单的模式，
这种模式下参数对函数的影响都是显式的。
这样不仅使程序变得更易读，
而且会让我们自由地向一些给定的函数分配子任务时不用担心其身份信息影响行为。

你现在应该已经明白了写一个Go程序所需要的所有语言特性信息。
在后面两章节中，我们会回顾一些之前的实例和工具，
支持我们写出更大规模的程序：如何将一个工程组织成一系列的包，
如何获取，构建，测试，性能测试，剖析，写文档，并且将这些包分享出去。

