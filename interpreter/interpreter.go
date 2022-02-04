package interpreter

import (
	"fmt"
	"lox/parser"
	"lox/token"
	"lox/types"
	"strconv"
)

type Interpreter struct {
	ast     parser.Node
	globals map[string]types.TypeValue
}

func New(parser *parser.Parser) (*Interpreter, error) {
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return &Interpreter{
		ast:     ast,
		globals: make(map[string]types.TypeValue),
	}, nil
}

func (i *Interpreter) Run() (types.TypeValue, error) {
	return i.visit(i.ast)
}

func (i *Interpreter) visit(node parser.Node) (types.TypeValue, error) {
	switch node.(type) {
	case parser.ProgramNode:
		return i.visitProgramNode(node.(parser.ProgramNode))
	case parser.LetStatementNode:
		return i.visitLetStatementNode(node.(parser.LetStatementNode))
	case parser.ExpressionStatementNode:
		return i.visitExpressionStatementNode(node.(parser.ExpressionStatementNode))
	case parser.PrintStatementNode:
		return i.visitPrintStatementNode(node.(parser.PrintStatementNode))
	case parser.ExpressionNode:
		return i.visit(node.(parser.ExpressionNode).Exp)
	case parser.BinaryOpNode:
		return i.visitBinaryOpNode(node.(parser.BinaryOpNode))
	case parser.UnaryOpNode:
		return i.visitUnaryOpNode(node.(parser.UnaryOpNode))
	case parser.NumberNode:
		return i.visitNumberNode(node.(parser.NumberNode))
	case parser.BooleanNode:
		return i.visitBooleanNode(node.(parser.BooleanNode))
	case parser.IdentifierNode:
		return i.visitIdentifierNode(node.(parser.IdentifierNode))
	case parser.NilNode:
		return i.visitNilNode(node.(parser.NilNode))
	}
	return types.TypeValue{}, fmt.Errorf("invalid node: %T", node)
}

func (i *Interpreter) visitProgramNode(node parser.ProgramNode) (types.TypeValue, error) {
	for _, node := range node.Declarations {
		_, err := i.visit(node)
		if err != nil {
			return types.TypeValue{}, err
		}
	}
	return types.TypeValue{}, nil
}

func (i *Interpreter) visitLetStatementNode(node parser.LetStatementNode) (types.TypeValue, error) {
	expression, err := i.visit(node.Value)
	if err != nil {
		return types.TypeValue{}, err
	}

	i.globals[node.Identifier.Token.Literal] = expression
	return types.TypeValue{}, nil
}

func (i *Interpreter) visitExpressionStatementNode(node parser.ExpressionStatementNode) (types.TypeValue, error) {
	// Evaluate the expression and discard the result (for now)
	_, err := i.visit(node.Exp)
	if err != nil {
		return types.TypeValue{}, err
	}

	return types.TypeValue{}, nil
}

func (i *Interpreter) visitPrintStatementNode(node parser.PrintStatementNode) (types.TypeValue, error) {
	result, err := i.visit(node.Exp)
	if err != nil {
		return types.TypeValue{}, err
	}

	fmt.Printf("%v\n", result.Value)

	return types.TypeValue{}, nil
}

func (i *Interpreter) visitBinaryOpNode(node parser.BinaryOpNode) (types.TypeValue, error) {
	left, err := i.visit(node.LHS)
	if err != nil {
		return types.TypeValue{}, err
	}

	right, err := i.visit(node.RHS)
	if err != nil {
		return types.TypeValue{}, err
	}

	switch node.Op.Type {
	case token.TT_PLUS:
		return left.Add(right)
	case token.TT_MINUS:
		return left.Subtract(right)
	case token.TT_DIVIDE:
		return left.Divide(right)
	case token.TT_MULTIPLY:
		return left.Multiply(right)
	case token.TT_EQ:
		return left.Equals(right)
	case token.TT_NEQ:
		return left.NotEquals(right)
	case token.TT_LT:
		return left.LessThan(right)
	case token.TT_LTE:
		return left.LessThanEq(right)
	case token.TT_GT:
		return left.GreaterThan(right)
	case token.TT_GTE:
		return left.GreaterThanEq(right)
	}
	return types.TypeValue{}, fmt.Errorf("invalid binary op: %s", node.Op.Type)
}

func (i *Interpreter) visitUnaryOpNode(node parser.UnaryOpNode) (types.TypeValue, error) {
	val, err := i.visit(node.Operand)
	if err != nil {
		return types.TypeValue{}, err
	}

	if node.Op.Type == token.TT_MINUS || node.Op.Type == token.TT_NOT {
		return val.Negate()
	}

	return types.TypeValue{}, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (i *Interpreter) visitNumberNode(node parser.NumberNode) (types.TypeValue, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return types.TypeValue{}, err
	}

	return types.TypeValue{Type: types.NUMBER, Value: val}, nil
}

func (i *Interpreter) visitBooleanNode(node parser.BooleanNode) (types.TypeValue, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return types.TypeValue{}, err
	}

	return types.TypeValue{Type: types.BOOL, Value: val}, nil
}

func (i *Interpreter) visitIdentifierNode(node parser.IdentifierNode) (types.TypeValue, error) {
	val, ok := i.globals[node.Token.Literal]

	if !ok {
		return types.TypeValue{}, fmt.Errorf("undeclared identifier: %s", node.Token.Literal)
	}
	return val, nil
}

func (i *Interpreter) visitNilNode(node parser.NilNode) (types.TypeValue, error) {
	return types.TypeValue{Type: types.NIL, Value: nil}, nil
}
