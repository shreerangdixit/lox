package run

import (
	"github.com/shreerangdixit/redes/eval"
)

func RunFile(file string) error {
	return eval.Import(eval.NewEvaluator(), eval.ModuleFromFile(file))
}
