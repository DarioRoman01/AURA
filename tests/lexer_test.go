package test_test

import (
	"katan/src"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/suite"
)

type LexerTests struct {
	suite.Suite
}

func (l *LexerTests) loadTokens(length int, source string) []src.Token {
	lexer := src.NewLexer(source)
	var tokens []src.Token

	for i := 0; i < length; i++ {
		tokens = append(tokens, lexer.NextToken())
	}

	return tokens
}

func (l *LexerTests) TestIllegalToken() {
	source := "¡¿@|&"
	tokens := l.loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.ILLEGAL, Literal: "¡"},
		{Token_type: src.ILLEGAL, Literal: "¿"},
		{Token_type: src.ILLEGAL, Literal: "@"},
		{Token_type: src.ILLEGAL, Literal: "|"},
		{Token_type: src.ILLEGAL, Literal: "&"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestOneCharacterOperator() {
	source := "+-/*<>!%="
	tokens := l.loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.PLUS, Literal: "+"},
		{Token_type: src.MINUS, Literal: "-"},
		{Token_type: src.DIVISION, Literal: "/"},
		{Token_type: src.TIMES, Literal: "*"},
		{Token_type: src.LT, Literal: "<"},
		{Token_type: src.GT, Literal: ">"},
		{Token_type: src.NOT, Literal: "!"},
		{Token_type: src.MOD, Literal: "%"},
		{Token_type: src.ASSING, Literal: "="},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestEOF() {
	source := "+"
	tokens := l.loadTokens(utf8.RuneCountInString(source)+1, source)

	expectedTokens := []src.Token{
		{Token_type: src.PLUS, Literal: "+"},
		{Token_type: src.EOF, Literal: ""},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestDilimiters() {
	source := "(){},;"
	tokens := l.loadTokens(len(source), source)

	expectedTokens := []src.Token{
		{Token_type: src.LPAREN, Literal: "("},
		{Token_type: src.RPAREN, Literal: ")"},
		{Token_type: src.LBRACE, Literal: "{"},
		{Token_type: src.RBRACE, Literal: "}"},
		{Token_type: src.COMMA, Literal: ","},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestAssingments() {
	source := "var cinco =  5;"
	tokens := l.loadTokens(5, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "cinco"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.INT, Literal: "5"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestFunctionDeclaration() {
	source := `
		var suma = funcion(x, y) {
			x + y;
		};
	`

	tokens := l.loadTokens(16, source)

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

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestFunctionCall() {
	source := "var resultado = suma(dos, tres);"
	tokens := l.loadTokens(10, source)

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

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestControllStatement() {
	source := `
		si (5 < 10) {
			regresa verdadero;
		} si_no {
			regresa falso;
		}
	`

	tokens := l.loadTokens(17, source)

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

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestTwoCharacterOperator() {
	source := `
		10 == 10;
		10 != 9;
		10 >= 9;
		10 <= 9;
		10 || 9;
		10 && 9;
		10++;
		10 += 1;
		10--;
		10 -= 1;
		10 *= 2;
		10 /= 2;
	`

	tokens := l.loadTokens(46, source)

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
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.PLUS2, Literal: "++"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.PLUSASSING, Literal: "+="},
		{Token_type: src.INT, Literal: "1"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.MINUS2, Literal: "--"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.MINUSASSING, Literal: "-="},
		{Token_type: src.INT, Literal: "1"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.TIMEASSI, Literal: "*="},
		{Token_type: src.INT, Literal: "2"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.INT, Literal: "10"},
		{Token_type: src.DIVASSING, Literal: "/="},
		{Token_type: src.INT, Literal: "2"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestUnderscoreVar() {
	source := "var num_12 = 5;"

	tokens := l.loadTokens(5, source)

	expectedTokens := []src.Token{
		{Token_type: src.LET, Literal: "var"},
		{Token_type: src.IDENT, Literal: "num_12"},
		{Token_type: src.ASSING, Literal: "="},
		{Token_type: src.INT, Literal: "5"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestString() {
	source := `
		"foo";
		"src is a great programing lenguage";
	`

	tokens := l.loadTokens(4, source)
	expectedTokens := []src.Token{
		{Token_type: src.STRING, Literal: "foo"},
		{Token_type: src.SEMICOLON, Literal: ";"},
		{Token_type: src.STRING, Literal: "src is a great programing lenguage"},
		{Token_type: src.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func TestLexerSuite(t *testing.T) {
	suite.Run(t, new(LexerTests))
}
