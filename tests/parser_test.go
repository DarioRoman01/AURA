package test_test

import (
	"aura/src/ast"
	l "aura/src/lexer"
	"aura/src/parser"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
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

type ParserTests struct {
	suite.Suite
}

func (p *ParserTests) InitParserTests(source string) (*parser.Parser, *ast.Program) {
	lexer := l.NewLexer(source)
	parser := parser.NewParser(lexer)
	program := parser.ParseProgam()

	return parser, &program
}

func (p *ParserTests) TestParseProgram() {
	source := "var x = 5;"
	_, program := p.InitParserTests(source)

	p.Assert().NotNil(program)

	p.Assert().IsType(&ast.Program{}, program)

	p.Assert().Implements((*ast.ASTNode)(nil), program)
}

func (p *ParserTests) TestLetStatements() {
	source := `
		var x = 5;
		var y = 10;
		var foo = 20;
		var bar = verdadero;
	`
	_, program := p.InitParserTests(source)

	p.Assert().Equal(4, len(program.Staments))
	expected := []PrefixTuple{
		{Operator: "x", Value: 5},
		{Operator: "y", Value: 10},
		{Operator: "foo", Value: 20},
		{Operator: "bar", Value: true},
	}

	for i, statement := range program.Staments {
		p.Assert().Equal("var", statement.TokenLiteral())
		p.Assert().IsType(&ast.LetStatement{}, statement.(*ast.LetStatement))

		letStatement := statement.(*ast.LetStatement)
		p.Assert().NotNil(letStatement.Name)
		p.testIdentifier(letStatement.Name, expected[i].Operator)

		p.Assert().NotNil(letStatement.Value)
		p.testLiteralExpression(letStatement.Value, expected[i].Value)
	}
}

func (p *ParserTests) TestNamesInLetStatements() {
	source := `
		var x = 5;
		var y = 10;
		var foo = 20;
	`
	_, program := p.InitParserTests(source)
	p.Assert().Equal(3, len(program.Staments))

	var names []string
	for _, stament := range program.Staments {
		stament := stament.(*ast.LetStatement)
		if !p.Assert().NotNil(stament.Name) {
			p.T().Fail()
		}

		if !p.Assert().Implements((*ast.Stmt)(nil), stament) {
			p.T().Fail()
		}

		names = append(names, stament.Name.Str())
	}

	expectedNames := []string{"x", "y", "foo"}
	if !p.Assert().Equal(expectedNames, names) {
		p.T().Fail()
	}

}

func (p *ParserTests) TestParseErrors() {
	source := "var x 5;"
	parser, _ := p.InitParserTests(source)
	if !p.Assert().Equal(1, len(parser.Errors())) {
		p.T().Fail()
	}
}

func (p *ParserTests) TestReturnStatement() {
	source := `
		regresa 5;
		regresa foo;
		regresa verdadero;
		regresa falso;
	`

	_, program := p.InitParserTests(source)

	if !p.Assert().Equal(4, len(program.Staments)) {
		p.T().Log("len of program statements are not 2")
		p.T().Fail()
	}

	expected := []PrefixTuple{
		{Operator: "regresa", Value: 5},
		{Operator: "regresa", Value: "foo"},
		{Operator: "regresa", Value: true},
		{Operator: "regresa", Value: false},
	}

	for i, statement := range program.Staments {
		p.Assert().Equal("regresa", statement.TokenLiteral())
		p.Assert().IsType(&ast.ReturnStament{}, statement.(*ast.ReturnStament))

		returnStament := statement.(*ast.ReturnStament)
		p.Assert().NotNil(returnStament.ReturnValue)
		p.testLiteralExpression(returnStament.ReturnValue, expected[i].Value)
	}
}

func (p *ParserTests) TestIdentifierExpression() {
	source := "foobar;"
	parser, program := p.InitParserTests(source)

	p.testProgramStatements(parser, program, 1)

	expressionStament := program.Staments[0].(*ast.ExpressionStament)
	if !p.Assert().NotNil(expressionStament.Expression) {
		p.T().Fail()
	}

	p.testLiteralExpression(expressionStament.Expression, "foobar")
}

func (p *ParserTests) TestClassStatement() {
	source := `
		clase Persona {
			constructor(nombre, carrera) {
				this.nombre = nombre;
				this.carrera = carrera;
			}

			saludar() {
				escribir("Hola soy ", this.nombre, " y estudio ", this.carrera)
			}
		}
	`
	parser, program := p.InitParserTests(source)
	p.Assert().Equal(1, len(program.Staments))
	p.Assert().Equal(0, len(parser.Errors()))

	p.Assert().IsType(&ast.ClassStatement{}, program.Staments[0].(*ast.ClassStatement))
	class := program.Staments[0].(*ast.ClassStatement)
	p.Assert().Equal(class.Name.Value, "Persona")
	p.Assert().Equal(len(class.Methods), 2)
}

func (p *ParserTests) TestIntegerExpressions() {
	source := "5;"
	parser, program := p.InitParserTests(source)

	p.testProgramStatements(parser, program, 1)
	expressionStament := program.Staments[0].(*ast.ExpressionStament)
	p.Assert().NotNil(expressionStament.Expression)
	p.testLiteralExpression(expressionStament.Expression, 5)
}

func (p *ParserTests) TestPrefixExpressions() {
	source := "!5; -15; !verdadero; !falso;"
	parser, program := p.InitParserTests(source)

	p.testProgramStatements(parser, program, 4)
	expectedExpressions := []PrefixTuple{
		{Operator: "!", Value: 5},
		{Operator: "-", Value: 15},
		{Operator: "!", Value: true},
		{Operator: "!", Value: false},
	}

	if len(program.Staments) == len(expectedExpressions) {
		for i, stament := range program.Staments {
			stament := stament.(*ast.ExpressionStament)
			p.Assert().IsType(&ast.Prefix{}, stament.Expression.(*ast.Prefix))

			prefix := stament.Expression.(*ast.Prefix)
			p.Assert().Equal(prefix.Operator, expectedExpressions[i].Operator)

			p.Assert().NotNil(prefix.Rigth)
			p.testLiteralExpression(prefix.Rigth, expectedExpressions[i].Value)
		}
	} else {
		p.T().Log("len of staments and expected expected expressions are not equal")
		p.T().Fail()
	}
}

func (p *ParserTests) TestLIstExpression() {
	source := "lista[1,2,3]"
	parser, program := p.InitParserTests(source)
	p.testProgramStatements(parser, program, 1)

	list := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.Array)
	p.Assert().IsType(&ast.Array{}, list)
	p.Assert().NotNil(list.Values)
	p.Assert().Equal(3, len(list.Values))
	p.testInteger(list.Values[0], 1)
	p.testInteger(list.Values[1], 2)
	p.testInteger(list.Values[2], 3)
}

func (p *ParserTests) TestArrayExpression() {
	source := `mapa{"a" => 1, "b" => 2, "c" => 3}`
	parser, program := p.InitParserTests(source)
	p.testProgramStatements(parser, program, 1)
	hashMap := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.MapExpression)
	p.Assert().IsType(&ast.MapExpression{}, hashMap)
	p.Assert().NotNil(hashMap.Body)
	p.Assert().Equal(3, len(hashMap.Body))

	p.Assert().NotNil(hashMap.Body)
	p.Assert().Equal(3, len(hashMap.Body))

	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}

	for idx, keyVal := range hashMap.Body {
		p.Assert().IsType(&ast.KeyValue{}, keyVal)
		p.Assert().IsType(&ast.StringLiteral{}, keyVal.Key.(*ast.StringLiteral))
		p.Assert().IsType(&ast.Integer{}, keyVal.Value.(*ast.Integer))

		key := keyVal.Key.(*ast.StringLiteral)
		value := keyVal.Value.(*ast.Integer)
		p.Assert().Equal(expectedKeys[idx], key.Value)
		p.Assert().Equal(expectedValues[idx], *value.Value)
	}
}

func (p *ParserTests) TestCallExpression() {
	source := "suma(1, 2 * 3, 4 + 5);"
	parser, program := p.InitParserTests(source)
	p.testProgramStatements(parser, program, 1)

	call := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.Call)
	p.Assert().IsType(&ast.Call{}, call)
	p.testIdentifier(call.Function, "suma")
	p.Assert().NotNil(call.Arguments)

	p.Assert().Equal(3, len(call.Arguments))
	p.testLiteralExpression(call.Arguments[0], 1)
	p.testInfixExpression(call.Arguments[1], 2, "*", 3)
	p.testInfixExpression(call.Arguments[2], 4, "+", 5)
}

func (p *ParserTests) TestFunctionLiteral() {
	source := "funcion(x, y) { x + y }"
	parser, program := p.InitParserTests(source)
	p.testProgramStatements(parser, program, 1)

	functionLiteral := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.Function)
	p.Assert().IsType(&ast.Function{}, functionLiteral)
	p.Assert().Equal(2, len(functionLiteral.Parameters))

	p.testLiteralExpression(functionLiteral.Parameters[0], "x")
	p.testLiteralExpression(functionLiteral.Parameters[1], "y")
	p.Assert().NotNil(functionLiteral.Body)

	p.Assert().Equal(1, len(functionLiteral.Body.Staments))
	body := functionLiteral.Body.Staments[0].(*ast.ExpressionStament)
	p.Assert().NotNil(body.Expression)
	p.testInfixExpression(body.Expression, "x", "+", "y")
}

func (p *ParserTests) TestFunctionParameter() {
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
		_, program := p.InitParserTests(test["input"].(string))
		function := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.Function)
		p.Assert().Equal(len(test["expected"].([]string)), len(function.Parameters))

		for idx, param := range test["expected"].([]string) {
			p.testLiteralExpression(function.Parameters[idx], param)
		}
	}
}

func (p *ParserTests) TestInfixExpressions() {
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
	parser, program := p.InitParserTests(source)

	p.testProgramStatements(parser, program, 12)
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
		stament := stamment.(*ast.ExpressionStament)
		p.Assert().NotNil(stament.Expression)
		p.Assert().IsType(&ast.Infix{}, stament.Expression.(*ast.Infix))
		p.testInfixExpression(
			stament.Expression,
			expectedOperators[i].Left,
			expectedOperators[i].Operator,
			expectedOperators[i].Rigth,
		)
	}
}

func (p *ParserTests) TestIfExpression() {
	source := "si (x > y) { z } si_no { w }"
	parser, program := p.InitParserTests(source)
	p.testProgramStatements(parser, program, 1)

	ifExpression := (program.Staments[0].(*ast.ExpressionStament)).Expression.(*ast.If)
	p.Assert().IsType(&ast.If{}, ifExpression)
	p.Assert().NotNil(ifExpression.Condition)

	p.testInfixExpression(ifExpression.Condition, "x", ">", "y")
	p.Assert().NotNil(ifExpression.Consequence)
	p.Assert().IsType(&ast.Block{}, ifExpression.Consequence)
	p.Assert().Equal(1, len(ifExpression.Consequence.Staments))

	consequenceStament := ifExpression.Consequence.Staments[0].(*ast.ExpressionStament)
	p.Assert().NotNil(consequenceStament.Expression)
	p.testIdentifier(consequenceStament.Expression, "z")

	p.Assert().NotNil(ifExpression.Alternative)
	p.Assert().IsType(&ast.Block{}, ifExpression.Alternative)
	p.Assert().Equal(1, len(ifExpression.Alternative.Staments))

	alternativeStament := ifExpression.Alternative.Staments[0].(*ast.ExpressionStament)
	p.Assert().NotNil(alternativeStament.Expression)
	p.testIdentifier(alternativeStament.Expression, "w")
}

func (p *ParserTests) TestBooleanExpressions() {
	source := "verdadero; falso;"
	parser, program := p.InitParserTests(source)

	p.testProgramStatements(parser, program, 2)
	expectedValues := []bool{true, false}

	for i, stament := range program.Staments {
		expressionStament := stament.(*ast.ExpressionStament)
		p.Assert().NotNil(expressionStament.Expression)
		p.testLiteralExpression(expressionStament.Expression, expectedValues[i])
	}
}

func (p *ParserTests) TestOperatorPrecedence() {
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
		parser, program := p.InitParserTests(source.source)

		p.testProgramStatements(parser, program, source.expectedCount)
		p.Assert().Equal(source.expected, program.Str())
	}
}

func (p *ParserTests) TestStringLiteral() {
	source := `"hello world!";`

	_, program := p.InitParserTests(source)
	expressionStatement := program.Staments[0].(*ast.ExpressionStament)
	stringLiteral := expressionStatement.Expression.(*ast.StringLiteral)

	p.Assert().IsType(&ast.StringLiteral{}, stringLiteral)
	p.Assert().Equal("hello world!", stringLiteral.Value)
}

func (p *ParserTests) testBoolean(expression ast.Expression, expectedValue bool) {
	boolean := expression.(*ast.Boolean)
	p.Assert().Equal(*boolean.Value, expectedValue)
	if expectedValue {
		p.Assert().Equal("verdadero", boolean.Token.Literal)
	} else {
		p.Assert().Equal("falso", boolean.Token.Literal)
	}
}

func (p *ParserTests) testInfixExpression(ex ast.Expression, expectedLeft interface{}, operator string, expectedRigth interface{}) {
	infix := ex.(*ast.Infix)
	p.Assert().NotNil(infix.Left)
	p.testLiteralExpression(infix.Left, expectedLeft)
	p.Assert().Equal(operator, infix.Operator)
	p.Assert().NotNil(infix.Rigth)
	p.testLiteralExpression(infix.Rigth, expectedRigth)
}

func (p *ParserTests) testProgramStatements(parser *parser.Parser, program *ast.Program, expectedStamenetCount int) {
	p.Assert().Equal(0, len(parser.Errors()))
	p.Assert().Equal(expectedStamenetCount, len(program.Staments))
	p.Assert().IsType(&ast.ExpressionStament{}, program.Staments[0].(*ast.ExpressionStament))
}

func (p *ParserTests) testLiteralExpression(expression ast.Expression, expectedValue interface{}) {
	switch expectedValue := expectedValue.(type) {
	case string:
		p.testIdentifier(expression, expectedValue)
	case int:
		p.testInteger(expression, expectedValue)
	case bool:
		p.testBoolean(expression, expectedValue)
	default:
		p.T().Log(fmt.Sprintf("unhandled type of expression, Got=%s", reflect.TypeOf(expectedValue).String()))
		p.T().Fail()
	}
}

func (p *ParserTests) testIdentifier(expression ast.Expression, expectedValue string) {
	p.Assert().IsType(&ast.Identifier{}, expression.(*ast.Identifier))

	identifier := expression.(*ast.Identifier)
	p.Assert().Equal(expectedValue, identifier.Str())
	p.Assert().Equal(expectedValue, identifier.TokenLiteral())
}

func (p *ParserTests) testInteger(expression ast.Expression, expectedValue int) {
	p.Assert().IsType(&ast.Integer{}, expression.(*ast.Integer))
	integer := expression.(*ast.Integer)
	p.Assert().Equal(expectedValue, *integer.Value)
	p.Assert().Equal(fmt.Sprint(expectedValue), integer.Token.Literal)
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserTests))
}
