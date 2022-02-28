package lex

var Keywords = map[string]TokenType{
	"var":      TT_VAR,
	"fun":      TT_FUNCTION,
	"if":       TT_IF,
	"else":     TT_ELSE,
	"true":     TT_TRUE,
	"false":    TT_FALSE,
	"return":   TT_RETURN,
	"while":    TT_WHILE,
	"nil":      TT_NIL,
	"break":    TT_BREAK,
	"continue": TT_CONTINUE,
	"defer":    TT_DEFER,
	"assert":   TT_ASSERT,
}
