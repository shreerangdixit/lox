package interpreter

import (
	"fmt"
	"github.com/shreerangdixit/lox/types"
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
	scopesVariables map[string]types.TypeValue
	enclosing       *Env
}

func NewEnv() *Env {
	return &Env{
		scopesVariables: make(map[string]types.TypeValue),
		enclosing:       nil,
	}
}

func NewEnvWithEnclosing(env *Env) *Env {
	return &Env{
		scopesVariables: make(map[string]types.TypeValue),
		enclosing:       env,
	}
}

func (e *Env) Declare(varName string, varValue types.TypeValue) error {
	if _, ok := e.scopesVariables[varName]; ok {
		return newEnvError(fmt.Sprintf("cannot redeclare variable %s", varName))
	}
	e.scopesVariables[varName] = varValue
	return nil
}

func (e *Env) Assign(varName string, varValue types.TypeValue) error {
	if _, ok := e.scopesVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(varName, varValue)
		}
		return newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	e.scopesVariables[varName] = varValue
	return nil
}

func (e *Env) Get(varName string) (types.TypeValue, error) {
	if _, ok := e.scopesVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(varName)
		}
		return types.NO_VALUE, newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	return e.scopesVariables[varName], nil
}
