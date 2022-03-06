package eval

import (
	"fmt"
	"os"

	"github.com/shreerangdixit/redes/ast"
	"github.com/shreerangdixit/redes/lex"
)

func Import(e *Evaluator, m Module) error {
	cmds, err := m.Commands()
	if err != nil {
		return err
	}

	root, err := ast.New(lex.New(string(cmds))).RootNode()
	if err != nil {
		if formatter, ok := NewFormatter(err, m.Source(), cmds); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	_, err = e.Evaluate(root)
	if err != nil {
		if formatter, ok := NewFormatter(err, m.Source(), cmds); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	return nil
}
