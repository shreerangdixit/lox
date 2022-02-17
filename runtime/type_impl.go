package runtime

import (
	"fmt"
)

// IEEE 754 floating point number type
// Implements the following interfaces
// Object
// Truthifier
// Negator
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
// Adder
// Subtractor
// Multiplier
// Divider
type Float64 struct{ Value float64 }

func NewFloat64(value float64) Float64          { return Float64{Value: value} }
func (f Float64) Type() ObjectType              { return TypeFloat64 }
func (f Float64) String() string                { return fmt.Sprintf("%v", f.Value) }
func (f Float64) Truthy() bool                  { return f.Value != 0 }
func (f Float64) Negate() (Object, error)       { return f.Multiply(NewFloat64(-1)) }
func (f Float64) LessThan(other Object) bool    { return f.Value < other.(Float64).Value }
func (f Float64) GreaterThan(other Object) bool { return f.Value > other.(Float64).Value }
func (f Float64) EqualTo(other Object) bool     { return f.Value == other.(Float64).Value }

func (f Float64) Add(other Object) (Object, error) {
	return NewFloat64(f.Value + other.(Float64).Value), nil
}

func (f Float64) Subtract(other Object) (Object, error) {
	return NewFloat64(f.Value - other.(Float64).Value), nil
}

func (f Float64) Multiply(other Object) (Object, error) {
	return NewFloat64(f.Value * other.(Float64).Value), nil
}

func (f Float64) Divide(other Object) (Object, error) {
	if other.(Float64).Value == 0 {
		return nil, fmt.Errorf("Divide by zero error")
	}
	return NewFloat64(f.Value / other.(Float64).Value), nil
}

// Boolean type
// Implements the following interfaces
// Object
// Truthifier
// Notter
type Bool struct{ Value bool }

var TRUE = NewBool(true)
var FALSE = NewBool(false)

func NewBool(value bool) Bool            { return Bool{Value: value} }
func (f Bool) Type() ObjectType          { return TypeBool }
func (f Bool) String() string            { return fmt.Sprintf("%v", f.Value) }
func (f Bool) EqualTo(other Object) bool { return f.Value == other.(Bool).Value }
func (f Bool) Truthy() bool              { return f.Value }
func (f Bool) Not() (Object, error)      { return NewBool(!f.Value), nil }

// String type
// Implements the following interfaces
// Object
// Sequence
// Truthifier
// Adder
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
type String struct{ Value string }

func NewString(value string) String            { return String{Value: value} }
func (f String) Type() ObjectType              { return TypeString }
func (f String) String() string                { return f.Value }
func (f String) Size() int                     { return len(f.Value) }
func (f String) Truthy() bool                  { return f.Size() > 0 }
func (f String) LessThan(other Object) bool    { return f.Value < other.(String).Value }
func (f String) GreaterThan(other Object) bool { return f.Value > other.(String).Value }
func (f String) EqualTo(other Object) bool     { return f.Value == other.(String).Value }

func (f String) Add(other Object) (Object, error) {
	return NewString(f.Value + other.(String).Value), nil
}

// Type information meta-type
// Implements the following interfaces
// Object
// Truthifier
// EqualToComparator
type Type struct{ Value ObjectType }

func NewType(value ObjectType) Type      { return Type{Value: value} }
func (f Type) Type() ObjectType          { return TypeType }
func (f Type) String() string            { return string(f.Value) }
func (f Type) Truthy() bool              { return true }
func (f Type) EqualTo(other Object) bool { return f.Value == other.(Type).Value }

// Heterogenous list type
// Implements the following interfaces
// Object
// Sequence
// Truthifier
// Adder
// TODO:
// LessThanComparator
// GreaterThanComparator
// EqualToComparator
type List struct{ Values []Object }

func NewList(values []Object) List { return List{Values: values} }
func (f List) Type() ObjectType    { return TypeList }
func (f List) String() string      { return fmt.Sprintf("%v", f.Values) }
func (f List) Size() int           { return len(f.Values) }
func (f List) Truthy() bool        { return f.Size() > 0 }

func (f List) Add(other Object) (Object, error) {
	l, ok := other.(List)
	if !ok {
		return nil, fmt.Errorf("cannot concatenate list with %s", other.Type())
	}

	return NewList(append(f.Values, l.Values...)), nil
}

// Nil type
// Implements the following interfaces
// Object
// Truthifier
// EqualToComparator
type Nil struct{}

func (f Nil) Type() ObjectType          { return TypeNil }
func (f Nil) String() string            { return "nil" }
func (f Nil) Truthy() bool              { return false }
func (f Nil) EqualTo(other Object) bool { return true }
