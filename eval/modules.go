package eval

import (
	"os"
	"path/filepath"
	"strings"
)

type FileModule struct {
	path string
	name string
}

func NewFileModule(path string) *FileModule {
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
	if !strings.HasSuffix(path, ".rds") {
		path += ".rds"
	}
	return &FileModule{
		path: path,
		name: name,
	}
}

func (m *FileModule) Name() string {
	return m.name
}

func (m *FileModule) Path() string {
	return m.path
}

func (m *FileModule) Data() (string, error) {
	data, err := os.ReadFile(m.path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type InMemoryModule struct {
	name string
	path string
	data string
}

func NewInMemoryModule(name, path, data string) *InMemoryModule {
	return &InMemoryModule{
		name: name,
		path: path,
		data: data,
	}
}

func (m *InMemoryModule) Name() string {
	return m.name
}

func (m *InMemoryModule) Path() string {
	return m.path
}

func (m *InMemoryModule) Data() (string, error) {
	return m.data, nil
}
