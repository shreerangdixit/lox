package ast

import (
	"github.com/shreerangdixit/lox/token"
)

type ProductionRuleHandler func() (Node, error)

type Tokenizer interface {
	NextToken() token.Token
}

type Ast struct {
	tok  Tokenizer
	curr token.Token
	prev token.Token
	next token.Token
}

func New(tok Tokenizer) *Ast {
	a := Ast{
		tok:  tok,
		curr: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		prev: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
		next: token.Token{Type: token.TT_ILLEGAL, Literal: "0"},
	}
	a.advance()
	return &a
}

func (a *Ast) RootNode() (Node, error) {
	return a.program()
}

// ------------------------------------
// Production rule handlers
// ------------------------------------

// program -> declaration* EOF ;
func (a *Ast) program() (Node, error) {
	declarations := make([]Node, 0, 100)
	for !a.consume(token.TT_EOF) {
		decl, err := a.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}
	return ProgramNode{
		Declarations: declarations,
	}, nil
}

// declaration -> varDecl
//             | statement ;
func (a *Ast) declaration() (Node, error) {
	if a.consume(token.TT_VAR) {
		return a.varDeclaration()
	}
	return a.statement()
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";" ;
func (a *Ast) varDeclaration() (Node, error) {
	atom, err := a.atom()
	if err != nil {
		return nil, err
	}

	identifier, ok := atom.(IdentifierNode)
	if !ok {
		return nil, newSyntaxError("Expected identifier after var", a.curr)
	}

	if !a.consume(token.TT_ASSIGN) {
		if !a.consume(token.TT_SEMICOLON) {
			return nil, newSyntaxError("expected a ; at the end of a declaration", a.curr)
		}

		return VarStmtNode{
			Identifier: identifier,
			Value:      NilNode{},
		}, nil
	}

	value, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_SEMICOLON) {
		return nil, newSyntaxError("expected a ; at the end of a declaration", a.curr)
	}

	return VarStmtNode{
		Identifier: identifier,
		Value:      value,
	}, nil
}

// statement -> exprStatement
//           | ifStatement
//           | block ;
func (a *Ast) statement() (Node, error) {
	if a.consume(token.TT_IF) {
		return a.ifStatement()
	} else if a.consume(token.TT_WHILE) {
		return a.whileStatement()
	} else if a.consume(token.TT_LBRACE) {
		return a.block()
	}
	return a.exprStatement()
}

// ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
func (a *Ast) ifStatement() (Node, error) {
	if !a.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for if condition", a.curr)
	}

	condExp, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for if condition", a.curr)
	}

	trueStmt, err := a.statement()
	if err != nil {
		return nil, err
	}

	var falseStmt Node = nil
	if a.consume(token.TT_ELSE) {
		falseStmt, err = a.statement()
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

// whileStatement -> "while" "(" expression ")" statement ;
func (a *Ast) whileStatement() (Node, error) {
	if !a.consume(token.TT_LPAREN) {
		return nil, newSyntaxError("expected opening '(' for 'while' condition", a.curr)
	}

	condition, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for 'while' condition", a.curr)
	}

	body, err := a.statement()
	if err != nil {
		return nil, err
	}

	return WhileStmtNode{
		Condition: condition,
		Body:      body,
	}, nil
}

// block -> "{" declaration* "}" ;
func (a *Ast) block() (Node, error) {
	declarations := make([]Node, 0, 100)

	for !a.check(token.TT_RBRACE) && !a.check(token.TT_EOF) {
		decl, err := a.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}

	if !a.consume(token.TT_RBRACE) {
		return nil, newSyntaxError("expected closing '}'", a.curr)
	}

	return BlockNode{
		Declarations: declarations,
	}, nil
}

// exprStatement -> expression ";" ;
func (a *Ast) exprStatement() (Node, error) {
	expr, err := a.expression()
	if err != nil {
		return nil, err
	}

	if a.consume(token.TT_SEMICOLON) {
		return ExpStmtNode{Exp: expr}, nil
	}
	return nil, newSyntaxError("expected a ; at the end of an expression statement", a.curr)
}

// expression -> assignment ( "?" assignment ":" assignment )? ;
func (a *Ast) expression() (Node, error) {
	exp, err := a.assignment()
	if err != nil {
		return nil, err
	}

	// Check ternary operator: <assignment> ? <assignment> : <assignment>
	if a.consume(token.TT_QUESTION) {
		trueExp, err := a.assignment()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_COLON) {
			return nil, newSyntaxError("expected ':'", a.curr)
		}

		falseExp, err := a.assignment()
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
func (a *Ast) assignment() (Node, error) {
	expr, err := a.logical_or()
	if err != nil {
		return nil, err
	}

	if a.consume(token.TT_ASSIGN) {
		if _, ok := expr.(IdentifierNode); !ok {
			return nil, newSyntaxError("expected an identifier for assignment", a.curr)
		}

		assign, err := a.assignment()
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
func (a *Ast) logical_or() (Node, error) {
	left, err := a.logical_and()
	if err != nil {
		return nil, err
	}

	for a.consume(token.TT_LOGICAL_OR) {
		right, err := a.equality()
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
func (a *Ast) logical_and() (Node, error) {
	left, err := a.equality()
	if err != nil {
		return nil, err
	}

	for a.consume(token.TT_LOGICAL_AND) {
		right, err := a.equality()
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
func (a *Ast) equality() (Node, error) {
	return a.binaryOp([]token.TokenType{token.TT_EQ, token.TT_NEQ}, a.comparison)
}

// comparison -> term ( ( "<" | "<=" | ">" | ">=" ) term )* ;
func (a *Ast) comparison() (Node, error) {
	return a.binaryOp([]token.TokenType{token.TT_LT, token.TT_LTE, token.TT_GT, token.TT_GTE}, a.term)
}

// term -> factor ( ( "+" | "-" ) factor )* ;
func (a *Ast) term() (Node, error) {
	return a.binaryOp([]token.TokenType{token.TT_PLUS, token.TT_MINUS}, a.factor)
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (a *Ast) factor() (Node, error) {
	return a.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY}, a.unary)
}

// unary -> ( "!" | "-" ) unary
//       | call ;
func (a *Ast) unary() (Node, error) {
	var node Node
	for a.consumeAny([]token.TokenType{token.TT_NOT, token.TT_MINUS}) {
		tok := a.curr

		n, err := a.unary()
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
	return a.call()
}

// call -> atom ( "(" arguments? ")" )*
//      | atom ( "[" expression "]" )* ;
func (a *Ast) call() (Node, error) {
	expr, err := a.atom()
	if err != nil {
		return nil, err
	}

	if a.check(token.TT_LPAREN) {
		for a.consume(token.TT_LPAREN) {
			expr, err = a.finishCall(expr)
			if err != nil {
				return nil, err
			}
		}
	} else if a.consume(token.TT_LBRACKET) {
		indexExpr, err := a.expression()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_RBRACKET) {
			return nil, newSyntaxError("expected closing ']' for index operation", a.curr)
		}

		expr = IndexOfNode{
			Sequence: expr,
			Index:    indexExpr,
		}
	}
	return expr, nil
}

func (a *Ast) finishCall(callee Node) (Node, error) {
	arguments := []Node{}
	var err error

	for !a.check(token.TT_RPAREN) {
		arguments, err = a.arguments()
		if err != nil {
			return nil, err
		}
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, newSyntaxError("expected closing ')' for function call", a.curr)
	}

	return CallNode{
		Callee:    callee,
		Arguments: arguments,
	}, nil
}

// arguments -> expression ( "," expression )* ;
func (a *Ast) arguments() ([]Node, error) {
	arguments := make([]Node, 0, 255)

	arg, err := a.expression()
	if err != nil {
		return nil, err
	}

	arguments = append(arguments, arg)
	for a.consume(token.TT_COMMA) {
		arg, err := a.expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, arg)
	}

	return arguments, nil
}

// atom -> NUMBER | STRING | "true" | "false" | "nil"
//      | "(" expression ")"
//      | "[" arguments? "]"
//      | IDENTIFIER ;
func (a *Ast) atom() (Node, error) {
	if a.consume(token.TT_NUMBER) {
		return NumberNode{Token: a.curr}, nil
	} else if a.consume(token.TT_STRING) {
		return StringNode{Token: a.curr}, nil
	} else if a.consumeAny([]token.TokenType{token.TT_TRUE, token.TT_FALSE}) {
		return BooleanNode{Token: a.curr}, nil
	} else if a.consume(token.TT_IDENTIFIER) {
		return IdentifierNode{Token: a.curr}, nil
	} else if a.consume(token.TT_NIL) {
		return NilNode{}, nil
	} else if a.consume(token.TT_LPAREN) {
		exp, err := a.expression()
		if err != nil {
			return nil, err
		}

		if a.consume(token.TT_RPAREN) {
			return ExpNode{Exp: exp}, nil
		}
		return nil, newSyntaxError("expected closing ')' after expression", a.curr)
	} else if a.consume(token.TT_LBRACKET) {
		return a.list()
	}
	return nil, newSyntaxError("expected a literal or an expression", a.curr)
}

func (a *Ast) list() (Node, error) {
	// List is empty []
	if a.consume(token.TT_RBRACKET) {
		return ListNode{
			Elements: make([]Node, 0),
		}, nil
	} else {
		arguments, err := a.arguments()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_RBRACKET) {
			return nil, newSyntaxError("expected closing ']' for list", a.curr)
		}

		return ListNode{
			Elements: arguments,
		}, nil
	}
}

// ------------------------------------
// Helpers
// ------------------------------------

func (a *Ast) binaryOp(tokenTypes []token.TokenType, fun ProductionRuleHandler) (Node, error) {
	left, err := fun()
	if err != nil {
		return nil, err
	}

	for a.consumeAny(tokenTypes) {
		tok := a.curr

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
func (a *Ast) check(tokType token.TokenType) bool {
	return a.checkAny([]token.TokenType{tokType})
}

// checkAny checks the next token if it matches any of the given types and returns true, otherwise it returns false
func (a *Ast) checkAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if a.next.Type == straw {
			return true
		}
	}
	return false
}

// consume consumes the next token if it matches the given type and returns true, otherwise it returns false
func (a *Ast) consume(tokType token.TokenType) bool {
	return a.consumeAny([]token.TokenType{tokType})
}

// consumeAny consumes the next token if it matches any of the given types and returns true, otherwise it returns false
func (a *Ast) consumeAny(tokTypes []token.TokenType) bool {
	for _, straw := range tokTypes {
		if a.next.Type == straw {
			a.advance()
			return true
		}
	}
	return false
}

func (a *Ast) advance() {
	if a.curr.Type != token.TT_EOF {
		a.prev = a.curr
		a.curr = a.next
		a.next = a.tok.NextToken()
	}
	if a.next.Type == token.TT_COMMENT {
		a.advance()
	}
}
