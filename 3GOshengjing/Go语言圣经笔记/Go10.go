第十章 包和工具
（p373）

现在随便一个小程序的实现都可能包含超过10000个函数。
然而作者一般只需要考虑其中很小的一部分和做很少的涉及，
因为绝大部分代码都是由他人编写的，它们通过类似包或模块的方式被重用。

Go语言有超过100个的标准包
（译注：可以用 go list std | wc -1 命令查看标准包的具体数目）
标准库为大多数的程序提供了必要的基础构件。
在Go的社区，有很多成熟的包被设计、共享、重用和改进，
目前互联网上已经发布了非常多的Go语言开源包，
他们可以通过http://godoc.org检索。
在本章，我们将演示如何使用已有的包和创建新的包。

Go还自带了工具箱，里面有用来简化工作区和包管理的小工具。
在本书开始的时候，我们已经见识过如何使用工具箱自带的工具来下载、构件和运行我们的演示程序了。
在本章，我们将看看这些工具的基本设计理论和尝试更多的功能，
例如打印工作区中的文档和查询相关的元数据等。
在下一章，我们将探讨包的单元测试用法。

10.1 包简介（p374）
任何包系统设计的目的都是为了简化大型程序的涉及和维护工作，
通过将一组相关的特性放进一个独立的单元以便于理解和更新，
在每个单元更新的同时保持和程序中其他单元的相对独立性。
这种模块化的特性允许每个包可以被其他的不同项目共享和重用，
在项目范围内、甚至全球范围统一的分发和复用。

每个包一般都定义了一个不同的名字空间用于它内部的每个标识符的访问。
每个名字空间关联到一个特定的包，让我们给类型、函数等选择简短明了的名字，
这样可以避免在我们使用它们的时候减少和其他部分名字的冲突。

每个包还通过控制包内名字的可见性和是否爆出来了实现封装特性。
通过限制包成员的可见性并隐藏包API的具体实现，
将允许包的维护者在不影响外部包用户的前提下调整包的内部实现。
通过限制包内变量的可见性，还可以强制用户通过某些特定函数来访问和更新内部变量，
这样可以保证内部变量的一致性和并发时的互斥约束。

当我们修改了一个源文件，我们必须诚信编译该源文件对应的包和所有依赖该包的其他包。
即使是从头构建，Go语言编译器的编译速度也明显快于其他编译语言。
Go语言的闪电般的编译速度主要得益于三个语言特性。
第一点，所有导入的包必须在每个文件的开头显式声明，
这样的话编译器就没有必要读取和分析整个源文件来判断包的依赖关系。
第二点，禁止包的环状依赖，因为没有循环依赖，包的依赖关系形成一个有向无环图，
每个包可以被独立编译，
而且很可能是被并发编译。
第三点，编译后包的目标文件不仅仅记录包本身的导出信息，
目标文件同时还记录了包的依赖关系。
因此，在编译一个包的时候，
编译器只需要读取每个直接导入包的目标文件，
而不需要遍历所有依赖的文件。
（译注:很多都是重复的间接依赖）

10.2 导入路径（p375）
每个包是由一个全局唯一的字符串所标识的导入路径定位。
出现在import语句中的导入路径也是字符串。
import(
	"fmt"
	"math/rand"
	"encoding/json"

	"golang.org/x/net/html"
	
	"github.com/go-sql-driver/mysql"
)

就像我们在2.6.1节提到过的，Go语言的规范并没有指明包的导入路径字符串的具体含义，
导入路径的具体含义是由构建工具来解释的。
在本章，我们经深入讨论Go语言工具箱的功能，
包括大家经常使用的构建测试等功能。
当然，也有第三方扩展的工具箱存在。
例如，Google公司内部的Go语言码农，他们就使用内部的多语言构建系统
（译注：Google公司使用的是类似Bazel的构建系统，
支持多种编程语言，目前该构建系统还不能完整支持Windows环境），
用不同的规则来处理包名字和定位包，
用不同的规则来处理单元测试等等，因为这样可以更紧密适配他们内部环境。

如果你计划分享或发布包，那么导入路径最好是全球唯一的。
为了避免冲突，所有非标准库包的导入路径建议以所在组织的互联网域名为前缀；
而且这样也有利于包的检索。
例如，上面的import语句导入了Go团队维护的HTML解析器和一个流行的第三方维护的MySQL驱动。

10.3包声明（p376）
在每个Go语音源文件的开头都必须有包声明语句。
包声明语句的主要目的是确定当前包被其他包导入时默认的标识符（也称为包名）。

例如，math/rand包的每个源文件的开头都包含 package rand 包声明语句，
所以当你导入这个包，你就可以用rand.Int、rand.Float64类似的方式访问包的成员。

package main

import(
	"fmt"
	"math/rand"
)

func main(){
	fmt.Println(rand.Int())
}

通常来说，默认的包名就是包导入路径名的最后一段，
因此即使两个包的导入路径不同，
他们依然可能有一个相同的包名。
例如，math/rand包和crypto/rand包的包名都是rand。
稍后我们将看到如何同时导入两个相同包名的包。

关于默认包名一般采用导入路径的最后一段的约定也有三种例外情况。
第一个例外，包对应一个可执行程序，也就是main包，
这时候main包本身的导入路径是无关紧要的。
名字为main的包是给go build（10.7.3）构建命令一个信息，
这个包编译完成之后必须调用连接器生成一个可执行程序。

第二个例外，包所在的目录中可能有一些文件名是以test.go为后缀的Go源文件
（译注：前面必须有其他的字符，因为以``前缀的源文件是被忽略的），
并且这些源文件声明的包名也是以_test为后缀名的。
这种目录可以包含两种包：一种普通包，加一种则是测试的外部扩展包。
所有以_test为后缀的测试外部扩展包都由go test命令独立编译，
普通包和测试的外部扩展包是相互独立的。
测试的外部扩展包一般用来避免测试代码中的循环导入依赖，
具体细节我们将在11.2.4节中介绍。

第三个例外，一些依赖版本号的管理工具会在导入路径后追加版本号信息，
例如“gopkg.in/yaml.v2”。
这种情况下包的名字并不包含版本号后缀，而是yaml。

10.4 导入声明（p377）
可以在一个Go语言源文件包声明语句之后，其他非导入声明语句之前，
包含零到多个导入包声明语句。
每个导入声明可以单独制定一个导入路径，
也可以通过圆括号同时导入多个路径。
下面两个导入形式是等价的，但是第二种形式更为常见。
import "fmt"
import "os"

import(
	"fmt"
	"os"
)

导入的包之间可以通过添加空行来分组；
通常将来如不同组织的包独自分组。
包的导入顺序无关紧要，
但是在每个分组中一般会根据字符串顺序排列。
（gofmt和goimports工具都可以将不同分组导入的包独立排序）

import(
	"fmt"
	"html/template"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/ipv4"
)

如果我们想同时导入两个有着名字相同的包，
例如math/rand包和crypto/rand包，
那么导入声明必须至少为一个同名包指定一个新的包名以避免冲突。
这叫做导入包的重命名。

import(
	"crypto/rand"
	mrand "math/rand" // alternative name mrand avoids conflict
)

导入包的重命名只影响当前的源文件。
其他的源文件如果导入了相同的包，
可以用导入包原本默认的名字或重命名为另一个完全不同的名字。

导入包重命名是一个有用的特性，它不仅仅为了解决名字冲突。
如果导入的一个包名很笨重，特别是在一些自动生成的代码中，
这时候用一个简短名称会更方便。
选择用简短名称重命名导入包时候最好统一，
以避免包名混乱。
选择另一个包名称可以帮助避免和本地普通变量名产生冲突。
例如，如果文件中哥已经有了一个名为path的变量，
那么我们可以将“path”标准包重命名为pathpkg。

每个导入声明语句都明确指定了当前包和被导入包之间的依赖关系。
如果遇到包循环导入的情况，//什么叫包循环导入？
Go语言的构建工具将报告错误。

10.5 包的匿名导入（p379）//要理解哦
如果只是导入一个包而并不使用导入的包将会导致一个编译错误。
但是有时候我们只是想利用导入包而产生的副作用：
他会计算包级变量的初始化表达式和执行导入包的init初始化函数（2.6.2）。
这时候我们需要抑制“unused import”编译错误，
我们可以用下划线 _ 来重命名导入的包。
像往常一样，下划线 _ 为空白标识符，并不能被访问。

import _ "image/png" // register PNG decoder

这个被称为包的匿名导入。
它通常是用来实现一个编译时机制，
然后通过在main主程序入口选择性地导入附加的包。
首先，让我们看看如何使用该特性，
然后再看看它是如何工作的。

标准库的image图像包含了一个 Decode 函数，
用于从 io.Reader 接口读取数据并解码图像，
它调用底层注册的图像解码器来完成任务，
然后返回image.Image类型的图像。
使用image.Decode很容易编写一个图像格式的转换工具，
读取一个格式的图像，然后编码为另一种图像格式：

示例代码ch10/jpeg

如果我们将 ch3/mandelbrot (3.3)	的输出导入到这个程序的标准输入，
它将解码输入的PNG格式图像，
然后转换为JPEG格式的图像输出。
$ go build gopl.io/ch3/mandelbrot
$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
Input format = png
// $ ./mandelbrot > ./jpeg >mandelbrot.jpg //将mandelbrot输出的东西，输入到jpeg里，然后再输出jpeg格式。
// 也可以获得相应的图片。


要注意image/png包的匿名导入语句（ _ image/png ）。
如果没有这一行语句，程序依然可以编译和运行，
但是他将不能正确识别和解码PNG格式的图像：

$ go build gopl.io/ch10/jpeg
$ ./mandelbrot | ./jpeg >mandelbrot.jpg
jpeg: image: unknown format

下面的代码演示了它的工作机制。
标准库还提供了GIF、PNG和JPEG等格式图像的解码器，
用户也可以提供自己的解码器，
但是为了保持程序体积较小，
很多解码器并没有被全部包含，
除非是明确需要支持的格式。
image.Decode函数在解码时会依次查询支持的格式列表。
每个格式驱动列表的每个入口指定了四件事情：
格式的名称；
一个用于描述这种图像数据开头部分模式的字符串，用于解码器检测识别；
一个Decode函数用于完成解码图像工作；
一个DecodeConfig函数用于解码图像的大小和颜色空间等信息。
每个驱动入口是通过调用image.RegisterFormat函数注册，
一般是在每个格式包的init初始化函数中调用，
例如，image/png包是这样注册的：

package png // image/png

func Decode(r io.Reader)(image.Image, error)
func DecodeConfig(r io.Reader)(image.Config, error)

func init(){
	const pngHeader = "\x89PNG\r\n\x1a\n"
	image.ResgisterFormat("png", pngHeader, Decode, DecodeConfig)
}

最终的效果是，主程序只需要匿名导入特定图像驱动包就可以用image.Decode解码对应格式的图像了。

数据库包database/sql也是采用了类似的技术，让用户可以根据自己需要选择导入必要的数据库驱动。
例如：
import(
	"database/sql"
	_ "github.com/lib/pq" // enable support for Postgres （类似于我只需要这个格式等）
	_ "github.com/go-sql-driver/mysql" // enable support for MySQL
)

db, err = sql.Open("postgres", dbname) // OK
db, err = sql.Open("mysql", dbname) // OK
db, err = sql.Open("sqlite3", dbname) // returns error: unknow driver "sqlite3"

10.6 包和命名（p382）
在本节中，我们将提供一些关于Go语言独特的包和成员命名的约定。

当创建一个包，一般要用短小的包名，但也不能太短导致难以理解。
标准库中最常用的包有bufio、bytes、flag、fmt、http、io、json、os、sort、sync和time等包。

它们的名字都简洁明了。
例如，不要将一个类似imageutil或ioutilis的通用包命名为util，
虽然它看起来很短小。
要尽量避免包名使用可能被经常用于局部变量的名字，
这样可能导致用户重命名导入包，例如前面看到的path包。

包名一般采用单数的形式。
标准库的bytes、errors和strings使用了复数形式，
这是为了避免和预定义的类型冲突，
同样还有go/types是为了避免和type关键字冲突。

要避免包名有其他的含义。
例如，2.5节中我们的温度转换包最初使用了temp包名，
虽然并没有持续多久。
但这是一个糟糕的尝试，
因为temp几乎是临时变量的同义词。
然后我们有一段时间使用了temperature作为包名，
虽然名字名没有表达包的真实用途。
最后我们改成了和strconv标准包类似的tempconv包名，
这个名字比之前的就好多了。

现在让我们看看如何命名包的成员。
由于是通过包的导入名字引入包里面的成员，
例如fmt.Println，同时包含了包名和成员名信息。
因此，我们一般并不需要关注Println的具体内容，
因为fmt包名已经包含了这个信息。
当设计一个包的时候，需要考虑包名和成员名两个部分如何很好地配合。
下面有一些例子：
bytes.Equal
flag.Int
http.Get
json.Marshal

我们可以看到一些常用的命名模式。
strings包提供了和字符串相关的诸多操作：
package strings

func Index(needle, haystack string) int

type Replacer struct{ /*...*/}
func NewReplacer(oldnew ...string) *Replacer

type Reader struct{/*...*/}
func NewReader(s string) *Reader

字符串string本身并没有出现在每个成员名字中。
因为用户会这样引用这些成员 string.Index 、 strings.Replacer等。

其他一些包，可能只描述了单一的数据类型，例如html/template和math/rand等，
只暴露一个主要的数据结构和它相关的方法，
还有一个以New命名的函数用于创建实例。

package rand // "math/rand"

type Rand struct{/*...*/}
func New(source Source) *Rand

这可能导致一些名字重复，例如template.Template或rand.Rand，
这就是为什么这些种类的包名往往特别短的原因之一。

在另一个极端，
还有像net/http包那样含有非常多的名字和种类不多的数据类型，
因为它们都是要执行一个复杂的符合任务。
尽管有将近二十种类型和更多的函数，
但是包中最重要的成员名字却是简单明了的：Get、Post、Handle、Error、Client、Server等。

10.7 工具（p384）
本章剩下的部分将讨论Go语言工具箱的具体功能，
包括如何下载、格式化、构建、测试和安装Go语言编写的程序。

Go语言的工具箱集合了一系列的功能的命令集。
它可以看作是一个包管理器（类似于Linux中的apt和rpm工具），
用于包的查询、计算的包依赖关系、从远程版本控制系统和下载它们等任务。
它也是一个构建系统，
计算文件的依赖关系，
然后调用编译器、汇编器和连接器构建程序，
虽然它故意被设计成没有标准的make命令那么复杂。
它也是一个单元测试和基准测试的驱动程序，
我们将在第11章讨论测试话题。

Go语言工具箱的命令有着类似“瑞士军刀”的风格，
带着一打子的子命令，
有一些我们经常用到，
例如get、run、build和fmt等。
你可以运行go或go help命令查看内置的帮助文档，
为了查询方便，我们列出最常用的命令：
$ go
...
	build 		compile packages and dependencies 编译包和依赖项
	clean 		remove object files 删除对象文件
	doc			show documentation for package or symbol 显示包或符号的文档
	env			print Go environment information 打印Go环境信息
	fmt			run gofmt on pockage sources 在pockage源上运行gofmt
	get			download and install packages and dependencies 下载并安装软件包和依赖项
	install		compile and install packages and dependencies 编译和安装包和依赖项
	list		list packages 列出程序包
	run			compile and run Go program 编译并运行Go程序
	test		test package 测试包
	version		print Go version 打印Go版本
	vet			run go tool vet on packages 在包上运行go tool vet

Use "go help [command]" for more information about a command.package main
...

为了达到零配置的设计目标，
Go语言的工具箱很多地方都依赖各种约定。
例如，根据给定的源文件的名称，
Go语言的工具可以找到源文件对应的包，
因为每个目录只包含了单一的包，
并且到的导入路径和工作区的目录结构是对应的。
给定一个包的导入路径，
Go语言的工具可以找到对应的目录中每个实体对应的源文件。
他还可以根据导入路径找到存储代码仓库的远程服务器的URL。

10.7.1工作区结构
对于大多数的Go语言用户，只需要配置一个名叫GOPATH的环境变量，
用来指定当前工作目录即可。
当需要切换到不同工作区的时候，
只要更新GOPATH就可以了。
例如，我们在编写本书时将GOPATH设置为 $HOME/gobook ：
$ export GOPATH=$HOME/gobook
$ go get gopl.io/...package main

当你用前面介绍的命令下载本书全部的例子源码之后，你的当前工作区的牡蛎结构应该是这样的:

GOPATH/
	src/
		gopl.io/
			.git/
			ch1/
				helloworld/
					main.go
				dup/
					main.go
				...
		golang.org/x/net/
			.git/
			html/
				parse.go
				node.go
				...
	bin/
		helloworld
		dup	
	pkg/
		darwin_amd64/
		...

GOPATH对应的工作区目录有三个子目录。
其中src子目录用于存储源代码。
每个包都被保存在与$GOPATH/src的相对路径为包导入路径的子目录中，
例如gopl.io/ch1/helloworld相对应的路径目录。
我们看到，一个GOPATH工作区的src目录中可能有多个独立的版本控制系统，
例如gopl.io和golang.org分别对应不同的Git仓库。
其中pkg子目录用于保存编译后的包的目标文件，
bin子目录用于保存编译后的可执行程序，
例如helloworld可执行程序。

第二个环境变量GOROOT用来指定Go的安装目录，
还有它自带的标准库的位置。
GOROOT的目录结构和GOPATH类似，
因此存放fmt包的源代码对应目录应该为$GOROOT/src/fmt
用户一般不需要设置GOROOT，
默认情况下Go语言安装工具会将其设置为安装的目录路径。

其中go env 命令用于查看Go语言工具涉及的所有环境变量的值，
包括未设置环境变量的默认值。
GOOS环境变量用于制定目标操作系统（例如android、linux、darwin或windows），
GOARCH环境变量用于制定处理器的类型，
例如amd64、386或arm等。
岁演GOPATH环境变量是唯一必须要设置的，但是其他环境变量也会偶尔用到。

示例：
$ go env
GOPATH="/home/gopher/gobook"
GOROOT="/usr/local/go"
GOARCH="amd64"
GOOS="darwin"
...

10.7.2下载包
使用Go语言工具箱的go命令，
不仅可以根据包导入路径找到本地工作区的包，
甚至可以从互联网上找到和更新包。

使用命令 go get 可以下载一个单一的包或者用 ... 下载整个子目录里面的每个包。
Go语言工具箱的go命令同时计算并下载所依赖的每个包，
这也是前一个例子中golang.org/x/net/html自动出现在本地工作区目录的原因。

一旦 go get 命令下载了包，然后就是安装包或包对应的可执行的程序。
我们将在下一节再关注它的细节，
现在只是展示整个下载过程是如何的简单。
第一个命令是获取golint工具，
他用于检测Go源代码的编程风格是否有问题。
第二个命令是用golint命令对2.6.2节的ch2/popcount包代码进行编码风格检查。
它友好地报告了忘记了包的文档：
$ go get github.com/golang/lint/golint
$ $GOPATH/bin/golint  地址/ch2/popcount
src/gopl.io/ch2/popcount/main.go:1:1:
	package comment should be of the form "Package popcount ..."

go get 命令支持当前流行的托管网站Github、Bitbucket和Launchpad，
可以直接向他们的版本控制系统请求代码。
对于其他的网站，你可能需要指定版本控制系统的具体路径和协议，
例如Git或Mercurial。
运行 go help importpath 获取相关信息。

go get 命令获取的代码是真实的本地存储仓库，
而不仅仅只是复制源文件，
因此你依然可以使用版本管理工具比较本地代码的变更或者切换到其他的版本。
例如golang.org/x/net包目录对应一个Git仓库：

$ cd $GOPATH/src/golang.org/x/net
$ git remote -v
origin https://go.googlesource.com/net (fetch)
origin https://go.googlesource.com/net (push)

需要注意的是导入路径含有的网址域名和本地Git仓库对应远程服务地址并不相同，
真实的Git地址是go.googlesource.com。
这其实是Go语言工具的一个特性，
可以让包用一个自定义的导入路径，
但是真实的代码确实由更通用的服务提供，
例如googlesource.com或github.com。
因为页面https://golang.org/x/net/html 包含了如下的元数据，
它告诉Go语言的工具当前包真实的Git仓库托管地址：

$ go build ...ch1/fetch
$ ./fetch https://goland.org/x/net/html | grep go-import
<meta name="go-import"
	content="golang.org/x/net git https://go.googlesource.com/net">

如果指定 -u 命令行标志参数，go get 命令将确保所有的包和依赖的版本都是最新的，
然后重新编译和安装他们。
如果不包含该标志参数的话，而且如果包已经在本地存在了，
那么代码将不会被自动更新。

go get -u 命令只是简单地保证每个包是最新版本，
如果是第一次下载包则是比较方便的；
但是对于发布程序则可能是不合适的，
因为本地程序可能需要对依赖的包做精确的版本依赖管理。
通常的解决方案是使用vendor的目录用于存储依赖包的固定版本的源代码，
对本地依赖的包的版本更新也是谨慎和持续可控的。
在Go1.5之前，一般需要修改包的导入路径，
所以复制后golang.org/x/net/html导入路径可能会变为
gopl.io/vendor/golang.org/x/net/html。
最新的Go语言命令已经支持vendor特性，但限于篇幅这里并不讨论vendor的具体细节。
不过可通过 go help gopath 命令查看Vendor的帮助文档。

10.7.3 构建包（p387）
go build 命令编译命令行参数指定的每个包。
如果包是一个库，
则忽略输出结果；这可以用于检测包的可以正确编译的。
如果包的名字是main，
go build将调用连接器在当前目录创建一个可执行程序；
以导入路径的最后一段作为可执行程序的名字。

因为每个目录只包含一个包，//这里指的是main包？
因此每个对应可执行程序或者叫Unix术语中的命令的包，
会要求放到一个独立的目录中。
这些目录有时候会放在名叫cmd目录的子目录下面，
例如用于提供Go文档服务的golang.org/x/tools/cmd/godoc命令就是放在cmd子目录（10.7.4）

每个包可以由它们的导入路径指定，
就像前面看到的那样，或者用一个相对目录的路径值指定，
相对路径必须以 . 或 .. 开头。
如果没有指定参数，那么默认指定为当前目录对应的包。
下面的命令用于构建同一个包，虽然它们的写法各不相同：

$ cd $GOPATH/src/gopl.io/ch1/helloworld
$ go build

或者
$ cd anywhere
$ go build ...ch1/helloworld

或者：
$ cd $GOPATH
$ go build ./src/...ch1/helloworld

但不能这样：
$ cd $GOPATH
$ go build src/...ch1/helloworld
Error: cannot find package "src/...ch1/helloworld"

也可以指定包的源文件列表，
这一般只用于构建一些小程序或做一些临时性的实验。
如果是main包，将会以第一个Go源文件的基础文件名作为最终的可执行程序的名字。
$ cat quoteargs.go
package main

import(
	"fmt"
	"os"
)

func main(){
	fmt.Printf("%q\n", os.Args[1:])
}
$ go build quoteargs.go
$ ./quoteargs one "two three" four\ five
["one" "two three" "four five"]

特别是对于这类一次性运行的程序，
我们希望尽快的构建并运行它。
go run 命令实际上是结合了构建和运行的两个步骤：
$ go run quoteargs.go one "two three" four\ five
["one" "two three" "four five"]

第一行的参数列表中，第一个不是以 .go 结尾的将作为可执行程序的参数运行。

默认情况下，
go build 命令构建指定的包和它依赖的包，
然后丢弃除了最后的可执行文件之外所有的中间编译结果。
依赖分析和编译过程虽然都是很快的，
但是随着项目增加到几十个包和成千上万行代码，
依赖关系分析和编译时间的消耗将变的可观，
有时候可能需要几秒钟，
即使这些依赖项没有改变。

go install 命令和 go build 命令很相似，
但是它会保存每个包的编译成功，
而不是将他们都丢弃。
被编译的包会被保存到$GOPATH/pkg目录下，
目录路径和src目录路径对应，
可执行程序被保存到$GOPATH/bin目录。
（很多用户会将$GOPATH/bin添加到可执行程序的搜索列表中）
还有，go install 命令和 go build 命令都不会重新编译没有发生变化的包，
这可以收后续构建更快捷。
为了方便编译依赖的包，go build -i 命令将安装每个目标锁依赖的包。

因为编译对应不同的操作系统平台和CPU架构，
go install 命令会将编译结果安装到GOOS和GOARCH对应的目录。
例如，买Mac子系统，golang.org/x/net/html包将被安装到 $GOPATH/pkg/darwin_amd64目录下的golong.org/x/net/html.a文件。

针对不同操作系统或CPU的交叉构架也是很简单的。
只需要设置好目标对应的GOOS和GOARCH，
然后运行构建命令即可。
下面交叉编译的程序将输出它在编译时操作系统和CPU类型：

示例代码ch10/cross

下面以64位和32位环境分别执行程序：
$ go build gopl.io/ch10/cross
$ ./cross
darwin amd64
$ GOARCH=386 go build gopl.io/ch10/cross
$ ./cross
darwin 386

有些包可能需要针对不同平台和处理器类型使用不同版本的代码文件，
以便于处理底层的可移植性问题或提供为一些特定代码提供优化。
如果一个文件名包含了一个操作系统或处理器类型名字，
例如net_lunux.go或asm_amd64.s ， 
Go语言的构建工具将只在对应的平台编译这些文件。
还有一个特别的构建注释可以提供更多的构建过程控制。
例如，文件中可能包含下面的注释：
// +build linux darwin 建立linux达尔文

在包声明和包注释的前面，该构键注释参数告诉 go build 只在编译程序对应的目标操作系统是Linux或者Max OS X时才编译这个文件。
下面的构键注释则表示不编译这个文件：
// +build ignore

更多细节，可以参考go/build包的构建约束部分的文档。
$ go doc go/build

10.7.4 包文档
Go语言的编码风格鼓励为每个包提供良好的文档。
包中每个导出的成员和包声明前都应该包含目的和用法说明的注释。

GO语言中包文档注释一般是完整的句子，
第一行是包的摘要说明，注释后紧跟着包声明语句。
注释中函数的参数或其他的标识符并不需要额外的引号或其他标记注明。
例如，下面是fmt.Fprintf的文档注释。
// Fprintf formats according to a format specifier and writers to w.
// It returns the number of bytes written and any write error envountered.
func Fprintf(w io.Writer, format string, a ...interface{})(int, error)

Fprint函数格式化的细节在fmt包文档中描述。
如果朱时候紧跟着包声明语句，
那注释对应整个包的文档。
包文档对应的注释只能有一个（译注：其实可以有多个，
他们会组合成一个包文档注释），
包注释可以出现在任何一个源文件中。
如果包的注释内容比较长，
一般会放到一个独立的源文件中；
fmt包注释就有300行之多。
这个专门用于保存包文档的源文件通常叫doc.go。

好的文档并不需要面面俱到，文档本身应该是简洁但不可忽略的。
事实上，Go语言的风格更喜欢简洁的文档，
并且文档也是需要向代码一样维护的。
对于一组声明语句，
可以用一个精炼的句子描述，
如果是显而易见的功能则并不需要注释。

在本书中，只要空间允许，
我们之前很多包声明都包含了注释文档，
但你可以从标准库中发现很多更好的例子。
有两个工具可以帮到你。

首先是go doc命令，该命令打印包的声明和每个成员的文档注释，
下面是整个包的文档：

$ go goc time
package time // import "time"

Package time provides functionality for measuring and displaying time.package main


const Nanosecond Duration = 1 ...
func After(d Duration) <- 1 ...
func Sleep(d Duration) <- chan Time
func Since(t Time) Duration
func Now() Time
type Duration int64
type Time struct{ ... }
...many more... 

或者是某个具体包成员的注释文档：
$ go goc time.Since
func Since(t Time) Duration

	Since returns the time elapsed since t.package main
	It is shorthand for time.Now().Sub(t) .

或者是某个具体包的一个方法的注释文档：
$ go doc time.Duration.Seconds
func (d Duration) Seconds() float64

	Seconds returns the duration as a floating-point number of seconds.

该命令并不需要输入完整的包导入路径或正确的大小写。
下面的命令将打印encoding/json包的 (*json.Decoder).Decode方法的文档：
$ go doc json.decode
func (dec *Decoder) Decode(v interface{}) error
	
	Decode reads the next JSON-encoded value from its input and stores
	it in the value pointed to by v.

第二个工具，名字也叫godoc，它提供可以相互交叉引用的HTML页面，
但是包含和go doc命令相同以及更多的信息。
10.1节演示了time包的文档，
11.6节将看到godoc演示可以交互的示例程序。
godoc的在线服务https://godoc.org，
包含了成千上万的开源包的检索工具。

你也可以在自己的工作区目录运行godoc服务。
运行下面的命令，然后在浏览器查看http://localhost:8000/pkg页面：
$ godoc -http :8000
其中 -analysis=type 和 -analysis=pointer 命令行标志参数用于打开文档和代码中关于静态分析的结果。

10.7.5 内部包（p392）
在Go语言程序中，包的封装机制是一个重要的特性。
没有导出的标识符只在同一个包内部可以访问，
而导出的标识符则是面向全宇宙都可见的。

有时候，一个中间的状态可能也是有用的，
对于一小部分信任的包是可见的，
但并不是对所有调用者都可见。
例如，当我们计划将一个大的包拆分为很多小的更容易维护的子包，
但是我们并不想将内部的子包结构也完全暴露出去。
同时，我们可能还希望在内部子包之间共享一些通用的处理包，
或者我们只是想实验一个新包的还并不稳定的接口，
暂时只暴露给一些受限制的用户使用。

为了满足这些需求，Go语言的构建工具对包含internal名字的路径的包导入路径做了特殊处理。
这种包叫internal包，
一个internal包只能被和internal目录有同一个父目录的包所导入。
例如，net/http/internal/chunked内部包只能被net/http/httputil或net/http包导入，
但是不能被net/url包导入。
不过net/url包却可以导入net/http/httputil包。 //他们都在同一个父包

net/http
net/http/internal/chunked
net/http/httputil
net/url

10.7.6 查询包（p393）
go list 命令可以查询可用包的信息。
其最简单的形式，可以测试包是否在工作区并打印他的导入路径：
$ go list github.com/go-sql-driver/mysql
github.com/go-sql-driver/mysql

go list 命令的参数还可以用 "..." 表示匹配任意的包的导入路径。
我们可以用它来列表工作区中的所有包：
$ go list ...
archive/tar
archive/zip
bufio
bytes
cmd/addr2line
cmd/api
...many more... 

或者是特定子目录下的所有包：
$ go list gopl.io/ch3...
gopl.io/ch3/basename1
gopl.io/ch3/basename2
gopl.io/ch3/comma
gopl.io/ch3/mandelbrot
gopl.io/ch3/netflag
gopl.io/ch3/printints
gopl.io/ch3/surface

或者是和某个主题相关的所有包：
$ go list ...xml... 
encoding/xml
gopl.io/ch7/xmlselect

go list 命令还可以获取每个包完整的元信息，
而不仅仅只是导入路径，
这些元信息可以以不同格式提供给用户。
其中 -json 命令行参数表示用JSON格式打印每个包的元信息。

$ go list -json hash
{
"Dir": "/home/gopher/go/src/hash",
"ImportPath": "hash",
"Name": "hash",
"Doc": "Package hash provides interfaces for hash functions.",
"Target": "/home/gopher/go/pkg/darwin_amd64/hash.a",
"Goroot": true,
"Standard": true,
"Root": "/home/gopher/go",
"GoFiles": [
"hash.go"
],
"Imports": [
"io"
],
"Deps": [
"errors",
"io",
"runtime",
"sync",
"sync/atomic",
"unsafe"
]
}

命令行参数 -f 则允许用户使用text/template包（4.6）的模板语言定义输出文本的格式。
下面的命令将打印strconv包的依赖的包，
然后用join模板函数将结果链接为一行，
连接时每个结果之间用一个空格分隔：
$ go list -f '{{join .Deps " "}}' strconv
errors math runtime unicode/utf8 unsafe

译注：上面的命令在Windows的命令行运行会遇到 
template：main:1 unclosed action的错误。
产生这个错误的原因是因为命令行对命令中的 " " 参数进行了转义处理。
可以按照下面的方法解决转义字符串的问题：
$ go list -f "{{join .Deps \" \"}}" strconv
下面的命令打印compress子目录下所有包的依赖包列表：
$ go list -f '{{.ImportPath}} -> {{join .Imports " "}}' compress/... 
compress/bzip2 -> bufio io sort
compress/flate -> bufio fmt io math sort strconv
compress/gzip -> bufio compress/flate errors fmt hash hash/crc32 io time
compress/lzw -> bufio errors fmt io
compress/zlib -> bufio compress/flate errors fmt hash hash/adler32 io

译注：Windows下有同样有问题，要避免转义字符串的干扰：
$ go list -f "{{.ImportPath}} -> {{join .Imports \" \"}}" compress/...

go list 命令对于一次性的交互式查询或自动化构建或测试脚本都很有帮助。
我们将在11.2.4节中再次使用它。
每个子命令的更多信息，包括可设置的字段和意义，
可以用go help list命令查看。

在本章，我们解释了Go语言工具中除了测试命令之外的所有重要的子命令。
在下一章，我们将看到如何用go test命令去运行Go语言程序中的测试代码。
