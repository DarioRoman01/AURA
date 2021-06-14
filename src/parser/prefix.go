package parser

import (
	"aura/src/ast"
	l "aura/src/lexer"
	"fmt"
	"strconv"
)

// parse boolean expression and check if true or false
func (p *Parser) parseBoolean() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	var value bool
	if p.currentToken.Token_type == l.TRUE {
		value = true
		return ast.NewBoolean(*p.currentToken, &value)
	}

	value = false
	return ast.NewBoolean(*p.currentToken, &value)
}

func (p *Parser) parseFor() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	forExpression := ast.NewFor(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	forExpression.Condition = p.parseRangeExpression()
	if !p.expepectedToken(l.RPAREN) {
		return nil
	}
	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	forExpression.Body = p.parseBlock()
	return forExpression
}

// parse a function declaration
func (p *Parser) parseFunction() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	function := ast.NewFunction(*p.currentToken, nil)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	function.Parameters = p.parseFunctionParameters()
	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	function.Body = p.parseBlock()
	return function
}

func (p *Parser) parseWhile() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	whileExpression := ast.NewWhile(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	p.advanceTokens()
	whileExpression.Condition = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	whileExpression.Body = p.parseBlock()
	return whileExpression
}

// parse a identifier
func (p *Parser) parseIdentifier() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return &ast.Identifier{Token: *p.currentToken, Value: p.currentToken.Literal}
}

// parse if expressions, check sintax and if there is an else in the expression
func (p *Parser) parseIf() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	ifExpression := ast.NewIf(*p.currentToken, nil, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	p.advanceTokens()
	ifExpression.Condition = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	ifExpression.Consequence = p.parseBlock()

	p.checkPeekTokenIsNotNil()
	if p.peekToken.Token_type == l.ELSE {
		p.advanceTokens()
		if !p.expepectedToken(l.LBRACE) {
			return nil
		}

		ifExpression.Alternative = p.parseBlock()
	}

	return ifExpression
}

// parse integer expressions
func (p *Parser) parseInteger() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	integer := ast.NewInteger(*p.currentToken, nil)

	val, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		message := fmt.Sprintf("no se pudo parsear %s como entero", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	integer.Value = &val
	return integer
}

// parse group expression like (5 + 5) / 2
func (p *Parser) parseGroupExpression() ast.Expression {
	p.advanceTokens()
	expression := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		return nil
	}

	return expression
}

// parse a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	prefixExpression := ast.NewPrefix(*p.currentToken, p.currentToken.Literal, nil)
	p.advanceTokens()
	prefixExpression.Rigth = p.parseExpression(PREFIX)
	return prefixExpression
}

func (p *Parser) parseStringLiteral() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return ast.NewStringLiteral(*p.currentToken, p.currentToken.Literal)
}

func (p *Parser) ParseArray() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	arr := ast.NewArray(*p.currentToken, nil)
	if !p.expepectedToken(l.LBRACKET) {
		return nil
	}

	arr.Values = p.ParseArrayValues()
	return arr
}

func (p *Parser) parseMap() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	mapExpress := ast.NewMapExpression(*p.currentToken, []*ast.KeyValue{})
	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	p.advanceTokens()
	keyVal := p.parseKeyValues()
	if keyVal == nil {
		return nil
	}

	mapExpress.Body = append(mapExpress.Body, keyVal)
	for p.peekToken.Token_type == l.COMMA {
		p.advanceTokens()
		p.advanceTokens()
		keyVal := p.parseKeyValues()
		if keyVal == nil {
			return nil
		}

		mapExpress.Body = append(mapExpress.Body, keyVal)
	}
	p.advanceTokens()
	return mapExpress
}
