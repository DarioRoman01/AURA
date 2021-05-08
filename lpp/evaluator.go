package lpp

func Evaluate(baseNode ASTNode) Object {
	switch node := baseNode.(type) {

	case Program:
		return evaluateStaments(node.Staments)

	case *ExpressionStament:
		CheckIsNotNil(node.Expression)
		return Evaluate(node.Expression)

	case *Integer:
		CheckIsNotNil(node.Value)
		return NewNumber(*node.Value)

	default:
		return nil
	}
}

func evaluateStaments(statements []Stmt) Object {
	var result Object

	for _, statement := range statements {
		result = Evaluate(statement)
	}

	return result
}

func CheckIsNotNil(val interface{}) {
	if val == nil {
		panic("expression or stament cannot be nil :(")
	}
}
