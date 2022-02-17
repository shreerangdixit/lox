package runtime

import (
	"fmt"
)

var NIL = Nil{}

type ObjectType string

const (
	TypeFloat64 ObjectType = "float64"
	TypeBool    ObjectType = "bool"
	TypeString  ObjectType = "string"
	TypeFunc    ObjectType = "function"
	TypeNil     ObjectType = "null"
	TypeType    ObjectType = "type"
	TypeList    ObjectType = "list"
)

// ------------------------------------
// Type interfaces
// ------------------------------------

type Object interface {
	Type() ObjectType
	String() string
}

type Callable interface {
	Object
	Arity() int
	Call(e *Evaluator, args []Object) (Object, error)
}

type Sequence interface {
	Object
	Size() int
}

// ------------------------------------
// Type declarations
// ------------------------------------

type Float64 struct{ Value float64 }

func NewFloat64(value float64) Float64 { return Float64{Value: value} }
func (f Float64) Type() ObjectType     { return TypeFloat64 }
func (f Float64) String() string       { return fmt.Sprintf("%v", f.Value) }

type Bool struct{ Value bool }

func NewBool(value bool) Bool   { return Bool{Value: value} }
func (f Bool) Type() ObjectType { return TypeBool }
func (f Bool) String() string   { return fmt.Sprintf("%v", f.Value) }

type String struct{ Value string }

func NewString(value string) String { return String{Value: value} }
func (f String) Type() ObjectType   { return TypeString }
func (f String) String() string     { return f.Value }

type Type struct{ Value ObjectType }

func NewType(value ObjectType) Type { return Type{Value: value} }
func (f Type) Type() ObjectType     { return TypeType }
func (f Type) String() string       { return string(f.Value) }

type List struct{ Values []Object }

func NewList(values []Object) List { return List{Values: values} }
func (f List) Type() ObjectType    { return TypeList }
func (f List) String() string      { return fmt.Sprintf("%v", f.Values) }
func (f List) Size() int           { return len(f.Values) }

type Nil struct{}

func (f Nil) Type() ObjectType { return TypeNil }
func (f Nil) String() string   { return "nil" }

// ------------------------------------
// Type functions
// ------------------------------------

func IsTruthy(o Object) bool {
	if o == NIL {
		return false
	} else if o.Type() == TypeFloat64 {
		return o.(Float64).Value != 0
	} else if o.Type() == TypeBool {
		return o.(Bool).Value
	} else if s, ok := o.(Sequence); ok {
		return s.Size() > 0
	}
	return true
}

func Add(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot add types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == TypeFloat64 {
		l := left.(Float64)
		r := right.(Float64)
		return NewFloat64(l.Value + r.Value), nil
	} else if left.Type() == TypeString {
		l := left.(String)
		r := right.(String)
		return NewString(l.Value + r.Value), nil
	} else if left.Type() == TypeList {
		l := left.(List)
		r := right.(List)
		return NewList(append(l.Values, r.Values...)), nil
	}

	return NIL, fmt.Errorf("Cannot add types %s and %s", left.Type(), right.Type())
}

func Subtract(left Object, right Object) (Object, error) {
	if left.Type() != right.Type() {
		return NIL, fmt.Errorf("Cannot subtract types %s and %s", left.Type(), right.Type())
	}

	if left.Type() == TypeFloat64 {
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

	if left.Type() == TypeFloat64 {
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

	if left.Type() == TypeFloat64 {
		l := left.(Float64)
		r := right.(Float64)
		return NewFloat64(l.Value * r.Value), nil
	}

	return NIL, fmt.Errorf("Cannot multiply types %s and %s", left.Type(), right.Type())
}

func Negate(o Object) (Object, error) {
	if o.Type() == TypeFloat64 {
		obj := o.(Float64)
		return NewFloat64(obj.Value * -1), nil
	} else if o.Type() == TypeBool {
		obj := o.(Bool)
		return NewBool(!obj.Value), nil
	}
	return NIL, fmt.Errorf("Cannot negate type %s", o.Type())
}

func EqualTo(left Object, right Object) Bool {
	if left.Type() != right.Type() {
		return NewBool(false)
	}

	if left.Type() == TypeFloat64 {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == TypeBool {
		l := left.(Bool)
		r := right.(Bool)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == TypeString {
		l := left.(String)
		r := right.(String)
		return NewBool(l.Value == r.Value)
	} else if left.Type() == TypeNil {
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

	if left.Type() == TypeFloat64 {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value < r.Value)
	} else if left.Type() == TypeString {
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

	if left.Type() == TypeFloat64 {
		l := left.(Float64)
		r := right.(Float64)
		return NewBool(l.Value > r.Value)
	} else if left.Type() == TypeString {
		l := left.(String)
		r := right.(String)
		return NewBool(l.Value > r.Value)
	}

	return NewBool(false)
}

func GreaterThanEq(left Object, right Object) Bool {
	return NewBool(GreaterThan(left, right).Value || EqualTo(left, right).Value)
}
