package evaluator

import (
	"aura/src/lexer"
	obj "aura/src/object"
	"aura/src/parser"
	"fmt"
	"math"
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
	// check that path exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, newError(fmt.Sprintf("La ruta %s no existe", path))
	}

	// check that the path is not a dir
	if fileInfo.IsDir() {
		return nil, newError("No se indico un archivo!")
	}

	// check that the file has the .aura extension
	if filepath.Ext(path) != ".aura" {
		return nil, newError(fmt.Sprintf(
			"El archivo %s no es un archivo aura",
			filepath.Base(path),
		))
	}

	// read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, newError(fmt.Sprintf("No se leer el archivo %s", filepath.Base(path)))
	}

	// parse and evaluate the file
	lexer := lexer.NewLexer(string(content))
	parser := parser.NewParser(lexer)
	env := obj.NewEnviroment(nil)
	program := parser.ParseProgam()

	// the file has syntax erros
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

// Check that given index is valid for a list call
func checkIndex(length int, index int) (int, *obj.Error) {
	if index >= length {
		return 0, indexOutOfRange(index, length)
	}

	if index < 0 {
		if math.Abs(float64(index)) > float64(length) {
			return 0, indexOutOfRange(index, length)
		}

		return length + index, nil
	}

	return index, nil
}
