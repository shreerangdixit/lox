package evaluate

import (
	"fmt"

	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/lex"
)

type EvalError struct {
	Node ast.Node
	Err  error
}

func NewEvalError(node ast.Node, err error) EvalError {
	return EvalError{
		Node: node,
		Err:  err,
	}
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e EvalError) ErrorType() string {
	return "runtime"
}

func (e EvalError) Begin() lex.Position {
	return e.Node.Begin()
}

func (e EvalError) End() lex.Position {
	return e.Node.End()
}

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
