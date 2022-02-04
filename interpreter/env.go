package interpreter

import (
	"fmt"
	"lox/types"
)

type EnvError struct {
	msg string
}

func (e EnvError) Error() string {
	return e.msg
}

func newEnvError(msg string) EnvError {
	return EnvError{msg: msg}
}

type Env struct {
	globals map[string]types.TypeValue
}

func NewEnv() *Env {
	return &Env{
		globals: make(map[string]types.TypeValue),
	}
}

func (e *Env) Declare(varName string, varValue types.TypeValue) error {
	if _, ok := e.globals[varName]; ok {
		return newEnvError(fmt.Sprintf("cannot redeclare variable %s", varName))
	}
	e.globals[varName] = varValue
	return nil
}

func (e *Env) Assign(varName string, varValue types.TypeValue) error {
	if _, ok := e.globals[varName]; !ok {
		return newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	e.globals[varName] = varValue
	return nil
}

func (e *Env) Get(varName string) (types.TypeValue, error) {
	if _, ok := e.globals[varName]; !ok {
		return types.NO_VALUE, newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	return e.globals[varName], nil
}
