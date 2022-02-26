package evaluator

import (
	"fmt"

	"github.com/shreerangdixit/lox/ast"
)

// Loop control flow exit due to `break;`
type BreakError struct{}

func NewBreakError() error {
	return BreakError{}
}

func (e BreakError) Error() string { return "break" }

// Loop control continue
type ContinueError struct{}

func NewContinueError() error {
	return ContinueError{}
}

func (e ContinueError) Error() string { return "continue" }

// Function control flow exit due to `return <exp>;`
type ReturnError struct {
	Value Object
}

func NewReturnError(value Object) error {
	return ReturnError{
		Value: value,
	}
}

func (e ReturnError) Error() string { return "return" }

// Assertion error
type AssertError struct {
	Exp ast.Node
}

func NewAssertError(exp ast.Node) error {
	return AssertError{
		Exp: exp,
	}
}

func (e AssertError) Error() string {
	return fmt.Sprintf("assert failed: %s", e.Exp)
}
