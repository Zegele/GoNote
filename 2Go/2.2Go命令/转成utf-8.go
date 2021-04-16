/*

方法有2个下载

1. 针对gbk转成utf-8的————text包
在git中输入：
gopm get -g -v golang.org/x/text

等待安装后，在GOPATH下的golang.org文件夹下，x文件夹下，会出现一个text文件夹。
另外，在user地址下的gomp--> Administrator--> .gomp-->repos-->golang.org-->x-->text文件夹


2.通用的转化为utf-8————net/html包
在git中输入：
gopm get -g -v golang.org/x/net/html

等待安装后，在GOPATH下的golang.org文件夹下，x文件夹下，会出现一个net文件夹。
另外，在user地址下的gomp--> Administrator--> .gomp-->repos-->golang.org-->x-->net文件夹


怎样查看是什么码？
有个charset是显示使用的码。
charset="utf-8">
charset="gbk"
如

*/