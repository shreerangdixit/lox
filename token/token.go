package token

import (
	"fmt"
)

type TokenType int

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
	TT_ILLEGAL TokenType = iota
	TT_EOF

	// Identifier + Literals
	TT_IDENTIFIER
	TT_NUMBER

	// Operators
	TT_ASSIGN
	TT_PLUS
	TT_MINUS
	TT_DIVIDE
	TT_MULTIPLY
	TT_NOT
	TT_EQ
	TT_NEQ
	TT_LT
	TT_LTE
	TT_GT
	TT_GTE

	// Delimiters
	TT_COMMA
	TT_SEMICOLON

	// Parens + Braces
	TT_LPAREN
	TT_RPAREN
	TT_LBRACE
	TT_RBRACE

	// Keywords
	TT_FUNCTION
	TT_LET
	TT_IF
	TT_ELSE
	TT_TRUE
	TT_FALSE
	TT_RETURN
)

func (t TokenType) String() string {
	switch t {
	case TT_ILLEGAL:
		return "ILLEGAL"
	case TT_EOF:
		return "EOF"
	case TT_IDENTIFIER:
		return "IDENT"
	case TT_NUMBER:
		return "NUMBER"
	case TT_ASSIGN:
		return "="
	case TT_PLUS:
		return "+"
	case TT_MINUS:
		return "-"
	case TT_DIVIDE:
		return "/"
	case TT_MULTIPLY:
		return "*"
	case TT_NOT:
		return "!"
	case TT_EQ:
		return "=="
	case TT_NEQ:
		return "!="
	case TT_LT:
		return "<"
	case TT_LTE:
		return "<="
	case TT_GT:
		return ">"
	case TT_GTE:
		return ">="
	case TT_COMMA:
		return ","
	case TT_SEMICOLON:
		return ","
	case TT_LPAREN:
		return "("
	case TT_RPAREN:
		return ")"
	case TT_LBRACE:
		return "{"
	case TT_RBRACE:
		return "}"
	case TT_FUNCTION:
		return "fun"
	case TT_LET:
		return "let"
	case TT_IF:
		return "if"
	case TT_ELSE:
		return "else"
	case TT_TRUE:
		return "true"
	case TT_FALSE:
		return "false"
	case TT_RETURN:
		return "return"
	default:
		return "<UNKNOWN>"
	}
}
