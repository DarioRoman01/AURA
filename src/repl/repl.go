package repl

import (
	"aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var EOF_TOKEN = l.Token{Token_type: l.EOF, Literal: ""}

func printParseErros(errors []string) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func StartRpl() {
	scanner := bufio.NewScanner(os.Stdin)
	var scanned []string

	fmt.Println("✨ Bienvenido a Aura✨")
	fmt.Println("escribe un comando para comenzar: ")

	for {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Error: %s", r)
				scanned = scanned[:len(scanned)-1]
			}
		}()

		fmt.Print(">>> ")
		scanner.Scan()
		source := scanner.Text()

		if source == "salir()" || source == "salir" {
			break
		} else if source == "limpiar()" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
			continue
		}

		scanned = append(scanned, source)
		lexer := l.NewLexer(strings.Join(scanned, " "))
		parser := p.NewParser(lexer)

		env := obj.NewEnviroment(nil)
		program := parser.ParseProgam()

		if len(parser.Errors()) > 0 {
			printParseErros(parser.Errors())
			scanned = scanned[:len(scanned)-1]
			continue
		}

		evaluated := evaluator.Evaluate(program, env)
		if strings.Contains(scanned[len(scanned)-1], "escribir") {
			scanned = scanned[:len(scanned)-1] // avoid to print the previus print
		}

		if evaluated != nil && evaluated != obj.SingletonNUll {
			fmt.Println(evaluated.Inspect())

			if _, isError := evaluated.(*obj.Error); isError {
				scanned = scanned[:len(scanned)-1] // delete error in scanned array
			}
		}
	}
}
