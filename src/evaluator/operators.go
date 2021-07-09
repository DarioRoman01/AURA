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

	case left.Type() == obj.FLOATING && right.Type() == obj.FLOATING:
		return evaluateFloatInfixExpression(operator, left, right)

	case left.Type() == obj.FLOATING && right.Type() == obj.INTEGERS:
		return evaluateLeftFloatInfixExp(operator, left, right)

	case left.Type() == obj.INTEGERS && right.Type() == obj.FLOATING:
		return evaluateRigthFloatInfixExp(operator, left, right)

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

// evaluate bool infix expressions
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

func evaluateLeftFloatInfixExp(operator string, left obj.Object, rigth obj.Object) obj.Object {
	leftVal := left.(*obj.Float).Value
	rigthVal := rigth.(*obj.Number).Value

	switch operator {
	case "+":
		return &obj.Float{Value: leftVal + float64(rigthVal)}
	case "-":
		return &obj.Float{Value: leftVal - float64(rigthVal)}
	case "*":
		return &obj.Float{Value: leftVal * float64(rigthVal)}
	case "/":
		return &obj.Float{Value: leftVal / float64(rigthVal)}
	case "+=":
		left.(*obj.Float).Value += float64(rigthVal)
		return left

	case "-=":
		left.(*obj.Float).Value -= float64(rigthVal)
		return left

	case "/=":
		left.(*obj.Float).Value /= float64(rigthVal)
		return left

	case "*=":
		left.(*obj.Float).Value *= float64(rigthVal)
		return left

	case ">":
		return toBooleanObject(leftVal > float64(rigthVal))
	case "<":
		return toBooleanObject(leftVal < float64(rigthVal))
	case "==":
		return toBooleanObject(leftVal == float64(rigthVal))
	case "!=":
		return toBooleanObject(leftVal != float64(rigthVal))
	case ">=":
		return toBooleanObject(leftVal >= float64(rigthVal))
	case "<=":
		return toBooleanObject(leftVal <= float64(rigthVal))

	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

func evaluateRigthFloatInfixExp(operator string, left obj.Object, rigth obj.Object) obj.Object {
	leftVal := left.(*obj.Number).Value
	rigthVal := rigth.(*obj.Float).Value

	switch operator {
	case "+":
		return &obj.Float{Value: float64(leftVal) + rigthVal}
	case "-":
		return &obj.Float{Value: float64(leftVal) - rigthVal}
	case "*":
		return &obj.Float{Value: float64(leftVal) * rigthVal}
	case "/":
		return &obj.Float{Value: float64(leftVal) / rigthVal}
	case ">":
		return toBooleanObject(float64(leftVal) > rigthVal)
	case "<":
		return toBooleanObject(float64(leftVal) < rigthVal)
	case "==":
		return toBooleanObject(float64(leftVal) == rigthVal)
	case "!=":
		return toBooleanObject(float64(leftVal) != rigthVal)
	case ">=":
		return toBooleanObject(float64(leftVal) >= rigthVal)
	case "<=":
		return toBooleanObject(float64(leftVal) <= rigthVal)

	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

// evaluate string infix expressions
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
	case ">":
		return toBooleanObject(leftVal > rigthVal)
	case "<":
		return toBooleanObject(leftVal < rigthVal)
	default:
		return unknownInfixOperator(
			obj.Types[left.Type()],
			operator,
			obj.Types[rigth.Type()],
		)
	}
}

// evaluate suffix expressions
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

func evaluateFloatInfixExpression(operator string, left, rigth obj.Object) obj.Object {
	leftVal := left.(*obj.Float).Value
	rigthVal := rigth.(*obj.Float).Value
	switch operator {
	case "+":
		return &obj.Float{Value: leftVal + rigthVal}
	case "-":
		return &obj.Float{Value: leftVal - rigthVal}
	case "*":
		return &obj.Float{Value: leftVal * rigthVal}
	case "/":
		return &obj.Float{Value: leftVal / rigthVal}
	case "+=":
		left.(*obj.Float).Value += rigthVal
		return left

	case "-=":
		left.(*obj.Float).Value -= rigthVal
		return left

	case "/=":
		left.(*obj.Float).Value /= rigthVal
		return left

	case "*=":
		left.(*obj.Float).Value *= rigthVal
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

// evluate infix integer expressions
func evaluateIntegerInfixExpression(operator string, left, rigth obj.Object) obj.Object {
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

// check that the character after - is a number and apply the operator
func evaluateMinusOperatorExpression(rigth obj.Object) obj.Object {
	switch num := rigth.(type) {
	case *obj.Number:
		num.Value = -num.Value
		return num

	case *obj.Float:
		num.Value = -num.Value
		return num

	default:
		return unknownPrefixOperator("-", obj.Types[rigth.Type()])
	}
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
