package main

import (
	"bufio"
	"fmt"
	"io"
	"lox/lexer"
	"lox/token"
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

		l := lexer.New(txt)
		for {
			nxt := l.NextToken()
			if nxt.Type == token.TT_EOF {
				break
			}
			fmt.Fprintf(out, "%+v\n", nxt)
		}
	}
}

func main() {
	startREPL(os.Stdin, os.Stdout)
}
