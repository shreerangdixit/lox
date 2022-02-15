package runtime

import (
	"fmt"
	"github.com/shreerangdixit/lox/parser"
	"github.com/shreerangdixit/lox/token"
	"strconv"
)

type Evaluator struct {
	env *Env
}

func NewEvaluator() *Evaluator {
	return &Evaluator{env: NewEnv()}
}

func (e *Evaluator) Evaluate(root parser.Node) (Object, error) {
	return e.eval(root)
}

func (e *Evaluator) eval(node parser.Node) (Object, error) {
	switch node := node.(type) {
	case parser.ProgramNode:
		return e.evalProgramNode(node)
	case parser.BlockNode:
		return e.evalBlockNode(node)
	case parser.LetStmtNode:
		return e.evalLetStmtNode(node)
	case parser.ExpStmtNode:
		return e.evalExpStmtNode(node)
	case parser.IfStmtNode:
		return e.evalIfStmtNode(node)
	case parser.PrintStmtNode:
		return e.evalPrintStmtNode(node)
	case parser.WhileStmtNode:
		return e.evalWhileStmtNode(node)
	case parser.AssignmentNode:
		return e.evalAssignmentNode(node)
	case parser.LogicalAndNode:
		return e.evalLogicalAndNode(node)
	case parser.LogicalOrNode:
		return e.evalLogicalOrNode(node)
	case parser.ExpNode:
		return e.eval(node.Exp)
	case parser.TernaryOpNode:
		return e.evalTernaryOpNode(node)
	case parser.BinaryOpNode:
		return e.evalBinaryOpNode(node)
	case parser.UnaryOpNode:
		return e.evalUnaryOpNode(node)
	case parser.IdentifierNode:
		return e.evalIdentifierNode(node)
	case parser.NumberNode:
		return e.evalNumberNode(node)
	case parser.BooleanNode:
		return e.evalBooleanNode(node)
	case parser.StringNode:
		return e.evalStringNode(node)
	case parser.NilNode:
		return e.evalNilNode(node)
	case parser.CallNode:
		return e.evalCallNode(node)
	}
	return NIL, fmt.Errorf("invalid node: %T", node)
}

func (e *Evaluator) evalProgramNode(node parser.ProgramNode) (Object, error) {
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalBlockNode(node parser.BlockNode) (Object, error) {
	// Reset environment at the end of block scope
	prev := e.env
	defer func() {
		e.env = prev
	}()

	// New environment at the beginning of block scope
	e.env = NewEnvWithEnclosing(e.env)
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalLetStmtNode(node parser.LetStmtNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}

	if err := e.env.Declare(node.Identifier.Token.Literal, value); err != nil {
		return NIL, err
	}
	return NIL, nil
}

func (e *Evaluator) evalExpStmtNode(node parser.ExpStmtNode) (Object, error) {
	return e.eval(node.Exp)
}

func (e *Evaluator) evalIfStmtNode(node parser.IfStmtNode) (Object, error) {
	value, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return e.eval(node.TrueStmt)
	} else {
		return e.eval(node.FalseStmt)
	}
}

func (e *Evaluator) evalPrintStmtNode(node parser.PrintStmtNode) (Object, error) {
	result, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	fmt.Printf("%s\n", result)

	return NIL, nil
}

func (e *Evaluator) evalWhileStmtNode(node parser.WhileStmtNode) (Object, error) {
	for {
		result, err := e.eval(node.Condition)
		if err != nil {
			return NIL, err
		}

		if !IsTruthy(result) {
			break
		}

		_, err = e.eval(node.Body)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalAssignmentNode(node parser.AssignmentNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}
	return NIL, e.env.Assign(node.Identifier.Token.Literal, value)
}

func (e *Evaluator) evalLogicalAndNode(node parser.LogicalAndNode) (Object, error) {
	left, err := e.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	right, err := e.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(left) && IsTruthy(right)), nil
}

func (e *Evaluator) evalLogicalOrNode(node parser.LogicalOrNode) (Object, error) {
	left, err := e.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(left) {
		return NewBool(true), nil
	}

	right, err := e.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(right)), nil
}

func (e *Evaluator) evalTernaryOpNode(node parser.TernaryOpNode) (Object, error) {
	value, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return e.eval(node.TrueExp)
	} else {
		return e.eval(node.FalseExp)
	}
}

func (e *Evaluator) evalBinaryOpNode(node parser.BinaryOpNode) (Object, error) {
	left, err := e.eval(node.LeftExp)
	if err != nil {
		return NIL, err
	}

	right, err := e.eval(node.RightExp)
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

func (e *Evaluator) evalUnaryOpNode(node parser.UnaryOpNode) (Object, error) {
	val, err := e.eval(node.Operand)
	if err != nil {
		return NIL, err
	}

	if node.Op.Type == token.TT_MINUS || node.Op.Type == token.TT_NOT {
		return Negate(val)
	}

	return NIL, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (e *Evaluator) evalIdentifierNode(node parser.IdentifierNode) (Object, error) {
	return e.env.Get(node.Token.Literal)
}

func (e *Evaluator) evalNumberNode(node parser.NumberNode) (Object, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return NIL, err
	}

	return NewFloat64(val), nil
}

func (e *Evaluator) evalBooleanNode(node parser.BooleanNode) (Object, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return NIL, err
	}

	return NewBool(val), nil
}

func (e *Evaluator) evalStringNode(node parser.StringNode) (Object, error) {
	return NewString(node.Token.Literal), nil
}

func (e *Evaluator) evalNilNode(node parser.NilNode) (Object, error) {
	return NIL, nil
}

func (e *Evaluator) evalCallNode(node parser.CallNode) (Object, error) {
	callee, err := e.eval(node.Callee)
	if err != nil {
		return NIL, fmt.Errorf("%s is not declared", node.Callee)
	}

	calleeValue, err := e.env.Get(callee.String())
	if err != nil {
		return NIL, fmt.Errorf("%s is not callable", callee.Type())
	}

	callable, ok := calleeValue.(Callable)
	if !ok {
		return NIL, fmt.Errorf("%s is not declared", calleeValue.Type())
	}

	if callable.Arity() != len(node.Arguments) {
		return NIL, fmt.Errorf(
			"incorrect number of arguments to %s - %d expected %d provided",
			callable,
			callable.Arity(),
			len(node.Arguments),
		)
	}

	argValues, err := e.makeCallArguments(node.Arguments)
	if err != nil {
		return NIL, err
	}

	return callable.Call(e, argValues)
}

func (e *Evaluator) makeCallArguments(argNodes []parser.Node) ([]Object, error) {
	argValues := make([]Object, 0, 255)
	for _, arg := range argNodes {
		argval, err := e.eval(arg)
		if err != nil {
			return []Object{}, err
		}

		argValues = append(argValues, argval)
	}
	return argValues, nil
}
