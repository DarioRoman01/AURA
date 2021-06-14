package evaluator

import (
	"aura/src/ast"
	b "aura/src/builtins"
	obj "aura/src/object"
	"fmt"
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

	return newError("Expression por invalida")
}

func evaluateRange(rangeExpress *ast.RangeExpression, env *obj.Enviroment) obj.Object {
	if list, isList := Evaluate(rangeExpress.Range, env).(*obj.List); isList {
		val, isVar := rangeExpress.Variable.(*ast.Identifier)
		if !isVar {
			return notAVariable(rangeExpress.Variable.Str())
		}
		iter := obj.NewIterator(list.Values[0], list.Values)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	if str, isStr := Evaluate(rangeExpress.Range, env).(*obj.String); isStr {
		val, isVar := rangeExpress.Variable.(*ast.Identifier)
		if !isVar {
			return notAVariable(rangeExpress.Variable.Str())
		}

		list := makeStringList(str.Value)
		iter := obj.NewIterator(list[0], list)
		env.SetItem(val.Value, iter.List[0])
		return iter
	}

	return notIterable(rangeExpress.Range.Str())
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

// recibe an expression and return the corresponding object type
func toBooleanObject(val bool) obj.Object {
	if val {
		return obj.SingletonTRUE
	}
	return obj.SingletonFALSE
}
