第四章 符合数据类型
1. 这章主要讨论四中类型————数组、slice、map和结构体，
最后会演示如何使用结构体来解码和编码到对应JSON格式的数据，
并且通过结合使用模板来生成HTML页面。

2. 数组和结构体是聚合类型；后半句没懂
数组是有同构的元素组成————每个元素都是完全相同的类型。
结构体则是由异构的元素组成的。
数组和结构体都是由固定内存大小的数据结构。
slice和map则是动态的数据结构，他们将根据需要动态增长。

4.1 数组

3. 数组是一个由固定长度的特定类型元素组成的序列。（p120）
一个数组可以由零个或多个元素组成。
因为数组的长度是固定的，因此在Go语言中很少直接使用数组。

4. Slice，它是可以增长和收缩动态序列，slice功能也更灵活。

5. 数则的每个元素可以通过索引下标来访问，索引下标的范围是从0开始到数组长度减1 的位置。
内置的len函数将返回数组中元素的个数。
var a [3]int // array of integers
fmt.Println(a[0]) // print the first element
fmt.Println(a[len(a)-1]) // print the last element, a[2]

// Print the indices and elements.
for i, v := range a {
	fmt.Printf("%d %d\n", i, v)
}

//Print the elements only.
for _, v := range a{
	fmt.Printf("%d\n", v)
}

6. 默认情况下，数组的每个元素都被初始化为算数类型对应的零值。对于数字类型来说就是0。
我们也可以使用数组字面语法用一组值来初始化数组：
var q [3]int = [3]int{1, 2, 3}
var r [3]int = [3]int{1, 2}
fmt.Println(r[2]) // "0"

7. 数组长度是数组类型的一个组成部分，因此[3]int 和[4]int 是两种不同的数组类型。
数组的长度必须是常量表达式，因为数组的长度需要在编译阶段确定。
q := [3]int{1, 2, 3}
q = [4]int{1, 2, 3, 4} // compile error: cannot assign [4]int to [3]int

8. 也可以制定一个索引和对应值列表的方式初始化：
type Currency int

const(
	USD Currency = iota //美元
	EUR
	GBP
	RMB
) 
symbol := [...]string{USD: "$", EUR:"€", GBP:"￡", RMB:"￥"}

fmt.Println(RMB, symbol[RMB], symbol[3]) //"3 ￥ ￥"

在这种形式的数组字面形式中，初始化索引的顺序是无关紧要的，
而且没用到的索引可以省略，和前面提到的规则一样，未制定初始值的元素将用零值初始化。

r := [...]int{99: -1} //表示第99位的值是-1，前面的0-98都是初始化0值。
定义了一个含有100个元素的数组r，最后一个元素被初始化为-1，其他元素都是用0初始化。

9. 如果一个数组的元素类型是可以相互比较的，那么数组类型也是可以相互比较的，
这时候我们可以直接通过==比较运算符来比较两个数组，
只有当两个数组的所有元素都是相等的时候数组才是相等的。
不相等比较运算符 != 遵循同样的规则。
a := [2]int{1, 2}
b := [...]int{1, 2}
c := [2]int{1, 3}
fmt.Println(a == b, a == c, b == c) // "true false false"
d := [3]int{1, 2}
fmt.Println(a == d) // compile error: cannot compare [2]int == [3]int

10. cypto/sha256包的Sum256函数对一个任意的字节slice类型的数据生成一个对应的消息摘要。（p121）
消息摘要有256bit大小，因此对应[32]byte 数组类型。
如果两个消息摘要是相同的，那么慢可以认为两个消息本身也是相同。（注：理论上有HASH码碰撞的情况，但实际应用可以基本忽略）
如果信息摘要不同，那么消息本身必然不同。

sha256示例代码中,Peintf函数的%x副词参数，它用于指定以十六进制的格式打印数组或slice全部的元素。
%t副词参数是用于打印布尔型数据，
%T副词参数是用于显示一个值对应的数据类型。


11.Go语言中，可以传入一个数组指针，这样函数通过指针对数组的任何修改都可以直接反馈到调用者。（最好先看看上一段文字，理解函数调用机制。） （p122）
func zero (ptr *[32]byte){ // ptr是[32]byte指针类型
	for i := range ptr {
		ptr[i] = 0
	}
}

12. 其实数组字面值[32]byte{}就可以生成一个32字节的数组。
而且每个数组的元素都是零值初始化，也就是0。
因此，我们可以将上面的的zero函数写的更简洁一点：
func zero(ptr *[32]byte){
	*ptr = [32]byte{}//直接把0值赋值给ptr指向内存的值。（注意：*在类型前，表示类型，*在变量前，表示值）
}

13. 数组的缺点：（p123）
虽然通过指针来传递数组参数是高效的，而且也允许在函数内部修改数组的值，但是数组依然是僵化的类型，因为数组的类型包含了僵化的长度信息。
上面的zero函数并不能接收指向[16]byte类型数组的指针，而且也没有任何添加或删除数组元素的方法。

4.2 Slice（p124）

14. slice代表变成的序列，序列中每个元素都有相同的类型。
一个slice类型一般写作[]T，其中T代表slice中元素的类型；slice的语法和数组很像，只是没有固定长度而已。

15. 一个slice是一个轻量级的数据结构，提供了访问数组子序列（或者全部）元素的功能，而且slice的底层确实引用一个数组对象。

一个slice由三个部分构成：指针、长度和容量。

指针：指向第一个slice元素对应的底层数组元素的地址，
要注意的是slice的第一个元素并不一定就是数组的第一个元素。

长度：对应slice中元素的数目；

容量：长度不能超过容量，容量一般是从slice的开始位置到底层数据的结尾位置。

内置的len和cap函数分别返回slice的长度和容量。

16. 多个slice之间可以共享底层的数据，并且引用的数组部分区间可能重叠。
months := [...]string{
	1:"January", 2:"February", 3:"March", 4:"April", 5:"May", 
	6:"June", 7:"July", 8:"August", 9:"September", 10:"October", 11:"November", 12:"December"}
//该数组（切片？）的第1个元素没有写，默认是零值。

17. s[i:j] (p124)
slice的切片操作s[i:j]，其中0<=i<=j<=cap(s)，用于创建一个新的slice，
引用s的第i个元素开始到第j-1个元素的子序列。
新的slice将只有j-i个元素。
如果i位置的索引被省略的话将使用0代替。
如果j位置的索引被省略的话将使用 len(s)代替。

因此，months[1:13]切片操作将引用全部有效的月份，和months[1:]操作等价；
months[:]切片操作则是引用整个数组。

Q2 := months[4:7]
summer := months[6:9]
fmt.Println(Q2) // ["April" "May" "June"]
fmt.Println(summer) // ["June" "July" "August"]

18. 相同元素测试（p125）
for _, s := range summer{
	for _, q := range Q2{ //for后面空格，不要连续
		if s == q{
			fmt.Printf("%S appears in both\n", s)
		}
	}
}

19. 如果切片操作超出 cap(s)的上限将导致一个panic异常，(p125)
但是超出 len(s)则是意味着扩展了slice，因为新的slice的长度会变大：
fmt.Println(summer[:20]) // panic: out of range

endlessSummer := summer[:5] // extent a slice (within capacity)
fmt.Println(endlessSummer) // "[June July August September October]"

20. 字符串的切片操作和[]byte字节类型切片的切片操作是类似的。(p126)
写作x[m:n],并且都是返回一个原始字节系列的子序列，底层都是共享之前的底层数组。
x[m:n]切片操作对于字符串则生成一个新字符串，如果x是 []byte 的话则生成一个新的[]byte。

21. 因为slice值包含指向第一个slice元素的指针，因此向函数传递slice将允许在函数内部修改底层数组的元素。(p126)
复制一个slice只是对底层的数组创建了一个新的slice别名（2.3.2）。

下面的reverse函数在原内存空间将[]int 类型的slice反转，而且它可以用于任意长度的slice。
（有示例代码）

22. slice并没有指明序列的长度，会隐式地创建一个额合适大小的数组，然后slice的指针指向底层的数组。（p126）
slice的字面值也可以按顺序制定初始化值序列，
或者是通过索引和元素值指定，
或者两种风格的混合语法初始化。

23. slice不可比较（p127）
和数组不同，slice之间不能比较，因此我们不能使用==操作符来判断两个slice是否含有全部相等元素。
不过标准库提供了高度优化的 bytes.Equal 函数来判断两个字节型slice是否相等（[]byte），
但是对于其他类型的slice，我们必须展开每个元素进行比较：
func equal(x, y []string) bool{
	if len(x) != len(y){
		return
	}
	for i := range x{
		if x[i]!=y[i]{
			return false
		}
	}
	return true
}

24. slice不可比较的原因（p127）详细看看
安全的做法是直接禁止slice之间的比较操作。

唯一合法的比较操作是nil比较：
if summer == nil{/*...*/}

一个零值的slice等于nil。
一个nil值的slice并没有底层数组。
一个nil值的slice的长度和容量都是0。
但是也有非nil值的slice的长度和容量也是0的。
例如：[]int{}或 make([]int, 3)[3:] 。
与任意类型的nil值一样，我们可以用[]int(nil)类型转换表达式来生成一个对应类型slice的nil值。（将nil转化为slice类型。）
var s []int // len(s) == 0 , s == nil
s = nil // len(s) == 0 , s == nil
s = []int(nil) // len(s) == 0 , s == nil
s = []int{} // len(s) == 0 , s != nil 类似已经预备了。比赛时的预备。。。开始！

25. 如果你需要测试一个slice是否是空的，使用 len(s)== 0 来判断，而不应该用 s==nil 来判断。

26. 除了和nil相等比较外，一个nil值的slice的行为和其他任意0长度的slice一样；
例如 reverse(nil)也是安全的。
除了文档已经说明的地方，所有的go语言函数应该以相同的方式对待nil值的slice和0长度的slice。

27. 内置的make函数创建一个而指定元素类型、长度和容量的slice。
容量部分可以省略，在这种情况下，容量将等于长度。
make([]T, len)
make([]T, len, cap) //same as make([]T,cap)[:len]
make([]T, cap)[:len]

28. 在底层，make创建了一个匿名的数组变量，然后返回一个slice；
只有通过返回的slice才能引用底层匿名的数组变量。
在第一种语句中，slice是整个数组的view。
在第二个语句中，slice只引用了底层数组的前len个元素，但是容量将包含整个的数组。
额外的元素是留给未来的增长用的。

4.2.1 append函数

29. 内置的append函数用于向slice追加元素：
var runes []rune //slice类型
for _, r := range "Hello, 世界"{
	RuneScape= append(runes, r) //给runes，增加r这个元素
}
fmt.Printf("%q\n", runes) //"['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"
fmt.Printf("%T\n", runes)

s := []rune("Hello, 世界")
fmt.Printf("%q\n", s)//%q 带双引号的字符串"abc"或带单引号的字符'c'

30. append原理
每次调用appendint函数，必须先检测slice底层数组是否有足够的容量来保存新添加的元素。（p129）

如果有足够的空间的话，直接扩展slice（依然在原有的底层数组之上），将新添加的y元素复制到xin扩展的空间，并返回slice。
因此，输入的x和输出的z共享相同的底层数组。

如果没有足够的增长空间的话，appendint函数则会先分配一个足够大的slice用于保存新的结果，
先将输入的x复制到新的空间，然后添加y元素。
结果z和输入的x引用的将是不同的底层数组。

31. copy
虽然通过循环复制（添加）元素更直接，不过内置的 copy 函数可以方便地将一个slice复制另一个相同类型的slice。
copy函数的第一个参数是要复制的目标slice，第二个参数是源slice，
目标和源的位置顺序和 dst=src 赋值语句是一致的。
两个slice可以共享同一个底层数组，甚至有重叠没有问题。
copy 函数将返回成功复制的元素的个数（我们这里没有用到），
等于两个slice中较小的长度，所以所以我们不用但因覆盖会超出目标slice的范围。

32. 为了提高内存使用效率，新分配的数组一般略大于保存x和y所需要的最低大小。
通过每次扩展数组时直接将长度翻倍从而避免了多次内存分配，也确保了添加单个元素操作的平均时间是一个常数时间。
这个程序演示了效果：
func main(){
	var x , y []int
	for i := 0; i<10; i++{
		y = appendInt(x, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
		x=y
	}
}

后面的分析结合ch4-append示例看。已经理解。

33. 通常是将append函数返回的结果直接赋值给出入的slice变量;(p131)
runes = append(runes, r)

更新slice变量不仅对调用append函数是必要的，实际上对应任何可能导致长度、容量或底层数组变化的操作都是必要的。(p131)
要正确地使用slice，需要记住尽管底层数组的袁术是间接访问的，但是slice对应结构体本身的指针、长度和容量部分是直接访问的。
要更新这些信息需要像上面例子那样一个显式的赋值操作。
从这个角度看，slice并不是一个纯粹的引用类型，它实际上是一个类似下面结构体的聚合类型：
type IntSlice struct{
	ptr *int
	len, cap int
}

34. 内置的append函数可以追加多个元素，甚至追加slice。
var x []int
x = append(x, 1)
x = append(x, 2, 3)
x = append(x, x...) // append the slice x (x...代表slice？？？)
fmt.Println(x)

35. ...int（p132）
在appendslice代码示例中，...int 省略号表示接收变长的参数为slice。
在5.7 详细解释这个特性。

4.2.2 slice内存技巧

36. 在原有slice内存空间上返回不包含空字符串的列表：(p132)
详见示例代码
如果有变动或更新，通常这样使用：data = nonempty(data)

37. 示例中的 append 函数技巧（p133）
这种方式重用一个slice一般都要求最多为一个输入值产生一个输出值。

38. slice模拟一个stack(堆)（p133）
最初给定的空slice对应一个空的stack，然后可以使用append函数将新的值压入stack：
stack = append(stack, v) // push v
stack的顶部位置对应slice的最后一个元素：
top := stack[len(stack)-1] // top of stack
通过收缩stack可以弹出栈顶的元素。去掉栈顶的元素
stack = stack[:len(stack)-1] // pop

39. 要删除slice中间的某个元素并保存原有元素顺序，可以通过内置的 copy 函数将后面的子slice向前依次移动一位完成。
func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func main(){
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s,2)) // "[5, 6, 8, 9]"
}

40. 如果删除元素后不用保持原来顺序的话，我们可以简单的用最后一个元素覆盖被删除的元素：
func remove2(slice []int, i int) []int {
	slice[i]=slice[len(slice)-1]
	return slice[:len(slice)-1]
}
func main(){
	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s,2)) //[5, 6, 9, 8]
}


4.3 map

41. 哈希表是一种巧妙并且实用的数据结构。 
是一个无序的key/value对的集合，其中所有的key都是不同的，
然后通过给定的key可以在常数时间复杂度内检索、更新或删除对应的value。

在Go语言中，一个map就是一个哈希表的引用，
map 类型可以写为 map[K]V,其中K和V分别对应key和value。
map中所有的key都是相同的类型，所有的value也有着相同的类型，
但是key和value可以是不同的类型。
其中K对应的key必须是支持==比较运算符的数据类型，所以map可以通过测试key是否相等来判断是否已经存在。

虽然浮点数类型也是支持相等运算符比较的，但是将浮点数用做key类型则是一个坏想法。
如第三章提到的，可能出现NaN和任何浮点数都不相等。
对于V对应的value数据类型则没有任何的限制。

42. 内置的make函数可以创建一个map：
ages := make(map[string]int) // mapping from strings to ints

我们也可以用map字面值的语法创建map，同时还可以制定一些最初的key/value：
ages := map[string]int{
	"alice": 31,
	"charlie": 43,
}
//map ["alice": 31 "charlie":43]

这相当于:
ages:= make(map[string]int)
ages["alice"] = 31
ages["charlie"] =34

43. 因此，另一种创建空的map的表达式是 map[string]int{}

44. 访问map（p135）
Map中的元素通过key对应的下标语法访问：
ages["alice"] = 32
fmt.Println(ages["alice"]) // "32"

45. 删除map中的元素（p135）
使用内置的delete函数可以删除元素：
delete(ages, "alice") // remove element ages["slice"]


46. 以上所有操作是安全的，即使这些元素不在map中也没有关系。
如果一个查找失败将返回value类型对应的零值，例如，即使map中不存在“bob”下面的代码也可以正常工作，因为ages["bob"]失败时将返回0。
ages["bob"] = ages["bob"]+1 
//等号右边的bob检索不到返回0值，但是加1后，总数是1。然后赋值给ages["bob"]，所以创建了bob这个key以及也有对应的value。

47. 而且 x += y 和 x++等简短赋值语法也可以用在map上，所以上面的代码可以改成;
ages["bob"] += 1

更简短的写法：
ages["bob"]++

48. 但是map中的元素并不是一个变量，因此我们不能对map的元素进行取指操作：
_ = &ages["bob"] // compile error: cannot take address of map element.
禁止对map元素取址的原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而可能导致之前的地址无效。

49. 想要遍历map中全部的key/value对的话，可以使用range风格的for循环实现，和之前的slice遍历语法类似。
下面的迭代语句将在每次迭代时设置name和age变量，他们对应下一个键/值对：
for name, age := range ages{
	fmt.Printf("%s\t%d\n", name, age)
}

50. map的迭代：随机
map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序。
在实践中，遍历的顺序是随机的，每一次遍历的顺序都不相同。
这是故意的，每次都使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现。

51. map的迭代：顺序
如果要顺序遍历key/value对，我们必须显式地对key进行排序，
可以使用sort包的Strings函数对字符串slice进行排序。
下面是常见的处理方法：
import "sort"

var names []string
for name := range ages{
	names = append(names, name)
}
sort.Strings(names)//Strings函数对slice（string类型的）进行排序
for _, name := range names{
	fmt.Printf("%s\t%d\n", name, ages[name])
}

因为我们一开始就知道names的最终大小，因此给slice分配一个合适的大小将会更有效。
下面的代码创建了一个空的slice，但是slice的容量刚好可以放下map中全部的key：
names := make([]string, 0, len(ages))//创建了一个容量为len(ages)的slice

在上面的第一个range循环中，我们只关心map中的key，所以我们忽略了第二个循环变量。
在第二个循环中，我们只关心names中的名字，所以我们使用_空白标识符来忽略第一个循环变量，也就是迭代slice是的索引。

52. map 类型的零值是nil，也就是没有引用任何哈希表。
var ages map[string]int
fmt.Println(ages == nil) //"true"
fmt.Println(len(ages) == 0) // "true"

53. map 上的大部分操作，包括查找、删除、len和range循环都可以安全工作在nil值的map上，
他们的行为和一个空的map类似。
但是向一个nil值的map存入元素将导致一个panic异常：
ages["carol"] = 21 //panic: assignment to entry in nil map

54. 在想map存数据前必须先创建map。
通过key作为索引下标来访问map将产生一个value。
如果key在map中是存在的，那么将得到与key对应的value；
如果key不存在，那么将得到value对应类型的零值，如上面的ages["bob"]。

但是有时候可能需要知道对应的元素是否真的是在map中，例如，如果元素类型是一个数字，你可以需要群分一个已经存在的0，和不逊在而返回零值的0，这一向下面这样测试：
age, ok := range ages["bob"]
if !ok{/* "bob" is not a key in this map; age == 0. */}

会经常看到将这两个结合起来使用，像这样:
if age, ok := ages["bob"]; !ok{/*...*/}
这种场景下，map的 下标语法 将产生两个值；就第二个是一个布尔值，用于报告元素是否真的存在。
布尔变量一般命名为ok，特别适合马上用于if条件判断部分。

55. map比较（p138）
和slice一样，map之间也不能进行相等比较；
唯一的例外是和nil进行比较。
要判断两个map是否包含相同的key和value，我们必须通过一个循环实现：
func equal(x, y map[string]int) bool{
	if len(x) != len(y){
		return false
	}
	for k, xv := range x{
		if yv, ok := y[k] || yv != xv{
			return false
		}
	}
	return true
}

56. 要注意我们是如何用 !ok 来区分元素确实和元素不同的。
我们不能简单地用 xv != y[k]判断，那样会导致在判断下面两个map是产生错误的结果：
// True if equal is written incorrectly
equal(map[string]int{"A": 0}, map[string]int{"B":42})
//由于检索到没有的元素会返回0值，所以如果没有排除元素是否存在，就直接比较，也许返回的0值，比较后居然相等，得出true，造成显然的错误。

57. Go语言中并没有提供一个set类型，但是map中的key也是不相同的，可以用map实现类似set的功能。（p138）
为了说明这一点，下面的dedup程序读取多行输入，但是只打印第一次出现的行。（1.3节中的dep的变体）
dedup程序 通过map来表示所有的输入行所对应的set集合，以确保已经在集合存在的黄不会被重复打印。

58. Go程序员将这种忽略value的map当做一个字符串集合，并非所有map[string]bool 类型value都是无关紧要的；
有些则可能会同时包含true和false的值。

59. 有时候我们需要一个map或set的key是slice类型，但是map的key必须是可比较的类型，
但是slice并不满足这个条件。
不过，我们可以通过两个步骤绕过这个限制。
第一步，定义一个辅助函数k，将slice转为map对应的string类型的key，
确保只有x和y相等时 k(x) == k(y)才成立。//用这个k(x) == k(y)进行对比，如果可比较根据比较结果，进行下一步。因为slice不能直接对比。
第二步，创建一个key为string类型的map，在每次对map操作是先用k辅助函数将slice转化为string类型。
下面的例子演示了如何使用map来记录提交相同的字符串列表的次数。
它使用了fmt.Sprintf函数将字符串列表转换为一个字符串以用于map的key，通过%q参数忠实的记录每个字符串元素的信息：
var m = make(map[string]int)
func k(list []string)string{return fmt.Sprintf("%q", list)}

func Add(list []string) {m[k(list)]++}
func Count(list []string) int {return m[k(list)]}

60. 使用同样的技术可以处理任何不可比较的key类型，而不仅仅是slice类型。
这种技术对于想使用自定义key比较函数的时候也很有用，例如在比较字符串的时候忽略大小写。
同时，辅助函数 k(x)也不一定是字符串类型，它可以返回任何可比较的类型，例如整数、数组或结构体等。

这是map的另一个例子，下面的程序用于统计输入中每个Unicode码点出现的次数。
虽然Unicode全部码点的数量巨大，但是出现在特定文档中的字符种类并没有多少，
使用map可以用比较自然的方式来跟踪那些出现过的字符的次数。


61. ReadRune方法执行UTF-8解码并返回三个值：
解码的rune字符的值，字符UTF-8编码后的长度，和一个错误值。

我们可预期的错误值已有对应文件结尾的io.EOF。
如果输入的是无效的UTF-8编码的字符，返回的将是unicode.ReplacementChar表示无效字符，
并且编码长度是1。

charcount程序同时打印不同UTF-8编码长度的字符数目。
对此，map并不是一个合适的数据结构；
因为UTF-8编码的长度总是从1到utf8.UTFMax（最大四个字节），
使用数组将更有效。

62. map的value类型也可以是一个聚合类型，比如是一个map或slice。
在代码中，graph的key类型是一个字符串，value类型map[string]bool代表一个字符串集合。
从概念上讲，graph将一个字符串类型的key映射到一组相关的字符串集合，它们指向新的graph的key。

63. 其中addEdge函数惰性初始化map是一个惯用方式，也就是说在每个值首次作为key时才初始化。
addEdge函数显示了如何让map的零值也能正常工作；
即使from到to的边不存在，graph[from][to]依然可以返回一个有意义的结果。 
返回的0值也对应false。所以是有意义的。

4.4 结构体（p143）
64. 结构体是一种聚合的数据类型，是由零个或多个任意类型的值聚合成的实体。
每个值称为结合体等成员。

type Employee struct{ //struct结构体类型
	ID int
	Name string
	Address string
	Dob time.Time
	Position string
	Salary int
	ManagerID int
}
var dilvert Employee

65.结构体变量的成员可以通过 . 操作符访问，
如dilbert.Name
因为dilbert是一个变量，它所有的成员也同样是变量，我们可以直接对每个成员赋值：

dilbert.Salary -= 5000 // demoted(降级), for writing too few lines of code 因为写的代码太少了。

66. 或者是对成员取地址，然后通过指针访问：
position := &dilbert.Position
*position = "Senior " + *position // promoted, for outsourcing（外包） to Elbonia 升职，因为外包给了E

67. 点操作符也可以和指向结构体的指针一起工作：（p143）
var employeeOfTheMonth *Empolyee = &dilbert //*Empolyee指针类型，初始值是&dilbert的地址。
employeeOfTheMonth.Positon += "(proactive team player)"

相当于下面：
(*employeeOfTheMonth).Position += "(proactive team player)"

下面的EmployeeByID函数将根据给定的员工ID返回对应的员工信息结构体的指针。
我们可以使用点操作符来访问它里面的成员：
func EmployeeByID(id int) *Employee{/*...*/}

fmt.Println(EmployeeByID(dilbert.ManagerID).Position) // "Pointy-haired boss"

id := dilbert.ID
EmployeeByID(id).Salary = 0 // fired for ... no real reason.
//EmployeeByID(id)返回的是该id对应的*Employee的所有信息，然后.Salary 则是该id对应的薪水的数据。

68. 后面的语句通过EmployeeByID返回的结构体指针更新了Employee结构体的成员。
如果将EmployeeByID函数的返回值从 *Employee 指针类型改为了Employee值类型，那么更新语句将不能编译通过，
因为在赋值语句的左边并不确定是一个变量（注：调用该函数返回的是值，并不是一个可取地址的变量）

69. 通常一行对应一个结构体成员，成员的名字在前类型在后，不过如果相邻的成员类型如果相同的话，可以被合并到一行。如下：
type Empolyee struct{
	ID int
	Name, Address string
	DoB time.Time
	Positon string
	Salary int
	ManagerID int
}

70. 结构体成员的输入顺序也有重要的意义。我们可以将Position成员合并（因为也是字符串类型），
或者是交换Name和Address出现的先后顺序，那样的话就是定义了不同的结构体类型。
通常，我们只是将相关的成员写到一起。

71. 如果结构体成员（结构体的值）名字是以大写字母开头的，那么该成员就是可导出的。
一个结构体可能同时包含可导出和不可导出的成员。

72. 结构体的类型往往是冗长的，因为它的每个成员可能都会占一行。
虽然我们每次都可以重写整个结构体成员，但是重复会令人厌烦。
因此，完整的结构体写法通常只在类型声明语句的地方出现，就像Employee类型声明语句那样。

73. 一个命名为S的结构体类型将不能再包含S类型的成员：（p144）
因为一个聚合的值不能包换它自身。（该限制同样适应于数组）
但是S类型的结构体可以包换 *s 指针类型的成员，这可以让我们创建地柜的数据结构，
比如链表和树结构等。
下面，使用一个二叉树来实现一个插入排序。
（一部分没看懂）

74. 结构体的零值是每个成员都对是零值。
通常会将零值作为最合理的默认值。
例如：对于bytes.Buffer 类型，结构体初始值就是一个随时可用的空缓存，还有在第9章sync.Mutex的零值
也是有效的未锁定状态。
有时候这种零值可用的特性是自然获得的，但是也有些类型需要一些额外的工作。

75. 如果结构体没有任何成员的话就是空结构体，写作struct{}。
它的大小为0，也不包含任何信息，但是有时候依然是有价值的。
有些Go语言程序员用map带模拟set数据结构时，用它来代替map中布尔类型的value,只是强调key的重要性，
但是因为节约的空间有限，而且语法比较复杂，所以我们通常避免这样的用法。
seen := make(map[string]struct{}) // set of strings  原来的做法查看p139
//...
if _, ok := seen[s]; !ok{ //如果seen[s]有值就ok, 如果没有值就!ok，继续执行下面的
	seen[s] = struct{}{} //struct{}{}结构体的结构体？
	//... first time seeing s...
//其实不需要struct{}有值，就是为了不需要重复出现s
}


4.4.1 结构体面值
76. 结构体值也可以用结构体面值表示，结构体面值可以制定每个成员的值。
type Point struct{	X,Y int }
p := Point{1,2}

77. 两种结构体面值的写法（p146）
（1）上面的是第一种写法，要求以结构体成员定义的顺序为每个结构体成员指定一个面值。
它要求代码和读代码的人要记住结构体的每个成员的类型和顺序，不过结构体成员有细微的调整就可能导致上述代码不能编译。
所以该方法一般在结构体的内部使用，或者是较小的结构体中使用，这些结构体的成员排列比较规则，
比如image.Point{x, y}

（2）以成员名和相应的值来初始化，可以包含部分或全部的成员。
如：
anim := gif.GIF{LoopCount: nframes}
这种更常用。
这种写法中，如果成员被忽略的话将默认用零值。因为，提供了成员的名字，所有成员出现的顺序并不重要。

两种不同形式的写法不能混合使用。
而且不能企图在外部包中用第一种顺序赋值的技巧来偷偷地初始化结构体中未导出的成员。
package p 
type T struct{a, b int} // a and b are not exported

package q 
import "p"
var _ = p.T{a:1, b:2} // compile error: can't reference a, b
var _ = p.T{1, 2} // compile error: can't reference a, b

虽然上面最后一行代码的编译错误信息中并没有显示提到未导出的成员，
但是这样企图隐式使用未导出成员的行为也是不允许的。

78. 结构体可以作为函数的参数和返回值。
例如，这个Scale函数将Point类型的值缩放后返回：
func Scale(p Point, factor int)Point{
	return Point{p.X * factor, p.Y * factor}
}

fmt.Println(Scale(Point{1, 2}, 5)) //"{5 10}"

79. 如果考虑效率的话，较大的结构体通常会用指针的方式传入和返回，
func Bonus(e *Employee, percent int) int{
	return e.Salary * percent / 100
}

80. 如果要在函数内部修改结构体成员的话，用指针传入是必须的；
因为在Go语言中，所有的函数参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量。
func AwardAnnualRaise(e *Employee){
	e.Salary = e.Salary *105/100
}

81. 因为结构体通常通过指针处理，（p147）
可以用下面的写法来创建并初始化一个结构体变量，
并返回结构体的地址：
pp := &Point(1, 2)

等价于：
pp := new(Point)
*pp = Point{1, 2}

不过&Point{1, 2}写法可以直接在表达式中使用，比如一个函数调用。

4.4.2 结构体比较（p147）
82. 如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的，
那样的话两个结构体将可以使用==或!=运算符进行比较。
相等比较运算符==将比较两个结构体的每个成员，
因此下面两个比较的表达式是等价的：
type Point struct{X, Y int}
p := Point{1, 2}
q := Point{2, 1}
fmt.Println(p.X == q.X && p.Y==q.Y) // "false"
fmt.Println(p == q) // "false"

可比较的结构体类型和其他可比较的类型一样，可以用于map的key类型。
type address struct{
	hostname string
	port int
}

hits := make(map[address]int) // map的key是一个结构体
hits[address{"golang.org", 443}]++

4.4.3 结构体嵌入和匿名成员（p148）
83. 我们将看到如何使用Go语言提供的不同寻常的结构体嵌入机制，
让一个命名的结构体包含另一个结构体类型的匿名成员，
这样就可以通过简单的点运算符x.f来访问你们成员链中嵌套的x.d.e.f成员。

考虑一个二维绘图程序，提供了一个各种图形的库，例如矩形，椭圆形，星形，轮形等几何形状。
这里是其中两个的定义：
type Circle struct{
	X, Y, Radius int
}

type Wheel struct{
	X , Y, Radius, Spokes int //spokes辐条
}

一个Circle代表的圆形类型包含了标准圆心的X和Y坐标信息，和一个Radius表示的半径信息。
一个Wheel除了包含Cirlde类型所有的成员外，还增加了Spokes表示径向辐条的数量。
我们可以这样创建一个wheel变量：
var w Wheel
w.X = 8
w.Y = 8
w.Radius = 5
w.Spoles = 20

随着库中几何数量的增多，我们一定会注意到它们之间的相似和重复之处，所以我们可能为了便于维护而将相同的属性独立出来：
type Point struct{
	X, Y int
}

type Circle struct{
	Center Point
	Radius int
}

type Wheel struct{
	Circle Circle //第一个Circle表示成员名，第二个表示类型，这个circle是指上面的的结构体。
	Spokes int
}

这样改动之后结构体类型变动的清晰了，但是这种修改同时也导致了访问每个成员变得繁琐：
var w Wheel
w.Circle.Center.X = 8
w.Circle.Center.Y = 8
w.Circle.Radius = 5
w.Spokes = 20

84. 匿名成员
Go 语言有一个特性让我们值声明一个成员对应的数据类型，而不是成员的名字；
这类成员就叫匿名成员。
匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。
下面的代码中，Circle和Wheel各自都有一个成员。
我们可以说Point类型被嵌入到了Circle结构体，同时Circle类型被嵌入到了Wheel结构体。
type Circle struct{
	Point //Point本身就是一个类型，结构体类型
	Radius int
}

type Wheel struct{
	Circle
	Spokes int
}

得益于匿名嵌入的特性，我们可以直接访问子属性而不需要给出完整的路径：
var w Wheel
w.X = 8 // equivalent to w.Circle.Point.X = 8
w.Y = 8 // equivalent to w.Circle.Point.Y = 8
w.Radius = 5 // equivalent to w.Circle.Radius = 5
w.Spokes = 20

在右边的注释中给出的显式形式访问这些子成员的语法依然有效。
因此匿名成员并不是真的无法访问了。
其中匿名成员Circle和Point都有自己的名字——就是命名的类型名字——但是这些名字在点操作符中是可选的。
我们在访问子成员的时候可以忽略任何匿名成员部分。

85. 不幸的是，结构体字面值并没有简短表示匿名成员的语法，因此下面的语句都不能编译通过：
w = Wheel{8, 8, 5, 20} // compile error: unknow fields
w = Wheel{X: 8, Y:8, Radius:5, Spokes:20} // compile error: unknow fields

86. 结构体字面值必须遵循形状类型声明时的结构，所以我们只能用下面的两种语法，他们彼此是等价的：（p150）
(参考示例代码)
老老实实赋值

87. 需要注意的是Printf函数中%v参数包含的#副词，它表示用和Go语言类似的语法打印值。
对于结构体类型来说，将包含每个成员的名字。

88. 因为匿名成员也有一个隐式的名字，因此不能同时包含两个类型相同个的匿名成员，
这会导致名字冲突。

89. 同时，因为成员的名字是有其类型隐式地决定的，所有匿名成员也有可见性的规则约束。
在上面的例子中，Point和Circle匿名成员都是导出的。
即使他们不导出（比如改成小鞋字母开头的point和circle），我们依然可以用简短形式访问匿名成员嵌套的成员。
w.X = 8 // equivalent to w.circle.point.X = 8

但是在包外部，因为circle和point没有导出，不能访问它们的成员。

90. 我们会看到匿名成员并不要求是结构体类型；其实任何命名的类型都可以作为结构体的匿名成员。
但是为什么要嵌入一个没有任何子成员类型的匿名成员类型呢？
答案是 匿名类型的方法集。（p151）
简短的 点 运算符语法可以用于选择匿名成员嵌套的成员，也可以用于访问他们的方法。
实际上，外层的结构体不仅仅是获得了匿名成员类型的所有成员，而且也获得了该类型导出的全部的方法。
这个机制可以用于将一个有简单行为的对象组合成有复杂行为的对象。
组合是Go语言中面向对象编程的核心，6.3专门讨论。


4.5 JSON（p152）
91. JavaScript对象表示法（JSON）是一种用于发送和接收结构化嘻嘻的标准协议。
在类似的协议中，JSON并不是唯一的一个标准协议。
XML（7.14）ANS.1和Google的Protocol Buffers都是类似的协议，并且有各自的特色，
但是由于简洁性、可读性和流行程度等原因，JSON是应用最广泛的一个。

Go语言对于这些标准格式的编码和解码都有良好的支持，
由标准库中的encoding/json、encoding/xml、encoding/asn1等包提供支持
（注：Protocol Buffers的支持由github.com/golang/protobuf包提供）
并且这类包都有着类似的API接口。
本节，我们将对重要的encoding/json包的用法做个概述。

JSON是对JavaScript中各种类型的值——字符串、数字、布尔值和对象——Unicode文本编码。
它可以用有效可读的方式表示第三章的技术数据类型和本章的数组、slice、结构体和map等聚合数据类型。

基本的JSON类型有数字（十进制或科学计数法）、布尔值（true或false）、字符串、其中字符串是以双引号包含的Unicode字符序列，支持和Go语言类似的反斜杠转义特性，
不过JSON使用的是\Uhhhh转义数字来表示一个UTF-16编码（注：UTF-16和UTF-8一样是一种变长的编码，有些Unicode码点较大的字符需要用4个字节表示；而且UTF-16还有大端和小端的问题）
而不是Go语言的rune（int32）类型。

这些基础类型可以通过JSON的数组和队形类型进行地柜组合。
一个JSON数组是一个有序的值系列，写在一个方括号中并以逗号分隔；一个JSON数组可以用于编码Go语言的数组和slice。
一个JSON对象是一个字符串到值的映射，写成一系列的name:value对形式，用花括号包含并以逗号分隔；
JSON的对象类型可以用于编码Go语言的map类型（key类型是字符串）和结构体。如：
boolean true
number -273.15
string "She said \"Hello, BF\""
array ["gold", "silver", "bronze"]
object {"year": 1980,
		"event":"archery",
		"medals":["gold", "silver", "bronze"]}

92. 考虑一个应用程序，该程序负责手机各种电影评论并提供反馈功能。它的Movie数据类型和一个典型的表示电影的值列表如下所示。
（在结构体声明中，Year和Color成员后面的字符串面值是结构体成员Tag（标签）；我们稍后会解释它的作用）
type Movie struct {
	Title  string
	Year   int      `json:"released"`
	Color  bool     `json:"color, omitempty"`
	Actors []string //数组类型？
}

var movies = []Movie{ //定义了movies是个Movie结构体类型的切片
	{Title: "Casabance", Year: 1942, Color: false,
		Actors: []sting{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}
这样的数据结构特别适合JSON格式，并且在两种之间相互 转换也很容易。（p153）
将一个Go语言中类似movie的结构体slice转为JSOON的过程叫编组（marshaling）。
编组通过调用json.Marshal函数完成：
data, err := json.Marshal(movies)
if err != nil{
	log.Fatalf("JSON marshaling failed: %s", err)	
}
fmt.Printf("%s\n", data)

Marshal函数返还一个编码后的字节slice，包含很长的字符串，并且没有空白缩进；我们将它折行以便于显示：
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingr
id Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Ac
tors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"
Actors":["Steve McQueen","Jacqueline Bisset"]}]

这种紧凑的表示形式虽然包含了全部的信息，但是很难阅读。
为了生成便于阅读的格式，
另一个json。MarshalIndent函数将产生整齐缩进的输出。
该函数有两个额外的字符串参数用于表示每一行输出的前缀和每一个层级的缩进：
data, err := json.MarshlIndent(movies, "", "   ")
if err != nil{
	log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)

生成的样子在（movie）代码中运行查看。

在编码时，默认使用GO语言结构体的成员名字作为JSON的对象（通过reflect反射技术，我们将在12.6讨论）
只有导出的结构体成员才会被编码，也就是我们为什么选择用大写字母开头的成员名称。

93. Tag（p154）
一个结构体成员 Tag 是和在编译阶段关联到该成员的原信息字符串：
Year int `json:"released"`  //对应的示例输出是 "released": 1967
Color bool `json:"color, omitempty"` //对应的示例输出是 "color": true
结构体的成员Tag可以是任意的字符串面值，但是通常是一系列用空格分隔的key:"value"键值对序列；
因为值中含以上双引号字符，因此成员Tag一般用原生字符串面值的形式书写。
json开头键名对应的值用于控制encoding/json包的编码和解码的行为，并且encoding/...下面其他的包也遵循这个约定。
成员Tag中json对应值的第一部分用于制定JSON对象的名字，
如将Go语言中的TotalCount成员对应到JSON中的total——count对象。
Color成员的Tag还带了一个额外的omitempty选项，表示当Go语言结构体成员为空或零值时不生产JSON对象（这里false为0值）
果然，Casablance是一个黑白电影，并没有输出Color成员。

{
	"Title": "Casablanca",
	"released": 1942,
	"Actors": [
		"Humphrey Bogart",
		"Ingrid Bergman"
	]
},

94. JSON解码为Go语言（p155）
编码的逆操作是解码，对应将JSON数据解码为Go语言的数据结构，Go语言中一般叫unmarshaling,
通过json.Unmarshal函数完成。
下面的代码将JSON格式的电影数据解码为一个结构体slice，结构体中只有Title成员。
通过定义合适的Go语言数据结构，我们可以选择性地解码JSON中感兴趣的成员。
当Unmarshal函数调用返回，slice将被只含有Title信息值填充，其他JSON成员将被忽略。
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil{
	log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles) // "[{Casablanca} {Cool Hand Luke} {Bullitt}]"

许多web服务器提供JSON接口，通过HTTP接口发送JSON格式请求并返回JSON合适的信息。
为了说明这一点，我们通过Github的issue查询福娃来演示类似的用法：

95. SearchIssues函数发出一个HTTP请求，然后解码返回的JSON格式的结果。
因为用户提供的查询条件可能包含类似？
和&之类的特殊字符，为了避免对URL造成冲突，
我们用url.QueryEscape来对查询中的特殊字符进行转义操作。

在这个例子中用了基于流式的解码器json.Decoder,
它可以从一个输入流解码JSON数据，尽管这不是必须的。

96. 我们调用Decode方法来填充变量。
这里有多种方法可以格式化结构。
下面是最简单的一种，以一个固定宽度打印每个issue，
但是在下一节我们将看到如果利用模板来输出复杂的格式。


97. GitHub的Web服务接口https://developer.hithub.com/v3/包含了更多的特性。


4.6 文本和HTML模板 (p160)

98. 上面是最简单的格式化，使用Printf就足够了。
但需要复杂的打印格式，这时候一般需要将格式化代码分离出来以便更安全地修改。
这些功能是由text/template和和html/template等模板包提供的，
他们提供了一个将变量值填充到一个文本或HTML格式的模板的机制。

99. 一个模板是一个字符串或一个文件，里面包含了一个或多个由双花括号包含的{{action}}对象。
大部分的字符串只是按面值打印，但是对于actions部分将除法其他的行为。
每个actions都包含了一个用模板语言书写的表达式，一个action虽然简短但是可以输出复杂的打印值，
模板语言包含通过选择结构体的成员、调用函数或方法、表达式控制流 if-else 语句和range循环语句，还有其他实例化模板等诸多特性。
下面是一个简单的模板字符串。

const templ = `{{.TotalCount}} issues:
{{range .Items}}-------------------------------------
Number: {{.Number}}
User: {{.User.Login}}
Title: {{.Title | printf "%.64s"}}
Age: {{.CreatedAt | daysAgo}} days
{{end}}`

这个模板先打印匹配到的issue总数，然后打印每个issue的编号、创建用户、标题还有存在的时间。
对于每一个action， 都有一个当前值的概念，对应点操作符，写作“ . ” 
当前值 . 最初被初始化为调用模板是的参数，在当前例子中对应github.IssuesSearchResult类型的变量。

模板中{{.TotalCount}}对应action将展开为结构体中TotalCount成员以默认的方式打印的值。

模板中{{range.Items}}和{{end}}对应一个循环action， 因此他们之间的内容可能会被展开多次，循环每次迭代的当前值对应当前的Items元素的值.package main


在一个action中，| 操作符表示将前一个表达式的结果作为后一个函数的输入，类似于UNIX中管道的概念。
在Title这一行的action中，第二个操作是一个printf函数，是一个基于fmt.Sprintf实现的内置函数，所有模板都可以直接使用。
对于Age部分，第二个动作是一个叫daysAgo的函数，通过time.Since函数将CreatedAt成员转换为过去的时间长度：
func daysAgo(t time.Time)int{
	return int(time.Since(t).Hours()/24)
}

需要注意的是createdAt的参数类型是time.Time，并不是字符串。
以同样的方式，我们可以通过定义一些方法来控制字符串的格式化（2.5），
一个类型同样可以定制自己的JSON编码和解码行为。
time.Time类型对应的JSON值是一个标准时间格式的字符串。

生成模板的输出需要两个处理步骤。
第一步是要分析模板并转为内部表示，然后基于指定的输入执行模板。
分析模板部分一般值需要执行一次。
下面的代码创建并分析上面定义的模板templ。
注意方法调用链的顺序：
template.New先创建并返回一个模板；
Funcs方法将daysAgo等自定义函数注册到模板中，并返回模板；
最后调用Parse函数分析模板。

report, err := template.New("report").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ)
if err != nil{
	log.Fatal(err)
}

100. 现在让我们转到html/trmplate模板包。
它使用和text/template包相同的API和模板语言，但是增加了一个将字符串自动转义特性，这可以避免输入字符串和HTML、JavaScript、CSS或URL语法产生冲突的问题。
这个特性还可以避免一些长期存在的安全问题，比如通过生成HTML注入攻击，通过构造一个含有恶意代码的问题标题，这些都可能让模板输出错误的输出，从而让他们控制页面。

下面的模板以HTML格式输出issue列表。
注意import语句的不同：
import "html/template"

$ go build gopl.io/ch4/issueshtml
$ ./issueshtml repo:golang/go commenter:gopherbot json encoder >issues.html

$ ./issueshtml repo:golang/go 3133 10535 >issues2.html
生成含有&和<字符的issue的html文件，内容如下：
10535 open dvyukov x/net/html: void element <link> has child nodes 
3133 closed ukai html/template: escape xmldesc as &lt;?xml 

html包已经自动将特殊字符转义，因此我们依然可以看到正确的字面值。
如果我们使用text/template包的haunted，这2个issue将会产生错误，其中 &lt; 四个字符会被但当做 < 字符处理，
同时<link> 字符串将会被当做一个链接元素处理，它们都会导致HTML文档结构的改变，从而导致有未知的风险。

101. 我们也可以通过对信任的HTML字符串使用template.HTML类型来抑制这种自动转义的行为。
还有很多采用类型命名的字符串类型分别对应信任的JavaScript、CSS和URL。
下面的程序演示了两个使用不同类型的相同字符串产生的不同结果：
A是一个普通字符串，
B是一个信任的template.HTML字符串类型。
参考示例：ch4/autoescape
输出后，A转义失效，B成功。
A: <b>Hello!</b>

B: Hello!

102. 我们这里讲述了目标系统中最基本的特性。
一如既往，如果想了解更多的信息，请自己查看包文档：
$ go doc text/template
$ go doc html/template


