package eval

import (
	"fmt"

	"github.com/shreerangdixit/redes/ast"
	"github.com/shreerangdixit/redes/lex"
)

type Option func(e *EvalError)

func WithInnerError(inner error) Option {
	return func(e *EvalError) {
		e.inner = inner
	}
}

type EvalError struct {
	node  ast.Node
	err   error
	inner error
}

func NewEvalError(node ast.Node, err error, opts ...Option) EvalError {
	e := EvalError{
		node: node,
		err:  err,
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func (e EvalError) Inner() error {
	return e.inner
}

func (e EvalError) ErrorType() string {
	return "runtime"
}

func (e EvalError) Begin() lex.Position {
	return e.node.Begin()
}

func (e EvalError) End() lex.Position {
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
