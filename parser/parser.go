package parser

import (
	"fmt"
	"github.com/shreerangdixit/lox/lexer"
	"github.com/shreerangdixit/lox/token"
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

type AssignmentNode struct {
	Identifier IdentifierNode
	Value      Node
}

type LetStmtNode struct {
	Identifier IdentifierNode
	Value      Node
}

type ExpStmtNode struct {
	Exp Node
}

type IfStmtNode struct {
	Exp       Node
	TrueStmt  Node
	FalseStmt Node
}

type PrintStmtNode struct {
	Exp Node
}

type WhileStmtNode struct {
	Condition Node
	Body      Node
}

type BlockNode struct {
	Declarations []Node
}

type ExpNode struct {
	Exp Node
}

type TernaryOpNode struct {
	Exp      Node
	TrueExp  Node
	FalseExp Node
}

type BinaryOpNode struct {
	LeftExp  Node
	Op       token.Token
	RightExp Node
}

type UnaryOpNode struct {
	Op      token.Token
	Operand Node
}

type LogicalAndNode struct {
	LHS Node
	RHS Node
}

type LogicalOrNode struct {
	LHS Node
	RHS Node
}

type BooleanNode struct {
	Token token.Token
}

type NumberNode struct {
	Token token.Token
}

type StringNode struct {
	Token token.Token
}

func (n NilNode) String() string        { return "nil" }
func (n ProgramNode) String() string    { return fmt.Sprintf("+%s", n.Declarations) }
func (n IdentifierNode) String() string { return fmt.Sprintf("%s", n.Token) }
func (n AssignmentNode) String() string { return fmt.Sprintf("%s=%s", n.Identifier, n.Value) }
func (n LetStmtNode) String() string    { return fmt.Sprintf("let %s=%s", n.Identifier, n.Value) }
func (n BlockNode) String() string      { return fmt.Sprintf("{%+s}", n.Declarations) }
func (n ExpStmtNode) String() string    { return fmt.Sprintf("%s", n.Exp) }
func (n IfStmtNode) String() string {
	return fmt.Sprintf("if(%s) %s else %s", n.Exp, n.TrueStmt, n.FalseStmt)
}
func (n PrintStmtNode) String() string { return fmt.Sprintf("%s", n.Exp) }
func (n WhileStmtNode) String() string { return fmt.Sprintf("while(%s)%s", n.Condition, n.Body) }
func (n ExpNode) String() string       { return fmt.Sprintf("%s", n.Exp) }
func (n TernaryOpNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", n.Exp, n.TrueExp, n.FalseExp)
}
func (n LogicalAndNode) String() string { return fmt.Sprintf("%s && %s", n.LHS, n.RHS) }
func (n LogicalOrNode) String() string  { return fmt.Sprintf("%s || %s", n.LHS, n.RHS) }
func (n BinaryOpNode) String() string   { return fmt.Sprintf("%s %s %s", n.LeftExp, n.Op, n.RightExp) }
func (n UnaryOpNode) String() string    { return fmt.Sprintf("%s%s", n.Op, n.Operand) }
func (n BooleanNode) String() string    { return fmt.Sprintf("%s", n.Token) }
func (n NumberNode) String() string     { return fmt.Sprintf("%s", n.Token) }
func (n StringNode) String() string     { return fmt.Sprintf("%s", n.Token) }

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

// program -> declaration* EOF ;
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

// declaration -> letDecl
//             | statement ;
func (p *Parser) declaration() (Node, error) {
	if p.consume(token.TT_LET) {
		return p.letDeclaration()
	}
	return p.statement()
}

// letDecl -> "let" IDENTIFIER ( "=" expression )? ";" ;
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

		return LetStmtNode{
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

	return LetStmtNode{
		Identifier: identifier,
		Value:      value,
	}, nil
}

// statement -> exprStatement
//           | ifStatement
//           | printStatement
//           | block ;
func (p *Parser) statement() (Node, error) {
	if p.consume(token.TT_IF) {
		return p.ifStatement()
	} else if p.consume(token.TT_PRINT) {
		return p.printStatement()
	} else if p.consume(token.TT_WHILE) {
		return p.whileStatement()
	} else if p.consume(token.TT_LBRACE) {
		return p.block()
	}
	return p.exprStatement()
}

// ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() (Node, error) {
	if !p.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for if condition", p.curr)
	}

	condExp, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for if condition", p.curr)
	}

	trueStmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	var falseStmt Node = nil
	if p.consume(token.TT_ELSE) {
		falseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return IfStmtNode{
		Exp:       condExp,
		TrueStmt:  trueStmt,
		FalseStmt: falseStmt,
	}, nil
}

// printStatement -> "print" expression ";" ;
func (p *Parser) printStatement() (Node, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.consume(token.TT_SEMICOLON) {
		return PrintStmtNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of a print statement", p.curr)
}

// whileStatement -> "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() (Node, error) {
	if !p.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for 'while' condition", p.curr)
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for 'while' condition", p.curr)
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return WhileStmtNode{
		Condition: condition,
		Body:      body,
	}, nil
}

// block -> "{" declaration* "}" ;
func (p *Parser) block() (Node, error) {
	declarations := make([]Node, 0, 100)

	for !p.check(token.TT_RBRACE) && !p.check(token.TT_EOF) {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}

	if !p.consume(token.TT_RBRACE) {
		return nil, newSyntaxError("expected closing '}'", p.curr)
	}

	return BlockNode{
		Declarations: declarations,
	}, nil
}

// exprStatement -> expression ";" ;
func (p *Parser) exprStatement() (Node, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if p.consume(token.TT_SEMICOLON) {
		return ExpStmtNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of an expression statement", p.curr)
}

// expression -> assignment ( "?" assignment ":" assignment )? ;
func (p *Parser) expression() (Node, error) {
	exp, err := p.assignment()
	if err != nil {
		return nil, err
	}

	// Check ternary operator: <assignment> ? <assignment> : <assignment>
	if p.consume(token.TT_QUESTION) {
		trueExp, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if !p.consume(token.TT_COLON) {
			return nil, newSyntaxError("expected ':'", p.curr)
		}

		falseExp, err := p.assignment()
		if err != nil {
			return nil, err
		}

		return TernaryOpNode{
			Exp:      exp,
			TrueExp:  trueExp,
			FalseExp: falseExp,
		}, nil
	}

	return exp, nil
}

// assignment -> IDENTIFIER "=" assignment
//            | logical_or ;
func (p *Parser) assignment() (Node, error) {
	expr, err := p.logical_or()
	if err != nil {
		return nil, err
	}

	if p.consume(token.TT_ASSIGN) {
		if _, ok := expr.(IdentifierNode); !ok {
			return nil, newSyntaxError("expected an identifier for assignment", p.curr)
		}

		assign, err := p.assignment()
		if err != nil {
			return nil, err
		}

		return AssignmentNode{
			Identifier: expr.(IdentifierNode),
			Value:      assign,
		}, nil
	}
	return expr, nil
}

// logical_or -> logical_and ( "||" logical_and )*
func (p *Parser) logical_or() (Node, error) {
	left, err := p.logical_and()
	if err != nil {
		return nil, err
	}

	for p.consume(token.TT_LOGICAL_OR) {
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		left = LogicalOrNode{
			LHS: left,
			RHS: right,
		}
	}
	return left, nil
}

// logical_and -> equality ( "&&" equality )* ;
func (p *Parser) logical_and() (Node, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.consume(token.TT_LOGICAL_AND) {
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		left = LogicalAndNode{
			LHS: left,
			RHS: right,
		}
	}
	return left, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_EQ, token.TT_NEQ}, p.comparison)
}

// comparison -> term ( ( "<" | "<=" | ">" | ">=" ) term )* ;
func (p *Parser) comparison() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_LT, token.TT_LTE, token.TT_GT, token.TT_GTE}, p.term)
}

// term -> factor ( ( "+" | "-" ) factor )* ;
func (p *Parser) term() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_PLUS, token.TT_MINUS}, p.factor)
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Node, error) {
	return p.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY}, p.unary)
}

// unary -> ( "!" | "-" ) unary
//       | atom ;
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

// atom -> NUMBER | STRING | "true" | "false" | "nil"
//      | "(" expression ")"
//      | IDENTIFIER ;
func (p *Parser) atom() (Node, error) {
	if p.consume(token.TT_NUMBER) {
		return NumberNode{Token: p.curr}, nil
	} else if p.consume(token.TT_STRING) {
		return StringNode{Token: p.curr}, nil
	} else if p.consumeAny([]token.TokenType{token.TT_TRUE, token.TT_FALSE}) {
		return BooleanNode{Token: p.curr}, nil
	} else if p.consume(token.TT_IDENTIFIER) {
		return IdentifierNode{Token: p.curr}, nil
	} else if p.consume(token.TT_NIL) {
		return NilNode{}, nil
	} else if p.consume(token.TT_LPAREN) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}

		if p.consume(token.TT_RPAREN) {
			return ExpNode{Exp: exp}, nil
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
			LeftExp:  left,
			Op:       tok,
			RightExp: right,
		}
	}
	return left, nil
}

// check checks the next token if it matches the given type and returns true, otherwise it returns false
func (p *Parser) check(tokType token.TokenType) bool {
	return p.checkAny([]token.TokenType{tokType})
}

// checkAny checks the next token if it matches any of the given types and returns true, otherwise it returns false
func (p *Parser) checkAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if p.next.Type == straw {
			return true
		}
	}
	return false
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
