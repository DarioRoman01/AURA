package lexer

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

var (
	charRegex   = regexp.MustCompile(`^[a-záéíóúA-ZÁÉÍÓÚñÑ_]$`)
	numRegex    = regexp.MustCompile(`^\d$`)
	wSpaceRegex = regexp.MustCompile(`^\s$`)
)

// Represents the lexer of the programming lenguage
type Lexer struct {
	source        string // represents the source code
	character     string // represents the current character
	read_position int    // represents the next position in the current source
	position      int    // represents the current postion in source read
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

// read next token and assing a token type to the token
func (l *Lexer) NextToken() Token {
	l.skipWhiteSpaces()
	var token Token

	if l.isLetter(l.character) {
		literal := l.readIdentifier()
		token_type := LookUpTokenType(literal)
		return NewToken(token_type, literal)

	} else if l.isNumber(l.character) {
		literal := l.readNumber()
		if l.character == "." {
			literal += l.character
			l.readCharacter()
			literal += l.readNumber()
			return NewToken(FLOAT, literal)
		}
		return NewToken(INT, literal)
	}

	switch l.character {
	case "":
		token = NewToken(EOF, l.character)
	case "(":
		token = NewToken(LPAREN, l.character)
	case ")":
		token = NewToken(RPAREN, l.character)
	case "{":
		token = NewToken(LBRACE, l.character)
	case "}":
		token = NewToken(RBRACE, l.character)
	case "[":
		token = NewToken(LBRACKET, l.character)
	case "]":
		token = NewToken(RBRACKET, l.character)
	case ":":
		token = NewToken(COLON, l.character)
	case ",":
		token = NewToken(COMMA, l.character)
	case "%":
		token = NewToken(MOD, l.character)
	case ";":
		token = NewToken(SEMICOLON, l.character)
	case ".":
		token = NewToken(DOT, l.character)

	case "=":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(EQ)
		} else if l.peekCharacter() == ">" {
			token = l.makeTwoCharacterToken(ARROW)
		} else {
			token = NewToken(ASSING, l.character)
		}

	case "+":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(PLUSASSING)
		} else if l.peekCharacter() == "+" {
			token = l.makeTwoCharacterToken(PLUS2)
		} else {
			token = NewToken(PLUS, l.character)
		}

	case "<":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(LTOREQ)
		} else {
			token = NewToken(LT, l.character)
		}

	case ">":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(GTOREQ)
		} else {
			token = NewToken(GT, l.character)
		}

	case "|":
		if l.peekCharacter() == "|" {
			token = l.makeTwoCharacterToken(OR)
		} else {
			token = NewToken(ILLEGAL, l.character)
		}

	case "&":
		if l.peekCharacter() == "&" {
			token = l.makeTwoCharacterToken(AND)
		} else {
			token = NewToken(ILLEGAL, l.character)
		}

	case "-":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(MINUSASSING)
		} else if l.peekCharacter() == "-" {
			token = l.makeTwoCharacterToken(MINUS2)
		} else {
			token = NewToken(MINUS, l.character)
		}

	case "/":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(DIVASSING)
		} else {
			token = NewToken(DIVISION, l.character)
		}

	case "*":
		if l.peekCharacter() == "*" {
			token = l.makeTwoCharacterToken(EXPONENT)
		} else if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(TIMEASSI)
		} else {
			token = NewToken(TIMES, l.character)
		}

	case "!":
		if l.peekCharacter() == "=" {
			token = l.makeTwoCharacterToken(NOT_EQ)
		} else {
			token = Token{Token_type: NOT, Literal: l.character}
		}

	case `"`:
		literal := l.readString()
		token = NewToken(STRING, literal)

	default:
		token = NewToken(ILLEGAL, l.character)
	}

	l.readCharacter()
	return token
}

// check if current character is letter
func (l *Lexer) isLetter(char string) bool {
	return charRegex.MatchString(char)
}

// check if current character is number
func (l *Lexer) isNumber(char string) bool {
	return numRegex.MatchString(char)
}

// create a two character token
func (l *Lexer) makeTwoCharacterToken(tokenType TokenType) Token {
	prefix := l.character
	l.readCharacter()
	suffix := l.character
	return NewToken(tokenType, fmt.Sprintf("%s%s", prefix, suffix))
}

// read the current character and advance to  the next character
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
	for l.isLetter(l.character) || l.isNumber(l.character) {
		l.readCharacter()
	}

	return l.source[initialPosition:l.position]
}

// read a sequence of digits characters
func (l *Lexer) readNumber() string {
	initialPosition := l.position
	for l.isNumber(l.character) {
		l.readCharacter()
	}
	return l.source[initialPosition:l.position]
}

// read string will read a string literal
func (l *Lexer) readString() string {
	l.readCharacter()
	initialPosition := l.position

	for l.character != `"` && l.read_position <= utf8.RuneCountInString(l.source) {
		l.readCharacter()
	}

	return l.source[initialPosition:l.position]
}

// return the next of character of the current string
func (l *Lexer) peekCharacter() string {
	if l.read_position >= utf8.RuneCountInString(l.source) {
		return ""
	}

	return string(l.source[l.read_position])
}

// skip all whitespaces
func (l *Lexer) skipWhiteSpaces() {
	for wSpaceRegex.MatchString(l.character) {
		l.readCharacter()
	}
}
