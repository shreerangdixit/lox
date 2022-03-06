package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shreerangdixit/redes/ast"
	"github.com/shreerangdixit/redes/eval"
	"github.com/shreerangdixit/redes/lex"
)

func RunFile(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	script, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	root, err := ast.New(lex.New(string(script))).RootNode()
	if err != nil {
		if formatter, ok := NewFormatter(err, Source(file), Commands(string(script))); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	e := eval.NewEvaluator()
	_, err = e.Evaluate(root)
	if err != nil {
		if formatter, ok := NewFormatter(err, Source(file), Commands(string(script))); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	return nil
}
