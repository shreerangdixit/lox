package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/token"
)

type SyntaxError struct {
	err string
	tok token.Token
}

func NewSyntaxError(err string, tok token.Token) SyntaxError {
	return SyntaxError{
		err: err,
		tok: tok,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.tok)
}
