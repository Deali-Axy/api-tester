package goext

import "reflect"

// NullCoalescing null coalescing operator
//
// if the first value is null, use this other value.
func NullCoalescing[T any](firstVal T, otherValue T) T {
	if reflect.ValueOf(firstVal).IsNil() {
		return otherValue
	}

	return firstVal
}
