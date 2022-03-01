package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/lex"
)

type ProductionRuleHandler func() (Node, error)

type Tokenizer interface {
	NextToken() lex.Token
}

type Ast struct {
	tok  Tokenizer
	curr lex.Token
	prev lex.Token
	next lex.Token
}

func New(tok Tokenizer) *Ast {
	a := Ast{
		tok:  tok,
		curr: lex.Token{Type: lex.TT_ILLEGAL, Literal: "0"},
		prev: lex.Token{Type: lex.TT_ILLEGAL, Literal: "0"},
		next: lex.Token{Type: lex.TT_ILLEGAL, Literal: "0"},
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
	begin := a.curr.BeginPosition

	declarations := make([]Node, 0, 100)
	for !a.consume(lex.TT_EOF) {
		decl, err := a.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}

	end := a.curr.EndPosition

	return ProgramNode{
		Declarations: declarations,
		BeginPos:     begin,
		EndPos:       end,
	}, nil
}

// declaration -> funDecl
//             | varDecl
//             | statement ;
func (a *Ast) declaration() (Node, error) {
	if a.consume(lex.TT_FUNCTION) {
		return a.funDeclaration()
	} else if a.consume(lex.TT_VAR) {
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

	begin := a.curr.BeginPosition

	if a.check(lex.TT_LPAREN) { // Anonymous function (generate identifier)
		identifier = IdentifierNode{
			Token: lex.Token{
				Type:    lex.TT_IDENTIFIER,
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

	if !a.consume(lex.TT_LBRACE) {
		return nil, NewSyntaxError("expected opening '{' for function body", a.curr)
	}

	body, err := a.block()
	if err != nil {
		return nil, err
	}

	if _, ok := body.(BlockNode); !ok {
		return nil, NewSyntaxError("expected function body to be a block", a.curr)
	}

	end := a.curr.BeginPosition

	funcNode := FunctionNode{
		Identifier: identifier.(IdentifierNode),
		Parameters: parameters,
		Body:       body.(BlockNode),
		BeginPos:   begin,
		EndPos:     end,
	}

	// Function evalutaion call directly follows declaration
	if a.check(lex.TT_LPAREN) {
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
	if !a.consume(lex.TT_LPAREN) {
		return nil, NewSyntaxError("expected opening '(' for parameters", a.curr)
	}

	if a.consume(lex.TT_RPAREN) { // Function arity = 0
		return []IdentifierNode{}, nil
	}

	params := make([]IdentifierNode, 0, 255)

	param, err := a.parameter()
	if err != nil {
		return nil, err
	}

	params = append(params, param)

	for a.consume(lex.TT_COMMA) {
		param, err := a.parameter()
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	if !a.consume(lex.TT_RPAREN) {
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
	begin := a.curr.BeginPosition

	atom, err := a.atom()
	if err != nil {
		return nil, err
	}

	identifier, ok := atom.(IdentifierNode)
	if !ok {
		return nil, NewSyntaxError("Expected identifier after var", a.curr)
	}

	if !a.consume(lex.TT_ASSIGN) {
		end := a.curr.BeginPosition
		return VarStmtNode{
			Identifier: identifier,
			Value:      NilNode{},
			BeginPos:   begin,
			EndPos:     end,
		}, nil
	}

	value, err := a.expression()
	if err != nil {
		return nil, err
	}

	end := a.curr.BeginPosition

	return VarStmtNode{
		Identifier: identifier,
		Value:      value,
		BeginPos:   begin,
		EndPos:     end,
	}, nil
}

// statement -> exprStatementNode
//           | ifStatement
//           | whileStatement
//           | breakStatement
//           | continueStatement
//           | returnStatement
//           | deferStatement
//           | assertStatement
//           | block ;
func (a *Ast) statement() (Node, error) {
	if a.consume(lex.TT_IF) {
		return a.ifStatement()
	} else if a.consume(lex.TT_WHILE) {
		return a.whileStatement()
	} else if a.consume(lex.TT_BREAK) {
		return a.breakStatement()
	} else if a.consume(lex.TT_CONTINUE) {
		return a.continueStatement()
	} else if a.consume(lex.TT_RETURN) {
		return a.returnStatement()
	} else if a.consume(lex.TT_DEFER) {
		return a.deferStatement()
	} else if a.consume(lex.TT_ASSERT) {
		return a.assertStatement()
	} else if a.consume(lex.TT_LBRACE) {
		return a.block()
	} else {
		return a.expStatement()
	}

}

// exprStatementNode -> expression
func (a *Ast) expStatement() (Node, error) {
	begin := a.curr.BeginPosition

	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	end := a.curr.BeginPosition

	return ExpStmtNode{
		Exp:      exp,
		BeginPos: begin,
		EndPos:   end,
	}, nil
}

// ifStatement -> "if" "(" expression ")" statement ( "else" statement )? ;
func (a *Ast) ifStatement() (Node, error) {
	if !a.consume(lex.TT_LPAREN) {
		return nil, NewSyntaxError("expected opening '(' for if condition", a.curr)
	}

	begin := a.curr.BeginPosition

	condExp, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(lex.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for if condition", a.curr)
	}

	trueStmt, err := a.statement()
	if err != nil {
		return nil, err
	}

	var falseStmt Node = nil
	if a.consume(lex.TT_ELSE) {
		falseStmt, err = a.statement()
		if err != nil {
			return nil, err
		}
	}

	end := a.curr.BeginPosition

	return IfStmtNode{
		Exp:       condExp,
		TrueStmt:  trueStmt,
		FalseStmt: falseStmt,
		BeginPos:  begin,
		EndPos:    end,
	}, nil
}

// whileStatement -> "while" "(" expression ")" statement ;
func (a *Ast) whileStatement() (Node, error) {
	if !a.consume(lex.TT_LPAREN) {
		return nil, NewSyntaxError("expected opening '(' for 'while' condition", a.curr)
	}

	begin := a.curr.BeginPosition

	condition, err := a.expression()
	if err != nil {
		return nil, err
	}

	if !a.consume(lex.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for 'while' condition", a.curr)
	}

	body, err := a.statement()
	if err != nil {
		return nil, err
	}

	end := a.curr.BeginPosition

	return WhileStmtNode{
		Condition: condition,
		Body:      body,
		BeginPos:  begin,
		EndPos:    end,
	}, nil
}

// breakStatement -> "break" ;
func (a *Ast) breakStatement() (Node, error) {
	return BreakStmtNode{
		Token:    a.curr,
		BeginPos: a.curr.BeginPosition,
		EndPos:   a.curr.EndPosition,
	}, nil
}

// continueStatement -> "continue" ;
func (a *Ast) continueStatement() (Node, error) {
	return ContinueStmtNode{
		Token:    a.curr,
		BeginPos: a.curr.BeginPosition,
		EndPos:   a.curr.EndPosition,
	}, nil
}

// returnStatement -> "return" expression ;
func (a *Ast) returnStatement() (Node, error) {
	begin := a.curr.BeginPosition

	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	end := a.curr.BeginPosition

	return ReturnStmtNode{
		Exp:      exp,
		BeginPos: begin,
		EndPos:   end,
	}, nil
}

// deferStatement -> "defer" funcCall ;
func (a *Ast) deferStatement() (Node, error) {
	begin := a.curr.BeginPosition

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

	end := a.curr.BeginPosition

	return DeferStmtNode{
		Call:     call.(CallNode),
		BeginPos: begin,
		EndPos:   end,
	}, nil
}

// assertStatement -> "assert" expression ;
func (a *Ast) assertStatement() (Node, error) {
	begin := a.curr.BeginPosition

	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	end := a.curr.BeginPosition

	return AssertStmtNode{
		Exp:      exp,
		BeginPos: begin,
		EndPos:   end,
	}, nil
}

// block -> "{" declaration* "}" ;
func (a *Ast) block() (Node, error) {
	begin := a.curr.BeginPosition

	declarations := make([]Node, 0, 100)

	for !a.check(lex.TT_RBRACE) && !a.check(lex.TT_EOF) {
		decl, err := a.declaration()
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, decl)
	}

	if !a.consume(lex.TT_RBRACE) {
		return nil, NewSyntaxError("expected closing '}'", a.curr)
	}

	end := a.curr.BeginPosition

	return BlockNode{
		Declarations: declarations,
		BeginPos:     begin,
		EndPos:       end,
	}, nil
}

// expression -> assignment ( "?" assignment ":" assignment )? ;
func (a *Ast) expression() (Node, error) {
	begin := a.curr.BeginPosition

	exp, err := a.assignment()
	if err != nil {
		return nil, err
	}

	// Check ternary operator: <assignment> ? <assignment> : <assignment>
	if a.consume(lex.TT_QUESTION) {
		trueExp, err := a.assignment()
		if err != nil {
			return nil, err
		}

		if !a.consume(lex.TT_COLON) {
			return nil, NewSyntaxError("expected ':'", a.curr)
		}

		falseExp, err := a.assignment()
		if err != nil {
			return nil, err
		}

		end := a.curr.BeginPosition

		return TernaryOpNode{
			Exp:      exp,
			TrueExp:  trueExp,
			FalseExp: falseExp,
			BeginPos: begin,
			EndPos:   end,
		}, nil
	}

	return exp, nil
}

// assignment -> IDENTIFIER "=" assignment
//            | logicalOr ;
func (a *Ast) assignment() (Node, error) {
	begin := a.curr.BeginPosition

	expr, err := a.logicalOr()
	if err != nil {
		return nil, err
	}

	if a.consume(lex.TT_ASSIGN) {
		if _, ok := expr.(IdentifierNode); !ok {
			return nil, NewSyntaxError("expected an identifier for assignment", a.curr)
		}

		assign, err := a.assignment()
		if err != nil {
			return nil, err
		}

		end := a.curr.BeginPosition

		return AssignmentNode{
			Identifier: expr.(IdentifierNode),
			Value:      assign,
			BeginPos:   begin,
			EndPos:     end,
		}, nil
	}
	return expr, nil
}

// logicalOr -> logicalAnd ( "||" logicalAnd )*
func (a *Ast) logicalOr() (Node, error) {
	begin := a.curr.BeginPosition

	left, err := a.logicalAnd()
	if err != nil {
		return nil, err
	}

	for a.consume(lex.TT_LOGICAL_OR) {
		right, err := a.equality()
		if err != nil {
			return nil, err
		}

		end := a.curr.BeginPosition

		left = LogicalOrNode{
			LHS:      left,
			RHS:      right,
			BeginPos: begin,
			EndPos:   end,
		}
	}
	return left, nil
}

// logicalAnd -> equality ( "&&" equality )* ;
func (a *Ast) logicalAnd() (Node, error) {
	begin := a.curr.BeginPosition

	left, err := a.equality()
	if err != nil {
		return nil, err
	}

	for a.consume(lex.TT_LOGICAL_AND) {
		right, err := a.equality()
		if err != nil {
			return nil, err
		}

		end := a.curr.BeginPosition

		left = LogicalAndNode{
			LHS:      left,
			RHS:      right,
			BeginPos: begin,
			EndPos:   end,
		}
	}
	return left, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (a *Ast) equality() (Node, error) {
	return a.binaryOp([]lex.TokenType{lex.TT_EQ, lex.TT_NEQ}, a.comparison)
}

// comparison -> term ( ( "<" | "<=" | ">" | ">=" ) term )* ;
func (a *Ast) comparison() (Node, error) {
	return a.binaryOp([]lex.TokenType{lex.TT_LT, lex.TT_LTE, lex.TT_GT, lex.TT_GTE}, a.term)
}

// term -> factor ( ( "+" | "-" ) factor )* ;
func (a *Ast) term() (Node, error) {
	return a.binaryOp([]lex.TokenType{lex.TT_PLUS, lex.TT_MINUS}, a.factor)
}

// factor -> unary ( ( "/" | "*" | "%" ) unary )* ;
func (a *Ast) factor() (Node, error) {
	return a.binaryOp([]lex.TokenType{lex.TT_DIVIDE, lex.TT_MULTIPLY, lex.TT_MODULO}, a.unary)
}

// unary -> ( "!" | "-" ) unary
//       | call ;
func (a *Ast) unary() (Node, error) {
	begin := a.curr.BeginPosition

	var node Node
	for a.consumeAny([]lex.TokenType{lex.TT_NOT, lex.TT_MINUS}) {
		tok := a.curr

		n, err := a.unary()
		if err != nil {
			return nil, err
		}

		end := a.curr.BeginPosition

		node = UnaryOpNode{
			Op:       tok,
			Operand:  n,
			BeginPos: begin,
			EndPos:   end,
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

	if a.check(lex.TT_LPAREN) {
		expr, err = a.funcCall(expr)
		if err != nil {
			return nil, err
		}
	} else if a.check(lex.TT_LBRACKET) {
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
	for a.consume(lex.TT_LPAREN) {
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

	begin := a.curr.BeginPosition

	for !a.check(lex.TT_RPAREN) {
		arguments, err = a.arguments()
		if err != nil {
			return nil, err
		}
	}

	if !a.consume(lex.TT_RPAREN) {
		return nil, NewSyntaxError("expected closing ')' for function call", a.curr)
	}

	end := a.curr.BeginPosition

	return CallNode{
		Callee:    callee,
		Arguments: arguments,
		BeginPos:  begin,
		EndPos:    end,
	}, nil
}

// indexCall -> atom ( "[" expression "]" )* ;
func (a *Ast) indexCall(atom Node) (Node, error) {
	begin := a.curr.BeginPosition
	expr := atom
	for a.consume(lex.TT_LBRACKET) {
		indexExpr, err := a.expression()
		if err != nil {
			return nil, err
		}

		if !a.consume(lex.TT_RBRACKET) {
			return nil, NewSyntaxError("expected closing ']' for index operation", a.curr)
		}

		end := a.curr.BeginPosition

		expr = IndexOfNode{
			Sequence: expr,
			Index:    indexExpr,
			BeginPos: begin,
			EndPos:   end,
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
	for a.consume(lex.TT_COMMA) {
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
	if a.consume(lex.TT_NUMBER) {
		return NumberNode{
			Token:    a.curr,
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	} else if a.consume(lex.TT_STRING) {
		return StringNode{
			Token:    a.curr,
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	} else if a.consumeAny([]lex.TokenType{lex.TT_TRUE, lex.TT_FALSE}) {
		return BooleanNode{
			Token:    a.curr,
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	} else if a.consume(lex.TT_IDENTIFIER) {
		return IdentifierNode{
			Token:    a.curr,
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	} else if a.consume(lex.TT_NIL) {
		return NilNode{
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	} else if a.consume(lex.TT_LPAREN) {
		return a.nestedExpressionNode()
	} else if a.consume(lex.TT_LBRACE) {
		return a.mapNode()
	} else if a.consume(lex.TT_LBRACKET) {
		return a.listNode()
	} else if a.consume(lex.TT_FUNCTION) {
		return a.funDeclaration()
	} else if a.consume(lex.TT_COMMENT) {
		return CommentNode{
			Token:    a.curr,
			BeginPos: a.curr.BeginPosition,
			EndPos:   a.curr.EndPosition,
		}, nil
	}

	return nil, NewSyntaxError("expected a literal or an expression", a.curr)
}

func (a *Ast) nestedExpressionNode() (Node, error) {
	begin := a.curr.BeginPosition
	exp, err := a.expression()
	if err != nil {
		return nil, err
	}

	if a.consume(lex.TT_RPAREN) {
		end := a.curr.BeginPosition
		return ExpNode{
			Exp:      exp,
			BeginPos: begin,
			EndPos:   end,
		}, nil
	}
	return nil, NewSyntaxError("expected closing ')' after expression", a.curr)
}

// map -> "{" keyValuePairs? "}" ;
func (a *Ast) mapNode() (Node, error) {
	begin := a.curr.BeginPosition
	if a.consume(lex.TT_RBRACE) { // Map is empty {}
		end := a.curr.BeginPosition
		return MapNode{
			Elements: make([]KeyValueNode, 0),
			BeginPos: begin,
			EndPos:   end,
		}, nil
	} else {
		kvps, err := a.keyValuePairs()
		if err != nil {
			return nil, err
		}

		if !a.consume(lex.TT_RBRACE) {
			return nil, NewSyntaxError("expected closing '}' for map", a.curr)
		}

		end := a.curr.BeginPosition

		return MapNode{
			Elements: kvps,
			BeginPos: begin,
			EndPos:   end,
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

	for a.consume(lex.TT_COMMA) {
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
	begin := a.curr.BeginPosition

	key, err := a.expression()

	if err != nil {
		return KeyValueNode{}, err
	}

	if !a.consume(lex.TT_COLON) {
		return KeyValueNode{}, NewSyntaxError("expected ':' for map key-value pair", a.curr)
	}

	value, err := a.expression()

	if err != nil {
		return KeyValueNode{}, err
	}

	end := a.curr.BeginPosition

	return KeyValueNode{
		Key:      key,
		Value:    value,
		BeginPos: begin,
		EndPos:   end,
	}, nil
}

// list -> "[" arguments? "]" ;
func (a *Ast) listNode() (Node, error) {
	begin := a.curr.BeginPosition
	if a.consume(lex.TT_RBRACKET) { // List is empty []
		end := a.curr.BeginPosition
		return ListNode{
			Elements: make([]Node, 0),
			BeginPos: begin,
			EndPos:   end,
		}, nil
	} else {
		arguments, err := a.arguments()
		if err != nil {
			return nil, err
		}

		if !a.consume(lex.TT_RBRACKET) {
			return nil, NewSyntaxError("expected closing ']' for list", a.curr)
		}

		end := a.curr.BeginPosition

		return ListNode{
			Elements: arguments,
			BeginPos: begin,
			EndPos:   end,
		}, nil
	}
}

// ------------------------------------
// Helpers
// ------------------------------------

func (a *Ast) binaryOp(tokenTypes []lex.TokenType, fun ProductionRuleHandler) (Node, error) {
	begin := a.curr.BeginPosition

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

		end := a.curr.BeginPosition

		left = BinaryOpNode{
			LeftExp:  left,
			Op:       tok,
			RightExp: right,
			BeginPos: begin,
			EndPos:   end,
		}
	}
	return left, nil
}

// check checks the next token if it matches the given type and returns true, otherwise it returns false
func (a *Ast) check(tokType lex.TokenType) bool {
	return a.checkAny([]lex.TokenType{tokType})
}

// checkAny checks the next token if it matches any of the given types and returns true, otherwise it returns false
func (a *Ast) checkAny(tokTypes []lex.TokenType) bool {
	for _, straw := range tokTypes {
		if a.next.Type == straw {
			return true
		}
	}
	return false
}

// consume consumes the next token if it matches the given type and returns true, otherwise it returns false
func (a *Ast) consume(tokType lex.TokenType) bool {
	return a.consumeAny([]lex.TokenType{tokType})
}

// consumeAny consumes the next token if it matches any of the given types and returns true, otherwise it returns false
func (a *Ast) consumeAny(tokTypes []lex.TokenType) bool {
	for _, straw := range tokTypes {
		if a.next.Type == straw {
			a.advance()
			return true
		}
	}
	return false
}

func (a *Ast) advance() {
	if a.curr.Type != lex.TT_EOF {
		a.prev = a.curr
		a.curr = a.next
		a.next = a.tok.NextToken()
	}
}
