package lpp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var EOF_TOKEN = Token{Token_type: EOF, Literal: ""}

func printParseErros(errors []string) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func StartRpl() {
	scanner := bufio.NewScanner(os.Stdin)
	var scanned []string

	for {
		fmt.Print(">> ")
		scanner.Scan()
		source := scanner.Text()

		if source == "salir()" {
			break
		}

		scanned = append(scanned, source)
		lexer := NewLexer(strings.Join(scanned, " "))
		parser := NewParser(lexer)

		env := NewEnviroment(nil)
		program := parser.ParseProgam()

		if len(parser.Errors()) > 0 {
			printParseErros(parser.Errors())
		}

		evaluated := Evaluate(program, env)
		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}
	}
}
