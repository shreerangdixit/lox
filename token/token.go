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
	return fmt.Sprintf("[Type:%s Literal:%s]", t.Type, t.Literal)
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
	TT_INTEGER    = "TT_INTEGER"
	TT_FLOAT      = "TT_FLOAT"

	// Operators
	TT_ASSIGN          = "TT_ASSIGN"
	TT_PLUS            = "TT_PLUS"
	TT_MINUS           = "TT_MINUS"
	TT_DIVIDE          = "TT_DIVIDE"
	TT_MULTIPLY        = "TT_MULTIPLY"
	TT_EQUALITY        = "TT_EQUALITY"
	TT_NOT             = "TT_NOT"
	TT_NOT_EQUAL       = "TT_NOT_EQUAL"
	TT_LESS_THAN       = "TT_LESS_THAN"
	TT_LESS_THAN_EQ    = "TT_LESS_THAN_EQ"
	TT_GREATER_THAN    = "TT_GREATER_THAN"
	TT_GREATER_THAN_EQ = "TT_GREATER_THAN_EQ"

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
