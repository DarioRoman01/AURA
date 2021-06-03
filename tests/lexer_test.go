package test_test

import (
	"katan/src"
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func loadTokens(rangee int, source string) []src.Token {
	lexer := src.NewLexer(source)
	var tokens []src.Token

	for i := 0; i < rangee; i++ {
		tokens = append(tokens, lexer.NextToken())
	}

	return tokens
}

func TestIllegalToken(t *testing.T) {
	source := "¡¿@|&"
	tokens := loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.ILLEGAL, Literal: "¡"},
		{Token_type: src.ILLEGAL, Literal: "¿"},
		{Token_type: src.ILLEGAL, Literal: "@"},
		{Token_type: src.ILLEGAL, Literal: "|"},
		{Token_type: src.ILLEGAL, Literal: "&"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in illegal tokens")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestOneCharacterOperator(t *testing.T) {
	source := "+=-/*<>!%"
	tokens := loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.PLUS, Literal: "+"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.MINUS, Literal: "-"},
		{Token_type: src.DIVISION, Literal: "/"},
		{Token_type: src.TIMES, Literal: "*"},
		{Token_type: src.LT, Literal: "<"},
		{Token_type: src.GT, Literal: ">"},
		{Token_type: src.NOT, Literal: "!"},
		{Token_type: src.MOD, Literal: "%"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal on character operator")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestEOF(t *testing.T) {
	source := "+"
	tokens := loadTokens(utf8.RuneCountInString(source)+1, source)

	expectedTokens := []src.Token{
		{Token_type: src.PLUS, Literal: "+"},
		{Token_type: src.EOF, Literal: ""},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in EOF")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestDilimiters(t *testing.T) {
	source := "(){},;"
	tokens := loadTokens(len(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.LPAREN, Literal: "("},
		{Token_type: src.RPAREN, Literal: ")"},
		{Token_type: src.LBRACE, Literal: "{"},
		{Token_type: src.RBRACE, Literal: "}"},
		{Token_type: src.COMMA, Literal: ","},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in Dilimiters")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestAssingments(t *testing.T) {
	source := "var cinco =  5;"
	tokens := loadTokens(5, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "cinco"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.INT, Literal: "5"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in assingment")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestFunctionDeclaration(t *testing.T) {
	source := `
		var suma = funcion(x, y) {
			x + y;
		};
	`

	tokens := loadTokens(16, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "suma"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.FUNCTION, Literal: "funcion"},
		{Token_type: src.LPAREN, Literal: "("},
		{Token_type: src.IDENT, Literal: "x"},
		{Token_type: src.COMMA, Literal: ","},
		{Token_type: src.IDENT, Literal: "y"},
		{Token_type: src.RPAREN, Literal: ")"},
		{Token_type: src.LBRACE, Literal: "{"},
		{Token_type: src.IDENT, Literal: "x"},
		{Token_type: src.PLUS, Literal: "+"},
		{Token_type: src.IDENT, Literal: "y"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.RBRACE, Literal: "}"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in function declarion")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestFunctionCall(t *testing.T) {
	source := "var resultado = suma(dos, tres);"
	tokens := loadTokens(10, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "resultado"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.IDENT, Literal: "suma"},
		{Token_type: src.LPAREN, Literal: "("},
		{Token_type: src.IDENT, Literal: "dos"},
		{Token_type: src.COMMA, Literal: ","},
		{Token_type: src.IDENT, Literal: "tres"},
		{Token_type: src.RPAREN, Literal: ")"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in function call")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestControllStatement(t *testing.T) {
	source := `
		si (5 < 10) {
			regresa verdadero;
		} si_no {
			regresa falso;
		}
	`

	tokens := loadTokens(17, source)

	expectedTokens := []src.Token{
		{Token_type: src.IF, Literal: "si"},
		{Token_type: src.LPAREN, Literal: "("},
		{Token_type: src.INT, Literal: "5"},
		{Token_type: src.LT, Literal: "<"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.RPAREN, Literal: ")"},
		{Token_type: src.LBRACE, Literal: "{"},
		{Token_type: src.RETURN, Literal: "regresa"},
		{Token_type: src.TRUE, Literal: "verdadero"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.RBRACE, Literal: "}"},
		{Token_type: src.ELSE, Literal: "si_no"},
		{Token_type: src.LBRACE, Literal: "{"},
		{Token_type: src.RETURN, Literal: "regresa"},
		{Token_type: src.FALSE, Literal: "falso"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.RBRACE, Literal: "}"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in controll statement")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestTwoCharacterOperator(t *testing.T) {
	source := `
		10 == 10;
		10 != 9;
		10 >= 9;
		10 <= 9;
		10 || 9;
		10 && 9;
	`

	tokens := loadTokens(24, source)

	expectedTokens := []src.Token{
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.EQ, Literal: "=="},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.NOT_EQ, Literal: "!="},
		{Token_type: src.INT, Literal: "9"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.GTOREQ, Literal: ">="},
		{Token_type: src.INT, Literal: "9"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.LTOREQ, Literal: "<="},
		{Token_type: src.INT, Literal: "9"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.OR, Literal: "||"},
		{Token_type: src.INT, Literal: "9"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.AND, Literal: "&&"},
		{Token_type: src.INT, Literal: "9"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in two character operator")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestUnderscoreVar(t *testing.T) {
	source := "var num_12 = 5;"

	tokens := loadTokens(5, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "num_12"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.INT, Literal: "5"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in underscore variable")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestString(t *testing.T) {
	source := `
		"foo";
		"src is a great programing lenguage";
	`

	tokens := loadTokens(4, source)
	expectedTokens := []src.Token{
		{Token_type: src.STRING, Literal: "foo"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.STRING, Literal: "src is a great programing lenguage"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	assert.Equal(t, expectedTokens, tokens)
}
