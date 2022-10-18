package builtins

import (
	obj "aura/src/object"
	"fmt"
	"strings"
)

func makeOneArgList(arg obj.Object) obj.Object {
	if num, isNum := arg.(*obj.Number); isNum {
		if num.Value == 0 {
			return &obj.Error{Message: "el rango debe mayor a 0"}
		}

		list := &obj.List{Values: []obj.Object{}}
		for i := 0; i < num.Value; i++ {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
		return list
	}

	return unsoportedArgumentType("rango", obj.Types[arg.Type()])
}

func makeTwoArgList(start, end obj.Object) obj.Object {
	startVal, isNum := start.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: fmt.Sprintf("El valor de inicio debe ser un entero %s", start.Inspect())}
	}

	endVal, isNum := end.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: fmt.Sprintf("El valor de inicio debe ser un entero %s", end.Inspect())}
	}

	list := &obj.List{Values: []obj.Object{}}
	if startVal.Value > endVal.Value {
		for i := startVal.Value; i > endVal.Value; i-- {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
	} else {
		for i := startVal.Value; i < endVal.Value; i++ {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
	}

	return list
}

func makeTreArgList(start, end, pass obj.Object) obj.Object {
	startVal, isNum := start.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: fmt.Sprintf("El valor de inicio debe ser un entero %s", start.Inspect())}
	}

	endVal, isNum := end.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: fmt.Sprintf("El valor final debe ser un entero %s", end.Inspect())}
	}

	passVal, isNum := pass.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: fmt.Sprintf("Los pasos deben debe ser un entero %s", pass.Inspect())}
	}

	if passVal.Value < 1 {
		return &obj.Error{Message: fmt.Sprintf("Los pasos deben debe ser mayor a 0 %s", pass.Inspect())}
	}

	list := &obj.List{Values: []obj.Object{}}
	if startVal.Value > endVal.Value {
		for i := startVal.Value; i > endVal.Value; i -= passVal.Value {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
	} else {
		for i := startVal.Value; i < endVal.Value; i += passVal.Value {
			list.Values = append(list.Values, &obj.Number{Value: i})
		}
	}

	return list
}

func formatString(str string, args []obj.Object) string {
	for i := 0; i < len(args); i++ {
		str = strings.Replace(str, "{}", args[i].Inspect(), 1)
	}

	return str
}
