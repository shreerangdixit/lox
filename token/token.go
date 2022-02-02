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
	if t.Type == TT_IDENTIFIER || t.Type == TT_NUMBER {
		return fmt.Sprintf("%s:%s", t.Type, t.Literal)
	}
	return t.Literal
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
	TT_IDENTIFIER = "IDENT"
	TT_NUMBER     = "NUM"

	// Operators
	TT_ASSIGN   = "="
	TT_PLUS     = "+"
	TT_MINUS    = "-"
	TT_DIVIDE   = "/"
	TT_MULTIPLY = "*"
	TT_NOT      = "!"
	TT_EQ       = "=="
	TT_NEQ      = "!="
	TT_LT       = "<"
	TT_LTE      = "<="
	TT_GT       = ">"
	TT_GTE      = ">="

	// Delimiters
	TT_COMMA     = ","
	TT_SEMICOLON = ";"

	// Parens + Braces
	TT_LPAREN = "("
	TT_RPAREN = ")"
	TT_LBRACE = "{"
	TT_RBRACE = "}"

	// Keywords
	TT_FUNCTION = "FUNC"
	TT_LET      = "LET"
	TT_IF       = "IF"
	TT_ELSE     = "ELSE"
	TT_TRUE     = "TRUE"
	TT_FALSE    = "FALSE"
	TT_RETURN   = "RETURN"
)
