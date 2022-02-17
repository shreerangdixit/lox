package runtime

import (
	"fmt"
)

var NIL = Nil{}

type ObjectType string

const (
	TypeNumber ObjectType = "number"
	TypeBool   ObjectType = "bool"
	TypeString ObjectType = "string"
	TypeFunc   ObjectType = "function"
	TypeNil    ObjectType = "null"
	TypeType   ObjectType = "type"
	TypeList   ObjectType = "list"
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
	Size() Number
	Index(Number) (Object, error)
}

type Truthifier interface {
	Object
	Truthy() Bool
}

type Adder interface {
	Object
	Add(other Object) (Object, error)
}

type Subtractor interface {
	Object
	Subtract(other Object) (Object, error)
}

type Multiplier interface {
	Object
	Multiply(other Object) (Object, error)
}

type Notter interface {
	Object
	Not() (Object, error)
}

type Negator interface {
	Object
	Negate() (Object, error)
}

type Divider interface {
	Object
	Divide(other Object) (Object, error)
}

type LessThanComparator interface {
	Object
	LessThan(other Object) Bool
}

type GreaterThanComparator interface {
	Object
	GreaterThan(other Object) Bool
}

type EqualToComparator interface {
	Object
	EqualTo(other Object) Bool
}

// ------------------------------------
// Type functions
// ------------------------------------

func IsTruthy(o Object) bool {
	if truthy, ok := o.(Truthifier); ok {
		return truthy.Truthy().Value
	}
	return false
}

func Add(left Object, right Object) (Object, error) {
	if err := checkTypeCompat(left, right); err != nil {
		return NIL, err
	}

	if adder, ok := left.(Adder); ok {
		return adder.Add(right.(Adder))
	}

	return NIL, fmt.Errorf("Cannot add types %s", left.Type())
}

func Subtract(left Object, right Object) (Object, error) {
	if err := checkTypeCompat(left, right); err != nil {
		return NIL, err
	}

	if subtractor, ok := left.(Subtractor); ok {
		return subtractor.Subtract(right.(Subtractor))
	}

	return NIL, fmt.Errorf("Cannot subtract type %s", left.Type())
}

func Divide(left Object, right Object) (Object, error) {
	if err := checkTypeCompat(left, right); err != nil {
		return NIL, err
	}

	if divider, ok := left.(Divider); ok {
		return divider.Divide(right.(Divider))
	}

	return NIL, fmt.Errorf("Cannot divide type %s", left.Type())
}

func Multiply(left Object, right Object) (Object, error) {
	if err := checkTypeCompat(left, right); err != nil {
		return NIL, err
	}

	if multiplier, ok := left.(Multiplier); ok {
		return multiplier.Multiply(right.(Multiplier))
	}

	return NIL, fmt.Errorf("Cannot multiply type %s", left.Type())
}

func Negate(o Object) (Object, error) {
	if negator, ok := o.(Negator); ok {
		return negator.Negate()
	}
	return NIL, fmt.Errorf("Cannot negate type %s", o.Type())
}

func Not(o Object) (Object, error) {
	if notter, ok := o.(Notter); ok {
		return notter.Not()
	}
	return NIL, fmt.Errorf("Cannot not type %s", o.Type())
}

func EqualTo(left Object, right Object) Bool {
	if err := checkTypeCompat(left, right); err != nil {
		return FALSE
	}

	if eqto, ok := left.(EqualToComparator); ok {
		return eqto.EqualTo(right.(EqualToComparator))
	}

	return FALSE
}

func NotEqualTo(left Object, right Object) Bool {
	return NewBool(!EqualTo(left, right).Value)
}

func LessThan(left Object, right Object) Bool {
	if err := checkTypeCompat(left, right); err != nil {
		return FALSE
	}

	if lt, ok := left.(LessThanComparator); ok {
		return lt.LessThan(right.(LessThanComparator))
	}

	return FALSE
}

func LessThanEq(left Object, right Object) Bool {
	return NewBool(LessThan(left, right).Value || EqualTo(left, right).Value)
}

func GreaterThan(left Object, right Object) Bool {
	if err := checkTypeCompat(left, right); err != nil {
		return FALSE
	}

	if gt, ok := left.(GreaterThanComparator); ok {
		return gt.GreaterThan(right.(GreaterThanComparator))
	}

	return FALSE
}

func GreaterThanEq(left Object, right Object) Bool {
	return NewBool(GreaterThan(left, right).Value || EqualTo(left, right).Value)
}

func ItemAtIndex(o Object, idx Number) (Object, error) {
	if seq, ok := o.(Sequence); ok {
		return seq.Index(idx)
	}
	return NIL, fmt.Errorf("Cannot index type %s", o.Type())

}

// ------------------------------------
// Helpers
// ------------------------------------

func checkTypeCompat(left Object, right Object) error {
	if left.Type() != right.Type() {
		return fmt.Errorf("incompatible types %s and %s", left.Type(), right.Type())
	}
	return nil
}
