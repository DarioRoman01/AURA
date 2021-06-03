package test_test

import (
	"katan/src"
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
		{"1 >= 1", true},
		{"1 <= 1", true},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"verdadero == verdadero", true},
		{"verdadero == falso", false},
		{"verdadero != falso", true},
		{"verdadero || falso", true},
		{"verdadero && falso", false},
		{"verdadero && verdadero", true},
		{"verdadero || verdadero", true},
		{"falso || falso", false},
		{"falso && falso", false},
		{"falso == falso", true},
		{"(1 < 2) == verdadero", true},
		{"(1 < 2) == falso", false},
		{"(1 > 2) == verdadero", false},
		{"(1 > 2) == falso", true},
		{"(1 >= 2) == falso", true},
		{"(1 <= 2) == verdadero", true},
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
		{"si (4 % 2 == 0) { 10 } si_no { 20 }", 10},
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
		{source: "5 + verdadero;", expected: "Discrepancia de tipos: ENTERO + BOOLEANO"},
		{source: "5 + verdadero; 9;", expected: "Discrepancia de tipos: ENTERO + BOOLEANO"},
		{source: "-verdadero", expected: "Operador desconocido: -BOOLEANO"},
		{source: "verdadero + falso", expected: "Operador desconocido: BOOLEANO + BOOLEANO"},
		{source: "5; verdadero - falso; 10;", expected: "Operador desconocido: BOOLEANO - BOOLEANO"},
		{source: `
			si (10 > 7) {
				regresa verdadero + falso;
			}
		`,
			expected: "Operador desconocido: BOOLEANO + BOOLEANO",
		},
		{source: `
			si (5 < 2) {
				regresa 1;
			} si_no {
				regresa verdadero / falso;
			}
		`,
			expected: "Operador desconocido: BOOLEANO / BOOLEANO",
		},
		{source: "foobar;", expected: "Identificador no encontrado: foobar"},
		{source: `"foo" - "bar";`, expected: "Operador desconocido: TEXTO - TEXTO"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.Assert().IsType(&src.Error{}, evaluated.(*src.Error))
		evaluatedError := evaluated.(*src.Error)
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

func (e *EvaluatorTests) TestBuiltinFunctions() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{source: `longitud("");`, expected: 0},
		{source: `longitud("cuatro");`, expected: 6},
		{source: `longitud("hola mundo");`, expected: 10},
		{source: "longitud(1);", expected: "argumento para longitud no valido, se recibio ENTERO"},
		{
			source:   `longitud("uno", "dos");`,
			expected: "numero incorrecto de argumentos para longitud, se recibieron 2, se requieren 1",
		},
		{source: "tipo(1);", expected: "ENTERO"},
		{source: "tipo(verdadero)", expected: "BOOLEANO"},
		{source: `tipo("hello world")`, expected: "TEXTO"},
		{source: `entero("1")`, expected: 1},
		{source: `entero("hola")`, expected: "No se puede parsear como entero hola"},
		{source: "texto(12)", expected: "12"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if val, isInt := test.expected.(int); isInt {
			e.testIntegerObject(evaluated, val)
		} else {
			expected := test.expected.(string)

			if str, isStr := evaluated.(*src.String); isStr {
				e.testStringObject(str, expected)
			} else {
				e.testErrorObject(evaluated, expected)
			}
		}
	}
}

func (e *EvaluatorTests) testErrorObject(evlauated src.Object, expected string) {
	e.IsType(&src.Error{}, evlauated.(*src.Error))
	err := evlauated.(*src.Error)
	e.Equal(expected, err.Message)
}

func (e *EvaluatorTests) TestFunctionEvaluation() {
	source := "funcion(x) { x + 2; };"
	evaluated := e.evaluateTests(source)
	e.IsType(&src.Def{}, evaluated.(*src.Def))

	function := evaluated.(*src.Def)
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
			source:   `funcion() { regresa "src is awesome"; }()`,
			expected: "src is awesome",
		},
	}

	for _, test := range tests {
		evluated := e.evaluateTests(test.source)
		e.IsType(&src.String{}, evluated.(*src.String))
		stringObj := evluated.(*src.String)
		e.Equal(test.expected, stringObj.Value)
	}
}

func (e *EvaluatorTests) TestStringConcatenation() {
	tests := []struct {
		source   string
		expected string
	}{
		{source: `"foo" + "bar";`, expected: "foobar"},
		{source: `"hello," + " " + "world!";`, expected: "hello, world!"},
		{source: `
			var saludo = funcion(nombre) {
				regresa "Hola " + nombre + "!";
			};

			saludo("David");
		`,
			expected: "Hola David!",
		},
	}

	for _, test := range tests {
		evluated := e.evaluateTests(test.source)
		e.testStringObject(evluated, test.expected)
	}
}

func (e *EvaluatorTests) testStringObject(evaluated src.Object, expected string) {
	e.IsType(&src.String{}, evaluated.(*src.String))
	str := evaluated.(*src.String)
	e.Equal(expected, str.Value)
}

func (e *EvaluatorTests) TestStringComparison() {
	tests := []struct {
		source   string
		expected bool
	}{
		{`"a" == "a"`, true},
		{`"a" != "a"`, false},
		{`"a" == "b"`, false},
		{`"a" != "b"`, true},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testBooleanObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) testNullObject(eval src.Object) {
	e.Assert().Equal(src.SingletonNUll, eval.(*src.Null))
}

func (e *EvaluatorTests) evaluateTests(source string) src.Object {
	lexer := src.NewLexer(source)
	parser := src.NewParser(lexer)
	program := parser.ParseProgam()
	env := src.NewEnviroment(nil)
	evaluated := src.Evaluate(program, env)
	e.Assert().NotNil(evaluated)
	return evaluated
}

func (e *EvaluatorTests) testBooleanObject(object src.Object, expected bool) {
	e.Assert().IsType(&src.Bool{}, object.(*src.Bool))
	evaluated := object.(*src.Bool)
	e.Assert().Equal(expected, evaluated.Value)
}

func (e *EvaluatorTests) testIntegerObject(evaluated src.Object, expected int) {
	e.Assert().IsType(&src.Number{}, evaluated.(*src.Number))
	eval := evaluated.(*src.Number)
	e.Assert().Equal(expected, eval.Value)
}

func TestEvalutorSuite(t *testing.T) {
	suite.Run(t, new(EvaluatorTests))
}
