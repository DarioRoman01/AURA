package evaluator

import (
	"fmt"
	"katan/src/ast"
	b "katan/src/builtins"
	obj "katan/src/object"
	"unicode/utf8"
)

// evlauate given nodes of an ast
func Evaluate(baseNode ast.ASTNode, env *obj.Enviroment) obj.Object {
	switch node := baseNode.(type) {

	case ast.Program:
		return evaluateProgram(node, env)

	case *ast.ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression, env)

	case *ast.Array:
		return evaluateArray(node, env)

	case *ast.Integer:
		CheckIsNotNil(node.Value)
		return &obj.Number{Value: *node.Value}

	case *ast.Boolean:
		CheckIsNotNil(node.Value)
		return toBooleanObject(*node.Value)

	case *ast.Prefix:
		CheckIsNotNil(node.Rigth)
		rigth := Evaluate(node.Rigth, env)
		CheckIsNotNil(rigth)
		return evaluatePrefixExpression(node.Operator, rigth)

	case *ast.Infix:
		go CheckIsNotNil(node.Left)
		CheckIsNotNil(node.Rigth)
		left := Evaluate(node.Left, env)
		rigth := Evaluate(node.Rigth, env)
		go CheckIsNotNil(left)
		CheckIsNotNil(rigth)
		return evaluateInfixExpression(node.Operator, left, rigth)

	case *ast.Block:
		return evaluateBLockStaments(node, env)

	case *ast.While:
		return evaluateWhileExpression(node, env)

	case *ast.If:
		return evaluateIfExpression(node, env)

	case *ast.ReturnStament:
		CheckIsNotNil(node.ReturnValue)
		value := Evaluate(node.ReturnValue, env)
		CheckIsNotNil(value)
		return &obj.Return{Value: value}

	case *ast.For:
		CheckIsNotNil(node.Condition)
		CheckIsNotNil(node.Body)
		return evaluateFor(node, env)

	case *ast.MethodExpression:
		CheckIsNotNil(node.Method)
		CheckIsNotNil(node.Obj)
		return evaluateMethod(node, env)

	case *ast.RangeExpression:
		CheckIsNotNil(node.Range)
		CheckIsNotNil(node.Variable)
		return evaluateRange(node, env)

	case *ast.CallList:
		return evaluateCallList(node, env)

	case *ast.Suffix:
		CheckIsNotNil(node.Left)
		CheckIsNotNil(node.Operator)
		left := Evaluate(node.Left, env)
		return evaluateSuffixExpression(node.Operator, left)

	case *ast.Reassignment:
		CheckIsNotNil(node.Identifier)
		CheckIsNotNil(node.NewVal)
		return evaluateReassigment(node, env)

	case *ast.NullExpression:
		return obj.NullVAlue

	case *ast.MapExpression:
		CheckIsNotNil(node.Body)
		return evaluateMap(node, env)

	case *ast.LetStatement:
		CheckIsNotNil(node.Value)
		value := Evaluate(node.Value, env)
		CheckIsNotNil(node.Name)
		env.SetItem(node.Name.Value, value)
		return obj.SingletonNUll

	case *ast.Identifier:
		return evaluateIdentifier(node, env)

	case *ast.Function:
		CheckIsNotNil(node.Body)
		return obj.NewDef(node.Body, env, node.Parameters...)

	case *ast.Call:
		function := Evaluate(node.Function, env)
		CheckIsNotNil(node.Arguments)
		args := evaluateExpression(node.Arguments, env)
		CheckIsNotNil(function)
		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}

	default:
		return obj.SingletonNUll
	}
}

// generates a new function object
func applyFunction(fn obj.Object, args []obj.Object) obj.Object {
	if function, isFn := fn.(*obj.Def); isFn {
		extendedEnviron := extendFunctionEnviroment(function, args)
		evaluated := Evaluate(function.Body, extendedEnviron)
		CheckIsNotNil(evaluated)
		return unwrapReturnValue(evaluated)

	} else if builtin, isBuiltin := fn.(*obj.Builtin); isBuiltin {
		return builtin.Fn(args...)
	}

	return notAFunction(obj.Types[fn.Type()])
}

// unwrap the return value of a function
func unwrapReturnValue(object obj.Object) obj.Object {
	if obj, isReturn := object.(*obj.Return); isReturn {
		return obj.Value
	}
	return object
}

// create a new enviroment when a function is called
func extendFunctionEnviroment(fn *obj.Def, args []obj.Object) *obj.Enviroment {
	env := obj.NewEnviroment(fn.Env)
	for idx, param := range fn.Parameters {
		env.SetItem(param.Value, args[idx])
	}

	return env
}

func evaluateVarReassigment(variable *ast.Identifier, newVal ast.Expression, env *obj.Enviroment) obj.Object {
	_, exists := env.GetItem(variable.Value)
	if !exists {
		return unknownIdentifier(variable.Value)
	}

	env.Store[variable.Value] = Evaluate(newVal, env)
	return obj.SingletonNUll
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
	if err := hashMap.Get(string(hashMap.Serialize(key))); err != nil {
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

	default:
		return noSuchMethod(method.Inspect(), "list")
	}
}

func evaluateMapMethods(hashMap *obj.Map, method *obj.Method) obj.Object {
	switch method.MethodType {
	case obj.CONTAIS:
		return obj.NewBool(hashMap.Get(string(hashMap.Serialize(method.Value))) != obj.NullVAlue)

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

	return &obj.Error{
		Message: fmt.Sprintf("%s no tiene metodos", obj.Types[evaluated.Type()]),
	}
}

func evaluateFor(forLoop *ast.For, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(forLoop.Condition, env)
	if iter, isIter := evaluated.(*obj.Iterator); isIter {
		val := forLoop.Condition.(*ast.RangeExpression).Variable.(*ast.Identifier).Value
		for iter.Next() != nil {
			evaluated = Evaluate(forLoop.Body, env)
			if returnVal, isReturn := evaluated.(*obj.Return); isReturn {
				return returnVal
			}

			env.Store[val] = iter.Current
		}
		return obj.SingletonNUll
	}

	if err, isError := evaluated.(*obj.Error); isError {
		return err
	}

	return newError("syntax error")
}

func evaluateMap(mapa *ast.MapExpression, env *obj.Enviroment) obj.Object {
	mapObj := &obj.Map{Store: map[string]obj.Object{}}
	for _, keyVal := range mapa.Body {
		key := Evaluate(keyVal.Key, env)
		val := Evaluate(keyVal.Value, env)
		if err := mapObj.SetValues(key, val); err != nil {
			return &obj.Error{Message: "No se permiten llaves duplicadas"}
		}
	}

	return mapObj
}

func evaluateRange(rangeExpress *ast.RangeExpression, env *obj.Enviroment) obj.Object {
	if list, isList := Evaluate(rangeExpress.Range, env).(*obj.List); isList {
		val, isVar := rangeExpress.Variable.(*ast.Identifier)
		if !isVar {
			return newError("no es una variable")
		}
		iter := obj.NewIterator(list.Values[0], list.Values)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	if str, isStr := Evaluate(rangeExpress.Range, env).(*obj.String); isStr {
		val, isVar := rangeExpress.Variable.(*ast.Identifier)
		if !isVar {
			return newError("no es una variable")
		}

		list := makeStringList(str.Value)
		iter := obj.NewIterator(list[0], list)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	return newError("Rango invalido")
}

// check that the given value is not nil
func CheckIsNotNil(val interface{}) {
	defer handlePanic()
	if val == nil {
		panic("expression or stament cannot be nil :(")
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("syntax Error")
	}
}

// evluate a block statement
func evaluateBLockStaments(block *ast.Block, env *obj.Enviroment) obj.Object {
	var result obj.Object = nil
	for _, statement := range block.Staments {
		result = Evaluate(statement, env)
		if result != nil && result.Type() == obj.RETURNTYPE || result.Type() == obj.ERROR {
			return result
		}
	}

	return result
}

// evaluate an slice of expressions
func evaluateExpression(expressions []ast.Expression, env *obj.Enviroment) []obj.Object {
	var result []obj.Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, env)
		CheckIsNotNil(evaluated)
		result = append(result, evaluated)
	}

	return result
}

func evaluateArray(arr *ast.Array, env *obj.Enviroment) obj.Object {
	var list obj.List

	for _, val := range arr.Values {
		list.Values = append(list.Values, Evaluate(val, env))
	}

	return &list
}

// check if given identifier exists in the enviroment
func evaluateIdentifier(node *ast.Identifier, env *obj.Enviroment) obj.Object {
	value, exists := env.GetItem(node.Value)
	if !exists {
		builtint, exists := b.BUILTINS[node.Value]
		if !exists {
			return unknownIdentifier(node.Value)
		}

		return builtint
	}

	return value
}

// evaluate program node
func evaluateProgram(program ast.Program, env *obj.Enviroment) obj.Object {
	var result obj.Object
	for _, statement := range program.Staments {
		result = Evaluate(statement, env)

		if returnObj, isReturn := result.(*obj.Return); isReturn {
			return returnObj.Value
		} else if err, isError := result.(*obj.Error); isError {
			return err
		}
	}

	return result
}

// change the bool value of the object
func evaluateBangOperatorExpression(rigth obj.Object) obj.Object {
	switch {
	case rigth == obj.SingletonTRUE:
		return obj.SingletonFALSE

	case rigth == obj.SingletonFALSE:
		return obj.SingletonTRUE

	case rigth == nil:
		return obj.SingletonTRUE

	default:
		return obj.SingletonFALSE
	}
}

func evaluateWhileExpression(whileExpression *ast.While, env *obj.Enviroment) obj.Object {
	CheckIsNotNil(whileExpression.Condition)
	condition := Evaluate(whileExpression.Condition, env)

	CheckIsNotNil(condition)
	for isTruthy(condition) {
		evaluated := Evaluate(whileExpression.Body, env)
		if returnVal, isReturn := evaluated.(*obj.Return); isReturn {
			return returnVal
		}

		condition = Evaluate(whileExpression.Condition, env)
	}

	return obj.SingletonNUll
}

func evaluateCallList(call *ast.CallList, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.ListIdent, env)
	if list, isList := evaluated.(*obj.List); isList {
		evaluated := Evaluate(call.Index, env)
		num, isNumber := evaluated.(*obj.Number)
		if !isNumber {
			return &obj.Error{Message: "El indice debe ser un entero"}
		}

		if num.Value >= len(list.Values) {
			return &obj.Error{Message: "Indice fuera de rango"}
		}

		return list.Values[num.Value]
	}

	if hashMap, isMap := evaluated.(*obj.Map); isMap {
		evaluated := Evaluate(call.Index, env)
		return hashMap.Get(string(hashMap.Serialize(evaluated)))
	}

	if str, isStr := evaluated.(*obj.String); isStr {
		evaluated := Evaluate(call.Index, env)
		num, isNumber := evaluated.(*obj.Number)
		if !isNumber {
			return &obj.Error{Message: "El indice debe ser un entero"}
		}

		if num.Value >= utf8.RuneCountInString(str.Value) {
			return &obj.Error{Message: "Indice fuera de rango"}
		}

		return &obj.String{Value: string(str.Value[num.Value])}
	}

	return notAList(obj.Types[evaluated.Type()])
}

func evaluateIfExpression(ifExpression *ast.If, env *obj.Enviroment) obj.Object {
	CheckIsNotNil(ifExpression.Condition)
	condition := Evaluate(ifExpression.Condition, env)

	CheckIsNotNil(condition)
	if isTruthy(condition) {
		CheckIsNotNil(ifExpression.Consequence)
		return Evaluate(ifExpression.Consequence, env)

	} else if ifExpression.Alternative != nil {
		return Evaluate(ifExpression.Alternative, env)
	}

	return obj.SingletonNUll
}

// check that the current object is true or false
func isTruthy(object obj.Object) bool {
	switch {
	case object == obj.SingletonNUll:
		return false

	case object == obj.SingletonTRUE:
		return true

	case object == obj.SingletonFALSE:
		return false

	default:
		return true
	}
}

// evluate infix expressions between objects
func evaluateInfixExpression(operator string, left obj.Object, right obj.Object) obj.Object {
	switch {

	case left.Type() == obj.INTEGERS && right.Type() == obj.INTEGERS:
		return evaluateIntegerInfixExpression(operator, left, right)

	case left.Type() == obj.STRINGTYPE && right.Type() == obj.STRINGTYPE:
		return evaluateStringInfixExpression(operator, left, right)

	case left.Type() == obj.BOOLEAN && right.Type() == obj.BOOLEAN:
		return evaluateBoolInfixExpression(operator, left.(*obj.Bool), right.(*obj.Bool))

	case operator == "==":
		return toBooleanObject(left == right)

	case operator == "!=":
		return toBooleanObject(left != right)

	case left.Type() != right.Type():
		return typeMismatchError(
			obj.Types[left.Type()],
			operator,
			obj.Types[right.Type()],
		)

	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[right.Type()],
		)
	}

}

func evaluateBoolInfixExpression(operator string, left *obj.Bool, rigth *obj.Bool) obj.Object {
	switch operator {

	case "||":
		return toBooleanObject(left.Value || rigth.Value)

	case "&&":
		return toBooleanObject(left.Value && rigth.Value)

	case "==":
		return toBooleanObject(left == rigth)

	case "!=":
		return toBooleanObject(left != rigth)

	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

func evaluateStringInfixExpression(operator string, left obj.Object, rigth obj.Object) obj.Object {
	leftVal := left.(*obj.String).Value
	rigthVal := rigth.(*obj.String).Value

	switch operator {
	case "+":
		return &obj.String{Value: leftVal + rigthVal}
	case "==":
		return toBooleanObject(leftVal == rigthVal)
	case "!=":
		return toBooleanObject(leftVal != rigthVal)
	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

func evaluateSuffixExpression(operator string, left obj.Object) obj.Object {
	if num, isNumber := left.(*obj.Number); isNumber {
		switch operator {
		case "++":
			num.Value++
			return num

		case "--":
			num.Value--
			return num

		case "**":
			num.Value *= num.Value
			return num
		default:
			return &obj.Error{Message: "Operador desconocido para entero"}
		}
	}

	return &obj.Error{Message: "No es un numero"}
}

// evluate infix integer operations
func evaluateIntegerInfixExpression(operator string, left obj.Object, rigth obj.Object) obj.Object {
	leftVal := left.(*obj.Number).Value
	rigthVal := rigth.(*obj.Number).Value

	switch operator {
	case "+":
		return &obj.Number{Value: leftVal + rigthVal}
	case "-":
		return &obj.Number{Value: leftVal - rigthVal}
	case "*":
		return &obj.Number{Value: leftVal * rigthVal}
	case "/":
		return &obj.Number{Value: leftVal / rigthVal}
	case "%":
		return &obj.Number{Value: leftVal % rigthVal}
	case "+=":
		left.(*obj.Number).Value += rigthVal
		return left

	case "-=":
		left.(*obj.Number).Value -= rigthVal
		return left

	case "/=":
		left.(*obj.Number).Value /= rigthVal
		return left

	case "*=":
		left.(*obj.Number).Value *= rigthVal
		return left

	case ">":
		return toBooleanObject(leftVal > rigthVal)
	case "<":
		return toBooleanObject(leftVal < rigthVal)
	case "==":
		return toBooleanObject(leftVal == rigthVal)
	case "!=":
		return toBooleanObject(leftVal != rigthVal)
	case ">=":
		return toBooleanObject(leftVal >= rigthVal)
	case "<=":
		return toBooleanObject(leftVal <= rigthVal)

	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

// check that the character after - is a number
func evaluateMinusOperatorExpression(rigth obj.Object) obj.Object {
	if right, isNumber := rigth.(*obj.Number); isNumber {
		right.Value = -right.Value
		return right
	}

	return unknownPrefixOperator("-", obj.Types[rigth.Type()])
}

// evaluate prefix expressions
func evaluatePrefixExpression(operator string, rigth obj.Object) obj.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(rigth)

	case "-":
		return evaluateMinusOperatorExpression(rigth)

	default:
		return unknownPrefixOperator(operator, obj.Types[rigth.Type()])
	}
}

// generates a new error instance
func newError(message string) *obj.Error {
	return &obj.Error{Message: message}
}

// recibe an expression and return the corresponding object type
func toBooleanObject(val bool) obj.Object {
	if val {
		return obj.SingletonTRUE
	}
	return obj.SingletonFALSE
}

// utils functions to return errors
func typeMismatchError(left, operator, rigth string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("Discrepancia de tipos: %s %s %s", left, operator, rigth),
	}
}

func unknownPrefixOperator(operator, rigth string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("Operador desconocido: %s%s", operator, rigth),
	}
}

func unknownInfixOperator(left, operator, rigth string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("Operador desconocido: %s %s %s", left, operator, rigth),
	}
}

func unknownIdentifier(identifier string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("Identificador no encontrado: %s", identifier),
	}
}

func notAFunction(identifier string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("No es una funcion: %s", identifier),
	}
}

func notAList(identifier string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("No es una lista: %s", identifier),
	}
}

func notAVariable(identifier string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("No es una variable: %s", identifier),
	}
}

func noSuchMethod(method, datastruct string) *obj.Error {
	return &obj.Error{
		Message: fmt.Sprintf("%s no tiene un metodo %s", datastruct, method),
	}
}
