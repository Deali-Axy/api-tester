package goext

import (
	"errors"
	"fmt"
)

type CatchContext struct {
	err *error
}

type FinallyContext struct {
}

func Try(f func() error) *CatchContext {
	err := f()
	return &CatchContext{
		err: &err,
	}
}

func (c CatchContext) Catch(f func(error)) *FinallyContext {
	f(*c.err)
	return &FinallyContext{}
}

func (c FinallyContext) Finally(f func()) {
	f()
}

func tryCatchExample() {
	// The native methods in Go language.
	if err := func() error {
		return errors.New("test error")
	}(); err != nil {
		fmt.Println("catch error: " + err.Error())
	}

	// The try-catch approach.
	Try(func() error {
		return errors.New("test error")
	}).Catch(func(err error) {
		fmt.Println("catch error: " + err.Error())
	}).Finally(func() {
		fmt.Println("finally execute!")
	})
}
