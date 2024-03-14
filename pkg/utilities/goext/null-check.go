package goext

import (
	"reflect"
)

func IfNull(val interface{}) bool {
	return val == nil
}

// NullCoalescing null coalescing operator
//
// if the first value is null, use this other value.
func NullCoalescing[T any](firstVal T, otherValue T) T {
	if reflect.ValueOf(firstVal).IsNil() {
		return otherValue
	}

	return firstVal
}

// NullCoalescing2 null coalescing operator
//
// The version implemented using interface{}.
func NullCoalescing2(firstVal interface{}, otherValue interface{}) interface{} {
	if firstVal == nil {
		return otherValue
	}

	return firstVal
}
