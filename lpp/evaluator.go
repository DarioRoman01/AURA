package lpp

import "fmt"

var singletonTRUE = &Bool{Value: true}
var singletonFALSE = &Bool{Value: false}
var SingletonNUll = &Null{}

func Evaluate(baseNode ASTNode, env *Enviroment) Object {
	switch node := baseNode.(type) {

	case Program:
		return evaluateProgram(node, env)

	case *ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression, env)

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

	case *If:
		return evaluateIfExpression(node, env)

	case *ReturnStament:
		CheckIsNotNil(node.ReturnValue)
		value := Evaluate(node.ReturnValue, env)
		CheckIsNotNil(value)
		return &Return{Value: value}

	case *LetStatement:
		CheckIsNotNil(node.Value)
		value := Evaluate(node.Value, env)
		CheckIsNotNil(node.Name)
		env.SetItem(node.Name.value, value)
		return nil

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

	default:
		return SingletonNUll
	}
}

func applyFunction(fn Object, args []Object) Object {
	if _, isFn := fn.(*Def); !isFn {
		return newError(notAFunction(types[fn.Type()]))
	}

	function := fn.(*Def)
	extendedEnviron := extendFunctionEnviroment(function, args)
	evaluated := Evaluate(function.Body, extendedEnviron)
	CheckIsNotNil(evaluated)
	return unwrapReturnValue(evaluated)
}

func unwrapReturnValue(obj Object) Object {
	if _, isReturn := obj.(*Return); isReturn {
		newObj := obj.(*Return)
		return newObj.Value
	}
	return obj
}

func extendFunctionEnviroment(fn *Def, args []Object) *Enviroment {
	env := NewEnviroment(fn.Env)
	for idx, param := range fn.Parameters {
		env.SetItem(param.value, args[idx])
	}

	return env
}

func CheckIsNotNil(val interface{}) {
	if val == nil {
		panic("expression or stament cannot be nil :(")
	}
}

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

func evaluateExpression(expressions []Expression, env *Enviroment) []Object {
	var result []Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, env)
		CheckIsNotNil(evaluated)
		result = append(result, evaluated)
	}

	return result
}

func evaluateIdentifier(node *Identifier, env *Enviroment) Object {
	value, exists := env.GetItem(node.value)
	if !exists {
		return newError(unknownIdentifier(node.value))
	}

	return value
}

func evaluateProgram(program Program, env *Enviroment) Object {
	var result Object
	for _, statement := range program.Staments {
		result = Evaluate(statement, env)

		if _, isReturn := result.(*Return); isReturn {
			return result.(*Return).Value
		} else if _, isError := result.(*Error); isError {
			return result
		}
	}

	return result
}

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

func evaluateInfixExpression(operator string, left Object, right Object) Object {
	switch {

	case left.Type() == INTEGERS && right.Type() == INTEGERS:
		return evaluateIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return toBooleanObject(left == right)

	case operator == "!=":
		return toBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError(typeMismatchError(
			types[left.Type()],
			operator,
			types[right.Type()],
		))

	default:
		return newError(unknownInfixOperator(
			types[left.Type()],
			operator,
			types[right.Type()],
		))
	}

}

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
	case ">":
		return toBooleanObject(leftVal > rigthVal)
	case "<":
		return toBooleanObject(leftVal < rigthVal)
	case "==":
		return toBooleanObject(leftVal == rigthVal)
	case "!=":
		return toBooleanObject(leftVal != rigthVal)
	default:
		return newError(unknownInfixOperator(
			types[left.Type()],
			operator,
			types[rigth.Type()],
		))
	}
}

func evaluateMinusOperatorExpression(rigth Object) Object {
	if _, isNumber := rigth.(*Number); !isNumber {
		return newError(unknownPrefixOperator("-", types[rigth.Type()]))
	}

	right := rigth.(*Number)
	right.Value = -right.Value
	return right
}

func evaluatePrefixExpression(operator string, rigth Object) Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(rigth)

	case "-":
		return evaluateMinusOperatorExpression(rigth)

	default:
		return newError(unknownPrefixOperator(operator, types[rigth.Type()]))
	}
}

func newError(message string) *Error {
	return &Error{Message: message}
}

func toBooleanObject(val bool) Object {
	if val {
		return singletonTRUE
	}
	return singletonFALSE
}

func typeMismatchError(left, operator, rigth string) string {
	return fmt.Sprintf("Discrepancia de tipos: %s %s %s", left, operator, rigth)
}

func unknownPrefixOperator(operator, rigth string) string {
	return fmt.Sprintf("Operador desconocido: %s%s", operator, rigth)
}

func unknownInfixOperator(left, operator, rigth string) string {
	return fmt.Sprintf("Operador desconocido: %s %s %s", left, operator, rigth)
}

func unknownIdentifier(identifier string) string {
	return fmt.Sprintf("Identificador no encontrado: %s", identifier)
}

func notAFunction(identifier string) string {
	return fmt.Sprintf("No es una funcion: %s", identifier)
}
