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

type Node interface{}

type NilNode struct {
}

type ProgramNode struct {
	Declarations []Node
}

type IdentifierNode struct {
	Token token.Token
}

type LetStatementNode struct {
	Identifier IdentifierNode
	Value      Node
}

type ExpressionStatementNode struct {
	Exp Node
}

type PrintStatementNode struct {
	Exp Node
}

type ExpressionNode struct {
	Exp Node
}

type BinaryOpNode struct {
	LHS Node
	Op  token.Token
	RHS Node
}

type UnaryOpNode struct {
	Op      token.Token
	Operand Node
}

type NumberNode struct {
	Token token.Token
}

type BooleanNode struct {
	Token token.Token
}

func (n NilNode) String() string                 { return "nil" }
func (n ProgramNode) String() string             { return fmt.Sprintf("+%s", n.Declarations) }
func (n IdentifierNode) String() string          { return fmt.Sprintf("%s", n.Token) }
func (n LetStatementNode) String() string        { return fmt.Sprintf("let %s=%s", n.Identifier, n.Value) }
func (n ExpressionStatementNode) String() string { return fmt.Sprintf("%s", n.Exp) }
func (n PrintStatementNode) String() string      { return fmt.Sprintf("%s", n.Exp) }
func (n ExpressionNode) String() string          { return fmt.Sprintf("%s", n.Exp) }
func (n BinaryOpNode) String() string            { return fmt.Sprintf("%s %s %s", n.LHS, n.Op, n.RHS) }
func (n UnaryOpNode) String() string             { return fmt.Sprintf("%s%s", n.Op, n.Operand) }
func (n NumberNode) String() string              { return fmt.Sprintf("%s", n.Token) }
func (n BooleanNode) String() string             { return fmt.Sprintf("%s", n.Token) }

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
	return p.program()
}

// ------------------------------------
// Grammar rule functions
// ------------------------------------
func (p *Parser) program() (Node, error) {
	declarations := make([]Node, 0, 100)
	for !p.consume(token.TT_EOF) {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}
	return ProgramNode{
		Declarations: declarations,
	}, nil
}

func (p *Parser) declaration() (Node, error) {
	if p.consume(token.TT_LET) {
		return p.letDeclaration()
	}
	return p.statement()
}

func (p *Parser) letDeclaration() (Node, error) {
	atom, err := p.atom()
	if err != nil {
		return nil, err
	}

	identifier, ok := atom.(IdentifierNode)
	if !ok {
		return nil, newSyntaxError("Expected identifier after let", p.curr)
	}

	if !p.consume(token.TT_ASSIGN) {
		if !p.consume(token.TT_SEMICOLON) {
			return nil, newSyntaxError("expected a ; at the end of a declaration", p.curr)
		}

		return LetStatementNode{
			Identifier: identifier,
			Value:      NilNode{},
		}, nil
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.consume(token.TT_SEMICOLON) {
		return nil, newSyntaxError("expected a ; at the end of a declaration", p.curr)
	}

	return LetStatementNode{
		Identifier: identifier,
		Value:      value,
	}, nil
}

func (p *Parser) statement() (Node, error) {
	if p.consume(token.TT_PRINT) {
		return p.printStatement()
	}
	return p.exprStatement()
}

func (p *Parser) printStatement() (Node, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.consume(token.TT_SEMICOLON) {
		return PrintStatementNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of a print statement", p.curr)
}

func (p *Parser) exprStatement() (Node, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.consume(token.TT_SEMICOLON) {
		return ExpressionStatementNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of an expression statement", p.curr)
}

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

	for p.consumeAny([]token.TokenType{token.TT_NOT, token.TT_MINUS}) {
		tok := p.curr

		n, err := p.unary()
		if err != nil {
			return nil, err
		}

		node = UnaryOpNode{
			Op:      tok,
			Operand: n,
		}
	}
	if node == nil {
		return p.atom()
	}
	return node, nil
}

func (p *Parser) atom() (Node, error) {
	if p.consume(token.TT_NUMBER) {
		return NumberNode{Token: p.curr}, nil
	} else if p.consumeAny([]token.TokenType{token.TT_TRUE, token.TT_FALSE}) {
		return BooleanNode{Token: p.curr}, nil
	} else if p.consumeAny([]token.TokenType{token.TT_IDENTIFIER}) {
		return IdentifierNode{Token: p.curr}, nil
	} else if p.consumeAny([]token.TokenType{token.TT_NIL}) {
		return NilNode{}, nil
	} else if p.consume(token.TT_LPAREN) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}

		if p.consume(token.TT_RPAREN) {
			return ExpressionNode{Exp: exp}, nil
		}
		return nil, newSyntaxError("expected closing ')' after expression", p.curr)
	}
	return nil, newSyntaxError("expected a literal or an expression", p.curr)
}

// ------------------------------------
// Helpers
// ------------------------------------

func (p *Parser) binaryOp(tokenTypes []token.TokenType, fun GrammarRuleFunc) (Node, error) {
	left, err := fun()
	if err != nil {
		return nil, err
	}

	for p.consumeAny(tokenTypes) {
		tok := p.curr

		right, err := fun()
		if err != nil {
			return nil, err
		}

		left = BinaryOpNode{
			LHS: left,
			Op:  tok,
			RHS: right,
		}
	}
	return left, nil
}

// consume consumes the next token if it matches the given type and returns true, otherwise it returns false
func (p *Parser) consume(tokType token.TokenType) bool {
	return p.consumeAny([]token.TokenType{tokType})
}

// consumeAny consumes the next token if it matches any of the given types and returns true, otherwise it returns false
func (p *Parser) consumeAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if p.next.Type == straw {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() {
	if p.curr.Type != token.TT_EOF {
		p.prev = p.curr
		p.curr = p.next
		p.next = p.lex.NextToken()
	}
}
