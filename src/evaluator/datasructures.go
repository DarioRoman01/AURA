package evaluator

import (
	"aura/src/ast"
	obj "aura/src/object"
	"reflect"
	"strings"
)

func evaluateMap(mapa *ast.MapExpression, env *obj.Enviroment) obj.Object {
	mapObj := &obj.Map{Store: map[string]obj.Object{}}
	for _, keyVal := range mapa.Body {
		key := Evaluate(keyVal.Key, env)
		val := Evaluate(keyVal.Value, env)
		if err := mapObj.SetValues(key, val); err != nil {
			return newError("no se permiten llaves duplicadas")
		}
	}

	return mapObj
}

func evaluateArray(arr *ast.Array, env *obj.Enviroment) obj.Object {
	var list obj.List
	for _, val := range arr.Values {
		list.Values = append(list.Values, Evaluate(val, env))
	}

	return &list
}

func evaluateListReassigment(call *ast.CallList, list *obj.List, newVal ast.Expression, env *obj.Enviroment) obj.Object {
	index := Evaluate(call.Index, env)
	num, isNum := index.(*obj.Number)
	if !isNum {
		return &obj.Error{Message: "El indice debe ser un numero"}
	}
	if num.Value >= len(list.Values) {
		return &obj.Error{Message: "Indice fuera de rango"}
	}

	list.Values[num.Value] = Evaluate(newVal, env)
	return obj.SingletonNUll
}

func evaluateMapReassigment(hashMap *obj.Map, key obj.Object, value obj.Object) obj.Object {
	if err := hashMap.Get(key.Inspect()); err != nil {
		hashMap.SetValues(key, value)
		return obj.SingletonNUll
	}

	hashMap.UpdateKey(key, value)
	return obj.SingletonNUll
}

func evaluateReassigment(reassigment *ast.Reassignment, env *obj.Enviroment) obj.Object {
	// variable reassigment
	if variable, isVar := reassigment.Identifier.(*ast.Identifier); isVar {
		return evaluateVarReassigment(variable, reassigment.NewVal, env)
	}

	if callList, isCall := reassigment.Identifier.(*ast.CallList); isCall {
		evaluated := Evaluate(callList.ListIdent, env)
		// list reassigment
		if list, isList := evaluated.(*obj.List); isList {
			return evaluateListReassigment(callList, list, reassigment.NewVal, env)
		}

		// map reassigment
		if hashMap, isMap := evaluated.(*obj.Map); isMap {
			key := Evaluate(callList.Index, env)
			newVal := Evaluate(reassigment.NewVal, env)
			return evaluateMapReassigment(hashMap, key, newVal)
		}

		return notAList(evaluated.Inspect())
	}

	return notAVariable(reassigment.Identifier.TokenLiteral())
}

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

	default:
		return noSuchMethod(method.Inspect(), "list")
	}
}

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

		if strings.Contains(str.Value, val.Value) {
			return obj.SingletonTRUE
		}

		return obj.SingletonFALSE

	default:
		return noSuchMethod(method.Inspect(), "texto")
	}
}

func evaluateMethod(method *ast.MethodExpression, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(method.Obj, env)
	if list, isList := evaluated.(*obj.List); isList {
		listMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(listMethod.Inspect(), "list")
		}

		return evaluateListMethods(list, listMethod)
	}

	if hashMap, isMap := evaluated.(*obj.Map); isMap {
		mapMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(mapMethod.Inspect(), "mapa")
		}

		return evaluateMapMethods(hashMap, mapMethod)
	}

	if str, isStr := evaluated.(*obj.String); isStr {
		strMethod, isMethod := Evaluate(method.Method, env).(*obj.Method)
		if !isMethod {
			return noSuchMethod(strMethod.Inspect(), "texto")
		}

		return evaluateStringMethod(str, strMethod)
	}

	return noSuchMethod(method.Method.Str(), method.Obj.Str())
}
