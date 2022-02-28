package lex

import (
	"fmt"
)

type TokenType int

type Token struct {
	Type          TokenType
	Literal       string
	BeginPosition Position
	EndPosition   Position
}

func (t Token) String() string {
	if t.Type == TT_IDENTIFIER || t.Type == TT_NUMBER {
		return fmt.Sprintf("%s:%s (%s - %s)", t.Type, t.Literal, t.BeginPosition, t.EndPosition)
	}
	return fmt.Sprintf("%s (%s - %s)", t.Literal, t.BeginPosition, t.EndPosition)
}
