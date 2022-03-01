package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/shreerangdixit/lox/lex"
)

type PositionError interface {
	error
	ErrorType() string
	Begin() lex.Position
	End() lex.Position
}

func PrintError(err error, file string) {
	if err, ok := err.(PositionError); ok {
		print(err.Begin(), err.End(), err.ErrorType(), file, err)
	}
}

func print(begin lex.Position, end lex.Position, errtype string, file string, err error) {
	defer func() {
		os.Exit(1)
	}()

	fmt.Printf("%s:%d:%d %s error: %v\n", file, end.Line, end.Column, errtype, err)
	lines := readLines(file)
	fmt.Println(lines[end.Line-1])
	markColumns(begin.Column, end.Column)
	fmt.Println("")
}

func markColumns(beginCol, endCol int) {
	if beginCol < endCol {
		for i := 1; i < beginCol; i++ {
			fmt.Printf(" ")
		}
		for i := beginCol; i <= endCol; i++ {
			fmt.Printf("^")
		}
	} else {
		for i := 0; i <= endCol; i++ {
			fmt.Printf("^")
		}
	}
}

func readLines(file string) []string {
	fstr, e := os.ReadFile(file)
	if e != nil {
		panic(e)
	}

	str := string(fstr) + "\n" // Hack to ensure we can highlight errors on the last line
	return strings.Split(str, "\n")
}
