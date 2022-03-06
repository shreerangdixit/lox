package run

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/shreerangdixit/redes/ast"
	"github.com/shreerangdixit/redes/build"
	"github.com/shreerangdixit/redes/eval"
	"github.com/shreerangdixit/redes/lex"
)

const Logo = `
_____  ______ _____  ______  _____
|  __ \|  ____|  __ \|  ____|/ ____|
| |__) | |__  | |  | | |__  | (___
|  _  /|  __| | |  | |  __|  \___ \
| | \ \| |____| |__| | |____ ____) |
|_|  \_\______|_____/|______|_____/
`

type Repl struct {
	in     io.Reader
	out    io.Writer
	errout io.Writer
}

func NewRepl() *Repl {
	return &Repl{
		in:     os.Stdin,
		out:    os.Stdout,
		errout: os.Stderr,
	}
}

func (r *Repl) Start() {
	fmt.Fprintf(r.out, "%s\n", Logo)
	fmt.Fprintf(r.out, "%s", build.Info)

	scanner := bufio.NewScanner(r.in)
	e := eval.NewEvaluator()
	for {
		fmt.Fprintf(r.out, "redes >>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		cmd := scanner.Text()

		root, err := ast.New(lex.New(cmd)).RootNode()
		if err != nil {
			r.printErr(cmd, err)
			continue
		}

		// If the input is a single expression, eval and print the result
		// Otherwise run statements
		exp, ok := isSingleExpression(root)
		if !ok {
			_, err = e.Evaluate(root)
			if err != nil {
				r.printErr(cmd, err)
				continue
			}
		} else {
			val, err := e.Evaluate(exp)
			if err != nil {
				r.printErr(cmd, err)
				continue
			} else if val != eval.NIL {
				fmt.Fprintf(r.out, "%s\n", val)
			}
		}
	}
}

func (r *Repl) printErr(cmd string, err error) {
	if formatter, ok := NewFormatter(err, Source("<repl>"), Commands(cmd)); ok {
		fmt.Fprintf(r.out, "%s", formatter.Format())
		return
	}
	fmt.Fprintf(r.out, "%s\n", err)
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
