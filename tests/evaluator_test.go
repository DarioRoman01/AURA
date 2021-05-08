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
	}

	for _, test := range tests {
		evaluated := e.evaluateTests(test.source)
		e.testIntegerObject(evaluated, test.expected)

	}
}

func (e *EvaluatorTests) evaluateTests(source string) lpp.Object {
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()
	evaluated := lpp.Evaluate(program)
	e.Assert().NotNil(evaluated)
	return evaluated
}

func (e *EvaluatorTests) testIntegerObject(evaluated lpp.Object, expected int) {
	e.Assert().IsType(&lpp.Number{}, evaluated.(*lpp.Number))
	eval := evaluated.(*lpp.Number)
	e.Assert().Equal(expected, eval.Value)
}

func TestEvalutorSuite(t *testing.T) {
	suite.Run(t, new(EvaluatorTests))
}
