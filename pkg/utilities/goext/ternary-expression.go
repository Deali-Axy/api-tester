package goext

// If 三元表达式
//
// condition ? trueVal : falseVal
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
