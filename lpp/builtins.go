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

func unsoportedArgumentType(objType string) string {
	return fmt.Sprintf("argumento para longitud no valido, se recibio %s", objType)
}

func noRequiredArgs(funcName string, actual int) string {
	return fmt.Sprintf("La funcion %s no recibe arguments pero se recibieron %d", funcName, actual)
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

func Escribir(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	switch node := args[0].(type) {

	case *String:
		str := fmt.Sprintf("%s\n", node.Inspect())
		fmt.Println(str)

	case *Number:
		fmt.Println(node.Inspect())

	case *Bool:
		fmt.Println(node.Inspect())

	default:
		return &Error{Message: unsoportedArgumentType(types[node.Type()])}
	}

	return nil
}

func Recibir(args ...Object) Object {
	if len(args) > 0 {
		return &Error{Message: noRequiredArgs("recibir", len(args))}
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()

	return &String{Value: str}
}

func RecibirEntero(args ...Object) Object {
	if len(args) > 0 {
		return &Error{Message: noRequiredArgs("recibir", len(args))}
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	strInt := scanner.Text()

	number, err := strconv.Atoi(strInt)
	if err != nil {
		return &Error{Message: "No se puede parsear como entero"}
	}

	return &Number{Value: number}
}

var BUILTINS = map[string]*Builtin{
	"longitud":       NewBuiltin(Longitud),
	"escribir":       NewBuiltin(Escribir),
	"recibir":        NewBuiltin(Recibir),
	"recibir_entero": NewBuiltin(RecibirEntero),
}
