package test_test

import (
	"lpp/lpp"
	"reflect"
	"testing"
	"unicode/utf8"
)

func loadTokens(rangee int, source string) []lpp.Token {
	lexer := lpp.NewLexer(source)
	var tokens []lpp.Token

	for i := 0; i < rangee; i++ {
		tokens = append(tokens, lexer.NextToken())
	}

	return tokens
}

func TestIllegalToken(t *testing.T) {
	source := "¡¿@"
	tokens := loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []lpp.Token{
		{Token_type: lpp.ILLEGAL, Literal: "¡"},
		{Token_type: lpp.ILLEGAL, Literal: "¿"},
		{Token_type: lpp.ILLEGAL, Literal: "@"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in illegal tokens")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}

func TestOneCharacterOperator(t *testing.T) {
	source := "+=-/*<>!"
	tokens := loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []lpp.Token{
		{Token_type: lpp.PLUS, Literal: "+"},
		{Token_type: lpp.ASSING, Literal: "="},
		{Token_type: lpp.MINUS, Literal: "-"},
		{Token_type: lpp.DIVISION, Literal: "/"},
		{Token_type: lpp.TIMES, Literal: "*"},
		{Token_type: lpp.LT, Literal: "<"},
		{Token_type: lpp.GT, Literal: ">"},
		{Token_type: lpp.NOT, Literal: "!"},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.PLUS, Literal: "+"},
		{Token_type: lpp.EOF, Literal: ""},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.LPAREN, Literal: "("},
		{Token_type: lpp.RPAREN, Literal: ")"},
		{Token_type: lpp.LBRACE, Literal: "{"},
		{Token_type: lpp.RBRACE, Literal: "}"},
		{Token_type: lpp.COMMA, Literal: ","},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.LET, Literal: "var"},
		{Token_type: lpp.IDENT, Literal: "cinco"},
		{Token_type: lpp.ASSING, Literal: "="},
		{Token_type: lpp.INT, Literal: "5"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.LET, Literal: "var"},
		{Token_type: lpp.IDENT, Literal: "suma"},
		{Token_type: lpp.ASSING, Literal: "="},
		{Token_type: lpp.FUNCTION, Literal: "funcion"},
		{Token_type: lpp.LPAREN, Literal: "("},
		{Token_type: lpp.IDENT, Literal: "x"},
		{Token_type: lpp.COMMA, Literal: ","},
		{Token_type: lpp.IDENT, Literal: "y"},
		{Token_type: lpp.RPAREN, Literal: ")"},
		{Token_type: lpp.LBRACE, Literal: "{"},
		{Token_type: lpp.IDENT, Literal: "x"},
		{Token_type: lpp.PLUS, Literal: "+"},
		{Token_type: lpp.IDENT, Literal: "y"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
		{Token_type: lpp.RBRACE, Literal: "}"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.LET, Literal: "var"},
		{Token_type: lpp.IDENT, Literal: "resultado"},
		{Token_type: lpp.ASSING, Literal: "="},
		{Token_type: lpp.IDENT, Literal: "suma"},
		{Token_type: lpp.LPAREN, Literal: "("},
		{Token_type: lpp.IDENT, Literal: "dos"},
		{Token_type: lpp.COMMA, Literal: ","},
		{Token_type: lpp.IDENT, Literal: "tres"},
		{Token_type: lpp.RPAREN, Literal: ")"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lpp.Token{
		{Token_type: lpp.IF, Literal: "si"},
		{Token_type: lpp.LPAREN, Literal: "("},
		{Token_type: lpp.INT, Literal: "5"},
		{Token_type: lpp.LT, Literal: "<"},
		{Token_type: lpp.INT, Literal: "10"},
		{Token_type: lpp.RPAREN, Literal: ")"},
		{Token_type: lpp.LBRACE, Literal: "{"},
		{Token_type: lpp.RETURN, Literal: "regresa"},
		{Token_type: lpp.TRUE, Literal: "verdadero"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
		{Token_type: lpp.RBRACE, Literal: "}"},
		{Token_type: lpp.ELSE, Literal: "si_no"},
		{Token_type: lpp.LBRACE, Literal: "{"},
		{Token_type: lpp.RETURN, Literal: "regresa"},
		{Token_type: lpp.FALSE, Literal: "falso"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
		{Token_type: lpp.RBRACE, Literal: "}"},
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
	`

	tokens := loadTokens(8, source)

	expectedTokens := []lpp.Token{
		{Token_type: lpp.INT, Literal: "10"},
		{Token_type: lpp.EQ, Literal: "=="},
		{Token_type: lpp.INT, Literal: "10"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
		{Token_type: lpp.INT, Literal: "10"},
		{Token_type: lpp.NOT_EQ, Literal: "!="},
		{Token_type: lpp.INT, Literal: "9"},
		{Token_type: lpp.SEMICOLON, Literal: ";"},
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Log("tokens and expected tokens are not equal in two character operator")
		t.Log("expected: ", expectedTokens)
		t.Log("found: ", tokens)
		t.Fail()
	}
}
