package lpp

import (
	"fmt"
	"unicode/utf8"
)

func wrongNumberofArgs(found, actual int) string {
	return fmt.Sprintf("numero incorrecto de argumentos para longitud, se recibieron %d, se requieren %d", found, actual)
}

func unsoportedArgumentType(objType string) string {
	return fmt.Sprintf("argumento para longitud no valido, se recibio %s", objType)
}

func Longitud(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		return &Number{Value: utf8.RuneCountInString(arg.Value)}
	}

	return &Error{Message: unsoportedArgumentType(types[args[0].Type()])}
}

var BUILTINS = map[string]*Builtin{
	"longitud": NewBuiltin(Longitud),
}
