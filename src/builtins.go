package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	switch arg := args[0].(type) {

	case *String:
		return &Number{Value: utf8.RuneCountInString(arg.Value)}

	case *List:
		return &Number{Value: len(arg.Values)}

	default:
		return &Error{Message: unsoportedArgumentType("longitud", types[args[0].Type()])}

	}
}

func Escribir(args ...Object) Object {
	var buff strings.Builder

	for _, arg := range args {
		switch node := arg.(type) {

		case *String:
			buff.WriteString(node.Inspect())

		case *Number:
			buff.WriteString(node.Inspect())

		case *List:
			buff.WriteString(node.Inspect())

		case *Bool:
			buff.WriteString(node.Inspect())

		default:
			return &Error{Message: unsoportedArgumentType("escribir", types[node.Type()])}
		}
	}

	fmt.Println(buff.String())
	return SingletonNUll
}

func Recibir(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if len(args) == 0 {
		str := input()
		return &String{Value: str}
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())
		str := input()
		return &String{Value: str}
	}

	return &Error{Message: unsoportedArgumentType("recibir", types[args[0].Type()])}
}

func castInt(args ...Object) Object {
	if len(args) > 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*String); isString {
		return toInt(arg.Value)
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

	if len(args) == 0 {
		strInt := input()
		return toInt(strInt)
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())
		strInt := input()
		return toInt(strInt)
	}

	return &Error{
		Message: unsoportedArgumentType("recbir_entero", types[args[0].Type()]),
	}

}

func AddToList(args ...Object) Object {
	if len(args) < 2 || len(args) > 2 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arr, isArray := args[0].(*List); isArray {
		arr.Values = append(arr.Values, args[1])
		return arr
	}

	return &Error{
		Message: unsoportedArgumentType("añandir", types[args[0].Type()]),
	}
}

func Tipo(args ...Object) Object {
	if len(args) > 1 || len(args) < 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return &String{Value: types[args[0].Type()]}
}

func input() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()
	return str
}

func toInt(str string) Object {
	number, err := strconv.Atoi(str)
	if err != nil {
		return &Error{Message: fmt.Sprintf("No se puede parsear como entero %s", str)}
	}

	return &Number{Value: number}
}

var BUILTINS = map[string]*Builtin{
	"longitud":       NewBuiltin(Longitud),
	"escribir":       NewBuiltin(Escribir),
	"recibir":        NewBuiltin(Recibir),
	"recibir_entero": NewBuiltin(RecibirEntero),
	"tipo":           NewBuiltin(Tipo),
	"entero":         NewBuiltin(castInt),
	"texto":          NewBuiltin(castString),
	"insertar":       NewBuiltin(AddToList),
}
