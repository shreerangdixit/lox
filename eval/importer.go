package eval

import (
	"fmt"
	"os"

	"github.com/shreerangdixit/redes/ast"
	"github.com/shreerangdixit/redes/lex"
)

type Importer struct {
	importSet map[FileModule]struct{}
	latest    *FileModule
	eval      *Evaluator
}

func NewImporter(eval *Evaluator) *Importer {
	return &Importer{
		importSet: map[FileModule]struct{}{},
		latest:    nil,
		eval:      eval,
	}
}

func (i *Importer) Import(m *FileModule) error {
	// Ignore dup imports
	if _, ok := i.importSet[*m]; ok {
		return nil
	}

	cmds, err := m.Data()
	if err != nil {
		return err
	}

	root, err := ast.New(lex.New(string(cmds))).RootNode()
	if err != nil {
		if formatter, ok := NewFormatter(err, m); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	i.latest = m
	i.importSet[*m] = struct{}{}
	_, err = i.eval.Evaluate(root)
	if err != nil {
		if formatter, ok := NewFormatter(err, i.latest); ok {
			fmt.Fprintf(os.Stderr, "%s", formatter.Format())
			os.Exit(1)
		}
		return err
	}

	return nil
}
