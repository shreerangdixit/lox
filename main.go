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
		node, err := p.Parse()
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
		} else {
			fmt.Fprintf(out, "%s\n", node)
		}
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
