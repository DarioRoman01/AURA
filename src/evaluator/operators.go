package evaluator

import (
	"aura/src/ast"
	obj "aura/src/object"
	"fmt"
	"reflect"
)

// evluate infix expressions between objects
func evaluateInfixExpression(operator string, left obj.Object, right obj.Object, env *obj.Enviroment, leftNode ast.Expression) obj.Object {
	switch {

	case left.Type() == obj.INTEGERS && right.Type() == obj.INTEGERS:
		return evaluateIntegerInfixExpression(operator, left, right)

	case left.Type() == obj.FLOATING && right.Type() == obj.FLOATING:
		return evaluateFloatInfixExpression(operator, left, right)

	case left.Type() == obj.FLOATING && right.Type() == obj.INTEGERS:
		return evaluateLeftFloatInfixExp(operator, left, right)

	case left.Type() == obj.INTEGERS && right.Type() == obj.FLOATING:
		return evaluateRigthFloatInfixExp(operator, left, right, env, leftNode)

	case left.Type() == obj.STRINGTYPE && right.Type() == obj.STRINGTYPE:
		return evaluateStringInfixExpression(operator, left, right)

	case left.Type() == obj.BOOLEAN && right.Type() == obj.BOOLEAN:
		return evaluateBoolInfixExpression(operator, left.(*obj.Bool), right.(*obj.Bool))

	case operator == "==":
		return toBooleanObject(reflect.DeepEqual(left, right))

	case operator == "!=":
		return toBooleanObject(!reflect.DeepEqual(left, right))

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
		return obj.NewFloat(leftVal + float64(rigthVal))
	case "-":
		return obj.NewFloat(leftVal - float64(rigthVal))
	case "*":
		return obj.NewFloat(leftVal * float64(rigthVal))
	case "/":
		if leftVal == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
		return obj.NewFloat(leftVal / float64(rigthVal))

	case "+=":
		left.(*obj.Float).Value += float64(rigthVal)
		return left

	case "-=":
		left.(*obj.Float).Value -= float64(rigthVal)
		return left

	case "/=":
		if left.(*obj.Float).Value == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
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

func isIdentifier(exp ast.Expression) (*ast.Identifier, *obj.Error) {
	ident, isIdentifier := exp.(*ast.Identifier)
	if !isIdentifier {
		return nil, newError("no se puede user el /= sin una variable a la derecha si se usa un flotante a la izquierda")
	}
	return ident, nil
}

func evaluateRigthFloatInfixExp(operator string, left obj.Object, rigth obj.Object, env *obj.Enviroment, leftNode ast.Expression) obj.Object {
	leftVal := left.(*obj.Number).Value
	rigthVal := rigth.(*obj.Float).Value
	var variable string

	if operator == "+=" || operator == "-=" || operator == "/=" || operator == "*=" {
		ident, err := isIdentifier(leftNode)
		if err != nil {
			return err
		}

		variable = ident.Value
	}

	switch operator {
	case "+":
		return obj.NewFloat(float64(leftVal) + rigthVal)
	case "-":
		return obj.NewFloat(float64(leftVal) - rigthVal)
	case "*":
		return obj.NewFloat(float64(leftVal) * rigthVal)
	case "/":
		if leftVal == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
		return obj.NewFloat(float64(leftVal) / rigthVal)

	case "+=":
		env.SetItem(variable, obj.NewFloat(float64(leftVal)+rigthVal))
		obj, _ := env.GetItem(variable)
		return obj

	case "-=":
		env.SetItem(variable, obj.NewFloat(float64(leftVal)-rigthVal))
		obj, _ := env.GetItem(variable)
		return obj

	case "/=":
		if leftVal == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
		env.SetItem(variable, obj.NewFloat(float64(leftVal)/rigthVal))
		obj, _ := env.GetItem(variable)
		return obj

	case "*=":
		env.SetItem(variable, obj.NewFloat(float64(leftVal)*rigthVal))
		obj, _ := env.GetItem(variable)
		return obj

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
		return obj.NewFloat(leftVal + rigthVal)
	case "-":
		return obj.NewFloat(leftVal - rigthVal)
	case "*":
		return obj.NewFloat(leftVal * rigthVal)
	case "/":
		if leftVal == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
		return obj.NewFloat(leftVal / rigthVal)
	case "+=":
		left.(*obj.Float).Value += rigthVal
		return left

	case "-=":
		left.(*obj.Float).Value -= rigthVal
		return left

	case "/=":
		if left.(*obj.Float).Value == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}

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
		if leftVal == 0 && rigthVal == 0 {
			return newError("division entre 0 ")
		}
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
		if left.(*obj.Number).Value == 0 && rigthVal == 0 {
			return divisionByZeroError()
		}
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
