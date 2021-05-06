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

func TestLetStatements(t *testing.T) {
	assert := assert.New(t)
	source := `
		var x = 5;
		var y = 10;
		var foo = 20;
	`
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	program := parser.ParseProgam()

	if !assert.Equal(3, len(program.Staments)) {
		t.Log("len of program statements are not 3")
		t.Fail()
	}

	if !assert.Equal("x", program.Staments[0].LetStatement.TokenLiteral()) {
		t.Fail()
	}

	if !assert.Equal("y", program.Staments[1].LetStatement.TokenLiteral()) {
		t.Fail()
	}
	if !assert.Equal("foo", program.Staments[2].LetStatement.TokenLiteral()) {
		t.Fail()
	}

	for _, statement := range program.Staments {
		if !assert.Equal("var", statement.TokenLiteral()) {
			t.Log("token are not a variable")
			t.Fail()
		}

		if !assert.IsType(&lpp.LetStatement{}, statement.LetStatement) {
			t.Log("statement are not let statement type")
			t.Fail()
		}
	}
}
