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

type ExpressionType int

const (
	TYPE_NUMBER ExpressionType = iota
	TYPE_BOOLEAN
)

func (e ExpressionType) String() string {
	switch e {
	case TYPE_NUMBER:
		return "NUMBER"
	case TYPE_BOOLEAN:
		return "BOOLEAN"
	default:
		return "<UNKNOWN>"
	}
}

type ExpressionValue struct {
	Type         ExpressionType
	Float64Value float64
	BooleanValue bool
}

func (v ExpressionValue) Add(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != TYPE_NUMBER && right.Type != TYPE_NUMBER {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot add types %s and %s", v.Type, right.Type))
	}

	return ExpressionValue{
		Type:         TYPE_NUMBER,
		Float64Value: v.Float64Value + right.Float64Value,
	}, nil
}

func (v ExpressionValue) Subtract(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != TYPE_NUMBER && right.Type != TYPE_NUMBER {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot subtract types %s and %s", v.Type, right.Type))
	}

	return ExpressionValue{
		Type:         TYPE_NUMBER,
		Float64Value: v.Float64Value - right.Float64Value,
	}, nil
}

func (v ExpressionValue) Divide(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != TYPE_NUMBER && right.Type != TYPE_NUMBER {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot divide types %s and %s", v.Type, right.Type))
	}

	return ExpressionValue{
		Type:         TYPE_NUMBER,
		Float64Value: v.Float64Value / right.Float64Value,
	}, nil
}

func (v ExpressionValue) Multiply(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != TYPE_NUMBER && right.Type != TYPE_NUMBER {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot multiply types %s and %s", v.Type, right.Type))
	}

	return ExpressionValue{
		Type:         TYPE_NUMBER,
		Float64Value: v.Float64Value * right.Float64Value,
	}, nil
}

func (v ExpressionValue) Negate() (ExpressionValue, error) {
	if v.Type == TYPE_NUMBER {
		return ExpressionValue{
			Type:         TYPE_NUMBER,
			Float64Value: v.Float64Value * -1,
		}, nil
	} else if v.Type == TYPE_BOOLEAN {
		return ExpressionValue{
			Type:         TYPE_BOOLEAN,
			BooleanValue: !v.BooleanValue,
		}, nil
	}
	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot negate types %s", v.Type))
}

func (v ExpressionValue) Equals(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value == right.Float64Value}, nil
	} else if v.Type == TYPE_BOOLEAN {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.BooleanValue == right.BooleanValue}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v ExpressionValue) NotEquals(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value != right.Float64Value}, nil
	} else if v.Type == TYPE_BOOLEAN {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.BooleanValue != right.BooleanValue}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v ExpressionValue) LessThan(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value < right.Float64Value}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v ExpressionValue) LessThanEq(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value <= right.Float64Value}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v ExpressionValue) GreaterThan(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value > right.Float64Value}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}

func (v ExpressionValue) GreaterThanEq(right ExpressionValue) (ExpressionValue, error) {
	if v.Type != right.Type {
		return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
	}

	if v.Type == TYPE_NUMBER {
		return ExpressionValue{Type: TYPE_BOOLEAN, BooleanValue: v.Float64Value >= right.Float64Value}, nil
	}

	return ExpressionValue{}, newTypeError(fmt.Sprintf("Cannot compare types %s and %s", v.Type, right.Type))
}
