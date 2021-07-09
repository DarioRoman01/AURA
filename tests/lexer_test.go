package test

import (
	"aura/src/lexer"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/suite"
)

type LexerTests struct {
	suite.Suite
}

func (l *LexerTests) loadTokens(length int, source string) []lexer.Token {
	lex := lexer.NewLexer(source)
	var tokens []lexer.Token

	for i := 0; i < length; i++ {
		tokens = append(tokens, lex.NextToken())
	}

	return tokens
}

func (l *LexerTests) TestIllegalToken() {
	source := "¡¿@&"
	tokens := l.loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.ILLEGAL, Literal: "¡"},
		{Token_type: lexer.ILLEGAL, Literal: "¿"},
		{Token_type: lexer.ILLEGAL, Literal: "@"},
		{Token_type: lexer.ILLEGAL, Literal: "&"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestOneCharacterOperator() {
	source := "+-/*<>!%="
	tokens := l.loadTokens(utf8.RuneCountInString(source), source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.PLUS, Literal: "+"},
		{Token_type: lexer.MINUS, Literal: "-"},
		{Token_type: lexer.DIVISION, Literal: "/"},
		{Token_type: lexer.TIMES, Literal: "*"},
		{Token_type: lexer.LT, Literal: "<"},
		{Token_type: lexer.GT, Literal: ">"},
		{Token_type: lexer.NOT, Literal: "!"},
		{Token_type: lexer.MOD, Literal: "%"},
		{Token_type: lexer.ASSING, Literal: "="},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestFloat() {
	source := "5.5"

	tokens := l.loadTokens(utf8.RuneCountInString(source)-1, source)
	expectedTokens := []lexer.Token{
		{Token_type: lexer.FLOAT, Literal: "5.5"},
		{Token_type: lexer.EOF, Literal: ""},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestEOF() {
	source := "+"
	tokens := l.loadTokens(utf8.RuneCountInString(source)+1, source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.PLUS, Literal: "+"},
		{Token_type: lexer.EOF, Literal: ""},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestDilimiters() {
	source := "(){},;"
	tokens := l.loadTokens(len(source), source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.LPAREN, Literal: "("},
		{Token_type: lexer.RPAREN, Literal: ")"},
		{Token_type: lexer.LBRACE, Literal: "{"},
		{Token_type: lexer.RBRACE, Literal: "}"},
		{Token_type: lexer.COMMA, Literal: ","},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestAssingments() {
	source := "var cinco =  5;"
	tokens := l.loadTokens(5, source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.LET, Literal: "var"},
		{Token_type: lexer.IDENT, Literal: "cinco"},
		{Token_type: lexer.ASSING, Literal: "="},
		{Token_type: lexer.INT, Literal: "5"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lexer.Token{
		{Token_type: lexer.LET, Literal: "var"},
		{Token_type: lexer.IDENT, Literal: "suma"},
		{Token_type: lexer.ASSING, Literal: "="},
		{Token_type: lexer.FUNCTION, Literal: "funcion"},
		{Token_type: lexer.LPAREN, Literal: "("},
		{Token_type: lexer.IDENT, Literal: "x"},
		{Token_type: lexer.COMMA, Literal: ","},
		{Token_type: lexer.IDENT, Literal: "y"},
		{Token_type: lexer.RPAREN, Literal: ")"},
		{Token_type: lexer.LBRACE, Literal: "{"},
		{Token_type: lexer.IDENT, Literal: "x"},
		{Token_type: lexer.PLUS, Literal: "+"},
		{Token_type: lexer.IDENT, Literal: "y"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.RBRACE, Literal: "}"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestFunctionCall() {
	source := "var resultado = suma(dos, tres);"
	tokens := l.loadTokens(10, source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.LET, Literal: "var"},
		{Token_type: lexer.IDENT, Literal: "resultado"},
		{Token_type: lexer.ASSING, Literal: "="},
		{Token_type: lexer.IDENT, Literal: "suma"},
		{Token_type: lexer.LPAREN, Literal: "("},
		{Token_type: lexer.IDENT, Literal: "dos"},
		{Token_type: lexer.COMMA, Literal: ","},
		{Token_type: lexer.IDENT, Literal: "tres"},
		{Token_type: lexer.RPAREN, Literal: ")"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
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

	expectedTokens := []lexer.Token{
		{Token_type: lexer.IF, Literal: "si"},
		{Token_type: lexer.LPAREN, Literal: "("},
		{Token_type: lexer.INT, Literal: "5"},
		{Token_type: lexer.LT, Literal: "<"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.RPAREN, Literal: ")"},
		{Token_type: lexer.LBRACE, Literal: "{"},
		{Token_type: lexer.RETURN, Literal: "regresa"},
		{Token_type: lexer.TRUE, Literal: "verdadero"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.RBRACE, Literal: "}"},
		{Token_type: lexer.ELSE, Literal: "si_no"},
		{Token_type: lexer.LBRACE, Literal: "{"},
		{Token_type: lexer.RETURN, Literal: "regresa"},
		{Token_type: lexer.FALSE, Literal: "falso"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.RBRACE, Literal: "}"},
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

	expectedTokens := []lexer.Token{
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.EQ, Literal: "=="},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.NOT_EQ, Literal: "!="},
		{Token_type: lexer.INT, Literal: "9"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.GTOREQ, Literal: ">="},
		{Token_type: lexer.INT, Literal: "9"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.LTOREQ, Literal: "<="},
		{Token_type: lexer.INT, Literal: "9"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.OR, Literal: "||"},
		{Token_type: lexer.INT, Literal: "9"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.AND, Literal: "&&"},
		{Token_type: lexer.INT, Literal: "9"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.PLUS2, Literal: "++"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.PLUSASSING, Literal: "+="},
		{Token_type: lexer.INT, Literal: "1"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.MINUS2, Literal: "--"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.MINUSASSING, Literal: "-="},
		{Token_type: lexer.INT, Literal: "1"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.TIMEASSI, Literal: "*="},
		{Token_type: lexer.INT, Literal: "2"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.INT, Literal: "10"},
		{Token_type: lexer.DIVASSING, Literal: "/="},
		{Token_type: lexer.INT, Literal: "2"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestUnderscoreVar() {
	source := "var num_12 = 5;"

	tokens := l.loadTokens(5, source)

	expectedTokens := []lexer.Token{
		{Token_type: lexer.LET, Literal: "var"},
		{Token_type: lexer.IDENT, Literal: "num_12"},
		{Token_type: lexer.ASSING, Literal: "="},
		{Token_type: lexer.INT, Literal: "5"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func (l *LexerTests) TestString() {
	source := `
		"foo";
		"src is a great programing lenguage";
	`

	tokens := l.loadTokens(4, source)
	expectedTokens := []lexer.Token{
		{Token_type: lexer.STRING, Literal: "foo"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
		{Token_type: lexer.STRING, Literal: "src is a great programing lenguage"},
		{Token_type: lexer.SEMICOLON, Literal: ";"},
	}

	l.Assert().Equal(expectedTokens, tokens)
}

func TestLexerSuite(t *testing.T) {
	suite.Run(t, new(LexerTests))
}
