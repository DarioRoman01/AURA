package evaluator

import (
	"aura/src/ast"
	b "aura/src/builtins"
	obj "aura/src/object"
	"math"
	"unicode/utf8"
)

// evlauate given nodes of the ast
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

	case *ast.ClassStatement:
		class := obj.NewClass(node.Name, node.Params, node.Methods)
		env.SetItem(class.Name.Value, class)
		return obj.SingletonNUll

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
		return obj.SingletonNUll

	case *ast.Identifier:
		return evaluateIdentifier(node, env)

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
		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}

	default:
		return obj.SingletonNUll
	}
}

// generates a new function object
func applyFunction(fn obj.Object, args []obj.Object) obj.Object {
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
			evaluated = Evaluate(forLoop.Body, env)
			if returnVal, isReturn := evaluated.(*obj.Return); isReturn {
				// we break the loop because we have a return statement
				return returnVal
			}

			// we update the variable in the expression
			env.Store[val] = iter.Current
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
		iter := obj.NewIterator(list.Values[0], list.Values)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	if str, isStr := Evaluate(rangeExpress.Range, env).(*obj.String); isStr {
		// if the iter is a string we make a iterable with all the string characters
		list := makeStringList(str.Value)
		iter := obj.NewIterator(list[0], list)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	return notIterable(rangeExpress.Range.Str())
}

// extends the class enviroment with the methods and constructor params
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
		evaluated := Evaluate(call.Field, class.Env)
		if _, isErr := evaluated.(*obj.Error); isErr {
			return evaluated
		}

		return evaluated
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

// evaluate an import statement
func evaluateImportStatement(importStmt *ast.ImportStatement, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(importStmt.Path, env)
	if str, isStr := evaluated.(*obj.String); isStr {
		fileEnv, err := importEnv(str.Value)
		if err != nil {
			return err
		}

		env.Outer = fileEnv
		return obj.SingletonNUll
	}

	return newError("La direccion para import un archivo debe ser un string")
}

// evluate all the statements in a bock expression
func evaluateBLockStaments(block *ast.Block, env *obj.Enviroment) obj.Object {
	var result obj.Object
	for _, statement := range block.Staments {
		result = Evaluate(statement, env)
		if result != nil && result.Type() == obj.RETURNTYPE || result.Type() == obj.ERROR {
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
	value, exists := env.GetItem(node.Value)
	if !exists {
		// check if the identifier is a builtin function
		builtint, exists := b.BUILTINS[node.Value]
		if !exists {
			// the identifier doest not exists
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

func evaluateFunction(function *ast.Function, env *obj.Enviroment) obj.Object {
	if function.Name != nil {
		def := obj.NewDef(function.Body, env, function.Parameters...)
		env.SetItem(function.Name.Value, def)
		return obj.SingletonNUll
	}

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
		if returnVal, isReturn := evaluated.(*obj.Return); isReturn {
			return returnVal
		}

		condition = Evaluate(whileExpression.Condition, env)
	}

	return obj.SingletonNUll
}

// Evaluate a call to a datastructure like:
//		array[0];
func evaluateCallList(call *ast.CallList, env *obj.Enviroment) obj.Object {
	evaluated := Evaluate(call.ListIdent, env)
	switch object := evaluated.(type) {

	case *obj.List:
		evaluated := Evaluate(call.Index, env)
		num, isNumber := evaluated.(*obj.Number)
		if !isNumber {
			return &obj.Error{Message: "El indice debe ser un entero"}
		}

		length := len(object.Values)
		if num.Value >= length {
			return indexOutOfRange(num.Value, length)
		}

		if num.Value < 0 {
			if int(math.Abs(float64(num.Value))) > length {

				return indexOutOfRange(num.Value, length)
			}

			return object.Values[length+num.Value]
		}

		return object.Values[num.Value]

	case *obj.Map:
		evaluated := Evaluate(call.Index, env)
		return object.Get(evaluated.Inspect())

	case *obj.String:
		evaluated := Evaluate(call.Index, env)
		num, isNumber := evaluated.(*obj.Number)
		if !isNumber {
			return &obj.Error{Message: "El indice debe ser un entero"}
		}

		strLen := utf8.RuneCountInString(object.Value)
		if num.Value >= strLen {
			return indexOutOfRange(num.Value, strLen)
		}

		if num.Value < 0 {
			if int(math.Abs(float64(num.Value))) > strLen {

				return indexOutOfRange(num.Value, strLen)
			}

			return &obj.String{Value: string(object.Value[strLen+num.Value])}
		}

		return &obj.String{Value: string(object.Value[num.Value])}

	default:
		return cannotBeIndexed(obj.Types[evaluated.Type()])
	}
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
