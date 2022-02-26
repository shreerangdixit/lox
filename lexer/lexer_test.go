package lexer

import (
	"testing"

	"github.com/shreerangdixit/lox/token"
	"github.com/stretchr/testify/assert"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token.Token
	}{
		{
			name:  "operators_paren",
			input: "= / + - * , ( ) { } == ! != < <= > >= && ||",
			want: []token.Token{
				token.Token{Type: token.TT_ASSIGN, Literal: "="},
				token.Token{Type: token.TT_DIVIDE, Literal: "/"},
				token.Token{Type: token.TT_PLUS, Literal: "+"},
				token.Token{Type: token.TT_MINUS, Literal: "-"},
				token.Token{Type: token.TT_MULTIPLY, Literal: "*"},
				token.Token{Type: token.TT_COMMA, Literal: ","},
				token.Token{Type: token.TT_LPAREN, Literal: "("},
				token.Token{Type: token.TT_RPAREN, Literal: ")"},
				token.Token{Type: token.TT_LBRACE, Literal: "{"},
				token.Token{Type: token.TT_RBRACE, Literal: "}"},
				token.Token{Type: token.TT_EQ, Literal: "=="},
				token.Token{Type: token.TT_NOT, Literal: "!"},
				token.Token{Type: token.TT_NEQ, Literal: "!="},
				token.Token{Type: token.TT_LT, Literal: "<"},
				token.Token{Type: token.TT_LTE, Literal: "<="},
				token.Token{Type: token.TT_GT, Literal: ">"},
				token.Token{Type: token.TT_GTE, Literal: ">="},
				token.Token{Type: token.TT_LOGICAL_AND, Literal: "&&"},
				token.Token{Type: token.TT_LOGICAL_OR, Literal: "||"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "integers",
			input: "123 456 7890",
			want: []token.Token{
				token.Token{Type: token.TT_NUMBER, Literal: "123"},
				token.Token{Type: token.TT_NUMBER, Literal: "456"},
				token.Token{Type: token.TT_NUMBER, Literal: "7890"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "floats",
			input: "0.123 1.23",
			want: []token.Token{
				token.Token{Type: token.TT_NUMBER, Literal: "0.123"},
				token.Token{Type: token.TT_NUMBER, Literal: "1.23"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "bad_floats",
			input: ".123 1.23",
			want: []token.Token{
				token.Token{Type: token.TT_ILLEGAL, Literal: "."},
				token.Token{Type: token.TT_NUMBER, Literal: "123"},
				token.Token{Type: token.TT_NUMBER, Literal: "1.23"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "identifiers",
			input: "X Y Z aa bb cc_c d",
			want: []token.Token{
				token.Token{Type: token.TT_IDENTIFIER, Literal: "X"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "Y"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "Z"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "aa"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "bb"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "cc_c"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "d"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "mixed",
			input: " {(a = b * 5) (c = 10.5 / z)} ",
			want: []token.Token{
				token.Token{Type: token.TT_LBRACE, Literal: "{"},
				token.Token{Type: token.TT_LPAREN, Literal: "("},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "a"},
				token.Token{Type: token.TT_ASSIGN, Literal: "="},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "b"},
				token.Token{Type: token.TT_MULTIPLY, Literal: "*"},
				token.Token{Type: token.TT_NUMBER, Literal: "5"},
				token.Token{Type: token.TT_RPAREN, Literal: ")"},
				token.Token{Type: token.TT_LPAREN, Literal: "("},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "c"},
				token.Token{Type: token.TT_ASSIGN, Literal: "="},
				token.Token{Type: token.TT_NUMBER, Literal: "10.5"},
				token.Token{Type: token.TT_DIVIDE, Literal: "/"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "z"},
				token.Token{Type: token.TT_RPAREN, Literal: ")"},
				token.Token{Type: token.TT_RBRACE, Literal: "}"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "keywords",
			input: "var x = 10 y = fun foo(){} if else true false return",
			want: []token.Token{
				token.Token{Type: token.TT_VAR, Literal: "var"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "x"},
				token.Token{Type: token.TT_ASSIGN, Literal: "="},
				token.Token{Type: token.TT_NUMBER, Literal: "10"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "y"},
				token.Token{Type: token.TT_ASSIGN, Literal: "="},
				token.Token{Type: token.TT_FUNCTION, Literal: "fun"},
				token.Token{Type: token.TT_IDENTIFIER, Literal: "foo"},
				token.Token{Type: token.TT_LPAREN, Literal: "("},
				token.Token{Type: token.TT_RPAREN, Literal: ")"},
				token.Token{Type: token.TT_LBRACE, Literal: "{"},
				token.Token{Type: token.TT_RBRACE, Literal: "}"},
				token.Token{Type: token.TT_IF, Literal: "if"},
				token.Token{Type: token.TT_ELSE, Literal: "else"},
				token.Token{Type: token.TT_TRUE, Literal: "true"},
				token.Token{Type: token.TT_FALSE, Literal: "false"},
				token.Token{Type: token.TT_RETURN, Literal: "return"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "strings",
			input: "\"foo\" \"bar\" \"foo bar\"",
			want: []token.Token{
				token.Token{Type: token.TT_STRING, Literal: "foo"},
				token.Token{Type: token.TT_STRING, Literal: "bar"},
				token.Token{Type: token.TT_STRING, Literal: "foo bar"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "comments",
			input: "// my very very long comment",
			want: []token.Token{
				token.Token{Type: token.TT_COMMENT, Literal: "// my very very long comment"},
				token.Token{Type: token.TT_EOF, Literal: "0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for _, want_tok := range tt.want {
				got_tok := l.NextToken()
				assert.Equal(t, want_tok, got_tok)
			}
		})
	}
}
