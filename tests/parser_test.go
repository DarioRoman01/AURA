package test_test

import (
	"fmt"
	"lpp/lpp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PrefixTuple struct {
	Operator string
	Value    interface{}
}

type InfixTuple struct {
	Left     interface{}
	Operator string
	Rigth    interface{}
}

func InitParserTests(source string) (*lpp.Parser, *lpp.Program) {
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()

	return parser, &program
}

func TestParseProgram(t *testing.T) {
	assert := assert.New(t)
	source := "var x = 5;"
	_, program := InitParserTests(source)

	if !assert.NotNil(program) {
		t.Log("program is nil")
		t.Fail()
	}

	if !assert.IsType(&lpp.Program{}, program) {
		t.Log("program type is not Program!")
		t.Fail()
	}

	if !assert.Implements((*lpp.ASTNode)(nil), program) {
		t.Log("program type does not implement ast node interface")
		t.Fail()
	}
}

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)
	source := `
		var x = 5;
		var y = 10;
		var foo = 20;
		var bar = verdadero;
	`
	_, program := InitParserTests(source)

	if !assert.Equal(4, len(program.Staments)) {
		t.Log("len of program statements are not 3")
		t.Fail()
	}
	expected := []PrefixTuple{
		{Operator: "x", Value: 5},
		{Operator: "y", Value: 10},
		{Operator: "foo", Value: 20},
		{Operator: "bar", Value: true},
	}

	for i, statement := range program.Staments {
		assert.Equal("var", statement.TokenLiteral())
		assert.IsType(&lpp.LetStatement{}, statement.(*lpp.LetStatement))

		letStatement := statement.(*lpp.LetStatement)
		assert.NotNil(letStatement.Name)
		testIdentifier(t, letStatement.Name, expected[i].Operator)

		assert.NotNil(letStatement.Value)
		testLiteralExpression(t, letStatement.Value, expected[i].Value)
	}
}

func TestNamesInLetStatements(t *testing.T) {
	assert := assert.New(t)
	source := `
		var x = 5;
		var y = 10;
		var foo = 20;
	`
	_, program := InitParserTests(source)
	assert.Equal(3, len(program.Staments))

	var names []string
	for _, stament := range program.Staments {
		stament := stament.(*lpp.LetStatement)
		if !assert.NotNil(stament.Name) {
			t.Fail()
		}

		if !assert.Implements((*lpp.Stmt)(nil), stament) {
			t.Fail()
		}

		names = append(names, stament.Name.Str())
	}

	expectedNames := []string{"x", "y", "foo"}
	if !assert.Equal(expectedNames, names) {
		t.Fail()
	}

}

func TestParseErrors(t *testing.T) {
	source := "var x 5;"
	parser, _ := InitParserTests(source)

	if !assert.Equal(t, 1, len(parser.Errors())) {
		t.Fail()
	}
}

func TestReturnStatement(t *testing.T) {
	assert := assert.New(t)
	source := `
		regresa 5;
		regresa foo;
		regresa verdadero;
		regresa falso;
	`

	_, program := InitParserTests(source)

	if !assert.Equal(4, len(program.Staments)) {
		t.Log("len of program statements are not 2")
		t.Fail()
	}

	expected := []PrefixTuple{
		{Operator: "regresa", Value: 5},
		{Operator: "regresa", Value: "foo"},
		{Operator: "regresa", Value: true},
		{Operator: "regresa", Value: false},
	}

	for i, statement := range program.Staments {
		assert.Equal("regresa", statement.TokenLiteral())
		assert.IsType(&lpp.ReturnStament{}, statement.(*lpp.ReturnStament))

		returnStament := statement.(*lpp.ReturnStament)
		assert.NotNil(returnStament.ReturnValue)
		testLiteralExpression(t, returnStament.ReturnValue, expected[i].Value)
	}
}

func TestIdentifierExpression(t *testing.T) {
	source := "foobar;"
	parser, program := InitParserTests(source)

	testProgramStatements(t, parser, program, 1)

	expressionStament := program.Staments[0].(*lpp.ExpressionStament)
	if !assert.NotNil(t, expressionStament.Expression) {
		t.Fail()
	}

	testLiteralExpression(t, expressionStament.Expression, "foobar")
}

func TestIntegerExpressions(t *testing.T) {
	source := "5;"
	parser, program := InitParserTests(source)

	testProgramStatements(t, parser, program, 1)
	expressionStament := program.Staments[0].(*lpp.ExpressionStament)
	assert.NotNil(t, expressionStament.Expression)
	testLiteralExpression(t, expressionStament.Expression, 5)
}

func TestPrefixExpressions(t *testing.T) {
	source := "!5; -15; !verdadero; !falso;"
	parser, program := InitParserTests(source)

	testProgramStatements(t, parser, program, 4)
	expectedExpressions := []PrefixTuple{
		{Operator: "!", Value: 5},
		{Operator: "-", Value: 15},
		{Operator: "!", Value: true},
		{Operator: "!", Value: false},
	}

	if len(program.Staments) == len(expectedExpressions) {
		for i, stament := range program.Staments {
			stament := stament.(*lpp.ExpressionStament)
			assert.IsType(t, &lpp.Prefix{}, stament.Expression.(*lpp.Prefix))

			prefix := stament.Expression.(*lpp.Prefix)
			assert.Equal(t, prefix.Operator, expectedExpressions[i].Operator)

			assert.NotNil(t, prefix.Rigth)
			testLiteralExpression(t, prefix.Rigth, expectedExpressions[i].Value)
		}
	} else {
		t.Log("len of staments and expected expected expressions are not equal")
		t.Fail()
	}
}

func TestCallExpression(t *testing.T) {
	assert := assert.New(t)
	source := "suma(1, 2 * 3, 4 + 5);"
	parser, program := InitParserTests(source)
	testProgramStatements(t, parser, program, 1)

	call := (program.Staments[0].(*lpp.ExpressionStament)).Expression.(*lpp.Call)
	assert.IsType(&lpp.Call{}, call)
	testIdentifier(t, call.Function, "suma")
	assert.NotNil(call.Arguments)

	assert.Equal(3, len(call.Arguments))
	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}

func TestFunctionLiteral(t *testing.T) {
	assert := assert.New(t)
	source := "funcion(x, y) { x + y }"
	parser, program := InitParserTests(source)
	testProgramStatements(t, parser, program, 1)

	functionLiteral := (program.Staments[0].(*lpp.ExpressionStament)).Expression.(*lpp.Function)
	assert.IsType(&lpp.Function{}, functionLiteral)
	assert.Equal(2, len(functionLiteral.Parameters))

	testLiteralExpression(t, functionLiteral.Parameters[0], "x")
	testLiteralExpression(t, functionLiteral.Parameters[1], "y")
	assert.NotNil(functionLiteral.Body)

	assert.Equal(1, len(functionLiteral.Body.Staments))
	body := functionLiteral.Body.Staments[0].(*lpp.ExpressionStament)
	assert.NotNil(body.Expression)
	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestFunctionParameter(t *testing.T) {
	assert := assert.New(t)
	tests := []map[string]interface{}{
		{
			"input":    "funcion() {};",
			"expected": []string{},
		},
		{
			"input":    "funcion(x) {};",
			"expected": []string{"x"},
		},
		{
			"input":    "funcion(x, y, z) {};",
			"expected": []string{"x", "y", "z"},
		},
	}

	for _, test := range tests {
		_, program := InitParserTests(test["input"].(string))
		function := (program.Staments[0].(*lpp.ExpressionStament)).Expression.(*lpp.Function)
		assert.Equal(len(test["expected"].([]string)), len(function.Parameters))

		for idx, param := range test["expected"].([]string) {
			testLiteralExpression(t, function.Parameters[idx], param)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	source := `
		5 + 5;
		5 - 5;
		5 * 5;
		5 / 5;
		5 > 5;
		5 < 5;
		5 == 5;
		5 != 5;
		verdadero == verdadero;
		verdadero != verdadero;
		falso == verdadero;
		falso != verdadero;
		
	`
	parser, program := InitParserTests(source)

	testProgramStatements(t, parser, program, 12)
	expectedOperators := []InfixTuple{
		{Left: 5, Operator: "+", Rigth: 5},
		{Left: 5, Operator: "-", Rigth: 5},
		{Left: 5, Operator: "*", Rigth: 5},
		{Left: 5, Operator: "/", Rigth: 5},
		{Left: 5, Operator: ">", Rigth: 5},
		{Left: 5, Operator: "<", Rigth: 5},
		{Left: 5, Operator: "==", Rigth: 5},
		{Left: 5, Operator: "!=", Rigth: 5},
		{Left: true, Operator: "==", Rigth: true},
		{Left: true, Operator: "!=", Rigth: true},
		{Left: false, Operator: "==", Rigth: true},
		{Left: false, Operator: "!=", Rigth: true},
	}

	for i, stamment := range program.Staments {
		stament := stamment.(*lpp.ExpressionStament)
		assert.NotNil(t, stament.Expression)
		assert.IsType(t, &lpp.Infix{}, stament.Expression.(*lpp.Infix))
		testInfixExpression(t,
			stament.Expression,
			expectedOperators[i].Left,
			expectedOperators[i].Operator,
			expectedOperators[i].Rigth,
		)
	}
}

func TestIfExpression(t *testing.T) {
	assert := assert.New(t)
	source := "si (x > y) { z } si_no { w }"
	parser, program := InitParserTests(source)
	testProgramStatements(t, parser, program, 1)

	ifExpression := (program.Staments[0].(*lpp.ExpressionStament)).Expression.(*lpp.If)
	assert.IsType(&lpp.If{}, ifExpression)
	assert.NotNil(ifExpression.Condition)

	testInfixExpression(t, ifExpression.Condition, "x", ">", "y")
	assert.NotNil(ifExpression.Consequence)
	assert.IsType(&lpp.Block{}, ifExpression.Consequence)
	assert.Equal(1, len(ifExpression.Consequence.Staments))

	consequenceStament := ifExpression.Consequence.Staments[0].(*lpp.ExpressionStament)
	assert.NotNil(consequenceStament.Expression)
	testIdentifier(t, consequenceStament.Expression, "z")

	assert.NotNil(ifExpression.Alternative)
	assert.IsType(&lpp.Block{}, ifExpression.Alternative)
	assert.Equal(1, len(ifExpression.Alternative.Staments))

	alternativeStament := ifExpression.Alternative.Staments[0].(*lpp.ExpressionStament)
	assert.NotNil(alternativeStament.Expression)
	testIdentifier(t, alternativeStament.Expression, "w")
}

func TestBooleanExpressions(t *testing.T) {
	source := "verdadero; falso;"
	parser, program := InitParserTests(source)

	testProgramStatements(t, parser, program, 2)
	expectedValues := []bool{true, false}

	for i, stament := range program.Staments {
		expressionStament := stament.(*lpp.ExpressionStament)
		assert.NotNil(t, expressionStament.Expression)
		testLiteralExpression(t, expressionStament.Expression, expectedValues[i])
	}
}

func TestOperatorPrecedence(t *testing.T) {
	type TupleToTest struct {
		source        string
		expected      string
		expectedCount int
	}
	test_source := []TupleToTest{
		{source: "-a * b;", expected: "((- a) * b)", expectedCount: 1},
		{source: "!-a;", expected: "(! (- a))", expectedCount: 1},
		{source: "a + b / c;", expected: "(a + (b / c))", expectedCount: 1},
		{source: "3 + 4; -5 * 5;", expected: "(3 + 4) ((- 5) * 5)", expectedCount: 2},
		{source: "a + b * c + d / e - f;", expected: "(((a + (b * c)) + (d / e)) - f)", expectedCount: 1},
		{source: "1 + (2 + 3) + 4;", expected: "((1 + (2 + 3)) + 4)", expectedCount: 1},
		{source: "(5 + 5) * 2;", expected: "((5 + 5) * 2)", expectedCount: 1},
		{source: "2 / (5 + 5);", expected: "(2 / (5 + 5))", expectedCount: 1},
		{source: "-(5 + 5);", expected: "(- (5 + 5))", expectedCount: 1},
		{source: "-(5 + 5);", expected: "(- (5 + 5))", expectedCount: 1},
		{source: "a + suma(b * c) + d;", expected: "((a + suma((b * c))) + d)", expectedCount: 1},
		{source: "a + suma(b * c) + d;", expected: "((a + suma((b * c))) + d)", expectedCount: 1},
		{
			source:        "suma(a, b, 1, 2 * 3, 4 + 5, suma(6, 7 * 8))",
			expected:      "suma(a, b, 1, (2 * 3), (4 + 5), suma(6, (7 * 8)))",
			expectedCount: 1,
		},
	}

	for _, source := range test_source {
		parser, program := InitParserTests(source.source)

		testProgramStatements(t, parser, program, source.expectedCount)
		assert.Equal(t, source.expected, program.Str())
	}
}

func testBoolean(t *testing.T, expression lpp.Expression, expectedValue bool) {
	boolean := expression.(*lpp.Boolean)
	assert.Equal(t, *boolean.Value, expectedValue)
	if expectedValue {
		assert.Equal(t, "verdadero", boolean.Token.Literal)
	} else {
		assert.Equal(t, "falso", boolean.Token.Literal)
	}
}

func testInfixExpression(t *testing.T, ex lpp.Expression, expectedLeft interface{}, operator string, expectedRigth interface{}) {
	infix := ex.(*lpp.Infix)
	assert.NotNil(t, infix.Left)
	testLiteralExpression(t, infix.Left, expectedLeft)
	assert.Equal(t, operator, infix.Operator)
	assert.NotNil(t, infix.Rigth)
	testLiteralExpression(t, infix.Rigth, expectedRigth)
}

func testProgramStatements(t *testing.T, parser *lpp.Parser, program *lpp.Program, expectedStamenetCount int) {
	assert := assert.New(t)
	assert.Equal(0, len(parser.Errors()))
	assert.Equal(expectedStamenetCount, len(program.Staments))
	assert.IsType(&lpp.ExpressionStament{}, program.Staments[0].(*lpp.ExpressionStament))
}

func testLiteralExpression(t *testing.T, expression lpp.Expression, expectedValue interface{}) {
	switch expectedValue := expectedValue.(type) {
	case string:
		testIdentifier(t, expression, expectedValue)
	case int:
		testInteger(t, expression, expectedValue)
	case bool:
		testBoolean(t, expression, expectedValue)
	default:
		t.Log(fmt.Sprintf("unhandled type of expression, Got=%s", reflect.TypeOf(expectedValue).String()))
		t.Fail()
	}
}

func testIdentifier(t *testing.T, expression lpp.Expression, expectedValue string) {
	assert := assert.New((t))
	assert.IsType(&lpp.Identifier{}, expression.(*lpp.Identifier))

	identifier := expression.(*lpp.Identifier)
	assert.Equal(expectedValue, identifier.Str())
	assert.Equal(expectedValue, identifier.TokenLiteral())
}

func testInteger(t *testing.T, expression lpp.Expression, expectedValue int) {
	assert.IsType(t, &lpp.Integer{}, expression.(*lpp.Integer))
	integer := expression.(*lpp.Integer)
	assert.Equal(t, expectedValue, *integer.Value)
	assert.Equal(t, fmt.Sprint(expectedValue), integer.Token.Literal)
}
