package test_test

import (
	"aura/src/ast"
	l "aura/src/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetStatement(t *testing.T) {
	program := ast.NewProgram([]ast.Stmt{
		ast.LetStatement{
			Token: l.Token{
				Token_type: l.LET,
				Literal:    "var",
			},
			Name: ast.NewIdentifier(l.Token{
				Token_type: l.IDENT,
				Literal:    "mi_var",
			}, "mi_var"),
			Value: ast.NewIdentifier(l.Token{
				Token_type: l.IDENT,
				Literal:    "otra_var",
			}, "otra_var"),
		},
	})

	assert.Equal(t, "var mi_var = otra_var;", program.Str())
}

func TestReturnStatements(t *testing.T) {
	program := ast.NewProgram([]ast.Stmt{
		ast.LetStatement{
			Token: l.Token{
				Token_type: l.LET,
				Literal:    "var",
			},
			Name: ast.NewIdentifier(l.Token{
				Token_type: l.IDENT,
				Literal:    "x",
			}, "x"),
			Value: ast.NewIdentifier(l.Token{
				Token_type: l.INT,
				Literal:    "5",
			}, "5"),
		},
		ast.ReturnStament{
			Token: l.Token{
				Token_type: l.RETURN,
				Literal:    "regresa",
			},
			ReturnValue: ast.NewIdentifier(
				l.Token{
					Token_type: l.IDENT,
					Literal:    "x",
				}, "x"),
		},
	})

	assert.Equal(t, "var x = 5; regresa x;", program.Str())
}

func TestIntegerExpression(t *testing.T) {
	value := 5
	program := ast.NewProgram([]ast.Stmt{
		ast.LetStatement{
			Token: l.Token{
				Token_type: l.LET,
				Literal:    "var",
			},
			Name: ast.NewIdentifier(l.Token{
				Token_type: l.IDENT,
				Literal:    "x",
			}, "x"),
			Value: ast.NewInteger(l.Token{
				Token_type: l.INT,
				Literal:    "5",
			}, &value),
		},
	})

	assert.Equal(t, "var x = 5;", program.Str())
}
