package main

import (
	"bufio"
	"fmt"
	"github.com/shreerangdixit/lox/interpreter"
	"github.com/shreerangdixit/lox/lexer"
	"github.com/shreerangdixit/lox/parser"
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

func startREPL(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, "%s\n", Logo)

	scanner := bufio.NewScanner(in)
	ipt := interpreter.New()
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
		root, err := p.Parse()
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}

		// If the input is a single expression, evaluate and print the result
		// Otherwise run statements
		exp, ok := isSingleExpression(root)
		if !ok {
			_, err = ipt.Run(root)
			if err != nil {
				fmt.Fprintf(out, "%s\n", err)
			}
		} else {
			val, err := ipt.Run(exp)
			if err != nil {
				fmt.Fprintf(out, "%s\n", err)
			} else if val != interpreter.NIL {
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

func main() {
	startREPL(os.Stdin, os.Stdout)
}
