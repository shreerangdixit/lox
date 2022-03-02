package eval

import (
	"fmt"
	"math"
	"time"

	"github.com/shreerangdixit/lox/ast"
)

// ------------------------------------
// User function
// ------------------------------------

type UserFunction struct {
	node    ast.FunctionNode
	closure *Environment
}

func NewUserFunction(node ast.FunctionNode, closure *Environment) *UserFunction {
	return &UserFunction{
		node:    node,
		closure: closure,
	}
}

func (f *UserFunction) Type() ObjectType { return TypeFunc }
func (f *UserFunction) Name() string     { return f.node.Identifier.Token.Literal }
func (f *UserFunction) String() string   { return "<fun-" + f.Name() + ">" }
func (f *UserFunction) Arity() int       { return len(f.node.Parameters) }
func (f *UserFunction) Variadic() bool   { return false }
func (f *UserFunction) Call(e *Evaluator, args []Object) (Object, error) {
	// New environment for function call
	env := NewEnvironment().WithEnclosing(f.closure)
	for i := range args {
		// Bind function arguments to values
		err := env.Declare(f.node.Parameters[i].Token.Literal, args[i])
		if err != nil {
			return NIL, err
		}
	}

	val, err := e.evalBlockNodeWithEnv(f.node.Body, env)
	if err != nil {
		switch err := err.(type) {
		case ReturnError:
			return err.Value, nil
		default:
			return NIL, err
		}
	}
	return val, err
}

// ------------------------------------
// Native (in-built) functions
// ------------------------------------

type NativeFunctionHandler func(e *Evaluator, args []Object) (Object, error)

type NativeFunction struct {
	name     string
	arity    int
	variadic bool
	handler  NativeFunctionHandler
}

func NewNativeFunction(name string, arity int, variadic bool, handler NativeFunctionHandler) *NativeFunction {
	return &NativeFunction{
		name:     name,
		arity:    arity,
		variadic: variadic,
		handler:  handler,
	}
}

func (f *NativeFunction) Type() ObjectType                                 { return TypeFunc }
func (f *NativeFunction) Name() string                                     { return f.name }
func (f *NativeFunction) String() string                                   { return "<native-" + f.name + ">" }
func (f *NativeFunction) Arity() int                                       { return f.arity }
func (f *NativeFunction) Variadic() bool                                   { return f.variadic }
func (f *NativeFunction) Call(e *Evaluator, args []Object) (Object, error) { return f.handler(e, args) }

var NativeFunctions = []*NativeFunction{
	// Time
	NewNativeFunction("sleep", 1, false, sleepHandler),
	NewNativeFunction("time", 0, false, timeHandler),
	// Math
	NewNativeFunction("abs", 1, false, absHandler),
	NewNativeFunction("max", 2, false, maxHandler),
	NewNativeFunction("min", 2, false, minHandler),
	NewNativeFunction("avg", 1, false, avgHandler),
	NewNativeFunction("sqrt", 1, false, sqrtHandler),
	// Collections
	NewNativeFunction("len", 1, false, lenHandler),
	NewNativeFunction("append", 2, false, appendHandler),
	// IO
	NewNativeFunction("print", 0, true, printHandler),
	NewNativeFunction("println", 0, true, printlnHandler),
	// Misc
	NewNativeFunction("type", 1, false, typeHandler),
	NewNativeFunction("zen", 0, false, zenHandler),
}

func sleepHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("sleep() expects a number")
	}

	time.Sleep(time.Duration(arg.Value) * time.Millisecond)
	return NIL, nil
}

func timeHandler(e *Evaluator, args []Object) (Object, error) {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	return NewNumber(float64(ms)), nil
}

func absHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("abs() expects a number")
	}
	return NewNumber(math.Abs(arg.Value)), nil
}

func maxHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("max() expects a number")
	}

	arg2, ok := args[1].(Number)
	if !ok {
		return NIL, fmt.Errorf("max() expects a number")
	}

	return NewNumber(math.Max(arg1.Value, arg2.Value)), nil
}

func minHandler(e *Evaluator, args []Object) (Object, error) {
	arg1, ok := args[0].(Number)
	if !ok {
		return NIL, fmt.Errorf("min() expects a number")
	}

	arg2, ok := args[1].(Number)
	if !ok {
		return NIL, fmt.Errorf("min() expects a number")
	}

	return NewNumber(math.Min(arg1.Value, arg2.Value)), nil
}

func avgHandler(e *Evaluator, args []Object) (Object, error) {
	seq, ok := args[0].(Sequence)
	if !ok {
		return NIL, fmt.Errorf("avg() expects a sequence")
	}

	sum := NewNumber(0)
	for _, arg := range seq.Elements() {
		if num, ok := arg.(Number); ok {
			s, _ := sum.Add(num)
			sum = s.(Number)
		} else {
			return NIL, fmt.Errorf("avg() expects numbers")
		}
	}

	return sum.Divide(seq.Size())
}

func sqrtHandler(e *Evaluator, args []Object) (Object, error) {
	if num, ok := args[0].(Number); ok {
		return NewNumber(math.Sqrt(num.Value)), nil
	}
	return NIL, fmt.Errorf("sqrt() expects a number")
}

func typeHandler(e *Evaluator, args []Object) (Object, error) {
	arg := args[0]
	return NewType(arg.Type()), nil
}

func lenHandler(e *Evaluator, args []Object) (Object, error) {
	arg, ok := args[0].(Sequence)
	if !ok {
		return NIL, fmt.Errorf("len() expects a sequence")
	}
	return arg.Size(), nil
}

func appendHandler(e *Evaluator, args []Object) (Object, error) {
	seq, ok := args[0].(Sequence)
	if !ok {
		return NIL, fmt.Errorf("append() expectes a sequence")
	}
	if list, ok := args[1].(List); ok {
		retval := seq
		var err error
		for _, elem := range list.Elements() {
			retval, err = retval.Append(elem)
			if err != nil {
				return retval, err
			}
		}
		return retval, nil
	} else {
		return seq.Append(args[1])
	}
}

func printHandler(e *Evaluator, args []Object) (Object, error) {
	for _, obj := range args {
		fmt.Print(obj)
	}
	return NIL, nil
}

func printlnHandler(e *Evaluator, args []Object) (Object, error) {
	_, _ = printHandler(e, args)
	fmt.Println()
	return NIL, nil
}

func zenHandler(e *Evaluator, args []Object) (Object, error) {
	fmt.Println(`
				---------------
				The Zen of Lox
				---------------
			 Donut is better than Bagel.
			 Cat is better than Dog.
			 Gin is better than Beer.
			 Tarkovsky is better than Bergman.
			 Golang is better than almost everything else.

	Interpreters are slower than the time it takes to build them.
	Although speed counts, the principles you learn building them are invaluable.
	`)
	return NIL, nil
}
