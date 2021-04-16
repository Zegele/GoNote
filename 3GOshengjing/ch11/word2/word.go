// Package word provides utilities for word games.
package word

import (
	"unicode"
)

// IsPalindrome reports whether s reads the same forward and backward.
// Letter case is ignored, as are non-letters.
func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) { //参考https://www.cnblogs.com/golove/p/3273585.html
			//判断r是否为字母,汉字也是一个字母字符
			letters = append(letters, unicode.ToLower(r)) //转化为小写格式
		}
	}
	for i := range letters {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}

//!-
