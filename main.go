package main

import (
	"bufio"
	"fmt"
	"io"
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

		p := parser.New(lexer.New(txt))
		fmt.Fprintf(out, "%s\n", p.Expression())
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
