package lpp

type Parser struct {
	lexer        *Lexer
	currentToken *Token
	peekToken    *Token
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: nil,
		peekToken:    nil,
	}
}

func (p *Parser) ParseProgam() Program {
	program := Program{Staments: []Statement{}}

	for p.currentToken.Token_type != EOF {
		statement := p.parseStament()
		if statement != nil {
			program.Staments = append(program.Staments)
		}
	}
	return program
}

func (p *Parser) expepectedToken(tokenType TokenType) bool {
	return true
}

func (p *Parser) advanceTokens() {

}

func (p *Parser) parseLetSatement() *LetStatement {
	letStatemnt := NewLetStatement(*p.currentToken, nil, nil)

	if !p.expepectedToken(IDENT) {
		return nil
	}

	letStatemnt.name = NewIdentifier(*p.currentToken, p.currentToken.Literal)

	if !p.expepectedToken(ASSING) {
		return nil
	}

	// todo finish when i know how to parse expressions
	for p.currentToken.Token_type != SEMICOLON {
		p.advanceTokens()
	}

	return letStatemnt
}

func (p *Parser) parseStament() *Statement {
	if p.currentToken.Token_type == LET {
		return p.parseLetSatement()
	}

	return nil
}
