package test

import (
	"aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
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

func (e *EvaluatorTests) TestFloatEvaluation() {
	tests := []struct {
		souce    string
		expected float64
	}{
		{`5.5`, 5.5},
		{`-5.5`, -5.5},
		{`2.2 + 2.7`, 4.9},
		{`5 + 2.7`, 7.7},
		{"5.5 += 2.4", 7.9},
		{"5.5 -= 2.4", 3.1},
		{"5.5 * 5.5", 30.25},
		{"5.5 *= 5.5", 30.25},
		{"5.5 / 2", 2.75},
		{"5.5 /= 2", 2.75},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.souce)
		e.testFloatObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestArrayEvaluation() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{
			source:   "lista[2,4,5];",
			expected: []int{2, 4, 5},
		},
		{
			source:   "lista[45,65,34,7];",
			expected: []int{45, 65, 34, 7},
		},
		{
			source: `
			x := lista[45,65,34,7];
			x = x:map(|v| => texto(v));
			x;
			`,
			expected: []string{"45", "65", "34", "7"},
		},
		{
			source: `
			i := 0;
			x := lista[45,65,34,7];
			x:porCada(|v| => i += v);
			i;
			`,
			expected: 151,
		},
		{
			source: `
			x := lista[45,65,34,7, 8];
			x = x:filtrar(|v| => v % 2 == 0);
			x;
			`,
			expected: []int{34, 8},
		},
		{
			source: `
			x := lista[45,65,34,7, 8];
			x = x:contar(|v| => v % 2 == 0);
			x;
			`,
			expected: 2,
		},
		{
			source:   `lista["h", "o", "la"];`,
			expected: []string{"h", "o", "la"},
		},
		{
			source: `
			x := lista["h", "o", "la"];
			x = x:filtrar(|v| => (v == "h") || (v == "o"));
			x;
			`,
			expected: []string{"h", "o"},
		},
		{
			source:   `lista["a", "u", "r", "a"];`,
			expected: []string{"a", "u", "r", "a"},
		},
		{
			source: `
			x := "a u r a";
			x = x:separar(" ");
			x;
			`,
			expected: []string{"a", "u", "r", "a"},
		},
		{
			source: `
			x := lista["a", "u", "r", "a"];
			x = x:contar(|v| => v == 'a');
			x;
			`,
			expected: 2,
		},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		switch t := test.expected.(type) {
		case []string:
			e.testStringArrayObject(evaluated, t)

		case []int:
			e.testIntArrayObject(evaluated, t)

		default:
			e.testIntegerObject(evaluated, test.expected.(int))
		}
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
		{"!(5 > 2)", false},
		{"!(5 > 2 && 4 < 2)", true},
		{"!(5 > 2 || 4 < 2)", false},
		{"!!(5 > 2 || 4 < 2)", true},
		{"!!(5 > 2 && 4 < 2)", false},
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
		{"(1 < 2) && verdadero", true},
		{"(1 < 2) || verdadero", true},
		{"(1 < 2) && falso", false},
		{"(1 < 2) || falso", true},
		{"(1 < 2) == falso", false},
		{"(1 > 2) == verdadero", false},
		{"(1 > 2) && verdadero", false},
		{"(1 > 2) || verdadero", true},
		{"(1 > 2) && falso", false},
		{"(1 > 2) || falso", false},
		{"(1 > 2) == falso", true},
		{"(1 >= 2) == falso", true},
		{"(1 >= 2) && falso", false},
		{"(1 >= 2) || falso", false},
		{"(1 >= 2) || verdadero", true},
		{"(1 >= 2) && verdadero", false},
		{"(1 <= 2) == verdadero", true},
		{"(1 <= 2) && verdadero", true},
		{"(1 <= 2) || verdadero", true},
		{"(1 <= 2) || falso", true},
		{"(1 <= 2) && falso", false},
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
		{"si (4 > 2 && 5 < 10) { 10 } si_no { 20 }", 10},
		{"si (4 > 2 && 5 > 10) { 10 } si_no { 20 }", 20},
		{"si (4 > 2 || 5 > 10) { 10 } si_no { 20 }", 10},
		{"si (4 < 2 || 5 > 10) { 10 } si_no { 20 }", 20},
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
		{source: "5 + verdadero;", expected: "Discrepancia de tipos: entero + booleano"},
		{source: "5 + verdadero; 9;", expected: "Discrepancia de tipos: entero + booleano"},
		{source: "-verdadero", expected: "Operador desconocido: -booleano"},
		{source: "verdadero + falso", expected: "Operador desconocido: booleano + booleano"},
		{source: "5; verdadero - falso; 10;", expected: "Operador desconocido: booleano - booleano"},
		{source: `
			si (10 > 7) {
				regresa verdadero + falso;
			}
		`,
			expected: "Operador desconocido: booleano + booleano",
		},
		{source: `
			si (5 < 2) {
				regresa 1;
			} si_no {
				regresa verdadero / falso;
			}
		`,
			expected: "Operador desconocido: booleano / booleano",
		},
		{source: "foobar;", expected: "Identificador no encontrado: foobar"},
		{source: `"foo" - "bar";`, expected: "Operador desconocido: texto - texto"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.Assert().IsType(&obj.Error{}, evaluated)
		evaluatedError := evaluated.(*obj.Error)
		e.Assert().Equal(test.expected, evaluatedError.Message)
	}
}

func (e *EvaluatorTests) TestAssingmentEvaluation() {
	tests := []tuple{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"a := 5; b := a; c := a + b + 5; c;", 15},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestWhileLoop() {
	tests := []tuple{
		{`i := 0; mientras(i <= 10) { i++; }; i;`, 11},
		{`i := 0; mientras(i <= 3) { i++; }; i;`, 4},
		{`i := 0; mientras(i <= 5) { i++; }; i;`, 6},
		{`i := 0; mientras(i <= 4) { i++; }; i;`, 5},
		{`i := 0; mientras(i <= 1) { i++; }; i;`, 2},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestMaps() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{`m := mapa{"a" => 1, "b" => 2}; m["a"]`, 1},
		{`m := mapa{"a" => 1, "b" => 2}; m["b"]`, 2},
		{`m := mapa{1 => "hola", 2 => "mundo"}; m[1]`, "hola"},
		{`m := mapa{1 => "hola", 2 => " mundo"}; m[1] + m[2]`, "hola mundo"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if num, isInt := test.expected.(int); isInt {
			e.testIntegerObject(evaluated, num)
		} else {
			e.testStringObject(evaluated, test.expected.(string))
		}
	}
}

func (e *EvaluatorTests) TestMapMethods() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{`m := mapa{"a" => 1, "b" => 2}; m:contiene("a");`, true},
		{`m := mapa{"a" => 1, "b" => 2}; m:contiene("b");`, true},
		{`m := mapa{"a" => 1, "b" => 2}; m:contiene("d");`, false},
		{`m := mapa{"a" => 1, "b" => 2}; m:valores();`, []int{1, 2}},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if boolVal, isBool := test.expected.(bool); isBool {
			e.testBooleanObject(evaluated, boolVal)
		} else {
			e.testIntArrayObject(evaluated, test.expected.([]int))
		}
	}
}

func (e *EvaluatorTests) TestStringMethods() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{source: `s := "hola"; s:mayusculas();`, expected: "HOLA"},
		{source: `s := "HOLA"; s:minusculas();`, expected: "hola"},
		{source: `s := "hola"; s:contiene("g");`, expected: false},
		{source: `s := "hola"; s:contiene("h");`, expected: true},
		{source: `s := "h"; s:es_mayuscula();`, expected: false},
		{source: `s := "H"; s:es_mayuscula();`, expected: true},
		{source: `s := "h"; s:es_minuscula();`, expected: true},
		{source: `s := "H"; s:es_minuscula();`, expected: false},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if str, isStr := test.expected.(string); isStr {
			e.testStringObject(evaluated, str)
		} else {
			e.testBooleanObject(evaluated, test.expected.(bool))
		}
	}
}

func (e *EvaluatorTests) TestListMethods() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{"a := lista[2,3]; a:agregar(4); a:pop();", 4},
		{"a := lista[2,3,4,2,12]; a:agregar(17); a:pop();", 17},
		{"a := lista[2,3,4,2,12]; a:agregar(4); a:popIndice(1);", 3},
		{"a := lista[2,3,4,2,12]; a:agregar(4); a:popIndice(0);", 2},
		{"a := lista[2,3]; a:agregar(4); largo(a);", 3},
		{"a := lista[2,3]; largo(a);", 2},
		{"a := lista[2,3,4,2,12]; a:popIndice(0); largo(a);", 4},
		{"a := lista[2,3,4,2,12]; a:popIndice(0); a:contiene(3);", true},
		{"a := lista[2,3,4,12]; a:popIndice(0); a:contiene(2);", false},
		{"a := lista[2,3,4,2,12]; a:popIndice(0); a:agregar(25); a:contiene(25);", true},
		{"a := lista[1,2,3]; a = a:map(|x| => { x++; }); a[0];", 2},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if num, isNum := test.expected.(int); isNum {
			e.testIntegerObject(evaluated, num)
		} else {
			e.testBooleanObject(evaluated, test.expected.(bool))
		}
	}
}

func (e *EvaluatorTests) TestForLoop() {
	tests := []tuple{
		{`i := 0; por(n en rango(10)) { i++; }; i;`, 10},
		{`i := 0; por(n en rango(3)) { i++; }; i;`, 3},
		{`i := 0; por(n en rango(5)) { i++; }; i;`, 5},
		{`i := 0; por(n en rango(4)) { i++; }; i;`, 4},
		{`i := 0; por(n en rango(5, 10)) {i += n}; i;`, 35},
		{`i := 0; j := "hola"; por(k en j) { i++; }; i;`, 4},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestReassigment() {
	tests := []tuple{
		{"a := 5; a = 2; a;", 2},
		{`a := 20; a = 10; a;`, 10},
		{"a := 12; a = 23; a = 25; a;", 25},
		{"a := 32; a = 34; a = 5; a;", 5},
		{"a := 32; a = 34; a = 6; a;", 6},
		{`b := mapa{"a" => 1, "b" => 2}; b["a"] = 5; b["a"]`, 5},
		{`b := mapa{"a" => 1, "b" => 2}; b["a"] = 5; b["a"] = 12; b["a"]`, 12},
		{`b := mapa{"a" => 1, "b" => 2}; b["c"] = 32; b["c"];`, 32},
		{`c := lista[2,3,4]; c[0] = 10; c[0]`, 10},
		{`c := lista[2,3,4]; c[0] = 10; c[0] = 23; c[0]`, 23},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestCallList() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{"mi_lista := lista[1,23,4,5]; mi_lista[1];", 23},
		{"mi_lista := lista[1,23,4,5]; mi_lista[0];", 1},
		{"mi_lista := lista[1,23,4,5]; mi_lista[2];", 4},
		{"mi_lista := lista[1,23,4,5]; mi_lista[3];", 5},
		{"mi_lista := lista[1,23,4,5]; mi_lista[-1];", 5},
		{"mi_lista := lista[1,23,4,5]; mi_lista[-2];", 4},
		{"mi_lista := lista[1,23,4,5]; mi_lista[-3];", 23},
		{"mi_lista := lista[1,23,4,5]; mi_lista[-5];", "Indice fuera de rango indice: -5, longitud: 4"},
		{"mi_lista := lista[1,23,4,5]; mi_lista[3] + mi_lista[0];", 6},
		{"mi_lista := lista[1,23,4,5]; mi_lista[1] + mi_lista[2];", 27},
		{"mi_lista := lista[1,23,4,5]; mi_lista[100];", "Indice fuera de rango indice: 100, longitud: 4"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if expected, isNum := test.expected.(int); isNum {
			e.testIntegerObject(evaluated, expected)
		} else {
			e.testErrorObject(evaluated, test.expected.(string))
		}
	}
}

func (e *EvaluatorTests) TestOperators() {
	tests := []tuple{
		{"a := 10; a++; a;", 11},
		{"a := 10; a+=1; a;", 11},
		{"a := 10; a--; a;", 9},
		{"a := 10; a-=1; a;", 9},
		{"a := 10; a/=2; a;", 5},
		{"a := 10; a*=2; a;", 20},
		{"a := 10; a**; a;", 100},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestClassEvaluation() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{source: `
			clase Persona(name, age) {
				saludar() {
					regresa formatear("Hi im {} and i have {} years old", name, age);
				}
			}

			p := nuevo Persona("joe", 28);
			p.saludar();
		`,
			expected: "Hi im joe and i have 28 years old",
		},
		{source: `
			clase Persona(name, age) {
				saludar() {
					regresa formatear("Hi im {} and i have {} years old", name, age);
				}
			}

			p := nuevo Persona("joe", 28);
			p.name;
		`,
			expected: "joe",
		},
		{source: `
			clase Persona(name, age) {
				saludar() {
					regresa formatear("Hi im {} and i have {} years old", name, age);
				}
			}

			p := nuevo Persona("joe", 28);
			p.age;
		`,
			expected: 28,
		},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if str, isStr := test.expected.(string); isStr {
			e.testStringObject(evaluated, str)
		} else {
			e.testIntegerObject(evaluated, test.expected.(int))
		}
	}
}

func (e *EvaluatorTests) TestBuiltinFunctions() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{source: `largo("");`, expected: 0},
		{source: `largo("cuatro");`, expected: 6},
		{source: `largo("hola mundo");`, expected: 10},
		{source: "largo(1);", expected: "argumento para largo no valido, se recibio entero"},
		{
			source:   `largo("uno", "dos");`,
			expected: "numero incorrecto de argumentos para largo, se recibieron 2, se requieren 1",
		},
		{source: "tipo(1);", expected: "entero"},
		{source: "tipo(verdadero)", expected: "booleano"},
		{source: `tipo("hello world")`, expected: "texto"},
		{source: `entero("1")`, expected: 1},
		{source: `entero("hola")`, expected: "No se puede parsear como entero hola"},
		{source: "texto(12)", expected: "12"},
		{source: "a := lista[2,3,4]; largo(a);", expected: 3},
		{source: "a := lista[2,3,4]; tipo(a);", expected: "lista"},
		{source: `b := mapa{"a" => 2}; largo(b);`, expected: 1},
		{source: `b := mapa{"a" => 2}; tipo(b);`, expected: "mapa"},
		{source: `b := mapa{"a" => 2}; tipo(b["a"]);`, expected: "entero"},
		{source: `a := "aura"; a:contiene("a");`, expected: true},
		{source: `a := "aura"; a:contiene("au");`, expected: true},
		{source: `a := "aura"; a:contiene("ra");`, expected: true},
		{source: `a := "aura"; a:contiene("ara");`, expected: false},
		{source: `a := "aura"; a:contiene("b");`, expected: false},
		{source: `a := lista[2,34,45]; a:contiene(2);`, expected: true},
		{source: `a := lista[2,34,45]; a:contiene(342);`, expected: false},
		{source: `formatear("hola soy {}", "dario");`, expected: "hola soy dario"},
		{source: `formatear("hi im {}, i am {} years old", "dario", 19);`, expected: "hi im dario, i am 19 years old"},
		{source: `abs(-1)`, expected: 1},
		{source: "abs(1)", expected: 1},
		{source: "abs(-2.5)", expected: 2.5},
		{source: "abs(2.5)", expected: 2.5},
		{source: `abs("h")`, expected: "argumento para abs no valido, se recibio texto"},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		switch val := test.expected.(type) {
		case int:
			e.testIntegerObject(evaluated, val)

		case string:
			if str, isStr := evaluated.(*obj.String); isStr {
				e.testStringObject(str, val)
			} else {
				e.testErrorObject(evaluated, val)
			}

		case bool:
			e.testBooleanObject(evaluated, val)

		case float64:
			e.testFloatObject(evaluated, val)

		default:
			panic("Some weird shit happend!")
		}
	}
}

func (e *EvaluatorTests) TestImportStatement() {
	source := `
		importar "../examples/func.aura"
		result := sum(2,3)
		result
	`

	evaluated := e.evaluateTests(source)
	e.testIntegerObject(evaluated, 5)
}

func (e *EvaluatorTests) testErrorObject(evlauated obj.Object, expected string) {
	if !e.IsType(&obj.Error{}, evlauated) {
		e.T().FailNow()
	}

	err := evlauated.(*obj.Error)
	e.Equal(expected, err.Message)
}

func (e *EvaluatorTests) TestFunctionEvaluation() {
	source := "funcion(x) { x + 2; };"
	evaluated := e.evaluateTests(source)
	if !e.IsType(&obj.Def{}, evaluated) {
		e.T().FailNow()
	}

	function := evaluated.(*obj.Def)
	e.Equal(1, len(function.Parameters))
	e.Equal("x", function.Parameters[0].Str())
	e.Equal("(x + 2)", function.Body.Str())
}

func (e *EvaluatorTests) TestFunctionCalls() {
	tests := []tuple{
		{"funcion identidad(x) { x }; identidad(5);", 5},
		{`
			funcion identidad(x) {
				regresa x;
			};

			identidad(5);
		`, 5,
		},
		{`
			doble := |x| => {
				regresa 2 * x;
			};

			doble(5);
		`, 10,
		},
		{`
			funcion suma(x, y) => x + y;
			suma(3, 8);
		`, 11,
		},
		{`
			funcion suma(x, y) => x + y;
			suma(5 + 5, suma(10, 10));
		`, 30,
		},
		{"funcion(x) { x }(5)", 5},
		{"|x| => {x}(5)", 5},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)
	}
}

func (e *EvaluatorTests) TestTryExcept() {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{
			source: `
			funcion parse(x) {
				intentar {
					x := entero(x);
					regresa x;
				} excepto(e) {
					regresa "no se pudo parsear"
				}
			}
			parse("5");
			`,
			expected: 5,
		},
		{
			source: `
			funcion parse(x) {
				intentar {
					x := entero(x);
					regresa x;
				} excepto(e) {
					regresa "no se pudo parsear"
				}
			}
			parse("25")
			`,
			expected: 25,
		},
		{
			source: `
			funcion parse(x) {
				intentar {
					x := entero(x);
					regresa x;
				} excepto(e) {
					regresa "no se pudo parsear"
				}

				
			}
			parse("h")
			`,
			expected: "no se pudo parsear",
		},
		{
			source: `
			funcion parse(x) {
				intentar {
					x := entero(x);
					regresa x;
				} excepto(e) {
					regresa "no se pudo parsear"
				}
			}
			parse("hello world")
			`,
			expected: "no se pudo parsear",
		},
		{
			source: `
			funcion parse(x) {
				intentar {
					x := entero(x);
					regresa x;
				} excepto(e) {
					regresa formatear("Error de parseo: {}", e);
				}
			}
			parse("hello world")
			`,
			expected: "Error de parseo: Error: No se puede parsear como entero hello world",
		},
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		if str, isStr := test.expected.(string); isStr {
			e.testStringObject(evaluated, str)
		} else {
			e.testIntegerObject(evaluated, test.expected.(int))
		}
	}
}

func (e *EvaluatorTests) TestStringEvaluation() {
	tests := []struct {
		source   string
		expected string
	}{
		{source: `"hello world!"`, expected: "hello world!"},
		{
			source:   `funcion() { regresa "aura is awesome"; }()`,
			expected: "aura is awesome",
		},
		{source: `a := "hello"; a += " world!"; a;`, expected: "hello world!"},
		{
			source:   `a := "G"; por(i en rango(3)) { a += "o"; }; a += "al!"; a;`,
			expected: "Goooal!",
		},
		{source: `a := ""; por(i en "Hello aura") { a += i }; a;`, expected: "Hello aura"},
	}

	for _, test := range tests {
		evluated := e.evaluateTests(test.source)
		e.testStringObject(evluated, test.expected)
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
			funcion saludo(nombre) => "Hola " + nombre + "!";
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

func (e *EvaluatorTests) testStringObject(evaluated obj.Object, expected string) {
	if !e.IsType(&obj.String{}, evaluated) {
		e.T().FailNow()
	}

	str := evaluated.(*obj.String)
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

func (e *EvaluatorTests) testNullObject(eval obj.Object) {
	e.Assert().Equal(obj.SingletonNUll, eval)
}

func (e *EvaluatorTests) evaluateTests(source string) obj.Object {
	lexer := l.NewLexer(source)
	parser := p.NewParser(lexer)
	program := parser.ParseProgam()
	env := obj.NewEnviroment(nil)
	evaluated := evaluator.Evaluate(program, env)
	e.Assert().NotNil(evaluated)
	return evaluated
}

func (e *EvaluatorTests) testBooleanObject(object obj.Object, expected bool) {
	if !e.Assert().IsType(&obj.Bool{}, object) {
		e.T().FailNow()
	}

	evaluated := object.(*obj.Bool)
	e.Assert().Equal(expected, evaluated.Value)
}

func (e *EvaluatorTests) testIntArrayObject(object obj.Object, expected []int) {
	if !e.Assert().IsType(&obj.List{}, object) {
		e.T().FailNow()
	}

	evaluated := object.(*obj.List)
	e.Assert().Equal(len(expected), len(evaluated.Values))

	objList := evaluated.Values
	for i := 0; i < len(expected); i++ {
		if !e.Assert().IsType(&obj.Number{}, objList[i]) {
			e.T().FailNow()
		}

		e.Assert().Equal(expected[i], objList[i].(*obj.Number).Value)
	}
}

func (e *EvaluatorTests) testStringArrayObject(object obj.Object, expected []string) {
	if !e.Assert().IsType(&obj.List{}, object) {
		e.T().FailNow()
	}
	evaluated := object.(*obj.List)
	e.Assert().Equal(len(expected), len(evaluated.Values))

	objList := evaluated.Values
	for i := 0; i < len(expected); i++ {
		if !e.Assert().IsType(&obj.String{}, objList[i]) {
			e.T().FailNow()
		}

		e.Assert().Equal(expected[i], objList[i].(*obj.String).Value)
	}
}

func (e *EvaluatorTests) testIntegerObject(evaluated obj.Object, expected int) {
	if !e.Assert().IsType(&obj.Number{}, evaluated) {
		e.T().FailNow()
	}

	eval := evaluated.(*obj.Number)
	e.Assert().Equal(expected, eval.Value)
}

func (e *EvaluatorTests) testFloatObject(evaluated obj.Object, expected float64) {
	if !e.Assert().IsType(&obj.Float{}, evaluated) {
		e.T().FailNow()
	}

	val := evaluated.(*obj.Float)
	e.Assert().Equal(expected, val.Value)
}

func TestEvalutorSuite(t *testing.T) {
	suite.Run(t, new(EvaluatorTests))
}
