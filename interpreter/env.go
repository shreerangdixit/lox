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
	scopeVariables map[string]types.TypeValue
	enclosing      *Env
}

func NewEnv() *Env {
	return &Env{
		scopeVariables: make(map[string]types.TypeValue),
		enclosing:      nil,
	}
}

func NewEnvWithEnclosing(env *Env) *Env {
	return &Env{
		scopeVariables: make(map[string]types.TypeValue),
		enclosing:      env,
	}
}

func (e *Env) Declare(varName string, varValue types.TypeValue) error {
	if _, ok := e.scopeVariables[varName]; ok {
		return newEnvError(fmt.Sprintf("cannot redeclare variable %s", varName))
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Env) Assign(varName string, varValue types.TypeValue) error {
	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(varName, varValue)
		}
		return newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Env) Get(varName string) (types.TypeValue, error) {
	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(varName)
		}
		return types.NO_VALUE, newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	return e.scopeVariables[varName], nil
}
