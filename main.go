package main

import (
	"bufio"
	"fmt"
	"github.com/shreerangdixit/lox/interpreter"
	"github.com/shreerangdixit/lox/lexer"
	"github.com/shreerangdixit/lox/parser"
	"io"
	"os"
)

func startREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	ipt := interpreter.New()
	for {
		fmt.Printf("lox >>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		txt := scanner.Text()
		if txt == "bye" || txt == "quit" {
			break
		}

		p := parser.New(lexer.New(txt))
		root, err := p.Parse()
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}

		_, err = ipt.Run(root)
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
		}
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
