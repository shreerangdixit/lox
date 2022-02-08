package interpreter

import (
	"fmt"
	"github.com/shreerangdixit/lox/parser"
	"github.com/shreerangdixit/lox/token"
	"github.com/shreerangdixit/lox/types"
	"strconv"
)

type Interpreter struct {
	env *Env
}

func New() *Interpreter {
	return &Interpreter{env: NewEnv()}
}

func (i *Interpreter) Run(root parser.Node) (types.TypeValue, error) {
	return i.eval(root)
}

func (i *Interpreter) eval(node parser.Node) (types.TypeValue, error) {
	switch node.(type) {
	case parser.ProgramNode:
		return i.evalProgramNode(node.(parser.ProgramNode))
	case parser.BlockNode:
		return i.evalBlockNode(node.(parser.BlockNode))
	case parser.LetStatementNode:
		return i.evalLetStatementNode(node.(parser.LetStatementNode))
	case parser.ExpressionStatementNode:
		return i.evalExpressionStatementNode(node.(parser.ExpressionStatementNode))
	case parser.IfStatementNode:
		return i.evalIfStatementNode(node.(parser.IfStatementNode))
	case parser.PrintStatementNode:
		return i.evalPrintStatementNode(node.(parser.PrintStatementNode))
	case parser.AssignmentNode:
		return i.evalAssignmentNode(node.(parser.AssignmentNode))
	case parser.ExpressionNode:
		return i.eval(node.(parser.ExpressionNode).Exp)
	case parser.TernaryOpNode:
		return i.evalTernaryOpNode(node.(parser.TernaryOpNode))
	case parser.BinaryOpNode:
		return i.evalBinaryOpNode(node.(parser.BinaryOpNode))
	case parser.UnaryOpNode:
		return i.evalUnaryOpNode(node.(parser.UnaryOpNode))
	case parser.IdentifierNode:
		return i.evalIdentifierNode(node.(parser.IdentifierNode))
	case parser.NumberNode:
		return i.evalNumberNode(node.(parser.NumberNode))
	case parser.BooleanNode:
		return i.evalBooleanNode(node.(parser.BooleanNode))
	case parser.StringNode:
		return i.evalStringNode(node.(parser.StringNode))
	case parser.NilNode:
		return i.evalNilNode(node.(parser.NilNode))
	}
	return types.NO_VALUE, fmt.Errorf("invalid node: %T", node)
}

func (i *Interpreter) evalProgramNode(node parser.ProgramNode) (types.TypeValue, error) {
	for _, node := range node.Declarations {
		_, err := i.eval(node)
		if err != nil {
			return types.NO_VALUE, err
		}
	}
	return types.NO_VALUE, nil
}

func (i *Interpreter) evalBlockNode(node parser.BlockNode) (types.TypeValue, error) {
	// Reset environment at the end of block scope
	prev := i.env
	defer func() {
		i.env = prev
	}()

	// New environment at the beginning of block scope
	i.env = NewEnvWithEnclosing(i.env)
	for _, node := range node.Declarations {
		_, err := i.eval(node)
		if err != nil {
			return types.NO_VALUE, err
		}
	}
	return types.NO_VALUE, nil
}

func (i *Interpreter) evalLetStatementNode(node parser.LetStatementNode) (types.TypeValue, error) {
	value, err := i.eval(node.Value)
	if err != nil {
		return types.NO_VALUE, err
	}

	if err := i.env.Declare(node.Identifier.Token.Literal, value); err != nil {
		return types.NO_VALUE, err
	}
	return types.NO_VALUE, nil
}

func (i *Interpreter) evalExpressionStatementNode(node parser.ExpressionStatementNode) (types.TypeValue, error) {
	return i.eval(node.Exp)
}

func (i *Interpreter) evalIfStatementNode(node parser.IfStatementNode) (types.TypeValue, error) {
	value, err := i.eval(node.Exp)
	if err != nil {
		return types.NO_VALUE, err
	}

	if value.Type != types.BOOL {
		return types.NO_VALUE, fmt.Errorf("expected if condition to evaluate to boolean")
	}

	if value.Value.(bool) {
		return i.eval(node.True)
	} else {
		return i.eval(node.False)
	}
}

func (i *Interpreter) evalPrintStatementNode(node parser.PrintStatementNode) (types.TypeValue, error) {
	result, err := i.eval(node.Exp)
	if err != nil {
		return types.NO_VALUE, err
	}

	fmt.Printf("%v\n", result.Value)

	return types.NO_VALUE, nil
}

func (i *Interpreter) evalAssignmentNode(node parser.AssignmentNode) (types.TypeValue, error) {
	value, err := i.eval(node.Value)
	if err != nil {
		return types.NO_VALUE, err
	}
	return types.NO_VALUE, i.env.Assign(node.Identifier.Token.Literal, value)
}

func (i *Interpreter) evalTernaryOpNode(node parser.TernaryOpNode) (types.TypeValue, error) {
	value, err := i.eval(node.Exp)
	if err != nil {
		return types.NO_VALUE, err
	}

	if value.Type != types.BOOL {
		return types.NO_VALUE, fmt.Errorf("expected ternary condition to evaluate to boolean")
	}

	if value.Value.(bool) {
		return i.eval(node.TrueExp)
	} else {
		return i.eval(node.FalseExp)
	}
}

func (i *Interpreter) evalBinaryOpNode(node parser.BinaryOpNode) (types.TypeValue, error) {
	left, err := i.eval(node.LHS)
	if err != nil {
		return types.NO_VALUE, err
	}

	right, err := i.eval(node.RHS)
	if err != nil {
		return types.NO_VALUE, err
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
	return types.NO_VALUE, fmt.Errorf("invalid binary op: %s", node.Op.Type)
}

func (i *Interpreter) evalUnaryOpNode(node parser.UnaryOpNode) (types.TypeValue, error) {
	val, err := i.eval(node.Operand)
	if err != nil {
		return types.NO_VALUE, err
	}

	if node.Op.Type == token.TT_MINUS || node.Op.Type == token.TT_NOT {
		return val.Negate()
	}

	return types.NO_VALUE, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (i *Interpreter) evalIdentifierNode(node parser.IdentifierNode) (types.TypeValue, error) {
	return i.env.Get(node.Token.Literal)
}

func (i *Interpreter) evalNumberNode(node parser.NumberNode) (types.TypeValue, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return types.NO_VALUE, err
	}

	return types.TypeValue{Type: types.NUMBER, Value: val}, nil
}

func (i *Interpreter) evalBooleanNode(node parser.BooleanNode) (types.TypeValue, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return types.NO_VALUE, err
	}

	return types.TypeValue{Type: types.BOOL, Value: val}, nil
}

func (i *Interpreter) evalStringNode(node parser.StringNode) (types.TypeValue, error) {
	return types.TypeValue{Type: types.STRING, Value: node.Token.Literal}, nil
}

func (i *Interpreter) evalNilNode(node parser.NilNode) (types.TypeValue, error) {
	return types.TypeValue{Type: types.NIL, Value: nil}, nil
}
