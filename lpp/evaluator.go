package lpp

var singletonTRUE = &Bool{Value: true}
var singletonFALSE = &Bool{Value: false}
var SingletonNUll = &Null{}

func Evaluate(baseNode ASTNode) Object {
	switch node := baseNode.(type) {

	case Program:
		return evaluateProgram(node)

	case *ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression)

	case *Integer:
		CheckIsNotNil(node.Value)
		return &Number{Value: *node.Value}

	case *Boolean:
		CheckIsNotNil(node.Value)
		return toBooleanObject(*node.Value)

	case *Prefix:
		CheckIsNotNil(node.Rigth)
		rigth := Evaluate(node.Rigth)
		CheckIsNotNil(rigth)
		return evaluatePrefixExpression(node.Operator, rigth)

	case *Infix:
		go CheckIsNotNil(node.Left)
		CheckIsNotNil(node.Rigth)
		left := Evaluate(node.Left)
		rigth := Evaluate(node.Rigth)
		go CheckIsNotNil(left)
		CheckIsNotNil(rigth)
		return evaluateInfixExpression(node.Operator, left, rigth)

	case *Block:
		return evaluateBLockStaments(node)

	case *If:
		return evaluateIfExpression(node)

	case *ReturnStament:
		CheckIsNotNil(node.ReturnValue)
		value := Evaluate(node.ReturnValue)
		CheckIsNotNil(value)
		return &Return{Value: value}

	default:
		return SingletonNUll
	}
}

func CheckIsNotNil(val interface{}) {
	if val == nil {
		panic("expression or stament cannot be nil :(")
	}
}

func evaluateBLockStaments(block *Block) Object {
	var result Object = nil
	for _, statement := range block.Staments {
		result = Evaluate(statement)
		if result != nil && result.Type() == RETURNTYPE {
			return result
		}
	}

	return result
}

func evaluateProgram(program Program) Object {
	var result Object
	for _, statement := range program.Staments {
		result = Evaluate(statement)

		if _, isReturn := result.(*Return); isReturn {
			returnObj := result.(*Return)
			return returnObj.Value
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

func evaluateIfExpression(ifExpression *If) Object {
	CheckIsNotNil(ifExpression.Condition)
	condition := Evaluate(ifExpression.Condition)

	CheckIsNotNil(condition)
	if isTruthy(condition) {
		CheckIsNotNil(ifExpression.Consequence)
		return Evaluate(ifExpression.Consequence)
	} else if ifExpression.Alternative != nil {
		return Evaluate(ifExpression.Alternative)
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
	if left.Type() == INTEGERS && right.Type() == INTEGERS {
		return evaluateIntegerInfixExpression(operator, left, right)
	} else if operator == "==" {
		return toBooleanObject(left == right)
	} else if operator == "!=" {
		return toBooleanObject(left != right)
	}

	return SingletonNUll
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
		return SingletonNUll
	}
}

func evaluateMinusOperatorExpression(rigth Object) Object {
	if _, isNumber := rigth.(*Number); !isNumber {
		return SingletonNUll
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
		return SingletonNUll
	}
}

func toBooleanObject(val bool) Object {
	if val {
		return singletonTRUE
	}
	return singletonFALSE
}
