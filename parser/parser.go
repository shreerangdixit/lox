package parser

import (
	"fmt"
	"lox/lexer"
	"lox/token"
)

type GrammarRuleFunc func() Node

// ------------------------------------
// Nodes
// ------------------------------------

type Node interface {
}

type BinaryOpNode struct {
	Left  Node
	Token token.Token
	Right Node
}

func (n BinaryOpNode) String() string {
	return fmt.Sprintf("[%s %s %s]", n.Left, n.Token, n.Right)
}

type UnaryOpNode struct {
	Token token.Token
	Node  Node
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("[%s%s]", n.Token, n.Node)
}

type LiteralNode struct {
	Token token.Token
}

func (n LiteralNode) String() string {
	return fmt.Sprintf("%s", n.Token)
}

// ------------------------------------
// Parser
// ------------------------------------

type Parser struct {
	lex  *lexer.Lexer
	curr token.Token
	prev token.Token
	next token.Token
}

func New(lex *lexer.Lexer) *Parser {
	p := Parser{
		lex:  lex,
		curr: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		prev: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		next: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
	}
	p.advance()
	return &p
}

func (p *Parser) Parse() Node {
	return p.expression()
}

// ------------------------------------
// Grammar rule functions
// ------------------------------------

func (p *Parser) expression() Node {
	return p.equality()
}

func (p *Parser) equality() Node {
	return p.binaryOp([]token.TokenType{token.TT_EQ, token.TT_NEQ}, p.comparison)
}

func (p *Parser) comparison() Node {
	return p.binaryOp([]token.TokenType{token.TT_LT, token.TT_LTE, token.TT_GT, token.TT_GTE}, p.term)
}

func (p *Parser) term() Node {
	return p.binaryOp([]token.TokenType{token.TT_PLUS, token.TT_MINUS}, p.factor)
}

func (p *Parser) factor() Node {
	return p.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY}, p.unary)
}

func (p *Parser) unary() Node {
	var node Node = nil
	for p.next.Type == token.TT_NOT || p.next.Type == token.TT_MINUS {
		p.advance()
		tok := p.curr
		node = UnaryOpNode{
			Token: tok,
			Node:  p.unary(),
		}
	}
	if node == nil {
		return p.atom()
	}
	return node
}

func (p *Parser) atom() Node {
	if checkTokenType(p.next, []token.TokenType{token.TT_NUMBER, token.TT_TRUE, token.TT_FALSE}) {
		p.advance()
		return LiteralNode{
			Token: p.curr,
		}
	} else if p.next.Type == token.TT_LPAREN {
		p.advance()
		exp := p.expression()
		if p.next.Type == token.TT_RPAREN {
			p.advance()
			return exp
		}
		// Syntax error
		return nil
	}
	return nil
}

// ------------------------------------
// Helpers
// ------------------------------------

func (p *Parser) advance() {
	if p.curr.Type != token.TT_EOF {
		p.prev = p.curr
		p.curr = p.next
		p.next = p.lex.NextToken()
	}
}

func (p *Parser) binaryOp(tokenTypes []token.TokenType, fun GrammarRuleFunc) Node {
	left := fun()
	for checkTokenType(p.next, tokenTypes) {
		p.advance()
		tok := p.curr
		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: fun(),
		}
	}
	return left
}

func checkTokenType(needle token.Token, haystack []token.TokenType) bool {
	for _, straw := range haystack {
		if needle.Type == straw {
			return true
		}
	}
	return false
}
