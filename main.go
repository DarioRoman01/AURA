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

// validate that given path exists, have a file and the extension
// is .aura
func validatePath(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("la ruta %s no existe", path)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("la ruta indicada no contiene un archivo: %s", path)
	}

	extension := filepath.Ext(path)
	if extension != ".aura" {
		return fmt.Errorf(
			"el archivo %s no es una archivo aura valido",
			filepath.Base(path),
		)
	}

	return nil
}

// read the file in the path and evaluate the file
func ReadFile(path string) {
	defer func() {
		// we handle a posible panic in the parser
		// and the evaluator
		if r := recover(); r != nil {
			fmt.Printf("Error: %s", r)
			return
		}
	}()

	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("No se pudo leer el archivo")
		return
	}

	if len(source) == 0 {
		fmt.Println("El archivo esta vacio!")
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
		// we dont evaluate the program if has syntax errors
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
	if err := validatePath(filePath); err != nil {
		fmt.Println(err.Error())
		return
	}

	ReadFile(filePath)
}
