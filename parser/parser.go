package parser

import (
	"fmt"
	"lox/lexer"
	"lox/token"
)

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

type Node interface {
}

type BinaryOpNode struct {
	Left  Node
	Token token.Token
	Right Node
}

func (n BinaryOpNode) String() string {
	return fmt.Sprintf("BinaryOp[Left:%s Token:%s Right:%s]", n.Left, n.Token, n.Right)
}

type UnaryOpNode struct {
	Token token.Token
	Node  Node
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("UnaryOp[Token: %s Node:%s]", n.Token, n.Node)
}

type NumberNode struct {
	Token token.Token
}

func (n NumberNode) String() string {
	return fmt.Sprintf("Number[%s]", n.Token)
}

func (p *Parser) Expression() Node {
	return p.equality()
}

func (p *Parser) equality() Node {
	left := p.comparison()
	for p.next.Type == token.TT_EQ || p.next.Type == token.TT_NEQ {
		p.advance()
		tok := p.curr
		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: p.comparison(),
		}
	}
	return left
}

func (p *Parser) comparison() Node {
	left := p.term()
	for p.next.Type == token.TT_LT || p.next.Type == token.TT_LTE || p.next.Type == token.TT_GT || p.next.Type == token.TT_GTE {
		p.advance()
		tok := p.curr
		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: p.term(),
		}
	}
	return left
}

func (p *Parser) term() Node {
	left := p.factor()
	for p.next.Type == token.TT_PLUS || p.next.Type == token.TT_MINUS {
		p.advance()
		tok := p.curr
		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: p.factor(),
		}
	}
	return left
}

func (p *Parser) factor() Node {
	left := p.unary()
	for p.next.Type == token.TT_DIVIDE || p.next.Type == token.TT_MULTIPLY {
		p.advance()
		tok := p.curr
		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: p.unary(),
		}
	}
	return left
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
		return p.primary()
	}
	return node
}

func (p *Parser) primary() Node {
	if p.next.Type == token.TT_NUMBER {
		p.advance()
		return NumberNode{
			Token: p.curr,
		}
	}
	return nil
}

func (p *Parser) advance() {
	if p.curr.Type != token.TT_EOF {
		p.prev = p.curr
		p.curr = p.next
		p.next = p.lex.NextToken()
	}
}
