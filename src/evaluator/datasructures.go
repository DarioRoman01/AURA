package evaluator

import (
	"aura/src/ast"
	obj "aura/src/object"
	"reflect"
	"strings"
)

// evaluate a map object
func evaluateMap(mapa *ast.MapExpression, env *obj.Enviroment) obj.Object {
	mapObj := &obj.Map{Store: map[string]obj.Object{}}

	// we loop and evaluate all the key value pairs in the expression
	// and add them to the hashMap
	for _, keyVal := range mapa.Body {
		key := Evaluate(keyVal.Key, env)
		val := Evaluate(keyVal.Value, env)

		if err := mapObj.SetValues(key, val); err != nil {
			// duplicated keys
			return newError(err.Error())
		}
	}

	return mapObj
}

// evaluate an array expression
func evaluateArray(arr *ast.Array, env *obj.Enviroment) obj.Object {
	var list obj.List
	// we loop and evaluated all the values in the array expression to add them to the list
	for _, val := range arr.Values {
		list.Values = append(list.Values, Evaluate(val, env))
	}

	return &list
}

// evaluate a list reassigment by index like:
//		arr[0] = 2;
func evaluateListReassigment(call *ast.CallList, list *obj.List, newVal ast.Expression, env *obj.Enviroment) obj.Object {
	index := Evaluate(call.Index, env)
	num, isNum := index.(*obj.Number)
	if !isNum {
		// the index is not a number
		return newError("El indice debe ser un numero")
	}
	if num.Value >= len(list.Values) {
		// the index is grater than the length of the list
		return indexOutOfRange(num.Value, len(list.Values))
	}

	list.Values[num.Value] = Evaluate(newVal, env)
	return obj.SingletonNUll
}

// evaluate a HashMap reassigment
func evaluateMapReassigment(hashMap *obj.Map, key obj.Object, value obj.Object) obj.Object {
	// we dont care if the key doesnt exist
	// we just add the key value pair to the map
	hashMap.Store[key.Inspect()] = value
	return obj.SingletonNUll
}

// evaluate a reassigment expression
func evaluateReassigment(reassigment *ast.Reassignment, env *obj.Enviroment) obj.Object {
	switch exp := reassigment.Identifier.(type) {
	case *ast.Identifier:
		return evaluateVarReassigment(exp, reassigment.NewVal, env)

	case *ast.ClassFieldCall:
		evaluated := Evaluate(exp.Class, env)
		if _, isErr := evaluated.(*obj.Error); isErr {
			return evaluated
		}

		if class, isClass := evaluated.(*obj.ClassInstance); isClass {
			return evaluateFieldReassigment(exp, class, reassigment.NewVal)
		}

		return notAClass(evaluated.Inspect())

	case *ast.CallList:
		evaluated := Evaluate(exp.ListIdent, env)

		if list, isList := evaluated.(*obj.List); isList {
			return evaluateListReassigment(exp, list, reassigment.NewVal, env)
		}

		if hashMap, isMap := evaluated.(*obj.Map); isMap {
			key := Evaluate(exp.Index, env)
			newVal := Evaluate(reassigment.NewVal, env)
			return evaluateMapReassigment(hashMap, key, newVal)
		}

		return notAList(evaluated.Inspect())

	default:
		return notAVariable(reassigment.Identifier.TokenLiteral())
	}
}

// evaluate a list method if the method is valid will be applied else will return an error
func evaluateListMethods(list *obj.List, method *obj.Method) obj.Object {
	switch method.MethodType {
	case obj.POP:
		return list.Pop()

	case obj.APPEND:
		list.Add(method.Value)
		return obj.SingletonNUll

	case obj.REMOVE:
		index := method.Value.(*obj.Number)
		return list.RemoveAt(index.Value)

	case obj.CONTAIS:
		for _, val := range list.Values {
			if reflect.DeepEqual(val, method.Value) {
				return obj.SingletonTRUE
			}
		}

		return obj.SingletonFALSE

	case obj.MAP:
		var newList obj.List
		fn := method.Value.(*obj.Def)
		for _, val := range list.Values {
			newList.Values = append(newList.Values, applyFunction(fn, val))
		}

		return &newList

	case obj.FOREACH:
		fn := method.Value.(*obj.Def)
		for _, val := range list.Values {
			applyFunction(fn, val)
		}
		return obj.SingletonNUll

	case obj.FILTER:
		var newList obj.List
		fn := method.Value.(*obj.Def)
		for _, val := range list.Values {
			eval := applyFunction(fn, val)
			if isTruthy(eval) {
				newList.Values = append(newList.Values, val)
			}
		}

		return &newList

	case obj.COUNT:
		var count obj.Number
		fn := method.Value.(*obj.Def)
		for _, val := range list.Values {
			eval := applyFunction(fn, val)
			if isTruthy(eval) {
				count.Value++
			}
		}

		return &count

	default:
		return noSuchMethod(method.Inspect(), "list")
	}
}

// evaluate a map method if the method is valid will be applied else will return an error
func evaluateMapMethods(hashMap *obj.Map, method *obj.Method) obj.Object {
	switch method.MethodType {
	case obj.CONTAIS:
		if hashMap.Get(method.Value.Inspect()) != obj.NullVAlue {
			return obj.SingletonTRUE
		}

		return obj.SingletonFALSE

	case obj.VALUES:
		list := &obj.List{Values: []obj.Object{}}
		for _, val := range hashMap.Store {
			list.Values = append(list.Values, val)
		}
		return list

	default:
		return noSuchMethod(method.Inspect(), "mapa")
	}
}

// evaluate a string method if the method is valid will be applied else will return an error
func evaluateStringMethod(str *obj.String, method *obj.Method) obj.Object {
	switch method.MethodType {
	case obj.UPPER:
		str.Value = strings.ToUpper(str.Value)
		return str

	case obj.LOWER:
		str.Value = strings.ToLower(str.Value)
		return str

	case obj.CONTAIS:
		val, isStr := method.Value.(*obj.String)
		if !isStr {
			return newError("La funcion contiene solo puede recibir caracteres o cadenas")
		}

		return str.Contains(val.Value)

	case obj.ISUPPER:
		return str.IsUpper()

	case obj.ISLOWER:
		return str.IsLower()

	case obj.SPLIT:
		separator := method.Value.(*obj.String)
		return str.Split(separator.Value)

	default:
		return noSuchMethod(method.Inspect(), "texto")
	}
}

// evaluate a method expression
func evaluateMethod(method *ast.MethodExpression, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(method.Obj, env)
	// we check the type of the method object
	switch data := evaluated.(type) {

	case *obj.List:
		listMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(listMethod.Inspect(), "list")
		}

		return evaluateListMethods(data, listMethod)

	case *obj.Map:
		mapMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(mapMethod.Inspect(), "mapa")
		}

		return evaluateMapMethods(data, mapMethod)

	case *obj.String:
		strMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(strMethod.Inspect(), "texto")
		}

		return evaluateStringMethod(data, strMethod)

	default:
		// the object has no methods
		return noSuchMethod(method.Method.Str(), method.Obj.Str())
	}
}
