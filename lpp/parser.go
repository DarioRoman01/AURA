package lpp

import "fmt"

type Parser struct {
	lexer        *Lexer
	currentToken *Token
	peekToken    *Token
	errors       []string
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer:        lexer,
		currentToken: nil,
		peekToken:    nil,
	}

	parser.advanceTokens()
	parser.advanceTokens()

	return parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgam() Program {
	program := Program{Staments: []Stmt{}}

	for p.currentToken.Token_type != EOF {
		statement := p.parseStament()
		if statement != nil {
			program.Staments = append(program.Staments, statement)
		}

		p.advanceTokens()
	}
	return program
}

func (p *Parser) expepectedToken(tokenType TokenType) bool {
	if p.peekToken.Token_type == tokenType {
		p.advanceTokens()
		return true
	}

	p.expectedTokenError(tokenType)
	return false
}

func (p *Parser) expectedTokenError(tokenType TokenType) {
	if p.peekToken == nil {
		panic("peer token is nil")
	}

	err := fmt.Sprintf("se esperaba que el siguient token fuera %s pero se obtuvo %s", tokens[tokenType], tokens[p.peekToken.Token_type])
	p.errors = append(p.errors, err)
}

func (p *Parser) advanceTokens() {
	p.currentToken = p.peekToken
	nextToken := p.lexer.NextToken()
	p.peekToken = &nextToken

}
func (p *Parser) parseLetSatement() Stmt {
	stament := NewLetStatement(*p.currentToken, nil, nil)
	if !p.expepectedToken(IDENT) {
		return nil
	}

	stament.Name = NewIdentifier(*p.currentToken, p.currentToken.Literal)
	if !p.expepectedToken(ASSING) {
		return nil
	}

	// todo finish when i know how to parse expressions
	for p.currentToken.Token_type != SEMICOLON {
		p.advanceTokens()
	}

	return stament
}

func (p *Parser) parseReturnStatement() Stmt {
	stament := NewReturnStatement(*p.currentToken, nil)
	p.advanceTokens()

	// todo finish when i know how to parse expressions
	for p.currentToken.Token_type != SEMICOLON {
		p.advanceTokens()
	}

	return stament
}

func (p *Parser) parseStament() Stmt {
	if p.currentToken.Token_type == LET {
		return p.parseLetSatement()
	} else if p.currentToken.Token_type == RETURN {
		return p.parseReturnStatement()
	}

	return nil
}
