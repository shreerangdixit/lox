package run

import (
	"fmt"
	"strings"

	"github.com/shreerangdixit/redes/lex"
)

type Source string
type Commands string

type PositionError interface {
	error
	ErrorType() string
	Begin() lex.Position
	End() lex.Position
	Inner() error
}

type Formatter struct {
	source Source
	lines  []string
	err    PositionError
}

func NewFormatter(err error, source Source, contents Commands) (*Formatter, bool) {
	if err, ok := err.(PositionError); ok {
		// Unwind stack trace
		for err.Inner() != nil {
			var ok bool
			if err, ok = err.Inner().(PositionError); ok {
				continue
			}
			break
		}
		lines := strings.Split(string(contents), "\n")
		lines = append(lines, "\n") // Hack to ensure we can highlight errors on the last line
		return &Formatter{
			source: source,
			lines:  lines,
			err:    err,
		}, true
	}
	return nil, false
}

func (f *Formatter) Format() string {
	str := fmt.Sprintf("\n%s:%d:%d %s error: %v\n", f.source, f.err.End().Line, f.err.End().Column, f.err.ErrorType(), f.err)
	str += fmt.Sprintf("%s\n", f.lines[f.err.End().Line-1])
	str += fmt.Sprintf("%s\n", f.arrows())
	return str
}

func (f *Formatter) arrows() string {
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
