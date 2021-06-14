package parser

import (
	"aura/src/ast"
	l "aura/src/lexer"
)

func (p *Parser) parseMethod(left ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	method := ast.NewMethodExpression(*p.currentToken, left, nil)
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	method.Method = p.parseExpression(LOWEST)
	return method
}

// parse infix expressoins
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	infix := ast.Newinfix(*p.currentToken, nil, p.currentToken.Literal, left)
	precedence := p.currentPrecedence()
	p.advanceTokens()
	infix.Rigth = p.parseExpression(precedence)
	return infix
}

// parse function calls
func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	call := ast.NewCall(*p.currentToken, function)
	call.Arguments = p.parseCallArguments()
	return call
}

func (p *Parser) parseCallList(valueList ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	callList := ast.NewCallList(*p.currentToken, valueList, nil)
	p.advanceTokens()
	callList.Index = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RBRACKET) {
		return nil
	}

	return callList
}

func (p *Parser) parseReassigment(ident ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	reassignment := ast.NewReassignment(*p.currentToken, ident, nil)
	p.advanceTokens()
	reassignment.NewVal = p.parseExpression(LOWEST)
	return reassignment
}

func (p *Parser) parseKeyValues() *ast.KeyValue {
	p.checkCurrentTokenIsNotNil()
	keyVal := ast.NewKeyVal(*p.currentToken, nil, nil)
	keyVal.Key = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.ARROW) {
		return nil
	}

	p.advanceTokens()
	keyVal.Value = p.parseExpression(LOWEST)
	return keyVal
}

func (p *Parser) parseRangeExpression() ast.Expression {
	rangeExpress := ast.NewRange(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	rangeExpress.Variable = p.parseIdentifier()
	if !p.expepectedToken(l.IN) {
		return nil
	}

	p.advanceTokens()
	rangeExpress.Range = p.parseExpression(LOWEST)
	return rangeExpress
}
