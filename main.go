package main

import (
	e "aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"aura/src/repl"
	"fmt"
	"io/ioutil"
	"os"
)

type CommandLine struct{}

func Repl() {
	repl.StartRpl()
}

func ReadFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Archivo no encontrado!")
		return
	}

	if len(source) == 0 {
		return
	}

	lexer := l.NewLexer(string(source))
	parser := p.NewParser(lexer)
	env := obj.NewEnviroment(nil)
	program := parser.ParseProgam()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println(err)
		}

		return
	}

	evaluated := e.Evaluate(program, env)
	if evaluated != nil && evaluated != obj.SingletonNUll {
		fmt.Println(evaluated.Inspect())
	}
}

func main() {
	if len(os.Args) < 2 {
		Repl()
	}

	filePath := os.Args[1]
	ReadFile(filePath)
}
