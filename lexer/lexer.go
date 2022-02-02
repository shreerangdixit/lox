package lexer

import (
	"lox/token"
)

type Lexer struct {
	input      string
	readPos    int
	currentPos int
	ch         byte
}

func New(input string) *Lexer {
	lexer := Lexer{
		input: input,
	}
	lexer.advance()
	return &lexer
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	var tok token.Token
	switch l.ch {
	case '=':
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.TT_EQ, literal)
		} else {
			tok = newToken(token.TT_ASSIGN, string(l.ch))
		}
	case '!':
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.TT_NEQ, literal)
		} else {
			tok = newToken(token.TT_NOT, string(l.ch))
		}
	case '<':
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.TT_LTE, literal)
		} else {
			tok = newToken(token.TT_LT, string(l.ch))
		}
	case '>':
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(token.TT_GTE, literal)
		} else {
			tok = newToken(token.TT_GT, string(l.ch))
		}
	case '+':
		tok = newToken(token.TT_PLUS, string(l.ch))
	case '-':
		tok = newToken(token.TT_MINUS, string(l.ch))
	case '/':
		tok = newToken(token.TT_DIVIDE, string(l.ch))
	case '*':
		tok = newToken(token.TT_MULTIPLY, string(l.ch))
	case ',':
		tok = newToken(token.TT_COMMA, string(l.ch))
	case ';':
		tok = newToken(token.TT_SEMICOLON, string(l.ch))
	case '(':
		tok = newToken(token.TT_LPAREN, string(l.ch))
	case ')':
		tok = newToken(token.TT_RPAREN, string(l.ch))
	case '{':
		tok = newToken(token.TT_LBRACE, string(l.ch))
	case '}':
		tok = newToken(token.TT_RBRACE, string(l.ch))
	case 0:
		tok = newToken(token.TT_EOF, "0")
	default:
		if isDigit(l.ch) {
			tok = l.readNumberToken()
		} else if isLetter(l.ch) {
			tok = l.readIdentifierToken()
		} else {
			tok = newToken(token.TT_ILLEGAL, string(l.ch))
		}
	}
	l.advance()
	return tok
}

func (l *Lexer) advance() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.currentPos = l.readPos
	l.readPos += 1
}

func (l *Lexer) rewind() {
	l.currentPos -= 1
	l.readPos -= 1
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) readNumberToken() token.Token {
	defer l.rewind()
	startPos := l.currentPos
	decimalCount := 0
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			decimalCount += 1
		}
		if decimalCount > 1 {
			return newToken(token.TT_ILLEGAL, string(l.ch))
		}
		l.advance()
	}
	numStr := l.input[startPos:l.currentPos]
	return newToken(token.TT_NUMBER, numStr)
}

func (l *Lexer) readIdentifierToken() token.Token {
	defer l.rewind()
	startPos := l.currentPos
	for isLetter(l.ch) {
		l.advance()
	}
	value := l.input[startPos:l.currentPos]
	return newToken(token.LookupIdentifierType(value), value)
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.advance()
	}
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_'
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}
