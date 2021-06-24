package evaluator

import (
	obj "aura/src/object"
	"fmt"
)

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

	case "+=":
		left.(*obj.String).Value += rigthVal
		return left
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

	return newError(fmt.Sprintf("el operador %s solo puede ser aplicado en numeros", operator))
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
