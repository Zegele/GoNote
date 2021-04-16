在git中测试代码。(在ccmouse视频第8章，我纳闷在sublime中怎么测试？我只能在git中操作。)
go 的测试代码需要test格式，例如：name_test.go    //_test.go是固定的，这样go test才会识别到。

go test
会显示代码是否通过。

go test -coverprofile=c.out  //显示代码覆盖的占比。
我目前也不知道显示了覆盖占比能干嘛。
**这个命令会生成一个c.out的文件。**

less c.out  //这个命令可以查看c.out的内容，但是不怎么看得懂。

go tool cover -html=c.out  //可以通过网页显示出来代码的覆盖情况。

go test -bench .  //显示性能，虽然我依然不知道是干啥的。“. ”代表所有文件，当然你也可以写具体的文件名。

go test -bench . -cpuprofile cpu.out
//生成一个cup.out文件(具体文件名是什么自己可以设定)

go tool pprof cpu.out
//查看cpu.out 会有个交互式的命令，选择你想查看的方式等。
然后输入web，就弹出一个图片，可以查看哪里花费了较长时间运行。
（但是我输入没有出现。。。先放着。
我遇到这个问题：Failed to execute dot. Is Graphviz installed? Error: exec: "dot": executable file not found in$PATH
解决方式，安装Graphviz 
安装代码：sudo apt install graphviz

可参考：

1.[https://www.jianshu.com/p/ca1aec491c8c](https://www.jianshu.com/p/ca1aec491c8c)
2.[https://blog.csdn.net/weixin_42654444/article/details/82108055](https://blog.csdn.net/weixin_42654444/article/details/82108055)

如果web后报错 Failed to execute dot. Is Graphviz installed? Error: exec: "dot": executable file not found in %PATH%

是你电脑没有安装gvedit导致的

fq进入gvedit官网https://graphviz.gitlab.io/_pages/Download/Download_windows.html 下载稳定版

按照提示进行安装即可

安装完成后 设置环境变量path 后面加上gvedit安装路径的bin文件夹

这里是我的安装目录和path设置

![image.png](https://upload-images.jianshu.io/upload_images/4370290-9aab73ef99ec0e44.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

重启，报错解决
 ———————————————— 
版权声明：本文为CSDN博主「昨夜是今晨的开始」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/weixin_42654444/article/details/82108055
）