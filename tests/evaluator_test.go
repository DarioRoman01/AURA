package test_test

import (
	"lpp/lpp"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EvaluatorTests struct {
	suite.Suite
}

func (e *EvaluatorTests) TestIntegerEvaluation() {
	tests := []struct {
		source   string
		expected int
	}{
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
	tests := []struct {
		source   string
		expected int
	}{
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

func (e *EvaluatorTests) testNullObject(eval lpp.Object) {
	e.Assert().Equal(lpp.SingletonNUll, eval.(*lpp.Null))
}

func (e *EvaluatorTests) evaluateTests(source string) lpp.Object {
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()
	evaluated := lpp.Evaluate(program)
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
