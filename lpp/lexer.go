package lpp

import (
	"regexp"
	"unicode/utf8"
)

type Lexer struct {
	source        string
	character     string
	read_position int
	position      int
}

// create a new lexer
func NewLexer(source string) *Lexer {
	lexer := &Lexer{
		source:        source,
		character:     "",
		read_position: 0,
		position:      0,
	}

	lexer.readCharacter()
	return lexer
}

// read next token
func (l *Lexer) NextToken() Token {
	l.skipWhiteSpaces()
	var token Token

	if equal, _ := regexp.MatchString(`^=$`, l.character); equal {
		token = Token{Token_type: ASSING, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\+$`, l.character); equal {
		token = Token{Token_type: PLUS, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^$`, l.character); equal {
		token = Token{Token_type: EOF, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\($`, l.character); equal {
		token = Token{Token_type: LPAREN, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\)$`, l.character); equal {
		token = Token{Token_type: RPAREN, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\{$`, l.character); equal {
		token = Token{Token_type: LBRACE, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\}$`, l.character); equal {
		token = Token{Token_type: RBRACE, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^,$`, l.character); equal {
		token = Token{Token_type: COMMA, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^;$`, l.character); equal {
		token = Token{Token_type: SEMICOLON, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^<$`, l.character); equal {
		token = Token{Token_type: LT, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^>$`, l.character); equal {
		token = Token{Token_type: GT, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\-$`, l.character); equal {
		token = Token{Token_type: MINUS, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\/$`, l.character); equal {
		token = Token{Token_type: DIVISION, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\*$`, l.character); equal {
		token = Token{Token_type: TIMES, Literal: l.character}

	} else if equal, _ := regexp.MatchString(`^\!$`, l.character); equal {
		token = Token{Token_type: NOT, Literal: l.character}

	} else if l.isLetter(l.character) {
		literal := l.readIdentifier()
		token_type := LookUpTokenType(literal)
		return Token{Token_type: token_type, Literal: literal}

	} else if l.isNumber(l.character) {
		literal := l.readNumber()
		return Token{Token_type: INT, Literal: literal}

	} else {
		token = Token{Token_type: ILLEGAL, Literal: l.character}
	}

	l.readCharacter()
	return token
}

// check if current character is letter
func (l *Lexer) isLetter(char string) bool {
	isValid, _ := regexp.MatchString(`^[a-záéíóúA-ZÁÉÍÓÚñÑ_]$`, char)
	return isValid
}

// check if current character is number
func (l *Lexer) isNumber(char string) bool {
	isValid, _ := regexp.MatchString(`^\d$`, char)
	return isValid
}

// read current character.
func (l *Lexer) readCharacter() {
	if l.read_position >= utf8.RuneCountInString(l.source) {
		l.character = ""
	} else {
		l.character = string([]rune(l.source)[l.read_position])
	}

	l.position = l.read_position
	l.read_position++
}

// read character sequence
func (l *Lexer) readIdentifier() string {
	initialPosition := l.position
	for l.isLetter(l.character) {
		l.readCharacter()
	}

	return string(l.source[initialPosition:l.position])
}

// read number sequence of characters
func (l *Lexer) readNumber() string {
	initialPosition := l.position
	for l.isNumber(l.character) {
		l.readCharacter()
	}

	return string(l.source[initialPosition:l.position])
}

// skipp whitespaces
func (l *Lexer) skipWhiteSpaces() {
	m, _ := regexp.Compile(`^\s$`)
	for m.MatchString(l.character) {
		l.readCharacter()
	}
}
