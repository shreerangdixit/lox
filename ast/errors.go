package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/lexer"
)

type SyntaxError struct {
	err string
	tok lexer.Token
}

func NewSyntaxError(err string, tok lexer.Token) SyntaxError {
	return SyntaxError{
		err: err,
		tok: tok,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.tok)
}
