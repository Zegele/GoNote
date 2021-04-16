第六章 方法 
(p209)

面向对象编程（OOP）就成为了称霸工程界和教育界的编程范式，
所以之后几乎所有大规模被应用的语言都包含了对OOP的支持，go语言也不例外。

尽管没有被大众所接收的明确的OOP的定义，从我们的理解来讲，一个对象其实也就是一个简单的值或一个变量，在这个对象中会包含一些方法，
而一个方法则是一个一个和特殊类型关联的函数。
一个面向对象的程序会用方法来表达期属性和对应的操作，这样使用这个对象的用户就不需要直接去操作对象，而是借助方法来做这些事情。

在早些的章节中，我们已经使用了标准库提供的一些方法，比如time.Duration这个类型的Seconds方法：
conse day = 24 * time.Hour
fmt.Println(day.Seconds()) // "86400"

并且在2.5 节中，我们定义了一个自己的方法，Celsius类型的String方法：
func (c Celsius) String() string { return fmt.Sprintf("%gC°", c)}

在本章中，OOP编程的第一方面，我们会向你展示如何有效地定义和使用方法。
我们会覆盖到OOP编程的两个关键点，封装和组合。

6.1 方法和声明（p210）

1. 在函数声明时，在其名字之前放上一个变量，即是一个方法。
这个附加的参数会将该函数附加到这种类型上，即相当于为这种该类型定义了一个独占的方法。

下面来写我么第一个方法的例子，这个例子在package geometry下：

示例代码

上面的代码里那个附加的参数p，叫做方法的接收器（receiver），
早期的面向对象语言留下的遗产将调用一个方法称为“向一个对象发送消息”。

在Go语言中，我们并不会像其他语言那样用this或self作为接收器；
我们可以任意的选择接收器的名字。
由于接收器的名字经常会被使用到，所以保持其在方法间传递时的一致性和简短性是不错的主意。
这里的建议是可以使用期类型的第一个字母，比如这里使用了Point的首字母p。

在方法调用过程中，接收器参数一般会在方法名之前出现。
这和方法声明是一样的，都是接收器参数在方法名字之前。
下面是例子：
p := Point{1, 2}
q := Point{4, 8}
fmt.Println(Distance(p, q)) // "5", function call函数调用
fmt.Println(p.Distance(q)) // "5", method call方法调用

可以看到，上面的两个函数调用都是Distance，但是却没有发生冲突。
第一个Distance的调用实际上用的是包级别的函数geometry.Distance,
而第二个则是使用刚刚声明的Point，调用的是Point类下声明的Point.Distance方法。

这种p.Distance的表达式叫做选择器，因为他会选择合适的对应p这个对象Distance方法来执行。
选择器也会被用来选择一个struct类型的字段，比如p.X。
由于方法和字段都是在同一命名空间，所以如果我们在这里声明一个X方法的话，编译器会报错，
因为在调用p.X时会有歧义。

因为每种类型都有期方法的命名空间，我们在用Distance这个名字的时候，不同的Distance调用指向了不同类型里的Distance方法。
让我们来定义一个Path类型，这个Path代表一个线段的集合，并且也给这个Path定义一个叫Distance的方法。

Path是一个命名的slice类型，而不是Point那样的struce类型，然而我们依然可以为它定义方法。
在能够给任意类型定义方法这一点上，Go和很多其他的面向对象的语言不太一样。
因此在Go语言里，我们一些简单的数值、字符串、slice、map来定义的一些附加行为很方便。


2. 方法可以被声明到任意类型，只要不是一个指针，或者一个interface。

两个Distance方法有不同的类型。他们两个方法之间没有任何关系，尽管Path的Distance方法会在内部调用Point.Distance方法来计算每个连接邻接点的线段的长度。

示例代码：
perim ：= Path{
	{1, 1},
	{5, 1},
	{5, 4},
	{1, 1}
}

fmt.Println(perim.Distance()) // "12" 注意小括号不要忘了

上面两个对Distance名字的方法的调用中，编译器会根据方法的名字以及接收器来决定具体调用的是哪一个函数。
例子中path[i-1]数组中的类型是Point，因此Point.Distance这个方法被调用；
例子中的perim类型是Path，因此Distance调用的是Path.Distance。

对于一个给定的类型，其内部的方法都必须有唯一的方法名，但是不同的类型却可以有同样的方法名，比如我们这里Point和Path就都有Distance这个名字的方法；
所以我们没有必要非在方法名前加类型名来消除歧义，比如PathDistance。


3. 这里我们已经看到了方法比之前函数的一些好处：方法名可以简短。
当我们在包外调用的时候这种好处就会被放大，因为我们可以使用这个短名字，而可以省略掉包的名字，
下面是例子。

import ".../geometry"

perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
fmt.Println(geometry.Path.Distance(perim)) // "12", standalone function
fmt.Println(perim.Distance()) // "12", method of geometry.Path

如果我们要用方法去计算perim的distance，还需要去写全geometry的包名，和其函数名，但是因为Path这个变量定义了一个可以直接用的Distance方法，
所以我们可以直接写perim.Distance()。
相当于可以少打很多字。
因为在Go里包外调用函数需要带上包名，还是挺麻烦的。

6.2 基于指针对象的方法（p213）
4. 当调用一个函数时，会对其每一个参数值进行拷贝，如果一个而函数需要更新一个变量，
或者函数的其中一个参数实在太大我们希望能够避免进行这种默认的拷贝，
这种情况下我们就需要用到指针了。

5. 对应到我们这里用来更新接收器的对象的方法，
当这个接受者变量本身比较大时，我们就可以用其指针而不是对象来声明方法。
如下：
func (p *Point) ScaleBy(factor float64){
	p.X *= factor // p.X = p.X * factor
	p.Y *= factor
}

这个方法的名字是 (*Point).ScaleBy 。
这里的括号是必须的；没有括号的话这个表达式可能被理解为*（Point.SceleBy） 。

6. 在现实的程序里，一般会约定如果Point这个类有一个指针作为接收器的方法，
那么所有Point的方法都必须有一个指针接收器，即使是那些并不需要这个指针接收器的函数。
我们在这里打破了这个约定只是为了展示一下两种方法的异同而已。

只有类型（Point）和指向他们的指针（*Point），才是可能会出现在接收器声明里的两种接收器。
此外，为了避免歧义，在声明方法时，如果一个类型名本身是一个指针的话，是不允许其出现在接收器中，
如下面的例子：
type P *int
func (P) f(){/*...*/} // compile error: invalid reveiver type

想要调用指针类型方法 (*Point).ScaleBy ，只要提供一个Point类型的指针即可，像下面这样。

r := &Point{1, 2}
r.ScaleBy(2)
fmt.Println(*r) // "{2, 4}"
//指针的方法是这样使用的。
或者这样：
p := Point{1, 2}
pptr := &p
pptr := pptr.ScaleBy(2)
fmt.Println(p) // "{2, 4}"
或者这样：
p := Point{1, 2}
(&p).ScaleBy(2)
fmt.Println(p) // "{2, 4}"

不过后面的两种方法有些笨拙。幸运的是，go语言本身在这种地方会帮到我们。
如果接收器p是一个Point类型的变量，并且其方法需要一个Point指针作为接收器，
我们可以用下面这种简短的写法：

p.ScaleBy(2)

编译器会隐式地帮我们用&p去调用ScaleBy这个方法。这种简写方法只适用于“变量”，
包括struct里的字段比如p.X，以及array和slice内的元素比如perim[0]。
我们不能通过一个无法渠道地址的接收器来调用指针方法，比如临时变量的内存地址就无法获取得到：
Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point leteral

7. 但是我们可以用一个 *Point 这样的接收器来调用Point的方法，因为我们可以通过地址来找到这个变量，只要用解引用符号 *  来取到该变量即可。
编译器在这里也会给我们隐式地插入 * 这个操作符，所以下面两种写法等价的：
pptr.Distance(q)
(*pptr).Distance(q)

或者接收器形参是类型T，但接收器实参是类型*T，这种情况下编译器会隐式地为我们取变量的地址：
p.ScaleBy(2) // implicit (&p)
//这里需要*Point类型，但p是Point类型，Go自动隐式的取了变量的地址。

或者接收器形参是类型*T，实参是类型T。编译器会隐式地为我们调用解引用，取到指针指向的实际变量：
pptr.Distance(q) // implicit (*pptr)
//pptr是*Point类型，这里需要Point类型。
//自动隐式地找到了pptr指向的值，并引用。

如果类型T的所有方法都是用T类型自己来做接收器（而不是*T），那么拷贝这种类型的示例就是安全的；
调用他的任何一个方法也就会产生一个值的拷贝。
比如time.Duration的这个类型，在调用其他方法时就会被全部拷贝一份，包括在作为参数传入函数的时候。
但是如果一个方法使用指针作为接收器，你需要避免对其进行拷贝，因为这样可能会破坏掉该类型内部的不变性。
比如你对bytes.Buffer对象进行了拷贝，那么可能会引起原始对象和拷贝对象只是别名而已，
但实际上其指向的对象是一致的。
紧接着对拷贝后的变量进行修改可能会有让你意外的结果。
（译注：作者这里说的比较绕，
在声明一个method的receiver该是指针还是非指针类型时，你需要考虑两方面，
第一方面是这个对象本身是不是特别大，如果声明为非指针变量时，调用会产生一次拷贝；
第二方面是如果你用指针类型作为receiver，那么你一定要注意，这种指针类型指向的始终是一块内存地址。
）

6.2.1 Nil也是一个合法的接收器类型
8. 就像一些函数允许nil指针作为参数一样，方法 理论上也可以用nil作为其接收器，尤其当nil对于对象来说是合法的零值时，
比如map或者slice。
在下面的简单int链表的例子里，nil代表的是空链表：
// An IntList is a linked list of integers.
// A nil *IntList represents the empty list.
type IntList struct{
	Value int
	Tail *IntList
}

// Sum returns the sum of the list elements.
func (list *IntList) Sum() int {
	if list == nil{
		return 0
	}
	return list.Value + list.Tail.Sum()
}

当你定义一个允许nil作为接收器值的方法的类型时，在类型前面的注释中指出nil变量代表的意义是很有必要的，就像我们上面例子里做的这样。

下面是net/url包里Values类型定义的一部分。
package url

// Values maps a string key to a list of values.
type Values map[string][]string
// Get returns the first value associated with the given key.
// or "" if there are none.
func (v Values) Get(key string) string{
	if vs := v[key]; len(vs) > 0{
		return vs[0]
	}
	return ""
}

// Add adds the value to key.
// It appends to any existing values associated with key.
func (v Values) Add(key, value string){
	v[key] = append(v[key], value)
}

这个定义向外暴露了一个map的类型的变量，并且提供了一些能够简单操作这个map的方法。
这个map的value字段是一个string的slice，所以这个Values是一个多维map。
客户端使用这个变量的时候可以使用map固有的一些操作（make，切片，m[key]等等），也可以使用这里提供的操作方法，或者两者并用，都是可以的：

示例代码：

9. 对Get的最后一次调用中，nil接收器的行为即是一个空map的行为。
我们可以等价地将这个操作写成 Value(nil).Get("item"),但是如果你直接写nil.Get("item")的话是无法通过编译的，
因为nil的字面量编译器无法判断其准备类型。
所以相比之下，最后的那行m.Add的调用就会产生一个panic，因为他尝试更新一个空map。

由于url.Values是一个map类型，并且间接引用了其key/value对，
因此url.Values.Add对这个map里的元素做任何的更新、删除操作对调用方都是可见的。
实际上，就像在普通函数中一样，虽然可以通过引用来操作内部值，
但在方法想要修改引用本身是不会影响到原始值的，比如把他置为nil，或者让这个引用指向了其他的对象，
调用方都不会受影响。
（译注：因为传入的是存储了内存地址的变量，你改变这个变量是影响不了原始的变量的。）
理解：储存了内存地址的变量。 如果内存地址的变量是10，传入是储存了10，内存地址是变了的。


6.3 通过嵌入结构体来扩展类型（p218）
来看看ColoredPoint这个类型：
import "image/color"

type Point struct{X, Y float64}

type ColoredPoint struct{
	Point
	Color color.RGBA
}

10. 我们完全可以将ColoredPoint定义为一个有三个字段的struct，但是我们却将Point这个类型嵌入到ColorPoint来提供X和Y这两个字段。
像我们4.4节中看到的那样，内嵌可以使我们在定义ColorPoint时得到一种句法上的简写形式，并使其包含Point类型所具有的一切字段，然后在定义一些自己的。
如果我们想要的话，我们可以直接认为通过嵌入的字段就是ColoredPoint自身的字段，而完全不需要在调用时指出Point，比如下面这样。
var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y) // "2"

对于Point中的方法我们也有类似的用法，我们可以把ColoredPoint类型当作接收器来调用Point里的方法，即使ColoredPoint里没有声明这些方法;
red := color.RGBA{255, 0, 0, 255}
blue := color.RGBA{0, 0, 255, 255}
var p = ColoredPoint{Point{1, 1}, red}
var q = ColoredPoint{Point{5, 4}, blue}
fmt.Println(p.Distance(q.Point)) // "5"
//p.Distance(q.Point)
//这里p是ColoredPoint类型，但ColoredPoint类型里有Point类型所以这样直接使用了。
p.ScaleBy(2)
q.ScaleBy(2)
fmt.Println(p.Distance(q.Point)) // "10"

Point类的方法也被引入了ColoredPoint。
用这种方式，内嵌可以使我们定义字段特别多的复杂类型，
我们可以将字段先按小类型分组，然后定义小类型的方法，
之后再把他们组合起来。


读者如果对基于类来实现面向对象的语言比较熟悉的话，可能会倾向于将Point看作一个基类，
而ColoredPoint看作其自雷或者继承类，或者将ColoredPoint看作“is a” Point类型。
但这是错误的理解。
请注意上面例子中Distance方法的调用。
Distance有一个参数是Point类型，但q并不是一个Point类，所以尽管q有着Point这个内嵌类型，我们也必须要显式地选择它。
尝试直接传q的话你会看到下面这样的错误：
p.Distance(q) // compile error : cannot ues q (CorloredPoint) as Point

一个ColoredPoint并不是一个Point，但他“has a” Point，并且它有从Point类里引入的Distance和ScaleBy方法。
如果你喜欢从现实的角度来考虑问题，内嵌字段会指导编译器去生成额外的包装方法来委托已经声明好的方法，和下面的形式是等价的：
func (p ColoredPoint) Distance（q Point) float 64{
	return p.Point.Distance(q)
}

func (p *ColoredPoint) ScaleBy(factor float64){
	p.Point.ScaleBy(factor)
}

11. 当Point.Distance被第一个包装方法调用时，他的接收器是p.Point，而不是p。
当然了，在Point类的方法里，你说访问不到ColoredPoint的任何字段的。

12. 在类型中内嵌的匿名字段也可能是一个命名类型的指针，这种情况下字段和方法会被间接地引入到当前的类型中（译注：访问需要通过该指针真相的对象去取）。
添加这一层间接关系让我们可以共享通用的结构并动态地改变对象之间的关系。
下面这个ColoredPoint的声明内嵌了一个*Point指针。
type ColoredPoint struct{
	*Point
	Color color.RGBA
}

p:= ColoredPoint(&Point{1,1}, red)
q:= ColoredPoint(&Point{1,1}, blue)
fmt.Println(p.Distance(*q.Point)) //"5"
q.Point = p.Point
p.ScaleBy(2)
fmt.Println(*p.Point, *q.Point) // "{2, 2} {2, 2}"

13. 一个struct类型也可能会有多个匿名字段。我们将ColoredPoint定义为下面这样：
type ColoredPoint struct{
	Point
	color.RGBA
}

然后这种类型的值便会拥有Point和RGBA类型的所有方法，以及直接定义在ColoredPoint中的方法。
当编译器解析一个选择器到方法时，比如p.ScaleBy,他会首先去找直接定义在这个类型里的ScaleBy方法，然后找被ColoredPoint的内嵌字段们引入的方法，
然后去找Point和RGBA的内嵌字段引入的方法，然后一直递归向下找。
如果选择器有二义性的话编译器会报错，比如你在同一级里有两个同名的方法。

方法只能在命名类型（像Point）或者指向类型的指针上定义，但是多亏了内嵌，有些时候我们给匿名struct类型来定义方法也有了手段。

下面是一个小trick。这个例子展示了简单的cache，其使用两个包级别的变量来实现，一个mutex互斥量（9.2）和它所操作的cache：
var (
	mu sync.Mutex // guards mapping
	mapping = make(map[string]string)
)

func Lookup(key string) string{
	mu.Lock()
	v := mapping[key]
	mu.Unlock()
	return v
}

下面这个版本在功能上是一致的，但将两个包级别的变量放在cache这个struct一组内：
var cache = struct{
	sync.Mutex
	mapping map[string]string
}{
	mapping:make(map[string]string),
}

func Lookup(key string) string{
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}

我们给新的变量起了一个更具表达性的名字：cache。
因为sync.Mutex字段也被嵌入到了这个struct里，其Lock和Unlock方法也就都被引入到了这个匿名结构中了，
这让我们能够以一个简单明了的语法来对其进行枷锁解锁操作。

6.4 方法值和方法表达式 （p222）
14. 我们经常选择一个方法，并且在同一个表达式里执行，
比如常见的的p.Distance()形式，实际上将其分为两步来执行也是可能的。
p.Distance叫作选择器，选择器会返回一个方法“值” ->一个将方法（Point.Distance）绑定到特定接收器变量的函数。
这个函数可以不通过制定其接收器即可被调用；
即调用时不需要制定接收器（译注：因为已经在前文中制定过了），只要传入函数的参数即可：

p := Point{1, 2}
q := Point{4, 6}

distanceFromP := p.Distance // method value
fmt.Println(distanceFromP(q)) // "5"

var origin Point // {0, 0}
fmt.Println(distanceFromP(origin)) // "2.23..."

scaleP := p.ScaleBy // method value
scaleP(2) // p becomes (2, 4)
scaleP(3) // then (6, 12)
scaleP(10) // then (60, 120)

在一个包的API需要一个函数值、且调用方希望操作的是某一个绑定了对象的方法的话，
方法“值”会非常实用。
举例来说，下面例子中的time.AfterFunc这个函数的功能是在制定的延迟时间之后来执行一个另外的函数。
且这个函数操作的是一个Rocket对象r
type Rocket struct {/*...*/}
func (r *Rocket) Launch(){/*...*/}
r := new(Rocket)
time.AfterFunc(10 * time.Second, func(){r.Launch()})

直接用方法“值”传入AfterFunc的话可以更为简短：
time.AfterFunc(10 * time.Second, r.Launch)

译注：省掉了上面那个例子里的匿名函数。

15. 和方法“值”相关的还有方法表达式。当调用一个方法时，与调用给一个普通的函数相比，
我们必须要用选择器（p.Distance）语法类制定方法的接收器。

当T是一个类型时，方法表达式可能会写作T.f或者*T.f ，会返回一个函数“值”，
这种函数会将其第一个参数用作接收器，所以可以用通常（译注：不写选择器）的方式来对其进行调用：

p := Point{1, 2}
q := Point{4, 6}

distance := Point.Distance // method expression 方法表达式
fmt.Println(distance(p, q)) // "5"
fmt.Println("%T\n", distance) // "func(Point, Point) float64"

scale := (*Point).ScaleBy
scale(&p, 2)
fmt.Println(p) // "{2 4}"
fmt.Printf("%T\n", scale) // "func(*Point, float64)"

当你根据一个变量来决定调用同一个类型的哪个函数时，方法表达式就显得很有用了。
你可以根据选择来调用接收器各不相同的方法。
下面的例子，变量op代表Point类型的addition或者subtraction方法，
Path.TranslateBy方法会为其Path数组中的每一个Point来调用对应的方法：

type Point struct{ X, Y float64}

func (p Point) Add(q Point) Point{return Point{p.X + q.X, p.Y + q.Y}}
func (p Point) Sub(q Point) Point{return Point{p.X - q.X, p.Y - q.Y}}

type Path []Point

func (path Path) TranslateBy(offset Point, add bool){
	var op func(p, q Point) Point
	if add { // 这个怎么理解，如果add 是true则继续，否则else
		op = Point.Add
	}else{
		op = Point.Sub
	}
	for i := range path{
		// Call either path[i].Add(offset) or path[i].Sub(offset).
		path[i] = op(path[i], offset)
	}
}

6.5 示例:Bit数组（p224）
16. Go语言里的集合一般会用map[T]bool这种形式来表示，T代表元素类型。
集合用map类型来表示虽然非常灵活，但我们可以以一种更好的形式来表示它。
例如在数据流分析领域，集合元素同城是一个非负整数，集合会包含很多元素，并且结合会经常并集，交集操作，
这种情况下，bit数组会比map表现更加理想。
（译注：这里在补充下一个例子，比如我们执行一个http下载任务，把文件按照16kb一块划分为很多块，需要有一个全局变量来标识哪些块下载完成了，这种时候也需要用到bit数组）

一个big数组通常会用一个无符号数或称之为“字”的slice来表示，每一个元素的每一位都表示集合里的一个值。
当集合的第i位被设置时，我们才说这个集合包含元素i。
下面的这个程序展示了一个简单的bit数组类型，并且实现了三个函数来对这个bit数组来进行操作：

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

17. |= （参照第3章笔记，第80行）
a |= b
a = a|b //等价于上面的公示


18. 因为每一个字都有64个二进制位，所以为了定位x的bit位，我们用了x/64的商作为字的下标，
并且用x%64得到的值作为这个字内的bit的所在位置。
UnionWith这个方法里用到了bit位的“或”逻辑操作符号 | 来一次完成64个元素的 或计算。

当前这个实现还缺少了很多必要的特性，我们把其中一些作为练习题列在本小节之后。
但是有一个方法如果缺失的话我们的bit数组可能会比较难混：
将IntSet作为一个字符串来打印。
这个我们来实现它，让我们来给上面的例子添加一个String方法，类似2.5节中做的那样：

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte('}')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

这里留意一下String方法，是不是和3.5.4节中的intsToString方法很相似；
bytes.Buffer在String方法里经常这么用。
当你为一个复杂的类型定义了一个String方法时，fmt包就会特殊对待这种类型的值，
这样可以让这些类型在打印的时候看起来更加友好，而不是直接打印其原始的值。
fmt会直接调用用户定义的String方法。
这种机制依赖于接口和类型断言，在第7章中我们会详细介绍。
现在我们就可以在实战中直接用上名定义好的IntSet了：

var x, y IntSet

x.Add(1)
x.Add(144)
x.Add(9)
fmt.Println(x.String()) // "{1 9 144}"

y.Add(9)
y.Add(42)
fmt.Println(y.String()) // "{9 42}"

x.UnionWith(&y)
fmt.Println(x.String()) // "{1 9 42 144}"
fmt.Println(x.Has(9), x.Has(123)) // "true false"

19. 这里要注意：我们声明的String和Has两个方法都是以指针类型*IntSet来作为接收器的，
但实际上对于这两个类型来说，把接收器声明为指针类型也没什么必要。
不过另外两个函数就不是这样了，因为另外两个函数操作的是s.words对象，如果你不把接收器声明为指针对象，
那么实际操作的是拷贝对象，而不是原来的那个对象。
因此，因为我们的String方法定义在IntSet指针上，
所以当我们的变量是IntSet类型而不是IntSet指针时，可能会有下面这样让人意外的情况：

fmt.Println(&x) //"{1 9 42 144}"
fmt.Println(x.String()) // "{1 9 42 144}"
fmt.Println(x) //"{[4398046511618 0 65536]}"//如果经过String函数后，fmt会优先调用String后的值（查看关于String的函数）

在第一个Println中，我们打印一个*IntSet的指针，这个类型的指针确实有自定义的String方法。
第二Println，我们直接调用了x变量的 String()方法；
这种情况下编译器会隐式地在x前插入&操作符，这样相当于我们还是调用的IntSet指针的String方法。
在第三个Println中，因为IntSet类型没有String方法，所以Println方法会直接以原始的方式理解并打印。
所以在这种情况下&符号是不能忘的。在我们这种场景下，你把String方法绑定到IntSet对象上，而不是IntSet指针上可能会更合适一些，不过这也需要具体问题具体分析。


6.6 封装 （p228）
20. 一个对象的变量或者方法如果对调用方是不可见的话，一般就被定义为“封装”。
封装有时候也被叫做信息隐藏，同时也是面对对象编程最关键的一个方面。

Go语言只有一种控制可见性的手段：大写字母的标识符会从定义他们的包中被导出，
小写字母的则不会。
这种限制包内成员的方式同样适用于struct或者一个类型的方法。
因而如果我们想要封装一个对象，我们必须将其定义为一个struct。

这也就是前面的小节中IntSet被定义为struct类型的原因，尽管它只有一个字段：
type IntSet struct{
	words []uint64
}

当然，我们也可以把IntSet定义为一个slice类型，尽管这样我们就需要把代码中所有方法里用到的s.words用*s替换掉了：
type IntSet []uint64

尽管这个版本的IntSet在本质上是一样的，他也可以允许其他包中可以直接读取并编辑这个slice。
换句话说，相对*s这个表达式会出现在所有的包中，s.words只需要定义IntSet的包中出现（译注：所以还是推荐struct类型的IntSet）

21. 这种基于名字的手段使得在语言中最小的封装单元是Package， 而不是想其他语言一样的类型。
一个struct类型的字段对同一个包的所有代码都有可见性，无论你的代码是写在一个函数还是方法里。

22. 封装提供了三方面的优点。
首先，因为调用方不能直接修改对象的变量值，其只需要关注少量的语句并且只要弄懂少量的可能的值即可。

第二，隐藏实现的细节，可以防止调用方法依赖那些可能变化的具体实现，这样使设计包的程序员在不破坏对外的api情况下能得到更大的自由。
把byte.Buffer这个类型作为例子来考虑。
这个类型在做短字符串叠加的时候很常用，所以在涉及的时候可以做一些预先的优化，
比如提前预留一部分空间，来避免反复的内存分配。
又因为Buffer是一个struct类型，这些额外的空间可以用附加的字节数组来保存，且放在一个小写字母开头的字段中。
这样在外部的调用方只能看到性能的提升，但并不会得到这个附加变量。
Buffer和其增长算法我们列在这里，为了简洁性稍微做了一些精简：

type Buffer struct{
	buf []byte
	initial [64]byte
	/*...*/
}

// Grow expands the buffer's capacity, if necessary,
// to guarantee space for another n bytes. [...]
func (b *Buffer) Grow(n int){
	if b.buf == nil{
		b.buf = b.initial[:0] // use preallocated space initially
	}
	if len(b.buf)+n>cap(b.buf){
		buf := make([]byte, b.Len(), 2*cap(b.buf)+n)
		copy(buf, b.buf)
		b.buf = buf
	}
}

第三，是阻止了外部调用方法对对象内部的值任意地进行修改。
因为对象内部变量只可以被同一个包内的函数修改，所以包的作者可以让这些函数确保对象内部的一些值的不变性。
比如下面的Counter类型允许调用方来增加counter变量的值，并且允许将这个值reset为0，但是不允许随便设置这个值（压根就访问不到）：

type Counter struct{n int}
func (c *Counter) N() int {return c.n}
func (c *Counter) Increment() {c.n++}
func (c *Counter) Reset() {c.n = 0}

只用来访问或修改内部变量的函数被成为setter或者getter，例子如下，比如log包里的Logger类型对应的一些函数。
在命名一个getter方法时，我们通常会省略掉前面的Get前缀。
这种简洁上的偏好也可以推广到各种类型的前缀比如Fetch，Find或者Lookup。

package log 
type Logger struct{
	flags int
	prefix string
	// ...
}

func (l *Logger) Flags() int
func (l *Logger) SetFlags(flag int)
func (l *Logger) Prefix() string
func (l *Logger) SetPrefix(prefix string)

Go的编码风格不禁止直接导出字段。
当然，一旦进行了导出，就没有办法保证API兼容的情况下去除对其的导出，
所以在一开始的选择一定要经过深思熟虑并且要考虑到包内部的一些不变量的保证，
未来可能的变化，以及调用方的代码质量是否会因为包的一点修改而变差。

封装并不总是理想的。
虽然封装在有些情况是必要的，但有时候我们也需要暴露一些内部内容，
比如：time.Duration（持续时间）将其表现暴露为一个int64数字的纳秒，使得我们可以用一般的数值操作来对时间进行对比，
甚至可以定义这种类型的常量：
const day = 24 * time.Hour
fmt.Println(day.Seconds()) // "86400"

另一个例子，将IntSet和本章开头的geometry.Path进行对比。
Path被定义为一个slice类型，
这允许其调用slice的字面方法来对其内部的points用range进行迭代遍历；
在这一点上，IntSet是没有办法让你这么做的。

这两种类型决定性的不同：
geometry.Path的本章是一个坐标点的序列，不多也不少，我们可以预见到之后也并不会给他增加额外的字段，
所以在geometry包中将Path暴露为一个slice。
相比之下，IntSet仅仅是在这里用了一个[]uint64的slice。
这个类型还可以用[]uint类型来表示，或者我们甚至可以用其他完全不同的占用更小内存的东西来表示这个集合，
所以我们可能还会需要额外的字段来在这个类型中记录元素的个数。
也正是因为这些原因，我们让IntSet对调用方透明。

在这章中，我们学到了如何将方法与命名类型进行组合，并且知道了如何调用这些方法。
尽管方法对于OOP编程来说至关重要，但他们只是OOP编程里的半边天。
为了完成OOP，我们还需要接口。
Go里的接口会在下一章中介绍。