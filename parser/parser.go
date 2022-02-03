package parser

import (
	"fmt"
	"lox/lexer"
	"lox/token"
)

type GrammarRuleFunc func() (Node, error)

type SyntaxError struct {
	err string
	tok token.Token
}

func newSyntaxError(err string, tok token.Token) SyntaxError {
	return SyntaxError{
		err: err,
		tok: tok,
	}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.tok)
}

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

type NumberNode struct {
	Token token.Token
}

func (n NumberNode) String() string {
	return fmt.Sprintf("%s", n.Token)
}

type BooleanNode struct {
	Token token.Token
}

func (n BooleanNode) String() string {
	return fmt.Sprintf("%s", n.Token)
}

type ExpressionNode struct {
	Token token.Token
	Node  Node
}

func (n ExpressionNode) String() string {
	return fmt.Sprintf("[%s%s]", n.Token, n.Node)
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

func (p *Parser) Parse() (Node, error) {
	return p.expression()
}

// ------------------------------------
// Grammar rule functions
// ------------------------------------

func (p *Parser) expression() (Node, error) {
	return p.equality()
}

func (p *Parser) equality() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_EQ, token.TT_NEQ}, p.comparison)
}

func (p *Parser) comparison() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_LT, token.TT_LTE, token.TT_GT, token.TT_GTE}, p.term)
}

func (p *Parser) term() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_PLUS, token.TT_MINUS}, p.factor)
}

func (p *Parser) factor() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY}, p.unary)
}

func (p *Parser) unary() (Node, error) {
	var node Node = nil

	for p.next.Type == token.TT_NOT || p.next.Type == token.TT_MINUS {
		p.advance()
		tok := p.curr

		n, err := p.unary()
		if err != nil {
			return nil, err
		}

		node = UnaryOpNode{
			Token: tok,
			Node:  n,
		}
	}
	if node == nil {
		return p.atom()
	}
	return node, nil
}

func (p *Parser) atom() (Node, error) {
	if p.next.Type == token.TT_NUMBER {
		p.advance()
		return NumberNode{Token: p.curr}, nil
	} else if p.nextTokenMatches([]token.TokenType{token.TT_TRUE, token.TT_FALSE}) {
		p.advance()
		return BooleanNode{Token: p.curr}, nil
	} else if p.next.Type == token.TT_LPAREN {
		p.advance()

		exp, err := p.expression()
		if err != nil {
			return nil, err
		}

		if p.next.Type == token.TT_RPAREN {
			p.advance()
			return ExpressionNode{Token: p.curr, Node: exp}, nil
		}
		return nil, newSyntaxError("expected closing ')' after expression", p.curr)
	}
	return nil, newSyntaxError("expected a literal or an expression", p.curr)
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

func (p *Parser) binaryOp(tokenTypes []token.TokenType, fun GrammarRuleFunc) (Node, error) {
	left, err := fun()
	if err != nil {
		return nil, err
	}

	for p.nextTokenMatches(tokenTypes) {
		p.advance()
		tok := p.curr

		right, err := fun()
		if err != nil {
			return nil, err
		}

		left = BinaryOpNode{
			Left:  left,
			Token: tok,
			Right: right,
		}
	}
	return left, nil
}

func (p *Parser) nextTokenMatches(haystack []token.TokenType) bool {
	for _, straw := range haystack {
		if p.next.Type == straw {
			return true
		}
	}
	return false
}
