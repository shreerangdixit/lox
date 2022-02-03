package interpreter

import (
	"fmt"
	"lox/parser"
	"lox/token"
	"lox/types"
	"strconv"
)

type Interpreter struct {
	ast parser.Node
}

func New(parser *parser.Parser) (*Interpreter, error) {
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return &Interpreter{
		ast: ast,
	}, nil
}

func (i *Interpreter) Run() (types.ExpressionValue, error) {
	return i.visit(i.ast)
}

func (i *Interpreter) visit(node parser.Node) (types.ExpressionValue, error) {
	switch node.(type) {
	case parser.NumberNode:
		return i.visitNumberNode(node.(parser.NumberNode))
	case parser.BooleanNode:
		return i.visitBooleanNode(node.(parser.BooleanNode))
	case parser.BinaryOpNode:
		return i.visitBinaryOpNode(node.(parser.BinaryOpNode))
	case parser.UnaryOpNode:
		return i.visitUnaryOpNode(node.(parser.UnaryOpNode))
	}
	return types.ExpressionValue{}, fmt.Errorf("invalid node")
}

func (i *Interpreter) visitNumberNode(node parser.NumberNode) (types.ExpressionValue, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return types.ExpressionValue{}, err
	}

	return types.ExpressionValue{Type: types.TYPE_NUMBER, Value: val}, nil
}

func (i *Interpreter) visitBooleanNode(node parser.BooleanNode) (types.ExpressionValue, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return types.ExpressionValue{}, err
	}

	return types.ExpressionValue{Type: types.TYPE_BOOLEAN, Value: val}, nil
}

func (i *Interpreter) visitBinaryOpNode(node parser.BinaryOpNode) (types.ExpressionValue, error) {
	left, err := i.visit(node.Left)
	if err != nil {
		return types.ExpressionValue{}, err
	}

	right, err := i.visit(node.Right)
	if err != nil {
		return types.ExpressionValue{}, err
	}

	switch node.Token.Type {
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
	return types.ExpressionValue{}, fmt.Errorf("invalid binary op: %s", node.Token.Type)
}

func (i *Interpreter) visitUnaryOpNode(node parser.UnaryOpNode) (types.ExpressionValue, error) {
	val, err := i.visit(node.Node)
	if err != nil {
		return types.ExpressionValue{}, err
	}

	if node.Token.Type == token.TT_MINUS || node.Token.Type == token.TT_NOT {
		return val.Negate()
	}

	return types.ExpressionValue{}, fmt.Errorf("invalid unary op: %s", node.Token.Type)
}
