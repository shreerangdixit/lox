package ast

import (
	"fmt"

	"github.com/shreerangdixit/yeti/lex"
)

type SyntaxError struct {
	Err   string
	Token lex.Token
}

func NewSyntaxError(err string, tok lex.Token) SyntaxError {
	return SyntaxError{
		Err:   err,
		Token: tok,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Token)
}

func (e SyntaxError) Inner() error {
	return nil
}

func (e SyntaxError) ErrorType() string {
	return "syntax"
}

func (e SyntaxError) Begin() lex.Position {
	return e.Token.BeginPosition
}

func (e SyntaxError) End() lex.Position {
	return e.Token.EndPosition
}
