package strext

import "strings"

func IsNullOrWhiteSpace(str string) bool {
	// 使用 strings.TrimSpace() 去除字符串首尾的空白字符
	trimmed := strings.TrimSpace(str)
	// 检查剩余字符串的长度是否为0
	return len(trimmed) == 0
}
