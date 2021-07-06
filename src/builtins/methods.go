package builtins

import (
	obj "aura/src/object"
)

func add(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return wrongNumberofArgs("agregar", len(args), 1)
	}

	if num, isNumber := args[0].(*obj.Number); isNumber {
		return obj.NewMethod(num, obj.APPEND)
	}

	return unsoportedArgumentType("add", obj.Types[args[0].Type()])
}

func remove(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) == 0 {
		return wrongNumberofArgs("popIndice", len(args), 1)
	}

	if num, isNumber := args[0].(*obj.Number); isNumber {
		return obj.NewMethod(num, obj.REMOVE)
	}

	return unsoportedArgumentType("popIndice", obj.Types[args[0].Type()])
}

func pop(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return wrongNumberofArgs("pop", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.POP)
}

func contains(args ...obj.Object) obj.Object {
	if len(args) > 1 || len(args) < 1 {
		return wrongNumberofArgs("contiene", len(args), 1)
	}

	return obj.NewMethod(args[0], obj.CONTAIS)
}

func values(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return wrongNumberofArgs("valores", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.VALUES)
}

func toUppper(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return wrongNumberofArgs("mayusculas", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.UPPER)
}

func toLower(args ...obj.Object) obj.Object {
	if len(args) > 0 {
		return wrongNumberofArgs("minusculas", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.LOWER)
}

func isUpper(args ...obj.Object) obj.Object {
	if len(args) != 0 {
		return wrongNumberofArgs("es_mayuscula", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.ISUPPER)
}

func isLower(args ...obj.Object) obj.Object {
	if len(args) != 0 {
		return wrongNumberofArgs("es_minuscula", len(args), 0)
	}

	return obj.NewMethod(obj.SingletonNUll, obj.ISLOWER)
}
