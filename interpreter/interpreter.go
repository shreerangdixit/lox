package interpreter

import (
	_ "fmt"
	"lox/parser"
	"lox/token"
	"strconv"
)

type NumberValue struct {
	Token token.Token
	Value float64
}

func (n NumberValue) Add(right NumberValue) NumberValue {
	return NumberValue{Value: n.Value + right.Value}
}

func (n NumberValue) Subtract(right NumberValue) NumberValue {
	return NumberValue{Value: n.Value - right.Value}
}

func (n NumberValue) Divide(right NumberValue) NumberValue {
	return NumberValue{Value: n.Value / right.Value}
}

func (n NumberValue) Multiply(right NumberValue) NumberValue {
	return NumberValue{Value: n.Value * right.Value}
}

func (n NumberValue) Negate() NumberValue {
	return NumberValue{Value: n.Value * -1}
}

// type BooleanValue struct {
// 	Token token.Token
// 	Value bool
// }

// func (b BooleanValue) Equal(right BooleanValue) bool {
// 	return b.Value == right.Value
// }

// func (b BooleanValue) Negate() bool {
// 	return !b.Value
// }

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

func (i *Interpreter) Run() NumberValue {
	return i.visit(i.ast)
}

func (i *Interpreter) visit(node parser.Node) NumberValue {
	switch node.(type) {
	case parser.NumberNode:
		return i.visitNumberNode(node.(parser.NumberNode))
	// case parser.BooleanNode:
	// 	return i.visitBooleanNode(node.(parser.BooleanNode))
	case parser.BinaryOpNode:
		return i.visitBinaryOpNode(node.(parser.BinaryOpNode))
	case parser.UnaryOpNode:
		return i.visitUnaryOpNode(node.(parser.UnaryOpNode))
	}
	return NumberValue{}
}

func (i *Interpreter) visitNumberNode(node parser.NumberNode) NumberValue {
	val, _ := strconv.ParseFloat(node.Token.Literal, 10)
	return NumberValue{
		Token: node.Token,
		Value: val,
	}
}

// func (i *Interpreter) visitBooleanNode(node parser.BooleanNode) BooleanValue {
// 	val, _ := strconv.ParseBool(node.Token.Literal)
// 	return BooleanValue{
// 		Token: node.Token,
// 		Value: val,
// 	}
// }

func (i *Interpreter) visitBinaryOpNode(node parser.BinaryOpNode) NumberValue {
	left := i.visit(node.Left)
	right := i.visit(node.Right)
	switch node.Token.Type {
	case token.TT_PLUS:
		return left.Add(right)
	case token.TT_MINUS:
		return left.Subtract(right)
	case token.TT_DIVIDE:
		return left.Divide(right)
	case token.TT_MULTIPLY:
		return left.Multiply(right)
	}
	return NumberValue{}
}

func (i *Interpreter) visitUnaryOpNode(node parser.UnaryOpNode) NumberValue {
	val := i.visit(node.Node)
	if node.Token.Type == token.TT_MINUS {
		return val.Negate()
	}
	return NumberValue{}
}
