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

const Logo = `
.____    ________  ____  ___
|    |   \_____  \ \   \/  /
|    |    /   |   \ \     / 
|    |___/    |    \/     \ 
|_______ \_______  /___/\  \
        \/       \/      \_/
`

func startREPL(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, "%s\n", Logo)

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
