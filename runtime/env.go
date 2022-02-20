package runtime

import (
	"fmt"
)

var globals map[string]Object = make(map[string]Object)

func init() {
	// Declare native functions
	for _, f := range NativeFunctions {
		globals[f.String()] = f
	}
}

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
	scopeVariables map[string]Object
	enclosing      *Env
}

func NewEnv() *Env {
	env := Env{
		scopeVariables: make(map[string]Object),
		enclosing:      nil,
	}

	return &env
}

func (e *Env) WithEnclosing(env *Env) *Env {
	e.enclosing = env
	return e
}

func (e *Env) Declare(varName string, varValue Object) error {
	if _, ok := e.scopeVariables[varName]; ok {
		return newEnvError(fmt.Sprintf("cannot redeclare variable %s", varName))
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Env) Assign(varName string, varValue Object) error {
	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(varName, varValue)
		}
		return newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Env) Get(varName string) (Object, error) {
	if val, ok := globals[varName]; ok {
		return val, nil
	}

	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(varName)
		}
		return NIL, newEnvError(fmt.Sprintf("variable not declared %s", varName))
	}
	return e.scopeVariables[varName], nil
}
