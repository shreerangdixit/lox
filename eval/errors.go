package eval

import (
	"fmt"

	"github.com/shreerangdixit/yeti/ast"
	"github.com/shreerangdixit/yeti/lex"
)

type Option func(e *EvaluateError)

func WithInnerError(inner error) Option {
	return func(e *EvaluateError) {
		e.inner = inner
	}
}

type EvaluateError struct {
	node  ast.Node
	err   error
	inner error
}

func NewEvaluateError(node ast.Node, err error, opts ...Option) EvaluateError {
	e := EvaluateError{
		node: node,
		err:  err,
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

func (e EvaluateError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func (e EvaluateError) Inner() error {
	return e.inner
}

func (e EvaluateError) ErrorType() string {
	return "runtime"
}

func (e EvaluateError) Begin() lex.Position {
	return e.node.Begin()
}

func (e EvaluateError) End() lex.Position {
	return e.node.End()
}

// Loop control flow exit due to `break;`
type BreakError struct{}

func NewBreakError() BreakError {
	return BreakError{}
}

func (e BreakError) Error() string { return "break" }

// Loop control continue
type ContinueError struct{}

func NewContinueError() ContinueError {
	return ContinueError{}
}

func (e ContinueError) Error() string { return "continue" }

// Function control flow exit due to `return <exp>;`
type ReturnError struct {
	Value Object
}

func NewReturnError(value Object) ReturnError {
	return ReturnError{
		Value: value,
	}
}

func (e ReturnError) Error() string { return "return" }

// Assertion error
type AssertError struct {
	Exp ast.Node
}

func NewAssertError(exp ast.Node) AssertError {
	return AssertError{
		Exp: exp,
	}
}

func (e AssertError) Error() string {
	return fmt.Sprintf("assert failed: %s", e.Exp)
}
