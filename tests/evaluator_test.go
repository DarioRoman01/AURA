package test_test

import (
	"lpp/lpp"
	"testing"

	"github.com/stretchr/testify/suite"
)

type tuple struct {
	source   string
	expected int
}

type EvaluatorTests struct {
	suite.Suite
}

func (e *EvaluatorTests) TestIntegerEvaluation() {
	tests := []tuple{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5;", 10},
		{"5 - 10;", -5},
		{"2 * 2 * 2 * 2;", 16},
		{"50 / 2;", 25},
		{"2 * (5 - 3)", 4},
		{"(2 + 7) / 3", 3},
		{"50 / 2 * 2 + 10", 60},
		{"5 / 2", 2},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)

	}
}

func (e *EvaluatorTests) TestBangOperator() {
	tests := []struct {
		source   string
		expected bool
	}{
		{"!verdadero", false},
		{"!falso", true},
		{"!!verdadero", true},
		{"!!falso", false},
		{"!5", false},
		{"!!5", true},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testBooleanObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestBooleanEvaluation() {
	tests := []struct {
		source   string
		expected bool
	}{
		{"verdadero", true},
		{"falso", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 2", true},
		{"verdadero == verdadero", true},
		{"verdadero == falso", false},
		{"verdadero != falso", true},
		{"falso == falso", true},
		{"(1 < 2) == verdadero", true},
		{"(1 < 2) == falso", false},
		{"(1 > 2) == verdadero", false},
		{"(1 > 2) == falso", true},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testBooleanObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestIfElseEvaluation() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{"si (verdadero) { 10 }", 10},
		{"si (falso) { 10 }", nil},
		{"si (1) { 10 }", 10},
		{"si (1 < 2) { 10 }", 10},
		{"si (1 > 2) { 10 }", nil},
		{"si (1 < 2) { 10 } si_no { 20 }", 10},
		{"si (1 > 2) { 10 } si_no { 20 }", 20},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if _, isInt := test.expected.(int); isInt {
			e.testIntegerObject(evaluated, test.expected.(int))
		} else {
			e.testNullObject(evaluated)
		}
	}
}

func (e *EvaluatorTests) TestReturnEvaluation() {
	tests := []tuple{
		{"regresa 10;", 10},
		{"regresa 10; 9;", 10},
		{"regresa 2 * 5; 9;", 10},
		{"9; regresa 3 * 6; 9;", 18},
		{source: `
			si (10 > 1) {
				si (20 > 10) {
					regresa 1;
				}

				regresa 0;
			}
		`,
			expected: 1,
		},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestErrorhandling() {
	tests := []struct {
		source   string
		expected string
	}{
		{source: "5 + verdadero;", expected: "Discrepancia de tipos: INTEGERS + BOOLEAN"},
		{source: "5 + verdadero; 9;", expected: "Discrepancia de tipos: INTEGERS + BOOLEAN"},
		{source: "-verdadero", expected: "Operador desconocido: -BOOLEAN"},
		{source: "verdadero + falso", expected: "Operador desconocido: BOOLEAN + BOOLEAN"},
		{source: "5; verdadero - falso; 10;", expected: "Operador desconocido: BOOLEAN - BOOLEAN"},
		{source: `
			si (10 > 7) {
				regresa verdadero + falso;
			}
		`,
			expected: "Operador desconocido: BOOLEAN + BOOLEAN",
		},
		{source: `
			si (5 < 2) {
				regresa 1;
			} si_no {
				regresa verdadero / falso;
			}
		`,
			expected: "Operador desconocido: BOOLEAN / BOOLEAN",
		},
		{source: "foobar;", expected: "Identificador no encontrado: foobar"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.Assert().IsType(&lpp.Error{}, evaluated.(*lpp.Error))
		evaluatedError := evaluated.(*lpp.Error)
		e.Assert().Equal(test.expected, evaluatedError.Message)
	}
}

func (e *EvaluatorTests) TestAssingmentEvaluation() {
	tests := []tuple{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestFunctionEvaluation() {
	source := "funcion(x) { x + 2; };"
	evaluated := e.evaluateTests(source)
	e.IsType(&lpp.Def{}, evaluated.(*lpp.Def))

	function := evaluated.(*lpp.Def)
	e.Equal(1, len(function.Parameters))
	e.Equal("x", function.Parameters[0].Str())
	e.Equal("(x + 2)", function.Body.Str())
}

func (e *EvaluatorTests) TestFunctionCalls() {
	tests := []tuple{
		{"var identidad = funcion(x) { x }; identidad(5);", 5},
		{`
			var identidad = funcion(x) {
				regresa x;
			};

			identidad(5);
		`, 5,
		},
		{`
			var doble = funcion(x) {
				regresa 2 * x;
			};

			doble(5);
		`, 10,
		},
		{`
			var suma = funcion(x, y) {
				regresa x + y;
			};

			suma(3, 8);
		`, 11,
		},
		{`
			var suma = funcion(x, y) {
				regresa x + y;
			};

			suma(5 + 5, suma(10, 10));
		`, 30,
		},
		{"funcion(x) { x }(5)", 5},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestStringEvaluation() {
	tests := []struct {
		source   string
		expected string
	}{
		{source: `"hello world!"`, expected: "hello world!"},
		{
			source:   `funcion() { regresa "lpp is awesome"; }()`,
			expected: "lpp is awesome",
		},
	}

	for _, test := range tests {
		evluated := e.evaluateTests(test.source)
		e.IsType(&lpp.String{}, evluated.(*lpp.String))
		stringObj := evluated.(*lpp.String)
		e.Equal(test.expected, stringObj.Value)
	}
}

func (e *EvaluatorTests) testNullObject(eval lpp.Object) {
	e.Assert().Equal(lpp.SingletonNUll, eval.(*lpp.Null))
}

func (e *EvaluatorTests) evaluateTests(source string) lpp.Object {
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()
	env := lpp.NewEnviroment(nil)
	evaluated := lpp.Evaluate(program, env)
	e.Assert().NotNil(evaluated)
	return evaluated
}

func (e *EvaluatorTests) testBooleanObject(object lpp.Object, expected bool) {
	e.Assert().IsType(&lpp.Bool{}, object.(*lpp.Bool))
	evaluated := object.(*lpp.Bool)
	e.Assert().Equal(expected, evaluated.Value)
}

func (e *EvaluatorTests) testIntegerObject(evaluated lpp.Object, expected int) {
	e.Assert().IsType(&lpp.Number{}, evaluated.(*lpp.Number))
	eval := evaluated.(*lpp.Number)
	e.Assert().Equal(expected, eval.Value)
}

func TestEvalutorSuite(t *testing.T) {
	suite.Run(t, new(EvaluatorTests))
}
