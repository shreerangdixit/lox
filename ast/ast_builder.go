package ast

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

type CallNode struct {
	Callee    Node
	Arguments []Node
}

func (n NilNode) String() string        { return "nil" }
func (n ProgramNode) String() string    { return fmt.Sprintf("+%s", n.Declarations) }
func (n IdentifierNode) String() string { return n.Token.String() }
func (n AssignmentNode) String() string { return fmt.Sprintf("%s=%s", n.Identifier, n.Value) }
func (n LetStmtNode) String() string    { return fmt.Sprintf("let %s=%s", n.Identifier, n.Value) }
func (n BlockNode) String() string      { return fmt.Sprintf("{%s}", n.Declarations) }
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
func (n BooleanNode) String() string    { return n.Token.String() }
func (n NumberNode) String() string     { return n.Token.String() }
func (n StringNode) String() string     { return n.Token.String() }
func (n CallNode) String() string       { return fmt.Sprintf("func %s(%s)", n.Callee, n.Arguments) }

// ------------------------------------
// AstBuilder
// ------------------------------------

type AstBuilder struct {
	lex  *lexer.Lexer
	curr token.Token
	prev token.Token
	next token.Token
}

func New(lex *lexer.Lexer) *AstBuilder {
	b := AstBuilder{
		lex:  lex,
		curr: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		prev: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		next: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
	}
	b.advance()
	return &b
}

func (b *AstBuilder) RootNode() (Node, error) {
	return b.program()
}

// ------------------------------------
// Grammar rule functions
// ------------------------------------

// program -> declaration* EOF ;
func (b *AstBuilder) program() (Node, error) {
	declarations := make([]Node, 0, 100)
	for !b.consume(token.TT_EOF) {
		decl, err := b.declaration()
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
func (b *AstBuilder) declaration() (Node, error) {
	if b.consume(token.TT_LET) {
		return b.letDeclaration()
	}
	return b.statement()
}

// letDecl -> "let" IDENTIFIER ( "=" expression )? ";" ;
func (b *AstBuilder) letDeclaration() (Node, error) {
	atom, err := b.atom()
	if err != nil {
		return nil, err
	}

	identifier, ok := atom.(IdentifierNode)
	if !ok {
		return nil, newSyntaxError("Expected identifier after let", b.curr)
	}

	if !b.consume(token.TT_ASSIGN) {
		if !b.consume(token.TT_SEMICOLON) {
			return nil, newSyntaxError("expected a ; at the end of a declaration", b.curr)
		}

		return LetStmtNode{
			Identifier: identifier,
			Value:      NilNode{},
		}, nil
	}

	value, err := b.expression()
	if err != nil {
		return nil, err
	}

	if !b.consume(token.TT_SEMICOLON) {
		return nil, newSyntaxError("expected a ; at the end of a declaration", b.curr)
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
func (b *AstBuilder) statement() (Node, error) {
	if b.consume(token.TT_IF) {
		return b.ifStatement()
	} else if b.consume(token.TT_PRINT) {
		return b.printStatement()
	} else if b.consume(token.TT_WHILE) {
		return b.whileStatement()
	} else if b.consume(token.TT_LBRACE) {
		return b.block()
	}
	return b.exprStatement()
}

// ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
func (b *AstBuilder) ifStatement() (Node, error) {
	if !b.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for if condition", b.curr)
	}

	condExp, err := b.expression()
	if err != nil {
		return nil, err
	}

	if !b.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for if condition", b.curr)
	}

	trueStmt, err := b.statement()
	if err != nil {
		return nil, err
	}

	var falseStmt Node = nil
	if b.consume(token.TT_ELSE) {
		falseStmt, err = b.statement()
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
func (b *AstBuilder) printStatement() (Node, error) {
	expr, err := b.expression()
	if err != nil {
		return nil, err
	}

	if b.consume(token.TT_SEMICOLON) {
		return PrintStmtNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of a print statement", b.curr)
}

// whileStatement -> "while" "(" expression ")" statement ;
func (b *AstBuilder) whileStatement() (Node, error) {
	if !b.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for 'while' condition", b.curr)
	}

	condition, err := b.expression()
	if err != nil {
		return nil, err
	}

	if !b.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for 'while' condition", b.curr)
	}

	body, err := b.statement()
	if err != nil {
		return nil, err
	}

	return WhileStmtNode{
		Condition: condition,
		Body:      body,
	}, nil
}

// block -> "{" declaration* "}" ;
func (b *AstBuilder) block() (Node, error) {
	declarations := make([]Node, 0, 100)

	for !b.check(token.TT_RBRACE) && !b.check(token.TT_EOF) {
		decl, err := b.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}

	if !b.consume(token.TT_RBRACE) {
		return nil, newSyntaxError("expected closing '}'", b.curr)
	}

	return BlockNode{
		Declarations: declarations,
	}, nil
}

// exprStatement -> expression ";" ;
func (b *AstBuilder) exprStatement() (Node, error) {
	expr, err := b.expression()
	if err != nil {
		return nil, err
	}

	if b.consume(token.TT_SEMICOLON) {
		return ExpStmtNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of an expression statement", b.curr)
}

// expression -> assignment ( "?" assignment ":" assignment )? ;
func (b *AstBuilder) expression() (Node, error) {
	exp, err := b.assignment()
	if err != nil {
		return nil, err
	}

	// Check ternary operator: <assignment> ? <assignment> : <assignment>
	if b.consume(token.TT_QUESTION) {
		trueExp, err := b.assignment()
		if err != nil {
			return nil, err
		}

		if !b.consume(token.TT_COLON) {
			return nil, newSyntaxError("expected ':'", b.curr)
		}

		falseExp, err := b.assignment()
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
func (b *AstBuilder) assignment() (Node, error) {
	expr, err := b.logical_or()
	if err != nil {
		return nil, err
	}

	if b.consume(token.TT_ASSIGN) {
		if _, ok := expr.(IdentifierNode); !ok {
			return nil, newSyntaxError("expected an identifier for assignment", b.curr)
		}

		assign, err := b.assignment()
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
func (b *AstBuilder) logical_or() (Node, error) {
	left, err := b.logical_and()
	if err != nil {
		return nil, err
	}

	for b.consume(token.TT_LOGICAL_OR) {
		right, err := b.equality()
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
func (b *AstBuilder) logical_and() (Node, error) {
	left, err := b.equality()
	if err != nil {
		return nil, err
	}

	for b.consume(token.TT_LOGICAL_AND) {
		right, err := b.equality()
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
func (b *AstBuilder) equality() (Node, error) {
	return b.binaryOp([]token.TokenType{token.TT_EQ, token.TT_NEQ}, b.comparison)
}

// comparison -> term ( ( "<" | "<=" | ">" | ">=" ) term )* ;
func (b *AstBuilder) comparison() (Node, error) {
	return b.binaryOp([]token.TokenType{token.TT_LT, token.TT_LTE, token.TT_GT, token.TT_GTE}, b.term)
}

// term -> factor ( ( "+" | "-" ) factor )* ;
func (b *AstBuilder) term() (Node, error) {
	return b.binaryOp([]token.TokenType{token.TT_PLUS, token.TT_MINUS}, b.factor)
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (b *AstBuilder) factor() (Node, error) {
	return b.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY}, b.unary)
}

// unary -> ( "!" | "-" ) unary | call ;
func (b *AstBuilder) unary() (Node, error) {
	var node Node
	for b.consumeAny([]token.TokenType{token.TT_NOT, token.TT_MINUS}) {
		tok := b.curr

		n, err := b.unary()
		if err != nil {
			return nil, err
		}

		node = UnaryOpNode{
			Op:      tok,
			Operand: n,
		}
	}

	if node != nil {
		return node, nil
	}
	return b.call()
}

// call -> atom ( "(" arguments? ")" )* ;
func (b *AstBuilder) call() (Node, error) {
	expr, err := b.atom()
	if err != nil {
		return nil, err
	}

	for b.consume(token.TT_LPAREN) {
		expr, err = b.finishCall(expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (b *AstBuilder) finishCall(callee Node) (Node, error) {
	arguments := []Node{}
	var err error

	for !b.check(token.TT_RPAREN) {
		arguments, err = b.arguments()
		if err != nil {
			return nil, err
		}
	}

	if !b.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for function call", b.curr)
	}

	return CallNode{
		Callee:    callee,
		Arguments: arguments,
	}, nil
}

// arguments -> expression ( "," expression )* ;
func (b *AstBuilder) arguments() ([]Node, error) {
	arguments := make([]Node, 0, 255)

	arg, err := b.expression()
	if err != nil {
		return nil, err
	}

	arguments = append(arguments, arg)
	for b.consume(token.TT_COMMA) {
		arg, err := b.expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, arg)
	}

	return arguments, nil
}

// atom -> NUMBER | STRING | "true" | "false" | "nil"
//      | "(" expression ")"
//      | IDENTIFIER ;
func (b *AstBuilder) atom() (Node, error) {
	if b.consume(token.TT_NUMBER) {
		return NumberNode{Token: b.curr}, nil
	} else if b.consume(token.TT_STRING) {
		return StringNode{Token: b.curr}, nil
	} else if b.consumeAny([]token.TokenType{token.TT_TRUE, token.TT_FALSE}) {
		return BooleanNode{Token: b.curr}, nil
	} else if b.consume(token.TT_IDENTIFIER) {
		return IdentifierNode{Token: b.curr}, nil
	} else if b.consume(token.TT_NIL) {
		return NilNode{}, nil
	} else if b.consume(token.TT_LPAREN) {
		exp, err := b.expression()
		if err != nil {
			return nil, err
		}

		if b.consume(token.TT_RPAREN) {
			return ExpNode{Exp: exp}, nil
		}
		return nil, newSyntaxError("expected closing ')' after expression", b.curr)
	}
	return nil, newSyntaxError("expected a literal or an expression", b.curr)
}

// ------------------------------------
// Helpers
// ------------------------------------

func (b *AstBuilder) binaryOp(tokenTypes []token.TokenType, fun GrammarRuleFunc) (Node, error) {
	left, err := fun()
	if err != nil {
		return nil, err
	}

	for b.consumeAny(tokenTypes) {
		tok := b.curr

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
func (b *AstBuilder) check(tokType token.TokenType) bool {
	return b.checkAny([]token.TokenType{tokType})
}

// checkAny checks the next token if it matches any of the given types and returns true, otherwise it returns false
func (b *AstBuilder) checkAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if b.next.Type == straw {
			return true
		}
	}
	return false
}

// consume consumes the next token if it matches the given type and returns true, otherwise it returns false
func (b *AstBuilder) consume(tokType token.TokenType) bool {
	return b.consumeAny([]token.TokenType{tokType})
}

// consumeAny consumes the next token if it matches any of the given types and returns true, otherwise it returns false
func (b *AstBuilder) consumeAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if b.next.Type == straw {
			b.advance()
			return true
		}
	}
	return false
}

func (b *AstBuilder) advance() {
	if b.curr.Type != token.TT_EOF {
		b.prev = b.curr
		b.curr = b.next
		b.next = b.lex.NextToken()
	}
	if b.next.Type == token.TT_COMMENT {
		b.advance()
	}
}
