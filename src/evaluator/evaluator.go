package evaluator

import (
	"aura/src/ast"
	b "aura/src/builtins"
	obj "aura/src/object"
	"unicode/utf8"
)

// evlauate given nodes of the ast
func Evaluate(baseNode ast.ASTNode, env *obj.Enviroment) obj.Object {
	switch node := baseNode.(type) {

	case *ast.Program:
		return evaluateProgram(node, env)

	case *ast.ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression, env)

	case *ast.Array:
		return evaluateArray(node, env)

	case *ast.Integer:
		CheckIsNotNil(node.Value)
		return &obj.Number{Value: *node.Value}

	case *ast.FloatExp:
		CheckIsNotNil(node.Value)
		return &obj.Float{Value: node.Value}

	case *ast.Boolean:
		CheckIsNotNil(node.Value)
		return toBooleanObject(*node.Value)

	case *ast.Prefix:
		CheckIsNotNil(node.Rigth)
		rigth := Evaluate(node.Rigth, env)
		CheckIsNotNil(rigth)
		return evaluatePrefixExpression(node.Operator, rigth)

	case *ast.Infix:
		CheckIsNotNil(node.Left)
		CheckIsNotNil(node.Rigth)
		left := Evaluate(node.Left, env)
		rigth := Evaluate(node.Rigth, env)
		CheckIsNotNil(left)
		CheckIsNotNil(rigth)
		return evaluateInfixExpression(node.Operator, left, rigth, env, node.Left)

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

	case *ast.ClassStatement:
		class := obj.NewClass(node.Name, node.Params, node.Methods)
		env.SetItem(class.Name.Value, class)
		return obj.SingletonNUll

	case *ast.TryExp:
		CheckIsNotNil(node.Try)
		CheckIsNotNil(node.Catch)
		return evaluateTryExcept(node, env)

	case *ast.ArrowFunc:
		CheckIsNotNil(node.Body)
		return obj.NewDef(node.Body, env, node.Params...)

	case *ast.ClassCall:
		CheckIsNotNil(node.Arguments)
		return evaluateClassCall(node, env)

	case *ast.ClassFieldCall:
		CheckIsNotNil(node.Class)
		CheckIsNotNil(node.Field)
		return evaluateClassFieldCall(node, env)

	case *ast.LetStatement:
		CheckIsNotNil(node.Value)
		value := Evaluate(node.Value, env)
		CheckIsNotNil(node.Name)
		env.SetItem(node.Name.Value, value)
		return obj.SingletonNUll

	case *ast.AssigmentExp:
		CheckIsNotNil(node.Val)
		value := Evaluate(node.Val, env)
		CheckIsNotNil(node.Name)
		env.SetItem(node.Name.Value, value)
		return value

	case *ast.Identifier:
		return evaluateIdentifier(node, env)

	case *ast.BreakStatement:
		return obj.SingleTonBreak

	case *ast.ContinueStatement:
		return obj.SingletonContinue

	case *ast.ThorwExpression:
		CheckIsNotNil(node.Message)
		return newError(node.Message.Str())

	case *ast.TernaryIf:
		CheckIsNotNil(node.Condition)
		CheckIsNotNil(node.Consequence)
		CheckIsNotNil(node.Alternative)
		return evaluateTernaryIf(node, env)

	case *ast.Function:
		CheckIsNotNil(node.Body)
		return evaluateFunction(node, env)

	case *ast.ImportStatement:
		CheckIsNotNil(node.Path)
		return evaluateImportStatement(node, env)

	case *ast.Call:
		function := Evaluate(node.Function, env)
		CheckIsNotNil(node.Arguments)
		args := evaluateExpression(node.Arguments, env)
		CheckIsNotNil(function)
		return applyFunction(function, args...)

	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}

	default:
		return obj.SingletonNUll
	}
}

// generates a new function object
func applyFunction(fn obj.Object, args ...obj.Object) obj.Object {
	switch function := fn.(type) {
	case *obj.Def:
		extendedEnviron := extendFunctionEnviroment(function, args)
		evaluated := Evaluate(function.Body, extendedEnviron)
		CheckIsNotNil(evaluated)
		return unwrapReturnValue(evaluated)

	case *obj.Builtin:
		return function.Fn(args...)

	default:
		return notAFunction(obj.Types[fn.Type()])
	}
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

// Evaluate a variable reassigment
func evaluateVarReassigment(variable *ast.Identifier, newVal ast.Expression, env *obj.Enviroment) obj.Object {
	_, exists := env.GetItem(variable.Value)
	if !exists {
		// the veriable does not exists
		return unknownIdentifier(variable.Value)
	}

	env.Store[variable.Value] = Evaluate(newVal, env)
	return obj.SingletonNUll
}

// Evaluate a forloop expression
func evaluateFor(forLoop *ast.For, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(forLoop.Condition, env)
	if iter, isIter := evaluated.(*obj.Iterator); isIter {

		// this does not fail because if not a variable the error will be handle by
		// the evaluate iter function
		val := forLoop.Condition.(*ast.RangeExpression).Variable.(*ast.Identifier).Value
		for iter.Next() != nil {
			evaluated = Evaluate(forLoop.Body, iter.Env)
			switch node := evaluated.(type) {
			case *obj.Return:
				return node

			case *obj.BreakObj:
				return obj.SingletonNUll

			case *obj.ContinueObj:
				iter.Env.SetItem(val, iter.Current)
				continue

			default:
				iter.Env.SetItem(val, iter.Current)
			}
		}

		return obj.SingletonNUll
	}

	if err, isError := evaluated.(*obj.Error); isError {
		return err
	}

	return newError("Expression por invalida")
}

// evaluate an iter expression like:
//		for(i in range(10)):
func evaluateRange(rangeExpress *ast.RangeExpression, env *obj.Enviroment) obj.Object {
	val, isVar := rangeExpress.Variable.(*ast.Identifier)
	if !isVar {
		return notAVariable(rangeExpress.Variable.Str())
	}

	if list, isList := Evaluate(rangeExpress.Range, env).(*obj.List); isList {
		// if the iter is a list we make a iterable with the list
		iter := obj.NewIterator(list.Values[0], list.Values, obj.NewEnviroment(env))
		iter.Env.SetItem(val.Value, iter.Current)
		return iter
	}

	if str, isStr := Evaluate(rangeExpress.Range, env).(*obj.String); isStr {
		// if the iter is a string we make a iterable with all the string characters
		list := makeStringList(str.Value)
		iter := obj.NewIterator(list[0], list, obj.NewEnviroment(env))
		iter.Env.SetItem(val.Value, iter.Current)
		return iter
	}

	return notIterable(rangeExpress.Range.Str())
}

// extends the class enviroment with the methods and constructor arguments
func extendClassEnviroment(class *obj.Class, args []obj.Object, methods []*ast.ClassMethodExp, env *obj.Enviroment) *obj.Enviroment {
	classEnv := obj.NewEnviroment(env)
	for idx, param := range class.Params {
		classEnv.SetItem(param.Value, args[idx])
	}

	for _, method := range methods {
		classEnv.SetItem(method.Name.Value, obj.NewDef(method.Body, classEnv, method.Params...))
	}

	return classEnv
}

// evaluate a call to a new class
func evaluateClassCall(call *ast.ClassCall, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.Class, env)
	if _, isErr := evaluated.(*obj.Error); isErr {
		return evaluated
	}

	if class, isClass := evaluated.(*obj.Class); isClass {
		args := evaluateExpression(call.Arguments, env)
		classEnv := extendClassEnviroment(class, args, class.Methods, env)
		classInstance := obj.NewClassInstance(class.Name.Value, classEnv)
		return classInstance
	}

	return notAClass(call.Class.Value)
}

// evaluate a call to a class instance field or method
func evaluateClassFieldCall(call *ast.ClassFieldCall, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.Class, env)
	if _, isErr := evaluated.(*obj.Error); isErr {
		return evaluated
	}

	if class, isClass := evaluated.(*obj.ClassInstance); isClass {
		return Evaluate(call.Field, class.Env)
	}

	return notAClass(evaluated.Inspect())
}

// evaluate a class field reassigment
func evaluateFieldReassigment(call *ast.ClassFieldCall, class *obj.ClassInstance, newVal ast.Expression) obj.Object {
	ident, isIdent := call.Field.(*ast.Identifier)
	if !isIdent {
		return newError("Una funcion no puede ser reasignada")
	}

	_, exists := class.Env.GetItem(ident.Value)
	if !exists {
		return noSuchField(class.Name, ident.Value)
	}

	evaluated := Evaluate(newVal, class.Env)
	class.Env.SetItem(ident.Value, evaluated)
	return obj.SingletonNUll
}

// check that the given value is not nil
func CheckIsNotNil(val interface{}) {
	if val == nil {
		panic("Error de evaluacion! Se esperaba una expression pero se obtuvo nulo!")
	}
}

func evaluateTernaryIf(ternary *ast.TernaryIf, env *obj.Enviroment) obj.Object {
	condition := Evaluate(ternary.Condition, env)
	if isTruthy(condition) {
		return Evaluate(ternary.Consequence, env)
	} else {
		return Evaluate(ternary.Alternative, env)
	}
}

// evaluate an import statement
func evaluateImportStatement(importStmt *ast.ImportStatement, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(importStmt.Path, env)
	if str, isStr := evaluated.(*obj.String); isStr {
		fileEnv, err := importEnv(str.Value)
		if err != nil {
			return err
		}

		env.SetOuter(fileEnv)
		return obj.SingletonNUll
	}

	return newError("La direccion para importar un archivo debe ser un string")
}

// evluate all the statements in a bock expression
func evaluateBLockStaments(block *ast.Block, env *obj.Enviroment) obj.Object {
	var result obj.Object
	for _, statement := range block.Staments {
		result = Evaluate(statement, env)
		if result != nil && result.Type() == obj.RETURNTYPE || result.Type() == obj.ERROR {
			return result
		}

		if result != nil && result.Type() == obj.BREAK || result.Type() == obj.CONTINUE {
			return result
		}
	}

	return result
}

// evaluate a slice of expressions
func evaluateExpression(expressions []ast.Expression, env *obj.Enviroment) []obj.Object {
	var result []obj.Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, env)
		CheckIsNotNil(evaluated)
		result = append(result, evaluated)
	}

	return result
}

// check if given identifier exists in the enviroment
func evaluateIdentifier(node *ast.Identifier, env *obj.Enviroment) obj.Object {
	object, exists := env.GetItem(node.Value)
	if !exists {
		// check if the identifier is a builtin function
		builtint, exists := b.BUILTINS[node.Value]
		if !exists {
			// the identifier doest not exists
			return unknownIdentifier(node.Value)
		}

		return builtint
	}

	return object
}

// evaluate program node
func evaluateProgram(program *ast.Program, env *obj.Enviroment) obj.Object {
	var result obj.Object
	for _, statement := range program.Staments {
		result = Evaluate(statement, env)

		if returnObj, isReturn := result.(*obj.Return); isReturn {
			return returnObj.Value
		}

		if err, isError := result.(*obj.Error); isError {
			return err
		}
	}

	return result
}

// evaluate a function a expression
func evaluateFunction(function *ast.Function, env *obj.Enviroment) obj.Object {
	if function.Name != nil {
		// if the name is not nil we have a named function like func main() {}
		def := obj.NewDef(function.Body, env, function.Parameters...)
		env.SetItem(function.Name.Value, def)
		return obj.SingletonNUll
	}

	// we have an anonimous function like x := func(a, b) {}
	return obj.NewDef(function.Body, env, function.Parameters...)
}

// evaluate a while looop expression
func evaluateWhileExpression(whileExpression *ast.While, env *obj.Enviroment) obj.Object {
	CheckIsNotNil(whileExpression.Condition)
	condition := Evaluate(whileExpression.Condition, env)
	CheckIsNotNil(condition)

	// we loop an update the condition until the condition is not trythy
	for isTruthy(condition) {
		evaluated := Evaluate(whileExpression.Body, env)
		switch node := evaluated.(type) {
		case *obj.Return:
			return node

		case *obj.BreakObj:
			return obj.SingletonNUll

		case *obj.ContinueObj:
			condition = Evaluate(whileExpression.Condition, env)
			continue

		default:
			condition = Evaluate(whileExpression.Condition, env)
		}
	}

	return obj.SingletonNUll
}

// Evaluate a call to a datastructure like:
//		array[0];
func evaluateCallList(call *ast.CallList, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.ListIdent, env)
	switch object := evaluated.(type) {

	case *obj.List:
		return evaluateListCall(object, call, env)

	case *obj.Map:
		evaluated := Evaluate(call.Index, env)
		return object.Get(evaluated.Inspect())

	case *obj.String:
		return evaluateStringCall(object, call, env)

	default:
		return cannotBeIndexed(obj.Types[evaluated.Type()])
	}
}

func evaluateListCall(list *obj.List, call *ast.CallList, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.Index, env)
	num, isNumber := evaluated.(*obj.Number)
	if !isNumber {
		return newError("el indice debe ser un enetero")
	}

	length := len(list.Values)
	index, err := checkIndex(length, num.Value)
	if err != nil {
		return err
	}

	return list.Values[index]
}

func evaluateStringCall(str *obj.String, call *ast.CallList, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.Index, env)
	num, isNumber := evaluated.(*obj.Number)
	if !isNumber {
		return newError("El indice debe ser un entero")
	}

	strLen := utf8.RuneCountInString(str.Value)
	index, err := checkIndex(strLen, num.Value)
	if err != nil {
		return err
	}

	return &obj.String{Value: string(str.Value[index])}
}

// evaluate an if expression
func evaluateIfExpression(ifExpression *ast.If, env *obj.Enviroment) obj.Object {
	CheckIsNotNil(ifExpression.Condition)
	condition := Evaluate(ifExpression.Condition, env)
	CheckIsNotNil(condition)

	if isTruthy(condition) {
		// if the first condition is truthy we evaluate the consequence
		CheckIsNotNil(ifExpression.Consequence)
		return Evaluate(ifExpression.Consequence, env)

	} else if ifExpression.Alternative != nil {
		// if the condition is not truthy and the alternative
		// is not nil we evaluate the alternative
		return Evaluate(ifExpression.Alternative, env)
	}

	return obj.SingletonNUll
}

// evaluate a try excpet expression
func evaluateTryExcept(try *ast.TryExp, env *obj.Enviroment) obj.Object {
	eval := Evaluate(try.Try, env)
	if err, isErr := eval.(*obj.Error); isErr {
		newEnv := obj.NewEnviroment(env)
		newEnv.SetItem(try.Param.Value, err)
		return Evaluate(try.Catch, newEnv)
	}

	if returnVal, isReturn := eval.(*obj.Return); isReturn {
		if err, isErr := returnVal.Value.(*obj.Error); isErr {
			newEnv := obj.NewEnviroment(env)
			newEnv.SetItem(try.Param.Value, err)
			return Evaluate(try.Catch, newEnv)
		}

		return returnVal
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

// recibe an expression and return the corresponding object type
func toBooleanObject(val bool) obj.Object {
	if val {
		return obj.SingletonTRUE
	}
	return obj.SingletonFALSE
}
