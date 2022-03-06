package run

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shreerangdixit/redes/eval"
)

func RunFile(file string) error {
	abspath, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	basedir := strings.TrimSuffix(abspath, filepath.Base(abspath))
	if err := os.Chdir(basedir); err != nil {
		panic(err)
	}

	return eval.Import(eval.NewEvaluator(), eval.ModuleFromFile(abspath))
}
