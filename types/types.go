package types

import (
	"fmt"
)

type TypeError struct {
	msg string
}

func (e TypeError) Error() string {
	return e.msg
}

func newTypeError(msg string) TypeError {
	return TypeError{msg: msg}
}

type Type int

const (
	NUMBER Type = iota
	BOOL
	NIL
)

func (e Type) String() string {
	switch e {
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case NIL:
		return "NIL"
	default:
		return "<UNKNOWN>"
	}
}

var NO_VALUE = TypeValue{Type: NIL}

type TypeValue struct {
	Type  Type
	Value interface{}
}

func (v TypeValue) Add(right TypeValue) (TypeValue, error) {
	if v.Type != NUMBER || right.Type != NUMBER {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot add types %s and %s", v.Type, right.Type))
	}

	return TypeValue{Type: NUMBER, Value: v.Value.(float64) + right.Value.(float64)}, nil
}

func (v TypeValue) Subtract(right TypeValue) (TypeValue, error) {
	if v.Type != NUMBER || right.Type != NUMBER {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot subtract types %s and %s", v.Type, right.Type))
	}

	return TypeValue{Type: NUMBER, Value: v.Value.(float64) - right.Value.(float64)}, nil
}

func (v TypeValue) Divide(right TypeValue) (TypeValue, error) {
	if v.Type != NUMBER || right.Type != NUMBER {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot divide types %s and %s", v.Type, right.Type))
	}

	return TypeValue{Type: NUMBER, Value: v.Value.(float64) / right.Value.(float64)}, nil
}

func (v TypeValue) Multiply(right TypeValue) (TypeValue, error) {
	if v.Type != NUMBER || right.Type != NUMBER {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot multiply types %s and %s", v.Type, right.Type))
	}

	return TypeValue{Type: NUMBER, Value: v.Value.(float64) * right.Value.(float64)}, nil
}

func (v TypeValue) Negate() (TypeValue, error) {
	if v.Type == NUMBER {
		return TypeValue{Type: NUMBER, Value: v.Value.(float64) * -1}, nil
	} else if v.Type == BOOL {
		return TypeValue{Type: BOOL, Value: !v.Value.(bool)}, nil
	}
	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot negate type %s", v.Type))
}

func (v TypeValue) Equals(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) == right.Value.(float64)}, nil
	} else if v.Type == BOOL {
		return TypeValue{Type: BOOL, Value: v.Value.(bool) == right.Value.(bool)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v TypeValue) NotEquals(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) != right.Value.(float64)}, nil
	} else if v.Type == BOOL {
		return TypeValue{Type: BOOL, Value: v.Value.(bool) != right.Value.(bool)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v TypeValue) LessThan(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) < right.Value.(float64)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v TypeValue) LessThanEq(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) <= right.Value.(float64)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v TypeValue) GreaterThan(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) > right.Value.(float64)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v TypeValue) GreaterThanEq(right TypeValue) (TypeValue, error) {
	if v.Type != right.Type {
		return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == NUMBER {
		return TypeValue{Type: BOOL, Value: v.Value.(float64) >= right.Value.(float64)}, nil
	}

	return TypeValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}
