package ast

import (
	"fmt"

	"github.com/shreerangdixit/lox/lex"
)

// ------------------------------------
// Nodes
// ------------------------------------

type Node interface {
	String() string
	Begin() lex.Position
	End() lex.Position
}

type NilNode struct {
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n NilNode) String() string      { return "nil" }
func (n NilNode) Begin() lex.Position { return n.BeginPos }
func (n NilNode) End() lex.Position   { return n.EndPos }

type ProgramNode struct {
	Declarations []Node
	BeginPos     lex.Position
	EndPos       lex.Position
}

func (n ProgramNode) Begin() lex.Position { return n.BeginPos }
func (n ProgramNode) End() lex.Position   { return n.EndPos }

func (n ProgramNode) String() string {
	str := ""
	for _, decl := range n.Declarations {
		str += decl.String()
		str += "\n"
	}

	return str
}

type IdentifierNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n IdentifierNode) Begin() lex.Position { return n.BeginPos }
func (n IdentifierNode) End() lex.Position   { return n.EndPos }
func (n IdentifierNode) String() string      { return n.Token.Literal }

type AssignmentNode struct {
	Identifier IdentifierNode
	Value      Node
	BeginPos   lex.Position
	EndPos     lex.Position
}

func (n AssignmentNode) Begin() lex.Position { return n.BeginPos }
func (n AssignmentNode) End() lex.Position   { return n.EndPos }
func (n AssignmentNode) String() string      { return fmt.Sprintf("%s=%s", n.Identifier, n.Value) }

type VarStmtNode struct {
	Identifier IdentifierNode
	Value      Node
	BeginPos   lex.Position
	EndPos     lex.Position
}

func (n VarStmtNode) Begin() lex.Position { return n.BeginPos }
func (n VarStmtNode) End() lex.Position   { return n.EndPos }
func (n VarStmtNode) String() string      { return fmt.Sprintf("var %s=%s", n.Identifier, n.Value) }

type ExpStmtNode struct {
	Exp      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n ExpStmtNode) Begin() lex.Position { return n.BeginPos }
func (n ExpStmtNode) End() lex.Position   { return n.EndPos }
func (n ExpStmtNode) String() string      { return n.Exp.String() }

type IfStmtNode struct {
	Exp       Node
	TrueStmt  Node
	FalseStmt Node
	BeginPos  lex.Position
	EndPos    lex.Position
}

func (n IfStmtNode) Begin() lex.Position { return n.BeginPos }
func (n IfStmtNode) End() lex.Position   { return n.EndPos }

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
	BeginPos  lex.Position
	EndPos    lex.Position
}

func (n WhileStmtNode) Begin() lex.Position { return n.BeginPos }
func (n WhileStmtNode) End() lex.Position   { return n.EndPos }
func (n WhileStmtNode) String() string      { return fmt.Sprintf("while(%s)\n\t%s", n.Condition, n.Body) }

type BreakStmtNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n BreakStmtNode) Begin() lex.Position { return n.BeginPos }
func (n BreakStmtNode) End() lex.Position   { return n.EndPos }
func (n BreakStmtNode) String() string      { return n.Token.String() }

type ContinueStmtNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n ContinueStmtNode) Begin() lex.Position { return n.BeginPos }
func (n ContinueStmtNode) End() lex.Position   { return n.EndPos }
func (n ContinueStmtNode) String() string      { return n.Token.String() }

type ReturnStmtNode struct {
	Exp      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n ReturnStmtNode) Begin() lex.Position { return n.BeginPos }
func (n ReturnStmtNode) End() lex.Position   { return n.EndPos }

func (n ReturnStmtNode) String() string {
	if n.Exp != nil {
		return n.Exp.String()
	}
	return "return"
}

type DeferStmtNode struct {
	Call     CallNode
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n DeferStmtNode) Begin() lex.Position { return n.BeginPos }
func (n DeferStmtNode) End() lex.Position   { return n.EndPos }
func (n DeferStmtNode) String() string      { return fmt.Sprintf("defer %s", n.Call) }

type AssertStmtNode struct {
	Exp      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n AssertStmtNode) Begin() lex.Position { return n.BeginPos }
func (n AssertStmtNode) End() lex.Position   { return n.EndPos }
func (n AssertStmtNode) String() string      { return fmt.Sprintf("assert %s", n.Exp) }

type BlockNode struct {
	Declarations []Node
	BeginPos     lex.Position
	EndPos       lex.Position
}

func (n BlockNode) Begin() lex.Position { return n.BeginPos }
func (n BlockNode) End() lex.Position   { return n.EndPos }
func (n BlockNode) String() string      { return fmt.Sprintf("{\n\t%s\n}", n.Declarations) }

type ExpNode struct {
	Exp      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n ExpNode) Begin() lex.Position { return n.BeginPos }
func (n ExpNode) End() lex.Position   { return n.EndPos }
func (n ExpNode) String() string      { return n.Exp.String() }

type TernaryOpNode struct {
	Exp      Node
	TrueExp  Node
	FalseExp Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n TernaryOpNode) Begin() lex.Position { return n.BeginPos }
func (n TernaryOpNode) End() lex.Position   { return n.EndPos }

func (n TernaryOpNode) String() string {
	return fmt.Sprintf("%s ? %s : %s", n.Exp, n.TrueExp, n.FalseExp)
}

type BinaryOpNode struct {
	LeftExp  Node
	Op       lex.Token
	RightExp Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n BinaryOpNode) Begin() lex.Position { return n.BeginPos }
func (n BinaryOpNode) End() lex.Position   { return n.EndPos }
func (n BinaryOpNode) String() string      { return fmt.Sprintf("%s %s %s", n.LeftExp, n.Op, n.RightExp) }

type UnaryOpNode struct {
	Op       lex.Token
	Operand  Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n UnaryOpNode) Begin() lex.Position { return n.BeginPos }
func (n UnaryOpNode) End() lex.Position   { return n.EndPos }
func (n UnaryOpNode) String() string      { return fmt.Sprintf("%s%s", n.Op, n.Operand) }

type LogicalAndNode struct {
	LHS      Node
	RHS      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n LogicalAndNode) Begin() lex.Position { return n.BeginPos }
func (n LogicalAndNode) End() lex.Position   { return n.EndPos }
func (n LogicalAndNode) String() string      { return fmt.Sprintf("%s && %s", n.LHS, n.RHS) }

type LogicalOrNode struct {
	LHS      Node
	RHS      Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n LogicalOrNode) Begin() lex.Position { return n.BeginPos }
func (n LogicalOrNode) End() lex.Position   { return n.EndPos }
func (n LogicalOrNode) String() string      { return fmt.Sprintf("%s || %s", n.LHS, n.RHS) }

type BooleanNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n BooleanNode) Begin() lex.Position { return n.BeginPos }
func (n BooleanNode) End() lex.Position   { return n.EndPos }
func (n BooleanNode) String() string      { return n.Token.String() }

type NumberNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n NumberNode) Begin() lex.Position { return n.BeginPos }
func (n NumberNode) End() lex.Position   { return n.EndPos }
func (n NumberNode) String() string      { return n.Token.Literal }

type StringNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n StringNode) Begin() lex.Position { return n.BeginPos }
func (n StringNode) End() lex.Position   { return n.EndPos }
func (n StringNode) String() string      { return n.Token.Literal }

type ListNode struct {
	Elements []Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n ListNode) Begin() lex.Position { return n.BeginPos }
func (n ListNode) End() lex.Position   { return n.EndPos }

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
	BeginPos  lex.Position
	EndPos    lex.Position
}

func (n CallNode) Begin() lex.Position { return n.BeginPos }
func (n CallNode) End() lex.Position   { return n.EndPos }

func (n CallNode) String() string {
	return fmt.Sprintf("fun %s(%s)", n.Callee, n.Arguments)
}

type IndexOfNode struct {
	Sequence Node
	Index    Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n IndexOfNode) Begin() lex.Position { return n.BeginPos }
func (n IndexOfNode) End() lex.Position   { return n.EndPos }
func (n IndexOfNode) String() string      { return fmt.Sprintf("%s[%s]", n.Sequence, n.Index) }

type FunctionNode struct {
	Identifier IdentifierNode
	Parameters []IdentifierNode
	Body       BlockNode
	BeginPos   lex.Position
	EndPos     lex.Position
}

func (n FunctionNode) Begin() lex.Position { return n.BeginPos }
func (n FunctionNode) End() lex.Position   { return n.EndPos }

func (n FunctionNode) String() string {
	return fmt.Sprintf("fun %s(%s)\n%s", n.Identifier, n.Parameters, n.Body)
}

type KeyValueNode struct {
	Key      Node
	Value    Node
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n KeyValueNode) Begin() lex.Position { return n.BeginPos }
func (n KeyValueNode) End() lex.Position   { return n.EndPos }
func (n KeyValueNode) String() string      { return fmt.Sprintf("%s:%s", n.Key.String(), n.Value.String()) }

type MapNode struct {
	Elements []KeyValueNode
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n MapNode) Begin() lex.Position { return n.BeginPos }
func (n MapNode) End() lex.Position   { return n.EndPos }
func (n MapNode) String() string      { return fmt.Sprintf("{%s}", n.Elements) }

type CommentNode struct {
	Token    lex.Token
	BeginPos lex.Position
	EndPos   lex.Position
}

func (n CommentNode) Begin() lex.Position { return n.BeginPos }
func (n CommentNode) End() lex.Position   { return n.EndPos }
func (n CommentNode) String() string      { return n.Token.String() }
