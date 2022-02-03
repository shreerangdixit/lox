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
		fmt.Printf("lox >>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		txt := scanner.Text()
		if txt == "bye" || txt == "quit" {
			break
		}

		ip, err := interpreter.New(parser.New(lexer.New(txt)))
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}

		_, err = ip.Run()
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
