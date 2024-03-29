package builtins

import (
	obj "aura/src/object"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// use singleton pattern for the scanner
var (
	scanner = bufio.NewScanner(os.Stdin)
	writer  = bufio.NewWriter(os.Stdout)
)

// return an error indicating the the builtin has wrong number of args
func wrongNumberofArgs(funcName string, found, actual int) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("numero incorrecto de argumentos para %s, se recibieron %d, se requieren %d", funcName, found, actual),
	}
}

func unsoportedArgumentType(funcname, objType string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("argumento para %s no valido, se recibio %s", funcname, objType),
	}
}

// Longitud return the length of the object if is suported by the function
func Longitud(args ...obj.Object) obj.Object {
	if len(args) != 1 {
		return wrongNumberofArgs("largo", len(args), 1)
	}

	switch arg := args[0].(type) {

	case *obj.String:
		return &obj.Number{Value: utf8.RuneCountInString(arg.Value)}

	case *obj.List:
		return &obj.Number{Value: len(arg.Values)}

	case *obj.Map:
		return &obj.Number{Value: len(arg.Store)}

	default:
		return unsoportedArgumentType("largo", obj.Types[args[0].Type()])
	}
}

// same as println function
func Escribir(args ...obj.Object) obj.Object {
	var buff strings.Builder
	for _, arg := range args {
		buff.WriteString(arg.Inspect())
	}

	defer writer.Flush()
	writer.WriteString(buff.String() + "\n")
	return obj.SingletonNUll
}

// same as python input function
func Recibir(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return wrongNumberofArgs("recibir", len(args), 1)
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

	return unsoportedArgumentType("recibir", obj.Types[args[0].Type()])
}

// convert a string object to int object
func castInt(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return wrongNumberofArgs("entero", len(args), 1)
	}

	switch node := args[0].(type) {
	case *obj.Number:
		return node

	case *obj.String:
		return toInt(node.Value)

	case *obj.Float:
		return &obj.Number{Value: int(node.Value)}

	default:
		return unsoportedArgumentType("entero", obj.Types[args[0].Type()])
	}
}

// convert a string object to int object
func castString(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return wrongNumberofArgs("texto", len(args), 1)
	}

	return &obj.String{Value: args[0].Inspect()}
}

// convert a string or integer object to a float
func castFloat(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return wrongNumberofArgs("flotante", len(args), 1)
	}

	switch node := args[0].(type) {
	case *obj.Number:
		return &obj.Float{Value: float64(node.Value)}

	case *obj.String:
		val, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			return &obj.Error{Message: fmt.Sprintf("no se pudo parsear como flotante %s", node.Value)}
		}

		return &obj.Float{Value: val}

	default:
		return unsoportedArgumentType("flotante", obj.Types[args[0].Type()])
	}
}

func formatrArgs(args ...obj.Object) obj.Object {
	if len(args) <= 1 {
		return wrongNumberofArgs("format", 0, 100)
	}

	str, isStr := args[0].(*obj.String)
	if !isStr {
		return &obj.Error{
			Message: "el primer argumento para formatear debe ser un string",
		}
	}

	formated := formatString(str.Value, args[1:])
	return &obj.String{Value: formated}
}

func printF(args ...obj.Object) obj.Object {
	if len(args) <= 1 {
		return wrongNumberofArgs("format", 0, 100)
	}

	str, isStr := args[0].(*obj.String)
	if !isStr {
		return &obj.Error{
			Message: "el primer argumento para formatear debe ser un string",
		}
	}

	formated := formatString(str.Value, args[1:])
	defer writer.Flush()
	writer.WriteString(formated + "\n")
	return obj.SingletonNUll
}

// same as python range function
func rango(args ...obj.Object) obj.Object {
	switch len(args) {
	case 1:
		return makeOneArgList(args[0])

	case 2:
		return makeTwoArgList(args[0], args[1])

	case 3:
		return makeTreArgList(args[0], args[1], args[2])

	default:
		return wrongNumberofArgs("rango", len(args), 3)
	}
}

func slep(args ...obj.Object) obj.Object {
	if len(args) > 1 {
		return wrongNumberofArgs("dormir", len(args), 1)
	}

	switch arg := args[0].(type) {
	case *obj.Number:
		time.Sleep(time.Duration(arg.Value * int(time.Second)))
		return obj.SingletonNUll

	case *obj.Float:
		time.Sleep(time.Duration(arg.Value * float64(time.Second)))
		return obj.SingletonNUll

	default:
		return unsoportedArgumentType("dormir", obj.Types[arg.Type()])
	}
}

// return the type of the object
func Tipo(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) < 1 {
		return wrongNumberofArgs("tipo", len(args), 1)
	}

	return &obj.String{Value: obj.Types[args[0].Type()]}
}

func sum(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) < 1 {
		return wrongNumberofArgs("sum", len(args), 1)
	}

	if list, isList := args[0].(*obj.List); isList {
		var res obj.Number
		for _, value := range list.Values {
			switch item := value.(type) {
			case *obj.Number:
				res.Value += item.Value

			case *obj.Float:
				res.Value += int(item.Value)

			default:
				return unsoportedArgumentType("suma", obj.Types[value.Type()])
			}
		}

		return &res
	}

	return unsoportedArgumentType("suma", obj.Types[args[0].Type()])
}

// return the absolute value of the given number
func abs(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return wrongNumberofArgs("abs", len(args), 1)
	}

	switch node := args[0].(type) {
	case *obj.Number:
		return &obj.Number{Value: int(math.Abs(float64(node.Value)))}

	case *obj.Float:
		return &obj.Float{Value: math.Abs(node.Value)}

	default:
		return unsoportedArgumentType("abs", obj.Types[args[0].Type()])
	}
}

// input function to recibe input from console
func input(scan *bufio.Scanner) string {
	scan.Scan()
	str := scan.Text()
	return str
}

// perform the string to int conversion and handle the posibe errror
func toInt(str string) obj.Object {
	number, err := strconv.Atoi(str)
	if err != nil {
		return &obj.Error{Message: fmt.Sprintf("No se puede parsear como entero %s", str)}
	}

	return &obj.Number{Value: number}
}

var BUILTINS = map[string]*obj.Builtin{
	"largo":        obj.NewBuiltin(Longitud),
	"escribir":     obj.NewBuiltin(Escribir),
	"recibir":      obj.NewBuiltin(Recibir),
	"tipo":         obj.NewBuiltin(Tipo),
	"entero":       obj.NewBuiltin(castInt),
	"texto":        obj.NewBuiltin(castString),
	"rango":        obj.NewBuiltin(rango),
	"agregar":      obj.NewBuiltin(add),
	"pop":          obj.NewBuiltin(pop),
	"popIndice":    obj.NewBuiltin(remove),
	"contiene":     obj.NewBuiltin(contains),
	"valores":      obj.NewBuiltin(values),
	"mayusculas":   obj.NewBuiltin(toUppper),
	"minusculas":   obj.NewBuiltin(toLower),
	"dormir":       obj.NewBuiltin(slep),
	"es_mayuscula": obj.NewBuiltin(isUpper),
	"es_minuscula": obj.NewBuiltin(isLower),
	"formatear":    obj.NewBuiltin(formatrArgs),
	"escribirF":    obj.NewBuiltin(printF),
	"map":          obj.NewBuiltin(mapList),
	"porCada":      obj.NewBuiltin(forEach),
	"filtrar":      obj.NewBuiltin(filter),
	"contar":       obj.NewBuiltin(count),
	"separar":      obj.NewBuiltin(split),
	"abs":          obj.NewBuiltin(abs),
	"flotante":     obj.NewBuiltin(castFloat),
	"suma":         obj.NewBuiltin(sum),
}
