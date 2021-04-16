package popcount

//pc[i] is the population count of i.
var pc [256]byte //256是什么意思？

func init() {
	for i := range pc {//注意：range循环只使用了索引，省略了没有用到的值的部分。也可以这样写： for i, _ := range pc {}
		pc[i] = pc[i/2] + byte(i&1)
	}
	}
}

//init函数用于不是简单的初始化，而是对循环或计算，等较复杂的对象进行初始化。

/*
对于上面的复杂处理的初始化，也可以通过将初始化逻辑包装为一个匿名函数处理。
如下：
var pc [256]byte = func() (pc [256]byte){
	for i := range pc{
		pc[i] = pc[i/2]+byte(i&1)
	}
	return
}()


*/

//PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])

}

//pc表格用于处理每个8bit宽度的数字含二进制的1bit的bit个数。