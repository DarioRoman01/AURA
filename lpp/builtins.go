package lpp

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"
)

func wrongNumberofArgs(found, actual int) string {
	return fmt.Sprintf("numero incorrecto de argumentos para longitud, se recibieron %d, se requieren %d", found, actual)
}

func unsoportedArgumentType(funcname, objType string) string {
	return fmt.Sprintf("argumento para %s no valido, se recibio %s", funcname, objType)
}

func Longitud(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		return &Number{Value: utf8.RuneCountInString(arg.Value)}
	}

	return &Error{Message: unsoportedArgumentType("longitud", types[args[0].Type()])}
}

func Escribir(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	switch node := args[0].(type) {

	case *String:
		fmt.Println(node.Inspect())

	case *Number:
		fmt.Println(node.Inspect())

	case *Bool:
		fmt.Println(node.Inspect())

	default:
		return &Error{Message: unsoportedArgumentType("escribir", types[node.Type()])}
	}

	return nil
}

func Recibir(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		str := scanner.Text()

		return &String{Value: str}
	}

	return &Error{Message: unsoportedArgumentType("recibir", types[args[0].Type()])}
}

func castInt(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		number, err := strconv.Atoi(arg.Value)
		if err != nil {
			return &Error{Message: fmt.Sprintf("No se puede parsear como entero %s", arg.Value)}
		}

		return &Number{Value: number}
	}

	return &Error{Message: unsoportedArgumentType("recibir", types[args[0].Type()])}
}

func castString(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isNumber := args[0].(*Number); isNumber {
		strInt := strconv.Itoa(arg.Value)
		return &String{Value: strInt}
	}

	return &Error{Message: unsoportedArgumentType("recibir", types[args[0].Type()])}
}

func RecibirEntero(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		strInt := scanner.Text()

		number, err := strconv.Atoi(strInt)
		if err != nil {
			return &Error{Message: "No se puede parsear como entero"}
		}

		return &Number{Value: number}
	}

	return &Error{
		Message: unsoportedArgumentType("recbir_entero", types[args[0].Type()]),
	}

}

func Tipo(args ...Object) Object {
	if len(args) > 1 || len(args) < 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return &String{Value: types[args[0].Type()]}
}

var BUILTINS = map[string]*Builtin{
	"longitud":       NewBuiltin(Longitud),
	"escribir":       NewBuiltin(Escribir),
	"recibir":        NewBuiltin(Recibir),
	"recibir_entero": NewBuiltin(RecibirEntero),
	"tipo":           NewBuiltin(Tipo),
	"entero":         NewBuiltin(castInt),
	"texto":          NewBuiltin(castString),
}
