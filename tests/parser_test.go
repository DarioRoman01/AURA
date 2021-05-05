package test_test

import (
	"lpp/lpp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProgram(t *testing.T) {
	assert := assert.New(t)
	source := "var x = 5;"
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()

	if !assert.NotNil(program) {
		t.Log("program is nil")
		t.Fail()
	}

	if !assert.IsType(lpp.Program{}, program) {
		t.Log("program type is not Program!")
		t.Fail()
	}

	if !assert.Implements((*lpp.ASTNode)(nil), program) {
		t.Log("program type does not implement ast node interface")
		t.Fail()
	}
}
