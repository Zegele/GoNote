// package treesort provides insertion sort using an unbalanced binary tree.

package treesort

type tree struct {
	value       int
	left, right *tree
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

//没看懂。。。。。

func add(t *tree, value int) *tree {
	if t == nil {
		//Equivalent to return &tree{value: value}.
		t = new(tree) //new函数就是相当于创建了指针。查看2.3.3
		//创建了tree类型的匿名变量，
		//初始化为tree的零值，然后返回变量地址，返回的指针类型为*tree。
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}
