package run

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shreerangdixit/lox/ast"
	"github.com/shreerangdixit/lox/build"
	"github.com/shreerangdixit/lox/eval"
	"github.com/shreerangdixit/lox/lex"
)

const Logo = `
.____    ________  ____  ___
|    |   \_____  \ \   \/  /
|    |    /   |   \ \     / 
|    |___/    |    \/     \ 
|_______ \_______  /___/\  \
        \/       \/      \_/
`

func StartREPL(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, "%s\n", Logo)
	fmt.Fprintf(out, "%s", build.Info)

	scanner := bufio.NewScanner(in)
	e := eval.NewEvaluator()
	for {
		fmt.Printf("lox >>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		txt := scanner.Text()
		if txt == "bye" || txt == "quit" {
			break
		}

		root, err := ast.New(lex.New(txt)).RootNode()
		if err != nil {
			continue
		}

		// If the input is a single expression, eval and print the result
		// Otherwise run statements
		exp, ok := isSingleExpression(root)
		if !ok {
			_, err = e.Evaluate(root)
			if err != nil {
				if formatter, ok := NewFormatter(err, ScriptSource("<repl>"), ScriptContents(string(txt))); ok {
					fmt.Fprintf(out, "%s", formatter.Format())
					continue
				}
				fmt.Fprintf(out, "%s\n", err)
			}
		} else {
			val, err := e.Evaluate(exp)
			if err != nil {
				if formatter, ok := NewFormatter(err, ScriptSource("<repl>"), ScriptContents(string(txt))); ok {
					fmt.Fprintf(out, "%s", formatter.Format())
					continue
				}
				fmt.Fprintf(out, "%s\n", err)
			} else if val != eval.NIL {
				fmt.Fprintf(out, "%s\n", val)
			}
		}
	}
}

func isSingleExpression(node ast.Node) (ast.Node, bool) {
	programNode, ok := node.(ast.ProgramNode)
	if !ok {
		return nil, false
	}

	if len(programNode.Declarations) == 0 || len(programNode.Declarations) > 1 {
		return nil, false
	}

	expStat, ok := programNode.Declarations[0].(ast.ExpStmtNode)
	if !ok {
		return nil, false
	}

	return expStat.Exp, true
}
