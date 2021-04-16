第十一章 测试
（396-426）
Maurice Wikes,第一个存储程序计算机EDSAC的设计者，
1949年他在实验室爬楼梯时有一个顿悟。
在《计算机先驱回忆录》（Memoirs of a Computer Poineer）里，
他回忆到：“忽然间有一种醍醐灌顶的感觉，我整个后半生的美好时刻都将在寻找程序BUG中度过了”。
肯定从那之后的大部分正常的码农都会同情Wikes过份悲观的想法，
虽然也许不是没有人困惑于他对软件开发的难度的天真看法。

现在的程序已经远比Wikes时代的更大也更复杂，
也有许多技术可以让软件的复杂性可得到控制。
其中有两种技术在实践中证明是比较有效的。
第一种是代码在被正式部署前需要进行代码评审。
第二种则是测试，也就是本章的讨论主题。

我们说测试的时候一般是指自动化测试，
也就是写一些小的程序用来检测被测试代码（产品代码）的行为
和预期的一样，
这些通常都是精心设计的执行某些特定的功能或者通过随机性的输入要验证边界的处理。

软件测试是一个巨大的领域。
测试的任务可能已经占据了一些程序员的部分时间和另一些程序员的全部时间。
和软件测试技术相关的图书或博客文章有成千上万之多。
对于每一种主流的编程语言，
都会有一打的用于测试的软件包，
同时都足以说服那些想要编写有效测试的程序员重新学医一套全新的技能。

Go语言的测试技术是相对低级的。
它依赖一个go test测试命令和一组按照约定方式编写的测试函数，
测试命令可以运行这些测试函数。
编写相对轻量级的纯测试代码是有效的，
而且它很容易延伸到基准测试和示例文档。

在实践中，编写测试代码和编写程序本身并没有多大区别。
我们编写的每一个函数也是针对每个具体的任务。
我们必须小心处理边界条件，
思考合适的数据结构，
推断合适的输入应该产生什么样的结果输出。
编程测试代码和编写普通的Go代码过程是类似的；
它并不需要学习新的符号、规则和工具。

11.1.go test (p397)
go test 命令是一个按照一定的约定和组织的测试代码的驱动程序。
在包目录内，所有以_test.go为后缀名的源文件并不是go build构建包的一部分，
他们是go test测试的一部分。

在*_test.go文件中，有三种类型的函数：
测试函数、基准测试函数、示例函数。
1.测试函数：是以Test为函数名前缀的函数，
用于测试程序的一些逻辑行为是否正确；
go test命令会调用这些测试函数并报告测试结果是PASS或FAIL。
2.基准测试函数：是以Benchmark为函数名前缀的函数，
它们用于衡量一些函数的性能；
go test命令会多次运行基准函数以计算一个平均的执行时间。
3.示例函数：是以Example为函数名前缀的函数，
提供一个由编译器保证正确性的实例文档。
我们将在11.2节讨论测试函数的所有细节，
并在11.4节讨论基准测试函数的细节，
然后在11.6节讨论示例函数的细节。

go test命令会遍历所有*_test.go文件中符合上述命名规则的函数，
然后生成一个临时的main包用于调用相应的测试函数，
然后构建并运行、报告测试结果，
最后清理测试中生成的临时文件。

11.2 测试函数（p398）
1. 测试函数
每个测试函数必须导入testing包。测试函数有如下的签名：
func TestName(t *testing.T){
	//...
}

测试函数的名字必须以Test开头，可选的后缀名必须以大写字母开头：
func TestSin(t *testing.T){/*...*/}
func TestCos(t *testing.T){/*...*/}
func TestLog(t *testing.T){/*...*/}

其中t参数用于报告测试失败和附加的日志信息。
让我们定义一个示例包gopl.io/ch11/word1，
其中只有一个函数IsPalindrome用于检查一个字符串是否从前向后和从后向前读都是一样的。
（下面这个实现对于一个字符串是否是回文字符串前后重复测试了两次；
我们稍后会在讨论这个问题。）

示例代码 ch11/word1

在相同的目录下，word_test.go测试文件中包含了TestPalindrome和
TestNonPalindrome两个测试函数。
每一个都是测试IsPalindrome是否给出正确的结果，
并使用t.Error报告失败信息：

示例代码

go test命令如果没有参数指定包那么将默认采用当前目录对应的包
（和 go build 命令一样）。
我们可以用下面的命令构建和运行测试。
$ cd $GOPATH.../ch11/word1
$ go test
ok ...ch11/word1 0.008s

结果还比较满意，我们运行了这个程序，
不过没有提前退出是因为还没有遇到BUG报告。
不过一个法国名为"Noelle Eve Elleon"的用户会抱怨IsPalindrome函数不能识别“été”。
另外一个来自美国中部用户的抱怨则是不能识别"A man, a plan, a canal: Panama."。
执行特殊和小BUG报告为我们提供了新的更自然的测试用例。

示例代码

为了避免两次输入较长的字符串，
我们使用了提供了有类似Pringf格式化功能的Errorf函数来汇报错误结果。

当添加了这两个测试用例之后，
go test返回了测试失败的信息。
$ go test
--- FAIL: TestFrenchPalindrome (0.00s)
	word_test.go:28: IsPalindrome("été") = flase
--- FAIL: TestCanalPalindrome (0.00s)
	word_test.go:35 IsPalindrome("A man, a plan, a canal: Panama") = false
FAIL
FAIL	...ch11/word1 0.014s

先编写测试用例并观察测试用例而出发了和用户报告的错误相同的描述是一个好的测试习惯。
只有这样，我们才能定位我们要正真解决的问题。

先写测试用例的另外的好处是，
运行测试通常会比手工描述报告的处理更快，
这让我们可以进行快速地迭代。
如果测试集有很多运行缓慢的测试，
我们可以通过只选择运行某些特定的测试来加快测试速度。

2.参数-v
参数-v 可以用于打印每个测试函数的名字和运行时间：

$ go test -v
=== RUN TestPalindrome
--- PASS: TestPalindrome (0.00s)
=== RUN TestPalindrome
--- PASS: TestNonPalindrome (0.00s)
=== RUN TestFrenchPalindrome 
--- FAIL: TestFrenchPalindrome (0.00s)
	word_test.go:28: IsPalindromee("été") = false
=== RUN TestCanalPalindrome
--- FAIL: TestCanalPalindrome (0.00s)
word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
FAIL
exit status 1
FAIL gopl.io/ch11/word1 0.017s

3.参数-run
参数 -run 对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被go test测试命令运行：
$ go test -v -run="French|Canal"
=== RUN TestFrenchPalindrome
--- FAIL: TestFrenchPalindrome (0.00s)
word_test.go:28: IsPalindrome("été") = false
=== RUN TestCanalPalindrome
--- FAIL: TestCanalPalindrome (0.00s)
word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
FAIL
exit status 1
FAIL gopl.io/ch11/word1 0.014s

当然，一旦我们已经修复了失败的测试用例，
在我们提交代码更新之前，
我们应该以不带参数的 go test命令运行全部的测试用例，
以确保修复失败测试的同时没有引入新的问题。

我们现在的任务就是修复这些错误。
简要分析后发现第一个BUG的原因是我们采用了byte而不是rune序列，
所以像“été”中的é等非ASCII字符不能正确处理。
第二个BUG是因为没有忽略空格和字母的大小写导致的。

针对上述两个BUG，我们仔细重写了函数：

示例代码ch11/word2

同时我们也将之前的所有测试数据合并到了一个测试中的表格中：

示例代码

现在我们的新测试都通过了：
$ go test ch11/word2
PASS
ok ...ch11/word2 0.015s

这种表格驱动的测试在Go语言中很常见的。
我们很容易向表格添加新的测试数据，并且后面的测试逻辑也没有冗余，
这样我们可以有更多的经理完善错误信息。

失败测试的输出并不包括调用t.Errorf时刻的堆栈调用信息。
和其他编程语言或测试框架的assert断言不同，
t.Errorf调用也没有引起panic异常或停止测试的执行。
即使表格中前面的数据导致了测试的失败，
表格后面的测试数据依然会运行测试，
因此在一个测试中我们可能了解多个失败的信息。

4.停止测试
如果我们真的需要停止测试，
或者是因为初始化失败或可能是早先的错误导致了后续错误等原因，
我们可以使用t.Fatal或t.Fatalf停止当前测试函数。
他们必须在和测试函数同一个goroutine内调用。

5.失败信息
测试失败的信息一般的形式是"f(x) = y, want z"，
其中 f(x) 解释了失败的操作和对应的输出，
y是实际的运行结果，z是期望的正确的结果。
就像前面检查回文字符串的例子，
实际的函数用于 f(x)部分。
显示x是表格驱动型测试中比较重要的部分，
因为同一个断言可能对应不同的表格执行多次。
要避免无用和冗余的信息。
在测试类似IsPalindrome返回布尔类型的函数时，
可以忽略并没有额外信息的z部分。
如果x、y或z是y的长度，输出一个相关部分的简明总结即可。
测试从作者应该要努力帮助程序员诊断测试失败的原因。

11.2.1 随机测试

表格驱动的测试便于构造基于精心挑选的测试数据的测试用例。
另一种测试思路是随机测试，
也就是通过构造更广泛的随机输入来测试探索函数的行为。

那么对于一个随机的输入，
我们如何能知道希望的输出结果呢？
这里有两种处理策略。

6. 第一个是编写另一个对照函数，
使用简单和清晰的算法，
虽然效率较低但是行为和要测试的函数是一致的，
然后针对相同的随机输入检查两者的输出结果。

7.第二种是生成的随机输入的数据遵循特定的模式，
这样我们就可以知道期望的输出模式。

下面的例子使用的是第二种方法：randomPalindrome函数用于随机生成回文字符串。

示例代码

虽然随机测试会有不确定因素，
但是它也是至关重要的，
我们可以从失败测试的日志获取足够的信息。
在我们的例子中，
输入IsPalindrome的p参数将告诉我们真实的数据，
但是对于函数将接受更复杂的输入，
不需要保存所有的输入，
只要日志中简单地记录随机数种子即可（像上面的方式）。
有了这些随机初始化种子，
我们可以很容易修改测试代码以重现失败的随机测试。

通过使用当前时间作为随机种子，
在整个过程中每次运行测试命令时都将探索新的随机数据。
如果你使用的是定期运行的自动化测试集成系统，
随机测试将特别有价值。

11.2.2 测试一个命令（p405）
对于测试包 go test 是一个有用的工具，
但是稍加努力我们也可以用它来测试可执行程序。
如果一个包的名字是main，
那么在构建时会生成一个可执行程序，
不过main包可以作为一个包被测试其代码导入。

让我们为2.3.2节的echo程序编写一个测试。
我们现将程序拆分为两个函数：
echo函数完成真正的工作，
main函数用于处理命令行输入阐述和echo可能返回的错误。

示例代码

在测试中我们可以用各种参数和标志调用echo函数，
然后检测它的输出是否正确，
我们通过增加参数来减少echo函数对全局变量的依存。
我们还增加了一个全局名为out的变量来替代直接使用os.Stdout，
这样测试代码可以根据需要将out修改为不同的对象以便于检查。
下面就是echo_test.go文件中的测试代码：

示例代码

要注意的是测试代码和产品代码在同一个包。
虽然是main包，也有对应的main入口函数，
但是在测试的时候main包只是TestEcho测试函数导入的一个普通包，
里面mian函数并没有被导出，而是被忽略的。

通过将测试放到表格中，
我们很容易添加新的测试用例。
让我通过增加下面的测试用例来看看失败的情况是怎么样的：
{true, ",", []string{"a","b","c"}, "a b c\n"}, //NOTE: wrong expectation!
go test 输出如下：

$ go test ...ch11/echo
--- FAIL: TestEcho (0.00s)
	echo_test.go:31: echo(true, ",", ["a" "b" "c"]) = "a,b,c", want "a b c\n"
	FAIL
	FAIL	...ch11/echo	0.006s

错误信息描述了尝试的操作（使用Go类似语法），
实际的结果和期望的结果。
通过这样的错误信息，
你可以在检视代码之前就很容易定位错误的原因。

要注意的是在测试代码中并没有调用log.Fatal或os.Exit，
因为调用这类函数会导致程序提前退出；
调用这些函数的特权应该放在main函数中。
如果真的有意外的事情导致函数发生panic异常，
测试驱动应该尝试用recover捕获异常，
然后将当前测试当做失败处理。
如果是可预期的错误，
例如非法的用户输入、找不到文件或配置文件不当等应该通过返回一个非空的error的方式处理。
幸运的是（上面的意外只是一个插曲），
我们的echo示例是比较简单的也没有需要返回非空error的情况。

11.2.3 白盒测试(p407)
一种测试分类的方法和基于测试者是否需要了解测试对象的内部工作原理。
黑盒测试只需要测试包公开的文档和API行为，
内部实现对测试代码是透明的。（内部实现是对测试代码不可见的意思？）
相反，白盒测试有访问包内部函数和数据结构的权限，
因此可以做到一个普通客户端无法实现的测试。
例如，一个白盒测试可以在每个操作之后检测不变量的数据类型。
（白盒测试只是一个传统的名称，其实成为clear box测试会更准确）

黑盒和白盒这两种测试方法是互补的。
黑盒测试一般更健壮，
随着软件实现的完善测试代码很少需要更新。
他们可以帮助测试者了解真实客户的需求，
也可以帮助发现API设计的一些不足之处。
相反，白盒测试则可以对内部一些棘手的实现提供更多的测试覆盖。

我们已经看到两种测试的例子。
TestIsPalindrome测试仅仅使用导出的IsPalindrome函数，
因此这是一个黑盒测试。
TestEcho测试则调用了内部的echo函数，
并且更新了内部的out包级变量，
这两个都是未导出的，
因此这是白盒测试。

当我们准备TestEcho测试的时候，
我们修改了echo函数使用包级的out变量作为输出对象，
因此测试代码可以用另一个实现代替标准输出，
这样可以方便对比echo输出的数据。
使用类似的技术，
我们可以将产品代码的其他部分也替换为一个容易测试的伪对象。
使用伪对象的好处是我们可以方便配置，
容易预测，更可靠，也更容易观察。
同时也可以避免一些不良的副作用，
例如更新生产数据库或信用卡消费行为。

下面的代码演示了为用户提供网络存储的web服务中的配额检测逻辑。
当用户使用了超过90%的存储配额之后将发送提醒邮件。

示例代码ch11/storage1

我们想测试这个代码，
但是我们并不希望发送真实的邮件。
因此我们将邮件处理逻辑放到一个私有的notifyUser函数中。


示例代码ch11/storage2

现在我们可以在测试中用伪邮箱发送函数替代真实的邮件发送函数。
它只是简单记录要通知的用户和邮件的内容。

示例代码 test.go

这里有一个问题：
当测试函数返回后，
CheckQuota将不能正常工作，
因为notifyUsers依然使用的是测试函数的伪发送邮件函数
（当更新全局对象的时候总会有这种风险）。
我们必须修改测试代码回复notifyUsers原先的状态以便后续其他的测试没有影响，
要确保所有的执行路径都能恢复，
包括测试失败或panic异常的情形。
在这种情况下，
我们建议使用defer语句来延后执行处理恢复的代码。

示例代码 test.go

这种处理模式可以用来暂时保存和回复所有的全局变量，
包括命令行标志参数、调试选项和优化参数；
安装和移除导致生产代码产生一些调试信息的钩子函数；
还有有些诱导生产代码进入某些重要状态的改变，
比如超时、错误，甚至是一些刻意制造的并发行为等因素。

以这种方式使用全局变量是安全的，
因为go test命令并不会同时并发地执行多个测试。

11.2.4 扩展测试包
考虑下这两个包：net/url包，
提供了URL解析的功能；
net/http包，
提供了web服务和HTTP客户端功能。
如我们所料，上层的net/http包依赖下层的net/url包。
然后，net/url包中的一侧测试是演示不同URL和HTTP客户端的交互行为。
也就是说，
一个下层包的测试代码导入了上层的包。

图11.1 
net/http
 |    ^
 V	  |  (cycle!)
net/url


这样的行为在net/url包的测试代码中会导致包的循环依赖，
正如图11.1中向上箭头所示，
同时正如我们在10.1节所讲的，Go语言规范是进制包的循环依赖的。

不过我们可以通过测试扩展包的方式解决循环依赖的问题，
也就是在net/url包所在的目录声明一个独立的url_test测试扩展包。
其中测试扩展包名的 _test 后缀告诉go test工具它应该建立一个额外的包来运行测试。
我们将这个扩展测试包的导入路径视作是net/url_test会更容易理解，
但实际上它并不能被其他任何包导入。

因为测试扩展包是一个独立的包，
所以可以导入测试代码依赖的其他的辅助包；
包内的测试代码可能无法做到。
在设计层面，
测试扩展包是在它锁依赖的包的上层，如图11.2图所示。

			net/url_test
			 |     |
net/http  <--- 	   |
	|      		   |
    V			   |
net/url <-----------

通过会比循环导入依赖，
扩展测试包可以更灵活的编写测试，
特别是集成测试（需要测试多个组建之间的交互），
可以像普通应用程序那样自由地导入其他包。

我们可以用go list命令查看包对应目录中哪些Go源文件是产品代码，
哪些是包内测试，
还哪些是测试扩展包。
我们以fmt包作为例子：
Gofiles表示产品代码对应的Go源文件列表；
也就是go build命令要编译的部分。

$ go list -f = {{.GoFiles}} fmt //fmt包中，Go源文件列表
[doc.go format.go print.go scan.go]

TestGoFiles表示的是fmt包内部测试代码，
以_test.go为后缀文件名，
不过只在测试时被构建：

$ go list -f = {{.TestGoFiles}} fmt
[export_test.go]

包的测试代码通常都在这些文件中，
不过fmt包并非如此；
稍后我们在解释export_test.go文件的作用。

XTestGoFiles表示的是属于测试扩展包的测试代码，
也就是fmt_test包，
因此它们必须先导入fmt包。
同样，这些文件也只是在测试时被构建运行：

$ go list -f = {{.XTestGoFiles}} fmt
[fmt_test.go scan_test.go stringer_test.go]

有时候测试扩展包也需要访问被测试包内部的代码，
例如在一个为了避免循环导入而被独立到外部测试扩展包的白盒测试。
在这种情况下，
我们可以通过一些技巧解决：
我们在包内的一个_test.go文件中导出一个内部的实现给测试扩展包。
因为这些代码只有在测试时才需要，
因此一般会放在export_test.go文件中。

例如，fmt包的fmt.Scanf函数需要unicode.IsSpace函数提供的功能。
但是为了避免太多的依赖，
fmt包并没有导入包含巨大表格数据的unicode包；
相反fmt包有一个叫isSpace内部的简易实现。

为了确保fmt.isSpace和unicode.IsSpace函数的行为一致，
fmt包谨慎地包含了一个测试。
是一个测试扩展包内的白盒测试，
是无法直接访问到isSpace内部函数的，
因此fmt通过一个秘密出口导出了isSpace函数。
export_test.go文件就是专门用于测试扩展包的秘密出口。

package fmt

var IsSpace = isSpace

这个测试文件并没有定义测试代码；
它只是通过fmt.IsSpace简单导出了内部的isSpace函数，
提供给测试扩展包使用。
这个技巧可以广泛用于位于测试扩展包的白盒测试。

11.2.5编写有效的测试（p413）
许多Go语言新人会惊异于它的极简的测试框架。
很多其他语言的测试框架都提供了识别测试函数的机制
（通常使用反射或元数据），
通过设置一些“setup”和“teardown”的钩子函数来执行测试用例运行的初始化和之后的清理操作，
同时测试工具箱还提供了很多类似assert断言，
值比较函数，格式化输出错误信合和停止一个识别的测试等辅助函数
（通常使用异常机制）。
虽然这些机制可以是的测试非常简洁，
但是测试输出的日志却会像火星文一般难以理解。
此外，虽然测试最终也会输出PASS或FAIL的报告，
但是他们提供的信息格式却非常不利于代码维护者快速定位问题，
因为失败的信息的具体含义是非常隐晦的，
比如"assert: 0 == 1" 或成页的海量跟踪日志。

Go语言的 测试风格则形成鲜明对比。
它期望测试者自己完成大部分的工作，
定义函数避免重复，
就像普通编程那样。
编写测试并不是一个机械的填空过程；
一个测试也有自己的接口，
尽管他的维护者也是测试仅有的一个用户。
一个好的测试不应该引发其他无关的错误信息，
它只要清晰简洁地描述问题的症状即可，
有时候可能还需要一些上下文信息。
在理想情况下，
维护者可以在不看代码的情况下就能根据错误信息定位错误产生的原因。
一个好的测试不应该在遇到一点小错误时就立刻退出测试，
他应该尝试报告更多的相关的错误信息，
因为我们可能从多个失败测试的模式中发现错误产生的规律。

下面的断言函数比较两个值，然后生成一个通用的错误信息，
并停止程序。
它很方便使用也确实有效果，
但是当测试失败的时候，
打印的错误信息却几乎是没有价值的。
他并没有为快速解决问题提供一个很好的入口。

import(
	"fmt"
	"strings"
	"testing"
)

// A poor assertion function.
func assertEqual(x, y int){
	if x != y{
		panic(fmt.Sprintf("%d != %d", x, y))
	}
}
func TestSplit(t *testing.T){
	words := strings.Split("a:b:c", ":")
	assertEqual(len(words), 3)
	// ...
}

从这个意义上说，断言函数犯了过早抽象的错误：
仅仅测试两个整数是否相同，
而放弃了根据上下文提供更有意义的错误信息的做法。
我们可以根据具体的错误打印一个更有价值的错误信息，
就像下面例子那样。
测试在只有一次重复的模式出现时引入抽象。

func TestSplit(t *testing.T){
	s, sep := "a:b:c",":"
	words := strings.Split(s, sep)
	if got, want := len(words), 3; got ！= want {
		t.Errorf("Split(%q, %q) returned %d words, want %d",
		s, sep, got, want)
	}
	// ...
}

现在的测试不仅报告了调用的具体函数、他的输入和结果的意义；
并且打印的真实返回的值和期望返回的值；
并且即使断言失败依然会继续尝试运行更多的测试。
一旦我们写了这样结果的测试，
下一步自然不是用更多if语句来扩展测试用例，
我们可以用像IsPalindrome的表驱动测试那样来准备更多的s和sep测试用例。

前面的例子并不需要额外的辅助函数，
如果有可以使测试代码更简单的方法我们也乐意接受。
（我们将在13.3节看到一个类似reflect.DeepEqual辅助函数）
开始一个好的测试的关键是通过实现你真正想要的具体行为，
然后才是考虑简化测试代码。
最好的接口是直接从库的抽象接口开始，
针对公共接口编写一些测试函数。

11.2.6 避免的不稳定的测试（p414）
如果一个应用程序对于新出现的但有效的输入经常失败说明程序不够稳健；
同样如果一个测试仅仅因为声音变化就会导致失败也是不合逻辑的。
就像一个不够稳健的程序会挫败它的用户一样，
一个脆弱性测试同样会激怒它的维护者。
最脆弱的测试代码会在程序没有任何变化的时候产生不同的结果，
时好时坏，处理他们会耗费大量的时间但是并不会得到任何好处。

当一个测试函数产生一个复杂的输出和一个很长的字符串，
或一个精心设计的数据机构或一个文件，
他可以用于和预设的“golden”结果数据对比，
用这种简单方式写测试是诱人的。
但是随着项目的发展，输出的某些部分很可能会发生变化，
尽管很可能是一个改进的实现导致的。
而且不仅仅是输出部分，
函数复杂复制的输入部分可能也跟着变化了，
因此测试使用的输入也就不再有效了。

避免脆弱测试代码的方法是只检测你真正关心的属性。
保持测试代码的简洁和内部结构的稳定。
特别是对断言部分要有所选择。
不要检查字符串的全匹配，
但是寻找相关的子字符串，
因为某些子字符串在项目的发展中是比较稳定不变的。
通常编写一个从复杂的输出中提取必要精华信息以用于断言是值得的，
虽然这可能会带来很多前期的工作，
但是它可以帮助迅速及时修复因为项目演化而导致的不合逻辑的失败测试。

11.3 测试覆盖率（p415）
就其性质而言，测试不可能是完整的。
计算机科学家Edsger Dijkstra曾说过：“
测试可以显示存在缺陷，但是并不是说没有BUG。”
再多的测试也不能证明一个程序没有BUG。
在最好的情况下，测试可以增强我们的信心：
代码在我们测试的环境是可以正常工作的。

由测试驱动触发运行到的被测试函数的代码数目称为测试的覆盖率。
测试覆盖率并不能量化————甚至连最简单的动态程序也难以精确测量
————但是可以启发并帮助我们编写有效的测试代码。

这些帮助信息中语句的覆盖率是最简单和最广泛使用的。
语句的覆盖率是指在测试中至少被运行一次的代码占总代码数的比例。
在本节中，我们使用 go test 命令中集成的测试覆盖率工具，
来度量下面代码的测试覆盖率，
帮助我们识别测试和我们期望间的差距。

下面的代码是一个表格驱动的测试，用于测试第七章的表达式求值程序：

示例代码ch7/eval/coverage_test

首先，我们要确保所有的测试都正常通过：
$ go test -v -run=Coverage ...ch7/eval //-run=Coverage 等号不能有空格，Coverage首字母必须是大写，否则结果不对。
//另外，ch7/eval必须
=== RUN TestCoverage
--- PASS： TestCoverage (0.00s)
PASS
ok .../ch7/eval 0.011s

下面这个命令可以显示测试覆盖率工具的使用用法：
$ go tool cover
Usage of 'go tool cover':
Given a coverage profile produced by 'go test':
go test -coverprofile=c.out
Open a web browser displaying annotated source code:
go tool cover -html=c.out
...

go tool 命令运行Go工具链的底层可执行程序。
这些底层可执行程序放在$GOROOT/pkg/tool/${GOOS}_${GOARCH}目录。
因为有go build命令的原因，
我们很少直接调用这些底层工具。

现在我们可以用 -coverprofile标志参数重新运行测试：
$ go test -run=Coverage -coverprofile=c.out ...ch7/eval
ok ...ch7/eval  0.032s coverage:68.5% of statements

这个标志参数通过在测试代码中插入生成钩子来统计覆盖率数据。
也就是说，在运行每个测试前，
它会修改要测试代码的副本，
在每个词块都会设置一个布尔标志变量。
当被修改后的被测试代码运行退出时，
将统计日志数据写入c.out文件，
并打印一部分执行的语句的一个总结。
（如果你需要的是摘要，使用go test -cover）

如果使用了 -covermode=count 标志参数，
那么将在每个代码块插入一个计数器而不是布尔标志量。
在统计结果中记录了每个块的执行次数，
这可以用于衡量哪些是被频繁执行的热点代码。

为了收集数据，我们运行了测试覆盖率工具，
打印了测试日志，
生成一个HTML报告，
然后在浏览器中打开（图11.3）。
$ go tool cover -html=c.out

绿色的代码块被测试覆盖到了，
红色的则表示没有被覆盖到。
为了清晰起见，
我们将红色文本的背景设置成了阴影效果。
我们可以马上发现unary操作的Eval方法并没有被执行到。
如果我们针对这部分未被覆盖的代码添加下面的测试用例，
然后重新运行上面的命令，
那么我们将会看到红色部分的代码也变成绿色了：

{"-x * -x", eval.Env{"x":2},"4"}

不过两个panic语句依然是红色的。
这是没有问题的，因为这两个语句并不会被执行到。

实现100%的测试覆盖率听起来很美，
但是在具体实践中通常是不可行的，
也是不值得推荐的做法。
因为那只能说明代码被执行过而已，
并不意味着代码就是没有BUG的；
因为对于逻辑复杂的语句需要针对不同的输入执行多次。
有一些语句，
例如上面的panic语句则永远都不会被执行到。
另外，还有一些隐晦的错误在显示中很少遇到也很难编写对应的测试代码。
测试从本质上来说是一个比较务实的工作，
编写测试代码和编写应用代码的成本对比是需要考虑的。
测试覆盖率工具可以帮助我们快速识别测试薄弱的地方，
但是设计好的测试用例和编写应用代码一样需要严密的思考。

11.4 基准测试（p419）
基准测试是测量一个程序在固定工作负载下的性能。
在Go语言中，基准测试函数和普通测试函数写法类似，
但是以Benchmark为前缀名，
并且带有一个*testing.B类型的参数；
*testing.B参数除了提供和*testing.T类似的方法，
还有额外一些和性能测量相关的方法。
他还提供了一个整数N，
用于指定操作执行的循环测试。

下面是IsPalindrome函数的基准测试，
其中循环将执行N次。

import "testing"

func BenchmarkIsPalindrome(b *testing.B){
	for i := 0; i < b.N; i++{
		IsPalindrome("A man, a plan, a canal: Panama")
	}
}

我们用下面的命令运行基准测试。
和普通测试不同的是，
默认情况下不运行任何基准测试。
我们需要通过 -bench 命令行标志参数手工指定要运行的基准测试函数。
该参数是一个正则表达式，
用于匹配要执行的基准测试函数的名字，默认值是空的。
其中 "." 模式将可以匹配所有基准测试函数，
但是这里总共只有一个基准测试函数，
因此和 -bench=IsPalindrome 参数是等价的效果。

$ cd ...ch11/word2
$ go test -chench=. //匹配所有基本测试函数
PASS
BenchmarkIsPalindrome-8 1000000 	1035 ns/op
ok	...ch11/word2	2.179s

结果中基准测试名的数字后缀部分，
这里是8，
表示运行时对应的GOMAXPROCS的值，
这对于一些和并发相关的基准测试是重要的信息。

报告显示每次调用IsPalindrome函数花费1.035微妙，
是执行1，000，000次的平均时间。
因为基准测试驱动器开始时并不知道每个基准测试函数运行所花的时间，
他会尝试在真正运行基准测试前先尝试用较小的N运行测试来估算基准测试函数所需要的时间，
然后推断一个较大的时间保证稳定的测量结果。

循环在基准测试函数内实现，
而不是放在基准测试框架内实现，
这样可以让每个基准测试函数有机会在循环启动前执行初始化代码，
这样并不会显著影响每次迭代的平均运行时间。
如果还是担心初始化代码部分对测量时间带来干扰，
那么可以通过testing.B参数提供的方法来立领男士关闭或重置计时器，
不过这些一般很少会用到。

现在我们有了一个基准测试和普通测试，
我们可以很容易测试新的让程序运行更快的想法。
也许最明显的优化是在IsPalindrome函数中第二个循环的停止检查，
这样可以避免每个比较都做两次：

n := len(letters)/2
for i := 0; i < n; i++{
	if letters[i] != letters[len(letters)-1-i]{
		return false
	}
}
return true

不过很多情况下，
一个明显的优化并不一定就能让代码有预期的效果。
这个改进在基准测试中只带来4%的性能提升。
$ go test -bench=. 
PASS
BenchmarkIsPalindrome-8 1000000 992 ns/op
ok ...ch11/word2	2.093s

另一个改进想法是在开始为每个字符预先分配一个足够大的数组，
这样就可以避免在append调用时可能会导致内存的多次重新分配。
声明一个leetters数组变量，
并制定合适的大小，
像下面这样，

letters := make([]rune, 0, len(s))
for _, r := range s{
	if unicode.IsLetter(r){
		letters = append(letters, unicode.ToLower(r))
	}
}

这个改进提升性能越35%，
报告结果是基于2，000,000次迭代的平均运行时间统计。
$ go test -bench=.
PASS
BenchmarkIsPalindrome-8 2000000   697 ns/op
ok  ...ch11/word2	1.468s

如这个例子所示，
快的程序往往是伴随着较少的内存分配。（较少？？）
-benchmem 命令行标志参数将在报告中包含内存的分配数据统计。
我们可以比较优化前后内存的分配情况：
$ go test -bench=. -benchmem
PASS
BenchmarkIsPalindrome	1000000 	1026 ns/op	304 B/op	4 allocs/op

这是优化之后的结果：
$ go test -bench=. -benchmem
PASS
BenchmarkIsPallindrome	2000000		807 ns/op	128 B/op	1 allocs/op

用一次内存分配代替多次的内存分配节省了75%的分配调用次数和减少近一般的内存需求。

这个基准测试告诉我们所需的绝对时间依赖给定的具体操作，
两个不同的操作所需时间的差异也是和不同环境相关的。
例如，如果一个函数需要1ms处理1000个元素，
那么处理10000或1百万将需要多少时间呢？
这样的比较揭示了渐进增长函数的运行时间。
另一个例子： I/O缓存该设置为多大呢？
基准测试可以帮助我们选择较小的缓存但能带来满意的性能。
第三个例子：对于一个确定的工作哪种算法更好？
基准测试可以评估两种不同算法对于相同的输入在不同的场景和负载下的优缺点。

一般比较基准测试都是结构类似的代码。
他们通常是采用一个参数的函数，
从几个标志的基准测试函数入口调用，就像这样：
func benchmark(b *testing.B, size int){/*...*/}
func Benchmark10(b *testing.B) {benchmark(b, 10)}
func Benchmark100(b *testing.B) {benchmark(b, 100)}
func Benchmark1000(b *testing.B) {benchmark(b, 1000)}

通过函数参数来指定输入的大小，
但是参数变量对于每个具体的基准测试都是固定的。
要避免直接修改b.N来控制输入的大小。
除非你将它是作为一个固定大小的迭代计算输入，
否则基准测试的结果将毫无意义。

基准测试对于编写代码是很有帮助的，
但是即使工作完成了应当保存基准测试代码。
因为随着项目的发展，
或者是输入的增加，
或者是部署到新的操作系统或不同的处理器，
我们可以再次用基准测试来帮助我们改进设计。

11.5 剖析（p422）
测量基准对于衡量特定操作的性能是有帮助的，
但是当我们试图让程序跑的更快的时候，
我们通常并不知道从哪里开始优化。
每个码农都应该知道Donald Knuth在1974年的“Structured Programming with go to Statements”
上说的格言。虽然经常被解读为不重视性能的意思，
但是从原文我们可以看到不同的含义。

毫无疑问，效率会导致各种滥用。
程序员需要浪费大量的时间思考或者担心，
被部分程序的速度所干扰，
实际上这些尝试提升效率的行为可能产生强烈的负面影响，
特别是当调试和维护的时候。
我们不应该过度纠结于细节的优化，
应该说约97%的场景：
过早的优化是万恶之源。

我们当然不应该放弃那关键的3%的机会。
一个好的程序员不会因为这个理由而满足，
他们会明智地观察和识别那些是关键的代码；
但是只有在关键代码已经被确认的前提下才会进行优化。
对于判断哪些部分是关键代码是经常容易反经验性错误的地方，
因此程序员普通使用的测量工具，是的他们的直觉很不靠谱。

当我们想仔细观察我们程序的运行速度的时候，
最好的技术是如何识别关键代码。
自动化的剖析技术是基于程序执行期间一些抽样数据，
然后推断后面的执行状态；
最终产生一个运行时间的统计数据文件。

Go语言支持多种类型的剖析性能分析，
每一种关注不同的方面，
但他们都涉及到每个采样记录的感兴趣的一系列事件消息，
每个事件都包含函数调用时函数调用堆栈的信息。
内建的go test 工具对集中分析方式都提供了支持。

CPU分析文件标识了函数执行所需要的CPU时间。
当前运行的系统线程在每隔几毫秒都会遇到操作系统的中断事件，
每次中断时都会记录一个分析文件然后恢复正常的运行。

堆分析则记录了程序的内存使用情况。
每个内存分配操作都会除法内部平均内存分配例程，
每个512KB的内存事情都会触发一个事件。

阻塞分析则记录了goroutine最大的阻塞操作，
例如系统调用、管道发送和接收，
还有获取锁等。
分析库会记录每个goroutine被阻塞时的相关操作。

在测试环境下只需要一个标志参数就可以生成各种分析文件。
当一次使用多个标志参数时需要当心，
因为分析操作本身也可能会影响程序的运行。

$ go test -cpuprofile=cpu.out
$ go test -blockprofile=block.out
$ go test -memprofile=mem.out

对于一些非测试程序也很容易支持分析的特性，
具体的实现方式和程序是短时间运行的小工具还是长时间运行的服务会与很大不同，
因此Go的runtime运行时包提供了程序运行时控制分析特性的接口。

一旦我们已经收集到了用于分析的采样数据，
我们就可以使用pprof来分析这些数据。
这是Go工具箱自带的一个工具，
但并不是一个日常工具，
它对应 go tool pprof 命令。
该命令有许多特性和选项，
但是最重要的有两个，
就是生成这个概要文件的可执行程序和对应的分析日志文件。

为了提高分析效率和减少空间，分析日志本身并不包含函数的名字；
它只包含函数对应的地址。
也就是说pprof需要和分析日志对应的可执行程序。
虽然 go test 命令通常会丢弃临时用的测试程序，
但是在启用分析的时候会将测试程序保存为foo.test文件，
其中foo部分对应测试包的名字。

下面的命令演示了如何生成一个CPU分析文件。
我们选择net/http包的一个基准测试为例。
通常是基于一个已经确定了是关键代码的部分进行基准测试。
基准测试会默认包含单元测试，
这里我们用-run=NONE参数禁止单元测试。

$ go test -run=NONE -bench=ClientServerParallelTLS64 \
	-cpuprofile=cpu.log net/http
PASS
BenchmarkClientServerParallelTLS64-8	1000
	3141325 ns/op 143010 B/op 1747 allocs/op
ok net/http 3.395s
$ go tool pprof -text -nodecount=10 ./http.test cpu.log
2570ms of 3590ms total (71.59%)
Dropped 129 nodes (cum <= 17.95ms)
Showing top 10 nodes out of 166 (cum >= 60ms)
flat flat% sum% cum cum%
1730ms 48.19% 48.19% 1750ms 48.75% crypto/elliptic.p256ReduceDegree
230ms 6.41% 54.60% 250ms 6.96% crypto/elliptic.p256Diff
120ms 3.34% 57.94% 120ms 3.34% math/big.addMulVVW
110ms 3.06% 61.00% 110ms 3.06% syscall.Syscall
90ms 2.51% 63.51% 1130ms 31.48% crypto/elliptic.p256Square
70ms 1.95% 65.46% 120ms 3.34% runtime.scanobject
60ms 1.67% 67.13% 830ms 23.12% crypto/elliptic.p256Mul
60ms 1.67% 68.80% 190ms 5.29% math/big.nat.montgomery
50ms 1.39% 70.19% 50ms 1.39% crypto/elliptic.p256ReduceCarry
50ms 1.39% 71.59% 60ms 1.67% crypto/elliptic.p256Sum

参数 -test 用于指定输出格式，
在这里每行是一个函数，
根据使用CPU的时间长短来排序。
其中 -nodecount=10 标志参数限制了只输出前10行的结果。
对于严重的性能问题，
这个文本格式基本可以帮助查明原因了。

这个概要文件告诉我们，
HTTPS基准测试中，
crypto/elliptic.p256ReduceDegree
函数占用了将近一般的CPU资源。
相比之下，
如果一个概要文件中主要是runtime包的内存分配的函数，
那么减少内存消耗可能是一个值得尝试的优化策略。

对于一些更微妙的问题，
你可能需要使用pprof的图形显示功能。
这个需要安装GraphViz工具，
可以从http://www.graphviz.org下载。
参数 -web 用于生成一个有向图文件，
包含了CPU的使用和最热点的函数等信息。

这一节我们只是简单看了下Go语言的分析数据工具。
如果想了解更多，可以阅读Go官方博客的“Profiling Go Programs”一文。

11.6 示例函数（p425）
第三种 go test 特别处理的函数是示例函数，
以Example为函数名开头。
示例函数没有函数参数和返回值。
下面是IsPalindrome函数对应的示例函数：

func ExampleIsPalindrome(){
	fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println(IsPalindrome("palindrome"))
	// Output:
	// true
	// false
}

示例函数有三个用处。
最主要的一个是作为文档：
一个包的例子可以更简洁直观的方式来演示函数的用法，
币文字藐视更直接易懂，
特别是作为一个提醒或快速参考时。
一个示例函数也可以方便展示属于同一个接口的几种类型或函数直接的关系，
所有的文档都必须关联到一个地方，
就像一个类型或函数晟敏都统一到包一样。
同时，示例函数和注释并不一样，
示例函数是完整真实的Go代码，
需要接收编译器的编译时检查，
这样可以保证示例代码不会腐烂成不能使用的旧代码。

根据示例函数的后缀名部分，
godoc的web文档会将一个示例函数关联到某个具体函数或包本身，
因此ExampleIsPalindrome示例函数将是IsPalindrome函数文档的一部分，
Example示例函数将是包文档的一部分。

示例文档的第二个用处是在 go test 执行测试的时候也运行示例函数测试。
如果示例函数内含有类似上面例子中的 // Output: 格式的注释，
那么测试工具会执行这个示例函数，
然后检测这个示例函数的标准输出和注释是否匹配。

示例函数的第三个目的提供一个真实的演练场。
http://golang.org 就是由godoc提供的文档服务，
它使用了Go Playground提高的技术让用户可以在浏览器中在线编辑和运行每个示例函数，
就像图11.4所示的那样。
这通常是学习函数使用或Go语言特性最快捷的方式。

图11.4

本书最后的两章是讨论reflect和unsafe包，
一般的Go用户很少直接使用它们。
因此，如果你还没有写过任何真实的Go程序的话，
现在可以忽略剩余部分而直接编码了。



