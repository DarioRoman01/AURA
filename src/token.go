package src

import "fmt"

type TokenType int

const (
	Head TokenType = iota
	AND
	ASSING
	ARROW
	COLON
	COMMA
	DATASTRCUT
	DIVISION
	DIVASSING
	ELSE
	EOF
	EQ
	EXPONENT
	FALSE
	FOR
	FUNCTION
	GT     // grather than
	GTOREQ // grater than or equeal
	IDENT
	IF
	ILLEGAL
	IN
	INT
	LBRACE
	LBRACKET
	LET
	LPAREN
	LT     // less than
	LTOREQ // less than or equal
	MINUS  // -
	MINUSASSING
	NOT // !
	NOT_EQ
	MOD
	OR
	PLUS
	PLUSASSING
	RBRACE
	RBRACKET
	RETURN
	RPAREN
	SEMICOLON
	TIMES // *
	TIMEASSI
	STRING
	TRUE
	WHILE
	NULLT
	MAP
)

var tokens = [...]string{
	AND:         "&&",
	ASSING:      "=",
	COLON:       ":",
	COMMA:       ",",
	DIVISION:    "/",
	ELSE:        "si_no",
	EOF:         "EOF",
	EQ:          "==",
	FALSE:       "falso",
	FUNCTION:    "funcion",
	GT:          ">",
	IDENT:       "IDENT",
	IF:          "si",
	ILLEGAL:     "ILLEGAL",
	IN:          "en",
	INT:         "INT",
	LBRACE:      "{",
	LET:         "var",
	LPAREN:      "(",
	LT:          "<",
	MINUS:       "-",
	NOT:         "!",
	NOT_EQ:      "!=",
	MOD:         "%",
	OR:          "||",
	PLUS:        "+",
	RBRACE:      "}",
	RETURN:      "regresa",
	RPAREN:      ")",
	RBRACKET:    "]",
	LBRACKET:    "[",
	SEMICOLON:   ";",
	TIMES:       "*",
	STRING:      `"`,
	TRUE:        "verdaro",
	WHILE:       "mientras",
	NULLT:       "nulo",
	MAP:         "mapa",
	PLUSASSING:  "+=",
	MINUSASSING: "-=",
	TIMEASSI:    "*=",
	DIVASSING:   "/=",
	EXPONENT:    "**",
}

type Token struct {
	Token_type TokenType
	Literal    string
}

func (t *Token) PrintToken() string {
	return fmt.Sprintf("Token Type: %s, Literal: %s", tokens[t.Token_type], t.Literal)
}

// verify that given literal is a string
func LookUpTokenType(literal string) TokenType {
	keywords := map[string]TokenType{
		"falso":     FALSE,
		"funcion":   FUNCTION,
		"regresa":   RETURN,
		"si":        IF,
		"si_no":     ELSE,
		"var":       LET,
		"verdadero": TRUE,
		"en":        IN,
		"mientras":  WHILE,
		"por":       FOR,
		"lista":     DATASTRCUT,
		"nulo":      NULLT,
		"mapa":      MAP,
	}

	if TokenType, exists := keywords[literal]; exists {
		return TokenType
	}

	return IDENT
}
