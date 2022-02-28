package lex

import (
	"fmt"
)

type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("Pos %d:%d", p.Line, p.Column)
}
