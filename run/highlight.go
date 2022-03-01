package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/evaluate"
	"github.com/shreerangdixit/lox/lex"
)

func HighlightError(err error, file string) {
	switch err := err.(type) {
	case ast.SyntaxError:
		highlight(err.Token.BeginPosition, err.Token.EndPosition, file, err)
	case evaluate.EvalError:
		highlight(err.Node.Begin(), err.Node.End(), file, err)
	}
}

func highlight(begin lex.Position, end lex.Position, file string, err error) {
	defer func() {
		os.Exit(1)
	}()

	fstr, e := os.ReadFile(file)
	if e != nil {
		panic(e)
	}

	script := string(fstr) + "\n"

	lineNumber := 0
	fmt.Printf("%s:%d:%d %v\n", file, end.Line, end.Column, err)
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
