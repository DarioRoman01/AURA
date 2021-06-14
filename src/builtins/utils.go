package builtins

import (
	"fmt"
	obj "katan/src/object"
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

	if startVal.Value > endVal.Value {
		return &obj.Error{Message: "El valor de inicio no puede ser mayor al de el final"}
	}

	list := &obj.List{Values: []obj.Object{}}
	for i := startVal.Value; i < endVal.Value; i++ {
		list.Values = append(list.Values, &obj.Number{Value: i})
	}

	return list
}
