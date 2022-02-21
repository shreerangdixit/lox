package runtime

import (
	"fmt"
	"strconv"

	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/token"
)

type Evaluator struct {
	env *Environment
}

func NewEvaluator() *Evaluator {
	return &Evaluator{env: NewEnvironment()}
}

func (e *Evaluator) Evaluate(root ast.Node) (Object, error) {
	return e.eval(root)
}

func (e *Evaluator) eval(node ast.Node) (Object, error) {
	switch node := node.(type) {
	case ast.ProgramNode:
		return e.evalProgramNode(node)
	case ast.BlockNode:
		return e.evalBlockNode(node)
	case ast.VarStmtNode:
		return e.evalVarStmtNode(node)
	case ast.ExpStmtNode:
		return e.evalExpStmtNode(node)
	case ast.IfStmtNode:
		return e.evalIfStmtNode(node)
	case ast.WhileStmtNode:
		return e.evalWhileStmtNode(node)
	case ast.BreakStmtNode:
		return e.evalBreakStmtNode(node)
	case ast.ReturnStmtNode:
		return e.evalReturnStmtNode(node)
	case ast.AssignmentNode:
		return e.evalAssignmentNode(node)
	case ast.LogicalAndNode:
		return e.evalLogicalAndNode(node)
	case ast.LogicalOrNode:
		return e.evalLogicalOrNode(node)
	case ast.ExpNode:
		return e.eval(node.Exp)
	case ast.TernaryOpNode:
		return e.evalTernaryOpNode(node)
	case ast.BinaryOpNode:
		return e.evalBinaryOpNode(node)
	case ast.UnaryOpNode:
		return e.evalUnaryOpNode(node)
	case ast.IdentifierNode:
		return e.evalIdentifierNode(node)
	case ast.NumberNode:
		return e.evalNumberNode(node)
	case ast.BooleanNode:
		return e.evalBooleanNode(node)
	case ast.StringNode:
		return e.evalStringNode(node)
	case ast.ListNode:
		return e.evalListNode(node)
	case ast.MapNode:
		return e.evalMapNode(node)
	case ast.NilNode:
		return e.evalNilNode(node)
	case ast.CallNode:
		return e.evalCallNode(node)
	case ast.IndexOfNode:
		return e.evalIndexOfNode(node)
	case ast.FunctionNode:
		return e.evalFunctionNode(node)
	}
	return NIL, fmt.Errorf("invalid node: %T", node)
}

func (e *Evaluator) evalProgramNode(node ast.ProgramNode) (Object, error) {
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalBlockNode(node ast.BlockNode) (Object, error) {
	return e.evalBlockNodeWithEnv(node, NewEnvironment().WithEnclosing(e.env))
}

func (e *Evaluator) evalBlockNodeWithEnv(node ast.BlockNode, env *Environment) (Object, error) {
	// Reset environment at the end of block scope
	prev := e.env
	defer func() {
		e.env = prev
	}()

	// New environment at the beginning of block scope
	e.env = env
	for _, node := range node.Declarations {
		_, err := e.eval(node)
		if err != nil {
			return NIL, err
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalVarStmtNode(node ast.VarStmtNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}

	if err := e.env.Declare(node.Identifier.Token.Literal, value); err != nil {
		return NIL, err
	}
	return NIL, nil
}

func (e *Evaluator) evalExpStmtNode(node ast.ExpStmtNode) (Object, error) {
	return e.eval(node.Exp)
}

func (e *Evaluator) evalIfStmtNode(node ast.IfStmtNode) (Object, error) {
	value, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(value) {
		return e.eval(node.TrueStmt)
	} else if node.FalseStmt != nil {
		return e.eval(node.FalseStmt)
	} else {
		return NIL, nil
	}
}

func (e *Evaluator) evalWhileStmtNode(node ast.WhileStmtNode) (Object, error) {
	for {
		condition, err := e.eval(node.Condition)
		if err != nil {
			return NIL, err
		}

		if !IsTruthy(condition) {
			break
		}

		_, err = e.eval(node.Body)
		if err != nil {
			switch err := err.(type) {
			case BreakError:
				return NIL, nil
			default:
				return NIL, err
			}
		}
	}
	return NIL, nil
}

func (e *Evaluator) evalBreakStmtNode(node ast.BreakStmtNode) (Object, error) {
	return NIL, NewBreakError()
}

func (e *Evaluator) evalAssignmentNode(node ast.AssignmentNode) (Object, error) {
	value, err := e.eval(node.Value)
	if err != nil {
		return NIL, err
	}
	return NIL, e.env.Assign(node.Identifier.Token.Literal, value)
}

func (e *Evaluator) evalLogicalAndNode(node ast.LogicalAndNode) (Object, error) {
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

func (e *Evaluator) evalLogicalOrNode(node ast.LogicalOrNode) (Object, error) {
	left, err := e.eval(node.LHS)
	if err != nil {
		return NIL, err
	}

	if IsTruthy(left) {
		return TRUE, nil
	}

	right, err := e.eval(node.RHS)
	if err != nil {
		return NIL, err
	}

	return NewBool(IsTruthy(right)), nil
}

func (e *Evaluator) evalTernaryOpNode(node ast.TernaryOpNode) (Object, error) {
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

func (e *Evaluator) evalBinaryOpNode(node ast.BinaryOpNode) (Object, error) {
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
	case token.TT_MODULO:
		return Modulo(left, right)
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

func (e *Evaluator) evalUnaryOpNode(node ast.UnaryOpNode) (Object, error) {
	val, err := e.eval(node.Operand)
	if err != nil {
		return NIL, err
	}

	if node.Op.Type == token.TT_MINUS {
		return Negate(val)
	} else if node.Op.Type == token.TT_NOT {
		return Not(val)
	}

	return NIL, fmt.Errorf("invalid unary op: %s", node.Op.Type)
}

func (e *Evaluator) evalIdentifierNode(node ast.IdentifierNode) (Object, error) {
	return e.env.Get(node.Token.Literal)
}

func (e *Evaluator) evalNumberNode(node ast.NumberNode) (Object, error) {
	val, err := strconv.ParseFloat(node.Token.Literal, 10)
	if err != nil {
		return NIL, err
	}

	return NewNumber(val), nil
}

func (e *Evaluator) evalBooleanNode(node ast.BooleanNode) (Object, error) {
	val, err := strconv.ParseBool(node.Token.Literal)
	if err != nil {
		return NIL, err
	}

	return NewBool(val), nil
}

func (e *Evaluator) evalStringNode(node ast.StringNode) (Object, error) {
	return NewString(node.Token.Literal), nil
}

func (e *Evaluator) evalListNode(node ast.ListNode) (Object, error) {
	elements, err := e.evalNodes(node.Elements)
	if err != nil {
		return nil, err
	}

	return NewList(elements), nil
}

func (e *Evaluator) evalMapNode(node ast.MapNode) (Object, error) {
	m := NewMap()

	for _, kvp := range node.Elements {
		key, err := e.eval(kvp.Key)
		if err != nil {
			return NIL, err
		}

		value, err := e.eval(kvp.Value)
		if err != nil {
			return NIL, err
		}

		m, err = m.Add(key, value)
		if err != nil {
			return NIL, err
		}
	}
	return m, nil
}

func (e *Evaluator) evalNilNode(node ast.NilNode) (Object, error) {
	return NIL, nil
}

func (e *Evaluator) evalCallNode(node ast.CallNode) (Object, error) {
	calleeNode, err := e.eval(node.Callee)
	if err != nil {
		return NIL, err
	}

	callable, ok := calleeNode.(Callable)
	if !ok { // If the callee node itself isn't callable, check if it's value is callable
		calleeValue, err := e.env.Get(calleeNode.String())
		if err != nil {
			return NIL, fmt.Errorf("%s is not callable", calleeNode.Type())
		}

		callable, ok = calleeValue.(Callable)
		if !ok {
			return NIL, fmt.Errorf("%s is not callable", calleeValue.Type())
		}
	}

	if !callable.Variadic() && callable.Arity() != len(node.Arguments) {
		return NIL, fmt.Errorf(
			"incorrect number of arguments to %s - %d expected %d provided",
			callable,
			callable.Arity(),
			len(node.Arguments),
		)
	}

	argValues, err := e.evalNodes(node.Arguments)
	if err != nil {
		return NIL, err
	}

	return callable.Call(e, argValues)
}

func (e *Evaluator) evalIndexOfNode(node ast.IndexOfNode) (Object, error) {
	seq, err := e.eval(node.Sequence)
	if err != nil {
		return nil, err
	}

	idx, err := e.eval(node.Index)
	if err != nil {
		return nil, err
	}

	return ItemAtIndex(seq, idx)
}

func (e *Evaluator) evalFunctionNode(node ast.FunctionNode) (Object, error) {
	fun := NewUserFunction(node, e.env)
	return NIL, e.env.Declare(fun.String(), fun)
}

func (e *Evaluator) evalReturnStmtNode(node ast.ReturnStmtNode) (Object, error) {
	val, err := e.eval(node.Exp)
	if err != nil {
		return NIL, err
	}

	return NIL, NewReturnError(val)
}

func (e *Evaluator) evalNodes(argNodes []ast.Node) ([]Object, error) {
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
