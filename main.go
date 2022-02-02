package main

import (
	"bufio"
	"fmt"
	"io"
	"lox/interpreter"
	"lox/lexer"
	"lox/parser"
	"os"
)

func startREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(">> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		txt := scanner.Text()
		if txt == "bye" {
			break
		}

		interpreter, err := interpreter.New(parser.New(lexer.New(txt)))
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}
		fmt.Fprintf(out, "%v\n", interpreter.Run().Value)
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
