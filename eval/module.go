package eval

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Module string
type ModuleCommands string
type ModuleSource string

func ModuleFromFile(path string) Module {
	return Module(strings.TrimSuffix(path, filepath.Ext(path)))
}

func (m Module) Source() ModuleSource {
	curr, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return ModuleSource(filepath.Join(curr, fmt.Sprintf("%s.rds", string(m))))
}

func (m Module) Commands() (ModuleCommands, error) {
	modulePath := fmt.Sprintf("%s.rds", string(m))
	file, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	cmds, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return ModuleCommands(cmds), nil
}
