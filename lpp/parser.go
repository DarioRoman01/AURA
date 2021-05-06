package lpp

type Parser struct {
	lexer        *Lexer
	currentToken *Token
	peekToken    *Token
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

func (p *Parser) ParseProgam() Program {
	program := Program{Staments: []Statement{}}

	for p.currentToken.Token_type != EOF {
		statement := p.parseStament()
		if statement != nil {
			program.Staments = append(program.Staments, *statement)
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

	return false
}

func (p *Parser) advanceTokens() {
	p.currentToken = p.peekToken
	nextToken := p.lexer.NextToken()
	p.peekToken = &nextToken

}
func (p *Parser) parseLetSatement() *Statement {
	stament := NewStatement(*p.currentToken, &LetStatement{})
	if !p.expepectedToken(IDENT) {
		return nil
	}

	stament.LetStatement.name = NewIdentifier(*p.currentToken, p.currentToken.Literal)
	if !p.expepectedToken(ASSING) {
		return nil
	}

	// todo finish when i know how to parse expressions
	for p.currentToken.Token_type != SEMICOLON {
		p.advanceTokens()
	}

	return stament
}

func (p *Parser) parseStament() *Statement {
	if p.currentToken.Token_type == LET {
		return p.parseLetSatement()
	}

	return nil
}
