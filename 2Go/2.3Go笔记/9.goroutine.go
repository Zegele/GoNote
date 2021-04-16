goroutine的定义（视频9-2）

goroutine是协程？

1.任何函数只需在前面加上go，就能送给调度器运行。
2.不需要在定义时区分是不是异步函数。（异步是什么东西？）
3.调度器在合适的点进行切换。
goroutine可能的切换的点

I/O，select   函数调用（有时）
channel		  runtime.Gosched()
等待锁//这是啥？
以上这些只是参考，不能保证切换，不能保证在其他地方不切换。



4.使用-race来检测数据访问冲突。
go run -race name//name是被检测文件的名称。


