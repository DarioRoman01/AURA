package builtins

import (
	"bufio"
	"fmt"
	obj "katan/src/object"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var scanner = bufio.NewScanner(os.Stdin)

func wrongNumberofArgs(found, actual int) string {
	return fmt.Sprintf("numero incorrecto de argumentos para largo, se recibieron %d, se requieren %d", found, actual)
}

func unsoportedArgumentType(funcname, objType string) string {
	return fmt.Sprintf("argumento para %s no valido, se recibio %s", funcname, objType)
}

func Longitud(args ...obj.Object) obj.Object {
	if len(args) != 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	switch arg := args[0].(type) {

	case *obj.String:
		return &obj.Number{Value: utf8.RuneCountInString(arg.Value)}

	case *obj.List:
		return &obj.Number{Value: len(arg.Values)}

	case *obj.Map:
		return &obj.Number{Value: len(arg.Store)}

	default:
		return &obj.Error{Message: unsoportedArgumentType("largo", obj.Types[args[0].Type()])}
	}
}

func Escribir(args ...obj.Object) obj.Object {
	var buff strings.Builder

	for _, arg := range args {
		switch node := arg.(type) {

		case *obj.String:
			buff.WriteString(node.Inspect())

		case *obj.Number:
			buff.WriteString(node.Inspect())

		case *obj.List:
			buff.WriteString(node.Inspect())

		case *obj.Bool:
			buff.WriteString(node.Inspect())

		case *obj.Map:
			buff.WriteString(node.Inspect())

		default:
			return &obj.Error{Message: unsoportedArgumentType("escribir", obj.Types[node.Type()])}
		}
	}

	fmt.Println(buff.String())
	return obj.SingletonNUll
}

func Recibir(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if len(args) == 0 {
		str := input(scanner)
		return &obj.String{Value: str}
	}

	if arg, isString := args[0].(*obj.String); isString {
		fmt.Print(arg.Inspect())
		str := input(scanner)
		return &obj.String{Value: str}
	}

	return &obj.Error{Message: unsoportedArgumentType("recibir", obj.Types[args[0].Type()])}
}

func castInt(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isString := args[0].(*obj.String); isString {
		return toInt(arg.Value)
	}

	return &obj.Error{Message: unsoportedArgumentType("recibir", obj.Types[args[0].Type()])}
}

func castString(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if arg, isNumber := args[0].(*obj.Number); isNumber {
		strInt := strconv.Itoa(arg.Value)
		return &obj.String{Value: strInt}
	}

	return &obj.Error{Message: unsoportedArgumentType("recibir", obj.Types[args[0].Type()])}
}

func RecibirEntero(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if len(args) == 0 {
		strInt := input(scanner)
		return toInt(strInt)
	}

	if arg, isString := args[0].(*obj.String); isString {
		fmt.Print(arg.Inspect())
		strInt := input(scanner)
		return toInt(strInt)
	}

	return &obj.Error{
		Message: unsoportedArgumentType("recbir_entero", obj.Types[args[0].Type()]),
	}

}

func add(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if num, isNumber := args[0].(*obj.Number); isNumber {
		return obj.NewMethod(num, obj.APPEND)
	}

	return &obj.Error{Message: unsoportedArgumentType("add", obj.Types[args[0].Type()])}
}

func remove(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if num, isNumber := args[0].(*obj.Number); isNumber {
		return obj.NewMethod(num, obj.REMOVE)
	}

	return &obj.Error{Message: unsoportedArgumentType("add", obj.Types[args[0].Type()])}
}

func pop(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return obj.NewMethod(obj.SingletonNUll, obj.POP)
}

func contains(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) < 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return obj.NewMethod(args[0], obj.CONTAIS)
}

func values(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return obj.NewMethod(obj.SingletonNUll, obj.VALUES)
}

func rango(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	if num, isNum := args[0].(*obj.Number); isNum {
		if num.Value == 0 {
			return &obj.Error{Message: "el rango debe mayor a 0"}
		}

		list := &obj.List{Values: []obj.Object{}}
		for i := 0; i < num.Value; i++ {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
		return list
	}

	return &obj.Error{Message: unsoportedArgumentType("rango", obj.Types[args[0].Type()])}
}

func Tipo(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) < 1 {
		return &obj.Error{Message: wrongNumberofArgs(len(args), 1)}
	}

	return &obj.String{Value: obj.Types[args[0].Type()]}
}

func input(scan *bufio.Scanner) string {
	scan.Scan()
	str := scan.Text()
	return str
}

func toInt(str string) obj.Object {
	number, err := strconv.Atoi(str)
	if err != nil {
		return &obj.Error{Message: fmt.Sprintf("No se puede parsear como entero %s", str)}
	}

	return &obj.Number{Value: number}
}

var BUILTINS = map[string]*obj.Builtin{
	"largo":          obj.NewBuiltin(Longitud),
	"escribir":       obj.NewBuiltin(Escribir),
	"recibir":        obj.NewBuiltin(Recibir),
	"recibir_entero": obj.NewBuiltin(RecibirEntero),
	"tipo":           obj.NewBuiltin(Tipo),
	"entero":         obj.NewBuiltin(castInt),
	"texto":          obj.NewBuiltin(castString),
	"rango":          obj.NewBuiltin(rango),
	"agregar":        obj.NewBuiltin(add),
	"pop":            obj.NewBuiltin(pop),
	"popIndice":      obj.NewBuiltin(remove),
	"contiene":       obj.NewBuiltin(contains),
	"valores":        obj.NewBuiltin((values)),
}
