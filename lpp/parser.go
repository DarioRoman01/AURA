package lpp

import "fmt"

type PrefixParsFn func() Expression
type InfixParseFn func(Expression) *Expression

type PrefixParsFns map[TokenType]PrefixParsFn
type InfixParseFns map[TokenType]InfixParseFn

type Parser struct {
	lexer         *Lexer
	currentToken  *Token
	peekToken     *Token
	errors        []string
	prefixParsFns PrefixParsFns
	infixParseFns InfixParseFns
}

type Precedence int

const (
	HeadPrecendence Precedence = iota
	LOWEST                     = 1
	EQUEAL                     = 2
	LESSGRATER                 = 3
	SUM                        = 4
	PRODUCT                    = 5
	PREFIX                     = 6
	CALL                       = 7
)

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer:        lexer,
		currentToken: nil,
		peekToken:    nil,
	}

	parser.prefixParsFns = parser.registerPrefixFns()
	parser.infixParseFns = parser.registerInfixFns()
	parser.advanceTokens()
	parser.advanceTokens()
	return parser
}

func (p *Parser) advanceTokens() {
	p.currentToken = p.peekToken
	nextToken := p.lexer.NextToken()
	p.peekToken = &nextToken
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

func (p *Parser) parseExpression(Precedence) Expression {
	if p.currentToken == nil {
		panic("current token cannot be nil")
	}

	prefixParseFn, exist := p.prefixParsFns[p.currentToken.Token_type]
	fmt.Println(p.currentToken.Token_type)
	if !exist {
		fmt.Println("entro aqui")
		return nil
	}

	leftExpression := prefixParseFn()
	return leftExpression
}

func (p *Parser) parserExpressionStatement() *ExpressionStament {
	if p.currentToken == nil {
		panic("peek token cannot be bil")
	}

	fmt.Println(p.currentToken.Literal)
	expressionStament := NewExpressionStament(*p.currentToken, nil)
	expressionStament.Expression = p.parseExpression(LOWEST)

	if p.peekToken == nil {
		panic("peek token cannot be bil")
	}
	if p.peekToken.Token_type == SEMICOLON {
		p.advanceTokens()
	}

	return expressionStament
}

func (p *Parser) parseIdentifier() Expression {
	if p.currentToken == nil {
		panic("current token cannot be nil in parse identifier")
	}

	return &Identifier{token: *p.currentToken, value: p.currentToken.Literal}
}

func (p *Parser) parseLetSatement() Stmt {
	stament := NewLetStatement(*p.currentToken, nil, nil)
	if !p.expepectedToken(IDENT) {
		return nil
	}

	stament.Name = p.parseIdentifier().(*Identifier)
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

	return p.parserExpressionStatement()
}

func (p *Parser) registerInfixFns() InfixParseFns {
	inFixFns := make(InfixParseFns)
	return inFixFns
}

func (p *Parser) registerPrefixFns() PrefixParsFns {
	prefixFns := make(PrefixParsFns)
	prefixFns[IDENT] = p.parseIdentifier
	return prefixFns
}
