package src

import (
	"fmt"
)

// use singleton patern with true false and null
var (
	singletonTRUE  = &Bool{Value: true}
	singletonFALSE = &Bool{Value: false}
	SingletonNUll  = &Null{} //  this null is for functions that dont return anything
	NullVAlue      = &Null{} // this null is the null value
)

// evlauate given nodes of an ast
func Evaluate(baseNode ASTNode, env *Enviroment) Object {
	switch node := baseNode.(type) {

	case Program:
		return evaluateProgram(node, env)

	case *ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression, env)

	case *Array:
		return evaluateArray(node, env)

	case *Integer:
		CheckIsNotNil(node.Value)
		return &Number{Value: *node.Value}

	case *Boolean:
		CheckIsNotNil(node.Value)
		return toBooleanObject(*node.Value)

	case *Prefix:
		CheckIsNotNil(node.Rigth)
		rigth := Evaluate(node.Rigth, env)
		CheckIsNotNil(rigth)
		return evaluatePrefixExpression(node.Operator, rigth)

	case *Infix:
		go CheckIsNotNil(node.Left)
		CheckIsNotNil(node.Rigth)
		left := Evaluate(node.Left, env)
		rigth := Evaluate(node.Rigth, env)
		go CheckIsNotNil(left)
		CheckIsNotNil(rigth)
		return evaluateInfixExpression(node.Operator, left, rigth)

	case *Block:
		return evaluateBLockStaments(node, env)

	case *While:
		return evaluateWhileExpression(node, env)

	case *If:
		return evaluateIfExpression(node, env)

	case *ReturnStament:
		CheckIsNotNil(node.ReturnValue)
		value := Evaluate(node.ReturnValue, env)
		CheckIsNotNil(value)
		return &Return{Value: value}

	case *For:
		CheckIsNotNil(node.Condition)
		CheckIsNotNil(node.Body)
		return evaluateFor(node, env)

	case *MethodExpression:
		CheckIsNotNil(node.Method)
		CheckIsNotNil(node.Obj)
		return evaluateMethod(node, env)

	case *RangeExpression:
		CheckIsNotNil(node.Range)
		CheckIsNotNil(node.Variable)
		return evaluateRange(node, env)

	case *CallList:
		return evaluateCallList(node, env)

	case *Reassignment:
		CheckIsNotNil(node.Identifier)
		CheckIsNotNil(node.NewVal)
		return evaluateReassigment(node, env)

	case *NullExpression:
		return NullVAlue

	case *MapExpression:
		CheckIsNotNil(node.Body)
		return evaluateMap(node, env)

	case *LetStatement:
		CheckIsNotNil(node.Value)
		value := Evaluate(node.Value, env)
		CheckIsNotNil(node.Name)
		env.SetItem(node.Name.value, value)
		return SingletonNUll

	case *Identifier:
		return evaluateIdentifier(node, env)

	case *Function:
		CheckIsNotNil(node.Body)
		return NewDef(node.Body, env, node.Parameters...)

	case *Call:
		function := Evaluate(node.Function, env)
		CheckIsNotNil(node.Arguments)
		args := evaluateExpression(node.Arguments, env)
		CheckIsNotNil(function)
		return applyFunction(function, args)

	case *StringLiteral:
		return &String{Value: node.Value}

	default:
		return SingletonNUll
	}
}

// generates a new function object
func applyFunction(fn Object, args []Object) Object {
	if function, isFn := fn.(*Def); isFn {
		extendedEnviron := extendFunctionEnviroment(function, args)
		evaluated := Evaluate(function.Body, extendedEnviron)
		CheckIsNotNil(evaluated)
		return unwrapReturnValue(evaluated)

	} else if builtin, isBuiltin := fn.(*Builtin); isBuiltin {
		return builtin.Fn(args...)
	}

	return notAFunction(types[fn.Type()])
}

// unwrap the return value of a function
func unwrapReturnValue(obj Object) Object {
	if obj, isReturn := obj.(*Return); isReturn {
		return obj.Value
	}
	return obj
}

// create a new enviroment when a function is called
func extendFunctionEnviroment(fn *Def, args []Object) *Enviroment {
	env := NewEnviroment(fn.Env)
	for idx, param := range fn.Parameters {
		env.SetItem(param.value, args[idx])
	}

	return env
}

func evaluateVarReassigment(variable *Identifier, newVal Expression, env *Enviroment) Object {
	_, exists := env.GetItem(variable.value)
	if !exists {
		return unknownIdentifier(variable.value)
	}

	env.store[variable.value] = Evaluate(newVal, env)
	return SingletonNUll
}

func evaluateListReassigment(call *CallList, list *List, newVal Expression, env *Enviroment) Object {
	index := Evaluate(call.Index, env)
	num, isNum := index.(*Number)
	if !isNum {
		return &Error{"El indice debe ser un numero"}
	}
	if num.Value >= len(list.Values) {
		return &Error{"Indice fuera de rango"}
	}

	list.Values[num.Value] = Evaluate(newVal, env)
	return SingletonNUll
}

func evaluateMapReassigment(hashMap *Map, key Object, value Object) Object {
	if err := hashMap.Get(string(hashMap.Serialize(key))); err != nil {
		hashMap.SetValues(key, value)
		return SingletonNUll
	}

	hashMap.UpdateKey(key, value)
	return SingletonNUll
}

func evaluateReassigment(reassigment *Reassignment, env *Enviroment) Object {
	// variable reassigment
	if variable, isVar := reassigment.Identifier.(*Identifier); isVar {
		return evaluateVarReassigment(variable, reassigment.NewVal, env)
	}

	if callList, isCall := reassigment.Identifier.(*CallList); isCall {
		evaluated := Evaluate(callList.ListIdent, env)
		// list reassigment
		if list, isList := evaluated.(*List); isList {
			return evaluateListReassigment(callList, list, reassigment.NewVal, env)
		}

		// map reassigment
		if hashMap, isMap := evaluated.(*Map); isMap {
			key := Evaluate(callList.Index, env)
			newVal := Evaluate(reassigment.NewVal, env)
			return evaluateMapReassigment(hashMap, key, newVal)
		}

		return notAList(evaluated.Inspect())
	}

	return notAVariable(reassigment.Identifier.TokenLiteral())
}

func evaluateListMethods(list *List, method *Method) Object {
	switch method.MethodType {
	case POP:
		return list.Pop()

	case APPEND:
		list.Add(method.Value)
		return SingletonNUll

	case REMOVE:
		index := method.Value.(*Number)
		return list.RemoveAt(index.Value)

	default:
		return noSuchMethod(method.Inspect(), "list")
	}
}

func evaluateMapMethods(hashMap *Map, method *Method) Object {
	switch method.MethodType {
	case CONTAIS:
		return &Bool{hashMap.Get(string(hashMap.Serialize(method.Value))) != NullVAlue}

	case VALUES:
		list := &List{[]Object{}}
		for _, val := range hashMap.store {
			list.Values = append(list.Values, val)
		}
		return list

	default:
		return noSuchMethod(method.Inspect(), "mapa")
	}
}

func evaluateMethod(method *MethodExpression, env *Enviroment) Object {
	evaluated := Evaluate(method.Obj, env)
	if list, isList := evaluated.(*List); isList {
		listMethod, isMethod := Evaluate(method.Method, env).(*Method)
		if !isMethod {
			return noSuchMethod(listMethod.Inspect(), "list")
		}

		return evaluateListMethods(list, listMethod)
	}

	if hashMap, isMap := evaluated.(*Map); isMap {
		mapMethod, isMethod := Evaluate(method.Method, env).(*Method)
		if !isMethod {
			return noSuchMethod(mapMethod.Inspect(), "mapa")
		}

		return evaluateMapMethods(hashMap, mapMethod)
	}

	return &Error{fmt.Sprintf("%s no tiene metodos", types[evaluated.Type()])}
}

func evaluateFor(forLoop *For, env *Enviroment) Object {
	evaluated := Evaluate(forLoop.Condition, env)
	if iter, isIter := evaluated.(*Iterator); isIter {
		val := forLoop.Condition.(*RangeExpression).Variable.(*Identifier).value
		for iter.Next() != nil {
			evaluated = Evaluate(forLoop.Body, env)
			if returnVal, isReturn := evaluated.(*Return); isReturn {
				return returnVal
			}

			env.store[val] = iter.current
		}
		return SingletonNUll
	}

	if err, isError := evaluated.(*Error); isError {
		return err
	}

	return newError("syntax error")
}

func evaluateMap(mapa *MapExpression, env *Enviroment) Object {
	mapObj := &Map{map[string]Object{}}
	for _, keyVal := range mapa.Body {
		key := Evaluate(keyVal.Key, env)
		val := Evaluate(keyVal.Value, env)
		if err := mapObj.SetValues(key, val); err != nil {
			return &Error{"No se permiten llaves duplicadas"}
		}
	}

	return mapObj
}

func evaluateRange(rangeExpress *RangeExpression, env *Enviroment) Object {
	if list, isList := Evaluate(rangeExpress.Range, env).(*List); isList {
		val, isVar := rangeExpress.Variable.(*Identifier)
		if !isVar {
			return newError("no es una variable")
		}
		iter := NewIterator(list.Values[0], list.Values)
		env.SetItem(val.value, iter.list[0])
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
func evaluateBLockStaments(block *Block, env *Enviroment) Object {
	var result Object = nil
	for _, statement := range block.Staments {
		result = Evaluate(statement, env)
		if result != nil && result.Type() == RETURNTYPE || result.Type() == ERROR {
			return result
		}
	}

	return result
}

// evaluate an slice of expressions
func evaluateExpression(expressions []Expression, env *Enviroment) []Object {
	var result []Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, env)
		CheckIsNotNil(evaluated)
		result = append(result, evaluated)
	}

	return result
}

func evaluateArray(arr *Array, env *Enviroment) Object {
	var list List

	for _, val := range arr.Values {
		list.Values = append(list.Values, Evaluate(val, env))
	}

	return &list
}

// check if given identifier exists in the enviroment
func evaluateIdentifier(node *Identifier, env *Enviroment) Object {
	value, exists := env.GetItem(node.value)
	if !exists {
		builtint, exists := BUILTINS[node.value]
		if !exists {
			return unknownIdentifier(node.value)
		}

		return builtint
	}

	return value
}

// evaluate program node
func evaluateProgram(program Program, env *Enviroment) Object {
	var result Object
	for _, statement := range program.Staments {
		result = Evaluate(statement, env)

		if returnObj, isReturn := result.(*Return); isReturn {
			return returnObj.Value
		} else if err, isError := result.(*Error); isError {
			return err
		}
	}

	return result
}

// change the bool value of the object
func evaluateBangOperatorExpression(rigth Object) Object {
	switch {
	case rigth == singletonTRUE:
		return singletonFALSE

	case rigth == singletonFALSE:
		return singletonTRUE

	case rigth == nil:
		return singletonTRUE

	default:
		return singletonFALSE
	}
}

func evaluateWhileExpression(whileExpression *While, env *Enviroment) Object {
	CheckIsNotNil(whileExpression.Condition)
	condition := Evaluate(whileExpression.Condition, env)

	CheckIsNotNil(condition)
	if !isTruthy(condition) {
		return SingletonNUll
	}
	Evaluate(whileExpression.Body, env)
	return Evaluate(whileExpression, env)
}

func evaluateCallList(call *CallList, env *Enviroment) Object {
	evaluated := Evaluate(call.ListIdent, env)
	if list, isList := evaluated.(*List); isList {
		evaluated := Evaluate(call.Index, env)
		num, isNumber := evaluated.(*Number)
		if !isNumber {
			return &Error{Message: "El indice debe ser un entero"}
		}

		if num.Value >= len(list.Values) {
			return &Error{Message: "Indice fuera de rango"}
		}

		return list.Values[num.Value]
	}

	if hashMap, isMap := evaluated.(*Map); isMap {
		evaluated := Evaluate(call.Index, env)
		return hashMap.Get(string(hashMap.Serialize(evaluated)))
	}

	return notAList(types[evaluated.Type()])
}

func evaluateIfExpression(ifExpression *If, env *Enviroment) Object {
	CheckIsNotNil(ifExpression.Condition)
	condition := Evaluate(ifExpression.Condition, env)

	CheckIsNotNil(condition)
	if isTruthy(condition) {
		CheckIsNotNil(ifExpression.Consequence)
		return Evaluate(ifExpression.Consequence, env)

	} else if ifExpression.Alternative != nil {
		return Evaluate(ifExpression.Alternative, env)
	}

	return SingletonNUll
}

// check that the current object is true or false
func isTruthy(obj Object) bool {
	switch {
	case obj == SingletonNUll:
		return false

	case obj == singletonTRUE:
		return true

	case obj == singletonFALSE:
		return false

	default:
		return true
	}
}

// evluate infix expressions between objects
func evaluateInfixExpression(operator string, left Object, right Object) Object {
	switch {

	case left.Type() == INTEGERS && right.Type() == INTEGERS:
		return evaluateIntegerInfixExpression(operator, left, right)

	case left.Type() == STRINGTYPE && right.Type() == STRINGTYPE:
		return evaluateStringInfixExpression(operator, left, right)

	case left.Type() == BOOLEAN && right.Type() == BOOLEAN:
		return evaluateBoolInfixExpression(operator, left.(*Bool), right.(*Bool))

	case operator == "==":
		return toBooleanObject(left == right)

	case operator == "!=":
		return toBooleanObject(left != right)

	case left.Type() != right.Type():
		return typeMismatchError(
			types[left.Type()],
			operator,
			types[right.Type()],
		)

	default:
		return unknownInfixOperator(
			types[left.Type()],
			operator,
			types[right.Type()],
		)
	}

}

func evaluateBoolInfixExpression(operator string, left *Bool, rigth *Bool) Object {
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
			types[left.Type()],
			operator,
			types[rigth.Type()],
		)
	}
}

func evaluateStringInfixExpression(operator string, left Object, rigth Object) Object {
	leftVal := left.(*String).Value
	rigthVal := rigth.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rigthVal}
	case "==":
		return toBooleanObject(leftVal == rigthVal)
	case "!=":
		return toBooleanObject(leftVal != rigthVal)
	default:
		return unknownInfixOperator(
			types[left.Type()],
			operator,
			types[rigth.Type()],
		)
	}
}

// evluate infix integer operations
func evaluateIntegerInfixExpression(operator string, left Object, rigth Object) Object {
	leftVal := left.(*Number).Value
	rigthVal := rigth.(*Number).Value

	switch operator {
	case "+":
		return &Number{Value: leftVal + rigthVal}
	case "-":
		return &Number{Value: leftVal - rigthVal}
	case "*":
		return &Number{Value: leftVal * rigthVal}
	case "/":
		return &Number{Value: leftVal / rigthVal}
	case "%":
		return &Number{Value: leftVal % rigthVal}
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
			types[left.Type()],
			operator,
			types[rigth.Type()],
		)
	}
}

// check that the character after - is a number
func evaluateMinusOperatorExpression(rigth Object) Object {
	if right, isNumber := rigth.(*Number); isNumber {
		right.Value = -right.Value
		return right
	}

	return unknownPrefixOperator("-", types[rigth.Type()])
}

// evaluate prefix expressions
func evaluatePrefixExpression(operator string, rigth Object) Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(rigth)

	case "-":
		return evaluateMinusOperatorExpression(rigth)

	default:
		return unknownPrefixOperator(operator, types[rigth.Type()])
	}
}

// generates a new error instance
func newError(message string) *Error {
	return &Error{Message: message}
}

// recibe an expression and return the corresponding object type
func toBooleanObject(val bool) Object {
	if val {
		return singletonTRUE
	}
	return singletonFALSE
}

// utils functions to return errors
func typeMismatchError(left, operator, rigth string) *Error {
	return &Error{
		Message: fmt.Sprintf("Discrepancia de tipos: %s %s %s", left, operator, rigth),
	}
}

func unknownPrefixOperator(operator, rigth string) *Error {
	return &Error{
		Message: fmt.Sprintf("Operador desconocido: %s%s", operator, rigth),
	}
}

func unknownInfixOperator(left, operator, rigth string) *Error {
	return &Error{
		Message: fmt.Sprintf("Operador desconocido: %s %s %s", left, operator, rigth),
	}
}

func unknownIdentifier(identifier string) *Error {
	return &Error{
		Message: fmt.Sprintf("Identificador no encontrado: %s", identifier),
	}
}

func notAFunction(identifier string) *Error {
	return &Error{
		Message: fmt.Sprintf("No es una funcion: %s", identifier),
	}
}

func notAList(identifier string) *Error {
	return &Error{
		Message: fmt.Sprintf("No es una lista: %s", identifier),
	}
}

func notAVariable(identifier string) *Error {
	return &Error{
		Message: fmt.Sprintf("No es una variable: %s", identifier),
	}
}

func noSuchMethod(method, datastruct string) *Error {
	return &Error{
		fmt.Sprintf("%s no tiene un metodo %s", datastruct, method),
	}
}
