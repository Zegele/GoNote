##操作环境：Win7 64位

## 前提: 
1. 已安装GO并设置好环境变量后。（例子：我的GOROOT=C:/Go  ；GOPATH=D:/gogo；**Path中已经添加%GOPATH%\bin**）
2. 如果是第二次安装并配置Sublime则需要把之前的Sublime文件删除干净，包括系统盘的有关sublime文件都删除干净，以免出现额外的问题。
3. 已安装Git

##1. 下载Sublime Text3

下载连接：[https://www.sublimetext.com](https://www.sublimetext.com/)


##2. 安装Package Control

按下：ctrl+`，调出控制台，复制相关代码，粘贴进控制台，按回车键。

相关代码（只适用于Sublime Text3）：

```
import urllib.request,os,hashlib; h = '6f4c264a24d933ce70df5dedcf1dcaee' + 'ebe013ee18cced0ef93d5f746d80ef60'; pf = 'Package Control.sublime-package'; ipp = sublime.installed_packages_path(); urllib.request.install_opener( urllib.request.build_opener( urllib.request.ProxyHandler()) ); by = urllib.request.urlopen( 'http://packagecontrol.io/' + pf.replace(' ', '%20')).read(); dh = hashlib.sha256(by).hexdigest(); print('Error validating download (got %s instead of %s), please try manual install' % (dh, h)) if dh != h else open(os.path.join( ipp, pf), 'wb' ).write(by)
```
参考：[https://packagecontrol.io/installation#st3](https://packagecontrol.io/installation#st3)

安装完毕后，关闭sublime，再打开，点击Preferences，出现Package Control则安装成功。

##3.安装插件

**3.1 GoSublime**
**（作用：补全代码、提示错误、提示改动等功能）**


下载地址：[https://github.com/DisposaBoy/GoSublime](https://github.com/DisposaBoy/GoSublime)

方法1： 用git clone命令获取，应该会获得一个GoSublime的文件夹。把该文件夹拷贝到sublime的包位置，也就是打开sublime - Preferences - Browse Packages...，将该文件粘贴进来）。

方法2：（建议使用方法1）下载压缩包，解压后会得到一个长名字的文件夹，把文件名改为GoSublime，再把该文件夹移动到安装sublime包的位置，也就是打开sublime - Preferences - Browse Packages...，将文件粘贴进该位置）

然后，关闭sublime，再打开sublime。（我发现：这时GOPATH的bin文件夹内生成一个margo.sublime可运行文件）

情况1：如果出现该情况，就是需要调整Margo文件。

![image](https://upload-images.jianshu.io/upload_images/4370290-3796c9853db736ff.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
（如果没出现这个情况则说明正常）

依次按快捷键 ctrl+.,ctrl+x调出margo.go
复制以下代码，并替换原有代码。

```

package margo

import (
    "margo.sh/golang"
    "margo.sh/mg"
    "time"
)

// Margo is the entry-point to margo
func Margo(m mg.Args) {
    // See the documentation for `mg.Reducer`
    // comments beginning with `gs:` denote features that replace old GoSublime settings

    // add our reducers (margo plugins) to the store
    // they are run in the specified order
    // and should ideally not block for more than a couple milliseconds
    m.Use(
        // MOTD keeps you updated about new versions and important announcements
        //
        // It adds a new command `motd.sync` available via the UserCmd palette as `Sync MOTD`
        //
        // Interval can be set in order to enable automatic update fetching.
        //
        // When new updates are found, it displays the message in the status bar
        // e.g. `★ margo.sh/cl/18.09.14 ★` a url where you see the upcoming changes before updating
        //
        // It sends the following data to the url https://api.margo.sh/motd.json:
        // * current editor plugin name e.g. `?client=gosublime`
        //   this tells us which editor plugin's changelog to check
        // * current editor plugin version e.g. `?tag=r18.09.14-1`
        //   this allows us to determine if there any updates
        // * whether or not this is the first request of the day e.g. `?firstHit=1`
        //   this allows us to get an estimated count of active users without storing
        //   any personally identifiable data
        //
        // No other data is sent. For more info contact privacy at kuroku.io
        //
        &mg.MOTD{
            // Interval, if set, specifies how often to automatically fetch messages from Endpoint
            // Interval: 3600e9, // automatically fetch updates every hour
        },

        mg.NewReducer(func(mx *mg.Ctx) *mg.State {
            // By default, events (e.g. ViewSaved) are triggered in all files.
            // Replace `mg.AllLangs` with `mg.Go` to restrict events to Go(-lang) files.
            // Please note, however, that this mode is not tested
            // and saving a non-go file will not trigger linters, etc. for that go pkg
            return mx.SetConfig(mx.Config.EnabledForLangs(
                mg.AllLangs,
            ))
        }),

        // Add `go` command integration
        // this adds a new commands:
        // gs: these commands are all callable through 9o:
        // * go: Wrapper around the go command, adding linter support
        // * go.play: Automatically build and run go commands or run go test for packages
        //   with support for linting and unsaved files
        // * go.replay: Wrapper around go.play limited to a single instance
        //   by default this command is bound to ctrl+.,ctrl+r or cmd+.,cmd+r
        //
        // UserCmds are also added for `Go Play` and `Go RePlay`
        &golang.GoCmd{},

        // add the day and time to the status bar
        &DayTimeStatus{},

        // both GoFmt and GoImports will automatically disable the GoSublime version
        // you will need to install the `goimports` tool manually
        // https://godoc.org/golang.org/x/tools/cmd/goimports
        //
        // gs: this replaces settings `fmt_enabled`, `fmt_tab_indent`, `fmt_tab_width`, `fmt_cmd`
        //
        // golang.GoFmt,
        // or
        // golang.GoImports,

        // Configure general auto-completion behaviour
        &golang.MarGocodeCtl{
            // whether or not to include Test*, Benchmark* and Example* functions in the auto-completion list
            // gs: this replaces the `autocomplete_tests` setting
            ProposeTests: false,

            // Don't try to automatically import packages when auto-compeltion fails
            // e.g. when `json.` is typed, if auto-complete fails
            // "encoding/json" is imported and auto-complete attempted on that package instead
            // See AddUnimportedPackages
            NoUnimportedPackages: false,

            // If a package was imported internally for use in auto-completion,
            // insert it in the source code
            // See NoUnimportedPackages
            // e.g. after `json.` is typed, `import "encoding/json"` added to the code
            AddUnimportedPackages: false,

            // Don't preload packages to speed up auto-completion, etc.
            NoPreloading: false,

            // Don't suggest builtin types and functions
            // gs: this replaces the `autocomplete_builtins` setting
            NoBuiltins: false,
        },

        // Enable auto-completion
        // gs: this replaces the `gscomplete_enabled` setting
        &golang.Gocode{
            // show the function parameters. this can take up a lot of space
            ShowFuncParams: true,
        },

        // show func arguments/calltips in the status bar
        // gs: this replaces the `calltips` setting
        &golang.GocodeCalltips{},

        // use guru for goto-definition
        // new commands `goto.definition` and `guru.definition` are defined
        // gs: by default `goto.definition` is bound to ctrl+.,ctrl+g or cmd+.,cmd+g
        &golang.Guru{},

        // add some default context aware-ish snippets
        // gs: this replaces the `autocomplete_snippets` and `default_snippets` settings
        golang.Snippets,

        // add our own snippets
        // gs: this replaces the `snippets` setting
        MySnippets,

        // check the file for syntax errors
        // gs: this and other linters e.g. below,
        //     replaces the settings `gslint_enabled`, `lint_filter`, `comp_lint_enabled`,
        //     `comp_lint_commands`, `gslint_timeout`, `lint_enabled`, `linters`
        &golang.SyntaxCheck{},

        // Add user commands for running tests and benchmarks
        // gs: this adds support for the tests command palette `ctrl+.`,`ctrl+t` or `cmd+.`,`cmd+t`
        &golang.TestCmds{
            // additional args to add to the command when running tests and examples
            TestArgs: []string{},

            // additional args to add to the command when running benchmarks
            BenchArgs: []string{"-benchmem"},
        },

        // run `go install -i` on save
        // golang.GoInstall("-i"),
        // or
        // golang.GoInstallDiscardBinaries("-i"),
        //
        // GoInstallDiscardBinaries will additionally set $GOBIN
        // to a temp directory so binaries are not installed into your $GOPATH/bin
        //
        // the -i flag is used to install imported packages as well
        // it's only supported in go1.10 or newer

        // run `go vet` on save. go vet is ran automatically as part of `go test` in go1.10
        // golang.GoVet(),

        // run `go test -race` on save
        // golang.GoTest("-race"),

        // run `golint` on save
        // &golang.Linter{Name: "golint", Label: "Go/Lint"},

        // run gometalinter on save
        // &golang.Linter{Name: "gometalinter", Args: []string{
        //  "--disable=gas",
        //  "--fast",
        // }},
    )
}

// DayTimeStatus adds the current day and time to the status bar
type DayTimeStatus struct {
    mg.ReducerType
}

func (dts DayTimeStatus) RMount(mx *mg.Ctx) {
    // kick off the ticker when we start
    dispatch := mx.Store.Dispatch
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for range ticker.C {
            dispatch(mg.Render)
        }
    }()
}

func (dts DayTimeStatus) Reduce(mx *mg.Ctx) *mg.State {
    // we always want to render the time
    // otherwise it will sometimes disappear from the status bar
    now := time.Now()
    format := "Mon, 15:04"
    if now.Second()%2 == 0 {
        format = "Mon, 15 04"
    }
    return mx.AddStatus(now.Format(format))
}

// MySnippets is a slice of functions returning our own snippets
var MySnippets = golang.SnippetFuncs(
    func(cx *golang.CompletionCtx) []mg.Completion {
        // if we're not in a block (i.e. function), do nothing
        if !cx.Scope.Is(golang.BlockScope) {
            return nil
        }

        return []mg.Completion{
            {
                Query: "if err",
                Title: "err != nil { return }",
                Src:   "if ${1:err} != nil {\n\treturn $0\n}",
            },
        }
    },
)
```
关闭sublime后，重启。
没有出现以上情况则说明安装成功。

参考：[https://www.jianshu.com/p/1aaaf18f4adc](https://www.jianshu.com/p/1aaaf18f4adc)

在设置中添加以下代码：
```
{
	"env":{
		"GOPATH":"E:/gogo",
		"GOROOT":"C:/Go",
	}
}
```
（我没设置以上代码貌似也没出现什么问题。前提是之前的go安装和设置环境变量，我已经设置过了，应该sublime的默认设置会识别到。）

**3.2 Goimports（它是一座大坑，爬过这个大坑有巨大的成就感。不把编辑器弄舒服以后怎么好好写代码呢？！所以希望小白的我们可以通过这个过程磨炼下耐心。）**
**作用：自动添加包，非常有用！**

1. 由于网络原因无法直接下载goimports，需要先下载gopm包。
打开git Bash Here，输入如下代码：

```
go get -v github.com/gpmgo/gopm
```

安装完成后，你会发现在：
1.1 GOPATH\src文件夹中看到github文件夹，以及很多文件。
1.2 GOPATH\bin文件夹中看到gopm的运行文件。

2. 下载goimports
在git bash here中输入以下命令
```swift
gopm get -g -v -u golang.org/x/tools/cmd/goimports
```
安装完成后GOPATH\src文件夹中看到golang.org文件夹，里面有很多文件。

如果出现安装不了相应的包，或者显示网络没有链接，则都说明不正确。

出现这个为安装正常
![下载goimports](https://upload-images.jianshu.io/upload_images/4370290-6b3bbdf8a1bbd95e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
（当然一切的目标是下载到goimports包，如果你有办法下到包即可。可以参考链接：gopm.io，下载你需要的包，当然也有golang.org/x/tools包）

/*给新机安装时发现gomp下载不了。。。
用以下方法反而快速拥有goimports
1.创建文件夹  
mkdir $GOPATH/src/golang.org/x/ 
2. 进入文件夹 cd $GOPATH/src/golang.org/x/ 
3. 下载源码 git clone https://github.com/golang/tools.git 
4. 安装 go install golang.org/x/tools/cmd/goimports
安装完成后，在GOPATH/bin中你会看到goimports的可运行文件。

参考文章：https://blog.csdn.net/weixin_30709061/article/details/98612978
*/

3. 安装goimports

```swift
go install golang.org/x/tools/cmd/goimports
```

安装完成后，在GOPATH/bin中你会看到goimports的可运行文件。

4.测试运行
打开sublime，调出marge.go（快捷键按ctrl+. ctrl+x），按照图片进行修改。

![golang.GoImports](https://upload-images.jianshu.io/upload_images/4370290-a61c6e4495d29736.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

保存修改后，关闭。
再打开一个go文件，点击保存（ctrl+s），可以看到import中会自动添加调用的包，并且会调整代码格式，这真是个非常好的插件。
（灵感来自与这篇文章，《Go语言环境配置 Sublime Text + GoSublime+ gocode + MarGo组》，看来后隐约觉得margo.go很重要，所以抱着在margo里改改代码试试，结果还真成功了。虽然这篇文章的内容我没怎么看，但是给了我一个灵感，让我想到再去margo试试，这个很重要。可以仔细读下margo中关于goimports的描述。）

如果你的没反应，再试试：打开Sublime-Preferences-Package Settings-GoSublime-Settings-User，输入以下代码
```
{
    "fmt_cmd": ["goimports"],
    "env": {
        "GOPATH": "/your/gopath/here" //输入你的GOPATH，如D:/gogo
    }
}
```
（以上设置我也没用到，如果你的goimports没有起作用可以试试。）

参考链接：（最后貌似没用上，但你也需要多查看资料，多试试，应该得这样）

1. .bash_profile 文件改坏了，什么命令都用不了了（[https://blog.csdn.net/qq_16177481/article/details/55518267](https://blog.csdn.net/qq_16177481/article/details/55518267)
）
（这个是查错的资料）

2. go环境变量配置 (GOROOT和GOPATH)
（[https://www.jianshu.com/p/4e699ff478a5](https://www.jianshu.com/p/4e699ff478a5)
）

3. goimports
[https://blog.csdn.net/geyujiao0828/article/details/89472658](https://blog.csdn.net/geyujiao0828/article/details/89472658)


**3.3 Gotests**
打开sublime，Preferences-Package Control-Install Package-输入Gotests，点击安装。
（研究好了再补充这个test怎么用）