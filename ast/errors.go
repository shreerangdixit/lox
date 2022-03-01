package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/lex"
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
