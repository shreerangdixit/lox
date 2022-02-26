package ast

import (
	"fmt"

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

// declaration -> funDecl
//             | varDecl
//             | statement ;
func (a *Ast) declaration() (Node, error) {
	if a.consume(token.TT_FUNCTION) {
		return a.funDeclaration()
	} else if a.consume(token.TT_VAR) {
		return a.varDeclaration()
	} else {
		return a.statement()
	}
}

// funDecl  -> "fun" function ;
// function -> IDENTIFIER? "(" parameters? ")" block ( funcCall ";" )? ;
func (a *Ast) funDeclaration() (Node, error) {
	var identifier Node
	var err error

	if a.check(token.TT_LPAREN) { // Anonymous function (generate identifier)
		identifier = IdentifierNode{
			Token: token.Token{
				Type:    token.TT_IDENTIFIER,
				Literal: fmt.Sprintf("anon-%s", RandStringBytes(8)),
			},
		}
	} else {
		identifier, err = a.atom()
		if err != nil {
			return nil, err
		}
	}

	if _, ok := identifier.(IdentifierNode); !ok {
		return nil, NewSyntaxError("function name should be an identifier", a.curr)
	}

	parameters, err := a.parameters()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_LBRACE) {
		return nil, NewSyntaxError("expected opening '{' for function body", a.curr)
	}

	body, err := a.block()
	if err != nil {
		return nil, err
	}

	if _, ok := body.(BlockNode); !ok {
		return nil, NewSyntaxError("expected function body to be a block", a.curr)
	}

	funcNode := FunctionNode{
		Identifier: identifier.(IdentifierNode),
		Parameters: parameters,
		Body:       body.(BlockNode),
	}

	// Function evalutaion call directly follows declaration
	if a.check(token.TT_LPAREN) {
		funcResult, err := a.funcCall(funcNode)
		if err != nil {
			return nil, err
		}

		return funcResult, nil
	} else {
		return funcNode, nil
	}
}

// parameters -> IDENTIFIER ("," IDENTIFIER)* ;
func (a *Ast) parameters() ([]IdentifierNode, error) {
	if !a.consume(token.TT_LPAREN) {
		return nil, NewSyntaxError("expected opening '(' for parameters", a.curr)
	}

	if a.consume(token.TT_RPAREN) { // Function arity = 0
		return []IdentifierNode{}, nil
	}

	params := make([]IdentifierNode, 0, 255)

	param, err := a.parameter()
	if err != nil {
		return nil, err
	}

	params = append(params, param)

	for a.consume(token.TT_COMMA) {
		param, err := a.parameter()
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for parameters", a.curr)
	}

	return params, nil
}

func (a *Ast) parameter() (IdentifierNode, error) {
	param, err := a.atom()
	if err != nil {
		return IdentifierNode{}, err
	}

	if _, ok := param.(IdentifierNode); !ok {
		return IdentifierNode{}, NewSyntaxError("param should be an identifier", a.curr)
	}

	return param.(IdentifierNode), nil
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ;
func (a *Ast) varDeclaration() (Node, error) {
	atom, err := a.atom()
	if err != nil {
		return nil, err
	}

	identifier, ok := atom.(IdentifierNode)
	if !ok {
		return nil, NewSyntaxError("Expected identifier after var", a.curr)
	}

	if !a.consume(token.TT_ASSIGN) {
		return VarStmtNode{
			Identifier: identifier,
			Value:      NilNode{},
		}, nil
	}

	value, err := a.expression()
	if err != nil {
		return nil, err
	}

	return VarStmtNode{
		Identifier: identifier,
		Value:      value,
	}, nil
}

// statement -> expression
//           | ifStatement
//           | whileStatement
//           | breakStatement
//           | continueStatement
//           | returnStatement
//           | deferStatement
//           | block ;
func (a *Ast) statement() (Node, error) {
	if a.consume(token.TT_IF) {
		return a.ifStatement()
	} else if a.consume(token.TT_WHILE) {
		return a.whileStatement()
	} else if a.consume(token.TT_BREAK) {
		return a.breakStatement()
	} else if a.consume(token.TT_CONTINUE) {
		return a.continueStatement()
	} else if a.consume(token.TT_RETURN) {
		return a.returnStatement()
	} else if a.consume(token.TT_DEFER) {
		return a.deferStatement()
	} else if a.consume(token.TT_LBRACE) {
		return a.block()
	}
	return a.expression()
}

// ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
func (a *Ast) ifStatement() (Node, error) {
	if !a.consume(token.TT_LPAREN) {
		return nil, NewSyntaxError("expected opening '(' for if condition", a.curr)
	}

	condExp, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for if condition", a.curr)
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
		return nil, NewSyntaxError("expected opening '(' for 'while' condition", a.curr)
	}

	condition, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(token.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for 'while' condition", a.curr)
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

// breakStatement -> "break" ;
func (a *Ast) breakStatement() (Node, error) {
	return BreakStmtNode{
		Token: a.curr,
	}, nil
}

// continueStatement -> "continue" ;
func (a *Ast) continueStatement() (Node, error) {
	return ContinueStmtNode{
		Token: a.curr,
	}, nil
}

// returnStatement -> "return" expression ;
func (a *Ast) returnStatement() (Node, error) {
	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	return ReturnStmtNode{
		Exp: exp,
	}, nil
}

// deferStatement -> "defer" funcCall ;
func (a *Ast) deferStatement() (Node, error) {
	node, err := a.atom()
	if err != nil {
		return nil, err
	}

	call, err := a.funcCall(node)
	if err != nil {
		return nil, err
	}

	if _, ok := call.(CallNode); !ok {
		return nil, NewSyntaxError("invalid call node", a.curr)
	}

	return DeferStmtNode{
		Call: call.(CallNode),
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
		return nil, NewSyntaxError("expected closing '}'", a.curr)
	}

	return BlockNode{
		Declarations: declarations,
	}, nil
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
			return nil, NewSyntaxError("expected ':'", a.curr)
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
//            | logicalOr ;
func (a *Ast) assignment() (Node, error) {
	expr, err := a.logicalOr()
	if err != nil {
		return nil, err
	}

	if a.consume(token.TT_ASSIGN) {
		if _, ok := expr.(IdentifierNode); !ok {
			return nil, NewSyntaxError("expected an identifier for assignment", a.curr)
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

// logicalOr -> logicalAnd ( "||" logicalAnd )*
func (a *Ast) logicalOr() (Node, error) {
	left, err := a.logicalAnd()
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

// logicalAnd -> equality ( "&&" equality )* ;
func (a *Ast) logicalAnd() (Node, error) {
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

// factor -> unary ( ( "/" | "*" | "%" ) unary )* ;
func (a *Ast) factor() (Node, error) {
	return a.binaryOp([]token.TokenType{token.TT_DIVIDE, token.TT_MULTIPLY, token.TT_MODULO}, a.unary)
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

// call -> funcCall
//      | indexCall ;
func (a *Ast) call() (Node, error) {
	expr, err := a.atom()
	if err != nil {
		return nil, err
	}

	if a.check(token.TT_LPAREN) {
		expr, err = a.funcCall(expr)
		if err != nil {
			return nil, err
		}
	} else if a.check(token.TT_LBRACKET) {
		expr, err = a.indexCall(expr)
		if err != nil {
			return nil, err
		}
	}
	return expr, nil
}

// funcCall -> atom ( "(" arguments? ")" )* ;
func (a *Ast) funcCall(atom Node) (Node, error) {
	exp := atom
	var err error
	for a.consume(token.TT_LPAREN) {
		exp, err = a.finishCall(exp)
		if err != nil {
			return nil, err
		}
	}
	return exp, nil
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
		return nil, NewSyntaxError("expected closing ')' for function call", a.curr)
	}

	return CallNode{
		Callee:    callee,
		Arguments: arguments,
	}, nil
}

// indexCall -> atom ( "[" expression "]" )* ;
func (a *Ast) indexCall(atom Node) (Node, error) {
	expr := atom
	for a.consume(token.TT_LBRACKET) {
		indexExpr, err := a.expression()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_RBRACKET) {
			return nil, NewSyntaxError("expected closing ']' for index operation", a.curr)
		}

		expr = IndexOfNode{
			Sequence: expr,
			Index:    indexExpr,
		}
	}
	return expr, nil
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
//      | list
//      | map
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
		return a.nestedExpressionNode()
	} else if a.consume(token.TT_LBRACE) {
		return a.mapNode()
	} else if a.consume(token.TT_LBRACKET) {
		return a.listNode()
	} else if a.consume(token.TT_FUNCTION) {
		return a.funDeclaration()
	}

	return nil, NewSyntaxError("expected a literal or an expression", a.curr)
}

func (a *Ast) nestedExpressionNode() (Node, error) {
	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	if a.consume(token.TT_RPAREN) {
		return ExpNode{Exp: exp}, nil
	}
	return nil, NewSyntaxError("expected closing ')' after expression", a.curr)
}

// map -> "{" keyValuePairs? "}" ;
func (a *Ast) mapNode() (Node, error) {
	if a.consume(token.TT_RBRACE) { // Map is empty {}
		return MapNode{
			Elements: make([]KeyValueNode, 0),
		}, nil
	} else {
		kvps, err := a.keyValuePairs()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_RBRACE) {
			return nil, NewSyntaxError("expected closing '}' for map", a.curr)
		}

		return MapNode{
			Elements: kvps,
		}, nil
	}
}

// keyValuePairs -> expression ":" expression ( "," expression ":" expression )* ;
func (a *Ast) keyValuePairs() ([]KeyValueNode, error) {
	kvps := make([]KeyValueNode, 0, 255)

	kvp, err := a.keyValuePair()
	if err != nil {
		return nil, err
	}

	kvps = append(kvps, kvp)

	for a.consume(token.TT_COMMA) {
		kvp, err := a.keyValuePair()
		if err != nil {
			return nil, err
		}

		kvps = append(kvps, kvp)
	}

	return kvps, nil
}

// expression ":" expression
func (a *Ast) keyValuePair() (KeyValueNode, error) {
	key, err := a.expression()

	if err != nil {
		return KeyValueNode{}, err
	}

	if !a.consume(token.TT_COLON) {
		return KeyValueNode{}, NewSyntaxError("expected ':' for map key-value pair", a.curr)
	}

	value, err := a.expression()

	if err != nil {
		return KeyValueNode{}, err
	}

	return KeyValueNode{
		Key:   key,
		Value: value,
	}, nil
}

// list -> "[" arguments? "]" ;
func (a *Ast) listNode() (Node, error) {
	if a.consume(token.TT_RBRACKET) { // List is empty []
		return ListNode{
			Elements: make([]Node, 0),
		}, nil
	} else {
		arguments, err := a.arguments()
		if err != nil {
			return nil, err
		}

		if !a.consume(token.TT_RBRACKET) {
			return nil, NewSyntaxError("expected closing ']' for list", a.curr)
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
