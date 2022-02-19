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

type NilNode struct {
}

type ProgramNode struct {
	Declarations []Node
}

type IdentifierNode struct {
	Token token.Token
}

type AssignmentNode struct {
	Identifier IdentifierNode
	Value      Node
}

type VarStmtNode struct {
	Identifier IdentifierNode
	Value      Node
}

type ExpStmtNode struct {
	Exp Node
}

type IfStmtNode struct {
	Exp       Node
	TrueStmt  Node
	FalseStmt Node
}

type WhileStmtNode struct {
	Condition Node
	Body      Node
}

type BreakStmtNode struct {
	Token token.Token
}

type BlockNode struct {
	Declarations []Node
}

type ExpNode struct {
	Exp Node
}

type TernaryOpNode struct {
	Exp      Node
	TrueExp  Node
	FalseExp Node
}

type BinaryOpNode struct {
	LeftExp  Node
	Op       token.Token
	RightExp Node
}

type UnaryOpNode struct {
	Op      token.Token
	Operand Node
}

type LogicalAndNode struct {
	LHS Node
	RHS Node
}

type LogicalOrNode struct {
	LHS Node
	RHS Node
}

type BooleanNode struct {
	Token token.Token
}

type NumberNode struct {
	Token token.Token
}

type StringNode struct {
	Token token.Token
}

type ListNode struct {
	Elements []Node
}

type CallNode struct {
	Callee    Node
	Arguments []Node
}

type IndexOfNode struct {
	Sequence Node
	Index    Node
}

func (n NilNode) String() string        { return "nil" }
func (n ProgramNode) String() string    { return fmt.Sprintf("%s", n.Declarations) }
func (n IdentifierNode) String() string { return n.Token.String() }
func (n AssignmentNode) String() string { return fmt.Sprintf("%s=%s", n.Identifier, n.Value) }
func (n VarStmtNode) String() string    { return fmt.Sprintf("var %s=%s", n.Identifier, n.Value) }
func (n BlockNode) String() string      { return fmt.Sprintf("{%s}", n.Declarations) }
func (n ExpStmtNode) String() string    { return n.Exp.String() }
func (n IfStmtNode) String() string {
	if n.FalseStmt != nil {
		return fmt.Sprintf("if(%s) %s else %s", n.Exp, n.TrueStmt, n.FalseStmt)
	} else {
		return fmt.Sprintf("if(%s) %s", n.Exp, n.TrueStmt)
	}
}
func (n WhileStmtNode) String() string { return fmt.Sprintf("while(%s)%s", n.Condition, n.Body) }
func (n BreakStmtNode) String() string { return n.Token.String() }
func (n ExpNode) String() string       { return n.Exp.String() }
func (n TernaryOpNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", n.Exp, n.TrueExp, n.FalseExp)
}
func (n LogicalAndNode) String() string { return fmt.Sprintf("%s && %s", n.LHS, n.RHS) }
func (n LogicalOrNode) String() string  { return fmt.Sprintf("%s || %s", n.LHS, n.RHS) }
func (n BinaryOpNode) String() string   { return fmt.Sprintf("%s %s %s", n.LeftExp, n.Op, n.RightExp) }
func (n UnaryOpNode) String() string    { return fmt.Sprintf("%s%s", n.Op, n.Operand) }
func (n BooleanNode) String() string    { return n.Token.String() }
func (n NumberNode) String() string     { return n.Token.String() }
func (n StringNode) String() string     { return n.Token.String() }
func (n CallNode) String() string       { return fmt.Sprintf("func %s(%s)", n.Callee, n.Arguments) }
func (n ListNode) String() string       { return fmt.Sprintf("%s", n.Elements) }
func (n IndexOfNode) String() string    { return fmt.Sprintf("%s[%s]", n.Sequence, n.Index) }
