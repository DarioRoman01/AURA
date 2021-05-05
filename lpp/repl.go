package lpp

import (
	"bufio"
	"fmt"
	"os"
)

var EOF_TOKEN = Token{Token_type: EOF, Literal: ""}

func StartRpl() {
	reader := bufio.NewReader(os.Stdin)
	for source, _ := reader.ReadString('\n'); source != "salir()\n"; {
		lexer := NewLexer(source)
		for token := lexer.NextToken(); token != EOF_TOKEN; {
			fmt.Println(token.PrintToken())
			token = lexer.NextToken()
		}

		source, _ = reader.ReadString('\n')
	}
}
