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
	switch node := node.(type) {
	case parser.ProgramNode:
		return i.evalProgramNode(node)
	case parser.BlockNode:
		return i.evalBlockNode(node)
	case parser.LetStmtNode:
		return i.evalLetStmtNode(node)
	case parser.ExpStmtNode:
		return i.evalExpStmtNode(node)
	case parser.IfStmtNode:
		return i.evalIfStmtNode(node)
	case parser.PrintStmtNode:
		return i.evalPrintStmtNode(node)
	case parser.WhileStmtNode:
		return i.evalWhileStmtNode(node)
	case parser.AssignmentNode:
		return i.evalAssignmentNode(node)
	case parser.LogicalAndNode:
		return i.evalLogicalAndNode(node)
	case parser.LogicalOrNode:
		return i.evalLogicalOrNode(node)
	case parser.ExpNode:
		return i.eval(node.Exp)
	case parser.TernaryOpNode:
		return i.evalTernaryOpNode(node)
	case parser.BinaryOpNode:
		return i.evalBinaryOpNode(node)
	case parser.UnaryOpNode:
		return i.evalUnaryOpNode(node)
	case parser.IdentifierNode:
		return i.evalIdentifierNode(node)
	case parser.NumberNode:
		return i.evalNumberNode(node)
	case parser.BooleanNode:
		return i.evalBooleanNode(node)
	case parser.StringNode:
		return i.evalStringNode(node)
	case parser.NilNode:
		return i.evalNilNode(node)
	case parser.CallNode:
		return i.evalCallNode(node)
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

func (i *Interpreter) evalWhileStmtNode(node parser.WhileStmtNode) (Object, error) {
	for {
		result, err := i.eval(node.Condition)
		if err != nil {
			return NIL, err
		}

		if !IsTruthy(result) {
			break
		}

		_, err = i.eval(node.Body)
		if err != nil {
			return NIL, err
		}
	}
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
		return EqualTo(left, right), nil
	case token.TT_NEQ:
		return NotEqualTo(left, right), nil
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

func (i *Interpreter) evalCallNode(node parser.CallNode) (Object, error) {
	argValues := make([]Object, 0, 255)
	for _, arg := range node.Arguments {
		argval, err := i.eval(arg)
		if err != nil {
			return NIL, err
		}

		argValues = append(argValues, argval)
	}
	// TODO: Evaluate function call
	fmt.Printf("Calling %s with arguments %s\n", node.Callee, argValues)
	return NIL, nil
}
