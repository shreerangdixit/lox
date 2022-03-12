package eval

import (
	"fmt"
	"strings"

	"github.com/shreerangdixit/redes/lex"
)

type PositionError interface {
	error
	ErrorType() string
	Begin() lex.Position
	End() lex.Position
	Inner() error
}

type Module interface {
	Data() (string, error)
	Name() string
	Path() string
}

type ErrorFormatter struct {
	mod   Module
	lines []string
	err   PositionError
}

func NewErrorFormatter(err error, mod Module) (*ErrorFormatter, bool) {
	if err, ok := err.(PositionError); ok {
		// Unwind stack trace
		for err.Inner() != nil {
			var ok bool
			if err, ok = err.Inner().(PositionError); ok {
				continue
			}
			break
		}
		if data, e := mod.Data(); e == nil {
			lines := strings.Split(data, "\n")
			lines = append(lines, "\n") // Hack to ensure we can highlight errors on the last line
			return &ErrorFormatter{
				mod:   mod,
				lines: lines,
				err:   err,
			}, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (f *ErrorFormatter) Format() string {
	str := fmt.Sprintf("\n%s:%d:%d %s error: %v\n", f.mod.Path(), f.err.End().Line, f.err.End().Column, f.err.ErrorType(), f.err)
	endLine := f.err.End().Line - 1
	if endLine < 0 {
		endLine = 0
	}
	str += fmt.Sprintf("%s\n", f.lines[endLine])
	str += fmt.Sprintf("%s\n", f.arrows())
	return str
}

func (f *ErrorFormatter) arrows() string {
	str := ""
	beginCol := f.err.Begin().Column
	endCol := f.err.End().Column
	if beginCol < endCol {
		for i := 1; i < beginCol; i++ {
			str += " "
		}
		for i := beginCol; i <= endCol; i++ {
			str += "^"
		}
	} else {
		for i := 0; i <= endCol; i++ {
			str += "^"
		}
	}
	return str
}
