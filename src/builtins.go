package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type MethodsTypes int

const (
	Lhead MethodsTypes = iota
	POP
	APPEND
	REMOVE
	CONTAIS
)

var scanner = bufio.NewScanner(os.Stdin)

func wrongNumberofArgs(found, actual int) string {
	return fmt.Sprintf("numero incorrecto de argumentos para largo, se recibieron %d, se requieren %d", found, actual)
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

	case *Map:
		return &Number{Value: len(arg.store)}

	default:
		return &Error{Message: unsoportedArgumentType("largo", types[args[0].Type()])}
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

		case *Map:
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
		str := input(scanner)
		return &String{Value: str}
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())
		str := input(scanner)
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
		strInt := input(scanner)
		return toInt(strInt)
	}

	if arg, isString := args[0].(*String); isString {
		fmt.Print(arg.Inspect())
		strInt := input(scanner)
		return toInt(strInt)
	}

	return &Error{
		Message: unsoportedArgumentType("recbir_entero", types[args[0].Type()]),
	}

}

func add(args ...Object) Object {
	if len(args) > 1 || len(args) == 0 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if num, isNumber := args[0].(*Number); isNumber {
		return NewMethod(num, APPEND)
	}

	return &Error{Message: unsoportedArgumentType("add", types[args[0].Type()])}
}

func remove(args ...Object) Object {
	if len(args) > 1 || len(args) == 0 {
		return &Error{wrongNumberofArgs(len(args), 1)}
	}

	if num, isNumber := args[0].(*Number); isNumber {
		return NewMethod(num, REMOVE)
	}

	return &Error{Message: unsoportedArgumentType("add", types[args[0].Type()])}
}

func pop(args ...Object) Object {
	if len(args) > 0 {
		return &Error{wrongNumberofArgs(len(args), 1)}
	}

	return NewMethod(SingletonNUll, POP)
}

func contains(args ...Object) Object {
	if len(args) > 1 || len(args) < 1 {
		return &Error{wrongNumberofArgs(len(args), 1)}
	}

	return NewMethod(args[0], CONTAIS)
}

func rango(args ...Object) Object {
	if len(args) > 1 || len(args) == 0 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if num, isNum := args[0].(*Number); isNum {
		if num.Value == 0 {
			return &Error{Message: "el rango debe mayor a 0"}
		}

		list := &List{Values: []Object{}}
		for i := 0; i < num.Value; i++ {
			list.Values = append(list.Values, &Number{i})
		}
		return list
	}

	return &Error{Message: unsoportedArgumentType("rango", types[args[0].Type()])}
}

func Tipo(args ...Object) Object {
	if len(args) > 1 || len(args) < 1 {
		return &Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return &String{Value: types[args[0].Type()]}
}

func input(scan *bufio.Scanner) string {
	scan.Scan()
	str := scan.Text()
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
	"largo":          NewBuiltin(Longitud),
	"escribir":       NewBuiltin(Escribir),
	"recibir":        NewBuiltin(Recibir),
	"recibir_entero": NewBuiltin(RecibirEntero),
	"tipo":           NewBuiltin(Tipo),
	"entero":         NewBuiltin(castInt),
	"texto":          NewBuiltin(castString),
	"rango":          NewBuiltin(rango),
	"agregar":        NewBuiltin(add),
	"pop":            NewBuiltin(pop),
	"popIndice":      NewBuiltin(remove),
	"contiene":       NewBuiltin(contains),
}
