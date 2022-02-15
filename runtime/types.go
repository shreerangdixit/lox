package runtime

import (
	"fmt"
)

var NIL = Nil{}

type ObjectType string

const (
	FLOAT64_OBJ = "float64"
	BOOL_OBJ    = "bool"
	STRING_OBJ  = "string"
	FUNC_OBJ    = "function"
	NIL_OBJ     = "null"
)

type Object interface {
	Type() ObjectType
	String() string
}

type Callable interface {
	Object
	Arity() int
	Call(e *Evaluator, args []Object) (Object, error)
}

type Float64 struct{ Value float64 }
type Bool struct{ Value bool }
type String struct{ Value string }
type Nil struct{}

func NewFloat64(value float64) Float64 { return Float64{Value: value} }
func NewBool(value bool) Bool          { return Bool{Value: value} }
func NewString(value string) String    { return String{Value: value} }

func (f Float64) Type() ObjectType { return FLOAT64_OBJ }
func (f Float64) String() string   { return fmt.Sprintf("%v", f.Value) }
func (f Bool) Type() ObjectType    { return BOOL_OBJ }
func (f Bool) String() string      { return fmt.Sprintf("%v", f.Value) }
func (f String) Type() ObjectType  { return STRING_OBJ }
func (f String) String() string    { return f.Value }
func (f Nil) Type() ObjectType     { return NIL_OBJ }
func (f Nil) String() string       { return "nil" }

func IsTruthy(o Object) bool {
	if o == NIL {
		return false
	} else if o.Type() == FLOAT64_OBJ {
		return o.(Float64).Value != 0
	} else if o.Type() == BOOL_OBJ {
		return o.(Bool).Value
	}
	return true
}

func Add(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot add types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewFloat64(l.Value + r.Value), nil
	} else if left.Type() == STRING_OBJ && right.Type() == STRING_OBJ {
		l := left.(String)
		r := right.(String)
		return NewString(l.Value + r.Value), nil
	}

	return NIL, fmt.Errorf("Cannot add types %s and %s", left.Type(), right.Type())
}

func Subtract(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot subtract types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewFloat64(l.Value - r.Value), nil
	}

	return NIL, fmt.Errorf("Cannot subtract types %s and %s", left.Type(), right.Type())
}

func Divide(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot divide types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		if r.Value == 0 {
			return NIL, fmt.Errorf("Divide by zero error")
		}
		return NewFloat64(l.Value / r.Value), nil
	}

	return NIL, fmt.Errorf("Cannot divide types %s and %s", left.Type(), right.Type())
}

func Multiply(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot multiply types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewFloat64(l.Value * r.Value), nil
	}

	return NIL, fmt.Errorf("Cannot multiply types %s and %s", left.Type(), right.Type())
}

func Negate(o Object) (Object, error) {
	if o.Type() == FLOAT64_OBJ {
		obj := o.(Float64)
		return NewFloat64(obj.Value * -1), nil
	} else if o.Type() == BOOL_OBJ {
		obj := o.(Bool)
		return NewBool(!obj.Value), nil
	}
	return NIL, fmt.Errorf("Cannot negate type %s", o.Type())
}

func EqualTo(left Object, right Object) Bool {
	if left.Type() != right.Type() {
		return NewBool(false)
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == BOOL_OBJ && right.Type() == BOOL_OBJ {
		l := left.(Bool)
		r := right.(Bool)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == STRING_OBJ && right.Type() == STRING_OBJ {
		l := left.(String)
		r := right.(String)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == NIL_OBJ && right.Type() == NIL_OBJ {
		return NewBool(true)
	}

	return NewBool(false)
}

func NotEqualTo(left Object, right Object) Bool {
	return NewBool(!EqualTo(left, right).Value)
}

func LessThan(left Object, right Object) Bool {
	if left.Type() != right.Type() {
		return NewBool(false)
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value < r.Value)
	} else if left.Type() == STRING_OBJ && right.Type() == STRING_OBJ {
		l := left.(String)
		r := right.(String)
		return NewBool(l.Value < r.Value)
	}

	return NewBool(false)
}

func LessThanEq(left Object, right Object) Bool {
	return NewBool(LessThan(left, right).Value || EqualTo(left, right).Value)
}

func GreaterThan(left Object, right Object) Bool {
	if left.Type() != right.Type() {
		return NewBool(false)
	}

	if left.Type() == FLOAT64_OBJ && right.Type() == FLOAT64_OBJ {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value > r.Value)
	} else if left.Type() == STRING_OBJ && right.Type() == STRING_OBJ {
		l := left.(String)
		r := right.(String)
		return NewBool(l.Value > r.Value)
	}

	return NewBool(false)
}

func GreaterThanEq(left Object, right Object) Bool {
	return NewBool(GreaterThan(left, right).Value || EqualTo(left, right).Value)
}
