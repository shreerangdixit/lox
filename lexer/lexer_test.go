package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "operators_paren",
			input: "= / + - * , ( ) { } == ! != < <= > >= && ||",
			want: []Token{
				Token{Type: TT_ASSIGN, Literal: "="},
				Token{Type: TT_DIVIDE, Literal: "/"},
				Token{Type: TT_PLUS, Literal: "+"},
				Token{Type: TT_MINUS, Literal: "-"},
				Token{Type: TT_MULTIPLY, Literal: "*"},
				Token{Type: TT_COMMA, Literal: ","},
				Token{Type: TT_LPAREN, Literal: "("},
				Token{Type: TT_RPAREN, Literal: ")"},
				Token{Type: TT_LBRACE, Literal: "{"},
				Token{Type: TT_RBRACE, Literal: "}"},
				Token{Type: TT_EQ, Literal: "=="},
				Token{Type: TT_NOT, Literal: "!"},
				Token{Type: TT_NEQ, Literal: "!="},
				Token{Type: TT_LT, Literal: "<"},
				Token{Type: TT_LTE, Literal: "<="},
				Token{Type: TT_GT, Literal: ">"},
				Token{Type: TT_GTE, Literal: ">="},
				Token{Type: TT_LOGICAL_AND, Literal: "&&"},
				Token{Type: TT_LOGICAL_OR, Literal: "||"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "integers",
			input: "123 456 7890",
			want: []Token{
				Token{Type: TT_NUMBER, Literal: "123"},
				Token{Type: TT_NUMBER, Literal: "456"},
				Token{Type: TT_NUMBER, Literal: "7890"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "floats",
			input: "0.123 1.23",
			want: []Token{
				Token{Type: TT_NUMBER, Literal: "0.123"},
				Token{Type: TT_NUMBER, Literal: "1.23"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "bad_floats",
			input: ".123 1.23",
			want: []Token{
				Token{Type: TT_ILLEGAL, Literal: "."},
				Token{Type: TT_NUMBER, Literal: "123"},
				Token{Type: TT_NUMBER, Literal: "1.23"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "identifiers",
			input: "X Y Z aa bb cc_c d",
			want: []Token{
				Token{Type: TT_IDENTIFIER, Literal: "X"},
				Token{Type: TT_IDENTIFIER, Literal: "Y"},
				Token{Type: TT_IDENTIFIER, Literal: "Z"},
				Token{Type: TT_IDENTIFIER, Literal: "aa"},
				Token{Type: TT_IDENTIFIER, Literal: "bb"},
				Token{Type: TT_IDENTIFIER, Literal: "cc_c"},
				Token{Type: TT_IDENTIFIER, Literal: "d"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "mixed",
			input: " {(a = b * 5) (c = 10.5 / z)} ",
			want: []Token{
				Token{Type: TT_LBRACE, Literal: "{"},
				Token{Type: TT_LPAREN, Literal: "("},
				Token{Type: TT_IDENTIFIER, Literal: "a"},
				Token{Type: TT_ASSIGN, Literal: "="},
				Token{Type: TT_IDENTIFIER, Literal: "b"},
				Token{Type: TT_MULTIPLY, Literal: "*"},
				Token{Type: TT_NUMBER, Literal: "5"},
				Token{Type: TT_RPAREN, Literal: ")"},
				Token{Type: TT_LPAREN, Literal: "("},
				Token{Type: TT_IDENTIFIER, Literal: "c"},
				Token{Type: TT_ASSIGN, Literal: "="},
				Token{Type: TT_NUMBER, Literal: "10.5"},
				Token{Type: TT_DIVIDE, Literal: "/"},
				Token{Type: TT_IDENTIFIER, Literal: "z"},
				Token{Type: TT_RPAREN, Literal: ")"},
				Token{Type: TT_RBRACE, Literal: "}"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "keywords",
			input: "var x = 10 y = fun foo(){} if else true false return",
			want: []Token{
				Token{Type: TT_VAR, Literal: "var"},
				Token{Type: TT_IDENTIFIER, Literal: "x"},
				Token{Type: TT_ASSIGN, Literal: "="},
				Token{Type: TT_NUMBER, Literal: "10"},
				Token{Type: TT_IDENTIFIER, Literal: "y"},
				Token{Type: TT_ASSIGN, Literal: "="},
				Token{Type: TT_FUNCTION, Literal: "fun"},
				Token{Type: TT_IDENTIFIER, Literal: "foo"},
				Token{Type: TT_LPAREN, Literal: "("},
				Token{Type: TT_RPAREN, Literal: ")"},
				Token{Type: TT_LBRACE, Literal: "{"},
				Token{Type: TT_RBRACE, Literal: "}"},
				Token{Type: TT_IF, Literal: "if"},
				Token{Type: TT_ELSE, Literal: "else"},
				Token{Type: TT_TRUE, Literal: "true"},
				Token{Type: TT_FALSE, Literal: "false"},
				Token{Type: TT_RETURN, Literal: "return"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "strings",
			input: "\"foo\" \"bar\" \"foo bar\"",
			want: []Token{
				Token{Type: TT_STRING, Literal: "foo"},
				Token{Type: TT_STRING, Literal: "bar"},
				Token{Type: TT_STRING, Literal: "foo bar"},
				Token{Type: TT_EOF, Literal: "0"},
			},
		},
		{
			name:  "comments",
			input: "// my very very long comment",
			want: []Token{
				Token{Type: TT_COMMENT, Literal: "// my very very long comment"},
				Token{Type: TT_EOF, Literal: "0"},
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
