package repl

import (
	"aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// iterate trough parser errors and print them
func printParseErros(errors []string, writer *bufio.Writer) {
	for _, err := range errors {
		writer.WriteString(fmt.Sprintf("%s\n", err))
	}

	writer.Flush()
}

// clear the console
func clearConsole(writter *bufio.Writer) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("clear")
	} else {
		cmd = exec.Command("cls")
	}

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		writter.WriteString("No se pudo limpiar la consola\n")
	}
}

// Start the repl
func StartRpl() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	var scanned []string

	writer.WriteString("✨ Bienvenido a Aura✨\n")
	writer.WriteString("escribe un comando para comenzar: \n")

	for {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Error: %s", r)
				scanned = scanned[:len(scanned)-1]
			}
		}()

		writer.WriteString(">>> ")
		writer.Flush()

		scanner.Scan()
		source := scanner.Text()

		if source == "salir()" || source == "salir" {
			break
		} else if source == "limpiar()" || source == "limpiar" {
			clearConsole(writer)
			continue
		}

		scanned = append(scanned, source)
		lexer := l.NewLexer(strings.Join(scanned, " "))
		parser := p.NewParser(lexer)

		env := obj.NewEnviroment(nil)
		program := parser.ParseProgam()

		if len(parser.Errors()) > 0 {
			printParseErros(parser.Errors(), writer)
			scanned = scanned[:len(scanned)-1]
			continue
		}

		evaluated := evaluator.Evaluate(program, env)
		if strings.Contains(scanned[len(scanned)-1], "escribir") {
			scanned = scanned[:len(scanned)-1] // avoid to call the previus print
		}

		if evaluated != nil && evaluated != obj.SingletonNUll {
			writer.WriteString(evaluated.Inspect() + "\n")
			writer.Flush()
			if _, isError := evaluated.(*obj.Error); isError {
				scanned = scanned[:len(scanned)-1] // delete error in scanned array
			}
		}
	}
}
