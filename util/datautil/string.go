package datautil

import "unicode"

// IsLegalUserID 功能: 这是主要的校验函数。它接收一个字符串 str，并返回一个布尔值。
// 如果字符串中的所有字符都是字母、数字或下划线，则返回 true，否则返回 false。
func IsLegalUserID(str string) bool {
	for _, r := range str {
		if !IsAlphanumeric(r) && r != '_' {
			return false
		}
	}
	return true
}

// IsAlphanumeric 这是一个辅助函数。它接收一个 rune 类型的字符，判断该字符是否为字母或数字。
func IsAlphanumeric(b rune) bool {
	if !unicode.IsLetter(b) && !unicode.IsDigit(b) {
		return false
	}
	return true
}
