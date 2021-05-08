package lpp

import (
	"fmt"
	"strconv"
)

type PrefixParsFn func() Expression
type InfixParseFn func(Expression) Expression

type PrefixParsFns map[TokenType]PrefixParsFn
type InfixParseFns map[TokenType]InfixParseFn

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

var PRECEDENCES = map[TokenType]Precedence{
	EQ:       EQUEAL,
	NOT_EQ:   EQUEAL,
	LT:       LESSGRATER,
	GT:       LESSGRATER,
	PLUS:     SUM,
	MINUS:    SUM,
	DIVISION: PRODUCT,
	TIMES:    PRODUCT,
}

type Parser struct {
	lexer         *Lexer
	currentToken  *Token
	peekToken     *Token
	errors        []string
	prefixParsFns PrefixParsFns
	infixParseFns InfixParseFns
}

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

func (p *Parser) checkCurrentTokenIsNotNil() {
	if p.currentToken == nil {
		panic("current token cannot be nil")
	}
}

func (p *Parser) checkPeekTokenIsNotNil() {
	if p.peekToken == nil {
		panic("peek token cannot be nil")
	}
}

func (p *Parser) currentPrecedence() Precedence {
	p.checkCurrentTokenIsNotNil()
	precedence, exists := PRECEDENCES[p.currentToken.Token_type]
	if !exists {
		return LOWEST
	}

	return precedence
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
	p.checkCurrentTokenIsNotNil()
	err := fmt.Sprintf("se esperaba que el siguient token fuera %s pero se obtuvo %s", tokens[tokenType], tokens[p.peekToken.Token_type])
	p.errors = append(p.errors, err)
}

func (p *Parser) parseBoolean() Expression {
	p.checkCurrentTokenIsNotNil()
	var value bool
	if p.currentToken.Token_type == TRUE {
		value = true
		return NewBoolean(*p.currentToken, &value)
	}

	value = false
	return NewBoolean(*p.currentToken, &value)
}

func (p *Parser) parseBlock() *Block {
	p.checkCurrentTokenIsNotNil()
	blockStament := NewBlock(*p.currentToken)
	p.advanceTokens()

	for !(p.currentToken.Token_type == RBRACE) && !(p.currentToken.Token_type == EOF) {
		stament := p.parseStament()
		if stament != nil {
			blockStament.Staments = append(blockStament.Staments, stament)
		}

		p.advanceTokens()
	}

	return blockStament
}

func (p *Parser) parseExpression(precedence Precedence) Expression {
	p.checkCurrentTokenIsNotNil()
	prefixParseFn, exist := p.prefixParsFns[p.currentToken.Token_type]
	if !exist {
		message := fmt.Sprintf("no se encontro ninguna funcion para parsear %s", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	leftExpression := prefixParseFn()
	p.checkPeekTokenIsNotNil()

	for !(p.peekToken.Token_type == SEMICOLON) && precedence < p.peekPrecedence() {
		infixParseFn, exist := p.infixParseFns[p.peekToken.Token_type]
		if !exist {
			return leftExpression
		}

		p.advanceTokens()
		if leftExpression == nil {
			panic("cannot be nil")
		}
		leftExpression = infixParseFn(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parserExpressionStatement() *ExpressionStament {
	p.checkCurrentTokenIsNotNil()
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

func (p *Parser) parseGroupExpression() Expression {
	p.advanceTokens()
	expression := p.parseExpression(LOWEST)
	if !p.expepectedToken(RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) parseIdentifier() Expression {
	p.checkCurrentTokenIsNotNil()
	return &Identifier{token: *p.currentToken, value: p.currentToken.Literal}
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	p.checkCurrentTokenIsNotNil()
	infix := Newinfix(*p.currentToken, nil, p.currentToken.Literal, left)
	precedence := p.currentPrecedence()
	p.advanceTokens()
	infix.Rigth = p.parseExpression(precedence)
	return infix
}

func (p *Parser) parseIf() Expression {
	p.checkCurrentTokenIsNotNil()
	ifExpression := NewIf(*p.currentToken, nil, nil, nil)
	if !p.expepectedToken(LPAREN) {
		return nil
	}

	p.advanceTokens()
	ifExpression.Condition = p.parseExpression(LOWEST)
	if !p.expepectedToken(RPAREN) {
		return nil
	}

	if !p.expepectedToken(LBRACE) {
		return nil
	}

	ifExpression.Consequence = p.parseBlock()

	p.checkPeekTokenIsNotNil()
	if p.peekToken.Token_type == ELSE {
		p.advanceTokens()
		if !p.expepectedToken(LBRACE) {
			return nil
		}

		ifExpression.Alternative = p.parseBlock()
	}

	return ifExpression
}

func (p *Parser) parseInteger() Expression {
	p.checkCurrentTokenIsNotNil()
	integer := NewInteger(*p.currentToken, nil)
	val, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		message := fmt.Sprintf("no se pudo parsear %s como entero", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	integer.Value = &val
	return integer
}

func (p *Parser) parseLetSatement() Stmt {
	p.checkCurrentTokenIsNotNil()
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

func (p *Parser) parsePrefixExpression() Expression {
	p.checkCurrentTokenIsNotNil()
	prefixExpression := NewPrefix(*p.currentToken, p.currentToken.Literal, nil)
	p.advanceTokens()
	prefixExpression.Rigth = p.parseExpression(PREFIX)
	return prefixExpression
}

func (p *Parser) parseReturnStatement() Stmt {
	p.checkCurrentTokenIsNotNil()
	stament := NewReturnStatement(*p.currentToken, nil)
	p.advanceTokens()

	// todo finish when i know how to parse expressions
	for p.currentToken.Token_type != SEMICOLON {
		p.advanceTokens()
	}

	return stament
}

func (p *Parser) parseStament() Stmt {
	p.checkCurrentTokenIsNotNil()
	if p.currentToken.Token_type == LET {
		return p.parseLetSatement()
	} else if p.currentToken.Token_type == RETURN {
		return p.parseReturnStatement()
	}

	return p.parserExpressionStatement()
}

func (p *Parser) peekPrecedence() Precedence {
	p.checkPeekTokenIsNotNil()
	precedence, exists := PRECEDENCES[p.peekToken.Token_type]
	if !exists {
		return LOWEST
	}

	return precedence
}

func (p *Parser) registerInfixFns() InfixParseFns {
	inFixFns := make(InfixParseFns)
	inFixFns[PLUS] = p.parseInfixExpression
	inFixFns[MINUS] = p.parseInfixExpression
	inFixFns[DIVISION] = p.parseInfixExpression
	inFixFns[TIMES] = p.parseInfixExpression
	inFixFns[EQ] = p.parseInfixExpression
	inFixFns[NOT_EQ] = p.parseInfixExpression
	inFixFns[LT] = p.parseInfixExpression
	inFixFns[GT] = p.parseInfixExpression
	return inFixFns
}

func (p *Parser) registerPrefixFns() PrefixParsFns {
	prefixFns := make(PrefixParsFns)
	prefixFns[FALSE] = p.parseBoolean
	prefixFns[IDENT] = p.parseIdentifier
	prefixFns[IF] = p.parseIf
	prefixFns[INT] = p.parseInteger
	prefixFns[LPAREN] = p.parseGroupExpression
	prefixFns[MINUS] = p.parsePrefixExpression
	prefixFns[NOT] = p.parsePrefixExpression
	prefixFns[TRUE] = p.parseBoolean
	return prefixFns
}
