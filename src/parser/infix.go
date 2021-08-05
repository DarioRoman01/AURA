package parser

import (
	"aura/src/ast"
	l "aura/src/lexer"
)

// parse a method expression
func (p *Parser) parseMethod(left ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.IDENT) {
		// syntax error. we dont allow this -> obj:();
		return nil
	}

	method := p.parseExpression(LOWEST)
	return ast.NewMethodExpression(token, left, method)
}

// parse an infix expressoin
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	operator := p.currentToken.Literal
	precedence := p.currentPrecedence()

	p.advanceTokens()
	rigth := p.parseExpression(precedence)
	return ast.Newinfix(token, rigth, operator, left)
}

// parse a function call
func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	args := p.parseExpressions(l.RPAREN)
	return ast.NewCall(token, function, args...)
}

// parse a call list expression
func (p *Parser) parseCallList(valueList ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	p.advanceTokens()
	index := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RBRACKET) {
		// syntax error. we dont allow tihs -> lista[2,3,4,5;
		return nil
	}

	return ast.NewCallList(token, valueList, index)
}

// parse a ressigment expression
func (p *Parser) parseReassigment(ident ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	p.checkPeekTokenIsNotNil()
	p.advanceTokens()
	newVal := p.parseExpression(LOWEST)
	return ast.NewReassignment(token, ident, newVal)
}

// parse a key value expression
func (p *Parser) parseKeyValues() *ast.KeyValue {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	key := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.ARROW) {
		return nil
	}

	p.advanceTokens()
	value := p.parseExpression(LOWEST)
	return ast.NewKeyVal(token, key, value)
}

// parse a range expression
func (p *Parser) parseRangeExpression() ast.Expression {
	token := p.currentToken
	if !p.expepectedToken(l.IDENT) {
		// syntax error. we dont allow this -> por(en rango(10))
		return nil
	}
	variable := p.parseIdentifier()
	if !p.expepectedToken(l.IN) {
		// syntax error. we dont allow this -> por(i rango(10))
		return nil
	}

	p.advanceTokens()
	exp := p.parseExpression(LOWEST)
	return ast.NewRange(token, variable, exp)
}

// parse a class field or method call
func (p *Parser) parseClassFieldsCall(left ast.Expression) ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	p.checkPeekTokenIsNotNil()
	p.advanceTokens()
	field := p.parseExpression(LOWEST)
	return ast.NewClassFieldCall(token, left, field)
}

// parse an assigment expression like
//		x := 5;
func (p *Parser) parseAssigmentExp(left ast.Expression) ast.Expression {
	ident, isIdent := left.(*ast.Identifier)
	if !isIdent {
		return nil
	}

	token := p.currentToken
	p.advanceTokens()
	val := p.parseExpression(LOWEST)
	return ast.NewAssigmentExp(token, ident, val)
}

// parse an arrow function expression
func (p *Parser) parseArrowFunc() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	params := p.parseIdentifiers(l.BAR)
	if !p.expepectedToken(l.ARROW) {
		return nil
	}

	var body *ast.Block
	if p.peekToken.Token_type == l.LBRACE {
		p.advanceTokens()
		body = p.parseBlock()
	} else {
		p.advanceTokens()
		exp := p.parserExpressionStatement()
		body = ast.NewBlock(p.currentToken, exp)
	}

	return ast.NewArrowFunc(token, params, body)
}
