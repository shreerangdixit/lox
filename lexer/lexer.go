package lexer

type Lexer struct {
	input      string
	readPos    int
	currentPos int
	ch         byte
	line       int
	col        int
	tokBegin   Position
	tokEnd     Position
}

func New(input string) *Lexer {
	lexer := Lexer{
		input: input,
		line:  1,
		col:   0,
	}
	lexer.advance()
	return &lexer
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	var tok Token
	switch l.ch {
	case '=':
		l.tokenBegin()
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_EQ, literal)
		} else {
			tok = newToken(TT_ASSIGN, string(l.ch))
		}
		l.tokenEnd()
	case '!':
		l.tokenBegin()
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_NEQ, literal)
		} else {
			tok = newToken(TT_NOT, string(l.ch))
		}
		l.tokenEnd()
	case '<':
		l.tokenBegin()
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_LTE, literal)
		} else {
			tok = newToken(TT_LT, string(l.ch))
		}
		l.tokenEnd()
	case '>':
		l.tokenBegin()
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_GTE, literal)
		} else {
			tok = newToken(TT_GT, string(l.ch))
		}
		l.tokenEnd()
	case '&':
		l.tokenBegin()
		if l.peek() == '&' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_LOGICAL_AND, literal)
		} else {
			tok = newToken(TT_ILLEGAL, string(l.ch))
		}
		l.tokenEnd()
	case '|':
		l.tokenBegin()
		if l.peek() == '|' {
			ch := l.ch
			l.advance()
			literal := string(ch) + string(l.ch)
			tok = newToken(TT_LOGICAL_OR, literal)
		} else {
			tok = newToken(TT_ILLEGAL, string(l.ch))
		}
		l.tokenEnd()
	case '+':
		l.tokenBegin()
		tok = newToken(TT_PLUS, string(l.ch))
		l.tokenEnd()
	case '-':
		l.tokenBegin()
		tok = newToken(TT_MINUS, string(l.ch))
		l.tokenEnd()
	case '/':
		l.tokenBegin()
		if l.peek() == '/' {
			tok = l.readCommentToken()
		} else {
			tok = newToken(TT_DIVIDE, string(l.ch))
		}
		l.tokenEnd()
	case '*':
		l.tokenBegin()
		tok = newToken(TT_MULTIPLY, string(l.ch))
		l.tokenEnd()
	case '%':
		l.tokenBegin()
		tok = newToken(TT_MODULO, string(l.ch))
		l.tokenEnd()
	case ',':
		l.tokenBegin()
		tok = newToken(TT_COMMA, string(l.ch))
		l.tokenEnd()
	case '(':
		l.tokenBegin()
		tok = newToken(TT_LPAREN, string(l.ch))
		l.tokenEnd()
	case ')':
		l.tokenBegin()
		tok = newToken(TT_RPAREN, string(l.ch))
		l.tokenEnd()
	case '{':
		l.tokenBegin()
		tok = newToken(TT_LBRACE, string(l.ch))
		l.tokenEnd()
	case '}':
		l.tokenBegin()
		tok = newToken(TT_RBRACE, string(l.ch))
		l.tokenEnd()
	case '[':
		l.tokenBegin()
		tok = newToken(TT_LBRACKET, string(l.ch))
		l.tokenEnd()
	case ']':
		l.tokenBegin()
		tok = newToken(TT_RBRACKET, string(l.ch))
		l.tokenEnd()
	case '?':
		l.tokenBegin()
		tok = newToken(TT_QUESTION, string(l.ch))
		l.tokenEnd()
	case ':':
		l.tokenBegin()
		tok = newToken(TT_COLON, string(l.ch))
		l.tokenEnd()
	case 0:
		tok = newToken(TT_EOF, "0")
	default:
		if isDigit(l.ch) {
			l.tokenBegin()
			tok = l.readNumberToken()
			l.tokenEnd()
		} else if isLetter(l.ch) {
			l.tokenBegin()
			tok = l.readIdentifierToken()
			l.tokenEnd()
		} else if l.ch == '"' {
			l.tokenBegin()
			tok = l.readStringToken()
			l.tokenEnd()
		} else {
			l.tokenBegin()
			tok = newToken(TT_ILLEGAL, string(l.ch))
			l.tokenEnd()
		}
	}
	l.advance()
	tok.BeginPosition = l.tokBegin
	tok.EndPosition = l.tokEnd
	return tok
}

func (l *Lexer) tokenBegin() {
	l.tokBegin = Position{
		Line:   l.line,
		Column: l.col,
	}
}

func (l *Lexer) tokenEnd() {
	l.tokEnd = Position{
		Line:   l.line,
		Column: l.col,
	}
}

func (l *Lexer) advance() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.currentPos = l.readPos
	l.readPos += 1
	l.col += 1
}

func (l *Lexer) rewind() {
	l.currentPos -= 1
	l.readPos -= 1
	l.col -= 1
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) readNumberToken() Token {
	defer l.rewind()
	startPos := l.currentPos
	decimalCount := 0
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			decimalCount += 1
		}
		if decimalCount > 1 {
			return newToken(TT_ILLEGAL, string(l.ch))
		}
		l.advance()
	}
	numStr := l.input[startPos:l.currentPos]
	return newToken(TT_NUMBER, numStr)
}

func (l *Lexer) readIdentifierToken() Token {
	defer l.rewind()
	startPos := l.currentPos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	value := l.input[startPos:l.currentPos]
	return newToken(LookupIdentifierType(value), value)
}

func (l *Lexer) readStringToken() Token {
	l.advance()
	startPos := l.currentPos
	for l.ch != '"' {
		l.advance()
	}
	value := l.input[startPos:l.currentPos]
	return newToken(TT_STRING, value)
}

func (l *Lexer) readCommentToken() Token {
	defer l.rewind()
	startPos := l.currentPos
	for !isNewline(l.ch) && l.ch != 0 {
		l.advance()
	}
	value := l.input[startPos:l.currentPos]
	return newToken(TT_COMMENT, value)
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		if isNewline(l.ch) {
			l.lineFeed()
		}
		l.advance()
	}
}

func (l *Lexer) lineFeed() {
	l.col = 0
	l.line += 1
}

func newToken(tokenType TokenType, literal string) Token {
	return Token{
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

func isNewline(b byte) bool {
	return b == '\n' || b == '\r'
}
