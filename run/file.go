package run

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shreerangdixit/yeti/eval"
)

func File(file string) error {
	absPath, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	basedir := strings.TrimSuffix(absPath, filepath.Base(absPath))
	if err := os.Chdir(basedir); err != nil {
		panic(err)
	}

	e := eval.NewEvaluator()
	return e.Importer.Import(eval.NewFileModule(absPath))
}
