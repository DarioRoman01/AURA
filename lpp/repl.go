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
	reader := bufio.NewReader(os.Stdin)
	var scanned []string

	for source, _ := reader.ReadString('\n'); source != "salir()\n"; {
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

		source, _ = reader.ReadString('\n')
	}
}
