package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/lex"
)

type SyntaxError struct {
	err string
	tok lex.Token
}

func NewSyntaxError(err string, tok lex.Token) SyntaxError {
	return SyntaxError{
		err: err,
		tok: tok,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.tok)
}
