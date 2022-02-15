package runner

import (
	"bufio"
	"fmt"
	"github.com/shreerangdixit/lox/lexer"
	"github.com/shreerangdixit/lox/parser"
	"github.com/shreerangdixit/lox/runtime"
	"io"
	"os"
)

const Logo = `
.____    ________  ____  ___
|    |   \_____  \ \   \/  /
|    |    /   |   \ \     / 
|    |___/    |    \/     \ 
|_______ \_______  /___/\  \
        \/       \/      \_/
`

func RunFile(file string) error {
	script, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	p := parser.New(lexer.New(string(script)))

	ast, err := p.Parse()
	if err != nil {
		return err
	}

	e := runtime.NewEvaluator()
	_, err = e.Evaluate(ast)
	if err != nil {
		return err
	}

	return nil
}

func StartREPL() {
	startREPL(os.Stdin, os.Stdout)
}

func startREPL(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, "%s\n", Logo)

	scanner := bufio.NewScanner(in)
	e := runtime.NewEvaluator()
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

		p := parser.New(lexer.New(txt))
		ast, err := p.Parse()
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}

		// If the input is a single expression, evaluate and print the result
		// Otherwise run statements
		exp, ok := isSingleExpression(ast)
		if !ok {
			_, err = e.Evaluate(ast)
			if err != nil {
				fmt.Fprintf(out, "%s\n", err)
			}
		} else {
			val, err := e.Evaluate(exp)
			if err != nil {
				fmt.Fprintf(out, "%s\n", err)
			} else if val != runtime.NIL {
				fmt.Fprintf(out, "%s\n", val)
			}
		}
	}
}

func isSingleExpression(node parser.Node) (parser.Node, bool) {
	programNode, ok := node.(parser.ProgramNode)
	if !ok {
		return nil, false
	}

	if len(programNode.Declarations) == 0 || len(programNode.Declarations) > 1 {
		return nil, false
	}

	expStat, ok := programNode.Declarations[0].(parser.ExpStmtNode)
	if !ok {
		return nil, false
	}

	return expStat.Exp, true
}
