//!+test
package word

import "testing"

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") { //IsPalindrome("detartrated")的结果是true，
		//If !true {...} 意思是如果不是true，则执行下面
		//if只执行true的？也就是 if true会被永远执行。对的！！！
		t.Error(`IsPalindrome("1detartrated1") = false`)
	}
	if !IsPalindrome("kayak") {
		t.Error(`IsPalindrome("kayak") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
}

//!-test

// The tests below are expected to fail.
// See package ch11/word2 for the fix.

//!+more
func TestFrenchPalindrome(t *testing.T) {
	if !IsPalindrome("été") {
		t.Error(`IsPalindrome(été) = false`)
	}
}

func TestCanalPalindrome(t *testing.T) {
	input := "A man, a plan, a canal: Panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = flase`, input)
	}
}

//!-more
