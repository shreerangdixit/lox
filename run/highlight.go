package run

import (
	"fmt"
	"strings"

	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/evaluate"
	"github.com/shreerangdixit/lox/lex"
)

func HighlightError(err error, script string) {
	switch err := err.(type) {
	case ast.SyntaxError:
		highlight(err.Token.BeginPosition, err.Token.EndPosition, script)
	case evaluate.EvalError:
		highlight(err.Node.Begin(), err.Node.End(), script)
	}
}

func highlight(begin lex.Position, end lex.Position, script string) {
	lineNumber := 0
	fmt.Printf("Error on Line %d, Col %d\n", end.Line, end.Column)
	for _, line := range strings.Split(script, "\n") {
		lineNumber += 1
		if lineNumber == end.Line {
			fmt.Println(line)
			if begin.Column < end.Column {
				for i := 1; i < begin.Column; i++ {
					fmt.Printf(" ")
				}
				for i := begin.Column; i <= end.Column; i++ {
					fmt.Printf("^")
				}
			} else {
				for i := 0; i <= end.Column; i++ {
					fmt.Printf("^")
				}
			}
			fmt.Println("")
			return
		}
	}
}
