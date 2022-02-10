package interpreter

import (
	"fmt"
	"github.com/shreerangdixit/lox/parser"
	"github.com/shreerangdixit/lox/token"
	"strconv"
)

type Interpreter struct {
	env *Env
}

func New() *Interpreter {
	return &Interpreter{env: NewEnv()}
}

func (i *Interpreter) Run(root parser.Node) (Object, error) {
	return i.eval(root)
}

func (i *Interpreter) eval(node parser.Node) (Object, error) {
	switch node.(type) {
	case parser.ProgramNode:
		return i.evalProgramNode(node.(parser.ProgramNode))
	case parser.BlockNode:
		return i.evalBlockNode(node.(parser.BlockNode))
	case parser.LetStmtNode:
		return i.evalLetStmtNode(node.(parser.LetStmtNode))
	case parser.ExpStmtNode:
		return i.evalExpStmtNode(node.(parser.ExpStmtNode))
	case parser.IfStmtNode:
		return i.evalIfStmtNode(node.(parser.IfStmtNode))
	case parser.PrintStmtNode:
		return i.evalPrintStmtNode(node.(parser.PrintStmtNode))
	case parser.AssignmentNode:
		return i.evalAssignmentNode(node.(parser.AssignmentNode))
	case parser.LogicalAndNode:
		return i.evalLogicalAndNode(node.(parser.LogicalAndNode))
	case parser.LogicalOrNode:
		return i.evalLogicalOrNode(node.(parser.LogicalOrNode))
	case parser.ExpNode:
		return i.eval(node.(parser.ExpNode).Exp)
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
	return NIL, fmt.Errorf("invalid node: %T", node)
}

func (i *Interpreter) evalProgramNode(node parser.ProgramNode) (Object, error) {
	for _, node := range node.Declarations {
		_, err := i.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (i *Interpreter) evalBlockNode(node parser.BlockNode) (Object, error) {
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
			return NIL, err
		}
	}
	return NIL, nil
}

func (i *Interpreter) evalLetStmtNode(node parser.LetStmtNode) (Object, error) {
	value, err := i.eval(node.Value)
	if err != nil {
		return NIL, err
	}

	if err := i.env.Declare(node.Identifier.Token.Literal, value); err != nil {
		return NIL, err
	}
	return NIL, nil
}

func (i *Interpreter) evalExpStmtNode(node parser.ExpStmtNode) (Object, error) {
	return i.eval(node.Exp)
}

func (i *Interpreter) evalIfStmtNode(node parser.IfStmtNode) (Object, error) {
	value, err := i.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return i.eval(node.TrueStmt)
	} else {
		return i.eval(node.FalseStmt)
	}
}

func (i *Interpreter) evalPrintStmtNode(node parser.PrintStmtNode) (Object, error) {
	result, err := i.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	fmt.Printf("%s\n", result)

	return NIL, nil
}

func (i *Interpreter) evalAssignmentNode(node parser.AssignmentNode) (Object, error) {
	value, err := i.eval(node.Value)
	if err != nil {
		return NIL, err
	}
	return NIL, i.env.Assign(node.Identifier.Token.Literal, value)
}

func (i *Interpreter) evalLogicalAndNode(node parser.LogicalAndNode) (Object, error) {
	left, err := i.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	right, err := i.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(left) && IsTruthy(right)), nil
}

func (i *Interpreter) evalLogicalOrNode(node parser.LogicalOrNode) (Object, error) {
	left, err := i.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(left) {
		return NewBool(true), nil
	}

	right, err := i.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(right)), nil
}

func (i *Interpreter) evalTernaryOpNode(node parser.TernaryOpNode) (Object, error) {
	value, err := i.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return i.eval(node.TrueExp)
	} else {
		return i.eval(node.FalseExp)
	}
}

func (i *Interpreter) evalBinaryOpNode(node parser.BinaryOpNode) (Object, error) {
	left, err := i.eval(node.LeftExp)
	if err != nil {
		return NIL, err
	}

	right, err := i.eval(node.RightExp)
	if err != nil {
		return NIL, err
	}

	switch node.Op.Type {
	case token.TT_PLUS:
		return Add(left, right)
	case token.TT_MINUS:
		return Subtract(left, right)
	case token.TT_DIVIDE:
		return Divide(left, right)
	case token.TT_MULTIPLY:
		return Multiply(left, right)
	case token.TT_EQ:
		return Equals(left, right), nil
	case token.TT_NEQ:
		return NotEquals(left, right), nil
	case token.TT_LT:
		return LessThan(left, right), nil
	case token.TT_LTE:
		return LessThanEq(left, right), nil
	case token.TT_GT:
		return GreaterThan(left, right), nil
	case token.TT_GTE:
		return GreaterThanEq(left, right), nil
	}
	return NIL, fmt.Errorf("invalid binary op: %s", node.Op.Type)
}

func (i *Interpreter) evalUnaryOpNode(node parser.UnaryOpNode) (Object, error) {
	val, err := i.eval(node.Operand)
	if err != nil {
		return NIL, err
	}

	if node.Op.Type == token.TT_MINUS || node.Op.Type == token.TT_NOT {
		return Negate(val)
	}

	return NIL, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (i *Interpreter) evalIdentifierNode(node parser.IdentifierNode) (Object, error) {
	return i.env.Get(node.Token.Literal)
}

func (i *Interpreter) evalNumberNode(node parser.NumberNode) (Object, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return NIL, err
	}

	return NewFloat64(val), nil
}

func (i *Interpreter) evalBooleanNode(node parser.BooleanNode) (Object, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return NIL, err
	}

	return NewBool(val), nil
}

func (i *Interpreter) evalStringNode(node parser.StringNode) (Object, error) {
	return NewString(node.Token.Literal), nil
}

func (i *Interpreter) evalNilNode(node parser.NilNode) (Object, error) {
	return NIL, nil
}
