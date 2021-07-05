package evaluator

import (
	"aura/src/lexer"
	obj "aura/src/object"
	"aura/src/parser"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"
)

func makeStringList(str string) []obj.Object {
	list := make([]obj.Object, 0, utf8.RuneCountInString(str))
	for _, char := range str {
		list = append(list, &obj.String{Value: string(char)})
	}

	return list
}

// import the enviroment of other file parsing and evaluating the other file
func importEnv(path string) (*obj.Enviroment, *obj.Error) {
	// check that path exists and its a file
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, newError(fmt.Sprintf("La ruta %s no existe", path))
	}

	if fileInfo.IsDir() {
		return nil, newError("No se indico un archivo!")
	}

	if filepath.Ext(path) != ".aura" {
		return nil, newError(fmt.Sprintf(
			"El archivo %s no es un archivo aura",
			filepath.Base(path),
		))
	}

	content, _ := os.ReadFile(path)
	lexer := lexer.NewLexer(string(content))
	parser := parser.NewParser(lexer)
	env := obj.NewEnviroment(nil)
	program := parser.ParseProgam()
	if len(parser.Errors()) != 0 {
		return nil, newError(fmt.Sprintf(
			"el archivo %s contiene errores de syntaxis",
			filepath.Base(path),
		))
	}

	evaluated := Evaluate(program, env)
	if evaluated != nil {
		return env, nil
	}

	return nil, newError("La evaluacion fue nula")
}
