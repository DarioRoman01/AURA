package lexer

import "fmt"

// Represents all the posibles tokens that can exist
type TokenType int

const (
	Head TokenType = iota
	AND
	ASSING
	ARROW
	BAR
	CLASS
	COLON
	COLONASSING
	COMMA
	DATASTRCUT
	DIVISION
	DIVASSING
	DOT
	ELSE
	EOF
	EQ
	EXCEPT
	EXPONENT
	FALSE
	FLOAT
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
	MINUS2
	MINUSASSING
	NEW
	NOT // !
	NOT_EQ
	MOD
	OR
	PLUS
	PLUS2
	PLUSASSING
	RBRACE
	RBRACKET
	RETURN
	RPAREN
	SEMICOLON
	THIS
	TIMES // *
	TIMEASSI
	TRY
	STRING
	TRUE
	WHILE
	NULLT
	MAP
	IMPORT
)

// String representation of all tokens
var Tokens = [...]string{
	AND:         "&&",
	ASSING:      "=",
	ARROW:       "=>",
	BAR:         "|",
	CLASS:       "clase",
	COLON:       ":",
	COLONASSING: ":=",
	COMMA:       ",",
	DIVISION:    "/",
	DOT:         ".",
	ELSE:        "si_no",
	EOF:         "final del archivo",
	EQ:          "==",
	EXCEPT:      "excepto",
	FALSE:       "falso",
	FLOAT:       "flotante",
	FUNCTION:    "funcion",
	GT:          ">",
	IDENT:       "identificador",
	IF:          "si",
	ILLEGAL:     "ilegal",
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
	THIS:        "this",
	TRY:         "intenta",
	TIMEASSI:    "*=",
	DIVASSING:   "/=",
	EXPONENT:    "**",
	IMPORT:      "importar",
}

// Represents a Token in the programmig lenguage
type Token struct {
	Token_type TokenType // represents the type of the token
	Literal    string    // represents the literal of the token
}

// Generate a new Token instance
func NewToken(t TokenType, literal string) Token {
	return Token{Token_type: t, Literal: literal}
}

// print token info
func (t *Token) PrintToken() string {
	return fmt.Sprintf("Token Type: %s, Literal: %s", Tokens[t.Token_type], t.Literal)
}

// verify that given literal is a keyword or not
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
		"this":      THIS,
		"clase":     CLASS,
		"nuevo":     NEW,
		"importar":  IMPORT,
		"intentar":  TRY,
		"excepto":   EXCEPT,
	}

	if TokenType, exists := keywords[literal]; exists {
		return TokenType
	}

	return IDENT
}
