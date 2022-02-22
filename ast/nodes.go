package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/token"
)

// ------------------------------------
// Nodes
// ------------------------------------

type Node interface {
	String() string
}

type NilNode struct{}

func (n NilNode) String() string {
	return "nil"
}

type ProgramNode struct {
	Declarations []Node
}

func (n ProgramNode) String() string {
	str := ""
	for _, decl := range n.Declarations {
		str += decl.String()
		str += "\n"
	}

	return str
}

type IdentifierNode struct {
	Token token.Token
}

func (n IdentifierNode) String() string {
	return n.Token.Literal
}

type AssignmentNode struct {
	Identifier IdentifierNode
	Value      Node
}

func (n AssignmentNode) String() string {
	return fmt.Sprintf("%s=%s", n.Identifier, n.Value)
}

type VarStmtNode struct {
	Identifier IdentifierNode
	Value      Node
}

func (n VarStmtNode) String() string {
	return fmt.Sprintf("var %s=%s", n.Identifier, n.Value)
}

type ExpStmtNode struct {
	Exp Node
}

func (n ExpStmtNode) String() string {
	return n.Exp.String()
}

type IfStmtNode struct {
	Exp       Node
	TrueStmt  Node
	FalseStmt Node
}

func (n IfStmtNode) String() string {
	if n.FalseStmt != nil {
		return fmt.Sprintf("if(%s)\n\t%s\nelse\t%s", n.Exp, n.TrueStmt, n.FalseStmt)
	} else {
		return fmt.Sprintf("if(%s)\n\t%s", n.Exp, n.TrueStmt)
	}
}

type WhileStmtNode struct {
	Condition Node
	Body      Node
}

func (n WhileStmtNode) String() string {
	return fmt.Sprintf("while(%s)\n\t%s", n.Condition, n.Body)
}

type BreakStmtNode struct {
	Token token.Token
}

func (n BreakStmtNode) String() string {
	return n.Token.String()
}

type ContinueStmtNode struct {
	Token token.Token
}

func (n ContinueStmtNode) String() string {
	return n.Token.String()
}

type ReturnStmtNode struct {
	Exp Node
}

func (n ReturnStmtNode) String() string {
	if n.Exp != nil {
		return n.Exp.String()
	}
	return "return"
}

type BlockNode struct {
	Declarations []Node
}

func (n BlockNode) String() string {
	return fmt.Sprintf("{\n\t%s\n}", n.Declarations)
}

type ExpNode struct {
	Exp Node
}

func (n ExpNode) String() string {
	return n.Exp.String()
}

type TernaryOpNode struct {
	Exp      Node
	TrueExp  Node
	FalseExp Node
}

func (n TernaryOpNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", n.Exp, n.TrueExp, n.FalseExp)
}

type BinaryOpNode struct {
	LeftExp  Node
	Op       token.Token
	RightExp Node
}

func (n BinaryOpNode) String() string {
	return fmt.Sprintf("%s %s %s", n.LeftExp, n.Op, n.RightExp)
}

type UnaryOpNode struct {
	Op      token.Token
	Operand Node
}

func (n UnaryOpNode) String() string {
	return fmt.Sprintf("%s%s", n.Op, n.Operand)
}

type LogicalAndNode struct {
	LHS Node
	RHS Node
}

func (n LogicalAndNode) String() string {
	return fmt.Sprintf("%s && %s", n.LHS, n.RHS)
}

type LogicalOrNode struct {
	LHS Node
	RHS Node
}

func (n LogicalOrNode) String() string {
	return fmt.Sprintf("%s || %s", n.LHS, n.RHS)
}

type BooleanNode struct {
	Token token.Token
}

func (n BooleanNode) String() string {
	return n.Token.String()
}

type NumberNode struct {
	Token token.Token
}

func (n NumberNode) String() string {
	return n.Token.Literal
}

type StringNode struct {
	Token token.Token
}

func (n StringNode) String() string {
	return n.Token.Literal
}

type ListNode struct {
	Elements []Node
}

func (n ListNode) String() string {
	str := ""
	for _, elem := range n.Elements {
		str += elem.String()
		str += " "
	}
	return str
}

type CallNode struct {
	Callee    Node
	Arguments []Node
}

func (n CallNode) String() string {
	return fmt.Sprintf("fun %s(%s)", n.Callee, n.Arguments)
}

type IndexOfNode struct {
	Sequence Node
	Index    Node
}

func (n IndexOfNode) String() string {
	return fmt.Sprintf("%s[%s]", n.Sequence, n.Index)
}

type FunctionNode struct {
	Identifier IdentifierNode
	Parameters []IdentifierNode
	Body       BlockNode
}

func (n FunctionNode) String() string {
	return fmt.Sprintf("fun %s(%s)\n%s", n.Identifier, n.Parameters, n.Body)
}

type KeyValueNode struct {
	Key   Node
	Value Node
}

func (n KeyValueNode) String() string {
	return fmt.Sprintf("%s:%s", n.Key.String(), n.Value.String())
}

type MapNode struct {
	Elements []KeyValueNode
}

func (n MapNode) String() string {
	return fmt.Sprintf("{%s}", n.Elements)
}
