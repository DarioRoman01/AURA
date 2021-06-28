package main

import (
	e "aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"aura/src/repl"
	"fmt"
	"os"
	"path/filepath"
)

// validate that the given file in the path have the .aura extension
func validateFileExtension(path string) error {
	extension := filepath.Ext(path)
	if extension != ".aura" {
		return fmt.Errorf(
			"el archivo %s no es una archivo aura valido",
			filepath.Base(path),
		)
	}

	return nil
}

// read the file in th
func ReadFile(path string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s", r)
		}
	}()

	source, err := os.ReadFile(path)
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
		repl.StartRpl()
		return
	}

	filePath := os.Args[1]
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("La ruta %s no existe", filePath)
		return
	}

	if fileInfo.IsDir() {
		fmt.Printf("La ruta indicada es una carpeta: %s", filePath)
		return
	}

	err = validateFileExtension(filePath)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	ReadFile(filePath)
}
