package eval

import (
	"os"
	"path/filepath"
	"strings"
)

type Module struct {
	Path   string
	Name   string
	Parent string
}

func NewModule(path string) *Module {
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		path, err = filepath.Abs(filepath.Join(cwd, path))
		if err != nil {
			panic(err)
		}
	}

	pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
	name := filepath.Base(pathWithoutExt)
	parent := strings.TrimSuffix(path, filepath.Base(path))
	if !strings.HasSuffix(path, ".rds") {
		path += ".rds"
	}
	return &Module{
		Path:   path,
		Name:   name,
		Parent: parent,
	}
}

func (m *Module) Data() (string, error) {
	data, err := os.ReadFile(m.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
