package lpp

type Parser struct {
	lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) ParseProgam() Program {
	program := Program{Staments: []Statement{}}
	return program
}
