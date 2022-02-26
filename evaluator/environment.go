package evaluator

import (
	"fmt"
)

var globals map[string]Object = make(map[string]Object)

func init() {
	// Declare native functions
	for _, f := range NativeFunctions {
		globals[f.Name()] = f
	}
}

type Environment struct {
	scopeVariables map[string]Object
	enclosing      *Environment
}

func NewEnvironment() *Environment {
	env := Environment{
		scopeVariables: make(map[string]Object),
		enclosing:      nil,
	}

	return &env
}

func (e *Environment) WithEnclosing(env *Environment) *Environment {
	e.enclosing = env
	return e
}

func (e *Environment) Declare(varName string, varValue Object) error {
	if _, ok := e.scopeVariables[varName]; ok {
		return fmt.Errorf("cannot redeclare variable %s", varName)
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Environment) Assign(varName string, varValue Object) error {
	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(varName, varValue)
		}
		return fmt.Errorf("variable not declared %s", varName)
	}
	e.scopeVariables[varName] = varValue
	return nil
}

func (e *Environment) Get(varName string) (Object, error) {
	if val, ok := globals[varName]; ok {
		return val, nil
	}

	if _, ok := e.scopeVariables[varName]; !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(varName)
		}
		return NIL, fmt.Errorf("variable not declared %s", varName)
	}
	return e.scopeVariables[varName], nil
}
