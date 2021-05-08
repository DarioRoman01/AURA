package lpp

import (
	"bufio"
	"fmt"
	"os"
)

var EOF_TOKEN = Token{Token_type: EOF, Literal: ""}

func printParseErros(errors []string) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func StartRpl() {
	reader := bufio.NewReader(os.Stdin)
	for source, _ := reader.ReadString('\n'); source != "salir()\n"; {
		lexer := NewLexer(source)
		parser := NewParser(lexer)
		program := parser.ParseProgam()
		if len(parser.Errors()) > 0 {
			printParseErros(parser.Errors())
		}

		fmt.Println(program.Str())
		source, _ = reader.ReadString('\n')
	}
}
