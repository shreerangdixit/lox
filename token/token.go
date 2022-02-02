package token

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("[%s:%s]", t.Type, t.Literal)
}

var keywords = map[string]TokenType{
	"let":    TT_LET,
	"fun":    TT_FUNCTION,
	"if":     TT_IF,
	"else":   TT_ELSE,
	"true":   TT_TRUE,
	"false":  TT_FALSE,
	"return": TT_RETURN,
}

func LookupIdentifierType(v string) TokenType {
	if val, ok := keywords[v]; ok {
		return val
	}
	return TT_IDENTIFIER
}

const (
	TT_ILLEGAL = "ILLEGAL"
	TT_EOF     = "EOF"

	// Identifier + Literals
	TT_IDENTIFIER = "TT_IDENTIFIER"
	TT_NUMBER     = "TT_NUMBER"

	// Operators
	TT_ASSIGN   = "TT_ASSIGN"
	TT_PLUS     = "TT_PLUS"
	TT_MINUS    = "TT_MINUS"
	TT_DIVIDE   = "TT_DIVIDE"
	TT_MULTIPLY = "TT_MULTIPLY"
	TT_NOT      = "TT_NOT"
	TT_EQ       = "TT_EQ"
	TT_NEQ      = "TT_NEQ"
	TT_LT       = "TT_LT"
	TT_LTE      = "TT_LTE"
	TT_GT       = "TT_GT"
	TT_GTE      = "TT_GTE"

	// Delimiters
	TT_COMMA     = "TT_COMMA"
	TT_SEMICOLON = "TT_SEMICOLON"

	// Parens + Braces
	TT_LPAREN = "TT_LPAREN"
	TT_RPAREN = "TT_RPAREN"
	TT_LBRACE = "TT_LBRACE"
	TT_RBRACE = "TT_RBRACE"

	// Keywords
	TT_FUNCTION = "TT_FUNCTION"
	TT_LET      = "TT_LET"
	TT_IF       = "TT_IF"
	TT_ELSE     = "TT_ELSE"
	TT_TRUE     = "TT_TRUE"
	TT_FALSE    = "TT_FALSE"
	TT_RETURN   = "TT_RETURN"
)
