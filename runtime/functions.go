package runtime

import (
	"fmt"
	"math"
	"time"
)

type FunctionHandler func(e *Evaluator, args []Object) (Object, error)

type Function struct {
	name     string
	arity    int
	variadic bool
	handler  FunctionHandler
}

func (f Function) Type() ObjectType                                 { return TypeFunc }
func (f Function) String() string                                   { return f.name }
func (f Function) Arity() int                                       { return f.arity }
func (f Function) Variadic() bool                                   { return f.variadic }
func (f Function) Call(e *Evaluator, args []Object) (Object, error) { return f.handler(e, args) }

var NativeFunctions = []Function{
	{
		name:     "sleep",
		arity:    1,
		variadic: false,
		handler:  sleepHandler,
	},
	{
		name:     "time",
		arity:    0,
		variadic: false,
		handler:  timeHandler,
	},
	{
		name:     "abs",
		arity:    1,
		variadic: false,
		handler:  absHandler,
	},
	{
		name:     "max",
		arity:    2,
		variadic: false,
		handler:  maxHandler,
	},
	{
		name:     "min",
		arity:    2,
		variadic: false,
		handler:  minHandler,
	},
	{
		name:     "type",
		arity:    1,
		variadic: false,
		handler:  typeHandler,
	},
	{
		name:     "len",
		arity:    1,
		variadic: false,
		handler:  lenHandler,
	},
	{
		name:     "print",
		arity:    -1,
		variadic: true,
		handler:  printHandler,
	},
	{
		name:     "println",
		arity:    -1,
		variadic: true,
		handler:  printlnHandler,
	},
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

func printHandler(e *Evaluator, args []Object) (Object, error) {
	for _, obj := range args {
		fmt.Printf("%v", obj)
		fmt.Printf(" ")
	}
	return NIL, nil
}

func printlnHandler(e *Evaluator, args []Object) (Object, error) {
	_, _ = printHandler(e, args)
	fmt.Println()
	return NIL, nil
}
