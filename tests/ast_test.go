package test_test

import (
	"katan/src"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetStatement(t *testing.T) {
	program := src.NewProgram([]src.Stmt{
		src.LetStatement{
			Token: src.Token{
				Token_type: src.LET,
				Literal:    "var",
			},
			Name: src.NewIdentifier(src.Token{
				Token_type: src.IDENT,
				Literal:    "mi_var",
			}, "mi_var"),
			Value: src.NewIdentifier(src.Token{
				Token_type: src.IDENT,
				Literal:    "otra_var",
			}, "otra_var"),
		},
	})

	assert.Equal(t, "var mi_var = otra_var;", program.Str())
}

func TestReturnStatements(t *testing.T) {
	program := src.NewProgram([]src.Stmt{
		src.LetStatement{
			Token: src.Token{
				Token_type: src.LET,
				Literal:    "var",
			},
			Name: src.NewIdentifier(src.Token{
				Token_type: src.IDENT,
				Literal:    "x",
			}, "x"),
			Value: src.NewIdentifier(src.Token{
				Token_type: src.INT,
				Literal:    "5",
			}, "5"),
		},
		src.ReturnStament{
			Token: src.Token{
				Token_type: src.RETURN,
				Literal:    "regresa",
			},
			ReturnValue: src.NewIdentifier(
				src.Token{
					Token_type: src.IDENT,
					Literal:    "x",
				}, "x"),
		},
	})

	assert.Equal(t, "var x = 5; regresa x;", program.Str())
}

func TestIntegerExpression(t *testing.T) {
	value := 5
	program := src.NewProgram([]src.Stmt{
		src.LetStatement{
			Token: src.Token{
				Token_type: src.LET,
				Literal:    "var",
			},
			Name: src.NewIdentifier(src.Token{
				Token_type: src.IDENT,
				Literal:    "x",
			}, "x"),
			Value: src.NewInteger(src.Token{
				Token_type: src.INT,
				Literal:    "5",
			}, &value),
		},
	})

	assert.Equal(t, "var x = 5;", program.Str())
}
